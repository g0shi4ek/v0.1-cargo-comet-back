package repository

import (
	"context"

	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/domain"
	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/pkg/database"
	"gorm.io/gorm"
)

type CometsRepository struct {
	db *gorm.DB
}

func NewCometsRepository() *CometsRepository {
	postgresClient, err := database.NewPostgresClient()
	if err != nil {
		return nil
	}
	return &CometsRepository{
		db: postgresClient,
	}
}

func (r *CometsRepository) CreateComets(ctx context.Context, comet *domain.Comet) error {
	return r.db.WithContext(ctx).Create(comet).Error
}

func (r *CometsRepository) GetCometsByID(ctx context.Context, id int) (*domain.Comet, error) {
	var comet domain.Comet
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&comet).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &comet, err
}

func (r *CometsRepository) GetCometsByUserID(ctx context.Context, userID int) ([]*domain.Comet, error) {
	var comets []*domain.Comet
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&comets).Error
	return comets, err
}

func (r *CometsRepository) UpdateComets(ctx context.Context, comet *domain.Comet) error {
	return r.db.WithContext(ctx).Save(comet).Error
}

func (r *CometsRepository) DeleteComets(ctx context.Context, id int, userID int) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.Comet{}).Error
}

func (r *CometsRepository) CreateObservation(ctx context.Context, observation *domain.Observation) error {
	return r.db.WithContext(ctx).Create(observation).Error
}

func (r *CometsRepository) GetObservationByID(ctx context.Context, id int) (*domain.Observation, error) {
	var observation domain.Observation
	err := r.db.WithContext(ctx).
		Preload("Comet").
		Where("id = ?", id).
		First(&observation).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &observation, err
}

func (r *CometsRepository) GetUserObservationsByCometID(ctx context.Context, cometID int, userID int) ([]*domain.Observation, error) {
	var observations []*domain.Observation
	err := r.db.WithContext(ctx).
		Where("comet_id = ? AND user_id = ?", cometID, userID).
		Order("observed_at ASC").
		Find(&observations).Error
	return observations, err
}

func (r *CometsRepository) UpdateObservation(ctx context.Context, observation *domain.Observation) error {
	return r.db.WithContext(ctx).Save(observation).Error
}

func (r *CometsRepository) DeleteObservation(ctx context.Context, id int, userID int) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.Observation{}).Error
}
