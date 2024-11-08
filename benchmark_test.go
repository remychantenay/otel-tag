package oteltag_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	oteltag "github.com/remychantenay/otel-tag"
)

var (
	traceExporter *tracetest.InMemoryExporter
	tracer        trace.Tracer
)

func init() {
	traceExporter, tracer = setUpBenchmarkTracer()
}

func setUpBenchmarkTracer() (*tracetest.InMemoryExporter, trace.Tracer) {
	exporter := tracetest.NewInMemoryExporter()
	traceProvider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	tracer := traceProvider.Tracer("benchmark-tracer")

	return exporter, tracer
}

func BenchmarkSpanAttributes_With(b *testing.B) {
	traceExporter.Reset()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
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

			_, span := tracer.Start(
				context.Background(), "operation-name",
				trace.WithAttributes(oteltag.SpanAttributes(m)...),
			)
			defer span.End()
		}
	})
}

func BenchmarkSpanAttributes_Without(b *testing.B) {
	traceExporter.Reset()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, span := tracer.Start(
				context.Background(), "operation-name",
				trace.WithAttributes(
					attribute.String("val_str", "a_string"),
					attribute.Int("val_int", 42),
					attribute.Int64("val_int64", 42000000000),
					attribute.Float64("val_float32", 3.14),
					attribute.Float64("val_float64", 99.718281828),
					attribute.Bool("val_bool", true),
					attribute.StringSlice("val_str_slice", []string{"a_string_1", "a_string_2", "a_string_3"}),
					attribute.IntSlice("val_int_slice", []int{1, 2, 3}),
					attribute.Int64Slice("val_int64_slice", []int64{100000, 200000, 300000}),
					attribute.Float64Slice("val_float32_slice", []float64{1.1, 2.2, 3.3}),
					attribute.Float64Slice("val_float64_slice", []float64{4.4, 5.5, 6.6}),
					attribute.BoolSlice("val_bool_slice", []bool{true, false, true, false}),
				),
			)
			defer span.End()
		}
	})
}
