package services

import (
	"context"
	"fmt"
	"io"
	"sort"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

func (fs *FileService) DownloadFile(ctx context.Context, filename string, writer io.Writer) error {
	prefix := fmt.Sprintf("%s_part_", filename)
	var objectNames []string

	for _, client := range fs.StorageRepo.Clients {
		objectCh := client.ListObjects(ctx, fs.StorageRepo.BucketName, minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		})

		for object := range objectCh {
			if object.Err != nil {
				fs.Logger.Error("Ошибка при перечислении объектов", zap.Error(object.Err))
				continue
			}
			objectNames = append(objectNames, object.Key)
		}
	}

	objectNames = removeDuplicates(objectNames)

	if len(objectNames) == 0 {
		fs.Logger.Error("Части файла не найдены", zap.String("filename", filename))
		return fmt.Errorf("части файла %s не найдены", filename)
	}

	sort.Slice(objectNames, func(i, j int) bool {
		partI := extractPartIndex(objectNames[i])
		partJ := extractPartIndex(objectNames[j])
		return partI < partJ
	})

	for _, objectName := range objectNames {
		err := fs.StorageRepo.DownloadPart(ctx, objectName, writer)
		if err != nil {
			fs.Logger.Error("Не удалось скачать часть", zap.String("objectName", objectName), zap.Error(err))
			return fmt.Errorf("не удалось скачать часть %s: %w", objectName, err)
		}
		fs.Logger.Info("Часть скачана", zap.String("objectName", objectName))
	}

	fs.Logger.Info("Файл успешно скачан", zap.String("filename", filename))
	return nil
}
