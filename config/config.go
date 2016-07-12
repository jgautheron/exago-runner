package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/exago/envconfig"
)

var Config cfg

type cfg struct {
	LogLevel       string `envconfig:"LOG_LEVEL" default:"info"`
	Repository     string `envconfig:"REPOSITORY" required:"true"`
	RepositoryPath string
}

func InitializeConfig() {
	if err := envconfig.Process("", &Config); err != nil {
		log.Fatal(err)
	}
}
