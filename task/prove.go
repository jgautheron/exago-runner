// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package task

import (
	"time"

	"github.com/karolgorecki/goprove"
)

type proveRunner struct {
	Runner
}

// ProveRunner launches goprove
func ProveRunner(m *Manager) Runnable {
	return &proveRunner{
		Runner{Label: "Go Prove", Mgr: m},
	}
}

// Execute goprove
func (r *proveRunner) Execute() error {
	defer r.trackTime(time.Now())

	passed, failed := goprove.RunTasks(r.Manager().RepositoryPath(), []string{"projectBuilds"})

	r.Data = struct {
		Passed []map[string]interface{} `json:"passed"`
		Failed []map[string]interface{} `json:"failed"`
	}{
		passed, failed,
	}

	return nil
}
