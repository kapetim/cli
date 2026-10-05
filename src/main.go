// Command cli is the kapetim repository toolchain.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kapetim/cli/src/pkg/logx"
	"github.com/kapetim/cli/src/pkg/manifest"
	"github.com/kapetim/cli/src/pkg/run"
	"github.com/kapetim/cli/src/pkg/version"
)

const usage = `cli — kapetim repository toolchain

Usage:
  cli <command> [flags] [kinds...]

Commands:
  lint [kinds...]        run linters (markdown, shell, docker, ci)
  validate [kinds...]    run checks (tables, filenames, case)
  scan tables            write the table manifest
  version                print the version
  help                   show this help

Flags:
  --dir <path>           repository root (default .)
  --manifest <path>      table manifest, relative to --dir (default tables.json)`

type flags struct {
	dir      string
	manifest string
	rest     []string
}

func parseFlags(args []string) flags {
	f := flags{dir: ".", manifest: "tables.json"}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dir":
			if i+1 < len(args) {
				f.dir = args[i+1]
				i++
			}
		case "--manifest":
			if i+1 < len(args) {
				f.manifest = args[i+1]
				i++
			}
		default:
			f.rest = append(f.rest, args[i])
		}
	}
	return f
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(logx.ExitValidation)
	}
	cmd, args := os.Args[1], os.Args[2:]

	switch cmd {
	case "version", "--version", "-v":
		fmt.Println("cli " + version.Resolved())
	case "help", "-h", "--help":
		fmt.Println(usage)
	case "lint":
		f := parseFlags(args)
		if len(f.rest) == 0 {
			f.rest = []string{"all"}
		}
		os.Exit(run.Main(run.Options{RepoDir: f.dir, Manifest: f.manifest, Lint: f.rest}))
	case "validate":
		f := parseFlags(args)
		if len(f.rest) == 0 {
			f.rest = []string{"tables"}
		}
		os.Exit(run.Main(run.Options{RepoDir: f.dir, Manifest: f.manifest, Validate: f.rest}))
	case "scan":
		os.Exit(scan(parseFlags(args)))
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s\n", cmd, usage)
		os.Exit(logx.ExitValidation)
	}
}

func scan(f flags) int {
	if len(f.rest) == 0 || f.rest[0] != "tables" {
		logx.Fail("usage: cli scan tables [--dir <path>] [--manifest <path>]")
		return logx.ExitValidation
	}
	m, err := manifest.Scan(f.dir)
	if err != nil {
		logx.Fail("%v", err)
		return logx.ExitError
	}
	out := f.manifest
	if !filepath.IsAbs(out) {
		out = filepath.Join(f.dir, out)
	}
	if err := manifest.Save(out, m); err != nil {
		logx.Fail("%v", err)
		return logx.ExitError
	}
	logx.OK("wrote %s (%d file(s))", out, len(m))
	return logx.ExitOK
}
