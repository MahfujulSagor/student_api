package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type HTTPServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

type Config struct {
	Env              string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	DBPath           string `yaml:"db_path" env-required:"true"`
	HTTPServerConfig `yaml:"http_server"`
	LoggingConfig    `yaml:"logging"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "Path to configuration file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatal("Configuration file path must be provided via CONFIG_PATH environment variable or --config flag")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Configuration file does not exist:", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("Failed to read configuration file:", err.Error())
	}

	return &cfg
}
