package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	ssov1 "github.com/Swiffy12/protos/gen/go/sso"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api   ssov1.AuthClient
	AppId int32
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
	appId int32,
) (*Client, error) {
	const op = "grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api:   ssov1.NewAuthClient(cc),
		AppId: appId,
	}, nil
}

func (c *Client) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "grpc.IsAdmin"

	resp, err := c.api.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: userId,
	})
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetIsAdmin(), nil
}

func (c *Client) Register(ctx context.Context, email, password string) (int64, error) {
	const op = "grpc.Register"

	resp, err := c.api.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetUserId(), nil
}

func (c *Client) Login(ctx context.Context, email, password string, appId int32) (string, error) {
	const op = "grpc.Login"

	resp, err := c.api.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appId,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetToken(), nil
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(level), msg, fields...)
	})
}
