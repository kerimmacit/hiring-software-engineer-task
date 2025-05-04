package validation

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Validate(v interface{}) error {
	return validate.Struct(v)
}

func GetBaseValidator() *validator.Validate {
	if validate == nil {
		validate = validator.New()
	}
	return validate
}
