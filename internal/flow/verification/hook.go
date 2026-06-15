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

package verification

import "context"

// PreHookExecutor is implemented by hooks that run before a verification flow
// is created (e.g. rate limiting, bot detection).
type PreHookExecutor interface {
	ExecuteVerificationPreHook(ctx context.Context, f *Flow) error
}

// PostHookExecutor is implemented by hooks that run after an address is
// successfully verified (e.g. grant a feature flag, fire an analytics event).
type PostHookExecutor interface {
	ExecuteVerificationPostHook(ctx context.Context, f *Flow) error
}

// HooksProvider is satisfied by the Driver and supplies the hook lists.
type HooksProvider interface {
	PreVerificationHooks(ctx context.Context) []PreHookExecutor
	PostVerificationHooks(ctx context.Context) []PostHookExecutor
}

// HookExecutorProvider lets the handler resolve the executor without depending
// on a concrete type. Mirrors the pattern used by login and registration.
type HookExecutorProvider interface {
	VerificationHookExecutor() *HookExecutor
}

// HookExecutor runs the pre- and post-hook chains for verification flows.
type HookExecutor struct{ d HooksProvider }

func NewHookExecutor(d HooksProvider) *HookExecutor { return &HookExecutor{d: d} }

// PreHook runs all registered pre-verification hooks in order.
// Returns the first error encountered; remaining hooks are not run.
func (e *HookExecutor) PreHook(ctx context.Context, f *Flow) error {
	for _, h := range e.d.PreVerificationHooks(ctx) {
		if err := h.ExecuteVerificationPreHook(ctx, f); err != nil {
			return err
		}
	}
	return nil
}

// PostHook runs all registered post-verification hooks in order.
func (e *HookExecutor) PostHook(ctx context.Context, f *Flow) error {
	for _, h := range e.d.PostVerificationHooks(ctx) {
		if err := h.ExecuteVerificationPostHook(ctx, f); err != nil {
			return err
		}
	}
	return nil
}
