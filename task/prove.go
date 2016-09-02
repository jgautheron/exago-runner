// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package task

import (
	"encoding/json"
	"os/exec"
	"time"
)

type proveRunner struct {
	Runner
}

// ProveRunner launches goprove
func ProveRunner() Runnable {
	return &proveRunner{Runner{Label: "Go Prove", parallel: true}}
}

// Execute goprove
func (r *proveRunner) Execute() {
	defer r.trackTime(time.Now())

	checklist := map[string][]map[string]string{}
	cl, err := exec.Command("goprove", "-output", "json", "-exclude", "testPassing", ".").CombinedOutput()
	if err != nil {
		r.toRunnerError(err)
	}
	json.Unmarshal(cl, &checklist)

	r.RawOutput = string(cl)
	r.Data = checklist
}
