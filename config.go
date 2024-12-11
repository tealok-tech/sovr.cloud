package main

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
	return &c
}
