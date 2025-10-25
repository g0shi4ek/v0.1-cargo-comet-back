package database

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client     *minio.Client
	BucketName string
	Endpoint   string
}

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	UseSSL          bool
}

func MinioConfigFromEnv() (*Config, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("failed to load minio config")
	}

	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	useSSL, _ := strconv.ParseBool(os.Getenv("MINIO_USE_SSL"))

	cfg := Config{
		Endpoint:        endpoint,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		BucketName:      bucketName,
		UseSSL:          useSSL,
	}

	return &cfg, nil
}

func NewCometClient() (*MinioClient, error) {
	minioCfg, err := MinioConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load minio config")
	}

	client, err := minio.New(minioCfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioCfg.AccessKeyID, minioCfg.SecretAccessKey, ""),
		Secure: minioCfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %v", err)
	}

	return &MinioClient{
		Client:     client,
		BucketName: minioCfg.BucketName,
		Endpoint:   minioCfg.Endpoint,
	}, nil
}

func (m *MinioClient) UploadCometImage(ctx context.Context, imageData []byte, filename string) (string, error) {
	tariffImageName := m.generateTariffImageName(filename)

	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = "image/jpeg"
	}

	_, err := m.Client.PutObject(ctx, m.BucketName, tariffImageName,
		bytes.NewReader(imageData), int64(len(imageData)),
		minio.PutObjectOptions{
			ContentType: contentType,
		})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to MinIO: %v", err)
	}

	return m.getTariffImageURL(tariffImageName), nil
}

func (m *MinioClient) DeleteTariffImage(ctx context.Context, tariffImageName string) error {
	err := m.Client.RemoveObject(ctx, m.BucketName, "images/"+tariffImageName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %v", err)
	}
	return nil
}

func (m *MinioClient) getTariffImageURL(tariffImageName string) string {
	return fmt.Sprintf("http://%s/%s/%s", m.Endpoint, m.BucketName, tariffImageName)
}

func (m *MinioClient) generateTariffImageName(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("images/%d%s", timestamp, ext)
}