package services

import (
	"context"
	"fmt"
	"io"

	"go.uber.org/zap"
)

func (fs *FileService) UploadFile(ctx context.Context, filename string, data io.Reader, size int64) error {
	if size == -1 {
		return fs.uploadFileWithoutSize(ctx, filename, data)
	}

	var numParts int64
	if size > fs.PartSizeThreshold {
		numParts = (size + fs.PartSizeThreshold - 1) / fs.PartSizeThreshold
	} else {
		numParts = 1
	}

	for i := int64(0); i < numParts; i++ {
		currentPartSize := fs.PartSizeThreshold
		if i == numParts-1 {
			currentPartSize = size - fs.PartSizeThreshold*i
		}

		limitedReader := io.LimitReader(data, currentPartSize)
		objectName := fmt.Sprintf("%s_part_%d", filename, i)

		err := fs.StorageRepo.UploadPart(ctx, objectName, limitedReader, currentPartSize)
		if err != nil {
			fs.Logger.Error("Не удалось загрузить часть", zap.Int64("part", i), zap.Error(err))
			return fmt.Errorf("не удалось загрузить часть %d: %w", i, err)
		}

		fs.Logger.Info("Часть загружена", zap.String("objectName", objectName))
	}

	fs.Logger.Info("Файл успешно загружен", zap.String("filename", filename))
	return nil
}
