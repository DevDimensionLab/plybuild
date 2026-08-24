# Quality Upgrade Plan

Last measured checkpoint: 2026-08-24, commit `312d168`.

## Objective

Reach an honest L2 quality level while preserving the CLI and public Go API.
"Honest" means held baseline debt is worked down rather than accepted merely
because the ratchet verdict says PASS.

## Operating Model

- Work only in `/Users/perottochristensen/github/ply/upgrade-quality` on
  `codex/upgrade-quality`.
- Start fresh agent sessions with `./codex-dev-start.sh`; its mutable prompt is
  the sole next-task source and the tracked plan is persistent knowledge.
- Make one focused commit per measured quality move.
- Limit a checkpoint to three moves, then update this plan and the handover.
- At a checkpoint or session boundary, automatically prepare and locally
  commit the next authorized session prompt before stopping.
- Run the relevant focused tests before the full gate.
- Measure every checkpoint from a clean commit with the stored baseline.
- Require zero ratchet regressions and resolve actionable review findings.
- Do not push, merge, or remove the worktree without explicit approval.
- Follow `docs/design/agent-session-continuity.md` for task authority, prompt
  archives, restart staging, and dirty-worktree recovery.

The user has authorized the ordered roadmap through P8. The following block is
machine-readable launcher state; keep the order and vocabulary exact.

<!-- CODEX_AUTHORIZED_CHECKPOINTS_BEGIN -->
P2A|complete
P2B|complete
P3|active
P4|queued
P5|queued
P6|queued
P7|queued
P8|queued
<!-- CODEX_AUTHORIZED_CHECKPOINTS_END -->

## Measured State

| Signal | Baseline | P1B | P2A | P2B | Interpretation |
| --- | ---: | ---: | ---: | ---: | --- |
| Absolute L0 PASS | 2 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | Q0.3 and Q0.8 remain closed. |
| Test functions | 33 | 37 | 37 | 37 | Distribution uses an external executable contract. |
| Skipped tests | 2 | 0 | 0 | 0 | Q0.6 remains improved. |
| Packages with tests | 5 / 20 | 7 / 22 | 7 / 22 | 7 / 22 | Fifteen packages still have no test files. |
| Process-exiting calls outside `main` | 127 | 127 | 127 | 127 | P2B did not enter P3 process-boundary scope. |
| Direct external effects outside adapters | not trustworthy in upstream scan | 80 / 80 | 80 / 80 | 80 / 80 | Q1.3 fails; all five adapter paths are absent. |
| Declared seam swap tests | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | Formal ratchet PASS, debt unchanged. |
| Mutation harnesses | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | Formal ratchet PASS, debt unchanged. |
| Acceptance scripts | 0 / 4 | 0 / 4 | 4 / 4 | 4 / 4 | Q2.5-Q2.7 and Q2.10 pass; Q2.8/Q2.9 remain manual. |
| Production scripts without meta-tests | 1 / 1 | 0 / 1 | 0 / 7 | 0 / 7 | Q0.8 remains improved over a larger population. |
| Reachable manual L1/L2 rows | 0 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | Synthetic non-empty fixtures still prove schema reachability. |
| Baseline numeric debt leaves | 228 | 228 | 228 | 228 | Distribution work did not migrate the instrument. |

Authoritative report: `target/quality-audit/scorecard.json`.

## Checkpoints

### P0 - Recovery And Measurement

Status: complete.

- Created the feature worktree and repaired `make install`.
- Established the pinned, fail-closed quality apparatus and baseline.
- Moved config, Maven, and template fixtures out of the repository tree.
- Replaced both skipped Kibana tests with deterministic transport tests.
- Finished with one improved, six held, and zero regressed ratchets.

Commits: `3822f4a`, `674e0a4`, `80b43ba`, `25689c5`.

### P0A - Establish Restartable Agent Sessions

Status: complete as an operational prerequisite for P1.

- Added a contract-tested interactive Codex launcher with one mutable mission,
  inert tail data, a pinned stable-skeleton digest, and Bash 3.2 coverage.
- Made the launcher prompt and rolling handover the startup path, replacing the
  ignored local task file as an active source.
- Added byte-exact, dated prompt archives, a connected reciprocal archive graph,
  and an automatic handoff protocol.
- Kept dirty work recoverable without injecting raw Git output into a prompt.
- Kept the launcher contract portable to Docker/source archives by testing the
  checked-in graph inside synthetic Git state and asserting its source types.
- Forced the fresh session to the normal Codex service tier while retaining the
  rest of the user's local profile.

This apparatus does not consume one of P1's three measured quality moves. Its
contract test runs in `make test` and every later checkpoint gate.

Commit: `3e80994`.

### P1 - Close L0 And Establish The Daily Gate

Status: complete. The three-move limit was observed.

- Pinned the official golangci-lint binary contract at `v2.12.2`, retained Go
  1.18, split read-only `make lint` from `make format`, and brought the selected
  five-linter plus `gofmt` gate to zero issues.
- Added an independent negative meta-test for `scripts/search-replace.sh`; it
  exercises selection and replacement, exposes pipeline failures, and kills a
  no-op mutant.
- Added the non-publishing `make preflight` gate and a contract test that rejects
  missing, empty, incomplete, orphaned, and symlinked script populations.

Commits: `2972d11`, `e956da7`, `ab9a5c6`.

Exit: L0 is 8 of 8, scripts without meta-tests equals zero, the tree is
unchanged by the gate, and no ratchet regresses. Verified at `ab9a5c6`: two
ratchets improved, five held, and zero regressed.

### P1B - Make Manual L1/L2 Evidence Reachable

Status: complete in exactly three measured moves; required before claiming L1
or L2.

- Added schema-2 criterion receipts for Q1.6, Q1.7, Q1.9, Q2.4, Q2.8,
  and Q2.9. Each receipt is module-, commit-, measured-tree-, inventory-,
  instrument-, criterion-, and canonical-evidence-digest-bound.
- Added criterion-specific semantic checks and negative meta-tests for stale,
  duplicate, empty, wrong-kind, dirty-tree, wrong-digest, and false-PASS
  evidence. Synthetic non-empty test, mutation-harness, and acceptance-script
  populations prove all six rows are reachable.
- Proved automated verdict precedence independently. A receipt can modify only
  upstream `UNMEASURABLE`; automated PASS/FAIL remains authoritative, and a
  zero independently measured test/harness/script population cannot pass.
- Reproduced source commit `5635d50` with the schema-1 and schema-2
  instruments. The raw bodies, complete criterion objects, complete
  denominators, and all 228 numeric debt leaves match; Q3.9 remains PASS.
  `.quality/baseline/instrument-migration.json` records both source commits,
  complete instrument identities, evidence hashes, and scorecard hashes.

Commits: `99cebaa`, `4887222`, `eb987fd`.

Exit: satisfied. No current-project PASS was claimed for Q2.4, Q2.8, or Q2.9;
their measured mutation-harness and acceptance-script populations remain zero.

### P2A - Characterize Compatibility

Status: complete in three measured moves, plus one gate correction that removed
an accidental process-exit ratchet regression from the build-ignored exporter.

Move 1 pins `golang.org/x/exp/cmd/apidiff` at
`v0.0.0-20260709172345-9ea1abe57597`. The machine-readable v1.0.1 comparison
allows only the compatible addition of `cmd/ply`; no incompatible change is
allowed.

Move 2 stores a normalized, order-independent Cobra tree for `v1.0.1` and
requires an explicit delta allowlist. The current tree has zero deltas. Fresh
subprocess contracts preserve root, status, upgrade, and build help plus the
existing unknown-command stdout, stderr, and exit-1 behavior.

Move 3 adds four host acceptance verifiers for fresh install, read-only status,
targeted dependency upgrade, and local build. They use isolated homes/caches,
checked-in non-empty fixtures, and loopback Maven metadata. Their 26 independent
meta-controls kill no-op artifacts, wrong help/output, broad or read-only
writes, and bad-input exit/artefact regressions.

Commits: `224cbcc`, `06b2a24`, `5f4ba29`; gate correction `14764fd`.

1. Pin an API-diff tool and record the exported Go API against tag `v1.0.1`.
   Maintain an explicit compatibility allowlist rather than reviewing raw text.
2. Export the Cobra command/flag tree into normalized, order-independent data;
   add output and exit contracts for install, status, upgrade, and build.
3. Add early acceptance skeletons that run fresh artifacts against local
   fixture, cache, and loopback inputs. They must fail when the artifact or
   asserted behavior is absent; full Binary+Docker coverage finishes in P6.
`cmd.Execute()` is an explicit compatibility decision point: its exported
symbol and signature can remain, but its process-exiting behavior conflicts
with exit-only-main. Characterize external use first. If callers require that
side effect, stop for an approved migration rather than hiding an exit behind
an adapter or function variable.

Exit: public API and CLI deltas are machine-readable, core surface contracts
are executable before refactoring, and all three moves are measured. The clean
audit at `14764fd` exits 1 only for 16 documented findings, with two improved,
five held, and zero regressed ratchets.

### P2B - Make Distribution Non-Publishing By Default

Status: complete in one measured move.

The default v2 GoReleaser configuration now disables remote releases and has no
Homebrew or Snap sections. The standalone brew configuration was removed,
`make release-brew` fails closed before any executable is reached, and both
`make snapshot` and the ordinary `make release` target use the same
credential-cleared `release --snapshot --clean --skip=publish` invocation.

`test/makefile_distribution_test.sh` records the exact argv without executing
GoReleaser, rejects Homebrew formula, Homebrew cask, and Snap publisher mutants,
requires the remote-release disable, verifies credential isolation, and proves
the inactive brew target cannot reach the recorder. Preflight runs this
contract. No publisher or GoReleaser release was invoked.

Commit: `312d168`.

Exit: satisfied. The clean audit at `312d168` exits 1 for the same 16 documented
findings, with two improved, five held, and zero regressed ratchets. Actual
snapshot artifact execution remains P6 acceptance work.

### P3 - Remove Process Exits, Then Introduce Seams

Status: queued.

First reduce Q1.2. Add an error-returning execution path and move process exit
to the main boundary. Both `main.go` and `cmd/ply/main.go` exist and are used by
different build paths; contract-test that both delegate to the same error/exit
path, or explicitly retire one without changing produced artifacts. Preserve
exported `cmd.Execute()` and `cmd.RootCmd` symbols/signatures through the P2
compatibility decision. Only after this boundary is green, introduce thin
adapters one at a time using `.quality/inventory`:

1. `internal/adapter/process` for subprocess execution, not application exit.
2. `internal/adapter/httpclient` for HTTP request execution.
3. `internal/adapter/filesystem` for production filesystem mutation.

For each adapter, first add a characterization or argument-swap test, then move
one coherent flow. Defaults must be safe and recording doubles must preserve
the complete dependency struct. Start with the declared git, Maven, and cloud
seams because their argument order can cause destructive behavior.

Exit: process exits outside `main` reach zero, migrated call sites disappear
from Q1.3, their seam swaps are killed, and CLI/API contracts stay compatible.

### P4 - Finish Absolute L1

Status: queued.

- Add `internal/adapter/clock` and `internal/adapter/server`.
- Cover all eight declared seams with executable swap tests.
- Reduce logic packages without tests to zero. A pure type/constant package may
  be excluded only through an explicit, audited N/A rule.
- Record evidence for default doubles, partial-failure content, non-empty
  populations, and state-leak checks.

Exit: untested logic packages equals zero, process exits outside `main` equals
zero, direct effects outside adapters equals zero, direct time calls outside
`internal/adapter/clock` equal zero, a movable fake clock proves time control,
all eight seam swaps are covered, and every manual L1 row has valid non-empty
evidence.

### P5 - Build L2 Mutation Evidence

Status: queued.

Implement the eight named harnesses from `.quality/inventory`, one subject per
measured move and no more than three moves per checkpoint. Each harness must
declare its mutations, prove it can fail, include the methodology T1-T10
meta-controls, and report declared versus killed mutations. Surviving mutations
must be classified as reachability, observability, or controllability gaps and
drive a test or seam improvement.

Exit: all eight harnesses declare at least eight meaningful mutations, pass
T1-T10, and report `declared == killed`, `survived == 0`, and `unusable == 0`.

### P6 - Build L2 Acceptance Evidence

Status: queued.

Extend the P2A host acceptance scripts for `install`, `status`, `upgrade`, and
`build` through both a fresh GoReleaser snapshot binary and a fresh Docker
image. Host install exercises `make install` / `go install`; it is not the
`ply install` command group. Acceptance continues to use local fixtures, cache,
and loopback services so it measures the artifact without public
infrastructure. Real-boundary smokes are separately labeled. Homebrew and Snap
remain outside the active distribution matrix.

Add `make quality` only when it runs preflight, mutation meta/harnesses,
acceptance scripts, and the authoritative audit scoped to
`--only Q0.*,Q1.*,Q2.*`. The target succeeds only on audit exit 0 and rejects a
missing criterion population. A separate full report may exit 1 for L3 debt.

Exit: Q2.5-Q2.10 have executable evidence for all four core flows and a clean
`--only Q0.*,Q1.*,Q2.*` audit attains L2 without held material debt.

### P7 - Adopt A Maintained Go And Dependency Baseline

Status: queued after honest L2.

- Select and document Go 1.26 or 1.27 based on supported stable tooling at
  execution time; update the module declaration and CI together.
- Upgrade dependencies in small groups, with `go mod tidy`, build, tests, race,
  vet/lint, API/CLI diff, acceptance, and vulnerability scanning after each.
- Keep dependency-only commits separate from behavior changes.
- Migrate the exact-toolchain baseline by reproducing old and new measurements;
  preserve every debt value and record both instrument identities.

Exit: the declared toolchain matches the verified toolchain, dependency
upgrades have no unexplained output or API drift, and vulnerability findings
are resolved or explicitly risk-accepted.

### P8 - Domain Modernization

Status: queued after the core L2 flows.

- Migrate cloud configuration toward `ply-config` while retaining cache-first
  behavior and compatibility fixtures.
- Repair and modernize Spring behavior under dedicated characterization tests.
- Revisit inactive packaging only through a separate scope decision.

## Gate For Every Checkpoint

```sh
make preflight
make test
make test-install
make test-agent-start
go test ./... -count=1
go test -race ./... -count=1
go vet ./...
bash .quality/tools/test-quality-audit.sh
bash .quality/tools/quality-audit.sh . \
  --baseline .quality/baseline/scorecard.json
git status --short
```

Run the count-2 hermetic test with an empty HOME and isolated writable state:

```sh
root="$(mktemp -d /private/tmp/ply-hermetic.XXXXXX)"
mkdir -p "$root/home" "$root/cache" "$root/tmp" "$root/config"
modcache="$(go env GOMODCACHE)"
HOME="$root/home" XDG_CONFIG_HOME="$root/config" \
  GOCACHE="$root/cache" GOTMPDIR="$root/tmp" GOMODCACHE="$modcache" \
  GOENV=off GOWORK=off go test ./... -count=2
```

The full audit exits 1 while measured project findings remain. Exit 2 means the
audit or its evidence is broken and invalidates the checkpoint.

Before treating the gate as clean, require:

```sh
test -z "$(git status --porcelain=v1 --untracked-files=all)"
```

The ignored `.agent-task/` path is retained only for compatibility and is not a
startup source. Keep it absent during measurement because ignored files remain
visible to the audit's filesystem identity. The launcher and prompt archive are
tracked and contract-tested. A clean shared clone at the exact commit remains
the preferred checkpoint measurement.

## Risk Register

| Risk | Control |
| --- | --- |
| CLI or API drift during refactoring | Characterization and acceptance before moving a boundary. |
| Real network behavior differs from doubles | Label evidence precisely and add real-boundary smokes where feasible. |
| Docker assumptions remain untested | Require a daemon-backed build and command smoke before claiming support. |
| Large toolchain/dependency jump obscures failures | Upgrade in isolated, dependency-only commits. |
| Ratchet PASS hides unchanged debt | Track the underlying number in this plan and require reduction by L1/L2 exit. |
| Tests mutate fixtures or developer state | Central safe writers, empty-HOME runs, and tree identity checks. |
| Agent sessions lose or duplicate task context | One mutable launcher mission, one rolling handover, and one connected reciprocal archive graph. |
| Prompt mutation changes executable shell behavior | Comment-encoded mutable data after a stable `exit`, plus a pinned normalized skeleton and command-injection probes. |
