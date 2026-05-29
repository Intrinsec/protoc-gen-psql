# CR: Quiet stderr noise + buf source-relative regression + extension docs — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Silence expected per-message stderr lines (debug-gate), add a buf-based regression test locking source-relative output, document the `tableType`/`column` extensions, and prepare a `v0.1.0` CHANGELOG entry.

**Architecture:** Two informational `v.Logf` calls in `psqlify.go` become `v.Debugf` (gated by the existing `DEBUG_PGV` env). A new self-contained buf fixture tree under `tests/buf/` plus a `test-buf-generate` Makefile target and CI job assert (a) nested source-relative output and (b) clean stderr. README gains a worked example. No production path-derivation change — that already works.

**Tech Stack:** Go 1.25 / `protoc-gen-star/v2`, `protoc` (libprotoc 35), `buf` v1.50.0, Make, GitHub Actions.

**Spec:** `docs/superpowers/specs/2026-05-29-cr-quiet-noise-buf-regression-design.md`

**Branch:** `cr/quiet-noise-buf-regression` (already created; spec already committed).

---

### Task 1: Debug-gate expected/skip stderr lines (Issue 1)

**Files:**
- Modify: `psqlify.go:194` and `psqlify.go:205`

- [ ] **Step 1: Capture current (noisy) behaviour as the failing baseline**

Build and run the plugin against the annotation-free fixture, capturing stderr:

```bash
go build -o bin/protoc-gen-psql .
protoc -I . -I "$(test -d /opt/homebrew/include/google/protobuf && echo /opt/homebrew/include || echo /usr/include)" \
  --plugin=protoc-gen-psql=$(pwd)/bin/protoc-gen-psql --psql_out=/tmp/psql-noise-check \
  tests/no_generation.proto 2>/tmp/psql-stderr.txt; mkdir -p /tmp/psql-noise-check
grep -c "Unable to find an extension tableType" /tmp/psql-stderr.txt
```

Expected NOW: prints `1` (the noisy line for `NoPSQLGeneration`). This is the behaviour we are removing.

- [ ] **Step 2: Change line 205 from `v.Logf` to `v.Debugf`**

In `psqlify.go`, `VisitMessage`, the no-`tableType` branch:

```go
	if !ok {
		v.Debugf("Unable to find an extension tableType equal to DATA or RELATION. Skipping message: %s", m.Name().String())
		return nil, nil
	}
```

- [ ] **Step 3: Change line 194 from `v.Logf` to `v.Debugf`**

In `psqlify.go`, `VisitMessage`, the `disabled` branch:

```go
	if ok, err := m.Extension(psql.E_Disabled, &disabled); ok && err == nil && disabled {
		v.Debugf("Generation disabled for message %s", m.Name().String())
		return nil, nil
	}
```

Leave every other `v.Logf` (error paths) unchanged.

- [ ] **Step 4: Verify silent by default, visible under DEBUG_PGV**

```bash
go build -o bin/protoc-gen-psql .
WKT=$(test -d /opt/homebrew/include/google/protobuf && echo /opt/homebrew/include || echo /usr/include)
mkdir -p /tmp/psql-noise-check
# default: silent
protoc -I . -I "$WKT" --plugin=protoc-gen-psql=$(pwd)/bin/protoc-gen-psql \
  --psql_out=/tmp/psql-noise-check tests/no_generation.proto 2>/tmp/psql-stderr.txt
echo "default count: $(grep -c 'Unable to find an extension tableType' /tmp/psql-stderr.txt)"
# DEBUG_PGV: visible
DEBUG_PGV=1 protoc -I . -I "$WKT" --plugin=protoc-gen-psql=$(pwd)/bin/protoc-gen-psql \
  --psql_out=/tmp/psql-noise-check tests/no_generation.proto 2>/tmp/psql-stderr-dbg.txt
echo "debug count: $(grep -c 'Unable to find an extension tableType' /tmp/psql-stderr-dbg.txt)"
```

Expected: `default count: 0`, `debug count: 1`.

- [ ] **Step 5: Confirm existing generation still passes (no behavioural drift)**

Run: `make test-generate`
Expected: command exits 0, diff section prints nothing (generated == references).

- [ ] **Step 6: Commit**

```bash
git add psqlify.go
git commit -m "fix: debug-gate expected skip/disabled stderr lines

A message without (psql.tableType) is the normal case, not a warning.
Move the per-message skip and disabled lines from Logf to Debugf so they
are silent by default and re-enabled with DEBUG_PGV=1. Closes CR Issue 1.

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 2: Allow buf fixtures under git + ignore generated/staged files

**Files:**
- Modify: `.gitignore`

- [ ] **Step 1: Add ignore rules for the buf tree**

The repo ignores `*.pb.psql` except `tests/references/**`. The new buf fixture needs its
own references whitelisted and its generated output + staged extension copy ignored.
Append to `.gitignore`:

```gitignore
# buf regression fixture (Issue 2): keep references, ignore generated output + staged extension copy
!tests/buf/references/**/*.pb.psql
tests/buf/gen/
tests/buf/proto/psql/
```

- [ ] **Step 2: Verify the ignore logic**

```bash
git check-ignore -v tests/buf/gen/demo/v1/10_tables_widget.pb.psql || echo "NOT ignored (bad)"
mkdir -p tests/buf/references/gen/demo/v1
touch tests/buf/references/gen/demo/v1/10_tables_widget.pb.psql
git check-ignore tests/buf/references/gen/demo/v1/10_tables_widget.pb.psql && echo "ref IGNORED (bad)" || echo "ref tracked (good)"
rm -rf tests/buf
```

Expected: first line shows the gen path is ignored; second prints `ref tracked (good)`.

- [ ] **Step 3: Commit**

```bash
git add .gitignore
git commit -m "chore: gitignore rules for buf regression fixture

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 3: Add buf source-relative regression fixture (Issue 2)

**Files:**
- Create: `tests/buf/buf.yaml`
- Create: `tests/buf/buf.gen.yaml`
- Create: `tests/buf/proto/demo/v1/widget.proto`
- Create: `tests/buf/references/gen/demo/v1/10_tables_widget.pb.psql`

- [ ] **Step 1: Create the buf module config**

`tests/buf/buf.yaml`:

```yaml
version: v2
modules:
  - path: proto
```

- [ ] **Step 2: Create the buf generation config**

`tests/buf/buf.gen.yaml` — invokes the locally built binary, requests source-relative
output into `gen`. The `${PSQL_BIN}` placeholder is substituted by the Makefile (buf reads
`local` as an argv list; the Makefile writes the absolute path in via envsubst-style sed):

```yaml
version: v2
plugins:
  - local: PSQL_BIN_PLACEHOLDER
    out: gen
    opt: paths=source_relative
```

- [ ] **Step 3: Create the annotated nested proto**

`tests/buf/proto/demo/v1/widget.proto` — one DATA table with two columns, plus one plain
annotation-free message to prove stderr stays clean:

```proto
syntax = "proto3";

package demo.v1;

import "psql/psql.proto";

option go_package = "example.com/demo/v1;demov1";

message Widget {
    option (psql.tableType) = DATA;

    string uuid = 1 [
        (psql.column) = "uuid PRIMARY KEY DEFAULT gen_random_uuid()"
    ];

    string label = 2 [
        (psql.column) = "text NOT NULL"
    ];
}

// Annotation-free: must NOT emit any stderr line and must NOT generate output.
message Plain {
    string note = 1;
}
```

- [ ] **Step 4: Create the expected reference output**

`tests/buf/references/gen/demo/v1/10_tables_widget.pb.psql` (note: nested under
`demo/v1/`, proving source-relative; header path is the descriptor name `demo/v1/widget.proto`):

```
-- File: demo/v1/widget.proto
CREATE TABLE IF NOT EXISTS Widget (
	uuid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	label text NOT NULL
);
```

(Indentation inside the table is a single TAB per line, matching `_templates/create_table.tpl.psql`. File ends with a trailing newline after `);`.)

- [ ] **Step 5: Commit the fixture**

```bash
git add tests/buf/buf.yaml tests/buf/buf.gen.yaml tests/buf/proto/demo/v1/widget.proto tests/buf/references/gen/demo/v1/10_tables_widget.pb.psql
git commit -m "test: add buf source-relative regression fixture

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 4: Add `test-buf-generate` Makefile target

**Files:**
- Modify: `Makefile` (the `PROTO_FIXTURES` line ~48, and append a new target near `test-generate`)
- Modify: `Makefile` `test` target (~78) and `clean` target (~83) if present

- [ ] **Step 1: Exclude the buf tree from the protoc fixture list**

In `Makefile`, change the `PROTO_FIXTURES` definition so the buf protos are not fed to the
protoc run (they import via buf module roots, not protoc `-I .`):

```makefile
PROTO_FIXTURES := $(sort $(shell find tests -name '*.proto' -not -path 'tests/buf/*'))
```

- [ ] **Step 2: Add the `test-buf-generate` target**

Append after the `test-generate` target (after current line ~59). It stages a copy of the
extension proto so buf resolves `import "psql/psql.proto"` with no registry, substitutes the
built binary path into a temp gen config, runs `buf generate`, asserts clean stderr, diffs
against references, then cleans up:

```makefile
.PHONY: test-buf-generate
test-buf-generate: build
	@cp -f $(NAME)/$(NAME).proto tests/buf/proto/psql/$(NAME).proto 2>/dev/null || (mkdir -p tests/buf/proto/psql && cp -f $(NAME)/$(NAME).proto tests/buf/proto/psql/$(NAME).proto)
	@sed 's#PSQL_BIN_PLACEHOLDER#$(shell pwd)/bin/protoc-gen-$(NAME)#' tests/buf/buf.gen.yaml > tests/buf/buf.gen.local.yaml
	@cd tests/buf && rm -rf gen && buf generate proto --template buf.gen.local.yaml 2>buf.stderr.txt; status=$$?; \
		cat buf.stderr.txt; \
		if grep -q "Unable to find an extension tableType" buf.stderr.txt; then \
			echo "FAIL: skip-message noise leaked to stderr (Issue 1 regression)"; exit 1; \
		fi; \
		test $$status -eq 0 || { echo "FAIL: buf generate exited $$status"; exit $$status; }
	# Checking diff between generated file and reference file (empty == identical)
	@for i in `find tests/buf/gen -name '*.pb.psql' | sort`; do \
		rel=$${i#tests/buf/}; \
		diff tests/buf/references/$$rel $$i || { echo "FAIL: $$i differs from reference (source-relative regression)"; exit 1; }; \
	done
	@test -f tests/buf/gen/demo/v1/10_tables_widget.pb.psql || { echo "FAIL: expected nested source-relative output missing"; exit 1; }
	@rm -f tests/buf/buf.gen.local.yaml tests/buf/buf.stderr.txt
	@echo "test-buf-generate OK"
```

- [ ] **Step 3: Wire into the `test` aggregate and `clean`**

Update the `test` target to include the new target:

```makefile
test: build test-generate test-buf-generate test-integration
```

If a `clean` target removes generated psql, add the buf gen dir + temp files:

```makefile
	@rm -rf tests/buf/gen tests/buf/buf.gen.local.yaml tests/buf/buf.stderr.txt tests/buf/proto/psql
```

- [ ] **Step 4: Run the new target (the real verification)**

Run: `make test-buf-generate`
Expected: ends with `test-buf-generate OK`, exit 0. No `Unable to find` line printed. The
diff loop prints nothing.

- [ ] **Step 5: Negative check — prove the test actually guards source-relative**

Temporarily break the reference path to confirm the diff fails loudly, then restore:

```bash
mv tests/buf/references/gen/demo/v1/10_tables_widget.pb.psql /tmp/ref-backup.psql
make test-buf-generate; echo "exit=$?"   # expect non-zero (FAIL message)
mv /tmp/ref-backup.psql tests/buf/references/gen/demo/v1/10_tables_widget.pb.psql
make test-buf-generate                    # expect OK again
```

Expected: first run exits non-zero with a FAIL message; restored run prints `test-buf-generate OK`.

- [ ] **Step 6: Confirm protoc fixtures still pass after the PROTO_FIXTURES change**

Run: `make test-generate`
Expected: exit 0, diff section empty.

- [ ] **Step 7: Commit**

```bash
git add Makefile
git commit -m "test: add test-buf-generate target locking source-relative + clean stderr

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 5: Add buf CI job

**Files:**
- Modify: `.github/workflows/ci.yml` (add job after `test-generate` ~line 67; extend release `needs` ~line 160)

- [ ] **Step 1: Add the `test-buf-generate` job**

Insert after the `test-generate` job (mirror its shape; also install buf, pinned):

```yaml
  test-buf-generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - name: Install buf
        run: go install github.com/bufbuild/buf/cmd/buf@v1.50.0
      - name: make test-buf-generate
        run: make test-buf-generate
```

- [ ] **Step 2: Add the job to the release gate**

Find the release job's `needs:` list (currently `[lint, unit-test, test-generate, test-integration, govulncheck, build, sbom, ...]` near line 160) and add `test-buf-generate` to it.

- [ ] **Step 3: Validate the workflow YAML locally**

Run:
```bash
python3 -c "import yaml,sys; yaml.safe_load(open('.github/workflows/ci.yml')); print('yaml ok')"
```
Expected: `yaml ok`.

- [ ] **Step 4: Commit**

```bash
git add .github/workflows/ci.yml
git commit -m "ci: run test-buf-generate; gate release on it

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 6: Document the `tableType` / `column` extensions (Issue 3)

**Files:**
- Modify: `README.md` (after the "## How to use it" section, ~line 9)

- [ ] **Step 1: Add a worked-example subsection**

Insert after the "## How to use it" paragraph:

```markdown
### Minimal annotated example

A message becomes a PostgreSQL table only when it carries the `(psql.tableType)`
extension. Import the extension definitions, mark the message, and annotate columns:

```proto
import "psql/psql.proto";

message Tenant {
  option (psql.tableType) = DATA;                                  // -> 10_tables_*.pb.psql

  string id   = 1 [(psql.column) = "uuid PRIMARY KEY DEFAULT gen_random_uuid()"];
  string name = 2 [(psql.column) = "text NOT NULL"];
}
```

- `option (psql.tableType) = DATA;` marks the message as a data table; its columns land in
  `10_tables_<file>.pb.psql`. Use `RELATION` for relation tables (`20_relations_<file>.pb.psql`).
- Each field's `(psql.column)` string is emitted verbatim after the column name.
- Messages **without** `(psql.tableType)` are skipped silently (set `DEBUG_PGV=1` to log skips).
```

- [ ] **Step 2: Verify it renders (no broken fences)**

Run: `grep -n "Minimal annotated example" README.md`
Expected: one match. Eyeball the fenced block boundaries.

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "docs: document tableType/column extensions with a worked example

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 7: CHANGELOG entry for v0.1.0

**Files:**
- Modify: `CHANGELOG.md`

- [ ] **Step 1: Inspect the existing format**

Run: `head -30 CHANGELOG.md`
Expected: see the heading style (Keep-a-Changelog or similar) to match.

- [ ] **Step 2: Add the `v0.1.0` entry**

Following the file's existing format, add an entry above the latest release. Content:

```markdown
## [0.1.0]

### Changed
- Per-message "skip" and "generation disabled" lines are now debug-only
  (silent by default; set `DEBUG_PGV=1` to see them). Fixes verbose `buf generate`
  output on annotation-free messages.

### Added
- buf-based regression test (`make test-buf-generate`) locking source-relative
  output (`paths=source_relative`) and asserting clean stderr.
- README section documenting the `tableType` / `column` extensions.

### Notes
- `paths=source_relative` was already honoured since the `protoc-gen-star/v2`
  migration; the new test guards against regression.
```

(Match the exact bracket/date convention already used in the file — adjust heading to fit.)

- [ ] **Step 3: Commit**

```bash
git add CHANGELOG.md
git commit -m "docs: changelog entry for v0.1.0

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 8: Full verification gate

**Files:** none (verification only)

- [ ] **Step 1: Lint**

Run: `golangci-lint run` (or `make lint` if defined)
Expected: exit 0, no findings.

- [ ] **Step 2: Unit tests**

Run: `go test -mod=vendor ./... -count=1`
Expected: `ok` for all packages.

- [ ] **Step 3: Generation tests (protoc + buf)**

Run: `make test-generate && make test-buf-generate`
Expected: both exit 0; `test-buf-generate OK` printed; no diffs.

- [ ] **Step 4: Confirm no stray generated/temp files left in the tree**

Run: `git status --porcelain`
Expected: empty (all work committed; gen + temp files gitignored).

- [ ] **Step 5: Report**

Summarise captured output for each command above. No "done" claim without this evidence.
Tagging `v0.1.0` is a separate human action via the release workflow — not part of this plan.

---

## Self-Review

**Spec coverage:**
- Issue 1 (debug-gate) → Task 1. ✓
- Issue 2 (buf regression test) → Tasks 2–5 (gitignore, fixture, Makefile target, CI). ✓
- Issue 3 (README docs) → Task 6. ✓
- Release v0.1.0 CHANGELOG → Task 7. ✓
- Verification gate → Task 8. ✓
- Out-of-scope (no path-logic change, no pgs bump, no summary mode) → respected. ✓

**Placeholder scan:** `PSQL_BIN_PLACEHOLDER` is an intentional, defined token substituted by
the Makefile sed in Task 4 Step 2 — not an unfilled placeholder. No TBD/TODO remain.

**Type/name consistency:** Target name `test-buf-generate` is identical across Makefile
(Task 4), CI job + release `needs` (Task 5), and verification (Task 8). Fixture paths
(`tests/buf/...`, `demo/v1/widget.proto`, reference at
`tests/buf/references/gen/demo/v1/10_tables_widget.pb.psql`) are identical across Tasks 3, 4,
6 example, and ignore rules in Task 2. Env var `DEBUG_PGV` matches `main.go`.
