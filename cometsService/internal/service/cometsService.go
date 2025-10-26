package service

import (
	"context"
	"log"
	"time"

	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/domain"
	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/pkg/database"
)

type CometsService struct {
	cometRepo         domain.ICometsRepository
	orbitCalcClient   domain.IOrbitCalculationClient
	fileStorageClient domain.IFileStorageClient
}

func NewCometsService(
	cometRepo domain.ICometsRepository,
	orbitCalcClient domain.IOrbitCalculationClient,
) *CometsService {
	minio, err := database.NewMinioClient()
	if err != nil {
		return nil
	}
	return &CometsService{
		cometRepo:         cometRepo,
		orbitCalcClient:   orbitCalcClient,
		fileStorageClient: minio,
	}
}

// Observation methods
func (s *CometsService) CreateObservation(ctx context.Context, userID int, req *domain.CreateObservationRequest) (*domain.Observation, error) {
	observedAt, err := time.Parse(time.RFC3339, req.ObservedAt)
	if err != nil {
		return nil, domain.ErrInvalidInput
	}

	observation := &domain.Observation{
		UserID:         userID,
		CometID:        req.CometID,
		RightAscension: req.RightAscension,
		Declination:    req.Declination,
		ObservedAt:     observedAt,
	}

	if err := s.cometRepo.CreateObservation(ctx, observation); err != nil {
		return nil, err
	}

	if req.CometID != nil {
		if err := s.resetCalculationFlags(ctx, *req.CometID, userID); err != nil {
			log.Printf("Warning: failed to reset calculation flags: %v", err)
		}
	}

	return observation, nil
}

func (s *CometsService) GetObservation(ctx context.Context, id int) (*domain.Observation, error) {
	return s.cometRepo.GetObservationByID(ctx, id)
}

func (s *CometsService) GetUserObservations(ctx context.Context, userID int) ([]*domain.Observation, error) {
	// Получаем все кометы пользователя
	comets, err := s.cometRepo.GetCometsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Собираем все наблюдения для всех комет пользователя
	var allObservations []*domain.Observation
	for _, comet := range comets {
		observations, err := s.cometRepo.GetUserObservationsByCometID(ctx, comet.ID, userID)
		if err != nil {
			return nil, err
		}
		allObservations = append(allObservations, observations...)
	}

	return allObservations, nil
}

func (s *CometsService) GetUserObservationsByCometID(ctx context.Context, cometID int, userID int) ([]*domain.Observation, error) {
	return s.cometRepo.GetUserObservationsByCometID(ctx, cometID, userID)
}

func (s *CometsService) UpdateObservation(ctx context.Context, userID, id int, req *domain.UpdateObservationRequest) error {
	// Сначала получаем существующее наблюдение
	existingObservation, err := s.cometRepo.GetObservationByID(ctx, id)
	if err != nil {
		return err
	}
	if existingObservation == nil {
		return domain.ErrNotFound
	}

	// Проверяем права доступа
	if existingObservation.UserID != userID {
		return domain.ErrUnauthorized
	}

	observedAt, err := time.Parse(time.RFC3339, req.ObservedAt)
	if err != nil {
		return domain.ErrInvalidInput
	}

	observation := &domain.Observation{
		ID:             id,
		UserID:         userID,
		CometID:        existingObservation.CometID,
		RightAscension: req.RightAscension,
		Declination:    req.Declination,
		ObservedAt:     observedAt,
	}

	err = s.cometRepo.UpdateObservation(ctx, observation)
	if err != nil {
		return err
	}

	// Сбрасываем флаги расчетов у кометы, если наблюдение привязано к комете
	if existingObservation.CometID != nil {
		if err := s.resetCalculationFlags(ctx, *existingObservation.CometID, userID); err != nil {
			log.Printf("Warning: failed to reset calculation flags: %v", err)
		}
	}

	return nil
}

func (s *CometsService) DeleteObservation(ctx context.Context, id int, userID int) error {
	observation, err := s.cometRepo.GetObservationByID(ctx, id)
	if err != nil {
		return err
	}
	if observation == nil {
		return domain.ErrNotFound
	}

	// Проверяем права доступа
	if observation.UserID != userID {
		return domain.ErrUnauthorized
	}

	// Удаляем наблюдение
	if err := s.cometRepo.DeleteObservation(ctx, id, userID); err != nil {
		return err
	}

	// Сбрасываем флаги расчетов у кометы, если наблюдение было привязано к комете
	if observation.CometID != nil {
		if err := s.resetCalculationFlags(ctx, *observation.CometID, userID); err != nil {
			log.Printf("Warning: failed to reset calculation flags: %v", err)
		}
	}

	return nil
}

// Comet methods
func (s *CometsService) CreateComet(ctx context.Context, userID int, name string, fileData []byte, fileName string) (*domain.Comet, error) {
	var photoURL string
	var err error

	// Упрощенная проверка - len() для nil слайсов возвращает 0
	if len(fileData) > 0 {
		photoURL, err = s.fileStorageClient.UploadPhoto(ctx, userID, fileData, fileName)
		if err != nil {
			return nil, err
		}
	}

	comet := &domain.Comet{
		UserID:   userID,
		Name:     name,
		PhotoURL: photoURL,
	}

	if err := s.cometRepo.CreateComets(ctx, comet); err != nil {
		return nil, err
	}

	return comet, nil
}

func (s *CometsService) GetComet(ctx context.Context, id int) (*domain.Comet, error) {
	return s.cometRepo.GetCometsByID(ctx, id)
}

func (s *CometsService) GetUserComets(ctx context.Context, userID int) ([]*domain.Comet, error) {
	return s.cometRepo.GetCometsByUserID(ctx, userID)
}

func (s *CometsService) DeleteComet(ctx context.Context, id int, userID int) error {
	// Сначала получаем комету, чтобы узнать photoURL и проверить права
	comet, err := s.cometRepo.GetCometsByID(ctx, id)
	if err != nil {
		return err
	}
	if comet == nil {
		return domain.ErrNotFound
	}

	// Проверяем права доступа
	if comet.UserID != userID {
		return domain.ErrUnauthorized
	}

	// Удаляем фото из хранилища, если оно есть
	if comet.PhotoURL != "" {
		if err := s.fileStorageClient.DeletePhoto(ctx, comet.PhotoURL); err != nil {
			// Логируем ошибку, но не прерываем удаление кометы
			log.Printf("Warning: failed to delete photo for comet %d: %v", id, err)
		}
	}

	// Soft delete кометы (устанавливаем deleted_at)
	return s.cometRepo.DeleteComets(ctx, id, userID)
}

// Calculation methods
func (s *CometsService) CalculateOrbit(ctx context.Context, userID, cometID int) (*domain.CometOrbitResponse, error) {
	// Проверяем существование кометы и права доступа
	comet, err := s.cometRepo.GetCometsByID(ctx, cometID)
	if err != nil {
		return nil, err
	}
	if comet == nil {
		return nil, domain.ErrNotFound
	}

	if comet.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	// Получаем наблюдения для кометы
	observations, err := s.cometRepo.GetUserObservationsByCometID(ctx, cometID, userID)
	if err != nil {
		return nil, err
	}

	if len(observations) < 5 {
		return nil, domain.ErrNotEnoughObservations
	}

	// Вычисляем орбитальные элементы
	orbitalElements, err := s.orbitCalcClient.CalculateOrbit(ctx, observations)
	if err != nil {
		return nil, err
	}

	// Обновляем комету с новыми орбитальными элементами
	comet.SemiMajorAxis = orbitalElements.SemiMajorAxis
	comet.Eccentricity = orbitalElements.Eccentricity
	comet.RaanDeg = orbitalElements.RaanDeg
	comet.AscendingNodeLong = orbitalElements.AscendingNodeLong
	comet.ArgumentOfPerihelion = orbitalElements.ArgumentOfPerihelion
	comet.TrueAnomalyDeg = orbitalElements.TrueAnomalyDeg
	comet.OrbitActual = true // Устанавливаем флаг
	comet.CalculatedAt = time.Now()

	comet.CloseActual = false
	comet.MinApproachDate = nil
	comet.MinApproachDistance = nil
	if err := s.cometRepo.UpdateComets(ctx, comet); err != nil {
		return nil, err
	}

	// Формируем ответ
	response := &domain.CometOrbitResponse{
		ID:                   comet.ID,
		SemiMajorAxis:        &comet.SemiMajorAxis,
		Eccentricity:         &comet.Eccentricity,
		RaanDeg:              &comet.RaanDeg,
		AscendingNodeLong:    &comet.AscendingNodeLong,
		ArgumentOfPerihelion: &comet.ArgumentOfPerihelion,
		TrueAnomalyDeg:       &comet.TrueAnomalyDeg,
		OrbitActual:          comet.OrbitActual,
	}

	return response, nil
}

func (s *CometsService) CalculateCloseApproach(ctx context.Context, userID, cometID int) (*domain.CometDistanceResponse, error) {
	// Проверяем существование кометы и права доступа
	// Проверяем существование кометы и права доступа
	comet, err := s.cometRepo.GetCometsByID(ctx, cometID)
	if err != nil {
		return nil, err
	}
	if comet == nil {
		return nil, domain.ErrNotFound
	}

	if comet.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	if !comet.OrbitActual {
		return nil, domain.ErrOrbitNotCalculated
	}

	// Получаем наблюдения для кометы
	observations, err := s.cometRepo.GetUserObservationsByCometID(ctx, cometID, userID)
	if err != nil {
		return nil, err
	}

	if len(observations) < 5 {
		return nil, domain.ErrNotEnoughObservations
	}

	// Вычисляем сближение
	closeApproach, err := s.orbitCalcClient.CalculateCloseApproach(ctx, observations)
	if err != nil {
		return nil, err
	}

	// Обновляем комету с данными о сближении
	comet.MinApproachDate = &closeApproach.Date
	comet.MinApproachDistance = &closeApproach.Distance
	comet.CalculatedAt = time.Now()
	comet.CloseActual = true

	if err := s.cometRepo.UpdateComets(ctx, comet); err != nil {
		return nil, err
	}

	// Формируем ответ
	response := &domain.CometDistanceResponse{
		ID:                  comet.ID,
		MinApproachDate:     comet.MinApproachDate,
		MinApproachDistance: comet.MinApproachDistance,
		CalculatedAt:        comet.CalculatedAt,
		CloseActual:         comet.CloseActual,
	}

	return response, nil
}

// File upload methods
func (s *CometsService) UploadCometPhoto(ctx context.Context, userID, cometID int, fileData []byte, fileName string) (*domain.Comet, error) {
	comet, err := s.cometRepo.GetCometsByID(ctx, cometID)
	if err != nil {
		return nil, err
	}
	if comet == nil {
		return nil, domain.ErrNotFound
	}

	if comet.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	photoURL, err := s.fileStorageClient.UploadPhoto(ctx, userID, fileData, fileName)
	if err != nil {
		return nil, err
	}

	comet.PhotoURL = photoURL
	if err := s.cometRepo.UpdateComets(ctx, comet); err != nil {
		return nil, err
	}

	return comet, nil
}

func (s *CometsService) resetCalculationFlags(ctx context.Context, cometID int, userID int) error {
	comet, err := s.cometRepo.GetCometsByID(ctx, cometID)
	if err != nil {
		return err
	}
	if comet == nil {
		return nil // Комета уже удалена или не существует
	}

	// Проверяем права доступа
	if comet.UserID != userID {
		return domain.ErrUnauthorized
	}

	// Сбрасываем флаги только если они были true
	if comet.OrbitActual || comet.CloseActual {
		comet.OrbitActual = false
		comet.CloseActual = false
		return s.cometRepo.UpdateComets(ctx, comet)
	}

	return nil
}

// GetTrajectory получает траекторию кометы и Земли для визуализации
func (s *CometsService) GetTrajectory(ctx context.Context, userID, cometID int, startTime, endTime time.Time, numPoints int) (*domain.Trajectory, error) {
	// Проверяем существование кометы и права доступа
	comet, err := s.cometRepo.GetCometsByID(ctx, cometID)
	if err != nil {
		return nil, err
	}
	if comet == nil {
		return nil, domain.ErrNotFound
	}

	if comet.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	// Проверяем, что орбитальные элементы уже рассчитаны
	if !comet.OrbitActual {
		return nil, domain.ErrOrbitNotCalculated
	}

	// Получаем наблюдения для кометы
	observations, err := s.cometRepo.GetUserObservationsByCometID(ctx, cometID, userID)
	if err != nil {
		return nil, err
	}

	if len(observations) < 3 {
		return nil, domain.ErrNotEnoughObservations
	}

	// Получаем траекторию
	return s.orbitCalcClient.GetTrajectory(ctx, observations, startTime, endTime, numPoints)
}
