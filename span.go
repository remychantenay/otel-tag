package oteltag

import (
	"reflect"
	"strings"

	"go.opentelemetry.io/otel/attribute"

	"github.com/remychantenay/otel-tag/internal"
)

// SpanAttributes takes in a struct and spits out OpenTelemetry span attributes ([attribute.KeyValue])
// based on the struct tags. Only basic types are supported: string, int, int64, float64, bool,
// and slices of these types. Nested structs and pointers are ignored.
func SpanAttributes(res any) []attribute.KeyValue {
	structValue := reflect.ValueOf(res)
	structType := structValue.Type()

	if structType.Kind() != reflect.Struct {
		return nil
	}

	fieldCount := structValue.NumField()
	if fieldCount == 0 {
		return nil
	}

	attrs := make([]attribute.KeyValue, 0, fieldCount)
	for i := range fieldCount {
		attr := basicTypeToAttribute(structType, structValue, i)
		if !attr.Valid() {
			continue
		}

		attrs = append(attrs, attr)
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
