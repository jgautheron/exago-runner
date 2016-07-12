package task

import (
	"fmt"
	"os"
	"sync"

	"github.com/davecgh/go-spew/spew"

	. "github.com/exago/runner/config"
)

type manager struct {
	runners []Runnable
}

func NewManager() *manager {
	// Add repository path to configuration
	Config.RepositoryPath = fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), Config.Repository)

	return &manager{
		runners: []Runnable{
			DownloadRunner(),
			TestRunner(),
			CoverageRunner(),
		},
	}
}

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
			// Break at go get step if it can't complete
			if ru.HasError() && ru.Name() == downloadName {
				break
			}
		}
	}
	// Wait for all runners to complete.
	wg.Wait()

	// Loop each runners again
	for _, r := range m.runners {
		spew.Dump(r)
	}

	return nil
}
