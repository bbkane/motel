package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

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

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer func() {
		fmt.Println("calling cancel")
		cancel()
	}()

	tracerProvider := must(motel.NewTracerProviderFromEnv(ctx, motel.NewTracerProviderFromEnvArgs{
		AppName: cmdCtx.App.Name,
		Version: cmdCtx.App.Version,
	}))

	// set tracer provider globally
	otel.SetTracerProvider(tracerProvider)
	defer func() {
		// shutdown in a new context since the main context will be canceled by the time this runs
		err := tracerProvider.Shutdown(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		fmt.Println("calling ticker stop")
		ticker.Stop()
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("context done")
			return nil
		case <-ticker.C:
			fmt.Println("rolling dice")
			rolldice(ctx)
		}
	}
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
