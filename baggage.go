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
	structValue := reflect.ValueOf(res)
	structType := structValue.Type()
	fieldCount := structValue.NumField()

	if fieldCount == 0 {
		return nil
	}

	members := make([]baggage.Member, 0, fieldCount)
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

		member, zeroValue := internal.BaggageMember(field, value, tag)
		if !zeroValue || !omitEmpty {
			members = append(members, member)
		}
	}

	return members
}
