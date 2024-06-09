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
		User     string `yaml:"user" env-default:"posgtres"`
		Password string `yaml:"password" env-default:"posgtres"`

		Host string `yaml:"host" env-default:"localhost"`
		Port string `yaml:"5432" env-default:"5432"`

		Name string `yaml:"name" env-default:"posgtres"`
	}
)

const (
	CfgFilePath = "../../config.yaml"
)

var (
	Cfg Config
)

// GenURI returns constructed URI to be used during conn
func (p *Postgre) GenURI() string {
	const (
		format = "postgres://%s:%s@%s:%s/%s"
	)

	return fmt.Sprintf(format, p.User, p.Password, p.Host, p.Port, p.Name)
}

// Load loads configuration and returns a pointer to a Config structure and an error if any
func Load(cfgPath string) error {
	err := cleanenv.ReadConfig(cfgPath, &Cfg)
	if err != nil {
		return fmt.Errorf("loading config from %q: %v", cfgPath, err)
	}

	return nil
}
