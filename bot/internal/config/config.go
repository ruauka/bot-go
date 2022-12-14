// Package config - config
package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config .
type Config struct {
	TelegaTokenProd string `env:"TELEGA_TOKEN_PROD"`
	TelegaTokenDev  string `env:"TELEGA_TOKEN_DEV"`
	PgUser          string `env:"POSTGRES_USER"`
	PgPassword      string `env:"POSTGRES_PASSWORD"`
	PgDbName        string `env:"POSTGRES_DB"`
	Level           string `env:"LEVEL"`
}

// GetConfig - get env vars.
func GetConfig() *Config {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Println(err)
		help, err := cleanenv.GetDescription(cfg, nil)
		if err != nil {
			log.Println(err)
		}
		log.Println(help)
	}

	return cfg
}
