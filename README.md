# Golang Project Template

[![Checks][badge-checks]][actions]
[![Latest Release][badge-release]][releases]
[![Codecov][badge-codecov]][codecov]
[![Go Report Card][badge-grc]][grc]

[badge-checks]: https://img.shields.io/github/check-runs/ComixHe/go-template/main
[actions]: https://github.com/ComixHe/go-template/actions

[badge-release]: https://img.shields.io/github/release-date/ComixHe/go-template
[releases]: https://github.com/ComixHe/go-template/releases

[badge-codecov]: https://codecov.io/gh/ComixHe/go-template/graph/badge.svg
[codecov]: https://codecov.io/gh/ComixHe/go-template

[badge-grc]: https://goreportcard.com/badge/github.com/ComixHe/go-template
[grc]: https://goreportcard.com/report/github.com/ComixHe/go-template

A small, modern Go application template:

- [go.dev layout](https://go.dev/doc/modules/layout): `cmd/` for binaries, `internal/` for everything else
- No runtime dependencies, dev tools are pinned in `go.mod` with Go 1.24 `tool` directives
- Task-based workflow, minimal but strict lint set (43 linters), race and vulnerability gates
- Pinned Go/tool/action versions, REUSE license compliance, Dependabot, Codecov
- Tag-driven releases: draft GitHub release, SBOMs, checksums, Sigstore signing, provenance attestation

## Requirements

- Go 1.27.3+
- [Task](https://taskfile.dev/) 3.54+
- [golangci-lint](https://golangci-lint.run/) v2.15+
- Optional: `reuse`, GoReleaser v2.20.3+, `syft`, `cosign`

## Quick start

```bash
git clone https://github.com/ComixHe/go-template.git
cd go-template
task tidy
task build
./bin/demo
./bin/demo -version
```

## Layout

```text
cmd/demo/              # one directory per binary (package main)
internal/metadata/     # private packages, not importable by other modules
```

No `pkg/`, `api/`, `configs/`, `init/`, `scripts/` or `test/`: add them only when
the project has real content for them. Tests live next to the code, fixtures go
in `testdata/`. New command: `cmd/<name>/main.go`. New private package:
`internal/<name>/`.

## Tasks

| Task | Description |
| --- | --- |
| `task tidy` | `go mod tidy` + `go mod verify` |
| `task fmt` / `task fmt:check` | Format / verify formatting (gofumpt, gci, goimports, golines) |
| `task lint` / `task lint:ci` | golangci-lint with / without `--fix` |
| `task lint:actions` | actionlint over `.github/workflows/` |
| `task verify` | Pre-merge gate: lint, format, workflow lint, race tests, govulncheck |
| `task test` / `task test:race` | Tests / tests with race detector + coverage |
| `task test:asan` / `task test:msan` | Sanitizers (linux/amd64, cgo + clang) |
| `task vuln` | `govulncheck` (via `go tool`) |
| `task build` / `task build:all` | Build the first `cmd/` app / all apps |
| `task build -- <app>` | Build a specific app |
| `task install` | Install built binaries (`PREFIX`, default `/usr/local`) |
| `task clean` / `task distclean` | Remove artifacts / all caches |

`DEBUG=1 task build` produces a debuggable binary. `task generate` runs
`go generate` and is intentionally not part of `build`.

## Using this template

1. **Module path**

   ```bash
   go mod edit -module github.com/<owner>/<project>
   grep -rl 'github.com/ComixHe/go-template' --include='*.go' . \
     | xargs sed -i 's|github.com/ComixHe/go-template|github.com/<owner>/<project>|g'
   ```

   `Taskfile.yaml` and `.goreleaser.yaml` derive the module path automatically.

2. **Application name**

   ```bash
   git mv cmd/demo cmd/<app>
   ```

   Then update `builds.main` in `.goreleaser.yaml` (plain path, not a template).

3. **License and ownership** — update `REUSE.toml`, `LICENSES/`, README badges
   and the Codecov settings.

4. **First release** — see below. `cd.yaml` runs the tests and publishes a
   **draft** release for review.

## Metadata

`internal/metadata` holds `AppName`, `Version` and `Commit`,
injected through `-ldflags -X`. If you rename the package or its variables,
update **both** `METADATA_PKG`/`LDFLAGS` in `Taskfile.yaml` and `ldflags` in
`.goreleaser.yaml`. The module path itself is derived automatically.

## Release

A `v*` tag triggers `.github/workflows/cd.yaml`:

1. **verify** — tag name, REUSE, `task test:race`, `goreleaser check`
2. **release** — build, then create a draft GitHub release

Tag names must match `vMAJOR.MINOR.PATCH` (a `-prerelease` suffix is allowed).
Signing tags is recommended: `git tag -s v1.0.0 -m "Release v1.0.0"`.

Each release contains:

- `*_linux_{amd64,arm64,riscv64}.tar.gz` and a source tarball
- one SPDX SBOM per artifact (syft)
- `checksums.txt` (sha256) and `checksums.txt.sigstore.json` (cosign keyless)
- build provenance attested over `checksums.txt`, covering every artifact

Actions are pinned to full commit SHAs; Dependabot keeps the pins and version comments up to date.

Verify a download:

```bash
gh attestation verify --owner <owner> <artifact>
cosign verify-blob --bundle checksums.txt.sigstore.json \
  --certificate-identity-regexp 'https://github.com/<owner>/<repo>/.github/workflows/cd.yaml@refs/tags/.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  checksums.txt
sha256sum -c checksums.txt
```

Local release check (signing needs CI's OIDC identity, so skip it locally):

```bash
goreleaser check
goreleaser release --snapshot --clean --skip=sign,sbom
```

### Release immutability

For a supply-chain-safe release flow, enable both protections once per repository:

1. **Settings → General → Releases → Enable release immutability.** Applies to
   future releases.
2. A **tag ruleset** for `v*` with *Restrict updates* and *Restrict deletions*
   (creation stays allowed).

`cd.yaml` creates a draft with every asset attached; review it, then publish
manually. Publishing seals the release:

- the tag is locked to its commit and cannot be moved or deleted while the
  release exists
- release assets cannot be added, modified or deleted
- after deleting the release you can delete the tag, but that tag name can
  never be reused

If a draft must be regenerated, delete the draft (not the tag) and re-run the
workflow. A published tag is final: cut a new version instead of reusing it.

## Code quality

`.golangci.yaml` enables 43 linters with no exclusion presets, no per-path
silencing and no `uniq-by-line` deduplication. Complexity is measured with
cognitive complexity (`gocognit`), not cyclomatic complexity. Suppressions must
be narrow and explained:

```go
//nolint:errcheck // the error is reported by the caller
```

Action and vulnerability tool versions live in `go.mod` (`tool` block), so Dependabot updates them with the module graph.

Run `task lint` to fix what can be fixed and `task verify` before pushing.

## License

Code is Apache-2.0, documentation is CC-BY-4.0. See [LICENSES](LICENSES/).
