# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project uses
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.13] - 2026-05-20

Release-automation bring-up. No change to generated PostgreSQL output.

### Added

- `.goreleaser.yaml` and a `release` job in `.github/workflows/ci.yml`
  gated on `refs/tags/v*`. Tag pushes (`vX.Y.Z`) now publish a complete
  release without manual `gh release` steps.
- Per-platform binaries: linux/{amd64,arm64} + darwin/{amd64,arm64} as
  `protoc-gen-psql_<version>_<os>_<arch>.tar.gz`.
- Per-platform CycloneDX SBOM via `syft` (binary-based scan, not just
  module graph) as `sbom_<os>_<arch>.cdx.json`.
- `SHA256SUMS` signed with `cosign sign-blob` (keyless, sigstore OIDC)
  producing `SHA256SUMS.sig`.

The release job depends on the full CI matrix (lint, unit-test,
test-generate, test-integration, govulncheck, build, sbom,
license-check). A tag that fails any gate never publishes binaries.

## [0.0.12] - 2026-05-20

Toolchain and dependency refresh. No change to generated PostgreSQL output.

### Added

- Multi-directory regression fixture (`tests/sub/asset.proto`) sharing the
  base filename `asset.proto` with the existing top-level fixture. Locks in
  source-relative output behaviour: any future regression flattening the
  output filename to `<basename>.<ext>` would silently overwrite one of the
  two `10_tables_asset.pb.psql` outputs and is now caught by
  `make test-generate`.

### Changed

- Bumped `github.com/lyft/protoc-gen-star` from `v0.6.2` to `v2.0.4`
  (new `/v2` module path).
- Regenerated `psql/psql.pb.go` with `protoc-gen-go` `v1.36.11` (was
  `v1.27.1`). Generated code now imports from
  `google.golang.org/protobuf/types/descriptorpb` instead of the legacy
  `github.com/golang/protobuf/protoc-gen-go/descriptor` path.
- Pinned `golangci-lint` to `v2.12.2` in CI and pre-commit (was `v2.10.1`).
- `Makefile` `test-generate` target now walks the full `tests/` tree via
  `find tests -name '*.proto'` instead of `tests/*.proto`, so sub-directory
  fixtures are picked up automatically.

### Removed

- Dropped the deprecated `github.com/golang/protobuf` direct dependency.
  All protobuf reflection now goes through `google.golang.org/protobuf`.
- Removed `gitleaks` from GitHub Actions, GitLab CI, and pre-commit hooks.
  gitleaks v8.x ships under BSL-1.1 and the project has no commercial
  license. Carve-out documented in `AGENTS.md` (`wont-fix`). Code review
  plus generic pre-commit hooks remain the safeguard.

## [0.0.11] - 2026-05-20

First release after iagen-dev brownfield onboarding (tier B Go CLI). No
behavioural change to the generated PostgreSQL output.

### Added

- `AGENTS.md` with tier B / Go-CLI operational contract and explicit
  carve-outs for runtime-only standards that do not apply to a build-time
  code generator.
- `docs/DEVELOPMENT.md` -- build, test, release runbook.
- `docs/superpowers/plans/` -- one re-runnable precheck plan plus
  per-gap correction plans driven by `/isec-iagen_dev-onboard-project`.
- `.golangci.yml` (golangci-lint v2 schema) enabling errcheck, govet
  (enable-all minus fieldalignment), staticcheck, ineffassign, unused,
  misspell, gofmt.
- `.pre-commit-config.yaml` -- gitleaks, golangci-lint, hygiene hooks.
- `renovate.json` -- weekly schedule, auto-merge indirect patch updates,
  group protobuf stack, manual review for major direct bumps.
- `.github/workflows/ci.yml` and `.gitlab-ci.yml` -- lint (golangci-lint,
  gitleaks), test (unit with coverage, test-generate, test-integration),
  vuln (govulncheck), build, supply-chain (CycloneDX SBOM, go-licenses
  allow-list gate).
- Weekly cron (Mon 06:00 UTC) on GitHub Actions to re-run vuln checks
  against `master` even with no commits.
- Vendored dependencies (`vendor/`); all Go CI jobs use `-mod=vendor`.
- Extended unit tests for `appendSlices`, `getStringBufferWithHeader`,
  `allocateRoomToParameters`, `generateCascadeIdentifierNames`,
  `generateFromTemplate`. Test count 5 -> 17, helper coverage now 88-100 %.

### Changed

- `go.mod` directive `go 1.17` -> `go 1.25.0`; new `toolchain go1.26.3`.
- Dependencies refreshed:
  - `github.com/golang/protobuf` v1.5.2 -> v1.5.4 (deprecated; migration
    tracked).
  - `github.com/lyft/protoc-gen-star` v0.5.3 -> v0.6.2.
  - `google.golang.org/protobuf` v1.27.1 -> v1.36.11.
  - `github.com/spf13/afero` v1.6.0 -> v1.15.0.
  - `golang.org/x/text` v0.3.5 -> v0.37.0.
- `Makefile`: `install` / `bin/protoc-gen-psql` targets now use
  `-mod=vendor`; `protoc` invocations gain `-I $(PROTOC_WKT_INCLUDE)`
  so well-known types resolve on stock Debian / Ubuntu / Homebrew;
  `docker-compose` invocations switched to a `$(DOCKER_COMPOSE)` helper
  that prefers the v2 plugin (`docker compose`).
- README.md: new Contributing section pointing to `docs/DEVELOPMENT.md`
  and `AGENTS.md`.

### Fixed

- `psqlify.go`: unchecked error on `pgs.Walk(v, field)` (errcheck) --
  now wrapped with `v.CheckErr` matching the existing pattern.
- `psqlify.go`: variable-shadow warnings on `ok`/`err` in the
  Prefix/Constraint/Suffix extension blocks -- renamed to
  `okExt`/`errExt`.
- `psqlify.go`: misspellings in comments (`analyse` -> `analyze`,
  `consistant` -> `consistent`).

### Security

- 8 stdlib CVEs surfaced by govulncheck on `go1.26.2` were cleared by
  the toolchain bump to `go1.26.3` (all were in net / net.http /
  html/template / net/mail and not on this plugin's call graph; bumping
  removes them from the vuln scan output entirely).
- All Go-module CVEs reported against earlier dependency versions cleared
  by the dep refresh.

[Unreleased]: https://github.com/Intrinsec/protoc-gen-psql/compare/v0.0.13...HEAD
[0.0.13]: https://github.com/Intrinsec/protoc-gen-psql/compare/v0.0.12...v0.0.13
[0.0.12]: https://github.com/Intrinsec/protoc-gen-psql/compare/v0.0.11...v0.0.12
[0.0.11]: https://github.com/Intrinsec/protoc-gen-psql/compare/v0.0.10...v0.0.11
