package main

import (
	"context"
	"log"
	"time"

	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthGRPCClient gRPC клиент для сервиса авторизации
type AuthGRPCClient struct {
	client AuthServiceClient
	conn   *grpc.ClientConn
}

// NewAuthGRPCClient создает новый gRPC клиент для авторизации
func NewAuthGRPCClient(addr string) (*AuthGRPCClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := NewAuthServiceClient(conn)
	
	return &AuthGRPCClient{
		client: client,
		conn:   conn,
	}, nil
}

// VerifyToken проверяет JWT токен через gRPC
func (a *AuthGRPCClient) VerifyToken(ctx context.Context, token string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &VerifyTokenRequest{
		Token: token,
	}

	resp, err := a.client.VerifyToken(ctx, req)
	if err != nil {
		return 0, err
	}

	if !resp.Valid {
		return 0, domain.ErrUnauthorized
	}

	return int(resp.UserId), nil
}

// GetUserPermissions получает разрешения пользователя (заглушка)
func (a *AuthGRPCClient) GetUserPermissions(ctx context.Context, userID int) ([]string, error) {
	// В реальной реализации здесь должен быть gRPC вызов
	// Пока возвращаем базовые разрешения
	return []string{"read:comets", "write:comets", "read:observations", "write:observations"}, nil
}

// Close закрывает соединение
func (a *AuthGRPCClient) Close() error {
	return a.conn.Close()
}

// OrbitCalculationGRPCClient gRPC клиент для сервиса вычислений орбит
type OrbitCalculationGRPCClient struct {
	client CalculationServiceClient
	conn   *grpc.ClientConn
}

// NewOrbitCalculationGRPCClient создает новый gRPC клиент для вычислений
func NewOrbitCalculationGRPCClient(addr string) (*OrbitCalculationGRPCClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := NewCalculationServiceClient(conn)
	
	return &OrbitCalculationGRPCClient{
		client: client,
		conn:   conn,
	}, nil
}

// CalculateOrbit вычисляет орбитальные элементы через gRPC
func (o *OrbitCalculationGRPCClient) CalculateOrbit(ctx context.Context, observations []*domain.Observation) (*domain.OrbitalElements, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Преобразуем наблюдения в gRPC формат
	observationProtos := make([]*ObservationProto, len(observations))
	for i, obs := range observations {
		observationProtos[i] = &ObservationProto{
			Id:             int32(obs.ID),
			UserId:         int32(obs.UserID),
			CometId:        int32(*obs.CometID),
			RightAscension: obs.RightAscension,
			Declination:    obs.Declination,
			ObservedAt:     obs.ObservedAt.Format(time.RFC3339),
		}
	}

	req := &CalculateOrbitRequest{
		Observations: observationProtos,
	}

	resp, err := o.client.CalculateOrbit(ctx, req)
	if err != nil {
		return nil, err
	}

	// Преобразуем ответ в доменную структуру
	timeOfPerihelion, err := time.Parse(time.RFC3339, resp.OrbitalElements.TimeOfPerihelion)
	if err != nil {
		return nil, err
	}

	return &domain.OrbitalElements{
		SemiMajorAxis:        resp.OrbitalElements.SemiMajorAxis,
		Eccentricity:         resp.OrbitalElements.Eccentricity,
		Inclination:          resp.OrbitalElements.Inclination,
		AscendingNodeLong:    resp.OrbitalElements.AscendingNodeLong,
		ArgumentOfPerihelion: resp.OrbitalElements.ArgumentOfPerihelion,
		TimeOfPerihelion:     timeOfPerihelion,
	}, nil
}

// CalculateCloseApproach вычисляет сближение через gRPC
func (o *OrbitCalculationGRPCClient) CalculateCloseApproach(ctx context.Context, orbitalElements *domain.OrbitalElements) (*domain.CloseApproach, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req := &CalculateCloseApproachRequest{
		OrbitalElements: &OrbitalElementsProto{
			SemiMajorAxis:        orbitalElements.SemiMajorAxis,
			Eccentricity:         orbitalElements.Eccentricity,
			Inclination:          orbitalElements.Inclination,
			AscendingNodeLong:    orbitalElements.AscendingNodeLong,
			ArgumentOfPerihelion: orbitalElements.ArgumentOfPerihelion,
			TimeOfPerihelion:     orbitalElements.TimeOfPerihelion.Format(time.RFC3339),
		},
	}

	resp, err := o.client.CalculateCloseApproach(ctx, req)
	if err != nil {
		return nil, err
	}

	// Преобразуем ответ в доменную структуру
	date, err := time.Parse(time.RFC3339, resp.CloseApproach.Date)
	if err != nil {
		return nil, err
	}

	return &domain.CloseApproach{
		Date:     date,
		Distance: resp.CloseApproach.Distance,
	}, nil
}

// GetCalculationStatus получает статус расчета через gRPC
func (o *OrbitCalculationGRPCClient) GetCalculationStatus(ctx context.Context, requestID int) (*domain.CalculationRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &GetCalculationStatusRequest{
		RequestId: int32(requestID),
	}

	resp, err := o.client.GetCalculationStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	return &domain.CalculationRequest{
		ID:           int(resp.CalculationRequest.Id),
		UserID:       int(resp.CalculationRequest.UserId),
		CometID:      int(resp.CalculationRequest.CometId),
		Status:       resp.CalculationRequest.Status,
		ErrorMessage: resp.CalculationRequest.ErrorMessage,
	}, nil
}

// Close закрывает соединение
func (o *OrbitCalculationGRPCClient) Close() error {
	return o.conn.Close()
}