package app

import (
	"context"
	"fmt"

	"github.com/barnigator/eshop-seller-service/internal/config"
	"github.com/barnigator/eshop-seller-service/internal/grpc/handler"
	"github.com/barnigator/eshop-seller-service/internal/grpc/server"
	"github.com/barnigator/eshop-seller-service/internal/storage/memory"
	"github.com/barnigator/eshop-seller-service/internal/storage/postgres"
	"github.com/barnigator/eshop-seller-service/internal/usecase"
)

type App struct {
	cfg *config.Config
}

func New(cfg *config.Config) *App {
	return &App{cfg}
}

func (a *App) Run() error {
	fmt.Printf(
		"seller service starting: env=%s grpc_port=%d postgres_configured=%t\n",
		a.cfg.Env,
		a.cfg.GRPC.Port,
		a.cfg.Postgres.DSN != "",
	)

	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.GRPC.Timeout)
	defer cancel()

	pool, err := postgres.New(ctx, a.cfg.Postgres.DSN)
	if err != nil {
		return err
	}
	defer pool.Close()

	sellerRepo := memory.New()

	sellerUseCase := usecase.New(sellerRepo)

	sellerHandler := handler.New(sellerUseCase)

	grpcServer := server.New(a.cfg.GRPC.Port, sellerHandler)

	return grpcServer.Run()
}
