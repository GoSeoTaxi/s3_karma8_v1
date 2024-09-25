package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// Тут я хотел придумать, что каждая часть лежит только на своем сервере, но в итоге переключился на MinIO
func (fs *FileService) DetermineClientIndex(objectName string) int {
	partIndex := extractPartIndex(objectName)
	if partIndex == -1 {
		fs.Logger.Error("Не удалось извлечь индекс части из имени объекта", zap.String("objectName", objectName))
		return 0
	}

	clientIndex := int(partIndex % int64(len(fs.StorageRepo.Clients)))
	return clientIndex
}

func extractPartIndex(objectName string) int64 {
	parts := strings.Split(objectName, "_part_")
	if len(parts) != 2 {
		return -1 // Возвращаем -1 в случае ошибки
	}

	partIndex, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return -1
	}

	return partIndex
}

func removeDuplicates(input []string) []string {
	if len(input) == 0 {
		return input
	}

	sort.Strings(input)
	result := []string{input[0]}

	for i := 1; i < len(input); i++ {
		if input[i] != input[i-1] {
			result = append(result, input[i])
		}
	}

	return result
}

func (fs *FileService) uploadFileWithoutSize(ctx context.Context, filename string, data io.Reader) error {
	partIndex := int64(0)
	buf := make([]byte, fs.PartSizeThreshold)
	totalUploaded := int64(0)

	for {
		n, err := io.ReadFull(data, buf)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				if n > 0 {
					objectName := fmt.Sprintf("%s_part_%d", filename, partIndex)

					reader := bytes.NewReader(buf[:n])

					err = fs.StorageRepo.UploadPart(ctx, objectName, reader, int64(n))
					if err != nil {
						fs.Logger.Error("Не удалось загрузить часть", zap.Int64("part", partIndex), zap.Error(err))
						return fmt.Errorf("не удалось загрузить часть %d: %w", partIndex, err)
					}

					fs.Logger.Info("Часть загружена", zap.String("objectName", objectName), zap.Int64("size", int64(n)))
					partIndex++
					totalUploaded += int64(n)
				}
				break
			} else {
				fs.Logger.Error("Ошибка при чтении данных", zap.Error(err))
				return err
			}
		}

		objectName := fmt.Sprintf("%s_part_%d", filename, partIndex)

		reader := bytes.NewReader(buf[:n])

		err = fs.StorageRepo.UploadPart(ctx, objectName, reader, int64(n))
		if err != nil {
			fs.Logger.Error("Не удалось загрузить часть", zap.Int64("part", partIndex), zap.Error(err))
			return fmt.Errorf("не удалось загрузить часть %d: %w", partIndex, err)
		}

		fs.Logger.Info("Часть загружена", zap.String("objectName", objectName), zap.Int64("size", int64(n)))
		partIndex++
		totalUploaded += int64(n)
	}

	fs.Logger.Info("Файл успешно загружен", zap.String("filename", filename), zap.Int64("totalUploaded", totalUploaded))
	return nil
}
