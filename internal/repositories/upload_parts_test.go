package repositories

import (
	"context"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestListFileParts(t *testing.T) {
	mockClient := new(MockMinioClient)

	ctx := context.Background()
	bucketName := "test-bucket"
	fileName := "test_file"

	objectCh := make(chan minio.ObjectInfo, 3)
	objectCh <- minio.ObjectInfo{Key: fileName + "_part_1", Size: 2048}
	objectCh <- minio.ObjectInfo{Key: fileName + "_part_0", Size: 1024}
	objectCh <- minio.ObjectInfo{Key: "other_file_part_2", Size: 512}
	close(objectCh)

	mockClient.On("ListObjects", ctx, bucketName, mock.AnythingOfType("minio.ListObjectsOptions")).
		Return(objectCh).
		Once()

	repo := &StorageRepository{
		Clients:    []MinioClient{mockClient},
		Logger:     zap.NewNop(),
		BucketName: bucketName,
	}

	parts, err := repo.ListFileParts(ctx, fileName)
	assert.NoError(t, err)

	assert.Len(t, parts, 2)

	expectedParts := []struct {
		PartNumber int
		Size       int64
	}{
		{PartNumber: 0, Size: 1024},
		{PartNumber: 1, Size: 2048},
	}

	for i, part := range parts {
		assert.Equal(t, expectedParts[i].PartNumber, part.PartNumber, "PartNumber не совпадает для части %d", i)
		assert.Equal(t, expectedParts[i].Size, part.Size, "Size не совпадает для части %d", i)
	}

	for i := 1; i < len(parts); i++ {
		assert.GreaterOrEqual(t, parts[i].PartNumber, parts[i-1].PartNumber, "Части не отсортированы по PartNumber")
	}

	mockClient.AssertCalled(t, "ListObjects", ctx, bucketName, mock.AnythingOfType("minio.ListObjectsOptions"))

	mockClient.AssertExpectations(t)
}
