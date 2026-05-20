# AGENTS.md — protoc-gen-psql

Operational contract for AI agents (Claude Code, Copilot, Gemini) and human contributors. Source of truth for tooling, workflow, and out-of-scope items.

- **Tier**: B (shared internal tool, no runtime SLA)
- **Type**: 2 — Go CLI (`protoc` plugin binary)
- **Module**: `github.com/intrinsec/protoc-gen-psql`
- **Purpose**: generate PostgreSQL DDL files from protobuf message definitions

## Language

- **Go**: `1.17` declared in `go.mod`. Toolchain in CI must match or exceed. Bump tracked in `docs/superpowers/plans/2026-05-20-onboard-go-version-bump.md`.
- **Build artifact**: single static binary `bin/protoc-gen-psql` installed via `go install .`
- **Entry point**: `main.go` registers `PSQLify()` module against `protoc-gen-star`.

## Workflow Skills

Mandatory iagen-dev skill order for any non-trivial change:

1. `superpowers:brainstorming` — before adding new psql option or behaviour.
2. `superpowers:writing-plans` — when scope ≥ 3 files or any cross-cutting refactor.
3. `superpowers:test-driven-development` — every behaviour change needs a failing test first. Unit (`*_test.go`) or integration (Docker compose) as appropriate.
4. `superpowers:verification-before-completion` — re-run `make test-generate` and `make test-integration` before declaring done.
5. `isec-iagen_lint-go` — `golangci-lint run ./...` clean before commit.
6. `isec-iagen_govulncheck` — `govulncheck ./...` clean before tagging release.

## Linting

- Config: `.golangci.yml` (to be created — see `docs/superpowers/plans/2026-05-20-onboard-linting.md`).
- Required enabled: `errcheck`, `govet`, `staticcheck`, `ineffassign`, `unused`, `gosimple`, `gofmt`, `misspell`.
- Run: `golangci-lint run ./...`
- Gate: zero issues on `master`. Use `//nolint:<linter> // reason` with justification only.

## Vulnerability scanning

- Tool: `govulncheck` (`isec-iagen_govulncheck` skill).
- Cadence: every CI build + weekly cron.
- Gate: zero **reachable** vulnerabilities. Indirect/module-level CVEs documented in `docs/superpowers/plans/<date>-vuln-<id>.md` if accepted.

## Testing

- **Unit**: `psqlify_test.go` (extend coverage of `psqlify.go` renderer functions). `go test ./... -count=1 -short`.
- **Generation**: `make test-generate` — runs plugin against `tests/*.proto`, diffs against `tests/references/*.pb.psql`.
- **Integration**: `make test-integration` — applies generated SQL to PostgreSQL in `tests/docker-compose.tests.yml`, runs assertions in the `client` container.
- Reference files in `tests/references/` are the regression oracle — update intentionally, never accept silent diffs.

## Vendoring

- `vendor/` to be added — tier B Go requires reproducible builds offline.
- Tracked in `docs/superpowers/plans/2026-05-20-onboard-vendoring.md`.

## CI pipeline

- Target: GitLab CI via `.gitlab-ci.yml` (to be bootstrapped with `isec-iagen_gitlab-cicd-go` skill).
- Stages: lint → test → vuln → build → (optional) release.
- Tracked in `docs/superpowers/plans/2026-05-20-onboard-ci-pipeline.md`.

## Dependency policy

- Renovate (`renovate.json`) with grouped Go module updates, monthly schedule, auto-merge on patch level for indirect deps.
- Tracked in `docs/superpowers/plans/2026-05-20-onboard-dep-policy.md`.

## Secret scanning

- Pre-commit + CI gate via `gitleaks`.
- Baseline scan + ignore list reviewed quarterly.
- Tracked in `docs/superpowers/plans/2026-05-20-onboard-secret-scanning.md`.

## Documentation

- `README.md` — user-facing plugin usage and options reference.
- `docs/dev.md` — design rationale for cascade-update triggers (ADR-style, keep).
- `docs/DEVELOPMENT.md` — to be created: build, test, release runbook.
  - Tracked in `docs/superpowers/plans/2026-05-20-onboard-dev-guide.md`.

## Carve-outs

The following standard sections are explicitly excluded from this project. Re-evaluate
on tier change or quarterly review.

| Section | Reason | Note |
|---------|--------|------|
| Monitoring (Prometheus, Grafana, alerts) | not-applicable | Build-time CLI plugin, no runtime process to monitor. |
| Auth (Keycloak client) | not-applicable | No HTTP/gRPC surface, no user sessions. |
| Secrets management (Vault AppRole) | not-applicable | No runtime secrets consumed. Build-time stdin/stdout only. |
| Database (CNPG manifest + backups) | not-applicable | Plugin emits SQL files; does not own a database. |
| Container hardening | not-applicable | No production Dockerfile; `tests/Dockerfile.psqlclient` is test-only. |
| Release signing (cosign / SLSA provenance) | deferred | Tier B build tool; revisit if distributed beyond internal CI. |
| SBOM (CycloneDX/syft) | deferred | Revisit when project ships binaries outside internal CI. |
| License compliance gate | replaced-by:manual-review | Small dep tree (5 direct), reviewed via Renovate PRs. |
| Threat model | not-applicable | Pure code generator, no trust boundary at runtime. |
| Architecture review (`docs/DECISIONS.md`) | replaced-by:docs/dev.md | Existing ADR covers active design areas (cascade triggers). |

Carve-outs are honored by `dev-update-project`, `dev-arch-review`, and other
iagen-dev skills — they will not be re-suggested unless the user removes the row.
