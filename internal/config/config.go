package config

import (
	"os"

	"github.com/Levan1e/url-shortener-service/pkg/http"
	"github.com/Levan1e/url-shortener-service/pkg/postgres"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Postgres      *postgres.Config `yaml:"postgres"`
	Server        *http.Config     `yaml:"server"`
	MigrationsDir string           `yaml:"migrations_dir"`
}

func ParseConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	target := &Config{}
	err = yaml.NewDecoder(f).Decode(target)
	return target, err
}
