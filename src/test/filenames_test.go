package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kapetim/cli/src/pkg/config"
	"github.com/kapetim/cli/src/pkg/validate"
)

func writeFile(t *testing.T, dir, rel string) {
	t.Helper()
	p := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestFilenamesClean(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Defaults().Filenames
	writeFile(t, dir, "README.md")
	writeFile(t, dir, "src/notes.md")
	writeFile(t, dir, "docker/scripts-go.Dockerfile")
	writeFile(t, dir, "docker/Dockerfile")

	if errs := validate.Filenames(dir, cfg); len(errs) != 0 {
		t.Fatalf("expected clean, got %v", errs)
	}
}

func TestFilenamesViolations(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Defaults().Filenames
	writeFile(t, dir, "src/readme.md") // readme must be uppercase
	writeFile(t, dir, "src/Bad.md")    // uppercase in name
	writeFile(t, dir, "src/weird.exe") // extension not allowed
	writeFile(t, dir, "src/bad.Dockerfile")

	errs := validate.Filenames(dir, cfg)
	if len(errs) < 3 {
		t.Fatalf("expected >=3 errors, got %d: %v", len(errs), errs)
	}
}

func TestCaseCollisions(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Defaults().Filenames
	writeFile(t, dir, "notes.md")
	writeFile(t, dir, "N.md")

	// N.md is uppercase -> filenames flags it, and N.md/notes.md differ only
	// by case at the directory level only if names match; here just assert the
	// collision detector runs over an explicit pair.
	writeFile(t, dir, "dual/one.md")
	writeFile(t, dir, "dual/One.md")
	if errs := validate.Case(dir, cfg); len(errs) == 0 {
		t.Fatal("expected a case collision")
	}
}
