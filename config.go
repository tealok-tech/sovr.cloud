package main

import (
	"os"
)

type Config struct {
	RelyingPartyDisplayName string
	RelyingPartyID          string
	RelyingPartyOrigins     []string
}

func CreateConfig() *Config {
	c := Config{
		RelyingPartyDisplayName: "localhost",
		RelyingPartyID:          "localhost",
		RelyingPartyOrigins:     []string{"http://localhost:8080"},
	}
	val := os.Getenv("RELYING_PARTY_DISPLAY_NAME")
	if val != "" {
		c.RelyingPartyDisplayName = val
	}
	return &c
}
