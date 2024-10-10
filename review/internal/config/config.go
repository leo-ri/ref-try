package config

import (
	"coupon_service/internal/api"
	"log"

	"github.com/brumhard/alligotor"
)

type Config struct {
	API api.Config
}

func New() Config {
	cfg := Config{
		API: api.Config{
			Host: "localhost",
			Port: 8080,
		},
	}
	if err := alligotor.Get(&cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}
