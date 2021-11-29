package srvc_errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrInternalError = errors.New("internal error")
)

type ValidationErrors struct {
	Fields map[string]string
}

func (e ValidationErrors) Error() string {
	if len(e.Fields) == 0 {
		return ""
	}

	var s = make([]string, 0, len(e.Fields))
	for f, v := range e.Fields {
		s = append(s, fmt.Sprintf("%s: %s", f, v))
	}
	return strings.Join(s, "; ")
}

func unwrapValidationError(errs validation.Errors) map[string]interface{} {
	fields := map[string]interface{}{}
	for k, v := range errs {
		if verrs, ok := v.(validation.Errors); ok {
			fields[k] = unwrapValidationError(verrs)
			continue
		}
		fields[k] = v.Error()
	}
	return fields
}

func ValidationError(err error) error {
	var (
		vie  validation.InternalError
		errs validation.Errors
	)

	if errors.As(err, &vie) {
		return ErrInternalError
	} else if errors.As(err, &errs) {
		fields := MapToFlatMap(unwrapValidationError(errs), nil)
		v := ValidationErrors{Fields: make(map[string]string)}
		for k, f := range fields {
			v.Fields[k] = f.(string)
		}
		return v
	}

	return ErrInternalError
}

func MapToFlatMap(in map[string]interface{}, keyMapper func(fname string) string) map[string]interface{} {
	out := make(map[string]interface{})
	for key, value := range in {
		if keyMapper != nil {
			key = keyMapper(key)
		}
		if valueAsMap, ok := value.(map[string]interface{}); ok {
			sub := MapToFlatMap(valueAsMap, keyMapper)
			for subKey, subValue := range sub {
				out[key+"."+subKey] = subValue
			}
			continue
		}
		out[key] = value
	}
	return out
}
