package oteltag_test

type testModel struct {
	ValStr          string    `otel:"val_str,omitempty"`
	ValInt          int       `otel:"val_int,omitempty"`
	ValInt64        int64     `otel:"val_int64,omitempty"`
	ValFloat64      float64   `otel:"val_float64,omitempty"`
	ValBool         bool      `otel:"val_bool,omitempty"`
	ValStrSlice     []string  `otel:"val_str_slice,omitempty"`
	ValIntSlice     []int     `otel:"val_int_slice,omitempty"`
	ValInt64Slice   []int64   `otel:"val_int64_slice,omitempty"`
	ValFloat64Slice []float64 `otel:"val_float64_slice,omitempty"`
	ValBoolSlice    []bool    `otel:"val_bool_slice,omitempty"`
}
