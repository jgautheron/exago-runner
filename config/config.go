// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

// Config contains global configuration
var Config struct {
	LogLevel       string `envconfig:"LOG_LEVEL" default:"info"`
	Repository     string `envconfig:"REPOSITORY" required:"true"`
	RepositoryPath string
}

// InitializeConfig loads configuration using envconfig
func InitializeConfig() {
	if err := envconfig.Process("", &Config); err != nil {
		log.Fatal(err)
	}
}
