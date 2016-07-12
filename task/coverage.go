package task

import (
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type coverageRunner struct {
	Runner
}

// CoverageRunner is a runner used for testing Go projects
func CoverageRunner() Runnable {
	return &coverageRunner{Runner{Label: coverageName, Parallel: true}}
}

// Execute ...
func (r *coverageRunner) Execute() {
	defer r.trackTime(time.Now())

	ftmp, err := ioutil.TempFile(os.TempDir(), "")
	//defer os.Remove(ftmp.Name())
	if err != nil {
		r.Error = &RunnerError{
			RawOutput: err.Error(),
			Message:   err,
		}
	}

	fcov, err := ioutil.TempFile(os.TempDir(), "")
	//defer os.Remove(fcov.Name())
	if err != nil {
		r.Error = &RunnerError{
			RawOutput: err.Error(),
			Message:   err,
		}
	}

	out, err := exec.Command(
		"echo 'mode: count' > "+fcov.Name(),
		"&&",
		"go list -f '{{.ImportPath}}' ./... | grep -v vendor",
		"| xargs -n1 -I{}",
		"sh -c 'go test -covermode=count -coverprofile="+ftmp.Name()+" {}",
		"&&",
		"tail -n +2 "+ftmp.Name()+" >> "+fcov.Name(),
		"test",
		"-v",
		"./...",
	).CombinedOutput()

	if err != nil {
		r.Error = &RunnerError{
			RawOutput: string(out),
			Message:   err,
		}
		return
	}

	cov, err := ioutil.ReadAll(fcov)
	if err != nil {
		r.Error = &RunnerError{
			RawOutput: err.Error(),
			Message:   err,
		}
		return
	}

	r.RawOutput = string(cov)
}
