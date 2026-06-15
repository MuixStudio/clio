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

package login

import (
	"context"

	"github.com/muixstudio/clio/internal/flow"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/session"
)

// PreHookExecutor runs before the login flow is processed.
type PreHookExecutor interface {
	ExecuteLoginPreHook(ctx context.Context, f flow.Flow) error
}

// PostHookExecutor runs after authentication succeeds, before the session
// is returned to the client. Receives the freshly issued session so hooks
// can augment it (e.g. add AMR claims) or send notifications.
type PostHookExecutor interface {
	ExecuteLoginPostHook(ctx context.Context, f flow.Flow, i *identity.Identity, sess *session.Session) error
}

// HooksProvider is implemented by the Driver.
type HooksProvider interface {
	PreLoginHooks(ctx context.Context) []PreHookExecutor
	PostLoginHooks(ctx context.Context) []PostHookExecutor
}

// HookExecutorProvider allows the Handler to retrieve the HookExecutor.
type HookExecutorProvider interface {
	LoginHookExecutor() *HookExecutor
}

type HookExecutor struct{ d HooksProvider }

func NewHookExecutor(d HooksProvider) *HookExecutor { return &HookExecutor{d: d} }

func (e *HookExecutor) PreHook(ctx context.Context, f flow.Flow) error {
	for _, h := range e.d.PreLoginHooks(ctx) {
		if err := h.ExecuteLoginPreHook(ctx, f); err != nil {
			return err
		}
	}
	return nil
}

func (e *HookExecutor) PostHook(ctx context.Context, f flow.Flow, i *identity.Identity, sess *session.Session) error {
	for _, h := range e.d.PostLoginHooks(ctx) {
		if err := h.ExecuteLoginPostHook(ctx, f, i, sess); err != nil {
			return err
		}
	}
	return nil
}
