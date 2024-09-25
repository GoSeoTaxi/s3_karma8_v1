package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

func (fs *FileService) WaitForMinio(ctx context.Context, retries int, delay time.Duration) error {
	for i := 0; i < retries; i++ {
		allAvailable := true
		for idx, client := range fs.StorageRepo.Clients {
			_, err := client.ListBuckets(ctx)
			if err != nil {
				fs.Logger.Warn("MinIO не доступен, повторная попытка...", zap.Int("retry", i+1), zap.Error(err))
				allAvailable = false
				break
			}
			fs.Logger.Info("MinIO доступен", zap.Int("clientIndex", idx))
		}
		if allAvailable {
			fs.Logger.Info("Все MinIO сервисы доступны")
			return nil
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("MinIO сервисы недоступны после %d попыток", retries)
}
