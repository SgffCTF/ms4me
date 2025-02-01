package validator

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func GetDetailedError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		firstError := validationErrors[0]
		kind := firstError.Kind()
		switch firstError.Tag() {
		case "max":
			if kind == reflect.String {
				return fmt.Sprintf("field '%s' require maximum %s characters", firstError.Field(), firstError.Param())
			}
			return fmt.Sprintf("field '%s' should be less then %s", firstError.Field(), firstError.Param())
		case "min":
			if kind == reflect.String {
				return fmt.Sprintf("field '%s' require minimum %s characters", firstError.Field(), firstError.Param())
			}
			return fmt.Sprintf("field '%s' should be greater then %s", firstError.Field(), firstError.Param())
		case "required":
			return fmt.Sprintf("field '%s' is required", firstError.Field())
		default:
			return fmt.Sprintf("field '%s' is invalid", firstError.Field())
		}
	}
	return err.Error()
}
