# otel-tag
Simple package that extracts OpenTelemetry [span attributes](https://opentelemetry.io/docs/demo/telemetry-features/manual-span-attributes/) and [baggage members](https://opentelemetry.io/docs/concepts/signals/baggage/) from a struct based on tags.

> [!WARNING]
> Please bear in mind that this package does not intend to cater to all use cases and please everyone. The code could also do with some refactoring here and there, but it's doing the job.

## Shortcomings
- Does not support basic type pointers (only struct pointers).

## Usage
```go
package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"

	oteltag "github.com/remychantenay/otel-tag"
)

type User struct {
	ID        string `otel:"app.user.id"`
	Username  string `otel:"app.user.username"`
	IsPremium bool   `otel:"app.user.premium"`

	UserDetails UserDetails
}

type UserDetails struct {
	Website string `otel:"app.user.website,omitempty"`
	Bio     string // Not tagged, will be ignored.
}

func main() {
	tracer := otel.Tracer("example")
	
	user := User{
		ID:        "123",
		Username:  "john_doe",
		IsPremium: true,
		UserDetails: UserDetails{
			Website: "https://example.com",
			Bio:     "Software developer",
		},
	}

	// Span attributes.
	_, span := tracer.Start(context.Background(), "someOperation",
		trace.WithAttributes(oteltag.SpanAttributes(user)...),
	)
	defer span.End()

	// Baggage Members.
	members := oteltag.BaggageMembers(user)
	bag, _ := baggage.New(members...)
	ctx := baggage.ContextWithBaggage(context.Background(), bag)
}
```

## License
Apache License Version 2.0
