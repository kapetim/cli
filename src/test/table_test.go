package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kapetim/cli/src/pkg/table"
)

const sample = `# T

<!-- begin table -->
| A | B |
| --- | --- |
| 1 | 2 |
<!-- end table -->
`

func TestParse(t *testing.T) {
	tables := table.Parse(sample)
	if len(tables) != 1 {
		t.Fatalf("tables = %d, want 1", len(tables))
	}
	if len(tables[0]) != 2 {
		t.Fatalf("rows = %d, want 2", len(tables[0]))
	}
}

func TestScan(t *testing.T) {
	p := filepath.Join(t.TempDir(), "a.md")
	if err := os.WriteFile(p, []byte(sample), 0o644); err != nil {
		t.Fatal(err)
	}
	locs, err := table.Scan(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(locs) != 1 {
		t.Fatalf("located = %d, want 1", len(locs))
	}
	if locs[0].Lines != "L3-L7" {
		t.Errorf("lines = %s, want L3-L7", locs[0].Lines)
	}
	if locs[0].Cols != 2 || locs[0].Rows != 2 {
		t.Errorf("cols=%d rows=%d, want 2/2", locs[0].Cols, locs[0].Rows)
	}
}

func TestRenderedLength(t *testing.T) {
	if n := table.RenderedLength("**bold**"); n != 4 {
		t.Errorf("bold = %d, want 4", n)
	}
	if n := table.RenderedLength("[text](https://example.com)"); n != 4 {
		t.Errorf("link = %d, want 4", n)
	}
}

func TestNumbers(t *testing.T) {
	if v, ok := table.BRL("**R$ 1.234,50**"); !ok || v != 1234.50 {
		t.Errorf("BRL = %v/%v", v, ok)
	}
	if v, ok := table.Percent("103%"); !ok || v != 103 {
		t.Errorf("Percent = %v/%v", v, ok)
	}
	if !table.Approx(1.0, 1.4, 0.5) {
		t.Error("Approx should be true")
	}
}
