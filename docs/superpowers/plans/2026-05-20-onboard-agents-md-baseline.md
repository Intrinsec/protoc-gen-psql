# AGENTS.md baseline

## Goal

Establish the operational contract (`AGENTS.md`) for tier B / Go CLI with carve-outs, so future iagen-dev skill runs do not re-suggest excluded sections.

## Context

Tier B shared internal tool, type 2 (Go CLI / protoc plugin). No runtime SLA. `AGENTS.md` was absent at onboard time (2026-05-20). Carve-outs needed because runtime-oriented sections (monitoring, auth, secrets, DB) do not apply to a build-time code generator.

## File Structure

- Modify: none
- Create: `AGENTS.md` (already written during onboard — this plan documents and codifies maintenance)
- Test: `AGENTS.md` is human-readable; lint with `markdownlint` if available.

## Tasks

### Task 1: Verify AGENTS.md exists and matches baseline

- [ ] `test -f AGENTS.md && echo ok` — must print `ok`.
- [ ] Grep for the carve-outs table: `grep -c '^| ' AGENTS.md` ≥ 10 (header + ≥9 carve-out rows).
- [ ] Confirm sections present: `Language`, `Workflow Skills`, `Linting`, `Vulnerability scanning`, `Testing`, `Vendoring`, `CI pipeline`, `Dependency policy`, `Secret scanning`, `Documentation`, `Carve-outs`.

### Task 2: Quarterly refresh trigger

- [ ] On each quarterly review (1 Jan / 1 Apr / 1 Jul / 1 Oct), re-run `/isec-iagen_dev-update-project` to refresh drift.
- [ ] If tier or type changes, re-run `/isec-iagen_dev-onboard-project` to recompute carve-outs.

### Task 3: Record carve-out changes

- [ ] Any addition or removal of a row in the `## Carve-outs` table must be committed with a message `chore(agents): carve-out <section> — <reason>` and reference the discussion that drove the change.

## Verification (end-to-end)

- [ ] `grep -A1 '^## Carve-outs' AGENTS.md` shows the explanatory paragraph.
- [ ] No iagen-dev skill re-suggests a carved-out section (manual check next time a skill runs).

## Cross-references

- Standard: `dev-setup-project` template (tier B, type 2).
- Related skill: `dev-update-project` for routine drift refresh.
