package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/GoSeoTaxi/s3_karma8_v1/internal/services"
)

type DownloadHandler struct {
	FileService *services.FileService
	Logger      *zap.Logger
}

func NewDownloadHandler(fs *services.FileService, logger *zap.Logger) *DownloadHandler {
	return &DownloadHandler{
		FileService: fs,
		Logger:      logger,
	}
}

func (h *DownloadHandler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	if fileName == "" {
		http.Error(w, "fileName is required", http.StatusBadRequest)
		return
	}

	h.Logger.Info("Скачивание файла", zap.String("filename", fileName))

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")

	err := h.FileService.DownloadFile(r.Context(), fileName, w)
	if err != nil {
		h.Logger.Error("Ошибка при скачивании файла", zap.Error(err))
		http.Error(w, "Ошибка при скачивании файла", http.StatusInternalServerError)
		return
	}
}
