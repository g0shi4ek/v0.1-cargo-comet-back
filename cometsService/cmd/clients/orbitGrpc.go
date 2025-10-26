package clients

import (
	"context"
	"time"

	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/domain"
	cometorbit "github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/grpc/cometorbit/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RealOrbitCalculationClient struct {
	conn   *grpc.ClientConn
	client cometorbit.OrbitServiceClient
}

// NewRealOrbitCalculationClient создает новый клиент расчета орбиты
func NewRealOrbitCalculationClient(orbitServiceAddr string) (*RealOrbitCalculationClient, error) {
	conn, err := grpc.Dial(orbitServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := cometorbit.NewOrbitServiceClient(conn)
	return &RealOrbitCalculationClient{
		conn:   conn,
		client: client,
	}, nil
}

// CalculateOrbit вычисляет орбитальные элементы на основе наблюдений
func (c *RealOrbitCalculationClient) CalculateOrbit(ctx context.Context, observations []*domain.Observation) (*domain.OrbitalElements, error) {
	// Конвертируем наблюдения в формат gRPC
	grpcObservations := make([]*cometorbit.Observation, len(observations))
	for i, obs := range observations {
		grpcObservations[i] = &cometorbit.Observation{
			TimeUtc:      obs.ObservedAt.Format("2006-01-02 15:04:05"),
			RaDeg:        obs.RightAscension,
			DecDeg:       obs.Declination,
			IsHorizontal: obs.IsHorizontal, // предполагаем экваториальные координаты
		}
	}

	request := &cometorbit.ObservationsRequest{
		Observations: grpcObservations,
	}

	response, err := c.client.CalculateKeplerianElements(ctx, request)
	if err != nil {
		return nil, err
	}

	// Конвертируем ответ в доменный формат
	return &domain.OrbitalElements{
		SemiMajorAxis:        response.SemiMajorAxisAu,
		Eccentricity:         response.Eccentricity,
		RaanDeg:              response.RaanDeg,
		InclinationDeg:   	  response.InclinationDeg,
		ArgumentOfPerihelion: response.ArgOfPeriapsisDeg,
		TrueAnomalyDeg:       response.TrueAnomalyDeg,
	}, nil
}

// CalculateCloseApproach вычисляет ближайшее сближение с Землей
func (c *RealOrbitCalculationClient) CalculateCloseApproach(ctx context.Context, observations []*domain.Observation) (*domain.CloseApproach, error) {
	grpcObservations := make([]*cometorbit.Observation, len(observations))
	for i, obs := range observations {
		grpcObservations[i] = &cometorbit.Observation{
			TimeUtc:      obs.ObservedAt.Format("2006-01-02 15:04:05"),
			RaDeg:        obs.RightAscension,
			DecDeg:       obs.Declination,
			IsHorizontal: obs.IsHorizontal, // предполагаем экваториальные координаты
		}
	}

	request := &cometorbit.ObservationsRequest{
		Observations: grpcObservations,
	}

	response, err := c.client.GetClosestApproach(ctx, request)
	if err != nil {
		return nil, err
	}

	date, err := time.Parse("2006-01-02 15:04:05", response.TimeUtc)
	if err != nil {
		return nil, err
	}

	return &domain.CloseApproach{
		Date:     date,
		Distance: response.DistanceAu,
	}, nil
}

// GetTrajectory получает траекторию кометы и Земли для визуализации
func (c *RealOrbitCalculationClient) GetTrajectory(ctx context.Context, observations []*domain.Observation, startTime, endTime time.Time, numPoints int) (*domain.Trajectory, error) {
	// Конвертируем наблюдения в формат gRPC
	grpcObservations := make([]*cometorbit.Observation, len(observations))
	for i, obs := range observations {
		grpcObservations[i] = &cometorbit.Observation{
			TimeUtc:      obs.ObservedAt.Format("2006-01-02 15:04:05"),
			RaDeg:        obs.RightAscension,
			DecDeg:       obs.Declination,
			IsHorizontal: obs.IsHorizontal,
		}
	}

	request := &cometorbit.TrajectoryRequest{
		Observations: &cometorbit.ObservationsRequest{
			Observations: grpcObservations,
		},
		StartTimeUtc: startTime.Format("2006-01-02 15:04:05"),
		EndTimeUtc:   endTime.Format("2016-01-02 15:04:05"),
		NumPoints:    int32(numPoints),
	}

	response, err := c.client.GetTrajectory(ctx, request)
	if err != nil {
		return nil, err
	}

	// Конвертируем траекторию кометы
	cometTrajectory := make([]domain.TrajectoryPoint, len(response.CometTrajectory))
	for i, point := range response.CometTrajectory {
		time, err := time.Parse("2006-01-02 15:04:05", point.TimeUtc)
		if err != nil {
			return nil, err
		}
		cometTrajectory[i] = domain.TrajectoryPoint{
			Time: time,
			X:    point.XAu,
			Y:    point.YAu,
			Z:    point.ZAu,
		}
	}

	// Конвертируем траекторию Земли
	earthTrajectory := make([]domain.TrajectoryPoint, len(response.EarthTrajectory))
	for i, point := range response.EarthTrajectory {
		time, err := time.Parse("2006-01-02 15:04:05", point.TimeUtc)
		if err != nil {
			return nil, err
		}
		earthTrajectory[i] = domain.TrajectoryPoint{
			Time: time,
			X:    point.XAu,
			Y:    point.YAu,
			Z:    point.ZAu,
		}
	}

	return &domain.Trajectory{
		CometTrajectory: cometTrajectory,
		EarthTrajectory: earthTrajectory,
	}, nil
}

// Close закрывает соединение
func (c *RealOrbitCalculationClient) Close() error {
	return c.conn.Close()
}