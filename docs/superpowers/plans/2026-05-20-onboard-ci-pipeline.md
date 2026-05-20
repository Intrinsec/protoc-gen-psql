# CI pipeline (GitLab)

## Goal

Bootstrap `.gitlab-ci.yml` covering lint → test → vuln → build for Go protoc plugin, using `isec-iagen_gitlab-cicd-go` skill.

## Context

No CI configuration exists. Project lives on internal intrinsec GitLab (module path `github.com/intrinsec/...`). Tier B Go CLI requires CI gating on every MR. Skill `isec-iagen_gitlab-cicd-go` produces a working config that avoids four recurring failure modes (cache misses, cross-compile, govulncheck install, golangci-lint version).

## File Structure

- Create: `.gitlab-ci.yml`.
- Modify: `README.md` (add CI badge once pipeline runs).
- Test: push a MR with a trivial change, verify all jobs green.

## Tasks

### Task 1: Run the skill

- [ ] Invoke `/isec-iagen_gitlab-cicd-go`.
- [ ] Choose stages: `lint`, `test`, `vuln`, `build`. Skip `release` for now (deferred per carve-out).
- [ ] Image: `golang:1.22` (or matching `go.mod` declared version after `go-version-bump` plan completes).
- [ ] Pin `golangci-lint` to a stable v2 release (≥ 2.10.1, current installed).

### Task 2: Add `test-generate` and `test-integration`

- [ ] Append jobs:

  ```yaml
  test-generate:
    stage: test
    image: golang:1.22
    before_script:
      - apt-get update && apt-get install -y protobuf-compiler
    script:
      - make test-generate

  test-integration:
    stage: test
    image: docker:24
    services:
      - docker:24-dind
    script:
      - apk add make
      - make test-integration
  ```

- [ ] Ensure `test-integration` uses `CI_JOB_ID` already exposed by the Makefile.

### Task 3: MR template + branch protection

- [ ] Add `.gitlab/merge_request_templates/Default.md` with checklist (tests pass, lint clean, AGENTS.md respected).
- [ ] Enable “pipelines must succeed” on `master` branch protection.

## Verification (end-to-end)

- [ ] Push a test branch, open MR, verify all jobs run and pass.
- [ ] Pipeline badge added to `README.md`.
- [ ] Re-run `2026-05-20-onboard-precheck.md` Task 2 — same results locally and in CI.

## Cross-references

- Standard: `dev-setup-project` CI section.
- Related skill: `isec-iagen_gitlab-cicd-go`.
- Related plans: `vuln-scanning`, `vendoring`, `go-version-bump`.
