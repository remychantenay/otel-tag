package oteltag_test

import (
	"context"
	"slices"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	oteltag "github.com/remychantenay/otel-tag"
)

func TestSpanAttributes(t *testing.T) {
	const testOperationName = "span"

	setupTracer := func() (*tracetest.SpanRecorder, trace.Tracer) {
		spanRecorder := tracetest.NewSpanRecorder()
		traceProvider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
		tracer := traceProvider.Tracer("test-tracer")

		return spanRecorder, tracer
	}

	t.Run("when non-zero values - should add all attributes to span", func(t *testing.T) {
		const (
			wantSpanCount          = 1
			expectedAttributeCount = 10
		)

		spanRecorder, tracer := setupTracer()

		want := map[attribute.Key]attribute.Value{
			"val_str":           attribute.StringValue("a_string"),
			"val_int":           attribute.IntValue(42),
			"val_int64":         attribute.Int64Value(42000000000),
			"val_float64":       attribute.Float64Value(99.718281828),
			"val_bool":          attribute.BoolValue(true),
			"val_str_slice":     attribute.StringSliceValue([]string{"a_string_1", "a_string_2", "a_string_3"}),
			"val_int_slice":     attribute.IntSliceValue([]int{1, 2, 3}),
			"val_int64_slice":   attribute.Int64SliceValue([]int64{100000, 200000, 300000}),
			"val_float64_slice": attribute.Float64SliceValue([]float64{1.1, 2.2, 3.3}),
			"val_bool_slice":    attribute.BoolSliceValue([]bool{true, false, true, false}),
		}

		m := testModel{
			ValStr:          "a_string",
			ValInt:          42,
			ValInt64:        42000000000,
			ValFloat64:      99.718281828,
			ValBool:         true,
			ValStrSlice:     []string{"a_string_1", "a_string_2", "a_string_3"},
			ValIntSlice:     []int{1, 2, 3},
			ValInt64Slice:   []int64{100000, 200000, 300000},
			ValFloat64Slice: []float64{1.1, 2.2, 3.3},
			ValBoolSlice:    []bool{true, false, true, false},
		}

		func() {
			_, span := tracer.Start(
				context.Background(),
				testOperationName,
				trace.WithAttributes(oteltag.SpanAttributes(m)...),
			)
			defer span.End()
		}()

		spans := spanRecorder.Ended()
		if len(spans) != wantSpanCount {
			t.Errorf("\ngot %d spans\nwant %d", len(spans), wantSpanCount)
		}

		attrCount := len(spans[0].Attributes())
		if attrCount != expectedAttributeCount {
			t.Errorf("\ngot %d attributes\nwant %d", attrCount, expectedAttributeCount)
		}

		for k, v := range want {
			if !slices.Contains(spans[0].Attributes(), attribute.KeyValue{Key: k, Value: v}) {
				t.Errorf("\nmissing '%v' attribute with value %q", k, v.AsString())
			}
		}
	})

	t.Run("when zero values and omitted - should not add attributes to span", func(t *testing.T) {
		const (
			wantSpanCount          = 1
			expectedAttributeCount = 1
		)

		spanRecorder, tracer := setupTracer()

		want := map[attribute.Key]attribute.Value{
			"val_bool": attribute.BoolValue(true),
		}

		m := testModel{ValBool: true}

		func() {
			_, span := tracer.Start(
				context.Background(),
				testOperationName,
				trace.WithAttributes(oteltag.SpanAttributes(m)...),
			)
			defer span.End()
		}()

		spans := spanRecorder.Ended()
		if len(spans) != wantSpanCount {
			t.Errorf("\ngot %d spans\nwant %d", len(spans), wantSpanCount)
		}

		attrCount := len(spans[0].Attributes())
		if attrCount != expectedAttributeCount {
			t.Errorf("\ngot %d attributes\nwant %d", attrCount, expectedAttributeCount)
		}

		for k, v := range want {
			if !slices.Contains(spans[0].Attributes(), attribute.KeyValue{Key: k, Value: v}) {
				t.Errorf("\nmissing '%v' attribute with value %q", k, v.AsString())
			}
		}
	})

	t.Run("when zero values and not omitted - should add attributes to span", func(t *testing.T) {
		const (
			wantSpanCount          = 1
			expectedAttributeCount = 1
		)

		spanRecorder, tracer := setupTracer()

		want := map[attribute.Key]attribute.Value{
			"val_str_not_omitted": attribute.StringValue(""),
		}

		m := struct {
			ValStrNotOmitted string `otel:"val_str_not_omitted"`
		}{
			ValStrNotOmitted: "",
		}

		func() {
			_, span := tracer.Start(
				context.Background(),
				testOperationName,
				trace.WithAttributes(oteltag.SpanAttributes(m)...),
			)
			defer span.End()
		}()

		spans := spanRecorder.Ended()
		if len(spans) != wantSpanCount {
			t.Errorf("\ngot %d spans\nwant %d", len(spans), wantSpanCount)
		}

		attrCount := len(spans[0].Attributes())
		if attrCount != expectedAttributeCount {
			t.Errorf("\ngot %d attributes\nwant %d", attrCount, expectedAttributeCount)
		}

		for k, v := range want {
			if !slices.Contains(spans[0].Attributes(), attribute.KeyValue{Key: k, Value: v}) {
				t.Errorf("\nmissing '%v' attribute with value %q", k, v.AsString())
			}
		}
	})

	t.Run("when fields are not tagged - should not add attributes to span", func(t *testing.T) {
		const (
			wantSpanCount          = 1
			expectedAttributeCount = 0
		)

		spanRecorder, tracer := setupTracer()

		m := struct {
			ValStrIgnored string
		}{
			ValStrIgnored: "a_string",
		}

		func() {
			_, span := tracer.Start(
				context.Background(),
				testOperationName,
				trace.WithAttributes(oteltag.SpanAttributes(m)...),
			)
			defer span.End()
		}()

		spans := spanRecorder.Ended()
		if len(spans) != wantSpanCount {
			t.Errorf("\ngot %d spans\nwant %d", len(spans), wantSpanCount)
		}

		attrCount := len(spans[0].Attributes())
		if attrCount != expectedAttributeCount {
			t.Errorf("\ngot %d attributes\nwant %d", attrCount, expectedAttributeCount)
		}
	})

	t.Run("when structs in structs - should add all attributes to span", func(t *testing.T) {
		const (
			wantSpanCount          = 1
			expectedAttributeCount = 3
		)

		want := map[attribute.Key]attribute.Value{
			"struct1.val_str_1": attribute.StringValue("a_string"),
			"struct2.val_str_2": attribute.StringValue("b_string"),
			"struct3.val_str_3": attribute.StringValue("c_string"),
		}

		spanRecorder, tracer := setupTracer()

		type struct1 struct {
			ValStr1 string `otel:"struct1.val_str_1"`
		}

		type struct2 struct {
			ValStr2 string `otel:"struct2.val_str_2"`
		}

		type struct3 struct {
			S1      struct1
			S2      *struct2
			ValStr3 string `otel:"struct3.val_str_3"`
		}

		m := struct3{
			S1: struct1{
				ValStr1: "a_string",
			},
			S2: &struct2{
				ValStr2: "b_string",
			},
			ValStr3: "c_string",
		}

		func() {
			_, span := tracer.Start(
				context.Background(),
				testOperationName,
				trace.WithAttributes(oteltag.SpanAttributes(m)...),
			)
			defer span.End()
		}()

		spans := spanRecorder.Ended()
		if len(spans) != wantSpanCount {
			t.Errorf("\ngot %d spans\nwant %d", len(spans), wantSpanCount)
		}

		attrCount := len(spans[0].Attributes())
		if attrCount != expectedAttributeCount {
			t.Errorf("\ngot %d attributes\nwant %d", attrCount, expectedAttributeCount)
		}

		for k, v := range want {
			if !slices.Contains(spans[0].Attributes(), attribute.KeyValue{Key: k, Value: v}) {
				t.Errorf("\nmissing '%v' attribute with value %q", k, v.AsString())
			}
		}
	})
}
