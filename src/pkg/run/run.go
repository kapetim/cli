// Package run is the umbrella entrypoint: it wires config, linting, validation,
// logging and exit codes so a consumer's main is a few lines.
package run

import (
	"path/filepath"

	"github.com/kapetim/cli/src/pkg/config"
	"github.com/kapetim/cli/src/pkg/lint"
	"github.com/kapetim/cli/src/pkg/logx"
	"github.com/kapetim/cli/src/pkg/validate"
)

// Options selects what to run.
type Options struct {
	RepoDir  string
	Lint     []string
	Validate []string
	Manifest string
}

// Main runs the requested work and returns an exit code.
func Main(o Options) int {
	if o.RepoDir == "" {
		o.RepoDir = "."
	}
	if o.Manifest == "" {
		o.Manifest = "tables.json"
	}

	cfg, err := config.Load(o.RepoDir)
	if err != nil {
		logx.Fail("%v", err)
		return logx.ExitError
	}

	if len(o.Lint) > 0 {
		if errs := lint.Run(o.RepoDir, o.Lint, cfg); len(errs) > 0 {
			for _, e := range errs {
				logx.Fail("%v", e)
			}
			return logx.ExitValidation
		}
		logx.OK("lint clean")
	}

	for _, kind := range o.Validate {
		switch kind {
		case "tables", "all":
			path := resolve(o.RepoDir, o.Manifest)
			if errs := validate.Tables(o.RepoDir, path); len(errs) > 0 {
				for _, e := range errs {
					logx.Fail("%v", e)
				}
				return logx.ExitValidation
			}
			logx.OK("tables match %s", o.Manifest)
		default:
			logx.Fail("unknown validation: %s", kind)
			return logx.ExitValidation
		}
	}

	return logx.ExitOK
}

func resolve(dir, p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(dir, p)
}
