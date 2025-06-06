package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	StoragePath      string `yaml:"storage_path" env-required:"true"`
	TelegramAPIToken string `yaml:"telegram_api_token" env-required:"true"`
	WebhookURL       string `yaml:"web_hook" env-required:"true"`
	ListenAdr        string `yaml:"http_listen_adr" env-required:"true"`
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
