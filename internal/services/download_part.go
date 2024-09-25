package services

import (
	"context"
	"fmt"
	"io"
)

func (fs *FileService) DownloadPart(ctx context.Context, fileName string, partNumber int, writer io.Writer) error {

	parts, err := fs.StorageRepo.ListFileParts(ctx, fileName)
	if err != nil {
		return err
	}
	if partNumber >= len(parts) || partNumber < 0 {
		return fmt.Errorf("part %d not found", partNumber)
	}

	objectName := fmt.Sprintf("%s_part_%d", fileName, partNumber)
	return fs.StorageRepo.DownloadPart(ctx, objectName, writer)
}
