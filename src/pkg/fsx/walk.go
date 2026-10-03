// Package fsx walks repository files with the standard skips.
package fsx

import (
	"io/fs"
	"path/filepath"
	"strings"
)

var skip = map[string]bool{".git": true, "node_modules": true}

// Files returns every file under dir for which pred is true (nil = all).
func Files(dir string, pred func(string) bool) ([]string, error) {
	var out []string
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if skip[d.Name()] {
				return fs.SkipDir
			}
			return nil
		}
		if pred == nil || pred(p) {
			out = append(out, p)
		}
		return nil
	})
	return out, err
}

// MarkdownFiles returns every .md file under dir.
func MarkdownFiles(dir string) ([]string, error) {
	return Files(dir, func(p string) bool { return strings.EqualFold(filepath.Ext(p), ".md") })
}
