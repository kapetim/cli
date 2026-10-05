package test

import (
	"testing"

	"github.com/kapetim/cli/src/pkg/version"
)

func TestResolvedExplicit(t *testing.T) {
	orig := version.Version
	defer func() { version.Version = orig }()

	version.Version = "1.2.3"
	if got := version.Resolved(); got != "1.2.3" {
		t.Fatalf("Resolved() = %q, want 1.2.3", got)
	}
}

func TestResolvedDevFallsBack(t *testing.T) {
	orig := version.Version
	defer func() { version.Version = orig }()

	// With the default "dev" and no module version in the test binary, it must
	// not panic and must return something non-empty.
	version.Version = "dev"
	if got := version.Resolved(); got == "" {
		t.Fatal("Resolved() returned empty")
	}
}
