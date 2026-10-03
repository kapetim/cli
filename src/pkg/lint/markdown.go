package lint

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kapetim/cli/src/pkg/config"
	"github.com/kapetim/cli/src/pkg/fsx"
	"github.com/kapetim/cli/src/pkg/table"
)

// lintMarkdown applies the built-in markdown rules to every .md file:
// balanced table markers and rendered table-cell width.
func lintMarkdown(dir string, cfg config.Config) []error {
	files, err := fsx.MarkdownFiles(dir)
	if err != nil {
		return []error{err}
	}
	var errs []error
	for _, f := range files {
		rel, err := filepath.Rel(dir, f)
		if err != nil {
			rel = f
		}
		errs = append(errs, lintMarkdownFile(filepath.ToSlash(rel), f, cfg.Markdown.MaxCell)...)
	}
	return errs
}

func lintMarkdownFile(rel, path string, maxCell int) []error {
	data, err := os.ReadFile(path)
	if err != nil {
		return []error{err}
	}
	var errs []error
	in := false
	lines := strings.Split(string(data), "\n")
	for i, raw := range lines {
		s := strings.TrimSpace(raw)
		switch s {
		case table.BeginMarker:
			if in {
				errs = append(errs, Issue{rel, i + 1, "table", "nested begin marker"})
			}
			in = true
		case table.EndMarker:
			if !in {
				errs = append(errs, Issue{rel, i + 1, "table", "end without begin marker"})
			}
			in = false
			continue
		}
		if !in || s == "" || table.IsDelimiter(s) {
			continue
		}
		for _, cell := range table.SplitRow(s) {
			if n := table.RenderedLength(cell); n > maxCell {
				errs = append(errs, Issue{rel, i + 1, "table-cell",
					"cell renders to " + strconv.Itoa(n) + " chars (max " + strconv.Itoa(maxCell) + ")"})
			}
		}
	}
	if in {
		errs = append(errs, Issue{rel, len(lines), "table", "unclosed begin marker"})
	}
	return errs
}
