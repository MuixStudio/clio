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

package binding

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	ejson "encoding/json"

	"github.com/muixstudio/clio/internal/infra/errors"

	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/codec/json"
)

var (
	JSON binding.BindingBody = jsonBinding{}
)

type jsonBinding struct{}

func (jb jsonBinding) Name() string {
	return "json"
}

func (jb jsonBinding) Bind(req *http.Request, obj any) error {
	if req == nil || req.Body == nil {
		return errors.BadRequest("INVALID_REQUEST", "invalid request")
	}
	return jb.decodeJSON(req.Body, obj)
}

func (jb jsonBinding) BindBody(body []byte, obj any) error {
	return jb.decodeJSON(bytes.NewReader(body), obj)
}

func (jb jsonBinding) decodeJSON(r io.Reader, obj any) error {
	decoder := json.API.NewDecoder(r)
	if binding.EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if binding.EnableDecoderDisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}

	if err := decoder.Decode(obj); err != nil {
		var syntaxErr *ejson.SyntaxError
		if errors.As(err, &syntaxErr) {
			return errors.BadRequest(
				"INVALID_JSON",
				fmt.Sprintf("invalid JSON syntax at offset %d", syntaxErr.Offset),
			).WithCause(err).WithMetadata(map[string]string{
				"offset": fmt.Sprintf("%d", syntaxErr.Offset),
			})
		}
		var unmarshalErr *ejson.UnmarshalTypeError
		if errors.As(err, &unmarshalErr) {
			return errors.BadRequest(
				"INVALID_JSON_FORMAT",
				fmt.Sprintf("field %q expects %s", unmarshalErr.Field, unmarshalErr.Type),
			).WithCause(err)
		}
		if errors.Is(err, binding.ErrConvertMapStringSlice) ||
			errors.Is(err, binding.ErrConvertToMapString) {
			return errors.BadRequest("INVALID_JSON_FORMAT", "invalid JSON format").WithCause(err)
		}
		return err
	}

	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}
