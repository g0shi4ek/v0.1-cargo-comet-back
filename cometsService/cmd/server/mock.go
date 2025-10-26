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
		InclinationDeg:    58.42,
		ArgumentOfPerihelion: 111.33,
		TrueAnomalyDeg:       162.3,
	}, nil
}

func (m *MockOrbitCalculationClient) CalculateCloseApproach(ctx context.Context, observations []*domain.Observation) (*domain.CloseApproach, error) {
	// Имитация расчета сближения
	return &domain.CloseApproach{
		Date:     time.Now().Add(30 * 24 * time.Hour), // через 30 дней
		Distance: 0.2,                                 // 0.2 а.е.
	}, nil
}

func (m *MockOrbitCalculationClient) GetTrajectory(ctx context.Context, observations []*domain.Observation, startTime, endTime time.Time, numPoints int) (*domain.Trajectory, error) {
	// Имитация расчета траектории
	if len(observations) < 3 {
		return nil, domain.ErrNotEnoughObservations
	}

	// Генерируем mock траекторию кометы
	cometTrajectory := make([]domain.TrajectoryPoint, numPoints)
	duration := endTime.Sub(startTime)

	for i := 0; i < numPoints; i++ {
		pointTime := startTime.Add(time.Duration(i) * duration / time.Duration(numPoints-1))

		// Простая эллиптическая орбита для демонстрации
		angle := 2 * 3.14159 * float64(i) / float64(numPoints)
		radius := 2.0 + 1.5*float64(i)/float64(numPoints) // Радиус увеличивается со временем

		cometTrajectory[i] = domain.TrajectoryPoint{
			Time: pointTime,
			X:    radius * cos(angle),
			Y:    radius * sin(angle) * 0.5, // Немного сжатая по Y
			Z:    radius * sin(angle) * 0.3, // И по Z
		}
	}

	// Генерируем mock траекторию Земли (круговая орбита)
	earthTrajectory := make([]domain.TrajectoryPoint, numPoints)

	for i := 0; i < numPoints; i++ {
		pointTime := startTime.Add(time.Duration(i) * duration / time.Duration(numPoints-1))
		angle := 2 * 3.14159 * float64(i) / float64(numPoints)

		earthTrajectory[i] = domain.TrajectoryPoint{
			Time: pointTime,
			X:    cos(angle), // Круговая орбита радиусом 1 а.е.
			Y:    sin(angle),
			Z:    0, // Земля движется примерно в плоскости эклиптики
		}
	}

	return &domain.Trajectory{
		CometTrajectory: cometTrajectory,
		EarthTrajectory: earthTrajectory,
	}, nil
}

// Вспомогательные функции для тригонометрии
func cos(x float64) float64 {
	// Простая реализация косинуса через ряд Тейлора для демонстрации
	result := 1.0
	term := 1.0
	x2 := x * x
	for i := 1; i <= 6; i++ {
		term *= -x2 / float64(2*i*(2*i-1))
		result += term
	}
	return result
}

func sin(x float64) float64 {
	// Простая реализация синуса через ряд Тейлора для демонстрации
	result := x
	term := x
	x2 := x * x
	for i := 1; i <= 6; i++ {
		term *= -x2 / float64(2*i*(2*i+1))
		result += term
	}
	return result
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
