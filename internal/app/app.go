package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/GoSeoTaxi/s3_karma8_v1/internal/configs"
	"github.com/GoSeoTaxi/s3_karma8_v1/internal/handlers"
)

const maxBytes = (1024 * 1024 * 1024 * 10) + 1024 // 10 ГБ + немного на заголовок  максимального размера загружаемого файла

func Run(
	lifecycle fx.Lifecycle,
	appConfig *configs.AppConfig,
	logger *zap.Logger,
	uploadHandler *handlers.UploadHandler,
	downloadHandler *handlers.DownloadHandler,
	partsHandler *handlers.PartsHandler,
) {
	r := chi.NewRouter()

	r.With(maxBytesMiddleware(maxBytes)).Post("/"+appConfig.StorageConfig.MinioBucketName+"/upload", uploadHandler.HandleUpload)
	r.Get("/"+appConfig.StorageConfig.MinioBucketName+"/{fileName}/download", downloadHandler.HandleDownload)

	// Работа с частями файла
	r.Get("/"+appConfig.StorageConfig.MinioBucketName+"/{fileName}/parts", partsHandler.ListFileParts)
	r.Get("/"+appConfig.StorageConfig.MinioBucketName+"/{fileName}/parts/{partNumber}", partsHandler.DownloadPart)

	server := &http.Server{
		Addr:    appConfig.ServerAddress,
		Handler: r,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Запуск HTTP-сервера", zap.String("адрес", appConfig.ServerAddress))
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("Ошибка при запуске сервера", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Остановка HTTP-сервера")
			return server.Shutdown(ctx)
		},
	})
}
