package repositories

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

func (s *StorageRepository) UploadPart(ctx context.Context, objectName string, data io.Reader, size int64) error {
	var client MinioClient
	client = s.getRandomClient()

	_, err := client.PutObject(ctx, s.BucketName, objectName, data, size, minio.PutObjectOptions{})
	if err != nil {
		s.Logger.Error("Не удалось загрузить часть", zap.Error(err))
		return err
	}

	s.Logger.Info("Часть загружена", zap.String("objectName", objectName))
	return nil
}
