package oteltag_test

import (
	"context"
	"testing"

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
			expectedAttributeCount = 12
		)

		spanRecorder, tracer := setupTracer()

		m := model{
			ValStr:          "a_string",
			ValInt:          42,
			ValInt64:        42000000000,
			ValFloat32:      3.14,
			ValFloat64:      99.718281828,
			ValBool:         true,
			ValStrSlice:     []string{"a_string_1", "a_string_2", "a_string_3"},
			ValIntSlice:     []int{1, 2, 3},
			ValInt64Slice:   []int64{100000, 200000, 300000},
			ValFloat32Slice: []float32{1.1, 2.2, 3.3},
			ValFloat64Slice: []float64{4.4, 5.5, 6.6},
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

		// for i := range spans[0].Attributes() {
		// 	t.Log(spans[0].Attributes()[i].Value)
		// }
	})

	t.Run("when zero values and omitted - should not add attributes to span", func(t *testing.T) {
		const (
			wantSpanCount          = 1
			expectedAttributeCount = 1
		)

		spanRecorder, tracer := setupTracer()

		m := model{ValBool: true}

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

	t.Run("when zero values and not omitted - should add attributes to span", func(t *testing.T) {
		const (
			wantSpanCount          = 1
			expectedAttributeCount = 1
		)

		spanRecorder, tracer := setupTracer()

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

		for i := range spans[0].Attributes() {
			t.Log(spans[0].Attributes()[i].Key)
		}
	})
}
