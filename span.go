package oteltag

import (
	"reflect"
	"strings"

	"go.opentelemetry.io/otel/attribute"

	"github.com/remychantenay/otel-tag/internal"
)

// SpanAttributes takes in a struct and spits out OpenTelemetry span attributes ([attribute.KeyValue])
// based on the struct tags.
func SpanAttributes(res any) []attribute.KeyValue {
	return structToAttributes(res)
}

// structToAttributes returns a slice of [attribute.KeyValue] for a struct.
func structToAttributes(s any) []attribute.KeyValue {
	structValue := reflect.ValueOf(s)
	structType := structValue.Type()

	// Handle pointer to struct
	if structType.Kind() == reflect.Pointer {
		if structValue.IsNil() {
			return nil
		}
		structValue = structValue.Elem()
		structType = structValue.Type()
	}

	// Validate input is a struct
	if structType.Kind() != reflect.Struct {
		return nil
	}

	fieldCount := structValue.NumField()
	if fieldCount == 0 {
		return nil
	}

	attrs := make([]attribute.KeyValue, 0, fieldCount)
	for i := 0; i < fieldCount; i++ {
		field, fieldValue := structType.Field(i), structValue.Field(i)
		if field.Type.Kind() == reflect.Struct && fieldValue.IsValid() {
			attrs = append(attrs, structToAttributes(fieldValue.Interface())...)
		} else if field.Type.Kind() == reflect.Pointer && fieldValue.IsValid() && !fieldValue.IsNil() { // Known shortcoming, assuming a pointer can only be a struct.
			attrs = append(attrs, structToAttributes(fieldValue.Elem().Interface())...)
		} else {
			attr := basicTypeToAttribute(structType, structValue, i)
			if !attr.Valid() {
				continue
			}

			attrs = append(attrs, attr)
		}
	}

	return attrs
}

// basicTypeToAttribute returns an [attribute.KeyValue] for a basic type.
func basicTypeToAttribute(structType reflect.Type, structValue reflect.Value, index int) attribute.KeyValue {
	tag := internal.ExtractTag(structValue, index)

	if tag == "" {
		return attribute.KeyValue{}
	}

	var omitEmpty bool
	before, after, found := strings.Cut(tag, ",")
	if found {
		tag = before
		if after == flagOmitEmpty {
			omitEmpty = true
		}
	}

	field, fieldValue := structType.Field(index), structValue.Field(index)
	attr, zeroValue := internal.SpanAttribute(field, fieldValue, tag)
	if zeroValue && omitEmpty {
		return attribute.KeyValue{}
	}

	return attr
}
