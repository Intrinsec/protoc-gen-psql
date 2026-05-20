# Testing — extend unit coverage

## Goal

Raise unit test coverage of `psqlify.go` so behaviour changes break fast. Today only `psqlify_test.go` exists with limited scope.

## Context

Baseline: `go test ./... -count=1 -short` → 5 passed. Coverage not measured. `psqlify.go` is 15 KB of rendering logic (options parsing, file-prefix sorting, trigger generation). Integration tests via `make test-generate` cover end-to-end but are slow; unit tests should cover the small pure functions.

## File Structure

- Modify: `psqlify_test.go` (extend).
- Create: optional `psqlify_internal_test.go` if testing unexported functions.
- Test: `go test ./... -cover` should report `psqlify.go` ≥ 50% line coverage after this plan.

## Tasks

### Task 1: Measure current coverage

- [ ] `go test ./... -cover -count=1 -short 2>&1 | tee /tmp/coverage-before.txt` — record numbers.
- [ ] Identify uncovered functions: `go test ./... -coverprofile=/tmp/cov.out && go tool cover -func=/tmp/cov.out | sort -k3 -n`.

### Task 2: Add unit tests for high-value functions

Priority list (from `psqlify.go` skim):

- [ ] File-prefix naming: `00_init_…`, `10_tables_…`, `20_relations_…`, `99_final_…` — verify ordering rules.
- [ ] `auto_fill_on_update` trigger SQL renderer (golden-file style).
- [ ] `relay_cascade_update` renderer covering both `INSERT`/`UPDATE`/`DELETE` branches.
- [ ] `cascade_update_on_related_table` renderer.
- [ ] Constraint emission and idempotency check.

Each test:
- Construct minimal input struct or read fixture `*.proto`.
- Compare output string to expected snippet (use `testify/require.Equal` or stdlib `if got != want`).
- No network, no Docker.

### Task 3: Coverage gate

- [ ] Add CI job step: `go test ./... -coverprofile=cov.out && go tool cover -func=cov.out | grep total | awk '{ if ($3+0 < 50.0) { print "coverage below 50%"; exit 1 } }'`.
- [ ] Document threshold (50% initial, raise to 70% next quarter) in AGENTS.md `Testing` section.

## Verification (end-to-end)

- [ ] `go test ./... -count=1 -short` — all green.
- [ ] `go test ./... -cover` — total ≥ 50%.
- [ ] `make test-generate` and `make test-integration` still pass (no regression).

## Cross-references

- Standard: `dev-setup-project` testing section.
- Related plan: `2026-05-20-onboard-ci-pipeline.md` (adds `go test` step).
