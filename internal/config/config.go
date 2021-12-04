package config

import (
	"github.com/kovetskiy/ko"
	"gopkg.in/yaml.v2"
)

type Database struct {
	Name     string `yaml:"name" required:"true" env:"DATABASE_NAME"`
	Host     string `yaml:"host" required:"true" env:"DATABASE_HOST"`
	Port     string `yaml:"port" required:"true" env:"DATABASE_PORT"`
	User     string `yaml:"user" required:"true"`
	Password string `yaml:"password" required:"true"`
}

type Server struct {
	Port string `yaml:"port" required:"true"`
}

type Token struct {
	SecretKey string `yaml:"secret_key" required:"true"`
	ExpirationTimer int `yaml:"expiration_timer" required:"true"`
}

type SMS struct {
	ExpirationTimer int64 `yaml:"expiration_timer" required:"true"`
}

type Config struct {
	Database Database `yaml:"database" required:"true"`
	Server   Server   `yaml:"server" required:"true"`
	Token    Token    `yaml:"token" required:"true"`
	SMS      SMS      `yaml:"sms" required:"true"`
}

func Load(path string) (*Config, error) {
	config := &Config{}
	err := ko.Load(path, config, ko.RequireFile(false), yaml.Unmarshal)
	if err != nil {
		return nil, err
	}

	return config, nil
}
