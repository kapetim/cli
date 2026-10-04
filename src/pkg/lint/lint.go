// Package lint runs the repo's linters: a built-in markdown linter plus calls
// to native tools (shellcheck, hadolint, actionlint).
package lint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kapetim/cli/src/pkg/config"
	"github.com/kapetim/cli/src/pkg/fsx"
	"github.com/kapetim/cli/src/pkg/logx"
	"github.com/kapetim/cli/src/pkg/process"
)

// Issue is a single lint finding.
type Issue struct {
	File    string
	Line    int
	Rule    string
	Message string
}

func (i Issue) Error() string {
	if i.File == "" {
		return fmt.Sprintf("%s: %s", i.Rule, i.Message)
	}
	if i.Line == 0 {
		return fmt.Sprintf("%s: %s: %s", i.File, i.Rule, i.Message)
	}
	return fmt.Sprintf("%s:%d: %s: %s", i.File, i.Line, i.Rule, i.Message)
}

// Run executes the requested linters (defaults from cfg) over dir.
func Run(dir string, kinds []string, cfg config.Config) []error {
	if len(kinds) == 0 {
		kinds = cfg.Linters
	}
	for _, k := range kinds {
		if k == "all" {
			kinds = cfg.Linters
			break
		}
	}
	var errs []error
	for _, k := range kinds {
		switch k {
		case "markdown":
			errs = append(errs, lintMarkdown(dir, cfg)...)
		case "shell":
			errs = append(errs, lintShell(dir)...)
		case "docker":
			errs = append(errs, lintDocker(dir)...)
		case "ci":
			errs = append(errs, lintCI(dir)...)
		default:
			errs = append(errs, Issue{Rule: "cli", Message: "unknown linter: " + k})
		}
	}
	return errs
}

func lintShell(dir string) []error {
	if !process.Available("shellcheck") {
		logx.Skip("shellcheck not installed")
		return nil
	}
	files, err := fsx.Files(dir, func(p string) bool { return strings.HasSuffix(p, ".sh") })
	if err != nil {
		return []error{err}
	}
	if len(files) == 0 {
		return nil
	}
	return run(dir, "shellcheck", append([]string{"-x"}, files...))
}

func lintDocker(dir string) []error {
	if !process.Available("hadolint") {
		logx.Skip("hadolint not installed")
		return nil
	}
	files, err := fsx.Files(dir, func(p string) bool {
		base := filepath.Base(p)
		return base == "Dockerfile" || strings.HasSuffix(base, ".Dockerfile")
	})
	if err != nil {
		return []error{err}
	}
	if len(files) == 0 {
		return nil
	}
	return run(dir, "hadolint", files)
}

func lintCI(dir string) []error {
	if !process.Available("actionlint") {
		logx.Skip("actionlint not installed")
		return nil
	}
	if _, err := os.Stat(filepath.Join(dir, ".github", "workflows")); err != nil {
		return nil
	}
	return run(dir, "actionlint", nil)
}

func run(dir, bin string, args []string) []error {
	out, err := process.Run(bin, args, process.Opts{Dir: dir})
	if err != nil {
		return []error{err}
	}
	if out.Code != 0 {
		msg := strings.TrimSpace(string(out.Stdout) + "\n" + string(out.Stderr))
		return []error{Issue{Rule: bin, Message: msg}}
	}
	return nil
}
