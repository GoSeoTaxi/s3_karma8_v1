package services

import (
	"context"

	"github.com/GoSeoTaxi/s3_karma8_v1/internal/models"
)

func (fs *FileService) ListFileParts(ctx context.Context, fileName string) ([]models.Part, error) {
	return fs.StorageRepo.ListFileParts(ctx, fileName)
}
