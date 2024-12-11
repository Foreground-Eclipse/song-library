package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `env:"ENV"`
	HTTPServer
	DBData
}

// DBData is a struct that represents DB data in config
type DBData struct {
	DBHost     string `env:"DB_HOST"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
	DBPort     string `env:"DB_PORT"`
	DBSSLMode  string `env:"DB_SSLMODE"`
}

// HTTPServer is a struct that represents HTTP Server data in config
type HTTPServer struct {
	Address     string        `env:"HTTP_SERVER_ADDRESS"`
	Timeout     time.Duration `env:"HTTP_SERVER_TIMEOUT"`
	IdleTimeout time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT"`
}

// MustLoad loads the config
func MustLoad() *Config {
	configPath, err := filepath.Abs("./config/config.env")
	if err != nil {
		log.Fatalf("error getting absolute path to config file: %s", err)
	}
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesnt exists %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cant read config: %s", err)
	}

	return &cfg

}
