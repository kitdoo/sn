package config

import (
	"os"

	"github.com/go-ozzo/ozzo-validation/v4"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Log        *Logger     `yaml:"logger"`
	GRPC       *GRPC       `yaml:"grpc"`
	PostgreSql *PostgreSql `yaml:"postgresql"`
}

func Load(path string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Log, validation.Required),
		validation.Field(&c.GRPC, validation.Required),
		validation.Field(&c.PostgreSql, validation.Required),
	)
}
