package api

import (
	"fmt"
	"reflect"
)

func requiredFields(fields map[string]any) error {
	for k, v := range fields {
		if err := requiredField(k, v); err != nil {
			return err
		}
	}

	return nil
}

func requiredField(fieldName string, v any) error {
	if reflect.ValueOf(v).IsZero() {
		return fmt.Errorf("%s is required", fieldName)
	}

	return nil
}
