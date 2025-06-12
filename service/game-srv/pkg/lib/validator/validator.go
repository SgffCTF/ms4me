package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func GetDetailedError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		firstError := validationErrors[0]

		switch firstError.Tag() {
		case "required":
			return fmt.Errorf("Поле %s необходимо", firstError.Field())
		case "min":
			return fmt.Errorf("Поле %s должно содержать минимум %s символов", firstError.Field(), firstError.Param())
		case "gt":
			return fmt.Errorf("Поле %s должно быть больше %s", firstError.Field(), firstError.Param())
		default:
			return fmt.Errorf("Поле %s не верно", firstError.Field())
		}
	}
	return fmt.Errorf("Внутренняя ошибка")
}

func Validate(s interface{}) error {
	return validator.New().Struct(s)
}
