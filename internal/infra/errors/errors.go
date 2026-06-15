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

import (
	"errors"
	"fmt"
)

const (
	UnknownCode   = 500
	UnknownReason = "UNKNOWN_REASON"
)

type Error struct {
	// http response status code
	code int32

	reason  string
	message string

	metadata map[string]string

	cause error
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v cause = %v", e.code, e.reason, e.message, e.metadata, e.cause)
}

func (e *Error) WithCause(cause error) *Error {
	err := Clone(e)
	err.cause = cause
	return err
}

func (e *Error) WithMetadata(md map[string]string) *Error {
	err := Clone(e)
	err.metadata = md
	return err
}

func (e *Error) GetCode() int32 {
	return e.code
}

func (e *Error) GetMessage() string {
	return e.message
}

func (e *Error) GetReason() string {
	return e.reason
}

func (e *Error) GetMetadata() map[string]string {
	return e.metadata
}

// New returns an error object for the code, message.
func New(code int32, reason, message string) *Error {
	return &Error{
		code:    code,
		message: message,
		reason:  reason,
	}
}

func Clone(err *Error) *Error {
	if err == nil {
		return nil
	}
	metadata := make(map[string]string, len(err.metadata))
	for k, v := range err.metadata {
		metadata[k] = v
	}
	return &Error{
		cause:    err.cause,
		code:     err.code,
		reason:   err.reason,
		message:  err.message,
		metadata: metadata,
	}
}

func ToError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	ne := New(UnknownCode, UnknownReason, err.Error())
	return ne
}
