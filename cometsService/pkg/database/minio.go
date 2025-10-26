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
	Client         *minio.Client
	BucketName     string
	Endpoint       string
	PublicEndpoint string // URL для доступа из браузера
}

type Config struct {
	Endpoint        string
	PublicEndpoint  string
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

	// Публичный endpoint для доступа из браузера
	publicEndpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT")
	if publicEndpoint == "" {
		// По умолчанию используем localhost вместо имени docker-контейнера
		publicEndpoint = strings.Replace(endpoint, "minio", "localhost", 1)
	}

	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	useSSL, _ := strconv.ParseBool(os.Getenv("MINIO_USE_SSL"))

	cfg := Config{
		Endpoint:        endpoint,
		PublicEndpoint:  publicEndpoint,
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

	minioClient := &MinioClient{
		Client:         client,
		BucketName:     minioCfg.BucketName,
		Endpoint:       minioCfg.Endpoint,
		PublicEndpoint: minioCfg.PublicEndpoint,
	}

	// Создаем бакет, если он не существует
	if err := minioClient.ensureBucketExists(); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %v", err)
	}

	return minioClient, nil
}

// ensureBucketExists проверяет существование бакета и создает его при необходимости
func (m *MinioClient) ensureBucketExists() error {
	ctx := context.Background()
	
	// Проверяем, существует ли бакет
	exists, err := m.Client.BucketExists(ctx, m.BucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %v", err)
	}

	// Если бакет не существует, создаем его
	if !exists {
		err = m.Client.MakeBucket(ctx, m.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
		
		// Устанавливаем публичную политику для чтения
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, m.BucketName)
		
		err = m.Client.SetBucketPolicy(ctx, m.BucketName, policy)
		if err != nil {
			// Не критическая ошибка, логируем и продолжаем
			fmt.Printf("Warning: failed to set bucket policy: %v\n", err)
		}
	}

	return nil
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
	// Используем PublicEndpoint для доступа из браузера
	return fmt.Sprintf("http://%s/%s/%s", m.PublicEndpoint, m.BucketName, objectName)
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
	
	// Пробуем с публичным endpoint
	prefix := fmt.Sprintf("http://%s/%s/", m.PublicEndpoint, m.BucketName)
	if strings.HasPrefix(photoURL, prefix) {
		return strings.TrimPrefix(photoURL, prefix)
	}
	
	// Пробуем с внутренним endpoint (для обратной совместимости)
	prefix = fmt.Sprintf("http://%s/%s/", m.Endpoint, m.BucketName)
	if strings.HasPrefix(photoURL, prefix) {
		return strings.TrimPrefix(photoURL, prefix)
	}
	
	return photoURL // Если не удалось извлечь, возвращаем как есть
}
