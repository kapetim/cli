package validate

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kapetim/cli/src/pkg/config"
)

// lowerName is the permitted shape for every non-excepted name: lowercase
// letters, digits, '.', '_' and '-'.
var lowerName = regexp.MustCompile(`^[a-z0-9._-]+$`)

// Filenames checks naming, casing and extension rules across dir.
func Filenames(dir string, cfg config.Filenames) []error {
	var errs []error
	walk(dir, cfg, func(rel, base string, isDir bool) {
		switch {
		case strings.EqualFold(base, "readme.md"):
			if !readmeOK(base, cfg.Readme) {
				errs = append(errs, fmt.Errorf("%s: readme must be README.md (config: %s)", rel, readmeMode(cfg.Readme)))
			}
		case strings.EqualFold(base, "dockerfile") || strings.HasSuffix(base, ".Dockerfile"):
			if !cfg.Dockerfiles {
				errs = append(errs, fmt.Errorf("%s: dockerfile not allowed", rel))
			} else if !dockerfileOK(base) {
				errs = append(errs, fmt.Errorf("%s: dockerfile must be Dockerfile or <lowercase>.Dockerfile", rel))
			}
		case contains(cfg.Allow, base):
			// explicitly allowed uppercase standard name
		default:
			if !lowerName.MatchString(base) {
				errs = append(errs, fmt.Errorf("%s: only lowercase letters, digits, '.', '_' and '-' are allowed", rel))
			}
		}
		if isDir {
			return
		}
		if ext := strings.TrimPrefix(filepath.Ext(base), "."); ext != "" && !contains(cfg.Extensions, ext) {
			errs = append(errs, fmt.Errorf("%s: extension .%s is not allowed", rel, ext))
		}
	})
	return errs
}

// Case reports paths that collide when lowercased (case-insensitive clashes on
// macOS/Windows, or two files differing only by case).
func Case(dir string, cfg config.Filenames) []error {
	seen := map[string]string{}
	var errs []error
	walk(dir, cfg, func(rel, base string, isDir bool) {
		key := strings.ToLower(rel)
		if prev, ok := seen[key]; ok && prev != rel {
			errs = append(errs, fmt.Errorf("%s: case collision with %s", rel, prev))
			return
		}
		seen[key] = rel
	})
	return errs
}

func walk(dir string, cfg config.Filenames, fn func(rel, base string, isDir bool)) {
	skip := map[string]bool{".git": true, "node_modules": true}
	for _, s := range cfg.Skip {
		skip[s] = true
	}
	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || p == dir {
			return nil
		}
		if d.IsDir() && skip[d.Name()] {
			return fs.SkipDir
		}
		rel, relErr := filepath.Rel(dir, p)
		if relErr != nil {
			return nil
		}
		fn(filepath.ToSlash(rel), d.Name(), d.IsDir())
		return nil
	})
}

func readmeMode(mode string) string {
	if mode == "lower" || mode == "any" {
		return mode
	}
	return "upper"
}

func readmeOK(base, mode string) bool {
	switch readmeMode(mode) {
	case "lower":
		return base == "readme.md"
	case "any":
		return strings.EqualFold(base, "readme.md")
	default:
		return base == "README.md"
	}
}

func dockerfileOK(base string) bool {
	if base == "Dockerfile" {
		return true
	}
	prefix := strings.TrimSuffix(base, ".Dockerfile")
	return prefix != "" && lowerName.MatchString(prefix)
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
