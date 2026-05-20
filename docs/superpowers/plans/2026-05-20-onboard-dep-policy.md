# Dependency policy (Renovate)

## Goal

Add `renovate.json` so Go module updates land as MRs on a predictable cadence, with auto-merge on safe patch updates.

## Context

5 direct deps in `go.mod`. No update automation. 11 transitive CVEs found by govulncheck (none reachable) — Renovate keeps the indirect tree fresh and reduces the chance one becomes reachable. Tier B Go standard requires Renovate or Dependabot.

## File Structure

- Create: `renovate.json` in repo root.
- Modify: `.gitlab-ci.yml` — optional Renovate-self-hosted runner job (if intrinsec hosts its own bot).
- Test: validate JSON; install Renovate config validator locally.

## Tasks

### Task 1: Write `renovate.json`

- [ ] Create with this content:

  ```json
  {
    "$schema": "https://docs.renovatebot.com/renovate-schema.json",
    "extends": [
      "config:recommended",
      ":enableVulnerabilityAlertsWithLabel(security)",
      ":disableRateLimiting",
      ":semanticCommits"
    ],
    "schedule": ["before 6am on monday"],
    "timezone": "Europe/Paris",
    "labels": ["dependencies"],
    "packageRules": [
      {
        "matchManagers": ["gomod"],
        "matchDepTypes": ["indirect"],
        "matchUpdateTypes": ["patch"],
        "automerge": true,
        "automergeType": "branch"
      },
      {
        "matchManagers": ["gomod"],
        "matchDepTypes": ["direct"],
        "matchUpdateTypes": ["major"],
        "labels": ["dependencies", "needs-review"]
      },
      {
        "matchManagers": ["gomod"],
        "groupName": "go protobuf stack",
        "matchPackagePatterns": ["protobuf", "protoc-gen"]
      }
    ],
    "vulnerabilityAlerts": {
      "labels": ["security"],
      "automerge": false
    }
  }
  ```

- [ ] Validate: `npx --yes renovate-config-validator renovate.json`.

### Task 2: Enable Renovate on the repo

- [ ] If intrinsec hosts a Renovate runner: add this project to the runner's repo list.
- [ ] Otherwise: install Renovate GitLab app/bot on the project.

### Task 3: Document policy

- [ ] Append to AGENTS.md `Dependency policy` section: "Renovate runs Monday 06:00 Europe/Paris. Patch updates to indirect deps auto-merge. Major updates require human review and the `needs-review` label."

## Verification (end-to-end)

- [ ] `renovate-config-validator renovate.json` exits 0.
- [ ] First Renovate MR appears within one week of activation.
- [ ] Auto-merged patch updates do not break CI.

## Cross-references

- Standard: `dev-setup-project` dep-policy section.
- Related plan: `2026-05-20-onboard-vuln-scanning.md`.
