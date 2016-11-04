// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package task

import (
	"os"
	"time"

	"github.com/hotolab/cov"
)

type coverageRunner struct {
	Runner
	tempFile *os.File
}

// CoverageRunner is a runner used for testing Go projects
func CoverageRunner(m *Manager) Runnable {
	return &coverageRunner{
		Runner: Runner{Label: "Code Coverage", Mgr: m},
	}
}

// Execute gets all the coverage files and returns the output of
// hotolab/cov
func (r *coverageRunner) Execute() {
	defer r.trackTime(time.Now())
	rep, err := cov.ConvertRepository(r.Manager().Repository())
	if err != nil {
		r.toRunnerError(err)
		return
	}

	r.Data = rep
}
