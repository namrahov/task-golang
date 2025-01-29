package util

import (
	"context"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"task-golang/config"
)

// UploadFileToMinio uploads a file to the specified Minio bucket.
func UploadFileToMinio(ctx context.Context, uniqueName string, multipartFile multipart.File, size int64, minioClient *minio.Client) error {

	// Upload file to Minio
	_, err := minioClient.PutObject(ctx, config.Props.MinioBucket, uniqueName, multipartFile, size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
