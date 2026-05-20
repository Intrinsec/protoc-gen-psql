# Vulnerability scanning

## Goal

Wire `govulncheck` into CI and document the indirect-CVE policy. Baseline already passes (0 reachable), so this plan locks the gate.

## Context

Baseline: `govulncheck ./...` reports 0 reachable vulnerabilities. 3 indirect (imported packages) + 8 module-level (declared but unused) CVEs exist but are not on the call graph. Tier B policy: gate on **reachable** CVEs only; indirect ones tracked via Renovate (see `dep-policy` plan).

## File Structure

- Modify: `.gitlab-ci.yml` (created by `ci-pipeline` plan) — add `vuln` stage.
- Create: `docs/superpowers/vuln-exceptions.md` (only if a reachable CVE is ever accepted).
- Test: re-run `govulncheck ./...`.

## Tasks

### Task 1: Add `vuln` stage to CI

- [ ] After `ci-pipeline` plan completes, add a job:

  ```yaml
  vuln:
    stage: test
    image: golang:1.22
    script:
      - go install golang.org/x/vuln/cmd/govulncheck@latest
      - govulncheck ./...
  ```

- [ ] Job must fail the pipeline on any reachable CVE (default `govulncheck` exit code).

### Task 2: Weekly cron

- [ ] In GitLab → CI/CD → Schedules, add weekly cron (`0 6 * * 1`) that runs the pipeline on `master`.
- [ ] Output published to a Slack channel or email group (define in CI variables).

### Task 3: Document acceptance procedure

- [ ] If a future reachable CVE is accepted (e.g. blocked upstream fix), append to `docs/superpowers/vuln-exceptions.md`: CVE ID, package, justification, target re-evaluation date.
- [ ] Reference this file in AGENTS.md `Vulnerability scanning` section.

## Verification (end-to-end)

- [ ] Re-run `2026-05-20-onboard-precheck.md` Task 2 — `govulncheck ./...` exits 0.
- [ ] CI `vuln` job appears green on `master`.
- [ ] Weekly schedule visible in GitLab CI schedules list.

## Cross-references

- Standard: `dev-setup-project` vuln-scanning section.
- Related skill: `isec-iagen_govulncheck`, `isec-iagen_gitlab-cicd-go`.
- Related plan: `2026-05-20-onboard-ci-pipeline.md`, `2026-05-20-onboard-dep-policy.md`.
