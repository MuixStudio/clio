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
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entrans "github.com/go-playground/validator/v10/translations/en"
	"github.com/muixstudio/clio/internal/infra/errors"
)

type Validator struct {
	trans    ut.Translator
	validate *validator.Validate
}

func NewValidator() (*Validator, error) {
	v := validator.New()

	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ := uni.GetTranslator("en")
	if err := entrans.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, err
	}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})

	registerCustomValidations(v, trans)

	return &Validator{validate: v, trans: trans}, nil
}

func registerCustomValidations(v *validator.Validate, trans ut.Translator) {
	registerTag(v, trans, "identifier",
		"{0} may only contain letters, digits, and - _ . @",
		func(fl validator.FieldLevel) bool {
			for _, r := range fl.Field().String() {
				if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' && r != '.' && r != '@' {
					return false
				}
			}
			return true
		},
	)
	registerTag(v, trans, "password",
		"{0} may only contain printable ASCII characters (no spaces)",
		func(fl validator.FieldLevel) bool {
			for _, r := range fl.Field().String() {
				if r < 0x21 || r > 0x7E {
					return false
				}
			}
			return true
		},
	)
}

func registerTag(v *validator.Validate, trans ut.Translator, tag, msg string, fn validator.Func) {
	_ = v.RegisterValidation(tag, fn)
	_ = v.RegisterTranslation(tag, trans,
		func(ut ut.Translator) error { return ut.Add(tag, msg, true) },
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(tag, fe.Field())
			return t
		},
	)
}

type SliceValidationError []error

// Error concatenates all error elements in SliceValidationError into a single string separated by \n.
func (err SliceValidationError) Error() string {
	if len(err) == 0 {
		return ""
	}

	var b strings.Builder
	for i := range len(err) {
		if err[i] != nil {
			if b.Len() > 0 {
				b.WriteString("\n")
			}
			fmt.Fprintf(&b, "[%d]: %s", i, err[i].Error())
		}
	}
	return b.String()
}

//var _ StructValidator = (*defaultValidator)(nil)

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *Validator) ValidateStruct(obj any) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Pointer:
		if value.Elem().Kind() != reflect.Struct {
			return v.ValidateStruct(value.Elem().Interface())
		}
		return v.validateStruct(obj)
	case reflect.Struct:
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(SliceValidationError, 0)
		for i := range count {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

// validateStruct receives struct type
func (v *Validator) validateStruct(obj any) error {
	err := v.validate.Struct(obj)
	if err == nil {
		return nil
	}
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if ok {
		translated := errs.Translate(v.trans)
		msgs := make([]string, 0, len(errs))
		for _, e := range errs {
			msg, ok := translated[e.Namespace()]
			if !ok {
				continue
			}
			ns := e.Namespace()
			if idx := strings.Index(ns, "."); idx != -1 {
				ns = ns[idx+1:]
			}
			msg = strings.Replace(msg, e.Field(), ns, 1)
			msgs = append(msgs, msg)
		}
		return errors.BadRequest("INVALID_ARGUMENT", strings.Join(msgs, "; ")).WithCause(err)
	}

	return errors.InternalServerError("INTERNAL_SERVER_ERROR", "parsing failed").WithCause(err)
}

func (v *Validator) Engine() any {
	return v.validate
}
