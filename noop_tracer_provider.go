package motel

import (
	"context"

	"go.opentelemetry.io/otel/trace/noop"
)

type NoopTracerProvider struct {
	noop.TracerProvider
}

func (NoopTracerProvider) Shutdown(_ context.Context) error {
	return nil
}

func NewNoopTracerProvider() TracerProviderWithShutdown {
	return NoopTracerProvider{
		TracerProvider: noop.NewTracerProvider(),
	}
}
