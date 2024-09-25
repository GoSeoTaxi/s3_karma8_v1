package repositories

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioClient interface {
	ListObjects(ctx context.Context, bucketName string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error)
	MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error
	ListBuckets(ctx context.Context) ([]minio.BucketInfo, error)
}
