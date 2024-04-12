package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env-default:"dev"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	DataBase    `yaml:"database"`
	JWTToken    `yaml:"jwt_token"`
	Redis       `yaml:"redis"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DataBase struct {
	Name     string `yaml:"name" env-required:"true"`
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	Type     string `yaml:"type" env-required:"true"`
}

type JWTToken struct {
	Secret   string `yaml:"secret" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	TTL      int    `yaml:"ttl" env-default:"3600"`
}

type Redis struct {
	Address string `yaml:"address" env-default:"localhost:6379"`
}

func MustLoad() *Config {

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		fmt.Println("Config path env var not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		fmt.Println("Error opening config file:", err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		fmt.Println("Error reading config file:", err)
	}

	return &cfg
}
