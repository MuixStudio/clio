package binding

type StructValidator interface {
	ValidateStruct(any) error
}

var Validator StructValidator = &defaultValidator{}

func Validate(obj any) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}
