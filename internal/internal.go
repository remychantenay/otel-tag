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
		v := fieldValue.Bool()
		return attribute.Bool(attrKey, v), !v
	case reflect.Slice:
		switch field.Type.Elem().Kind() {
		case reflect.String:
			s := fieldValue.Interface().([]string)
			return attribute.StringSlice(attrKey, s), len(s) == 0
		case reflect.Int:
			s := fieldValue.Interface().([]int)
			return attribute.IntSlice(attrKey, s), len(s) == 0
		case reflect.Int64:
			s := fieldValue.Interface().([]int64)
			return attribute.Int64Slice(attrKey, s), len(s) == 0
		case reflect.Float64:
			s := fieldValue.Interface().([]float64)
			return attribute.Float64Slice(attrKey, s), len(s) == 0
		case reflect.Bool:
			s := fieldValue.Interface().([]bool)
			return attribute.BoolSlice(attrKey, s), len(s) == 0
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
		m, err := baggage.NewMemberRaw(memberKey, v)
		if err != nil {
			return baggage.Member{}, true
		}
		return m, v == ""
	case reflect.Int:
		v := fieldValue.Int()
		m, err := baggage.NewMemberRaw(memberKey, strconv.Itoa(int(v)))
		if err != nil {
			return baggage.Member{}, true
		}
		return m, v == 0
	case reflect.Int64:
		v := fieldValue.Int()
		m, err := baggage.NewMemberRaw(memberKey, strconv.FormatInt(v, 10))
		if err != nil {
			return baggage.Member{}, true
		}
		return m, v == 0
	case reflect.Float64:
		v := fieldValue.Float()
		m, err := baggage.NewMemberRaw(memberKey, strconv.FormatFloat(v, 'f', -1, 64))
		if err != nil {
			return baggage.Member{}, true
		}
		return m, v == 0.0
	case reflect.Bool:
		v := fieldValue.Bool()
		m, err := baggage.NewMemberRaw(memberKey, strconv.FormatBool(v))
		if err != nil {
			return baggage.Member{}, true
		}
		return m, !v
	case reflect.Slice:
		switch field.Type.Elem().Kind() {
		case reflect.String:
			s := fieldValue.Interface().([]string)
			m, err := baggage.NewMemberRaw(memberKey, strings.Join(s, ","))
			if err != nil {
				return baggage.Member{}, true
			}
			return m, len(s) == 0
		case reflect.Int:
			s := fieldValue.Interface().([]int)
			sStr := make([]string, len(s))
			for i, v := range s {
				sStr[i] = strconv.Itoa(int(v))
			}
			m, err := baggage.NewMemberRaw(memberKey, strings.Join(sStr, ","))
			if err != nil {
				return baggage.Member{}, true
			}
			return m, len(s) == 0
		case reflect.Int64:
			s := fieldValue.Interface().([]int64)
			sStr := make([]string, len(s))
			for i, v := range s {
				sStr[i] = strconv.FormatInt(v, 10)
			}
			m, err := baggage.NewMemberRaw(memberKey, strings.Join(sStr, ","))
			if err != nil {
				return baggage.Member{}, true
			}
			return m, len(s) == 0
		case reflect.Float64:
			s := fieldValue.Interface().([]float64)
			sStr := make([]string, len(s))
			for i, v := range s {
				sStr[i] = strconv.FormatFloat(v, 'f', -1, 64)
			}
			m, err := baggage.NewMemberRaw(memberKey, strings.Join(sStr, ","))
			if err != nil {
				return baggage.Member{}, true
			}
			return m, len(s) == 0
		case reflect.Bool:
			s := fieldValue.Interface().([]bool)
			sStr := make([]string, len(s))
			for i, v := range s {
				sStr[i] = strconv.FormatBool(v)
			}
			m, err := baggage.NewMemberRaw(memberKey, strings.Join(sStr, ","))
			if err != nil {
				return baggage.Member{}, true
			}
			return m, len(s) == 0
		}
	}

	return baggage.Member{}, true
}
