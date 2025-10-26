package database

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func NewMinioClient() (*MinioClient, error) {
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

func (m *MinioClient) UploadPhoto(ctx context.Context, userID int, fileData []byte, fileName string) (string, error) {
	// Генерируем уникальное имя файла с учетом userID
	objectName := m.generateImageName(userID, fileName)

	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	if contentType == "" {
		contentType = "image/jpeg"
	}

	_, err := m.Client.PutObject(ctx, m.BucketName, objectName,
		bytes.NewReader(fileData), int64(len(fileData)),
		minio.PutObjectOptions{
			ContentType: contentType,
		})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to MinIO: %v", err)
	}

	return m.getImageURL(objectName), nil
}

func (m *MinioClient) DeletePhoto(ctx context.Context, photoURL string) error {
	// Извлекаем objectName из полного URL
	objectName := m.extractObjectNameFromURL(photoURL)

	err := m.Client.RemoveObject(ctx, m.BucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %v", err)
	}
	return nil
}

func (m *MinioClient) GetPhotoURL(ctx context.Context, photoURL string) (string, error) {
	// Для MinIO URL обычно статический, но можно добавить логику генерации signed URL если нужно
	// Пока просто возвращаем тот же URL
	return photoURL, nil
}

// Вспомогательные методы

func (m *MinioClient) getImageURL(objectName string) string {
	return fmt.Sprintf("http://%s/%s/%s", m.Endpoint, m.BucketName, objectName)
}

func (m *MinioClient) generateImageName(userID int, originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Unix()
	// Добавляем userID в путь для организации файлов по пользователям
	return fmt.Sprintf("images/user_%d/%d%s", userID, timestamp, ext)
}

func (m *MinioClient) extractObjectNameFromURL(photoURL string) string {
	// Извлекаем objectName из полного URL
	// Пример: http://localhost:9000/comet-images/images/user_1/1234567890.jpg -> images/user_1/1234567890.jpg
	prefix := fmt.Sprintf("http://%s/%s/", m.Endpoint, m.BucketName)
	if strings.HasPrefix(photoURL, prefix) {
		return strings.TrimPrefix(photoURL, prefix)
	}
	return photoURL // Если не удалось извлечь, возвращаем как есть
}
