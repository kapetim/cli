// Package version holds the CLI version.
package version

// Version is the reported CLI version. It is overridden at build time with
// -ldflags "-X github.com/kapetim/cli/src/pkg/version.Version=<tag>".
var Version = "0.1.0"
