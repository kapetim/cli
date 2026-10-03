// Package render converts repository markdown to HTML.
//
// Planned implementation: GitHub-flavored markdown via goldmark, with a CSS
// asset for public rendering. This is a stub for now.
package render

import (
	"errors"
	"flag"
)

// Run renders markdown files to HTML. Placeholder until goldmark lands.
func Run(args []string) error {
	flags := flag.NewFlagSet("render", flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return err
	}
	return errors.New("render: not implemented yet (planned: goldmark GFM -> HTML)")
}
