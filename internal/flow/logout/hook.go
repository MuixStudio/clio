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

package logout

import "context"

// PreHookExecutor runs before the session is revoked.
// Use cases: rate-limit checks, admin overrides, audit initiation.
type PreHookExecutor interface {
	ExecuteLogoutPreHook(ctx context.Context, f *Flow) error
}

// PostHookExecutor runs after the session is revoked and the flow is completed.
// Use cases: audit logging, revoking refresh tokens, clearing cookies, analytics.
type PostHookExecutor interface {
	ExecuteLogoutPostHook(ctx context.Context, f *Flow) error
}

// HooksProvider is implemented by the Driver.
type HooksProvider interface {
	PreLogoutHooks(ctx context.Context) []PreHookExecutor
	PostLogoutHooks(ctx context.Context) []PostHookExecutor
}

// HookExecutorProvider allows the Handler to retrieve the HookExecutor.
type HookExecutorProvider interface {
	LogoutHookExecutor() *HookExecutor
}

type HookExecutor struct{ d HooksProvider }

func NewHookExecutor(d HooksProvider) *HookExecutor { return &HookExecutor{d: d} }

func (e *HookExecutor) PreHook(ctx context.Context, f *Flow) error {
	for _, h := range e.d.PreLogoutHooks(ctx) {
		if err := h.ExecuteLogoutPreHook(ctx, f); err != nil {
			return err
		}
	}
	return nil
}

func (e *HookExecutor) PostHook(ctx context.Context, f *Flow) error {
	for _, h := range e.d.PostLogoutHooks(ctx) {
		if err := h.ExecuteLogoutPostHook(ctx, f); err != nil {
			return err
		}
	}
	return nil
}
