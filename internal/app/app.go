package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/barnigator/eshop-seller-service/internal/config"
	"github.com/barnigator/eshop-seller-service/internal/grpc/handler"
	"github.com/barnigator/eshop-seller-service/internal/grpc/server"
	"github.com/barnigator/eshop-seller-service/internal/repository/postgres"
	"github.com/barnigator/eshop-seller-service/internal/usecase"
	"google.golang.org/grpc"
)

type App struct {
	cfg *config.Config
	log *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

func (a *App) Run() error {
	a.log.Info(
		"starting application",
		slog.Int("grpc_port", a.cfg.GRPC.Port),
		slog.Bool("postgres_configured", a.cfg.Postgres.DSN != ""),
	)

	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.App.Timeout)
	defer cancel()

	pool, err := postgres.NewPool(ctx, a.cfg.Postgres.DSN)
	if err != nil {
		return fmt.Errorf("initialize postgres pool: %w", err)
	}
	defer func() {
		a.log.Info("closing postgres pool")
		pool.Close()
		a.log.Info("postgres pool closed")
	}()

	a.log.Info("postgres connection established")

	repo := postgres.New(pool)

	useCase := usecase.New(repo, repo)

	handlers := handler.New(useCase, useCase)

	grpcServer := server.New(a.cfg.GRPC.Port, handlers, a.log)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- grpcServer.Run()
	}()

	a.log.Info(
		"starting grpc server",
		slog.Int("port", a.cfg.GRPC.Port),
	)

	shutdownCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err = <-serverErrors:
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return fmt.Errorf("grpc server stopped unexpectedly: %w", err)
		}

		return nil

	case <-shutdownCtx.Done():
		a.log.Info("shutdown signal received")
	}

	a.log.Info("stopping grpc server")
	grpcServer.Stop()
	a.log.Info("grpc server stopped")

	return nil
}
