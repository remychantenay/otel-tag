package internal

import (
	"reflect"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
)

// tagName used by this library.
const tagName = "otel"

// ExtractTag extract the tag's value of a given Struct field index.
func ExtractTag(value reflect.Value, index int) string {
	return value.Type().Field(index).Tag.Get(tagName)
}

// SpanAttribute creates and returns an OpenTelemetry span attribute for the provided field.
// Also returns a boolean that indicates whether or not the field's value is a zero-value.
func SpanAttribute(field reflect.StructField, fieldValue reflect.Value, attrKey string) (attribute.KeyValue, bool) {
	switch field.Type.Kind() {
	case reflect.String:
		v := fieldValue.String()
		return attribute.String(attrKey, v), v == ""
	case reflect.Int:
		v := fieldValue.Int()
		return attribute.Int(attrKey, int(v)), v == 0
	case reflect.Int64:
		v := fieldValue.Int()
		return attribute.Int64(attrKey, v), v == 0
	case reflect.Float64:
		v := fieldValue.Float()
		return attribute.Float64(attrKey, v), v == 0.0
	case reflect.Bool:
		return attribute.Bool(attrKey, fieldValue.Bool()), false
	case reflect.Slice:
		switch field.Type.Elem().Kind() {
		case reflect.String:
			if s, ok := fieldValue.Interface().([]string); ok {
				return attribute.StringSlice(attrKey, s), len(s) == 0
			}
		case reflect.Int:
			if s, ok := fieldValue.Interface().([]int); ok {
				return attribute.IntSlice(attrKey, s), len(s) == 0
			}
		case reflect.Int64:
			if s, ok := fieldValue.Interface().([]int64); ok {
				return attribute.Int64Slice(attrKey, s), len(s) == 0
			}
		case reflect.Float64:
			if s, ok := fieldValue.Interface().([]float64); ok {
				return attribute.Float64Slice(attrKey, s), len(s) == 0
			}
		case reflect.Bool:
			if s, ok := fieldValue.Interface().([]bool); ok {
				return attribute.BoolSlice(attrKey, s), len(s) == 0
			}
		}
	}

	return attribute.KeyValue{}, true
}

// BaggageMember creates and returns an OpenTelemetry baggageMember for the provided field.
// Also returns a boolean that indicates whether or not the field's value is a zero-value.
func BaggageMember(field reflect.StructField, fieldValue reflect.Value, memberKey string) (baggage.Member, bool) {
	switch field.Type.Kind() {
	case reflect.String:
		v := fieldValue.String()
		m, _ := baggage.NewMemberRaw(memberKey, v)
		return m, v == ""
	case reflect.Int:
		v := fieldValue.Int()
		m, _ := baggage.NewMemberRaw(memberKey, strconv.Itoa(int(v)))
		return m, v == 0
	case reflect.Int64:
		v := fieldValue.Int()
		m, _ := baggage.NewMemberRaw(memberKey, strconv.FormatInt(v, 10))
		return m, v == 0
	case reflect.Float64:
		v := fieldValue.Float()
		m, _ := baggage.NewMemberRaw(memberKey, strconv.FormatFloat(v, 'f', -1, 64))
		return m, v == 0.0
	case reflect.Bool:
		m, _ := baggage.NewMemberRaw(memberKey, strconv.FormatBool(fieldValue.Bool()))
		return m, false
	case reflect.Slice:
		switch field.Type.Elem().Kind() {
		case reflect.String:
			if s, ok := fieldValue.Interface().([]string); ok {
				m, _ := baggage.NewMemberRaw(memberKey, strings.Join(s, ","))
				return m, len(s) == 0
			}
		case reflect.Int:
			if s, ok := fieldValue.Interface().([]int); ok {
				sStr := make([]string, len(s))
				for i, v := range s {
					sStr[i] = strconv.Itoa(int(v))
				}
				m, _ := baggage.NewMemberRaw(memberKey, strings.Join(sStr, ","))
				return m, len(s) == 0
			}
		case reflect.Int64:
			if s, ok := fieldValue.Interface().([]int64); ok {
				sStr := make([]string, len(s))
				for i, v := range s {
					sStr[i] = strconv.FormatInt(v, 10)
				}
				m, _ := baggage.NewMemberRaw(memberKey, strings.Join(sStr, ","))
				return m, len(s) == 0
			}
		case reflect.Float64:
			if s, ok := fieldValue.Interface().([]float64); ok {
				sStr := make([]string, len(s))
				for i, v := range s {
					sStr[i] = strconv.FormatFloat(v, 'f', -1, 64)
				}
				m, _ := baggage.NewMemberRaw(memberKey, strings.Join(sStr, ","))
				return m, len(s) == 0
			}
		case reflect.Bool:
			if s, ok := fieldValue.Interface().([]bool); ok {
				sStr := make([]string, len(s))
				for i, v := range s {
					sStr[i] = strconv.FormatBool(v)
				}
				m, _ := baggage.NewMemberRaw(memberKey, strings.Join(sStr, ","))
				return m, len(s) == 0
			}
		}
	}

	return baggage.Member{}, true
}
