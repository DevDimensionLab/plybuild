# Quality Upgrade Plan

Last measured checkpoint: 2026-08-25, commit `7754575`.

## Objective

Reach an honest L2 quality level while preserving the CLI and public Go API.
"Honest" means held baseline debt is worked down rather than accepted merely
because the ratchet verdict says PASS.

## Operating Model

- Work only in `/Users/perottochristensen/github/ply/upgrade-quality` on
  `codex/upgrade-quality`.
- Start fresh agent sessions with `./codex-dev-start.sh`; its supervisor runs
  one non-interactive archived mission at a time, and its mutable prompt is the
  sole next-task source while the tracked plan is persistent knowledge.
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

| Signal | Baseline | P1B | P2A | P2B | P3.1 | P3.2 | P3.3 | P3.4 | P3.5 | P3.6 | P3.7 | P3.8 | P3.9 | P3.10 | P3.11 | P3.12 | P3.13 | P3.14 | P3.15 | P3.16 | P3.17 | P3.18 | P3.19 | P3.20 | P3.21 | P3.22 | P3.23 | P3.24 | P3.25 | P3.26 | P3.27 | P3.28 | P3.29 | P3.30 | Interpretation |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| Absolute L0 PASS | 2 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | 8 / 8 | Q0.3, Q0.6, and Q0.8 remain closed. |
| Test functions | 33 | 37 | 37 | 37 | 39 | 47 | 50 | 55 | 65 | 79 | 92 | 101 | 112 | 123 | 130 | 136 | 140 | 144 | 148 | 153 | 157 | 164 | 169 | 175 | 180 | 187 | 194 | 202 | 209 | 216 | 224 | 230 | 236 | 243 | Templates walk dependency selection, exact receiver-derived root, delivered callback order, paths, metadata and errors, matching, relative names, loaded projects, partial results, exact errors, nil results, safe defaults, and recorded populations are covered. |
| Skipped tests | 2 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | Q0.6 remains improved. |
| Packages with tests | 5 / 20 | 7 / 22 | 7 / 22 | 7 / 22 | 9 / 22 | 11 / 23 | 11 / 23 | 11 / 23 | 13 / 24 | 14 / 25 | 15 / 25 | 15 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | 16 / 25 | Nine packages have no test files; the new contracts stay within already-tested packages. |
| Process-exiting calls outside `main` | 127 | 127 | 127 | 127 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | Q1.2 stays closed; only the two executable `main` functions terminate. |
| Direct external effects outside adapters | not trustworthy in upstream scan | 80 / 80 | 80 / 80 | 80 / 80 | 80 / 80 | 77 / 79 | 76 / 78 | 76 / 78 | 74 / 79 | 71 / 80 | 67 / 77 | 66 / 76 | 64 / 74 | 62 / 72 | 58 / 68 | 57 / 67 | 56 / 66 | 55 / 65 | 54 / 64 | 52 / 62 | 50 / 60 | 49 / 60 | 48 / 60 | 47 / 60 | 46 / 60 | 44 / 59 | 43 / 58 | 42 / 58 | 41 / 57 | 40 / 56 | 39 / 56 | 38 / 55 | 37 / 54 | 36 / 53 | `GitCloudConfig.Templates`' direct tree walk moved into the declared filesystem adapter with no new production effect site. |
| Declared seam swap tests | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 2 / 8 | 3 / 8 | 4 / 8 | 5 / 8 | 6 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | 7 / 8 | Templates recording contracts do not claim a new inventory seam. |
| Executable mutation harnesses | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | 0 / 8 | All five P3 seam drivers remain non-executable and do not qualify as later P5 harnesses. |
| Acceptance scripts | 0 / 4 | 0 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | 4 / 4 | All four host flows pass through fresh artifacts. |
| Production scripts without meta-tests | 1 / 1 | 0 / 1 | 0 / 7 | 0 / 7 | 0 / 7 | 0 / 8 | 0 / 9 | 0 / 10 | 0 / 10 | 0 / 11 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | 0 / 12 | Q0.8 remains improved over the same population. |
| Reachable manual L1/L2 rows | 0 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | 6 / 6 | Synthetic non-empty fixtures prove schema reachability. |
| Baseline numeric debt leaves | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | 228 | The measurement instrument is unchanged. |

Authoritative report: `target/quality-audit/scorecard.json`.

P3.31 clean checkpoint:

| Signal | P3.31 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 250 | Seven filtered-walk recording contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 16 / 25 | No package denominator changed. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 35 / 52 | The filtered template walk now uses the filesystem adapter. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.32 clean checkpoint:

| Signal | P3.32 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 254 | Four markdown-write recording contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 16 / 25 | No package denominator changed. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 34 / 51 | The template markdown write now uses the filesystem adapter. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.33 clean checkpoint:

| Signal | P3.33 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 258 | Four project-config-write recording contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 16 / 25 | No package denominator changed. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 33 / 50 | The project config write now uses the filesystem adapter. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.34 measured implementation result:

| Signal | P3.34 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 262 | Four tips-show-read recording contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 16 / 25 | No package denominator changed. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 32 / 49 | The tips-show source read now uses the filesystem adapter. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

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

Status: active. Moves 1 through 33 are complete.

Move 1 reduced Q1.2 from 127 to zero. A syntax-aware boundary contract covers
both `main.go` files, requires them to delegate to `cmd.ExecuteE() error`, and
rejects `log.Fatal*` or `os.Exit` outside the two executable `main` functions.
The legacy exported `cmd.Execute()` and `cmd.RootCmd` signatures remain intact;
the compatible `ExecuteE` addition is explicit in the P2A API allowlist. Cobra
initialization, pre-run hooks, and run hooks now return errors, while the two
entry points alone print the returned error and select exit 1.

Commit: `5c6f2fa`.

The clean implementation gate passed the boundary, API, Cobra, subprocess,
four host acceptance, preflight, test, race, vet, lint, audit-meta, and
empty-HOME count-2 contracts. Q1.1 also improved from fifteen to thirteen
untested packages. The full audit exited 1, never 2. It identified one
handoff-only Q3.4 wording regression inherited from `2034f63`; finalization
rephrases that host probe before the next measured move.

Move 2 added the zero-value-safe `internal/adapter/process` boundary and moved
Git clone, pull, init, add, and commit execution behind it without changing the
exported shell API or the legacy `shell.Output` behavior. Recording contracts
preserve full arguments and working directories, reject an empty population,
and cover both declared Git swaps. The three direct `pkg/shell/git.go` process
sites left Q1.3; the adapter contributes two in-boundary process effects, so the
scanner reports 77 violations of 79 production effect sites and four missing
adapter paths. Q1.4 is 2 of 8.

The inventory-bound `scripts/mutate-file-shell` path is a deliberately
non-executable seam-test driver, paired with its Q0.8 meta-test. It carries the
two labels to the real Go tests but declares no mutations and does not satisfy
Q2.1; the executable harness ratio stays 0 of 8. The upstream filename-only
denominator nevertheless calls it one harness and therefore reports the
expected Q2.2/Q2.3 project findings until P5 replaces it with the real
file-shell harness.

Commit: `03d6242`.

The clean move-2 gate passed focused adapter/seam tests, API/CLI compatibility,
all four host flows, preflight, test, install, launcher, uncached and race
tests, vet, the 15-control audit meta-suite, and empty-HOME count-2. The full
audit exited 1 for 16 documented findings, never 2, with five improved, two
held, zero regressed, and one not-comparable ratchet. The prior four-improved /
three-held forecast described the 0-of-8 Q1.4 state; raising Q1.4 to 2 of 8
necessarily moves that ratchet from held to improved.

Move 3 reused `internal/adapter/process` for the coherent Maven command flow.
The exported `maven.RunOn` signature and returned callback are unchanged, while
the callback now passes the complete executable, argument vector, project
working directory, nil stderr, logger stdout writer, and dependency error
through `process.Runner`. Recording contracts prove exact ordering and values,
whole-dependency forwarding, safe no-process defaults, and a non-empty asserted
population.

The inventory-bound `scripts/mutate-maven-sorting` path remains deliberately
non-executable and is paired with a Q0.8 meta-test. It carries the immutable
`maven-process` label to the three real Go contracts without declaring P5
mutations. Q1.4 is therefore 3 of 8 while the exact Q2.1 validator remains 0 of
8. The two Maven process sites left Q1.3; preserving `logger.StdOut()` as the
process stdout exposes that existing `*os.File` as one filesystem capability,
so the scanner's exact net result is 76 violations of 78 production sites.

Commit: `f5ee37d`.

The clean move-3 gate passed focused adapter/Maven/seam tests, API/CLI and
subprocess compatibility, all four host flows, preflight, test, install,
launcher, uncached and race tests, vet, the 15-control audit meta-suite, and
empty-HOME count-2. The full audit exited 1 for 16 documented findings, never
2, with five improved, two held, zero regressed, and one not-comparable
ratchet.

Move 4 added a private, zero-value-safe interface boundary around only the Git
clone/pull dependency used by `GitCloudConfig.Refresh`. The exported
`CloudConfig` interface, `GitCloudConfig`, `OpenGitCloudConfig`, every caller,
the exact `<target>/.git` probe, logging, formatted errors, and cache-first
branch selection remain unchanged. Production delegates to the existing
`shell.GitClone` and `shell.GitPull` functions, so config does not reconstruct
Git argv or bypass the process adapter.

Five recording contracts prove the complete configured URL precedes the
complete target on clone, an existing `.git` path selects only pull, clone and
pull errors retain `shell.Output.FormatError`, the dependency's zero value
cannot mutate Git state, and an empty call population fails. Test setup uses
the guarded `internal/testutil` copy helper; an initial direct `os.MkdirAll`
fixture was rejected by Q0.6 and corrected before the implementation commit.

The inventory-bound `scripts/mutate-config-cloud` path is deliberately
non-executable and its Q0.8 meta-test binds only the immutable `cloud-clone`
label to the five Go contracts. It declares no P5 mutations. Q1.4 is therefore
4 of 8 while exact Q2.1 remains 0 of 8; the upstream filename-only denominator
now sees three non-executable `mutate-*` paths.

Commit: `ee5e9ab`.

The clean move-4 gate passed focused config/shell/process/seam tests, API/CLI
and subprocess compatibility, all four host flows, preflight, test, install,
launcher, uncached and race tests, vet, the 15-control audit meta-suite, and
empty-HOME count-2. The full audit exited 1 for 16 documented findings, never
2, with L0 8 of 8, five improved, two held, zero regressed, and one
not-comparable ratchet.

Move 5 added the zero-value-safe `internal/adapter/httpclient` boundary for
complete anonymous or basic-auth GET request values. The system dependency
preserves the existing anonymous `http.Get` path and authenticated
`NewRequest`/`SetBasicAuth`/fresh-client path. The exported `http.GetXml` and
`http.GetAuthXml` signatures still own XML decoding and the shared response
status, body-read, and close lifecycle; JSON/token, Wget/Wpost, and Kibana
requests remain outside the move.

Five Maven recording contracts prove the complete repository-derived metadata
URL, nil-auth anonymous selection, username before password, parsed XML,
dependency-error propagation, a safe no-request default, complete destination
delivery, and rejection of an empty population. Adapter and HTTP tests also
cover complete request forwarding, real anonymous/basic-auth transport,
success parsing, formatted status errors, read errors, and response closure.

The inventory-bound `scripts/mutate-maven-sorting` path remains deliberately
non-executable. Its Q0.8 meta-test now binds both the `maven-process` and
`maven-http` labels to eight Go contracts and still declares no P5 mutations.
Q1.4 is 5 of 8 and exact Q2.1 remains 0 of 8. The authenticated XML
`NewRequest` and `Client.Do` sites left Q1.3; the adapter contributes three
in-boundary HTTP sites while the direct anonymous helper remains for
out-of-scope JSON, yielding 74 violations of 79 production sites.

Commit: `ffc4e77`.

The clean move-5 gate passed focused Maven/HTTP/adapter/seam tests, API/CLI and
subprocess compatibility, all four host flows, preflight, test, install,
launcher, uncached and race tests, vet, the 15-control audit meta-suite, and
empty-HOME count-2. The full audit exited 1 for 16 documented findings, never
2, with L0 8 of 8, five improved, two held, zero regressed, and one
not-comparable ratchet.

Move 6 added the zero-value-safe `internal/adapter/filesystem` boundary and
routed only the non-existing-target template copy operations through it. The
exported `template.MergeTemplate`, `file.CopyOrMerge`, and `file.CopyFile`
signatures and callers are unchanged. Existing-target merge selection and
internals are untouched; the migrated branch preserves source read,
destination and directory probes, the legacy double directory probe, 0755
directory creation, source mode lookup, source-before-destination logging, and
destination write order. The legacy missing-directory error returned after a
failed `MkdirAll` is explicitly characterized rather than silently corrected.

Fourteen adapter, file, and template contracts prove complete source and
resolved destination delivery, missing-copy versus existing-merge selection,
bytes and mode, directory creation, dependency error order, logging, safe
defaults, whole-dependency forwarding, system behavior, and rejection of empty
recorded populations. The inventory-bound `scripts/mutate-template` path is a
deliberately non-executable seam driver with a Q0.8 meta-test. It binds only the
immutable `template-copy` label and declares no P5 mutations. Q1.4 is 6 of 8
and exact Q2.1 is 0 of 8; the upstream filename-only denominator sees four
non-executable `mutate-*` paths.

Commit: `204e222`.

The clean move-6 gate passed the 14 focused contracts, API/CLI and subprocess
compatibility, all four host flows, preflight, test, install, launcher,
uncached and race tests, vet, the 15-control audit meta-suite, and empty-HOME
count-2. Q0.6 records zero unsafe direct test writes, Q1.2 is zero, and Q1.3 is
71 violations of 80 production effect sites with filesystem valid and the
clock/server paths absent. The full audit exited 1 for 16 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, and one
not-comparable ratchet.

Move 7 extended the existing HTTP and filesystem adapters through the coherent
Spring initializer download-to-archive flow. The exported
`spring.DownloadInitializer` and `http.Wget` signatures and callers remain
unchanged. A private Spring dependency object accepts the already resolved
archive path and records download, unzip, and deletion follow-ups without
moving `os.Getwd`, `time.Now`, unzip internals, or archive deletion.

`Wget` now passes the complete anonymous URL through `httpclient.Dependencies`,
then creates/truncates and copies the response body through
`filesystem.Dependencies`. Thirteen recording and system contracts preserve
the encoded `<base>/starter.zip?<formData.Encode()>` URL before the complete
archive path, anonymous GET and redirect-capable transport selection, lack of
status rejection, exact body bytes, create/truncate, file-before-response
close order, errors, logging, unzip/delete selection, safe defaults, complete
dependency delivery, `spring-<Unix>.zip` naming, and a non-empty population.

The inventory-bound `scripts/mutate-spring` path is deliberately
non-executable. Its Q0.8 meta-test binds only the immutable `spring-download`
label to the 13 contracts and declares no P5 mutations. Q1.4 is therefore 7 of
8 while exact Q2.1 remains 0 of 8; the upstream filename-only denominator sees
five non-executable `mutate-*` paths.

Commit: `89d0f76`.

The clean move-7 gate passed focused Spring/HTTP/filesystem/seam tests,
API/CLI and subprocess compatibility, all four host flows, preflight, test,
install, launcher, uncached and race tests, vet, the 15-control audit
meta-suite, and empty-HOME count-2. Q0.6 records zero unsafe direct test writes,
Q1.2 is zero, and Q1.3 is 67 violations of 77 production effect sites with
clock and server absent. The full audit exited 1 for 16 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, and one
not-comparable ratchet.

Move 8 routed the anonymous request used by `http.GetJson` through the existing
`httpclient.Dependencies` value and the shared response lifecycle. The exported
signature and JSON unmarshal target remain unchanged. Standard anonymous
`http.Get` transport and redirects, status `>= 400` text, body read and close,
dependency, read, and unmarshal errors are recorded without affecting XML,
token HTTP, Wget/Wpost, or Kibana behavior.

The coherent Spring discovery flow now has a private, zero-value-safe interface
boundary that passes complete `httpclient.Request` and destination values. The
exported `GetRoot`, `GetDependencies`, and `Validate` functions preserve the
exact base and `<base>/dependencies` URLs, parsed response types, validation
selection and user-dependency order, invalid-dependency reporting, and all
later build follow-ups. Nine new top-level contracts bring the suite to 101
tests and reject an empty recorded discovery population.

The non-executable `scripts/mutate-spring` driver still carries its single
immutable `spring-download` label and now binds 22 focused download and
discovery contracts. It declares no P5 mutations, so Q1.4 remains 7 of 8 and
exact Q2.1 remains 0 of 8. Removing the last direct anonymous `GetJson`
request changes Q1.3 from 67 of 77 to 66 of 76 production effect sites.

Commit: `dee214c`.

The clean move-8 gate passed focused HTTP/Spring/adapter/seam tests, API/CLI and
subprocess compatibility, all four host flows, preflight, test, install,
launcher, uncached and race tests, vet, the 15-control audit meta-suite, and
empty-HOME count-2. The full audit exited 1 for 16 documented findings, never
2, with L0 8 of 8, five improved, two held, zero regressed, and one
not-comparable ratchet.

Move 9 extended `httpclient.Request` with one complete `BearerJSON` access-token
value. The system dependency retains anonymous `http.Get` and basic-auth
selection unchanged, while bearer JSON uses the existing fresh-client GET path
with exact `Authorization: Bearer <token>` and `Content-Type: application/json`
headers and standard redirects.

The exported `http.GetJsonWithAccessToken`, `bitbucket.With`, and
`bitbucket.QueryRepos` signatures remain unchanged. The token helper now routes
only request execution through `httpclient.Dependencies` while preserving its
distinct ignored-read-error, no-status-rejection, debug-log, JSON-target, and
body-close lifecycle. A private zero-value-safe Bitbucket interface receives
the complete request and decode destination; synchronization retains project
order, lowercase repository lookup, warning behavior, and later clone/pull
follow-ups.

Eleven new top-level contracts bring the suite to 112 tests and add
`pkg/bitbucket` to the tested population. They prove complete bearer request
delivery, GET and headers, redirects, non-success parsing, ignored read errors,
dependency/unmarshal errors, closure, logging, safe defaults, exact project and
repository URLs/limits, parsed response content, synchronization selection and
warnings, and rejection of an empty recorded population. No seam driver or
inventory label changed, so Q1.4 remains 7 of 8 and exact Q2.1 remains 0 of 8.
The two token request sites left Q1.3, producing 64 violations of 74 production
effect sites.

Commit: `789ae23`.

The clean move-9 gate passed focused Bitbucket/HTTP/adapter and adjacent Spring
and Maven tests, API/CLI and subprocess compatibility, all four host flows,
preflight, test, install, launcher, uncached and race tests, vet, the 15-control
audit meta-suite twice, and empty-HOME count-2. The full audit exited 1 for 16
documented findings, never 2, with L0 8 of 8, five improved, two held, zero
regressed, and one not-comparable ratchet.

Move 10 extended `httpclient.Request` with one narrow `POST` value containing
the complete body bytes and headers. The system dependency retains the exact
anonymous `http.Get` fast path and existing basic-auth and bearer JSON GET
paths. POST presence selects `http.MethodPost`, a fresh client, and the supplied
body and headers, preserving standard redirect behavior.

Exported `kibana.POST`, `KibanaFetchRequest`, `KibanaResponse`, and
`ExecuteKibanaQuery` signatures remain unchanged, as does private
`internalPOST`. Only its request execution now passes through
`httpclient.Dependencies`; it still owns the first `"size"` rewrite to 500,
four exact headers, response body read and close, no status rejection, newline
split, and header/result JSON decoding. A private zero-value-safe query
interface receives the complete `KibanaFetchRequest` through initial and
recursive callers while production delegates to the exported legacy `POST`
boundary. Retry selection and its 15-second clock remain unchanged.

Eleven new top-level contracts bring the suite to 123 tests. They prove exact
URL, POST method, rewritten body, all four headers, standard redirects,
non-success parsing, dependency/read/unmarshal errors, response closure, safe
defaults, complete caller delivery, recursive propagation, deterministic retry
structure, and rejection of empty recorded populations. No seam driver or
inventory label changed, so Q1.4 remains 7 of 8 and exact Q2.1 remains 0 of 8.
The two direct Kibana request-execution sites left Q1.3, producing 62
violations of 72 production effect sites.

Commit: `07ac6ce`.

The clean move-10 gate passed focused Kibana/HTTP/adapter and adjacent Spring,
Maven, and Bitbucket tests, API/CLI and subprocess compatibility, all four host
flows, preflight, test, install, launcher, uncached and race tests, vet, the
15-control audit meta-suite, and empty-HOME count-2. The full audit exited 1
for 16 documented findings, never 2, with L0 8 of 8, five improved, two held,
zero regressed, and one not-comparable ratchet.

Move 11 extended the existing narrow POST representation with an explicit
default-client selection used only by `Wpost`. Exported
`http.Wpost(downloadUrl, filePath string, formData url.Values) error` and every
caller remain unchanged. A private complete-dependency boundary passes the
full URL, `formData.Encode()` bytes, POST method, exact form content type, and
default-client choice through `httpclient.Dependencies`, then routes create and
copy through `filesystem.Dependencies` while Wpost retains response-before-file
creation and file-before-response deferred closure.

Seven new top-level contracts bring the suite to 130 tests. They prove the
debug log, complete request and dependency delivery, standard
`http.PostForm` client and redirect behavior, exact copied bytes,
create/truncate, absence of status rejection, request/create/copy errors after
recording, close order, ignored close errors, safe defaults, and rejection of
empty recorded populations. Anonymous, basic-auth, bearer JSON, Kibana POST,
Wget, and Spring behavior remain unchanged. No inventory entry, seam driver,
or mutation label changed, so Q1.4 remains 7 of 8 and exact Q2.1 remains 0 of
8. Exactly the four Wpost HTTP/filesystem effects left Q1.3, producing 58
violations of 68 production effect sites.

Commit: `a7eb3ef`.

The clean move-11 gate passed focused and adjacent adapter/HTTP/Spring/Kibana
contracts, API/CLI and subprocess compatibility, all four host flows,
preflight, test, install, launcher, uncached and race tests, vet, the 15-control
audit meta-suite, and empty-HOME count-2. The full audit exited 1 for 16
documented findings, never 2, with L0 8 of 8, five improved, two held, zero
regressed, one not-comparable ratchet, and zero dirty paths.

The requested operational continuity change completed at `1b85711` separately
from the numbered P3 quality moves. Normal launcher execution now
supervises one fresh `codex exec --json` turn per archived mission with exact
prompt bytes, normal service tier, workspace-write scope, concise terminal
progress, and uniquely named raw logs outside the worktree. Its structured
event parser and post-turn repository validation require both a successful
terminal stream and a clean committed reciprocal handoff before another turn;
failure, malformed or contradictory events, no progress, contract drift, and
signals stop fail-closed. The 62 Bash 3.2 launcher contracts pass, including
two-generation progression and eventual `COMPLETE`, without invoking the real
Codex executable.

Move 12 gave the private Bitbucket repository flow one complete dependency
value containing `filesystem.Dependencies` and a resolvable Git interface.
Production `With` selects `filesystem.System()`, repository existence delegates
to the adapter's established only-missing-means-absent rule, and clone/pull
execution still reaches the unchanged `shell.GitClone` and `shell.GitPull`
boundaries. The safe zero value refuses a Git operation instead of probing or
mutating the developer repository.

Six new top-level recording contracts bring the suite to 136 tests. They prove
the complete repository path and dependency value, missing selection of clone,
existing and other stat-error selection of pull, exact clone/pull paths and
logging, operation-error propagation, safe defaults, and non-empty recorded
stat and Git populations. No adapter contract, inventory label, seam driver, or
mutation harness changed. Q1.3 moves from 58 of 68 to 57 of 67 while Q1.2 stays
zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

Commit: `e13a036`.

The clean move-12 gate passed focused Bitbucket/filesystem/process/shell
contracts, API/CLI and subprocess compatibility, all four host flows,
preflight, test, install, launcher, uncached and race tests, vet, the 15-control
audit meta-suite, and empty-HOME count-2. The full audit exited 1 for 16
documented findings, never 2, with L0 8 of 8, five improved, two held, zero
regressed, one not-comparable ratchet, and zero dirty paths.

Move 13 gave exported `file.Exists(string) bool` a private complete dependency
value containing `filesystem.Dependencies`. The exported function remains the
production wrapper and selects `filesystem.System()`; its private helper
delegates the probe to `filesystem.Exists`, preserving the legacy rule that
only a missing-path error returns false and every other result returns true.
The dependency's zero value is safe and returns false without probing a real
path.

Four new top-level recording contracts bring the suite to 140 tests. They
prove the complete path and dependency value, missing, existing, and other
stat-error results, production system selection, safe defaults, and a
non-empty recorded stat population. No adapter contract, inventory label,
seam driver, mutation harness, public API, or caller changed. Q1.3 moves from
57 of 67 to 56 of 66 while Q1.2 stays zero, Q1.4 stays 7 of 8, and exact Q2.1
stays 0 of 8.

Commit: `f59a3f0`.

The clean move-13 gate passed focused file/filesystem and all relevant caller
contracts, API/CLI and subprocess compatibility, all four host flows,
preflight, test, install, launcher, uncached and race tests, vet, the 15-control
audit meta-suite, and empty-HOME count-2. The full audit exited 1 for 16
documented findings, never 2, with L0 8 of 8, five improved, two held, zero
regressed, one not-comparable ratchet, and zero dirty paths.

Move 14 gave exported
`file.Overwrite(lines []string, filePath string) error` a private complete
dependency value containing `filesystem.Dependencies`. The exported function
remains the production wrapper and selects `filesystem.System()`; its private
helper delegates the single write to `filesystem.WriteFile` with the unchanged
`strings.Join(lines, "\n")` bytes and mode `0644`. The dependency's zero value
returns `filesystem.ErrNoFilesystem` before mutation.

Four new top-level recording contracts bring the suite to 144 tests. They
prove complete dependency delivery, path, empty/single/multiple line bytes,
mode, dependency-error propagation, production system selection, safe
defaults, and a non-empty recorded write population. No adapter contract,
inventory label, seam driver, mutation harness, public API, or caller changed.
Q1.3 moves from 56 of 66 to 55 of 65 while Q1.2 stays zero, Q1.4 stays 7 of 8,
and exact Q2.1 stays 0 of 8.

Commit: `61714a5`.

The clean move-14 gate passed focused file/filesystem and all relevant caller
contracts, API/CLI and subprocess compatibility, all four host flows,
preflight, test, install, launcher, uncached and race tests, vet, the 15-control
audit meta-suite, and empty-HOME count-2. One repeated launcher run missed the
partial raw log in its nested signal fixture; the prior standalone and
preflight runs and the immediate complete `make test` rerun passed all 62
controls. The full audit exited 1 for 16 documented findings, never 2, with L0
8 of 8, five improved, two held, zero regressed, one not-comparable ratchet,
and zero dirty paths.

Move 15 gave exported `file.CreateFile(path, content string) error` a private
complete dependency value containing `filesystem.Dependencies`. The exported
function remains the production wrapper and selects `filesystem.System()`;
its private helper delegates the write to `filesystem.WriteFile` with the
unchanged `[]byte(content)` bytes and mode `0644`. The dependency's zero value
returns `filesystem.ErrNoFilesystem` before mutation.

Four new top-level recording contracts bring the suite to 148 tests. They
prove complete dependency delivery, path, representative empty and multiline
content bytes, mode, dependency-error propagation, production system
selection, safe defaults, and a non-empty recorded write population. No
adapter contract, inventory label, seam driver, mutation harness, public API,
or caller changed. Q1.3 moves from 55 of 65 to 54 of 64 while Q1.2 stays zero,
Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

Commit: `60e5aac`.

The clean move-15 gate passed focused file/filesystem and all relevant caller
contracts, API/CLI and subprocess compatibility, all four host flows,
preflight, test, install, launcher, uncached and race tests, vet, the 15-control
audit meta-suite, and empty-HOME count-2. The full audit exited 1 for 16
documented findings, never 2, with L0 8 of 8, five improved, two held, zero
regressed, one not-comparable ratchet, and zero dirty paths.

Move 16 gave exported `file.CreateDirectory(path string) error` a private
complete dependency value containing `filesystem.Dependencies`. The exported
function remains the production wrapper and selects `filesystem.System()`;
its private helper delegates the probe and recursive creation to
`filesystem.Stat` and `filesystem.MkdirAll` with mode `0755`. Missing paths
alone select creation, existing paths and other stat errors return nil, and a
creation failure still returns the original missing-path stat error. The
dependency's zero value returns nil without probing or mutating a real path.

Five new top-level recording contracts bring the suite to 153 tests. They
prove complete dependency delivery and path, missing-only selection of exactly
one recursive creation, mode, exact legacy errors, existing and non-missing
stat-error behavior, production system selection, safe defaults, and non-empty
recorded stat and creation populations. No adapter contract, inventory label,
seam driver, mutation harness, public API, or caller changed. Q1.3 moves from
54 of 64 to 52 of 62 while Q1.2 stays zero, Q1.4 stays 7 of 8, and exact Q2.1
stays 0 of 8.

Commit: `e054082`.

The clean move-16 gate passed focused file/filesystem and all relevant caller
contracts, API/CLI and subprocess compatibility, all four host flows,
preflight, test, install, launcher, uncached and race tests, vet, the 15-control
audit meta-suite, and empty-HOME count-2. One repeated `make test` run hit the
previously observed partial-raw-log signal-fixture flake; the standalone,
preflight, and immediate complete `make test` rerun passed all 62 launcher
controls. The full audit exited 1 for 16 documented findings, never 2, with L0
8 of 8, five improved, two held, zero regressed, one not-comparable ratchet,
and zero dirty paths.

Move 17 gave exported `file.Open(filePath string) ([]byte, error)` a private
complete dependency value containing `filesystem.Dependencies`. The exported
function remains the production wrapper and selects `filesystem.System()`;
its private helper delegates exactly one complete-path read to
`filesystem.ReadFile`. Successful reads return the dependency's exact complete
bytes. Every dependency error discards any partial data and returns a non-nil
empty `[]byte{}` with that exact error. The dependency's zero value returns
`filesystem.ErrNoFilesystem` without accessing the supplied path.

Four new top-level recording contracts bring the suite to 157 tests. They
prove complete dependency delivery and path, exact empty, multiline, and
binary bytes, exact error identity with non-nil empty results even after
partial data, production system selection, safe defaults, and a non-empty
recorded read population. `OpenFile` and every other file operation remain
unchanged. No adapter contract, inventory label, seam driver, mutation
harness, public API, or caller changed. Q1.3 moves from 52 of 62 to 50 of 60
while Q1.2 stays zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

Commit: `27c0d1a`.

The move-17 implementation gate passed focused file/filesystem and all
relevant caller contracts, API/CLI and subprocess compatibility, all four
host flows, preflight, test, install, launcher, uncached and race tests, vet,
the 15-control audit meta-suite, and empty-HOME count-2. The first complete
preflight run hit the previously observed partial-raw-log signal-fixture
flake; the standalone launcher contract and immediate complete preflight
rerun passed all 62 launcher controls. The focused seven-criterion audit
exited 1 for documented findings, never 2, with four improved, two held, zero
regressed, and one not-comparable ratchet. The clean full audit at `27c0d1a`
exited 1 for 16 documented findings, never 2, with L0 8 of 8, five improved,
two held, zero regressed, one not-comparable ratchet, and zero dirty paths.

Move 18 extended the existing zero-value-safe filesystem adapter with only an
`OpenFile(path, flags, mode) (*os.File, error)` operation. Exported
`file.OpenFile(fileName string) (*os.File, error)` remains the production
wrapper and selects `filesystem.System()`. Its private complete dependency now
drives the unchanged existence probe, missing-only empty write with mode
`0644`, create-error short circuit, and final append-open with exact flags
`os.O_APPEND|os.O_WRONLY` and mode `0644`. Existing paths and non-missing stat
errors still skip creation, and the exact dependency file and error pair are
returned.

Seven new top-level recording contracts bring the suite to 164 tests. They
prove the complete dependency and path, exact flags and modes, exact returned
file pointer and error, missing/existing/non-missing-stat-error sequencing,
empty creation and create-error short circuit, production system selection,
safe defaults without path access, and non-empty recorded populations. Every
complete filesystem recording double implements the adapter extension. No
inventory label, seam driver, mutation harness, public API, or caller changed,
and `pkg/shell` retains its separate direct `os.OpenFile`. Q1.3 moves from 50
of 60 to 49 of 60 while Q1.2 stays zero, Q1.4 stays 7 of 8, and exact Q2.1
stays 0 of 8.

Commit: `c2f3597`.

The move-18 implementation gate passed focused file/filesystem and adjacent
complete-double contracts, API/CLI and subprocess compatibility, all four host
flows, preflight, test, install, launcher, uncached and race tests, vet, the
15-control audit meta-suite, and empty-HOME count-2. The focused
seven-criterion audit exited 1 for documented findings, never 2, with four
improved, two held, zero regressed, and one not-comparable ratchet. The clean
full audit at `c2f3597` exited 1 for 16 documented findings, never 2, with L0
8 of 8, five improved, two held, zero regressed, one not-comparable ratchet,
and zero dirty paths.

Move 19 extended the existing zero-value-safe filesystem adapter with only a
`Remove(path string) error` operation. Exported
`file.DeleteSingleFile(filePath string) error` remains the production wrapper
and selects `filesystem.System()`. Its private complete dependency delegates
the complete path unchanged and returns the dependency's exact error without
probing, normalizing, wrapping, retrying, or adding selection behavior. The
dependency's zero value returns `filesystem.ErrNoFilesystem` without accessing
the supplied path, and the system implementation performs the single
`os.Remove(path)` operation.

Five new top-level recording contracts bring the suite to 169 tests. They
prove complete dependency selection and delivery, the exact path, successful
nil error, exact dependency-error identity, safe defaults without path access,
and a non-empty recorded remove population. Every complete filesystem
recording double implements the adapter extension. No inventory label, seam
driver, mutation harness, public API, or caller changed, and `DeleteAll`,
`ClearDir`, and every other file operation remain unchanged. Q1.3 moves from
49 of 60 to 48 of 60 while Q1.2 stays zero, Q1.4 stays 7 of 8, and exact Q2.1
stays 0 of 8.

Commit: `9e2d669`.

The move-19 implementation gate passed focused file/filesystem and adjacent
complete-double and caller contracts, API/CLI and subprocess compatibility,
all four host flows, preflight, test, install, launcher, uncached and race
tests, vet, the 15-control audit meta-suite, and empty-HOME count-2. The focused
seven-criterion audit exited 1 for documented findings, never 2, with four
improved, two held, zero regressed, and one not-comparable ratchet. The clean
full audit at `9e2d669` exited 1 for 16 documented findings, never 2, with L0
8 of 8, five improved, two held, zero regressed, one not-comparable ratchet,
and zero dirty paths.

Move 20 extended the existing zero-value-safe filesystem adapter with only a
`RemoveAll(path string) error` operation. Exported
`file.DeleteAll(dirPath string) error` remains the production wrapper and
selects `filesystem.System()`. Its private complete dependency delegates the
complete path unchanged and returns the dependency's exact error without
probing, normalizing, wrapping, retrying, or adding selection behavior. The
dependency's zero value returns `filesystem.ErrNoFilesystem` without accessing
the supplied path, and the system implementation performs the single
`os.RemoveAll(path)` operation.

Six new top-level recording contracts bring the suite to 175 tests. They prove
complete dependency selection and delivery, the exact path, successful nil
error, exact dependency-error identity, safe defaults without path access, a
non-empty recorded recursive-remove population, and the system adapter's
recursive removal behavior. Every complete filesystem recording double
implements the adapter extension. No inventory label, seam driver, mutation
harness, public API, or caller changed, and `ClearDir`, `Move`, and every other
file operation remain unchanged. Q1.3 moves from 48 of 60 to 47 of 60 while
Q1.2 stays zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

Commit: `a4deb76`.

The move-20 implementation gate passed focused file/filesystem, adjacent
complete-double, and caller contracts, API/CLI and subprocess compatibility,
all four host flows, preflight, test, install, launcher, uncached and race
tests, vet, the 15-control audit meta-suite, and empty-HOME count-2. The first
standalone launcher run and the first complete `make test` run each hit the
previously observed partial-raw-log signal-fixture flake; their immediate
complete reruns passed all 62 launcher controls. The focused seven-criterion
audit exited 1 for documented findings, never 2, with four improved, two held,
zero regressed, and one not-comparable ratchet. The clean full audit at
`a4deb76` exited 1 for 16 documented findings, never 2, with L0 8 of 8, five
improved, two held, zero regressed, one not-comparable ratchet, and zero dirty
paths.

Move 21 extended the existing zero-value-safe filesystem adapter with only a
`Rename(source, destination string) error` operation. Exported
`file.Move(source, destination string) error` remains the production wrapper
and selects `filesystem.System()`. Its private complete dependency delegates
both complete paths unchanged and returns the dependency's exact error without
probing, normalizing, wrapping, retrying, copying, removing, or adding selection
behavior. The dependency's zero value returns `filesystem.ErrNoFilesystem`
without accessing either supplied path, and the system implementation performs
the single `os.Rename(source, destination)` operation.

Five new top-level recording contracts bring the suite to 180 tests. They
prove complete dependency selection and delivery, the exact source and
destination paths, successful nil error, exact dependency-error identity, safe
defaults without path access, and a non-empty recorded rename population.
Every complete filesystem recording double implements the adapter extension.
No inventory label, seam driver, mutation harness, public API, or caller
changed, and `DeleteAll`, `DeleteSingleFile`, `ClearDir`, and every other file
operation remain unchanged. Q1.3 moves from 47 of 60 to 46 of 60 while Q1.2
stays zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

Commit: `acda4e3`.

The move-21 implementation gate passed focused file/filesystem, adjacent
complete-double, and caller contracts, API/CLI and subprocess compatibility,
all four host flows, preflight, test, install, all 62 launcher controls,
uncached and race tests, vet, the 15-control audit meta-suite, and empty-HOME
count-2. The focused seven-criterion audit exited 1 for documented findings,
never 2, with four improved, two held, zero regressed, and one not-comparable
ratchet. The clean full audit at `acda4e3` exited 1 for 16 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths.

Move 22 extended the existing zero-value-safe filesystem adapter with only a
`Glob(pattern string) ([]string, error)` operation and reused its existing
`RemoveAll(path string) error` operation. Exported
`file.ClearDir(dirPath string, excludes []string) error` remains the production
wrapper and selects `filesystem.System()`. Its private complete dependency
delegates the exact `filepath.Join(dirPath, "*")` pattern, consumes the ordered
glob results unchanged, runs the complete ordered substring-exclusion loop,
and delegates each selected removal path unchanged. Existing skip and removal
logs remain in the same order, the first removal error still short-circuits,
and exact glob and removal errors are returned without wrapping, retrying,
preflighting, normalizing, sorting, deduplicating, or continuing.

Seven new top-level recording contracts bring the suite to 187 tests. They
prove complete dependency selection and delivery, the exact glob pattern,
ordered results and removals, exclusion and logging decisions, successful nil
error, exact glob- and removal-error identities, first-error short-circuiting,
safe defaults without developer-path access, and non-empty recorded glob and
removal populations. Every complete filesystem recording double implements
the adapter extension. No inventory label, seam driver, mutation harness,
public API, or caller changed, and `Move`, `DeleteAll`, `DeleteSingleFile`, and
every other file operation remain unchanged. Q1.3 moves from 46 of 60 to 44 of
59: the two direct `ClearDir` effects disappear, while the system `Glob` enters
the adapter population and the reused, already-counted system `RemoveAll`
causes the total effect-site denominator to contract by one. Q1.2 stays zero,
Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

Commit: `838daa1`.

The move-22 implementation gate passed focused file/filesystem, adjacent
complete-double, and caller contracts, API/CLI and subprocess compatibility,
all four host flows, preflight, test, install, uncached and race tests, vet,
the 15-control audit meta-suite, and empty-HOME count-2. The first standalone
launcher run and the first complete `make test` run each hit the previously
observed partial-raw-log signal-fixture flake; their immediate complete reruns
passed all 62 launcher controls. The focused seven-criterion audit exited 1
for documented findings, never 2, with four improved, two held, zero regressed,
and one not-comparable ratchet. The clean full audit at `838daa1` exited 1 for
16 documented findings, never 2, with L0 8 of 8, five improved, two held, zero
regressed, one not-comparable ratchet, and zero dirty paths.

Move 23 gave exported
`file.Render(inputFilePath string, outputFilePath string, r interface{}) error`
a private complete filesystem dependency and retained it as the production
wrapper selecting `filesystem.System()`. The private composition delegates the
exact input path to the existing adapter `ReadFile` operation before template
parsing, then delegates the exact output path to the existing adapter `Create`
operation only after a successful read and parse. It preserves the existing
`template.Must` parse panic, creation and truncation behavior, exact rendered
bytes, exact returned read, create, and execution errors, nil success, and the
existing lack of a close without wrapping, retrying, preflighting, normalizing,
reordering, cleaning up, or adding selection behavior.

Eight top-level recording contracts replace the old worktree-mutating render
contract and bring the suite to 194 tests. They prove complete system
dependency selection and delivery, exact input and output paths,
read-before-create ordering, no create or write after an input error, no create
after a template parse panic, exact create- and writer-error identities, exact
rendered bytes, successful nil error, safe defaults without developer-path
access, the unchanged lack of a close, and non-empty recorded read, create, and
write populations. The contracts perform no real filesystem mutation. No
adapter operation, complete filesystem recording double, inventory label,
seam driver, mutation harness, public API, or caller changed, and every other
file operation remains unchanged. Q1.3 moves from 44 of 59 to 43 of 58 because
the direct render create disappears while the already-counted adapter `Create`
operation is reused. Q1.2 stays zero, Q1.4 stays 7 of 8, and exact Q2.1 stays
0 of 8.

Commit: `224a691`.

The move-23 implementation gate passed focused file/filesystem and both caller
contracts, API/CLI and subprocess compatibility, all four host flows,
preflight, test, install, the standalone 62-control launcher contract,
uncached and race tests, vet, the 15-control audit meta-suite, and empty-HOME
count-2. The first complete `make test` run hit the previously observed
partial-raw-log signal-fixture flake; the standalone, preflight, and immediate
complete `make test` rerun passed all 62 launcher controls. The focused
seven-criterion audit exited 1 for documented findings, never 2, with four
improved, two held, zero regressed, and one not-comparable ratchet. The clean
full audit at `224a691` exited 1 for 16 documented findings, never 2, with L0
8 of 8, five improved, two held, zero regressed, one not-comparable ratchet,
and zero dirty paths.

Move 24 extended the existing zero-value-safe filesystem adapter with only a
`Walk(root string, callback filepath.WalkFunc) error` operation. Exported
`file.FindFirst(fileSuffix string, dir string) (result string, err error)`
remains the production wrapper and selects `filesystem.System()`. Its private
complete dependency delegates the exact root and callback through the adapter.
The unchanged callback still ignores `fi` and `errIn`, matches only
`strings.HasSuffix(path, fileSuffix)`, assigns the exact matching path before
returning `io.EOF`, stops at the first match, and normalizes only a final walk
error equal to `io.EOF` while preserving every other exact error and named
result.

Eight new top-level adapter and file recording contracts bring the suite to
202 tests. They prove complete dependency selection and delivery, exact root
and callback values, ordered callback paths, ignored callback metadata and
errors, case-sensitive exact suffix matching, first-match result and stop,
no-match empty result and nil error, exact non-EOF walk errors with the current
partial result, final `io.EOF` normalization, safe defaults without
developer-path access, system `filepath.Walk` behavior, and non-empty walk and
callback populations. The file contracts perform no real filesystem mutation.
Every complete filesystem recording double implements the adapter extension.
No inventory label, seam driver, mutation harness, public API, caller, or other
file operation changed. Q1.3 moves from 43 of 58 to 42 of 58 because the direct
find-first walk moves behind the adapter while the system adapter contributes
the corresponding already-declared boundary effect. Q1.2 stays zero, Q1.4
stays 7 of 8, and exact Q2.1 stays 0 of 8.

Commit: `549685d`.

The move-24 implementation gate passed focused filesystem/file contracts, both
caller packages, adjacent complete-double packages, API/CLI and subprocess
compatibility, all four host flows, preflight, test, install, the standalone
62-control launcher contract, uncached and race tests, vet, the 15-control
audit meta-suite, and empty-HOME count-2. The first complete `make test` run
hit the previously observed partial-raw-log signal-fixture flake; the
standalone, preflight, and immediate complete `make test` rerun passed all 62
launcher controls. The focused seven-criterion audit exited 1 for documented
findings, never 2, with four improved, two held, zero regressed, and one
not-comparable ratchet. The clean full audit at `549685d` exited 1 for 16
documented findings, never 2, with L0 8 of 8, five improved, two held, zero
regressed, one not-comparable ratchet, and zero dirty paths.

Move 25 routes only exported
`file.FindAll(suffix string, excludes []string, dir string) (result []string, err error)`'s
direct `filepath.Walk` operation through the existing filesystem adapter. The
exported production wrapper selects `filesystem.System()`, and its private
complete dependency delegates the exact root and callback through
`filesystem.Walk`. The callback body is unchanged: it ignores `fi` and
`errIn`, matches only `strings.HasSuffix(path, suffix) &&
!SuffixIn(path, excludes)`, preserves `SuffixIn`'s ordered
`strings.Contains` decisions, appends every matching path in callback order,
and always returns nil. Named results, partial accumulation, exact final
errors, and normalization of only a final error equal to `io.EOF` remain
unchanged.

Seven new top-level recording contracts bring the suite to 209 tests. They
prove complete system dependency selection, exact root and callback delivery,
ordered callback paths, ignored callback metadata and errors, case-sensitive
suffix decisions, exact ordered substring exclusion decisions including an
empty exclusion, ordered accumulation of every match, nil results for no
match, exact non-EOF errors with the current partial result, final `io.EOF`
normalization, safe zero-value behavior without developer-path access, and
non-empty recorded walk and callback populations. The contracts perform no
real filesystem mutation. No adapter operation, complete filesystem recording
double, inventory label, seam driver, mutation harness, public API, caller, or
other file operation changed. Q1.3 moves from 42 of 58 to 41 of 57. Q1.2 stays
zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

Commit: `04cfe44`.

The move-25 implementation gate passed focused file/filesystem and caller
contracts, adjacent complete-double packages, API/CLI and subprocess
compatibility, all four host flows, preflight, test, install, the standalone
62-control launcher contract, uncached and race tests, vet, the 15-control
audit meta-suite, and empty-HOME count-2. The focused seven-criterion audit
exited 1 for documented findings, never 2, with four improved, two held, zero
regressed, and one not-comparable ratchet. The clean full audit at `04cfe44`
exited 1 for 16 documented findings, never 2, with L0 8 of 8, five improved,
two held, zero regressed, one not-comparable ratchet, and zero dirty paths.

Move 26 routes only exported
`file.GrepRecursive(targetDir string, keyword string) (files []string, err error)`'s
direct `filepath.Walk` operation through the existing filesystem adapter. The
exported production wrapper selects `filesystem.System()`, and its private
complete dependency delegates the exact root and callback through
`filesystem.Walk`. The callback body is unchanged: it ignores `fi` and
`errIn`, calls `Grep(path, keyword)` exactly once for every delivered path,
returns that exact error when non-nil, appends the exact path only when `hit`
is true, and otherwise returns nil. `Grep` and `OpenLines` retain exact
case-sensitive `strings.Contains` matching and suppressed read failures.
Named results, ordered partial accumulation, exact final errors, and
normalization of only a final error equal to `io.EOF` remain unchanged.

Seven new top-level recording contracts bring the suite to 216 tests. They
prove complete system dependency selection, exact root and callback delivery,
ordered callback paths, ignored callback metadata and errors, exact
case-sensitive keyword decisions against immutable checked-in fixtures,
suppressed per-path read failures, ordered accumulation of every hit including
duplicate callback paths, nil results for no match, exact non-EOF errors with
the current partial result, final `io.EOF` normalization, safe zero-value
behavior without developer-path access, and non-empty recorded walk and
callback populations. The contracts perform no real filesystem mutation. No
adapter operation, complete filesystem recording double, inventory label,
seam driver, mutation harness, public API, caller, `Grep`, `OpenLines`, or
other file operation changed. Q1.3 moves from 41 of 57 to 40 of 56. Q1.2
stays zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

Commit: `d03e96d`.

The move-26 implementation gate passed focused file/filesystem and config
caller contracts, adjacent complete-double packages, API/CLI and subprocess
compatibility, all four host flows, preflight, test, install, the standalone
62-control launcher contract, uncached and race tests, vet, the 15-control
audit meta-suite, and empty-HOME count-2. The focused seven-criterion audit
exited 1 for documented findings, never 2, with four improved, two held, zero
regressed, and one not-comparable ratchet. The clean full audit at `d03e96d`
exited 1 for 16 documented findings, never 2, with L0 8 of 8, five improved,
two held, zero regressed, one not-comparable ratchet, and zero dirty paths.

Move 27 extends only the existing zero-value-safe filesystem adapter with
`ReadDir(path string) ([]fs.FileInfo, error)`. Its system implementation
delegates directly to `ioutil.ReadDir(path)`, retaining the exact returned
metadata, error, and filename sorting. Exported
`file.RemoveIntellijFiles(targetDir string, recursive bool, dryRun bool) (string, error)`
retains its signature and recursive selection. Only its non-recursive private
composition selects `filesystem.System()` and delegates the exact `targetDir`
through `filesystem.ReadDir`.

The non-recursive loop body remains unchanged. It iterates entries in their
delivered order and keeps two ordered, independent checks: a non-directory
whose exact name contains `.iml`, then a directory whose exact name contains
`.idea`. Each match still uses `Path("%s/%s", targetDir, f.Name())`, calls
`logAndDelete`, stops on its first error, and contributes to the exact
`Iml files: %d, .idea dirs: %d` report. An initial read error still returns an
empty report and the exact error. Recursive cleanup, `FindAll`,
`GrepRecursive`, `logAndDelete`, `DeleteAll`, and every other operation remain
unchanged.

Eight new top-level adapter and IDE recording contracts bring the suite to 224
tests. They prove complete system dependency selection, exact target delivery,
ordered entries and metadata, `ioutil.ReadDir` filename sorting, exact read
errors, exact ordered `.iml` file and `.idea` directory substring decisions,
ignored opposite types, case sensitivity, exact dry-run paths and logs, counts
and report, safe zero behavior without developer-path access, unchanged
recursive selection, and non-empty recorded directory-read and entry
populations. The only fixture mutation is confined to guarded temporary
directories. Every complete filesystem recording double implements the one
adapter extension. Q0.6 records 22 guarded safe-writer sites, 17 write and 5
copy, with zero unsafe direct writes. Q1.3 moves from 40 of 56 to 39 of 56;
Q1.2 stays zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

Commit: `d6cb593`.

The move-27 implementation gate passed focused filesystem/file, command-caller,
and adjacent complete-double contracts, API/CLI and subprocess compatibility,
all four host flows, preflight, test, install, the standalone 62-control
launcher contract, uncached and race tests, vet, the 15-control audit
meta-suite, and empty-HOME count-2. The first standalone launcher run and a
later complete `make test` run each hit the previously observed
partial-raw-log signal-fixture flake; their immediate complete reruns passed
all 62 controls, and the full preflight also passed all 62. The focused
seven-criterion audit exited 1 for documented findings, never 2, with four
improved, two held, zero regressed, and one not-comparable ratchet. The clean
full audit at `d6cb593` exited 1 for 16 documented findings, never 2, with L0
8 of 8, five improved, two held, zero regressed, one not-comparable ratchet,
and zero dirty paths.

Move 28 routes only exported
`GitCloudConfig.GitHookFiles(path string) ([]string, error)`'s direct
`ioutil.ReadDir(root)` operation through the existing filesystem adapter. The
exported production wrapper selects `filesystem.System()`, and its private
complete dependency delegates the exact receiver-derived
`file.Path("%s/%s", gitCfg.Implementation().Dir(), path)` root through
`filesystem.ReadDir`. The loop body is unchanged: it iterates entries in their
delivered order, appends only each exact `f.Name()` when `!f.IsDir()`, and
returns filename-only ordered results. Production sorting remains the adapter's
established `ioutil.ReadDir` behavior. An initial read error still returns an
unchanged nil result and the exact error, while empty and all-directory
populations still return a nil result and nil error.

Six new top-level config recording contracts bring the suite to 230 tests.
They prove complete system dependency selection and delivery, the exact
receiver-derived root, delivered entry order and metadata, exact read-error
identity, exact non-directory classification including ignored directories,
filename-only ordered results, nil results, safe zero behavior without
developer-path access, and non-empty recorded directory-read and entry
populations. The production-composition contracts perform no real filesystem
mutation. No adapter operation, complete recording double, inventory label,
seam driver, mutation harness, public API, caller, `Examples`, `Templates`, or
other config method changed. Q0.6 remains at 22 guarded safe-writer sites, 17
write and 5 copy, with zero unsafe direct writes. Q1.3 moves from 39 of 56 to
38 of 55; Q1.2 stays zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

Commit: `82e631d`.

The move-28 implementation gate passed focused config/filesystem and command
caller contracts, API/CLI and subprocess compatibility, all four host flows,
preflight, test, install, the standalone 62-control launcher contract,
uncached and race tests, vet, the 15-control audit meta-suite, and empty-HOME
count-2. The focused seven-criterion audit exited 1 for documented findings,
never 2, with four improved, two held, zero regressed, and one not-comparable
ratchet. The clean full audit at `82e631d` exited 1 for 16 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths.

Move 29 routes only exported
`GitCloudConfig.Examples() ([]string, error)`'s direct
`ioutil.ReadDir(examplesDir)` operation through the existing filesystem
adapter. The exported production wrapper selects `filesystem.System()`, and
its private complete dependency delegates the exact receiver-derived
`file.Path("%s/examples", gitCfg.Implementation().Dir())` root through
`filesystem.ReadDir`. The loop body is unchanged: it iterates entries in their
delivered order, appends only each exact `item.Name()` when `item.IsDir()`, and
returns directory-name-only ordered results. Production sorting remains the
adapter's established `ioutil.ReadDir` behavior. An initial read error still
returns an unchanged nil result and the exact error, while empty and all-file
populations still return a nil result and nil error.

Six new top-level config recording contracts bring the suite to 236 tests.
They prove complete system dependency selection and delivery, the exact
receiver-derived examples root, delivered entry order and metadata, exact
read-error identity, exact directory classification including ignored files,
directory-name-only ordered results, nil results, safe zero behavior without
developer-path access, and non-empty recorded directory-read and entry
populations. The production-composition contracts perform no real filesystem
mutation. No adapter operation, complete recording double, inventory label,
seam driver, mutation harness, public API, caller, `Templates`,
`GitHookFiles`, or other config method changed. Q0.6 remains at 22 guarded
safe-writer sites, 17 write and 5 copy, with zero unsafe direct writes. Q1.3
moves from 38 of 55 to 37 of 54; Q1.2 stays zero, Q1.4 stays 7 of 8, and exact
Q2.1 stays 0 of 8.

Commit: `cb94f81`.

The move-29 implementation gate passed focused config/filesystem and command
caller contracts, API/CLI and subprocess compatibility, all four host flows,
preflight, test, install, the standalone 62-control launcher contract,
uncached and race tests, vet, the 15-control audit meta-suite, and empty-HOME
count-2. The focused seven-criterion audit exited 1 for documented findings,
never 2, with four improved, two held, zero regressed, and one not-comparable
ratchet. The clean full audit at `cb94f81` exited 1 for 16 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths.

Move 30 routes only exported
`GitCloudConfig.Templates() (templates []CloudTemplate, err error)`'s direct
`filepath.Walk(root, callback)` operation through the existing filesystem
adapter. The exported production wrapper selects `filesystem.System()`, and
its private complete dependency delegates the exact receiver-derived
`file.Path("%s/templates", gitCfg.Implementation().Dir())` root and unchanged
callback through `filesystem.Walk`. The adapter's system implementation still
uses `filepath.Walk`, preserving lexical system order, delivered callback
paths, metadata and incoming errors, callback decisions, and the exact final
walk error.

Seven new top-level config recording contracts bring the suite to 243 tests.
They prove complete system dependency selection and delivery, the exact
receiver-derived templates root, non-empty recorded walk-root and callback
populations, delivered callback order, paths, metadata and incoming errors,
incoming walk-error suppression without metadata dereference, exact current
and legacy project-file matching, ignored entries, exact relative template
names and loaded projects, exact project-load and walk errors, ordered partial
results, nil results, and safe zero behavior without developer-path access.
The production-composition contracts perform no real filesystem mutation. No
adapter operation, complete recording double, inventory label, seam driver,
mutation harness, public API, caller, `Examples`, `GitHookFiles`, project
loading, or other config method changed. Q0.6 remains at 22 guarded
safe-writer sites, 17 write and 5 copy, with zero unsafe direct writes. Q1.3
moves from 37 of 54 to 36 of 53; Q1.2 stays zero, Q1.4 stays 7 of 8, and exact
Q2.1 stays 0 of 8.

Commit: `dfcfa75`.

The move-30 implementation gate passed focused config/filesystem and command
caller contracts, API/CLI and subprocess compatibility, all four host flows,
preflight, test, install, the standalone 62-control launcher contract,
uncached and race tests, vet, the 15-control audit meta-suite, and empty-HOME
count-2. One full `make test` attempt hit the documented launcher
signal-interruption partial-raw-log timing flake; its immediate complete rerun
passed. The focused seven-criterion audit exited 1 for documented findings,
never 2, with four improved, two held, zero regressed, and one not-comparable
ratchet. The clean full audit at `dfcfa75` exited 1 for 16 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths.

Move 31 routes only private
`template.filteredFilesFromTemplate(sourceDir string, filter []string)
(files []string, err error)`'s direct `filepath.Walk(sourceDir, callback)`
operation through the existing filesystem adapter. The unchanged private
production entry selects `filesystem.System()` for a private complete
dependency, which delegates the exact supplied `sourceDir` and unchanged
callback through `filesystem.Walk`. The adapter's system implementation still
uses `filepath.Walk`, preserving lexical system order, delivered callback
paths, metadata and incoming errors, callback decisions, and the exact final
walk error.

Seven new top-level template recording contracts bring the suite to 250
tests. They prove complete system dependency selection and delivery, the exact
supplied source root, non-empty recorded walk-root and callback populations,
delivered callback order, paths, metadata and incoming errors, exact directory
handling and metadata observations, the existing nil-info incoming-error
panic, non-nil-info incoming-error processing, slash-based root detection,
ordered filters, non-root `pom.xml` handling, substring matches, the global
`.render` exception, exact debug logging, ordered and nil results, exact final
walk errors, ordered partial results, and safe zero behavior without
developer-path access. Production-composition contracts perform no real
filesystem mutation. No adapter operation, complete recording double,
inventory label, seam driver, mutation harness, public API, caller, merge,
copy, search, render, config, file, or Maven operation changed. Q0.6 remains at
22 guarded safe-writer sites, 17 write and 5 copy, with zero unsafe direct
writes. Q1.3 moves from 36 of 53 to 35 of 52; Q1.2 stays zero, Q1.4 stays 7 of
8, and exact Q2.1 stays 0 of 8.

Commit: `170b0ae`.

The move-31 implementation gate passed focused template/filesystem and
relevant config, file, Maven, and command caller contracts, API/CLI and
subprocess compatibility, all four host flows, preflight, test, install, the
standalone 62-control launcher contract, uncached and race tests, vet, the
15-control audit meta-suite, and empty-HOME count-2. One full `make test`
attempt hit the documented nested launcher signal-interruption partial-raw-log
timing flake; its immediate complete rerun passed. The focused seven-criterion
audit exited 1 for documented findings, never 2, with four improved, two held,
zero regressed, and one not-comparable ratchet. The clean full audit at
`170b0ae` exited 1 for 16 documented findings, never 2, with L0 8 of 8, five
improved, two held, zero regressed, one not-comparable ratchet, and zero dirty
paths.

Move 32 routes only exported
`template.SaveTemplateListMarkdown(gitCfg config.CloudConfig,
markdownDocument string) (string, error)`'s direct
`os.WriteFile(readmePath, []byte(markdownDocument), 0644)` operation through
the existing filesystem adapter. The exported production entry retains its
signature and selects `filesystem.System()` only for a private complete
markdown-write dependency. The private composition retains the exact
`gitCfg.Implementation().Dir() + "/" + TemplatesDir + "/README.md"` path
expression and delegates that path, `[]byte(markdownDocument)`, and `0644`
through `filesystem.WriteFile` before returning the exact `readmePath, err`.

Four new top-level template recording contracts bring the suite to 254 tests.
They prove complete system dependency selection and delivery, the exact
receiver-derived README path without cleaning or normalization, non-empty
recorded write populations, exact empty, non-ASCII, and arbitrary document
bytes, exact `0644` mode, one write attempt, exact returned path on success and
failure, exact write-error identity, receiver evaluation order and count, and
safe zero behavior without developer-path access. The production-composition
contracts perform no real filesystem mutation. No adapter operation, complete
recording double, `templateCopyDependencies`, inventory label, seam driver,
mutation harness, public API, caller, list rendering, filtered walking, merge,
copy, search, render, config, file, or Maven operation changed. Q0.6 remains at
22 guarded safe-writer sites, 17 write and 5 copy, with zero unsafe direct
writes. Q1.3 moves from 35 of 52 to 34 of 51; Q1.2 stays zero, Q1.4 stays 7 of
8, and exact Q2.1 stays 0 of 8.

Commit: `8d49345`.

The move-32 implementation gate passed focused template/filesystem and
relevant config, file, Maven, and command caller contracts, API/CLI and
subprocess compatibility, all four host flows, preflight, test, install, the
standalone 62-control launcher contract, uncached and race tests, vet, the
15-control audit meta-suite, and empty-HOME count-2. The focused
seven-criterion audit exited 1 for documented findings, never 2, with four
improved, two held, zero regressed, and one not-comparable ratchet. The clean
full audit at the committed handoff exited 1 for 15 documented findings, never
2, with L0 8 of 8, five improved, two held, zero regressed, one not-comparable
ratchet, and zero dirty paths. The exact README-path contract also makes Q3.2
pass, reducing the documented finding count from 16 to 15.

Move 33 routes only exported
`(*config.ProjectConfiguration).WriteTo(targetFile string) error`'s direct
`os.WriteFile(targetFile, data, 0644)` operation through the existing
filesystem adapter. The exported method and `ProjectConfig` contract remain
unchanged. The production entry selects `filesystem.System()` only for a
private complete project-config-write dependency, and the private composition
retains the exact first log call, then the exact
`json.MarshalIndent(config, "", "    ")` evaluation and early error return,
then one `filesystem.WriteFile` attempt with the exact caller target, returned
bytes, and `0644` mode.

Four new top-level config contracts bring the suite to 258 tests. They prove
complete system dependency selection and delivery, the exact caller-supplied
target without cleaning or normalization, non-empty recorded write
populations, exact JSON bytes for nil, zero-value, and representative non-nil
receivers, empty and nil collections, non-ASCII fields and map-key
serialization, exact `0644` mode, one write attempt, exact write-error
identity, the exact log before the write, unchanged receiver state, and safe
zero behavior without developer-path access. The production-composition
contracts perform no real filesystem mutation. No adapter operation, other
complete recording double, inventory label, seam driver, mutation harness,
public API, caller, project initialization, `SortAndWritePom`, cloud config,
template, file, Maven, or command behavior changed. Q0.6 remains at 22 guarded
safe-writer sites, 17 write and 5 copy, with zero unsafe direct writes. Q1.3
moves from 34 of 51 to 33 of 50; Q1.2 stays zero, Q1.4 stays 7 of 8, and exact
Q2.1 stays 0 of 8.

Commit: `7754575`.

The move-33 implementation gate passed focused config/filesystem and relevant file,
template, Maven, context, and command caller contracts, API/CLI and subprocess
compatibility, all four host flows, preflight, test, install, the standalone
62-control launcher contract, uncached and race tests, vet, the 15-control
audit meta-suite, and empty-HOME count-2. The focused seven-criterion audit
exited 1 for documented findings, never 2, with four improved, two held, zero
regressed, and one not-comparable ratchet. The full audit from clean commit
`7754575` exited 1 for 15 documented findings, never 2, with L0 8 of 8, five
improved, two held, zero regressed, one not-comparable ratchet, and zero dirty
paths.

Move 34 routes only `tipsShowCmd.RunE`'s direct
`os.ReadFile(tipsPath)` operation through the existing filesystem adapter. All
command objects and registrations remain unchanged. For a non-empty argument
list, the production entry still evaluates exact `name := args[0]`, then exact
`file.Path("%s/%s.md", tips.LocalDir(ctx.CloudConfig), name)`, and now selects
`filesystem.System()` only for a private complete tips-show-read dependency
before one adapter `ReadFile` attempt with that exact caller-composed path.

Four new top-level command contracts bring the suite to 262 tests. They prove
complete system dependency selection and delivery, exact caller-composed path
delivery without cleaning or normalization, non-empty recorded read
populations, exact empty, representative, non-ASCII, and arbitrary returned
bytes, one read attempt, exact dependency-error identity at the private
boundary, safe zero behavior without developer-path access, and no unrelated
adapter operation. The production-composition contracts neither read the real
filesystem nor run the command. No adapter operation, established complete
recording double, inventory label, seam driver, mutation harness, public API,
command registration, missing-argument branch, list or sync branch, profile,
path composition, render/output/log/cloud sequencing, config, template, file,
or Maven behavior changed. Q0.6 remains at 22 guarded safe-writer sites, 17
write and 5 copy, with zero unsafe direct writes. Q1.3 moves from 33 of 50 to
32 of 49; Q1.2 stays zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

The move-34 implementation gate passed focused command/filesystem and relevant
tips, context, config, file, template, and Maven caller package tests, API/CLI
and subprocess compatibility, all four host flows, full preflight, test,
install, the standalone 62-control launcher contract, uncached and race tests,
vet, and the 15-control audit meta-suite. The focused seven-criterion audit
exited 1 for documented findings, never 2, with four improved, two held, zero
regressed, and one not-comparable ratchet. The clean full checkpoint audit and
implementation commit identity are recorded by the following session handoff.

With the process, HTTP, and filesystem boundaries green, continue one coherent
flow at a time using `.quality/inventory`:

1. `internal/adapter/process` for subprocess execution, not application exit.
2. `internal/adapter/httpclient` for HTTP request execution.
3. `internal/adapter/filesystem` for production filesystem mutation.

For each production-effect move, first add characterization or recording
tests, then move one coherent flow. Defaults must be safe and recording doubles
must preserve the complete dependency struct. The seven declared P3 swaps
through Spring download are complete, and Spring discovery, Bitbucket token
JSON, Kibana POST, Wpost, Bitbucket clone/pull selection, file existence, file
overwrite, file creation, directory creation, exported file reads, exported
append opens, exported single-file deletes, exported recursive deletes,
exported file moves, and exported clear-directory selection and removals now
use the adapters without adding inventory seams. Exported `file.Render` now
uses the same filesystem dependency for its input read and output creation,
and exported `file.FindFirst` and `file.FindAll` now use the filesystem
adapter's narrow `Walk` operation, without adding inventory seams.
`GitCloudConfig.GitHookFiles` and `GitCloudConfig.Examples` now reuse the
filesystem adapter's `ReadDir` operation, while `GitCloudConfig.Templates` and
`template.filteredFilesFromTemplate` reuse its `Walk` operation, without
adding inventory seams. `template.SaveTemplateListMarkdown` now reuses
`WriteFile` without changing its exact receiver-derived path, document bytes,
mode, returned path, or error. `(*config.ProjectConfiguration).WriteTo` now
reuses `WriteFile` without changing the exported method or `ProjectConfig`
contract, exact log and `json.MarshalIndent` sequencing, caller target,
serialized bytes, `0644` mode, or exact error. `tipsShowCmd.RunE` now reuses
`ReadFile` without changing command registration, exact name and path
composition, read error wrapping, returned bytes, or later command sequencing.
Continue with the next isolated production effect while leaving project
initialization, `SortAndWritePom`, cloud config, template/file/Maven behavior,
every other tips branch, command sequencing, the adapter, public API, and every
other completed effect unchanged.

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
