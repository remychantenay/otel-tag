package oteltag

import (
	"reflect"
	"strings"

	"go.opentelemetry.io/otel/baggage"

	"github.com/remychantenay/otel-tag/internal"
)

// BaggageMembers takes in a struct and spits out OpenTelemetry baggage members
// based on the struct tags.
func BaggageMembers(res any) []baggage.Member {
	return structToBaggageMembers(res)
}

// structToBaggageMembers returns a slice of [baggage.Member] for a struct.
func structToBaggageMembers(s any) []baggage.Member {
	structValue := reflect.ValueOf(s)
	structType := structValue.Type()
	fieldCount := structValue.NumField()

	if fieldCount == 0 {
		return nil
	}

	members := make([]baggage.Member, 0, fieldCount)
	for i := 0; i < fieldCount; i++ {
		field, fieldValue := structType.Field(i), structValue.Field(i)
		if field.Type.Kind() == reflect.Struct && fieldValue.IsValid() {
			members = append(members, structToBaggageMembers(fieldValue.Interface())...)
		} else if field.Type.Kind() == reflect.Pointer && fieldValue.IsValid() { // Known shortcoming, assuming a pointer can only be a struct.
			members = append(members, structToBaggageMembers(fieldValue.Elem().Interface())...)
		} else {
			member := basicTypeToBaggageMember(structType, structValue, i)
			if member.Key() == "" {
				continue
			}

			members = append(members, member)
		}
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
