// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package task

import (
	"time"

	"github.com/karolgorecki/goprove"

	. "github.com/hotolab/exago-runner/config"
)

type proveRunner struct {
	Runner
}

// ProveRunner launches goprove
func ProveRunner() Runnable {
	return &proveRunner{Runner{Label: "Go Prove"}}
}

// Execute goprove
func (r *proveRunner) Execute() {
	defer r.trackTime(time.Now())

	passed, failed := goprove.RunTasks(Config.RepositoryPath, []string{"projectBuilds"})

	r.Data = struct {
		Passed []map[string]interface{} `json:"passed"`
		Failed []map[string]interface{} `json:"failed"`
	}{
		passed, failed,
	}
}
