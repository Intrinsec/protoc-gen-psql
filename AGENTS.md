# AGENTS.md — protoc-gen-psql

Operational contract for AI agents (Claude Code, Copilot, Gemini) and human contributors. Source of truth for tooling, workflow, and out-of-scope items.

- **Tier**: B (shared internal tool, no runtime SLA)
- **Type**: 2 — Go CLI (`protoc` plugin binary)
- **Module**: `github.com/intrinsec/protoc-gen-psql`
- **Purpose**: generate PostgreSQL DDL files from protobuf message definitions
- **Go**: `go 1.25.0` directive, `toolchain go1.26.3`. Entry point `main.go` registers `PSQLify()` against `protoc-gen-star`. Build artifact: `bin/protoc-gen-psql`.

## Language

Default agent response: English, even if user writes French.
French response only when user explicitly asks French.

Caveman compression **mandatory** for all conversational responses (default level: `full`).
Code blocks, commit messages, PR descriptions, security warnings, irreversible-action
confirmations stay normal prose (Caveman auto-clarity rules). Do not disable Caveman unless
user says "stop caveman" or "normal mode".

HARD RULE: all code, comments, identifiers, doc strings, commit messages, ADRs, technical
docs in English — every project type, regardless of team spoken language.
User-facing strings + UI copy exempt — match audience language.

## Workflow Skills (mandatory)

Every agent session in this repo must load + apply these skill packs:

- **superpowers** — process discipline (`brainstorming`, `writing-plans`, `executing-plans`,
  `test-driven-development`, `systematic-debugging`, `verification-before-completion`,
  `requesting-code-review`).
- **caveman** — response compression (see Language section).

Pack missing? Install per iagen-dev `INSTALL.md` before work.

### Session start gate

Before any response, clarification, repository inspection, shell command, or file edit:
run `superpowers:using-superpowers` first, then run `caveman` so compression is active
for every response. Use `superpowers:using-superpowers` to decide which additional
skills apply, then follow the selected skill workflows.

### Plan-writing mandatory before non-trivial implementation

Any feature, refactor, bugfix touching more than one function, or agent cannot reason in
one pass:

1. Run `superpowers:brainstorming` — clarify intent + requirements.
2. Run `superpowers:writing-plans` — persist plan at `docs/superpowers/plans/<short-name>.md`
   (commit to git).
3. Execute via `superpowers:executing-plans` (single-session) or
   `superpowers:subagent-driven-development` (parallelisable steps).
4. Gate completion with `superpowers:verification-before-completion` — no "done" claim
   without evidence (test output, lint output, build output).

**Trivial edits exception:** typos, single-line config tweaks, self-evident one-liners
skip steps 1–3 but still verify before claiming done.

### Bug fixes go through systematic-debugging

Any bug, failing test, unexpected behaviour → `superpowers:systematic-debugging` first.
No symptom patching without root cause.

### Code review before merge

Before merge or PR for non-trivial work: run `superpowers:requesting-code-review`.

## Code Quality

After modifying any Go file: run `golangci-lint run ./...` before marking work complete.
Fix all lint errors, re-run until clean. Lint errors = task not done.
`gofmt` non-negotiable — zero diff allowed. Run `gofmt -w .` if in doubt.

### Project-specific linter config

- Config: `.golangci.yml` (golangci-lint v2 schema, pinned `version: "2"`).
- Enabled linters: `errcheck`, `govet` (`enable-all` minus `fieldalignment`),
  `staticcheck`, `ineffassign`, `unused`, `misspell`. `gofmt` under `formatters`.
- Excluded paths: `psql/psql.pb.go`, `tests/*.pb.go`, `vendor/`.
- Gate: zero issues on `master`. `//nolint:<linter> // reason` allowed with
  written justification only.

## Vulnerability Scanning

After modifying `go.mod` / `go.sum`: run `govulncheck ./...` before marking work complete.
(`vendor/` is honored via `GOFLAGS=-mod=vendor` in env; govulncheck has no `-mod` CLI flag.)
Fix called vulns: `go get <module>@<fixed>`, `go mod tidy`, re-vendor if applicable,
re-run until clean. Imported-only vulns: report to user.
Called vulns remaining = task not done.

CI runs `govulncheck` on every build plus a weekly Monday-06:00-UTC cron against
`master` so newly-published CVEs surface even with no commits.

## Dependency Management

Any `go.mod` change → run `go mod tidy` then `go mod vendor`.
`vendor/` must be committed — never gitignored.
CI uses `go build -mod=vendor`. Never `go get` inside Docker build without re-vendoring after.

**`vendor/` is read-only.** Never edit files under `vendor/` by hand — not to patch a bug,
not to silence a lint warning, not to "just try something". Upstream-only fixes:
`go get <module>@<fixed-version>` + `go mod tidy` + `go mod vendor`. If upstream lacks
a needed fix, fork the module, point `replace` at the fork in `go.mod`, then re-vendor.
Hand-edits to `vendor/` get blown away on the next `go mod vendor` and silently mask
supply-chain provenance.

### Dependency upgrade policy (Renovate)

`renovate.json` schedules dependency MRs Monday 06:00 Europe/Paris. Rules:

- Indirect deps, patch updates → auto-merge (branch automerge).
- Direct deps, major updates → labelled `needs-review`, manual gate.
- Protobuf stack (`protobuf`, `protoc-gen*`) grouped.
- Vulnerability alerts → labelled `security`, manual review.

## Generated Code

Generated sources are **read-only**. Never hand-edit files produced by a code generator:

- `protoc` / `buf` outputs for gRPC + Protobuf (typically `*.pb.go`, `*_grpc.pb.go`,
  often under `gen/`, `pb/`, or `proto/`)
- `mockgen` / `moq` mocks
- `sqlc`, `ent`, `gqlgen`, `wire_gen.go`, `oapi-codegen`, `swag` outputs
- any file with a `// Code generated ... DO NOT EDIT.` header

To change generated code: change the source of truth (`.proto`, `.sql`, schema, interface)
then re-run the generator (`buf generate`, `go generate ./...`, `sqlc generate`, etc.).
Commit the regenerated files alongside the source change in the same commit.

Generated files are committed (not gitignored) so builds and reviews are reproducible.
CI must regenerate and `git diff --exit-code` to catch drift between source + output.

In this project: `psql/psql.pb.go` regenerated from `psql/psql.proto` via `make psql/psql.pb.go`.
`tests/*.pb.go` regenerated by `make test-generate` and diffed against `tests/references/*.pb.psql`.

## Testing & Architecture

Red-Green-Refactor: failing test first, then implementation.
DI via constructors — no package-level globals, no `init()` side effects.
Small, focused interfaces at call site. Never inject concrete type where interface suffices.
Push I/O (DB, HTTP, filesystem) to edges. Domain logic side-effect-free, testable without
external services.

### Project-specific test layout

- **Unit**: `psqlify_test.go` covers pure helper functions (`generateIdentifierName`,
  `allocateRoomToParameters`, `appendSlices`, `getStringBufferWithHeader`,
  `generateCascadeIdentifierNames`, `generateFromTemplate`). Run via
  `go test -mod=vendor ./... -count=1`.
- **Generation**: `make test-generate` invokes plugin against `tests/*.proto`,
  diffs against `tests/references/*.pb.psql`. Empty diff = pass.
- **Integration**: `make test-integration` applies generated SQL to a PostgreSQL
  container defined by `tests/docker-compose.tests.yml`; assertions run in the
  `client` container.

Reference files in `tests/references/` are the regression oracle — update intentionally,
never accept silent diffs.

## Project Layout

Layered layout for non-trivial Go services. `internal/` = compiler-enforced boundary —
packages inside not importable outside module → domain logic private by construction.

```
cmd/<binary>/main.go        # entrypoint + dependency wiring
internal/domain/            # entities, value objects, core interfaces
internal/usecase/           # business logic, orchestrates domain + ports
internal/repository/        # DB / external-API implementations
internal/delivery/http/     # HTTP handlers, DTOs, middleware
pkg/                        # only if code is intentionally exported
```

Two layers (handler + store) OK for small CRUD. Full four-layer split when domain
complexity justifies. No `usecase` passthrough files that forward calls.

Tests next to code (`foo.go` + `foo_test.go`). Cross-package integration tests under
`test/` at module root.

> **Applicability note for this project:** protoc-gen-psql is a single-package
> `main` (visitor + helpers in `psqlify.go`, schema package in `psql/`). The
> layered layout above does not apply at current scope. Reconsider only if the
> plugin grows into multiple commands or shared libraries.

## Dependency Injection

Pick one DI mechanism, keep uniform across service.

| Mechanism | Use when | Trade-off |
|-----------|----------|-----------|
| Manual (explicit constructors in `main.go`) | Default for most services | Verbose when graph grows past ~50 wiring lines |
| Google Wire (compile-time codegen) | `main.go` wiring unreadable or diverges per env | Extra build step, generated code to sync |
| Uber Dig (runtime reflection) | Avoid | Errors surface at runtime, undermines Go compile-time safety |

Default = manual. Switch to Wire only if manual wiring demonstrably unmaintainable.
Do not adopt Dig.

> **Applicability note for this project:** plugin instantiates a single visitor
> (`initPSQLVisitor`) with explicit constructor arguments. No DI framework
> needed at current scope.

## Error Handling

"Crash early, let orchestrator recover" model:
- Transient errors (network, timeout): retry 1–3× exponential backoff, log each retry
  at WARN. Retries exhausted → log ERROR with full context, exit non-zero.
- Structural errors (missing config, unavailable critical dep): crash immediately at
  startup. No retry.
Never swallow errors silently. Every error includes enough context for diagnosis
without accessing running pod.

> **Project-specific:** the plugin is invoked by `protoc` and has no network calls.
> Use `pgs.DebuggerCommon.CheckErr` / `Failf` to surface errors back to protoc.
> Never silently drop the result of `pgs.Walk` or `Extension` calls.

## Local Development

`docker compose up` starts complete local env with all deps (database, cache, mock
services). No manual setup steps required.
Dev compose must not need staging cluster or production secrets.

> **Project-specific:** `tests/docker-compose.tests.yml` brings up an ephemeral
> PostgreSQL plus a client container that applies the generated SQL. Use
> `make test-integration` rather than invoking compose directly. See
> `docs/DEVELOPMENT.md` for the full prerequisite matrix and troubleshooting.

## Logging

`slog` (stdlib) with JSON handler — never `fmt.Println` or `log.Printf`.
Every log entry includes minimum: `timestamp`, `level`, `msg`, `service`.
Logs to stdout only — never files.

> **Project-specific:** the plugin is a short-lived `protoc` subprocess. It uses
> `pgs.DebuggerCommon` (`v.Logf`, `v.Debug`, `v.Failf`) which routes through
> `protoc-gen-star`. New logging from plugin code should go via that interface
> — adding `slog` here would conflict with protoc's stdout protocol.

## Documentation Coherence

After any meaningful change (feature, bugfix touching public behaviour, API surface,
config schema, CLI flags, deps with user impact): verify `README.md` + `docs/**` still
match shipped reality before mark task done. Out-of-sync doc = task not done.

Pre-release sweep (mandatory before every tag, all release levels):
- README accurate — install steps, quickstart, examples runnable as-is.
- All in-repo doc links + references resolve (no dead anchors, no stale paths).
- Public API docs match shipped surface (endpoints, flags, env vars).
- Migration notes present for breaking changes.
- Screenshots / diagrams reflect current UI + architecture.
- `CHANGELOG.md` matches release scope (see Changelog section).

Release blocked if sweep fails. Doc fix = same MR as code change, never separate
follow-up.

## Changelog

Maintain user-friendly `CHANGELOG.md` at repo root. Format: Keep-a-Changelog
(https://keepachangelog.com/en/1.1.0/) + SemVer.

Every user-visible change → entry under `## [Unreleased]` with one type:
`Added` / `Changed` / `Deprecated` / `Removed` / `Fixed` / `Security`.

Wording rules (user-facing, not commit log):
- End-user perspective. No commit hash, no internal module name, no implementation
  detail.
- Bad: "refactor auth middleware to use JWT v2 lib".
- Good: "Sessions survive backend restarts; existing tokens stay valid".
- Breaking changes prefix `**BREAKING:**` + 1-3 line migration note.

Release cut process:
1. Rename `## [Unreleased]` → `## [X.Y.Z] - YYYY-MM-DD`.
2. Create fresh empty `## [Unreleased]` block at top.
3. Tag matches header version exactly (`vX.Y.Z`).
4. Bump compare links at file bottom.

Internal-only changes (refactor with zero user impact, test infra, CI config) skip
CHANGELOG entry — but if in doubt, log under `Changed`.

## CI pipeline

- Primary: GitHub Actions (`.github/workflows/ci.yml`) — repository lives on
  `github.com/Intrinsec/protoc-gen-psql`.
- Mirror: `.gitlab-ci.yml` retained so the project can be moved to or mirrored on
  intrinsec GitLab without re-bootstrap.
- Jobs (both surfaces): golangci-lint, unit-test (coverage), test-generate,
  test-integration, govulncheck, build, sbom (cyclonedx-gomod), license-check
  (go-licenses allow-list).
- Weekly schedule: Monday 06:00 UTC, full pipeline against `master`.
- All Go jobs use `-mod=vendor` for offline-reproducible builds.

## SBOM

Generated on every CI build with `cyclonedx-gomod` → `sbom.cdx.json`. Uploaded as a
30-day CI artifact (`sbom` job).

Rationale: although this is a build-time plugin, it ends up embedded in consumers'
build pipelines, so the dependency closure is part of *their* supply chain. SBOM
lets downstream teams answer "do we ship lib X?" without re-running their own scan.

## License compliance

`go-licenses check ./...` runs in CI against an explicit allow-list:
`Apache-2.0`, `BSD-2-Clause`, `BSD-3-Clause`, `MIT`, `ISC`, `MPL-2.0`.
Job fails on any disallowed license — review and either replace the dep or extend
the allow-list (with a written rationale in the commit message).
License inventory exported as `licenses.csv` artifact (30-day retention).

## Documentation

- `README.md` — user-facing plugin usage and options reference.
- `docs/DEVELOPMENT.md` — build, test, release runbook (replaces the older inline
  Makefile-spelunking dance).
- `docs/dev.md` — design rationale for cascade-update triggers (ADR-style, keep).
- `docs/superpowers/plans/` — iagen-dev plans (precheck + per-gap corrections).
- `CHANGELOG.md` — Keep-a-Changelog history, SemVer.

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
| Secret scanning (gitleaks) | wont-fix | gitleaks v8.x is BSL-1.1; no commercial license available. Rely on code review + pre-commit `check-yaml` / generic hooks. |

Carve-outs are honored by `dev-update-project`, `dev-arch-review`, and other
iagen-dev skills — they will not be re-suggested unless the user removes the row.
