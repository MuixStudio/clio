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

package errors

// ── 4xx ──────────────────────────────────────────────────────────────────────

// BadRequest returns a 400 error.
func BadRequest(reason, message string) *Error {
	return New(400, reason, message)
}

// Unauthorized returns a 401 error.
func Unauthorized(reason, message string) *Error {
	return New(401, reason, message)
}

// Forbidden returns a 403 error.
func Forbidden(reason, message string) *Error {
	return New(403, reason, message)
}

// NotFound returns a 404 error.
func NotFound(reason, message string) *Error {
	return New(404, reason, message)
}

// Conflict returns a 409 error.
func Conflict(reason, message string) *Error {
	return New(409, reason, message)
}

// Gone returns a 410 error.
func Gone(reason, message string) *Error {
	return New(410, reason, message)
}

// ── 5xx ──────────────────────────────────────────────────────────────────────

// InternalServerError returns a 500 error.
func InternalServerError(reason, message string) *Error {
	return New(500, reason, message)
}

// BadGateway returns a 502 error.
func BadGateway(reason, message string) *Error {
	return New(502, reason, message)
}
