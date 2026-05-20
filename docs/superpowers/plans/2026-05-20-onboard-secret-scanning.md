# Secret scanning (gitleaks)

## Goal

Add `gitleaks` pre-commit hook and CI job to prevent accidental secret commits. Baseline-scan the existing history.

## Context

No secret scanning today. Project history is small but contains test fixtures (`tests/*.proto`, `tests/docker-compose.tests.yml`) that could plausibly include credentials. Tier B standard requires gitleaks or trufflehog.

## File Structure

- Create: `.gitleaks.toml` (allowlist) — only if baseline scan finds false positives.
- Create: `.pre-commit-config.yaml` (or extend `.lefthook.yml` if added by linting plan).
- Modify: `.gitlab-ci.yml` — add `secret-scan` job.
- Test: baseline scan + intentional fake secret to verify hook blocks it.

## Tasks

### Task 1: Baseline scan

- [ ] Install: `go install github.com/zricethezav/gitleaks/v8@latest` or `apt install gitleaks`.
- [ ] Run: `gitleaks detect --source . --no-banner --redact --report-format=json --report-path=/tmp/gitleaks-baseline.json`.
- [ ] Review findings. For each false positive, add to `.gitleaks.toml`:

  ```toml
  [[rules.allowlist]]
  description = "<reason>"
  paths = ['''<glob>''']
  ```

- [ ] If true positives are found: STOP. Rotate the secret first, then purge from git history (`git filter-repo`), then resume.

### Task 2: Pre-commit hook

- [ ] Create `.pre-commit-config.yaml`:

  ```yaml
  repos:
    - repo: https://github.com/gitleaks/gitleaks
      rev: v8.18.4
      hooks:
        - id: gitleaks
  ```

- [ ] Install: `pip install pre-commit && pre-commit install`.
- [ ] Document in `docs/DEVELOPMENT.md` under "Prerequisites".

### Task 3: CI job

- [ ] Add to `.gitlab-ci.yml`:

  ```yaml
  secret-scan:
    stage: test
    image: zricethezav/gitleaks:latest
    script:
      - gitleaks detect --source . --no-banner --redact --exit-code 1
  ```

- [ ] Job must fail pipeline on any new leak.

### Task 4: Quarterly review

- [ ] Each quarter, re-baseline + review `.gitleaks.toml` allowlist for entries that can be removed.

## Verification (end-to-end)

- [ ] `gitleaks detect --source . --no-banner --redact` exits 0.
- [ ] Inserting a fake AWS access key in a test branch triggers pre-commit failure.
- [ ] CI `secret-scan` job green on `master`.

## Cross-references

- Standard: `dev-setup-project` secret-scanning section.
