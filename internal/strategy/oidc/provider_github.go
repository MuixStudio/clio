/*
 * Copyright 2026 MuixStudio
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package oidc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Compile-time assertion that ProviderGitHub satisfies Provider.
var _ Provider = (*ProviderGitHub)(nil)

const (
	githubAuthURL      = "https://github.com/login/oauth/authorize"
	githubTokenURL     = "https://github.com/login/oauth/access_token"
	githubUserURL      = "https://api.github.com/user"
	githubUserEmailURL = "https://api.github.com/user/emails"
	githubIssuer       = "https://github.com/login/oauth/access_token" // stable token endpoint URL
)

// GitHubConfig holds the OAuth2 application credentials registered at
// https://github.com/settings/applications/new.
type GitHubConfig struct {
	ClientID     string
	ClientSecret string
	// Scopes defaults to ["read:user", "user:email"] when empty.
	Scopes []string
}

// ProviderGitHub implements Provider for GitHub OAuth2.
// It handles the code↔token exchange and normalizes the GitHub user
// profile into the shared Claims struct.
type ProviderGitHub struct {
	cfg *GitHubConfig
}

// NewProviderGitHub constructs a GitHub provider.
// Default scopes ("read:user user:email") are applied when cfg.Scopes is empty.
func NewProviderGitHub(cfg *GitHubConfig) *ProviderGitHub {
	if len(cfg.Scopes) == 0 {
		cfg.Scopes = []string{"read:user", "user:email"}
	}
	return &ProviderGitHub{cfg: cfg}
}

func (p *ProviderGitHub) ID() string { return "github" }

// AuthCodeURL returns the GitHub authorization URL the client should be redirected to.
func (p *ProviderGitHub) AuthCodeURL(state string) string {
	v := url.Values{
		"client_id": {p.cfg.ClientID},
		//"redirect_uri": {p.cfg.RedirectURI},
		"scope": {strings.Join(p.cfg.Scopes, " ")},
		"state": {state},
	}
	return githubAuthURL + "?" + v.Encode()
}

// Exchange trades the one-time authorization code for normalized Claims.
// Internally: code → access token → /user (+ /user/emails fallback).
func (p *ProviderGitHub) Exchange(ctx context.Context, code string) (*Claims, error) {
	token, err := p.exchangeCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("github: token exchange: %w", err)
	}
	claims, err := p.fetchClaims(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("github: fetch user: %w", err)
	}
	return claims, nil
}

// =========================================================================
// GitHub API helpers
// =========================================================================

// exchangeCode posts the authorization code to GitHub's token endpoint.
func (p *ProviderGitHub) exchangeCode(ctx context.Context, code string) (string, error) {
	body := url.Values{
		"client_id":     {p.cfg.ClientID},
		"client_secret": {p.cfg.ClientSecret},
		"code":          {code},
		//"redirect_uri":  {p.cfg.RedirectURI},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, githubTokenURL,
		strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Error != "" {
		return "", fmt.Errorf("%s: %s", result.Error, result.ErrorDesc)
	}
	return result.AccessToken, nil
}

type githubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type githubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

// fetchClaims calls /user (and /user/emails if the email is private) and
// maps the response into the normalized Claims struct.
func (p *ProviderGitHub) fetchClaims(ctx context.Context, token string) (*Claims, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubUserURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub /user returned HTTP %d", resp.StatusCode)
	}

	var u githubUser
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}

	// GitHub omits email when the user marks it private; fall back to /user/emails.
	if u.Email == "" {
		u.Email, _ = p.fetchPrimaryEmail(ctx, token)
	}

	return &Claims{
		Subject:  fmt.Sprintf("%d", u.ID), // numeric ID is stable and immutable
		Issuer:   githubIssuer,
		Email:    u.Email,
		Name:     u.Name,
		Nickname: u.Login,
		Picture:  u.AvatarURL,
	}, nil
}

// fetchPrimaryEmail fetches the verified primary email from /user/emails.
// Not fatal if it fails — identity can exist without an email.
func (p *ProviderGitHub) fetchPrimaryEmail(ctx context.Context, token string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubUserEmailURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []githubEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	return "", nil
}
