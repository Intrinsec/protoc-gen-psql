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

- Primary: GitHub Actions (`.github/workflows/ci.yml`) -- the repository lives on `github.com/Intrinsec/protoc-gen-psql`.
- Mirror: `.gitlab-ci.yml` retained so the project can be moved to or mirrored on intrinsec GitLab without re-bootstrap.
- Jobs (both surfaces): golangci-lint, gitleaks, unit-test (with coverage), test-generate, test-integration, govulncheck, build.
- All Go jobs use `-mod=vendor` to keep builds offline-reproducible.

## Dependency policy

- Renovate (`renovate.json`) with grouped Go module updates, monthly schedule, auto-merge on patch level for indirect deps.
- Tracked in `docs/superpowers/plans/2026-05-20-onboard-dep-policy.md`.

## Secret scanning

- Pre-commit + CI gate via `gitleaks`.
- Baseline scan + ignore list reviewed quarterly.
- Tracked in `docs/superpowers/plans/2026-05-20-onboard-secret-scanning.md`.

## SBOM

- Generated on every CI build with `cyclonedx-gomod` -> `sbom.cdx.json`.
- Uploaded as an artifact by the CI `sbom` job, retained 30 days.
- Rationale: although this is a build-time plugin, it ends up embedded in
  consumers' build pipelines, so the dependency closure is part of *their*
  supply chain. SBOM lets downstream teams answer "do we ship lib X?"
  without re-running their own scan.

## License compliance

- `go-licenses check ./...` runs in CI against an explicit allow-list:
  `Apache-2.0`, `BSD-2-Clause`, `BSD-3-Clause`, `MIT`, `ISC`, `MPL-2.0`.
- Job fails on any disallowed license -- review and either replace the
  dep or extend the allow-list (with a written rationale in the commit
  message).
- License inventory exported as `licenses.csv` artifact (30-day retention).

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
| Threat model | not-applicable | Pure code generator, no trust boundary at runtime. |
| Architecture review (`docs/DECISIONS.md`) | replaced-by:docs/dev.md | Existing ADR covers active design areas (cascade triggers). |

Carve-outs are honored by `dev-update-project`, `dev-arch-review`, and other
iagen-dev skills — they will not be re-suggested unless the user removes the row.
