package main

import (
	"github.com/codegangsta/cli"

	"github.com/exago/runner/task"
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
