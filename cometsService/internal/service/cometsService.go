package service

import (
	"context"
	"time"

	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/domain"
)

type CometsService struct {
	cometRepo         domain.ICometsRepository
	orbitCalcClient   domain.IOrbitCalculationClient
	fileStorageClient domain.IFileStorageClient
}

func NewCometsService(
	cometRepo domain.ICometsRepository,
	orbitCalcClient domain.IOrbitCalculationClient,
	fileStorageClient domain.IFileStorageClient,
) *CometsService {
	return &CometsService{
		cometRepo:         cometRepo,
		orbitCalcClient:   orbitCalcClient,
		fileStorageClient: fileStorageClient,
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
		PhotoURL:       req.PhotoURL,
	}

	if err := s.cometRepo.CreateObservation(ctx, observation); err != nil {
		return nil, err
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

	return s.cometRepo.UpdateObservation(ctx, observation)
}

func (s *CometsService) DeleteObservation(ctx context.Context, id int, userID int) error {
	return s.cometRepo.DeleteObservation(ctx, id, userID)
}

// Comet methods
func (s *CometsService) CreateComet(ctx context.Context, userID int, req *domain.CreateCometRequest) (*domain.Comet, error) {
	// добавить добавление изображения
	comet := &domain.Comet{
		UserID: userID,
		Name:   req.Name,
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

	if len(observations) < 3 {
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
	comet.Inclination = orbitalElements.Inclination
	comet.AscendingNodeLong = orbitalElements.AscendingNodeLong
	comet.ArgumentOfPerihelion = orbitalElements.ArgumentOfPerihelion
	comet.TimeOfPerihelion = orbitalElements.TimeOfPerihelion
	comet.CalculatedAt = time.Now()

	if err := s.cometRepo.UpdateComets(ctx, comet); err != nil {
		return nil, err
	}

	// Формируем ответ
	response := &domain.CometOrbitResponse{
		ID:                   comet.ID,
		SemiMajorAxis:        &comet.SemiMajorAxis,
		Eccentricity:         &comet.Eccentricity,
		Inclination:          &comet.Inclination,
		AscendingNodeLong:    &comet.AscendingNodeLong,
		ArgumentOfPerihelion: &comet.ArgumentOfPerihelion,
		TimeOfPerihelion:     &comet.TimeOfPerihelion,
	}

	return response, nil
}

func (s *CometsService) CalculateCloseApproach(ctx context.Context, userID, cometID int) (*domain.CometDistanceResponse, error) {
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
	if comet.SemiMajorAxis == 0 || comet.Eccentricity == 0 {
		return nil, domain.ErrNotEnoughObservations
	}

	// Подготавливаем орбитальные элементы для расчета сближения
	orbitalElements := &domain.OrbitalElements{
		SemiMajorAxis:        comet.SemiMajorAxis,
		Eccentricity:         comet.Eccentricity,
		Inclination:          comet.Inclination,
		AscendingNodeLong:    comet.AscendingNodeLong,
		ArgumentOfPerihelion: comet.ArgumentOfPerihelion,
		TimeOfPerihelion:     comet.TimeOfPerihelion,
	}

	// Вычисляем сближение
	closeApproach, err := s.orbitCalcClient.CalculateCloseApproach(ctx, orbitalElements)
	if err != nil {
		return nil, err
	}

	// Обновляем комету с данными о сближении
	comet.MinApproachDate = &closeApproach.Date
	comet.MinApproachDistance = &closeApproach.Distance
	comet.CalculatedAt = time.Now()

	if err := s.cometRepo.UpdateComets(ctx, comet); err != nil {
		return nil, err
	}

	// Формируем ответ
	response := &domain.CometDistanceResponse{
		ID:                  comet.ID,
		MinApproachDate:     comet.MinApproachDate,
		MinApproachDistance: comet.MinApproachDistance,
		CalculatedAt:        comet.CalculatedAt,
	}

	return response, nil
}

func (s *CometsService) GetCalculationStatus(ctx context.Context, userID, requestID int) (*domain.CalculationRequest, error) {
	return s.orbitCalcClient.GetCalculationStatus(ctx, requestID)
}

// File upload methods
func (s *CometsService) UploadObservationPhoto(ctx context.Context, userID int, fileData []byte, fileName string) (string, error) {
	return s.fileStorageClient.UploadPhoto(ctx, userID, fileData, fileName)
}
