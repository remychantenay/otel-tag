package oteltag_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/baggage"

	oteltag "github.com/remychantenay/otel-tag"
)

func TestBaggageMembers(t *testing.T) {
	t.Run("when non-zero values - should add all members to baggage", func(t *testing.T) {
		const wantMemberCount = 10

		want := map[string]string{
			"val_str":           "a_string",
			"val_int":           "42",
			"val_int64":         "42000000000",
			"val_float64":       "99.718281828",
			"val_bool":          "true",
			"val_str_slice":     "a_string_1,a_string_2,a_string_3",
			"val_int_slice":     "1,2,3",
			"val_int64_slice":   "100000,200000,300000",
			"val_float64_slice": "1.1,2.2,3.3",
			"val_bool_slice":    "true,false,true,false",
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

		members := oteltag.BaggageMembers(m)
		bag, _ := baggage.New(members...)
		ctx := baggage.ContextWithBaggage(context.Background(), bag)
		bag = baggage.FromContext(ctx)

		memberCount := len(bag.Members())
		if memberCount != wantMemberCount {
			t.Errorf("\ngot %d members\nwant %d", memberCount, wantMemberCount)
		}

		for k, v := range want {
			member := bag.Member(k)
			if member.Value() != v {
				t.Errorf("\ngot %q for member %q\nwant %q", member.Value(), k, v)
			}
		}
	})

	t.Run("when zero values and omitted - should not add members to baggage", func(t *testing.T) {
		const wantMemberCount = 1

		want := map[string]string{
			"val_bool": "true",
		}

		m := testModel{ValBool: true}

		members := oteltag.BaggageMembers(m)
		bag, _ := baggage.New(members...)
		ctx := baggage.ContextWithBaggage(context.Background(), bag)
		bag = baggage.FromContext(ctx)

		memberCount := len(bag.Members())
		if memberCount != wantMemberCount {
			t.Errorf("\ngot %d members\nwant %d", memberCount, wantMemberCount)
		}

		for k, v := range want {
			member := bag.Member(k)
			if member.Value() != v {
				t.Errorf("\ngot %q for member %q\nwant %q", member.Value(), k, v)
			}
		}
	})

	t.Run("when zero values and not omitted - should add members to baggage", func(t *testing.T) {
		const wantMemberCount = 1

		want := map[string]string{
			"val_str_not_omitted": "",
		}

		m := struct {
			ValStrNotOmitted string `otel:"val_str_not_omitted"`
		}{
			ValStrNotOmitted: "",
		}

		members := oteltag.BaggageMembers(m)
		bag, _ := baggage.New(members...)
		ctx := baggage.ContextWithBaggage(context.Background(), bag)
		bag = baggage.FromContext(ctx)

		memberCount := len(bag.Members())
		if memberCount != wantMemberCount {
			t.Errorf("\ngot %d members\nwant %d", memberCount, wantMemberCount)
		}

		for k, v := range want {
			member := bag.Member(k)
			if member.Value() != v {
				t.Errorf("\ngot %q for member %q\nwant %q", member.Value(), k, v)
			}
		}
	})

	t.Run("when fields are not tagged - should not add members to baggage", func(t *testing.T) {
		const wantMemberCount = 0

		m := struct {
			ValStrIgnored string
		}{
			ValStrIgnored: "a_string",
		}

		members := oteltag.BaggageMembers(m)
		bag, _ := baggage.New(members...)
		ctx := baggage.ContextWithBaggage(context.Background(), bag)
		bag = baggage.FromContext(ctx)

		memberCount := len(bag.Members())
		if memberCount != wantMemberCount {
			t.Errorf("\ngot %d members\nwant %d", memberCount, wantMemberCount)
		}
	})
}
