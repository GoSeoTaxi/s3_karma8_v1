package services

import (
	"go.uber.org/zap"

	"github.com/GoSeoTaxi/s3_karma8_v1/internal/repositories"
)

type FileService struct {
	StorageRepo       *repositories.StorageRepository
	Logger            *zap.Logger
	PartSizeThreshold int64
}

func NewFileService(storageRepo *repositories.StorageRepository, logger *zap.Logger, partSizeThreshold int64) *FileService {
	return &FileService{
		StorageRepo:       storageRepo,
		Logger:            logger,
		PartSizeThreshold: partSizeThreshold,
	}
}
