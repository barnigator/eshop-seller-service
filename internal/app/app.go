package app

import (
	"fmt"

	"github.com/barnigator/eshop-seller-service/internal/config"
	"github.com/barnigator/eshop-seller-service/internal/grpc/handler"
	"github.com/barnigator/eshop-seller-service/internal/grpc/server"
)

type App struct {
	cfg *config.Config
}

func New(cfg *config.Config) *App {
	return &App{cfg}
}

func (a *App) Run() error {
	fmt.Printf(
		"seller service starting: env=%s grpc_port=%d\n",
		a.cfg.Env,
		a.cfg.GRPC.Port,
	)

	sellerHandler := handler.New()

	grcpServer := server.New(a.cfg.GRPC.Port, sellerHandler)

	return grcpServer.Run()
}
