package internal

import (
	"reflect"
)

// tagName used by this library.
const tagName = "otel"

// ExtractTag extract the tag's value of a given Struct field index.
func ExtractTag(value reflect.Value, index int) string {
	return value.Type().Field(index).Tag.Get(tagName)
}
