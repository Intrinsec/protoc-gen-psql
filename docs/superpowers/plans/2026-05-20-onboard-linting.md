# Linting

## Goal

Add `.golangci.yml`, fix the one outstanding `errcheck` violation, and make `golangci-lint run ./...` exit 0 with zero issues.

## Context

Baseline: 1 issue — `errcheck` on `pgs.Walk(v, field)` in `psqlify.go`. No `.golangci.yml`, so the default rule set is in effect. Tier B Go requires an explicit, version-pinned linter config so behaviour is reproducible across contributors and CI.

## File Structure

- Modify: `psqlify.go` (fix unchecked error from `pgs.Walk`).
- Create: `.golangci.yml`.
- Test: re-run `golangci-lint run ./...` until clean.

## Tasks

### Task 1: Generate `.golangci.yml` via skill

- [ ] Run `/isec-iagen_lint-go-config` and accept the tier-B / Go-CLI default.
- [ ] If skill emits choices, enable at minimum: `errcheck`, `govet`, `staticcheck`, `ineffassign`, `unused`, `gosimple`, `gofmt`, `misspell`. Disable nothing aggressive (`gochecknoinits` is OK to keep off — `main` package registers a global module).
- [ ] Pin `version: "2"` (golangci-lint v2 schema, matching installed 2.10.1).
- [ ] Commit: `chore(lint): add .golangci.yml (tier B baseline)`.

### Task 2: Fix `errcheck` in `psqlify.go`

- [ ] Read `psqlify.go` around the `pgs.Walk(v, field)` call site.
- [ ] Capture the returned error and surface it via `v.Push(...)` / `v.Fail(...)` or wrap with `pgs.CheckErr(...)` per `protoc-gen-star` idiom.
- [ ] Re-run `golangci-lint run ./...` — must exit 0 with `0 issues`.
- [ ] Re-run `go test ./... -count=1` — must still pass 5/5.
- [ ] Re-run `make test-generate` — diff against `tests/references/` must remain empty (no behaviour change).

### Task 3: Wire into pre-commit

- [ ] Add `golangci-lint run ./...` invocation to the `pre-commit` hook (or `.lefthook.yml` / `.pre-commit-config.yaml` if either exists; otherwise document in `docs/DEVELOPMENT.md`).

## Verification (end-to-end)

- [ ] Re-run `2026-05-20-onboard-precheck.md` Task 2 — `golangci-lint run ./...` reports `0 issues`.
- [ ] `.golangci.yml` committed; `cat .golangci.yml | head -5` shows `version: "2"`.
- [ ] AGENTS.md section `Linting` still references this config.

## Cross-references

- Standard: `dev-setup-project` linting section.
- Related skill: `isec-iagen_lint-go-config`, `isec-iagen_lint-go`.
