# Onboard Precheck — protoc-gen-psql

## Goal

Re-runnable diagnostic. Run before and after correction plans to verify state.
Records the baseline captured during onboard on **2026-05-20**.

## Baseline (frozen)

| Check | Tool version | Result | Exit code |
|-------|--------------|--------|-----------|
| `golangci-lint --version` | 2.10.1 (go1.26.0) | installed | 0 |
| `govulncheck -version` | v1.2.0 (DB 2026-05-07) | installed | 0 |
| `go version` | go1.26.2 linux/amd64 | installed (module declares 1.17) | 0 |
| `golangci-lint run ./...` | — | 1 issue: `errcheck` on `pgs.Walk(v, field)` in `psqlify.go` | 0 (default config tolerant) |
| `govulncheck ./...` | — | 0 reachable; 3 in imports + 8 in modules (not called) | 0 |
| `go test ./... -count=1 -short` | — | 5 passed in 2 packages | 0 |
| `make test-generate` | — | not run during baseline (requires `protoc` + plugin build) | — |
| `make test-integration` | — | not run during baseline (requires Docker) | — |

## Tasks

### Task 1: Verify tooling installed

- [ ] `golangci-lint --version` exits 0; version ≥ 1.55. Run `/isec-iagen_lint-go-install` if missing.
- [ ] `govulncheck -version` exits 0. Run `/isec-iagen_govulncheck-install` if missing.
- [ ] `go version` exits 0; version ≥ go.mod declared.
- [ ] `protoc --version` exits 0 (required by `make test-generate`).
- [ ] `docker compose version` exits 0 (required by `make test-integration`).

### Task 2: Run all mandatory checks

- [ ] `golangci-lint run ./...` — record issue count, compare to baseline (1).
- [ ] `govulncheck ./...` — record reachable CVE count, compare to baseline (0).
- [ ] `go test ./... -count=1 -short` — record pass/fail (baseline: 5 passed).
- [ ] `make test-generate` — must exit 0 with empty diff against `tests/references/`.
- [ ] `make test-integration` — must exit 0 (client container reports OK).

### Task 3: Report deltas

- [ ] Append run result + delta vs baseline to this plan file under `## Run history` (date, executor, each check result, regressions if any).

## Run history

<!-- Append a block per run -->

### 2026-05-20 by Stany MARCEL (Claude Opus 4.7, onboard/baseline-precheck) -- initial baseline

- Tooling — golangci-lint 2.10.1 ✓, govulncheck v1.2.0 ✓, go 1.26.2 ✓, protoc **missing** ✗, docker compose v5.0.2 ✓.
- golangci-lint: 1 issue (delta vs baseline: 0). `errcheck` on `pgs.Walk(v, field)` in `psqlify.go`.
- govulncheck: 0 reachable (delta: 0). 3 indirect + 8 module-level CVEs unchanged.
- go test: 5 passed in 2 packages (delta: 0).
- make test-generate: **blocked** — `protoc: command not found`, `make: *** [Makefile:29: psql/psql.pb.go] Error 127`.
- make test-integration: **skipped** — build depends on protoc-generated `psql/psql.pb.go`.
- Notes: Install `protoc` (apt `protobuf-compiler` or `brew install protobuf`) to unblock generation + integration paths. Tracked as prerequisite gate in `2026-05-20-onboard-dev-guide.md`.

### 2026-05-20 by Stany MARCEL (Claude Opus 4.7, onboard/baseline-precheck) -- post-corrections

After applying go-version-bump + dep-update, linting, dev-guide, dep-policy,
secret-scanning, vendoring, ci-pipeline correction plans on the same branch.

- Tooling — golangci-lint 2.10.1 ✓, govulncheck v1.2.0 ✓, go 1.26.2 ✓ (module now declares 1.25.0), gitleaks **missing** locally ✗ (CI uses container), protoc still missing locally ✗ (CI installs it).
- golangci-lint: **0 issues** (delta vs baseline: **-1**). Strict v2 config now enforced (govet enable-all minus fieldalignment, errcheck, staticcheck, ineffassign, unused, misspell, gofmt).
- govulncheck: 0 reachable (delta: 0). 3 indirect + 5 module-level CVEs (delta: **-3** after dep bumps).
- go test: 5 passed in 2 packages (delta: 0). Coverage: 11.0% total -- raise via `testing` correction plan.
- make test-generate: still blocked locally (protoc absent). CI job `test-generate` installs `protobuf-compiler` and exercises this path.
- make test-integration: still blocked locally (no Docker DinD here). CI job `test-integration` runs against ephemeral PostgreSQL.
- Notes: vendor/ now committed; all Go jobs use `-mod=vendor`. Pipeline configured for both GitHub Actions (primary) and GitLab CI (mirror).

<!--
### 2026-MM-DD by <name>
- golangci-lint: <count> (delta vs baseline: <±N>)
- govulncheck: <count> (delta: <±N>)
- go test: <pass/fail>
- make test-generate: <ok/diff>
- make test-integration: <ok/fail>
- Notes: ...
-->
