// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package task

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/hotolab/cov"

	. "github.com/hotolab/exago-runner/config"
)

const coverMode = "count"

var (
	modeRegex = regexp.MustCompile("mode: [a-z]+\n")
)

type coverageRunner struct {
	Runner
	ignore   []string
	tempFile *os.File
}

// CoverageRunner is a runner used for testing Go projects
func CoverageRunner() Runnable {
	return &coverageRunner{
		Runner: Runner{Label: "Code Coverage", parallel: true},
		ignore: []string{".git", "vendor"},
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

	r.lookupTestFiles()
	if r.Err != nil {
		return
	}

	raw, err := ioutil.ReadAll(r.tempFile)
	if err != nil {
		r.toRunnerError(err)
		return
	}

	rep, err := cov.ConvertProfile(r.tempFile.Name())
	if err != nil {
		r.toRunnerError(err)
		return
	}

	r.RawOutput = string(raw)
	r.Data = rep
}

// processPackage executes go test command with coverage and outputs
// errors and output into channels so they are combined later in a single
// file and passed to cov for getting the expected JSON output
func (r *coverageRunner) processPackage(wg *sync.WaitGroup, fullPath, relPath string, out chan<- string, errs chan<- string) {
	defer wg.Done()

	// Create temporary file to output the file coverage
	// this file is trashed after processing
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		errs <- err.Error()
		return
	}
	defer os.Remove(tmp.Name())

	cmd, err := exec.Command("go", "test", "-covermode="+coverMode, "-coverprofile="+tmp.Name(), relPath).CombinedOutput()
	if err != nil {
		errs <- string(cmd)
		return
	}

	// Get file contents
	b, err := ioutil.ReadFile(tmp.Name())
	if err != nil {
		errs <- err.Error()
		return
	}

	out <- string(b)
}

// lookupTestFiles crawls the filesystem from the repository path
// and finds test files using glob, if a package doesn't have tests
// it is automatically skipped.
func (r *coverageRunner) lookupTestFiles() {
	out := make(chan string)
	errs := make(chan string)

	wg := &sync.WaitGroup{}

	walker := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		rel := ""
		if path != Config.RepositoryPath {
			rel = strings.Replace(path, Config.RepositoryPath+"/", "", 1)
		}
		// Check if path is ignored
		if r.isIgnored(rel) {
			return filepath.SkipDir
		}
		// Rebuild relative path
		rel = "./" + rel
		if files, err := filepath.Glob(rel + "*_test.go"); len(files) == 0 || err != nil {
			if err != nil {
				return err
			}
			// No test file
			log.Debugf("No test files in directory %s, skipping", rel)
			return nil
		}
		// Process package
		wg.Add(1)
		go r.processPackage(wg, path, rel, out, errs)

		return nil
	}

	// Start the crawler
	if err := filepath.Walk(Config.RepositoryPath, walker); err != nil {
		r.toRunnerError(err)
		return
	}

	// Wait for all routines to complete
	go func() {
		wg.Wait()
		close(out)
		close(errs)
	}()

	select {
	// Get errors (if any) and convert them to a runner error
	case err, ok := <-errs:
		if ok {
			r.toRunnerError(errors.New(err))
			return
		}
	// Get content of the output channel and write it to the temp file
	// attached to the runner
	case buff, ok := <-out:
		if ok {
			buff = modeRegex.ReplaceAllString(buff, "")
			buff = "mode: " + coverMode + "\n" + buff

			if err := ioutil.WriteFile(r.tempFile.Name(), []byte(buff), 0644); err != nil {
				r.toRunnerError(err)
				return
			}
		}
	}
}

// isIgnored checks if a path is ignored
func (r *coverageRunner) isIgnored(path string) bool {
	for _, i := range r.ignore {
		if i == path {
			return true
		}
	}
	return false
}
