package motel

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type TracerProviderWithShutdown interface {
	trace.TracerProvider
	Shutdown(ctx context.Context) error
}

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

type NewTracerProviderFromEnvArgs struct {
	AppName string
	Version string
}

func NewTracerProviderFromEnv(ctx context.Context, args NewTracerProviderFromEnvArgs) (TracerProviderWithShutdown, error) {

	tracerResource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(args.AppName),
		semconv.ServiceVersion(args.Version),
	)

	// NOTE: https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#exporter-selection
	// defines OTEL_TRACES_EXPORTER, but grepping https://github.com/open-telemetry/opentelemetry-go/tree/main doesn't reveal that...
	// so we're doing our own
	var exporter sdktrace.SpanExporter
	var exporterErr error
	exporterType := os.Getenv("MOTEL_TRACES_EXPORTER")
	switch exporterType {
	// case "file":
	// 	// TODO - probably wrap the stdouttrace one but open a file and shut it down...
	case "none":
		// TODO: this needs to be done with the trace provider, not the exporter...
		// https://github.com/open-telemetry/opentelemetry-go/blob/main/trace/noop_test.go
		tp := NewNoopTracerProvider()
		return tp, nil
	case "otlpgrpc":
		// https://github.com/open-telemetry/opentelemetry-go/blob/main/exporters/otlp/otlptrace/otlptracegrpc/example_test.go
		exporter, exporterErr = otlptracegrpc.New(ctx)
	case "otlphttp":
		// https://github.com/open-telemetry/opentelemetry-go/blob/main/exporters/otlp/otlptrace/otlptracehttp/example_test.go
		exporter, exporterErr = otlptracehttp.New(ctx)
	// case "prettyprint":
	// 	// TODO - write this one like logos
	case "stdout", "": // "" if not set
		exporter, exporterErr = stdouttrace.New(
			stdouttrace.WithWriter(os.Stdout),
		)
	default:
		return nil, fmt.Errorf("unreachable: unknown exporter %s", exporterType)
	}
	if exporterErr != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", exporterErr)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(tracerResource),
	)

	return tp, nil
}
