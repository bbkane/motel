package motel

import (
	"context"
	"io"
	"os"

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

func NewFileExporter(file io.WriteCloser, opts ...stdouttrace.Option) (*FileExporter, error) {
	opts = append(opts, stdouttrace.WithWriter(file))
	exp, err := stdouttrace.New(
		opts...,
	)
	if err != nil {
		return nil, err
	}
	return &FileExporter{
		file:     file,
		Exporter: exp,
	}, nil

}

func NewFileExporterFromEnv() (*FileExporter, error) {
	filePath := os.Getenv("MOTEL_TRACES_FILE_EXPORTER_FILE_PATH")
	if filePath == "" {
		return nil, errors.New("MOTEL_TRACES_FILE_EXPORTER_FILE_PATH environment variable is not set")
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return NewFileExporter(file)
}
