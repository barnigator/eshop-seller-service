package main

import (
	"log/slog"
	"os"

	"github.com/barnigator/eshop-seller-service/internal/app"
	"github.com/barnigator/eshop-seller-service/internal/config"
	"github.com/barnigator/eshop-seller-service/internal/logger"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.MustLoad()
	log := logger.New(cfg.Env)

	application := app.New(cfg, log)

	if err := application.Run(); err != nil {
		log.Error(
			"application stopped with error",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}
}
