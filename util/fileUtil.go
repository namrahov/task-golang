package util

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"strings"
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

// DeleteFileFromMinio deletes a file from the specified Minio bucket.
func DeleteFileFromMinio(ctx context.Context, filePath string, minioClient *minio.Client) error {
	// Extract the object name from the file path
	pathParts := splitFilePath(filePath)
	if len(pathParts) < 2 {
		return fmt.Errorf("invalid file path format: %s", filePath)
	}
	objectName := pathParts[1]

	// Remove the object from Minio
	err := minioClient.RemoveObject(ctx, config.Props.MinioBucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Helper function to split file path
func splitFilePath(filePath string) []string {
	return strings.Split(filePath, "/")
}
