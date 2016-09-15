// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package task

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

const (
	downloadName     = "download"
	testName         = "test"
	coverageName     = "coverage"
	proveName        = "goprove"
	thirdPartiesName = "thirdparties"
)

// RunnerError is the struct containing processing errors
type RunnerError struct {
	RawOutput string
	Message   error
}

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
	Err *RunnerError `json:"error"`

	// Whether runner should execute in parallel or not
	parallel bool

	// breakOnError will stop execution chain if true
	breakOnError bool
}

// Runnable
type Runnable interface {
	Name() string
	Execute()
	CanParallelize() bool
	HasError() bool
	BreakOnError() bool
}

// Name returns the name of the runner
func (r *Runner) Name() string {
	return r.Label
}

// Execute launches the runner
func (r *Runner) Execute() {
}

// CanParallelize tells if runner can be parallelized
func (r *Runner) CanParallelize() bool {
	return r.parallel
}

// HasError tells if runner had some errors during processing
func (r *Runner) HasError() bool {
	return r.Err != nil
}

// BreakOnError tells if runner should break on error
func (r *Runner) BreakOnError() bool {
	return r.breakOnError
}

// toRunnerError converts a golang error to a Runner error
func (r *Runner) toRunnerError(err error) {
	log.Println(err)
	r.Err = &RunnerError{
		RawOutput: err.Error(),
		Message:   err,
	}
}

// trackTime measures time elapsed given the time passed to the func
func (r *Runner) trackTime(start time.Time) {
	r.ExecutionTime = time.Since(start)
}
