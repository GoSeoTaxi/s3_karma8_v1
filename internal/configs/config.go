package configs

import (
	"os"
	"strconv"
	"strings"

	"github.com/GoSeoTaxi/s3_karma8_v1/internal/repositories"
)

type AppConfig struct {
	StorageConfig     repositories.StorageConfig
	ServerAddress     string
	PartSizeThreshold int64
}

func LoadConfig() *AppConfig {
	serverAddress := getEnv("SERVER_ADDRESS", ":8080")
	minioEndpoints := getEnv("MINIO_ENDPOINTS", "minio1:9000,minio2:9000,minio3:9000,minio4:9000,minio5:9000,minio6:9000")
	endpoints := strings.Split(minioEndpoints, ",")
	minioBucketName := getEnv("MINIO_BUCKET_NAME", "files")
	accessKey := getEnv("MINIO_ACCESS_KEY", "minioadmin")
	secretKey := getEnv("MINIO_SECRET_KEY", "minioadmin")
	useSSL := getEnvAsBool("MINIO_USE_SSL", false)
	partSizeThreshold := getEnvAsInt64("PART_SIZE_THRESHOLD", 1024*1024)

	return &AppConfig{
		ServerAddress:     serverAddress,
		PartSizeThreshold: partSizeThreshold,
		StorageConfig: repositories.StorageConfig{
			Endpoints:       endpoints,
			AccessKey:       accessKey,
			SecretKey:       secretKey,
			UseSSL:          useSSL,
			MinioBucketName: minioBucketName,
		},
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	if valueStr, exists := os.LookupEnv(name); exists {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultVal
}

func getEnvAsInt64(name string, defaultVal int64) int64 {
	if valueStr, exists := os.LookupEnv(name); exists {
		if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			return value
		}
	}
	return defaultVal
}
