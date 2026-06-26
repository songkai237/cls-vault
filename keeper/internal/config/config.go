package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	RPCURL          string
	PrivateKey      string
	StrategyAddress string
	PollInterval    time.Duration
	DryRun          bool
	ChainID         int64
	MinIdleAmount0  string // optional override; empty = use strategy min swap
	MinIdleAmount1  string
}

func Load() (Config, error) {
	cfg := Config{
		RPCURL:          strings.TrimSpace(os.Getenv("RPC_URL")),
		PrivateKey:      strings.TrimSpace(os.Getenv("PRIVATE_KEY")),
		StrategyAddress: strings.TrimSpace(os.Getenv("STRATEGY_ADDRESS")),
		PollInterval:    12 * time.Second,
		MinIdleAmount0:  strings.TrimSpace(os.Getenv("MIN_IDLE_AMOUNT0")),
		MinIdleAmount1:  strings.TrimSpace(os.Getenv("MIN_IDLE_AMOUNT1")),
	}

	if v := strings.TrimSpace(os.Getenv("POLL_INTERVAL")); v != "" {
		sec, err := strconv.Atoi(v)
		if err != nil || sec <= 0 {
			return cfg, fmt.Errorf("invalid POLL_INTERVAL: %q", v)
		}
		cfg.PollInterval = time.Duration(sec) * time.Second
	}

	if v := strings.TrimSpace(os.Getenv("DRY_RUN")); v == "1" || strings.EqualFold(v, "true") {
		cfg.DryRun = true
	}

	if v := strings.TrimSpace(os.Getenv("CHAIN_ID")); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return cfg, fmt.Errorf("invalid CHAIN_ID: %w", err)
		}
		cfg.ChainID = id
	}

	if cfg.RPCURL == "" {
		return cfg, fmt.Errorf("RPC_URL is required")
	}
	if cfg.PrivateKey == "" {
		return cfg, fmt.Errorf("PRIVATE_KEY is required")
	}
	if cfg.StrategyAddress == "" {
		return cfg, fmt.Errorf("STRATEGY_ADDRESS is required")
	}

	return cfg, nil
}
