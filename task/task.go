// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package task

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	. "github.com/hotolab/exago-runner/config"
)

// Manager contains all registered runnables
type Manager struct {
	runners map[string]Runnable
}

// NewManager instantiates a runnable manager
// the manager has the responsibility to execute all runners
// and decide whether a runner should run in parallel processing or not
func NewManager() *Manager {
	// Add repository path to configuration
	Config.RepositoryPath = fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), Config.Repository)

	return &Manager{
		runners: map[string]Runnable{
			downloadName:     DownloadRunner(),
			testName:         TestRunner(),
			coverageName:     CoverageRunner(),
			proveName:        ProveRunner(),
			thirdPartiesName: ThirdPartiesRunner(),
		},
	}
}

// ExecuteRunners launches the runners
func (m *Manager) ExecuteRunners() error {
	// Execute download runner synchronously
	dlr := m.runners[downloadName]
	// Execute synchronously
	dlr.Execute()
	// Exit early if we can't download
	if dlr.HasError() {
		return nil
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
