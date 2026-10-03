// Package manifest reads, writes, scans and diffs the table manifest
// (`tables.json`) that records every marker-delimited table and its location.
package manifest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/kapetim/cli/src/pkg/fsx"
	"github.com/kapetim/cli/src/pkg/table"
)

// Table is one table's expected location and shape.
type Table struct {
	Begin int    `json:"begin"`
	End   int    `json:"end"`
	Lines string `json:"lines"`
	Cols  int    `json:"cols"`
	Rows  int    `json:"rows"`
}

// Manifest maps a repo-relative file path to its tables.
type Manifest map[string][]Table

// Scan builds the manifest from the current state of dir.
func Scan(dir string) (Manifest, error) {
	files, err := fsx.MarkdownFiles(dir)
	if err != nil {
		return nil, err
	}
	m := Manifest{}
	for _, f := range files {
		locs, err := table.Scan(f)
		if err != nil {
			return nil, err
		}
		if len(locs) == 0 {
			continue
		}
		rel, err := filepath.Rel(dir, f)
		if err != nil {
			return nil, err
		}
		ts := make([]Table, 0, len(locs))
		for _, l := range locs {
			ts = append(ts, Table{Begin: l.Begin, End: l.End, Lines: l.Lines, Cols: l.Cols, Rows: l.Rows})
		}
		m[filepath.ToSlash(rel)] = ts
	}
	return m, nil
}

// Load reads a manifest; a missing file yields an empty manifest.
func Load(path string) (Manifest, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Manifest{}, nil
		}
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	if m == nil {
		m = Manifest{}
	}
	return m, nil
}

// Save writes the manifest (sorted, stable).
func Save(path string, m Manifest) error {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0o644)
}

// Diff reports differences between the scanned manifest (got) and the expected
// one (want).
func Diff(got, want Manifest) []string {
	keys := map[string]bool{}
	for k := range got {
		keys[k] = true
	}
	for k := range want {
		keys[k] = true
	}
	names := make([]string, 0, len(keys))
	for k := range keys {
		names = append(names, k)
	}
	sort.Strings(names)

	var out []string
	for _, k := range names {
		g, w := got[k], want[k]
		if len(g) == 0 {
			out = append(out, k+": table missing")
			continue
		}
		if len(w) == 0 {
			out = append(out, k+": unexpected table")
			continue
		}
		if len(g) != len(w) {
			out = append(out, k+": table count "+strconv.Itoa(len(g))+" != "+strconv.Itoa(len(w)))
			continue
		}
		for i := range w {
			if g[i].Lines != w[i].Lines {
				out = append(out, k+": table "+strconv.Itoa(i+1)+" at "+g[i].Lines+" != "+w[i].Lines)
			}
			if g[i].Cols != w[i].Cols {
				out = append(out, k+": table "+strconv.Itoa(i+1)+" cols "+strconv.Itoa(g[i].Cols)+" != "+strconv.Itoa(w[i].Cols))
			}
		}
	}
	return out
}
