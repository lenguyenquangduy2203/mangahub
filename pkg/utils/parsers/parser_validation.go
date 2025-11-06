package parsers

import (
	"fmt"
	"mangahub/pkg/models/dtos"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ParseValidationErrors(err error) []dtos.ErrorDetail {
	var details []dtos.ErrorDetail

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			msg := translateValidationError(e)
			details = append(details, dtos.ErrorDetail{
				Field: strings.ToLower(e.Field()),
				Issue: msg,
			})
		}
	} else {
		// fallback for non-validation errors
		details = append(details, dtos.ErrorDetail{
			Field: "unknown",
			Issue: err.Error(),
		})
	}

	return details
}

func translateValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s characters", e.Field(), e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", e.Field(), e.Param())
	default:
		// fallback for unknown tag
		return fmt.Sprintf("%s is invalid (%s)", e.Field(), e.Tag())
	}
}
