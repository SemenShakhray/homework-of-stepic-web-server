package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	SecretJWT   string        `yaml:"secret_jwt" env-required:"true"`
}

func MustLoad() Config {
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		log.Fatal("file config don't exists")
	}

	var cfg Config

	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		log.Fatal("failed parse config", err)
	}

	return cfg
}
