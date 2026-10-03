// Package logx is the shared output + exit-code contract.
package logx

import (
	"fmt"
	"os"
)

// Exit codes are the contract CI relies on.
const (
	ExitOK         = 0
	ExitError      = 1
	ExitValidation = 2
	ExitNetwork    = 3
)

// OK reports success (stdout).
func OK(format string, args ...any) { fmt.Printf("[OK] "+format+"\n", args...) }

// Skip reports a skipped step (stdlib out).
func Skip(format string, args ...any) { fmt.Printf("[skip] "+format+"\n", args...) }

// Info is a neutral note.
func Info(format string, args ...any) { fmt.Printf("[cli] "+format+"\n", args...) }

// Fail reports a failure (stderr).
func Fail(format string, args ...any) { fmt.Fprintf(os.Stderr, "[FAIL] "+format+"\n", args...) }
