// Package metadata provides information about the build, including the application name,
// version number and git commit hash. These values are injected at build time using LDFLAGS.
package metadata

import "fmt"

var (
	// AppName is the name of the application binary.
	AppName = "unknown"
	// Version is the semantic version of the build, e.g. v1.2.3.
	Version = "v0.0.0-dev"
	// Commit is the git commit the binary was built from.
	Commit = "none"
)

// String returns a human-readable summary of the build metadata.
func String() string {
	return fmt.Sprintf("%s %s (%s)", AppName, Version, Commit)
}
