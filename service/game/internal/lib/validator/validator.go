package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetDetailedError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		firstError := validationErrors[0]

		switch firstError.Tag() {
		case "required":
			return status.Error(codes.InvalidArgument, fmt.Sprintf("Поле %s необходимо", firstError.Field()))
		case "min":
			return status.Error(codes.InvalidArgument, fmt.Sprintf("Поле %s должно содержать минимум %s символов", firstError.Field(), firstError.Param()))
		case "gt":
			return status.Error(codes.InvalidArgument, fmt.Sprintf("Поле %s должно быть больше %s", firstError.Field(), firstError.Param()))
		default:
			return status.Error(codes.InvalidArgument, fmt.Sprintf("Поле %s не верно", firstError.Field()))
		}
	}
	return status.Error(codes.Internal, "Внутренняя ошибка")
}

func Validate(s interface{}) error {
	return validator.New().Struct(s)
}
