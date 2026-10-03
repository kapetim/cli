package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kapetim/cli/src/pkg/manifest"
	"github.com/kapetim/cli/src/pkg/validate"
)

func TestValidateTables(t *testing.T) {
	dir := t.TempDir()
	docPath := filepath.Join(dir, "doc.md")
	if err := os.WriteFile(docPath, []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Scan(dir)
	if err != nil {
		t.Fatal(err)
	}
	mp := filepath.Join(dir, "tables.json")
	if err := manifest.Save(mp, m); err != nil {
		t.Fatal(err)
	}

	if errs := validate.Tables(dir, mp); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}

	// drift: table removed
	if err := os.WriteFile(docPath, []byte("# D\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if errs := validate.Tables(dir, mp); len(errs) == 0 {
		t.Fatal("expected errors after drift")
	}
}
