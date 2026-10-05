package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/97Dmitry/go-todolist/internal/apperrors"
	"github.com/go-playground/validator/v10"
)

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Report JSON field names in validation errors instead of Go field names.
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}

		return name
	})

	return v
}

// DecodeAndValidate decodes a single JSON object from the request body into
// dest and validates it with `validate` struct tags.
func DecodeAndValidate(r *http.Request, dest any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dest); err != nil {
		return decodeError(err)
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err != nil {
			return decodeError(err)
		}

		return apperrors.InvalidArguments("request body must contain a single JSON object")
	}

	if err := validate.Struct(dest); err != nil {
		validationErrors, ok := errors.AsType[validator.ValidationErrors](err)
		if !ok {
			return fmt.Errorf("validate request: %w", err)
		}

		return fmt.Errorf("validate request: %w: %w", apperrors.InvalidArguments(validationMessage(validationErrors)), err)
	}

	return nil
}

func decodeError(err error) error {
	if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
		return fmt.Errorf("decode request body: %w", err)
	}

	var message string

	if typeErr, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
		message = fmt.Sprintf("field %q has invalid type", typeErr.Field)
	} else if field, ok := strings.CutPrefix(err.Error(), "json: unknown field "); ok {
		message = "unknown field " + field
	} else if errors.Is(err, io.EOF) {
		message = "request body is empty"
	} else {
		message = "request body is not valid JSON"
	}

	return fmt.Errorf("decode request body: %w: %w", apperrors.InvalidArguments(message), err)
}

func validationMessage(validationErrors validator.ValidationErrors) string {
	messages := make([]string, 0, len(validationErrors))

	for _, fieldErr := range validationErrors {
		messages = append(messages, fieldMessage(fieldErr))
	}

	return strings.Join(messages, "; ")
}

func fieldMessage(fieldErr validator.FieldError) string {
	field := fieldErr.Field()
	unit := ""
	if fieldErr.Kind() == reflect.String {
		unit = " characters"
	}

	switch fieldErr.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "min":
		return fmt.Sprintf("%s must be at least %s%s", field, fieldErr.Param(), unit)
	case "max":
		return fmt.Sprintf("%s must be at most %s%s", field, fieldErr.Param(), unit)
	default:
		return field + " is invalid"
	}
}
