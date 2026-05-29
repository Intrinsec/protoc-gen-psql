# Design — CR: quiet stderr noise, lock buf source-relative, document extensions

**Date:** 2026-05-29
**Reporter:** SIRP project (Intrinsec / ia-gen-lab)
**Target release:** `v0.1.0`

## Context

A change request from the SIRP project flagged three items against `v0.0.13`:

1. **Verbose stderr** — every annotation-free message prints a `[psql]` info line.
2. **`paths=source_relative` allegedly ignored** — flat output instead of nested.
3. **Missing docs** — no README example of the `tableType` + `column` extensions.

### Investigation findings

- **Issue 1 reproduces.** `psqlify.go:205` uses `v.Logf` (always written to stderr) for
  messages lacking the `tableType` extension. It also fires for the imported
  `psql/psql.proto`'s own messages.
- **Issue 2 does NOT reproduce on master.** `v0.0.13` is already on
  `protoc-gen-star/v2` (commit `16dd248`). The module derives its output path from
  `f.InputPath()` (the proto descriptor name), so output is inherently source-relative.
  Verified with both `protoc -I proto` and real `buf generate`: output nests correctly at
  `gen/.../<pkg>/<v>/10_tables_<file>.pb.psql`. The reporter likely hit an older
  pre-v2 `@latest`. No code change required — only a regression test + release.
- **Issue 3 is a real gap.** README documents options but has no end-to-end annotated
  message example.

## Scope

### 1. Issue 1 — debug-gate expected/skip lines (silent default)

In `psqlify.go`, switch two **informational** log lines from `v.Logf` to `v.Debugf`:

- `VisitMessage` L194: `Generation disabled for message %s`
- `VisitMessage` L205: `Unable to find an extension tableType equal to DATA or RELATION. Skipping message: %s`

Both describe the *normal* case (a message that is not a psql table), not an error.
`v.Debugf` is gated by `pgs.DebugEnv("DEBUG_PGV")` (already wired in `main.go`), so output
is silent by default and re-enabled with `DEBUG_PGV=1`.

**Leave untouched** every genuine error `Logf` (extension-retrieval failures, write
failures, the `tableType` decode error at L200). Those are real signals.

Net effect: annotation-free messages and imported-extension messages produce zero stderr.

### 2. Issue 2 — buf source-relative regression test

Lock the working behaviour against regression. No production code change.

New fixture tree `tests/buf/`:

- `tests/buf/buf.yaml` — v2, single local module rooted at `proto`.
- `tests/buf/buf.gen.yaml` — v2, one `local` plugin (the built binary), `out: gen`,
  `opt: paths=source_relative`.
- `tests/buf/proto/demo/v1/widget.proto` — a nested, annotated proto (`tableType = DATA`,
  one `column`) plus one plain annotation-free message (to assert silence).
- `tests/buf/proto/psql/` — copy of the `psql` extension protos so buf resolves the import
  with no registry dependency.
- `tests/buf/references/gen/demo/v1/10_tables_widget.pb.psql` — expected nested output.

New Makefile target `test-buf-generate`:

1. `build` the plugin.
2. Run `buf generate` from `tests/buf`, capturing stderr.
3. Assert stderr contains **no** `Unable to find` line (covers Issue 1 regression too).
4. `diff` generated tree against `tests/buf/references/` (asserts the nested
   `demo/v1/` path — proves source-relative).
5. Clean generated output.

**Fix existing `test-generate`:** the buf fixtures must not leak into the protoc run.
Change `PROTO_FIXTURES` to exclude the buf tree:

```makefile
PROTO_FIXTURES := $(sort $(shell find tests -name '*.proto' -not -path 'tests/buf/*'))
```

Keep the existing protoc multi-dir fixture (`tests/sub/`) as-is.

CI (`.github/workflows/ci.yml`): add a `test-buf-generate` job mirroring `test-generate`,
installing buf via `go install github.com/bufbuild/buf/cmd/buf@<pinned-version>`, and add it
to the release job's `needs` gate list.

### 3. Issue 3 — README extensions example

Add a worked-example section near "How to use it" showing:

```proto
import "psql/psql.proto";

message Tenant {
  option (psql.tableType) = DATA;
  string id = 1 [(psql.column) = "uuid primary key"];
}
```

with a one-line note on what each annotation does and which output file it lands in.

### 4. Release

Add a `v0.1.0` entry to `CHANGELOG.md` (quiet stderr, buf regression coverage, docs).
Tagging the release is a human action via the existing release workflow — this work does
**not** push tags.

## Out of scope

- No change to path-derivation logic (already correct).
- No summary-line logging mode (debug-gate chosen instead).
- No bump of `protoc-gen-star` (already v2).

## Testing

- `make test-generate` — existing protoc fixtures still pass (buf tree excluded).
- `make test-buf-generate` — new: nested output + clean stderr.
- Manual: `DEBUG_PGV=1 buf generate` still shows the skip lines (debug path intact).

## Verification gate

No "done" claim without captured output of `make test-generate`, `make test-buf-generate`,
`golangci-lint`, and `go test ./...`.
