package main

import (
	"fmt"
	"os"

	"github.com/barnigator/eshop-seller-service/internal/app"
	"github.com/barnigator/eshop-seller-service/internal/config"
)

func main() {
	cfg := config.MustLoad()

	application := app.New(cfg)

	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run application: %v\n", err)
		os.Exit(1)
	}
}
