# DEVELOPMENT

Build, test, and release runbook for `protoc-gen-psql`.

For the project's operational contract (tier, type, carve-outs) see `AGENTS.md`.
For cascade-trigger design rationale see `docs/dev.md`.

## Prerequisites

| Tool | Minimum | Install |
|------|---------|---------|
| Go | matches `go.mod` (currently `go 1.25` directive, `toolchain go1.26.3`) | https://go.dev/dl/ |
| `protoc` | 3.20 | `apt install protobuf-compiler` or `brew install protobuf` |
| Docker + Compose v2 | 24.x | https://docs.docker.com/engine/install/ |
| `golangci-lint` | 2.10.1 | https://golangci-lint.run/usage/install/ |
| `govulncheck` | v1.2.0 | `go install golang.org/x/vuln/cmd/govulncheck@latest` |
| `pre-commit` (optional) | 3.x | `pip install pre-commit && pre-commit install` |

## Build

```sh
make build
```

Produces `bin/protoc-gen-psql`. The Makefile also produces a temporary
`bin/protoc-gen-go` to regenerate `psql/psql.pb.go` from `psql/psql.proto`.

To install on `$PATH`:

```sh
make install
```

This runs `go install -v .` and drops the binary in `$GOBIN`.

## Generate `psql/psql.pb.go`

The plugin's own protobuf option schema lives in `psql/psql.proto` and is
compiled to `psql/psql.pb.go`. This file is committed to the repository.

Regenerate after editing `psql/psql.proto`:

```sh
make psql/psql.pb.go
```

## Run tests

```sh
# Unit tests
go test ./... -count=1

# Generation tests (diff plugin output vs tests/references/)
make test-generate

# Integration tests (apply generated SQL to a PostgreSQL container)
make test-integration

# Everything
make test
```

`make test-generate` regenerates `tests/*.pb.psql` and diffs them against the
oracle files in `tests/references/`. A non-empty diff means the change altered
plugin output -- review carefully before updating references.

## Update reference files

When an intentional change to generated SQL is made:

```sh
make test-generate            # generates and diffs
# inspect diff, confirm it is the intended change
cp tests/*.pb.psql tests/references/
git add tests/references/
```

Commit the reference update in the same change that altered the generator.

## Lint and vuln scan

```sh
golangci-lint run ./...
govulncheck ./...
```

Both must exit 0 with no findings (lint: 0 issues, vuln: 0 reachable). See
`AGENTS.md` for the policy on module-level CVEs not on the call graph.

## Pre-commit hooks

Install once:

```sh
pip install pre-commit
pre-commit install
```

Hooks run gitleaks (secret scanning), golangci-lint, and pre-commit-hooks
sanity checks. Configuration in `.pre-commit-config.yaml`.

## Clean

```sh
make clean            # remove generated *.pb.psql
make distclean        # also remove bin/ and psql/psql.pb.go
```

## Release (deferred)

Release signing and SLSA provenance are currently carved out (see `AGENTS.md`).
Re-evaluate when binaries ship beyond internal CI.

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| `protoc: command not found` | `protobuf-compiler` not installed | `apt install protobuf-compiler` |
| `make: *** [Makefile:29: psql/psql.pb.go] Error 127` | missing protoc | install protoc, re-run `make build` |
| `bin/protoc-gen-go` install fails | offline / GOPROXY unreachable | enable network, or use vendored deps with `-mod=vendor` |
| Integration tests hang | leftover containers from previous run | `docker compose -p psql-local -f tests/docker-compose.tests.yml down -v` |
| `golangci-lint` reports unexpected issues | version mismatch with `.golangci.yml` schema | reinstall matching version (config pins schema `version: "2"`) |
| `make test-generate` diff non-empty after small refactor | reference file out of date or behaviour regression | inspect diff; update `tests/references/` only if change is intentional |
