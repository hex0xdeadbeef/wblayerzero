package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Environment `yaml:"env"`
		Postgre     `yaml:"database"`
	}

	Environment struct {
		Status string `yaml:"status" env-default:"prod"`
	}

	Postgre struct {
		Name string `yaml:"name" env-default:"posgtres"`

		Host string `yaml:"host" env-default:"localhost"`
		Port string `yaml:"5432" env-default:"5432"`

		User     string `yaml:"user" env-default:"posgtres"`
		Password string `yaml:"password" env-default:"posgtres"`
	}
)

const (
	CfgFilePath = "../../config.yaml"
)

// Load loads configuration and returns a pointer to a Config structure and an error if any	
func Load(cfgPath string) (*Config, error) {
	var (
		cfg Config
	)

	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("loading config from %q: %v", cfgPath, err)
	}

	return &cfg, nil
}
