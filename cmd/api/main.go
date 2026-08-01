package main

import (
	"github.com/nghiem-pham/bookmark-management/internal/api"
	"github.com/rs/zerolog/log"
)

func main() {
	// create app config
	cfg, err := api.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	engine, err := api.NewEngine(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init engine")
	}

	if err := engine.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
