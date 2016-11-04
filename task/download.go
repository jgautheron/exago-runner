// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package task

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"time"
)

type downloadRunner struct {
	Runner
}

// DownloadRunner is a runner used for downloading Go projects
// from remote repositories such as Github, Bitbucket etc.
func DownloadRunner(m *Manager) Runnable {
	return &downloadRunner{
		Runner: Runner{Label: "Go Get", Mgr: m},
	}
}

// Execute, downloads a Go repository using the go get command
// too bad, we can't do this as a library :/
func (r *downloadRunner) Execute() {
	defer r.trackTime(time.Now())

	// Return early if repository is already in the GOPATH
	if _, err := os.Stat(r.Manager().RepositoryPath()); err == nil {
		r.toRepoDir()
		return
	}

	// Go get the package
	p := []string{"get", "-d", "-t"}
	if r.Manager().Shallow() {
		p = append(p, "-s")
	}
	rep := r.Manager().Repository()
	if r.Manager().Reference() != "" {
		rep += ":" + r.Manager().Reference()
	}
	p = append(p, rep+"/...")

	out, err := exec.Command("go", p...).CombinedOutput()
	if err != nil {
		// If we can't download, stop execution as BreakOnError is true with this runner
		r.toRunnerError(errors.New(string(out)))
		return
	}

	r.RawOutput = string(out)

	// cd into repository
	r.toRepoDir()
}

func (r *downloadRunner) toRepoDir() {
	// Change directory
	err := os.Chdir(r.Manager().RepositoryPath())
	if err != nil {
		log.Fatal(err)
	}
}
