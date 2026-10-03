package test

import (
	"strings"
	"testing"

	"github.com/kapetim/cli/src/pkg/process"
)

func TestAvailable(t *testing.T) {
	if process.Available("definitely-not-a-real-binary-xyzzy") {
		t.Fatal("expected unavailable")
	}
}

func TestRun(t *testing.T) {
	out, err := process.Run("sh", []string{"-c", "printf hi"}, process.Opts{})
	if err != nil {
		t.Fatal(err)
	}
	if out.Code != 0 {
		t.Fatalf("code = %d", out.Code)
	}
	if strings.TrimSpace(string(out.Stdout)) != "hi" {
		t.Fatalf("stdout = %q", out.Stdout)
	}
}

func TestRunNonZero(t *testing.T) {
	out, err := process.Run("sh", []string{"-c", "exit 3"}, process.Opts{})
	if err != nil {
		t.Fatal(err)
	}
	if out.Code != 3 {
		t.Fatalf("code = %d, want 3", out.Code)
	}
}
