package backup

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/pkg/errors"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type fileBackuper struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewFileBackuper(ctx context.Context, filename string) (*fileBackuper, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0o664)
	if err != nil {
		return nil, err
	}
	fb := &fileBackuper{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}
	go fb.waitAndClose(ctx)
	return fb, nil
}

func (b *fileBackuper) Backup(metrics []domain.Metric) error {
	// Clean up file
	err := b.file.Truncate(0)
	if err != nil {
		return errors.Wrap(err, "[file backuper] error when truncating file before doing backup")
	}
	// Reset internal file offset
	_, err = b.file.Seek(0, 0)
	if err != nil {
		return errors.Wrap(err, "[file backuper] error when truncating file before doing backup")
	}

	// Write metrics to file
	err = b.encoder.Encode(metrics)
	if err != nil {
		return errors.Wrap(err, "[file backuper] error when backing up metrics")
	}
	return nil
}

func (b *fileBackuper) Restore() ([]domain.Metric, error) {
	// Check if file is empty
	stat, err := b.file.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "[file backuper] error when reading file stat")
	}
	size := stat.Size()
	if size == 0 {
		// File is empty, nothing to restore from
		return nil, nil
	}

	// Try to read JSON from file
	metrics := make([]domain.Metric, 0)
	err = b.decoder.Decode(&metrics)
	if err != nil && !errors.Is(err, io.EOF) { // Ignore io.EOF (empty file)
		return nil, errors.Wrap(err, "[file backuper] error when restoring metrics")
	}
	return metrics, nil
}

func (b *fileBackuper) waitAndClose(ctx context.Context) {
	<-ctx.Done()
	logger.New(ctx).Debugf("[file backuper] context cancelled, closing file")

	err := b.file.Close()
	if err != nil {
		logger.New(ctx).Errorf("[file backuper] error when closing file: %s", err.Error())
	}
}
