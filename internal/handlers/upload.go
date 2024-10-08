package handlers

import (
	"context"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/GoSeoTaxi/s3_karma8_v1/internal/services"
)

type UploadHandler struct {
	FileService *services.FileService
	Logger      *zap.Logger
}

func NewUploadHandler(fs *services.FileService, logger *zap.Logger) *UploadHandler {
	return &UploadHandler{
		FileService: fs,
		Logger:      logger,
	}
}

func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	multipartReader, err := r.MultipartReader()
	if err != nil {
		h.Logger.Error("Ошибка при разборе формы", zap.Error(err))
		http.Error(w, "Ошибка при разборе формы", http.StatusBadRequest)
		return
	}

	for {
		select {
		case <-r.Context().Done():
			h.Logger.Info("Загрузка файла отменена клиентом")
			http.Error(w, "Загрузка отменена", http.StatusRequestTimeout)
			return
		default:
		}

		part, err := multipartReader.NextPart()
		if err == io.EOF {
			// Все части обработаны
			break
		}
		if err != nil {
			h.Logger.Error("Ошибка при чтении части", zap.Error(err))
			http.Error(w, "Ошибка при чтении части", http.StatusInternalServerError)
			return
		}

		if part.FormName() != "file" {
			continue
		}

		filename := part.FileName()
		if filename == "" {
			h.Logger.Error("Не указано имя файла")
			http.Error(w, "Не указано имя файла", http.StatusBadRequest)
			return
		}

		pr, pw := io.Pipe()
		uploadErrChan := make(chan error, 1)

		go func() {
			defer func() { _ = pr.Close() }()
			err := h.FileService.UploadFile(r.Context(), filename, pr, -1)
			uploadErrChan <- err
		}()

		_, err = ioCopyContext(r.Context(), pw, part)

		if err != nil {
			_ = pw.CloseWithError(err)
			h.Logger.Error("Ошибка при копировании данных в pipe", zap.Error(err))
			http.Error(w, "Ошибка при загрузке файла", http.StatusInternalServerError)
			return
		}
		_ = pw.Close()

		uploadErr := <-uploadErrChan
		if uploadErr != nil {
			h.Logger.Error("Ошибка при загрузке файла", zap.Error(uploadErr))
			http.Error(w, "Ошибка при загрузке файла", http.StatusInternalServerError)
			return
		}

		h.Logger.Info("Файл успешно загружен", zap.String("filename", filename))
		break
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Файл успешно загружен!\n"))
}

func ioCopyContext(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	done := make(chan error, 1)

	var n int64

	go func() {
		var err error
		n, err = io.Copy(dst, src)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return n, ctx.Err()
	case err := <-done:
		return n, err
	}
}
