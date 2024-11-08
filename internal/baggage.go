package internal

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/baggage"
)

// BaggageMember creates and returns an OpenTelemetry baggageMember for the provided field.
// Also returns a boolean that indicates whether or not the field's value is a zero-value.
func BaggageMember(field reflect.StructField, value reflect.Value, memberKey string) (baggage.Member, bool) {

	switch field.Type.Kind() {
	case reflect.String:
		v := value.String()
		m, _ := baggage.NewMemberRaw(memberKey, v)
		return m, v == ""
	case reflect.Int:
		v := value.Int()
		m, _ := baggage.NewMemberRaw(memberKey, strconv.Itoa(int(v)))
		return m, v == 0
	case reflect.Int64:
		v := value.Int()
		m, _ := baggage.NewMemberRaw(memberKey, strconv.FormatInt(v, 10))
		return m, v == 0
	case reflect.Float32, reflect.Float64:
		v := value.Float()
		m, _ := baggage.NewMemberRaw(memberKey, strconv.FormatFloat(v, 'f', -1, 64))
		return m, v == 0.0
	case reflect.Bool:
		m, _ := baggage.NewMemberRaw(memberKey, strconv.FormatBool(value.Bool()))
		return m, false
	case reflect.Slice:
		switch field.Type.Elem().Kind() {
		case reflect.String:
			if s, ok := value.Interface().([]string); ok {
				m, _ := baggage.NewMemberRaw(memberKey, strings.Join(s, ","))
				return m, len(s) == 0
			}
		case reflect.Int:
			if s, ok := value.Interface().([]int); ok {
				sStr := make([]string, len(s))
				for i, v := range s {
					sStr[i] = strconv.Itoa(int(v))
				}
				m, _ := baggage.NewMemberRaw(memberKey, strings.Join(sStr, ","))
				return m, len(s) == 0
			}
		case reflect.Int64:
			if s, ok := value.Interface().([]int64); ok {
				sStr := make([]string, len(s))
				for i, v := range s {
					sStr[i] = strconv.FormatInt(v, 10)
				}
				m, _ := baggage.NewMemberRaw(memberKey, strings.Join(sStr, ","))
				return m, len(s) == 0
			}
		case reflect.Float32:
			if s, ok := value.Interface().([]float32); ok {
				sStr := make([]string, len(s))
				for i, v := range s {
					sStr[i] = fmt.Sprintf("%f", v)
				}
				m, _ := baggage.NewMemberRaw(memberKey, strings.Join(sStr, ","))
				return m, len(s) == 0
			}
		case reflect.Float64:
			if s, ok := value.Interface().([]float64); ok {
				sStr := make([]string, len(s))
				for i, v := range s {
					sStr[i] = strconv.FormatFloat(v, 'f', -1, 64)
				}
				m, _ := baggage.NewMemberRaw(memberKey, strings.Join(sStr, ","))
				return m, len(s) == 0
			}
		case reflect.Bool:
			if s, ok := value.Interface().([]bool); ok {
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
