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

// Package config loads and validates mini-identity's runtime configuration
// using Viper with mapstructure decoding.
//
// All values can be overridden via environment variables:
//
//	SERVE_PORT=8080 DSN=postgres://... ./mini-identity
//
// Env var names are the dot-path uppercased with dots replaced by underscores.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/muixstudio/clio/internal/infra/database"
	"github.com/muixstudio/clio/internal/infra/logger"
	"github.com/spf13/viper"
)

type (
	// Config is the top-level configuration struct.
	// Loaded once at startup; passed to driver.New as the single source of truth.
	Config struct {
		Env string `mapstructure:"env"` // Env, e.g., "production" or "development"

		Log         logger.Options    `mapstructure:"log"`      // Log configuration
		Database    database.DBConfig `mapstructure:"database"` // Database configuration
		Serve       ServeConfig       `mapstructure:"serve"`
		Session     SessionConfig     `mapstructure:"session"`
		SelfService SelfServiceConfig `mapstructure:"selfservice"`
		Courier     CourierConfig     `mapstructure:"courier"`
	}

	// ServeConfig controls the HTTP listener.
	ServeConfig struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}

	// SessionConfig controls session lifetime and token settings.
	SessionConfig struct {
		Lifespan            time.Duration `mapstructure:"lifespan"`
		AccessTokenLifespan time.Duration `mapstructure:"access_token_lifespan"`
		JWTSecret           string        `mapstructure:"jwt_secret"`
	}

	// SelfServiceConfig holds self-service flow and method configuration.
	SelfServiceConfig struct {
		Flows   FlowsConfig   `mapstructure:"flows"`
		Methods MethodsConfig `mapstructure:"methods"`
	}

	// FlowsConfig holds per-flow-type settings.
	FlowsConfig struct {
		Registration RegistrationFlowsConfig `mapstructure:"registration"`
		Login        LoginFlowsConfig        `mapstructure:"login"`
		Logout       LogoutFlowsConfig       `mapstructure:"logout"`
		Recovery     RecoveryFlowsConfig     `mapstructure:"recovery"`
		Verification VerificationFlowsConfig `mapstructure:"verification"`
	}

	// FlowConfig is the base configuration embedded by every flow type.
	FlowConfig struct {
		// Lifespan is how long a created flow is valid before it expires.
		Lifespan time.Duration `mapstructure:"lifespan"`
		After    struct {
			DefaultBrowserReturnUrl string `mapstructure:"default_browser_return_url"`
		} `mapstructure:"after"`
	}

	RegistrationFlowsConfig struct {
		FlowConfig `mapstructure:",squash"`
	}
	LoginFlowsConfig struct {
		FlowConfig `mapstructure:",squash"`
	}
	LogoutFlowsConfig struct {
		FlowConfig `mapstructure:",squash"`
	}
	RecoveryFlowsConfig struct {
		FlowConfig `mapstructure:",squash"`
	}
	VerificationFlowsConfig struct {
		FlowConfig `mapstructure:",squash"`
	}

	// MethodsConfig holds per-strategy settings.
	MethodsConfig struct {
		Password PasswordMethodConfig `mapstructure:"password"`
		OIDC     OIDCMethodConfig     `mapstructure:"oidc"`
		Code     CodeMethodConfig     `mapstructure:"code"`
	}

	// PasswordMethodConfig controls the password strategy.
	PasswordMethodConfig struct {
		// Enabled activates password login (and optionally registration).
		Enabled bool `mapstructure:"enabled"`
		Config  struct {
			// RegistrationEnabled, when false, disables password-based sign-up
			// while keeping password login active for existing accounts.
			RegistrationEnabled bool `mapstructure:"registration_enabled"`
		} `mapstructure:"config"`
	}

	// OIDCMethodConfig controls the OIDC/OAuth2 strategy.
	OIDCMethodConfig struct {
		// Enabled activates the strategy for both login and registration.
		Enabled bool `mapstructure:"enabled"`
		Config  struct {
			Providers []OIDCProvider `mapstructure:"providers"`
		} `mapstructure:"config"`
	}

	// OIDCProvider describes a single OAuth2 / OIDC provider.
	OIDCProvider struct {
		// ID is a unique slug used as the provider key (e.g. "github").
		ID string `mapstructure:"id"`
		// Provider selects the built-in implementation ("github", etc.).
		Provider     string `mapstructure:"provider"`
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
	}

	// CodeMethodConfig controls the passwordless email-code strategy.
	CodeMethodConfig struct {
		// Enabled activates the strategy for both login and registration.
		Enabled bool `mapstructure:"enabled"`
	}

	// CourierConfig controls how outgoing messages (verification codes) are sent.
	CourierConfig struct {
		// Strategy is "console" (dev, prints to stdout) or "smtp" (production).
		Strategy string     `mapstructure:"strategy"`
		SMTP     SMTPConfig `mapstructure:"smtp"`
	}

	// SMTPConfig is used when CourierConfig.Strategy == "smtp".
	SMTPConfig struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		From     string `mapstructure:"from"`
	}

	Provider interface {
		Config() *Config
	}
)

// Addr returns the host:port string for http.ListenAndServe.
func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Serve.Host, c.Serve.Port)
}

// =========================================================================
// Loader
// =========================================================================

// Load reads the YAML config file at path, applies defaults, and returns a
// validated *Config. Environment variables take precedence over file values.
//
// Supported duration formats: "15m", "24h", "720h", etc.
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)

	// AutomaticEnv maps env vars to config keys.
	// Example: SERVE_PORT=8080 overrides serve.port.
	// SetEnvKeyReplacer converts "." to "_" so DATABASE_HOST maps to database.host.
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// ---- Defaults -------------------------------------------------------
	// General
	v.SetDefault("env", "development")

	// HTTP server
	v.SetDefault("serve.host", "0.0.0.0")
	v.SetDefault("serve.port", 9899)

	// Logging
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "console")
	v.SetDefault("log.maxsize", 100)
	v.SetDefault("log.maxbackups", 3)
	v.SetDefault("log.maxage", 7)
	v.SetDefault("log.compress", false)

	// Database
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.time_zone", "UTC")
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.conn_max_lifetime", "1h")
	v.SetDefault("database.conn_max_idle_time", "10m")
	v.SetDefault("database.slow_threshold", "200ms")
	v.SetDefault("database.parameterized_queries", true)
	v.SetDefault("database.prepare_stmt", true)
	v.SetDefault("database.disable_foreign_key_constraint_when_migrating", false)
	v.SetDefault("database.singular_table", false)

	// Session
	v.SetDefault("session.lifespan", "24h")
	v.SetDefault("session.access_token_lifespan", "15m")
	v.SetDefault("session.jwt_secret", "change-me-in-production")

	// Self-service flows
	v.SetDefault("selfservice.flows.registration.lifespan", "15m")
	v.SetDefault("selfservice.flows.login.lifespan", "15m")
	v.SetDefault("selfservice.flows.logout.lifespan", "5m")
	v.SetDefault("selfservice.flows.recovery.lifespan", "15m")
	v.SetDefault("selfservice.flows.verification.lifespan", "15m")

	// Self-service methods
	v.SetDefault("selfservice.methods.password.enabled", true)
	v.SetDefault("selfservice.methods.password.config.registration_enabled", true)
	v.SetDefault("selfservice.methods.oidc.enabled", false)
	v.SetDefault("selfservice.methods.code.enabled", false)

	// Courier
	v.SetDefault("courier.strategy", "console")
	v.SetDefault("courier.smtp.port", 587)
	// ----------------------------------------------------------------------

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal: %w", err)
	}

	return &cfg, nil
}
