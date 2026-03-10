# Golang Project Template

[![Checks][badge-checks]][actions]
[![Latest Release][badge-release-date]][releases]
[![Codecov][badge-codecov]][codecov]
[![Go Report Card][badge-go-report-card]][go-report-card]
[![License][badge-license]][license]

[badge-checks]: https://img.shields.io/github/check-runs/ComixHe/go-template/main
[actions]: https://github.com/ComixHe/go-template/actions

[badge-release-date]: https://img.shields.io/github/release-date/ComixHe/go-template
[releases]: https://github.com/ComixHe/go-template/releases

[badge-codecov]: https://codecov.io/github/ComixHe/go-template/graph/badge.svg?token=YOUR_TOKEN_HERE
[codecov]: https://codecov.io/github/ComixHe/go-template

[badge-go-report-card]: https://goreportcard.com/badge/github.com/ComixHe/go-template
[go-report-card]: https://goreportcard.com/report/github.com/ComixHe/go-template

[badge-license]: https://img.shields.io/github/license/ComixHe/go-template
[license]: https://github.com/ComixHe/go-template/blob/main/LICENSES/Apache-2.0.txt

A modern, production-ready Go project template with comprehensive tooling, best practices, and automation for building high-quality applications.

## Features

- **Modern Go Project Structure** - Well-organized directory layout following Go best practices
- **Complete Development Toolchain** - Pre-configured tools for building, testing, and linting
- **Automated Build & Release** - Integrated GoReleaser for seamless multi-platform releases
- **Code Quality Assurance** - 40+ linters via golangci-lint with comprehensive rules
- **License Compliance** - REUSE specification for automated license management
- **Coverage Tracking** - Codecov integration for test coverage monitoring
- **Metadata Injection** - Build-time version and commit information injection

## Tech Stack

- **Go** - Programming language
- **Task** - Build tool (Make alternative)
- **golangci-lint** - Static code analysis (40+ linters)
- **GoReleaser** - Automated release tool
- **REUSE** - License compliance tool
- **Codecov** - Code coverage tracking

## Project Structure

```text
.
├── api/              # API definitions (OpenAPI, protobuf, etc.)
├── cmd/              # Application entry points
│   └── demo/         # Demo application
├── configs/          # Configuration files
├── init/             # Initialization scripts (systemd, etc.)
├── internal/         # Internal packages (not importable by other projects)
│   └── metadata/     # Build metadata package
├── pkg/              # Public packages (importable by other projects)
├── scripts/          # Helper scripts
├── test/             # Test utilities and fixtures
├── LICENSES/         # License files
├── .golangci.yaml    # golangci-lint configuration
├── .goreleaser.yaml  # GoReleaser configuration
├── Taskfile.yaml     # Task build configuration
└── go.mod            # Go module definition
```

### Directory Descriptions

- **`cmd/`** - Contains the main entry points for your applications. Each subdirectory represents a separate executable.
- **`internal/`** - Private application and library code. This code cannot be imported by other projects.
- **`pkg/`** - Public library code that can be imported by other projects.
- **`api/`** - API definitions, protocol buffer files, OpenAPI specifications, etc.
- **`configs/`** - Configuration files for different environments.
- **`init/`** - System initialization scripts (systemd units, init scripts, etc.).
- **`scripts/`** - Helper scripts for development, deployment, and maintenance.
- **`test/`** - Test utilities, test fixtures, and integration test helpers.

## ⚠️ Important: Metadata Package Configuration

The `internal/metadata` package is a core component of this template. It's used to inject build-time information such as version, Git commit hash, and build date into your application.

This package is referenced by the following tools:

- **Taskfile.yaml** - Injects metadata via LDFLAGS during build
- **.goreleaser.yaml** - Injects metadata via LDFLAGS during release

### ⚠️ Critical: Synchronization Required

If you modify any of the following, you **must** update both `Taskfile.yaml` and `.goreleaser.yaml`:

1. **Changed the `internal/metadata` package path**
2. **Changed the module name** (the module path in `go.mod`)
3. **Changed variable names** in the metadata package (`AppName`, `Version`, `Commit`, `BuildDate`)

### What to Update

**In `Taskfile.yaml`:**

- `MODULE_PATH` variable
- `METADATA_PKG` variable
- `LDFLAGS` `-X` parameters

**In `.goreleaser.yaml`:**

- `ldflags` configuration

### Example

If you rename `internal/metadata` to `internal/buildinfo`, you need to update:

```yaml
# Taskfile.yaml
METADATA_PKG: '{{.MODULE_PATH}}/internal/buildinfo'
LDFLAGS: >-
  -X '{{.METADATA_PKG}}.Version=v{{.SEMVER}}'
  -X '{{.METADATA_PKG}}.Commit={{.GIT_COMMIT}}'
  -X '{{.METADATA_PKG}}.BuildDate={{.BUILD_DATE}}'

# .goreleaser.yaml
ldflags: &build_defaults_ldflags
  - -X {{ .ModulePath }}/internal/buildinfo.Version={{ .Version }}
  - -X {{ .ModulePath }}/internal/buildinfo.Commit={{ .FullCommit }}
  - -X {{ .ModulePath }}/internal/buildinfo.BuildDate={{ .Date }}
```

## Quick Start

### Prerequisites

- Go 1.26.1 or later
- Task (install via `go install github.com/go-task/task/v3/cmd/task@latest`)
- golangci-lint (install via `curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin`)

### Installation

Clone the repository

```bash
git clone https://github.com/ComixHe/go-template.git
cd go-template
```

Install dependencies

```bash
task tidy
```

Build the project

```bash
task build
```

Run the demo application

```bash
./bin/demo
```

Check version information

```bash
./bin/demo -v
```

## Development Guide

This project uses [Task](https://taskfile.dev/) as the build tool. All development commands are defined in `Taskfile.yaml`.

### Dependency Management

```bash
task tidy
```

Runs `go mod tidy` and `go mod verify` to clean up and verify dependencies.

### Lint Code

```bash
task lint
```

Runs comprehensive code checks including:

- REUSE license compliance verification
- golangci-lint with auto-fix enabled

```bash
task lint:ci
```

Runs code checks for CI environments (no auto-fix, 5-minute timeout).

```bash
task lint:license
```

Verifies REUSE license compliance only.

### Code Generation

```bash
task generate
```

Safely generates code (disables GOWORK to prevent module file corruption).

### Building

```bash
task default
```

Runs the standard development pipeline: tidy → lint → build:all

```bash
task build:all
```

Builds all applications found in the `cmd/` directory.

```bash
task build
```

Builds the default application (project name).

```bash
task build -- <appname>
```

Builds a specific application from the `cmd/` directory.

**Example:**

```bash
task build -- demo
```

### Testing

```bash
task test
```

Runs all tests with:

- Randomized test order
- 2-minute timeout
- Verbose output

```bash
task test:race
```

Runs all tests with race detector enabled:

- Generates coverage report (`coverage.txt`)
- Atomic coverage mode
- 5-minute timeout

```bash
task test:memory
```

Runs memory safety tests:

- Address Sanitizer (ASan)
- Memory Sanitizer (MSan)
- Requires CGO_ENABLED=1 and clang

```bash
task test:all
```

Runs all test suites: test → test:race → test:memory

### Installation and Cleanup

```bash
task install
```

Installs all built binaries to the system (default: `/usr/local/bin`).

**Custom installation path:**

```bash
PREFIX=/opt/myapp task install
```

```bash
task clean
```

Cleans build artifacts and test cache:

- Removes `bin/` directory
- Clears test cache
- Runs `go clean -v`

```bash
task distclean
```

Cleans all caches:

- Everything from `clean`
- Build cache
- Module cache
- Fuzz cache

## Build and Release

### Local Build

Build for your current platform:

```bash
task build
```

The binary will be placed in `bin/` directory with embedded metadata.

### Release with GoReleaser

This template is configured for automated releases using GoReleaser.

**Test release locally:**

```bash
goreleaser release --snapshot --clean
```

**Create a GitHub release:**

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
goreleaser release
```

### Supported Platforms

The template is configured to build for:

- **OS**: Linux
- **Architectures**: amd64, arm64, riscv64
- **AMD64 versions**: v3, v4
- **ARM64 versions**: v8.0, v9.0

## Code Quality

This project uses [golangci-lint](https://golangci-lint.run/) with 40+ linters configured in `.golangci.yaml`. Key linters include:

- **Error handling**: errcheck, errorlint, errname, wrapcheck
- **Security**: gosec
- **Performance**: prealloc, perfsprint, makezero
- **Style**: revive, gocritic, misspell
- **Complexity**: cyclop, gocognit
- **Testing**: testifylint

Run `task lint` to check your code and auto-fix issues where possible.

## License

This project is dual-licensed:

- **Code**: [Apache-2.0](LICENSES/Apache-2.0.txt)
- **Documentation**: [CC-BY-4.0](LICENSES/CC-BY-4.0.txt)

See the [LICENSES](LICENSES/) directory for full license texts.

## Acknowledgments

This template follows best practices from the Go community and incorporates tools and configurations from various open-source projects.
