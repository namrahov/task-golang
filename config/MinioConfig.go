package config

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// NewMinioClient initializes and returns a Minio client
func NewMinioClient() (*minio.Client, error) {
	client, err := minio.New(Props.MinioUrl, &minio.Options{
		Creds:  credentials.NewStaticV4(Props.MinioAccessKey, Props.MinioSecretKey, ""),
		Secure: Props.MinioUseSsl,
	})
	if err != nil {
		fmt.Printf("Failed to initialize Minio client: %v", err)
		return nil, err
	}
	return client, nil
}
