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

package registration

import (
	"context"

	"github.com/muixstudio/clio/internal/flow"
	"github.com/muixstudio/clio/internal/identity"
)

// PreHookExecutor runs before the registration flow begins processing.
// Use it for: rate limiting, bot detection, IP allow-listing, ...
// Accepts flow.Flow so hooks stay decoupled from the concrete *Flow type.
type PreHookExecutor interface {
	ExecuteRegistrationPreHook(ctx context.Context, f flow.Flow) error
}

// PostHookExecutor runs after the identity is created but before the session
// is issued. Use it for: welcome emails, audit logging, analytics events, ...
type PostHookExecutor interface {
	ExecuteRegistrationPostHook(ctx context.Context, f flow.Flow, i *identity.Identity) error
}

// HooksProvider is implemented by the Driver.
type HooksProvider interface {
	PreRegistrationHooks(ctx context.Context) []PreHookExecutor
	PostRegistrationHooks(ctx context.Context) []PostHookExecutor
}

// HookExecutorProvider allows the Handler to retrieve the HookExecutor.
type HookExecutorProvider interface {
	RegistrationHookExecutor() *HookExecutor
}

// HookExecutor runs pre- and post-hooks in order, stopping on first error.
// It is the only place where hooks are called — the Handler never calls hooks directly.
type HookExecutor struct{ d HooksProvider }

func NewHookExecutor(d HooksProvider) *HookExecutor { return &HookExecutor{d: d} }

func (e *HookExecutor) PreHook(ctx context.Context, f flow.Flow) error {
	for _, h := range e.d.PreRegistrationHooks(ctx) {
		if err := h.ExecuteRegistrationPreHook(ctx, f); err != nil {
			return err
		}
	}
	return nil
}

func (e *HookExecutor) PostHook(ctx context.Context, f flow.Flow, i *identity.Identity) error {
	for _, h := range e.d.PostRegistrationHooks(ctx) {
		if err := h.ExecuteRegistrationPostHook(ctx, f, i); err != nil {
			return err
		}
	}
	return nil
}
