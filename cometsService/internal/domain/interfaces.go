package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound              = errors.New("resource not found")
	ErrNotEnoughObservations = errors.New("not enough observations for orbit calculation")
	ErrUnauthorized          = errors.New("unauthorized access")
	ErrInvalidInput          = errors.New("invalid input data")
	ErrOrbitNotCalculated    = errors.New("orbit not calculated for this comet")
)

type ICometsRepository interface {
	CreateComets(ctx context.Context, comet *Comet) error
	GetCometsByID(ctx context.Context, id int) (*Comet, error)
	GetCometsByUserID(ctx context.Context, userID int) ([]*Comet, error)
	UpdateComets(ctx context.Context, comet *Comet) error
	DeleteComets(ctx context.Context, id int, userID int) error

	CreateObservation(ctx context.Context, observation *Observation) error
	GetObservationByID(ctx context.Context, id int) (*Observation, error)
	GetUserObservationsByCometID(ctx context.Context, cometID int, userID int) ([]*Observation, error)
	UpdateObservation(ctx context.Context, observation *Observation) error
	DeleteObservation(ctx context.Context, id int, userID int) error
}

// ICometsService интерфейс для сервиса комет и наблюдений
type ICometsService interface {
	// Observation methods
	CreateObservation(ctx context.Context, userID int, req *CreateObservationRequest) (*Observation, error)
	GetObservation(ctx context.Context, id int) (*Observation, error)
	GetUserObservations(ctx context.Context, userID int) ([]*Observation, error)
	GetUserObservationsByCometID(ctx context.Context, cometID int, userID int) ([]*Observation, error)
	UpdateObservation(ctx context.Context, userID, id int, req *UpdateObservationRequest) error
	DeleteObservation(ctx context.Context, id int, userID int) error

	// Comet methods
	CreateComet(ctx context.Context, userID int, name string, fileData []byte, fileName string) (*Comet, error)
	GetComet(ctx context.Context, id int) (*Comet, error)
	GetUserComets(ctx context.Context, userID int) ([]*Comet, error)
	DeleteComet(ctx context.Context, id int, userID int) error

	// Calculation methods
	CalculateOrbit(ctx context.Context, userID, cometID int) (*CometOrbitResponse, error)
	CalculateCloseApproach(ctx context.Context, userID, cometID int) (*CometDistanceResponse, error)
	GetTrajectory(ctx context.Context, userID, cometID int, startTime, endTime time.Time, numPoints int) (*Trajectory, error)

	// File upload methods
	UploadCometPhoto(ctx context.Context, userID, cometID int, fileData []byte, fileName string) (*Comet, error)
}

// AuthClient интерфейс для сервиса авторизации
type IAuthClient interface {
	VerifyToken(token string) (bool, int32, error)
}

// OrbitCalculationClient интерфейс для сервиса расчетов орбит
type IOrbitCalculationClient interface {
	CalculateOrbit(ctx context.Context, observations []*Observation) (*OrbitalElements, error)
	CalculateCloseApproach(ctx context.Context, observations []*Observation) (*CloseApproach, error)
	GetTrajectory(ctx context.Context, observations []*Observation, startTime, endTime time.Time, numPoints int) (*Trajectory, error)
}

// FileStorageClient интерфейс для сервиса хранения файлов
type IFileStorageClient interface {
	UploadPhoto(ctx context.Context, userID int, fileData []byte, fileName string) (string, error)
	DeletePhoto(ctx context.Context, photoURL string) error
	GetPhotoURL(ctx context.Context, photoURL string) (string, error)
}
