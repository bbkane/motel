# motel

A library to help me set up OTEL trace exports with environment variables. Similar to [autoexport](https://pkg.go.dev/go.opentelemetry.io/contrib/exporters/autoexport), but tailored to what I care about.

## Project status (2025-06-26)

`motel` is brand new and super unstable. I'll probably change a lot of things.

## Use

Initialize in code:

```go
	tracerProvider, err := motel.NewTracerProviderFromEnv(ctx, motel.NewTracerProviderFromEnvArgs{
		AppName: "myapp",
		Version: "v0.0.1"
	})
```

Set environment variables:

```bash
MOTEL_SPAN_PROCESSOR=batch|sync
MOTEL_TRACES_EXPORTER=file|none|otlpgrpc|otlphttp|stderr|stdout
MOTEL_TRACES_FILE_EXPORTER_FILE_PATH=tmp.jsonl  # required if MOTEL_TRACES_EXPORTER=file
```

See https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/ for other environment variables.

## Examples

### `file` + local [logdy](https://logdy.dev/)

Create a file and start following it with logy using  [this config](https://github.com/bbkane/dotfiles/blob/master/logdy/dot-config/logdy.json):

```bash
$ touch tmp.jsonl && logdy follow ./tmp.jsonl --config ~/Git-GH/dotfiles/logdy/dot-config/logdy.json
```

Open the [local logs page](http://127.0.0.1:8080/) and wait for logs.

In a new terminal, run the motel example:

```bash
$ MOTEL_TRACES_EXPORTER=file MOTEL_TRACES_FILE_EXPORTER_FILE_PATH=tmp.jsonl  go run ./internal/cmd/motel run
rolling dice
rolling dice
^Ccontext done
calling ticker stop
calling cancel
```

Head back to the browser to view traces coming in (might take a minute if using `MOTEL_SPAN_PROCESSOR=batch`)

![logdy-demo](./README.assets/logdy-demo.png)

### `stdout` (easiest and least useful)

```bash
$ MOTEL_TRACES_EXPORTER=stdout go run ./internal/cmd/motel run
rolling dice
^Ccontext done
calling ticker stop
{"Name":"rolldice","SpanContext":{"TraceID":"990f2f9167e20ecc96d1e777061520fa","SpanID":"6be766f137dba10e","TraceFlags":"01","TraceState":"","Remote":false},"Parent":{"TraceID":"00000000000000000000000000000000","SpanID":"0000000000000000","TraceFlags":"00","TraceState":"","Remote":false},"SpanKind":1,"StartTime":"2025-08-03T16:23:12.24662-07:00","EndTime":"2025-08-03T16:23:12.246738958-07:00","Attributes":[{"Key":"hello","Value":{"Type":"STRING","Value":"world"}},{"Key":"key","Value":{"Type":"STRING","Value":"value"}},{"Key":"intkey","Value":{"Type":"INT64","Value":32}}],"Events":[{"Name":"exception","Attributes":[{"Key":"exception.type","Value":{"Type":"STRING","Value":"*errors.errorString"}},{"Key":"exception.message","Value":{"Type":"STRING","Value":"the oopsie error"}}],"DroppedAttributeCount":0,"Time":"2025-08-03T16:23:12.246735-07:00"}],"Links":null,"Status":{"Code":"Error","Description":"oopsie"},"DroppedAttributes":0,"DroppedEvents":0,"DroppedLinks":0,"ChildSpanCount":0,"Resource":[{"Key":"service.name","Value":{"Type":"STRING","Value":"motel"}},{"Key":"service.version","Value":{"Type":"STRING","Value":"(devel)"}}],"InstrumentationScope":{"Name":"go.bbkane.com/motel/internal/cmd/motel","Version":"","SchemaURL":"","Attributes":null},"InstrumentationLibrary":{"Name":"go.bbkane.com/motel/internal/cmd/motel","Version":"","SchemaURL":"","Attributes":null}}
calling cancel
```

### `otlphttp` + locally hosted [OpenObserve](https://openobserve.ai)

While I'm using locally hosted OpenObserve since it's easy to run locally, this should also work with OTEL vendors in the cloud: [Grafana](https://grafana.com/products/cloud/),  [Honeycomb](https://www.honeycomb.io/), [SigNoz](https://signoz.io/) and others I'm sure I'm forgetting right now

Download OpenObserve (I use the OSS version) from the [downloads page](https://openobserve.ai/downloads/).

Run it with some environemental variables:

```bash
ZO_ROOT_USER_EMAIL=root@example.com \
ZO_ROOT_USER_PASSWORD=pass \
./openobserve
```

Open the [Data sources > Traces page](http://127.0.0.1:5080/web/ingestion/custom/traces/opentelemetry?org_identifier=default) and turn the OTLP HTTP sources into environment variables (note: the headers seems to change each time), then run the example in a new terminal window:

```bash
$ MOTEL_TRACES_EXPORTER=otlphttp \
OTEL_EXPORTER_OTLP_ENDPOINT='http://127.0.0.1:5080/api/default' \
OTEL_EXPORTER_OTLP_HEADERS='Authorization=Basic cm9vdEBleGFtcGxlLmNvbTptODM4MDVDc2NZU1BZbTJP' \
go run ./internal/cmd/motel run
rolling dice
rolling dice
rolling dice
^Ccontext done
calling ticker stop
calling cancel
```

View traces from the [Traces](http://127.0.0.1:5080/web/traces?org_identifier=default) page:

![openobserve-demo](./README.assets/openobserve-demo.png)

## Notes

See [Go Project Notes](https://www.bbkane.com/blog/go-project-notes/) for notes on development tooling.
