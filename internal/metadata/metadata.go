// Package metadata provides information about the build, including the application Name,
// version number, git commit hash, and build date. These variables are typically injected at
// build time using LDFLAGS.
package metadata

import (
	"fmt"
)

var (
	AppName   = "unknown"
	Version   = "v0.0.0-dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func String() string {
	return fmt.Sprintf("%s %s (%s) built at %s", AppName, Version, Commit, BuildDate)
}
