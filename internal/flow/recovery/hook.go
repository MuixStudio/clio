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

package recovery

import (
	"context"

	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/session"
)

// PreHookExecutor runs before a recovery flow is created (rate limiting, etc.).
type PreHookExecutor interface {
	ExecuteRecoveryPreHook(ctx context.Context, f *Flow) error
}

// PostHookExecutor runs after recovery succeeds and a session is issued.
// Receives the flow, the recovered identity, and the new session — mirrors
// the login PostHookExecutor signature so the same hooks can be reused.
type PostHookExecutor interface {
	ExecuteRecoveryPostHook(ctx context.Context, f *Flow, i *identity.Identity, sess *session.Session) error
}

// HooksProvider is satisfied by the Driver.
type HooksProvider interface {
	PreRecoveryHooks(ctx context.Context) []PreHookExecutor
	PostRecoveryHooks(ctx context.Context) []PostHookExecutor
}

// HookExecutorProvider lets the handler resolve the executor without depending
// on a concrete type.
type HookExecutorProvider interface {
	RecoveryHookExecutor() *HookExecutor
}

// HookExecutor runs the pre- and post-hook chains for recovery flows.
type HookExecutor struct{ d HooksProvider }

func NewHookExecutor(d HooksProvider) *HookExecutor { return &HookExecutor{d: d} }

func (e *HookExecutor) PreHook(ctx context.Context, f *Flow) error {
	for _, h := range e.d.PreRecoveryHooks(ctx) {
		if err := h.ExecuteRecoveryPreHook(ctx, f); err != nil {
			return err
		}
	}
	return nil
}

func (e *HookExecutor) PostHook(ctx context.Context, f *Flow, i *identity.Identity, sess *session.Session) error {
	for _, h := range e.d.PostRecoveryHooks(ctx) {
		if err := h.ExecuteRecoveryPostHook(ctx, f, i, sess); err != nil {
			return err
		}
	}
	return nil
}
