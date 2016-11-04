// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package task

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
)

// Manager contains all registered runnables
type Manager struct {
	runners        map[string]Runnable
	repository     string
	repositoryPath string
	shallow        bool
	reference      string
}

// NewManager instantiates a runnable manager
// the manager has the responsibility to execute all runners
// and decide whether a runner should run in parallel processing or not
func NewManager(r string) *Manager {
	if strings.TrimSpace(r) == "" {
		log.Fatal("Repository is required")
	}

	m := &Manager{
		repository:     r,
		repositoryPath: fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), r),
	}

	m.runners = map[string]Runnable{
		downloadName:     DownloadRunner(m),
		testName:         TestRunner(m),
		coverageName:     CoverageRunner(m),
		proveName:        ProveRunner(m),
		thirdPartiesName: ThirdPartiesRunner(m),
	}

	return m
}

// DoShallow sets shallow flag to true
func (m *Manager) DoShallow() {
	m.shallow = true
}

// Shallow returns shallow flag
func (m *Manager) Shallow() bool {
	return m.shallow
}

// UseReference sets reference flag
func (m *Manager) UseReference(r string) {
	m.reference = r
}

// Reference returns reference
func (m *Manager) Reference() string {
	return m.reference
}

// RepositoryPath returns repository path
func (m *Manager) RepositoryPath() string {
	return m.repositoryPath
}

// Repository returns repository (e.g. :vcs/:owner/:package+)
func (m *Manager) Repository() string {
	return m.repository
}

// ExecuteRunners launches the runners
func (m *Manager) ExecuteRunners() error {
	// Execute download runner synchronously
	dlr := m.runners[downloadName]
	// Execute synchronously
	dlr.Execute()
	// Exit early if we can't download
	if dlr.HasError() {
		return dlr.Error()
	}

	var wg sync.WaitGroup
	for n, ru := range m.runners {
		// Skip download runner
		if n == downloadName {
			continue
		}
		// Increment the WaitGroup counter.
		wg.Add(1)
		go func(r Runnable) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			// Execute the runner
			r.Execute()
		}(ru)
	}
	// Wait for all runners to complete.
	wg.Wait()

	// And printout JSON
	printOutput(m.runners)

	return nil
}

func printOutput(o interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(o)
}
