// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package main

import (
	"github.com/codegangsta/cli"

	"github.com/hotolab/exago-runner/task"
)

// RunnerCommand starts the task runner
func RunnerCommand() cli.Command {
	return cli.Command{
		Name:  "runner",
		Usage: "Start the runners",
		Action: func(c *cli.Context) error {
			m := task.NewManager()
			return m.ExecuteRunners()
		},
	}
}
