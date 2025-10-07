package oteltag

import (
	"reflect"
	"strings"

	"go.opentelemetry.io/otel/baggage"

	"github.com/remychantenay/otel-tag/internal"
)

// BaggageMembers takes in a struct and spits out OpenTelemetry baggage members
// based on the struct tags. Only basic types are supported: string, int, int64, float64, bool,
// and slices of these types. Nested structs and pointers are ignored.
func BaggageMembers(res any) []baggage.Member {
	structValue := reflect.ValueOf(res)
	structType := structValue.Type()

	if structType.Kind() != reflect.Struct {
		return nil
	}

	fieldCount := structValue.NumField()
	if fieldCount == 0 {
		return nil
	}

	members := make([]baggage.Member, 0, fieldCount)
	for i := range fieldCount {
		member := basicTypeToBaggageMember(structType, structValue, i)
		if member.Key() == "" {
			continue
		}

		members = append(members, member)
	}

	return members
}

// basicTypeToBaggageMember returns an [attribute.Member] for a basic type.
func basicTypeToBaggageMember(structType reflect.Type, structValue reflect.Value, index int) baggage.Member {
	tag := internal.ExtractTag(structValue, index)

	if tag == "" {
		return baggage.Member{}
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
	attr, zeroValue := internal.BaggageMember(field, fieldValue, tag)
	if zeroValue && omitEmpty {
		return baggage.Member{}
	}

	return attr
}
