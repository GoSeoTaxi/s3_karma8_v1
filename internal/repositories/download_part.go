package repositories

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

func (s *StorageRepository) DownloadPart(ctx context.Context, objectName string, writer io.Writer) error {
	var client MinioClient
	client = s.getRandomClient()

	obj, err := client.GetObject(ctx, s.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		s.Logger.Error("Не удалось получить объект", zap.Error(err))
		return err
	}
	defer func() { _ = obj.Close() }()

	_, err = io.Copy(writer, obj)
	if err != nil {
		s.Logger.Error("Не удалось скопировать объект", zap.Error(err))
		return err
	}

	s.Logger.Info("Часть скачана", zap.String("objectName", objectName))
	return nil
}
