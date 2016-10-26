package configs

import (
	"github.com/caarlos0/env"
)

type Config struct {
	GOENV                 string `env:"GOENV" envDefault:"development"`
	MONGO_PRIMARY_URL     string `env:"MONGO_PRIMARY_URL" envDefault:"localhost:27017"`
	DEFAULT_DATABASE_NAME string `env:"DEFAULT_DATABASE_NAME" envDefault:"venus"`
	REDIS_HOST            string `env:"REDIS_HOST"`
	REDIS_PORT            string `env:"REDIS_PORT"`
}

var AppConfig = Config{}

func init() {
	err := env.Parse(&AppConfig)
	if err != nil {
		panic(err)
	}
}
