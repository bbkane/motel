package motel

import (
	"context"
	"io"

	"errors"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
)

type FileExporter struct {
	file io.WriteCloser
	*stdouttrace.Exporter
}

func (fe *FileExporter) Shutdown(ctx context.Context) error {
	shutdownErr := fe.Exporter.Shutdown(ctx)
	fileCloseErr := fe.file.Close()
	if shutdownErr != nil || fileCloseErr != nil {
		return errors.Join(errors.New("failed to shutdown FileExporter"), shutdownErr, fileCloseErr)
	}
	return nil
}

func NewFileExporter(file io.WriteCloser) (*FileExporter, error) {
	// TODO: be able to set opts so I can set no timestamps and test this thing.
	exp, err := stdouttrace.New(
		stdouttrace.WithWriter(file),
	)
	if err != nil {
		return nil, err
	}
	return &FileExporter{
		file:     file,
		Exporter: exp,
	}, nil

}

// TODO: NewFileExporterFromEnv() (*FileExporter, error) {}
