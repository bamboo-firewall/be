package httpbase

import "github.com/go-playground/validator/v10"

var (
	defaultValidator = validator.New()
)

func RegisterValidator(tag string, fn validator.Func) error {
	return defaultValidator.RegisterValidation(tag, fn)
}

func RegisterStructValidation(fn validator.StructLevelFunc, types ...interface{}) {
	defaultValidator.RegisterStructValidation(fn, types...)
}
