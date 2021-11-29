package config

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

type PostgreSql struct {
	Host     string `yaml:"host" default:"localhost"`
	Port     int64  `yaml:"port" default:"27017"`
	Database string `yaml:"database" default:"control"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (m *PostgreSql) Validate() (err error) {
	return validation.ValidateStruct(m,
		validation.Field(&m.Host, validation.Required),
		validation.Field(&m.Port, validation.Required),
		validation.Field(&m.Database, validation.Required),
	)
}
