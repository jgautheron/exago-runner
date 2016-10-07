// Copyright (c) 2016, Hotolab. All rights reserved.
//
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package task

import (
	"os/exec"
	"regexp"
	"strings"
	"time"

	. "github.com/hotolab/exago-runner/config"
)

type thirdPartiesRunner struct {
	Runner
}

// ThirdPartiesRunner launches go list to find all dependencies
func ThirdPartiesRunner() Runnable {
	return &thirdPartiesRunner{Runner{Label: "Go List (finds all 3rd parties)"}}
}

// Execute go list
func (r *thirdPartiesRunner) Execute() {
	defer r.trackTime(time.Now())

	list, err := exec.Command("go", "list", "-f", `'{{ join .Deps ", " }}'`, "./...").CombinedOutput()
	if err != nil {
		r.toRunnerError(err)
	}

	r.Data = r.parseListOutput(string(list))
}

func (r *thirdPartiesRunner) parseListOutput(output string) (out []string) {
	reg := regexp.MustCompile(`(?m)([\w\d\-]+)\.([\w]{2,})\/([\w\d\-]+)\/([\w\d\-\.]+)(\.v\d+)?`)
	// See https://danott.co/posts/json-marshalling-empty-slices-to-empty-arrays-in-go.html
	out = make([]string, 0)
	uniq := map[string]bool{}
	sl := strings.Split(output, ",")
	for _, v := range sl {
		v = strings.TrimSpace(v)
		m := reg.FindAllString(v, -1)

		// Only interested in third parties
		if len(m) == 0 {
			continue
		}

		// Match only the last path found in the path
		// That way we support imports made this way:
		// github.com/heroku/hk/Godeps/_workspace/src/code.google.com/p/go-uuid/uuid
		lastMatch := m[len(m)-1]
		if lastMatch == Config.Repository {
			continue
		}

		uniq[lastMatch] = true
	}

	for pkg := range uniq {
		out = append(out, pkg)
	}

	return out
}
