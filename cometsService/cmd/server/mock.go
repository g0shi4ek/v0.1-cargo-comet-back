package main

import (
	"context"
	"time"

	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/domain"
)

// MockOrbitCalculationClient заглушка для клиента расчетов орбит
type MockOrbitCalculationClient struct{}

func NewMockOrbitCalculationClient() *MockOrbitCalculationClient {
	return &MockOrbitCalculationClient{}
}

func (m *MockOrbitCalculationClient) CalculateOrbit(ctx context.Context, observations []*domain.Observation) (*domain.OrbitalElements, error) {
	// Имитация расчета орбитальных элементов на основе наблюдений
	if len(observations) < 5 {
		return nil, domain.ErrNotEnoughObservations
	}

	return &domain.OrbitalElements{
		SemiMajorAxis:        17.8,
		Eccentricity:         0.967,
		RaanDeg:              162.3,
		AscendingNodeLong:    58.42,
		ArgumentOfPerihelion: 111.33,
		TrueAnomalyDeg:       162.3,
	}, nil
}

func (m *MockOrbitCalculationClient) CalculateCloseApproach(ctx context.Context, orbitalElements []*domain.Observation) (*domain.CloseApproach, error) {
	// Имитация расчета сближения
	return &domain.CloseApproach{
		Date:     time.Now().Add(30 * 24 * time.Hour), // через 30 дней
		Distance: 0.2,                                 // 0.2 а.е.
	}, nil
}

// MockFileStorageClient заглушка для клиента хранения файлов
type MockFileStorageClient struct{}

func NewMockFileStorageClient() *MockFileStorageClient {
	return &MockFileStorageClient{}
}

func (m *MockFileStorageClient) UploadPhoto(ctx context.Context, userID int, fileData []byte, fileName string) (string, error) {
	return "https://storage.example.com/photos/" + fileName, nil
}

func (m *MockFileStorageClient) DeletePhoto(ctx context.Context, photoURL string) error {
	return nil
}

func (m *MockFileStorageClient) GetPhotoURL(ctx context.Context, photoURL string) (string, error) {
	return photoURL, nil
}

// MockAuthClient заглушка для клиента авторизации
type MockAuthClient struct{}

func NewMockAuthClient() *MockAuthClient {
	return &MockAuthClient{}
}

func (m *MockAuthClient) VerifyToken(token string) (bool, int32, error) {
	// В реальной реализации здесь должна быть проверка JWT токена
	// Для демонстрации возвращаем фиксированный userID
	if token == "" {
		return false, 0, domain.ErrUnauthorized
	}
	return true, 1, nil // Возвращаем userID = 1 для всех валидных токенов
}
