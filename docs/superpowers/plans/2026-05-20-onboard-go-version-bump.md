# Go version bump

## Goal

Bump `go.mod` declared Go version from `1.17` to `1.22` (minimum supported), pin toolchain, regenerate dependent files.

## Context

`go.mod` declares `go 1.17`. Local toolchain is `go1.26.2`. The 1.17→1.22 gap covers 5 years of language changes (generics, `errors.Join`, `slices`/`maps` packages, `min`/`max` builtins). Tier B should track at most N-2 minor versions behind current stable (currently 1.24 stable → N-2 = 1.22).

## File Structure

- Modify: `go.mod` (line 3: `go 1.17` → `go 1.22`).
- Modify: `go.sum` (refreshed by `go mod tidy`).
- Modify: `.gitlab-ci.yml` (`image: golang:1.22`).
- Modify: `vendor/` (regenerate after bump — see vendoring plan).
- Test: full build + tests.

## Tasks

### Task 1: Bump `go.mod`

- [ ] Edit `go.mod` line 3 to `go 1.22`.
- [ ] Optionally add `toolchain go1.22.0` line below the `go 1.22` line to pin toolchain across machines.
- [ ] `go mod tidy` — refresh `go.sum`.

### Task 2: Verify builds

- [ ] `make build` — exits 0.
- [ ] `go test ./... -count=1` — 5 passed.
- [ ] `make test-generate` — empty diff against references.
- [ ] `golangci-lint run ./...` — same issue count (1 from baseline; addressed in linting plan).
- [ ] `govulncheck ./...` — 0 reachable.

### Task 3: Coordinate with vendoring + CI plans

- [ ] If `vendoring` plan already ran: re-run `go mod vendor` to refresh vendor tree under new Go version.
- [ ] If `ci-pipeline` plan already ran: update CI `image:` references from any older `golang:1.X` to `golang:1.22`.

## Verification (end-to-end)

- [ ] `head -3 go.mod` shows `go 1.22`.
- [ ] All checks in Task 2 exit 0.
- [ ] Re-run `2026-05-20-onboard-precheck.md` Task 2.

## Cross-references

- Related plans: `vendoring`, `ci-pipeline`.
- Standard: `dev-setup-project` language section.
