# Dev guide

## Goal

Write `docs/DEVELOPMENT.md` — runbook for build, test, release, troubleshoot. Distinct from `README.md` (user-facing) and `docs/dev.md` (ADR on cascade triggers).

## Context

No dev runbook exists. Newcomers must reverse-engineer the Makefile to learn how to build and test. Tier B Go CLI standard requires a single, terse `DEVELOPMENT.md`.

## File Structure

- Create: `docs/DEVELOPMENT.md`.
- Modify: `README.md` — link to `docs/DEVELOPMENT.md` in a "Contributing" section.
- Test: a fresh checkout + following `DEVELOPMENT.md` step-by-step must produce a working binary.

## Tasks

### Task 1: Draft skeleton

- [ ] Create `docs/DEVELOPMENT.md` with sections:

  - **Prerequisites** — Go ≥ go.mod, `protoc` ≥ 3.20, Docker (for integration tests).
  - **Build** — `make build`. Output: `bin/protoc-gen-psql`.
  - **Install on $PATH** — `make install`. Note: requires `bin/` removed from `GOBIN`.
  - **Generate `psql/psql.pb.go`** — `make psql/psql.pb.go`. Triggers automatically by `make build`.
  - **Run unit tests** — `go test ./... -count=1 -short`.
  - **Run generation tests** — `make test-generate`. Diffs against `tests/references/`.
  - **Run integration tests** — `make test-integration`. Spins up PostgreSQL via `tests/docker-compose.tests.yml`, runs assertions in `client` container.
  - **Update reference files** — when an intentional change to generated SQL output is made: `cp tests/*.pb.psql tests/references/` then review diff and commit.
  - **Clean** — `make clean` (artifacts), `make distclean` (artifacts + generated `.pb.go`).
  - **Lint** — `golangci-lint run ./...`.
  - **Vuln scan** — `govulncheck ./...`.
  - **Release** — tag with `git tag vX.Y.Z && git push --tags` (deferred until release-signing plan is unfrozen).
  - **Troubleshooting**:
    - `protoc: command not found` → install `protobuf-compiler` (apt) or `brew install protobuf`.
    - `bin/protoc-gen-go` build fails → check `go install google.golang.org/protobuf/cmd/protoc-gen-go` network access.
    - Integration tests hang → `docker compose -p psql-local -f tests/docker-compose.tests.yml down -v` then retry.

### Task 2: Link from README

- [ ] Append section to `README.md`:

  ```markdown
  ## Contributing

  See `docs/DEVELOPMENT.md` for build, test, and release procedures.
  See `AGENTS.md` for the iagen-dev operational contract.
  ```

### Task 3: Validate by dry run

- [ ] Fresh shell, `cd` into project, follow `DEVELOPMENT.md` top to bottom. Any step that fails must be fixed in the doc, not papered over.

## Verification (end-to-end)

- [ ] `test -f docs/DEVELOPMENT.md && wc -l docs/DEVELOPMENT.md` — > 40 lines.
- [ ] README.md references `docs/DEVELOPMENT.md`.
- [ ] All commands in `Build`/`Test` sections run clean.

## Cross-references

- Standard: `dev-setup-project` dev-guide section.
- Existing related doc: `docs/dev.md` (cascade triggers ADR; keep separate).
