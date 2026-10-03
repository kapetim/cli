package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kapetim/cli/src/pkg/manifest"
)

const doc = `# D

<!-- begin table -->
| A | B |
| --- | --- |
| 1 | 2 |
<!-- end table -->
`

func TestScanSaveLoadDiff(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "doc.md"), []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}

	m, err := manifest.Scan(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(m["doc.md"]) != 1 {
		t.Fatalf("tables for doc.md = %d, want 1", len(m["doc.md"]))
	}

	path := filepath.Join(dir, "tables.json")
	if err := manifest.Save(path, m); err != nil {
		t.Fatal(err)
	}
	loaded, err := manifest.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if diffs := manifest.Diff(m, loaded); len(diffs) != 0 {
		t.Fatalf("round-trip diff: %v", diffs)
	}

	// drift: a table was removed
	got, err := manifest.Scan(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "doc.md"), []byte("# D\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	changed, err := manifest.Scan(dir)
	if err != nil {
		t.Fatal(err)
	}
	if diffs := manifest.Diff(changed, got); len(diffs) == 0 {
		t.Fatal("expected diff after removing the table")
	}
}
