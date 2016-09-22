// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package task

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/hotolab/cov"
)

const (
	coverMode = "count"
)

type coverageRunner struct {
	Runner
	tempFile *os.File
}

// CoverageRunner is a runner used for testing Go projects
func CoverageRunner() Runnable {
	return &coverageRunner{
		Runner: Runner{Label: "Code Coverage", parallel: true},
	}
}

// Execute gets all the coverage files and returns the output of
// hotolab/cov
func (r *coverageRunner) Execute() {
	// Create temporary directory to output coverage files
	file, err := ioutil.TempFile("", "exago-coverage")
	if err != nil {
		r.toRunnerError(err)
		return
	}

	// temp file will be removed after processing
	defer os.Remove(file.Name())

	r.tempFile = file

	err = r.lookupTestFiles()
	if err != nil {
		r.toRunnerError(err)
		return
	}

	rep, err := cov.ConvertProfile(r.tempFile.Name())
	if err != nil {
		r.toRunnerError(err)
		return
	}

	r.Data = rep
}

// processPackage executes go test command with coverage and outputs
// errors and output into channels so they are combined later in a single
// file and passed to cov for getting the expected JSON output
func (r *coverageRunner) processPackage(rel string) (string, error) {
	// Create temporary file to output the file coverage
	// this file is trashed after processing
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())

	log.Debugf("go test -covermode=%s -coverprofile=%s %s", coverMode, tmp.Name(), rel)
	_, err = exec.Command("go", "test", "-covermode="+coverMode, "-coverprofile="+tmp.Name(), rel).CombinedOutput()
	if err != nil {
		return "", nil
	}

	// Get file contents
	b, err := ioutil.ReadFile(tmp.Name())
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// lookupTestFiles crawls the filesystem from the repository path
// and finds test files using glob, if a package doesn't have tests
// it is automatically skipped.
func (r *coverageRunner) lookupTestFiles() error {
	pkgs, err := r.packageList()
	if err != nil {
		return err
	}

	// Bufferize channel
	tasks := make(chan string, 64)
	var (
		wg      sync.WaitGroup
		errBuff bytes.Buffer
		outBuff bytes.Buffer
	)

	// Create as much threads as we have CPUs
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			for pkg := range tasks {
				res, err := r.processPackage(pkg)
				if err != nil {
					errBuff.WriteString(err.Error())
					return
				}
				outBuff.WriteString(res)
			}
			wg.Done()
		}()
	}

	for _, pkg := range pkgs {
		tasks <- pkg
	}

	// Close worker channel
	close(tasks)

	// Wait for the workers to finish
	wg.Wait()

	// Get errors (if any) and convert them to a runner error
	errs := errBuff.String()
	if errs != "" {
		return errors.New(errs)
	}

	// Get content of the buffer and write it
	// to the temp file attached to the runner
	out := outBuff.String()
	out = regexp.MustCompile("mode: [a-z]+\n").ReplaceAllString(out, "")
	out = "mode: " + coverMode + "\n" + out

	log.Debug(out)

	if err := ioutil.WriteFile(r.tempFile.Name(), []byte(out), 0644); err != nil {
		return err
	}

	return nil
}

// packageList returns a list of Go-like files or directories from PWD,
func (r *coverageRunner) packageList() ([]string, error) {
	cmd, err := exec.Command("sh", "-c", `go list -f '{{.ImportPath}}' ./... | grep -v vendor | grep -v Godeps`).CombinedOutput()
	if err != nil {
		return nil, err
	}

	pl := strings.Split(string(cmd), "\n")

	log.Debug(pl)

	return pl, nil
}
