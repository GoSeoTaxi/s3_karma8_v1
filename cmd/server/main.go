package main

import (
	"context"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/GoSeoTaxi/s3_karma8_v1/internal/app"
	"github.com/GoSeoTaxi/s3_karma8_v1/internal/configs"
	"github.com/GoSeoTaxi/s3_karma8_v1/internal/handlers"
	"github.com/GoSeoTaxi/s3_karma8_v1/internal/repositories"
	"github.com/GoSeoTaxi/s3_karma8_v1/internal/services"
)

const tryWaitForMinio = 2
const secWaitForMinio = 1 * time.Second

func main() {

	fx.New(
		fx.Provide(
			configs.LoadConfig,
			zap.NewProduction,
			provideStorageConfig,
			providePartSizeThreshold,
			repositories.NewStorageRepository,
			services.NewFileService,
			handlers.NewUploadHandler,
			handlers.NewDownloadHandler,
			handlers.NewPartsHandler,
		),
		fx.Invoke(registerHooks),
		fx.Invoke(app.Run),
	).Run()
}

func provideStorageConfig(appConfig *configs.AppConfig) *repositories.StorageConfig {
	return &appConfig.StorageConfig
}

func providePartSizeThreshold(appConfig *configs.AppConfig) int64 {
	return appConfig.PartSizeThreshold
}

func registerHooks(lc fx.Lifecycle, fs *services.FileService, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Запуск приложения. Ожидание доступности MinIO сервисов...")
			err := fs.WaitForMinio(ctx, tryWaitForMinio, secWaitForMinio)
			if err != nil {
				logger.Fatal("MinIO сервисы недоступны", zap.Error(err))
				return err
			}
			logger.Info("MinIO сервисы успешно доступны")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Остановка приложения...")
			return nil
		},
	})
}
