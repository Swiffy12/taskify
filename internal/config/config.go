package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:"3000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
	Storage     `yaml:"postgresql"`
	Clients     ClientsConfig `yaml:"clients"`
	AppSecret   string        `yaml:"app_secret" env-required:"true" env:"APP_SECRET"`
	AppId       int32         `yaml:"app_id" env-default:2`
}

type Storage struct {
	DBHost     string `yaml:"host" env-required:"true"`
	DBPort     string `yaml:"port" env-required:"true"`
	DBUser     string `yaml:"user" env-required:"true"`
	DBPassword string `yaml:"password" env-required:"true"`
	DBName     string `yaml:"name" env-required:"true"`
}

type Client struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retries_count"`
}

type ClientsConfig struct {
	SSO Client `yaml:"sso"`
}

func MustLoad() *Config {
	configPath := os.Getenv("TASKIFY_CONFIG_PATH")

	if configPath == "" {
		log.Fatal("TASKIFY_CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config file %s", configPath)
	}

	return &cfg
}
