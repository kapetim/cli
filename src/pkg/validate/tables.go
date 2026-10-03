// Package validate runs repository checks.
package validate

import (
	"errors"

	"github.com/kapetim/cli/src/pkg/manifest"
)

// Tables scans dir and diffs the result against the manifest at path.
func Tables(dir, path string) []error {
	want, err := manifest.Load(path)
	if err != nil {
		return []error{err}
	}
	got, err := manifest.Scan(dir)
	if err != nil {
		return []error{err}
	}
	diffs := manifest.Diff(got, want)
	errs := make([]error, 0, len(diffs))
	for _, d := range diffs {
		errs = append(errs, errors.New(d))
	}
	return errs
}
