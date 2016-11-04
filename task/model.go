// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package task

import (
	"errors"
	"time"
)

const (
	downloadName     = "download"
	testName         = "test"
	coverageName     = "coverage"
	proveName        = "goprove"
	thirdPartiesName = "thirdparties"
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

	// Error returns details about the error
	Err string `json:"error"`

	// Mgr holds the manager instance
	Mgr *Manager `json:"-"`
}

// Runnable
type Runnable interface {
	Name() string
	Execute()
	HasError() bool
	Error() error
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

// HasError tells if runner had some errors during processing
func (r *Runner) HasError() bool {
	return r.Err != ""
}

// Error returns a golang error
func (r *Runner) Error() error {
	if !r.HasError() {
		return nil
	}
	return errors.New(r.Err)
}

// toRunnerError converts a golang error to a Runner error
func (r *Runner) toRunnerError(err error) {
	r.Err = err.Error()
}

// trackTime measures time elapsed given the time passed to the func
func (r *Runner) trackTime(start time.Time) {
	r.ExecutionTime = time.Since(start)
}
