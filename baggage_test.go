package oteltag_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/baggage"

	oteltag "github.com/remychantenay/otel-tag"
)

func TestBaggageMembers(t *testing.T) {
	t.Run("when non-zero values - should add all members to baggage", func(t *testing.T) {
		const wantMemberCount = 12

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

		members := oteltag.BaggageMembers(m)
		bag, _ := baggage.New(members...)
		ctx := baggage.ContextWithBaggage(context.Background(), bag)
		bag = baggage.FromContext(ctx)

		memberCount := len(bag.Members())
		if memberCount != wantMemberCount {
			t.Errorf("\ngot %d members\nwant %d", memberCount, wantMemberCount)
		}

		for i := range bag.Members() {
			t.Log(bag.Members()[i].Key())
		}
	})

	t.Run("when zero values and omitted - should not add members to baggage", func(t *testing.T) {
		const wantMemberCount = 1

		m := model{ValBool: true}

		members := oteltag.BaggageMembers(m)
		bag, _ := baggage.New(members...)
		ctx := baggage.ContextWithBaggage(context.Background(), bag)
		bag = baggage.FromContext(ctx)

		memberCount := len(bag.Members())
		if memberCount != wantMemberCount {
			t.Errorf("\ngot %d members\nwant %d", memberCount, wantMemberCount)
		}
	})

	t.Run("when zero values and not omitted - should add members to baggage", func(t *testing.T) {
		const wantMemberCount = 1

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

		for i := range bag.Members() {
			t.Log(bag.Members()[i].Key())
		}
	})
}
