// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package task

import (
	"os"
	"os/exec"
	"time"

	. "github.com/hotolab/exago-runner/config"
)

type downloadRunner struct {
	Runner
}

// DownloadRunner is a runner used for downloading Go projects
// from remote repositories such as Github, Bitbucket etc.
func DownloadRunner() Runnable {
	return &downloadRunner{Runner{Label: "Go Get", breakOnError: true}}
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

	// Go get the package
	out, err := exec.Command("go", "get", "-d", "-t", Config.Repository+"/...").CombinedOutput()
	if err != nil {
		// If we can't download, stop execution as BreakOnError is true with this runner
		r.toRunnerError(err)
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
		r.toRunnerError(err)
	}
}
