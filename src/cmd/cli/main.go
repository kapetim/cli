// Command cli is the kapetim repository toolchain.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/kapetim/cli/src/internal/render"
	"github.com/kapetim/cli/src/internal/validate"
	"github.com/kapetim/cli/src/internal/version"
)

const usage = `cli — kapetim repository toolchain

Usage:
  cli <command> [flags] [paths...]

Commands:
  validate   validate markdown and repository files
  render     render markdown to HTML (planned)
  version    print the version
  help       show this help

Run "cli <command> --help" for command flags.`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
	cmd, args := os.Args[1], os.Args[2:]

	var err error
	switch cmd {
	case "validate":
		err = validate.Run(args)
	case "render":
		err = render.Run(args)
	case "version", "--version", "-v":
		fmt.Println("cli " + version.Version)
	case "help", "-h", "--help":
		fmt.Println(usage)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s\n", cmd, usage)
		os.Exit(2)
	}
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, "cli: "+err.Error())
		os.Exit(1)
	}
}
