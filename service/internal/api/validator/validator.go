package validator

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator структура для кастомного валидатора, использующего библиотеку validator
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate метод, который Echo вызывает для валидации данных
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
