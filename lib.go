package motel

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type TracerProviderWithShutdown interface {
	trace.TracerProvider
	Shutdown(ctx context.Context) error
}

type NewTracerProviderFromEnvArgs struct {
	AppName string
	Version string
}

// NewTracerProviderFromEnv creates a new TracerProvider based on environment variables.
//
// Processor environment variables (defaults to "batch"):
//   - MOTEL_SPAN_PROCESSOR=batch|sync
//
// Exporter environment variables (defaults to "none"):
//   - MOTEL_TRACES_EXPORTER=file|none|otlpgrpc|otlphttp|stdout
//   - MOTEL_TRACES_FILE_EXPORTER_FILE_PATH (required if MOTEL_TRACES_EXPORTER=file)
//
// See https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/ for other environment variables.
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
	case "file":
		exporter, exporterErr = NewFileExporterFromEnv()
	case "none", "": // "" means not set, so we default to noop
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
	case "stdout":
		exporter, exporterErr = stdouttrace.New(
			stdouttrace.WithWriter(os.Stdout),
		)
	default:
		return nil, fmt.Errorf("unreachable: unknown exporter %s", exporterType)
	}
	if exporterErr != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", exporterErr)
	}

	var spanProcessor sdktrace.TracerProviderOption
	switch os.Getenv("MOTEL_SPAN_PROCESSOR") {
	case "sync":
		spanProcessor = sdktrace.WithSyncer(exporter)
	case "batch", "":
		spanProcessor = sdktrace.WithBatcher(exporter)
	default:
		return nil, fmt.Errorf("unreachable: unknown span processor %s", os.Getenv("MOTEL_SPAN_PROCESSOR"))
	}

	tp := sdktrace.NewTracerProvider(
		spanProcessor,
		sdktrace.WithResource(tracerResource),
	)

	return tp, nil
}
