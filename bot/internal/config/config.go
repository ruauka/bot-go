// Package config - config
package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config .
type Config struct {
	TelegaToken string `env:"TELEGA_TOKEN"`
	PgUser      string `env:"POSTGRES_USER"`
	PgPassword  string `env:"POSTGRES_PASSWORD"`
	PgDbName    string `env:"POSTGRES_DB"`
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
