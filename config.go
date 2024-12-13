package main

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	RelyingPartyDisplayName string   `env:"RELYING_PARTY_DISPLAY_NAME"`
	RelyingPartyID          string   `env:"RELYING_PARTY_ID"`
	RelyingPartyOrigins     []string `env:"RELYING_PARTY_ORIGINS"`
	SessionSecret           string   `env:"SESSION_SECRET"`
}

func CreateConfig() *Config {
	var c Config
	ctx := context.Background()
	if err := envconfig.Process(ctx, &c); err != nil {
		panic(err)
	}
	return &c
}
