package repositories

import (
	"context"
	"math/rand"
	"strconv"
	"strings"

	"github.com/minio/minio-go/v7"
)

func (s *StorageRepository) getRandomClient() MinioClient {
	index := rand.Intn(len(s.Clients))
	return s.Clients[index]
}

func extractPartIndex(objectName, fileName string) int {
	partStr := strings.TrimPrefix(objectName, fileName+"_part_")
	index, err := strconv.Atoi(partStr)
	if err != nil {
		return -1
	}
	return index
}

func ensureBucketExists(client MinioClient, bucketName string) error {
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
