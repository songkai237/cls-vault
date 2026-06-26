package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/clsVault/keeper/internal/config"
	"github.com/clsVault/keeper/internal/keeper"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	k, err := keeper.New(cfg)
	if err != nil {
		log.Fatalf("keeper init: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Printf("keeper started poll=%s dryRun=%v", cfg.PollInterval, cfg.DryRun)
	if err := k.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("keeper stopped: %v", err)
	}
	log.Printf("keeper stopped")
}
