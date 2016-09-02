// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package task

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/hotolab/cov"
)

type coverageRunner struct {
	Runner
}

// CoverageRunner is a runner used for testing Go projects
func CoverageRunner() Runnable {
	return &coverageRunner{Runner{Label: "Code Coverage", parallel: true}}
}

// Execute Measures the code coverage accurately using cov
func (r *coverageRunner) Execute() {
	defer r.trackTime(time.Now())

	ftmp, err := ioutil.TempFile(os.TempDir(), "")
	defer os.Remove(ftmp.Name())

	if err != nil {
		r.toRunnerError(err)
	}

	fcov, err := ioutil.TempFile(os.TempDir(), "")
	defer os.Remove(fcov.Name())

	if err != nil {
		r.toRunnerError(err)
	}

	cmd := fmt.Sprintf(
		"echo 'mode: count' > %s && go list -f '{{.ImportPath}}' ./... | grep -v vendor | xargs -n1 -I{} go test -covermode=count -coverprofile=%s {} && tail -n +2 %s >> %s",
		fcov.Name(),
		ftmp.Name(),
		ftmp.Name(),
		fcov.Name(),
	)

	_, err = exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		r.toRunnerError(err)
		return
	}

	raw, err := ioutil.ReadAll(fcov)
	if err != nil {
		r.toRunnerError(err)
		return
	}

	rep, err := cov.ConvertProfile(fcov.Name())
	if err != nil {
		r.toRunnerError(err)
		return
	}

	r.RawOutput = string(raw)
	r.Data = rep
}
