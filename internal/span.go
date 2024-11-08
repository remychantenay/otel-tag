package internal

import (
	"reflect"

	"go.opentelemetry.io/otel/attribute"
)

// OTelAttribute creates and returns an OpenTelemetry span attribute for the provided field.
// Also returns a boolean that indicates whether or not the field's value is a zero-value.
func OTelAttribute(field reflect.StructField, value reflect.Value, attrKey string) (attribute.KeyValue, bool) {

	switch field.Type.Kind() {
	case reflect.String:
		v := value.String()
		return attribute.String(attrKey, v), v == ""
	case reflect.Int:
		v := value.Int()
		return attribute.Int(attrKey, int(v)), v == 0
	case reflect.Int64:
		v := value.Int()
		return attribute.Int64(attrKey, v), v == 0
	case reflect.Float32, reflect.Float64:
		v := value.Float()
		return attribute.Float64(attrKey, v), v == 0.0
	case reflect.Bool:
		return attribute.Bool(attrKey, value.Bool()), false
	case reflect.Slice:
		switch field.Type.Elem().Kind() {
		case reflect.String:
			if s, ok := value.Interface().([]string); ok {
				return attribute.StringSlice(attrKey, s), len(s) == 0
			}
		case reflect.Int:
			if s, ok := value.Interface().([]int); ok {
				return attribute.IntSlice(attrKey, s), len(s) == 0
			}
		case reflect.Int64:
			if s, ok := value.Interface().([]int64); ok {
				return attribute.Int64Slice(attrKey, s), len(s) == 0
			}
		case reflect.Float32:
			if s, ok := value.Interface().([]float32); ok {
				s64 := make([]float64, len(s))
				for i, v := range s {
					s64[i] = float64(v)
				}

				return attribute.Float64Slice(attrKey, s64), len(s64) == 0
			}
		case reflect.Float64:
			if s, ok := value.Interface().([]float64); ok {
				return attribute.Float64Slice(attrKey, s), len(s) == 0
			}
		case reflect.Bool:
			if s, ok := value.Interface().([]bool); ok {
				return attribute.BoolSlice(attrKey, s), len(s) == 0
			}
		}
	}

	return attribute.KeyValue{}, true
}
