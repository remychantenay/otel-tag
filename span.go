package oteltag

import (
	"reflect"
	"strings"

	"go.opentelemetry.io/otel/attribute"

	"github.com/remychantenay/otel-tag/internal"
)

// SpanAttributes takes in a struct and spits out OpenTelemetry span attributes
// based on the struct tags.
func SpanAttributes(res any) []attribute.KeyValue {
	structValue := reflect.ValueOf(res)
	structType := structValue.Type()
	fieldCount := structValue.NumField()

	if fieldCount == 0 {
		return nil
	}

	attrs := make([]attribute.KeyValue, 0, fieldCount)
	for i := 0; i < fieldCount; i++ {
		tag := internal.ExtractTag(structValue, i)

		if len(tag) == 0 || tag == valIgnore {
			continue
		}

		var omitEmpty bool
		before, after, found := strings.Cut(tag, ",")
		if found {
			tag = before
			if after == valOmitEmpty {
				omitEmpty = true
			}
		}

		field, value := structType.Field(i), structValue.Field(i)

		attr, zeroValue := internal.OTelAttribute(field, value, tag)
		if !zeroValue || !omitEmpty {
			attrs = append(attrs, attr)
		}
	}

	return attrs
}
