// Package process runs external tools (linters) and reports their output.
package process

import (
	"bytes"
	"os/exec"
)

// Opts configures a command run.
type Opts struct {
	Dir string
	Env []string
}

// Out is the captured result of a run.
type Out struct {
	Stdout []byte
	Stderr []byte
	Code   int
}

// Available reports whether a binary is on PATH.
func Available(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// Run executes name with args. A non-zero exit is returned in Out.Code (no Go
// error); only a failure to start returns an error.
func Run(name string, args []string, opts Opts) (Out, error) {
	cmd := exec.Command(name, args...)
	if opts.Dir != "" {
		cmd.Dir = opts.Dir
	}
	if opts.Env != nil {
		cmd.Env = opts.Env
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	out := Out{Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}
	if err == nil {
		return out, nil
	}
	if ee, ok := err.(*exec.ExitError); ok {
		out.Code = ee.ExitCode()
		return out, nil
	}
	return out, err
}
