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

package orgctx

import (
	"context"

	"github.com/google/uuid"
)

type contextKey int

const orgIDKey contextKey = iota

func WithOrgID(ctx context.Context, id uuid.NullUUID) context.Context {
	return context.WithValue(ctx, orgIDKey, id)
}

func OrgIDFromCtx(ctx context.Context) (uuid.NullUUID, bool) {
	id, ok := ctx.Value(orgIDKey).(uuid.NullUUID)
	return id, ok
}
