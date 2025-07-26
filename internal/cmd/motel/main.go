package main

import (
	"context"
	"errors"

	"go.bbkane.com/motel"
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/wargcore"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"go.opentelemetry.io/otel/trace"
)

//nolint:gochecknoglobals  // this is a global tracer for the package
var tracer = otel.Tracer("go.bbkane.com/motel/internal/cmd/motel")
var version string

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}

func rolldice(ctx context.Context) {
	_, span := tracer.Start(
		ctx,
		"rolldice",
		trace.WithAttributes(attribute.String("hello", "world")),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("key", "value"),
		attribute.Int("intkey", 32),
	)
	span.SetStatus(codes.Error, "oopsie")
	err := errors.New("the oopsie error")
	span.RecordError(err)
}

func run(cmdCtx wargcore.Context) error {
	ctx := context.Background()

	tracerProvider, err := motel.NewTracerProviderFromEnv(ctx, motel.NewTracerProviderFromEnvArgs{
		AppName: cmdCtx.App.Name,
		Version: cmdCtx.App.Version,
	})
	panicOn(err)

	// set globally
	otel.SetTracerProvider(tracerProvider)

	rolldice(ctx)

	err = tracerProvider.Shutdown(ctx)
	panicOn(err)
	return nil
}

func buildApp() wargcore.App {
	app := warg.New(
		"motel",
		version,
		section.New(
			"Example Go CLI",
			section.NewCommand(
				"run",
				"Send some traces",
				run,
			),
			section.CommandMap(warg.VersionCommandMap()),
		),
		warg.GlobalFlagMap(warg.ColorFlagMap()),
	)
	return app
}

func main() {
	app := buildApp()
	app.MustRun()
}
