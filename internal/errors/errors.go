package errors

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ParseFieldErrors parses Field and Detail from fieldErrs, returning
// a comma separated string of the field errors.
func ParseFieldErrors(fieldErrs field.ErrorList) string {
	if fieldErrs == nil {
		return ""
	}

	var errs []string
	for _, err := range fieldErrs {
		errs = append(errs, fmt.Sprintf("%v for %s; %s", err.Type, err.Field, err.Detail))
	}
	if len(errs) > 0 {
		return strings.Join(errs, ". ")
	}

	return ""
}
