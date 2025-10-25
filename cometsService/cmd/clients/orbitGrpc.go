package clients

import (
	"context"
	"time"

	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/domain"
	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/grpc/cometorbit/proto"
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
		RaanDeg:              response.InclinationDeg,
		AscendingNodeLong:    response.RaanDeg,
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

// Close закрывает соединение
func (c *RealOrbitCalculationClient) Close() error {
	return c.conn.Close()
}