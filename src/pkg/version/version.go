// Package version holds the CLI version.
package version

import (
	"runtime/debug"
	"strings"
)

// Version is the reported CLI version. Released builds override it at build
// time with -ldflags "-X github.com/kapetim/cli/src/pkg/version.Version=<tag>".
// It defaults to "dev".
var Version = "dev"

// Resolved returns the version to report: the ldflags-injected Version when
// set, otherwise the module version recorded in the build info (so consumers
// that import the cli module report the real version too).
func Resolved() string {
	if Version != "" && Version != "dev" {
		return Version
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return Version
	}
	if v := moduleVersion(info.Main.Path, info.Main.Version); v != "" {
		return v
	}
	for _, d := range info.Deps {
		if v := moduleVersion(d.Path, d.Version); v != "" {
			return v
		}
	}
	return Version
}

func moduleVersion(path, v string) string {
	if path != "github.com/kapetim/cli" {
		return ""
	}
	if v == "" || v == "(devel)" {
		return ""
	}
	return strings.TrimPrefix(v, "v")
}
