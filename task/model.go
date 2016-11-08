// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package task

import "time"

const (
	downloadName     = "download"
	testName         = "test"
	coverageName     = "coverage"
	proveName        = "goprove"
	thirdPartiesName = "thirdparties"
	locName          = "golocc"
)

// Runner is the struct holding all informations about the runner
type Runner struct {
	// Label is the name of the task runner
	// This is the only field that must be set
	Label string `json:"label"`

	// Data holds the specialized object associated to the task
	// runner i.e. specialized object for Goprove and Gotest
	Data interface{} `json:"data"`

	// RawOutput is the process's standard output and error.
	// It is used for system commands output and can be empty
	// for library calls.
	RawOutput string `json:"raw_output"`

	// ExecutionTime is the time that task took to complete
	ExecutionTime time.Duration `json:"execution_time"`

	// Mgr holds the manager instance
	Mgr *Manager `json:"-"`
}

// Runnable interface
type Runnable interface {
	Name() string
	Execute() error
	Manager() *Manager
}

// Manager returns the current manager
func (r *Runner) Manager() *Manager {
	return r.Mgr
}

// Name returns the name of the runner
func (r *Runner) Name() string {
	return r.Label
}

// Execute launches the runner
func (r *Runner) Execute() {
}

// trackTime measures time elapsed given the time passed to the func
func (r *Runner) trackTime(start time.Time) {
	r.ExecutionTime = time.Since(start)
}
