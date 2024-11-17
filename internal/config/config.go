package config

import (
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Tokens           string  `env:"TOKENS,required"`
	BitkubApiBaseUrl string  `env:"BITKUB_API_BASE_URL,required"`
	BitkubApiKey     string  `env:"BITKUB_API_KEY,required"`
	BitkubApiSecret  string  `env:"BITKUB_API_SECRET,required"`
	StartTimestamp   *uint64 `env:"START_TIMESTAMP"`
}

func initEnv() {
	envPath := os.Getenv("DOTENV_PATH")
	if envPath == "" {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	} else {
		err := godotenv.Load(envPath)
		if err != nil {
			panic(err)
		}
	}
}

func NewConfig() *Config {
	initEnv()
	config := &Config{}
	if err := env.Parse(config); err != nil {
		panic(err)
	}

	return config
}
