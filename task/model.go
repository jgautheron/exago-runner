package task

import "time"

const (
	downloadName = "download"
	testName     = "test"
	coverageName = "coverage"
)

type RunnerError struct {
	RawOutput string
	Message   error
}

// Runner
type Runner struct {
	// Label is the name of the task runner
	// This is the only field that must be set
	Label string

	// Data holds the specialized object associated to the task
	// runner i.e. specialized object for Goprove and Gotest
	Data interface{}

	// RawOutput is the process's standard output and error.
	// It is used for system commands output and can be empty
	// for library calls.
	RawOutput string

	// ExecutionTime is the time that task took to complete
	ExecutionTime time.Duration

	// Error returns details about the error
	Error *RunnerError

	// Whether runner should execute in parallel or not
	Parallel bool
}

// Runnable
type Runnable interface {
	Name() string
	Execute()
	CanParallelize() bool
	HasError() bool
}

func (r *Runner) Name() string {
	return r.Label
}

func (r *Runner) Execute() {
}

func (r *Runner) CanParallelize() bool {
	return r.Parallel
}

func (r *Runner) HasError() bool {
	return r.Error != nil
}

func (r *Runner) trackTime(start time.Time) {
	r.ExecutionTime = time.Since(start)
}
