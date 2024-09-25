package repositories

import (
	"context"
	"sort"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"

	"github.com/GoSeoTaxi/s3_karma8_v1/internal/models"
)

func (s *StorageRepository) ListFileParts(ctx context.Context, fileName string) ([]models.Part, error) {
	var allParts []models.Part

	prefix := fileName + "_part_"

	var client MinioClient
	client = s.getRandomClient()

	objectCh := client.ListObjects(ctx, s.BucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: false,
	})

	for object := range objectCh {
		if object.Err != nil {
			s.Logger.Error("Ошибка при получении списка объектов", zap.Error(object.Err))
			return nil, object.Err
		}

		partNumber := extractPartIndex(object.Key, fileName)
		if partNumber == -1 {
			s.Logger.Warn("Не удалось извлечь номер части из имени объекта", zap.String("objectName", object.Key))
			continue
		}

		part := models.Part{
			FileName:   fileName,
			PartNumber: partNumber,
			Size:       object.Size,
		}

		allParts = append(allParts, part)
	}

	sort.Slice(allParts, func(i, j int) bool {
		return allParts[i].PartNumber < allParts[j].PartNumber
	})

	s.Logger.Info("Получен список частей файла", zap.String("fileName", fileName), zap.Int("totalParts", len(allParts)))
	return allParts, nil
}
