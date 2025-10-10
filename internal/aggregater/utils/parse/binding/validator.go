package binding

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entrans "github.com/go-playground/validator/v10/translations/en"
)

type defaultValidator struct {
	once       sync.Once
	validate   *validator.Validate
	translator ut.Translator
}

type SliceValidationError []error

// Error concatenates all error elements in SliceValidationError into a single string separated by \n.
func (err SliceValidationError) Error() string {
	if len(err) == 0 {
		return ""
	}

	var b strings.Builder
	for i := 0; i < len(err); i++ {
		if err[i] != nil {
			if b.Len() > 0 {
				b.WriteString("\n")
			}
			b.WriteString("[" + strconv.Itoa(i) + "]: " + err[i].Error())
		}
	}
	return b.String()
}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *defaultValidator) ValidateStruct(obj any) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		if value.Elem().Kind() != reflect.Struct {
			return v.ValidateStruct(value.Elem().Interface())
		}
		return v.validateStruct(obj)
	case reflect.Struct:
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(SliceValidationError, 0)
		for i := 0; i < count; i++ {
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
func (v *defaultValidator) validateStruct(obj any) error {
	v.lazyinit()
	err := v.validate.Struct(obj)
	var VErrs validator.ValidationErrors
	ok := errors.As(err, &VErrs)
	if ok {
		var b strings.Builder
		for i, e := range VErrs {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(e.Translate(v.translator))
		}
		return errors.New(b.String())
	}
	return err
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {

		enLocale := en.New()
		uni := ut.New(enLocale, enLocale)
		v.translator, _ = uni.GetTranslator("en")

		v.validate = validator.New()
		entrans.RegisterDefaultTranslations(v.validate, v.translator)
		v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	})
}
