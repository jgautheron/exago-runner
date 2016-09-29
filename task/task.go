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

// manager contains all registered runnables
type manager struct {
	runners map[string]Runnable
}

// NewManager instantiates a runnable manager
// the manager has the responsibility to execute all runners
// and decide whether a runner should run in parallel processing or not
func NewManager() *manager {
	// Add repository path to configuration
	Config.RepositoryPath = fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), Config.Repository)

	return &manager{
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
func (m *manager) ExecuteRunners() error {
	var wg sync.WaitGroup
	for _, ru := range m.runners {
		if ru.CanParallelize() {
			// Increment the WaitGroup counter.
			wg.Add(1)
			go func(r Runnable) {
				// Decrement the counter when the goroutine completes.
				defer wg.Done()
				// Execute the runner
				r.Execute()
			}(ru)
		} else {
			// Execute synchronously
			ru.Execute()
			// Break if runner has error and should break on error
			// this only applies to non parallel runners
			if ru.HasError() && ru.BreakOnError() {
				break
			}
		}
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
