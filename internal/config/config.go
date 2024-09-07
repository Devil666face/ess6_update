package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	https = "https"
)

type Config struct {
	WebIP         string `env:"WEB_IP" env-default:"127.0.0.1"`
	PortHTTP      uint   `env:"PORT_HTTP" env-default:"8000"`
	PortHTTPS     uint   `env:"PORT_HTTPS" env-default:"4443"`
	AllowHost     string `env:"ALLOW_HOST" env-default:"localhost"`
	ConnectHTTP   string
	HTTPSRedirect string
	ConnectHTTPS  string
}

func New() (*Config, error) {
	cfg := Config{}
	if err := cleanenv.ReadConfig("default.env", &cfg); err != nil {
		log.Println(err)
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, fmt.Errorf("env variable not found: %w", err)
		}
	}
	cfg.ConnectHTTP = fmt.Sprintf("%s:%d", cfg.WebIP, cfg.PortHTTP)
	cfg.ConnectHTTPS = fmt.Sprintf("%s:%d", cfg.WebIP, cfg.PortHTTPS)
	cfg.HTTPSRedirect = fmt.Sprintf("%s://%s:%d", https, cfg.AllowHost, cfg.PortHTTPS)
	return &cfg, nil
}
