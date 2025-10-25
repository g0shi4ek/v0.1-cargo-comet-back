package clients

import (
	"context"

	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/grpc/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Реализация реального клиента аутентификации
type RealAuthClient struct {
	conn   *grpc.ClientConn
	client auth.AuthServiceClient
}

// NewRealAuthClient создает новый клиент аутентификации
func NewRealAuthClient(authServiceAddr string) (*RealAuthClient, error) {
	conn, err := grpc.Dial(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := auth.NewAuthServiceClient(conn)
	return &RealAuthClient{
		conn:   conn,
		client: client,
	}, nil
}

// VerifyToken проверяет токен через gRPC сервис аутентификации
func (c *RealAuthClient) VerifyToken(token string) (bool, int32, error) {
	ctx := context.Background()
	request := &auth.VerifyTokenRequest{
		Token: token,
	}

	response, err := c.client.VerifyToken(ctx, request)
	if err != nil {
		return false, 0, err
	}

	return response.Valid, response.UserId, nil
}

// Close закрывает соединение
func (c *RealAuthClient) Close() error {
	return c.conn.Close()
}