// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/hotolab/exago-runner/task"
)

var (
	App *cli.App
)

// Initialize commandline app.
func init() {
	App = cli.NewApp()

	// For fancy output on console
	App.Name = "exago runner"
	App.Usage = `Check -h`
	App.Author = "Hotolab <dev@hotolab.com>"

	// Version is injected at build-time
	App.Version = ""
	App.Action = func(c *cli.Context) error {
		m := task.NewManager(c.Args().Get(0))

		if c.String("ref") != "" {
			m.UseReference(c.String("ref"))
		}

		out := m.ExecuteRunners()
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(out)

		return nil
	}

	App.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "ref",
			Usage: "reference passed when cloning (branch or SHA1)",
		},
	}

	InitializeLogging(os.Getenv("LOG_LEVEL"))
}

func main() {
	if err := App.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// InitializeLogging sets logrus log level.
func InitializeLogging(logLevel string) {
	// If log level cannot be resolved, exit gracefully
	if logLevel == "" {
		log.SetLevel(log.InfoLevel)
		return
	}
	// Parse level from string
	lvl, err := log.ParseLevel(logLevel)

	if err != nil {
		log.WithFields(log.Fields{
			"passed":  logLevel,
			"default": "fatal",
		}).Warn("Log level is not valid, fallback to default level")
		log.SetLevel(log.FatalLevel)
		return
	}

	log.SetLevel(lvl)
	log.WithFields(log.Fields{
		"level": logLevel,
	}).Debug("Log level successfully set")
}
