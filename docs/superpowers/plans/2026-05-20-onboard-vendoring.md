# Vendoring

## Goal

Add `vendor/` directory so builds are reproducible offline and dependency state is explicit in the repo. Required by tier B Go standard.

## Context

`vendor/` is absent. Project depends on `protoc-gen-star`, `google.golang.org/protobuf`, `golang/protobuf`, `afero`, `x/text` (5 direct + transitive). Vendoring locks the exact bytes used by CI and removes external-fetch failure modes.

## File Structure

- Modify: `.gitignore` (do not exclude `vendor/`).
- Modify: `Makefile` (add `vendor` target, use `-mod=vendor` in build).
- Create: `vendor/` (entire tree).
- Test: build with `-mod=vendor`; run `make test-generate`.

## Tasks

### Task 1: Populate `vendor/`

- [ ] `go mod tidy` — clean unused entries first.
- [ ] `go mod vendor` — populates `vendor/`.
- [ ] `git add vendor/ go.mod go.sum` — commit verbatim.

### Task 2: Update `.gitignore`

- [ ] Read `.gitignore` current content (4 lines: `bin/`, `*.pb.psql`, `!tests/references/*.pb.psql`, `tests/*.pb.go`).
- [ ] Confirm no entry excludes `vendor/`. If one is added later, remove it.

### Task 3: Update `Makefile`

- [ ] Add target:

  ```makefile
  .PHONY: vendor
  vendor:
  	@go mod tidy
  	@go mod vendor
  ```

- [ ] Modify `install` and `bin/protoc-gen-$(NAME)` targets to pass `-mod=vendor` to `go install`. Example:

  ```makefile
  @GOBIN=$(shell pwd)/bin go install -mod=vendor .
  ```

- [ ] Run `make build` — must succeed without network access (verify by `unshare -n make build` if possible, or temporarily `export GOFLAGS=-mod=vendor`).

### Task 4: CI cache

- [ ] In `.gitlab-ci.yml`, ensure jobs reference `-mod=vendor` and skip module download where possible.

## Verification (end-to-end)

- [ ] `ls vendor/modules.txt` exists.
- [ ] `go build -mod=vendor ./...` exits 0.
- [ ] `make test-generate` still produces empty diff against references.
- [ ] `git status` clean after `go mod vendor` (no untracked changes outside `vendor/`).

## Cross-references

- Standard: `dev-setup-project` vendoring section (tier B Go).
- Related plan: `2026-05-20-onboard-go-version-bump.md` (bump before vendoring to avoid re-vendoring).
