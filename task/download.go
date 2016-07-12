package task

import (
	"os"
	"os/exec"
	"time"

	. "github.com/exago/runner/config"
)

type downloadRunner struct {
	Runner
}

// DownloadRunner is a runner used for downloading Go projects
// from remote repositories such as Github, Bitbucket etc.
func DownloadRunner() Runnable {
	return &downloadRunner{Runner{Label: downloadName}}
}

// Execute, downloads a Go repository using the go get command
// too bad, we can't do this as a library :/
func (r *downloadRunner) Execute() {
	defer r.trackTime(time.Now())

	// Return early if repository is already in the GOPATH
	if _, err := os.Stat(Config.RepositoryPath); err == nil {
		r.toRepoDir()
		return
	}

	out, err := exec.Command("go", "get", "-d", "-t", Config.Repository+"/...").CombinedOutput()
	if err != nil {
		r.Error = &RunnerError{
			RawOutput: string(out),
			Message:   err,
		}
		return
	}

	r.RawOutput = string(out)

	// cd into repository
	r.toRepoDir()
}

func (r *downloadRunner) toRepoDir() {
	// Change directory
	err := os.Chdir(Config.RepositoryPath)
	if err != nil {
		r.Error = &RunnerError{
			RawOutput: err.Error(),
			Message:   err,
		}
	}
}
