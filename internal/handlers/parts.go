package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/GoSeoTaxi/s3_karma8_v1/internal/services"
)

type PartsHandler struct {
	fileService *services.FileService
	logger      *zap.Logger
}

func NewPartsHandler(fs *services.FileService, logger *zap.Logger) *PartsHandler {
	return &PartsHandler{
		fileService: fs,
		logger:      logger,
	}
}

func (h *PartsHandler) ListFileParts(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	if fileName == "" {
		http.Error(w, "fileName is required", http.StatusBadRequest)
		return
	}

	parts, err := h.fileService.ListFileParts(r.Context(), fileName)
	if err != nil {
		h.logger.Error("Не удалось получить список частей файла", zap.Error(err))
		http.Error(w, "Failed to list file parts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(parts)
}

func (h *PartsHandler) DownloadPart(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	if fileName == "" {
		http.Error(w, "fileName is required", http.StatusBadRequest)
		return
	}

	partNumberStr := chi.URLParam(r, "partNumber")
	if partNumberStr == "" {
		http.Error(w, "part_number is required", http.StatusBadRequest)
		return
	}

	partNumber, err := strconv.Atoi(partNumberStr)
	if err != nil {
		http.Error(w, "Invalid partNumber", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName+"_part_"+partNumberStr)
	w.Header().Set("Content-Type", "application/octet-stream")

	err = h.fileService.DownloadPart(r.Context(), fileName, partNumber, w)
	if err != nil {
		h.logger.Error("Не удалось скачать часть файла", zap.Error(err))
		http.Error(w, "Failed to download file part", http.StatusInternalServerError)
		return
	}
}
