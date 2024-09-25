package repositories

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type StorageConfig struct {
	Endpoints       []string
	AccessKey       string
	SecretKey       string
	UseSSL          bool
	MinioBucketName string
}

type StorageRepository struct {
	Clients    []MinioClient
	Logger     *zap.Logger
	BucketName string
}

func NewStorageRepository(config *StorageConfig, logger *zap.Logger) (*StorageRepository, error) {
	var clients []MinioClient

	for _, endpoint := range config.Endpoints {
		client, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
			Secure: config.UseSSL,
		})
		if err != nil {
			return nil, fmt.Errorf("не удалось создать клиент MinIO для endpoint %s: %w", endpoint, err)
		}
		clients = append(clients, client)
	}

	repo := &StorageRepository{
		Clients:    clients,
		Logger:     logger,
		BucketName: config.MinioBucketName,
	}

	for _, client := range clients {
		err := ensureBucketExists(client, config.MinioBucketName)
		if err != nil {
			return nil, fmt.Errorf("не удалось проверить или создать бакет на MinIO: %w", err)
		}
	}

	return repo, nil
}
