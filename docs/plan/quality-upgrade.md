# Quality Upgrade Plan

Last accepted quality checkpoint: 2026-09-05, commit `ca19dcd`.

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
- Put all disposable agent/test state beneath `$CODEX_SESSION_SCRATCH_ROOT`.
  Preserve only compact tracked decisions and digests; never create independent
  retained `/private/tmp/ply-*` roots.

The user has authorized the ordered roadmap through P8. The following block is
machine-readable launcher state; keep the order and vocabulary exact.

<!-- CODEX_AUTHORIZED_CHECKPOINTS_BEGIN -->
P2A|complete
P2B|complete
P3|complete
P4|complete
P5|complete
P6|complete
P7|active
P8|queued
<!-- CODEX_AUTHORIZED_CHECKPOINTS_END -->

On 2026-09-06, the operator removed accumulated reproducible scratch, audit,
cache, and Ply Docker artifacts after they consumed roughly 763 GiB. All
external evidence paths recorded later in this historical plan are receipts,
not live retention requirements. Recreate current evidence only inside the
launcher's bounded scratch tree; the launcher and individual heavy harnesses
must clean it automatically on exit.

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

P3.34 clean checkpoint:

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

P3.35 clean checkpoint:

| Signal | P3.35 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 266 | Four structurizr-output-write recording contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 17 / 25 | `pkg/structurizr` now has focused contracts. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 31 / 48 | The structurizr output write now uses the filesystem adapter. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.36 clean checkpoint:

| Signal | P3.36 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 270 | Four Maven graph-style-write recording contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 17 / 25 | The contracts stay in the already-tested Maven package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 30 / 47 | The Maven graph-style output write now uses the filesystem adapter. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.37 clean checkpoint:

| Signal | P3.37 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 274 | Four local-config-touch-write recording contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 17 / 25 | The contracts stay in the already-tested config package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 29 / 46 | The local-config touch output write now uses the filesystem adapter. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.38 clean checkpoint:

| Signal | P3.38 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 278 | Four local-config-update-write recording contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 17 / 25 | The contracts stay in the already-tested config package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 28 / 45 | The local-config update output write now uses the filesystem adapter. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.39 clean checkpoint:

| Signal | P3.39 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 282 | Four local-config-directory-stat recording contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 17 / 25 | The contracts stay in the already-tested config package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 27 / 44 | The local-config directory stat now uses the filesystem adapter. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.40 clean checkpoint:

| Signal | P3.40 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 286 | Four local-config-touch-create recording contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 17 / 25 | The contracts stay in the already-tested config package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 25 / 42 | The local-config touch create and concrete file Close classification left the scanner population. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.41 clean checkpoint:

| Signal | P3.41 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 290 | Four local-config-update-create recording contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 17 / 25 | The contracts stay in the already-tested config package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 23 / 40 | The local-config update create and concrete file Close classification left the scanner population. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.42 clean checkpoint:

| Signal | P3.42 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 296 | Two adapter and four local-config-directory-create contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 17 / 25 | The contracts stay in already-tested packages. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 22 / 40 | The direct config mkdir becomes injected while the exact system adapter call remains in the population. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.43 clean checkpoint:

| Signal | P3.43 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 304 | Two adapter and six tips-list contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 18 / 25 | Direct tips tests reduce the uncovered package count by one. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 21 / 40 | The direct tips read becomes injected while the exact system adapter call remains in the population. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.44 clean checkpoint:

| Signal | P3.44 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 308 | Four focused profile-editor contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 18 / 25 | The contracts stay in already-tested packages. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 19 / 38 | Both direct profile-editor process sites leave the production population. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.45 clean checkpoint:

| Signal | P3.45 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 316 | Three process-start and five browser-launcher contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 19 / 25 | Direct webservice tests reduce the uncovered package count by one. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 16 / 36 | The three browser command/start chains leave while one exact system-adapter start site enters the production population. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.46 clean checkpoint:

| Signal | P3.46 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 324 | Three filesystem-adapter and five Spring archive-path contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 19 / 25 | The contracts stay in already-tested packages. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 15 / 36 | The Spring working-directory caller leaves while one exact system-adapter site remains in the production population. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.47 clean checkpoint:

| Signal | P3.47 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 330 | Six focused shell unzip contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 19 / 25 | The contracts stay in an already-tested package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 13 / 34 | Both direct unzip directory-creation caller sites leave while the existing system-adapter site remains singular. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.48 clean checkpoint:

| Signal | P3.48 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 332 | Two focused shell unzip file-open contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 19 / 25 | The contracts stay in an already-tested package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 12 / 33 | The direct unzip file-open caller site leaves while the existing system-adapter site remains singular. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.49 clean checkpoint:

| Signal | P3.49 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 336 | Four focused shell Run process contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Packages with tests | 19 / 25 | The contracts stay in an already-tested package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 11 / 32 | The direct public-Run process-construction site leaves while the existing system-adapter site remains singular. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |

P3.50 clean checkpoint:

| Signal | P3.50 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 339 | Three focused shell unzip copy contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The archive fixture remains guarded and the copy double writes only in memory. |
| Guarded safe-writer sites | 25 | The population remains 20 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in an already-tested package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 11 / 32 | The direct `io.Copy` call is gone; the scanner still classifies the replacement adapter request because its opened file and archive reader arguments retain filesystem-effect provenance. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 64 Markdown files | 0 | Q3.4 remains held. |
| Focused comparable ratchets | 4 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.51 clean checkpoint:

| Signal | P3.51 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 346 | Seven focused filesystem-adapter and private shell unzip archive-open contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | Both archive fixtures remain guarded below `t.TempDir()`. |
| Guarded safe-writer sites | 26 | The population is 21 write and 5 copy sites after the focused system archive-open fixture. |
| Packages with tests | 19 / 25 | The contracts stay in already-tested packages. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 10 / 32 | The one direct `zip.OpenReader` caller site is replaced by one in-boundary adapter implementation site. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 65 Markdown files | 0 | Q3.4 remains held. |
| Focused comparable ratchets | 4 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.52 clean checkpoint:

| Signal | P3.52 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 352 | Six focused filesystem-adapter and private shell unzip output-close contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The archive fixture remains guarded below `t.TempDir()` and the injected pipe files do not target the repository. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in already-tested packages. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 10 / 32 | The direct output-file Close is gone; the scanner still classifies the replacement adapter request because its opened-file argument retains filesystem-effect provenance. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 66 Markdown files | 0 | Q3.4 remains held. |
| Focused comparable ratchets | 4 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.53 clean checkpoint:

| Signal | P3.53 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 358 | Six focused filesystem-adapter and private shell unzip entry-reader-close contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The archive fixture remains guarded below `t.TempDir()`. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in already-tested packages. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 10 / 32 | The direct entry-reader Close and its adapter request are both absent from the scanner's reported violation set. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Acceptance scripts | 4 / 4 | All host flows pass. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 68 Markdown files | 0 | Q3.4 remains held. |
| Focused comparable ratchets | 4 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.54 clean checkpoint (product `3bd07e9`, audit apparatus `ce736a2`):

| Signal | P3.54 | Interpretation |
| --- | ---: | --- |
| Test functions | 360 | Two focused private shell unzip archive-close contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The archive fixture remains guarded below `t.TempDir()`. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in an already-tested package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 10 / 32 | The concrete deferred archive Close and its adapter request are both absent from the scanner's reported violation set. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 72 Markdown files | 0 | Q3.4 remains held. |
| Complete comparable ratchets | 5 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.55 clean checkpoint (product `3972fc6`):

| Signal | P3.55 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 364 | Four focused private plugin-diagrams export-process contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The focused contracts use only temporary working directories and launch no external program. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in an already-tested package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 9 / 31 | The isolated `structurizr-cli export` construction/execution site moved through the existing process adapter. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 73 Markdown files | 0 | Q3.4 remains held. |
| Complete comparable ratchets | 5 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.56 clean checkpoint (product `6f4bc37`):

| Signal | P3.56 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 369 | Five focused private plugin-diagrams open-process contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The focused contracts use only a recorder and source AST; they launch no external program. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in an already-tested package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 8 / 30 | The isolated macOS `open` construction/execution site moved through the existing process adapter. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 74 Markdown files | 0 | Q3.4 remains held. |
| Complete comparable ratchets | 5 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.57 clean checkpoint (product `9164e06`):

| Signal | P3.57 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 375 | Six focused private plugin-diagrams Graphviz process/write contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The focused contracts record in memory and launch no external program. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in an already-tested package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 7 / 29 | The isolated Graphviz `dot` construction/execution site moved through the existing process adapter. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 75 Markdown files | 0 | Q3.4 remains held. |
| Complete comparable ratchets | 5 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.58 clean checkpoint (product `4376e05`):

| Signal | P3.58 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 380 | Five focused process/Maven recording contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The focused contracts record in memory and launch no external program. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in already-tested packages. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 15 / 37 | The Maven caller-side `logger.StdOut()` selection moved behind the process dependency; the unchanged import-aware scanner now follows the new standard-output capability through existing `process.System()` selections. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 76 Markdown files | 0 | Q3.4 remains held. |
| Complete comparable ratchets | 5 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.59 clean checkpoint (product `9f56714`):

| Signal | P3.59 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 384 | Four focused filesystem and shell archive-entry-open contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The focused contracts use guarded temporary archives and recording boundaries; they launch no external program. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in already-tested packages. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 15 / 37 | The unchanged scanner does not catalog `archive/zip.File.Open`, so the exact population holds while the operation moves behind the filesystem dependency. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 76 Markdown files | 0 | Q3.4 remains held. |
| Complete comparable ratchets | 5 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.60 clean checkpoint (product `8503601`):

| Signal | P3.60 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 386 | One focused process-selector contract and one shell Run logging/order contract were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The focused contracts record in memory, launch no external program, touch no network, and write no repository fixture. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in already-tested packages. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 14 / 36 | Public shell Run now selects the exact system runner without inheriting the unused `os.Stdout` capability. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 78 Markdown files | 0 | Q3.4 remains held. |
| Complete comparable ratchets | 5 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.61 clean checkpoint (product `5b678ab`):

| Signal | P3.61 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 388 | Two focused shell Git selector, request, order, and legacy-error contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The focused contracts record in memory, launch no external program, touch no network, and write no repository fixture. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in an already-tested package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 10 / 32 | The four public shell Git wrappers now select the exact system runner without inheriting the unused `os.Stdout` capability. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 79 Markdown files | 0 | Q3.4 remains held. |
| Complete comparable ratchets | 5 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.62 clean checkpoint (product `b0d324a`):

| Signal | P3.62 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 388 | Four existing focused profile-editor selector, request, stream, and error contracts were strengthened. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The focused contracts record in memory, launch no external program, touch no network, and write no repository fixture. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in an already-tested package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 9 / 31 | The private profile-editor selector now carries the exact system runner without inheriting the unused process-dependency `os.Stdout` capability. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 80 Markdown files | 0 | Q3.4 remains held. |
| Complete comparable ratchets | 5 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.63 clean checkpoint (product `1bce06f`):

| Signal | P3.63 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 388 | Four existing focused plugin-diagrams selector and process-error contracts were strengthened; three recording doubles now preserve complete caller-owned dependencies. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The focused contracts record in memory, launch no external program, touch no network, and write no repository fixture. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in an already-tested package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 8 / 30 | The three private plugin-diagrams selectors now carry the exact system runner without giving their one Structurizr caller composition the unused process-dependency `os.Stdout` capability. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 81 Markdown files | 0 | Q3.4 remains held. |
| Complete comparable ratchets | 5 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.64 clean checkpoint (product `0b96f10`):

| Signal | P3.64 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 388 | Five existing focused browser-launcher selector, platform, error, and public-composition contracts were strengthened; the recording double now preserves the complete caller-owned dependency. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The focused contracts record in memory, launch no external program, touch no network, and write no repository fixture. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in an already-tested package. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 7 / 29 | The private browser-launcher selector now carries the exact system runner without giving public `OpenBrowser` the unused process-dependency `os.Stdout` capability. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 81 Markdown files | 0 | Q3.4 remains held. |
| Complete comparable ratchets | 5 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.65 clean checkpoint (product `40cece5`):

| Signal | P3.65 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 388 | Two existing focused filesystem-open and Unzip file-open contracts were strengthened without changing the test denominator. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The focused contracts record in memory, launch no external program, touch no network, and write no repository fixture. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in already-tested packages. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 5 / 27 | Unzip's unchanged destination open now returns through `filesystem.File`, so the following adapter `Copy` and `Close` requests no longer inherit concrete `*os.File` provenance. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 83 Markdown files | 0 | Q3.4 remains held. |
| Complete comparable ratchets | 5 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P3.66 clean checkpoint (product `c627a3e`):

| Signal | P3.66 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 388 | Existing focused process-standard-output and Maven recording contracts were strengthened without changing the test denominator. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The focused contracts record in memory, launch no external program, touch no network, and write no repository fixture. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 19 / 25 | The contracts stay in already-tested packages. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 4 / 27 | Public Maven `RunOn` now selects the exact system runner and standard output without inheriting concrete `os.Stdout` provenance through the complete process dependency; the adapter helper keeps the production-effect denominator at 27. |
| Declared seam swap tests | 7 / 8 | This move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 84 Markdown files | 0 | Q3.4 remains held. |
| Complete comparable ratchets | 5 improved, 2 held, 0 regressed | Q1.3 is the single not-comparable selected criterion because the declared clock and server adapters remain absent. |

P4.1 clean checkpoint (product `91422eb`):

| Signal | P4.1 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 394 | Four clock-adapter contracts and two additional Spring archive-path recording contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Unsafe direct test writes | 0 | The focused contracts record in memory, launch no external program, touch no network, and write no repository fixture. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites. |
| Packages with tests | 20 / 26 | The new clock adapter is a tested package; the six packages without test files are unchanged. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 3 / 27 | Spring's archive timestamp now crosses the exact declared clock adapter; the clock adapter is present and only the server adapter is missing. |
| Declared seam swap tests | 7 / 8 | This focused clock move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Production scripts without meta-tests | 0 / 12 | Q0.8 remains improved. |
| Reachable manual L1/L2 rows | 6 / 6 | Evidence reachability is unchanged. |
| Baseline numeric debt leaves | 228 | The measurement instrument is unchanged. |
| Claim phrases in 85 Markdown files | 0 | Q3.4 remains held. |
| Complete comparable ratchets | 5 improved, 2 held, 0 regressed | Q1.3 remains the single not-comparable selected criterion because the declared server adapter is absent. |

P4.2 clean checkpoint (product `4c0d97a`):

| Signal | P4.2 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 401 | Seven exact clock-sleep and Kibana retry contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites, with zero unsafe direct test writes. |
| Packages with tests | 20 / 26 | The six packages without test files are unchanged. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 2 / 27 | Kibana retry sleep crosses the clock adapter; only the two webservice server calls remain. |
| Declared seam swap tests | 7 / 8 | The focused retry move adds no inventory seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Claim phrases in 86 Markdown files | 0 | Q3.4 remains held. |
| Full audit | exit 1, 14 non-passing | The audit never exits 2 and records zero dirty paths. |
| Complete ratchets | 5 improved, 2 held, 0 regressed, 1 non-comparable | Q1.3 remains non-comparable only because the server adapter is absent. |

P4.3 clean checkpoint (product `9601d29`):

| Signal | P4.3 | Interpretation |
| --- | ---: | --- |
| Absolute L0 PASS | 8 / 8 | Every absolute gate remains closed. |
| Test functions | 414 | Thirteen exact server-adapter and webservice start/stop contracts were added. |
| Skipped tests | 0 | Q0.6 remains improved. |
| Guarded safe-writer sites | 26 | The population holds at 21 write and 5 copy sites, with zero unsafe direct test writes. |
| Packages with tests | 21 / 27 | The new server adapter is tested; six packages still have no test files. |
| Process-exiting calls outside `main` | 0 | Q1.2 stays closed. |
| Direct external effects outside adapters | 0 / 27 | All five declared adapter paths are present and valid. |
| Declared seam swap tests | 7 / 8 | The server effect move does not claim the separate loopback-address seam. |
| Executable mutation harnesses | 0 / 8 | P5 remains queued. |
| Claim phrases in 87 Markdown files | 0 | Q3.4 remains held. |
| Full audit | exit 1, 13 non-passing | One prior Q1.3 finding closes; the audit never exits 2 and records zero dirty paths. |
| Complete ratchets | 6 improved, 2 held, 0 regressed, 0 non-comparable | Every selected ratchet is now comparable. |

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

Status: complete. Moves 1 through 66 are complete.

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
regressed, and one not-comparable ratchet.

Commit: `2566438`.

The clean full audit from commit `2566438` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 35 routes only `structurizr.RunWithOutputToFile`'s direct
`os.WriteFile(outputFile, out.Bytes(), 0644)` operation through the existing
filesystem adapter. The exported function and every caller remain unchanged.
It still creates `var out bytes.Buffer`, then `var stderr bytes.Buffer`, assigns
the command's stdout and stderr buffers, invokes `command.Run()` once, and
returns the exact non-nil command error before any write. On success,
production selects `filesystem.System()` only for a private complete
structurizr-output-write dependency and makes one adapter `WriteFile` attempt
with the exact caller path, exact captured stdout bytes, and `0644` mode before
discarding the write error and returning nil.

Four new top-level structurizr contracts bring the suite to 266 tests and the
tested population to 17 of 25 packages. They prove complete system dependency
selection and delivery, exact arbitrary caller path delivery without cleaning
or normalization, non-empty recorded write populations, exact empty,
representative, non-ASCII, and arbitrary stdout bytes, exact `0644` mode, one
write attempt, exact write-error identity at the private boundary, safe zero
behavior without developer-path access, and no unrelated adapter operation.
The production-composition contract neither mutates the real filesystem nor
runs a command. No adapter operation, established complete recording double,
inventory label, seam driver, mutation harness, public API, caller, command
construction, command execution, buffer assignment, process error handling,
plugin-diagram behavior, config, template, file, Maven, or tips behavior
changed. Q0.6 remains at 22 guarded safe-writer sites, 17 write and 5 copy,
with zero unsafe direct writes. Q1.3 moves from 32 of 49 to 31 of 48; Q1.1
moves from 9 to 8 untested packages; Q1.2 stays zero, Q1.4 stays 7 of 8, and
exact Q2.1 stays 0 of 8.

The move-35 implementation gate passed focused structurizr/filesystem and
relevant process, command, file, template, config, Maven, tips, and context
package tests, API/CLI and subprocess compatibility, all four host flows, full
preflight, test, install, the standalone 62-control launcher contract,
uncached and race tests, vet, and the 15-control audit meta-suite. The focused
seven-criterion audit exited 1 for documented findings, never 2, with four
improved, two held, zero regressed, and one not-comparable ratchet.

Commit: `527a8b9`.

The clean full audit from commit `527a8b9` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 36 routes only exported `maven.WriteGraphStyles`' direct
`os.WriteFile(stylesFile, jsonStyles, 0644)` operation through the existing
filesystem adapter. The exported function and every caller remain unchanged.
It still evaluates exact `json.MarshalIndent(styles, "", "    ")` first and
returns its error before composing exact
`file.Path("%s/target/dependency-graph-styles.json", projectPath)`. The existing
`file.Exists(stylesFile)` selection still returns nil without a write when the
path is present. When absent, production selects `filesystem.System()` only
for a private complete graph-style-write dependency and attempts the existing
adapter `WriteFile` operation once with the exact caller-composed path, exact
marshaled bytes, and mode `0644`, returning the exact write error.

Four new top-level Maven contracts bring the suite to 270 tests. They prove
complete system dependency selection and delivery, exact arbitrary
caller-composed paths without cleaning or normalization, non-empty recorded
write populations, exact empty, representative, non-ASCII, and arbitrary
output bytes, exact `0644` mode, one write attempt, exact write-error identity
at the private boundary, safe zero behavior without developer-path access, and
no unrelated adapter operation. The production-composition contracts neither
mutate the real filesystem nor run Maven, dot, or another command. No adapter
operation, established complete recording double, inventory label, seam
driver, mutation harness, public API, caller, JSON marshaling, path
construction, existence selection, Graph filtering or style mutation, Graph
arguments, `RunOn`, Maven/dot execution, command behavior, structurizr, tips,
config, template, file, HTTP, Bitbucket, Wpost, or other completed effect
changed. Q0.6 remains at 22 guarded safe-writer sites, 17 write and 5 copy,
with zero unsafe direct writes. Q1.3 moves from 31 of 48 to 30 of 47; Q1.1
stays at 8 untested packages; Q1.2 stays zero, Q1.4 stays 7 of 8, and exact
Q2.1 stays 0 of 8.

The move-36 implementation gate passed focused Maven/filesystem/process and
relevant command, file, template, config, tips, structurizr, and context
package tests, API/CLI and subprocess compatibility, all four host flows, full
preflight, test, install, the standalone 62-control launcher contract,
uncached and race tests, vet, and the 15-control audit meta-suite. The first
standalone launcher run hit the documented nested signal-interruption partial-
raw-log timing flake; its immediate complete rerun passed all 62 controls. The
focused seven-criterion audit exited 1 for documented findings, never 2, with
four improved, two held, zero regressed, and one not-comparable ratchet.

Commit: `886dff0`.

The clean full audit from commit `886dff0` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 37 routes only `LocalConfigDir.TouchFile`'s direct
`os.WriteFile(configFilePath, d, 0644)` operation through the existing
filesystem adapter. The exported method and every caller remain unchanged. It
still evaluates `CheckOrCreateConfigDir()` and its early error, exact
`FilePath()`, zero-value `LocalConfiguration{}`, exact
`defaultCloudConfigUrl` assignment, `yaml.Marshal(&config)` and its early
error, and exact creation log in the same order. It then preserves the
`os.Create(configFilePath)` attempt and its early error before production
selects `filesystem.System()` only for a private complete touch-write
dependency. The adapter receives the exact caller-composed path, exact YAML
bytes, and mode `0644` once. Its exact write error remains the result without a
Close attempt; successful writes still return the exact `f.Close()` result.

Four new top-level config contracts bring the suite to 274 tests. They prove
complete system dependency selection and delivery, exact arbitrary
caller-composed paths without cleaning or normalization, non-empty recorded
write populations, exact empty, representative, non-ASCII, and arbitrary
output bytes, exact `0644` mode, one write attempt, exact write-error identity
at the private boundary, safe zero behavior without developer-path access, and
no unrelated adapter operation. The production-composition contracts do not
mutate the real filesystem or run another command. No adapter operation,
established complete recording double, inventory label, seam driver, mutation
harness, public API, caller, directory evaluation, path construction, config
initialization or defaults, YAML marshaling, logging, file creation or closing,
`UpdateLocalConfig`, context/profile behavior, ProjectConfig, CloudConfig,
Maven, structurizr, tips, template, file, HTTP, Bitbucket, Wpost, or other
completed effect changed. Q0.6 remains at 22 guarded safe-writer sites, 17
write and 5 copy, with zero unsafe direct writes. Q1.3 moves from 30 of 47 to
29 of 46; Q1.1 stays at 8 untested packages; Q1.2 stays zero, Q1.4 stays 7 of
8, and exact Q2.1 stays 0 of 8.

The move-37 implementation gate passed focused config/filesystem/process and
relevant command, context, file, template, Maven, tips, structurizr,
Bitbucket, HTTP, Kibana, Spring, and shell package tests, API/CLI and subprocess
compatibility, all four host flows, full preflight, test, install, the
standalone 62-control launcher contract, uncached and race tests, vet, and the
15-control audit meta-suite. The focused seven-criterion audit exited 1 for
documented findings, never 2, with four improved, two held, zero regressed,
and one not-comparable ratchet. The first handoff-state launcher run hit the
documented nested signal-interruption partial-raw-log timing flake; its
immediate complete rerun passed all 62 controls.

Commit: `89f43aa`.

The clean full audit from commit `89f43aa` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 38 routes only `LocalConfigDir.UpdateLocalConfig`'s direct
`os.WriteFile(configFilePath, d, 0644)` operation through the existing
filesystem adapter. The exported method and every caller remain unchanged. It
still evaluates `CheckOrCreateConfigDir()` and its early error, exact
`FilePath()`, the exact caller-supplied `config`, `yaml.Marshal(&config)` and
its early error, and exact update log in the same order. It then preserves the
`os.Create(configFilePath)` attempt and its early error before production
selects `filesystem.System()` only for a private complete update-write
dependency. The adapter receives the exact caller-composed path, exact YAML
bytes, and mode `0644` once. Its exact write error remains the result without a
Close attempt; successful writes still return the exact `f.Close()` result.

Four new top-level config contracts bring the suite to 278 tests. They prove
complete system dependency selection and delivery, exact arbitrary
caller-composed paths without cleaning or normalization, non-empty recorded
write populations, exact empty, representative, non-ASCII, and arbitrary
output bytes, exact `0644` mode, one write attempt, exact write-error identity
at the private boundary, safe zero behavior without developer-path access, and
no unrelated adapter operation. The production-composition contracts do not
mutate the real filesystem or run another command. No adapter operation,
established complete recording double, inventory label, seam driver, mutation
harness, public API, caller, directory evaluation, path construction,
caller-supplied config, YAML marshaling, logging, file creation or closing,
`TouchFile`, context/profile behavior, ProjectConfig, CloudConfig, Maven,
structurizr, tips, template, file, HTTP, Bitbucket, Wpost, or other completed
effect changed. Q0.6 remains at 22 guarded safe-writer sites, 17 write and 5
copy, with zero unsafe direct writes. Q1.3 moves from 29 of 46 to 28 of 45;
Q1.1 stays at 8 untested packages; Q1.2 stays zero, Q1.4 stays 7 of 8, and
exact Q2.1 stays 0 of 8.

The move-38 implementation gate passed focused config/filesystem/process and
relevant command, context, file, template, Maven, tips, structurizr,
Bitbucket, HTTP, Kibana, Spring, and shell package tests, API/CLI and subprocess
compatibility, all four host flows, full preflight, test, install, the
standalone 62-control launcher contract, uncached and race tests, vet, the
15-control audit meta-suite, and empty-HOME count-2. The first complete
`make test` run hit the documented nested signal-interruption partial-raw-log
timing flake; its immediate complete rerun passed all 62 launcher controls.
The focused seven-criterion audit exited 1 for documented findings, never 2,
with four improved, two held, zero regressed, and one not-comparable ratchet.

Commit: `fd45ebf`.

The clean full audit from commit `fd45ebf` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 39 routes only `LocalConfigDir.CheckOrCreateConfigDir`'s direct
`os.Stat(dir)` operation through the existing filesystem adapter. The exported
method and every caller remain unchanged. It still evaluates exact
`localCfgDir.Implementation().Path`, makes one stat attempt with that arbitrary
string, and applies the same `os.IsNotExist(err)` decision. Only a not-exist
error attempts exact `os.Mkdir(dir, 0755)` and returns its exact non-nil error;
nil stat results, existing-directory results, and arbitrary other stat errors
still return nil without a directory creation attempt. Production selects
`filesystem.System()` only for a private complete local-config-directory-stat
dependency, and the adapter preserves the exact `fs.FileInfo` and error.

Four new top-level config contracts bring the suite to 282 tests. They prove
complete system dependency selection and delivery, exact arbitrary directory
paths without cleaning or normalization, a non-empty recorded stat population,
one stat attempt, exact `fs.FileInfo` and error identity for nil, existing,
not-exist, and arbitrary error results, safe zero behavior without
developer-path access, and no unrelated adapter operation. The production-
composition contract does not mutate the real filesystem or run another
command. No adapter operation, established complete recording double,
inventory label, seam driver, mutation harness, public API, caller, directory
evaluation, `os.IsNotExist`, `os.Mkdir`, directory mode or error selection,
file creation or closing, `TouchFile`, `UpdateLocalConfig`, context/profile
behavior, ProjectConfig, CloudConfig, Maven, structurizr, tips, template, file,
HTTP, Bitbucket, Wpost, or other completed effect changed. Q0.6 remains at 22
guarded safe-writer sites, 17 write and 5 copy, with zero unsafe direct writes.
Q1.3 moves from 28 of 45 to 27 of 44; Q1.1 stays at 8 untested packages; Q1.2
stays zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

The move-39 implementation gate passed focused config/filesystem/process and
relevant command, context, file, template, Maven, tips, structurizr,
Bitbucket, HTTP, Kibana, Spring, and shell package tests, API/CLI and subprocess
compatibility, all four host flows, full preflight, test, install, the
standalone 62-control launcher contract, uncached and race tests, vet, the
15-control audit meta-suite, and empty-HOME count-2. A standalone launcher run
and a later separate `make test-agent-start` run each hit the documented nested
signal-interruption partial-raw-log timing flake; their immediate complete
reruns and the full preflight run passed all 62 launcher controls. The focused
seven-criterion audit exited 1 for documented findings, never 2, with four
improved, two held, zero regressed, and one not-comparable ratchet.

Commit: `d2e330b`.

The clean full audit from commit `d2e330b` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 40 routes only `LocalConfigDir.TouchFile`'s direct
`os.Create(configFilePath)` operation through the existing filesystem adapter.
The exported method and every caller remain unchanged. It still evaluates
`CheckOrCreateConfigDir()` and its early error, exact `FilePath()`, the
zero-value `LocalConfiguration{}`, exact `defaultCloudConfigUrl` assignment,
`yaml.Marshal(&config)` and its early error, and the exact creation log in the
same order. Production then selects `filesystem.System()` only for a private
complete local-config-touch-create dependency. The adapter makes one create
attempt with the exact caller-composed path and preserves the exact
`filesystem.File` and error.

Four new top-level config contracts bring the suite to 286 tests. They prove
complete system dependency selection and delivery, exact arbitrary file paths
without cleaning or normalization, a non-empty recorded create population,
one create attempt, exact file and error identity for nil, successful,
arbitrary-error, and unusual combined results, safe zero behavior without
developer-path access, the dedicated file double's exact Close result, and no
unrelated adapter operation. The production-composition contract does not
mutate the real filesystem or run another command. A non-nil create error
still returns before the completed touch write; create success still invokes
the existing touch-write composition once with the exact path, YAML bytes, and
mode `0644`; a write error still returns without Close; and write success still
returns the exact final `f.Close()` result. No adapter operation, established
complete recording double, inventory label, seam driver, mutation harness,
public API, caller, directory evaluation, path construction, config
initialization, default URL, YAML marshaling, logging, touch write, update
behavior, context/profile behavior, ProjectConfig, CloudConfig, Maven,
structurizr, tips, template, file, HTTP, Bitbucket, Wpost, or other completed
effect changed. Q0.6 remains at 22 guarded safe-writer sites, 17 write and 5
copy, with zero unsafe direct writes. Q1.3 moves from 27 of 44 to 25 of 42 as
the direct create and its concrete `*os.File` Close classification leave while
the same lifecycle call remains; Q1.1 stays at 8 untested packages; Q1.2 stays
zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

The move-40 implementation gate passed focused config/filesystem/process and
relevant command, context, file, template, Maven, tips, structurizr,
Bitbucket, HTTP, Kibana, Spring, and shell package tests, API/CLI and subprocess
compatibility, all four host flows, full preflight, test, install, the
standalone 62-control launcher contract, uncached and race tests, vet, the
15-control audit meta-suite, and empty-HOME count-2. Complete preflight and
test runs encountered the documented nested signal-interruption partial-raw-
log timing flake before complete reruns passed all 62 launcher controls. The
focused seven-criterion audit exited 1 for documented findings, never 2, with
four improved, two held, zero regressed, and one not-comparable ratchet.

Commit: `b232dda`.

The clean full audit from commit `b232dda` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 41 routes only `LocalConfigDir.UpdateLocalConfig`'s direct
`os.Create(configFilePath)` operation through the existing filesystem adapter.
The exported method and every caller remain unchanged. It still evaluates
`CheckOrCreateConfigDir()` and its early error, exact `FilePath()`, the exact
caller-supplied config, `yaml.Marshal(&config)` and its early error, and the
exact update log in the same order. Production then selects
`filesystem.System()` only for a private complete local-config-update-create
dependency. The adapter makes one create attempt with the exact
caller-composed path and preserves the exact `filesystem.File` and error.

Four new top-level config contracts bring the suite to 290 tests. They prove
complete system dependency selection and delivery, exact arbitrary file paths
without cleaning or normalization, a non-empty recorded create population,
one create attempt, exact file and error identity for nil, successful,
arbitrary-error, and unusual combined results, safe zero behavior without
developer-path access, the dedicated file double's exact Close result, and no
unrelated adapter operation. The production-composition contract does not
mutate the real filesystem or run another command. A non-nil create error
still returns before the completed update write; create success still invokes
the existing update-write composition once with the exact path, YAML bytes,
and mode `0644`; a write error still returns without Close; and write success
still returns the exact final `f.Close()` result. No adapter operation,
established complete recording double, inventory label, seam driver, mutation
harness, public API, caller, directory evaluation, path construction,
caller-supplied config, YAML marshaling, logging, update write, touch behavior,
context/profile behavior, ProjectConfig, CloudConfig, Maven, structurizr,
tips, template, file, HTTP, Bitbucket, Wpost, or other completed effect
changed. Q0.6 remains at 22 guarded safe-writer sites, 17 write and 5 copy,
with zero unsafe direct writes. Q1.3 moves from 25 of 42 to 23 of 40 as the
direct create and its concrete `*os.File` Close classification leave while the
same lifecycle call remains; Q1.1 stays at 8 untested packages; Q1.2 stays
zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

The move-41 implementation gate passed focused config/filesystem/process and
relevant command, context, file, template, Maven, tips, structurizr,
Bitbucket, HTTP, Kibana, Spring, and shell package tests, API/CLI and subprocess
compatibility, all four host flows, full preflight, test, install, the
standalone 62-control launcher contract, uncached and race tests, vet, the
15-control audit meta-suite, and empty-HOME count-2. The first standalone
launcher run hit the documented nested signal-interruption partial-raw-log
timing flake; its immediate complete rerun, full preflight, and full test run
passed all 62 launcher controls. The focused seven-criterion audit exited 1
for documented findings, never 2, with four improved, two held, zero
regressed, and one not-comparable ratchet.

Commit: `6928a72`.

The clean full audit from commit `6928a72` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 42 routes only `LocalConfigDir.CheckOrCreateConfigDir`'s direct
`os.Mkdir(dir, 0755)` operation through one new distinct filesystem adapter
`Mkdir` operation. The exported method and every caller remain unchanged. It
still evaluates exact `localCfgDir.Implementation().Path`, makes the completed
one-attempt adapter `Stat` call with that arbitrary string, and applies the
same exact `os.IsNotExist(err)` decision. Nil stat results, existing-directory
results, and arbitrary non-not-exist stat errors still return nil without a
mkdir. Only a not-exist error selects `filesystem.System()` for a private
complete local-config-directory-create dependency and makes one single-level
mkdir attempt with the exact `dir` and mode `0755`; its exact error remains the
method result and success returns nil.

Six new top-level contracts bring the suite to 296 tests. The two filesystem
adapter contracts prove exact arbitrary paths and modes, one recorded attempt,
non-empty populations, exact error identity, safe zero behavior, and the
system implementation's direct single-level semantics: successful creation,
existing-directory failure, missing-parent failure, and no recursive parent
creation. Four local-config-directory-create contracts prove complete system
dependency selection and delivery, the exact receiver-selected arbitrary path
without cleaning or normalization, exact `0755`, one attempt, nil and
arbitrary error results, safe zero behavior without developer-path access,
non-empty populations, and no unrelated adapter operation. The production-
composition contract does not mutate the real filesystem or run another
command. All 34 pre-existing complete `filesystem.FileSystem` test doubles
implement the new method only to remain complete and reject it as unrelated;
the dedicated adapter and local-config doubles record it. Established
`MkdirAll`, stat, every other adapter operation, inventory, seam drivers,
mutation harnesses, public API, callers, local-config create/write/Close
lifecycles, context/profile behavior, ProjectConfig, CloudConfig, Maven,
structurizr, tips, template, file, HTTP, Bitbucket, Wpost, and every other
completed effect remain unchanged. Q0.6 stays at 22 guarded safe-writer sites,
17 write and 5 copy, with zero unsafe direct writes. Q1.3 improves from 23 of
40 to 22 of 40: the direct config mkdir becomes injected, while the required
direct `os.Mkdir` system implementation remains inside the exact declared
adapter and therefore keeps the scanner population at 40. Q1.1 stays at 8
untested packages; Q1.2 stays zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0
of 8.

The move-42 implementation gate passed focused config/filesystem/process and
relevant command, context, file, template, Maven, tips, structurizr,
Bitbucket, HTTP, Kibana, Spring, and shell package tests, API/CLI and subprocess
compatibility, all four host flows, full preflight, test, install, the
standalone 62-control launcher contract, uncached and race tests, vet, the
15-control audit meta-suite, and empty-HOME count-2. The complete preflight
rerun used contract-correct external Go and golangci-lint caches and tool
discovery after two runner-only setup failures; every repository control then
passed. The focused seven-criterion audit exited 1 for documented findings,
with four improved, two held, zero regressed, and one not-comparable ratchet.

Commit: `b2d37cc`.

The clean full audit from commit `b2d37cc` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 43 routes only `tips.List`'s direct `os.ReadDir(LocalDir(gitCfg))`
operation through one new distinct filesystem adapter `ReadDirEntries`
operation. The exported `List` signature, `LocalDir`, `TipsDir`, every caller,
and `tipsShowCmd` remain unchanged. `List` still evaluates the exact
`LocalDir(gitCfg)` path through its established `file.Path("%s/%s", ...)`
behavior, selects `filesystem.System()` only for a private complete tips-list
dependency, and makes one read attempt. A read error still returns the nil
named result and that exact error without inspecting delivered entries. A
successful read still inspects entries once in delivered order, excludes
directories before reading their names, includes only non-directories whose
exact names have the case-sensitive `.md` suffix, preserves included entry
identity, and returns nil for empty or all-filtered populations.

Eight new top-level contracts bring the suite to 304 tests and direct tests in
`pkg/tips`, improving package coverage to 18 of 25. The two filesystem adapter
contracts, together with the extended zero-value contract, prove exact
arbitrary paths, one recorded attempt, non-empty populations, exact nil,
non-nil empty, representative, partial-entry, and arbitrary-error results,
ordered entry identity, exact `ErrNoFilesystem`, and the system
implementation's direct `os.ReadDir` filename sorting, names, directory
classification, and missing-path failure. Six tips-list contracts prove
complete system dependency selection, exact `LocalDir` evaluation and path,
one read attempt, exact error identity and nil-on-error behavior without entry
inspection, delivered-order filtering and entry identity, short-circuit
`IsDir`/`Name` use without `Info` or `Type`, nil empty/all-filtered results,
non-empty populations, safe zero behavior, and no unrelated adapter operation.
The production-composition contract does not read the real filesystem or run
another command. All 35 pre-existing complete `filesystem.FileSystem` test
doubles implement the new method only to remain complete: the dedicated
adapter double records it and the other 34 reject it as unrelated; the new
tips-list double records only its focused operation. Established `ReadDir`,
its `ioutil.ReadDir` mapping and `[]fs.FileInfo` result, every other adapter
operation, inventory, seam drivers, mutation harnesses, public API, CLI,
callers, tips-show reads and rendering, config, local-config, Maven,
structurizr, template, file, HTTP, Bitbucket, Wpost, and every other completed
effect remain unchanged.

Q0.6 stays closed with zero skipped tests and zero unsafe direct test writes;
its safe-writer population becomes 24 guarded sites, 19 write and 5 copy,
because the system adapter contract adds two guarded temporary-fixture writes.
Q1.1 improves from 8 to 7 untested packages. Q1.3 improves from 22 of 40 to 21
of 40: the direct tips read becomes injected, while the required direct
`os.ReadDir` system implementation remains inside the exact declared adapter
and therefore keeps the scanner population at 40. Q1.2 stays zero, Q1.4 stays
7 of 8, and exact Q2.1 stays 0 of 8.

The move-43 implementation gate passed focused tips/filesystem/process and
relevant command, context, config, file, template, Maven, structurizr,
Bitbucket, HTTP, Kibana, Spring, and shell package tests, API/CLI and subprocess
compatibility, all four host flows and their meta-contracts, full preflight,
test, install, the standalone 62-control launcher contract, uncached and race
tests, vet, the 15-control audit meta-suite, and empty-HOME count-2. The full
preflight rerun used an external golangci-lint cache after one runner-only
cache-permission failure; every repository control then passed. The focused
seven-criterion audit exited 1 for documented findings, never 2, with four
improved, two held, zero regressed, and one not-comparable ratchet.

Commit: `85f4c2b`.

The clean full audit from commit `85f4c2b` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 44 extends only the existing process request with `Stdin io.Reader`, maps
that value directly to `exec.Cmd.Stdin`, and routes only the profile edit
branch's direct `exec.Command` construction and `cmd.Run` through the process
adapter. The process adapter's name, arguments, directory, stdout, stderr,
runner, dependency, execution, system, error, and safe-zero behavior remain
unchanged. The profile branch still reads `EDITOR` once, replaces only an
empty value with exact `vim`, evaluates the arbitrary local-config path once,
then selects one private complete production composition containing
`process.System()`, exact `os.Stdin`, and exact `os.Stdout`. It makes one
adapter execution attempt with the exact editor as `Name`, the exact path as
the only argument, an empty directory, exact stdin/stdout identities, and nil
stderr. The exact process error still returns before sync, reset, or print;
success preserves every later branch, condition, order, and result.

Four new top-level profile-editor contracts bring the suite to 308 tests. They
prove complete production dependency and stream selection, exact non-empty
and empty `EDITOR` behavior, the `vim` fallback, arbitrary local-config path
bytes, exact stream identities, empty directory, nil stderr, one non-empty
recorded attempt, exact error identity, safe zero behavior, and rejection of
an empty recording population without starting a real or unrelated process.
The extended process contracts prove exact stdin reader identity reaches a
recording runner without reads or extra attempts, direct system execution
receives the delivered input, and all prior command fields, exact errors, and
safe zero behavior remain intact. No command object, registration, public API,
inventory, mutation harness, or other process, filesystem, HTTP, tips, config,
Maven, structurizr, Bitbucket, Wpost, local-config, or supervisor operation
changed.

Q0.6 stays closed at 24 guarded safe-writer sites, 19 write and 5 copy, with
zero skipped tests and zero unsafe direct test writes. Q1.1 holds at 7 of 25
untested packages, Q1.2 stays zero, Q1.3 improves from 21 of 40 to 19 of 38 as
the two direct profile process sites leave, Q1.4 stays 7 of 8, and exact Q2.1
stays 0 of 8.

The move-44 implementation gate passed focused process, command, config,
context, Maven, tips, filesystem, file, template, structurizr, Bitbucket,
HTTP, Kibana, Spring, shell, local-config, and caller package tests; API/CLI
and subprocess compatibility; all four host flows and their meta-contracts;
full preflight, test, install, the standalone 62-control launcher contract,
uncached and race tests, vet, the 15-control audit meta-suite, and empty-HOME
count-2. Generated compatibility and audit reports and all Go and linter caches
remained outside the measured tree. The focused seven-criterion audit exited 1
for documented findings, never 2, with four improved, two held, zero regressed,
and one not-comparable ratchet.

Commit: `0e10282`.

The clean full audit from commit `0e10282` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 45 extends only the existing complete process request with boolean
`Start`, maps every established command field exactly as before, and selects
the direct `exec.Cmd.Start()` result only when that value is true; false keeps
the existing single `exec.Cmd.Run()` attempt. `Runner`, `Dependencies`,
`Execute`, `System`, safe-zero behavior, name, arguments, directory, stdin,
stdout, stderr, and every established process caller remain unchanged. Only
the three direct `exec.Command(...).Start()` branches in `OpenBrowser` move
behind one private complete browser-launcher composition containing
`process.System()` and exact `runtime.GOOS`. Linux still requests `xdg-open`
with the URL as its only argument, Windows still requests `rundll32` with
exact first argument `url.dll,FileProtocolHandler` and the URL second, and
Darwin still requests `open` with the URL as its only argument. Every
supported request has start mode true, empty directory, nil streams, one
adapter attempt, and the exact start error. Unsupported platforms still
return exact text `unsupported platform` without a process attempt.

Eight new top-level contracts bring the suite to 316 tests and add direct
tests to `pkg/webservice`, improving package coverage to 19 of 25. The process
adapter contracts prove complete start-mode identity reaches a recording
runner, false retains direct synchronous-run exit behavior, true returns
without observing the child exit, exact direct start errors remain unwrapped,
and zero dependencies remain safe. Five private browser-launcher contracts
prove complete production process and runtime selection, every exact platform
request, arbitrary URL bytes without rewriting, exact error identity for all
three supported branches, one non-empty recorded attempt, empty directory,
nil streams, unsupported-platform no-attempt behavior, safe zero behavior,
and rejection of an empty recording population without opening a real browser
or starting a server. `OpenBrowser(string) error`, its callers, every server
operation and global, public API, CLI, inventory, mutation harness, profile,
Maven, Git, filesystem, HTTP, tips, config, structurizr, Bitbucket, Wpost,
local-config, supervisor, and every completed effect remain unchanged.

Q0.6 stays closed at 24 guarded safe-writer sites, 19 write and 5 copy, with
zero skipped tests and zero unsafe direct test writes. Q1.1 improves from 7 to
6 of 25 untested packages, Q1.2 stays zero, Q1.3 improves from 19 of 38 to 16
of 36 as the three browser caller sites leave and the exact adapter start site
enters, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

The move-45 implementation gate passed focused process, webservice, command,
context, config, HTTP, Maven, shell, structurizr, profile, tips, filesystem,
file, template, Bitbucket, Kibana, Spring, local-config, and caller package
tests; API/CLI and subprocess compatibility; all four host flows and their
meta-contracts; full preflight, the standalone 62-control launcher contract,
uncached and race tests, vet, the 15-control audit meta-suite, and empty-HOME
count-2. Generated compatibility and audit reports and all Go and linter
caches remained outside the measured tree. The focused seven-criterion audit
exited 1 for documented findings, never 2, with four improved, two held, zero
regressed, and one not-comparable ratchet.

Commit: `9fcdfb5`.

The clean full audit from commit `9fcdfb5` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 46 adds only `WorkingDirectory() (string, error)` to the existing complete
filesystem interface. Its forwarding helper returns the exact safe-zero pair
`"", filesystem.ErrNoFilesystem`, otherwise makes one dependency attempt and
returns the exact string and error. The system implementation directly returns
`os.Getwd()` without normalization, fallback, retry, logging, or wrapping.
Every established filesystem method, dependency selection, error sentinel,
system operation, and caller remains unchanged; all pre-existing complete
filesystem doubles reject the unrelated operation.

Only the direct `os.Getwd()` in private `pkg/spring.archivePath` moves behind
one private complete composition selecting `filesystem.System()`.
`archivePath() (string, error)` remains the production entry. A working-
directory error returns the same empty path and exact error before the direct
`time.Now().Unix()` call. Success passes the arbitrary directory bytes
unchanged to the established `file.Path("%s/spring-%d.zip", curDir, now)` call.
The filename, timestamp unit, ordering, caller graph, Spring discovery,
download, unzip, delete, and demo-file behavior remain unchanged. No file,
HTTP, process, shell, download, unzip, delete, or clock seam was added.

Eight new top-level contracts bring the suite to 324 tests. Three focused
filesystem-adapter contracts prove exact arbitrary working-directory values
and errors pass through one recording attempt, the system result matches the
direct `os.Getwd()` operation, safe-zero identity remains exact, and an empty
recorded population fails. Five focused private Spring contracts prove
complete production selection, exact arbitrary directory composition, one
non-empty attempt, exact error identity and empty error result, safe-zero
behavior, direct Unix filename construction, rejection of an empty population,
and absence of unrelated filesystem operations. The contracts do not create a
real file or invoke HTTP, a process, download, unzip, delete, or a clock seam.

Q0.6 stays closed at 24 guarded safe-writer sites, 19 write and 5 copy, with
zero skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25,
Q1.2 stays zero, Q1.3 improves from 16 of 36 to 15 of 36 as the Spring caller
site leaves and one exact adapter site enters, Q1.4 stays 7 of 8, and exact
Q2.1 stays 0 of 8.

The move-46 implementation gate passed focused Spring, filesystem, process,
command, context, config, HTTP, Maven, shell, structurizr, profile, browser,
tips, file, template, Bitbucket, Kibana, local-config, and caller package tests;
API/CLI and subprocess compatibility; all four host flows and their
meta-contracts; full preflight, the standalone 62-control launcher contract,
uncached and race tests, vet, the 15-control audit meta-suite, and empty-HOME
count-2. Generated compatibility and audit reports and all Go and linter caches
remained outside the measured tree. One standalone launcher attempt and one
`make test` attempt hit the documented nested signal-interruption partial-raw-
log timing flake; each immediate complete rerun passed all 62 controls. The
focused seven-criterion audit exited 1 for documented findings, never 2, with
four improved, two held, zero regressed, and one not-comparable ratchet.

Commit: `d982f63`.

The clean full audit from commit `d982f63` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 47 adds one private complete shell unzip composition containing only
`filesystem.Dependencies` and selects `filesystem.System()` only in the public
production `Unzip` entry. The public signature and caller graph remain
unchanged. The established body stays in its original order behind the private
composition, and only its two direct `os.MkdirAll` calls now use the existing
`filesystem.MkdirAll` forwarding helper with exact `os.ModePerm`. Directory
entries still append their exact joined path before one creation attempt and
continue directly after success. File entries still append their exact joined
path before one parent-directory attempt on `filepath.Dir(fpath)`, and a
failure still returns before direct file open. Archive open and close, zip-slip
validation, debug logging, entry order, file open flags and mode, entry open,
copy, both closes, partial results, and every exact error and precedence remain
unchanged. No filesystem interface or adapter implementation changed.

Six focused top-level shell contracts bring the suite to 330 tests. They prove
complete production filesystem selection, exact directory-entry and file-
parent paths, exact `os.ModePerm`, one attempt for every reached legacy site,
joined filename append and archive entry order, exact error identity and
partial filenames on both failure paths, safe zero dependencies returning
exact `filesystem.ErrNoFilesystem` before any later mutation, rejection of an
empty recording population, and absence of unrelated adapter operations. The
single guarded archive-fixture write is confined to `t.TempDir()` and no test
changes process, network, working-directory, or repository state.

Q0.6 stays closed with zero skipped tests and zero unsafe direct test writes;
its guarded safe-writer population improves from 24 to 25 sites, 20 write and
5 copy. Q1.1 stays 6 of 25, Q1.2 stays zero, Q1.3 improves from 15 of 36 to 13
of 34 as the two direct unzip directory-creation caller sites leave and the
existing adapter implementation remains singular, Q1.4 stays 7 of 8, and
exact Q2.1 stays 0 of 8.

The move-47 implementation gate passed focused shell/filesystem/Spring and
relevant process, command, context, config, HTTP, Maven, structurizr, profile,
browser, tips, file, template, Bitbucket, Kibana, local-config, and caller
package tests; API/CLI and subprocess compatibility; all four host flows and
their meta-contracts; full preflight, test, install, the standalone 62-control
launcher contract, uncached and race tests, vet, the pinned linter with zero
issues, the 15-control audit meta-suite, and empty-HOME count-2. Generated
compatibility and audit reports and all Go and linter caches remained outside
the measured tree. One standalone launcher attempt, the first complete
preflight attempt, and the first two complete `make test` attempts hit the
documented nested signal-interruption partial-raw-log timing flake; complete
reruns passed all 62 controls. The focused seven-criterion audit exited 1 for
documented findings, never 2, with four improved, two held, zero regressed,
and one not-comparable ratchet.

Commit: `b9fe209`.

The clean full audit from commit `b9fe209` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 48 reuses the private complete shell unzip composition and changes only
the file-entry branch's direct `os.OpenFile` call to the existing
`filesystem.OpenFile` forwarding helper. The exact joined path,
`os.O_WRONLY|os.O_CREATE|os.O_TRUNC` flags, entry mode, public signature,
production `filesystem.System()` selection, archive open and traversal order,
filename append order, both completed directory-creation boundaries,
directory short-circuit, parent-directory selection, partial results, and
error identity and precedence remain unchanged. Direct entry open, copy,
file close, entry close, and every other unzip effect also remain unchanged.
No filesystem interface, adapter implementation, caller, or inventory entry
changed.

Two focused top-level shell contracts bring the suite to 332 tests. The
focused recording double now records only the already-established MkdirAll
operation and the existing OpenFile operation. The contracts prove the exact
joined file path, flags, entry mode, parent-create-before-open order, filename
append and archive entry order, exact injected open-error identity and partial
results, no later file-entry operation after failure, directory-entry no-open
short-circuit, rejection of an empty OpenFile recording population, and
absence of unrelated adapter operations. The guarded archive fixture remains
under `t.TempDir()` and invokes no process, network request, working-directory
change, or repository write.

Q0.6 holds at 25 guarded safe-writer sites, 20 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2
stays zero, Q1.3 improves from 13 of 34 to 12 of 33 as the one direct unzip
file-open caller site leaves and the existing adapter implementation remains
singular, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

The move-48 implementation gate passed focused shell/filesystem/Spring and
relevant process, command, context, config, HTTP, Maven, structurizr, profile,
browser, tips, file, template, Bitbucket, Kibana, local-config, and caller
package tests; API/CLI and subprocess compatibility; all four host flows and
their meta-contracts; full preflight, test, install, the standalone 62-control
launcher contract, uncached and race tests, vet, the pinned linter with zero
issues, the 15-control audit meta-suite, and empty-HOME count-2. Generated
compatibility and audit reports and all Go and linter caches remained outside
the measured tree. The first complete preflight attempt hit the documented
nested signal-interruption partial-raw-log timing flake; its complete rerun
passed all 62 controls. The focused seven-criterion audit exited 1 for
documented findings, never 2, with four improved, two held, zero regressed,
and one not-comparable ratchet.

Commit: `2a684a0`.

The clean full audit from commit `2a684a0` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 49 changes only public `shell.Run`'s direct `exec.Command` construction
and synchronous `cmd.Run` execution flow to a private complete composition
containing `process.Dependencies`, with `process.System()` selected only in
the public production entry. The private flow builds the exact existing
`process.Command`, keeps command-name and variadic-argument order, wires
separate stdout and stderr buffers after the established debug log, and makes
one synchronous `process.Execute` attempt with empty directory and stdin and
`Start` false. The public signature, `Output` type and methods, returned
stdout/stderr bytes, callers, and the legacy behavior that discards the
process error and leaves `Output.Err` nil remain unchanged. The obsolete
uncalled private `run(*exec.Cmd)` helper is removed. No process interface,
adapter implementation, caller, other shell operation, or inventory entry
changed.

Four focused private shell contracts bring the suite to 336 tests. Their
recording process double records only the complete existing process request,
copies the ordered argument slice, and writes guarded in-memory stdout and
stderr bytes. The contracts prove complete production process selection,
exact command-name and ordered arbitrary argument bytes, one request, the
synchronous defaults, distinct writer identity and delivered arbitrary
bytes, exact injected dependency-error observation through the legacy nil
`Output.Err` contract, safe-zero behavior, rejection of an empty recording
population, and absence of unrelated process requests. They launch no
process, touch no network, change no working directory, and write no file.

Q0.6 holds at 25 guarded safe-writer sites, 20 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2
stays zero, Q1.3 improves from 12 of 33 to 11 of 32 as the one direct public
Run construction site leaves and the existing adapter implementation remains
singular, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8.

The move-49 implementation gate passed focused shell/process/Maven and
relevant command, context, config, filesystem, HTTP, Spring, structurizr,
profile, browser, tips, file, template, Bitbucket, Kibana, local-config, and
caller package tests; API/CLI and subprocess compatibility; all four host
flows and their meta-contracts; full preflight, test, install, the standalone
62-control launcher contract, uncached and race tests, vet, the pinned linter
with zero issues, the 15-control audit meta-suite, and empty-HOME count-2.
Generated compatibility and audit reports and all Go and linter caches
remained outside the measured tree. One standalone launcher attempt and the
first complete `make test` attempt hit the documented nested
signal-interruption partial-raw-log timing flake; complete reruns passed all
62 controls. The focused seven-criterion audit exited 1 for documented
findings, never 2, with four improved, two held, zero regressed, and one
not-comparable ratchet.

Commit: `da7eebf`.

The clean full audit from commit `da7eebf` exited 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passed.

Move 50 changes only the file-entry payload-transfer expression in private
shell unzip flow from direct `io.Copy(outFile, rc)` to the existing
`filesystem.Copy(dependencies.Files, outFile, rc)` forwarding helper. The
public signature, complete private dependency composition, production
`filesystem.System()` selection, archive open and close, traversal, zip-slip
check, filename append order, directory short-circuit, parent creation,
output-file open path, flags, and mode, entry open, one transfer attempt per
reached file entry, ignored transfer count and error, output-file then entry
close attempts, later traversal, partial results, and every earlier error and
close error identity and precedence remain unchanged. No filesystem interface,
adapter implementation, dependency selection, caller, other unzip branch, or
inventory entry changed.

Three focused top-level shell contracts bring the suite to 339 tests. The
private unzip filesystem double now records only its existing MkdirAll and
OpenFile requests plus the established Copy request, retaining the exact
opened destination identity and reading delivered source bytes only into
memory. The contracts prove exact ordered arbitrary archive-entry bytes, one
copy request for every reached file entry, parent-create, file-open, entry-open,
then copy order, arbitrary ignored copy counts and errors, continued output
and entry closes and later entry traversal, no destination write by the
double, no copy for directories or earlier failures, rejection of an empty
copy population, complete dependency preservation, and absence of unrelated
filesystem requests. The guarded archive fixture stays below `t.TempDir()`;
the contracts launch no process, touch no network, and change no working
directory.

Q0.6 holds at 25 guarded safe-writer sites, 20 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2
stays zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8. The focused audit
measures Q1.3 at 11 of 32 rather than the expected 10 of 31: the direct
`io.Copy` identity leaves, but the import-aware scanner continues to classify
the adapter request at the same caller line because its opened destination and
archive reader arguments retain filesystem-effect provenance. Changing that
scanner, the adapter contract, archive-open injection, or either close is
outside this move, so the measured ratchet holds without regression.

The move-50 implementation gate passes focused shell/filesystem/Spring and
relevant process, command, context, config, HTTP, Maven, structurizr, profile,
browser, tips, file, template, Bitbucket, Wpost, local-config, Kibana, and
caller package tests; API/CLI and subprocess compatibility; all four host
flows and their meta-contracts; full preflight, test, install, the standalone
62-control launcher contract, uncached and race tests, vet, the pinned linter
with zero issues, the 15-control audit meta-suite, and empty-HOME count-2.
Generated audit reports and Go and linter caches stay outside the measured
tree; generated compatibility reports are removed before the clean checkpoint
audit. The valid focused seven-criterion audit exits 1 for the documented
Q1.3 finding, with four improved, two held, zero regressed, and one
not-comparable ratchet.

Commit: `89918bd`.

The clean full audit from commit `89918bd` exits 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passes.

Move 51 changes only the initial archive-open request in private shell unzip
flow from direct `zip.OpenReader(src)` to
`filesystem.OpenZipReader(dependencies.Files, src)`. The existing filesystem
interface gains only `OpenZipReader(string) (*zip.ReadCloser, error)`, its
zero-safe forwarding helper, and the exact `zip.OpenReader(path)` system
implementation; every complete filesystem double receives the mechanical
method addition. The public `Unzip` signature, complete private dependency
composition, production `filesystem.System()` selection, exact source path,
returned archive identity, early archive-open return, deferred archive close
expression and placement, traversal, zip-slip check, filename append order,
directory short-circuit, parent creation, output-file open path, flags, and
mode, entry open, one payload attempt per reached file, output-file and entry
close attempts, later traversal, partial results, and all established error
and close precedence remain unchanged. No other filesystem operation, caller,
inventory entry, or production effect changes.

Seven focused top-level contracts bring the suite to 346 tests. The adapter
contracts prove zero-safe behavior, exact path forwarding, exact returned
reader identity and dependency error, rejection of an empty archive-open
population, and system entry order, names, and bytes. The private unzip double
records the one archive-open request before its already-established MkdirAll,
OpenFile, and Copy requests. Its contracts use an injected reader with a
deliberately nonexistent supplied source path to prove exact source delivery,
one request, injected-reader traversal and deferred close, ordered filenames,
the exact archive-open error with nil partial filenames, complete dependency
preservation, zero-safe failure before a developer path is opened, rejection
of an empty recording population, and absence of unrelated filesystem
requests. Both guarded archive fixtures stay below `t.TempDir()`; the focused
contracts launch no process, touch no network, and change no working
directory.

Q0.6 holds with zero skipped tests and zero unsafe direct test writes; the
guarded safe-writer population is 26 sites, 21 write and 5 copy, after the
focused system archive fixture. Q1.1 stays 6 of 25, Q1.2 stays zero, Q1.4
stays 7 of 8, and exact Q2.1 stays 0 of 8. Q1.3 improves from 11 of 32 to 10
of 32 because the direct public-package archive open becomes one exact
filesystem-adapter implementation site. The scanner and inventory remain
unchanged.

The move-51 implementation gate passes focused shell/filesystem/Spring and
relevant process, command, context, config, HTTP, Maven, structurizr, profile,
browser, tips, file, template, Bitbucket, Wpost, local-config, Kibana, and
caller package tests; API/CLI and subprocess compatibility; all four host
flows and their meta-contracts; full preflight, test, install, the standalone
62-control launcher contract, uncached and race tests, vet, the pinned linter
with zero issues, the 15-control audit meta-suite, and empty-HOME count-2.
Generated audit reports, compatibility reports, Go caches, and linter caches
stay outside the measured tree. The valid focused seven-criterion audit exits
1 for the documented Q1.3 finding, never 2, with four improved, two held, zero
regressed, and one not-comparable ratchet.

Commit: `2f0a072`.

The clean full audit from commit `2f0a072` exits 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passes.

Move 52 changes only the output-file close request in private shell unzip flow
from direct `err = outFile.Close()` to
`err = filesystem.Close(dependencies.Files, outFile)`. The existing filesystem
interface gains only `Close(filesystem.File) error`, its zero-safe forwarding
helper, and the exact `file.Close()` system implementation; every complete
filesystem double receives the mechanical method addition. The public
`Unzip(string, string) ([]string, error)` signature, complete private
dependency composition, production `filesystem.System()` selection, exact
opened output identity, archive open and deferred close, traversal, zip-slip
check, filename append order, directory short-circuit, parent creation,
output-file open path, flags, and mode, entry open, one payload attempt per
reached file, the output-close assignment and following error branch, entry
close, later traversal, partial results, and all established error and close
precedence remain unchanged. No archive-open, archive-close, entry-reader
close, Copy, OpenFile, MkdirAll, entry-open, caller, other unzip branch,
inventory entry, or other production effect changes.

Six focused top-level contracts bring the suite to 352 tests. The adapter
contracts prove zero-safe behavior without closing a developer-supplied file,
exact file identity and one request, exact dependency error, exact system
close invocation and error, and rejection of an empty close population. The
private unzip double records the exact opened output identity after its
already-established Copy request. Its contracts prove one output-close request
per reached file, Copy-before-close order, exact close-error and partial
filenames, successful entry close and later traversal, error short-circuit
before entry close and later entries, complete dependency preservation,
rejection of an empty output-close population, and absence of unrelated
filesystem requests. The guarded archive fixture stays below `t.TempDir()`;
the focused contracts launch no process, touch no network, and change no
working directory.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2
stays zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8. The focused audit
measures Q1.3 at 10 of 32: the direct output-file `Close` identity leaves, but
the import-aware scanner continues to classify the adapter request at the same
caller line because its opened-file argument retains filesystem-effect
provenance. The exact system implementation is inside the declared filesystem
adapter, so the denominator and violation count both hold. The scanner and
inventory remain unchanged.

The move-52 implementation gate passes focused shell/filesystem/Spring and
relevant process, command, context, config, HTTP, Maven, structurizr, profile,
browser, tips, file, template, Bitbucket, Wpost, local-config, Kibana, and
caller package tests; API/CLI and subprocess compatibility; all four host
flows and their meta-contracts; full preflight, test, install, the standalone
62-control launcher contract, uncached and race tests, vet, the pinned linter
with zero issues, the 15-control audit meta-suite, and empty-HOME count-2.
Generated audit reports, compatibility reports, Go caches, and linter caches
stay outside the measured tree. The valid focused seven-criterion audit exits
1 for the documented Q1.3 finding, never 2, with four improved, two held, zero
regressed, and one not-comparable ratchet.

Commit: `0128cd0`.

The clean full audit from commit `0128cd0` exits 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passes.

Move 53 changes only the entry-reader close request in private shell unzip flow
from direct `err = rc.Close()` to
`err = filesystem.CloseReader(dependencies.Files, rc)`. The existing filesystem
interface gains only `CloseReader(io.ReadCloser) error`, its zero-safe forwarding
helper, and the exact `reader.Close()` system implementation; every complete
filesystem double receives the mechanical method addition. The public
`Unzip(string, string) ([]string, error)` signature, complete private dependency
composition, production `filesystem.System()` selection, exact opened entry-reader
identity, archive open and deferred close, traversal, zip-slip check, filename
append order, directory short-circuit, parent creation, output-file open path,
flags, and mode, entry open, Copy request and ignored results, output-file Close,
the entry-close assignment and following error branch, later traversal, partial
results, and all established error and close precedence remain unchanged. No
archive-open, archive-close, output-file Close, Copy, OpenFile, MkdirAll,
entry-open, caller, other unzip branch, inventory entry, or other production
effect changes.

Six focused contracts bring the suite to 358 tests. The adapter contracts prove
zero-safe behavior without closing a developer-supplied reader, exact reader
identity and one dependency request, exact dependency error, exact system close
invocation and error, and rejection of an empty reader-close population. The
private unzip double records the exact opened reader already delivered as the
Copy source. Its contracts prove one entry-reader-close request per reached file
after Copy and successful output-file Close, exact request order and identity,
exact close-error and partial filenames, suppression of later traversal after an
entry-close error, suppression of entry Close after an output-close error,
successful reader Close and later traversal, complete dependency preservation,
rejection of an empty entry-close population, and absence of unrelated filesystem
requests. The guarded archive fixture stays below `t.TempDir()`; the focused
contracts launch no process, touch no network, and change no working directory.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2 stays
zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8. The focused audit measures
Q1.3 at 10 of 32: the direct interface-typed entry-reader `Close` was absent
from the prior violation set, and the import-aware scanner also omits the new
CloseReader adapter request. The exact system implementation is inside the
declared filesystem adapter, so the denominator and violation count both hold.
The scanner and inventory remain unchanged.

The move-53 implementation gate passes focused shell/filesystem/Spring and
relevant process, command, context, config, HTTP, Maven, structurizr, profile,
browser, tips, file, template, Bitbucket, Wpost, local-config, Kibana, and caller
package tests; API/CLI and subprocess compatibility; all four host flows and
their meta-contracts; full preflight, test, install, the standalone 62-control
launcher contract, uncached and race tests, vet, the pinned linter with zero
issues, the 15-control audit meta-suite, and empty-HOME count-2. Generated audit
and compatibility reports and all Go and linter caches stay outside the
measured tree. The valid focused seven-criterion audit exits 1 for the documented
Q1.3 finding, never 2, with four improved, two held, zero regressed, and one
not-comparable ratchet.

Commit: `70bee0e`.

The clean full audit from commit `70bee0e` exits 1 for 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
not-comparable ratchet, and zero dirty paths. Empty-HOME count-2 also passes.

Move 54 changes only the deferred concrete archive-reader close request in
private shell unzip flow from direct `_ = r.Close()` to
`_ = filesystem.CloseReader(dependencies.Files, r)` inside the same anonymous
defer immediately after successful archive open. The existing filesystem
`CloseReader` parameter is widened from `io.ReadCloser` to `io.Closer` in the
interface, zero-safe helper, exact system implementation, and every complete
filesystem double so it accepts the exact `*zip.ReadCloser` returned by
`OpenZipReader`. The public `Unzip(string, string) ([]string, error)` signature,
complete private dependency composition, production `filesystem.System()`
selection, exact archive identity, defer placement, archive and entry
traversal, zip-slip check, filenames, entry order and bytes, directory and file
operations, entry-reader and output-file closes, partial results, and every
established error and close precedence remain unchanged. Archive close is
attempted exactly once after each successful archive open and never after an
open failure; its exact result remains ignored without changing named results
or an earlier error.

Two focused top-level contracts bring the suite to 360 tests, while the
existing complete private Unzip contracts now record the archive-close request.
They prove the exact opened archive identity, one deferred request after
successful traversal and after every reached entry-reader close, execution on
zip-slip and every established later error return, no request after archive-open
failure, ignored injected archive-close errors without result or precedence
changes, complete dependency preservation, distinction from entry-reader-close
requests, rejection of an empty archive-close population, and absence of
unrelated filesystem operations. The guarded archive fixture stays below
`t.TempDir()`; the focused contracts launch no process, touch no network, and
change no working directory.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2 stays
zero, Q1.4 stays 7 of 8, and exact Q2.1 stays 0 of 8. The focused audit measures
Q1.3 at 10 of 32: the prior concrete archive-reader Close was absent from its
reported violation set, and the import-aware scanner also omits the replacement
adapter request. The scanner and inventory remain unchanged. The selected
seven-criterion audit exits 1 for documented findings, never 2, with four
improved, two held, zero regressed, and one non-comparable ratchet.

Focused shell/filesystem/Spring and all relevant caller-package tests, API/CLI
and subprocess compatibility, the 62-control launcher contract, Make
meta-contracts, complete `make test`, install, uncached and race tests, vet,
all four host flows and their meta-contracts, the repaired 15-control audit
meta-suite, the clean full audit, and empty-HOME count-2 pass. The full audit
exits 1 for 15 documented findings, never 2, with L0 8 of 8, five improved,
two held, zero regressed, one non-comparable ratchet, and zero dirty paths.

Product commit: `3bd07e9`. Audit-apparatus commit: `ce736a2`.

The authoritative clean full audit from `ce736a2` measures the unchanged
`3bd07e9` product tree. Its structured scorecard SHA-256 is
`c7ff3bd2cf32f8615c0fb1329c980bafba9b5e988cf4e785f4b8dc608d3f4243`;
the measured commit tree is `20266d8c4ff215679c981f9e4bf5d005fb435656`.

The inherited failure was reproduced before repair: exact predecessor
`70bee0e` and implementation `3bd07e9` replicas passed T1-T14 and failed T15
identically after historical tests created ignored measured-tree paths. A
focused red contract then proved that T15 had to distinguish audit-test effects
from source present before execution without excluding or deleting either.

T15 now verifies pristine old/new replicas before execution, rejects a tracked
source-drift probe, and records every post-test path, mode, raw digest, and
checkout-root-independent digest in separate clean-HOME effect replicas. It
runs the historical raw audit in equivalent execution replicas and feeds old
and new structured parsers equivalent pristine replicas of exact `5635d50`
plus the authorized inventory overlay. This preserves dirty-path truth while
preventing test-created output from masquerading as pre-measurement source.

Orchestration alone exposed a second historical precondition: an empty HOME
made the untouched historical tests pass, while stored baseline debt records
their failure. The directly related reproduction recipe and migration metadata
therefore pin an apparatus-owned directory-shaped legacy active-profile fixture
with marker SHA-256
`cb95f24c35d3987f8aba51231aade19580ffe9242324804fcad2ccff350d1c9a`.
The parser, historical checkout, baseline scorecards and raw debt, instrument
identities, inventory, scanner, and product remain unchanged.

Complete repaired meta-suites pass all 15 controls over the unchanged
implementation-content tree and over exact predecessor `70bee0e` with only the
four apparatus paths overlaid. Old and new structured baseline hashes reproduce
as `d420887d73aabf496ff276fcc55d13ad379ac49c9322fad58808b5e28fdba7df`
and `5fb3226009cfbf0d29f63fa03592157cce4efcec6e38583b64a86f6288e89490`;
the normalized stored raw body reproduces as
`cf23c9dca987acd3a966693f933c4c7eca9f7d11c51f0fbabfe1d697f3d7497f`.
All 228 numeric debt leaves and Q3.9 match.

Move 55 changes only the ignored `structurizr-cli export` process request in
private plugin-diagrams flow. `structurizrCmd.RunE` retains the mandatory
workspace lookup and passes the complete private
`pluginDiagramsExportDependencies` value to `runStructurizrDiagrams`.
Production selects `process.System()` explicitly. After the unchanged exact
`.structurizr/` deletion and before the unchanged dot-file discovery, the flow
now calls `_ = process.Execute` with executable `structurizr-cli`, ordered
arguments `export`, `-w`, the caller-supplied workspace, `-format`, `dot`,
`-output`, `.structurizr/`, and zero values for directory, streams, and Start.
The one synchronous result remains ignored. Graphviz `dot`, macOS `open`,
output capture and writing, iteration, errors, paths, printed text, Cobra
registration, public API, adapter, scanner, inventory, and audit apparatus are
unchanged.

Four focused contracts bring the suite to 364 tests. They prove complete
production process selection, the exact command and ordered arbitrary workspace
bytes, one request, synchronous zero-value directory and streams, the exact
ignored injected process error, continued discovery, safe-zero behavior,
rejection of an empty recorded population, and absence of unrelated export,
`dot`, or `open` requests. The focused tests use temporary working directories,
launch no external program, touch no network, and write no repository fixture.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2 stays
zero, Q1.3 improves from 10 of 32 to 9 of 31, Q1.4 stays 7 of 8, exact Q2.1
stays 0 of 8, and Q3.4 stays zero phrases across 73 Markdown files. The scanner
and inventory remain unchanged.

The clean move-55 gate passes focused cmd/process/structurizr and relevant
caller tests; API/CLI and subprocess compatibility; build, complete and
uncached tests, race, vet, pinned lint, all 62 launcher controls, Make and all
production-script meta-contracts, all four host acceptance flows, the repaired
15-control audit meta-suite, and empty-HOME count-2. The first launcher attempts
with an explicitly nested temporary root hit the documented partial-raw-log
signal timing flake; the standalone default-path contract and complete Make run
passed all 62 controls. The clean full audit exits 1 for the same 15 documented
findings, never 2, with L0 8 of 8, five improved, two held, zero regressed, one
non-comparable ratchet, and zero dirty paths.

Product commit: `3972fc6`.

The authoritative clean full audit from `3972fc6` has structured scorecard
SHA-256
`ced6611a0990d55285007f9b9d19a990bbbe2bef9ee2fbe18084420fe9956fa6`;
the measured commit tree is `395bbad053237f8769be804e211c7cb0f49550cc`.

Move 56 changes only the ignored macOS `open` process request after each
successful Graphviz conversion in private plugin-diagrams flow.
`structurizrCmd.RunE` retains the mandatory workspace lookup and now passes the
complete private `pluginDiagramsOpenDependencies` value alongside the completed
export dependency. Production selects `process.System()` explicitly.

After the unchanged direct `dot` request succeeds and before the range advances,
the private `openStructurizrDiagram` helper calls `_ = process.Execute` once with
executable `open`, the exact single `outputPngFile` argument, and zero values for
directory, streams, and Start. The synchronous result remains ignored, and the
next discovered file is attempted unconditionally. Workspace lookup, deletion,
export, discovery, iteration, output path construction, printed text, direct
Graphviz execution and error, output capture and writing, Cobra behavior,
public API, `pkg/structurizr`, adapter, scanner, inventory, and audit apparatus
are unchanged.

Five focused private contracts bring the suite to 369 tests. They prove complete
production process selection, the exact command and arbitrary output-path bytes,
one request per helper invocation, synchronous zero-value directory and streams,
the exact directly ignored injected process error, safe-zero behavior, rejection
of an empty recorded population, placement only after the direct Graphviz error
gate, unconditional continued iteration, complete dependency delivery, and
absence of export, `dot`, or another process request. They launch no external
program, touch no network, and write no repository fixture.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2 stays
zero, Q1.3 improves from 9 of 31 to 8 of 30, Q1.4 stays 7 of 8, exact Q2.1
stays 0 of 8, and Q3.4 stays zero phrases across 74 Markdown files. The scanner
and inventory remain unchanged.

The clean move-56 gate passes focused cmd/process/structurizr and relevant
caller tests; API/CLI and subprocess compatibility; build, complete and
uncached tests, race, vet, pinned lint, all 62 launcher controls, Make and all
production-script meta-contracts, all four host acceptance flows, the repaired
15-control audit meta-suite, and empty-HOME count-2. The clean full audit exits
1 for the same 15 documented findings, never 2, with L0 8 of 8, five improved,
two held, zero regressed, one non-comparable ratchet, and zero dirty paths.

Product commit: `6f4bc37`.

The authoritative clean full audit from `6f4bc37` has structured scorecard
SHA-256
`efaaee797a34d182d03d6830d0af4c342b7e744180775942119a77c7e0f88090`;
the measured commit tree is `cf7c1af5b9ecbbbb7b218cfa6951354d6bf94854`.

Move 57 changes only the direct Graphviz `dot` conversion request in private
plugin-diagrams flow. `structurizrCmd.RunE` retains the mandatory workspace
lookup and now passes the complete private
`pluginDiagramsGraphvizDependencies` value alongside the completed export and
open dependencies. Production selects `process.System()` and
`filesystem.System()` explicitly.

For every discovered dot file, the private `convertStructurizrDiagram` helper
calls `process.Execute` once with executable `dot`, ordered arguments containing
the input file and `-Tpng`, empty directory, nil stdin, distinct stdout and
stderr buffers, and synchronous Start false. The exact process error returns
before any output write or open attempt. On success, the helper passes the
captured stdout bytes, exact output PNG path, and mode `0644` to
`filesystem.WriteFile` once, ignores the exact write result, discards captured
stderr, and leaves the following completed ignored open request in place.
Workspace lookup, deletion, export, discovery, iteration, output path
construction, printed text, open behavior, Cobra behavior, public API,
`pkg/structurizr`, both adapters, scanner, inventory, and audit apparatus are
unchanged.

Six focused private contracts bring the suite to 375 tests. They prove complete
production dependency selection, exact process command and stream identities,
one process request, exact process-error precedence, suppression of write and
open after process failure, exact successful output bytes/path/mode, one write
attempt, exact ignored write error, following open placement, continued
iteration, safe-zero behavior, rejection of empty process and write
populations, complete dependency delivery, and absence of export, open, or
another unrelated process/write request. They launch no external program,
touch no network, and write no repository fixture.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2 stays
zero, Q1.3 improves from 8 of 30 to 7 of 29, Q1.4 stays 7 of 8, exact Q2.1
stays 0 of 8, and Q3.4 stays zero phrases across 75 Markdown files. The scanner
and inventory remain unchanged.

The clean move-57 gate passes focused cmd/process/filesystem/structurizr and
relevant caller tests; API/CLI and subprocess compatibility; build, complete
and uncached tests, race, vet, pinned lint, all 62 launcher controls, Make and
all production-script meta-contracts, all four host acceptance flows, the
repaired 15-control audit meta-suite, and empty-HOME count-2. The clean full
audit exits 1 for the same 15 documented findings, never 2, with L0 8 of 8,
five improved, two held, zero regressed, one non-comparable ratchet, and zero
dirty paths.

Product commit: `9164e06`.

The authoritative clean full audit from `9164e06` has structured scorecard
SHA-256
`9f584097422a17fd6eea005dda97182e042aad05a34837f0461ce0933248f700`;
the measured commit tree is `78fef4303e29d2b8d79434727bb8f7d9e0fd070a`.

Move 58 changes only Maven command's conditional standard-output capability
selection. `process.Dependencies` now carries the runner and standard-output
writer together, `process.System()` selects its existing runner and exact
`os.Stdout` identity together, and the narrow `process.Stdout` helper returns
that injected writer only when enabled. `pkg/maven/command.go` retains its
global logrus level predicate but forwards the complete process dependency to
the helper instead of selecting output through `logger.StdOut()`.

Every Maven command retains the exact executable and ordered argument bytes,
`project.Path` directory, nil stdin and stderr, synchronous Start false, one
process attempt, exact dependency error, and unchanged log text before process
execution. Debug and trace retain exact production `os.Stdout`; panic, fatal,
error, warn, and info retain nil stdout. The zero dependency remains a safe
no-op with nil output. `maven.RunOn`, its returned callback, `Repository`,
`pkg/logger`, plugin diagrams, shell unzip, all other callers, public API/CLI,
scanner, inventory, audit apparatus, and every completed effect are unchanged.

Five new contracts bring the suite to 380 tests. They prove exact system
`os.Stdout`, enabled injected-writer identity, disabled and zero-dependency nil
output, every global logrus level, the complete Maven command, one request,
exact error, complete dependency preservation, exact log ordering, rejection
of an empty population, and absence of another process request. Relevant
process doubles now start from the complete system dependency and replace only
the runner. Recording boundaries launch no external program, touch no network,
and write no repository fixture.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2 stays
zero, regenerated exact Q1.3 is 15 of 37, Q1.4 stays 7 of 8, exact Q2.1 stays 0
of 8, and Q3.4 stays zero phrases across 76 Markdown files. The caller-side
`logger.StdOut()` site moved behind the process dependency; the unchanged
import-aware scanner follows the new `os.Stdout` capability through existing
`process.System()` selections, accounting for the changed exact population.
Clock and server remain absent, so Q1.3 remains the single non-comparable
ratchet. The scanner and inventory are unchanged.

The clean move-58 gate passes focused process/Maven/logger and relevant caller
tests; API/CLI and subprocess compatibility; build, complete and uncached
tests, race, vet, pinned lint, all 62 launcher controls, Make and all
production-script meta-contracts, all four host acceptance flows, the repaired
15-control audit meta-suite, and empty-HOME count-2. The clean full audit exits
1 for the same 15 documented findings, never 2, with L0 8 of 8, five improved,
two held, zero regressed, one non-comparable ratchet, and zero dirty paths.

Product commit: `4376e05`.

The authoritative clean full audit from `4376e05` has structured scorecard
SHA-256
`09cb1f364cbf55243f557d77c400b4befa7672d5b0a9378ad9e7d7cd6447ce32`;
the measured commit tree is `550e10792b81c5b1fa3bc4714a451a1afb586541`.

Move 59 changes only shell Unzip's direct archive-entry reader selection in
private `unzipWithDependencies`. The existing complete filesystem dependency
now exposes `OpenZipEntry(*zip.File) (io.ReadCloser, error)`, whose zero-safe
helper returns the established `ErrNoFilesystem` and whose system
implementation invokes `Open` on the exact supplied entry. Production replaces
only `rc, err := f.Open()` with the helper call on the existing injected
filesystem dependency.

The exported `Unzip` signature, archive source and destination, one archive
open, exact deferred archive-close identity and ignored result, entry order,
log text and placement, zip-slip check, partial-filename append point,
directory handling, parent creation, destination flags and entry mode remain
unchanged. Each file entry still opens its destination before one entry-open
attempt, copies once with the exact returned reader while ignoring count and
error, closes output before entry reader, returns the same close errors, and
suppresses every later operation after the same failure points. A zero
dependency performs no archive or entry open or filesystem mutation and
returns the exact established error.

Four focused top-level contracts bring the suite to 384 tests, while the
complete shell recorder now captures entry-open identity, reader, error, and
operation placement. They prove the exact `*zip.File`, exact returned reader
and error, safe zero behavior, the complete successful file-entry sequence,
entry-open failure after destination open, suppression of copy, closes, and
later traversal on that failure, complete dependency preservation, rejection
of an empty entry-open population, and absence of an unrelated filesystem
request. Every complete filesystem double and caller mechanically preserves
the extended dependency. Focused tests use only the existing guarded temporary
archive helper, launch no external process, touch no network, and write no
repository fixture.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2 stays
zero, regenerated exact Q1.3 stays 15 of 37 because the unchanged scanner does
not catalog `archive/zip.File.Open`, Q1.4 stays 7 of 8, exact Q2.1 stays 0 of 8,
and Q3.4 stays zero phrases across 76 Markdown files. The scanner, inventory,
baseline, and audit apparatus are unchanged.

The clean move-59 gate passes focused filesystem/shell and all relevant caller
tests; API/CLI and subprocess compatibility; complete tests, race, vet, pinned
lint, all 62 launcher controls, Make and production-script meta-contracts, all
four host acceptance flows, the repaired 15-control audit meta-suite, and
empty-HOME count-2. The first complete preflight attempt reached only the
sandbox-blocked default golangci-lint cache; its isolated-cache rerun passed.
The clean full audit exits 1 for the same 15 documented findings, never 2, with
L0 8 of 8, five improved, two held, zero regressed, one non-comparable ratchet,
and zero dirty paths.

Product commit: `9f56714`.

The authoritative clean full audit from `9f56714` has structured scorecard
SHA-256
`41b51871f6436851613af04c86af93b6fb7f7f0e7f74833c47d8e3d45ddfb10e`;
the measured commit tree is `18f303931958d94ba72b9ed861faa24b94e29925`.

Move 60 changes only public shell `Run`'s production process-capability
selection. The existing process adapter now exports `SystemRunner()`, which
returns the exact existing `systemRunner` with a nil standard-output field.
Private `systemRunDependencies()` selects that runner-only dependency instead
of `process.System()`. The complete `process.System()` implementation remains
the same runner plus exact `os.Stdout`, preserving Maven's completed debug
standard-output behavior and every other caller.

The exported `Run(name string, args ...string) Output` signature, command name
and ordered argument bytes, debug log text and placement, distinct stdout and
stderr buffer identities, nil stdin and working directory, synchronous
execution, one exact runner request, returned output bytes, legacy nil
`Output.Err` after a dependency error, safe zero dependency, and every return
path remain unchanged. No process execution, retry, fallback, wrapping,
logging, cleanup, or global-state behavior changed.

Two focused top-level contracts bring the suite to 386 tests. The process
contract proves the runner-only selector carries the exact system runner and
nil stdout while `process.System()` still carries exact `os.Stdout`. The shell
contracts prove the complete arbitrary command and stdout/stderr identities,
exact dependency error behavior, complete caller-owned dependency preservation,
exact log-before-process ordering, one non-empty request population, absence of
a second request, and the safe zero behavior. The recording boundaries launch
no external program, touch no network, and write no repository fixture.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2 stays
zero, regenerated exact Q1.3 improves from 15 of 37 to 14 of 36, Q1.4 stays 7
of 8, exact Q2.1 stays 0 of 8, and Q3.4 stays zero phrases across 78 Markdown
files. The unchanged import-aware scanner no longer follows `os.Stdout`
provenance through this one runner-only dependency. Clock and server remain
absent, so Q1.3 remains the single non-comparable ratchet. The scanner,
inventory, baseline, and audit apparatus are unchanged.

The clean move-60 gate passes focused process/shell and all relevant process
caller tests; API/CLI and fresh subprocess compatibility; build, complete and
uncached tests, race, vet, pinned lint, all 62 launcher controls, Make and all
production-script meta-contracts, all four host acceptance flows, the repaired
15-control audit meta-suite, and empty-HOME count-2. Repeated complete Make
attempts hit the previously observed partial-raw-log signal-fixture flake; an
immediate standalone launcher run and the final complete `make test` rerun
passed all 62 controls. The clean full audit exits 1 for the same 15 documented
findings, never 2, with L0 8 of 8, five improved, two held, zero regressed, one
non-comparable ratchet, and zero dirty paths.

Product commit: `8503601`.

The authoritative clean full audit from `8503601` has structured scorecard
SHA-256
`09ef869b92d3d77a6deb0e5e4f1cd9bdeec1d72596ace95648d6165cf10797e1`;
the measured commit tree is `480c44f20723b48c74243ab4b8c5618c8ddf81fd`.

Move 61 changes only the production process-capability selections used by the
four public shell Git execution wrappers. Private
`systemGitDependencies()` returns the exact existing `process.SystemRunner()`;
`GitClone`, `GitPull`, `GitInit`, and `GitAddAndCommit` select that one complete
runner-only dependency instead of `process.System()`. `runGit`, `GitDirty`,
`GitIsRepo`, public shell Run, Maven, every other caller, and both process
adapter selectors remain unchanged.

The four exported signatures, executable `git`, exact ordered arguments,
debug log text and placement, distinct stdout and stderr buffers and bytes,
nil stdin and directory, synchronous execution, exact one-request and
two-request populations, existing add/commit request selection, direct runner
error behavior including nil returned `Output.Err`, safe zero dependency, and
all return paths remain unchanged. There is no fallback, retry, wrapping,
additional logging, execution, cleanup, or global-state change.

Two focused top-level contracts bring the suite to 388 tests. They prove exact
system-runner identity with nil dependency stdout, complete arbitrary commands
and buffer identities, exact request sequences and log/process placement,
legacy error behavior, safe zero behavior for every wrapper, complete
caller-owned dependency preservation, non-empty populations, and the absence
of another process request. The recording boundary launches no external
program, touches no network, and writes no repository fixture.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2 stays
zero, regenerated exact Q1.3 improves from 14 of 36 to 10 of 32, Q1.4 stays 7
of 8, exact Q2.1 stays 0 of 8, and Q3.4 stays zero phrases across 79 Markdown
files. The unchanged import-aware scanner no longer follows `os.Stdout`
provenance through those four runner-only production dependency selections.
Clock and server remain absent, so Q1.3 remains the single non-comparable
ratchet. The scanner, inventory, baseline, and audit apparatus are unchanged.

The clean move-61 gate passes focused process/shell and all relevant process
caller tests; API/CLI and fresh subprocess compatibility; build, complete and
uncached tests, race, vet, pinned lint, all 62 launcher controls, Make and all
production-script meta-contracts, all four host acceptance flows, the repaired
15-control audit meta-suite, and empty-HOME count-2. The clean full audit exits
1 for the same 15 documented findings, never 2, with L0 8 of 8, five improved,
two held, zero regressed, one non-comparable ratchet, and zero dirty paths.

Product commit: `5b678ab`.

The authoritative clean full audit from `5b678ab` has structured scorecard
SHA-256
`0632a71b59dc2793c2cf308d0514569b698cdead3f02fef83a87cbe029fccfff`;
the measured commit tree is `010970642433199467db2fa3b6eaaf7c9ab7bece`.

Move 62 changes only the private profile-editor production process-capability
selection. `systemProfileEditorDependencies()` now returns the exact existing
`process.SystemRunner()` instead of `process.System()`, so its complete process
dependency has a nil `Stdout`. Its separate command-stream fields remain exact
`os.Stdin` and `os.Stdout`. `runProfileEditor`, profile commands and
configuration, Maven, shell Run and Git, plugin diagrams, browser, Unzip, every
other caller, and both process adapter selectors remain unchanged.

The arbitrary editor executable, single exact config-path argument, empty
directory, dependency-selected stdin and stdout, nil stderr, synchronous
execution, exact one-request population, direct runner error, safe zero
dependency, `EDITOR` selection, `vim` default, Cobra registration and flags,
later profile operations and errors, and every public path remain unchanged.
There is no fallback, retry, wrapping, additional logging, execution, cleanup,
or global-state change.

Four strengthened focused top-level contracts keep the suite at 388 tests.
They prove exact system-runner identity with nil process-dependency stdout,
separate exact production terminal streams, complete arbitrary editor and
config-path requests and injected stream identities, direct errors, safe zero
behavior, complete caller-owned dependency preservation, a non-empty
one-request population, and absence of another request. The recording boundary
launches no external program, touches no network, and writes no repository
fixture.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2 stays
zero, regenerated exact Q1.3 improves from 10 of 32 to 9 of 31, Q1.4 stays 7 of
8, exact Q2.1 stays 0 of 8, and Q3.4 stays zero phrases across 80 Markdown
files. The unchanged import-aware scanner no longer follows `os.Stdout`
provenance through the private profile-editor runner-only production dependency.
Clock and server remain absent, so Q1.3 remains the single non-comparable
ratchet. The scanner, inventory, baseline, and audit apparatus are unchanged.

The clean move-62 gate passes focused profile/process and all relevant process
caller tests; API/CLI and fresh subprocess compatibility; build, complete and
uncached tests, race, vet, pinned lint, all 62 launcher controls, Make and all
production-script meta-contracts, all four host acceptance flows, the repaired
15-control audit meta-suite, and empty-HOME count-2. The first complete Make
attempt hit the established nested partial-raw-log signal-fixture flake; its
immediate standalone launcher run and the final complete `make test` rerun
passed all 62 controls. The clean full audit exits 1 for the same 15 documented
findings, never 2, with L0 8 of 8, five improved, two held, zero regressed, one
non-comparable ratchet, and zero dirty paths.

Product commit: `b0d324a`.

The authoritative clean full audit from `b0d324a` has structured scorecard
SHA-256
`e9d5e45b477cde1791a0e73b0d5cf1332036fec371ccf29f4d9b0827a838ab46`;
the measured commit tree is `86b97ce8e9cbdbc2259a21a28e87c3678552c06d`.

Move 63 changes only the three private plugin-diagrams production
process-capability selections. `systemPluginDiagramsExportDependencies()`,
`systemPluginDiagramsGraphvizDependencies()`, and
`systemPluginDiagramsOpenDependencies()` now carry the exact existing
`process.SystemRunner()` instead of `process.System()`, so their complete
process dependencies have nil `Stdout`. The Graphviz selector's separate
complete exact `filesystem.System()` dependency remains unchanged. Every
execution helper, command and filesystem field, return path, and other caller
remains unchanged.

The ignored synchronous `structurizr-cli export` request and error, workspace
and temporary-directory arguments, deletion-before-export and
discovery-after-export placement, and discovery result and error behavior are
unchanged. The synchronous `dot` request, ordered input and format arguments,
distinct stdout and stderr buffers and bytes, direct process error,
no-write-after-process-error behavior, exact `0644` stdout write, ignored write
error, and process-then-write order are unchanged. The ignored synchronous
macOS `open` request and error, output path, placement after each successful
conversion, and continued iteration are unchanged. Safe zero dependencies,
non-empty populations, exact one-request attempts, absence of another request
or write, Cobra registration, flags, and every public behavior remain intact.
There is no fallback, retry, wrapping, logging, execution, cleanup, or
global-state change.

Four strengthened focused top-level contracts keep the suite at 388 tests.
They prove exact system-runner identity with nil process-dependency stdout for
all three selectors, the Graphviz selector's unchanged exact system filesystem
dependency, complete requests, arguments, stream identities and bytes, writes,
modes, attempts and sequences, ignored and direct errors, safe zero behavior,
complete caller-owned dependency preservation, non-empty populations, and the
absence of another request or write. Three recording doubles preserve complete
caller-owned process and filesystem dependency values. They launch no external
program, touch no network, and write no repository fixture.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2 stays
zero, regenerated exact Q1.3 improves from 9 of 31 to 8 of 30, Q1.4 stays 7 of
8, exact Q2.1 stays 0 of 8, and Q3.4 stays zero phrases across 81 Markdown
files. The three selectors feed one Structurizr caller composition, so the
unchanged import-aware scanner removes one provenance site. Clock and server
remain absent, so Q1.3 remains the single non-comparable ratchet. The scanner,
inventory, baseline, and audit apparatus are unchanged.

The clean move-63 gate passes focused plugin-diagrams, process, filesystem, and
all relevant process caller tests; API/CLI and fresh subprocess compatibility;
build, complete and uncached tests, race, vet, pinned lint, all 62 launcher
controls, Make and all production-script meta-contracts, all four host
acceptance flows, the repaired 15-control audit meta-suite, and empty-HOME
count-2. An initial launcher invocation and one complete preflight invocation
hit the established nested partial-raw-log signal-fixture flake; their
immediate standalone or complete unchanged reruns passed all 62 controls. The
first complete preflight also required its documented writable lint cache under
the sandbox. The clean full audit exits 1 for the same 15 documented findings,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one
non-comparable ratchet, and zero dirty paths.

Product commit: `1bce06f`.

The authoritative clean full audit from `1bce06f` has structured scorecard
SHA-256
`5cf5c2c171749e42dee5dfea6f5986b09dce96d00510560f4e625c7a6e4e3643`;
the measured commit tree is `ebd29be397a8fbdda257ad94e3c8863fbaae84fa`.

Move 64 changes only the private browser-launcher production
process-capability selection. `systemBrowserLauncherDependencies()` now
carries the exact existing `process.SystemRunner()` instead of
`process.System()`, so its complete process dependency has a nil `Stdout`. Its
separate exact `runtime.GOOS` selection, `openBrowser`, exported `OpenBrowser`,
its caller, command fields, return path, and every other process caller remain
unchanged.

The exact asynchronous `xdg-open` Linux request, `rundll32` Windows request
with ordered `url.dll,FileProtocolHandler` and URL arguments, and `open` macOS
request are unchanged. Each arbitrary URL byte, empty directory, nil stdin,
stdout, and stderr command streams, `Start: true`, direct process start error,
one exact request attempt, absence of another request, exact unsupported-
platform error without a process attempt, safe zero dependency, and every
public behavior remain intact. There is no fallback, retry, wrapping, logging,
synchronous execution, cleanup, or global-state change.

Five strengthened focused top-level contracts keep the suite at 388 tests.
They prove exact system-runner identity with nil process-dependency stdout,
unchanged exact runtime platform selection, complete platform requests, URL
arguments, stream identities, asynchronous attempts, direct errors, safe zero
behavior, complete caller-owned dependency preservation, non-empty population,
absence of another request, and the exported wrapper's exact composition. The
recording double preserves the complete caller-owned browser dependency. The
contracts launch no external program, touch no network, and write no repository
fixture.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2 stays
zero, regenerated exact Q1.3 improves from 8 of 30 to 7 of 29, Q1.4 stays 7 of
8, exact Q2.1 stays 0 of 8, and Q3.4 stays zero phrases across 81 Markdown
files. The unchanged import-aware scanner no longer follows `os.Stdout`
provenance through the private browser-launcher runner-only production
dependency. Clock and server remain absent, so Q1.3 remains the single non-
comparable ratchet. The scanner, inventory, baseline, and audit apparatus are
unchanged.

The clean move-64 gate passes focused browser-launcher, process, and relevant
caller tests; API/CLI and fresh subprocess compatibility; build, complete and
uncached tests, race, vet, pinned lint, all 62 launcher controls, Make and all
production-script meta-contracts, all four host acceptance flows, the repaired
15-control audit meta-suite, and empty-HOME count-2. One complete preflight
invocation hit the established nested partial-raw-log signal-fixture flake; its
immediate complete unchanged rerun passed all 62 controls. Compatibility used
the exact pinned `apidiff` path after the first shell environment lacked that
executable on `PATH`. The clean full audit exits 1 for the same 15 documented
findings, never 2, with L0 8 of 8, five improved, two held, zero regressed, one
non-comparable ratchet, and zero dirty paths.

Product commit: `0b96f10`.

The authoritative clean full audit from `0b96f10` has structured scorecard
SHA-256
`f817974080ab4b1641155e24e7d2b7f4d0fddf1ae7c06d7c4922320a70f5a842`;
the measured commit tree is `8a7180aa20082d504a47ce8bdb4a12d405c9bf50`.

Move 65 changes only the private filesystem adapter forwarding surface and
shell Unzip's destination-file open call. New `filesystem.OpenFileAsFile`
delegates once to the complete existing `FileSystem.OpenFile` operation with
the exact path, flags, and mode while returning its exact file through the
existing `filesystem.File` interface and preserving the exact error. The
complete `FileSystem` interface, system implementation, exported
`filesystem.OpenFile`, public `file.OpenFile`, and every other operation and
caller are unchanged.

`unzipWithDependencies` now uses that helper only for its one destination
open. Exact source and destination bytes, archive open and deferred close,
zip-slip boundary and error, archive directory metadata classification, joined
paths, recursive directories and modes, ordered filenames, destination flags,
entry modes and opens, ignored copy results, output and reader closes, attempts,
operation order, direct errors, partial filenames, error precedence, traversal,
safe zero behavior, system selection, exported `Unzip`, and its caller remain
unchanged. There is no selection, fallback, retry, wrapping, logging, cleanup,
or global-state change.

Two strengthened focused top-level contracts keep the suite at 388 tests. They
prove safe zero behavior, exact complete dependency delivery and preservation,
path, flags, mode, returned `filesystem.File` identity, exact error, non-empty
population, absence of another filesystem operation, the exact private helper
placement, and every existing Unzip archive, path, directory, entry, copy,
close, order, partial-result, error-precedence, and traversal guarantee. The
contracts launch no external program, touch no network, and write no repository
fixture.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2 stays
zero, regenerated exact Q1.3 improves from 7 of 29 to 5 of 27, Q1.4 stays 7 of
8, exact Q2.1 stays 0 of 8, and Q3.4 stays zero phrases across 83 Markdown
files. The unchanged import-aware scanner no longer follows concrete
`*os.File` provenance through Unzip's adapter `Copy` and `Close` requests.
Clock and server remain absent, so Q1.3 remains the single non-comparable
ratchet. The scanner, inventory, baseline, and audit apparatus are unchanged.

The clean move-65 gate passes focused filesystem, Unzip, Spring, and relevant
caller tests; API/CLI and fresh subprocess compatibility; build, complete and
uncached tests, race, vet, pinned lint, all 62 launcher controls, Make and all
production-script meta-contracts, all four host acceptance flows, the repaired
15-control audit meta-suite, and empty-HOME count-2. The clean full audit exits
1 for the same 15 documented findings, never 2, with L0 8 of 8, five improved,
two held, zero regressed, one non-comparable ratchet, and zero dirty paths.

Product commit: `40cece5`.

The authoritative clean full audit from `40cece5` has structured scorecard
SHA-256
`1bae6abec78d2306c2205bdf1f96b7bffdabcef7e82a232a2aedafe536f3bf56`;
the measured commit tree is `c1009903d65b062786604f71d8f1fa89d2dfe61d`.

Move 66 changes only the process adapter's standard-output forwarding surface
and Maven's private production dependency selection. New
`process.SystemStdout` returns the exact existing `Stdout(System(), true)`
composition. It performs no fallback, wrapping, buffering, copying, logging,
cleanup, retry, or global-state mutation, and the complete process dependency,
runner, system selectors, execution path, and implementation remain unchanged.

Private `systemRunOnDependencies` carries the exact existing
`process.SystemRunner()` and that exact system standard output. Public
`maven.RunOn` replaces only its complete `process.System()` selection with that
private selector; private `runOn`, every callback and caller, and the return
path are unchanged. Every arbitrary command and argument byte and order,
repository value, exact project directory, info log before execution,
conditional stdout selection for debug and trace only, nil stdin and stderr,
synchronous one-attempt execution, and direct process error remain exact.

Strengthened focused contracts keep the suite at 388 tests. They prove exact
`os.Stdout` identity through the existing process composition, exact system
runner identity, complete caller-owned dependency delivery and preservation,
every log and command field, the complete log-level matrix, callback
construction and invocation timing, one exact attempt, direct error, safe zero
behavior, non-empty populations, exported public composition, and absence of
another process operation. The contracts launch no external program, touch no
network, and write no repository fixture.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 25, Q1.2 stays
zero, regenerated exact Q1.3 improves from 5 of 27 to 4 of 27, Q1.4 stays 7 of
8, exact Q2.1 stays 0 of 8, and Q3.4 stays zero phrases across 84 Markdown
files. The unchanged import-aware scanner no longer follows concrete
`os.Stdout` provenance through Maven's complete public `RunOn` dependency.
Clock and server remain absent, so Q1.3 remains the single non-comparable
ratchet. The scanner, inventory, baseline, and audit apparatus are unchanged.

The clean move-66 gate passes focused process, Maven, and relevant caller
tests; API/CLI and fresh subprocess compatibility; build, complete and uncached
tests, race, vet, pinned lint, all 62 launcher controls, Make and all
production-script meta-contracts, all four host acceptance flows, the repaired
15-control audit meta-suite, and empty-HOME count-2. Two complete Make attempts
hit the established nested partial-raw-log signal-fixture flake before an
unchanged complete `make test` rerun passed. The clean full audit exits 1 for
the same 15 documented findings, never 2, with L0 8 of 8, five improved, two
held, zero regressed, one non-comparable ratchet, and zero dirty paths.

Product commit: `c627a3e`.

The authoritative clean full audit from `c627a3e` has structured scorecard
SHA-256
`1023a95ee7e4d7a10e95ae647821892dae71b8815f35062ee5e4b8cbd9c207d0`;
the measured commit tree is `fc8ae247df072a17c61d3e68a8359dfb4f1d7765`.

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
`structurizr.RunWithOutputToFile` now reuses `WriteFile` without changing
command construction or execution, buffer capture, early process errors,
caller path, stdout bytes, `0644` mode, discarded write errors, or callers.
`maven.WriteGraphStyles` now reuses `WriteFile` without changing JSON marshal
sequencing, its early error, exact caller-composed target, existence selection,
no-write-when-present behavior, output bytes, mode, exact error, Graph behavior,
or command execution.
`LocalConfigDir.TouchFile` now reuses `WriteFile` without changing directory
evaluation, path selection, config initialization, default cloud URL, YAML
marshal sequencing, logging, file creation, exact output bytes, mode, write
error, close behavior, `UpdateLocalConfig`, or callers.
The isolated `structurizr-cli export`, Graphviz `dot`, and ignored macOS `open`
requests in `cmd/plugin_diagrams.go` now use the process adapter; the Graphviz
stdout write also uses the filesystem adapter. Maven command's conditional
standard-output selection now uses the complete process dependency. Shell
Unzip's archive-entry reader selection now uses the complete filesystem
dependency without changing any other archive operation. Public shell `Run`
and the four public shell Git execution wrappers now use the process adapter's
exact runner-only system dependency without inheriting its unused
standard-output capability. The private profile-editor selector now also uses
the exact runner-only system dependency while retaining its separate exact
`os.Stdin` and `os.Stdout` command streams. The three private plugin-diagrams
selectors now also use the exact runner-only system dependency while retaining
their exact export, Graphviz, open, write, error, discovery, and iteration
behavior and the Graphviz selector's complete exact filesystem dependency.
The private browser-launcher selector now also uses the exact runner-only system
dependency while retaining runtime platform selection, the exact asynchronous
commands and URLs for Linux, Windows, and macOS, direct start errors,
unsupported-platform behavior, safe zero behavior, and public `OpenBrowser`
behavior. Shell Unzip's destination open now also returns through the existing
`filesystem.File` interface while retaining the same complete
`FileSystem.OpenFile` operation, path, flags, mode, returned file and error,
archive and entry opens, path validation, directories, copy and close calls,
order, ignored and direct errors, partial filenames, traversal, exported
`Unzip`, complete adapter interface and system implementation, exported
`filesystem.OpenFile`, public `file.OpenFile`, and every public behavior.
The P3 exit is satisfied: process exits outside `main` are zero, every P3
production-effect move is recorded, its migrated call sites no longer appear
in Q1.3, and the CLI/API contracts remain compatible. P4 is now active. Begin
it with one focused clock-adapter move for Spring's archive-path timestamp;
that move is complete, and Kibana sleep and server behavior remain later
coherent P4 moves.

Exit: process exits outside `main` reach zero, migrated call sites disappear
from Q1.3, their seam swaps are killed, and CLI/API contracts stay compatible.

### P4 - Finish Absolute L1

Status: complete.

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

Move 1 introduces only the declared `internal/adapter/clock` current-time
boundary and migrates only Spring's private archive-path timestamp selection.
`clock.Now` returns the exact injected `time.Time`; its zero dependency returns
the deterministic exact zero time, and `clock.System` selects a private system
clock whose only read directly returns `time.Now()`. It performs no fallback,
truncation, rounding, timezone conversion, caching, monotonic rewriting,
logging, cleanup, retry, or global-state mutation.

Private `archivePathDependencies` retains the exact complete filesystem system
dependency and adds only the exact complete system clock. The existing working-
directory request and direct error still happen before any clock read. After
directory success, exactly one clock read supplies the unchanged `Unix()`
second to the unchanged `spring-%d.zip` name and slash composition. Named
returns, arbitrary directory bytes, pre-epoch and subsecond conversion,
`archivePath`, every caller, archive creation, download, unzip, deletion,
logging, errors, and public behavior remain exact.

Six focused contracts bring the suite to 394 tests across 20 of 26 packages.
They prove deterministic safe zero behavior, a direct system-time value within
a bounded before/after observation, exact arbitrary injected times, complete
dependency delivery and preservation, directory-before-clock order, no clock
read after a directory error, one read after success, exact path and returns,
non-empty populations, public composition, and absence of another operation.
They launch no external program, touch no network, and write no repository
fixture.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 is now 6 of 26 because
the tested clock adapter adds one package, Q1.2 stays zero, regenerated exact
Q1.3 improves from 4 of 27 to 3 of 27, Q1.4 stays 7 of 8, exact Q2.1 stays 0
of 8, and Q3.4 stays zero phrases across 85 Markdown files. The declared clock
adapter is now present; the server adapter remains absent, so Q1.3 remains the
single non-comparable ratchet. The scanner, inventory, baseline, and audit
apparatus are unchanged.

The clean move-1 gate passes focused clock, Spring, and relevant caller tests;
API/CLI and fresh subprocess compatibility; build, complete and uncached tests,
race, vet, pinned lint, all 62 launcher controls, Make and all production-script
meta-contracts, all four host acceptance flows, the repaired 15-control audit
meta-suite, and empty-HOME count-2. The first API compatibility invocation
lacked the exact pinned `apidiff` executable and passed after rebuilding that
pin outside the worktree. The first preflight invocation used the sandbox-
blocked default lint cache and its isolated-cache rerun passed. The first
complete `make test` invocation hit the established nested partial-raw-log
signal-fixture flake; its immediate unchanged complete rerun passed all 62
controls. The clean full audit exits 1 for 14 documented non-passing criteria,
never 2, with L0 8 of 8, five improved, two held, zero regressed, one non-
comparable ratchet, and zero dirty paths.

Product commit: `91422eb`.

The authoritative clean full audit from `91422eb` has structured scorecard
SHA-256
`60bc1cf4994c63444d9b841f1dfd971d99283b21aad02cf3a859305c41fc62c4`;
the measured commit tree is `eedf149499ead042a28352a4c35b9d827f41a58f`.

Move 2 extends only the existing clock adapter with exact duration sleep and
migrates only Kibana public `POST`'s fixed retry delay. `clock.Sleep` passes the
exact requested `time.Duration` to its injected sleeper once; its zero
dependency is a deterministic no-op. The private system clock now implements
both existing exact `Now` and exact `Sleep`, with production sleep delegating
its argument directly to `time.Sleep` once and performing no fallback,
normalization, clamping, rounding, conversion, caching, logging, cleanup,
retry, or global-state mutation.

Private `postDependencies` owns the complete existing `internalPOST` operation
and the complete clock dependency, and its production selector chooses exactly
`internalPOST` plus `clock.System()`. Public `POST` composes only that selector
and the private dependency-taking operation. One first request still happens
before hit-list selection. A non-empty first hit list returns the exact first
error and response with no print, sleep, or retry. An empty first hit list,
including with a non-nil first error, still prints `sleep and retry`, sleeps
once for exactly 15 seconds, retries the same complete request once, and
returns the exact second error and response. The first result is discarded and
there is no third request. `internalPOST`, `internalPost`, request execution,
parsing, response bodies, callers, and public APIs are unchanged.

Seven net new focused contracts bring the suite to 401 tests across 20 of 26
packages. They prove safe zero sleep, arbitrary exact duration delivery, direct
system delegation, complete caller-owned dependency delivery and preservation,
first-request ordering, both hit-selection branches, retry despite a first
error, exact print/sleep/retry order and duration, exact returned identities,
at most two requests, non-empty recorded populations, public composition, and
absence of another operation. They perform no real sleep, launch no external
program, touch no network, and write no repository fixture.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 stays 6 of 26, Q1.2
stays zero, regenerated exact Q1.3 improves from 3 of 27 to 2 of 27, Q1.4
stays 7 of 8, exact Q2.1 stays 0 of 8, and Q3.4 stays zero phrases across 86
Markdown files. Both remaining Q1.3 violations are the unchanged server calls
in `pkg/webservice/api.go`; the declared server adapter alone remains absent.
The scanner, inventory, baseline, and audit apparatus are unchanged.

The clean move-2 gate passes focused clock, Kibana, and relevant caller tests;
API/CLI and fresh subprocess compatibility; build, complete and uncached tests,
race, vet, pinned lint, all 62 launcher controls, Make and all production-script
meta-contracts, all four host acceptance flows, the repaired 15-control audit
meta-suite, and empty-HOME count-2. One preflight launcher signal-fixture run
failed at its documented partial-raw-log race; its immediate unchanged complete
rerun passed all 62 controls and the full preflight. The clean full audit exits
1 for 14 documented non-passing criteria, never 2, with L0 8 of 8, five
improved, two held, zero regressed, one non-comparable ratchet, and zero dirty
paths.

Product commit: `4c0d97a`.

The authoritative clean full audit from `4c0d97a` has structured scorecard
SHA-256
`be0142a87d2ccb5ff1a3d9414ba9d376843f7ab8b92e5ba0f75f512913e0b0db`;
the measured commit tree is `2bff9b81b137130ba026772a37526bbb24595c32`.
The focused seven-criterion scorecard SHA-256 is
`10cf5c8891a57a0ac883b9b685afd8ad719d14adffadf5dcfa5ff0db6631e63a`.

Move 3 adds the declared `internal/adapter/server` boundary and migrates only
webservice's exact `ListenAndServe` and `Shutdown` effects. The unchanged
package-level `*http.Server` still has address `:7999`; the same four handler
paths and callbacks register in the same order before one listen attempt. A
non-nil listen error still reaches `log.Print` exactly once, nil prints
nothing, and stop still derives one context from `context.Background()` with
exactly five seconds, defers cancellation, attempts one shutdown, and ignores
its exact error.

The adapter's production operations call the corresponding method once on the
exact caller-selected server and return its exact error. Safe-zero dependencies
select no server, perform no network operation, and return deterministic nil.
Private start and stop compositions retain complete caller-owned operation,
server-selector, handler, logger, background, and timeout dependencies. The
server selector is dereferenced only inside the declared adapter, preserving
the exact global identity while allowing the unchanged import-aware scanner to
recognize the boundary.

Thirteen focused contracts bring the suite to 414 tests across 21 of 27
packages. They record exact server and context identities, direct system
delegation, complete dependency preservation, handler/listen and
background/timeout/shutdown/cancel order, nil/non-nil logging, ignored shutdown
errors, safe-zero behavior, public composition, non-empty populations, and the
absence of an unrelated operation. They open no socket, make no request,
perform no wait, launch no program, mutate no production server, and write no
repository fixture.

Q0.6 holds at 26 guarded safe-writer sites, 21 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 is 6 of 27, Q1.2 is
zero, exact Q1.3 improves from 2 of 27 to 0 of 27 with all five adapter paths
valid, Q1.4 stays 7 of 8, exact Q2.1 stays 0 of 8, and Q3.4 is zero phrases
across 87 Markdown files. The scanner, inventory, baseline, audit apparatus,
server construction, and address policy are unchanged.

The clean move-3 gate passes focused server/webservice and caller tests;
API/CLI and fresh subprocess compatibility; build, complete tests, race, vet,
pinned lint, `make test`, complete preflight, all 62 launcher controls, Make
and production-script contracts, all four host acceptance flows, the repaired
15-control audit meta-suite, and empty-HOME count-2. An initial preflight
invocation used the wrong tool variable names before the correctly pinned
complete rerun passed. The clean full audit exits 1 for 13 documented
non-passing criteria, never 2, with L0 8 of 8, six improved, two held, zero
regressed, zero non-comparable ratchets, and zero dirty paths.

Product commit: `9601d29`.

The authoritative clean full audit from `9601d29` has structured scorecard
SHA-256
`010923f3718d99723c99c76a6b541669cb7f99d7a4a8eb119044705be005cdfa`;
the measured commit tree is `fd78ef577b886a423d2964f4352503d2aac8992b`.
The focused seven-criterion scorecard SHA-256 is
`d52e2d34017dac46134a4eacc2fca628e4e39bed995995ab7837217796ca65ed`.

Move 4 adds only `pkg/sorting/sort_test.go`; `pkg/sorting/sort.go`, its API,
and every caller are unchanged. Nine focused characterization contracts prove
exact slice length, value-receiver and shared-slice behavior, in-place swap and
full sort, strict less-than comparison, complete concatenated key construction,
all scope weights, matching `1-` and nonmatching `100-` group prefixes,
case-sensitive substring selection, empty sort-key behavior, group and artifact
tie-breaks, and non-empty table populations. The package reaches 100% statement
coverage without a request, socket, external program, wait, global-state
mutation, or repository fixture write.

The existing implementation compares its formatted keys lexically. Therefore
an unknown scope's textual `100:` prefix sorts before the empty scope's `10:`
prefix. The new contracts deliberately preserve that legacy result; changing
it to a preferred numeric order is a later product decision, not part of this
coverage move. Dependency fields outside scope, group ID, and artifact ID also
remain outside the comparison key.

The suite now has 423 tests across 22 of 27 packages. Q0.6 holds at 26 guarded
safe-writer sites, 21 write and 5 copy, with zero skipped tests and zero unsafe
direct test writes. Q1.1 improves from 6 of 27 to 5 of 27, Q1.2 remains zero,
exact Q1.3 remains 0 of 27 with all five adapter paths valid, Q1.4 remains 7 of
8, exact Q2.1 remains 0 of 8, and Q3.4 remains zero phrases across 88 Markdown
files. The scanner, inventory, baseline, audit apparatus, mutation harnesses,
and production code are unchanged.

The clean move-4 gate passes focused sorting, config, and Maven caller tests at
100% sorting coverage; API/CLI, CLI surface, and fresh subprocess compatibility;
build, complete and uncached tests, race, vet, pinned lint, `make test`, all 62
launcher controls, Make and production-script contracts, all four host
acceptance flows, the repaired 15-control audit meta-suite, and empty-HOME
count-2. The first integrated preflight invocation supplied the pinned linter
as a Make command-line override, which leaked into and invalidated the
fake-linter negative control; the unchanged rerun supplied the same pin through
the environment and passed the complete preflight.

Product commit: `cbb620a`.

The authoritative clean full audit from `cbb620a` has structured scorecard
SHA-256
`4f7550ad8fd184f09237c5f8ae339a1d746ab8c2b366706d4cd9c5feb79ae0af`;
the measured commit tree is `18d1ff3d026dc358f7f854b338fd7d5e2d411c68`.
It exits 1 for 13 documented findings, never 2, with L0 8 of 8, six
improved ratchets, two held, zero regressed, zero non-comparable, and zero dirty
paths. The focused seven-criterion scorecard SHA-256 is
`7f6fbf4354a7b7040c350ba0786b468a054d2febae6b7a3c9031507e16471723`.
Continue P4 with one test-only characterization move for `pkg/resources`.

Move 5 adds only `pkg/resources/resources_test.go`; production
`pkg/resources/resources.go`, its two exported function signatures, config and
file boundaries, callers, and observable filesystem behavior are unchanged.
Three focused characterization contracts prove the exact one-time
`Implementation` and `Dir` selection with no `FilePath` call, representative,
empty, and trailing-separator `LocalDir` path bytes, literal resource/filename
slash composition, successful empty, arbitrary, non-UTF-8, and nested resource
reads, exact byte-to-string preservation, direct missing `*os.PathError` shape,
and non-empty table populations. All fixtures live below `t.TempDir()`, their
parents are prepared first, and their files use the central guarded writer.

The suite now has 426 tests across 23 of 27 packages. Q0.6 has 27 guarded safe-
writer sites, 22 write and 5 copy, with zero skipped tests and zero unsafe direct
test writes. Q1.1 improves from 5 of 27 to 4 of 27, Q1.2 remains zero, exact
Q1.3 remains 0 of 27 with all five adapter paths valid, Q1.4 remains 7 of 8,
exact Q2.1 remains 0 of 8, and Q3.4 remains zero phrases across 89 Markdown
files. The scanner, inventory, baseline, audit apparatus, mutation harnesses,
and production code are unchanged.

The clean move-5 gate passes focused resources and template caller tests at
100% resources statement coverage; API/CLI, CLI surface, and fresh subprocess
compatibility; build, complete and uncached tests, race, vet, pinned lint,
`make preflight`, `make test`, all 62 launcher controls, Make and production-
script contracts, all four host acceptance flows, the repaired 15-control audit
meta-suite, and empty-HOME count-2. The first focused coverage command and first
preflight invocation used sandbox-blocked default build/lint caches; unchanged
isolated-cache reruns passed. The first `make test` invocation hit the documented
nested partial-raw-log signal-fixture race; its immediate unchanged complete
rerun passed all controls.

Product commit: `64c6189`.

The authoritative clean full audit from `64c6189` has structured scorecard
SHA-256
`2bbbc911efefe7427cd3413b73b0a96583a1111460e57f5cdfb9b145b6818a95`;
the measured commit tree is `88da045036be8849a8a3ffb327e01bb5f23c81b5`.
It exits 1 for 13 documented findings, never 2, with L0 8 of 8, six
improved ratchets, two held, zero regressed, zero non-comparable, and zero dirty
paths. The focused seven-criterion scorecard SHA-256 is
`9eedfafb6977fae6a0e27b878174b49ecf3536b7d8997b2746df41f722f34430`.
Continue P4 with one test-only characterization move for
`pkg/webservice/templates`.

Move 6 adds only `pkg/webservice/templates/templates_test.go`; all four
production template sources, exported `Generate` and `Upgrade` untyped string
constants, private header/body/footer fragments, API callers, parser choices,
HTML, form behavior, and observable output remain unchanged. Three focused
characterization contracts prove exact complete lengths and SHA-256 identities
for both public templates and all four private fragments, exact shared-header +
private-body + shared-footer composition, leading and trailing bytes, form
actions, input names, template actions and counts, CDN references, and
successful parsing through the callers' exact `text/template` Generate and
`html/template` Upgrade engines. Every table rejects an empty population before
iteration. The package contains no statements, so Go correctly reports
`[no statements]` coverage while the new test file closes its Q1.1 package gap.

The existing unquoted Generate value attributes and caller-selected
`text/template` parsing are deliberately preserved. Preferred quoting or HTML
escaping is a later product decision requiring explicit authority; this move
does not change bytes to make that preference pass.

The suite now has 429 tests across 24 of 27 packages. Q0.6 holds at 27 guarded
safe-writer sites, 22 write and 5 copy, with zero skipped tests and zero unsafe
direct test writes. Q1.1 improves from 4 of 27 to 3 of 27, Q1.2 remains zero,
exact Q1.3 remains 0 of 27 with all five adapter paths valid, Q1.4 remains 7 of
8, exact Q2.1 remains 0 of 8, and Q3.4 remains zero phrases across 90 Markdown
files. The scanner, inventory, baseline, audit apparatus, mutation harnesses,
production templates, API globals, and callers are unchanged.

The clean move-6 gate passes focused templates, API, and webservice caller
tests; API/CLI, CLI surface, and fresh subprocess compatibility; build,
complete and uncached tests, race, vet, pinned lint, `make preflight`, all 62
launcher controls, Make and production-script contracts, all four host
acceptance flows, the repaired 15-control audit meta-suite, and empty-HOME
count-2. The first `make test` invocation hit the documented nested partial-
raw-log signal-fixture race at launcher control 26; its immediate unchanged
complete rerun passed all controls.

Product commit: `8eeeb2a`.

The authoritative clean full audit from `8eeeb2a` has structured scorecard
SHA-256
`213df9ece4eb7c231ec3faa8f76e9180e9215a6ef0fc7d83039911940224c7e7`;
the measured commit tree is `9c7322b185f485b4202cec0691ce9fd39c8dbdc0`.
It exits 1 for 13 documented findings, never 2, with L0 8 of 8, six
improved ratchets, two held, zero regressed, zero non-comparable, and zero dirty
paths. The focused seven-criterion scorecard SHA-256 is
`aaa5507955ca84ef255432c998ee2c3092776caeb737754736e78752bc49a57a`.
Continue P4 with one test-only characterization move for
`pkg/webservice/api`.

Move 7 adds only `pkg/webservice/api/handlers_test.go`; all three production
API files, exported handler signatures, `GenerateOptions`, `GOptions`,
`CallbackChannel`, `CurrentProject`, routes, templates, callers, HTML, form
parsing, response text, callback timing, config, Spring, and observable behavior
are unchanged. Five focused characterization contracts bring the suite to 434
tests across 25 of 27 packages and cover all four handlers at 100% statement
coverage.

Generate GET is executed through the existing `text/template` handler path and
pins the exact representative 2,957-byte response with SHA-256
`9a9eb38ab026e751b06dee004e84d8db9a427454c5e33f4b864ebc83635399c1`,
status 200, one cloud-template selection, and the existing raw, unquoted project,
template, dependency-name, and dependency-ID bytes. Upgrade GET is executed
through the existing `html/template` path and pins status 200 plus the exact
788-byte response with SHA-256
`ae9194d9a6734d1fe568692e11327fc1b31b3f0d797fd558fd6ceba3ffa58ecc`.
Preferred validation, attribute quoting, or escaping remains a later product
decision requiring explicit scope.

Generate POST proves body parsing, body-over-query first-value scalar selection,
exact overwrite of the six existing scalar project fields, preservation of
every other field, and ordered duplicate-preserving append of every submitted
`templates` and `dependencies` value. Both POST handlers return exact status
200 and `OK`, return before a test-owned unbuffered callback receiver exists,
then deliver exactly one `true` callback. Upgrade POST preserves its existing
choice not to parse even a malformed form body. The unbuffered rendezvous
releases each only sender; no sleep, timed wait, socket, external program,
logger mutation, fixture write, or goroutine leak is introduced. Each subtest
restores the exact `GOptions`, `CallbackChannel`, and `CurrentProject` values,
and the outer contracts verify restored pointer/channel identities and values.
Recorded callbacks and every iterated table reject empty populations.

Q0.6 holds at 27 guarded safe-writer sites, 22 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 improves from 3 of 27 to
2 of 27, Q1.2 remains zero, exact Q1.3 remains 0 of 27 with all five adapter
paths valid, Q1.4 remains 7 of 8, exact Q2.1 remains 0 of 8, and Q3.4 remains
zero phrases across 91 Markdown files. The scanner, inventory, baseline, audit
apparatus, mutation harnesses, production code, and all completed contracts are
unchanged.

The clean move-7 gate passes focused API, webservice, templates, and command
caller tests at 100% API-handler coverage; API/CLI, CLI surface, and fresh
subprocess compatibility; build, complete and uncached tests, race, vet, pinned
lint, `make preflight`, `make test`, all 62 launcher controls, Make and
production-script contracts, all four host acceptance flows, the repaired
15-control audit meta-suite, and empty-HOME count-2. Three initial standalone
launcher invocations hit the documented partial-raw-log signal-fixture race:
twice at control 26 and once inside the nested control-50 run; the fourth
unchanged complete run passed all 62 controls. The first `make test` invocation
hit the same race at control 26; its unchanged complete rerun passed.

Product commit: `4da64d2`.

The authoritative clean full audit from `4da64d2` has structured scorecard
SHA-256
`0682df76893e0eac05302870ea986b1f18cab32e7513a7e0e91ab85da5f16c38`;
the measured commit tree is `d23edc06e2f958c43f2add2a43c540f48c7145e9`.
It exits 1 for 13 documented findings, never 2, with L0 8 of 8, six
improved ratchets, two held, zero regressed, zero non-comparable, and zero dirty
paths. The focused seven-criterion scorecard SHA-256 is
`54b314c71a0e6bce29b84647d347d15b643ecf3c8baf982876411ef7b72cb428`.
Continue P4 with one test-only characterization move for `pkg/logger`.

Move 8 adds only `pkg/logger/logger_test.go`; both production logger files,
all exported signatures, package globals and pointer identities, hook
installation, levels, fields, formatter and output selection, callers, logging
bytes, and observable behavior are unchanged. Six focused characterization
contracts bring the suite to 440 tests across 26 of 27 packages and cover
`pkg/logger` at 100% statement coverage.

The contracts pin the exact ordered `Info` and `Warn` collector levels,
append-only capture of the exact entry pointers without copying, and the
collector hook installed only on the private package logger. They preserve and
characterize the existing private-versus-standard logrus state split:
`DebugLogger` returns an empty allocated field map and enables debug only on
the private logger; `Context` returns empty data; `ExternalError` keeps its
representative exact multiline and empty text without wrapping; JSON selection
uses the exact zero `logrus.JSONFormatter`; field-logger transitions are exact;
and `StdOut` depends only on standard-logrus debug enablement and returns the
exact `os.Stdout` pointer or nil. `LogEntries` exposes the collector's exact
slice backing and entry pointers without copying. Preferred logger ownership,
isolation, copying, formatting, or output design remains a later product
decision requiring explicit scope.

Each mutating subtest snapshots and restores the exact `fieldLogger`,
`collector`, private `log`, and touched private and standard-logrus level,
formatter, output, hooks, and caller-reporting state; its outer contract proves
the restored identities and values. No parallel test, socket, external
program, sleep, timed wait, fixture write, hook leak, or entry leak is added.
Every iterated table and recorded-entry population rejects an empty input.

Q0.6 holds at 27 guarded safe-writer sites, 22 write and 5 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 improves from 2 of 27 to
1 of 27, Q1.2 remains zero, exact Q1.3 remains 0 of 27 with all five adapter
paths valid, Q1.4 remains 7 of 8, exact Q2.1 remains 0 of 8, and Q3.4 remains
zero phrases across 92 Markdown files. The scanner, inventory, baseline, audit
apparatus, mutation harnesses, production code, and all completed contracts are
unchanged.

The clean move-8 gate passes focused logger and relevant command and package
caller tests at 100% logger coverage; API/CLI, CLI surface, and fresh
subprocess compatibility; build, complete and uncached tests, race, vet,
pinned lint, `make preflight`, `make test`, all 62 launcher controls, Make and
production-script contracts, all four host acceptance flows, the repaired
15-control audit meta-suite, and empty-HOME count-2. The first focused coverage
invocation used the sandbox-blocked default Go build cache; its unchanged
isolated-cache rerun passed. Initial launcher and `make test` invocations hit
the documented partial-raw-log signal-fixture race at control 26 or the nested
control-50 run; unchanged complete reruns passed all controls.

Product commit: `67344a8`.

The authoritative clean full audit from `67344a8` has structured scorecard
SHA-256
`38e938f060312840cb244a575536aa613266b172de98542cb2bb2e0998da576b`;
the measured commit tree is `e85dc852f5875c812509b8f7e1b08520d7f0251d`.
It exits 1 for 13 documented findings, never 2, with L0 8 of 8, six
improved ratchets, two held, zero regressed, zero non-comparable, and zero dirty
paths. The focused seven-criterion scorecard SHA-256 is
`5c640e0717d17f3200b1acfed83171f1cc824529ae2af18ab9b31a7b4fa77ee6`.
Continue P4 with one test-only characterization move for `pkg/context`.

Move 9 adds only `pkg/context/context_test.go`; both production Context files,
all exported fields and method signatures, the private package logger and its
initialization, config and project ownership, file and Maven calls, callers,
logging, errors, writes, globals, and observable behavior are unchanged. Nine
focused characterization contracts bring the suite to 449 tests across all 27
packages and cover `pkg/context` at 97.6% statement coverage.

The contracts pin package logger initialization plus exact `SetLogger`
assignment and restoration; zero and directly owned public Context state;
non-recursive and recursive project discovery; existing-project append,
filesystem traversal, and partial-result order; missing-root and invalid-project
behavior; empty project handling; cloud-default merging and errors; project-type
skips; stealth, dry-run, dirty-repository, job, nil-job, and write behavior;
root-project selection; profile assignment, creation, retention, and failure;
configured authenticated and anonymous repositories; legacy default repository
selection; and representative exact log levels, messages, and errors. Value
copies, shared pointer fields, callback order, early returns, continuations, and
non-empty iterated and recorded populations are explicit.

All test-created directories and files remain below `t.TempDir()` and use the
two central guarded writers. Every logger, environment, homedir cache, and
working-directory mutation restores its exact original identity and value. No
parallel test around process state, socket, network access, external program,
sleep, timed wait, repository fixture, leaked hook, entry, file, environment,
or process-state change is added. The two uncovered statements are legacy
error logs that current interfaces cannot induce deterministically:
`file.FindAll` suppresses callback walk errors, and `maven.DefaultRepository`
derives its home from `os/user.Current().HomeDir` and offers no injectable
failure seam. Preferred error propagation and a caller-owned Maven home or
repository seam remain later product decisions requiring explicit scope.

Q0.6 improves to 29 guarded safe-writer sites, 23 write and 6 copy, with zero
skipped tests and zero unsafe direct test writes. Q1.1 improves from 1 of 27 to
0 of 27, Q1.2 remains zero, exact Q1.3 remains 0 of 27 with all five adapter
paths valid, Q1.4 remains 7 of 8, exact Q2.1 remains 0 of 8, and Q3.4 remains
zero phrases across 93 Markdown files. The scanner, inventory, baseline, audit
apparatus, mutation harnesses, production code, and all completed contracts are
unchanged.

The clean move-9 gate passes focused context and command caller tests at 97.6%
context coverage; API/CLI, CLI surface, and fresh subprocess compatibility;
build, complete and uncached tests, race, vet, pinned lint, `make preflight`,
`make test`, all 62 launcher controls, Make and production-script contracts,
all four host acceptance flows, the repaired 15-control audit meta-suite, and
empty-HOME count-2. The first focused audit exposed three direct fixture
directory writes; replacing those test-only setup calls with the existing
guarded copy writer and amending the single product commit restored exact Q0.6
before the authoritative measurements. The first `make test` invocation hit
the documented partial-raw-log signal-fixture race at launcher control 26; its
immediate unchanged standalone launcher run and complete Make rerun passed.

Product commit: `746a5ab`.

The authoritative clean full audit from `746a5ab` has structured scorecard
SHA-256
`48a58e58abf53254ec318de93266a683f60e6ac4a194e02b187c76c170a823b2`;
the measured commit tree is `7bfabae149bffc5dfd2ad59f90abb08e29c982f2`.
It exits 1 for 13 documented findings, never 2, with L0 8 of 8, six
improved ratchets, two held, zero regressed, zero non-comparable, and zero dirty
paths. The focused seven-criterion scorecard SHA-256 is
`f937c9be84913156be11dfa9a937565959cfcde6ba1b9389bea7b2f337e0d7c6`.
Move 10 changes only the private package-level interactive HTTP server address
from `:7999` to exact IPv4 loopback `127.0.0.1:7999`. The port constant, exact
`*http.Server` identity, selector and adapter, four handler paths and order,
single listen attempt, nil/non-nil error logging, background shutdown context,
five-second timeout, deferred cancellation, ignored shutdown error, endpoint
URIs, browser request, callback blocking, public signatures, and command caller
flow remain unchanged. No contract opens a socket, launches a browser, or
starts the real server.

Four focused contracts bring the suite to 453 tests across all 27 packages.
They prove the exact address and unchanged start/stop selection, endpoint URIs,
server/browser/callback composition, standalone and project blocking flow,
interactive build option delivery and validation order, and interactive
upgrade guards and callback selection. The new regular non-executable
`scripts/mutate-interactive-build` P4 seam driver contains the eighth inventory
label exactly once, declares no mutation operation, and runs 21 named command,
webservice, server-adapter, API, and callback contracts. Its reciprocal
meta-contract proves the driver shape, label, absence of premature mutation,
successful JSON execution, non-empty population, and every named test run.

Q0.6 remains 29 guarded safe-writer sites, 23 write and 6 copy, with zero
skipped tests and zero unsafe direct test writes. Q0.8 remains reciprocal for
all 13 production scripts. Q1.1 remains 0 of 27, Q1.2 remains zero, exact Q1.3
remains 0 of 27 with all five adapter paths valid, Q1.4 improves from 7 of 8 to
8 of 8, exact Q2.1 remains 0 of 8, and Q3.4 remains zero phrases across 94
Markdown files. The scanner, inventory, baseline, audit apparatus, executable
mutation harness population, and acceptance behavior are unchanged.

The clean move-10 gate passes focused coverage and shuffled tests; API/CLI,
surface, and fresh subprocess compatibility; build, complete tests, race, vet,
pinned lint, `make preflight`, `make test`, all 62 launcher controls, Make and
all production-script contracts, all four host acceptance flows, the repaired
15-control audit meta-suite, and empty-HOME count-2. Unchanged launcher tests
intermittently hit the documented partial-raw-log signal-fixture race during
two broader invocations; unchanged standalone and complete Make reruns passed
all controls.

Product commit: `c382caa`.

The authoritative clean full audit from `c382caa` has structured scorecard
SHA-256
`b4c7ddc7d163a833663f8b7cd876e94aadf4c1a875b945252a46441e354c6b72`;
the measured commit tree is `6e44c014c73fb8c99c6c90f20a64a41f314fce43`.
It exits 1 for 13 documented findings, never 2, with L0 8 of 8, six
improved ratchets, two held, zero regressed, zero non-comparable, and zero dirty
paths. The focused eight-criterion scorecard SHA-256 is
`b63afac3b5917dd45d74e19ace4fbf588fe3c871356db4abf1d0c428dea6f703`.
The evidence review after move 10 measured the clean continuity commit
`18beee0a806bcc962cb2f7fbdd5952e43a2cb633`, whose parent is exact product
commit `c382caa`. Its commit tree is
`85a6c67f38b8f764b849dcfac4da9eec071d22d0`; the clean status digest is
`6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`,
inventory SHA-256 is
`4cec690f46b9595d70bce9a08c164aa56d46ac5bb3d69bfe84f83d9006c06e8d`,
and the complete instrument identity is recorded in the rolling handover.
Initial and final clean no-evidence scorecards are byte-identical with
SHA-256
`37ef9bd20cfddf450f381bbc1afe5bd8ea5a037f5a8b0f50713572d22e613f2b`.
They exit 1 for the same 13 documented findings, never 2, with L0 8 of 8,
L1 six PASS and three UNMEASURABLE, six improved ratchets, two held, zero
regressed, zero non-comparable, and zero dirty paths.

The complete Q1.6 production population contains 57 dependency structs with
74 fields: five adapter structs and 52 private composition structs. Read-only
structure review found safe defaults, argument-recording contracts, and
whole-struct delivery for that population; 245 focused top-level dependency
and recording contracts pass. No Q1.6 receipt was emitted because the evidence
mission required all three manual populations to be valid before creating the
external schema-2 document.

Q1.7 cannot truthfully receive a receipt. The explicit scope includes
collection-item failure paths that log and continue, plus operations that
return a populated partial aggregate with an error; it excludes fail-fast
pre-population errors, intentionally ignored errors without a reported partial
result, non-error business warnings, sequential non-collection phases, and
test support. Seventeen operations are in scope. Nine have exact failure-
content coverage, but eight do not: `cmd`'s `initCmd` project loop;
`Bitbucket.SynchronizeAllRepos`' clone/pull warning branch;
`spring.DeleteDemoFiles`; `Repository.upgradeDependencies`;
`Repository.upgradePluginsOnModel`; `maven.RemoveDeprecated`;
`template.MergeTemplates`; and `GitCloudConfig.ValidTemplatesFrom`. The
existing 23 focused top-level partial-failure contracts pass, but they are a
covered subset and cannot establish the missing eight operations.

Q1.9 also cannot truthfully receive a receipt. All 169 syntactic Go test
`range` statements were enumerated. The complete local-table subpopulation has
64 loops over an identifier named `tests`: only 13 functions contain an
executable `len(tests)` guard and 51 can false-green when that table is empty.
Additional non-table expectation and helper populations have the same gap.
The existing 69 named empty-population contracts pass, but they do not cover
the failed table population. The exact external review reports have SHA-256
`c4978c74184f94a8369f9f3e4b8c1a48e90db15ab9b6db24749dce45f7e72c25`
for all 169 ranges and
`d832e212662a36c472b1f8cac2d1d8f0e57f87e98dbb4551ab20352c1e521c1c`
for the 64-table classification.

No manual-evidence document or receipt digest was created, and no focused or
full manual-evidence audit was run. P4 therefore remains active, Q1.6, Q1.7,
and Q1.9 remain UNMEASURABLE, and P5-P8 remain queued. The next coherent move
repairs only Q1.9's executable empty-population assertions, starting with all
51 confirmed table blockers and completing classification of the remaining
test iterations before any receipt is reconsidered.

Move 11 completes that Q1.9-only test repair. Commit `3a8f0bc` changes 40
`_test.go` files and no production file, adding executable fail-on-empty
assertions without changing test data, order, subtest names, side effects, or
public behavior. A fresh AST inventory still contains exactly 169 syntactic
test ranges. Executable-control-flow review classifies all 169 as 72 range
sites that required and now have a new assertion, 83 that were already
independently guarded, and 14 fixture/support sites whose emptiness is valid or
whose result population is rejected elsewhere. All 64 loops over a local
identifier named `tests` now have a direct executable non-empty guard; the 51
previous blockers and 21 additional expectation, recording, and helper range
sites are repaired. The final inventory SHA-256 is
`05b2fc916c1c9d94a44322fcbb997d3e142a9ee70fe11b40de0b43e84a88b343`;
the complete path/function/collection/classification/guard record SHA-256 is
`5f8cb9f00dec919769fde376ae240ae9d5a281063f2a48508d7439f9786eb005`.

The exact affected packages, all 69 named empty-population contracts, complete
uncached tests, race, vet, API/CLI and subprocess compatibility, pinned lint,
launcher and Make contracts, complete preflight, all four host acceptance
flows, the 15-control audit meta-suite, and empty-HOME count-2 pass. The clean
no-evidence audit measures commit
`3a8f0bcd030baf787a29440ee8e4e4a087edb33b`, tree
`75dee3dcc1cb2fd91f7eb359e96b31cc4be3d986`, and zero dirty paths. Its
scorecard SHA-256 is
`08e5ff9d450a6fa3b8824d82c720d339cf0d44003a490d01d3a412681081a615`;
it exits 1 for the same 13 documented non-passing criteria, never 2, with 453
tests across all 27 packages, L0 8 of 8, L1 six PASS and Q1.6/Q1.7/Q1.9
UNMEASURABLE, six improved ratchets, two held, zero regressed, and zero
non-comparable. Q1.9 is independently supported by the complete classified
population, but remains formally UNMEASURABLE because no manual-evidence
receipt was authorized while Q1.7 is blocked. P4 therefore remains active for
the eight previously classified Q1.7 gaps; P5-P8 remain queued.

Move 12 completes every Q1.7 contract exercisable through the production
control points that existed at its start. Test-only commit `7be93e7` changes
seven `_test.go` files and no production file. Ten new top-level contracts
bring the suite from 453 to 463 tests. They assert the exact ordered
failure content and continuation or populated partial result for `initCmd`;
`Bitbucket.SynchronizeAllRepos` repository-query, clone, and pull failures;
the deletions in `spring.DeleteDemoFiles`; dependency version,
maximum-version, parse, and metadata failures; plugin version, parse,
metadata, and release failures; both `maven.RemoveDeprecated` partial/error
and replacement-warning paths; `template.MergeTemplates`; and
`GitCloudConfig.ValidTemplatesFrom`.

Fresh executable-flow classification keeps the complete population at 17
operations. Sixteen operations now have exact content plus continuation or
ordered-partial-result coverage. One included branch remains blocked:
`DeleteDemoFiles`' discovery warning. That public function hard-codes
`file.FindFirst`, whose system `filepath.Walk` callback ignores every incoming
walk error and returns only `nil` or `io.EOF`; final `io.EOF` is normalized to
nil. Therefore no current test seam can make the checked `err` non-nil. This
mission forbade the production injection seam required to exercise the branch,
so no Q1.7 PASS claim or manual receipt was made. The complete 17-row record
SHA-256 is
`4b247f3343a651374b96245d6c15b563f573a152799f3329713c578ace863e57`;
the exact 33-contract focused population manifest SHA-256 is
`43391020bb74aa696a750c68cbdefad40ceee4de5928c4f547356f581ddcae22`.

The exact affected packages, complete 33-contract Q1.7 population, complete
uncached tests, race, vet, API/CLI and subprocess compatibility, pinned lint,
launcher and Make contracts, complete preflight, all four host acceptance
flows, the 15-control audit meta-suite, and empty-HOME count-2 pass. The known
unchanged launcher partial-raw-log signal fixture missed during initial
broader runs; unchanged standalone, complete `make test`, and complete
preflight reruns passed all 62 controls. The clean no-evidence audit measures
commit `7be93e7e822f05bdb007e4e5a1a00402e620a6cb`, tree
`0f0955f8bbe3e49f5c8d1b4f98454dcb04ea62f1`, and zero dirty paths. Its
scorecard SHA-256 is
`b2ccc93dbb045f06404f955ac4b122a614a654f2c0d983d48941fe7026b352ec`;
it exits 1 for the same 13 documented non-passing criteria, never 2, with 463
tests across all 27 packages, L0 8 of 8, L1 six PASS and Q1.6/Q1.7/Q1.9
UNMEASURABLE, six improved ratchets, two held, zero regressed, and zero
non-comparable. P4 remains active for the one Q1.7 controllability repair and
later combined manual evidence; P5-P8 remain queued.

Move 13 adds the single private controllability seam authorized for the last
Q1.7 branch. Public `spring.DeleteDemoFiles(string,
config.ProjectConfiguration)` keeps its exact signature and delegates to an
unexported helper that accepts only the discovery function; production passes
`file.FindFirst`. Suffix selection, discovery arguments, warning spelling,
fixed demo-file order, deletion behavior, and all other side effects remain
unchanged. No dependency struct, deletion injection, exported API, audit,
scanner, parser, inventory, baseline, acceptance, mutation-harness, or manual
evidence change was made.

Focused implementation commit `4f1f45f` adds one exact discovery-failure
contract. It injects a sentinel error, proves the exact `.kt` suffix and
`src/test/kotlin` lookup path, asserts warning
`Unable to find testfile, fileSuffix=.kt`, and proves continuation deletes
`HELP.md`, `mvnw`, and `mvnw.cmd`. The suite increases from 463 to 464
top-level tests. Fresh executable-control-flow review keeps the complete Q1.7
population at 17 operations and now supports every included failure branch
with exact content plus continuation or exact ordered partial-result coverage.
The final 17-row record SHA-256 is
`01cb1abc6293f780fd0caffd6a37da7398ea6d49e9b6b7589ae15ff23d83b9eb`;
the exact 34-contract population manifest SHA-256 is
`0772d019160fe45a8eeefce5fcadba7821fc18da896374d19cb2b48a1dbb05e3`.

Focused Spring and complete 34-contract Q1.7 tests, complete uncached tests,
race, vet, API/CLI and CLI-surface compatibility, fresh subprocess contracts,
pinned lint, launcher and Make contracts, complete preflight, all four host
acceptance flows, the 15-control audit meta-suite, and empty-HOME count-2 pass.
The clean no-evidence audit measures commit
`4f1f45f3f758f2f09dd7d20967efe8ef74a0c613`, tree
`7326302c171604dc315e2f79cd6871c55935aa60`, and zero dirty paths. Its
scorecard SHA-256 is
`ce4190dc03b7bc37aa285569f737b2a59a38c15c5e99db125e340c64e6301b5e`;
it exits 1 for the same 13 documented non-passing criteria, never 2, with 464
tests across all 27 packages, L0 8 of 8, L1 six PASS and Q1.6/Q1.7/Q1.9
UNMEASURABLE, six improved ratchets, two held, zero regressed, and zero
non-comparable. Q1.6, Q1.7, and Q1.9 are now independently supported by their
complete classified populations but remain formally UNMEASURABLE because this
move did not authorize schema-2 evidence. P4 remains active for the combined
manual-evidence checkpoint; P5-P8 remain queued.

Move 14 completes the combined schema-2 manual-evidence checkpoint from clean
continuity commit `88a95ad6effe7c2198d9505963dc19022bcabc10`, whose exact
parent is `4f1f45f3f758f2f09dd7d20967efe8ef74a0c613`. The measured
commit tree is `120d8db0b83a7724b0026f7c18a77a706e9f1cda`; the clean
status SHA-256 is
`6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`;
the inventory SHA-256 remains
`4cec690f46b9595d70bce9a08c164aa56d46ac5bb3d69bfe84f83d9006c06e8d`.
Every helper, record, report, cache, and the evidence document remained outside
the worktree.

Fresh Go-structure review confirms 57 production dependency structs and 74
named fields for Q1.6. All 57 receipt subjects have a default double and
argument recorder for every field and whole-struct delivery; the 245 exact
focused contracts run and pass. The fresh struct inventory, review record, and
contract manifest have SHA-256 values
`38ff8cb2ad78f9ebf41fcfc5aeba2f9c5ddcf6a7dcf467f14f5fc1886a6345da`,
`2dff55671c9b545a28de2e7ff7bcb614eb3064d01b2436ce58a20ce9a875f45b`,
and `6e045b462101ae8feab21ee6d6e117cb308fdaaec7a6615e70de215588914955`.
The Q1.6 receipt SHA-256 is
`85db15d942b850f99baa518428043533ceee7b94d911b2de94b6e45a5f3f0ad7`.

Fresh executable-flow review keeps Q1.7 at 17 production operations and
enumerates 42 distinct included failure branches. Every branch has exact
failure content plus continuation or an exact ordered partial result, with no
exit-only assertion. All 34 exact named contracts run and pass. The
classification and manifest retain SHA-256 values
`01cb1abc6293f780fd0caffd6a37da7398ea6d49e9b6b7589ae15ff23d83b9eb`
and `0772d019160fe45a8eeefce5fcadba7821fc18da896374d19cb2b48a1dbb05e3`.
The Q1.7 receipt SHA-256 is
`e94221777b1f8ea0300cdc286fbb7086440d1038b5b9e9da0f7d845d18845e2d`.

Fresh AST enumeration corrects the inherited Q1.9 total: current HEAD has 183,
not 169, syntactic test `range` sites. The 14 additional sites came from later
Q1.7 tests. The complete current classification contains 168
assertion-required ranges and 15 fixture-construction, delivery, or state-copy
support exclusions; all required ranges have executable empty-population
failures. The 69 exact named guard contracts run and pass, and the complete
34-contract Q1.7 run covers all later sites. The AST inventory, complete
classification, required-population record, and named-guard manifest have
SHA-256 values
`417d9add0e20ae9791a5d5a106ccb5dc1accdccacc9d99355092b25a2cf10787`,
`8326a5ef3037f9cdd72e16a470ead8c2d4e4de96e960215af8002dce5ccc5c26`,
`107d5953139a2959a0777047183893e894ff5743480979c04e2cb6e9bd5d08b7`,
and `bfd7210785bc672192575e372b098ebb1cbf226b5793d0ac8c15065be0e6295b`.
The Q1.9 receipt SHA-256 is
`6d4bf0c4429fa7a731078b84bc850d401a1c636f4cf72bb8d110dfe4700d9616`.

The external schema-2 document contains only the three authorized criterion
receipts and has SHA-256
`c97c0336998508fe410fc3f56b65cc41762281850ba5822d53c623dc74be34c0`.
The clean no-evidence, focused manual-evidence, and full authoritative
scorecards have SHA-256 values
`0d038c41c160c8851b405e32f69bb9dcf03729749fca5f423ef18435b2cc58e2`,
`63ede0de00d52d950f10845ee8d6e16c1e948e73c8dd183c368ef0e87cc4c729`,
and `52490a43d26c5b6c1e6031460352a9fbc6de31bba52128301d3d3e1f31d86bc6`.
The focused audit exits 0. The full audit exits 1 for the nine documented
P5-P8 findings, never 2, and records all nine L1 rows PASS, L0 8 of 8, six
improved ratchets, two held, zero regressed, zero non-comparable, and zero
dirty paths. Every non-manual criterion object is identical to the no-evidence
audit.

API/CLI and CLI-surface compatibility, fresh entry/subprocess contracts,
pinned golangci-lint 2.12.2, complete tests across 27 packages, race, vet,
launcher and Make contracts, all four host acceptance flows, exact empty-HOME
count-2, the standalone 15-control audit meta-suite, and complete preflight all
pass. Two launcher invocations missed the known signal-fixture partial-log
timing assertion, once under concurrent heavy gates and once in the final
source-archive nested generation; unchanged isolated reruns and complete
preflight pass all 62 controls. Ordinary and ignored status are empty after the
gate. P4 is complete; P5 is active for the first bounded `cli-context`
mutation-harness move, while P6-P8 remain queued.

### P5 - Build L2 Mutation Evidence

Status: complete.

P5.1 clean checkpoint (first subject `cli-context`, implementation `1dbc163`):

- Added the regular executable `scripts/mutate-cli-context` and its executable
  `scripts/test-mutate-cli-context` falsifiability meta-test. The inventory is
  unchanged; exact executable subject coverage improves from 0 of 8 to 1 of 8.
- The harness declares ten unique, deterministic production mutations across
  `cmd` and `pkg/context`: CLI exit status, interactive build endpoint,
  interactive upgrade endpoint, recursive discovery branch, flattened-POM
  exclusion, target-directory exclusion, stealth selection, per-project
  dry-run gate, root-project selection, and configured Maven-repository
  selection. Every search syntax occurs exactly once.
- One clean external control ran each exact selection. Every mutant used a
  fresh external Git archive, cache, HOME, config, and temp root; compiled
  separately; and ran its non-empty exact named test population under
  `go test -json`. A mutation counted as killed only when every selected test
  emitted a run and terminal action and at least one selected test emitted a
  failure. Final totals are `declared=10`, `killed=10`, `survived=0`, and
  `unusable=0`; no production or test repair was required.
- T1-T10 fail closed on empty and duplicate manifests, unauthorized paths,
  zero and multiple replacements, empty/inexact test selection, a broken
  clean control, an uncompiled or unexercised mutant, falsified accounting, an
  unclassified survivor, repository-local artifacts, and non-deterministic
  manifest/totals. The production run report SHA-256 is
  `53d1489a6db15f18cd2acf38ff49ab3fb0db545a6c8e63639c90311f9bf53d5c`;
  the T1-T10 meta-log SHA-256 is
  `f1d801fccc97af714d483f5a3c84202fe9868fa19b17038ada0d56dbd6f3e67e`.
- The no-evidence Q2.1-Q2.4 structured view records the exact Q2.1 numerator
  1 of 8, Q2.2/Q2.3 debt only in the remaining P3/P4 seam drivers, and Q2.4
  `UNMEASURABLE` until supplied the run receipt; its scorecard SHA-256 is
  `7500813f38b080a06494b62d70a3def51ddd35bf75d6f804e382ddb2133c4545`.
  The current external schema-2 document adds only the truthful `cli-context`
  Q2.4 receipt to refreshed Q1.6/Q1.7/Q1.9 receipts. Its SHA-256 is
  `0a6659c3b5a5234ca19134b4e2305de43a094fcb02032f0c6e106f5bbecaabbd`;
  Q2.4 evidence-object SHA-256 is
  `70e2101ca2719d07022dd42e8db56b9ffc76c444473460907a216b9f6f4b124f`.
  The focused four-row audit exits 0 with scorecard SHA-256
  `c693cc912f8c936c4438661aae48102ec09d05fe1ec4a78ad52f7c3d40abf7b8`.
- The full authoritative audit exits 1, never 2, with scorecard SHA-256
  `20528637cc01f651f6411484575bf3cec7f1a1ac6cccc9f1b30cdf9b989838ac`:
  L0 is 8 of 8, L1 is 9 of 9, Q2.1 improves from 0 of 8 to 1 of 8,
  Q2.4 passes, seven ratchets improve, one holds, none regress, dirty paths
  are empty, and eight P5-P8 rows remain non-passing.
- API/CLI and subprocess compatibility, pinned golangci-lint 2.12.2, complete
  tests across all 27 packages, race, vet, launcher and Make contracts, all
  four host acceptance flows, exact empty-HOME count-2, the standalone
  15-control audit meta-suite, and complete preflight all pass. API and CLI
  report SHA-256 values remain
  `ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
  and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

P5.2 clean checkpoint (second subject `config-cloud`, implementation
`e5b4a26`):

- Replaced the non-executable P3 `config-cloud` seam driver with the regular
  executable `scripts/mutate-config-cloud` harness and replaced its meta-test
  with an executable T1-T10 falsifiability test. The inventory and exact seam
  label `3. cloud clone keeps URL before target directory` are unchanged;
  exact executable subject coverage improves from 1 of 8 to 2 of 8.
- The harness declares ten unique deterministic mutations in
  `pkg/config/cloud.go`: clone argument ordering, existing-repository branch,
  pull and clone error gates, Git-hook file classification, template config
  matching and relative naming, example directory classification, default
  service-environment selection, and valid-template deduplication. Every
  production search syntax occurs exactly once and every declaration binds an
  exact non-empty named `pkg/config` test population.
- One clean external control ran every exact selection. Every mutant used a
  fresh external Git archive, cache, HOME, config, and temp root; compiled
  separately; and ran its selected JSON test population with exact run and
  terminal-action validation. Final totals are `declared=10`, `killed=10`,
  `survived=0`, and `unusable=0`; no production or test repair was required.
- T1-T10 fail closed on empty/duplicate manifests, unauthorized paths,
  zero/multiple replacements, empty or inexact test selection, a broken clean
  control, an uncompiled/unexercised mutant, false accounting, an unclassified
  survivor, repository-local artifacts, and non-deterministic manifests or
  totals. The run report SHA-256 is
  `db76fdf8c624c4326483ae71fa9ec3e7e0f94de4d6d3ae8a4dff185c9a7f23f7`;
  the independent T1-T10 meta-log SHA-256 is
  `12b1197d521681b70dea0b481f6b7d8bb9daba2ca21dc68a635868177d43ffec`.
- The no-evidence Q2.1-Q2.4 view records exact Q2.1 coverage 2 of 8 and has
  scorecard SHA-256
  `5a17bd88fc7a48065883d72b1c6383b037fc2b47c18e47390918ba7cd782ed24`.
  The refreshed external schema-2 document covers both completed P5 subjects
  and has SHA-256
  `7665b6d6d15da8dbd6b2da3e20a45902cb4e1076683fb0e1acd52c17b85c4648`;
  its Q2.4 evidence-object SHA-256 is
  `496ae9cd500b896080a14874dd5ff34bc0bdf521a401cbb7044e5f87604184b5`.
  The focused audit exits 0 with scorecard SHA-256
  `7ef785e2a17b5371358cc18741a78fe7cbca65c587725970970710d42a7a6b0c`.
- The full authoritative audit exits 1, never 2, with scorecard SHA-256
  `4a7089ea117a97bdf5b265f3712c534954417e00ff352b24a5874a2d710bf648`:
  L0 is 8 of 8, L1 is 9 of 9, Q2.1 is 2 of 8, Q2.4 passes, seven
  ratchets improve, one holds, none regress, and dirty paths are empty. Eight
  P5-P8 rows remain non-passing.
- API/CLI and subprocess compatibility, pinned golangci-lint 2.12.2, complete
  tests across all 27 packages, race, vet, `make test`, all 62 launcher
  controls, Make contracts, all four host acceptance flows, exact empty-HOME
  count-2, the standalone 15-control audit meta-suite, and complete preflight
  all pass. API and CLI report SHA-256 values remain
  `ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
  and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

P5.3 clean checkpoint (third subject `maven-sorting`, implementation
`1a40e11`):

- Replaced the non-executable P3 `maven-sorting` seam driver with the regular
  executable `scripts/mutate-maven-sorting` harness and replaced its meta-test
  with an executable T1-T10 falsifiability test. The inventory and exact seam
  labels `2. Maven command keeps executable before arguments` and
  `4. Maven metadata keeps username before password` are unchanged; exact
  executable subject coverage improves from 2 of 8 to 3 of 8.
- The harness declares ten unique deterministic mutations across
  `pkg/maven/command.go`, `pkg/maven/query.go`, and `pkg/sorting/sort.go`:
  Maven executable/argument ordering; metadata username/password and coordinate
  path ordering; authentication and error gates; dependency swap/comparison
  directions; sort-key field ordering; group-match gating; and compile-scope
  priority. Every production search syntax occurs exactly once and every
  declaration binds an exact non-empty named test population.
- One clean external control ran every exact selection. Every mutant used a
  fresh external Git archive, cache, HOME, config, and temp root; compiled its
  changed package separately; and ran its selected JSON test population with
  exact run and terminal-action validation. Final totals are `declared=10`,
  `killed=10`, `survived=0`, and `unusable=0`; no production or test repair was
  required.
- T1-T10 fail closed on empty/duplicate manifests, unauthorized paths,
  zero/multiple replacements, empty or inexact test selection, a broken clean
  control, an uncompiled/unexercised mutant, false accounting, an unclassified
  survivor, repository-local artifacts, and non-deterministic declarations or
  totals. The run report SHA-256 is
  `6b360f604802c047b4946b474c35b2860d47bda4581b3cb33f5af45653dc111e`;
  the independent T1-T10 meta-log SHA-256 is
  `fda3185a71ebd842f3e924ae507ad0daa9d59d4b30da3699d114c0ff0d2d53b3`.
- The no-evidence Q2.1-Q2.4 view records exact Q2.1 coverage 3 of 8 and has
  scorecard SHA-256
  `f6e6c4ecaf67549fde6c913f848a1630c79237659a000c540b308b5f8fb8e724`.
  The refreshed external schema-2 document covers all three completed P5
  subjects and has SHA-256
  `a279e8ac3402cacad23eb757bc3adbf7e9e1e105115135f7183ebf5353ff1bb4`;
  its Q2.4 evidence-object SHA-256 is
  `53d27607c73ae97ea00dcfce384a373c5c421fcd1e900433521fbfbb1e4bf689`.
  The focused audit exits 0 with scorecard SHA-256
  `07d9f680bba06cb88c589fe87962feaf74ca6add94fda8b933679d6a98bc313c`.
- The full authoritative audit exits 1, never 2, with scorecard SHA-256
  `122deef2c884c1e3f6edb3052b20d47c0fad82ba5a4dc44f0b3696afe947d9f4`:
  L0 is 8 of 8, L1 is 9 of 9, Q2.1 is 3 of 8, Q2.4 passes, seven
  ratchets improve, one holds, none regress, and dirty paths are empty. Eight
  P5-P8 rows remain non-passing.
- API/CLI and subprocess compatibility, pinned golangci-lint 2.12.2, complete
  tests across all 27 packages, race, vet, `make test`, all 62 launcher
  controls, Make contracts, all four host acceptance flows, exact empty-HOME
  count-2, the standalone 15-control audit meta-suite, and complete preflight
  all pass. API and CLI report SHA-256 values remain
  `ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
  and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

P5.4 clean checkpoint (fourth subject `template`, implementation `da1cc88`):

- Replaced the non-executable P3 `template` seam driver with the regular
  executable `scripts/mutate-template` harness and replaced its meta-test with
  an executable T1-T10 falsifiability test. The inventory and exact seam label
  `6. template copy keeps source before destination` are unchanged; exact
  executable subject coverage improves from 3 of 8 to 4 of 8.
- The harness declares ten unique deterministic mutations in
  `pkg/template/template.go`: source-root selection, target path component
  ordering, copy argument ordering, copy error gating, walk directory
  classification, root-POM detection, nested-POM exemption, render-file
  exemption, returned file ordering, and markdown write bytes. Every production
  search syntax occurs exactly once and every declaration binds an exact
  non-empty named `pkg/template` test population.
- One clean external control ran every exact selection. Every mutant used a
  fresh external Git archive, cache, HOME, config, and temp root; compiled the
  changed package separately; and ran its selected JSON test population with
  exact run and terminal-action validation. Final totals are `declared=10`,
  `killed=10`, `survived=0`, and `unusable=0`; no production or test repair was
  required.
- T1-T10 fail closed on empty/duplicate manifests, unauthorized paths,
  zero/multiple replacements, empty or inexact test selection, a broken clean
  control, an uncompiled/unexercised mutant, false accounting, an unclassified
  survivor, repository-local artifacts, and non-deterministic declarations or
  totals. The run report SHA-256 is
  `e769cb49b2ae6fd3da15b207acab8ef11c402f5bff171dc6e150623b93279e95`;
  the independent T1-T10 meta-log SHA-256 is
  `cde7a149274b7f4b968a40b1ec5e4f7769c5d6d4243edcb40783a4fa19ba91bd`.
- The no-evidence Q2.1-Q2.4 view records exact Q2.1 coverage 4 of 8 and has
  scorecard SHA-256
  `1f29a6fdd16484ec3c4741772dde20ce3b076025532db90f1231506f7ee137d3`.
  The refreshed external schema-2 document covers all four completed P5
  subjects and has SHA-256
  `f235776be8c9fc454906415a1f983aeebf6ef29adbd4a8665a83e3bed76fe050`;
  its Q2.4 evidence-object SHA-256 is
  `bd2899300363740c159c07b41dfa76201371ec5e74b29ac545e20033ebb0a433`.
  The focused audit exits 0 with scorecard SHA-256
  `28732aa8d367b86cfb74a3e25d42a9909d01e8fc4564d4488d720b04675c4e8f`.
- The full authoritative audit exits 1, never 2, with scorecard SHA-256
  `fcfddf9b5794dabb1d2bdbbc19b04d594522bf7efbd7db7432a35882a6ee797c`:
  L0 is 8 of 8, L1 is 9 of 9, Q2.1 is 4 of 8, Q2.4 passes, seven
  ratchets improve, one holds, none regress, and dirty paths are empty. Eight
  P5-P8 rows remain non-passing.
- API/CLI and entry/subprocess compatibility, pinned golangci-lint 2.12.2,
  complete tests across all 27 packages, race, vet, `make test`, all 62
  launcher controls, Make contracts, all four host acceptance flows, exact
  empty-HOME count-2, the standalone 15-control audit meta-suite, and complete
  preflight all pass. API and CLI report SHA-256 values remain
  `ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
  and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

P5.5 clean checkpoint (fifth subject `file-shell`, implementation `7bb94ac`):

- Replaced the non-executable P3 `file-shell` seam driver with the regular
  executable `scripts/mutate-file-shell` harness and replaced its meta-test
  with an executable T1-T10 falsifiability test. The inventory and exact seam
  labels `1. git clone keeps URL before target directory` and
  `7. git commit keeps target directory before message` are unchanged; exact
  executable subject coverage improves from 4 of 8 to 5 of 8.
- The harness declares ten unique deterministic mutations across
  `pkg/file/file.go`, `pkg/shell/command.go`, and `pkg/shell/git.go`: clone,
  pull, commit, command, unzip, and move argument ordering; the Git-add error
  gate; first-match suffix selection; all-match exclusion selection; and
  clear-directory exclusion gating. Every production search syntax occurs
  exactly once and every declaration binds an exact non-empty named
  `pkg/file` or `pkg/shell` test population.
- One clean external control ran every exact selection. Every mutant used a
  fresh external Git archive, cache, HOME, config, and temp root; compiled its
  changed package separately; and ran its selected JSON test population with
  exact run and terminal-action validation. Final totals are `declared=10`,
  `killed=10`, `survived=0`, and `unusable=0`; no production or test repair was
  required.
- T1-T10 fail closed on empty/duplicate manifests, unauthorized paths,
  zero/multiple replacements, empty or inexact test selection, a broken clean
  control, an uncompiled/unexercised mutant, false accounting, an unclassified
  survivor, repository-local artifacts, and non-deterministic declarations or
  totals. The run report SHA-256 is
  `ea12b122e6210f2b2c8c64a7792ce74a9886e2273e9ddac6176df6552afd1c29`;
  the independent T1-T10 meta-log SHA-256 is
  `759ad7a28b3d3ab06d42d42e8e1679e0df516d6abdfd90f3fbac1e92e52bfaad`.
- The no-evidence Q2.1-Q2.4 view records exact Q2.1 coverage 5 of 8 and has
  scorecard SHA-256
  `fcbb44afd1ec5a92ce8d4cd6dd5b8b6a12de363fbc96b33dab07e345c87a4338`.
  The refreshed external schema-2 document covers all five completed P5
  subjects and has SHA-256
  `c5ec79633274ad618b7c78568e38891b9955352e020053c8b3d987ad6348fd5e`;
  its Q2.4 evidence-object SHA-256 is
  `732083b2add74a30bf5bed70ef9446682bbba20ebb4d1dfc89d0b904db81b922`.
  The focused audit exits 0 with scorecard SHA-256
  `5261b60d5e20338c5f1a88870b60c55e54411200a7fa84a0d80791dd96597001`.
- The full authoritative audit exits 1, never 2, with scorecard SHA-256
  `7a42d1f40ec96641fce8dab1188a32a73a2f6f790cedea60e98e5b05b1833502`:
  L0 is 8 of 8, L1 is 9 of 9, Q2.1 is 5 of 8, Q2.4 passes, seven
  ratchets improve, one holds, none regress, and dirty paths are empty. Eight
  P5-P8 rows remain non-passing.
- API/CLI and entry/subprocess compatibility, pinned golangci-lint 2.12.2,
  complete tests across all 27 packages, race, vet, `make test`, all 62
  launcher controls, Make contracts, all four host acceptance flows, exact
  empty-HOME count-2, the standalone 15-control audit meta-suite, and complete
  preflight all pass. API and CLI report SHA-256 values remain
  `ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
  and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

P5.6 clean checkpoint (sixth subject `spring`, implementation `9c20e52`):

- Replaced the non-executable P3 `spring` seam driver with the regular
  executable `scripts/mutate-spring` harness and replaced its meta-test with an
  executable T1-T10 falsifiability test. The inventory and exact seam label
  `5. Spring download keeps URL before archive path` are unchanged; exact
  executable subject coverage improves from 5 of 8 to 6 of 8.
- The harness declares ten unique deterministic mutations in
  `pkg/spring/io.go`: download and unzip argument ordering; download and unzip
  error gates; archive deletion selection; injected Unix-second ordering; root
  and dependency discovery URLs; empty-dependency validation; and valid
  dependency equality. Every production search syntax occurs exactly once and
  every declaration binds an exact non-empty named `pkg/spring` test
  population.
- One clean external control ran every exact selection. Every mutant used a
  fresh external Git archive, cache, HOME, config, and temp root; compiled the
  changed package separately; and ran its selected JSON test population with
  exact run and terminal-action validation. Final totals are `declared=10`,
  `killed=10`, `survived=0`, and `unusable=0`; no production or test repair was
  required.
- T1-T10 fail closed on empty/duplicate manifests, unauthorized paths,
  zero/multiple replacements, empty or inexact test selection, a broken clean
  control, an uncompiled/unexercised mutant, false accounting, an unclassified
  survivor, repository-local artifacts, and non-deterministic declarations or
  totals. The run report SHA-256 is
  `0aa2bdcb32cf20372f8a7b5c232dcefc93668259eb4ed76c45f8ca339b1dbf25`;
  the independent T1-T10 meta-log SHA-256 is
  `f373190fa04abba73b61453b1693f892aced845d931ce267b3674ea84961fa06`.
- The no-evidence Q2.1-Q2.4 view records exact Q2.1 coverage 6 of 8 and has
  scorecard SHA-256
  `a05b6401db516d8dc09e213b6545aa080a58fdb48ca710de29a2d1dff9e49d1f`.
  The refreshed external schema-2 document covers all six completed P5
  subjects and has SHA-256
  `1eaabad8e769ad3bdc8eedfcf145aa521cb982da1d64697c06535df1e4b5ff14`;
  its Q2.4 evidence-object SHA-256 is
  `61731b25435e57c08c0c813d34d41cfce24644defabb13483fbc5cf88ee1194c`.
  The focused audit exits 0 with scorecard SHA-256
  `ab347e874db50c8df2682c30ad86f95f47c109bed71342dc131e90c3bb645a0c`.
- The full authoritative audit exits 1, never 2, with scorecard SHA-256
  `50ea67476a0e01cf7e31fa58e8a6c4590934d6691939b5c304e835a8af6aaf67`:
  L0 is 8 of 8, L1 is 9 of 9, Q2.1 is 6 of 8, Q2.4 passes, seven
  ratchets improve, one holds, none regress, and dirty paths are empty. Eight
  P5-P8 rows remain non-passing.
- API/CLI and entry/subprocess compatibility, pinned golangci-lint 2.12.2,
  complete tests across all 27 packages, race, vet, `make test`, all 62
  launcher controls, Make contracts, all four host acceptance flows, exact
  empty-HOME count-2, the standalone 15-control audit meta-suite, and complete
  preflight all pass. API and CLI report SHA-256 values remain
  `ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
  and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

P5.7 clean implementation checkpoint (seventh subject `http`, implementation
`5829939`):

- Added the regular executable `scripts/mutate-http` and its executable
  `scripts/test-mutate-http` T1-T10 falsifiability test. The inventory remains
  unchanged; `http` has no separate seam label. Exact executable subject
  coverage improves from 6 of 8 to 7 of 8.
- The harness declares ten deterministic, unique mutations in
  `pkg/http/client.go`: anonymous JSON request URL; response dependency and
  status gates; response-body source; basic-auth credential order; bearer URL
  component order and token selection; Wget request URL and destination path;
  and Wpost request URL. Every production search syntax occurs exactly once
  and every declaration binds an exact non-empty named `pkg/http` test.
- One clean external control ran every exact selection. Every mutant used a
  fresh external Git archive, cache, HOME, config, and temp root; compiled the
  changed package separately; and ran its selected JSON test population with
  exact run and terminal-action validation. Final totals are `declared=10`,
  `killed=10`, `survived=0`, and `unusable=0`; no production or test repair was
  required.
- T1-T10 fail closed on empty/duplicate manifests, unauthorized paths,
  zero/multiple replacements, empty or inexact test selection, a broken clean
  control, an uncompiled/unexercised mutant, false accounting, an unclassified
  survivor, repository-local artifacts, and non-deterministic declarations or
  totals. The run report SHA-256 is
  `865d29bee9baad2053e603e188d8ccdd5688ecbcbe053979f586d535898645c7`;
  the independent T1-T10 meta-log SHA-256 is
  `6b932bf00ed9ae0f14152b821ee5cfcd4636c09c739fe0b4d21ef7974b900270`.
- The no-evidence Q2.1-Q2.4 view records exact Q2.1 coverage 7 of 8 and has
  scorecard SHA-256
  `46a5070d9538576053e5b8889fac01b87f90b1bcfa9e5bc9326c7ab334a98ed2`.
  Q2.2 and Q2.3 now identify only the legacy `interactive-build` recorder.
  The refreshed external schema-2 document covers all seven completed P5
  subjects and has SHA-256
  `c4f5d9c8c97aa653afd1b4ef8095f1e1bd35873bc40acc5ab9e6d9f308544713`;
  its Q2.4 evidence-object SHA-256 is
  `6a462402ad3de49c3f14b3907664e04edd39af7b360421b6497a170a1a620de8`.
  The focused audit exits 0 with scorecard SHA-256
  `a9034e35989bf4d1508edd5deda4b59fc4ffc34172db12052996a27c39d56383`.
- The full authoritative audit exits 1, never 2, with scorecard SHA-256
  `c3d8cebe64fd4aebc10596852ec9ece963937bc079858e9fab0300a41d8fb968`:
  L0 is 8 of 8, all nine L1 rows pass, Q2.1 is 7 of 8, Q2.4 passes,
  seven ratchets improve, and one Q3.4 documentation-phrase ratchet regresses.
  Its three findings are literal `currently` phrases in the required tracked
  rolling handover and byte-exact current prompt archive; the implementation
  introduced no product or apparatus regression, and immutable prompt history
  cannot be rewritten truthfully. Nine P5-P8 rows are non-passing.
- API/CLI and entry/subprocess compatibility, pinned golangci-lint 2.12.2,
  complete tests across all 27 packages, race, vet, `make test`, all 62
  launcher controls, Make contracts, all four host acceptance flows, exact
  empty-HOME count-2, the standalone 15-control audit meta-suite, and complete
  preflight all pass. The first complete preflight hit the previously recorded
  final signal-fixture timing miss after 25 launcher controls; the unchanged
  rerun passed the entire gate.

P5.8 clean implementation checkpoint (eighth subject `interactive-build`,
implementation `58c5224`):

- Replaced the P3 recording driver with the regular executable
  `scripts/mutate-interactive-build` harness and replaced its meta-test with an
  executable T1-T10 falsifiability test. The inventory and exact seam label
  `8. interactive server binds only to loopback` are unchanged. Exact
  executable subject coverage improves from 7 of 8 to 8 of 8.
- The harness declares ten deterministic, unique mutations across
  `cmd/build.go`, `pkg/webservice/init.go`, and `pkg/webservice/api.go`:
  project configuration, cloud configuration, root response, and generation
  endpoint wiring; interactive URI host and asynchronous server start;
  loopback bind address; generate route; listen-error gate; and Darwin browser
  platform selection. Every production search syntax occurs exactly once and
  every declaration binds an exact non-empty named test population.
- One clean unmodified external control ran every exact selection. Every mutant
  used a fresh external Git archive, cache, HOME, config, and temp root;
  compiled every changed package separately; and ran its selected JSON test
  population with exact run and terminal-action validation. Final totals are
  `declared=10`, `killed=10`, `survived=0`, and `unusable=0`; no production or
  test repair was required.
- T1-T10 fail closed on empty/duplicate manifests, unauthorized paths,
  zero/multiple replacements, empty or inexact test selection, a broken clean
  control, an uncompiled/unexercised mutant, false accounting, an unclassified
  survivor, repository-local artifacts, and non-deterministic declarations or
  totals. The retained run report SHA-256 is
  `109619d5d336a52f01d5f751d3b89b5bcfa85f2ba512cdc5073fbc58e1954d5f`;
  the independent T1-T10 meta-log SHA-256 is
  `22a0ac95820514fec641d72289c6ef1a3a56023cc6471f4c2e154087199c9f81`.
- The no-evidence Q2 view records Q2.1 at 8 of 8 and automated PASS verdicts
  for Q2.1-Q2.3; Q2.4 remains correctly unmeasurable without the external
  receipt. Its scorecard SHA-256 is
  `4154c02e7520e2021cbcb4ffee777b00d1448c6f451717691d5a5bceb9f8c221`.
  The refreshed schema-2 document covers all eight completed subjects and has
  SHA-256
  `b12d4c262ecd16c9529887968cc33324552b08a7cd1d300ff0f1dff6e34acbba`;
  its Q2.4 evidence-object SHA-256 is
  `2565dd9824e224520004774521164f0d24ff771b134f7662125a337ae27dd09d`.
  The focused audit exits 0 with scorecard SHA-256
  `a5478b88f4dbd99dfa5561eafa09117c649a3f3998b79b34818b54b3970d6077`.
- The full authoritative audit exits 1, never 2, with scorecard SHA-256
  `9032b0e3763c64f820564f1608fd6ef52b0a751311c4dcbcd9a2fca8b1eef972`:
  L0 is 8 of 8, L1 is 9 of 9, Q2.1-Q2.4 pass, seven ratchets improve,
  one Q3.4 documentation-phrase ratchet regresses, dirty paths are empty, and
  seven P6-P8 rows remain non-passing.
- API/CLI and entry/subprocess compatibility, pinned golangci-lint 2.12.2,
  complete tests across all 27 packages, race, vet, `make test`, all 62
  launcher controls, Make contracts, all four host acceptance flows, exact
  empty-HOME count-2, the standalone 15-control audit meta-suite, and complete
  preflight all pass. The first complete preflight hit the known launcher
  signal-retention timing control after 25 passes; the unchanged complete
  rerun passed the entire gate.

P5 is complete: all eight named inventory subjects have regular executable
harnesses, at least eight meaningful mutations, T1-T10 controls, and exact
zero-survivor/zero-unusable totals. P6 is the next active phase.

P6 snapshot-binary checkpoint (implementation `06e3ac4`):

- Added the regular executable `scripts/accept-snapshot`, its executable
  20-control falsifiability test, the private
  `test/acceptance/snapshot_acceptance.py` helper, and focused Make/distribution
  wiring. The orchestration requires a clean committed source, a fresh external
  evidence root, exact pinned GoReleaser identity, the established
  credential-cleared `release --snapshot --clean --skip=publish` invocation,
  and metadata-derived discovery of exactly one host executable. No
  `.goreleaser.yml`, production Go, Go test, inventory, audit, parser, scanner,
  baseline, dependency, API, CLI, or completed P5 file changed.
- A real GoReleaser v2.17.1 run at `06e3ac4` produced exactly one
  `darwin/arm64` binary record,
  `dist/plybuild_darwin_arm64_v8.0/ply`. The retained executable is 19,389,170
  bytes with SHA-256
  `28d7012a81e9bdab2b776a8be64f84d12fc9e2f239e52dcbfc2932f7469d2e4f`.
  The report and evidence-manifest SHA-256 values are
  `6e759f75ae40c27c6312cfc4f4c1e8a8c4eab56119a64a88fa4c4c62dec40f61`
  and `bb33ae40ccab0992284318a123c9462fa86ae6a1e5dd2c7b881d0b302c4ab9a1`.
- Host install remained a separate `make install` / `go install ./cmd/ply`
  contract into an initially missing external `GOBIN`; `verify-install` ran
  with its artifact override unset. The exact snapshot path then drove the
  existing status, upgrade, and build verifiers. Each emitted exactly one
  terminal PASS and two non-help behavioral trace records against the stable
  snapshot identity.
- The 20 controls fail closed for no GoReleaser call, wrong or publishing argv,
  leaked credentials, stale output, missing/ambiguous/symlinked/non-executable,
  wrong-platform, misplaced, wrong-build-platform, or foreign artifacts,
  skipped/duplicate verifiers, artifact substitution, verifier failure,
  duplicate terminal PASS, and repository-local output. The focused meta-log
  SHA-256 is
  `9804c5b0d54f5ada0dcfbd1ee401559dc026c8aa93778b2c604afb64adffa9b9`.
- API/CLI and entry/subprocess compatibility, pinned golangci-lint 2.12.2,
  complete uncached tests across 27 packages, race, vet, `make test`, all 62
  launcher controls, Make contracts, all four existing host acceptance flows,
  exact empty-HOME count-2, the standalone 15-control audit meta-suite, and an
  unchanged complete preflight rerun pass. The first preflight hit the known
  launcher signal-retention timing control after 25 assertions; its retained
  log records the classified fixture flake and the isolated unchanged rerun
  passed the full gate.
- The current external schema-2 document validates at `06e3ac4`. Its SHA-256 is
  `7a09ec592025564f6600b3edf21fd0c2d56dc28fd73a92de2b3bb80024f32e70`;
  the exact Q1.6/Q1.7/Q1.9/Q2.4 audit exits 0. The full authoritative audit
  exits 1, never 2, with scorecard SHA-256
  `9b75a367fc07ddc698787103f80360227e6f03827009c62c4a80a59dcd547b4d`:
  L0 is 8/8, L1 is 9/9, L2 is 8/10, seven ratchets improve, one documentation
  ratchet regresses from immutable prompt history, and the same seven rows are
  non-passing. The focused Q2.5-Q2.10 view has scorecard SHA-256
  `367819708c9959d6b6bd6085e7a362124c95a2f9bb744db48a7294a1ffa0d482`:
  Q2.5, Q2.6, Q2.7, and Q2.10 pass while Q2.8 and Q2.9 remain honestly manual.

This is one bounded P6 acceptance target, not P6 exit. Fresh daemon-backed
Docker evidence is next; `make quality` remains deferred until both artifact
populations exist and the scoped L2 gate can be implemented without masking
Q2.8 or Q2.9.

P6 Docker daemon attempt (launch continuity `9fa49c3`):

- The branch, exact ancestry, ordinary and ignored cleanliness, reciprocal
  archive graph, launcher check, and active authorized P6 checkpoint all
  passed before the probe. Required Docker, distribution, snapshot,
  acceptance, audit, compatibility, fixture, and design inputs were read.
- The installed Docker 29.4.0 client and buildx 0.33.0 select the
  `desktop-linux` context on a macOS 15.3.1 arm64 host, but no server answered:
  `docker version` and `docker info` exited 1. Docker Desktop 4.70.0 could not
  be started from the managed sandbox because its required
  `~/Library/Containers/com.docker.docker` log state is outside the writable
  roots; LaunchServices, Apple Events, and UI control were unavailable.
- No alternate local daemon is installed. No Docker build or run occurred, so
  no Dockerfile defect, image identity, platform, executable population, or
  verifier runtime can be claimed. No implementation shape was selected and no
  production, test, Make, Dockerfile, or distribution file changed.
- Exact external blocker evidence is retained at
  `/private/tmp/ply-docker-blocker.KRsWoi`; its evidence-manifest SHA-256 is
  `991d7bf4a7a739b9d971da7776a8be7e447d2aad6c931a7cf72c1bbdaeb682ae`.

Docker acceptance remains the next P6 move. Docker Desktop must be started
outside the managed workspace sandbox and `docker version` must report a real
Server identity before the resume session edits or implements anything.

P6 Docker daemon resume attempt (launch continuity `ddea23b`):

- Exact ancestry, ordinary and ignored cleanliness, reciprocal archive links,
  the launcher check, and the authorized P6 checkpoint passed again.
- The fail-closed probe stopped before Dockerfile inspection or implementation:
  Docker 29.4.0 returned JSON `Server: null`; `docker version`, `docker info`,
  and `docker desktop status` exited 1 for `desktop-linux`. Buildx 0.33.0 could
  name the builder but reported no driver or supported platforms because the
  daemon was unreachable. The managed sandbox again could not open Docker
  Desktop's required host log.
- No Docker build/run, login, push, publication, design selection, image claim,
  or repository implementation change occurred. The 42-file external manifest
  at `/private/tmp/ply-docker-probe.IpwJLG` verifies and has SHA-256
  `7adb0412028e616c4041a527e3a215d514a69f42a65d5f77d7da34e181ed614e`.

Docker acceptance remains blocked on a daemon started outside the managed
sandbox. The next session must repeat the same Client+Server and `docker info`
gate before reading the Dockerfile or selecting an implementation design.

P6 Docker daemon second resume attempt (launch continuity `20dbf71`):

- Exact ancestry, ordinary and ignored cleanliness, reciprocal archive links,
  the launcher check, and the authorized P6 checkpoint passed again.
- The fail-closed probe stopped before Dockerfile inspection or implementation:
  Docker 29.4.0 returned JSON `Server: null`; `docker version`, `docker info`,
  and `docker desktop status` exited 1 for `desktop-linux`. Buildx 0.33.0 named
  the builder but exposed no driver or supported platforms because the daemon
  was unreachable. The managed sandbox again could not open Docker Desktop's
  required host log.
- No Docker build/run, login, push, publication, design selection, image claim,
  or repository implementation change occurred. The 49-file external manifest
  at `/private/tmp/ply-docker-probe.ffXBUb` verifies and has SHA-256
  `16ed6fed1d517f8585815032363e2ba8b2a57fa6f1e811ad4bde4da892ec2a52`.

Docker acceptance remains blocked on a daemon started outside the managed
sandbox. The next session must repeat the same Client+Server and `docker info`
gate before reading the Dockerfile or selecting an implementation design.

P6 Docker daemon third resume attempt (launch continuity `fde001c`):

- Exact ancestry, ordinary and ignored cleanliness, reciprocal archive links,
  the launcher check, and the authorized P6 checkpoint passed again.
- The fail-closed probe stopped before Dockerfile inspection or implementation:
  Docker 29.4.0 returned JSON `Server: null`; `docker version`, `docker info`,
  and `docker desktop status` exited 1 for `desktop-linux`. Buildx 0.33.0 named
  the builder but exposed no driver or supported platforms because the daemon
  was unreachable. The managed sandbox again could not open Docker Desktop's
  required host log.
- No Docker build/run, login, push, publication, design selection, image claim,
  or repository implementation change occurred. The 91-entry external manifest
  at `/private/tmp/ply-docker-probe.Fd9y12` verifies and has SHA-256
  `164f99c7643074bafd9bc1234756e6b51af7139f8dfd96a8142dc4747efb992b`.

Docker acceptance remains blocked on a daemon started outside the managed
sandbox. The next session must repeat the same Client+Server and `docker info`
gate before reading the Dockerfile or selecting an implementation design.

P6 Docker daemon fourth resume attempt (launch continuity `36156e4`):

- Exact ancestry, ordinary and ignored cleanliness, reciprocal archive links,
  the launcher check, and the authorized P6 checkpoint passed again.
- The fail-closed probe stopped before Dockerfile inspection or implementation:
  Docker 29.4.0 returned JSON `Server: null`; `docker version`, `docker info`,
  and `docker desktop status` exited 1 for `desktop-linux`. Buildx 0.33.0 named
  the builder but exposed no driver or supported platforms because the daemon
  was unreachable. The managed sandbox again could not open Docker Desktop's
  required host log.
- No Docker build/run, login, push, publication, design selection, image claim,
  or repository implementation change occurred. The 91-entry external manifest
  at `/private/tmp/ply-docker-probe.g26Ld5` verifies and has SHA-256
  `9c91e10b69eb9b7bf223f135840d11c23775d1e1f51a63715d48a07a34f72a80`.

Docker acceptance remains blocked on a daemon started outside the managed
sandbox. The next session must repeat the same Client+Server and `docker info`
gate before reading the Dockerfile or selecting an implementation design.

P6 Docker-image checkpoint (implementation `7e37293`):

- The mandatory external probe passed before Dockerfile inspection or editing:
  Docker client 29.4.0 reached Docker Desktop 4.70.0 / Engine 29.4.0 through
  `desktop-linux`; `docker version` and `docker info` exited 0; and the daemon,
  builder driver, BuildKit, host, and supported platform identities were
  recorded. The verified probe-manifest SHA-256 is
  `219851ff9c5ee84d8a59ac541445b17f8f28aaf0f095a2e363afa48bb593050e`.
- Executed clean builds classified two local Dockerfile defects: the builder
  needed `python3` to execute the existing launcher contract, and forced
  `GOARCH=amd64` created the wrong executable for a native arm64 image. The
  smallest repair added `python3` and removed forced target variables. The
  accepted non-publishing path used an external empty Docker config and the
  daemon builder with `DOCKER_BUILDKIT=0`; no credential, login, push, publish,
  registry, or remote release operation ran.
- Added executable `scripts/accept-docker`, its 26-control meta-test, private
  Docker acceptance helper, and focused Make/distribution wiring. The path
  requires clean committed archive source, a unique absent tag, exactly one
  no-cache daemon build, immutable and unambiguous identity, platform and
  configuration inspection, and a regular non-symlink executable `/bin/ply`.
  It rejects absent or wrong build/run calls, publication, stale or ambiguous
  identity, tag retargeting, substitution, wrong platform/configuration,
  invalid executable content, skipped or duplicate verifiers, verifier
  failure, duplicate terminal PASS, and repository-local output.
- The authoritative clean build used tag
  `ply-acceptance:7e3729321fdd-20260831t090408-47800` and produced immutable
  `linux/arm64` image
  `sha256:a3b1f58b9861f23f555cdcfc51b6fc82d3452d9425c3314c48997fec5bd475ff`,
  entrypoint `/bin/ply`, size 275,720,311. `/bin/ply` is regular and executable,
  size 23,797,696, SHA-256
  `972bf6fca791fae76eecf09cd9cfcb04ffa7a994405985fa090b304b48ad58f8`.
  Exactly one build, ten runs, and zero publication operations were recorded.
- Host install remained a separate initially missing external `GOBIN` and
  `make install` / `go install ./cmd/ply` contract. `verify-install` passed
  with its artifact override unset. Status, upgrade, and build then each
  emitted one terminal PASS, two non-help behavioral receipts, and three
  runtime receipts proving the immutable image ID, not the tag or a host
  executable, ran. The 506-entry Docker evidence manifest and report SHA-256
  values are
  `6128c6a1c9471a0ffb64b57f519053582986fcb59907135d041c5ad1a6ec2bee`
  and `f9ae547a1bd98fdcbf49b1f4e222d49d6da5bc83ebfe7bec515372ca039c8b09`.
- Current snapshot acceptance reran at `7e37293` with GoReleaser v2.17.1,
  exact non-publishing argv, separate host install, and all three exact
  verifiers. Docker meta 26/26, snapshot meta 20/20, API/CLI and
  entry/subprocess compatibility, pinned lint, complete tests/race/vet,
  launcher and Make contracts, complete preflight, existing host acceptance,
  the independent 15-control audit meta-suite, and empty-HOME count-2 all
  pass. The verified 534-entry gate-manifest SHA-256 is
  `b7d1e6d82eada86df900b22556a74239ab793142836fb0d288971b7a3be3abbf`.
- With no manual evidence supplied, the clean full audit exits 1, never 2, and
  its scorecard SHA-256 is
  `c7c4c9a9cff2831bedd6c5417660551d69771f44e625339d8b01f55b29e1daec`:
  L0 is 8/8, L1 has six PASS and three UNMEASURABLE, L2 has seven PASS and
  three UNMEASURABLE, seven ratchets improve, one regresses, and dirty paths
  are empty. The focused Q2.5-Q2.10 scorecard SHA-256 is
  `4a76baac730083a683580492e7e1eb4b6a7918641775fc1fdaf9a9e87c5bfe3c`:
  Q2.5, Q2.6, Q2.7, and Q2.10 pass while Q2.8 and Q2.9 remain honestly
  unmeasurable.

P6 quality exit checkpoint (implementation `097a9f1`):

- Added the scoped `make quality` apparatus and its executable Make contract.
  The target preserves the existing preflight, eight mutation meta-tests and
  harnesses, four host verifiers, fresh snapshot acceptance, fresh Docker
  acceptance, and authoritative audit entry points. It requires explicit
  external regular tools, schema-2 evidence, output, Go caches, and module
  cache; rejects local or stale state; and proves an exact ordered 21-stage
  population before accepting only `--only Q0.*,Q1.*,Q2.*` at clean L2.
- The canonical external schema-2 document is retained at
  `/private/tmp/ply-p6-quality-review.097a9f1.IfKUcH/manual-evidence-schema-2.json`.
  Its SHA-256 is
  `1601eaa449b2435908024aa8675445d5d95d84cc59e03db8a4529ad573a55bdb`.
  It binds six current receipts to commit `097a9f1`, tree `8f9616f`, the
  inventory and audit instruments, every declared mutation subject, and all
  four core verifier scripts. The focused Q1.6/Q1.7/Q1.9/Q2.4/Q2.8/Q2.9
  audit exits 0.
- The exact complete apparatus exits 0 under
  `/private/tmp/ply-p6-quality-gate.097a9f1-final.lOV5NN`. Its verified
  15,849-entry manifest SHA-256 is
  `b89ac6645f814d40d2444ddca3a808919afacf40c873b576a170bc0e5a4091a6`.
  All 27 Q0-Q2 criteria pass, L2 is attained, manual evidence is valid with
  six receipts, mutation and acceptance denominators are 8 and 4, and held,
  regressed, not-comparable, and dirty-path counts are zero. The scorecard
  SHA-256 is
  `bfe32a5e7eb90e9c7e9cc756a47f38bfefd68866e754cbb030653d409fce5322`.
- The independent regression gate is retained at
  `/private/tmp/ply-p6-regression-gate.097a9f1.YGeuN9`. Its verified
  15,458-entry manifest SHA-256 is
  `85fc11883c3ecd9ce444ae23777ef34c36944273d03fbc1f869b17fa477c7de1`.
  Fresh snapshot and Docker acceptance, compatibility, pinned lint, complete
  tests/race/vet, launcher and Make contracts, complete preflight, host
  acceptance, audit meta 15/15, focused and Q0-Q2 audits, and empty-HOME
  count-2 pass. The separate full report exits the expected 1 at L2 only for
  four queued L3 rows; it is not an exit gate.

P6 and the first bounded P7 maintained-toolchain baseline move are complete.
P7 remains active for dependency groups; P8 remains queued.

The eight named harnesses from `.quality/inventory` were implemented one
subject per measured move. Each declares its mutations, proves it can fail,
includes the methodology T1-T10 meta-controls, and reports declared versus
killed mutations. No survivor remained to require a reachability,
observability, or controllability repair.

Exit: all eight harnesses declare at least eight meaningful mutations, pass
T1-T10, and report `declared == killed`, `survived == 0`, and `unusable == 0`.

### P6 - Build L2 Acceptance Evidence

Status: complete.

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

Status: active for the prepared bounded evaluation of graph-selected
transitive exact `github.com/modern-go/reflect2 v1.0.2`.
Option 1 now explicitly retains exact indirect/production-loaded unqualified
mapstructure v1.5.0, exact direct/runtime-relevant unqualified go-homedir
v1.1.0 and Promptui v0.9.0, exact direct-indirect/runtime-relevant
unqualified emoji/v2 v2.2.12, and exact inherited/unloaded unqualified kr/text
v0.2.0, kr/pty v1.1.1, kr/pretty v0.3.1,
kr/logfmt
v0.0.0-20140226030751-b84e30acd515, kr/fs v0.1.0,
go-windows-terminal-sequences v1.0.1, gotool v1.0.0, and errcheck v1.5.0 under
separate target-specific non-transferable exceptions. Exact
inherited/unloaded httprouter v1.2.0, jtolds/gls v4.20.0+incompatible,
go-junit-report v0.9.1, json-iterator v1.1.12, and clockwork v0.1.0 remain
separately retained and unqualified under their own exceptions. Earlier
demangle, strcase, and affected/not-secure memberlist decisions remain
separate under their own guards. Completed dependency groups remain final
through accepted Google UUID v1.4.0 and qualified go-cleanhttp v0.5.2.
Product source and dependency metadata remain unchanged. Kr/pty v1.1.4 is the
highest behavior-qualified Go-1.18-compatible exact-path stable, but no
genuine supported tidy-stable project owner requests it; v1.1.5-v1.1.8 fail
ordinary documented Start behavior. Neither exact-path kr/text stable
qualifies, and its matching adjacent-overlong-word fix remains unreleased.
The mapstructure, go-homedir, Promptui, emoji/v2, kr/text, and kr/pty option-1
decisions are final under their separate exact expiry guards below. No
canonical mapstructure, go-homedir, kyokomi/emoji/v2, or Promptui stable
qualifies; their decisions retain exact selected v1.5.0, v1.1.0, v2.2.12,
and v0.9.0 without source or dependency change. Modern-go/concurrent release
`1.0.3` qualifies with a genuine supported exact-path owner, and the selected
pseudo-version already names its exact commit/source, so no selection or
metadata changed. The modern-go/reflect2 evaluation successor is prepared but
was not executed. P8 remains queued.

Toolchain baseline move (2026-09-01):

- Selected exact Go 1.26.7. Go 1.26 and Go 1.27 are both maintained, but the
  pinned golangci-lint 2.12.2 binary was built with Go 1.26.2 and its primary
  support policy does not claim targets newer than its build Go line. Pinned
  GoReleaser 2.17.1 was built with Go 1.26.5. Go 1.27.0 was therefore rejected
  for this move rather than upgrading quality-tool dependencies out of scope.
- The official Go archive, actions/go-versions manifest, Docker Official Image
  source, and Docker registry manifest all resolve Go 1.26.7. The module's
  distinct preferred-toolchain directive, Docker builder, deactivated release
  workflow, active documentation, and exact-toolchain baseline now declare one
  version. The module retains `go 1.18` as its characterized language and
  compatibility floor, following Go's documented separation of the `go` and
  `toolchain` directives. Raising that floor provisionally introduced ten
  pinned-lint findings from newer language-analysis semantics; changing source
  or lint policy would not be toolchain-only. The split declaration instead
  preserves the 0-issue pinned-lint contract. The deactivated lint workflow has
  no independent Go selector; inactive Snap packaging remains a separate P8
  scope decision.
- Old Go 1.26.2 reproduced the stored baseline at scorecard SHA-256
  `5fb3226009cfbf0d29f63fa03592157cce4efcec6e38583b64a86f6288e89490`.
  Candidate Go 1.26.7 produced
  `9f044510ae95d1df85b6c0868323c11453f97015c16fc46721d70c9da0b3a463`;
  the only structured difference is `tool.go_build.version`. The raw body,
  every criterion and denominator, and all 228 numeric debt leaves match.
- Before and after the declaration change, both tools selected the same 233
  modules and 143 non-standard packages, emitted byte-identical core help, and
  reproduced the same existing 207-line `go mod tidy -diff`. That untidy
  historical sum/indirect-requirement debt belongs to the queued dependency
  work and was not edited here. The unchanged projection is pinned in
  `toolchain-migration.json`; no dependency metadata was rewritten in this
  toolchain-only move.
- Audit-meta T15 now composes the existing schema-1/schema-2 instrument
  reproduction with the pre-P7/current Go 1.26.2/1.26.7 scorecards. It requires
  byte-exact archived/current scorecards, permits only the recorded Go version
  field between toolchain identities, and independently compares all 228
  numeric debt leaves across both instrument and toolchain axes.
- The focused implementation is commit
  `16ecb67eebb6450226d3264f30a64a14391b5245`, exact parent
  `9852ed15b6e22ec0263bc4f81a973f4e088d7cf5`, clean tree
  `94262fbfd426692fdc44f9eb21f85065ed40a6e2`. No direct or indirect
  dependency version, production Go behavior, API/CLI surface, acceptance or
  mutation population, publisher, registry, credential, or release behavior
  changed.
- Current schema-2 review evidence is retained at
  `/private/tmp/ply-p7-quality-review.16ecb67.Ugm5Bh`; its verified
  22,002-entry manifest SHA-256 is
  `9fc89ac3d4296c2ad7f9a30b6f9decc308113a6b5008afecd12b741983a7e6eb`.
  The canonical evidence SHA-256 is
  `f3ea505bc6bd414f74bcdbc5b5f6b747bfadf7dc48ddc16ad9977937b0c4958d`,
  and the focused six-receipt audit exits 0.
- The exact complete quality apparatus exits 0 at
  `/private/tmp/ply-p7-quality-gate.16ecb67-final.hdvr9q`. Its verified
  246,560-entry manifest SHA-256 is
  `15136fa28368eb4ada0b81dfddc4f36daa19f17b7ab00d3a4ea85d496066cb8e`.
  Its exact 21-stage ledger ends in `audit:Q0.*,Q1.*,Q2.*`; all 27 scoped
  criteria pass at L2 with valid manual evidence, 8/8 mutation and 4/4
  acceptance populations, and zero held, regressed, not-comparable, or dirty
  counts. The Q0-Q2 scorecard SHA-256 is
  `0cf6f16797ec73ed8d51944a77dc0cde4de85e174526b915afb58542e98d4cec`.
- The independent 36-stage regression is retained at
  `/private/tmp/ply-p7-regression-gate.16ecb67-final.4nMFvg`; its verified
  67,726-entry manifest SHA-256 is
  `48be0c3b828c60702ff4f06c513b8720dba4259c9df39d99cf23b0f655dffded`.
  Toolchain declarations, graph selection, API/CLI and entry/subprocess
  compatibility, pinned lint, tests/race/vet, launcher and Make contracts,
  complete preflight, host/snapshot/Docker meta and acceptance, standalone
  audit meta, focused and Q0-Q2 audits, and empty-HOME count-2 pass. The full
  report exits the expected 1 only for queued Q3.1, Q3.3, Q3.4, and Q3.7.
  Fresh-cache replay requires the candidate SDK directory first in `PATH` as
  well as exact `GO`, no ambient `GOFLAGS`, and module-cache warming from an
  external source archive; these conditions keep literal `go` subprocesses on
  1.26.7 without rewriting the repository's historical `go.sum`.
- The complete cross-root evidence index is
  `/private/tmp/ply-p7-complete-move.16ecb67.tGLtNi/evidence-index.json`,
  SHA-256
  `68e9a7cf649a02b882703bd079d37fafab69caba166a97efa164616308b9578b`;
  its verified one-entry manifest SHA-256 is
  `8da5e8a79e26cbb8e4279150f820c8e0bf88da1fb032860015b079d8da718607`.

Terminal dependency group (2026-09-02):

- Selected direct `golang.org/x/term v0.29.0` and its exact MVS coupling
  `golang.org/x/sys v0.30.0`. Both declare Go 1.18. Their immediately newer
  releases raise the module floor to Go 1.23, while the latest releases declare
  Go 1.25, so v0.29.0/v0.30.0 is the highest pair compatible with the retained
  `go 1.18` language floor and preferred Go 1.26.7 toolchain.
- The external old/candidate comparison retains 233 selected modules, 143
  non-standard packages, and 3,551 graph edges. Only the authorized pair and
  its version-bearing edges change. The pre-existing 207-line tidy projection
  becomes 211 lines only for pair checksums; no unrelated selection or checksum
  enters the candidate.
- Exact host and Windows vulnerability scans preserve 22 reachable findings.
  GO-2026-5024 remains module-only for `golang.org/x/sys/windows` before v0.44.0
  and is not symbol-reachable in either selection. Requiring its fixed version
  would raise the language floor and was not mixed into this dependency group.
- The focused implementation is commit
  `0f93a527d1f95cc043a695101746d359c74c9c54`, exact parent
  `e25a98546c20f8707ffe4f71b1d544fc8fb4944f`, clean tree
  `ebbc91e08907b37175faf45ddc3283e8bd0ed771`. Exact
  `go get golang.org/x/term@v0.29.0` changes only `go.mod` and `go.sum`; no
  source, public API/CLI, toolchain, build/release input, quality contract,
  acceptance/mutation population, artifact contract, publisher, registry,
  credential, packaging, or P8 code changes.
- Selection evidence is retained at
  `/private/tmp/ply-p7-terminal-selection.ctFhJk`; its verified 76,252-entry
  manifest SHA-256 is
  `c5cb5e2706ecccfc8eea7c80cadaa747c458e605b2919ab80cd378b263f84aca`.
  Commit-bound schema-2 evidence is retained at
  `/private/tmp/ply-p7-terminal-quality-review.0f93a52.nzAZ5z`; its verified
  252,440-entry manifest SHA-256 is
  `5213d578d639a85bae51c2b85cba73ff72fc380e950afa23dabc1728e3370d4a`.
- Exact `make quality` exits 0 at
  `/private/tmp/ply-p7-terminal-quality-gate.0f93a52.V0TBoc`. Its verified
  268,572-entry manifest SHA-256 is
  `e8027bcb9e2a30ff2fbbf3e3bb8c1a1072abe98a7f55ff6a1fb313c4200d91df`.
  The Q0-Q2 scorecard SHA-256 is
  `1ec4203b099e49017f67a906f19a2ceeb69f1e86b7ca33bd6f8cd1bf23a27cc5`:
  all 27 rows pass at L2 with valid manual evidence, 8/8 mutation and 4/4
  acceptance populations, and zero held, regressed, not-comparable, or dirty
  counts.
- The independent 40-stage regression is retained at
  `/private/tmp/ply-p7-terminal-regression-gate.0f93a52-final.hIYQFq`; its
  verified 68,961-entry manifest SHA-256 is
  `1a4414491a83262c4e7684dcc50d7c8372310a66d41e67cf0b9b9b30aca4964a`.
  Focused callers, pinned lint, tests/race/vet, API/CLI and entry/subprocess
  compatibility, launcher and Make contracts, complete preflight, host and
  fresh snapshot/Docker meta and acceptance, audit meta, focused and Q0-Q2
  audits, vulnerability equality, and empty-HOME count-2 pass. The full report
  exits the expected 1 only for queued Q3.1, Q3.3, Q3.4, and Q3.7.

Logrus dependency group (2026-09-03):

- Selected direct `github.com/sirupsen/logrus v1.9.3`. The originally queued
  assertion that it was the highest Go-1.18-compatible release was corrected:
  v1.9.4 exists and declares Go 1.17, but its module requirements raise
  `github.com/stretchr/testify v1.8.1` to v1.10.0. Version v1.9.3 declares Go
  1.13 and is the highest release that preserves both the retained Go 1.18
  floor and this session's one-module MVS closure. The v1.10.x line declares
  Go 1.23.
- The external old/candidate comparison retains 233 selected modules, 143
  production non-standard dependencies, 197 including tests, and 3,551 graph
  edges. Only logrus changes. The 211-line historical tidy projection becomes
  217 lines only for the two v1.9.3 checksums.
- GO-2025-4188 affects logrus writer-scanner APIs before v1.9.3. Those symbols
  were not reachable here, and the upgrade removes the module-level finding:
  host and Windows reachable findings remain 22 while module findings improve
  from 34 to 33 with no addition.
- The focused implementation is commit
  `efc47ca02f9e02e3d31466ef4ddbb58c90772b9c`, exact parent
  `30a0bfd2a7aa47d8ef6be4c9f7411280f886ce83`, clean tree
  `f8b4d04e75c91981455df97d510fa715be267a47`. Exact
  `go get github.com/sirupsen/logrus@v1.9.3` changes only `go.mod` and `go.sum`.
- Selection evidence is retained at
  `/private/tmp/ply-p7-logrus-selection.30a0bfd.vNfo6P`; its verified
  44,364-entry manifest SHA-256 is
  `678daa930e157e38ffd9b1f08e88eebefc9d00cc150ba8ee153da4820a4f49ad`.
  Schema-2 review evidence is retained at
  `/private/tmp/ply-p7-logrus-quality-review.efc47ca.1rtWsz`; its verified
  42,550-entry manifest SHA-256 is
  `436f8ef7d05a66aea0cf0473ae95c5f5f58a3bd952de8a6b4db9a9def82eaa29`.
- Exact `make quality` exits 0 at
  `/private/tmp/ply-p7-logrus-quality-gate.efc47ca.1944215232`; its verified
  234,632-entry manifest SHA-256 is
  `49c0d0158c7ed5248fa27c33258fe3fe191024f02c684b3050444e567848e40c`.
  The Q0-Q2 scorecard SHA-256 is
  `52b675981b8beb2ca37ea7c56f96e99a76064e9dae8d4c510c749203483b20c0`:
  all 27 rows pass at L2 with 8/8 mutation and 4/4 acceptance populations and
  zero held, regressed, not-comparable, or dirty counts.
- The independent 40-stage regression is retained at
  `/private/tmp/ply-p7-logrus-regression-gate.efc47ca-final.mAc8VU`; its
  verified 68,942-entry manifest SHA-256 is
  `3afd845878b5b3b9874e13c0759aadc7e43d8df2aa5db3b32725c63c635c7c2b`.
  All required compatibility, test, lint, preflight, acceptance, audit,
  vulnerability, hermeticity, and clean-tree stages pass. The full report
  exits the expected 1 only for queued Q3.1, Q3.3, Q3.4, and Q3.7.

  A Cobra-session recheck found that this selection manifest file still has
  the recorded digest and 44,364 entries, but one mutable Go telemetry count
  file beneath its retained HOME no longer matches. The other 44,363 material
  files remain covered by the manifest. No repository state changed; the
  independently rebuilt Cobra selection below supersedes that older selection
  root for current dependency decisions.

Cobra dependency group (2026-09-03):

- Selected direct `github.com/spf13/cobra v1.6.1` -> `v1.10.1` with exactly
  `github.com/cpuguy83/go-md2man/v2 v2.0.2` -> `v2.0.6` and
  `github.com/spf13/pflag v1.0.5` -> `v1.0.9`. Primary repository and module
  evidence confirms that Cobra v1.10.1 declares Go 1.15, go-md2man v2.0.6 and
  pflag v1.0.9 declare Go 1.12, and all are compatible with preferred Go
  1.26.7 and the retained main-module `go 1.18` declaration.
- Cobra v1.10.2 is the latest release, also declares Go 1.15, and replaces the
  existing YAML module path with `go.yaml.in/yaml/v3 v3.0.4`. It therefore
  grows selection from 233 to 234 modules and graph edges from 3,551 to 3,552.
  Version v1.10.1 is the highest release preserving the smaller three-module,
  existing-module closure.
- Old and selected states retain exactly 233 selected modules, 3,551 graph
  edges, and the byte-identical 429-package test population. The historical
  217-line tidy projection becomes 238 lines only for authorized closure
  checksums. Build, complete tests, pinned lint, public help, API/CLI,
  snapshot/Docker meta, and distribution contracts pass. Host/Windows symbol
  findings remain 22 and host module findings remain 33, with exact old/new ID
  equality and no primary Go-vulnerability entry for the three closure
  modules.
- The focused implementation is commit
  `e14ed5ecd0856893c282f37995077542845fb563`, exact parent
  `b10076effe353edfae1020e0585175332cc2ea5d`, clean tree
  `f17f7cf92e3c52be38f1a5180e283023a76aa75a`. Exact
  `go get github.com/spf13/cobra@v1.10.1` changes only `go.mod` and `go.sum`,
  with seven insertions and two deletions; dependency metadata was not
  hand-edited and no behavior, API/CLI, toolchain, release, quality, baseline,
  population, packaging, publisher, registry, credential, or P8 source changed.
- Independent selection evidence is retained at
  `/private/tmp/ply-p7-cobra-selection.b10076e.94XcUY`; its verified
  43,656-entry manifest SHA-256 is
  `5e8a8cce3fe2495620987bce5bc9ebb98ac6db871aa4d7e14875a00f45f05e8f`.
  Its selection-summary SHA-256 is
  `d93ceaa82a5cc2128574b075bb6ee2d7bb26a1d2d0525e42d55919ff5e18b135`.
- Commit-bound schema-2 review evidence is retained at
  `/private/tmp/ply-p7-cobra-quality-review.e14ed5e.SQfQ5c`; its verified
  449-entry manifest SHA-256 is
  `a2becf7b8f273162520bccdad21a2ecaeaca867008fd3d1c1e79d1d9b1d9691e`.
  The canonical evidence SHA-256 is
  `ff89b536faa0ca7f2edf4e3835858d62e62a54697aba961e7ae6526a321fc5b9`,
  every declared subject/file digest was independently refreshed, and all six
  receipts validate.
- Exact `make quality` exits 0 at
  `/private/tmp/ply-p7-cobra-quality-gate.e14ed5e.1788465147.56571`. Its
  verified 234,708-entry manifest SHA-256 is
  `a59de1e99bdb27b1bcd1e364f2a5feb1d9ba7c1dfec1a1fc381fd93327749e71`.
  The Q0-Q2 scorecard SHA-256 is
  `3756e807d0fb56aee5f49b4a43d60e80526c6ab7f1b8456b388cccedfb20b9b8`:
  all 27 rows pass at L2 with valid manual evidence, 8/8 mutation and 4/4
  acceptance populations, 80/80 killed mutations, and zero held, regressed,
  not-comparable, or dirty counts. Its Docker artifact's internal evidence
  manifest SHA-256 is
  `26127866f5f1d2cce2c4a4a3d36ae2770b0df8f415f86c8b13126b1fd8525a3a`.
- The independent 40-stage regression is retained at
  `/private/tmp/ply-p7-cobra-regression-gate.e14ed5e.H7N4M5`; its verified
  97,882-entry manifest SHA-256 is
  `ec18888b35837e6353825cecec7edeb4bd1daa25a7cc428c4ace567130fb80f1`.
  Focused callers, compatibility, lint, tests/race/vet, launcher and Make
  contracts, complete preflight, host and fresh snapshot/Docker acceptance,
  audit meta, focused and Q0-Q2 audits, vulnerability equality, empty-HOME
  count-2, and final cleanliness pass. The full report exits the expected 1,
  never 2, only for Q3.1, Q3.3, Q3.4, and Q3.7; its scorecard SHA-256 is
  `0f9011741f9b421abdf30f42da16deb4b19b3c6641ceb0a886ca0ce86bb60b8c`.
  An initial Docker attempt inherited an isolated HOME that hid Docker
  Desktop's buildx plugin and stopped before build; it is retained alongside
  the passing canonical retry, which used the operator HOME only for plugin
  discovery and a fresh external Docker configuration for the actual build.

Pflag dependency group (2026-09-04):

- Selected existing indirect `github.com/spf13/pflag v1.0.9` -> latest
  `v1.0.10` as a one-module move. Primary Go proxy and repository evidence
  confirms the 2025-09-02 release, exact tag commit
  `0491e5702ad2bb108bc519a5221bcc0f52aa9564`, verified GitHub commit
  signature, Go 1.12 declaration, and checksum pair. The release's only
  production edit restores its declared Go 1.12 compatibility by replacing an
  `errors.Is` comparison with direct `ErrHelp` equality; the remaining release
  edits are documentation or tests.
- Both states retain exactly 233 selected modules, 3,551 graph edges, and the
  byte-identical 429-package population. Only pflag changes. The historical
  tidy projection grows from 238 to 240 lines for the v1.0.10 checksum pair.
  Candidate and committed build, complete tests, race, vet, pinned lint,
  public help, API/CLI and entry/subprocess compatibility, host/snapshot/Docker
  acceptance, audit contracts, and empty-HOME count-2 pass. Darwin symbol,
  Darwin module, and Windows symbol vulnerability ID sets remain exactly
  22/33/22, and the primary Go vulnerability module index has no pflag entry.
- The focused implementation is commit
  `360b2f3c792ca131d259f40852619f2840cefdf1`, exact parent
  `4f0cc642e9c45960b641133c32a0ffc6b1b3a94b`, clean tree
  `a1b941041f5ae685613dd44f4f2db20092cd472c`. Exact
  `go get github.com/spf13/pflag@v1.0.10` changes only `go.mod` and `go.sum`,
  with three insertions and one deletion; dependency metadata was not
  hand-edited.
- Independent selection evidence is retained at
  `/private/tmp/ply-p7-pflag-selection.4f0cc64.58RgL0`. Its final verified
  44,617-entry manifest SHA-256 is
  `b64cad473c92a5551e997010f88243d46fafa2ddc65c11db5d2d6bb90786298c`,
  and selection-summary SHA-256 is
  `e47aa25c01cfaab27409b83b6cbf809e81a315e55725744c17a447df8f96bf79`.
  The original 43,102-entry manifest, SHA-256
  `c7712aafaf9d8440cd2661ce6ffaa4a638321ddbdb9719b49a187420276b430d`,
  verified at seal time; follow-on candidate use later rewrote 27 Go
  build-cache action indexes. The original manifest, passing verification log,
  later mismatch log, and final post-cache manifest are all retained; stable
  selection outputs did not change.
- Commit-bound schema-2 review evidence is retained at
  `/private/tmp/ply-p7-pflag-quality-review.360b2f3.57B9Z6`; its verified
  447-entry manifest SHA-256 is
  `caeace41486ddc6797b5fc737ccd5a4757d2052eb3a5c3451f7da0bcc39f04b8`.
  The canonical evidence SHA-256 is
  `d3c5d725e99fb633186277cd66fc9dc1f17b1c8f0705926b70c449f990526c69`,
  every declared source digest was refreshed, and all six receipts validate.
- Exact `make quality` exits 0 at
  `/private/tmp/ply-p7-pflag-quality-gate.360b2f3.1788472549.40415`. Its
  verified 234,799-entry manifest SHA-256 is
  `47971b2e6cf9399f085e233d9a264f9d22c2e5f45dd024375fbb8e151cbd36c6`.
  The Q0-Q2 scorecard SHA-256 is
  `f50ee942007f4cf5b76b0d87f610469f9d95dc04c6363779e2d6e83a7c632934`:
  all 27 rows pass at L2 with valid manual evidence, 8/8 mutation and 4/4
  acceptance populations, 80/80 killed mutations, and zero held, regressed,
  not-comparable, or dirty counts. Its Docker artifact built once, ran ten
  containers, published zero times, and has internal evidence-manifest
  SHA-256 `3d0f1e342dbd9e218d371e29914f7aafdec93d81920b9cfd3f06a3f77ff094d5`.
- The independent 40-stage regression is retained at
  `/private/tmp/ply-p7-pflag-regression-gate.360b2f3.gb4Rvk`; its verified
  92,211-entry manifest SHA-256 is
  `ec5eb0c20975138fd96b4de3c10f855c0ceb11a5f1152039be2586025beec857`.
  All required selection, compatibility, lint, test/race/vet, launcher/Make,
  complete preflight, host/snapshot/Docker, audit, vulnerability, empty-HOME,
  and clean-tree stages pass. One preflight attempt retained a transient
  launcher signal-fixture failure; the isolated 62-check launcher suite and
  the complete preflight retry passed before acceptance. The full report exits
  the expected 1, never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7; its
  scorecard SHA-256 is
  `4f72c988701326e4a21749306a6155fa849ac3c31ba4e0b293ec6a9f3f92bf08`.

Viper v1.16.0 closure decision (2026-09-04):

- Rejected the direct `github.com/spf13/viper v1.15.0` -> `v1.16.0`
  candidate before implementation. The official Go module reference defines a
  module's `go` directive as the minimum Go version required to use it. P7's
  established compatibility rule has selected dependency versions whose
  module declarations do not exceed the retained `go 1.18` floor; pinned Go
  1.26.7 build success is necessary but cannot replace that floor decision.
- Independent exact-Go replay reproduced the retained closure byte-for-byte:
  233 -> 234 selected modules, 3,551 -> 3,561 graph edges, exactly 31 changed
  selections, only `github.com/google/s2a-go` added, the same exact 429-package
  population, and a 240 -> 270-line historical tidy projection. Candidate
  `go.mod`, `go.sum`, selected modules, graph, packages, and tidy projection
  match the retained probe exactly.
- Reading the target `go.mod` for every changed selection found a stronger
  floor conflict than the earlier handoff singled out. Eleven selected target
  modules declare Go 1.19 and `github.com/stretchr/testify v1.8.3` declares Go
  1.20. The Go-1.19 population is `cloud.google.com/go`,
  `cloud.google.com/go/compute`, `cloud.google.com/go/longrunning`,
  `github.com/googleapis/enterprise-certificate-proxy`,
  `github.com/googleapis/gax-go/v2`, `github.com/hashicorp/consul/api`,
  `go.etcd.io/etcd/api/v3`, `go.etcd.io/etcd/client/pkg/v3`,
  `go.etcd.io/etcd/client/v3`, `google.golang.org/api`, and
  `google.golang.org/genproto` at their exact candidate versions.
- Fresh primary evidence confirms Viper v1.16.0's verified annotated tag
  `c19a006378aa5ee373af48ddee7fdf78d62c5c06`, peeled commit
  `21a7fd828ed231bbe62068d6aafa5aa9f85dc79e`, 2023-05-30 release, Go 1.17
  declaration, and checksum pair
  `h1:rGGH0XDZhdUOryiDWjmIvUSWpbNqisK8Wk0Vyefw8hc=` /
  `h1:yg78JgCJcbrQOvV9YLXgkLaZqUidkY9K+Dd1FofRzQg=`. Latest Viper remains
  v1.21.0 and declares Go 1.23.0.
- Exact Go 1.26.7 build, complete tests, pinned golangci-lint 2.12.2, public
  help, API/CLI reports, module verification, toolchain and snapshot/Docker
  meta-contracts all pass for old and candidate states. API and CLI report
  SHA-256 values remain
  `ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
  and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
  These technical passes explain why no behavior regression was observed; they
  do not erase the selected modules' declared minimum versions.
- Fresh Go vulnerability data preserves exact 22/33/22 Darwin-symbol,
  Darwin-module, and Windows-symbol ID populations. The primary module index
  contains no Viper or Testify entry. No production, test, `go.mod`, `go.sum`,
  API/CLI, toolchain, build/release, quality, baseline, acceptance/mutation,
  packaging, publisher, registry, credential, or P8 file changed.
- The sealed decision and next-group selection root is
  `/private/tmp/ply-p7-viper-decision.360b2f3.Wdvj52`. Its verified
  69,347-entry manifest SHA-256 is
  `0b1e054f4fcfae3a630e5192603cb02f53d2bdd3b7623c3fe14bdb3b74eb757f`;
  `decision-summary.json` SHA-256 is
  `cd8cc895ca5979514e5e584b023ac70d9967cf4951ecf81ddb760267e10cbd7d`.

Uniseg dependency group (2026-09-04):

- Selected existing indirect `github.com/rivo/uniseg v0.4.4` -> latest
  `v0.4.7` as an exact one-module move. Fresh Go proxy and repository evidence
  confirms the 2024-02-08 module publication, lightweight tag at exact commit
  `03509a98a092b522b2ff0de13e53513d18b3b837`, Go 1.18 declaration, and
  checksum pair `h1:WUdvkW8uEhrYfLC4ZzdpI2ztxP1I582+49Oc5Mq64VQ=` /
  `h1:FN3SvrM+Zdj16jyLfmOkMNblXMcoc8DfTHruCPUcx88=`. The bounded changes cover
  Unicode 15 data, segmentation and line-breaking corrections, East Asian
  ambiguous-width configuration, transition/property-search improvements,
  and the v0.4.7 variation-selector width correction.
- Both states retain exactly 233 selected modules, 3,551 graph edges, and the
  exact 429-package population. Only Uniseg changes. The historical tidy
  projection grows from 240 to 244 lines. The caller chain remains
  `plybuild/cmd` -> `go-term-markdown` -> `go-term-text` -> `go-runewidth` ->
  `uniseg`; existing rendering, public help, API, CLI, and artifact contracts
  fully cover the indirect behavior, so no focused source contract was needed.
- The focused implementation is commit
  `a5ff9cf16d5daf2ed6a7578d23a859cfd73b5df3`, exact parent
  `5be563b7b8b09ec74d2aefeb32da3d80fbc16281`, clean tree
  `8c46c2d4a04780043080196b36b86722818aec97`. Exact
  `go get github.com/rivo/uniseg@v0.4.7` changed only `go.mod` and `go.sum`;
  production behavior, public API/CLI, toolchain and main-module language
  declarations, build/release inputs, quality contracts and populations,
  packaging, publishers, registries, credentials, and P8 code did not change.
- Independent selection evidence is retained at
  `/private/tmp/ply-p7-uniseg-selection.5be563b.sDPWuq`. Its verified
  48,985-entry manifest SHA-256 is
  `f0d543dc54efeace23901edf108a919663e657d948965efc9665fd5d42242b6c`,
  and selection-summary SHA-256 is
  `83621fbd0af20b47cbdaa3eb040b4e216dd675387d795379bb8aa32236efd6c0`.
  Old and candidate build, complete tests, pinned lint, byte-identical help,
  API/CLI, artifact-meta, and exact 22/33/22 Darwin-symbol, Darwin-module, and
  Windows-symbol vulnerability populations pass; the primary Go vulnerability
  index has no Uniseg entry.
- Commit-bound schema-2 evidence is retained at
  `/private/tmp/ply-p7-uniseg-quality-review-retry3.a5ff9cf.zhM7ns`. Its
  verified 15-entry manifest SHA-256 is
  `bf1f6e38061fb2f380a85afe91b6eba912ac6a4d3f7e2f8ba2bbb102b391b5b3`,
  canonical evidence SHA-256 is
  `042cd5dcd67b110f23206da3aca58fedb18274a5776ab3ec23b1154e5292a6e3`,
  and all six focused receipts pass.
- Exact `make quality` exits 0 at
  `/private/tmp/ply-p7-uniseg-quality-parent-retry3.a5ff9cf.sNA3P4/quality-gate`.
  Its verified 234,985-entry manifest SHA-256 is
  `f3d5ed9f507d688effb66889aff69ff9160b35099035a3697bd62df0a7982f7e`.
  The exact 21-stage ledger passes; Q0-Q2 scorecard SHA-256
  `b6accd036627e90713fe24e9507b0d80389915501514c2834c467ed0aca715af`
  has all 27 rows at L2, valid schema-2 evidence, 80/80 killed mutations, 8/8
  mutation and 4/4 acceptance populations, and zero held, regressed,
  not-comparable, or dirty counts.
- The independent 40-stage regression is retained at
  `/private/tmp/ply-p7-uniseg-regression-gate.a5ff9cf.W5MLDG`; its verified
  98,323-entry manifest SHA-256 is
  `cd07dc2ff68317e08b1ce647546daff0e7c3cec9f09d591fea473e1e4a9a81dd`.
  Selection, checksums, compatibility, lint, tests/race/vet, launcher and Make,
  complete preflight, host and fresh snapshot/Docker acceptance, audit meta,
  focused and Q0-Q2 audits, vulnerability equality, empty-HOME count-2, and
  cleanliness pass. The separate full report SHA-256 is
  `8a3f113d91f5e3c4a227efed24e412ee9c1e8893bb25f83af0ffaa4377c19038`;
  it exits expected 1, never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7.

Colorable/go-isatty dependency group (2026-09-04):

- Selected existing indirect `github.com/mattn/go-colorable v0.1.13` -> latest
  `v0.1.15` with its exact MVS companion
  `github.com/mattn/go-isatty v0.0.17` -> `v0.0.20`. Fresh Go proxy and
  repository evidence confirms Colorable's 2026-05-29 publication,
  lightweight tag at exact commit
  `8bf39a204f13f0cfcf86ab9b297c3d6e0668e54a`, Go 1.18 declaration, and
  checksum pair `h1:+u9SLTRGnXv73cEsnsmoZBom+dMU88B2M0aDcWy0/jY=` /
  `h1:6LmQG8QLFO4G5z1gPvYEzlUgJ2wF+stgPZH1UqBm1s8=`. Go-isatty v0.0.20 has
  exact tag commit `a7c02353c47bc4ec6b30dc9628154ae4fe760c11`, declares Go 1.15,
  and has checksum pair `h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=` /
  `h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=`.
- Both states retain exactly 233 selected modules, 3,551 graph edges, and the
  exact 429-package population. Only Colorable and go-isatty change; the
  historical tidy projection grows from 244 to 252 lines. The caller chain is
  `plybuild/cmd` -> `go-term-markdown` -> `fatih/color` -> `go-colorable`.
  Existing behavior, help, API/CLI, and artifact contracts cover the indirect
  boundary, so no production or test source contract was added.
- Exact `go get github.com/mattn/go-colorable@v0.1.15` produced focused
  implementation commit `33e187317c6c1be79ef8c5b64ddb2a0f8caec71f`, exact parent
  `e1fd64e4a7bebd6dabf529675b1560172f12f1d5`, clean tree
  `6942bd34879b9935382b44ffab682606cf5ef5c9`. Only `go.mod` and `go.sum`
  changed; all behavior, API/CLI, toolchain, language-floor, distribution,
  quality, packaging, publication, credential, and P8 contracts are unchanged.
- Independent selection evidence is retained at
  `/private/tmp/ply-p7-colorable-selection.e1fd64e.QlohU1`. Its verified
  49,635-entry manifest SHA-256 is
  `b8ec3adb51fc916820b300f0faf4ff0eddc7e1ae51a9f839b96ceed9c57cebb7`,
  and selection-summary SHA-256 is
  `484e11cb3c6929dfd0a50c050a49008df4ed3908133ae3d87a3a3839f783e07c`.
- Commit-bound schema-2 evidence is retained at
  `/private/tmp/ply-p7-colorable-quality-review.33e1873.eGDQR5`. Its verified
  21-entry manifest SHA-256 is
  `f6294878de4b3f19a8731b0b94485f8a203ab20137e37fd3e866f27a0afde24b`;
  the manual evidence SHA-256 is
  `cfbe47aa455636613ec5f03c4ede4a12a4169e2a0da9483cf6fb0519e9a3ce66`,
  and all six focused receipts pass.
- Exact `make quality` exits 0 at
  `/private/tmp/ply-p7-colorable-quality-parent.33e1873.31b10h/quality-gate-retry3`.
  Its verified 235,059-entry manifest SHA-256 is
  `6695ee2f2a181c46fbcf953518c2fd9581902f3f3be28bd761602eecab0785ae`.
  The Q0-Q2 scorecard SHA-256 is
  `55cf7f1a4867f54bba31108f2d7ef104f1b0901b44a5ac29aea11e50aec007e6`:
  all 27 rows pass at L2 with 80/80 killed mutations, 8/8 mutation and 4/4
  acceptance populations, and zero held, regressed, not-comparable, or dirty
  counts. Two earlier quality roots are retained as noncanonical: an initially
  over-narrow PATH exposed the launcher signal-log race; the isolated contract
  and the full corrected-PATH retry passed.
- The independent 40-stage regression is retained at
  `/private/tmp/ply-p7-colorable-regression-gate.33e1873.1Yb6le`; its verified
  98,451-entry manifest SHA-256 is
  `cdf2963affe5e2395ab8c7911efeaa525e4f39e700ede3945e484cfacb63802c`.
  Selection, checksums, compatibility, lint, tests/race/vet, launcher and Make,
  complete preflight, host and fresh snapshot/Docker acceptance, audit meta,
  focused and Q0-Q2 audits, vulnerability equality, empty-HOME count-2, and
  cleanliness pass. The full scorecard SHA-256 is
  `c6e7a71ac8cf187654a83f043a3967d53f8e51035c69bf717909117cd7b8ff1b`;
  it exits expected 1, never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7.
  A sparse-checksum expansion from an aggregate coverage invocation is retained
  as an external 208-line patch and was exactly reversed before every accepted
  gate.

Ini v1.67.3 closure decision (2026-09-04):

- Rejected `gopkg.in/ini.v1 v1.67.0` -> latest `v1.67.3` before
  implementation. Exact replay changes three selections: ini -> v1.67.3,
  `github.com/stretchr/testify v1.8.1` -> v1.11.1, and
  `github.com/stretchr/objx v0.5.0` -> v0.5.2. Objx v0.5.2 declares Go 1.20,
  exceeding the retained Go 1.18 floor; passing Go 1.26.7 technical checks
  cannot replace that compatibility decision.
- The candidate retains 233 modules and the exact 429 packages, grows the graph
  from 3,551 to 3,564 edges, and grows tidy from 252 to 267 lines. Build,
  complete tests, pinned lint, byte-identical help, and exact 22/33/22
  vulnerability parity pass. No tracked file changed.
- The sealed decision root is
  `/private/tmp/ply-p7-next-selection.33e1873.mRehC2`. Its verified
  39,374-entry manifest SHA-256 is
  `c7af82cbe2327de51dc93fcf48c61e65e6b38e5d51a3a7556c07d6d753016181`;
  decision-summary SHA-256 is
  `04081e0b5f0b3324e49ca321f5edf0a2be5b9b938bce31f68c30aabb89c8a90d`.

Cobra/YAML dependency group (2026-09-04):

- Upgraded direct `github.com/spf13/cobra v1.10.1` to latest `v1.10.2` with
  newly selected exact MVS companion `go.yaml.in/yaml/v3 v3.0.4`. Independent
  external replay changed exactly those two selections, moved 233 -> 234
  modules and 3,551 -> 3,552 graph edges, retained the exact 429-package
  population, and grew the historical tidy projection from 252 to 254 lines.
- Current proxy and repository evidence retains Cobra v1.10.2 at lightweight
  tag commit `88b30ab89da2d0d0abb153818746c5a2d30eccec`, Go 1.15, and checksum pair
  `h1:DMTTonx5m65Ic0GOoRY2c16WCbHxOOw6xxezuLaBpcU=` /
  `h1:7C1pvHqHw5A4vrJfjNwvOdzYu0Gml16OCs2GRiTUUS4=`. YAML v3.0.4 remains exact
  tag commit `c3552c15f996075a7634df5159d9161c67bf3d76`, Go 1.16, and checksum pair
  `h1:tfq32ie2Jv2UxXFdLJdh3jXuOzWiL1fo0bu/FbuKpbc=` /
  `h1:DhzuOOF2ATzADvBadXxruRBLzYTpT36CKvDb3+aBEFg=`. YAML v3.0.5 became the
  standalone latest release after the original handoff measurement, but did
  not alter Cobra's bounded exact v3.0.4 MVS closure.
- Exact `go get github.com/spf13/cobra@v1.10.2` produced dependency-only commit
  `620226258841f3692dd918cca54094ad74ca5017`, parent `1f4cb3958768a2673d16b0812ab5fd5dfae055e0`,
  and tree `51e0409a972e9bff20a0c3046063c8a595637a08`. Only `go.mod` and `go.sum`
  changed; the retained `gopkg.in/yaml.v2 v2.4.0` and
  `gopkg.in/yaml.v3 v3.0.1` selections did not move.
- The independent selection root is
  `/private/tmp/ply-p7-cobra-yaml-selection.1f4cb39.TwFCo3`; its verified
  24,805-entry manifest SHA-256 is
  `0eb9ffc0cb5ff418c7304cb41ca3b97f7d9d3c1d5d34898d1ef69f742c350626`,
  and selection-summary SHA-256 is
  `4fa616418c4a3cc4599ea408b959d7ad2a6f1ce217138cc4b66e5edfa0defe82`.
- Commit-bound schema-2 review evidence is retained at
  `/private/tmp/ply-p7-cobra-yaml-quality-review.6202262.lbyesK`; its verified
  19-entry manifest SHA-256 is
  `f282581202b89aaf40188203cb1463588665e14cd232f22e39c6c278fe0dbaff`,
  and manual evidence SHA-256 is
  `d217d09871a23a60fdf80aac74147ff39bfaf860279bec92a2ceac35593226a9`.
- Exact `make quality` exits 0 at the accepted root
  `/private/tmp/ply-p7-cobra-yaml-quality-parent.6202262.ITUpqp/quality-gate`.
  Its verified 235,159-entry manifest SHA-256 is
  `98d22ad621a2c44505cc472f087c7d89583ec73495b45161a06036a0de6127cc`;
  scorecard SHA-256
  `a262ae769245684cfa87650a1c68a46e6df567abe77dc91d594e8aad488429d6`
  records 27/27 PASS at L2, 80/80 killed mutations, 8/8 mutation and 4/4
  acceptance populations, and zero held, regressed, not-comparable, or dirty
  counts.
- The independent regression root is
  `/private/tmp/ply-p7-cobra-yaml-regression-gate.6202262.EfnWY9`; its verified
  58,127-entry manifest SHA-256 is
  `2e17e4301037c8411c2fcc066f15635793aa112f06904fcda8d8ff9edb144077`.
  All 40 stages pass. Full scorecard SHA-256
  `f6146700f8f396dc3c34a3e0816d782b8e361e6765a154eb060295da0901a85c`
  exits expected 1, never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7.
  Vulnerability populations remain exact 22/33/22.

YAML v3.0.5 dependency group (2026-09-04):

- Upgraded selected indirect `go.yaml.in/yaml/v3 v3.0.4` to latest v3.0.5 as
  an exact one-selection move. Current proxy and repository evidence confirms
  the 2026-07-26 publication, lightweight unsigned tag at
  `e16c7af9361b241fa02d91582fb59ce4954d8afc`, Go 1.16 declaration, and
  checksum pair `h1:N6y/pJk8buWs9NY5ERU2HSMfm+IuD/OtfdAnq6kESPw=` /
  `h1:HVTZu1O7/Vkt2N+BFy8Zza+lnLsABggaTM2ZpNIGuKg=`. The 10-commit delta
  changes production files only for documentation-comment formatting; it also
  removes the test-only check dependency, ports tests to standard testing,
  retracts invalid module-path tags, and adds CodeQL CI.
- Both states retain exactly 234 selected modules and the byte-identical
  429-package population. Only YAML changes. The explicit pin adds one graph
  edge, 3,552 -> 3,553, and grows the historical tidy projection from 254 to
  271 lines. No loaded main-module package imports `go.yaml.in/yaml/v3`, so
  projected tidy removes the pin and its two checksum lines; the projection
  was recorded rather than applied.
- Exact `go get go.yaml.in/yaml/v3@v3.0.5` produced dependency-only commit
  `cfdcb370e98ba4e4536f0bf08d5b1abbb01f1856`, parent
  `fdd9986d55d0aee27679ccc3675710e03f0d1e1a`, clean tree
  `7244802bfe8e229ab128c122c0b6cb1a49fc7729`. Only `go.mod` and `go.sum`
  changed with three insertions and no deletions. Cobra remains v1.10.2, and
  retained `gopkg.in/yaml.v2 v2.4.0` and `gopkg.in/yaml.v3 v3.0.1` do not
  move.
- Independent selection evidence is retained at
  `/private/tmp/ply-p7-yaml-selection.fdd9986.cf5t7R`. Its verified
  31,004-entry manifest SHA-256 is
  `1274e8d593cc6dd526ef2c857168212204660653eb6045835ae7943cfc94d986`;
  selection-summary SHA-256 is
  `6dcd8a236217e7b620b7408a4d79270d36003e1364563ead6b875d361950825b`.
- Commit-bound schema-2 evidence is retained at
  `/private/tmp/ply-p7-yaml-quality-review.cfdcb37.3xOd0m`. Its verified
  22-entry manifest SHA-256 is
  `ea61323dafe61f94c4d60846817e4c9e7918a28bf086ce0bb8ea5cac01045dd0`;
  manual evidence SHA-256 is
  `19d5ea397a9db90479596f40db04ca3ab3932b245ccc8334eccb9dd3fdf4e818`,
  and all six focused rows pass.
- Exact `make quality` exits 0 at
  `/private/tmp/ply-p7-yaml-quality-parent.cfdcb37.xxdcLh/quality-gate`. Its
  verified 235,241-entry manifest SHA-256 is
  `919b44aa5ec811e2df49dc7784d02c7f70820e52939e8b8545d8347859db0417`.
  Q0-Q2 scorecard SHA-256
  `9f8e6e70705a7c3d27456a78c8cdf7bf23f5d7ca47dbdfabb83371f15c886c95`
  records 27/27 PASS at L2, 80/80 killed mutations, 8/8 mutation and 4/4
  acceptance populations, and zero held, regressed, not-comparable, or dirty
  counts.
- The independent regression root is
  `/private/tmp/ply-p7-yaml-regression-gate.cfdcb37.dISqmD`. Its verified
  59,161-entry manifest SHA-256 is
  `2f0206fcc142c17858347fd814ffbb4b996d20c3a3b70f9a335aa141c7cc3cdd`.
  All 40 stages pass, vulnerability populations remain exact 22/33/22, and
  regression-summary SHA-256 is
  `8552e0df7ed87a8ad1a327b72e1967a3adb295039126edbf97a5f312e56ecd95`.
  Full scorecard SHA-256
  `1fbccda6eccc11fe9f3628a6dcad906d4d7fab920c7acccedf5c6d2bcf799b68`
  exits expected 1, never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7.

Go-md2man v2.0.7 dependency group (2026-09-04):

- Upgraded selected indirect `github.com/cpuguy83/go-md2man/v2 v2.0.6` to
  latest v2.0.7 as an exact one-selection move. Current proxy and repository
  evidence confirms the 2025-04-24 publication, lightweight tag at verified
  commit `061b6c7cbecd6752049221aa15b7a05160796698`, Go 1.12 declaration, and
  checksum pair `h1:zbFlGlXEAKlwXpmvle3d8Oe3YnkKIK4xSRTd3sHPnBo=` /
  `h1:oOW0eioCTA6cOiMLiUPZOpcVxMig6NIQQ7OS05n1F4g=`. Existing selected
  Blackfriday v2.1.0 has no `go` directive.
- Both states retain exactly 234 selected modules and the byte-identical
  429-package test population. Only go-md2man changes. Exact `go get` adds
  explicit indirect requirements for go-md2man v2.0.7 and already-selected
  Blackfriday v2.1.0, grows the graph 3,553 -> 3,556 edges, and newly records
  Blackfriday's module checksum without moving its selection. The historical
  tidy projection grows 271 -> 278 lines. No loaded main-module package imports
  either module, so tidy would remove the two explicit requirements and their
  newly materialized checksum lines; the projection was recorded, not applied.
- The 17-commit delta fixes roff table rendering for long final-row cells, adds
  its regression test, refreshes CI/lint/docs, and otherwise changes only
  package documentation and lint-comment formatting.
- Exact `go get github.com/cpuguy83/go-md2man/v2@v2.0.7` produced
  dependency-only commit `55dc69dda20c4d8f96b6dbcdd70cfef76467cc85`, parent
  `6fe18fe01326f77bee375d3f8834b44a0416907d`, and clean tree
  `321e39594a4bbe0b7c6afda4e9e5a07cf395fe47`. Only `go.mod` and `go.sum`
  changed with five insertions and no deletions. Cobra remains v1.10.2, YAML
  remains v3.0.5, and both retained gopkg.in YAML selections do not move.
- Independent implementation selection evidence is retained at
  `/private/tmp/ply-p7-md2man-selection.6fe18fe.fBmfEU`. Its verified
  44,219-entry manifest SHA-256 is
  `77c1ec0a2fd1b60e4369d87c3c387efc39bcb1297cc6cf698fac6af3b42d6807`;
  selection-summary SHA-256 is
  `7c5cab15a7a0d813df60b176666fae4b9bd536c2f081566032509da18c1873dd`.
- Commit-bound schema-2 evidence is retained at
  `/private/tmp/ply-p7-md2man-quality-review.55dc69d.BdRjpS`. Its verified
  22-entry manifest SHA-256 is
  `af36ebc39220021cb78b7995f975d68c35b46b52e94dc854d8f5bce17e8fe6c5`;
  manual evidence SHA-256 is
  `e8b1e54604b073ea71d0bf2a568f0e69b60cd22ef90a530eb1cfa233d561c68f`,
  and all six focused rows pass.
- Exact `make quality` exits 0 at
  `/private/tmp/ply-p7-md2man-quality-parent.55dc69d.7VciIk/quality-gate`.
  Its verified 235,335-entry manifest SHA-256 is
  `79765081255f37f2223a4cfa10b37e900cdb2842f3b60589cac985b747ad64ed`.
  Q0-Q2 scorecard SHA-256
  `798ab4be7e9642cd822bb855e5010f6477fa3cd3524c06422da2feef45643c71`
  records 27/27 PASS at L2, 80/80 killed mutations, 8/8 mutation and 4/4
  acceptance populations, and zero held, regressed, not-comparable, or dirty
  counts.
- The independent regression root is
  `/private/tmp/ply-p7-md2man-regression-gate.55dc69d.bD0Kbx`. Its verified
  68,729-entry manifest SHA-256 is
  `bd07a6975b0a495c872bde02c15a3c745e5527768baaf4edded40362af072fa3`.
  All 40 stages pass; regression-summary SHA-256 is
  `e7b852b82883a12d37a1977f4d170dcb16f19520af5d4c05a5215aa9f80cba63`.
  Full scorecard SHA-256
  `ff12e563966050557bc2755e1e60541eebb8cef01cccc66db8cd0ea3b7e7d0b3`
  exits expected 1, never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7.
  Vulnerability populations remain exact 22/33/22.

Regexp2 v1.12.0 dependency group (2026-09-04):

- Upgraded selected indirect `github.com/dlclark/regexp2 v1.8.1` to latest
  v1.12.0 as an exact one-selection move. Current proxy and repository evidence
  confirms its 2026-04-18 publication, lightweight unsigned tag at commit
  `3d5df45b703801b3fe51eb3f5c0dd302e8b0d676`, Go 1.13 declaration, and
  checksum pair `h1:0j4c5qQmnC6XOWNjP3PIXURXN2gWx76rd3KvgdPkCz8=` /
  `h1:DHkYz0B9wPfa6wondMfaivmHpzrQ3v9q8cnmRbL6yW8=`.
- Both states retain exactly 234 selected modules, 3,556 graph edges, and the
  byte-identical 429-package test population. Only regexp2 changes. Exact
  `go get` replaces the indirect requirement and adds the v1.12.0 checksum
  pair while retaining the old v1.8.1 pair. The historical tidy projection
  grows 278 -> 282 lines, retains the explicit v1.12.0 pin, and would remove
  the old selected checksum pair; the projection was recorded, not applied.
- Regexp2 and regexp2/syntax are loaded through
  `plybuild/cmd -> go-term-markdown -> Chroma -> regexp2`. The 12-commit delta
  fixes lazy-loop stack corruption and termination, timeout false positives,
  runner text retention, and ECMAScript plus Singleline dot behavior; it adds
  single-letter Unicode property syntax, timeout test controls, text
  marshaling, regressions, and documentation. The Go vulnerability module
  index has no regexp2 entry.
- Exact `go get github.com/dlclark/regexp2@v1.12.0` produced dependency-only
  commit `3fd6684941c57bf004ac7d09b73f5e716c12689b`, parent
  `d043241b0f092235edb63353a83daf1f57996bd5`, and clean tree
  `8d08fa15c7bd85a229380d8336ca25a3312f19b1`. Only `go.mod` and `go.sum`
  changed with three insertions and one deletion. Production behavior, public
  API/CLI, toolchain and language declarations, distribution, quality,
  packaging, publication, credential, and P8 contracts are unchanged.
- The authoritative pre-implementation selection root remains
  `/private/tmp/ply-p7-next-selection.55dc69d.8kz8Cw`; its verified 49,090-entry
  manifest SHA-256 is
  `a2d03742f02396ede46ac73a19b00191812d84cfcad2b416dcf1f19b460b3a10`.
  The independent implementation replay is retained at
  `/private/tmp/ply-p7-regexp2-replay.d043241.GjPoeU`; its verified
  46,723-entry manifest SHA-256 is
  `5b614df2c43b7769643aab1bd3b96a7bcbb959efc3f63ba3df1ac02d00c5f1fa`,
  and selection-summary SHA-256 is
  `6df19722cc55cf8099f99ccd18b3396836ab58c0dc4ab072eb23c8ad71af5024`.
- Commit-bound schema-2 review evidence is retained at
  `/private/tmp/ply-p7-regexp2-quality-review.3fd6684.l0sSof`. Its verified
  24-entry manifest SHA-256 is
  `62d0e6c4e10db44fcd646bd3683568cf77711b2ece9fb04dfb6bd7dca68a5a62`;
  manual evidence SHA-256 is
  `6a39bf49fe11affd7e75bf09451fd42bb81ec3fdb7a06e2bdc124caa5dc6fa5f`,
  focused scorecard SHA-256 is
  `59d4a096d6cdf4a1d79d8a0876c853a17bffd88285910fb0985ac3226d2e0a56`,
  and all six focused receipts pass.
- Exact `make quality` exits 0 at
  `/private/tmp/ply-p7-regexp2-quality-parent.3fd6684.1BZGVe/quality-gate`.
  Its verified 235,425-entry manifest SHA-256 is
  `9bcc902d7eac1085e8f453989faf4b4d432522eb5a97928454177faba2280115`.
  Q0-Q2 scorecard SHA-256
  `0e7139dc25cc2cd125807536e24b18e27f91f48084cd5ddf644381ee292dd2a8`
  records 27/27 PASS at L2, 80/80 killed mutations, 8/8 mutation and 4/4
  acceptance populations, and zero held, regressed, not-comparable, or dirty
  counts.
- The independent regression root is
  `/private/tmp/ply-p7-regexp2-regression-gate.3fd6684.YQdxua`; its verified
  70,973-entry manifest SHA-256 is
  `2443aa2ba177642c2416f534ab66b83ac1be498080396e75c0ec87e40f3fcce7`,
  and regression-summary SHA-256 is
  `dcce4656a650406e9ee4dbcfb31ea8a69991267c917ce8f96b9c7e94170ab8b7`.
  All 40 stages pass, including complete preflight, host and fresh snapshot/
  Docker acceptance, audit meta, empty-HOME count-2, and exact 22/33/22
  vulnerability parity. Full scorecard SHA-256
  `c95201c96297a198b7bec9b5a031443e42c0fe4be5f17a4d383362b9dce5f527`
  exits expected 1, never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7.

Current queue decisions and next bounded P7 group (2026-09-14):

- Reject `github.com/fatih/color v1.14.1` -> latest v1.19.0 because the latest
  module declares Go 1.25.0, and reject floor-compatible v1.16.0-v1.18.0
  because they break go-term-markdown v0.1.4's exact ANSI reset bytes. Accept
  v1.15.0 as the highest qualified stable release and implement it in exact
  dependency-only commit `6ca672e`.
- Retain `github.com/fsnotify/fsnotify v1.6.0`. V1.7.0 fails its contained-Go
  repeat suite on descriptor exhaustion, v1.8.0 has a deterministic 2024-only
  native-test failure plus a multiple-write timing miss, and v1.9.0 has an
  intermittent Darwin kqueue `bad file descriptor` failure in two independent
  complete Go 1.18 repeats. V1.9.0 already contains an explicitly incomplete
  mitigation; further v1.10 kqueue descriptor fixes require Go 1.23 and cannot
  preserve Ply's retained floor.
- Evaluate `github.com/ghodss/yaml v1.0.0` next. It is selected/latest and the
  proxy's only version, has no `go` directive or declared requirements, and is
  present only through grpc-gateway's MVS edge; no package is loaded by Ply.
- Reject latest indirect `github.com/gomarkdown/markdown`
  `v0.0.0-20260824154242-13c5cf49db8d`. Its exact one-selection closure,
  Go 1.12 declaration, 234-module/3,556-edge/429-package population, and
  22/33/22 -> 20/30/20 vulnerability improvement reproduce, but loaded
  behavior does not: five previously passing go-term-markdown v0.1.4 renderer
  cases panic with `Unknown node type *ast.ReferenceDefinition`.
- Commit `d18ffa2a47071de47ab9ad4dd61f0a01158ff97a` deliberately adds reference
  definitions to the AST for round trip, while the retained consumer's type
  switch panics on unknown nodes. Focused contract commit `8a6ef5f` makes that
  loaded reference-definition behavior executable.
- Accept the immediate pre-breaking fallback as canonical pseudoversion
  `v0.0.0-20260824151336-45814d58469f`, bound to unsigned commit
  `45814d58469f9a462c899674ee58ca440b432dbf`, Go 1.12, and checksum pair
  `h1:21LNG7BIMF2dePpJYpuNzDCaHakI3TYHtOmNRkrNlzA=` /
  `h1:JDGcbDT52eL4fju3sZ4TeHGsQwhG9nbDV21aMyhwPoA=`. It is the immediate parent
  of the breaking AST commit, remains below rejected latest, and contains the
  fixed commits for GO-2023-2074, GO-2024-3205, and GO-2026-5208.
- Old and fallback states retain 234 selected modules, 3,556 graph edges, 429
  complete packages, and four loaded gomarkdown packages. Only gomarkdown and
  the main edge change; tidy projects 282 -> 285 lines. All 30 previously
  passing consumer renderer cases pass, including the five reference cases,
  with the same four fixture failures and zero regressions. Vulnerability
  populations improve exactly 22/33/22 -> 20/30/20 with no additions.
- Exact `go get` produced dependency-only commit `d4f0538`, parent `18d67bb`,
  changing only `go.mod` and `go.sum`. Exact `make quality` passes 27/27 Q0-Q2
  rows at L2, 80/80 mutations, all host/snapshot/Docker acceptance, and zero
  held, regressed, not-comparable, or dirty counts. Independent regression
  passes every required gate; full audit exits expected 1 only for queued L3
  rows.
- Accepted decision, commit-review, exact-quality, and regression evidence is
  sealed respectively at
  `/private/tmp/ply-p7-gomarkdown-fallback-replay.18d67bb.U5kpjt`,
  `/private/tmp/ply-p7-gomarkdown-fallback-quality-review.d4f0538.sBeoKD`,
  `/private/tmp/ply-p7-gomarkdown-fallback-quality-parent.d4f0538-final.2kKFpS`,
  and
  `/private/tmp/ply-p7-gomarkdown-fallback-regression-gate.d4f0538.6YaNVS`.
  Their final verified manifest SHA-256 values are
  `b438197999acb18ad6955d66c5d38e0af92df8b0693ecba7445f5a960b1b1fa8`,
  `cdd7691693373d726fdec2e954aacdc11bcae0bbffde3a471edd56c34de78ab6`,
  `70ef46136c13b7e5276ac1ef1af62514013902e93c9e115551d863f0f3a8e40b`,
  and `daba32728fe7cf75d594371d5e04f0be74101f3a38cd007aaf4e6a51210f8e4f`;
  complete evidence identities are indexed in the rolling handover.
- Accept selected indirect `github.com/google/btree v1.0.0` -> latest v1.1.3.
  Primary proxy, checksum-database, and repository evidence binds its
  lightweight unsigned tag to commit
  `aeba20f7a1e1315badec4eca4fdc9f754f5f880a` at
  2024-08-21T16:26:17Z, declaring Go 1.18, with checksum pair
  `h1:CVpQJjYgC4VbzxeGVHfvZrv1ctoYCAI8vbl07Fcxlyg=` /
  `h1:qOPhT0dTNdNzV6Z/lhRX0YXUafgPLFUh+gZMl761Gm4=`. All nine proxy files
  match the tag-commit archive, and the primary vulnerability module index has
  no btree entry.
- Both states select 234 modules and byte-identical 429-package complete test
  populations. Exactly one selection changes. Graph edges grow 3,556 -> 3,557
  only because the exact selection adds the main-module edge to btree v1.1.3;
  historical transitive edges are unchanged. No btree package loads, and
  `go mod why -m` says the main module does not need it. Tidy projects
  285 -> 299 lines and was recorded, not applied.
- Exact `go get github.com/google/btree@v1.1.3` produced dependency-only commit
  `2e2f8e0`, parent `6f3f39d`, tree `b70ada2`, changing only `go.mod` and
  `go.sum`. Candidate module tests/race/vet and focused clone tests pass.
  Repository build/tests/race/vet, pinned lint, help/API/CLI identity, exact
  quality, empty-HOME count-2, and clean-tree gates pass; vulnerability ID
  populations remain exactly 20/30/20.
- Accepted selection, review, exact-quality, and regression evidence is sealed
  respectively at `/private/tmp/ply-p7-btree-selection.6f3f39d.oPbcTr`,
  `/private/tmp/ply-p7-btree-quality-review.2e2f8e0.kUNRFX`,
  `/private/tmp/ply-p7-btree-quality-parent.2e2f8e0-final.6v4pBz`, and
  `/private/tmp/ply-p7-btree-regression-gate.2e2f8e0.FhWpnk`. Their verified
  manifest SHA-256 values are
  `bc86252747c48e4389fd1b60c8c4843f8003c7342a4f7cf79b66b7f1bf4d8006`,
  `52aa4d4743e2cc70fa9f1d17a2e068941d7dbb67f90ee25e8b2dc6c9652c5ee1`,
  `2395eb672221bac15a70b564dbce09bf28546d93876111dd22fff6500b1fc7cf`,
  and `9cca917b12dd761c58bf91652e78b3e999f55eeb1ffc39f9f3bbf56a52aef463`.
  Exact quality passes 21/21 stages, 27/27 rows at L2, and 80/80 mutations;
  full audit exits expected 1 only for Q3.1, Q3.3, Q3.4, and Q3.7.
- A go-cmp-session prerequisite replay corrected one historical btree evidence
  claim: the btree regression manifest file retains its recorded digest and
  24,197 entries, but 163 mutable cache/HOME files no longer replay after
  preliminary version enumeration appended 161 module-cache `@v/list` files
  and updated one sumdb latest record and one telemetry counter after sealing.
  Its stable regression summary retains SHA-256
  `f348bd0b9e0f68c698c43f466904b66c1c3db6eab5e783b14364aaf40b6b1ea2`.
  The btree decision remains supported by the other three completely verified
  manifests and was superseded for current regression claims by the freshly
  sealed go-cmp gate below.
- Reject latest indirect `github.com/google/go-cmp v0.7.0`. Independent Go
  proxy, checksum-database, and upstream resolution binds it to lightweight
  tag commit `9b12f366a942ebc7254abc7f32ca05068b455fb7`, published
  2025-01-14T18:15:44Z, with module declaration `go 1.21` and checksum pair
  `h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=` /
  `h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=`. That declaration exceeds
  the retained Go 1.18 compatibility floor, so latest was not selected.
- Accept `github.com/google/go-cmp v0.5.9` -> v0.6.0 as the highest compatible
  release in the same module group. The only higher listed release is rejected
  v0.7.0. Version v0.6.0 resolves to lightweight tag commit
  `c3ad8435e7bef96af35732bc0789e5a2278c6d5f`, published
  2023-08-31T17:32:40Z, declares `go 1.13`, and has checksum pair
  `h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=` /
  `h1:17dUlkBOakJ0+DkrSSNjCkIjxS6bF9zb3elmeNGIjoY=`. Both releases have 48
  proxy files byte-identical to their tag-commit archives. Both tags are
  lightweight; GitHub reports valid commit signatures, while local commit
  verification was unavailable because `gpg` is not installed. The primary
  Go vulnerability module index has no go-cmp record.
- The six-commit v0.5.9 -> v0.6.0 history removes purego fallbacks, exercises
  Go 1.20, pins workflow inputs, adds identifier links, and introduces
  `cmpopts.EquateComparable`. Both states retain 234 selected modules, 3,557
  graph edges, 429 complete test packages, and five loaded go-cmp packages.
  Exactly one selection and the matching main edge change. The loaded path is
  `plybuild/cmd -> mvn-pom-mutator/pkg/pom -> go-cmp/cmp`; focused old/new
  consumers pass at count 10. Tidy remains an unapplied 299 -> 308-line
  projection whose only candidate-specific metadata is the version/checksum
  replacement.
- Exact `go get github.com/google/go-cmp@v0.6.0` produced dependency-only
  commit `c314bcb440b5f94871d71a249bca7ae7f87d5543`, parent `cbdb0a9`, tree
  `9a212b3`, changing only `go.mod` and `go.sum`. Candidate module complete
  tests/race and repository build/tests/race/vet/lint pass. Help, API, and CLI
  reports remain byte-identical. Vulnerability ID populations remain exactly
  20 Darwin symbol, 30 Darwin module, and 20 Windows symbol.
- Accepted go-cmp selection, review, exact-quality, and regression roots are
  `/private/tmp/ply-p7-go-cmp-selection.cbdb0a9.kbzc4C`,
  `/private/tmp/ply-p7-go-cmp-quality-review.c314bcb.XeWg4T`,
  `/private/tmp/ply-p7-go-cmp-quality-parent.c314bcb.zxFKyw`, and
  `/private/tmp/ply-p7-go-cmp-regression-gate.c314bcb-final.rMW0l5`. Their
  fully verified manifest populations/SHA-256 values are respectively
  29,812/`dcd18baaa45eb8ad59b268fa7c356164dabc5d4c9836423c3fd1c3aa6d0d1797`,
  9,487/`24ede5a473b9ac9f47839b121d25e9a8d921369c5ba149bf42abd902b5271727`,
  252,077/`6d8104b43fdcb838ec8c0955870798be38a3b4a1f9a7e73f22e683ee566a9f69`,
  and 30,255/`00b736b53494a0c4f27ee7503ecef1c7c7bf74f00bdaab85f7afa9eab19ed3c5`.
  Exact quality passes 21/21 stages, 27/27 Q0-Q2 rows at L2, 80/80 mutations,
  all host/snapshot/Docker acceptance, and zero held, regressed,
  not-comparable, or dirty counts. Independent preflight and regression pass;
  full audit exits expected 1, never 2, only for Q3.1, Q3.3, Q3.4, and Q3.7.
- During the first quality run, midnight cleanup removed standard-library
  sources from the external SDK while leaving its binary intact. The same
  official Go 1.26.7 archive was re-fetched, verified at SHA-256
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`,
  and overlaid at the mandated path; the binary still hashes to the pinned
  `9da68c65...`, and the successful quality/regression runs used that restored
  SDK. Recovery evidence is sealed at
  `/private/tmp/ply-p7-go1.26.7-recovery.qvk4zs` with two-entry manifest
  SHA-256 `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.
- Evaluate selected indirect `github.com/hashicorp/hcl v1.0.0` with no
  dependency change. The Go proxy's canonical `@latest` and `go list -m -u`
  both retain stable v1.0.0; the latter has no `Update` field. Its annotated
  unsigned tag object `5592d2526badd60c172ffa95c6a3b209bea9d1ee`
  points to unsigned commit `8cb6e5b959231cc1119e43259c4a608f9c51a241`,
  published 2018-08-26T00:51:36Z. The 196 proxy files are byte-identical to
  that commit, the module has no `go` directive, and the sumdb pair is
  `h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=` /
  `h1:E5yfLk+7swimpb2L/Alb/PJmXilQ/rhwaUYs4T20WEQ=`.
- Resolve all eight higher v1 strings as application-targeted prereleases:
  `v1.0.1-vault` through `v1.0.1-vault-7` and
  `v1.0.1-nomad-1`. All lack a `go` directive or declare Go 1.14/1.15 and
  every proxy archive matches its upstream commit, but Go's stable-release
  query does not select them. The original `v1.0.1-vault` tag is now absent;
  its 196-file proxy source exactly identifies retained commit `809e678`.
  Vault-7 is the highest string only when targeted prereleases are included;
  its 21-commit branch adds Vault-specific unused-key, nested-JSON, and
  duplicate-key behavior. Nomad-1 is a separate 25-commit targeted line.
- Upstream HCL latest v2.24.0 is module `github.com/hashicorp/hcl/v2`, not an
  update of the selected path. It resolves to commit `6b506809` and declares
  Go 1.23.0, so it both requires an out-of-scope import/API migration and
  exceeds the retained Go 1.18 floor.
- Old and exact-v1.0.0 candidate projections are byte-identical: 234 selected
  modules, 3,557 graph edges, 429 complete test packages, ten loaded HCL
  packages, three main consumer packages, the same checksum pair, and the
  same 308-line unapplied tidy projection. Exact
  `go get github.com/hashicorp/hcl@v1.0.0` changes zero go.mod/go.sum lines.
  The path is `plybuild/cmd -> viper -> viper/internal/encoding/hcl -> hcl`.
  Viper's HCL codec and all three Ply consumers pass at count 10 in both
  states. Darwin-symbol/Darwin-module/Windows-symbol vulnerability populations
  independently remain exact and identical 20/30/20; the primary Go module
  index contains no HCL v1 or v2 record.
- Reject an HCL selection move before repository quality execution because
  required complete module tests fail. Under exact Go 1.26.7 both
  `go test ./...` and `go test -race ./...` exit 1 when default vet rejects
  `hcl/parser/parser_test.go:243`: `%s` is applied to `*ast.LiteralType`.
  Diagnostic `-vet=off` complete tests and race pass all 12 packages, proving
  the failure is historical upstream test drift, but disabling the required
  check is not an acceptance. No repository metadata or implementation commit
  was made and the post-implementation gate was correctly not run.
- Rejected/no-change HCL evidence is sealed at
  `/private/tmp/ply-p7-hcl-selection.99f508d.CFvyqD`: its fully verified
  21,427-entry manifest SHA-256 is
  `ad2954c3ebcf8bfe40e5e1684fcfcf662610555e6ed113515750db1a91c1acff`,
  and selection-summary SHA-256 is
  `c20b0d6354c35adf0e05916610d03a93c0210bb613ae889250316bffc34132d9`.
- Retain selected indirect `github.com/inconshreveable/mousetrap v1.1.0`
  unchanged. Fresh Go proxy and upstream evidence lists only stable v1.0.0,
  v1.0.1, and v1.1.0; canonical `@latest`, `go list -m -u`, the default
  branch, and the exact selected tag all resolve v1.1.0 at commit
  `4e8053ee7ef85a6bd26368364a6d27f1641c1d21`, published
  2022-11-27T22:01:53Z. It declares Go 1.18, has no requirements, and its
  checksum pair is `h1:wN+x4NVGpMsO7ErUn/mUI3vEoE6Jt13X2s0bqwp9tc8=` /
  `h1:vpF70FUmC8bwa3OWnCshd2FqLfsEA9PFc4w1p2J65bw=`.
- All three stable release archives are byte-identical to their upstream tag
  commits. The tags are lightweight; v1.0.0 and v1.0.1 commits are unsigned,
  while GitHub validates v1.1.0's embedded commit signature and local
  verification lacks `gpg`. There are no prereleases, retract directives, v2
  tags/module path, later default-branch commits, or other floor-compatible
  candidates above v1.1.0. Noncanonical alias tags v1.0 and v1.1 point to the
  exact v1.0.0 and v1.1.0 commits.
- Exact `go get github.com/inconshreveable/mousetrap@v1.1.0` emits no output
  and changes zero go.mod/go.sum bytes. Old and candidate states retain 234
  modules, 3,557 graph edges, 429 native and 433 Windows complete-test
  packages, identical checksum populations, and the same 308-line unapplied
  tidy projection. With no module requirements, the minimal selection, edge,
  checksum, and metadata closure is empty.
- Mousetrap is not loaded on Darwin. Windows loads one package through
  `plybuild/cmd -> cobra -> mousetrap`. Old and candidate Ply cmd consumers
  pass at count 10, both states cross-build for windows/amd64, and writable
  exact Cobra source passes all packages at count 10 and compiles its Windows
  test binary. A direct read-only module-cache diagnostic fails only because a
  Cobra test writes beside its source; the byte-identical writable replay
  passes and leaves its source manifest unchanged.
- Mousetrap's complete one-package tests, race, and vet pass, with no upstream
  test files. Repository build, complete tests/race/vet, pinned lint, CLI
  surface, and help/API/CLI identity pass. Vulnerability populations remain
  exact and identical 20/30/20, and the primary Go vulnerability index has no
  mousetrap entry. No dependency commit or post-implementation quality run was
  made for the exact no-op selection.
- No-change evidence is sealed at
  `/private/tmp/ply-p7-mousetrap-selection.516ae39.PzdXjl`: its fully verified
  29,283-entry manifest SHA-256 is
  `71f337b114d1383cd103ef53856a5e767d519eda0b5635e9598b26654a78a152`,
  and selection-summary SHA-256 is
  `64371a02034a3265fa879330e05bcdb477c5c2560de16fa9e57c62e4a5c2e08a`.
- Accept selected indirect `github.com/lucasb-eyer/go-colorful v1.2.0` ->
  canonical latest and highest floor-compatible v1.4.1. The proxy lists eight
  stable releases through v1.4.1, with no prerelease, retraction, or v2 path;
  `@latest`, `go list -m -u`, and the default branch agree. Version v1.4.1 is
  lightweight unsigned tag commit `315b48282c63bac7b48ba128d0c87b7f827b2285`
  at 2026-08-02T08:53:53Z, declares Go 1.12, has no requirements, and has
  checksum pair `h1:1EO+WB73+EH8EVbzlrG3KLAfEypQWVHIBqlTf+2hNss=` /
  `h1:R4dSotOR9KMtayYi1e77YzuveK+i7ruzyGqttikkLy0=`. All 45 proxy files match
  the exact upstream tag archive.
- Old and selected states retain 234 modules, 3,557 graph edges, and 429
  complete-test packages. Exactly one selection and one main edge change;
  go.sum adds only the v1.4.1 pair while retaining v1.2.0. Tidy projects
  308 -> 310 lines and remains unapplied. One package loads through
  `plybuild/cmd -> go-term-markdown -> ansimage -> go-colorful`.
- Candidate complete tests/race, focused MakeColor/HSV behavior, loaded
  consumers, repository build/tests/race/vet/lint, help/API/CLI identity,
  complete preflight, acceptance, audits, empty-HOME count-2, and clean-tree
  gates pass. Standalone candidate vet reproduces five legacy unkeyed-literal
  findings only in doc examples; default test vet and repository vet pass.
  Vulnerability populations remain exact 20/30/20 with no new IDs.
- Exact `go get github.com/lucasb-eyer/go-colorful@v1.4.1` produced
  dependency-only commit `dc4f27c`, parent `f40b032`, tree `a5f3282`, changing
  only `go.mod` and `go.sum`. Exact quality passes 21/21 stages, 27/27 Q0-Q2
  rows at L2, 80/80 mutations, host/snapshot/Docker acceptance, and zero held,
  regressed, not-comparable, or dirty counts. Full audit exits expected 1,
  never 2, only for Q3.1, Q3.3, Q3.4, and Q3.7.
- Accepted selection, review, exact-quality, and regression roots are
  `/private/tmp/ply-p7-go-colorful-selection.f40b032.KHS4f6`,
  `/private/tmp/ply-p7-go-colorful-quality-review.dc4f27c.gs9tGv`,
  `/private/tmp/ply-p7-go-colorful-quality-final.dc4f27c.B7eOLZ`, and
  `/private/tmp/ply-p7-go-colorful-regression-gate.dc4f27c.X2Y4Aj`. Their
  verified manifest populations/SHA-256 values are
  47,633/`a47facc21231accc5c2ed73980c9e3a7542538c6a41d7778de5cfbb410312adc`,
  13/`6444e55370fb37a6c2c06c2c1c280b8fcbd34fc1f94f110a70bc7d067b05673d`,
  236,146/`924500cb4512fa9451cf218b119ce8969e32e7f02a67558807d067e3f41017a6`,
  and 19,011/`2cf393025cf7768ef90fbfa8c0faf378284c3521cb118220dd74ba6fc4735d0e`.
- A properties-session prerequisite replay corrected one inherited
  go-colorful evidence claim: the selection manifest retains its recorded
  digest and 47,633 entries, but the single Go telemetry counter beneath its
  retained HOME no longer matches after later Go commands updated it. The
  other 47,632 entries verify. The stable selection summary retains SHA-256
  `07467cdfd752a8ee9c64b1b6632370c8e0d75b6e3715090931ec8831b2080658`,
  and the completely verified review, exact-quality, and regression roots
  supersede that mutable selection entry for current acceptance claims.
- Reject every selection change for indirect
  `github.com/magiconair/properties v1.8.7` and retain that exact version as
  documented unresolved Go-floor debt. Fresh proxy evidence lists exactly 34
  stable releases, from v1.0.0 through unusual but ordinary same-major
  v1.18.11. `@latest`, `go list -m -u`, the GitHub Release and tag, and the
  default branch all resolve v1.18.11 at commit
  `33f415127a06719be057e464f58f9885c30657aa`,
  2026-07-23T07:36:47Z, with zero later commits. All 34 GitHub Release objects
  are stable; there are no prerelease tags, retractions, v2 tags, or
  major-path changes.
- Releases through v1.8.1 have no Go directive, v1.8.2 through v1.8.6 declare
  Go 1.13, and selected v1.8.7 through canonical latest v1.18.11 declare Go
  1.19. Canonical latest and every higher release therefore exceed the
  retained Go 1.18 floor. Highest compatible v1.8.6 is annotated tag object
  `e97104ee21cdf9cc3ac33992fc829cdc3664d238`, peeled commit
  `869a5592420f4ff6ebf74500b5c658b29d973172` at
  2022-02-23T08:50:35Z and tag time 2022-02-23T08:51:39Z, with checksum pair
  `h1:5ibWZ6iY0NctNGWo87LalDlEZ6R41TqbbDamhfG/Qzo=` /
  `h1:y3VJvCyxH9uVvJTWEGAELF3aiYNyPKd5NZ3oSwXrF60=`. Its 24-file proxy source
  is byte-identical to the tag commit. GitHub reports the tag unverified for
  `bad_email` and the commit signature valid; local verification lacks `gpg`.
- Exact `go get github.com/magiconair/properties@v1.8.6` is not a bounded
  one-module closure: Viper v1.15.0 requires properties v1.8.7, so the command
  downgrades direct Viper to v1.14.0, changes 20 selected versions, and removes
  one module. Modules become 234 -> 233; graph edges become 3,557 -> 3,558
  through 68 removals and 69 additions; complete-test packages remain the
  byte-identical 429. The go.mod projection changes properties and Viper,
  go.sum remains byte-identical because both older checksum pairs already
  exist, and tidy projects 310 -> 312 lines with a 163-line projection delta.
  None of it was applied.
- One properties package loads through `plybuild/cmd -> viper ->
  viper/internal/encoding/javaproperties -> properties`; its codec uses
  `NewProperties`, `WriteComment`, and `Load`. Writable Viper v1.15.0/v1.14.0
  codec tests and three Ply consumer packages pass at count 10. Candidate
  v1.8.6 complete tests and vet pass, but complete race fails twice because
  `assert/TestPanicPanicsAndDoesNotPanic` expects its diagnostic at line 65
  while race instrumentation reports line 66. Selected v1.8.7 reproduces the
  same upstream self-test defect. These closure and self-test failures stop
  selection before repository-wide or post-implementation quality gates.
- Govulncheck v1.7.0 preserves exact old/candidate Darwin-symbol,
  Darwin-module, and Windows-symbol ID populations 20/30/20. The primary Go
  vulnerability module index contains no properties entry. No dependency
  metadata or implementation commit was made.
- Properties decision evidence is sealed at
  `/private/tmp/ply-p7-properties-selection.469049f.iT6qcB`; its fully
  verified 17,240-entry manifest SHA-256 is
  `941e805f6edf4e7c4d07a8912dcf8e5203dd2cea48530aa3fe6995cb19da097f`,
  and decision-summary SHA-256 is
  `23b8556153b7ca9d606d5156d50d7d7791b4978ec2368ac3171ef6370b4f2c2d`.
- Accept selected indirect `github.com/mattn/go-runewidth v0.0.14` -> v0.0.17
  as the highest release compatible with the retained Go 1.18 floor. The Go
  proxy lists 29 stable tags through canonical latest v0.0.29, with no
  prereleases, retractions, v1 path, or later default-branch commits. Releases
  v0.0.1-v0.0.4 have no Go directive, v0.0.5-v0.0.17 declare Go 1.9,
  v0.0.18-v0.0.25 declare Go 1.20, and v0.0.26-v0.0.29 declare Go 1.23.
  Canonical latest and all 11 other releases above v0.0.17 therefore exceed
  the floor.
- Version v0.0.17 is lightweight tag commit
  `94c0db1df07f7755a1e3e962becc01cb5fe8a086` at
  2025-09-25T15:58:59Z. GitHub verifies its commit signature; local signature
  verification is unavailable without `gpg`. Its 16-file proxy archive is
  byte-identical to the tag commit, and its checksum pair is
  `h1:78v8ZlW0bP43XfmAfPsdXcoNCelfMHsDmd/pkENfrjQ=` /
  `h1:Jdepj2loyihRzMpdS35Xk/zdY8IAYHsh153qUoGf23w=`. Canonical latest
  v0.0.29 is lightweight tag commit
  `218f489f6718aaf2c68285d71dfc363e13387075` at
  2026-09-02T01:31:19Z and also has byte-identical proxy/tag source.
- Exact `go get github.com/mattn/go-runewidth@v0.0.17` changes one selection,
  retains 234 modules, 3,557 graph edges, and the byte-identical 429-package
  complete-test population. The main edge and runewidth's unchanged
  `rivo/uniseg v0.2.0` requirement edge are relabeled; MVS continues to select
  uniseg v0.4.7. Go.sum adds only the v0.0.17 checksum pair while retaining
  the historical pair. Tidy remains unapplied and projects 310 -> 313 lines.
- One runewidth package loads through `plybuild/cmd -> go-term-markdown ->
  go-term-text -> go-runewidth`; the direct consumer uses `StringWidth`,
  `Truncate`, `FillRight`, and `RuneWidth`. Candidate complete tests/race/vet,
  Windows test compilation, focused old/new consumer behavior, repository
  build/tests/race/vet/lint, Windows build, help/API/CLI identity, and exact
  vulnerability parity 20/30/20 all pass. The primary Go vulnerability module
  index has no runewidth entry.
- The dependency-only implementation is commit
  `ca19dcd4da9320112e6c9f0c5db507ccbe8b88e5`, exact parent
  `5d0fa3ae0fff98784c349cc23e001018eb9637e1`, tree
  `273c9b635e988ddd686d6c83846f7320bb9e4e33`; it changes only `go.mod` and
  `go.sum`. Exact quality passes 21/21 stages, all 27 Q0-Q2 rows at L2, 80/80
  mutations, host/snapshot/Docker acceptance, and zero held, regressed,
  not-comparable, or dirty counts. The full audit exits expected 1, never 2,
  only for Q3.1, Q3.3, Q3.4, and Q3.7.
- Selection, review, exact-quality, and regression evidence is sealed at
  `/private/tmp/ply-p7-go-runewidth-selection.5d0fa3.dmbKoB`,
  `/private/tmp/ply-p7-go-runewidth-quality-review.ca19dcd.UgdZ2N`,
  `/private/tmp/ply-p7-go-runewidth-quality-final.ca19dcd.t6LExP`, and
  `/private/tmp/ply-p7-go-runewidth-regression-gate.ca19dcd.O6Uxyq`. Their
  fully verified manifest populations/SHA-256 values are
  28,870/`2d6986055d24d4937dcbbaa8fa1e46da6a7bcc1e94f481125343284c96687193`,
  15/`f885b995dc6192cb7ef2d8c9241f93cb44a885deda7faec7e0415f68ce21d764`,
  236,329/`28fbbde41059dd806dfea5653f273a390b6b1fa2f60ebeec2e3d654cc792e2ed`,
  and 4,888/`b83091c1c8a52ec2176c4c3fd7a4c3c5e456adf7a61ee2c6b6f7f97706eb4894`.
- Retain selected indirect `github.com/mitchellh/mapstructure v1.5.0`
  unchanged. The Go proxy lists exactly 17 stable releases from v1.0.0 through
  v1.5.0; proxy `@latest` and `go list -m -u` both select v1.5.0, with no
  `Update` field. There are no prereleases, retractions, original-repository
  v2 tags, or GitHub Release objects. Releases through v1.1.2 have no `go`
  directive and v1.2.0-v1.5.0 declare Go 1.14, so canonical latest is also the
  highest stable release compatible with the retained Go 1.18 floor.
- Version v1.5.0 is lightweight tag commit
  `ab69d8d93410fce4361f4912bb1ff88110a81311` at
  2022-04-20T22:31:31Z. GitHub verifies its commit signature; local
  verification is unavailable without `gpg`, and the lightweight tag has no
  tag object to sign. Its 13-file proxy archive is byte-identical to the tag
  commit, and its checksum pair is
  `h1:jeMsZIYE/09sWLaz43PL7Gy6RuMjD2eJVyuac5Z2hdY=` /
  `h1:bFUtVrKA4DC2yAKiSyO/QUcy7e+RRV2QTWOzhPopBRo=`. The archived original
  repository has 12 later untagged commits and an unreleased 1.5.1 changelog,
  but neither qualifies as a stable release; the module has no formal
  `Deprecated` marker.
- The maintained fork's v1.6.0 is published from
  `github.com/go-viper/mapstructure` while declaring the original module path,
  so adopting it requires an explicit `replace` and is not a canonical
  original-path update. `github.com/go-viper/mapstructure/v2 v2.5.0` is a
  distinct module/import path requiring source migration. Neither lineage is
  part of this bounded no-change group.
- Exact `go get github.com/mitchellh/mapstructure@v1.5.0` exits 0 without
  output and changes zero go.mod/go.sum bytes. Old and candidate states remain
  byte-identical at 234 selected modules, 3,557 graph edges, 429 native
  complete-test packages, five mapstructure checksum lines, and a 313-line
  unapplied tidy projection. With no module requirements, the minimal
  selection, edge, checksum, package, and metadata closure is empty.
- One package loads through `plybuild/cmd -> viper -> mapstructure`. Viper's
  focused unmarshalling tests pass in ten fresh processes in both states, its
  source manifests remain unchanged, and the three Ply consumer packages pass
  at count 10. Mapstructure complete tests, count-10 tests, race, and vet pass.
  Repository build/tests/race/vet, Windows build, pinned lint, help/API/CLI
  identity, CLI surface, and the full preflight pass. The final preflight has
  all 80 mutation kills and 15 audit meta-controls; no post-implementation
  quality gate was run because no selection changed.
- Exact Darwin-symbol/Darwin-module/Windows-symbol vulnerability populations
  remain 20/30/20 with identical ID sets. The primary Go vulnerability index
  has no original mapstructure-path entry. Two records associated with the
  maintained fork paths do not apply to the selected original module path.
- No-change decision evidence is sealed at
  `/private/tmp/ply-p7-mapstructure-selection.f974b58.GbAMTJ`; its fully
  verified 46,352-entry manifest SHA-256 is
  `f15fbc74ed6efd7ee0b3c5ce91b0a9ded4104c5d363aca1272848ec5fca25ecf`,
  and selection-summary SHA-256 is
  `8910174c8f8fa387a32cd49871dfbd37668ff4a2fbf1c51c2dae6fdc2612eff3`.
  Dependency metadata remained byte-identical and no implementation commit
  was manufactured.
- Retain selected indirect `github.com/pelletier/go-toml v1.9.5` unchanged.
  The v1 proxy lists exactly 28 stable releases from v0.1.0 through v1.9.5;
  proxy `@latest`, exact `go list -m -versions -retracted`, `go list -m -u`,
  the selected tag, and its stable GitHub Release agree on v1.9.5. There are
  no v1 prereleases or retractions. Releases through v1.2.0 have no `go`
  directive and v1.3.0-v1.9.5 declare Go 1.12, so selected v1.9.5 is both
  canonical latest and the highest stable release compatible with Go 1.18.
- V1.9.5 is lightweight tag commit
  `fed1464066413075eac02cd4dc368b5221845541`, committed
  2022-01-05T14:17:32Z and published as a stable GitHub Release
  2022-04-21T23:21:51Z. The lightweight ref has no tag object or tag
  signature. GitHub validates the embedded commit signature; local
  verification is unavailable because `gpg` is absent. Its checksum pair is
  `h1:4yBQzkHv+7BHq2PQUZF3Mx0IYxG7LsP222s7Agd3ve8=` /
  `h1:u1nR/EPcESfeI/szUZKdtJ0xRNbUoANCkoOuaOx1Y+c=`. All 65 module-ZIP-eligible
  files match the exact tag commit; Go correctly excludes six files belonging
  to the nested benchmark module.
- The repository remains active and unarchived for its default `v2` branch,
  while `master` says v1 will receive no updates and recommends v2. The v1
  module has no formal `Deprecated` marker. The only commit after v1.9.5 on
  `master` changes that README notice; it is not a later v1 release.
  `github.com/pelletier/go-toml/v2` is a distinct module and import path, not
  an in-place v1 selection.
- Exact `go get github.com/pelletier/go-toml@v1.9.5` exits 0 without output or
  go.mod/go.sum changes. Old and candidate states remain byte-identical at 234
  selected modules, 3,557 graph edges, 429 native complete-test packages,
  three v1 checksum lines, and a 313-line unapplied tidy projection. V1.9.5
  has zero requirements, so the exact selection, edge, checksum, package, and
  metadata closure is empty.
- No v1 package loads and `go mod why -m` says the main module does not need
  it, so no focused loaded behavior applies. Five packages from the distinct
  `/v2` module do load through Viper. V1.9.5 complete tests, count-10 tests,
  and race pass across six packages. Standalone `go vet ./...` diagnoses 132
  legacy unkeyed `Position` literals in query test files; this is not a
  candidate regression because exact selection is unchanged, default complete
  tests pass, and repository vet is clean.
- Repository build, complete tests/race/vet, Windows build, pinned lint, CLI
  surface, byte-identical help/API/CLI, empty-HOME count-2, and the final full
  preflight pass. Two earlier preflights retain the established launcher
  signal-fixture timing diagnostic; isolated direct and Make launcher suites
  pass 62/62 before the final preflight passes all repository, 80 mutation,
  and 15 audit-meta controls. No post-implementation quality run was required
  because dependency metadata did not change.
- Govulncheck v1.7.0 preserves exact and identical Darwin-symbol/
  Darwin-module/Windows-symbol populations 20/30/20. The primary Go
  vulnerability module index has no go-toml v1 or v2 entry. No dependency
  implementation commit was manufactured. Decision evidence is sealed at
  `/private/tmp/ply-p7-go-toml-selection.09cc40b.jGVQz6`, whose fully verified
  37,719-entry manifest SHA-256 is
  `e3e516511d60caf7d2004ee1feb29261ecc45b35339dac8e44d19f286606e1b8`;
  selection-summary SHA-256 is
  `85deac527b69ebbfd2608a970684b3f915563d48a538644f3da357961adc268e`.
- Accept selected indirect `github.com/pelletier/go-toml/v2 v2.0.7` -> v2.2.2
  as the highest stable release compatible with the retained Go 1.18 floor.
  Fresh proxy and upstream evidence enumerates 33 semantic versions: 23
  stable and ten prereleases, with no retractions. Canonical latest v2.4.3
  declares Go 1.21.0 and is rejected; all 15 stable releases through v2.2.2
  declare Go 1.16. The v1 module is a distinct path, the proxy has no v3
  versions, and a non-semver GitHub prerelease named `latest` is not a module
  version.
- Candidate v2.2.2 is lightweight tag commit
  `a3d5a0bb530b5206c728eed9cb57323061922bcb` at
  2024-04-29T10:02:54Z. Its GitHub Release was published
  2024-05-01T15:13:03Z. All 33 tags are lightweight refs with no tag object or
  tag signature; GitHub verifies every embedded commit signature. All 97 files
  in the candidate module ZIP match the tag commit byte-for-byte. Its checksum pair is
  `h1:aYUidT7k73Pcl9nb2gScu7NSrKCSHIDE89b3+6Wq+LM=` /
  `h1:1t835xjRzz80PqgE6HHgN2JOsmgYu/h4qDAS4n929Rs=`.
- Exact pinned get changes go-toml/v2 v2.0.7 -> v2.2.2, testify v1.8.1 ->
  v1.9.0, and objx v0.5.0 -> v0.5.2. V2.2.2 directly requires testify
  v1.9.0, which requires objx v0.5.2; objx's v1.8.4 testify edge contributes
  only its go.mod checksum under MVS. Modules remain 234, complete packages
  remain 429, graph edges change 3,557 -> 3,565 through exactly six removed
  and 14 added edges, and go.sum adds exactly six lines. The unapplied tidy
  projection changes 313 -> 321 lines.
- Five `/v2` packages load through `plybuild/cmd -> viper ->
  viper/internal/encoding/toml -> pelletier/go-toml/v2`. Viper encode/decode
  tests pass at count 10 in both states. Candidate complete tests at count 1
  and count 10 and race pass across 16 packages. One standalone module-vet
  malformed test-tag diagnostic is normalized-identical in v2.0.7 and v2.2.2
  and dates to 2021 commit `0d20a845`; repository vet remains clean.
- Dependency-only implementation commit
  `8cce879d285f08a3f8710a1ea35358fb057a9bfc`, exact parent `34f1314`, tree
  `8a88c3ff919b5275599d8461c1a176457997d2e7`, changes only `go.mod` and
  `go.sum`. Exact quality passes 21/21 stages, all 27 Q0-Q2 rows at L2, 80/80
  mutations, host/snapshot/Docker acceptance, and zero held, regressed,
  not-comparable, or dirty counts. The full audit exits expected 1, never 2,
  only for Q3.1, Q3.3, Q3.4, and Q3.7. Help/API/CLI identity and exact
  20/30/20 vulnerability populations are preserved; the primary vulnerability
  index has no go-toml record. During documentation handoff, two Make launcher
  attempts reproduced the established nested signal-fixture timing diagnostic;
  an isolated run and the final Make entry-point run pass all 62 controls.
- Selection, review, exact-quality, and regression evidence is sealed at
  `/private/tmp/ply-p7-go-toml-v2-selection-final.34f1314.gn2ns2`,
  `/private/tmp/ply-p7-go-toml-v2-selection.34f1314.s0l6bA/quality-review-final`,
  `/private/tmp/ply-p7-go-toml-v2-selection.34f1314.s0l6bA/exact-quality-final`,
  and
  `/private/tmp/ply-p7-go-toml-v2-selection.34f1314.s0l6bA/regression-gate-final`.
  Their fully verified manifest populations/SHA-256 values are
  17,093/`ce1a47d85cd2fa0ff3bde19c1953c0dc73914017145abe2f4c5a27c4cf685ff1`,
  8/`328672c0410e2bfffa058fdf5ec389f4c8640bd73caccde22603b0ad57c9b070`,
  236,599/`85f24a00baaf63d62bb6a3a07dc97bffcaa98b8fed1417b7d6f8425ba8d98d3d`,
  and 3,743/`9601ce867fe389eb1efca59b2ae578052ba9b48670d1112cef3dc66298026ec7`.
- Reject selected indirect `github.com/spf13/afero v1.9.4` -> v1.10.0
  despite v1.10.0 being the highest stable release compatible with the retained
  Go 1.18 floor. The proxy lists 37 main-module versions, all stable, with no
  prereleases, retractions, or formal deprecation markers. Releases through
  v1.2.2 have no Go directive, v1.3.0-v1.8.2 declare Go 1.13,
  v1.9.0-v1.10.0 declare Go 1.16, v1.11.0 declares Go 1.19, v1.12.0
  declares Go 1.21, and v1.13.0-v1.15.0 declare Go 1.23.0. Canonical latest
  v1.15.0 therefore exceeds the floor.
- Candidate v1.10.0 is lightweight tag commit
  `ee6eef77ef4a6c73b07a4bc070a9a2f076fd121e` at
  2023-09-22T14:18:35Z. It has no tag object or tag signature, while GitHub
  verifies its commit signature. Its checksum pair is
  `h1:EaGW2JJh15aKOejeuJ+wpFSHnbd7GE6Wvp3TsNhb6LY=` /
  `h1:UBogFpq8E9Hx+xc5CNTTEpTnuHVmXDwZcZcE1eb/UhQ=`; all 66 proxy files
  match the exact tag commit. Canonical latest v1.15.0 is annotated tag object
  `18d690e34969d06817fa791ccf69194ebd4a5e8d`, peeled commit
  `399bb34ad9fd8a252ad1d8bfaef96279b66dc774`, with GitHub-verified tag
  and commit signatures. The active upstream is neither archived nor a fork;
  four gcsfs/sftpfs tags are nested modules, no v2 path exists, and 77 later
  master commits are unreleased.
- Projected exact get changes Afero plus x/crypto
  `v0.0.0-20220525230936-793ad666bf5e` ->
  `v0.0.0-20220722155217-630584e8d5aa`. Afero directly raises x/crypto,
  whose candidate declares Go 1.17; x/net and x/text requirements remain
  below retained selections. This minimal closure keeps 234 modules and the
  byte-identical 429-package population, changes graph edges 3,565 -> 3,572,
  and adds four checksum lines. The unapplied tidy projection changes
  321 -> 323 lines.
- Three Afero packages load through `plybuild/cmd -> spf13/viper ->
  spf13/afero`. Focused Viper filesystem/configuration behavior, Ply
  `./cmd`, and the candidate root package pass at count 10. Candidate
  complete count-1 tests, race, and vet pass across seven packages.
- Required complete count-10 module tests deterministically fail in tarfs
  `TestRead` with 27 identical diagnostics. Selected v1.9.4 normalizes to
  the exact failure; the test is byte-identical, dates to 2020 commit
  `a4ea980f`, and passes in ten fresh count-1 processes. Although historical
  rather than a candidate regression, it violates the explicit module
  self-test stop rule. No dependency metadata or implementation commit was
  made, and repository-wide candidate gates after the stop were not run.
- Govulncheck v1.7.0 preserves exact and identical Darwin-symbol,
  Darwin-module, and Windows-symbol ID populations 20/30/20. The primary Go
  vulnerability index has no Afero record.
- Rejection evidence is sealed at
  `/private/tmp/ply-p7-afero-selection.5a46539.udvGSh`; its fully verified
  79,866-entry manifest SHA-256 is
  `f0ad8ce41436a5f8bf48222106e79f57b32a9533b6a9a15b290615a54f162405`,
  and decision-summary SHA-256 is
  `c3301a1c2b13f36956551ed1eb09b1a83bf2ac51cf236c6614f46eaf70ea65de`.
- Accept selected indirect `github.com/spf13/cast v1.5.0` -> v1.5.1 as the
  highest stable release compatible with the retained Go 1.18 floor. The Go
  proxy lists 17 stable releases and no prereleases or retractions. V1.5.0 and
  v1.5.1 declare Go 1.18; the immediately newer v1.6.0 declares Go 1.19.
  Canonical latest v1.10.0 declares Go 1.21.0 and therefore exceeds the floor.
- Candidate v1.5.1 is unsigned lightweight tag commit
  `bbed5559a59db54120dda79106f85f9ec0afb038` at
  2023-05-15T08:28:46Z and was published as a stable GitHub Release. Its
  checksum pair is
  `h1:R+kOtfhWQE6TVQzY+4D7wJLBgkdVasCEFxSUBYBYIlA=` /
  `h1:b9PdjNptOpzXr7Rq1q9gJML/2cdGQAo69NKzQ10KN48=`. All 13 proxy files
  match the tag commit byte-for-byte. Canonical latest v1.10.0 is annotated
  tag object `3efe057bdf41dab273f8e65893fbcb89c1a4dc3b`, peeled commit
  `fc73346bfc4e6597bc520fb6eea04360299e77d2`, with GitHub-verified tag and
  commit signatures. The active upstream is neither archived nor a fork; no
  later major module exists, while 32 post-v1.10.0 master commits are
  unreleased.
- Exact pinned get changes cast plus its minimal test-only MVS closure:
  quicktest v1.14.3 -> v1.14.4, kr/pretty v0.3.0 -> v0.3.1, and
  rogpeppe/go-internal v1.6.1 -> v1.9.0. Their declarations are Go 1.13,
  Go 1.12, and Go 1.17. Existing go-cmp v0.6.0, kr/text v0.2.0, and pkg/diff
  selections already dominate or retain the candidate requirements. Modules
  remain 234, complete packages remain byte-identical at 429, graph edges
  change 3,565 -> 3,564 because v1.5.1 drops cast's x/xerrors edge, and go.sum
  adds exactly five lines. The unapplied tidy projection changes 321 -> 332
  lines.
- One cast package loads through `plybuild/cmd -> spf13/viper -> spf13/cast`.
  Viper root plus its dotenv, ini, and javaproperties encoders use 17 cast
  symbols. Ten fresh focused Viper runs, encoder count-10 tests, and Ply
  `./cmd` count-10 tests pass in both states. Candidate complete count-1 and
  count-10 tests, race, vet, and module verification pass for its sole package.
- Dependency-only implementation commit
  `cf4fd4933251b2ca7844b6420fd3cf2a86f049a5`, exact parent
  `17f7277fd2512ba06db76d1acb0af5f8b623243c`, tree
  `485690c010cef3b02385d8f2471d3ef53611abe9`, changes only `go.mod` and
  `go.sum`. Repository build, tests/race/vet, Windows build, pinned lint, CLI
  surface, launcher, byte-identical help/API/CLI, and empty-HOME count-2 pass.
  Exact Darwin-symbol/Darwin-module/Windows-symbol vulnerability populations
  remain identical at 20/30/20, and the primary Go vulnerability index has no
  cast record.
- Exact quality passes all 21 stages, all 27 Q0-Q2 rows at L2, both 80/80
  mutation populations, host/snapshot/Docker acceptance, and zero held,
  regressed, not-comparable, dirty, or publication counts. The independent
  focused audit is 27/27 PASS; the full audit exits expected 1, never 2, only
  for Q3.1, Q3.3, Q3.4, and Q3.7.
- Selection, review, exact-quality, and regression evidence is sealed at
  `/private/tmp/ply-p7-cast-selection.17f7277.9nRKpS`, its
  `quality-review-final` child,
  `/private/tmp/ply-p7-cast-exact-quality-final.cf4fd49.wlsG4F/output`, and
  `/private/tmp/ply-p7-cast-regression-gate.cf4fd49.sAkB56`. Their fully
  verified manifest populations/SHA-256 values are
  82,586/`dfa699f75d0b222d7d64e8fd7a09d05306de98ef4ac5cbf6d253fd33b2dddac0`,
  12/`850d40f319d13024275b35909e8ff84ce2f9c5ce5959cf8d707461a3b5e5d9ce`,
  236,782/`b3d410a70d89e2176daf1a013f97d53e704ef3887f02d67bd6c2564ab0a6d02a`,
  and 607/`97cafb6719f470ab7d057a603b522f0d10418deda2d997e6d5723883efd2e24a`.
  Decision, quality, and regression summary SHA-256 values are
  `13c7ab63019a891d66508638b6ade40155ffb340548f36a17e761b204da89c87`,
  `c8c2bc454c816656297ab1622da50b534c84b090c0a2ddac6366f065cca9d89d`,
  and `533ccdf2704153e7d4f467e3ceb75f932202a1367ec38dc1c3f8df51ec3a458f`.
- Retain selected indirect `github.com/spf13/jwalterweatherman v1.1.0`
  unchanged. The Go proxy lists only stable v1.0.0 and v1.1.0, with no
  prereleases or retractions; both omit a Go directive. Proxy `@latest`, exact
  `go list -m -versions -retracted`, the lightweight tag, and the sole stable
  GitHub Release agree on v1.1.0, so it is both canonical latest stable and
  the highest stable release compatible with the retained Go 1.18 floor.
- V1.1.0 is lightweight tag commit
  `94f6ae3ed3bceceafa716478c5fbf8d29ca601a1` at
  2018-10-28T14:53:47Z; GitHub verifies the commit signature and the tag has
  no independent tag object. Its checksum pair is
  `h1:ue6voC5bR5F8YxI5S67j9i582FU4Qvo2bmqnqMYADFk=` /
  `h1:aNWZUN0dPAAO/Ljvb5BEdw96iTZ0EXowPYD95IqWIGo=`. All nine proxy files
  match the exact tag commit; ZIP SHA-256 is
  `43cc5f056caf66dc8225dca36637bfc18509521b103a69ca76fbc2b6519194a3`
  and file-manifest SHA-256 is
  `dd5824d06389cec49ebdfed59947ae4d96bf1e9ef0ced0bab30d7088820451ee`.
- The active, unarchived, undisabled upstream is not a fork and has no later
  major path. Master commit `91990e8269243f4f5d234f9c24e0a73c4a0df47f`
  is three commits after v1.1.0; its proxy pseudo-version declares Go 1.20,
  but that line is unreleased and excluded from stable selection. V1.0.0 is
  an unsigned annotated tag over a GitHub-verified commit. Local signature
  verification remains unavailable because `gpg` is absent.
- Exact `go get github.com/spf13/jwalterweatherman@v1.1.0` changes zero
  go.mod/go.sum bytes. Both projections are byte-identical at 234 modules,
  3,564 graph edges, 429 complete packages, the same checksum pair, and the
  same 332-line unapplied tidy diff. Jwalterweatherman's requirements on
  go-spew v1.1.1 and go-difflib v1.0.0 remain selected; its testify v1.2.2
  edge loses to selected v1.9.0. The minimal selection, edge, checksum, and
  metadata closure is empty.
- One package loads through `plybuild/cmd -> spf13/viper ->
  spf13/jwalterweatherman`. Viper's root package uses the TRACE, DEBUG, INFO,
  WARN, and ERROR loggers. Candidate module tests at count 1 and count 10,
  race, vet, and verification pass; focused Viper read/merge/override/alias
  behavior in ten fresh processes and Ply `./cmd` count-10 tests pass in both
  states. Repository build, complete tests/race/vet, pinned lint, CLI surface,
  byte-identical help/API/CLI, complete preflight, and empty-HOME count-2 pass.
- Exact Darwin-symbol/Darwin-module/Windows-symbol vulnerability populations
  remain identical at 20/30/20. The 1,392-entry primary vulnerability index
  has no jwalterweatherman record. No dependency metadata or implementation
  commit was made, so no post-implementation exact-quality run was required.
- No-change evidence is sealed at
  `/private/tmp/ply-p7-jwalterweatherman-selection.2dd8fdf.h4EKeQ`; its fully
  verified 76,374-entry manifest SHA-256 is
  `54040aefec38df7415f56da66af04d48c8256dd46ae6de8e06a7c4bbe5ce4d8e`,
  and decision-summary SHA-256 is
  `342330b73509caf414a406dda3a14c026bfc59512d09bf6393f710127974dc52`.
- Retain selected indirect `github.com/subosito/gotenv v1.4.2` unchanged.
  Fresh proxy evidence lists eleven stable releases through v1.6.0, no
  prereleases or retractions, and no later major path. Proxy `@latest`, the
  version list, and the signed upstream tag establish v1.6.0 as both canonical
  latest stable and the highest stable release compatible with Go 1.18 by its
  declaration. The repository is active, unarchived, undisabled, and not a
  fork; six later master commits are unreleased, and current master declares
  Go 1.22.
- V1.6.0 is signed annotated tag object
  `b1e4a564ab83bcd77e00d5687055911884986d98`, peeled to signed commit
  `14a05352a5cf0f66fd7cbce114374f56065891f0` at
  2023-08-15T12:05:45Z. Its checksum pair is
  `h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=` /
  `h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSWPKKo0FU=`. All 21 proxy files
  match the exact tag commit; ZIP SHA-256 is
  `142db3dd2328e744c157e85cf3291d027013b79f92a45984f860fe38bc0f1f8d`
  and file-manifest SHA-256 is
  `a0044f44e956b351d4482487488db8e95e5e9d3210a42524a50b1274b5bfc7c9`.
- Exact v1.6.0 projection would move gotenv v1.4.2 -> v1.6.0,
  `golang.org/x/text` v0.7.0 -> v0.12.0, `golang.org/x/mod`
  v0.6.0-dev.0.20220419223038-86c51ed26bb4 -> v0.8.0, and
  `golang.org/x/tools` v0.1.12 -> v0.6.0. Every selected declaration remains
  Go 1.17 or Go 1.18.
  Modules stay 234, graph edges rise 3,564 -> 3,568, complete packages rise
  429 -> 434 through five x/text packages, and go.sum gains exactly the
  gotenv and x/text checksum pairs, four lines. The unapplied tidy projection
  changes 332 -> 335 lines. This is the explained minimal MVS closure.
- Exactly one gotenv package loads through `plybuild/cmd -> spf13/viper ->
  spf13/viper/internal/encoding/dotenv -> subosito/gotenv`; Viper's dotenv
  codec calls `gotenv.StrictParse`. Ten-run focused Viper internal codec,
  public dotenv read/write, and Ply `./cmd` behavior pass with normalized
  old/candidate identity.
- Reject v1.6.0 because its complete module tests fail `TestScanner` under
  exact Go 1.26.7 at count 1, count 10, and race: trailing LF, CR, and CRLF
  fixtures expect four scanner tokens but produce three. Verification and vet
  pass. Selected v1.4.2 fails equivalently; the upstream correction exists
  only on the unreleased Go-1.22 master line. The mandatory self-test stop rule
  therefore forbids a dependency edit and makes downstream repository gates
  inapplicable. A first attempt lacked its external GOTMPDIR; the canonical
  replay created that directory and exposed the actual test failure.
- Exact Darwin-symbol/Darwin-module/Windows-symbol vulnerability populations
  remain identical at 20/30/20. The fresh 1,392-entry primary vulnerability
  index has no gotenv record. Rejection evidence is sealed at
  `/private/tmp/ply-p7-gotenv-selection.3a8cadc.Dai0IG`; its fully verified
  48,662-entry manifest SHA-256 is
  `ab2bd356b863d529085cff9b918d004aa2a03a121a7ea360a66767a2820cc88c`,
  and decision-summary SHA-256 is
  `014c050a3ccefd6fa501f4c5d04aded2ee93d02ad4039a9935d6f4d478c21e11`.
- Reject selected indirect `golang.org/x/image v0.5.0` -> v0.16.0 even though
  v0.16.0 is the highest stable release whose complete changed closure
  preserves the retained Go 1.18 floor. Fresh proxy evidence enumerates 45
  stable semantic releases, v0.1.0 through canonical latest v0.45.0, with no
  prereleases, retractions, deprecation marker, later major path, or GitHub
  Release objects. V0.45.0 declares Go 1.25.0. Releases through v0.24.0
  directly declare at most Go 1.18, but v0.17.0 is the first whose x/text
  requirement selects an x/tools pseudo-version declaring Go 1.19; v0.16.0's
  x/text v0.15.0 closure instead selects Go-1.18 x/tools v0.6.0 and x/mod
  v0.8.0. One post-v0.45.0 master commit is unreleased and declares Go 1.26.0.
- Candidate v0.16.0 is lightweight unsigned tag commit
  `55c4ab6bd625a2e8433671ec9f9b6c46daddf2cf` at
  2024-05-05T12:58:27Z. Its checksum pair is
  `h1:9kloLAKhUufZhA12l5fwnx2NZW39/we1UhBesW433jw=` /
  `h1:ugSZItdV4nOxyqp56HmXwH0Ry0nBCpjnZdpDaIHdoPs=`. All 253 proxy files
  match the exact authoritative upstream tag; ZIP SHA-256 is
  `2ccf3619375bf633937c1b72ccd9a426953aead2ec9b2ccb0ad39e089ea6ac8a`.
  The authoritative Go repository and unarchived, undisabled, non-fork GitHub
  mirror have identical master and all 45 tag refs. Every tag is a lightweight
  commit ref with no independent tag signature; relevant commits are unsigned.
- Exact projected get changes x/image v0.5.0 -> v0.16.0 and x/text v0.7.0 ->
  v0.15.0 in go.mod, adds exactly their two checksum pairs, and selects x/mod
  v0.8.0 and x/tools v0.6.0 while the pruned legacy graph exposes goldmark
  v1.3.5 instead of v1.4.13. Existing x/sys v0.30.0 dominates v0.5.0. Modules
  remain 234, complete packages remain byte-identical at 429, graph edges
  change 3,564 -> 3,548 through 22 removals and six additions, and the
  unapplied tidy projection changes 332 -> 381 lines. No module or package is
  added; this is the explained minimal MVS selection/edge/checksum closure.
- Eight x/image packages load through `plybuild/cmd -> go-term-markdown ->
  pixterm/pkg/ansimage -> x/image/bmp`; ansimage blank-imports bmp, tiff, and
  webp, and imaging imports bmp and tiff. All eight loaded decoder packages
  pass count-10 in both projections with normalized-identical output. The
  consumer image-render fixture reproduces the same retained gomarkdown
  expected-text mismatch in both states and has normalized-identical output,
  so it exposes no x/image behavior delta.
- Reject v0.16.0 because its required complete module tests fail under exact
  Go 1.26.7 at count 1, count 10, and race: `draw.TestScaleDown` rejects four
  golden images after Go 1.26 changed JPEG decoding. Vet also fails with 80
  unkeyed `basicfont.Range` literals. Selected v0.5.0 has the same failure
  classes. Upstream's explicit JPEG-tolerance fix first appears in v0.32.0 and
  its vet fix first appears in v0.34.0; both releases declare Go 1.24.0 and
  cannot repair the floor-compatible line. The mandatory module self-test and
  vet stop rules forbid dependency edits and make downstream repository and
  post-implementation quality gates inapplicable.
- A projected v0.16.0 would reduce exact Darwin-symbol/Darwin-module/
  Windows-symbol vulnerability populations 20/30/20 -> 18/28/18 by removing
  GO-2023-1989 and GO-2023-1990. The fresh 1,392-entry primary vulnerability
  index has 13 x/image records. That improvement does not override a mandatory
  gate failure. No dependency metadata or implementation commit was made.
  Rejection evidence is sealed at
  `/private/tmp/ply-p7-x-image-selection.9e9ef26.9J9x1F`; its fully verified
  52,533-entry manifest SHA-256 is
  `84972465aef5f19f888102c7991aadd67446146a7056658b006d005666bd82bf`,
  and decision-summary SHA-256 is
  `2b70234a63a199ab6cf11050212645162a78c4836aa9aafa4cb00dad764eb936`.
- Reject selected indirect `golang.org/x/net v0.7.0` -> v0.25.0 even though
  v0.25.0 is the highest stable release whose complete changed closure
  preserves the retained Go 1.18 floor. Fresh proxy evidence enumerates 58
  stable releases through canonical latest v0.58.0, with no prereleases,
  retractions, deprecation marker, nested or later major path, or GitHub
  Release objects. V0.58.0 declares Go 1.25.0. Releases through v0.35.0
  directly declare at most Go 1.18, but v0.26.0 is the first whose x/text
  requirement selects an x/tools pseudo-version declaring Go 1.19; v0.25.0's
  x/text v0.15.0 closure instead selects Go-1.18 x/tools v0.6.0 and x/mod
  v0.8.0. Twenty-two post-v0.58.0 master commits are unreleased and current
  master declares Go 1.26.0.
- Candidate v0.25.0 is lightweight unsigned tag commit
  `d27919b57fa8dd03198f85ca9e675e1a09babd7d` at
  2024-05-06T16:24:48Z. Its checksum pair is
  `h1:d/OCCoBEUq33pjydKrGQhw7IlUPI2Oylr+8qLx49kac=` /
  `h1:JkAGAh7GEvH74S6FOH42FLoXpXbE/aqXSrIQjXgsiwM=`. All 778 proxy files
  match the exact authoritative upstream tag; ZIP SHA-256 is
  `7fd8464681c3011736f2c75beb20f88fff553a17f4f574325bce5ca5dc1fcf83`.
  The authoritative Go repository and unarchived, undisabled, non-fork GitHub
  mirror have identical master and all 58 lightweight unsigned tag refs.
- Exact projected get changes x/net v0.7.0 -> v0.25.0 and x/text v0.7.0 ->
  v0.15.0 in go.mod, adds exactly their two checksum pairs, and selects
  x/crypto v0.23.0, x/mod v0.8.0, and x/tools v0.6.0 while existing x/sys
  v0.30.0 and x/term v0.29.0 dominate lower requirements. Modules and complete
  packages remain 234 and 429, graph edges change 3,564 -> 3,568 through five
  removals and nine additions, and the unapplied tidy projection changes 332
  -> 335 lines. This is the explained minimal MVS selection/edge/checksum
  closure.
- Two x/net packages load through `plybuild/cmd -> go-term-markdown ->
  x/net/html`. Both packages and the Ply markdown contract pass count-10 in
  both projections with normalized-identical output. The upstream consumer
  test reproduces the same pre-existing ANSI expected-text mismatch in both
  states, so it exposes no x/net behavior delta.
- Reject v0.25.0 because its required complete module tests fail under exact
  Go 1.26.7 at count 1, count 10, and race: `route.TestRouteMessage` cannot
  create an AF_ROUTE raw socket in the measured Darwin sandbox and reports
  `operation not permitted`. Selected v0.7.0 fails the same route test, and
  the fatal socket behavior remains on unreleased master. Candidate vet and
  verification pass, but the explicit self-test stop rule forbids dependency
  edits and makes downstream repository and post-implementation quality gates
  inapplicable.
- The projection preserves exact Darwin and Windows symbol populations at
  20/20, including the same nine symbol-reachable x/net IDs and traces. Darwin
  module findings improve 30 -> 27 by removing GO-2023-1988, GO-2023-2102,
  and GO-2024-2687, but that module-only improvement does not override the
  failed quality contract. No dependency metadata or implementation commit was
  made. Rejection evidence is sealed at
  `/private/tmp/ply-p7-x-net-selection.e5ea7ae.cP57X8`; its fully verified
  71,887-entry manifest SHA-256 is
  `88e4be87483322792ad3da0064bcef4ced5ee21c52821bb1f80c143f80a3aaf1`,
  and decision-summary SHA-256 is
  `6e32698abcde67604e159488ba21f6e780fd61dfba2ed9299644b5f011940dff`.
- Reject selected indirect `golang.org/x/text v0.7.0` -> v0.15.0 even though
  v0.15.0 is the highest stable release whose complete changed closure
  preserves the retained Go 1.18 floor. Fresh proxy evidence enumerates 49
  stable releases through canonical latest v0.41.0, with no prereleases,
  retractions, deprecation marker, fork, archive/disable state, or nested or
  later major path. V0.41.0 declares Go 1.25.0. Releases v0.14.0-v0.22.0
  directly declare Go 1.18, but v0.16.0 first requires x/tools pseudo-version
  `v0.21.1-0.20240508182429-e35e4ccd0d2d`, which declares Go 1.19 and remains
  required through v0.22.0. V0.15.0 instead requires x/tools v0.6.0, x/mod
  v0.8.0, and x/sys v0.5.0, all at Go 1.18 or lower. Twelve post-v0.41.0
  master commits are unreleased, and master declares Go 1.26.0.
- Candidate v0.15.0 is lightweight unsigned tag commit
  `8d533a0c40adec778a7d09ac6c8aa640d3c883f4` at
  2024-04-15T18:14:38Z. Its checksum pair is
  `h1:h1V/4gjBv8v9cjcR6+AR5+/cIYK5N/WAgiv4xlsEtAk=` /
  `h1:18ZOQIKpY8NJVqYksKHtTdi31H5itFRjB5/qKTNYzSU=`. All 542 proxy files
  match the authoritative tag; ZIP SHA-256 is
  `13faee7e46c8a18c8a28f3eceebf15db6d724b9a108c3c0482a6d2e58ba73a73`.
  Selected v0.7.0 is lightweight unsigned tag commit
  `71a9c9afc4cd710b9412f7f99f0d8e35b10e488a` at
  2023-01-31T16:01:06Z; all 530 proxy files match. Authoritative and mirror
  master/tag refs are identical.
- Exact projected get changes only x/text v0.7.0 -> v0.15.0 in go.mod, adds
  exactly its two checksum lines, and changes MVS selections x/mod
  `v0.6.0-dev.0.20220419223038-86c51ed26bb4` -> v0.8.0 and x/tools v0.1.12 ->
  v0.6.0. Existing x/sys v0.30.0 dominates v0.5.0, and x/sync stays v0.1.0.
  Modules and complete packages remain 234 and 429; graph edges change 3,564
  -> 3,567 through one removal and four additions; go.sum projects 1,031 ->
  1,033 lines; the unapplied tidy projection changes 332 -> 333 lines. This
  is the explained minimal selection/edge/checksum closure.
- Three x/text packages load through `plybuild/cmd -> spf13/viper ->
  spf13/afero -> x/text/runes`: runes, transform, and unicode/norm. Afero uses
  `transform.Chain`, normalization forms, rune removal, and `transform.String`.
  Loaded-package, Afero, and Ply cmd focused tests pass count 10 in both
  projections with normalized-identical output.
- Reject v0.15.0 because complete module tests at count 1/count 10, race, and
  vet all fail under exact Go 1.26.7. Three stale Example identifiers fail
  vet, and message/pipeline panics in x/tools v0.6.0's SSA builder on a Go 1.26
  range-over-function construct; repeated tests additionally expose
  `cases.TestShortBuffersAndOverflow` persistence failures. Standalone vet
  adds 20 unkeyed literals, four unreachable statements, and one unused
  currency string result. Selected v0.7.0 reproduces the failure classes.
  Example repairs first ship in v0.18.0, whose closure already declares Go
  1.19. Verification passes and source is unchanged, but the mandatory module
  gate forbids dependency edits and makes downstream repository and
  post-implementation gates inapplicable.
- Exact Darwin-symbol/Darwin-module/Windows-symbol vulnerability populations
  remain 20/30/20. Both states have the same sole x/text module finding,
  GO-2026-5970, fixed in v0.39.0; neither has an x/text symbol-reachable
  finding or trace. No dependency metadata or implementation commit was made.
  Rejection evidence is sealed at
  `/private/tmp/ply-p7-x-text-selection.6ab5945.3qn5Bj`; its fully verified
  82,629-entry manifest SHA-256 is
  `48c872da4ab796ef1115003fcf0a226a07bf5791a163d36bc9c946c9d9f92ed6`,
  and decision-summary SHA-256 is
  `50d6bee90ea7cea060e23400e820b1537742ec2167bee04f188070e2e5eb39c9`.
- Retain selected `gopkg.in/yaml.v2 v2.4.0` without dependency metadata
  edits. Fresh proxy evidence enumerates exactly 14 stable releases through
  canonical latest v2.4.0, with no prereleases or retractions. V2.4.0 declares
  Go 1.15 and is also the highest stable release compatible with the retained
  Go 1.18 floor for this exact module path. Its checksum pair is
  `h1:D8xgwECY7CYvx+Y2n4sBz93Jn9JRvxdiyyo8CTfuKaY=` /
  `h1:RDklbk79AGWmwhnvt/jBztapEOGDOx6ZbXqjP6csGnQ=`.
- V2.4.0 is lightweight unsigned tag commit
  `7649d4548cb53a614db133b2a8ac1f31859dda8c` at
  2020-11-17T15:46:20Z. All 24 proxy files match both the archived original
  repository and active successor mirror; ZIP SHA-256 is
  `ede49e27c4cca6cdd2ec719aed8ea4d363710cceb3d411e7a786fbdec0d391fd`.
  Later v2.4.1-v2.4.4 tags use distinct module path `go.yaml.in/yaml/v2`, and
  the v3/v4 paths are not candidates for this group.
- Exact projected `go get gopkg.in/yaml.v2@v2.4.0` changes zero go.mod/go.sum
  bytes, selections, edges, or checksums. Both states retain 234 modules,
  3,564 edges, 429 complete packages, 1,031 checksum lines, and the same
  332-line unapplied tidy projection. The retained later gocheck selection
  dominates YAML v2's declared 2016 pseudo-version, so the changed closure is
  empty. The sole loaded package is consumed by `pkg/config`, whose actual
  Marshal/Unmarshal behavior passes count 10 with normalized-identical output.
- The module verifies and passes count-1 and race tests, but independent old
  and no-op-candidate replays both fail count 10 with nine persistent
  `S.TestLineWrapping` failures after `FutureLineWrap` changes irreversible
  global state. Vet also fails identically with 27 legacy malformed test
  struct tags. The mandatory module stop rule therefore forbids an
  implementation commit and makes downstream repository candidate gates
  inapplicable.
- Exact Darwin-symbol/Darwin-module/Windows-symbol vulnerability populations
  remain 20/30/20 with identical IDs and no YAML v2 finding or trace. The
  fresh primary index's three YAML v2 records are all fixed before v2.4.0.
  No dependency commit was made. No-change evidence is sealed at
  `/private/tmp/ply-p7-yaml-v2-selection.1969442.0BpgSy`; its fully verified
  29,219-entry manifest SHA-256 is
  `8dd34cc54f56e4bb1370b2f7f1f124aa8b501a4dcca177effe4f6377c3919e14`,
  and decision-summary SHA-256 is
  `61d8132e932d4613aea071c5f3d55a9a38c2873e2428011b83231611517885d6`.
- Retain selected `gopkg.in/yaml.v3 v3.0.1` without dependency metadata
  edits. Fresh Go proxy and checksum-database evidence enumerates exactly two
  stable releases, v3.0.0 and v3.0.1, with no prereleases or retractions.
  Both releases omit a `go` declaration and require only
  `gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405`, which also omits a
  Go declaration. V3.0.1 is therefore both canonical latest and the highest
  stable release compatible with the retained Go 1.18 floor for this exact
  module path; this conclusion comes from declarations and the empty changed
  closure, not from its modern-toolchain test result. Its checksum pair is
  `h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=` /
  `h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=`.
- V3.0.1 is lightweight unsigned tag commit
  `f6f7691b1fdeb513f56608cd2c32c51f8194bf51`, tree
  `1cb2e60a039c6b3cfdfbd33cc2d049bf68c7db00`, at
  2022-05-27T08:35:30Z. All 24 proxy files match both the archived original
  `go-yaml/yaml` tag and the active `yaml/go-yaml` successor mirror tag; the
  proxy ZIP SHA-256 is
  `aab8fbc4e6300ea08e6afe1caea18a21c90c79f489f52c53e2f20431f1a9a015`.
  Both GitHub repositories report `fork: false`, while the successor describes
  itself as the maintained YAML-org fork. The original v3 branch has exactly
  one later unreleased README-only commit, which marks the archived project
  unmaintained; the selected module has no formal deprecation marker and
  neither repository publishes GitHub Release objects.
- Successor tags v3.0.2-v3.0.5 belong to distinct module path
  `go.yaml.in/yaml/v3`; that path retracts its invalid v3.0.0-v3.0.1 tags and
  retains accepted v3.0.5. `gopkg.in/yaml.v2 v2.4.0` is a distinct retained
  major, and `go.yaml.in/yaml/v4` has only v4.0.0-rc.1 through rc.6. None is
  an upgrade candidate for the exact `gopkg.in/yaml.v3` group. Relevant
  v3.0.0-v3.0.1 history is one commit adding nil-token parser guards and an
  invalid-input regression; the GO-2022-0603 fix commit precedes v3.0.0.
- Exact projected `go get gopkg.in/yaml.v3@v3.0.1` exits zero with empty
  stdout/stderr and zero-byte go.mod/go.sum diffs. Both states retain 234
  selected modules, 3,564 graph edges, 429 complete packages, 1,031 go.sum
  lines, and an identical 332-line unapplied tidy projection. Selected
  `gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15` already dominates
  YAML v3's declared 2016 pseudo-version, so the minimal changed closure has
  zero selections, edges, or checksums. Project module verification passes.
- Exactly one YAML v3 package loads. Real paths include `plybuild/cmd ->
  pkg/config -> gopkg.in/yaml.v3`; `cmd/tips.go` and `cmd/build_options.go`
  call `GitCloudConfig.GlobalCloudConfig`, whose implementation uses
  `yaml.Unmarshal`. An external-only focused fixture calls that production
  path and validates decoded fields and `SourceFor`; ten old and candidate
  repetitions pass with normalized-identical output.
- The module source remains byte-identical, verifies, lists one package, and
  passes complete count-1, count-10, and race tests in two independent old
  and no-op-candidate replays. Mandatory vet fails reproducibly in both states
  with exactly 32 identical legacy malformed struct-tag diagnostics across
  decode, encode, and node tests. The stop rule therefore makes repository
  build, complete tests/race/vet, pinned lint, help/API/CLI, snapshot/Docker,
  quality, and audit acceptance inapplicable. No dependency implementation
  commit was made.
- Govulncheck v1.7.0 preserves byte-identical old/candidate
  Darwin-symbol/Darwin-module/Windows-symbol populations at 20/30/20 with no
  YAML v3 finding or trace. The fresh 1,392-entry primary index contains only
  GO-2022-0603 for this module, fixed at
  `v3.0.0-20220521103104-8f96da9f5d5e`; both stable releases follow that fix.
  No-change evidence is sealed at
  `/private/tmp/ply-p7-yaml-v3-selection.cdb5b4c.bbbTdr`; its fully verified
  40,121-entry manifest SHA-256 is
  `84bd12b7b4cbff806abcbff213b4b7bc5bac230a7371c3fe0b3e2def643ed741`,
  and decision-summary SHA-256 is
  `95d5aa003b77880ecf72770ad66c061e0013dfd022fc222eb600860bc9a7799f`.
- Retain selected
  `gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15` without dependency
  metadata edits. Fresh proxy evidence exposes no listed stable or prerelease
  semantic versions. Proxy `@latest` and exact Go resolution select
  `v1.0.0-20201130134442-10cb98267c6c` at 2020-11-30T13:44:42Z. It is a
  pseudo-version, not a stable release; there are no Git tags, retractions, or
  GitHub Release objects. The candidate declares Go 1.11 and is canonical
  latest and the highest exact-path version compatible with Go 1.18.
- Candidate commit `10cb98267c6cb43ea9cd6793f29ff4089c306974`, tree
  `b5ee34e90064a88d7614bf6f6714d32bfac33e66`, is the exact v1 head in both
  fresh gopkg and GitHub clones. All seven resolved pseudo-version snapshots
  match upstream Git sources. The repository is active and not a fork; its
  separate `master` head is an older divergent 2014 line, not an upgrade
  candidate. The candidate has no tag; GitHub reports its embedded commit
  signature verified and valid, while local cryptographic verification was
  unavailable because `gpg` is absent.
- Exact projected `go get` advances Check and adds creack/pty v1.1.9 and
  pkg/diff's 2021 pseudo-version: 234 -> 236 selected modules and 3,564 ->
  3,574 graph edges. Complete packages remain 429. Four indirect requirements
  and eight go.sum lines are projected, with no removals; all changed-closure
  declarations are absent or at most Go 1.17. The unapplied tidy projection
  grows from 332 to 353 lines. Project module verification passes.
- Check loads in zero complete project packages in either state and has zero
  repository imports. The `go mod why` path crosses `gopkg.in/yaml.v2.test`,
  so this is a dependency-test requirement rather than a main-module consumer.
- Two independent candidate module replays verify, list one package, and pass
  count-1 and race tests, but both count-10 runs fail because global suite
  registration accumulates to 16 suites rather than eight. Both vet runs also
  fail on the same four unkeyed go/printer and go/ast composite literals. The
  mandatory module stop rule rejects the candidate before downstream
  repository, quality, snapshot/Docker, and audit gates. No dependency commit
  was made.
- Old and candidate govulncheck outputs and exact IDs remain byte-identical at
  20/30/20 for Darwin symbol, Darwin module, and Windows symbol scans. The
  fresh 1,392-entry primary index has no Check vulnerability record, finding,
  or trace. Rejection evidence is sealed at
  `/private/tmp/ply-p7-check-v1-selection.a8e707c.DD2k28`; its fully verified
  18,182-entry manifest SHA-256 is
  `d45958adc2d3be460e85b3f6e379053a8b011af7cd221d4026d9cb0e1bfd57df`,
  and decision-summary SHA-256 is
  `93d0b2833dc78ea2cddb7ea0454713ef96a1ee1879d2ab38692d44bf9e8c7e51`.
- Retain selected `gopkg.in/errgo.v2 v2.1.0` without dependency metadata
  edits. The Go proxy lists exactly three stable releases, v2.0.0, v2.0.1,
  and v2.1.0, with no prereleases or retractions. Proxy `@latest`, exact
  `go list ...@latest`, `@v2`, `@master`, and authoritative gopkg source
  metadata resolve v2.1.0 at 2018-08-23T07:15:05Z. It is canonical latest
  and the highest qualified stable release compatible with Go 1.18.
- V2.1.0 is annotated unsigned tag object
  `635edbc13741bd819966931e8b599c690b4f074d`, peeled merge commit
  `f768c5ab0476c50e978b039312180859c10fe8c0`, tree
  `cb540a0ae1ad20359a13e238d412878b955e6099`. GitHub reports the embedded
  commit signature verified and valid; local cryptographic verification was
  unavailable because `gpg` is absent. Its checksum pair is
  `h1:0vLT13EuvQ0hNvakwLuFZ/jYrLp5F3kcWHXdRggjCE8=` /
  `h1:hNsd1EY+bozCKY1Ytp96fpM3vjJbqLJn88ws8XvfDNI=`. The three proxy
  archives match their exact tag commits byte-for-byte.
- The upstream repository is enabled, unarchived, and not a fork, but the v2
  line has had no commits since 2018. It publishes no GitHub Release objects
  and carries no module deprecation marker. Two post-v2.1.0 v2-branch commits
  remove gocheck and rewrite tests; exact commit resolution produces
  `v2.1.1-0.20180823084306-81e25171e4fb`. That pseudo-version is an
  unreleased branch head, not a stable release and not selected by proxy
  `@latest`, `@master`, or authoritative gopkg metadata, so it is excluded
  from the qualified candidate set.
- V2.1.0 has no Go directive and declares only kr/pretty v0.1.0 and the 2018
  Check pseudo-version. Existing MVS selections already dominate them with
  kr/pretty v0.3.1 at Go 1.12 and the retained 2019 Check pseudo-version with
  no Go directive; selected go-internal v1.9.0 declares Go 1.17. The empty
  changed-selection closure therefore preserves the retained Go 1.18 floor
  by declaration.
- Exact projected `go get gopkg.in/errgo.v2@v2.1.0` exits zero but changes no
  selected version: it only adds a redundant indirect main-module requirement,
  its graph edge, and the already-selected module's full checksum. Modules
  remain 234 and complete packages remain 429, while graph edges project
  3,564 -> 3,565, go.sum 1,031 -> 1,032 lines, and the unapplied tidy diff
  332 -> 344 lines. This metadata-only pin is not an authorized changed
  selection and was not manufactured in the repository.
- Errgo loads in zero complete project packages in either state, repository
  source has zero Errgo imports, and `go mod why -m` says the main module does
  not need it. It is historical unloaded MVS graph debt, selected through an
  old go-internal requirement edge rather than a main or dependency-test
  consumer path.
- Two independent v2.1.0 replays from the proxy archive and exact upstream tag
  verify without source mutation, list two packages, and pass count-1,
  count-10, race, and vet. Their five-module standalone test apparatus is
  separate from the unchanged project closure. Old and projected repository
  build/tests/race/vet, Windows build, pinned lint, CLI surface, and
  byte-identical help/API/CLI reports pass. Empty-HOME count-2 and the final
  full preflight pass all 62 launcher, 80 mutation, and 15 audit-meta
  controls. Earlier runs retain the established primary and nested launcher
  signal-fixture timing diagnostics; direct and Make launcher suites pass
  62/62. One outer Make command-line tool override contaminated a nested lint
  meta-test before the corrected environment-only invocation passed.
- Old and projected govulncheck v1.7.0 outputs and exact IDs remain
  byte-identical at 20/30/20 for Darwin symbol, Darwin module, and Windows
  symbol scans. The fresh 1,392-entry primary vulnerability index contains no
  Errgo record, finding, or trace. No dependency implementation commit or
  post-implementation quality run was required. No-change evidence is sealed
  at `/private/tmp/ply-p7-errgo-v2-selection.1899ff7.HoEmPo`; its fully
  verified 152,333-entry manifest SHA-256 is
  `e9a3c3b1b9c311a1dd96f1804167d57359c2af02615dde2ffb544690dac45ab2`,
  and decision-summary SHA-256 is
  `aa8a30c6f9e261a5c29be3bd59d87cb96e5b07f85582b6f1d55a4f78fcf3a3ec`.
- Retain selected `gopkg.in/resty.v1 v1.12.0` without dependency metadata
  edits. Fresh proxy evidence lists ten v1 semantic tags and resolves
  `@latest`, `@v1`, `@master`, and the exact query to stable v1.12.0 at
  2019-02-28T07:26:48Z, with no prerelease or retraction. V1.7.0 and v1.8.0
  are listed but cannot be downloaded for this path because they declare
  `github.com/go-resty/resty`; v1.0-v1.6 are non-canonical shorthand tags.
  V1.12.0 is canonical latest and the highest exact-path stable release
  compatible with the retained Go 1.18 floor by declaration.
- Selected v1.12.0 is lightweight tag and signed commit
  `fa5875c0caa5c260ab78acec5a244215a730247f`, tree
  `5029acc2e860c8e9495d3b46fdc5d40d429808b8`. GitHub verifies the embedded
  commit signature as valid; there is no separate tag object or tag
  signature, and local verification is unavailable because `gpg` is absent.
  Its checksum pair is `h1:CuXP0Pjfw9rOuY6EP+UvtNvt5DSqHpIxILZKT/quCZI=` /
  `h1:mDo4pnntr5jdWRML875a/NmxYqAlA73dVijT2AXvQQo=`. The 30-file proxy
  archive and exact upstream tag are byte-identical.
- The upstream repository is enabled, unarchived, and not a fork. Its active
  default branch is v3, while the v1 branch ends exactly at v1.12.0. A
  pre-release v1 commit resolves only to the older unreleased pseudo-version
  `v1.11.1-0.20190110224454-0ecc38d58bec`. V2.17.2 and v3.0.0-rc.3 use
  different declared module paths and require Go 1.23; asking the v1 path for
  their commits is rejected for module-path mismatch. They are not v1 upgrade
  candidates.
- V1.12.0 has no Go directive and requires only x/net's 2018 pseudo-version,
  also without a Go directive. The project already selects x/net v0.7.0 at Go
  1.17, so the changed-selection closure is empty and preserves Go 1.18 by
  declaration.
- Exact projected `go get gopkg.in/resty.v1@v1.12.0` exits zero without
  changing any selected version. It only adds a redundant indirect
  requirement, main-module graph edge, and full module checksum. Modules stay
  234 and complete packages stay 429; graph edges project 3,564 -> 3,565,
  go.sum 1,031 -> 1,032 lines, and the unapplied tidy diff 332 -> 344 lines.
  That metadata-only pin is not an authorized selection change and was not
  manufactured in the repository.
- Resty loads in zero complete project packages in both states, repository Go
  source has zero imports, and `go mod why -m` says the main module does not
  need it. The only graph path is historical go.mod debt through
  mvn-pom-mutator v0.2.3, whose source also has zero Resty imports.
- The exact proxy module source verifies and does not mutate, and its
  standalone two-module apparatus needs no test-only requirement. Count-1 and
  full count-2 tests fail the same six legacy error-string assertions under Go
  1.26.7; race reproduces those assertions without a data-race diagnostic,
  while vet passes. The mandatory module stop rule rejects even the redundant
  projection. Old and projected repository build/tests/race/vet, Windows
  build, pinned lint, help/API/CLI identity, empty-HOME count-2, direct 62-check
  launcher suite, and final preflight all pass. The final preflight passes 62
  launcher checks, 80 mutation controls, and 15 audit meta-controls. No
  dependency implementation or post-selection quality/snapshot/Docker run was
  required.
- Govulncheck v1.7.0 preserves byte-identical old/projected
  Darwin-symbol/Darwin-module/Windows-symbol populations at 20/30/20. The
  fresh 1,392-entry primary index has no exact `gopkg.in/resty.v1` record; its
  only Resty entry is GO-2023-2328 for the alternate v2 module path. Retention
  evidence is sealed at
  `/private/tmp/ply-p7-resty-v1-selection.a855007.Gg7QC9`; its fully verified
  110,457-entry manifest SHA-256 is
  `72d739b05c507fab058737de6ab967b534e5fcdab1c9588582188532f1a926c2`,
  and decision-summary SHA-256 is
  `008b8b9396b58c93e696a5d495037ac974224149fce9b34badb6c7786723a109`.
- Retain selected `gopkg.in/alecthomas/kingpin.v2 v2.2.6` without dependency
  metadata edits. The Go proxy lists exactly 40 stable v2 semantic versions,
  with no prereleases or retractions, and resolves `@latest`, `@v2`, and the
  authoritative gopkg `@master` to v2.4.0 at 2023-09-30T22:59:49Z. V2.4.0
  is therefore canonical latest as a stable tag, but it declares
  `github.com/alecthomas/kingpin/v2`; exact gopkg-path `go get` rejects it.
  V2.3.0 has an invalid non-`.v2` GitHub module declaration, and v2.3.1,
  v2.3.2, and v2.4.0 use the alternate GitHub `/v2` path. Proxy module
  metadata for the 36 stable tags through v2.2.6 declares the exact gopkg
  path. Selected v2.2.6 is the highest qualified stable exact-path release and
  imposes no declaration above the retained Go 1.18 floor.
- Selected v2.2.6 is lightweight tag and unsigned commit
  `947dcec5ba9c011838740e680966fd7087a71d0d`, tree
  `7b5bf03129479f410ca39ea2647dd190bddee927`. Its checksum pair is
  `h1:jMFz6MfLP0/4fUyZle81rXUoxOBFi19VUFKVDOQfozc=` /
  `h1:FMv+mEhP44yOT+4EoQTLFTRgOQ1FBLkstjWtayDeSgw=`; proxy ZIP SHA-256 is
  `638080591aefe7d2642f2575b627d534c692606f02ea54ba89f42db112ba8839`.
  All 39 proxy files match the exact tag commit. All 40 v2 tags are
  lightweight commit refs without tag signatures; local OpenPGP status is
  `N` for every tag commit. GitHub reports v2.2.6 unsigned and v2.4.0's
  commit signature valid.
- Exact-path pseudo-versions at pre-module commit `c2ca6a1e4f86`, module
  conversion commit `102f372a17d4`, and divergent `v3-unstable` head
  `95d230a53780` resolve and download, but are unreleased commits rather than
  stable candidates; the last is also below the selected semantic line.
  Current GitHub master `177e1b9ba430` resolves to unreleased
  `v2.4.1-0.20260824075000-177e1b9ba430` on the alternate module path and is
  rejected by exact gopkg-path get. The unarchived, undisabled, non-fork
  repository now says `CONTRIBUTIONS ONLY` and names the GitHub `/v2` module
  as current stable; the exact gopkg module itself has no deprecation marker.
- Exact projected
  `go get gopkg.in/alecthomas/kingpin.v2@v2.2.6` changes no selected version.
  It only adds redundant indirect requirements on the already-selected
  template, units, and Kingpin modules, their three main-module graph edges,
  and their three existing full checksums. Modules remain 234 and complete
  packages remain 429; graph edges project 3,564 -> 3,567, go.sum lines
  1,031 -> 1,034, and the unapplied tidy diff 332 -> 353 lines. The changed
  selection closure is empty, so this metadata-only projection is not an
  authorized implementation and was not manufactured in the repository.
- Kingpin loads in zero complete project packages, has zero repository Go
  imports, and `go mod why -m` says the main module does not need it. It is
  historical graph debt through old Prometheus tsdb/common requirements,
  reached from mvn-pom-mutator rather than a loaded project consumer path.
- The exact proxy source and upstream tag remain source-identical after test.
  Their separate seven-module test apparatus is not part of project MVS.
  Count-1, ten fresh count-1 processes, race, and vet pass, but both exact
  sources fail `go test ./... -count=10`: repetitions 2-10 fail
  `TestRequiredArgWithEnvarMissingErrors` and
  `TestRequiredWithEnvarMissingErrors`, 18 failures total, because sibling
  tests leave `TEST_ARG_ENVAR` and `TEST_ENVAR` set in process-global state.
  The mandatory repeated module-suite gate therefore independently rejects
  even the redundant exact-selected projection.
- Old and projected repository verification, build, complete tests/race/vet,
  Windows build, pinned lint, CLI surface, API/CLI reports, and public help
  pass or remain byte-identical. Empty-HOME count-2 passes. A corrected clean
  committed projection and external compatibility reports produce a zero-exit
  full preflight with all 62 launcher controls, 80/80 killed mutants, and 15
  audit meta-controls. Five aggregate attempts and one isolated attempt retain
  the established signal-retention timing failure; a later attempt also
  proved why compatibility reports must stay external before the unchanged
  corrected run passed. Changed-selection quality, snapshot/Docker, and
  focused/full audits remain inapplicable because no selection changed and the
  mandatory module gate failed.
- Govulncheck v1.7.0 preserves byte-identical old/projected outputs and exact
  Darwin-symbol/Darwin-module/Windows-symbol ID populations at 20/30/20. The
  fresh 1,392-entry primary index has no Kingpin record, finding, or trace.
  No-change evidence is sealed at
  `/private/tmp/ply-p7-kingpin-v2-selection.c440130.2WGm6X`; its fully verified
  97,185-entry manifest SHA-256 is
  `ce50700fe0df4e208ce1a679d7ef8c9b97fdca45398843d392351673fa70837c`,
  and decision-summary SHA-256 is
  `d9080ecc8983da371aacfe1ef234a1290e7e577738c4b5dc9d99078365bba54f`.
- Accept selected `github.com/alecthomas/units`
  `v0.0.0-20190717042225-c3de453c63f4` -> canonical-latest
  `v0.0.0-20240927000941-0f3dac36c52b`. The proxy lists no semantic or
  prerelease tags; proxy `@latest` and exact Go `@master` agree on default-
  branch commit `0f3dac36c52b29c22285af9a6e6593035dadd74c` at
  2024-09-27T00:09:41Z. It is an unreleased pseudo-version, not a stable
  release. The enabled, unarchived, non-fork repository has no tags, releases,
  retractions, alternate path, or deprecation marker. A newer Renovate branch
  and unmerged pull-request ref are not qualified exact-path latest versions.
- Candidate `go.mod` declares Go 1.15 and requires only already-selected
  testify v1.9.0 at Go 1.17. Proxy/upstream source identity is byte-exact;
  both commits have valid GitHub signature verification. Candidate count-1,
  count-10, race, and vet pass in both sources. The old source also passes
  with its exact separately recorded test-only assert closure, which is not
  project MVS closure.
- Exact get changes one selection, adds only main -> Units and Units ->
  already-selected testify graph edges, and adds the candidate checksum pair.
  Measurements move 234 -> 234 selected modules, 3,564 -> 3,566 edges,
  429 -> 429 packages, 1,031 -> 1,033 go.sum lines, and 332 -> 341 unapplied
  tidy-diff lines. Units has zero loaded packages, imports, or consumer paths;
  it is historical unloaded MVS graph debt.
- Dependency-only commit `1c874fda104baa36c3a32aaef2d2d2be9689a942`,
  exact parent `ae5e731121c17338d08dba4c74a99a1c95a36c49`, tree
  `744d61deae7d7c3dc25fe1fd1e13a1fc94798981`, changes only `go.mod` and
  `go.sum`. Repository build/tests/race/vet, Windows build, pinned lint,
  empty-HOME count-2, and byte-identical help/API/CLI reports pass. Primary
  vulnerability evidence has no Units record and preserves exact 20/30/20
  populations and traces.
- Exact `make quality` exits 0 at `/private/tmp/ply-p7-uq11`; its verified
  237,860-entry manifest SHA-256 is
  `e8b4fdc59ec12c42afa21286939a1a5ca28a4fef29c88655f73890b580564328`.
  All 27 Q0-Q2 rows attain L2, 80/80 mutations are killed, host/snapshot/
  Docker acceptance passes, and held, regressed, not-comparable, and dirty
  counts are zero. Full audit exits expected 1, never 2, only for established
  queued Q3.1, Q3.3, Q3.4, and Q3.7. Selection evidence is sealed at
  `/private/tmp/ply-p7-units-selection.ae5e731.smI4xp`, 72,950 entries and
  manifest SHA-256
  `2cdb64e74a40304c86ef35f6b8e479c49190d4c2d0e6678770bee4e6486493ec`;
  decision-summary SHA-256 is
  `e30405aeac0eecf6356570f8a49677de1715ad9a7aee15e8efdc9739f55c4cd0`.
- Accept selected `github.com/alecthomas/assert`
  `v0.0.0-20170929043011-405dbfeb8e38` -> canonical-latest stable v1.0.0.
  The exact-path proxy lists only this stable release, and exact Go `@latest`,
  `@v1`, and `@master` agree. It is lightweight tag and unsigned commit
  `73444aca37e09619baa21040dad4527f857f9abb`, tree
  `d09edea32b1be9b1e73428797782d491949b50b7`, at
  2021-12-01T05:52:01Z. The enabled, unarchived, non-fork repository has no
  v1 prerelease, retraction, deprecation marker, or GitHub Release object.
  Stable v2.11.0 and a newer unreleased master pseudo-version belong to the
  alternate `/v2` module path and are not exact-v1 candidates.
- Candidate v1.0.0 declares Go 1.17 and requires colour v0.1.0, repr's 2021
  pseudo-version, and go-diff v1.2.0 plus lower isatty/x/sys versions. Exact
  get changes those four selections. Their selected declarations are absent,
  Go 1.15, Go 1.12, Go 1.15, and Go 1.18 or lower, so the complete closure
  preserves the retained Go 1.18 floor. All four source archives match exact
  upstream commits after documented proxy normalization for tracked symlinks.
  Assert v1.0.0 changes module/Hermit metadata but no production Go source
  relative to the old pseudo-version.
- The minimal closure keeps 234 modules and 429 packages, changes graph edges
  3,566 -> 3,580 through 14 additions and no removals, adds exactly four
  checksum pairs for go.sum 1,033 -> 1,041, and changes the unapplied tidy
  projection 341 -> 354 lines. The changed modules load in zero old/new
  complete packages; assert has no repository imports and `go mod why` says
  the main module does not need it. It is historical unloaded graph debt
  through go-term-markdown's old Chroma requirement.
- Candidate proxy and upstream source pass verify/list/count-1/count-10/race/
  vet without source mutation. The sole native package has no test files;
  `_example` is excluded by `./...`, and no separate test-only apparatus is
  needed. Old, projected, and committed repository build/tests/race/vet,
  Windows build, pinned lint, empty-HOME count-2, and help/API/CLI contracts
  pass with byte-identical outputs. Primary vulnerability evidence has no
  assert/closure record and preserves exact 20/30/20 IDs and traces.
- Exact get produced dependency-only commit
  `0781fd6cdb625623c6d726743bf58352113d0ccb`, parent
  `8f24778c9ddafdbf42583595400093450e4112ec`, tree
  `690e939eaacc57af08868c7f86ce3fec5dbcf034`, changing only `go.mod` and
  `go.sum`. Exact `make quality` exits 0 with all 27 Q0-Q2 rows at L2, 80/80
  mutations killed, fresh host/snapshot/Docker acceptance, and zero held,
  regressed, not-comparable, or dirty counts. Standalone audit meta and the
  focused audit pass; full audit exits expected 1 only for queued Q3.1, Q3.3,
  Q3.4, and Q3.7.
- Assert selection evidence is sealed at
  `/private/tmp/ply-p7-assert-selection.8f24778.fEzhH8`, 46,700 entries and
  manifest SHA-256
  `072c8a42889b85e1fa86d0e5701db7f31f8840c255f06de5a7633403db621782`;
  decision-summary SHA-256 is
  `bb13fa01a123a32399806db46f7e61813e6549b7ca1c5a42c25a09254684d17f`.
  Exact quality evidence at `/private/tmp/ply-p7-assert-quality.0781fd6.q1`
  has 237,950 entries and verified manifest SHA-256
  `04b202f6084740762589a042c8847c2411da5a54ee6572d3fa54790c4a3301e1`.
- Accept selected `github.com/alecthomas/repr`
  `v0.0.0-20210801044451-80ca428c5142` -> canonical-latest stable v0.5.4.
  Proxy `@latest` and exact Go `@latest`, `@v0`, and `@master` agree on
  lightweight unsigned tag commit `9b3680b9bb4e172c4fe539347ac5d0ddf5de0aa9`,
  tree `d5c923ba26338b41955d8fc02e88948c460777a0`, at
  2026-07-15T12:04:01Z. The exact path has ten stable tags, no prerelease,
  retraction, deprecation, alternate path, or GitHub Release object. A newer
  Renovate head is an unreleased non-default-branch pseudo-version, not the
  exact-path latest decision.
- Candidate v0.5.4 declares Go 1.18 and has no requirements; the old version
  declares Go 1.15 and also has none. Proxy/checksum-database identities agree,
  and proxy/upstream regular files match after accounting for five tracked
  symlinks omitted by the proxy. Both source forms pass verify, list, count-1,
  count-10, race, and vet without mutation and need no test-only apparatus.
- Exact get changes only Repr's selection and relabels one existing main graph
  edge. Measurements remain 234 modules, 3,580 edges, and 429 packages;
  go.sum moves 1,041 -> 1,043 with exactly the candidate checksum pair and the
  unapplied tidy projection moves 354 -> 356 lines. Repr is absent from the
  complete dependency-test population and repository imports; it is unloaded
  historical MVS graph debt.
- Dependency-only commit `6f4d02eb9c86ec2df8a488a85ed973aababe1f38`,
  exact parent `53475076e1c79f6d2181877e3238a1a5d389c246`, tree
  `82e7b1f5c503659082207481b8339e8113e38c10`, changes only `go.mod` and
  `go.sum`. Repository build/tests/race/vet, Windows build, pinned lint,
  empty-HOME count-2, help/API/CLI, launcher/Make/preflight, host/snapshot/
  Docker, and exact Q0-Q2 quality gates pass. All 27 rows attain L2, 80/80
  mutations are killed, and held, regressed, not-comparable, and dirty counts
  are zero. Primary old/candidate vulnerability IDs and traces are identical
  at exact 20/30/20 populations.
- Repr evidence has 447 verified entries; manifest SHA-256 is
  `aba1a6ac45d6b8a51a664e0d11f40ccccf82017240c91beaa59ec36866e26711`
  and decision-summary SHA-256 is
  `900e96c7001a89857ca7c16e1f2572ea0b90db9f2ab4d4e4e6a1c3153808bbcf`.
  Preserve its separately warmed caches, BSD-mktemp adapter, Docker-config,
  and scratch-relative manifest corrections. Historical golangci-lint and
  govulncheck binary hashes are nonportable build receipts; pinned fresh
  binaries were independently obtained/rebuilt and passed every gate.
- Reject selected `github.com/alecthomas/kong`
  `v0.2.1-0.20190708041108-0548c6b1afae` -> v1.3.0 even though v1.3.0 is the
  highest stable release whose complete changed closure preserves the retained
  Go 1.18 floor. The proxy lists 75 stable tags and no prerelease or
  retraction. Canonical-latest v1.16.1 and every release from v1.4.0 declare
  Go 1.20; newer master resolves only as unreleased Go-1.20 pseudo-version
  `v1.16.2-0.20260828071222-a60008c6dae2`.
- Candidate v1.3.0 is lightweight unsigned tag commit
  `7bbb0b76ada1610f18cf71c54cca74209da88bd8`, tree
  `96a61d79e4745b05738a01906c43a7b8e5820b5a`, at
  2024-11-01T01:25:41Z. Its checksum pair is
  `h1:YJKuU6/TV2XOBtymafSeuzDvLAFR8cYMZiXVNLhAO6g=` /
  `h1:IDc8HyiouDdpdiEiY81iaEJM8rSIW6LzX8On4FCO0bE=`. All 64 proxy regular
  files match the exact tag after accounting for five omitted Hermit symlinks
  and the seven-file nested example module.
- Exact projection adds assert/v2 v2.11.0 and gotextdiff v1.0.3 while accepted
  Repr v0.5.4 dominates the candidate's v0.4.0 request. Their declarations are
  Go 1.18, Go 1.16, and Go 1.18. Modules project 234 -> 236, graph edges 3,580
  -> 3,584, complete packages remain 429, go.sum projects 1,043 -> 1,045 only
  for Kong's checksum pair, and the unapplied tidy diff grows 356 -> 359 lines.
  Kong remains absent from all complete packages and repository imports;
  `go mod why` says the main module does not need it.
- Proxy and exact-tag candidate sources independently verify and pass complete
  count-1, count-10, and race tests without mutation. Mandatory vet fails in
  both with byte-identical 77-line malformed-struct-tag diagnostics across
  four test files. This stop-rule failure makes downstream repository,
  snapshot/Docker, exact-quality, and audit gates inapplicable. No dependency
  metadata or implementation commit was made.
- Fresh primary vulnerability evidence has no Kong record and preserves exact
  byte-identical old/candidate 20/30/20 Darwin-symbol/Darwin-module/Windows-
  symbol IDs and traces. Kong decision evidence has 211 verified entries;
  manifest SHA-256 is
  `6148c409545f453a78ffdc3a0a10934b75e082c1681f2805c29c8253a23d01ae`,
  and decision-summary SHA-256 is
  `028c0fd0025a5f5de433b65d63a5258e18c6e959b440446e9834b03925f98b36`.
- Retain selected exact-path `github.com/alecthomas/chroma v0.10.0`. The
  proxy lists 30 stable v0 tags, no prerelease or retraction, and exact
  `@latest`/`@v0` resolve selected v0.10.0 at 2022-01-12T10:49:38Z. It is a
  lightweight tag at unsigned commit
  `36bdd4b98823bd1d7be96767cde3dd575e60b406`, tree
  `d27058989b845352d49b45e7538f7b0004c8d651`, and has a GitHub Release
  object. Exact-path `@master` resolves only the later unreleased pseudo-
  version `v0.10.1-0.20220126230913-d491f1b5c1d2`. The following commit
  changes the module path to `/v2`; current `/v2` latest v2.27.0 declares Go
  1.25, while current repository master is the distinct `/v3` alpha lineage.
  Neither alternate module path is an in-place v0 update.
- V0.10.0 declares Go 1.13. Its standalone test graph selects eight modules
  across nine edges and has no declaration above Go 1.13. All 696 proxy
  regular files match its exact tag; eight nested-module files and seven
  tracked symlinks are correctly omitted. Its checksum pair is
  `h1:7XDcGkCQopCNKjZHfYrNLraA+M7e0fMiJ/Mfikbfjek=` /
  `h1:jtJATyUxlIORhUOFNA9NZDWGAQ8wpxQQqNSB4rjA/1s=` and proxy ZIP SHA-256 is
  `beb07b996ee33bc052fe039c93d1c0726e61bcc4819ca39f7bf63304f2ae8c49`.
- Exact selected-version get is a true no-op. Old and replay states remain
  234 modules, 3,580 edges, 429 complete packages, 1,043 go.sum lines, and a
  356-line unapplied tidy projection. The empty changed-selection closure
  adds no version, edge, checksum, or Go-floor requirement, so no dependency
  implementation was manufactured.
- Exactly 33 Chroma packages load through `plybuild/cmd -> go-term-markdown ->
  chroma`; repository source has no direct import. The consumer uses Chroma's
  lexer selection/coalescing, tokenization, Pygments style, and fallback/TTY8
  formatting symbols for code blocks. The project Markdown contract passes
  ten times, but the direct consumer's v0.7.1-era code-block golden fails all
  ten project-MVS repetitions because current selected Chroma/color output
  changes ANSI intensity and padding.
- Proxy and exact-tag sources independently pass verify, list, complete
  count-1/count-10, and race without mutation. Mandatory vet fails in both
  with the same normalized 7,175 unkeyed-`Rule` diagnostics across 206 files;
  sorted diagnostic SHA-256 is
  `cbfbdc99bce5e357f83b20eec1a08644ddb86a16ed90f3a2c176cb501a68604e`.
  The module and loaded-behavior stop rules make changed-selection quality,
  snapshot/Docker, and audit gates inapplicable.
- Fresh primary vulnerability evidence has no Chroma record and preserves
  exact old/replay 20/30/20 Darwin-symbol/Darwin-module/Windows-symbol IDs and
  normalized traces. Chroma decision evidence has 4,335 verified entries;
  manifest SHA-256 is
  `bbd6a41e0cff3248bacc2fef6630ed3679a3299484cb21446fc2841a39a742d5`,
  and decision-summary SHA-256 is
  `476f926b06cd8fbfce8ac22c9e0ed1102c8f5c8f8e2070bcb1636a825225927f`.
- Retain selected exact-path `github.com/alecthomas/colour v0.1.0`. It is the
  only stable proxy version, and proxy `@latest` plus exact Go `@latest`,
  `@v0`, `@master`, and `@v0.1.0` all resolve it at
  2019-11-01T02:47:59Z. The lightweight tag is unsigned commit
  `a1c6bd85eba7190e4d2959ecd15831d0a25b37b9`, tree
  `ef4bddf202f91763747ee4828270fce2741373ca`, parent
  `60882d9e27213e8552dcff6328914fe4c2b44bc9`. The public, enabled,
  unarchived, undisabled, non-fork repository has one branch, one tag, 12
  commits, no newer default-branch commit, and no GitHub Release object.
  Stable qualification comes from the sole semver tag and proxy publication;
  there is no prerelease, retraction, or authoritative alternate module path.
- The historical pseudo-version
  `v0.0.0-20160524082231-60882d9e2721` is the exact parent and not a newer
  candidate. V0.1.0 adds `^S` strikethrough. Its checksum pair is
  `h1:nOE9rJm6dsZ66RGWYSFrXw461ZIt9A6+nHgL7FRrDUk=` /
  `h1:QO9JBoKquHd+jz9nshCh40fOfO+JzsoXy8qTHF68zU0=` and agrees with the
  checksum database. Proxy ZIP SHA-256 is
  `74d51002731fa104943b62ee11fb61b14c517e75a4a3983bfb03976b6c75349b`;
  all five regular files match the exact tag without omission.
- The upstream tag predates `go.mod`; proxy-synthesized metadata declares
  only the exact path, with no Go version or requirements, so the declared
  complete closure is Colour alone. Identical isolated Go-1.18 source-test
  apparatus selects imported go-isatty v0.0.20 at Go 1.15 and x/sys v0.6.0
  at Go 1.17. That apparatus is separate from project MVS, which already
  selects go-isatty v0.0.20 and x/sys v0.30.0 at no more than Go 1.18.
- Exact selected-version get is byte-identical. Old/candidate states remain
  234 modules, 3,580 graph edges, 429 complete packages, 1,043 go.sum lines,
  and a 356-line unapplied tidy projection. Module, graph, package, go.mod,
  go.sum, and tidy diffs are empty; the changed-selection closure adds no
  version, edge, checksum, or floor requirement. No dependency implementation
  commit was created.
- Colour loads in zero project packages, repository source has no Colour
  import, and `go mod why` says the main module does not need it. Existing
  edges are main -> Colour v0.1.0 and Assert v1.0.0 -> Colour v0.1.0; old
  Chroma v0.7.1 requests the parent pseudo-version, which MVS already
  upgrades. An explicitly external project-MVS fixture passes count-10/race/
  vet for formatter, stripper, forced ANSI, nonterminal TTY, string-printer,
  reset, underline, strikethrough, and escaped-caret behavior.
- Proxy and exact-tag copies independently verify/list and pass complete
  count-1/count-10/race/vet without source mutation. Repository verify/build/
  tests/race/vet, Windows build, pinned lint, byte-identical public help and
  API/CLI reports, API/CLI compatibility, and CLI surface pass. Full preflight
  passes 62 launcher checks, distribution/lint/install/toolchain contracts,
  snapshot/Docker meta-contracts, 80/80 mutations, verification meta-tests,
  and 15 audit controls. Changed-selection host/snapshot/Docker/quality runs
  are inapplicable because no selection changed.
- Fresh primary vulnerability evidence has 1,392 module records and no Colour
  record or trace. Old/candidate normalized findings and traces are identical
  at exact 20/30/20 Darwin-symbol/Darwin-module/Windows-symbol populations.
  Colour decision evidence has 322 verified entries; manifest SHA-256 is
  `8277a8e7d5b4810a769547d2adbb5e0779d4e52c09b6b8e30c209383fc718543`,
  and decision-summary SHA-256 is
  `6932da454c33f6580d12349122346e1441e7277c81264ed091c7322880132592`.
  Preserve the pre-module source-test apparatus, separately warmed v1.0.1
  compatibility cache, and scratch-only BSD-mktemp/preflight correction.
- Retain selected exact-path
  `github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751`
  without a dependency edit. The exact proxy list has no stable or prerelease
  version; `@latest` and `@master` select this default-branch pseudo-version at
  2019-07-18T01:26:54Z, while `@v0` has no match. The public, enabled,
  unarchived, non-fork repository has one branch, zero tags, zero Releases,
  and no deprecation marker. The only newer source ref is a closed, unmerged
  pull request. `/v2`, `/v3`, and gopkg.in alternate probes find no module.
- Selected commit `fb15b899a75114aa79cc930e33c46b577cc664b1`, tree
  `9658e953ba71f92dcf44f2d39cc5f90a27a0b88b`, is the master head. Its
  embedded GitHub web-flow signature is valid and independently verifies with
  fingerprint `5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`; there is no tag signature.
  All 22 proxy regular files match upstream. Selected ZIP SHA-256 is
  `25e3be7192932d130d0af31ce5bcddae887647ba4afcfb32009c3b9b79dbbdb3`;
  checksum pair is
  `h1:JYp7IbQjafoB+tBA3gMyHYHrpOtNuDiK/uB5uXxq5wM=` /
  `h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=` and sumdb agrees.
- Raw metadata declares only the exact module path, with no Go directive or
  requirements. The complete declared closure is Template alone and cannot
  raise Go 1.18. Go's implicit 1.16 version when the source is tested as a
  main module is not an upstream declaration. The selected commit only adds
  `go.mod` relative to the historical 2016 pseudo-version; source and tests
  are unchanged.
- Exact selected-version get changes no selection. It projects 234 -> 234
  modules, 3,580 -> 3,581 edges, 429 -> 429 packages, 1,043 -> 1,044 sum
  lines, and 356 -> 358 tidy-diff lines solely by adding a redundant indirect
  requirement, a main edge, and the full checksum; tidy removes the
  requirement. Template has zero loaded project packages and repository
  imports, and `go mod why` says the main module does not need it. Its three
  historical graph requests arrive through Prometheus common/tsdb modules.
- Proxy and exact-commit sources independently verify/list and remain
  unchanged. Mandatory count-1, count-10, and race-enabled complete tests fail
  identically in `TestJSEscaping`: current `unicode.IsPrint` emits U+FDFF
  literally while the old test expects an escape. Vet passes. This stop-rule
  failure makes repository, compatibility, snapshot/Docker, quality, and audit
  gates for a changed selection inapplicable, and no dependency implementation
  commit was created.
- Fresh primary vulnerability evidence has 1,392 module records and no
  Template record or trace. Old/candidate exact findings and traces remain
  identical at 20/30/20 Darwin-symbol/Darwin-module/Windows-symbol
  populations. Template evidence has 1,702 verified entries; manifest SHA-256
  is `2bfa08ca73085154ab5f0833814872efe481088703e6d7480edf1e2a40853068`,
  and decision-summary SHA-256 is
  `4acc26e68fa5ad432a44d64ade0d399464f9b07eb2c6244b29711e048ae24170`.
- Retain selected exact-path `github.com/antihax/optional v1.0.0` without a
  dependency edit. It is the sole stable proxy version, canonical `@latest`,
  and highest stable Go-1.18-floor-compatible candidate. The public, enabled,
  unarchived, non-fork repository has one branch, one lightweight tag, and one
  non-draft, non-prerelease GitHub Release. The tag is commit
  `c3f0ba9c1a592b971d66b2787679af55b5c58f21`, tree
  `b9328a8aa4526004bb36928dbc136c3acb8eec3a`, at
  2019-10-10T23:37:20Z. The lightweight tag is unsigned; the associated
  GitHub web-flow commit signature independently verifies with fingerprint
  `5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`.
- Exact `@master` is only unreleased pseudo-version
  `v1.0.1-0.20220101210036-407d38fabb55`, commit
  `407d38fabb5592e58b0841e8adb93edd65ee1319`, tree
  `3dffdae3ba5f11e0b140c2edd085e5d1a69df40d`, at
  2022-01-01T21:00:36Z. Its two post-tag commits change only `README.md`; all
  21 Go files and `go.mod` are identical to v1.0.0. There are no prereleases,
  retractions, or authoritative `/v2`, `/v3`, gopkg.in, or renamed paths.
- V1.0.0 declares Go 1.13 and has no requirements, so its complete declared
  closure is Optional alone. Proxy and exact-tag forms match at all 23 regular
  files and independently pass verify/list, complete count-1/count-10/race
  tests, and vet without mutation. The package has no test files and no test-
  only requirements. Its checksum pair is
  `h1:xK2lYat7ZLaVVcIuj82J8kIro4V6kDe0AUDFboUCwcg=` /
  `h1:uupD/76wgC+ih3iEmQUL+0Ugr19nfwCT1kdvxnR2qWY=` and sumdb agrees.
- Exact selected-version get changes no selection. It projects 234 -> 234
  modules, 3,580 -> 3,581 edges, 429 -> 429 packages, 1,043 -> 1,044 sum
  lines, and 356 -> 358 tidy-diff lines only by adding a redundant indirect
  requirement, main edge, and full checksum. Optional loads in zero project
  packages, Ply has no import, and `go mod why` says the main module does not
  need it. The original edge is grpc-gateway v1.16.0 -> Optional v1.0.0; its
  generated examples and an external Go-1.18 fixture cover the seven wrapper
  types and their used methods without mislabeling them as loaded Ply code.
- Untouched/projected repository verify/build/count-1/count-10/race/vet,
  Windows build, pinned lint, byte-identical help and API/CLI reports,
  compatibility, CLI surface, and clean full preflight pass. Fresh primary
  vulnerability evidence has no Optional record and preserves exact 20/30/20
  old/projected IDs and traces. No dependency implementation commit was
  created. Optional evidence has 540 verified entries; manifest SHA-256 is
  `345a29c62bfadf332f14aa234b1265a9dc9142739d212534d83ace531fb39afa`,
  and decision-summary SHA-256 is
  `945b7ca8fe9d9faa605a6a80d8d416655925f95cd727071babf813bcd4caa27b`.
  Preserve its compatibility-cache, source/signature, checksum-delta,
  vulnerability-order, and scratch-local BSD-mktemp/preflight corrections.
- Advance exact-path `github.com/armon/circbuf` from selected pseudo-version
  `v0.0.0-20150827004946-bbbad097214e` to canonical latest and master-head
  pseudo-version `v0.0.0-20190214190532-5111143e8da2`. The exact stable list is
  empty and the public, enabled, unarchived, non-fork repository has one
  branch, nine commits, zero tags, and zero Releases. There are no
  prereleases, retractions, deprecation markers, or authoritative alternate
  module paths. This is an unreleased pseudo-version, not a stable release.
- The old unsigned commit is
  `bbbad097214e2918d8543d5201d12bfd7bca254d` at
  2015-08-27T00:49:46Z. The candidate is commit
  `5111143e8da2e98b4ea6a8f32b9065ea1821c191`, tree
  `2ab2d9cf2632ab7b549f7da7f081dbe868a697db`, at
  2019-02-14T19:05:32Z. Its GitHub web-flow signature independently verifies
  with fingerprint
  `5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`; the key is currently expired.
  The only source delta is a one-line `go.mod`; every Go source and test file
  is identical.
- Neither version declares a Go version or requirements. The candidate's
  complete minimal closure is Circbuf alone and preserves Go 1.18. Proxy and
  exact-commit forms independently pass verify/list, count-1/count-10/race
  tests, and vet without mutation; exact Go 1.18.10 also passes count-10/race/
  vet. Historical/current Serf consumers and an external Go-1.18 fixture cover
  `NewBuffer`, `Write`, `TotalWritten`, `Size`, `String`, and `Bytes`. Circbuf
  loads in zero Ply packages and `go mod why` says the main module does not
  need it.
- Exact candidate get changes only Circbuf's version and adds an indirect
  requirement, one main graph edge, and the candidate checksum pair. Project
  measurements change 234 -> 234 modules, 3,580 -> 3,581 edges, 429 -> 429
  packages, 41 -> 41 loaded modules, 1,043 -> 1,045 sum lines, and 356 -> 361
  tidy-diff lines. Tidy removes the explicit pin and its two checksum lines;
  inherited debt is unchanged. Implementation commit is
  `3be2183ee310ccdc358ce4ed372c0785de25b88b`, parent
  `7a0ca4e2caba6d2fff20a9c169b181f9187b40dc`, tree
  `a51846840b5f252a30359530dfc950f811398431`, changing only `go.mod` and
  `go.sum`.
- Candidate/repository/consumer, Windows, pinned-lint, help/API/CLI,
  compatibility, preflight, host/snapshot/Docker, and empty-HOME gates pass.
  Exact quality exits 0 with all 27 Q0-Q2 rows at L2, 80/80 mutations killed,
  and zero held/regressed/not-comparable/dirty counts. Full audit exits 1 only
  for queued L3 rows. Fresh primary vulnerability results have no Circbuf
  record or trace and preserve exact 20/30/20 old/candidate populations.
  Circbuf evidence has 430 verified entries; manifest SHA-256 is
  `2684f8c269376320bcfca3404bfbf10c638c81753e6af36b5be59c4df2370d62`
  and decision-summary SHA-256 is
  `17441163b7924e2cd61cb167b739562da7672fa5abfa10089da36ef906f07ff9`.
- Retain exact-path `github.com/armon/consul-api` at already-selected canonical
  latest and `master` head pseudo-version
  `v0.0.0-20180202201655-eb2c6b5be1b6`, commit
  `eb2c6b5be1b66bab83016e0b05f01b8d5496ffbd`, tree
  `aeb2299aaf107d0823ce91f057798119b821e81b`, at
  2018-02-02T20:16:55Z. The exact stable proxy list is empty and `@v0` has no
  match. The public, enabled, unarchived, non-fork repository has one branch,
  46 commits, zero tags, and zero Releases. This is an unreleased pseudo-
  version, not a stable release. Its GitHub web-flow commit signature is
  cryptographically valid for fingerprint
  `5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`; the key is currently expired.
- All 21 regular proxy files match the exact commit archive and Git tree at
  normalized manifest SHA-256
  `358c1cfd7c63c682fc02e069ef0ab20f494bc44d9789e62ca39e62fc4026cbee`.
  The checksum pair is
  `h1:G1bPvciwNyF7IUmKXNt9Ak3m6u9DE1rF+RmtIkBpVdA=` /
  `h1:grANhF5doyWs3UAsr3K4I6qtAmlQcZDesFNEHPZAzj8=` and sumdb agrees. The
  synthetic module declares neither a Go version nor requirements, so its
  complete declared minimal closure is Consul API alone and preserves Go 1.18.
  Its README explicitly deprecates this source for distinct path
  `github.com/hashicorp/consul/api`; that is an out-of-scope migration whose
  current v1.34.4 latest declares Go 1.26, not an exact-path upgrade.
- Exact Go 1.26.7 proxy and commit forms verify/list but identically fail
  complete count-1/count-10/race tests with 33/330/33 refused connections to
  an external Consul agent at 127.0.0.1:8500. Vet also fails on two
  `testing.T.Fatalf` calls from non-test goroutines. Exact Go 1.18.10 compiles
  the package but reproduces those failures. A historical Crypt consumer and
  external Go-1.18 HTTP fixture exercise the actual configuration, client, KV,
  query, metadata, key, and value symbols; the fixture passes count-10/race/
  vet. Consul API loads in zero Ply packages and Ply imports none of it.
- Exact selected-version get changes no selection and projects only a redundant
  indirect requirement, one main graph edge, and one full checksum, all removed
  by tidy. No metadata or dependency implementation commit was created.
  Accepted measurements remain 234 modules, 3,581 edges, 429 packages, 41
  loaded modules, 1,045 sum lines, 361 tidy-diff lines, and exactly 29 added
  checksum lines relative to accepted go-cmp commit c314bcb. Repository/full
  quality gates were inapplicable after the no-change result and dependency
  stop-rule failures.
- Fresh primary vulnerability data has no Consul API record or trace and
  preserves identical old/projected 20-ID/22-trace Darwin symbol, 30-ID Darwin
  module, and 20-ID/22-trace Windows symbol populations. Consul API evidence
  has 315 verified entries; manifest SHA-256 is
  `1dbe7063e72dade7bc31e9c8967da78a60ef97859b68562ffa1a07b75a0b3b0b`
  and decision-summary SHA-256 is
  `0f5413813a179949c2dbce2a029becfceb9f650f7a603ac8b19755c3f4d48733`.
- Retain exact-path `github.com/armon/go-metrics v0.4.0` without dependency
  metadata edits. The proxy lists 22 stable versions and no prereleases;
  `@latest` and `@v0` resolve v0.6.1, while `@master` resolves unreleased
  pseudo-version `v0.6.2-0.20260907064447-465585286d74`. Both declare Go
  1.25.0 and `module github.com/hashicorp/go-metrics`, so neither is an exact-
  path or retained-floor candidate. Proxy v0.4.2 is the first version with the
  renamed declaration; it aliases v0.5.0 commit
  `aee7470331bc2a027cb2711f759f6f527a557ed8`, is no longer a current tag, and
  exact get rejects its path mismatch. The distinct HashiCorp path is a
  migration, not an in-place update.
- Highest exact-path candidate v0.4.1 is a stable GitHub Release at commit
  `b6d5c860c07ef6eeec89f4a662c7b452dd4d0c93`, tree
  `d582c4e222a01bfe89e45a22b64273970b2008a8`, dated
  2022-09-08T12:00:52Z. Selected v0.4.0 is commit
  `129ee86de65934631a7fdbeb8c5aa0ec08bfdb6c`, tree
  `dc8cd4f53669534104e66bf06f5ebaaf26d6feb4`, dated
  2022-05-25T15:01:32Z. Both lightweight tags have no tag-object signature;
  both merge commits cryptographically verify with expired GitHub web-flow
  fingerprint `5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`. The sole retraction is
  v0.3.11 for an undocumented metrics-sink breaking change.
- Versions v0.4.0 and v0.4.1 both declare Go 1.12 and the same 11 requirements. The
  candidate's complete standalone closure has 52 modules, four packages, and
  a maximum explicit Go declaration of 1.12, preserving Go 1.18. All 31 proxy
  files match the tag source. Candidate checksum pair is
  `h1:hR91U9KYmb6bLBYLQjyM+3j+rcd/UhE+G78SFnF8gJA=` /
  `h1:E6amYzXo6aW1tqzoZGT755KkbgrJsSdpwZ+3JqfkOG4=` and sumdb agrees.
- Candidate qualification fails. Under exact Go 1.26.7, proxy and exact-tag
  forms do not all pass even count-1 or race, both fail count-10 in flaky UDP
  Statsite/DogStatsD tests, and both fail vet with 46 diagnostics: 29
  `testing.T.Fatalf` calls from non-test goroutines, 12 unkeyed labels, one
  copied lock, and four empty appends. Exact Go 1.18.10 builds and passes one
  race run but reproduces count-10 failure and 44 vet diagnostics. Historical
  selected Serf and Memberlist consumers compile and vet against v0.4.1 under
  Go 1.18; focused coordinate and awareness tests pass repeated and race runs.
  Neither consumer nor Go Metrics is loaded by Ply.
- Exact candidate get changes only v0.4.0 -> v0.4.1, but also materializes
  three indirect requirements, 14 graph edges, and four full checksum lines.
  It projects 234 -> 234 modules, 3,581 -> 3,595 edges, 429 -> 429 packages,
  41 -> 41 loaded modules, zero -> zero loaded Go Metrics packages, 1,045 ->
  1,049 sum lines, and 361 -> 386 tidy-diff lines. Tidy removes all projected
  metadata and restores the inherited v0.4.0 selection. Because module stop
  rules fail, no dependency implementation commit or downstream quality claim
  was made; accepted Circbuf measurements remain unchanged.
- Fresh primary vulnerability data contains 1,392 records and no old- or new-
  path Go Metrics record or trace. Original/candidate projections retain
  identical 20-ID/22-trace Darwin and Windows reachable findings and 30
  Darwin module IDs. Go Metrics evidence has 336 verified entries; manifest
  SHA-256 is
  `5afb2b5fb26824c5e4c1db7497b8a6e6dbd24bc20194df7ca7fc5f4385e0dde1`
  and decision-summary SHA-256 is
  `c5b47a2a7d6f881e7a5aad5895d197027556e238c683bca685f45fff3b3c0592`.
- Retain exact-path `github.com/armon/go-radix v1.0.0` without dependency
  metadata. The proxy lists only stable v1.0.0; exact `@latest` and `@v1`
  resolve it at 2018-08-24T02:57:28Z. Exact `@master` resolves unreleased
  pseudo-version `v1.0.1-0.20221118154546-54df44f2176c` at
  2022-11-18T15:45:46Z. Go-import metadata identifies the public, enabled,
  unarchived, non-fork exact-path repository, which has default `master`, one
  branch, one lightweight tag, and zero GitHub Releases. Stable tag v1.0.0 is
  release-qualified but is not a GitHub Release; master is neither.
- Tag v1.0.0 identifies commit
  `1a2de0c21c94309923825da3df33a4381872c795`, tree
  `8c6d01daaee6076244d5f41247608c75a8ad4224`, parent
  `7fddfc383310abc091d79a27f116d30cf0424032`. Its lightweight tag has no
  tag-object signature; the commit verifies with fingerprint
  `7A01BBD67E7E8ADD50E00714744E147AA52F5B0A`, while GitHub currently reports
  `unknown_key`. Master merge commit
  `54df44f2176c4a553657a4f0dbe6fdb108288be3`, tree
  `87f38e748e5fc5c602ba3296793b25a782fbf02d`, has a GitHub-valid signature;
  local cryptographic verification identifies expired web-flow fingerprint
  `5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`. Its ten post-release commits
  change three files with 86 insertions and 13 deletions.
- Selected/master checksum pairs are
  `h1:F4z6KzEeeQIMeLFa97iZU6vupzoecKdU5TX24SNppXI=` /
  `h1:ufUuZ+zHj4x4TnLV4JWEpy2hxWSpsRywHrMgIH9cCH8=` and
  `h1:651/eoCRnQ7YtSjAnSzRucrJz+3iGEFt+ysraELS81M=` with the same `go.mod`
  checksum; sumdb agrees. Both seven-file proxy archives match their exact Git
  sources. Neither declares a Go version or requirements, so each complete
  minimal module closure is one module/one package. Proxy/exact-Git forms pass
  count-1/count-10/race/vet under exact Go 1.26.7, and proxy forms pass under
  exact Go 1.18.10; this execution, not missing directives, proves the floor.
- Selected Mitchellh CLI v1.1.0 is the actual graph consumer and exercises
  `New`, `Insert`, `Get`, `Walk`, `WalkPrefix`, and `LongestPrefix`; focused
  consumer tests pass repeated/race/vet against stable and master. Serf only
  records an indirect edge. Radix loads in zero Ply packages and `go mod why`
  says the main module does not need it.
- Exact selected-version get changes no selection and projects only a redundant
  main requirement/edge and full checksum: 234 modules, 3,582 edges, 429
  packages, 41 loaded modules, zero loaded Radix packages, 1,046 sum lines,
  and 367 tidy-diff lines. Exact master get changes only the Radix selection
  but similarly projects 3,582 edges and its two checksums, for 1,047 sum and
  369 tidy-diff lines. Tidy removes either pin and restores inherited v1.0.0.
  Master passes execution but is not release-qualified, so no projection was
  applied and accepted measurements remain 234/3,581/429/41/1,045/361.
- Baseline/master repository verify/build/count-1/count-10/race/vet,
  Windows, pinned lint, help, API/CLI, complete corrected preflight, and
  empty-HOME checks pass. Changed-selection-only quality and acceptance gates
  are inapplicable and not claimed. Fresh primary vulnerability results have
  no Radix record/trace and preserve identical 20-ID/22-trace Darwin and
  Windows reachable findings plus 30 Darwin module IDs. Radix evidence has
  390 verified entries; manifest SHA-256 is
  `6be6b857e77c7097b402ea5f4fe0ce849e15382b7167a004448cd85d82926eaf`
  and decision-summary SHA-256 is
  `493e832f2d0e6456bb64462104e1bd5b80f4b6eb32d401b9b098acdd0dd95e6a`.
- Retain exact-path `github.com/beorn7/perks v1.0.1` without dependency
  metadata. The exact proxy lists stable v1.0.0 and v1.0.1; exact `@latest`,
  `@v1`, `@master`, the exact-path default branch, and tag v1.0.1 all resolve
  selected commit `37c8de3658fcb183f997c4e13e8337516ab753e6` at
  2019-07-31T12:00:54Z. No prerelease, retraction, deprecation, `/v2` module,
  or later exact-path master commit exists.
- Go-import identifies the enabled, unarchived exact-path fork as the
  canonical published `github.com/beorn7/perks` source. Its GitHub
  parent/source `bmizerany/perks` is repository ancestry, not module identity:
  parent go-import declares distinct path `github.com/bmizerany/perks`, and
  its untagged 2023 master is divergent and has no `go.mod`. The fork has five
  branches, two signed annotated tags, and zero GitHub Releases; the parent
  has three branches, zero tags, and zero Releases. A stable tag is not a
  GitHub Release object.
- Tag v1.0.0 targets unsigned commit
  `4b2b341e8d7715fae06375aa633dbb6e91b3fb46` and verifies locally with
  fingerprint `5C69F212D616C4340FA8DD8504ABA6153ADA0C25`; GitHub reports
  `unknown_key`. Tag v1.0.1 targets unsigned commit
  `37c8de3658fcb183f997c4e13e8337516ab753e6` and verifies locally with
  fingerprint `A100A34F34DEC17EE5EEF14C851C3DA17D748D03`; GitHub reports valid.
  V1.0.1 changes only `go.mod`, lowering its declaration from Go 1.12 to 1.11.
- Both versions have no requirements, so each complete minimal closure is one
  module/three packages and preserves Go 1.18. All 15 proxy files match the
  corresponding exact tag. Final proxy/tag matrices pass verify/list,
  count-1, repeated count-10, race, and vet under Go 1.26.7 and Go 1.18.10;
  v1.0.1 also passes topk count-100. One earlier v1.0.0 count-10 run exposed a
  non-reproduced TopK equal-count/map-order failure; it remains recorded.
- Selected Prometheus client_golang v1.4.0 is the real consumer and uses
  `quantile.Stream`, `NewTargeted`, `Insert`, `Count`, `Query`, and `Reset`.
  Summary tests pass count-1/count-10/race under Go 1.18.10. Six historical
  consumer test-vet int-to-string diagnostics are isolated and recorded;
  mandatory Perks and project vet pass.
- Exact selected-version get changes no selection and projects only a
  redundant indirect requirement, main edge, and full checksum: accepted
  234/3,581/429/41/1,045/361 becomes
  234/3,582/429/41/1,046/370. Tidy removes those two metadata lines and
  restores inherited selection. Perks loads in zero Ply packages and
  `go mod why` says the main module does not need it. No projection or
  dependency implementation commit was made.
- Baseline/projected repository verification, build, complete count-1/
  count-10/race tests, vet, Windows, pinned lint, public help, API/CLI,
  corrected full preflight, and empty-HOME count-2 pass. Fresh primary data
  has no exact or parent Perks record and preserves 20-ID/22-trace Darwin and
  Windows reachable findings plus 30 Darwin module IDs. Perks evidence has
  477 verified entries; manifest SHA-256 is
  `aa77d2ab9cf676ecfb7ef544c1c5db76ecda86f580ff8ff8f03b8c78038677f0`
  and decision-summary SHA-256 is
  `6db4ca7260bde6b4affa48c62adb381bda20687f39f013c4b1cdc25397418754`.
- Advance exact-path `github.com/bgentry/speakeasy` from inherited v0.1.0 to
  canonical stable latest/default-branch head v0.2.0. Exact `@latest`, `@v0`,
  `@master`, tag, and public enabled/unarchived/non-fork repository master all
  resolve commit `760eaf8b681647364e7a400b856e0921248728a5` at
  2022-09-10T01:20:23Z. The lightweight tag has no tag-object signature; the
  commit's GitHub web-flow signature verifies cryptographically. The sole
  non-draft, non-prerelease GitHub Release was published much later, at
  2024-06-27T20:45:36Z, and is explicitly distinct from tag/commit time.
- Selected v0.1.0 is unsigned commit
  `4aabc24848ce5fd31929f7d1e4ea74d3709c14cd` at
  2017-04-17T20:07:03Z. Its signed annotated tag was created
  2017-06-15T22:05:56Z and verifies with fingerprint
  `757FD463E177A2F1CD1C89038B6EDBF713E83E69`; the key is now expired. Both
  eight-file proxy archives match exact tag source and both sumdb pairs
  verify. No prerelease, retraction, deprecation, `/v2` module, fork redirect,
  or later default-branch commit exists.
- Neither root module declares Go or requirements. The complete minimal closure
  is nevertheless proved as one module/two packages using only the standard
  library. Proxy and exact-Git forms of both versions pass verify/list,
  count-1, two count-10 runs, race, and vet under exact Go 1.26.7 and Go
  1.18.10. This execution, not missing directives, proves the retained floor.
- Mitchellh CLI v1.1.0 is the real graph consumer and calls only
  `speakeasy.Ask("")` through `BasicUi.AskSecret`. Its focused tests pass
  repeated/race/vet under Go 1.18. A real scratch PTY fixture waits for ECHO
  suppression before writing, proves the secret is absent from terminal
  output, and verifies ECHO restoration for old and candidate versions.
- Exact candidate get changes only Speakeasy v0.1.0 -> v0.2.0. Modules remain
  234, complete packages 429, loaded modules 41, and loaded packages 197.
  Graph edges grow 3,581 -> 3,582 only for the main candidate edge; `go.sum`
  grows 1,045 -> 1,047 only for its pair; tidy projection grows 361 -> 371
  lines and would remove the pin/pair and restore inherited v0.1.0. Speakeasy
  loads in zero Ply packages and `go mod why` says the main module does not
  need it.
- Exact Go 1.26.7 get produced dependency-only commit
  `41f9561f6ea2f5b6395c5c4d9bcc56a54533133a`, parent
  `1f55aaa31280a4610ed66c1c37ddda953bfa6a8e`, tree
  `a50d4f256fe742e472f1f7e9cf0589a596e971b0`, changing only `go.mod` and
  `go.sum` with three insertions. Repository verify/build/count-1/count-10/
  race/vet, Windows, pinned lint, byte-identical help/API/CLI, complete
  preflight, host/fresh snapshot/fresh Docker acceptance, audit meta, focused
  evidence, and empty-HOME count-2 all pass.
- Exact `make quality` passes 21/21 stages, 27/27 Q0-Q2 rows at L2, 80/80
  mutations, and zero held/regressed/not-comparable/dirty counts; scorecard
  SHA-256 is
  `48decac359a9ebab23e59c29680e63682ab6d1a65141d13911644143a8db2a01`.
  Full audit exits expected 1 only for Q3.1/Q3.3/Q3.4/Q3.7. Fresh primary
  vulnerability data has no Speakeasy record or trace and preserves exact
  20-ID/22-trace Darwin and Windows reachable populations plus 30 Darwin
  module IDs.
- Speakeasy evidence has 459 verified entries; manifest SHA-256 is
  `17bb756dfd0e81c39da3616f6f81edbcc59e299b295d422c1a701be603c02cc4`,
  and decision-summary SHA-256 is
  `d3d5678a31f494e194321951086ccdb7579c70f39a0b7cc311bc4f9c925666f0`.
- Retain exact-path `github.com/bketelsen/crypt` at selected prerelease-form
  pseudo-version `v0.0.3-0.20200106085610-5cbc8cc4026c` without dependency
  metadata. The proxy lists stable v0.0.1 through v0.0.5; exact `@latest`,
  `@v0`, `@master`, repository master, and tag v0.0.5 resolve commit
  `60c5f2086f0eae50f5275599096bcc6090d12cf8` at
  2021-10-08T10:39:19Z. The exact-path repository is an enabled, unarchived
  fork of the distinct xordataexchange lineage. Master has no later commit;
  its 13 other feature/Dependabot branches remain unreleased.
- V0.0.3 through v0.0.5 have non-draft, non-prerelease GitHub Releases, but
  all five tags are lightweight and have no tag signatures. The selected
  pseudo-version plus v0.0.3/v0.0.4/v0.0.5 merge commits independently verify
  GitHub web-flow signatures with expired fingerprint
  `5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`; v0.0.1/v0.0.2 commits are
  unsigned. Every proxy archive matches its exact Git commit and all six
  checksum pairs independently agree with sumdb.
- V0.0.5 declares Go 1.12. Its complete closure selects 151 modules across
  2,384 Go-1.26.7 graph edges, eight Crypt packages, and no declaration above
  Go 1.17, so it preserves the retained Go 1.18 floor. Proxy and exact-Git
  forms pass verify/list, count-1, two count-10 runs, and race tests under
  exact Go 1.26.7; a fully contained Go-1.18.10 replay also passes. Mandatory
  vet fails in both SDKs on the same seven unkeyed `backend.Response` literals
  across mock, etcd, Consul, and Firestore backends. Stable v0.0.3 and v0.0.4
  reproduce those seven diagnostics, so no stable version higher than the
  selected pseudo-version qualifies.
- Historical Viper v1.7.1's remote consumer compiles and vets against both
  selected and v0.0.5. Its exact config-manager constructors plus
  `ConfigManager`, `Response`, `Get`, and `Watch` are exercised by the
  consumer compile and Crypt's mock-backed Set/Get/List/Watch tests. Crypt is
  selected in Ply only through mvn-pom-mutator v0.2.3, whose source imports no
  Crypt package. Ply loads zero Crypt packages and `go mod why` says the main
  module does not need it.
- Exact candidate get changes only Crypt. It projects 234/234 modules,
  3,582/3,700 graph edges, 429/429 complete packages, 41/41 loaded modules,
  197/197 loaded packages, 1,047/1,065 sum lines, and 371/461 tidy-diff lines.
  The 118 new edges are candidate/transitive declarations; the 18 checksum
  additions are its pair plus 16 transitive module hashes. Tidy removes the
  pin and every candidate checksum, restoring the inherited pseudo-version.
  Exact selected-version get adds only a redundant main edge, requirement, and
  full checksum, all tidy-removable; neither projection was applied.
- The candidate Ply projection passes verify/build/count-1/count-10/race/vet,
  Windows, pinned lint, byte-identical public help, identical API/CLI reports,
  compatibility, and empty-HOME count-2. Changed-selection quality and
  acceptance remain inapplicable after the module stop rule. Fresh primary
  vulnerability evidence has 1,392 records, no Crypt record/trace, and exact
  old/candidate 20-ID/22-trace Darwin and Windows reachable populations plus
  30 Darwin module IDs. Crypt evidence has 10,694 verified entries; manifest
  SHA-256 is
  `e0320cf4ef063c83cee9a5131fdad9ceae8b90730502617ab48d2f16c759c93b`
  and decision-summary SHA-256 is
  `6356ad761738c8d9e664693a0c550bbe4305db500bf9198e567ef7a3e1d33a75`.
- Retain exact-path `github.com/census-instrumentation/opencensus-proto` at
  selected v0.3.0 without dependency metadata. Exact `@latest`/`@v0` are
  stable v0.4.1 at 2022-09-23T17:40:20Z; archived repository `master` instead
  resolves later unreleased `v0.2.2-0.20230502190750-1664cc961550` from the
  divergent v0.2.1 line. Primary tag evidence corrects selected v0.3.0 to
  commit `4aa53e15cbf1a47bc9087e6cfdca214c1eea4e89`, tree
  `ea0ad81b63231a53d01a5c0c90afec09692d87fb`, at
  2020-07-21T05:46:08Z. V0.4.1 is commit
  `e53624a87b9b9b919147a9b4626c669a869ebb34` on the separate release line;
  stable precedence and later default-branch commit time are not
  interchangeable.
- V0.4.1's complete standalone closure has 28 modules/64 edges and selects
  genproto `9e6da59bd2fc`, whose Go 1.19 directive violates the retained floor.
  V0.4.0's 36-module/208-edge closure preserves Go 1.18 but its committed
  generated gateway code imports grpc-gateway/v2 while `go.mod` requires v1;
  proxy and exact-Git package/test/race/vet gates fail under both Go SDKs.
  V0.4.x also makes eight incompatible exported gateway-handler signature
  changes versus v0.3.0. No higher stable release qualifies.
- Selected v0.3.0's seven zero-test packages pass verify/list/count-1/two
  count-10/race/vet in the accepted selected closure under Go 1.26.7 and
  Go 1.18.10. A separate marshal/unmarshal and bufconn TraceService fixture
  passes for v0.3.0/v0.4.1. Historical Viper v1.10.1 and Sagikazarmark Crypt
  v0.4.0 reach only `gen-go/trace/v1` through Firestore/gRPC xDS/Envoy; their
  focused tests and race pass for v0.3.0/v0.4.0/v0.4.1. Crypt's invariant
  unkeyed-literal vet finding is its own historical debt.
- Baseline remains 234 modules, 3,582 edges, 429 complete packages, 41 loaded
  modules, 197 loaded packages, 1,047 checksum lines, and a 371-line tidy
  projection. Exact selected get changes no selection and adds only a
  redundant requirement/main edge/full checksum. V0.4.0 projects 234 modules,
  3,591 edges, 1,049 checksum lines, and 377 tidy lines; only OpenCensus Proto
  changes. V0.4.1 projects 235 modules because grpc-gateway/v2 enters MVS, with
  the same edge/checksum/tidy totals. Neither projection was applied.
- OpenCensus Proto is unloaded by Ply. The v0.4.0 project projection still
  passes mod verify/build/count-1/count-10/race/vet, Windows, pinned lint,
  byte-identical help/API/CLI, compatibility, CLI surface, and empty-HOME
  count-2. Primary vulnerability results are invariant at 20-ID/22-trace
  Darwin/Windows symbol, 22-ID Darwin package, and 30-ID Darwin module
  populations, with no target record or trace. An initial prohibited
  scratch-only `go mod download all` cache warm touched no measured worktree
  and was superseded by a fresh exact-get/graph/package-list cache plus fully
  offline compatibility replay. Evidence has 459,623 verified
  entries; manifest SHA-256 is
  `55c7090ab788c633f20666ba9d70e8f4d5f4bcb446b9e5ef8f530cd78e3b074a`
  and decision-summary SHA-256 is
  `f7f4ed0dd015d7d7da58141a48f63537b8552c34761c05bc90de5b3efa22564a`.
- Accept exact-path `github.com/cespare/xxhash/v2 v2.3.0`, canonical latest,
  in dependency-only commit
  `e5d6252825d7a1822c01819b9144050f345a6ad4`, tree
  `5213ba55981d77d7c8061312915e237e80d29af8`. Exact Go 1.26.7 `go get`
  changes only XXHash v2.1.2 -> v2.3.0 and adds one `go.mod` requirement plus
  the candidate checksum pair. The proxy lists six stable v2 releases, no
  retractions or deprecation, and selected/latest declare Go 1.11. Every proxy
  ZIP matches exact Git source. The public repository is enabled, unarchived,
  non-fork, defaults to `main`, and has annotated unsigned tags but no GitHub
  Release objects. Latest tag commit is
  `998dce232f17418a7a5721ecf87ca714025a3243`; the only later main commit is a
  CI-only update and was not selected.
- Selected/candidate three-package closures pass verify/list/count-1/two
  count-10/vet/pure-Go and supported root race tests under exact Go 1.26.7 and
  Go 1.18.10. The dynamic non-race-plugin/race-host incompatibility is
  invariant while normal repeated plugin tests pass. Known-vector, streaming,
  marshal, and candidate seeded fixtures pass. Pinned apidiff finds only two
  compatible seeded API additions. Darwin-arm64, Windows/Linux amd64 assembly,
  pure-Go/appengine/non-gc fallbacks, cross-builds, and static Linux execution
  all qualify. Historical grpc xDS `Sum64String` consumer tests pass under both
  versions and SDKs; Ply itself loads no XXHash package.
- The accepted project is 234 modules, 3,583 edges, 429 complete-test packages,
  41 loaded modules, 197 loaded packages, 1,049 checksum lines, and 381
  unapplied tidy lines. Only the main edge and v2.3.0 checksum pair differ;
  tidy would remove the pin and restore inherited v2.1.2. Relative to accepted
  go-cmp, `go.sum` adds 33 lines and removes zero. Project verify/build/
  count-1/count-10/race/vet/Windows/pinned lint, help/API/CLI, offline replay,
  and empty-HOME count-2 pass.
- Exact `make quality` passes all 21 stages with 27/27 Q0-Q2 PASS at L2,
  80/80 mutants killed, 4/4 acceptance, six valid manual receipts, and zero
  held/regressed/not-comparable/dirty counts. Separate audit meta passes 15
  controls; full audit exits 1 only for queued L3 Q3.1/Q3.3/Q3.4/Q3.7 and
  attains L2. Fresh vulnerability results are invariant with no XXHash record
  or trace. Evidence has 608 entries; manifest SHA-256 is
  `6f57be6cbf177c6617ebf62e0eea7b38f6b3b41eb00873060e160ddf02f1cf34`
  and decision-summary SHA-256 is
  `4e2f0559dbefd98f22e8f1efe34a1e538ecb8e1dd2da0749af8478c88c0a1a20`.
- Retain exact-path `github.com/chzyer/logex v1.2.1` without a dependency edit.
  The proxy lists v1.1.1 through v1.1.10 plus v1.2.0/v1.2.1, and selected is
  already canonical stable latest at 2022-04-24T13:13:51Z. Short pre-module
  Git tags v1.0/v1.1 are not equivalent proxy releases. Selected lightweight
  tag/commit `2f95bdde8c3c97bfbf6d016fcc410669a895b9e7` is unsigned; the public,
  enabled, unarchived, non-fork repository defaults to master. Its five newer
  commits form only an unreleased pseudo-version, declare Go 1.21 at head, and
  touch tests/CI/module/build input rather than production Go.
- All twelve proxy archives byte-match exact Git tags; v1.2.1 is one pure-Go,
  standard-library-only package declaring Go 1.15 with no platform-specific
  files. Independent API, behavior, cross-build, floor, and consumer fixtures
  pass under exact Go 1.26.7 and Go 1.18.10. The selected native `TestLogex`
  has two stale caller-line expectations under both SDKs; unreleased commit
  `3e09012` repairs only those tests. This recorded release gap and the absence
  of a higher stable release preclude a changed selection, not continued MVS
  retention.
- The exact dependency-test path is Ply -> Promptui v0.9.0 -> Readline v1.5.1
  -> Readline test -> chzyer/test v1.0.0 -> Logex. Ply loads no Logex package.
  Readline's full/focused suites and targeted Logex calls pass; its two vet
  warnings and chzyer/test's nil-comparison panic are invariant historical
  consumer debt. Complete consumer closures top out at Go 1.17.
- Current measurements remain 234 modules, 3,583 edges, 429 complete-test
  packages, 41 loaded modules, 197 loaded packages, 1,049 checksum lines, and
  381 unapplied tidy lines. Exact selected get adds only a redundant indirect
  pin/main edge and no checksum, so it was not applied. Project gates and
  current/replay help/API/CLI reports pass. Fresh vulnerability populations
  remain 20-ID/22-trace Darwin/Windows symbol, 22-ID Darwin package, and
  30-ID Darwin module, with no Logex record or trace. Evidence has 636 entries;
  manifest SHA-256 is
  `aa33eeb8ea36b17f6830ae7f063fbe1c39ec9ac5d3c9e2417fa2fe844ceb8516`
  and decision-summary SHA-256 is
  `6720115c3316e14b89a3bedcb40a59faba4a4f898727051177a682bb1435b49d`.
- Retain exact-path `github.com/chzyer/readline v1.5.1` without a dependency
  edit. The exact proxy lists v1.5.0/v1.5.1 and selected is canonical stable
  latest at 2022-07-15T12:48:48Z. Historical v1.0-v1.4 GitHub Releases/tags
  lack patch components and `go.mod`; they are proxy-absent and resolve by
  commit only as pseudo-versions. Distinct gopkg.in vanity paths and
  unreleased main/dev_v2 pseudo-versions do not change the exact selection.
- Selected is unsigned annotated tag object
  `704f339125f222987e1fde71641f3185f6eda206` targeting unsigned commit
  `7f93d88cd5ffa0e805d58d2f9fc3191be15ec668`. Main is three commits newer
  but not released. Its exact parent is
  `fcb4d7d9a9f653462a7adf557fb1f931f00391f2`, correcting the incoming
  near-match. Proxy/Git source is byte-identical; the selected four-module
  closure tops out at Go 1.17.
- Readline's pure-Go Windows kernel32 and Unix/Linux/BSD/AIX/Solaris terminal
  splits cross-build under Go 1.18.10 and 1.26.7. Native proxy/Git tests pass
  count-1, two count-10 runs, and race under both. Vet's invariant nonstandard
  `WriteTo` and `ReadRune` signatures are historical selected APIs, not test
  failures or a candidate regression. V1.5.0 -> v1.5.1 has three compatible
  API additions.
- Ply's loaded production path is `plybuild/cmd -> Promptui v0.9.0 ->
  Readline`; Ply uses Prompt twice, while Promptui also offers Select. Fresh
  historical-consumer projections pass tests/race/vet/Windows under both SDKs.
  Readline's chzyer/test -> Logex path is dependency-test-only. Exact selected
  get is wholly inert, so project measurements remain 234 modules, 3,583
  edges, 429 complete-test packages, 41 loaded modules, 197 loaded packages,
  1,049 checksum lines, and 381 unapplied tidy lines. Project gates and the
  established vulnerability populations pass unchanged with no Readline
  record or trace. Evidence has 745 entries; manifest SHA-256 is
  `8ab036d5f96a8d92ebf682ed0fef1d5f116621dff18d670bb200172581158a00`
  and decision-summary SHA-256 is
  `435eafc93ae6df466970eb1af57a8127e5f8c389863c82fffadec236ca6bbf8a`.
- Retain exact-path `github.com/danwakefield/fnmatch
  v0.0.0-20160403171240-cbb64ac3d964` without a dependency edit. Its exact
  proxy list is empty and `@latest` is the selected pseudo-version; the
  synthesized module file declares the exact path but no Go directive or
  requirements. Sumdb records source/module hashes
  `h1:y5HC9v93H5EPKqaS1UYVg1uYah5Xf51mBfIoWehClUQ=` and
  `h1:Xd9hchkHSWYkEqJwUGisez3G1QY8Ryz0sdWrLPMGjLk=`. Proxy and exact Git
  source match at selected/master commit `cbb64ac3d964b81592e64f957ad53df015803288`.
- The public enabled/unarchived/non-fork repository defaults to master, has no
  tags or GitHub Releases, and contains only four master commits. Its later
  activity timestamp comes from unmerged PR refs. Fork tags, PR commits,
  inspired modules, and the original gist ancestry use alternate identities
  or are unreleased; no higher exact stable candidate exists.
- The standard-library-only, pure-Go one-module closure passes native tests,
  two repeated passes, race, independent matching fixtures, and relevant
  cross-builds under exact Go 1.26.7 and Go 1.18.10. This establishes the floor
  without inferring it from the directive-free synthesized module file. Vet's
  selected unreachable statement and the actual empty-string plus FNM_PERIOD
  panic are historical release gaps; unmerged PR #1 is not selectable.
- Only the root pin and historical Chroma v0.7.1 -> Fnmatch edge remain.
  Chroma v0.7.1 consumer fixtures pass, but selected final Chroma v0.10.0 uses
  `filepath.Match` and no Fnmatch package is loaded. Exact selected get is
  inert. Project metrics, gates, compatibility reports, and vulnerability
  populations remain unchanged, with no Fnmatch record or trace. Evidence has
  72 entries; manifest SHA-256 is
  `97130774da6c5f883f1c8fd0afa8190490be6aa3f8c998af4deaa9b7c79d44f1`
  and decision-summary SHA-256 is
  `96f4a6d499ac460e99d3c4b7b9395b2ed826374d79fca4602e89c359b676eae1`.
- Retain exact-path `github.com/disintegration/imaging v1.6.2` without a
  dependency edit. Its 15-version exact proxy list ends at selected/latest
  stable v1.6.2. The stable GitHub Release and lightweight tag resolve to
  verified commit `acabd8315e63bfcaac97d52d68a7a0b88d2eea93`, while later
  master resolves only as unreleased
  `v1.6.3-0.20201218193011-d40f48ce0f09`. No higher exact stable release,
  prerelease, redirect, fork, `/v2`, or alternate-path identity qualifies.
- The complete standalone closure is Imaging, its declared x/image
  pseudo-version, and x/text v0.3.0; proxy and exact-Git native suites,
  repeated/race/vet checks, independent image/property fixtures, and five
  cross-build targets pass under exact Go 1.26.7 and Go 1.18.10. Imaging is
  pure Go and the closure preserves the retained floor. Its own resource API
  does not cap input bytes or decoded pixels, so callers must bound input.
- Ply loads Imaging through `plybuild/cmd -> go-term-markdown ->
  pixterm/ansimage -> Imaging`. Independent consumer fixtures cover five image
  formats, malformed data, and all scale modes under both SDKs. They also
  preserve an ansimage-only two-pixel no-dither empty-output defect; this is
  not an Imaging failure. Project MVS's x/image v0.5.0 selection remains
  distinct and unchanged.
- Exact selected Imaging get is inert. Project counts, 381-line tidy
  projection, compatibility hashes, no-op quality checks, and 20-ID/22-trace
  vulnerability populations remain unchanged; the primary 1,392-record index
  has no Imaging record or trace. Evidence has 1,713 entries; manifest SHA-256
  is `046746e0c4004d62ebac4838dac739ce37a0d4576a0fae3e5d1db987e4d47308`
  and decision-summary SHA-256 is
  `507f403c1289ff6d698beffb31eea6c3a3c07835efcf5609d6bde475bb4cdc5e`.
- Retain exact nested module `github.com/eliukblau/pixterm/pkg/ansimage
  v0.0.0-20191210081756-9fb6cf8c2f75` without a dependency edit. Its exact
  proxy list is empty and exact `@latest` fails, but selected remains
  fetchable with its checksum pair. Root-project v1.3.0 points at selected;
  because the tag lacks the nested path prefix it does not create an exact
  nested stable release. Later root releases and the consolidated root module
  are different identities.
- Correct the incoming history claim: selected is the only commit that adds
  `pkg/ansimage/go.mod`, not the only commit containing it. Six later
  side-branch/merge commits retain it and resolve only as unreleased nested
  pseudo-versions. Five have selected source; one changes comment URLs only.
  Latest fetchable `v0.0.0-20191221044037-630511e42559` is byte-identical to
  selected. No higher qualified exact stable release exists.
- The five-module closure tops out at Go 1.13. Ansimage is one pure-Go file
  with no native tests. Independent constructor, format, scale, alpha,
  malformed-input, render, and cross-platform fixtures pass exact Go 1.26.7
  and Go 1.18.10, including the Markdown production path for six registered
  formats. Output is 24-bit ANSI only and terminal size is caller-supplied.
- Preserve selected render gaps: no-dither skips the first pixel pair, so a
  two-pixel-high image emits no ANSI while Markdown still reports success;
  four pixels render only the second pair. Dither omits its last aggregate
  row. Non-zero origins, unchecked `SetMaxProcs`, unbounded decode/input, and
  URL/consumer body lifetime remain historical selected behavior. They are
  invariant in every fetchable later exact pseudo-version, not Imaging drift.
- Exact selected get is inert. Project measurements remain
  234/3,583/429/41/197 with 1,049 checksum lines and 381 tidy-diff lines;
  project gates, compatibility hashes, preflight, 80/80 mutation controls,
  acceptance, and empty-HOME tests pass. Fresh vulnerability totals remain
  20-ID/22-trace symbol, 22 package, and 30 module findings. Ansimage has no
  direct record but is a call frame in 12 existing traces across 11 x/image
  IDs under project MVS; its declared 2019 closure has zero findings. Evidence
  has 73 entries; manifest SHA-256 is
  `95d6ca0ac0d767c94f78b3a0a5c32f1008cd2e7b05a3ffc5f639b9af021dd112`
  and decision-summary SHA-256 is
  `3f39b1aa5f16d3a1b398ce329e4a73621931b09d3b191dd7c0bbadc5743de0ed`.
- Upgrade exact-path `github.com/fatih/color v1.14.1` to highest qualified
  stable v1.15.0. V1.14.1-v1.18.0 complete closures preserve Go 1.18, but
  v1.16.0-v1.18.0 change green-bold and blue-background-italic reset bytes and
  fail 20 exact go-term-markdown v0.1.4 golden subtests. Latest v1.19.0 is
  independently ineligible because it declares Go 1.25.0.
- Exact `go get github.com/fatih/color@v1.15.0` produced dependency-only
  commit `6ca672ef38688b7f6f505cf0cb273d07c4c2ba9a`, changing only `go.mod`
  and `go.sum`. Project MVS retains go-colorable v0.1.15, go-isatty v0.0.20,
  and x/sys v0.30.0. Measurements are 234 modules, 3,583 edges, 429 complete
  packages, 41 loaded modules, 197 loaded packages, 1,051 checksum lines, and
  a 383-line unapplied tidy projection.
- All dependency, consumer, repository, compatibility, lint, acceptance,
  snapshot/Docker, audit-meta, vulnerability, and empty-HOME gates pass.
  Exact `make quality` reports all 27 Q0-Q2 rows PASS at L2 with zero held,
  regressed, not-comparable, or dirty counts; scorecard SHA-256 is
  `dae9e51e26353f72d9026e1d6eecbef697bcfa3b06cdc1905164c6d20057dafb`.
  Old/candidate vulnerability populations remain 20 IDs/22 traces, 22 package
  findings, and 30 module findings, with no Fatih Color record or trace.
- Retain exact-path `github.com/fsnotify/fsnotify v1.6.0` without a dependency
  edit. The proxy exposes 40 stable semantic versions; retracted v1.5.0 and
  v1.5.3 leave the incoming 38 eligible versions. Proxy, sumdb, go-import,
  exact Git source, stable release identities, signatures, and default-main
  ancestry agree through v1.10.1.
- V1.6.0-v1.9.0 complete minimal closures preserve Go 1.18. V1.7.0 fails its
  contained-Go count-10 suite after raising descriptor limits and reaching
  `EMFILE`; v1.8.0 fails a deterministic 2024-only test in 2026 and observes
  a separate multiple-write miss; v1.9.0 fails two independent complete Go
  1.18 repeats with an intermittent Darwin kqueue `bad file descriptor`.
  Later passing focused/full repeats characterize intermittence and do not
  erase it. Mitigation commit `0023e08` is already in v1.9.0 and explicitly
  says it does not completely fix the longstanding problem. V1.10 has further
  kqueue descriptor changes, but v1.10.0-v1.10.1 are floor-ineligible.
- Viper full/focused/race/vet consumers and independent watcher lifecycle,
  event, path, non-recursion, close, error, and concurrency fixtures pass old
  and candidate under both SDKs. Ply does not invoke Viper WatchConfig.
  Linux/inotify, Darwin/BSD kqueue, Windows, FEN, unsupported-target, pure-Go,
  CI, API, MVS, platform-build, and descriptor/overflow differences are
  separately characterized in the answered archive.
- Retained project measurements remain 234 modules, 3,583 graph edges, 429
  complete-test entries, 41 loaded modules, 197 loaded packages, 1,051 sum
  lines, and 383 unapplied tidy-diff lines. V1.9.0 would add only
  `fsnotify/internal`, two checksum lines, and the exact main edge while
  retaining x/sys v0.30.0; v1.10.1 also moves main Go 1.18 to 1.23. No
  projection was applied.
- Project verify/build/count-1/count-10/race/vet/Windows/lint, compatibility,
  byte-identical help, and empty-HOME checks pass. Changed-selection-only
  quality work is inapplicable; the accepted 27/27 Q0-Q2 L2 scorecard remains
  `dae9e51e26353f72d9026e1d6eecbef697bcfa3b06cdc1905164c6d20057dafb`.
  Fresh vulnerability populations remain 20 IDs/22 traces, 22 package
  findings, and 30 module findings with no Fsnotify record, finding, or trace.
  The 899-entry evidence manifest SHA-256 is
  `e03b23b7a574cb7d4ec4213fc7a28d785fc5a70517a864a338fa64f227db11db`;
  decision-summary SHA-256 is
  `6ce37518efc54801eea666954aa6f850142909d99c347cad21d75aa408133a34`.
- Retain exact-path `github.com/ghodss/yaml v1.0.0` without a dependency edit.
  Proxy and stable latest expose only v1.0.0. Its lightweight unsigned tag and
  commit are `0ca9ea5df5451ffdf184b4428c902747c2c11cd7`, tree
  `252e285a136d503d2913ad39cb3b1669b9a76999`, with checksum pair
  `h1:wQHKEahhL6wmXdzwWG11gIVCkOv05bNOh+Rxn0yngAk=` /
  `h1:4dBDuWmgqj2HViK6kFavaiC9ZROes6MMH2rRYeMEF04=`. Canonical Git source and
  the seven-file proxy ZIP are byte-identical. There is no prerelease,
  retraction, deprecation, redirect, or exact-path alternate stable release.
- Current master `d8423dcdf3440d0a5baffc6f90a11e4128545620` is 16 commits after v1.0.0
  but resolves only as unreleased
  `v1.0.1-0.20220118164431-d8423dcdf344`; maintained `sigs.k8s.io/yaml` is a
  distinct fork/module. Neither was substituted. Master is not a quality
  upgrade under Go 1.26 because its native tests fail via a malformed example
  vet error while the release's unreachable-code vet finding remains.
- The otherwise empty release module file resolves a complete source/test
  closure through `gopkg.in/yaml.v2 v2.4.0` and historical
  `gopkg.in/check.v1`; it preserves Go 1.18 and its native tests pass exact Go
  1.26.7 and contained Go 1.18.10. Release count-1, two count-10 repeats,
  race, independent behavior/concurrency fixtures, and production cross-builds
  pass both SDKs. Historical debt is explicit: vet reports one unreachable
  return and the release test file does not compile on Linux/386 due to an
  overflowing untyped MaxInt64, while the production package does compile.
- V1.0.0 has no strict API. Duplicate YAML keys are last-wins, unknown struct
  fields and later documents are ignored, integer/bool keys stringify,
  composite/null keys fail, and JSON/YAML types, tags, aliases, and invalid
  binary content lose information. Distinct keys `1` and `"1"` collide after
  JSON stringification and the survivor is map-iteration-dependent. Master
  strict APIs reject duplicates; unknown fields require a separate option.
- Ghodss YAML exists only through grpc-gateway v1.16.0's declared edge;
  `go mod why -m` reports it unneeded and no package is loaded. Retention keeps
  project measurements at 234/3,583/429/41/197, 1,051 sum lines, and 383
  tidy-diff lines. Exact selected `go get` would manufacture a redundant main
  edge and source sum without loading a package, so it was not applied.
- Project dependency, build/test/race/vet/lint, compatibility, CLI,
  empty-HOME, complete preflight, 80/80 mutation, and 62-control launcher
  gates pass. Fresh vulnerability populations remain 20 IDs/22 traces, 22
  package findings, and 30 module findings with no Ghodss YAML record,
  finding, or trace. The accepted 27/27 Q0-Q2 L2 scorecard remains
  `dae9e51e26353f72d9026e1d6eecbef697bcfa3b06cdc1905164c6d20057dafb`.
  The 92-entry evidence manifest SHA-256 is
  `cfb5bfcc497db42efebda5323b9e187aa63520192a7354631fd93391970bd771`;
  decision-summary SHA-256 is
  `724f7f2ac4394296cd36be360df90539deed1ecaf0b92386c6027b0f2c45f7b9`.
- Retain exact historical root module `github.com/go-gl/glfw
  v0.0.0-20190409004039-e6da0acd62b1` without a dependency edit. The root
  proxy list and root release population are empty; both repository tags are
  nested-v3.4 prereleases. Exact latest declares Go 1.19, while the last
  Go-1.18-compatible default-branch candidate still fails default
  Darwin/arm64 v3.1 compilation. The separately selected nested v3.3 module
  remains untouched.
- Selected has three Cgo-only packages and no non-standard-library module
  closure or real upstream tests. With scratch-provisioned GLFW 3.0.4 for
  v3.0, all packages pass Darwin/amd64 count-1, two count-10 repeats, and race
  under both SDKs. V3.0/v3.2 pass arm64, but v3.1 selects no default client
  library there; selected also retains two `reflect.SliceHeader` vet findings.
  Linux/Windows/FreeBSD need target C toolchains and native libraries; no
  pure-Go or headless fallback exists.
- Retention preserves project measurements at 234 modules, 3,583 edges, 429
  complete-test entries, 41 loaded modules, 197 loaded module-backed packages,
  1,051 sum lines, and 383 tidy-diff lines. The root exists only through
  x/exp's declaration; `go mod why` says it is unneeded and no GLFW package is
  loaded. Exact selected `go get` would only add a redundant explicit edge and
  source checksum, so it was not applied.
- Project build/test/race/vet/lint, Windows build, compatibility, core help,
  empty-HOME, 62-control launcher, 80/80 mutation, and 15-control audit-meta
  gates pass. Fresh vulnerability populations remain 20 IDs/22 traces, 22
  package findings, and 30 module findings with no GLFW record or trace. The
  accepted 27/27 Q0-Q2 L2 scorecard remains
  `dae9e51e26353f72d9026e1d6eecbef697bcfa3b06cdc1905164c6d20057dafb`.
  The 220-entry evidence manifest SHA-256 is
  `73a24b7220ca0c8df778188840fdbb5ed31195b8509f4aa61916ca3931f8f696`;
  decision-summary SHA-256 is
  `3e4009ecf940c8627c5fdbea7daa6ec94d100e70529a9bb3dbddb265a5143e4e`.
- Upgrade exact-path `github.com/go-logfmt/logfmt v0.4.0` to highest qualified
  stable v0.6.0. The exact proxy exposes eight stable releases and no
  prereleases; canonical go-import, proxy, sumdb, Git tag/release, tree, and
  default-main ancestry identities agree. Latest v0.6.1 declares Go 1.21 and
  is ineligible. V0.6.0 declares Go 1.17, has no module requirements, and
  preserves the retained Go 1.18 floor through its complete minimal closure.
- Selected v0.4.0 alone requires historical `github.com/kr/logfmt`; its
  closure and every serious go-logfmt candidate pass native tests, two
  count-10 repeats, race, vet, source verification, independent behavior
  fixtures, and applicable cross-builds under exact Go 1.26.7 and contained
  Go 1.18.10. V0.6.0 adds only compatible `NewDecoderSize`. V0.6.1 fixes DEL
  quoting/key validation, but v0.6.0 preserves selected behavior and the
  project loads no go-logfmt package; the fix does not override the Go-floor
  rejection.
- Exact `go get github.com/go-logfmt/logfmt@v0.6.0` produced dependency-only
  commit `3d4cfdbae0a67e757d37022be7eeedaf32c72772`, changing only `go.mod`
  and `go.sum` with three insertions. Project measurements are 234 modules,
  3,584 graph edges, 429 complete-test entries, 41 loaded modules, 197 loaded
  module-backed packages, 1,053 checksum lines, and a 392-line unapplied tidy
  projection. Relative to go-cmp commit `c314bcb`, sums are +37/-0. No
  unrelated selection moved; historical kr/logfmt remains through TSDB.
- Project dependency, build/test/repeat/race/vet/lint, Windows, compatibility,
  launcher, Make, empty-HOME, snapshot, Docker, and vulnerability gates pass.
  Exact `make quality` reports 27/27 Q0-Q2 PASS at L2, 80/80 killed mutations,
  seven ratchet improvements, and zero held/regressed/not-comparable rows;
  scorecard SHA-256 is
  `0cff5be4fb1609d296664f13a2de5567c7bdb8f574e61b6b05c0d92f44c4eb8f`.
  Fresh vulnerability populations remain 20 IDs/22 traces, 22 package
  findings, and 30 module findings with no go-logfmt or kr/logfmt record,
  finding, or trace. The 595-entry evidence manifest SHA-256 is
  `ac964598ef5933db2136fb7797738ad5595c81d94e092cf0a58c7e4e9266e4b9`;
  decision-summary SHA-256 is
  `d72e3498eaafb37925188496a131d744d717b2de33cd4aaffeb62e57f7f88dc2`.
- Upgrade exact-path `github.com/go-stack/stack v1.8.0` to highest stable
  v1.8.1. The exact proxy exposes ten stable releases and no prereleases;
  canonical go-import, proxy, sumdb, Git tag/release, tree, and default-master
  identities agree. V1.8.1 declares Go 1.17, has no module requirements, and
  preserves the complete Go 1.18 floor. Production, tests, examples, and
  exported API are byte-identical to selected; only CI and `go.mod` changed.
- Both releases pass source verification, native tests, two count-10 repeats,
  race, vet, independent capture/format/marshal/trim/concurrency fixtures, and
  24 cross-builds under exact Go 1.26.7 and contained Go 1.18.10. Their shared
  `-trimpath` `TrimRuntime` failure is path-sensitive historical behavior,
  not candidate regression, and no Go Stack package is loaded by Ply.
- Exact `go get github.com/go-stack/stack@v1.8.1` produced dependency-only
  commit `647d4fd71b226fbb1e916b4238b4c7ad87cd7975`, changing only `go.mod`
  and `go.sum` with three insertions. Project measurements are 234 modules,
  3,585 graph edges, 429 complete-test entries, 41 loaded modules, 197 loaded
  module-backed packages, 1,055 checksum lines, and a 396-line unapplied tidy
  projection. Relative to go-cmp commit `c314bcb`, sums are +39/-0.
- Exact 21-stage `make quality` exits zero with 27/27 Q0-Q2 PASS at L2,
  80/80 killed mutations, seven improvements, and zero held/regressed/
  not-comparable/dirty counts. Scorecard SHA-256 is
  `4e1e1b50c2d16666d18278efd0b9ded86df668503e32184c49a1024053dd227c`.
  Fresh vulnerability populations remain 20 IDs/22 traces, 22 package
  findings, and 30 module findings with no Go Stack record or trace. The
  817-entry evidence manifest SHA-256 is
  `e3edb91f49616d36c93ea08e85eeba8ce5730d3f61dea8b668b4da49405a8421`;
  decision-summary SHA-256 is
  `26fca7f6b5b683f4d16fb6b18a9aa88e412948ffde04c7d71764a1a94c077b72`.
- Upgrade exact-path `github.com/godbus/dbus/v5 v5.0.4` to highest
  qualified stable v5.1.0. The proxy lists eleven tags with no prereleases;
  v5.0.0/v5.0.1 are invalid list-only `/v5` anomalies. Canonical go-import,
  proxy, sumdb, Git, and GitHub identities agree. V5.2.x declares Go 1.20,
  requires `golang.org/x/sys v0.27.0`, and fails exact Go 1.18.10 source
  compilation; v5.1.0 declares Go 1.12, has no requirements, and preserves
  the complete Go 1.18 floor.
- V5.1.0 exposes only compatible API additions and passes private-daemon
  authentication, message/variant/signature, export, match/signal,
  cancellation/deadline, invalid-input, determinism, concurrency, repeat,
  race, and applicable cross-platform fixtures under both SDKs. It fixes the
  candidate-observed transmission of already-cancelled calls. Upstream retains
  a Darwin test build-tag omission, two test-only vet findings, cgo
  requirements on FreeBSD/DragonFly, and an immediate-close nonce-TCP stress
  limitation on Go 1.18; formal suites and race pass.
- Exact `go get github.com/godbus/dbus/v5@v5.1.0` produced dependency-only
  commit `6472dce617eb80484ed022ae8a53cc350c8be6fe`, changing only
  `go.mod` and `go.sum` with three insertions. Project measurements are
  234 modules, 3,586 graph edges, 429 complete-test entries, 41 loaded
  modules, 197 loaded module-backed packages, 1,057 checksum lines, and a
  400-line unapplied tidy projection. Relative to go-cmp commit `c314bcb`,
  sums are +41/-0. MVS retains the inherited go-systemd v5.0.4 edge, while
  the exact root edge selects v5.1.0; no Godbus package is loaded and no
  unrelated selection moved.
- Exact 21-stage `make quality` exits zero with 27/27 Q0-Q2 PASS at L2,
  80/80 killed mutations, seven improvements, and zero held/regressed/
  not-comparable/dirty counts. Scorecard SHA-256 is
  `1756c66aecfdd5c6d246c894d9a63630f39dc5dd789532c488e8a3ddbb247bcb`.
  Fresh vulnerability populations remain 20 IDs/22 traces, 22 package
  findings, and 30 module findings with no Godbus record or trace. The
  576-entry evidence manifest SHA-256 is
  `a665eda7b6763d2e4b0815fc19269ea3d1424eee65e19761ea897a13dd075ac9`;
  decision-summary SHA-256 is
  `8973eda3a30a7a1ec4311178f4ee4c701ba114d53ff7454ab8e947f8fcd754d2`.
- Retain exact-path `github.com/gogo/protobuf v1.3.2` without a dependency
  edit. It is the canonical latest of eight stable releases, declares Go 1.15,
  and its 12-module, 168-package closure preserves Go 1.18. Proxy, sumdb,
  go-import, tagged Git, and public unarchived non-fork repository identities
  agree; the 770-path proxy/Git manifests are byte-identical. Master is
  deprecated and unreleased, while alternate forks/modules are not exact-path
  candidates and have higher floors.
- V1.3.2 is the first release fixed for GO-2021-0053. Direct v1.3.1 harnesses
  produce exact module/package/symbol findings reaching the generated
  unmarshal plugin; v1.3.2 and the project produce no Gogo finding or trace.
  Public API is unchanged from v1.3.1. Independent golden-wire, round-trip,
  deterministic map, extension/oneof, JSON/text, invalid-input, nil/resource,
  buffer-reuse, and concurrency fixtures pass race count-20 under exact Go
  1.26.7 and contained Go 1.18.10.
- Recorded upstream qualifications are three omitted generator goldens plus
  modern gofmt drift, randomized nil GoString output, one test-only vet
  finding, and an inherited mixed binary/JSON marshal race on the same mutable
  generated message. Corrected applicable suites, commands, vet, purego,
  Linux, Windows, and js/wasm builds pass under both SDKs. No Gogo package is
  loaded by Ply, and relevant independent-message/project behavior passes.
- MVS retains v1.3.2 through Viper/etcd declarations even though `go mod why`
  is negative. Exact selected `go get` would add only a redundant indirect
  edge and source checksum without changing a selection, so it was not
  applied. Project measurements remain 234 modules, 3,586 graph edges, 429
  complete-test entries, 41 loaded modules, 197 loaded module-backed packages,
  1,057 sum lines, and the accepted 400-line tidy projection. Project
  build/test/repeat/race/vet/lint, compatibility, cross-build, empty-HOME,
  Go-1.18, and launcher/Make contract components pass. The accepted 27/27
  Q0-Q2 L2 scorecard remains
  `1756c66aecfdd5c6d246c894d9a63630f39dc5dd789532c488e8a3ddbb247bcb`.
  The 2,163-entry evidence manifest SHA-256 is
  `b92185ef4d63011df510369fb75a8f35a07d02f11c464efeedc6b0739d38f433`;
  decision-summary SHA-256 is
  `964fa19e8793454bc1e1d2db71da02cf3a5cc7c4b706e0555ce46dc9b6a12deb`.
- Upgrade exact-path `github.com/golang/protobuf v1.5.2` to highest
  qualified stable v1.5.3. The proxy exposes 18 stable releases and four
  v1.4.0 release candidates. Canonical proxy, sumdb, go-import, tagged Git,
  and public unarchived non-fork repository identities agree; the module is
  deprecated in favor of distinct successor `google.golang.org/protobuf`.
  V1.5.3 preserves Go 1.9 and Google Protobuf v1.26.0 and has no public API
  change from v1.5.2. Latest v1.5.4 was rejected because it removes public
  descriptor symbol `Default_FileOptions_PhpGenericServices`.
- V1.5.3's complete four-module closure, 21 packages, native repeats, race,
  production vet, commands, Linux/Windows/js-wasm builds, and independent
  wire/APIv1-APIv2/descriptor/extension/oneof/JSON/text/invalid-input/
  concurrency fixtures pass under exact Go 1.26.7 and Go 1.18.10. Full
  upstream vet reports only unkeyed composite literals in historical test
  fixtures. The inherited runtime has no general binary recursion cap, so
  callers must bound hostile deeply nested messages.
- Exact `go get github.com/golang/protobuf@v1.5.3` produced dependency-only
  commit `6870e029474b30ecd149d33b6c344c52cced384c`, changing only `go.mod`
  and `go.sum` with three insertions. It changes only that MVS selection;
  Google Protobuf remains v1.28.1 and every unrelated module is identical.
  The project now has 234 modules, 3,589 graph edges, 429 complete-test
  entries, 41 loaded modules, 197 loaded module-backed packages, zero Golang
  Protobuf packages, 1,059 sum lines, and a 404-line tidy projection.
- Fresh direct scans contain no Golang Protobuf module, package, or symbol
  finding, and project populations remain 30 Darwin module findings, 22
  Darwin package findings, and 20 IDs/22 reachable traces on Darwin and
  Windows. Exact 21-stage quality exits zero with 27/27 Q0-Q2 PASS at L2 and
  80/80 mutants killed. Scorecard SHA-256 is
  `65d833a10bb1cd6b48147d3446f8bc21ef3145cc9b9fc265d9ea2618d02a8807`;
  the 327-entry evidence manifest and decision summary hash to
  `5d116bb51bfea923778db9b2d6c03bf6c7df2078d966398c5c71a81a68b60dc7`
  and `fec25073702fd5dd3d9ce18fd4aa9bf5b3d9a642bed32bdd690f421778bb324e`.
- Upgrade exact-path `github.com/golang/snappy v0.0.3` to highest qualified
  stable v1.0.0. Proxy, sumdb, go-import, tagged Git, and the public,
  unarchived, non-fork repository agree; the five proxy versions are all
  stable. V1.0.0 is the sole GitHub Release. The serious candidates have no
  Go directive or requirements, their one-module standard-library-only
  closures pass exact Go 1.26.7 and Go 1.18.10, and apidiff reports only the
  compatible v0.0.4 `Reader.ReadByte` addition.
- Raw/framed golden, checksum, malformed/truncated/oversized, short/failing
  I/O, flush/close/reset, reuse, determinism, and concurrency fixtures pass
  every serious candidate under both SDKs, as do native repeats, race, vet,
  `noasm`/`appengine`, the command tool, and 108 cross-builds. All releases
  can allocate the declared raw decoded length before validating the complete
  block, so callers must prebound hostile raw inputs; framed decoding remains
  bounded. No Snappy vulnerability finding or trace exists.
- Exact Go 1.26.7 `go get github.com/golang/snappy@v1.0.0` produced
  dependency-only commit `372f8e988d9bb86a4f426ae6f953d2e8226c05fe`, changing only
  `go.mod` and `go.sum` with three insertions. It adds the exact root edge and
  changes no other module selection. The project now has 234 modules, 3,590
  graph edges, 429 complete-test entries, 41 loaded modules, 197 loaded
  module-backed packages, zero Snappy packages, 1,061 sum lines, and a
  407-line tidy projection.
- Exact 21-stage quality preserves 27/27 Q0-Q2 PASS at L2, kills 80/80
  mutants, and records seven improvements with zero held, regressed,
  non-comparable, or dirty counts. Scorecard SHA-256 is
  `7320f432090c578973bcfe4abc8736ef68a83f492a5bd057922ec150ce2a03c5`;
  the 677-entry evidence manifest and decision summary hash to
  `e4c81eed9a21f737f3a8427360908835b7f195f9981bd1be9a5bdbd953091449`
  and `3d606c0a6ca3f0cd523c4603de9ded354311d1ac6135c9f679ee58371d3ebc79`.
- Upgrade exact-path `github.com/google/martian/v3 v3.2.1` to highest
  API-compatible qualified stable v3.3.2. Canonical proxy, sumdb, go-import,
  tagged Git, and the public archived non-fork repository agree. The proxy
  exposes five stable v3 versions with no prerelease or retraction. Master
  equals v3.3.3, and malformed `3.2`/`v3.2` tags do not version this module.
  V3.3.3 is rejected because it removes exported `martian.Init` and the
  inherited `cmd/marbl -v` flag.
- Selected and candidate have byte-identical Go-1.11 module files and
  37-module complete closures. Their 46-package, 428-entry Go-1.26.7 and
  365-entry Go-1.18.10 test closures pass native tests, serial repeats, race,
  production vet, commands, and 1,380 cross-build contexts. Apidiff is empty.
  Independent modifier, context, HTTP proxy, MITM/CONNECT, HAR, malformed-I/O,
  lifecycle, reuse, and concurrency fixtures pass both SDKs.
- Retained qualifications include non-idempotent proxy close, unbounded MITM
  certificate/HAR/gRPC input growth, single-read HTTP/2 preface forwarding,
  permissive truncated terminal gRPC frames, malformed body-range slice panic,
  and traffic-shape map nondeterminism. These are inherited by all serious
  candidates and require caller input bounds and coordinated shutdown.
- Exact Go 1.26.7 `go get github.com/google/martian/v3@v3.3.2` produced
  dependency-only commit `4644476aef78e24d0c7e6a1140baba4c3226a1cc`, changing
  only `go.mod` and `go.sum` with three insertions. It moves no unrelated
  selection. Historical Cloud Go vertices explain selection despite negative
  `go mod why` and zero loaded Martian packages. The project now has 234
  modules, 3,597 graph edges, 429 complete-test entries, 41 loaded modules,
  197 loaded module-backed packages, 1,063 sum lines, and a 418-line tidy
  projection.
- Fresh project vulnerability populations remain 30 module findings, 22
  Darwin package findings, 23 Windows package findings, and 20 IDs/22 traces
  on both symbol platforms, with no Martian record, finding, or trace. Exact
  21-stage quality preserves 27/27 Q0-Q2 PASS at L2 and kills 80/80 mutants.
  Scorecard SHA-256 is
  `094a7668e43eb0a5b563bcd637b3b5a301e1bb5f5053903ad8934e2b2906ae7f`;
  the 509-entry evidence manifest and decision summary hash to
  `614921201697ce6f3accdf0c6e619c86ddb3502f23c3b03ee6314a2347a9a649`
  and `374c235a5ebeeb7a97254e8eab57b3dc60224fc7f103611beceafff701f5756e`.
- Retain exact-path
  `github.com/google/pprof v0.0.0-20210720184732-4bb14d4b1be1` without a
  dependency edit. Proxy, Git, and GitHub expose no tag or release. Proxy
  latest is the protected `main` head, has no release qualification, declares
  Go 1.25.0, and incompatibly changes `Profile.Aggregate`,
  `driver.ObjTool.Open`, two exported embedded-asset packages, and a public
  constant. Selected and latest proxy/Git manifests match their exact commits.
- Selected's six-module, 18-package closure declares Go 1.14 and builds under
  contained Go 1.18.10. Native tests with package vet disabled, independent
  repeats, race, focused profile round-trip/merge/malformed/concurrency
  fixtures, local/HTTP command flows, deterministic top/DOT reports, and
  Darwin/Linux/Windows/FreeBSD builds pass both SDKs. Inherited qualifications
  are whole-buffer input growth, unauthenticated timeout-free web serving,
  lost gzip-close errors, an unclosed saved profile, global non-reentrant
  driver state, one legacy vet warning, a count-2 global-flag test panic, and
  unsupported js/wasm through historical Readline.
- Historical Cloud Go v0.90.0 supplies pprof through the exact retained graph
  chain even though MVS selects Cloud Go v0.105.0. `go mod why -m` remains
  negative and no pprof package loads. The project stays at 234 modules,
  3,597 edges, 429 complete-test entries, 41 loaded modules, 197 loaded
  module-backed packages, 1,063 sum lines, and a 418-line tidy projection.
  Exact latest would force main Go 1.25.0 and move only pprof, Demangle, and
  X/Sys while leaving the loaded population unchanged.
- Fresh direct vulnerability scans are zero and have no pprof record. Project
  base/latest populations are identical at 30 module findings, 22 Darwin
  package findings, 23 Windows package findings, and 20 IDs/22 traces per
  symbol platform. Exact Go 1.26.7 project gates and complete contained
  preflight pass; the Go 1.18.10 projection retains only the two accepted
  shell error-wording differences. Changed-selection quality is inapplicable,
  so accepted 27/27 Q0-Q2 L2 scorecard SHA-256 remains
  `094a7668e43eb0a5b563bcd637b3b5a301e1bb5f5053903ad8934e2b2906ae7f`.
  The 188-entry evidence manifest and decision summary hash to
  `e6090a35f4f669bfff0eccf3b179def466dcbc18645ec61ed5d2e81f92e4f0f7`
  and `20804b312d90e5390df74dde4c1b252e68d8b1183897e7e562f2c8a133ac5d08`.
- Upgrade exact-path `github.com/google/renameio v0.1.0` to the highest
  qualified stable v1.0.1. Proxy, sumdb, go-import metadata, and tagged Git
  agree on the public active non-fork repository. Root v2 tags instead
  declare `github.com/google/renameio/v2`; they are not exact-path
  candidates. All serious root releases have standard-library-only closures
  compatible with Go 1.18, and apidiff finds only the compatible v1
  `renameio/maybe` package addition.
- Selected and v1.0.1 pass native verification, tests, repeats, race, vet, a
  309-line atomic-write fixture under both SDKs, and all 48 production
  cross-builds. V1.0.0 is rejected because `maybe` cannot compile on Windows;
  v1.0.1 supplies the repaired non-atomic Windows fallback. The atomic API
  syncs file content before close/rename but does not fsync the parent
  directory, replaces symlinks and hard-link destinations rather than
  preserving their inode identity, and is last-writer-wins for independent
  concurrent writers.
- Exact Go 1.26.7 `go get github.com/google/renameio@v1.0.1` produced
  dependency-only commit `394ec36ac5cb6712640024f0a47af28bf3631906`,
  changing only `go.mod` and `go.sum` with three insertions. Historical
  HTools vertices explain MVS selection despite negative `go mod why` and
  zero loaded Renameio packages. The project now has 234 modules, 3,598
  edges, 429 complete-test entries, 41 loaded modules, 197 loaded
  module-backed entries, 1,065 sum lines, and a 428-line tidy projection.
- Fresh direct selected/v1.0.1 Renameio vulnerability scans are zero; v1.0.0
  module/package scans are zero and its Windows symbol load hits the known
  compile defect. Project base/v1.0.1 populations are identical at 30 module
  findings, 22 Darwin and 23 Windows package findings, and 20 IDs/22 traces
  per symbol platform. Exact 21-stage quality preserves 27/27 Q0-Q2 PASS at
  L2 and kills 80/80 mutants.
  Scorecard SHA-256 is
  `6be51aec069011e99a80c2b5b5097bb79a67abe1c0078589ccd60198ab75c2ce`;
  the 7,798-entry evidence manifest and decision summary hash to
  `4d4fd72b84786270f069b18934eb831d31ecc801e1d348e6e912f756fc114f84`
  and `07fcaace013997f5d4513de49c5a59b3359b7b270c4e6634b39729c8a45a6096`.
- Upgrade exact-path `github.com/google/uuid v1.1.2` to highest qualified
  stable v1.4.0. Proxy, sumdb, go-import, lightweight Git tags, thirteen
  GitHub Releases, and the public active non-fork BSD-3-Clause repository
  agree. All thirteen stable releases have one-package standard-library-only
  closures with no Go directive or requirements and pass complete source/test
  closure gates under exact Go 1.26.7 and contained Go 1.18.10.
- Apidiff through v1.4.0 reports compatible additions only. Independent
  parse/format, version/variant, v1/v2/v4, name-hash, entropy, clock sequence,
  node identity, concurrency, pool, serialization, and SQL contracts pass
  both SDKs. The upstream same-process repeat panic through v1.4.0 is a
  classified test-order leak from a finite global random reader; independent
  fresh-process repeats pass.
- Reject v1.5.0 and v1.6.0 because their new UUIDv6 generator and `Time()`
  extractor do not round trip under either SDK. V1.5.0 also duplicates fixed
  input V7 values; v1.6.0 repairs V7 monotonicity but retains the V6 defect.
  Master `53dda83` fixes V6 but is unreleased and was not promoted.
- Exact Go 1.26.7 `go get github.com/google/uuid@v1.4.0` produced
  dependency-only commit
  `cf53bc64eeb69471d35c7536d196bf1da15f3973`, changing only `go.mod`
  and `go.sum` with three insertions. Fourteen historical gRPC vertices
  explain selection despite negative `go mod why` and zero loaded UUID
  packages. The project now has 234 modules, 3,599 edges, 429 complete-test
  entries, 41 loaded modules, 197 loaded module-backed packages, 1,067 sum
  lines, and a 432-line tidy projection.
- Fresh direct selected/candidate UUID vulnerability scans are zero and
  project populations remain 30 module findings, 22 Darwin and 23 Windows
  package findings, and 20 IDs/22 traces per symbol platform. Exact 21-stage
  quality preserves 27/27 Q0-Q2 PASS at L2 and kills 80/80 mutants.
  Scorecard SHA-256 is
  `576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`;
  the 3,206-entry evidence manifest and decision summary hash to
  `0796b997d528e8b9d25d35cbd90a04cc93bf89fbbb3cdacc27fcd46546d00c6d`
  and `522b3fc0387507a9612daab87242625d2bfde7aabd74a45ebd24b058d883c31b`.
- Retain exact-path `github.com/googleapis/gax-go/v2 v2.7.0` without a
  dependency edit under the two explicit 2026-09-12 product decisions. Fresh
  proxy, sumdb, go-import, tagged Git, and the public active non-fork
  repository agree. Of 42 v2-looking stable proxy tags, v2.0.0 does not
  contain `v2/go.mod`; the other 41 exact-path releases are valid. No
  redirect, alternate path, fork, deprecation, retraction, prerelease, or
  arbitrary branch head qualifies as a better release.
- V2.5.1 is the highest stable release declaring Go 1.18; every v2.6.0-or-
  newer release declares at least Go 1.19, v2.12.5 declares Go 1.20, and
  latest v2.24.1 declares Go 1.25. Six serious releases pass exact Go 1.26.7
  verification, complete closure tests/repeats, race, vet, cross-builds, and
  source identity checks. Deep selected/v2.5.1 gates also pass contained Go
  1.18.10. V2.12.4/v2.12.5 fail technically in their gRPC generic-atomic
  closure under Go 1.18, and latest is rejected by the parser. Technical
  success does not erase selected's three Go 1.19 closure declarations; the
  authorized inherited-floor exception remains necessary.
- Selected exposes the `gax` retry/call/stream utilities and `apierror`.
  Independent option-ordering, retry/no-retry, deadline/cancellation,
  gRPC/HTTP conversion, backoff/timer, error-detail, header/content,
  JSON-stream, nil/panic, concurrency, and platform contracts pass both SDKs.
  Retained qualifications include mutable non-concurrent Retryer/Backoff,
  package-global jitter, first-pause over-cap and negative-duration panic
  boundaries, odd-header/nil-reader panics, and non-idempotent stream close.
- `apierror.ParseError(err, false)` can construct an `APIError` whose required
  `Error()` panics. Upstream repair `22c16e7bff` first shipped in v2.12.5,
  whose Go 1.20 floor is ineligible. The authorized behavior exception is
  bounded by the freshly revalidated zero-loaded-package invariant. If GAX is
  loaded or directly imported later, the exception expires and that owning
  checkpoint must stop for a fresh dependency and product decision before
  merge.
- Viper v1.15.0 supplies shortest path main -> Viper -> GAX v2.7.0, while
  historical Cloud/Google/Viper/Crypt vertices declare v2.0.4-v2.1.1.
  `go mod why -m` is negative and exactly zero GAX packages load. Exact
  v2.5.1 would downgrade Viper and unrelated selections; later candidates
  broadly upgrade unrelated selections, and latest raises the main Go line.
  A redundant v2.7.0 root edge has no selection purpose. The unchanged project
  stays at 234 modules, 3,599 edges, 429 complete-test entries, 41 loaded
  modules, 197 module-backed packages, 1,067 sum lines, and a 432-line tidy
  projection. Raw candidate module/edge/sum counts are 234/4,088/1,158 for
  v2.5.1, 346/3,777/1,077 for explicit selected, 248/3,674/1,084 for
  v2.12.5, and 260/3,763/1,106 for latest; corresponding tidied counts are
  231/3,555/950, base-identical, 234/3,561/949, and 234/3,582/958, always
  with zero loaded GAX packages.
- Fresh primary vulnerability data has 1,398 module records and no exact GAX
  advisory. Project populations remain 30 module findings, 22 Darwin and 23
  Windows package findings, and 20 IDs/22 traces per symbol platform, with no
  GAX occurrence. Direct selected-source scans record inherited reachable
  `x/net` HTTP/2 and protobuf protojson findings; they do not enter the
  unchanged project, but remain part of the future-loading stop condition.
- Exact Go 1.26.7 project, compatibility, lint, empty-HOME, cross-build,
  preflight, host/snapshot/Docker, mutation, and audit gates pass. Applicable
  Go 1.18.10 gates pass with only the two accepted shell error-wording
  differences. The unchanged authoritative scorecard remains 27/27 Q0-Q2
  PASS at L2 with SHA-256
  `576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.
  The 863-entry selected-evidence manifest and decision summary hash to
  `8b8ea0d9b842b8b3eecc1d7c7b8f8754c24662a74dd21d9d1cef0304bbd18b5b`
  and `26a2a6a7ad4d06a2c059959d7eaf8df64ee9baeb1688c03927b08b1fd650bfe1`.
- Retain exact-path `github.com/googleapis/google-cloud-go-testing
  v0.0.0-20200911160855-bcd43fbb19e8` without a dependency edit. Fresh
  proxy, sumdb, `go-import`, Git, and GitHub evidence resolves the public
  archived non-fork repository. It has zero tags and releases, its release
  document explicitly says it has no releases, and the stable proxy list is
  empty. Selected is branch `bcb-to-fb` tip `bcd43fbb19e8`; proxy latest is
  archived `master` tip `1c9a4c676720`. Strict fsck and ancestry pass.
  Selected is two commits behind latest, and only README archive wording plus
  SECURITY.md changed; proxy/Git manifests and module sums agree.
- Selected/latest declare Go 1.11 and have identical requirements, source,
  tests, and exported API. Their isolated closure contains 36 modules and
  337/336 edges under Go 1.26.7/1.18.10. The imported production/test closure
  uses 14 external modules, peaks at Go 1.11, and passes both SDKs; the whole
  graph including tools peaks at Go 1.12. The five packages are root docs plus
  BigQuery, Datastore, Pub/Sub, and Storage interface adapters. Independent
  forwarding, context/callback, aliasing, iterator/stream, nil/panic,
  concurrency, and platform fixtures pass. Three malformed upstream example
  suffixes fail modern vet only; a scratch-only name correction makes default
  tests/vet/cross-builds pass and does not affect runtime qualification.
- Afero v1.9.4 supplies the sole current incoming edge and shortest path main
  -> Afero -> selected; historical Afero v1.8.2 declares the same version.
  `go mod why -m` is negative, and zero target or GAX packages load. Raw
  selected/latest projections change no unrelated selection or loaded
  population; both temporarily add the target and five tools requirements,
  while tidy makes all projections byte-identical and restores selected.
  Project measurements remain 234 modules, 3,599 edges, 429 complete-test
  entries, 41 loaded modules, 197 module-backed packages, 1,067 sum lines, and
  a 432-line tidy projection. Fresh primary data has no target advisory and
  preserves project vulnerability populations. Exact project gates and a
  fresh 27/27 Q0-Q2 L2 run pass; accepted scorecard remains
  `576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.
  The 788-entry evidence manifest and decision summary hash to
  `34e3574d2ced8aef523df1607300134664907cf7634cb40f89c3f615a151f106`
  and `9652c69c7971189749d2fd4590dc39b49a11f6f159cb80c2287767251312db2a`.
- Stop Enterprise Certificate Proxy for a fresh product decision without a
  dependency edit. Fresh exact-path proxy, sumdb, go-import, Git, and GitHub
  evidence resolves the public active non-fork Apache-2.0 repository and 29
  stable-semver versions. Git tag `0.3.10` lacks the required `v` and does not
  version the module. GitHub marks v0.1.0, v0.2.0, v0.2.1, and v0.3.15 as
  prereleases/public previews despite stable proxy version syntax.
- V0.2.0 is the highest declaration-eligible tag. V0.1.0/v0.2.0 declare Go
  1.18, v0.2.1-v0.3.4 declare Go 1.19, and later releases require Go 1.23 or
  newer. V0.1.0/v0.2.0/v0.2.1 have standard-library-only minimal closures:
  202/221 production/test entries under Go 1.26.7 and 137/157 under Go
  1.18.10. Selected v0.2.1 cannot be retained under the current contract
  because it has no Go-1.19 exception.
- Authoritative v0.2.0/v0.2.1 module zips lose the executable bit on the
  upstream signer shell fixture and fail raw tests under both SDKs. A scratch-
  only mode correction makes tests, repeats, race, vet, native builds, and
  broad cross-builds pass, but does not repair the release. V0.2.0 and v0.2.1
  have identical exported Go API; v0.1.0 is API-incompatible.
- Both v0.2 clients have unbounded signer startup/RPC waits, incomplete failed-
  startup cleanup, non-idempotent close, caller aliasing, typed-nil panics, and
  error-string-dependent lifecycle behavior. V0.2.1 additionally permanently
  discards the process-global logger in its default path. The C ABI trusts raw
  pointer/length pairs, permits out-of-bounds access, conflates errors as zero,
  truncates short buffers, and lets invalid inputs escape as uncaught Go
  panics. Independent valid-flow, malformed-flow, concurrency, process,
  configuration, crypto, and C-shared characterization ran under both SDKs.
- Viper v1.15.0 supplies the sole project edge to selected v0.2.1; the target
  and GAX each have zero loaded packages. An exact v0.2.0 request necessarily
  downgrades Viper to v1.14.0, GAX to v2.6.0, and changes 20 selections total.
  Its raw projection is 233 modules/3,601 edges/1,069 sums with unchanged
  429/197/41 package/module loading counts; tidy leaves 21 selection
  differences from base tidy. V0.1.0 downgrades Viper to v1.13.0 and tidy
  removes the target. These are not authorized target-only changes.
- Fresh primary vulnerability data contains no target record. Direct serious-
  candidate scans are zero; base and v0.2.0 project populations are identical
  at 30 module findings, 22 Darwin and 23 Windows package findings, and 20
  IDs/22 traces on both symbol platforms. Exact Go 1.26.7 project and preflight
  gates pass; compatible Go 1.18.10 gates retain only the two accepted shell
  wording failures. No metadata changed, so accepted quality remains 27/27
  Q0-Q2 PASS at L2 with scorecard
  `576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.
  The 562-entry evidence manifest and decision summary hash to
  `59868d9547ece628b950faf9200377108a0aab55fbbfab731d5f0b5fbcd10033`
  and `276efebfe75500789ce6501630e8ecdb0dbbf8241ae20c74f57bccfbc348180c`.
- Product direction supplied 2026-09-13: the user explicitly selected option 1
  with the recommended bounds. Retain exact selected, unloaded Enterprise
  Certificate Proxy v0.2.1 without dependency metadata changes. Preserve the
  main Go 1.18 floor and accept v0.2.1's Go 1.19 declaration only for this exact
  version and unchanged inherited Viper v1.15.0 edge. Separately accept only the
  already documented public-preview/release-packaging, configuration/file-
  lifetime, process-lifecycle, close, nil/panic, aliasing, global-logger, and C
  ABI behavior/safety findings. The invariant was revalidated from clean
  decision HEAD `f671bc0`: selected v0.2.1 has exactly one incoming edge from
  Viper v1.15.0, `go mod why -m` remains negative, and the 429-entry complete
  test load contains zero target packages and zero GAX packages. `go.mod` and
  `go.sum` remain unchanged, with SHA-256 values
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
  and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Direct import/loading, runtime reachability, a target version or incoming-
  edge change, or a new vulnerability/advisory or independent disqualifier
  expires the exceptions and requires a fresh dependency/product decision
  before merge. Do not add a root edge, patch/replace the target, change or
  independently audit Viper/GAX, raise the Go floor, manufacture a dependency
  commit, or request this same decision again while the invariants hold.
- Stop exact-path `github.com/gopherjs/gopherjs` for a fresh bounded product
  decision without a dependency edit. Fresh proxy, sumdb, go-import, strict
  Git, and GitHub evidence resolves the public active non-fork BSD-2-Clause
  repository and thirteen proxy versions. Selected is unsigned 2018 commit
  `0766667cb4d1cfb8d5fde1fe210ae41ead3cf589`, not a tagged release, and its
  proxy module file is synthetic. Stable v1.17.2 is the highest release whose
  complete imported closure preserves Go 1.18; the v1.18 releases are semver
  prereleases and stable v1.20/v1.21 releases raise the floor.
- Reject selected because it has no declared closure, misses a vendored test
  package in the module zip, resolves to post-floor tools under modern MVS,
  cannot complete its tests, and deliberately fails compilation outside Go
  1.11. V1.12.80 similarly requires Go 1.12. V1.17.2's 157-module graph peaks
  at Go 1.17 and its actual production/test closure uses 29 external packages
  across 14 modules, but complete testing under Go 1.18.10 reproducibly fails
  seven Darwin/arm64 generated programs through unimplemented
  `internal/abi.FuncPCABI0`. A trivial generated program crashes on first
  output. Go 1.26.7 adds AST and VERSION-parser failures. Its supported vet and
  race scopes pass, and Linux JS/source maps are deterministic, but those
  partial passes do not qualify the release.
- Reject v1.18.0-beta3 because it is a prerelease and its command cannot link
  under required Go 1.26.7 through `go/build.parseGoEmbed`, despite a complete
  Go 1.18.10 suite pass. Reject v1.20.2 and v1.21.0 as floor-ineligible.
  Pinned apidiff also records incompatible build/compiler API changes between
  each release line. Independent session-state fixtures prove nil panic and
  caller-options mutation boundaries under both SDKs.
- The target's only path is main -> direct mvn-pom-mutator v0.2.3 -> GoConvey
  v1.6.4 -> selected. `go mod why -m` is negative and zero target packages
  load. Disposable v1.17.2 and beta3 requests add closure modules but retain
  zero target loading and tidy byte-identically back to the base selection.
  V1.20.2/v1.21.0 alter unrelated `x/*` selections; v1.21.0 also raises the
  main Go line. The unchanged project remains 234 modules, 3,599 edges,
  429/197/41 complete-test/module-backed/loaded-module entries, 1,067 sum
  lines, and the established 432-line tidy projection. ECP and GAX separately
  remain at zero loaded packages.
- Fresh primary data has 1,398 records and no exact GopherJS advisory. Direct
  candidate package scans report inherited Logrus and `x/sys` advisories but
  zero reachable symbol or test-symbol traces. Project populations remain 30
  module findings, 22 Darwin and 23 Windows package findings, and 20 IDs/22
  traces on both symbol platforms, with no target occurrence. Exact Go 1.26.7
  project/repeat/race/vet/lint/cross-build/API/CLI/full-preflight gates pass;
  the Go 1.18.10 projection retains only the two accepted shell wording
  failures. No metadata changed, so accepted quality remains 27/27 Q0-Q2 PASS
  at L2 with scorecard
  `576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.
- Product direction recorded 2026-09-13: the user explicitly selected option 1
  with the recommended bounds. Retain exact selected, inherited, unloaded
  GopherJS `v0.0.0-20181017120253-0766667cb4d1` without dependency metadata
  changes. Accept only its already documented unqualified pseudo-release,
  unsigned/untagged identity, proxy-synthesized module file, incomplete modern
  source/test closure, missing vendored test package, post-floor resolution,
  and known compiler/command/generated-runtime incompatibilities. No new or
  independently discovered defect is accepted.
- Every guard was revalidated from clean decision HEAD `ae2a9d3` with freshly
  unpacked, hash-verified exact Go 1.26.7 first in `PATH`, `GOENV=off`,
  `GOWORK=off`, `GOTOOLCHAIN=local`, `LC_ALL=C`, `LANG=C`, and no ambient
  `GOFLAGS`. Selected GopherJS and its sole incoming GoConvey v1.6.4 edge are
  unchanged; `go mod why -m` remains negative; repository source has zero
  direct target imports; and the 429-entry complete test load has 197 module-
  backed packages across 41 modules with zero GopherJS, Enterprise Certificate
  Proxy, or GAX packages. GopherJS therefore remains runtime-unreachable in the
  current project. Project selection remains 234 modules and 3,599 edges.
  `go.mod` and `go.sum` remain unchanged at SHA-256
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
  and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`;
  `go.sum` remains 1,067 lines. Fresh primary data still has 1,398 records,
  Last-Modified 2026-09-10T16:28:28Z, and no exact GopherJS entry.
- The GopherJS exception is non-transferable and valid only while zero target
  packages load, the module remains runtime-unreachable, its exact version and
  sole GoConvey v1.6.4 edge remain unchanged, and no new advisory or independent
  disqualifier appears. Direct import/loading, runtime reachability, a version
  or incoming-edge change, or a new advisory or independent disqualifier
  expires the exception and requires a fresh dependency/product decision
  before merge. Do not add a direct edge, change GopherJS, mvn-pom-mutator, or
  GoConvey, authorize a broader parent/removal group, raise the Go floor, select
  a prerelease/fork/replacement/patch/alternate path, move unrelated modules,
  manufacture a dependency commit, or request the same decision again while
  the guards hold. Enterprise Certificate Proxy and GAX retain only their own
  separate guarded exceptions.
- Stop exact-path `github.com/gorilla/websocket` for a fresh bounded product
  decision without a dependency edit. Fresh proxy, sumdb, `go-import`, strict
  Git, and GitHub evidence resolves the canonical public active unarchived
  non-fork BSD-2-Clause repository and eleven proxy versions. Serious stable
  releases are v1.4.2, v1.5.0, v1.5.1, and v1.5.3. GitHub marks v1.5.2 a
  prerelease; the proxy preserves original commit `1bddf2e0`, while its forge
  tag was overwritten to `9ec25ca` and unreleased main later retracts it.
- V1.4.2, v1.5.0, and v1.5.3 declare Go 1.12 and have standard-library-only
  complete minimal closures under exact Go 1.26.7 and Go 1.18.10. V1.5.1 and
  v1.5.2 declare Go 1.20, add an `x/net` closure, and have broken default
  vendor mode because their proxy zips retain `vendor/modules.txt` without the
  vendored packages. V1.5.2 also imports Go-1.20-only
  `http.NewResponseController` and fails Go 1.18 compilation. Thus v1.5.3 is
  the highest stable floor-preserving release, but it is not behavior/security
  qualified.
- Fresh primary GO-2026-6278/GHSA-w67g-5rqw-f597 evidence says weak WebSocket
  client-mask PRNG affects versions before v1.5.3 and labels v1.5.3 fixed.
  Inspected release source contradicts that range: v1.5.3 explicitly reverts
  `newMaskKey` to `math/rand`, and a seed-controlled independent fixture proves
  repeatable mask keys. Security commit `d67f41855da42d7bccd9ef050c49f7e54e783b95`
  switches to `crypto/rand` only after v1.5.3 on unreleased main. Original
  proxy v1.5.2 contains that fix, but is prerelease and floor-ineligible. No
  published exact-path stable release both preserves Go 1.18 and contains the
  fix.
- The reviewed govulncheck v1.8.0 snapshot updated
  2026-09-10T14:48:42Z omits the unreviewed GO record, so it reports no target
  occurrence. The primary index has 1,398 module records, Last-Modified
  2026-09-10T16:28:28Z, and one affected selected module. Project
  package/symbol/test-symbol/reachable target results are zero only because no
  target package loads. Source evidence and the fresh primary record take
  precedence over the reviewed-snapshot omission.
- Every package, example, test, benchmark, generated/build-tag file, exported
  API, and relevant connection boundary was inspected. The independent fixture
  covers real handshake success/failure, headers/auth/cookies, origin and
  subprotocol policy, proxy/TLS-visible paths, compression, fragmentation and
  control-frame interleaving, malformed/oversized frames, close/deadline/error
  behavior, prepared-message copying/concurrent reuse, pools, deterministic
  masking, and cleanup under both SDKs. A repeated whole-source Go 1.26 race
  selection finds an upstream cross-test lifecycle defect (`cstHandler` logs
  after its test returns) in both v1.4.2 and v1.5.3; isolated supported-use
  fixture race passes.
- The sole incoming selected edge is direct mvn-pom-mutator v0.2.3;
  `go mod why -m` is negative, repository source has no target import, and zero
  target packages load. Exact v1.4.2/v1.5.0/v1.5.3 projections add a redundant
  direct requirement but no unrelated selection and tidy byte-identically to
  base. V1.5.1/v1.5.2 request unrelated `x/*` upgrades that remain after tidy.
  The unchanged project stays at 234 modules, 3,599 edges, 429/197/41 load
  entries/module-backed packages/loaded modules, 1,067 sum lines, and a
  432-line tidy projection. `go.mod`/`go.sum` SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
  and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- Exact Go 1.26.7 project verify/load/build/repeats/race/vet/lint/hermetic/
  cross-build/API/CLI/full-preflight gates pass. Applicable Go 1.18.10 gates
  retain only the two accepted `pkg/shell` wording failures. No metadata
  changed; accepted quality remains 27/27 Q0-Q2 PASS at L2 with scorecard
  `576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.
  Selected evidence and decision-summary hashes are
  `cc15144b88edbf330e1300dda430eead2419087883d97788edb7d781bb48b912`
  and `4f36e9a15e7e17cf9a09b79ef4f875e2dbecbb7f171682d803aca9514e8d7f91`.
- Product direction recorded 2026-09-14: the user explicitly selected option 1
  with the recommended bounds. Retain exact selected, inherited, unloaded
  Gorilla WebSocket v1.4.2 without dependency metadata changes. Accept only
  GO-2026-6278's documented weak `math/rand` client-mask behavior and the
  documented upstream full-source Go 1.26 cross-test lifecycle race. No new or
  independently discovered defect is accepted.
- Every guard was revalidated from clean decision HEAD `778d578` with freshly
  unpacked, hash-verified exact Go 1.26.7 first in `PATH`, `GOENV=off`,
  `GOWORK=off`, `GOTOOLCHAIN=local`, `LC_ALL=C`, `LANG=C`, and no
  ambient `GOFLAGS`. Selected v1.4.2 and its sole incoming mvn-pom-mutator
  v0.2.3 edge are unchanged; `go mod why -m` remains negative; repository
  source has zero direct target imports; and the 429-entry complete test load
  has 197 module-backed packages across 41 modules with zero Gorilla WebSocket,
  GopherJS, Enterprise Certificate Proxy, or GAX packages. Gorilla WebSocket
  therefore remains runtime-unreachable. Project selection remains 234 modules
  and 3,599 edges. `go.mod` and `go.sum` remain unchanged at SHA-256
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
  and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`;
  `go.sum` remains 1,067 lines. Exact Go 1.26.7 module verification, build,
  full count-1 tests, full race tests, vet, and the 62-control launcher
  lifecycle suite pass.
- Fresh primary data still has 1,398 records, index SHA-256
  `cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
  and Last-Modified 2026-09-10T16:28:28Z. Its exact module entry contains only
  GO-2020-0019, fixed in v1.4.1, and selected-affecting GO-2026-6278. The
  latter remains unwithdrawn with the same pre-v1.5.3 published boundary and
  record SHA-256
  `bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.
  No new advisory or independent disqualifier appeared.
- The Gorilla WebSocket exception is non-transferable and valid only while
  exact v1.4.2 and its sole mvn-pom-mutator v0.2.3 edge remain unchanged, zero
  target packages load, the module remains runtime-unreachable, and no new
  advisory or independent disqualifier appears. Direct import/loading, runtime
  reachability, a version or incoming-edge change, or a new advisory or
  independent disqualifier expires the exception and requires a fresh
  dependency/product decision before merge. Do not add a direct edge, select
  v1.5.3 or another version to silence the advisory range, change
  mvn-pom-mutator or GoConvey, reopen the GopherJS edge guard, raise the Go
  floor, authorize a parent/removal group, select a patch/fork/replacement/
  unreleased commit/alternate path, move unrelated modules, manufacture a
  dependency commit, or request this same decision again while the guards
  hold. GopherJS, Enterprise Certificate Proxy, and GAX retain only their own
  separate guarded exceptions.

- Stop exact-path `github.com/grpc-ecosystem/go-grpc-middleware` for a fresh
  bounded product decision without changing dependency metadata. Fresh proxy,
  sumdb, `go-import`, strict Git, and GitHub evidence resolves the canonical
  public active unarchived non-fork Apache-2.0 repository and seven stable
  exact-path releases, v1.0.0 through proxy-latest v1.4.0. Current main is the
  different `/v2` module requiring Go 1.24; five post-v1.4.0 v1-branch commits
  are unreleased and not candidates.
- V1.4.0 is the highest stable release whose complete minimal production/test
  closure preserves Go 1.18, but no stable release qualifies. Every v1 release
  reproduces the same production retry cancellation defect under exact Go
  1.26.7 and Go 1.18.10: an already-canceled parent with zero backoff still
  reaches the invoker three times and returns Unavailable instead of Canceled.
  Every release also discards the cancel function from per-call retry timeout
  contexts; v1.1.0-v1.4.0 fail modern vet on that resource leak.
- Selected v1.0.0 has no release `go.mod` or deterministic complete test
  closure. V1.1.0-v1.2.1 suites fail on CN-only certificates without SANs;
  v1.2.2 has a Go 1.26 retry-suite server-state race. V1.4.0 full repeated
  suites panic when unclosed client work reaches a subtest-bound global logger,
  and `logging/settable` incorrectly nests variadic arguments. V1.3.0 passes
  repeats/race/cross-builds but retains the universal retry defects.
- Pinned API comparison records selected-to-later protobuf/marshaller
  incompatibilities and a v1.3.0-to-v1.4.0 change of all logging/kit public
  logger types from `github.com/go-kit/kit/log` to `github.com/go-kit/log`.
  Independent positive interceptor fixtures pass under both SDKs; separate
  fixtures reproduce the universal cancellation defect and v1.4.0 logger
  forwarding defect. Their SHA-256 values are
  `17bf692b527e101a0a3048e3d8da13b559a414db85541fb6cf98e409c8d1ab76`
  and `6c3d2ba53a405da762493a3a2db8480ff70a64cc1761a93f23bbc1116c3dd0c4`.
- Selected v1.0.0 exists solely through mvn-pom-mutator v0.2.3;
  `go mod why -m` is negative, source imports are zero, and zero target packages
  load. Exact v1.4.0 `go get` keeps loading at zero but changes four unrelated
  existing selections, adds 116 graph modules, and moves the graph from
  234/3,599 to 350/3,784 modules/edges. Tidy returns byte-identically to the
  base projection and restores inherited v1.0.0. Those unrelated MVS changes
  require a separate authorization and were not applied.
- Fresh primary vulnerability data has 1,398 records, index SHA-256
  `cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
  Last-Modified 2026-09-10T16:28:28Z, and no exact target record. Reviewed base
  and v1.4.0 project scans are identical at 30 module findings, 22 package
  findings, and 20 IDs/22 symbol traces, with zero target occurrence. Isolated
  v1.4.0 inherits old gRPC/x/net findings through testing helpers; none is an
  exact middleware advisory.
- The unchanged project stays at 234 modules, 3,599 edges, 429 complete-test
  entries, 197 module-backed packages, 41 loaded modules, 1,067 sum lines, and
  a 432-line tidy projection. Exact Go 1.26.7 full preflight, repeats, race,
  vet, lint, hermetic, cross-build, API, CLI, lifecycle, and audit controls pass;
  Go 1.18.10 retains only the two accepted `pkg/shell` wording failures.
  Accepted quality remains 27/27 Q0-Q2 PASS at L2. Evidence-manifest and
  decision-summary SHA-256 values are
  `bab1aa974f1182ba6d953831dcfccc0257b6f32f5064c477c68df170b7113b8e`
  and `a07237e54d3fac720966f6fc33767d6e98d61db6a2efc8a44de02f78316d3150`.
- Product direction recorded 2026-09-14: the user explicitly selected option 1
  with the recommended bounds. Retain exact selected, inherited, unloaded gRPC
  middleware v1.0.0 without dependency metadata changes. Accept only the
  already documented missing release module metadata/deterministic closure,
  TLS test-certificate failures in later candidates, universal retry
  cancellation and discarded-timer-cancel resource defects, and related
  recorded source/test qualification findings. The exception is target-
  specific and non-transferable and remains valid only while exact v1.0.0 and
  its sole mvn-pom-mutator v0.2.3 edge remain unchanged, zero target packages
  load, the module remains runtime-unreachable, and no new advisory or
  independent disqualifier appears. Direct import/loading, runtime
  reachability, a version or incoming-edge change, or a new advisory or
  independent defect expires the exception and requires a fresh dependency/
  product decision before merge. Do not add a direct edge, select v1.4.0 or
  another release, change mvn-pom-mutator or GoConvey, raise the Go floor,
  authorize a patch/replacement or parent/removal design, move unrelated
  selections, manufacture a dependency commit, transfer another exception, or
  request this same decision again while the guards hold.
- Every gRPC middleware and pre-existing exception guard was revalidated from
  clean decision HEAD `1444e88`, tree
  `71824e9b30e68396630e6bf1aca15d70d2a0fb64`, with a freshly unpacked,
  hash-verified exact Go 1.26.7 first in `PATH`, `GOENV=off`, `GOWORK=off`,
  `GOTOOLCHAIN=local`, `LC_ALL=C`, `LANG=C`, and no ambient `GOFLAGS`.
  Selected middleware v1.0.0 and its sole mvn-pom-mutator v0.2.3 edge remain
  unchanged. Gorilla WebSocket retains its sole mvn-pom-mutator edge, GopherJS
  its sole GoConvey v1.6.4 edge, and Enterprise Certificate Proxy and GAX their
  Viper v1.15.0 edges. All five `go mod why -m` results remain negative,
  repository source imports are zero, and the 429-entry complete test load has
  197 module-backed packages across 41 modules with zero packages from every
  guarded target. All remain runtime-unreachable.
- The no-change project remains 234 modules, 3,599 graph edges, and 1,067 sum
  lines. `go.mod` and `go.sum` remain byte-identical at SHA-256
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
  and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact Go 1.26.7 module verification, build, count-1, race, and vet pass. The
  62-control launcher lifecycle suite passes on an immediate unchanged rerun
  after the first nested source-archive signal-retention probe hit the already
  documented timing flake. Fresh primary vulnerability data remains byte-
  identical at 1,398 records, index
  SHA-256
  `cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
  and Last-Modified 2026-09-10T16:28:28Z. It has no exact middleware,
  GopherJS, Enterprise Certificate Proxy, or GAX record; Gorilla retains only
  GO-2020-0019 and selected-affecting unwithdrawn GO-2026-6278 at record
  SHA-256
  `bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.
  No new advisory or independent disqualifier appeared.

- Stop exact-path `github.com/grpc-ecosystem/go-grpc-prometheus` for a fresh
  bounded product decision without changing dependency metadata. Fresh proxy,
  sumdb, `go-import`, strict Git, and GitHub evidence resolves canonical
  `https://github.com/grpc-ecosystem/go-grpc-prometheus.git`, a public archived
  non-fork Apache-2.0 repository. The proxy exposes only stable v1.2.0. The
  v1.0/v1.1 tags are not semantic module versions; the 27-commit-later master
  and `draft-v2.0.0` branch are unreleased root-module identities. Master
  deprecates the project in favor of the different maintained module path
  `github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus`.
- V1.2.0 is the signed annotated tag object
  `502116f1a0a0c1140aab04fd3787489209b357d3`, commit
  `c225b8c3b01faf2899099b768856a9e916e5087b`, tree
  `5548a81e1bd9d5476b5b6d07ba43b0dca6f9c435`, released
  2018-06-04T12:28:56Z. The 25-file proxy archive is byte-identical to Git;
  strict Git verification passes and the Apache-2.0 license SHA-256 is
  `332c4f0a1a657a8937ad0caf7a335a31ec72343821481e364d8894f568c8b010`.
  The release contains 14 Go files/five packages/two tests/two generated files,
  no benchmarks/fuzz/testdata/build-tag branches/cgo/embed, and no `go.mod` or
  `go.sum`; the proxy synthesizes only the module declaration.
- No complete deterministic release closure qualifies. Current source-time
  resolution has 72 modules/135 edges, reaches Go 1.26 dependencies, and
  cannot compile native tests after removal of `prometheus.Handler`; Go 1.18.10
  cannot compile the modern closure. A 28-module historical reconstruction
  from post-release master preserves Go 1.18 and passes build/vet/race/cross-
  build, but uses metadata never present in the release and both SDKs' two
  count-10 native repeats fail because server-stream tests inspect metrics
  before an asynchronous server-handler update. The project-selected closure
  builds production under both SDKs but reaches a Go 1.19 Genproto dependency
  and its native tests also fail on the removed handler API.
- The independent 379-line behavior fixture, SHA-256
  `d6d2366bc2add1d9e7fac686de0b1b556a767841698599e061371ddab377058f`,
  passes count-1, two count-10 repeats, race, vet, cancellation/deadline cases,
  and 128-way supported concurrency under both SDKs. It proves that unary
  client success omits received-message counting while errors increment it,
  repeated terminal stream receives double-count handled metrics, and raw
  cancellation/deadline errors preserve identity but classify as Unknown.
  Abandoned streams/send failures/panics can leave started without handled;
  init globally registers eight counter vectors; histogram enablement ignores
  duplicate-registration errors; global enablement is unsynchronized; and nil
  callback/options/descriptor/info/streamer boundaries panic. The unary defect
  was fixed only after v1.2.0 on unreleased master.
- Selected v1.2.0 has the sole incoming mvn-pom-mutator v0.2.3 edge. Its
  `go mod why -m` is negative, source imports are zero, and zero target packages
  occur in the 429-entry load, so it is runtime-unreachable. Exact disposable
  `go get` nevertheless adds a direct root, grows the graph from
  234/3,599/1,067 modules/edges/sum lines to 346/3,755/1,080, and changes four
  unrelated existing selections: BigQuery, Datastore, Pub/Sub, and Envoy
  control-plane. Loading stays unchanged; tidy converges byte-identically to
  base and removes the root. No graph edit was authorized or applied.
- Fresh primary vulnerability data has 1,398 records, index SHA-256
  `cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
  Last-Modified 2026-09-10T16:28:28Z, no exact target record, and an empty OSV
  exact target/version response. Project scans remain 30 module findings,
  22 package findings, and 20 called IDs/22 symbol traces with zero target
  occurrence. Historical closure findings belong to inherited old Prometheus,
  gRPC, and x/net paths rather than an exact target advisory.
- The unchanged project remains 234 modules, 3,599 edges, 429 complete-test
  entries, 197 module-backed packages, 41 loaded modules, 1,067 sum lines, and
  a 432-line tidy projection. `go.mod`/`go.sum` SHA-256 remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact Go 1.26.7 full preflight, repeats, race, vet, lint, hermetic, cross-
  build, API/CLI, lifecycle, mutation, and audit-meta gates pass. Go 1.18.10
  retains only the two accepted `pkg/shell` wording failures. Accepted quality
  remains 27/27 Q0-Q2 PASS at L2. One final lifecycle rerun hit the already
  documented nested source-archive signal-retention timing flake at control 51;
  the immediate unchanged rerun passed all 62 controls. This is launcher test-
  design/timing evidence, not a target failure. The 263-entry evidence manifest
  and decision summary SHA-256 values are
  `b0fede88c14009f0932c1cd6698cda3c2da1fd96acd4e14c7c10729543a83f61`
  and `c5dfbf19a58b61f6ce0fba615bf0c6ba517dc3475eab524ddf2bc0fe3ae028be`.
- Product direction supplied 2026-09-14: the user explicitly selected option 1
  with the recommended bounds. Retain exact selected, inherited, unloaded gRPC
  Prometheus v1.2.0 without dependency metadata changes. Accept only the
  already documented missing release module metadata and deterministic complete
  closure, asynchronous native stream-test failures, unary client received-
  message inversion, repeated terminal stream-receive counting, raw
  cancellation/deadline classification as Unknown, and related archived,
  global-state, lifecycle, and qualification findings. The exception is target-
  specific and non-transferable and remains valid only while exact v1.2.0 and
  its sole mvn-pom-mutator v0.2.3 edge remain unchanged, zero target packages
  load, the module remains runtime-unreachable, and no new advisory or
  independent disqualifier appears. Direct import/loading, runtime
  reachability, a version or incoming-edge change, or a new advisory or
  independent defect expires the exception and requires a fresh dependency/
  product decision before merge. Do not add a direct edge, select an unreleased
  master/draft/tag or different maintained module path, change mvn-pom-mutator
  or GoConvey, raise the Go floor, authorize parent/removal or replacement/patch
  work, move unrelated selections, manufacture a dependency commit, transfer
  another exception, or request this same decision again while the guards hold.
- The option 1 decision was recorded from clean decision HEAD `547460d`, tree
  `c63bb971527dbe435f3878bb1b98311bbc49575f`, using a freshly unpacked,
  hash-verified exact Go 1.26.7 distribution first in `PATH`, `GOENV=off`,
  `GOWORK=off`, `GOTOOLCHAIN=local`, `LC_ALL=C`, `LANG=C`, and no ambient
  `GOFLAGS`. Gating selections remain exact gRPC Prometheus v1.2.0, gRPC
  middleware v1.0.0, Gorilla WebSocket v1.4.2, GopherJS
  `v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy v0.2.1,
  and GAX v2.7.0. Each has exactly one incoming edge at its selected version:
  the three gRPC/Gorilla targets from mvn-pom-mutator v0.2.3, GopherJS from
  GoConvey v1.6.4, and the two Google targets from Viper v1.15.0.
- All six `go mod why -m` results remain negative, repository Go source has
  zero guarded imports, and the unchanged 429-entry complete test load has 197
  module-backed packages across 41 modules with zero guarded target packages.
  All six targets remain runtime-unreachable. The project remains 234 selected
  modules, 3,599 graph edges, and 1,067 sum lines; `go.mod`/`go.sum` SHA-256
  values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact Go 1.26.7 module verification, build, count-1 tests, race, vet, and the
  62-control launcher lifecycle suite pass. No dependency implementation was
  created, and accepted quality remains 27/27 Q0-Q2 PASS at L2.
- Fresh primary vulnerability data remains byte-identical at 1,398 records,
  index SHA-256
  `cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
  and Last-Modified 2026-09-10T16:28:28Z. It has no exact gRPC Prometheus,
  gRPC middleware, GopherJS, Enterprise Certificate Proxy, or GAX record, and
  their exact-version OSV queries remain empty. Gorilla retains only its
  existing index records; selected-affecting GO-2026-6278 remains unwithdrawn
  at SHA-256
  `bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.
  No new advisory or independent disqualifier appeared. The next bounded P7
  group is selected exact-path `github.com/grpc-ecosystem/grpc-gateway
  v1.16.0`; its evaluation is a separate successor mission and was not begun
  in this decision-recording move.

- Stop exact-path `github.com/grpc-ecosystem/grpc-gateway` for a fresh bounded
  product decision without changing dependency metadata. Fresh proxy, sumdb,
  `go-import`, strict Git, and GitHub evidence resolves the canonical public
  active unarchived non-fork BSD-3-Clause repository and 55 proxy versions: 53
  stable and two prereleases. V1.16.0 is the latest exact-path stable release,
  lightweight tag/commit `094a6fe78b3ca888297d090185cdf30f0e42e157`,
  parent `e0a026aeb20c3eb0f32cbfdc537c9966634895b9`, tree
  `2283b306f85e2b97858d8575954fa0db23db4b5d`, dated
  2020-10-28T10:29:51Z with a GitHub-verified commit signature. The proxy zip
  SHA-256 is
  `377b03aef288b34ed894449d3ddba40d525dd7fb55de6e79045cdf499e7fe565`;
  proxy and Git sources are byte-identical, strict Git verification passes,
  sumdb agrees, and the license SHA-256 is
  `a15b1d1b168954c92ff7fb1620382418f7c72f4f4d251ee791d1098ad68ab0c4`.
  The later v1 branch is unreleased. Current main is the different
  `github.com/grpc-ecosystem/grpc-gateway/v2` module path, requires Go 1.26,
  and is not an eligible exact-path candidate.
- V1.16.0 declares Go 1.14 and its complete isolated production/test closure
  preserves Go 1.18. Exact Go 1.26.7/Go 1.18.10 resolve 33 modules, 136/135
  graph edges, 29 packages, 310/247 production entries, and 342/279
  complete-test entries; 12 external modules load and the maximum declared
  floor is Go 1.14. The release has 348 files, 179 Go files, 32 native test
  files, 164 tests, six benchmarks, one old go-fuzz target, 30 generated files,
  four commands, no Go examples/testdata, and no cgo/embed/symlinks. Pinned
  v1.15.2 -> v1.16.0 API comparison has zero incompatible changes and nine
  compatible Swagger-option additions.
- Native count-1, race, broad cross-builds, and targeted test compilation pass
  under both SDKs. Both independent count-10 repeats fail under both SDKs
  because codegenerator reuses a consumed buffer and runtime/Swagger tests
  leak package-global timeout/flag state. Vet fails because production
  `runtime/context.go` discards the cancel from `context.WithTimeout`, leaking
  timer/context resources; two additional findings are test-only `Fatalf`
  calls from goroutines. Every one of the 53 stable exact-path releases
  contains the production discarded-cancel defect. Adjacent serious candidate
  v1.15.2 reproduces the repeat and vet failures, so no stable exact-path
  release qualifies.
- The independent runtime fixture SHA-256 is
  `8f097645c3f31bde4b6e4a1b7addd7f92473d89bfe2647862a676d0df3ad3d17`.
  It passes count-10 and race under both SDKs while the source vet defect
  remains, and covers routing/path/method behavior, query handling and its
  global parser setter, marshaling, metadata/binary headers, malformed input,
  unary/stream/status translation, cancellation/deadlines, cleanup,
  concurrency, and nil/panic boundaries. Generator protocol output is
  byte-identical across repeats and SDKs at gateway response SHA-256
  `ff0829e0e0d3a76ee14a72d0e3e3941b85be7e6d70f5159c91e5b022615dc28c`
  and Swagger response SHA-256
  `d4e950e6163c8769b44d574fe8350e74b5f35435144048433f0272808a79d85b`;
  malformed input exits 255 with timestamped, stack-bearing nondeterministic
  diagnostics.
- Selected v1.16.0 has exactly two incoming selected-version edges: etcd/api/v3
  v3.5.1 and OTLP v0.7.0. `go mod why -m` is negative, repository imports are
  zero, and zero target packages occur in the 429-entry complete project load,
  so it is runtime-unreachable. Exact disposable `go get v1.16.0` only adds a
  direct root/source checksum, keeps 234 modules and v1.16.0 selected, grows
  edges 3,599 -> 3,600 and sum lines 1,067 -> 1,068, and leaves loading zero;
  tidy removes the root. Exact v1.15.2 instead cascades 234 -> 160 modules and
  3,599 -> 2,247 edges, removes the direct mvn-pom-mutator requirement, and
  makes the project unloadable. Neither projection was applied.
- Fresh primary vulnerability data remains 1,398 records at SHA-256
  `cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`
  and Last-Modified 2026-09-10T16:28:28Z, with no exact Gateway record and an
  empty exact OSV response. Govulncheck v1.8.0 reports 33 inherited module IDs,
  17 imported-package IDs, and eight called/reachable IDs in the v1.16.0
  production and test closure: GO-2020-0036, GO-2022-0956, GO-2023-1571,
  GO-2023-2153, GO-2024-2687, GO-2025-3372, GO-2026-4762, and GO-2026-6061.
  V1.15.2 is identical. Normal project scans remain 30 module IDs, 22 package
  IDs, and 20 called IDs with zero Gateway occurrence because no target package
  loads.
- The unchanged project remains 234 modules, 3,599 edges, 429 complete-test
  entries, 197 module-backed packages, 41 loaded modules, 1,067 sum lines, and
  a 432-line tidy projection. `go.mod`/`go.sum` SHA-256 remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact Go 1.26.7 project verify, repeats, race, vet, lint, hermetic,
  cross-build, and API/CLI gates pass. The Go 1.18.10 projection loads 366
  complete-test entries and retains only the two accepted Darwin `pkg/shell`
  closed-file wording assertions; every compatible package passes repeats and
  race. Focused quality remains 27/27 Q0-Q2 PASS at L2 with scorecard SHA-256
  `04039eb917cc8aa674c093e866b8fbfb303ab562a28dcbaf2414dba1c9d01c3a`,
  zero held/regressed/non-comparable rows, and seven improved rows. All eight
  mutation meta-suites, all 80/80 live kills, host/snapshot acceptance, and all
  15 audit controls pass. Docker acceptance is unavailable because the host
  Docker CLI lacks buildx. Full quality-wrapper/nested lifecycle runs retain
  the known signal-retention timing race; the final handoff lifecycle contract
  passes all 62 checks and every applicable quality stage passes independently.
  The 7,701-entry evidence manifest and decision
  summary SHA-256 values are
  `873610fe31b0da11ba8e91bfd74665be47428ef60c1b589244324ff9b7dc13d7`
  and `b73b6153778b43978c2c81738c246d0823a6731a995068ed5c0f8409e0dee354`.
- Product direction supplied 2026-09-14: the user explicitly selected option 1
  with its recommended bounds. Retain exact selected, inherited, unloaded gRPC
  Gateway v1.16.0 without dependency metadata changes. Accept only the
  documented universal discarded-cancel resource defect, native repeated-suite
  global and test-state failures, eight inherited called/reachable
  vulnerability IDs in the release-native closure, malformed-generator
  diagnostic nondeterminism, and related recorded qualification findings. The
  exception is target-specific and non-transferable and remains valid only
  while exact v1.16.0 and both selected-version incoming edges from etcd/api/v3
  v3.5.1 and OTLP v0.7.0 remain unchanged, zero target packages load, the module
  remains runtime-unreachable, and no new advisory or independent disqualifier
  appears. Direct import/loading, runtime reachability, a version or incoming-
  edge change, or a new advisory or independent defect expires the exception
  and requires a fresh dependency/product decision before merge. Do not add a
  direct edge, select another v1 release, change either parent, move to `/v2`,
  raise the Go floor, authorize parent-removal or architecture work, move
  unrelated selections, manufacture a dependency commit, transfer another
  exception, or request this same decision again while the guards hold.
- The Gateway option 1 decision was recorded from clean decision HEAD
  `b9e4b64d00d5593bb1a9d45e4506e55fd6eeb1f4`, tree
  `a5a2d5d6d7dafd058b6df74142463fae112561f3`, with a freshly unpacked,
  hash-verified exact Go 1.26.7 distribution first in `PATH`, `GOENV=off`,
  `GOWORK=off`, `GOTOOLCHAIN=local`, `LC_ALL=C`, `LANG=C`, and no ambient
  `GOFLAGS`. Gateway remains exact v1.16.0 through exactly two selected-version
  incoming edges from etcd/api/v3 v3.5.1 and OTLP v0.7.0. The earlier guarded
  selections and their sole recorded mvn-pom-mutator v0.2.3, GoConvey v1.6.4,
  or Viper v1.15.0 edges also remain exact.
- All seven `go mod why -m` results remain negative, repository Go source has
  zero guarded imports, and the unchanged 429-entry complete test load has 197
  module-backed packages across 41 modules with zero Gateway or earlier
  guarded-target packages. All seven targets remain runtime-unreachable. The
  project remains 234 selected modules, 3,599 graph edges, and 1,067 sum lines;
  `go.mod`/`go.sum` SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact Go 1.26.7 module verification, build, count-1 tests, race, vet, and the
  launcher lifecycle contract pass. No dependency implementation was created,
  and accepted quality remains 27/27 Q0-Q2 PASS at L2.
- Fresh primary vulnerability data remains byte-identical at 1,398 records,
  index SHA-256
  `cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
  and Last-Modified 2026-09-10T16:28:28Z. It has no exact Gateway, gRPC
  Prometheus, gRPC middleware, GopherJS, Enterprise Certificate Proxy, or GAX
  record, and exact target/version OSV queries for those six modules remain
  empty. Gorilla retains only its existing entries; selected-affecting
  GO-2026-6278 remains unwithdrawn at SHA-256
  `bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.
  No new advisory or independent disqualifier appeared. The next bounded P7
  group is selected exact-path `github.com/hashicorp/consul/api v1.18.0`; its
  evaluation is a separate successor mission and was not begun in this
  decision-recording move.

- Stop exact-path `github.com/hashicorp/consul/api` for a fresh bounded product
  decision without changing dependency metadata. Fresh proxy, sumdb,
  `go-import`, Git, and GitHub evidence resolves canonical public active
  unarchived non-fork `https://github.com/hashicorp/consul.git`, module
  subdirectory `api`, and MPL-2.0 at selected source. The proxy lists 96 v1
  entries: 84 stable and 12 prerelease/suffix entries. Selected v1.18.0 is
  lightweight tag `api/v1.18.0`, commit
  `13836d5ca84c71b35f201da06dd75f5ad6699c36`, parent
  `18dffc51de5587f8fba5f188a636d1864bf2b98a`, repository tree
  `1960701737d9ae23b0dc7d24a5087a4aa0541f83`, API subtree
  `b088d6696716b2cb8b61c824299dca34cab203f4`, dated
  2022-11-30T18:59:52Z with a verified commit signature. Proxy zip SHA-256 is
  `0dc6cfca8c71b05b3ba859726378d2ee611c15304fc85c2c030e3366ee068062`;
  proxy and Git source are byte-identical and sumdb agrees.
- All 83 retrievable stable v1 go.mod files were inspected. V1.18.0 is the
  highest exact-path stable release whose declaration and complete imported
  closure preserve Go 1.18; v1.18.1 is the first Go 1.19 release. Latest
  v1.34.5 declares Go 1.26.7 and has five incompatible API changes. The
  distinct `/api/v2 v2.0.0` path declares Go 1.26 and is out of scope. No
  redirect, mutated tag, prerelease, branch tip, alternate path, or major line
  was promoted.
- The selected release has 78 module files, 74 Go files, 40 production files,
  34 test files, two packages, and 217 tests. Its module archive is not
  standalone: go.mod replaces the SDK with `../sdk`, and native tests need 14
  certificate fixtures from `../test/client_certs`. Exact tagged API, SDK, and
  test subtrees resolve 56 modules under both SDKs; no imported module declares
  above Go 1.17. Verification, production build, vet, and Darwin/Linux/Windows
  production cross-builds pass under both SDKs, as do complete watch-package
  count-one, two count-ten repeats, race, and test cross-compilation.
- No stable candidate qualifies. Selected v1.18.0's root API test binary cannot
  link on Darwin or Linux under exact Go 1.26.7 because old x/net references
  `syscall.recvmsg`; v1.17.0 reproduces it. Its `consulent` test branch does
  not compile under either SDK because two default constants are absent. Four
  TLS tests fail with official Consul v1.14.2 because tagged certificates
  expired in 2023, though the other 213 tests pass repeats/race. Every selected
  response-metadata call site also discards parse errors, so malformed Consul
  metadata headers are silently accepted with defaults.
- The independent consumer fixture SHA-256 is
  `40e09f7b5f6003da6baf17cb9aa45d60cf8319d53732559186f2aeb6e0c7a90d`.
  It passes verify, count-one, two count-ten repeats, race, and vet under both
  SDKs while covering configuration, request/options encoding, auth/token/TLS,
  cancellation/deadline and error identity, no-retry behavior, cleanup,
  malformed input, nil/panic boundaries, and concurrency. Pinned API diffs
  show zero incompatibilities to adjacent v1.17.0/v1.18.1, while v1.34.5 has
  five.
- Selected v1.18.0 has the sole Viper v1.15.0 incoming edge. Its `go mod why
  -m` result is negative, repository imports are zero, zero target packages
  load, and it is runtime-unreachable. Exact disposable v1.18.0 `go get`
  manufactures eight indirect roots and 26 sum lines but tidy removes them
  byte-identically to base tidy. V1.17.0 downgrades Viper; v1.18.1 exceeds the
  floor and upgrades x/net/x/text; v1.34.5 raises the main Go declaration and
  moves unrelated selections. No projection was applied.
- Fresh primary data remains 1,398 records at SHA-256
  `cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`
  and Last-Modified 2026-09-10T16:28:28Z, with no target record; exact OSV
  queries for all serious candidates are empty. Base and selected project
  scans are identical at 30 module, 22 package, and 20 called IDs with zero
  target occurrence. The source production closure has inherited module-only
  GO-2026-5024; its test closure adds GO-2022-0603 in a loaded yaml test package,
  but neither closure reaches a vulnerable symbol.
- The unchanged project remains 234 modules, 3,599 edges, 429 complete-test
  entries, 197 module-backed packages, 41 loaded modules, 1,067 sum lines, and
  the recorded 432-line tidy projection. `go.mod`/`go.sum` SHA-256 remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact Go 1.26.7 applicable gates and the Go 1.18-compatible population pass;
  Go 1.18 full runs retain only the two accepted shell wording assertions.
  Accepted quality remains 27/27 Q0-Q2 PASS at L2. Buildx is now present, but
  Docker acceptance reaches the known Python 3.14 timestamp-control failure
  after building/testing the image. One completed-handoff lifecycle run passes
  all 62 controls; three later final-text runs pass outer controls 1-50 and
  reproduce only the known nested signal-log timing failure at control 51.
  These are target-independent. All earlier guards remain valid.
  The 9,780-entry evidence manifest and decision-summary SHA-256 values are
  `966263870f7529dcdd4713eba204dc1a8a00e17901a5f702557467b3c79ac46e`
  and `5957399af9cf4564b9864e4b6dea8c84e8912db3846e446c953345dadbddb08a`.
  Product direction supplied 2026-09-14: the user explicitly selected option 1
  with the recommended bounds. Retain exact selected, inherited, unloaded
  Consul API v1.18.0 without dependency metadata changes. Accept only the
  documented non-standalone release closure, Go 1.26 Unix native-test link
  incompatibility, broken `consulent` test branch, expired TLS fixtures,
  silently discarded response-metadata errors, inherited closure-only
  vulnerability findings, documented nil/panic and mutation boundaries, and
  related completed qualification findings. The exception is target-specific
  and non-transferable and remains valid only while exact v1.18.0 and its sole
  Viper v1.15.0 selected-version incoming edge remain unchanged, zero target
  packages load, the module remains runtime-unreachable, and no new advisory or
  independent disqualifier appears. Direct import/loading, runtime
  reachability, a version or incoming-edge change, or a new advisory or
  independent defect expires the exception and requires a fresh dependency/
  product decision before merge. Do not add a direct edge, select another v1
  release, move to `/api/v2`, change Viper, raise the Go floor, authorize parent
  modernization/removal or replacement architecture work, move unrelated
  selections, manufacture a dependency commit, transfer another exception, or
  request this same decision again while the guards hold.
- The Consul API option 1 decision was recorded from clean decision HEAD
  `474909583bfe6efbc5ae49b9ff1d68ec7402ee41`, parent
  `0ae1c16a2f35f66e128d45309967fce654d2e904`, and tree
  `a4656d4ae2e772914a72c8a8dd3f6548bba2d88c`. Consul API remains exact
  v1.18.0 through exactly one selected-version incoming edge from Viper
  v1.15.0. Gateway and every earlier guarded selection also retain their exact
  versions and recorded selected-version incoming edges.
- All eight guarded `go mod why -m` results remain negative, repository Go
  source has zero guarded imports, and the unchanged 429-entry complete test
  load has 197 module-backed packages across 41 modules with zero Consul API or
  earlier guarded-target packages. All eight targets remain runtime-
  unreachable. The project remains 234 selected modules, 3,599 graph edges,
  and 1,067 sum lines; `go.mod`/`go.sum` SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact Go 1.26.7 module verification, build, count-one tests, race, vet, and
  the launcher lifecycle contract pass. No dependency implementation was
  created, and accepted quality remains 27/27 Q0-Q2 PASS at L2.
- Fresh primary vulnerability data remains byte-identical at 1,398 records,
  index SHA-256
  `cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
  and Last-Modified 2026-09-10T16:28:28Z. It has no exact Consul API, Gateway,
  gRPC Prometheus, gRPC middleware, GopherJS, Enterprise Certificate Proxy, or
  GAX record, and exact target/version OSV queries for those seven modules
  remain empty. Gorilla retains only GO-2020-0019 and unwithdrawn GO-2026-6278;
  the latter's record SHA-256 remains
  `bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.
  No new advisory or independent disqualifier appeared.
- Stop exact-path `github.com/hashicorp/consul/sdk` for a fresh bounded product
  decision without changing dependency metadata. Fresh proxy, sumdb, Git, and
  GitHub evidence resolves canonical public active unarchived non-fork
  `https://github.com/hashicorp/consul.git`, module subdirectory `sdk`, and
  MPL-2.0 candidate source. The proxy lists 35 entries: 30 stable releases and
  five release candidates. Every proxy version maps to an exact `sdk/` Git tag;
  serious tags are lightweight and have no GitHub Release object. Strict Git
  verification passes. No redirect, fork, alternate path, different major
  line, prerelease, non-versioning tag, or branch head was promoted.
- Selected v0.8.0 is tag `sdk/v0.8.0`, commit
  `b1ee900870b9f04377d64936d4af452a5799b72c`, parent
  `a4a43460e51d66bc562fbea476844b6c13418543`, repository tree
  `8e328b47681316d07b7abbd89dff743dfbf1293d`, SDK subtree
  `dcbc915786cb1a668b9ca0f10bfd129b4526b4a8`, and proxy zip SHA-256
  `cf29fff6c000ee67eda1b8cacec9648d06944e3cdbb80e2e22dc0165708974c6`.
  Highest floor candidate v0.13.0 is commit
  `e297e3e75e357c532ced05123a4e9b48f5419393`, parent
  `5ce0132e9b13e9afa668ef7af91005ce2883f11f`, repository tree
  `7cad367c383f5e93a20f641ad0e342d6dac9d405`, and SDK subtree
  `7c45f4545bc9530db00fc92a15aaeecf84697216`.
- All 30 stable declarations and complete imported closures were inspected.
  V0.1.0 through v0.13.0 declare Go 1.12; v0.13.1 is the first Go 1.19
  release. V0.13.0's complete production/test closure declares no higher than
  Go 1.17 and runs under contained Go 1.18.10, so it is the highest real floor-
  eligible candidate. Latest v0.18.2 declares Go 1.26.7. No floor-ineligible
  release was promoted.
- No stable floor-eligible candidate qualifies. Every v0.1.0-v0.13.0 release
  implements the exported retry `Counter` with an equality-only stop, so a
  negative `Count` never terminates. Selected v0.8.0 and highest-candidate
  v0.13.0 independently reproduce retry attempts after a timer deadline and a
  nil `Stop` panic; freeport duplicate returns and a zero-block panic; malformed
  iptables UID/port/CIDR acceptance; and a `TempFile` descriptor that remains
  writable after cleanup. Retry has no cancellation boundary, freeport uses
  process-lifetime global state, and Linux iptables mutation is non-atomic and
  loses wrapped error identity.
- Selected test-server helpers also ignore their ready-timeout setting, use
  HTTP without request contexts or client deadlines, leak successful service/
  check response bodies, concatenate unescaped paths, and terminate callers
  through `log.Fatal`. V0.8.0 contains an additional double-`Wait`/unbuffered-
  goroutine cleanup defect. These are exact target behaviors, not inherited
  closure or harness findings.
- Both v0.8.0 and v0.13.0 expose `freeport`, `iptables`, `testutil`, and
  `testutil/retry`; there are no commands, examples, benchmarks, fuzz targets,
  testdata trees, generated files, cgo, embeds, or go:generate directives.
  Exact source passes verify, build, vet, and production/test cross-builds under
  exact Go 1.26.7 and Go 1.18.10. Darwin native tests require only a scratch
  overlay for the managed sandbox's denied `/usr/sbin/sysctl`; with it both
  candidates pass count-one, two count-ten repeats, and race under both SDKs.
  Independent fixtures reproduce every listed defect under both SDKs.
- Pinned API comparison finds incompatible testing-interface and token changes
  from v0.8.0 to v0.13.0, only a compatible `R.Logf` addition through v0.13.1,
  and numerous incompatible retry/testing/iptables changes through v0.18.2.
- Selected v0.8.0 has exactly one selected-version incoming edge, from the
  historical Consul API v1.12.0 graph vertex. Its why result is negative,
  repository imports are zero, zero target packages occur in the 429-entry
  complete project load, and it is runtime-unreachable. Disposable exact gets
  retain zero load but manufacture a direct indirect root. V0.13.0 also adds
  go-version; v0.13.1 exceeds the floor; latest raises the main Go declaration
  and moves unrelated selections. Base and the first three candidate tidy
  projections converge to selected v0.8.0 and the same base tidy result. No
  projection was applied.
- Fresh primary data remains 1,398 records at SHA-256
  `cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`
  and Last-Modified 2026-09-10T16:28:28Z, with no Consul SDK record. Exact OSV
  queries for v0.8.0, v0.13.0, v0.13.1, and v0.18.2 are empty. Isolated older
  closures contain only inherited x/sys findings: v0.8.0 has GO-2026-5024 and
  GO-2022-0493, while v0.13.0/v0.13.1 have GO-2026-5024; no Darwin called
  symbol is reached. Project scans are identical for every candidate with zero
  target package, symbol, test-symbol, or reachable trace.
- The unchanged project remains 234 modules, 3,599 edges, 429 complete-test
  entries, 197 module-backed packages, 41 loaded modules, 1,067 sum lines, and
  the recorded 432-line tidy projection. `go.mod`/`go.sum` SHA-256 remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact Go 1.26.7 applicable gates pass, all earlier exception guards remain
  valid, no dependency implementation was created, and accepted quality
  remains 27/27 Q0-Q2 PASS at L2. Product direction supplied 2026-09-14: the
  user explicitly selected option 1 with the recommended bounds. Retain exact
  selected, inherited, unloaded Consul SDK v0.8.0 without dependency metadata
  changes. Accept only the completed retry, freeport, iptables, TempFile, test-
  server, API, inherited closure-only x/sys, and related qualification
  findings. The exception is SDK-specific and non-transferable and remains
  valid only while exact v0.8.0 and its sole selected-version incoming edge
  from the historical Consul API v1.12.0 graph vertex remain unchanged, zero
  SDK packages load, the module remains runtime-unreachable, and no new
  advisory or independent disqualifier appears. Direct import/loading, runtime
  reachability, a version or incoming-edge change, or a new advisory or
  independent defect expires the exception and requires a fresh dependency/
  product decision before merge. Do not add a direct SDK edge, select v0.13.0
  or a later release, change or remove the historical parent, raise the Go
  floor, authorize dependency modernization, alter Viper, move unrelated
  modules, manufacture a dependency commit, transfer another exception, or
  request this same decision again while the guards hold.
- The option 1 decision was recorded after guard-only revalidation from clean
  decision HEAD `2346fcecd1f9f464697d51c3fa223b6709e84972`, parent
  `0f9204ef2a2f53cb53495c2302d16f0b36428d6d`, and tree
  `88ab43a3fc01557456853e8a11d2836818b52fcf`. Consul SDK remains exact
  v0.8.0 through exactly one selected-version incoming edge from the historical
  Consul API v1.12.0 graph vertex. Every earlier guarded selection also retains
  its exact version and recorded selected-version incoming edge.
- All nine guarded `go mod why -m` results remain negative, repository Go
  source has zero guarded imports, and the unchanged 429-entry complete test
  load has 197 module-backed packages across 41 modules with zero SDK or
  earlier guarded-target packages. The production load also contains zero
  guarded packages. All nine targets remain runtime-unreachable. The project
  remains 234 selected modules, 3,599 graph edges, and 1,067 sum lines;
  `go.mod`/`go.sum` SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  No dependency implementation or metadata change was created.
- Fresh primary vulnerability data remains byte-identical at 1,398 records,
  index SHA-256
  `cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
  and Last-Modified 2026-09-10T16:28:28Z. It has no exact SDK or other new
  guarded-target record. Exact SDK v0.8.0 OSV results remain empty. Gorilla
  retains only its recorded module entries, including selected-affecting
  GO-2026-6278; that record remains unwithdrawn at SHA-256
  `bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`
  with its existing GHSA alias. No new advisory or independently observed
  disqualifier appeared during the guard-only revalidation.
- The bounded exact-path Hashicorp Errwrap evaluation is complete. The proxy
  exposes only stable v1.0.0 and v1.1.0, both with stdlib-only source/test
  closures that preserve Go 1.18; v1.1.0 is therefore the highest floor-
  eligible stable candidate. Canonical `go-import`, proxy, sumdb, Git, and
  GitHub evidence resolves public, active, unarchived, non-fork
  `https://github.com/hashicorp/errwrap`, MPL-2.0, no GitHub Release objects,
  no retractions or module deprecation, and no `/v2` module line. V1.0.0 is
  lightweight tag/commit
  `8a6fb523712970c966eefc6b39ed2c5e74880354`, parent
  `d6c0cd88035724dd42e0f335ae30161c20575ecc`, tree
  `9863613ad8fe960290d1631d84ede660af5f0746`, dated
  2018-08-24T00:39:10Z; its commit signature is not forge-verified. V1.1.0 is
  lightweight tag/verified merge commit
  `7b00e5db719c64d14dd0caaacbd13e76254d02c0`, parents v1.0.0 and
  `96a78ad11c51762df122b738e6f4f30f58e03d8d`, tree
  `aefd62cf8e9549e154a65afb8b4524b7b6ed5f2e`, dated
  2020-07-14T15:51:01Z. Proxy and Git bytes match, strict Git verification
  passes, and v1.0.0 is an ancestor of v1.1.0. Unreleased master now declares
  Go 1.24 and is not a stable candidate.
- Each release contains one package, one production file, one test file, and
  no command, example, benchmark, fuzz target, testdata, generated file,
  build-tag/platform branch, cgo, embed, go:generate directive, symlink, or
  non-stdlib dependency. Both verify, build, pass native count-one, two
  count-ten repeats, race, vet, and production/test cross-compilation for
  Darwin amd64, Linux amd64/arm64, Windows amd64, and js/wasm under exact Go
  1.26.7 and contained Go 1.18.10. Pinned API diff finds the same nine
  exported functions plus `WalkFunc` and `Wrapper` in both releases with no
  incompatible API change; v1.1.0 adds standard single-error `Unwrap` behavior
  and deprecates `Wrapf` without changing the exported declaration set.
- Neither stable release qualifies. The exported type-matching helpers compare
  `reflect.Type.String()` rather than concrete type identity, so distinct
  types with the same package/type spelling in different import paths falsely
  match. Both releases also fail to traverse Go's standard `Unwrap() []error`
  graph, including `errors.Join`, under Go 1.26.7; v1.0.0 additionally does not
  expose its chain through standard `Unwrap() error`. Independent fixtures
  reproduce these defects under both SDKs where applicable and characterize
  deterministic formatting/order and identity, explicit nil outer-error and
  nil callback panics, outer/inner aliasing, isolated returned slices,
  supported immutable concurrent reads, and allocations. The module has no
  global state or resource ownership; recursive walking has no cycle guard.
- Selected v1.0.0 exists through exactly one selected-version incoming edge
  from `github.com/hashicorp/go-multierror v1.1.0`; the graph also retains the
  historical v1.0.0 parent edge. Its shortest path begins at
  mvn-pom-mutator v0.2.3, passes historical Viper v1.10.1, Serf v0.9.6, and
  go-multierror v1.1.0. `go mod why -m` is negative, repository imports are
  zero, and zero target packages occur in either project load, so Errwrap is
  runtime-unreachable. Disposable exact gets manufacture only a direct root
  and target checksum rows: v1.0.0 retains 234 modules/429 complete-test
  entries and adds one graph edge and one sum line; v1.1.0 does the same with
  two sum lines. Both preserve every unrelated selection, remain unloaded,
  and tidy back to selected v1.0.0 and the base projection. No projection was
  applied.
- Fresh primary vulnerability data remains 1,398 records at SHA-256
  `cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`
  and Last-Modified 2026-09-10T16:28:28Z, with no Errwrap record; exact
  v1.0.0/v1.1.0 OSV responses are empty. Exact Go 1.26.7 source scans find no
  module, package, symbol, or test-symbol vulnerability. Go 1.18.10 source
  scans report only that SDK's old standard-library findings, not Errwrap.
  Base and v1.1.0 project scans are byte-identical, with zero target occurrence
  at module, package, symbol, test-symbol, or reachable-trace level. Earlier
  guarded advisories remain unchanged.
- The project remains byte-for-byte unchanged at 234 modules, 3,599 graph
  edges, 429 complete-test entries, 197 module-backed packages across 41
  loaded modules, 1,067 sum lines, and the recorded 432-line unapplied tidy
  projection. `go.mod`/`go.sum` SHA-256 remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact Go 1.26.7 applicable gates pass, contained Go 1.18.10 retains only its
  two accepted Darwin shell wording failures, all nine earlier exception
  guards remain valid, and accepted quality remains 27/27 Q0-Q2 PASS at L2.
  No dependency implementation was created. On 2026-09-15 the user explicitly
  selected option 1 with the recommended
  bounds. Retain exact selected, inherited, unloaded Errwrap v1.0.0 without
  changing `go.mod` or `go.sum`. Accept only the completed concrete-type
  collision, absent standard single/multi-error traversal, nil/panic,
  aliasing, allocation, recursion, API, and related qualification findings.
  The exception is Errwrap-specific and non-transferable. It is valid only
  while exact v1.0.0 and its sole selected-version incoming edge from
  go-multierror v1.1.0 remain unchanged, zero Errwrap packages load, the module
  remains runtime-unreachable, and no new advisory or independent defect
  appears. Direct import/loading, runtime reachability, a version or incoming-
  edge change, or a new advisory or independent defect expires the exception
  and requires a fresh Errwrap dependency and product decision before merge.
  Do not add a direct edge, select v1.1.0 or another version, remove or change
  the historical parent chain, change go-multierror, Serf, Viper or mvn-pom-
  mutator, raise the Go floor, authorize replacement or modernization, move
  unrelated selections, manufacture a dependency commit, transfer another
  exception, or request this same choice again while all guards hold.
  The decision was recorded after guard-only revalidation from clean decision
  HEAD `f77cb8a75c6813b28656c85b654d17d7df464731`, parent
  `a2f1bc98dc6b978d802d42f7c36bae96627bda33`, tree
  `9dd5c04d0cc3cdebf482b831faef8809ea25643e`. Exact v1.0.0 retains its
  sole selected-version go-multierror v1.1.0 incoming edge; the historical
  go-multierror v1.0.0 request also remains. All ten guarded versions and
  recorded incoming edges remain unchanged. All ten why results remain
  negative, repository Go source has zero guarded-path occurrences, and
  production and 429-entry complete-test loads have zero guarded packages.
  The complete-test load retains 197 module-backed entries across 41 modules,
  so all guarded modules remain runtime-unreachable.
- Guard-only project measurements remain 234 selected modules, 3,599 graph
  edges, and 1,067 sum lines. `go.mod`/`go.sum` SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact Go 1.26.7 module verification, build, count-one tests, race tests, and
  vet pass. The launcher lifecycle check recognizes the reciprocal successor.
  No dependency implementation or metadata change was created.
- Fresh primary vulnerability data now has 1,399 records at SHA-256
  `033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`
  and Last-Modified 2026-09-15T18:54:35Z. It has no Errwrap or other new
  guarded-target record. Exact Errwrap v1.0.0/v1.1.0 OSV results remain
  empty. Gorilla retains only GO-2020-0019 and unwithdrawn GO-2026-6278;
  that record's SHA-256 remains
  `bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.
  No new advisory or independently observed defect appeared.

- The bounded exact-path Hashicorp go-cleanhttp evaluation is complete. Retain
  qualified exact selected `github.com/hashicorp/go-cleanhttp v0.5.2` without
  changing dependency metadata. It is the highest exact-path stable release,
  its complete minimal production/test closure preserves Go 1.18, and its
  relevant source, API, behavior, tests, and project effects satisfy the
  existing contracts. No exception or product-risk acceptance was needed.
- Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves canonical
  public active unarchived non-fork
  `https://github.com/hashicorp/go-cleanhttp.git`, MPL-2.0, one protected
  `master` branch, and exactly three stable versions: v0.5.0, v0.5.1, and
  v0.5.2. V0.5.2 is the sole GitHub Release and proxy `@latest`, published
  2021-02-03T18:51:13Z. There is no prerelease, `/v2` line, alternate module
  path, redirect, retraction, deprecation, or qualified later stable identity.
  Current unreleased master declares Go 1.24 and was not promoted.
- The three release tags are lightweight, unsigned commit objects and form
  linear ancestry v0.5.0 -> v0.5.1 -> v0.5.2 -> master. V0.5.0 is commit
  `e8ab9daed8d1ddd2d3c4efba338fe2eeae2e4f18`, tree
  `4e73aae6cb561e640cd9619ec84785d956638d1b`, and proxy zip SHA-256
  `f5ca0da7f2432a3e961064529f0af8de1113b854bf239a880d16f89424538b44`.
  V0.5.1 is commit `eda1e5db218aad1db63ca4642c8906b26bcf2744`, tree
  `d7c6740e1c415649cbf335bed109358e84b7e552`, and zip SHA-256
  `e3cc9964b0bc80c6156d6fb064abcb62ff8c00df8be8009b6f6d3aefc2776a23`.
  Selected v0.5.2 is commit
  `6d9e2ac5d828e5f8594b97f88c4bde14a67bb6d2`, parent
  `d3fcbee8e1810ecee4bdbf415f42f84cfd0e3361`, tree
  `c1160f09cedce00dc3ef7b06169ac45cad1c12e8`, and zip SHA-256
  `e9f3dcfcb33172ba499b4f8e888169252d7f1e072082182124a6e2053523f7df`.
  Proxy and Git release bytes agree exactly.
- Every release contains exactly seven regular files, one package, three
  production Go files, and one test file. There are no commands, examples,
  benchmarks, fuzz targets, testdata, generated files, build tags, platform
  branches, cgo, embeds, go:generate directives, symlinks, or non-standard-
  library dependencies. V0.5.0/v0.5.1 have no Go directive and v0.5.2 declares
  Go 1.13. Exact Go 1.26.7 resolves 184 production and 209 complete-test
  entries; contained Go 1.18.10 resolves 123 and 147. The target is the sole
  external selected module in each isolated closure.
- Selected v0.5.2 passes verification, build, native count-one, two independent
  count-ten repeats, race, vet, and production/test cross-builds for Darwin
  amd64/arm64, Linux amd64/arm64, Windows amd64, and js/wasm under both SDKs.
  V0.5.0/v0.5.1 pass build, vet, and every cross-build but reproducibly fail
  count-one, both repeats, and race under both SDKs because their own tests put
  raw LF, CR, and NUL bytes in request URLs; current Go rejects those URLs as
  invalid. V0.5.2 percent-encodes those test inputs, retains v0.5.1's nil-
  request/nil-next handler fix, and adds `ForceAttemptHTTP2`.
- Pinned API comparison is identical in both directions across all releases:
  `DefaultTransport`, `DefaultPooledTransport`, `DefaultClient`,
  `DefaultPooledClient`, `PrintablePathCheckHandler`, and exported
  `HandlerInput.ErrStatus`. All three export blobs have SHA-256
  `e96e3d2d100e788eac2ba58bf3e1b07237bedfd3d2ab5969fca30da56b946c64`.
- Independent fixtures under both SDKs verify fresh transport/client identity,
  exact transport defaults, no package-owned global mutation, deterministic
  printable/non-printable path routing, Unicode and decoded-control handling,
  HTTP/1 keepalive disablement/reuse and idle cleanup, HTTP/2 negotiation,
  immutable-handler concurrency, and allocations. They characterize the
  supported boundaries: a zero `ErrStatus` is mutated to 400 at construction;
  the handler retains that pointer; later caller mutation changes behavior and
  concurrent mutation is unsupported; nil request and nil next are no-ops;
  a non-nil request with nil URL and an invalid nonzero status panic through
  standard-library preconditions; malformed UTF-8 path bytes become printable
  replacement runes; pooled callers own idle-resource cleanup. The fixture
  SHA-256 is
  `1b3aa948358fdbfd136041fc6c930ba904ff8bc6260bd915ea9af50f40dd73a2`.
  The open tuning discussion about `GOMAXPROCS+1` remains configurable through
  the returned transport and is not a correctness disqualifier.
- Selected v0.5.2 exists through three selected-version graph edges: Viper
  v1.15.0, the historical Viper v1.10.1 vertex, and
  `sagikazarmark/crypt v0.4.0`. Lower historical parents request v0.5.1 and
  v0.5.0. The shortest path is main -> Viper v1.15.0 -> go-cleanhttp v0.5.2.
  `go mod why -m` is negative, repository imports are zero, and production and
  complete-test loads contain zero target packages, so it is runtime-
  unreachable. An exact disposable v0.5.2 get preserves 234 modules, 429
  complete-test entries, every unrelated selection, and zero load; it adds
  only a direct indirect root edge and one source checksum line, yielding
  3,600 edges and 1,068 sum lines. Its go.mod/go.sum SHA-256 pair is
  `55698b6242da86f1e456b0dcf9cadc52f72f2b419500a6196c04b0117281b063` /
  `e1b800f07e83f6d7b9e566da23d0de7841f65ccb9866fdd547b90214f0d88cd5`.
  Tidy returns exactly to the base projection. No direct edge was added.
- Exact disposable v0.5.0/v0.5.1 gets force broad unrelated downgrades, remove
  the required mvn-pom-mutator root, and fail project package loading. They
  leave only 182/180 total selected modules and 2,311/2,953 graph edges.
  Their tidy results restore the missing parent but converge to a different
  historical Viper projection. Those prohibited projections were not applied.
- Fresh primary vulnerability data contains 1,399 records at SHA-256
  `033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`
  and Last-Modified 2026-09-15T18:54:35Z with no go-cleanhttp record. Exact OSV
  responses for all three releases are empty. Go 1.26.7 isolated module,
  package, symbol, and test-symbol scans have zero findings. Go 1.18.10 reports
  only that old SDK's standard-library population: 89 IDs and respectively
  90/154/186/335 finding paths; target frames in seven production IDs/32 paths
  and 45 test IDs/178 paths call into the vulnerable standard library, but no
  advisory is assigned to go-cleanhttp. Base and disposable v0.5.2 project
  scans have identical normalized populations at 30 IDs and 30/52/74 Darwin
  plus 75 Windows finding paths, with zero target SBOM package, symbol, test-
  symbol, or reachable-trace occurrence.
- The project remains byte-for-byte unchanged at 234 modules, 3,599 graph
  edges, 429 complete-test entries, 197 module-backed packages across 41
  loaded modules, 1,067 sum lines, and the recorded 432-line tidy projection.
  `go.mod`/`go.sum` SHA-256 remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact Go 1.26.7 applicable gates, all 17 script meta-tests, and all 15 audit
  meta-controls pass. Contained Go 1.18.10 retains only the two accepted Darwin
  shell wording failures. Accepted quality remains 27/27 Q0-Q2 PASS at L2,
  and no dependency implementation or metadata commit was created. The 412-
  entry source/evaluation evidence manifest SHA-256 is
  `40b0a6fc11f4378aed62586bfc8bdd94bef33b0a0ba88051dff51912807137c1`.
- Retain this qualified result while exact v0.5.2 and its three selected-
  version incoming edges remain unchanged, zero target packages load, it stays
  runtime-unreachable, and no new advisory or independently disqualifying
  behavior appears. Direct import/loading, runtime reachability, a target
  version or incoming-edge change, or a new advisory/defect requires a fresh
  bounded go-cleanhttp decision. This guard is target-specific and transfers
  no earlier exception.

- The bounded exact-path Hashicorp go-hclog evaluation is complete: no stable
  release satisfies the current behavior and compatibility contracts. On
  2026-09-16 the user explicitly selected option 1 with the recommended bounds,
  retaining exact selected, inherited, unloaded v1.2.0 under the target-specific
  exception recorded below. No dependency metadata changed.
- Canonical identity is public active unarchived non-fork MIT repository
  `https://github.com/hashicorp/go-hclog.git`. The proxy exposes 31 stable
  releases from v0.7.0 through `@latest` v1.6.3, no retraction, deprecation,
  redirect, or `/v2` line. Selected v1.2.0 is commit
  `b6b55671f4e5b82443139ee3e9f4417603c4cd72`, tree
  `45f0da0a3b7eb001522fc9891504cee03d951ef0`; latest v1.6.3 is commit
  `d12136aa2e51933c460084f5083b6d5bb9d41960`, tree
  `26c5c9247ad5b3a6f930ea2b1fc0a1ee532a7804`. Proxy/Git source matches and
  strict verification passes. Non-semver aliases, nested-module tags, branch
  heads, and unreleased Go-1.25 main were not promoted.
- Every stable release declares no more than Go 1.13. Imported selected/latest
  source/test closures declare no more than Go 1.17 and execute under exact Go
  1.26.7 and contained Go 1.18.10. Serious releases contain one package, 12
  production files, seven tests, two platform color branches, and one
  benchmark, with no commands, examples, fuzz, testdata, generated/cgo/embed/
  generator files, or symlinks. Serious native repeated/race/vet and applicable
  cross-build gates pass under both SDKs.
- V1.3.1 is the last API-compatible candidate: v1.2.1/v1.2.2 are API-identical
  to selected, and v1.3.0/v1.3.1 add only
  `LoggerOptions.ColorHeaderAndFields`. V1.4.0+ incompatibly add
  `Logger.GetLevel` to the exported interface.
- The independent fixture SHA-256 is
  `b62af8023aa4e9d2b84a20db69627992998bfa29e7c268832025b345dcf5bb26`.
  Selected and every compatible patch fail five contracts under both SDKs:
  NaN silently drops JSON records; caller fields overwrite `@message` and
  `@level`; absent-sink deregistration underflows and disables a future sink;
  self-deregistration deadlocks under the held sink mutex; and
  `SetDefault(nil)` violates the documented non-nil `FromContext` result.
  Source history retains all five through latest. Latest additionally races in
  `SyncParentLevel` concurrent level/epoch access under both SDKs. Determinism,
  routing, exact error identity, nil/panic, mutation/aliasing, allocations,
  supported concurrency, and caller-owned cleanup were characterized.
- Selected v1.2.0 has one selected-version Viper v1.15.0 edge; historical
  requests are v1.0.0 and v0.12.0. Its why result is negative, repository Go
  occurrences and target package loads are zero, and it is runtime-unreachable.
  Disposable exact serious gets preserve 234 modules, every unrelated
  selection, 429 complete-test entries, and zero load; tidy restores selected
  and the exact base projection. Exact v1.1.0 forces Viper v1.15.0 -> v1.10.1,
  so lower releases require a separate parent decision. No projection applied.
- Fresh primary data remains 1,399 records at SHA-256
  `033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`
  and Last-Modified 2026-09-15T18:54:35Z with no target record; exact serious
  OSV responses are empty. Selected has zero Go-1.26 source findings. V1.3.1+
  has module-only GO-2026-5024 through old `x/sys` with no vulnerable package,
  called symbol, target frame, or project reachability. Project scans retain 30
  IDs and zero target occurrence.
- The unchanged project remains 234 modules, 3,599 graph edges, 429 complete-
  test entries, 197 module-backed entries across 41 loaded modules, 1,067 sum
  lines, and the 432-line tidy projection. Exact module hashes remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Applicable Go-1.26/project/API/CLI/lint/acceptance/meta/audit controls pass;
  Go 1.18 retains only the two accepted shell wording failures. An independent
  lifecycle run passes 62/62; full preflight retains its known nested timing
  race. Accepted quality remains 27/27 Q0-Q2 PASS at L2. The 1,328-entry
  evidence manifest SHA-256 is
  `e6647301b469e2c915ca07bd1f2b0f7856e6b55df75c76a350943cc780783a9f`.
  Product direction recorded 2026-09-16: the user explicitly selected option 1
  with the recommended bounds. Retain exact selected, inherited, unloaded
  go-hclog v1.2.0 without changing `go.mod` or `go.sum`. Accept only the five
  completed common behavior findings, characterized nil/panic, aliasing,
  allocation, concurrency, resource, API, closure-vulnerability, and related
  qualification findings. The exception is go-hclog-specific and non-
  transferable. It remains valid only while exact v1.2.0 and its sole selected-
  version incoming edge from Viper v1.15.0 remain unchanged, zero go-hclog
  packages load, the module remains runtime-unreachable, and no new advisory
  or independent defect appears. Direct import/loading, runtime reachability,
  a version or incoming-edge change, or a new advisory or independent defect
  expires the exception and requires a fresh go-hclog dependency and product
  decision before merge. Do not add a direct edge, select v1.3.1 or another
  version, promote v1.4.0+, downgrade or otherwise change Viper, remove/change
  a parent, patch/fork source, authorize a wrapper or replacement architecture,
  raise the Go floor, move unrelated selections, alter another guarded
  dependency, manufacture a dependency commit, transfer another exception, or
  request this same choice again while all guards hold.
- Guard-only decision revalidation from clean HEAD
  `d536b6832a85123dc49cf4fbcf4876f4083458d4`, parent
  `f7f3a65d645ce4ce4caf8d1ac3e1da77ccd19727`, tree
  `cab76bf3e4948d54fdbf032ba0c4a788c7cafda7`, preserves exact go-hclog
  v1.2.0 and its sole selected-version Viper v1.15.0 incoming edge. All 12
  guarded modules retain their recorded selections and incoming edges, all 12
  `go mod why -m` results remain negative, repository imports remain zero, and
  the 355-entry production and 429-entry complete-test loads contain zero
  guarded packages. The project remains 234 modules, 3,599 graph edges, 197
  module-backed complete-test entries across 41 loaded modules, and exact
  `go.mod`/`go.sum` SHA-256 values
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Runtime unreachability therefore remains intact.
- Fresh primary vulnerability data remains byte-identical at 1,399 records,
  SHA-256
  `033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`,
  and Last-Modified 2026-09-15T18:54:35Z. Go-hclog v1.2.0 and every guarded
  target except Gorilla have empty exact-version OSV responses; Gorilla retains
  only the recorded GO-2026-6278/GHSA-w67g-5rqw-f597 result and the primary
  index's existing GO-2020-0019 entry. No new advisory or independent defect
  appeared. The following bounded exact-path go-immutable-radix evaluation
  supersedes the former queue boundary without changing this decision.

- The bounded exact-path Hashicorp go-immutable-radix evaluation is complete:
  no stable release satisfies every behavior and current-API contract. Exact
  selected v1.3.1 remains inherited and unloaded, but is not qualified or
  accepted. No dependency metadata changed; P7 stops for the prepared bounded
  product decision.
- Canonical identity is public active unarchived non-fork MPL-2.0 repository
  `https://github.com/hashicorp/go-immutable-radix.git`. The exact v1 proxy
  exposes v1.0.0-v1.3.1, with v1.3.1 at `@latest`; there is no retraction,
  module deprecation, redirect, or newer exact-path stable. V2 releases and
  current master belong to distinct module path `/v2`. Selected v1.3.1 is
  verified commit `49d1d02c49a783de548d1ba8ae8fde466a20b9e6`, parent
  `f63f49c0b598a5ead21c5015fb4d08fe7e3c21ea`, tree
  `b7c0058c779f65f8fb81f1db6e17023cbe8bfe0c`; proxy/Git bytes match.
- All five releases omit a `go` directive and resolve the same minimal target,
  production `golang-lru v0.5.0`, and test-only `go-uuid v1.0.0` closure. The
  imported closure preserves Go 1.18. Every release passes verification,
  build, count-one, two count-ten repeats, race, vet, and Darwin AMD64, Linux
  AMD64/ARM64, Windows AMD64, and js/wasm production/test cross-compilation
  under exact Go 1.26.7 and Go 1.18.10.
- Selected source has one package, six production files, three tests, and no
  command, example, benchmark, modern fuzz target, testdata, generated file,
  platform/build-tag branch, cgo, embed, generator, or symlink. API additions
  through v1.3.0 are compatible; v1.3.1 makes `ReverseIterator` non-comparable.
  V1.0.0 lacks six current APIs and is not a compatible downgrade.
- The behavior fixture SHA-256 is
  `778d70a8c17156da17a0c6f32cf3f71b1b0033e4abb646d28516b4d3b089b32b`.
  It passes deterministic ordering, CRUD/snapshots, prefix deletion,
  transaction cloning, watch lifecycle, arbitrary-byte keys, immutable
  concurrent reads, nil/panic, allocation, and aliasing checks under both
  SDKs. Keys alias caller/internal byte storage and must remain immutable;
  values retain caller identity. Transactions are not thread-safe; independent
  clones are. Package watch/LRU resource boundaries are characterized.
- The minimal blocker fixture SHA-256 is
  `8e9476f3989e6ab4bed98ec99cc4409401b0d727982362b5de2374e4b92b6318`.
  On an empty tree, missing `Iterator.SeekPrefix` followed by
  `Iterator.SeekLowerBound` nil-dereferences and panics under both SDKs.
  V1.1.0-v1.3.1 all fail and open upstream issue #50 records the defect.
  V1.1.0-v1.3.0 also retain the older issue-#28 prefix-key panic; v1.0.0 lacks
  the API and removes current exports. No exact-path stable release qualifies.
- Selected v1.3.1 has selected-version incoming edges from Viper v1.15.0,
  historical Viper v1.10.1, and `sagikazarmark/crypt v0.4.0`. Its why result
  and source import search are negative and target package loads are zero, so
  it is runtime-unreachable. A disposable v1.3.1 get preserves all 234
  selections and zero load but manufactures target/golang-lru roots and two
  checksum lines; tidy returns to the exact base projection. Exact
  v1.0.0-v1.3.0 gets downgrade Viper to v1.9.0, move golang-lru, remove
  mvn-pom-mutator, collapse the graph to 180 modules, and break project loads.
  No projection was applied.
- Fresh primary vulnerability data remains 1,399 records at SHA-256
  `033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`
  and Last-Modified 2026-09-15T18:54:35Z with no target record; all five exact
  OSV responses are empty. Go 1.26.7 isolated scans are empty. Go 1.18.10
  reports only old-SDK standard-library advisories. Project base and disposable
  v1.3.1 scan populations are identical with zero target package, symbol,
  test-symbol, or reachable trace.
- All twelve guarded modules retain exact selections and incoming edges,
  negative why results, zero repository imports and package loads, runtime
  unreachability, and prior advisory state. The project remains 234 modules,
  3,599 edges, 355 production entries, 429 complete-test entries, 197 module-
  backed entries across 41 modules, 1,067 sum lines, and the 432-line tidy
  projection. Module hashes remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Applicable exact-Go/project/API/CLI/lint/acceptance/meta/audit controls pass,
  including all 17 script meta-tests and 80/80 mutation kills. Go 1.18 retains
  only the two accepted shell wording failures. The known launcher signal/log
  timing race remains target-independent; `--check` and launcher-auto pass.
  Accepted quality remains 27/27 Q0-Q2 PASS at L2. The 435-entry evidence
  manifest SHA-256 is
  `8a30047f3bf22c7ad7e33974aee395c26dc5c10cbd37c9f6dd03daafc1e4b939`.
  Product direction recorded 2026-09-16: the user-authorized automatic
  controller explicitly selected option 1 with the recommended bounds. Retain
  exact selected, inherited, unloaded go-immutable-radix v1.3.1 without
  changing `go.mod` or `go.sum`. Accept only the completed empty-iterator
  reseek panic and the characterized API, nil/panic, key/value aliasing,
  allocation, concurrency, global-state, resource, closure-vulnerability, and
  related findings. The exception is go-immutable-radix-specific and non-
  transferable. It remains valid only while exact v1.3.1 and its three
  selected-version incoming edges from Viper v1.15.0, historical Viper
  v1.10.1, and `sagikazarmark/crypt v0.4.0` remain unchanged, zero target
  packages load, the module remains runtime-unreachable, and no new advisory
  or independent defect appears. Direct import/loading, runtime reachability,
  a version or incoming-edge change, or a new advisory or independent defect
  expires the exception and requires a fresh go-immutable-radix dependency and
  product decision before merge. Do not add a direct edge, select another v1
  release, promote `/v2`, downgrade or otherwise change Viper, change/remove a
  parent, patch/fork source, authorize a wrapper or replacement architecture,
  raise the Go floor, move unrelated selections, alter another guarded
  dependency, manufacture a dependency commit, transfer another exception, or
  request this same choice again while all guards hold.
- Guard-only decision revalidation from clean HEAD
  `7c6d298b71e723a335dc295a1dc42a9b3ead3c93`, parent
  `a3f71d8fa31e787fb85f6a3e2aed9d832d454493`, tree
  `951851771ba2ab15e1d786bbba635a529115b560`, preserves exact
  go-immutable-radix v1.3.1 and its three selected-version incoming edges from
  Viper v1.15.0, historical Viper v1.10.1, and
  `sagikazarmark/crypt v0.4.0`. All 13 guarded modules retain their recorded
  selections and incoming edges, all 13 `go mod why -m` results remain
  negative, repository imports remain zero, and the 355-entry production and
  429-entry complete-test loads contain zero guarded packages. The project
  remains 234 modules, 3,599 graph edges, 197 module-backed complete-test
  entries across 41 loaded modules, and exact `go.mod`/`go.sum` SHA-256 values
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Runtime unreachability therefore remains intact.
- Fresh primary vulnerability data remains byte-identical at 1,399 records,
  SHA-256
  `033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`,
  and Last-Modified 2026-09-15T18:54:35Z. Go-immutable-radix v1.3.1 and every
  guarded target except Gorilla have empty exact-version OSV responses;
  Gorilla retains only the recorded GO-2026-6278/GHSA-w67g-5rqw-f597 result
  and the primary index's existing GO-2020-0019 entry. No new advisory or
  independent defect appeared.
- The bounded exact-path `github.com/hashicorp/go-msgpack` evaluation is
  complete without a qualifying release or repository metadata change. Fresh
  identity evidence resolves Hashicorp's public active unarchived MIT fork of
  `ugorji/go`. Selected v0.5.3 is commit
  `be3a5be7ee2202386d02936a19ae4fbde1c77800`, parents
  `2e9170ac1d8fb32e1e645d8364e4d8f21b530bb3` and
  `3cecd6b54d4f8710db6f370c02e007298827cd34`, tree
  `92376aea28f7d767e816d4f24cc830bba62de0d0`. Highest non-retracted stable
  v0.5.5 is commit `ad60660ecf9c5a1eae0ca32182ed72bab5807961`, tree
  `51aaf54c47643758233939e0ef0f010610605846`. Strict Git verification passes
  and proxy/Git bytes agree. V0.5.0-v0.5.2 expose no root package; v1.1.5 and
  v1.1.6 are retracted after a compatibility-breaking upstream merge, v1.1.6
  requires Go 1.19, and `/v2` is a distinct Go-1.24-or-later module.
- V0.5.3-v0.5.5 have a standard-library-only complete source/test closure,
  preserve Go 1.18, pass production/test cross-compilation under exact Go
  1.26.7 and Go 1.18.10, and expose byte-identical API snapshots. All three
  fail native count-one, two independent count-ten repeats, and race under
  both SDKs because upstream test initialization calls `flag.Parse()` before
  testing registers `-test.paniconexit0`. V0.5.3/v0.5.4 also fail vet on a
  malformed legacy build comment; v0.5.5 fixes the comment but retains the
  early parse. No candidate passes the complete-suite contract.
- The independent behavior fixture SHA-256 is
  `a788002acaf4eb5902c81e0c10d7c577ea1dd1b5ac873d2ce2f45e04969db864`.
  Selected v0.5.3 panics while updating an existing map value and incorrectly
  omits a non-nil pointer to `false`; v0.5.4 fixes the map panic, and v0.5.5
  fixes both. V0.5.5 passes fixture native/repeated/race/vet gates under both
  SDKs but its upstream suite remains blocked. Struct encoding is
  deterministic; map encoding is not. Nil/panic, aliasing, allocation,
  immutable-handle concurrency, protected global-cache, malformed-input,
  unbounded decode-resource, error-identity, and RPC cleanup boundaries are
  characterized.
- Selected v0.5.3 has four exact incoming requests from memberlist
  v0.1.3/v0.3.0 and Serf v0.8.2/v0.9.6. Its why result and repository import
  search are negative; production and complete-test loads contain zero target
  packages, so it is runtime-unreachable. Disposable exact
  v0.5.3/v0.5.4/v0.5.5 gets preserve all 234 selections and every unrelated
  version, manufacture only a direct indirect target root/one graph edge, and
  retain zero load. Tidy removes the root and returns every copy to the base
  projection. No projection was applied.
- Fresh primary data now has 1,402 records at SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
  Last-Modified 2026-09-17T17:29:18Z, with no go-msgpack record. Every serious
  exact-version OSV response is empty. Go 1.26.7 isolated scans are empty; Go
  1.18.10 reports only old-SDK standard-library advisories. Base and v0.5.5
  project scan populations are identical with zero target trace. All 13 prior
  guarded selections remain intact. The 521-entry evaluation manifest
  SHA-256 is
  `100fa819df480ac176571f5d4fe672afd582554c9feb1fa8d9acc0fb50dc6a76`.
- Applicable exact-Go project verification, native/repeated/race/vet, pinned
  lint, API/CLI, empty-HOME, four cross-build, host/snapshot/Docker acceptance,
  full preflight, all 17 script meta-tests, 80/80 mutation kills, and all 15
  audit controls pass. No changed-selection scorecard applies, so accepted
  quality remains 27/27 Q0-Q2 PASS at L2.
- Product direction recorded 2026-09-19: option 1 was explicitly selected with
  the recommended bounds. Retain exact selected, inherited, unloaded
  go-msgpack v0.5.3 without changing `go.mod` or `go.sum`. Accept only the
  completed map-update panic, pointer-`false` omission, malformed-build-
  comment/vet failure, upstream early-test-flag failure, and the characterized
  API, nil/panic, byte aliasing, allocation, concurrency, global-cache,
  resource, deterministic, malformed-input, error-identity, RPC-cleanup, and
  vulnerability findings. The exception is go-msgpack-specific and non-
  transferable and accepts no new or independently discovered defect.
- The exception remains valid only while exact v0.5.3 and all four requests
  from memberlist v0.1.3/v0.3.0 and Serf v0.8.2/v0.9.6 remain unchanged, zero
  target packages load, the module remains runtime-unreachable, and no new
  advisory or independent defect appears. Direct import/loading, runtime
  reachability, a target version or incoming-request change, or a new advisory
  or independent defect expires the exception and requires a fresh go-msgpack
  dependency and product decision before merge. Do not add a direct edge,
  select v0.5.4/v0.5.5, promote a retracted v1 release or `/v2`, change/remove
  a parent, patch/fork source, authorize a wrapper or replacement architecture,
  raise the Go floor, move unrelated selections, alter another guarded
  dependency, manufacture a dependency commit, transfer another exception, or
  request this same choice again while all guards hold.
- Guard-only decision revalidation from clean HEAD
  `69017090043510b4754e8e81feab762a8011fe90`, parent
  `13f63a2c7639e373b261ca74176955b2be6461b2`, tree
  `0db64c05ce2166f40821aeb296e10576124ed349`, preserves exact go-msgpack
  v0.5.3 and all four recorded requests. All 14 guarded selections and incoming
  edges remain unchanged; all 14 `go mod why -m` results remain negative;
  repository imports are zero; and the 355-entry production and 429-entry
  complete-test loads contain zero guarded packages. The project remains 234
  modules, 3,599 graph edges, 197 module-backed complete-test entries across
  41 loaded modules, and exact `go.mod`/`go.sum` SHA-256 values
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Runtime unreachability remains intact.
- Fresh primary vulnerability data remains byte-identical at 1,402 records,
  SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
  and Last-Modified 2026-09-17T17:29:18Z. Exact go-msgpack v0.5.3 and every
  guarded target except Gorilla have empty exact-version OSV responses.
  Gorilla retains only the recorded GO-2026-6278/GHSA-w67g-5rqw-f597 result
  and existing GO-2020-0019 primary-index entry. No new advisory or independent
  defect appeared.
- Exact Go 1.26.7 module verification, build, count-one tests, race tests, vet,
  and launcher lifecycle check pass under the established `umask 022` mode.
  No dependency implementation or metadata commit was created; no changed-
  selection scorecard applies, and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. P7 continues only with the prepared bounded exact-path go-multierror
  evaluation; it was not executed in this decision session.

Hashicorp go-multierror evaluation and blocked product boundary (2026-09-19):

- The exact module resolves to Hashicorp's public, active, unarchived,
  non-fork MPL-2.0 repository. The only exact-path stable releases are v1.0.0,
  selected v1.1.0, and latest v1.1.1. There are no release objects,
  retractions, deprecations, redirects, alternate module lines, or qualified
  newer tags. Current main is unreleased and was not promoted. Proxy and
  tagged-Git source bytes agree for every release.
- All three candidates have one package and a complete minimal module closure
  of the target plus exact `github.com/hashicorp/errwrap v1.0.0`. Imported
  production and test source preserves Go 1.18. Each release passes source
  verification, native count-one, two independent count-ten repeats, race,
  vet, and five production/test cross-build targets under exact Go 1.26.7 and
  contained Go 1.18.10.
- V1.0.0 lacks `Group` and `(*Error).Unwrap`. V1.1.0 and v1.1.1 have no API
  diff; v1.1.1 fixes only the v1.1.0 typed-nil `WrappedErrors` panic among the
  independently characterized boundaries. Formatting/order, mutation,
  aliasing, allocation, nil/panic behavior, ordinary error traversal,
  concurrency, Group lifecycle, global state, and resource ownership are
  recorded in the answered archive.
- No stable release qualifies. `Prefix` delegates to Errwrap v1.0.0 and loses
  standard error identity: under both SDKs v1.1.0 and v1.1.1 fail all three
  independent assertions for `errors.Is` traversal through a prefixed plain
  error and multierror member and `errors.As` traversal through a prefixed
  member. V1.0.0 has the same design. Upstream issue 56 and unmerged PR 155
  corroborate a `%w` repair, but no release contains it.
- Selected inherited v1.1.0 has six requests: Serf v0.9.6 requests v1.1.0;
  memberlist v0.1.3/v0.3.0, mitchellh/cli v1.0.0/v1.1.0, and
  posener/complete v1.2.3 request v1.0.0. Its why result and repository import
  search are negative, production and complete-test loads contain zero target
  packages, and it is runtime-unreachable. All 14 earlier guards remain intact.
- Disposable exact v1.1.0/v1.1.1 gets preserve all 234 selections and zero
  target load but manufacture main-module Multierror and Errwrap roots.
  V1.1.1 adds a new target-to-Errwrap graph edge, expiring the separate Errwrap
  guard. Exact v1.0.0 removes mvn-pom-mutator and makes the project unloadable.
  Tidy returns the viable projections to the base selection; none was applied.
- Fresh 1,402-record primary vulnerability data has no target record. Exact
  OSV results for all candidates and Go 1.26.7 isolated scans are empty. Go
  1.18.10 reports only its old standard-library population. Base and v1.1.1
  project scan outputs are byte-identical with zero target advisory or trace.
- Every applicable no-change project gate passes, including full preflight,
  all 17 script/meta stages, 80/80 mutation kills, host/snapshot/Docker
  acceptance, and all 15 audit controls. The authoritative clean-tree result
  is 27/27 Q0-Q2 PASS at L2 with a 469-entry selected-evidence manifest at
  SHA-256
  `7db323255b0c8361371237d1e3e7405e697e3ad7f53167678f756535dc66b7f3`.
  No source or dependency metadata changed.
- P7 is blocked on the prepared bounded go-multierror product decision. The
  recommended option is to retain exact selected, inherited, unloaded v1.1.0
  under a target-specific exception and exact selection/request/load/runtime/
  advisory guards. An exact v1.1.1 choice would require a later dependency
  implementation and a fresh, separately sequenced Errwrap decision. No
  choice is inferred by this evaluation.

Hashicorp go-multierror product decision (2026-09-19):

- Authorized option 1 was selected exactly: retain exact selected, inherited,
  unloaded `github.com/hashicorp/go-multierror v1.1.0` without changing
  `go.mod` or `go.sum`. The go-multierror-specific, non-transferable exception
  accepts only the completed `Prefix` `errors.Is`/`errors.As` identity loss,
  typed-nil `WrappedErrors` panic, and completed API, formatting/order,
  nil/panic, mutation, aliasing, allocation, concurrency, Group lifecycle,
  resource/global-state, MVS, vulnerability, and related qualification
  findings. It accepts no new or independently discovered defect.
- Guard-only revalidation from clean HEAD
  `a6dba3eee16d9c3c622fdc68c075aa45595c27cc`, parent
  `ea7bc704e1160eca51b3b35d00411db28babcac8`, tree
  `99c8ac8c6c063ec11cc66232af33f69d48c7e054`, passed under exact Go 1.26.7
  binary SHA-256
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
  The prior handoff changed exactly the launcher, answered evaluation archive,
  new decision archive, rolling handover, and roadmap; the ordinary and
  ignored worktree status was clean and the reciprocal archive chain passed.
- Go-multierror remains exact v1.1.0 through all six recorded requests: Serf
  v0.9.6 requests v1.1.0; memberlist v0.1.3/v0.3.0, mitchellh/cli
  v1.0.0/v1.1.0, and posener/complete v1.2.3 request v1.0.0. All 15 guarded
  selections and recorded incoming edges remain unchanged, all 15 why results
  remain negative, repository imports remain zero, and the 355-entry
  production and 429-entry complete-test loads contain zero guarded packages.
  The complete-test load retains 197 module-backed entries across 41 loaded
  modules, so every guarded module remains runtime-unreachable.
- The separate Errwrap decision remains intact: exact Errwrap v1.0.0 retains
  its sole selected-version incoming edge from go-multierror v1.1.0; the
  historical go-multierror v1.0.0 request remains; and no main-module Errwrap
  root or new v1.1.1-to-Errwrap graph edge exists.
- The project remains 234 modules, 3,599 graph edges, 1,067 sum lines, and the
  recorded 432-line unapplied tidy projection. `go.mod`/`go.sum` SHA-256
  values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- Fresh primary vulnerability data remains byte-identical at 1,402 records,
  SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
  and Last-Modified 2026-09-17T17:29:18Z. Exact go-multierror v1.1.0 and every
  guarded target except Gorilla retain empty exact-version OSV responses;
  Gorilla retains only its recorded results. No new advisory or independently
  observed defect appeared.
- The exception remains valid only while exact v1.1.0 and all six incoming
  requests remain unchanged, target imports and production/complete-test loads
  remain zero, runtime unreachability remains intact, the separate Errwrap
  v1.0.0 selection and selected-version edge remain unchanged without a new
  root or v1.1.1 edge, and no new advisory or independent defect appears.
  Direct import/loading, runtime reachability, a target version or request
  change, an Errwrap-guard change, or a new advisory or independent defect
  expires the exception and requires the owning fresh dependency and product
  decision before merge.
- No dependency implementation or metadata commit was created. Exact Go
  1.26.7 module verification, build, count-one tests, race tests, vet, and the
  launcher lifecycle check pass; no changed-selection scorecard applies and
  accepted quality remains 27/27 Q0-Q2 PASS at L2. P7 continues only with the
  prepared bounded exact-path go-retryablehttp v0.5.3 evaluation; it was not
  executed in this decision session.

Hashicorp go-retryablehttp evaluation and blocked product boundary
(2026-09-19):

- Fresh exact-path identity evidence resolves Hashicorp's public, active,
  unarchived, non-fork MPL-2.0 repository. The proxy exposes 23 stable releases
  from v0.5.0 through latest v0.7.8, with no prerelease, retraction,
  deprecation, redirect, alternate path, `/v2` line, or GitHub Release object.
  Current main is 17 commits beyond v0.7.8 and remains unreleased. Proxy and
  Git manifests agree for every release; relevant tags are unsigned.
- Selected v0.5.3 is verified commit
  `357460732517ec3b57c05c51443296bdd6df1874`. Last declared-floor-compatible
  v0.7.5 is `4165cf8897205a879a06b20d1ed0a2a76fbb6a17`. First secure v0.7.7
  contains fix commit `a99f07beb3c5faaa0a283617e6eb6bcf25f5049a`; latest v0.7.8 is
  `e1f5485fe84728709b857cb89e17088894c301d6`.
- V0.5.0-v0.6.2 have no `go` directive, v0.6.3-v0.7.5 declare Go 1.13,
  v0.7.6/v0.7.7 declare Go 1.19, and v0.7.8 declares Go 1.23. V0.7.5 is the
  highest release preserving the declared Go 1.18 floor. Every eligible
  release is affected by reviewed GO-2024-2947/GHSA-v6v8-xj6m-xwqh; first-
  fixed v0.7.7 is floor-ineligible. Actual Go 1.18 compilation of later
  releases does not qualify their unsupported declared floor.
- Every release after v0.5.3 is exported-API-incompatible because
  `Client.Logger` changes from `Logger` to `interface{}`. Floor-compatible
  releases pass contained Go 1.18.10 isolated tests. Exact Go 1.26.7
  v0.6.3-v0.7.5 tests fail modern wrapped-x509 retry classification;
  v0.5.3/v0.6.2 retain a test-only vet failure; v0.7.7 passes complete exact-
  Go native/repeated/race/vet and five cross-build targets.
- The behavior fixture SHA-256 is
  `cf4ceef2094699c773cc5a44afb90ed992d4ac3ae222a3ee27ff79d95bea1a1c`.
  Selected v0.5.3 fails x509 unknown-authority, unsupported-scheme, and
  redirect-limit permanent classification and loses final transport-error
  `errors.Is` identity. Body replay, response closing, cancellation,
  status/backoff, hooks, nil/panic, mutation, aliasing, allocation,
  concurrency, global client, and resource boundaries are characterized.
- Fresh primary data retains 1,402 records at SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
  and Last-Modified 2026-09-17T17:29:18Z. Exact OSV results contain both
  target IDs through v0.7.6 and are empty for v0.7.7/v0.7.8. The independent
  credential-disclosure fixture SHA-256 is
  `2b15b5812d21c3f91f9279105923423f6fbb311568c56ec8d3651e91fd952bcb`;
  consumer module/package/symbol/test-symbol scans corroborate the affected
  and fixed ranges.
- Selected v0.5.3 exists solely through go-metrics v0.3.10. Its why/import
  results are negative and production/complete-test loads contain zero target
  packages, so it is runtime-unreachable. V0.5.3/v0.6.2/v0.7.5 projections
  move no unrelated selection. V0.7.6/v0.7.7 also upgrade guarded go-hclog
  v1.2.0 to v1.6.3 and fatih/color v1.15.0 to v1.16.0; v0.7.8 also raises
  the main Go directive to 1.23. No projection was applied.
- All 15 earlier guards remain exact and runtime-unreachable. The project
  remains 234 modules, 3,599 graph edges, 355/429 package entries, 197 module-
  backed complete-test entries across 41 modules, and 1,067 sum lines.
  `go.mod`/`go.sum` SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- Exact Go project, compatibility, lint, cross-build, host/snapshot/Docker,
  all 17 script/meta, all 80 mutation, and all 15 audit controls pass. The
  authoritative exact 21-stage result is 27/27 Q0-Q2 PASS at L2 with zero
  held/regressed/non-comparable/dirty counts; scorecard SHA-256 is
  `de13154181ae319fd80a29a98df78535762d70c2c01438f9e3cb735c2f624e81`.
  The 547-entry complete-evidence manifest SHA-256 is
  `12f36165ac7c0ddb4a69a5e432c3b12ea57c977a9e6d7c7a338a359d36b4faa2`.
- No stable release qualifies. P7 is blocked on the prepared bounded product
  decision. Option 1 is recommended: retain exact inherited, unloaded v0.5.3
  under a new target-specific exception accepting only the completed advisory,
  credential-disclosure, behavior, API, MVS, and related findings while exact
  selection and the sole go-metrics request, zero imports/load, runtime
  unreachability, and no-new-advisory/defect guards hold. V0.7.7 instead needs
  separate Go-floor, exported-API, go-hclog, and color decisions; parent-chain
  remediation is a third separately scoped option. No decision is inferred.

Hashicorp go-retryablehttp product decision (2026-09-19):

- Authorized option 1 was selected exactly: retain exact selected, inherited,
  unloaded `github.com/hashicorp/go-retryablehttp v0.5.3` without changing
  `go.mod` or `go.sum`. The target-specific, non-transferable exception accepts
  only GO-2024-2947/GHSA-v6v8-xj6m-xwqh and the independently reproduced
  Basic-auth URL credential disclosure; the x509 unknown-authority,
  unsupported-scheme, and redirect-limit permanent-retry classification
  defects; final transport-error `errors.Is` identity loss; the test-only vet
  defect; and the completed API, body, cancellation, status, backoff, hook,
  nil/panic, mutation, aliasing, allocation, concurrency, global-state,
  resource, MVS, vulnerability, and related findings. It accepts no new or
  independently discovered defect.
- Guard-only revalidation began from clean handoff HEAD
  `a4c24f29e25eb21cc229cde674ffdacdecc78f38`, parent
  `24b1f33671362309526328170e67e68aaf753a3a`, tree
  `ccca83e0256147f7e51d14eb07a2fb2472e2aabe`, under exact Go 1.26.7 binary
  SHA-256
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
  The handoff changed exactly the launcher, answered evaluation archive, new
  decision archive, rolling handover, and roadmap; reciprocal archive history,
  Google UUID implementation ancestry, and the launcher lifecycle check pass.
- Go-retryablehttp remains exact v0.5.3 through the sole exact go-metrics
  v0.3.10 request. All 15 earlier selections/requests plus the target guard
  remain unchanged; all 16 why results are negative; repository imports are
  zero; and the 355-entry production and 429-entry complete-test loads contain
  zero guarded packages. The complete-test load retains 197 module-backed
  entries across 41 modules, so every guarded module remains runtime-
  unreachable.
- The project remains 234 modules, 3,599 graph edges, 1,067 sum lines, and the
  recorded 432-line unapplied tidy projection. `go.mod`/`go.sum` SHA-256 values
  remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  No parent, Go-floor, API, source, or unrelated selection changed.
- Fresh primary vulnerability data remains byte-identical at 1,402 records,
  SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
  and Last-Modified 2026-09-17T17:29:18Z. Exact-version OSV remains empty for
  14 earlier guards; Gorilla retains only its recorded finding; and exact
  target v0.5.3 retains exactly its two accepted active identifiers. No new
  advisory or independently observed defect appeared.
- The exception remains valid only while exact v0.5.3 and the sole go-metrics
  v0.3.10 request remain unchanged, no direct root or repository import is
  added, production and complete-test target loads remain zero, runtime
  unreachability and every earlier owning guard remain intact, and no new
  target advisory or independent defect appears. Direct import/loading,
  runtime reachability, a selection/request change, a new direct root, an
  earlier guard change, or a new advisory or independent defect expires the
  exception and requires the owning fresh dependency and product decision
  before merge.
- No dependency implementation or metadata commit was created. Exact Go
  1.26.7 module verification, build, count-one tests, race tests, vet, and the
  launcher lifecycle check pass; no changed-selection scorecard applies and
  accepted quality remains 27/27 Q0-Q2 PASS at L2. P7 continues only with the
  prepared bounded exact-path `github.com/hashicorp/go-rootcerts v1.0.2`
  evaluation; it was not executed in this decision session.

Hashicorp go-rootcerts evaluation and blocked product boundary (2026-09-20):

- Fresh exact-path evidence resolves Hashicorp's public, active, unarchived,
  non-fork MPL-2.0 repository. The proxy exposes exactly stable v1.0.0,
  v1.0.1, and latest v1.0.2, with no prerelease, retraction, module
  deprecation, redirect, alternate exact path, `/v2` line, or GitHub Release
  object. Current master is 34 commits beyond v1.0.2; its unreleased
  pseudo-version declares Go 1.23 and is not promoted.
- V1.0.0 is unsigned commit
  `63503fb4e1eca22f9ae0f90b49c5d5538a0e87eb`; v1.0.1 is unsigned commit
  `df8e78a645e18d56ed7bb9ae10ffb8174ab892e2`; v1.0.2 is unsigned annotated
  tag `dafef5a8f4d570697f148431a2004909428ca75b` over GitHub-verified commit
  `98fadc2a5ba2ad2a534a179b352ecdfd1f4259aa`. Proxy-ZIP SHA-256 values are
  `4393b0b9cd741e00de5624d5124cf054bf50c57231d4b1caff84c8a4d16c6a47`,
  `3f558b1a436ed6fb15872383545109227f9552bf5daa95583e9402bbd3a24fff`, and
  `864a48e642e87a273fb5ef60bb3575bd74a7090510f93143163fa6700be31948`.
  Sumdb identities and every regular proxy/Git source byte agree. Standard Go
  ZIP rules omit two tracked symlink fixtures.
- V1.0.0 has no Go directive; v1.0.1/v1.0.2 declare Go 1.12. The complete
  minimal closure adds only go-homedir v1.0.0 or v1.1.0, whose modules omit a
  Go directive, so every stable candidate preserves Go 1.18. Tagged Git
  sources pass verification, native/repeated/race/vet, and five cross-builds
  under exact Go 1.26.7 and contained Go 1.18.10. Proxy-archive tests for all
  three tags fail only the symlink-fixture test because the fixture directory
  is absent.
- Pinned API comparison finds no v1.0.0-to-v1.0.1 change. V1.0.2 adds
  `AppendCertificate` and `Config.CACertificate`; its `[]byte` field makes
  `Config` non-comparable. There is no later stable candidate.
- Independent fixtures reproduce two disqualifying behaviors in every stable
  release under both SDKs. File/path/ConfigureTLS errors stringify underlying
  I/O causes with `%s`, losing `errors.Is` identity. Darwin system-root loading
  can return a non-nil empty pool and nil error when successful keychain
  commands emit no PEM because the false `AppendCertsFromPEM` result is
  ignored. Upstream tests miss the empty-pool condition, and current master
  retains both defects. Valid certificate construction, precedence, malformed
  input, nil, mutation/aliasing, allocation, concurrency, resource, and global-
  state behavior is otherwise characterized. Fixture SHA-256 values are
  `512a4a1df5b96659c975ca6863c8172832a37a4fe9038da015165b28cad1acd1` and
  `4f0c62af97695e5015db568558cf401e85d03d41efca8954e9d6c3f727d97768`.
- Fresh 1,402-record primary vulnerability data remains at SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
  and Last-Modified 2026-09-17T17:29:18Z. Exact OSV, GitHub advisory, and
  govulncheck module/package/symbol/test-symbol scans are empty for all three
  tags. Base and redundant-v1.0.2 project scan finding sets are identical with
  no target frame or reachable trace.
- Selected v1.0.2 exists through requests from Viper v1.15.0, historical
  Viper v1.10.1, `sagikazarmark/crypt v0.4.0`, and historical Consul API
  v1.12.0; historical Consul API v1.1.0 requests v1.0.0. The target why and
  repository import results are negative; both package loads contain zero
  target packages. A redundant exact v1.0.2 root changes only one edge and one
  checksum line before tidy converges to the base projection. V1.0.1/v1.0.0
  instead force Viper v1.8.1, remove mvn-pom-mutator, disrupt guards, and make
  project packages unloadable. No projection was applied.
- The project remains 234 modules, 3,599 graph edges, 355 production and 429
  complete-test entries, 197 module-backed complete-test entries across 41
  modules, 1,067 sum lines, and the recorded 432-line tidy projection. All 16
  earlier guards plus the target retain negative why/import/load results and
  their exact advisory states. `go.mod`/`go.sum` SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- Exact-Go project native/repeated/race/vet, lint, API/CLI, empty-HOME, cross-
  build, full preflight, all 17 script/meta, all eight mutation meta-stages,
  80/80 mutation kills, host/snapshot/Docker acceptance, and all 15 audit
  controls pass. Contained Go 1.18 retains only the two accepted Darwin shell
  wording failures; 26 unaffected packages and applicable gates pass. No
  changed-selection scorecard applies, so accepted quality remains 27/27 Q0-
  Q2 PASS at L2 with zero held/regressed/non-comparable/dirty counts.
- No stable release qualifies. P7 is blocked on the prepared go-rootcerts
  product decision. Recommended option 1 retains exact inherited, unloaded
  v1.0.2 under a new target-specific exception bounded to the completed
  behavior/archive/API findings, exact five requests, zero imports/load,
  runtime unreachability, every earlier guard, and no new advisory or defect.
  Patch/fork/replacement and parent-chain pruning are separately scoped
  alternatives. No choice is inferred by this evaluation.

Hashicorp go-rootcerts product decision (2026-09-20):

- Authorized option 1 was selected exactly: retain exact selected, inherited,
  unloaded `github.com/hashicorp/go-rootcerts v1.0.2` without changing
  `go.mod` or `go.sum`. The new target-specific, non-transferable exception
  accepts only the independently reproduced loss of underlying filesystem-
  error identity for file, path, and ConfigureTLS failures; the Darwin silent-
  empty system-root pool when successful keychain commands emit no PEM; the
  proxy-archive symlink-fixture test failure; the v1.0.2 `Config`
  comparability change; and the completed certificate-pool, PEM, file/path,
  precedence, system-root, environment, nil/panic, mutation, aliasing,
  allocation, concurrency, global-state, resource, API, MVS, vulnerability,
  and related qualification findings. No new or independently discovered
  defect is accepted.
- Guard-only revalidation from clean handoff HEAD
  `c530f7d8e39aa0181899d7cb2893f6a16555ebef`, parent
  `888412a62d768cb11af6d9e132600cc4a34f0679`, tree
  `e0eeec6c821ba393d8ab30751cf411296b8ec0f6`, preserves exact v1.0.2
  and all five recorded requests: Viper v1.15.0, historical Viper v1.10.1,
  `sagikazarmark/crypt v0.4.0`, and historical Consul API v1.12.0 request
  v1.0.2; historical Consul API v1.1.0 requests v1.0.0. All 17 target-plus-
  earlier why results remain negative, repository imports are zero, and the
  355-entry production and 429-entry complete-test loads contain zero guarded
  packages. The complete-test load retains 197 module-backed entries across 41
  modules, so every guarded target remains runtime-unreachable.
- The project remains 234 modules, 3,599 graph edges, 1,067 sum lines, and the
  recorded 432-line unapplied tidy projection. `go.mod`/`go.sum` SHA-256
  values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  No parent, Go-floor, exported API, source, or unrelated selection changed.
- Fresh primary vulnerability data remains byte-identical at 1,402 records,
  SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
  and Last-Modified 2026-09-17T17:29:18Z. Exact-version OSV remains empty for
  go-rootcerts and every guarded target except Gorilla and go-retryablehttp,
  which retain only their recorded accepted identifiers. The target GitHub
  advisory feed remains empty. No new advisory or independent defect appeared.
- The exception remains valid only while exact go-rootcerts v1.0.2 and all five
  incoming requests and versions remain unchanged, no direct main-module edge
  or repository import is added, production and complete-test target loads
  remain zero, runtime unreachability holds, every earlier selection/request/
  load/runtime/advisory guard remains intact, and no new target advisory or
  independent defect appears. Direct import/loading, runtime reachability, a
  version or request change, a new direct root, an earlier owning-guard change,
  or a new advisory or independent defect expires the exception and requires
  the owning fresh dependency and product decision before merge.
- No dependency implementation or metadata commit was created. Exact Go
  1.26.7 module verification, build, count-one tests, race tests, vet, and the
  62-control launcher lifecycle suite pass; no changed-selection scorecard
  applies and accepted quality remains 27/27 Q0-Q2 PASS at L2. P7 continues
  only with the prepared bounded exact-path
  `github.com/hashicorp/go-sockaddr v1.0.0` evaluation; it was not executed in
  this decision session.

Hashicorp go-sockaddr evaluation (2026-09-20):

- No exact-path stable `github.com/hashicorp/go-sockaddr` release qualifies.
  Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves the
  public active unarchived non-fork MPL-2.0 Hashicorp repository and only
  stable v1.0.0-v1.0.7. There is no prerelease, retraction, module
  deprecation, redirect, alternate exact path, or `/v2` line. V1.0.7 is
  latest and the sole GitHub Release; every version tag is lightweight.
  Current master is 63 commits beyond v1.0.7, declares Go 1.25, and is not a
  candidate. Proxy and Git regular-file bytes agree.
- V1.0.0-v1.0.3 omit a `go` directive and are the only serious floor-eligible
  releases. V1.0.4 declares Go 1.19.0 and v1.0.5-v1.0.7 declare Go 1.19.
  V1.0.1-v1.0.3 have a complete 13-module graph and 11 actually imported
  external modules whose source preserves Go 1.18. Selected v1.0.0 omits four
  command dependencies from its released `go.mod`, so read-only listing,
  native tests, race, vet, command build, whole-module API export, and all
  seven cross-target loads fail under both SDKs.
- V1.0.1-v1.0.3 load and cross-build under exact Go 1.26.7 and contained Go
  1.18.10, but native repeats/race retain upstream Darwin interface, route,
  and ordering failures; Go 1.18 cannot parse the exported
  `HashiCorpDefault2016` template; and vet rejects every release's unreachable
  panic. V1.0.3 and later additionally invert Windows PowerShell detection.
- Independent fixtures under both SDKs reproduce four defects in v1.0.0-
  v1.0.3 and v1.0.7: IPv6 `ContainsAddress` rejects an interior host; no-match
  interface helpers return empty output with nil error; all six exported
  attribute inventories expose mutable global slices and race; and `Host()`
  exposes the global IPv6 host-mask backing array, allowing caller corruption
  of later addresses. Valid address, Unix, RFC, JSON, CLI, enumeration,
  copy/non-aliasing, allocation, error, nil/panic, global-state, environment,
  command/resource, and concurrency boundaries are otherwise characterized.
- Root exported APIs are identical through v1.0.3; v1.0.7 only adds compatible
  error sentinels. V1.0.1-v1.0.3 command smokes pass under both SDKs, while
  v1.0.0 cannot build the command read-only.
- Selected v1.0.0 remains inherited through exactly memberlist v0.1.3 and
  v0.3.0. Target why and import results are negative and both project loads
  contain zero target packages. A redundant v1.0.0 root adds only one edge
  and checksum line. Exact v1.0.1-v1.0.3 projections also add
  `go-wordwrap v1.0.0` and upgrade unrelated `columnize` to
  v2.1.0+incompatible; v1.0.7 moves 19 module lines including guarded Errwrap.
  No projection was applied.
- Fresh 1,402-record primary vulnerability data retains SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
  with no target record. Exact OSV, GitHub advisory, and isolated
  module/package/symbol/test-symbol scans are empty. Base and disposable
  v1.0.3 project results are identical with no target trace. All 18 target-
  plus-earlier guards remain negative/unloaded/runtime-unreachable with their
  recorded advisory states.
- Exact Go 1.26.7 project verification, build, repeats, race, vet, pinned
  lint, API/CLI, empty-HOME, cross-builds, canonical preflight, all script/meta
  populations, 80/80 live mutation kills, host/snapshot/Docker acceptance,
  and all 15 audit meta-controls pass. Contained Go 1.18 has only the two
  accepted Darwin shell wording failures; all 26 unaffected packages and
  applicable gates pass. No changed-selection scorecard applies and accepted
  quality remains 27/27 Q0-Q2 PASS at L2.
- Dependency metadata remained unchanged. The evaluation stopped on the sole
  prepared go-sockaddr decision without inferring acceptance from physical
  selection or zero reachability.

Hashicorp go-sockaddr product decision (2026-09-20):

- Authorized option 1 was selected exactly: retain exact selected, inherited,
  unloaded `github.com/hashicorp/go-sockaddr v1.0.0` without changing
  `go.mod` or `go.sum`. The new target-specific, non-transferable exception
  accepts only the incomplete released module closure and command
  unbuildability; the reproduced IPv6 containment, no-match nil-error,
  mutable global-slice/race, and exposed IPv6 mask defects; the characterized
  upstream Darwin, Go 1.18 template, vet, sorting, route-command, API/CLI,
  nil/panic, mutation, aliasing, allocation, concurrency, global-state,
  resource, environment, MVS, vulnerability, and related findings. No new or
  independently discovered defect is accepted.
- Guard-only revalidation from clean handoff HEAD
  `5c957ce042e2c36589e0ec5f00fca5261ce76352`, parent
  `a9289a1684e101399810bbec6294b62e4a41df92`, tree
  `654abb0f7bf3261c0b2224f99d61ec968af826f2`, preserves exact v1.0.0 and
  its two requests from memberlist v0.1.3/v0.3.0 with no direct root. All 18
  target-plus-earlier why results remain negative, repository imports are
  zero, and the 355-entry production and 429-entry complete-test loads contain
  zero guarded packages. The complete-test load retains 197 module-backed
  entries across 41 modules, so every guarded target remains runtime-
  unreachable.
- The project remains 234 modules, 3,599 graph edges, 1,067 sum lines, and the
  recorded 432-line unapplied tidy projection. `go.mod`/`go.sum` SHA-256
  values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  No parent, Go-floor, API, source, direct root, or unrelated selection
  changed.
- Fresh primary vulnerability data remains byte-identical at 1,402 records,
  SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
  and Last-Modified 2026-09-17T17:29:18Z. Exact-version OSV remains empty for
  go-sockaddr and every guarded target except Gorilla and go-retryablehttp,
  which retain only their recorded accepted identifiers. The target GitHub
  repository and exact global-advisory results remain empty. No new advisory
  or independent defect appeared.
- The exception remains valid only while exact v1.0.0 and both memberlist
  requests remain unchanged, no direct root or repository import is added,
  production and complete-test target loads remain zero, runtime
  unreachability holds, every earlier selection/request/load/runtime/advisory
  guard remains intact, and no new target advisory or independent defect
  appears. Direct import/loading, runtime reachability, a version/request
  change, a new direct root, an earlier owning-guard change, or a new advisory
  or independent defect expires the exception and requires the owning fresh
  dependency and product decision before merge.
- No dependency implementation or metadata commit was created. Exact Go
  1.26.7 module verification and applicable no-change project and launcher
  gates pass; no changed-selection scorecard applies and accepted quality
  remains 27/27 Q0-Q2 PASS at L2. P7 continues only with the prepared bounded
  exact-path `github.com/hashicorp/go-syslog v1.0.0` evaluation; it was not
  executed in this decision session.

Hashicorp go-syslog evaluation and product decision (2026-09-20):

- No exact-path stable `github.com/hashicorp/go-syslog` release qualifies.
  Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves the public
  active unarchived non-fork MIT Hashicorp repository. Exact v1.0.0 is the only
  proxy version and `@latest`; there is no other stable or prerelease version,
  GitHub Release, retraction, module deprecation, redirect, alternate exact
  path, or `/v2` line. Six branches and returned forks are not candidates.
- V1.0.0 is a lightweight tag at verified signed commit
  `8d1874e3e8d1862b74e0536851e218c4571066a5`, parent
  `326bf4a7f709d263f964a6a96558676b103f3534`, tree
  `a81f242fd29ca869a5c2fb666ff9e4dd7a0af570`. Proxy ZIP and `go.mod`
  SHA-256 values are
  `a0ca8b61ea365e9ecdca513b94f200aef3ff68b4c95d9dabc88ca25fcb33bce6` /
  `628ea5197c3c1c9f1e04c7354d9323d12ccb79b1f8a851f817e3bd04f6828aa6`;
  proxy and Git regular files match and sumdb verifies the release. Current
  master `40240a543e78c0ccbdb437b470650ea7acaf3da8` is 48 commits beyond the
  tag, declares Go 1.23.0, and is unreleased and floor-ineligible. Its exported
  API is unchanged and it fixes only one of the three disqualifiers.
- V1.0.0 has no `go` directive or external dependency; its complete minimal
  production/test closure is target plus standard library and preserves Go
  1.18. The release has one package, four production files, no native tests,
  commands, examples, benchmarks, fuzz targets, testdata, generated files,
  cgo, or embeds. Native verification, listing, repeats, race, vet, and
  applicable cross-builds pass under exact Go 1.26.7 and Go 1.18.10. The
  exported API comprises `Priority`, eight severity constants, `Syslogger`,
  `NewLogger`, and `DialLogger`; there is no CLI and the v1.0.0-to-master API
  diff is empty.
- Independent fixture SHA-256
  `6d9eb0b8fbfefc1556630519301609b6d343cb37ec2b3b5c46d6d6c2378d46d0`
  passes behavior repeats, race, and vet under both SDKs. It characterizes
  UDP/TCP formatting and priority, local Unix-daemon behavior, framing,
  dial/write/close ownership, error identity, malformed input, nil/panic,
  mutation/non-aliasing, allocation, concurrency, global state, resources,
  and environment/network interaction. Embedded newlines pass through
  verbatim and may create multiple newline-framed records.
- Three independent defects disqualify the only release. External fixture
  SHA-256
  `b50b035d93a2244fbe575ee2d0dbb34ec1d70b9ca501bde3817f517781f8b58a`
  cannot compile because `NewLogger` and `DialLogger` are absent on AIX and
  `js/wasm` under both SDKs and on `wasip1/wasm` under Go 1.26.7, though the
  package itself nominally builds there. `writeString` ignores a
  `SetWriteDeadline` error, writes anyway, returns nil, loses error identity,
  and defeats the blocking-avoidance guarantee. `DialLogger` accepts invalid
  severity `Priority(8)` and silently maps it to `LOG_EMERG`. Both runtime
  blockers reproduce under both SDKs; master fixes only the deadline error.
- Selected v1.0.0 exists through exactly Serf v0.8.2 and v0.9.6. Its why and
  import results are negative and both loads contain zero target packages. A
  redundant root adds only one edge and checksum line, changes no selection or
  load, and is removed by tidy; base and candidate projections converge byte-
  identically. All 18 earlier guarded selections, requests, negative why/
  imports, zero loads, runtime unreachability, and advisory states remain
  exact. No projection was applied.
- Fresh primary vulnerability data remains 1,402 records at SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
  with no target entry. Exact OSV and GitHub advisory queries are empty.
  Exact-Go isolated scans have zero findings; old-SDK findings belong only to
  its standard library. Base and redundant-root project scans are identical at
  30 module, 22 package, and 20 symbol/test-symbol IDs, with zero target frame
  or trace.
- The project remains 234 modules, 3,599 edges, 355 production entries, 429
  complete-test entries, 197 module-backed entries across 41 loaded modules,
  1,067 sum lines, and the recorded 432-line tidy projection. `go.mod` and
  `go.sum` SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact-Go native/repeat/race/vet, lint, API/CLI, empty-HOME, cross-build,
  canonical preflight, every script/meta population, 80/80 mutation kills,
  host/snapshot/Docker acceptance, and all 15 audit controls pass. No changed-
  selection scorecard applies; accepted quality remains 27/27 Q0-Q2 L2.
- Dependency metadata remained unchanged. Product direction supplied
  2026-09-20: the user-authorized automatic controller explicitly selected
  option 1 with the recommended bounds. Retain exact selected, inherited,
  unloaded go-syslog v1.0.0 without changing product source, `go.mod`, or
  `go.sum`. Accept only the completed constructor-availability, deadline-error,
  priority-validation, newline-framing, API/behavior, platform, MVS,
  vulnerability, and related recorded findings. The exception is go-syslog-
  specific and non-transferable; it accepts no new or independently discovered
  defect. It remains valid only while exact v1.0.0 and both exact requests from
  Serf v0.8.2 and v0.9.6 remain unchanged, no direct main-module root or
  repository import is added, production and complete-test target loads remain
  zero, the module remains runtime-unreachable, every earlier guard remains
  intact, and no new advisory or independent defect appears. Direct import or
  loading, runtime reachability, a target version or incoming-request change, a
  new direct root, an earlier owning-guard change, or a new advisory or
  independent defect expires the exception and requires the owning fresh
  dependency and product decision before merge. Do not change product source
  or dependency metadata, add a direct edge, select an unreleased branch,
  change or remove either Serf parent, patch or fork source, authorize a wrapper
  or replacement architecture, raise the Go floor, move unrelated selections,
  alter another guarded dependency, manufacture a dependency commit, transfer
  another exception, or request this same choice again while all guards hold.
  The Go-1.18-compatible patch/fork/replacement and Serf parent/graph-removal
  alternatives were not authorized.
- Guard-only decision revalidation began from clean ordinary and ignored state
  at controller-authorized HEAD
  `37e05491d3a49fa2bb4f1be2f9be21c365b44884`, parent
  `0a223a920e48d1f52ad69c21bc1af5b1b6ae6ffe`, tree
  `65115211adcc75764f15f4d5093010b7ded45d1b`. The evaluation handoff
  `0a223a9` has exact parent `2ee06f1`, tree `d021d00`, and the recorded
  five-file change; the controller commit has the recorded four-file change.
  The branch, clean ordinary/ignored status, reciprocal 216-archive chain,
  launcher check, and exact Google UUID implementation ancestry pass.
- Exact Go 1.26.7 binary SHA-256
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  ran first in `PATH` with `GOENV=off`, `GOWORK=off`, `GOTOOLCHAIN=local`, no
  ambient `GOFLAGS`, `LC_ALL=C`, `LANG=C`, and `umask 022`. Exact go-syslog
  v1.0.0 retains both Serf v0.8.2/v0.9.6 requests and no direct root. All 19
  guarded selections and recorded requests remain exact; all 19 why/import
  results are negative; production and complete-test loads contain zero
  guarded packages; and 197 module-backed complete-test entries remain across
  41 loaded modules. Runtime unreachability therefore remains intact.
- The project remains 234 modules, 3,599 graph edges, 355 production entries,
  429 complete-test entries, 1,067 sum lines, and the recorded 432-line
  unapplied tidy projection. `go.mod`/`go.sum` retain SHA-256 values
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact module verification passes; no source, dependency metadata, parent,
  Go-floor, direct-root, or unrelated selection changed.
- Fresh primary vulnerability data remains byte-identical at 1,402 records,
  SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
  and Last-Modified 2026-09-17T17:29:18Z. Exact-version OSV remains empty for
  go-syslog and every guarded target except Gorilla and go-retryablehttp,
  which retain only their recorded accepted identifiers. Go-syslog's exact
  global GitHub and repository advisory results remain empty. No new advisory
  or independently observed defect appeared.
- The controller-authorized option 1 is recorded exactly under those guards.
  No dependency implementation or metadata commit was created, no changed-
  selection scorecard applies, and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. P7 continues only with the prepared bounded evaluation of selected
  exact-path `github.com/hashicorp/go-uuid v1.0.1`; it was not executed in
  this decision session.

Hashicorp go-uuid evaluation (2026-09-20):

- No exact-path stable `github.com/hashicorp/go-uuid` release qualifies.
  Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves
  Hashicorp's public active unarchived non-fork MPL-2.0 repository. The proxy
  exposes only stable v1.0.0, v1.0.1, v1.0.2, and v1.0.3; v1.0.3 is latest.
  There is no prerelease, GitHub Release object, retraction, module
  deprecation, redirect, alternate path, or `/v2` line. Ten branch heads and
  returned forks are not stable exact-path candidates.
- All four tags are lightweight. V1.0.0 is verified-signed commit
  `de160f5c59f693fed329e73e291bb751fe4ea4dc`; v1.0.1, v1.0.2, and v1.0.3
  are unsigned commits `4f571afc59f3043a65f8fe6bf46d887b10a01d43`,
  `6195a4f20692188e0a389def25559731d1e130f9`, and
  `67cd70b3bc50aeab4f2a02f869573616871a9a2a`. Their trees are respectively
  `453b27e64e2d644f0671c8723c0f9dd8e43ae127`,
  `ef50db991c26b953506884cab49165fa25dd7bca`,
  `8816e60600c501a4313b2ff3e81ff92bd00be69f`, and
  `3503fcc9c17dfa29bc601a50af54732413964dc2`. Proxy regular files match Git
  and sumdb verifies every release. V1.0.3 changes only license copyright from
  v1.0.2.
- Current master `f405b577e09f44ce7e9b484af2911859f9dab5a6`, tree
  `6440a80c2e8439d284fb19162f395e0256c4b76a`, is 28 commits past v1.0.3,
  declares Go 1.18, and is unreleased. It changes randomness error formatting
  from `%v` to `%w` but retains negative-size panic behavior. No branch, fork,
  redirect, alternate path, or unreleased commit was promoted.
- Each release contains one root package, no external module dependency, and
  a complete minimal closure of target plus standard library. There are no
  commands, examples, fuzz targets, testdata, generated files, build-tag
  branches, cgo, or embeds. Under exact Go 1.26.7 and contained Go 1.18.10 all
  stable releases pass source verification, listing, native count-one tests,
  two independent count-ten repeats, race, vet, and production/test
  cross-builds for Darwin amd64, Linux amd64/arm64, Windows amd64, and
  `js/wasm`.
- V1.0.0/v1.0.1 export four functions: random bytes, random UUID text,
  formatting, and parsing. V1.0.2/v1.0.3 compatibly add two caller-reader
  functions. Pinned API comparison finds no other change. The package is
  intentionally random-byte UUID text rather than RFC UUID generation; it
  does not set version or variant bits.
- Independent fixtures SHA-256
  `c2e61456d0cf65695c8a3b409cc9240da5d4941334223d5c6bc2f0bb304c843b`
  and
  `d192c87b504f65cf4bc68b48173cb38a8e9f82f861f32527d3d1a23996341727`
  characterize formatting/parsing, validation, mutation/non-aliasing, random
  behavior, caller-owned readers, allocation, concurrency, nil/panic, globals,
  resources, and environment interaction under both SDKs. Non-blocking
  repeats, race, and vet pass. V1.0.3 measures six format, one parse, and one
  random-byte allocation under both SDKs.
- Every stable release is disqualified by randomness-error identity loss:
  errors are formatted with `%v`, so `errors.Is` cannot recover the cause.
  Under Go 1.18, v1.0.0/v1.0.1 return the non-wrapping formatted error; under
  Go 1.26 their `crypto/rand.Read` route fatally terminates the process for an
  injected failing global reader. V1.0.2/v1.0.3 use `io.ReadFull` for their
  reader APIs but retain `%v` identity loss under both SDKs. All stable
  releases also panic for negative byte counts despite returning an error.
  Unreleased master passes error-identity fixtures through `%w` but retains the
  negative-size panic.
- MVS selects inherited v1.0.1 through six exact requests from Consul API
  v1.1.0/v1.12.0, Consul SDK v0.1.1/v0.8.0, and Serf v0.8.2/v0.9.6; both
  go-immutable-radix v1.0.0/v1.3.1 request v1.0.0. Why/import and both target
  loads are zero. Direct v1.0.1 adds only an edge and checksum; direct
  v1.0.2/v1.0.3 change only target selection and add an edge. Tidy removes
  every manufactured root and restores the exact base projection. Direct
  v1.0.0 breaks a guarded selection by removing mvn-pom-mutator v0.2.3 and
  makes loads fail. No projection was applied.
- Fresh primary vulnerability data remains 1,402 records at SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
  with no target record. Exact OSV and repository/global GitHub advisory
  queries are empty for every stable version. Exact-Go isolated scans are
  empty; old-Go reachable findings belong only to its standard library. Base
  and v1.0.1/v1.0.2/v1.0.3 project scans are identical, with zero target
  assignment or trace. Earlier guarded findings remain exact.
- The project remains 234 modules, 3,599 edges, 355 production entries, 429
  complete-test entries, 197 module-backed entries across 41 loaded modules,
  1,067 sum lines, and the 432-line unapplied tidy projection. All 20
  target-plus-earlier guarded selections and requests remain exact; all 20
  why/import results are negative and both loads contain zero guarded
  packages. `go.mod`/`go.sum` SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- The canonical exact-Go project gate passes verify, build, native/repeat/race/
  vet/lint, API/CLI, empty-HOME, four cross-builds, full preflight, every
  script/meta population, all eight mutation meta-stages with 80/80 kills,
  host/snapshot/Docker acceptance, and all 15 audit controls. The fresh
  clean-tree scorecard is 27/27 Q0-Q2 PASS at L2, SHA-256
  `62ee293e3c4253f4ed5c71df3c858a3db996636d0b6e968d8532b7fe2825ce99`,
  with the exact 21-line ledger. Documented cache, bare-`mktemp`, physical-tool,
  Make-inheritance, mode-mask, and launcher-timing reproductions are harness
  boundaries, not target findings.
- No dependency implementation or metadata commit was created. Product
  direction supplied 2026-09-20: the user-authorized automatic controller
  explicitly selected option 1 with the recommended bounds. Retain exact
  selected, inherited, unloaded go-uuid v1.0.1 without changing product
  source, `go.mod`, or `go.sum`. Accept only the completed randomness-error
  identity, Go 1.26 fatal injected-reader, negative-size panic, non-RFC
  generation, API/behavior, MVS, vulnerability, and related recorded findings.
  The exception is go-uuid-specific and non-transferable; it accepts no new or
  independently discovered defect. It remains valid only while exact v1.0.1,
  all six exact v1.0.1 requests from Consul API v1.1.0/v1.12.0, Consul SDK
  v0.1.1/v0.8.0, and Serf v0.8.2/v0.9.6, and both exact v1.0.0 requests from
  go-immutable-radix v1.0.0/v1.3.1 remain unchanged, no direct main-module
  root or repository import is added, production and complete-test target
  loads remain zero, the module remains runtime-unreachable, every earlier
  guard remains intact, and no new advisory or independent defect appears.
  Direct import or loading, runtime reachability, a target version or incoming-
  request change, a new direct root, an earlier owning-guard change, or a new
  advisory or independent defect expires the exception and requires the owning
  fresh dependency and product decision before merge. Do not change product
  source or dependency metadata, add a direct edge, select v1.0.3 or unreleased
  master, change or remove any Consul API, Consul SDK, Serf, or go-immutable-
  radix parent, patch or fork source, authorize a wrapper or replacement
  architecture, raise the Go floor, move unrelated selections, alter another
  guarded dependency, manufacture a dependency commit, transfer another
  exception, or request this same choice again while all guards hold. The
  Go-1.18-compatible patch/fork/replacement and parent-chain removal
  alternatives were not authorized. Revalidate only the guards, record the
  decision, answer its archive, and prepare one next bounded P7 mission without
  executing it in the decision turn.
- Guard-only decision revalidation began from clean ordinary and ignored state
  at controller-authorized HEAD
  `2774d4e54971b833c3eae3342ba158f21c1b139a`, parent
  `44f5328e66718401ae8d2574d550acb7c420bec8`, tree
  `a7657441d8b821f82f2ace8b1cb3165971c569df`. The preceding handoff retained
  its exact five-file delta and the controller commit its exact four-file
  delta; branch, ancestry, reciprocal 218-archive chain, latest Google UUID
  implementation identity, clean status, and launcher check passed.
- Exact Go 1.26.7 binary SHA-256
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  ran first in `PATH` under the recorded environment. Go-uuid retains exact
  v1.0.1 through all six exact v1.0.1 requests and both exact v1.0.0 requests,
  no direct root, a negative why result, zero repository imports, and zero
  production/test loads. All 20 guarded selections and requests remain exact;
  all 20 why results, repository imports, and guarded loads remain zero. The
  complete-test load retains 197 module-backed entries across 41 modules.
- The project remains 234 modules, 3,599 graph edges, 355 production entries,
  429 complete-test entries, and 1,067 sum lines. `go.mod`/`go.sum` retain
  SHA-256 values
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact module verification passes; no source, dependency metadata, parent,
  Go-floor, direct-root, or unrelated selection changed.
- Fresh primary vulnerability data remains byte-identical at 1,402 records,
  SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
  and Last-Modified 2026-09-17T17:29:18Z. Exact OSV results remain empty for
  go-uuid and every guard except the recorded Gorilla and go-retryablehttp
  identifiers. Go-uuid's exact global GitHub and repository advisory feeds
  remain empty. No new advisory or independently observed defect appeared.
- The controller-authorized option 1 is recorded exactly under those guards.
  No dependency implementation or metadata commit was created, no changed-
  selection scorecard applies, and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. P7 continues only with the prepared bounded evaluation of selected
  exact-path `github.com/hashicorp/go.net v0.0.1`; it was not executed in this
  decision session.

Hashicorp go.net evaluation (2026-09-20):

- No exact-path stable `github.com/hashicorp/go.net` release qualifies. Fresh
  proxy, sumdb, `go-import`, Git, and GitHub evidence resolves Hashicorp's
  public active unarchived non-fork BSD-3-Clause repository. V0.0.1 is the
  only proxy version, only stable tag, and proxy latest; there is no
  prerelease, GitHub Release, retraction, module deprecation, redirect,
  alternate exact path, or `/v2` line. Branches and returned forks are not
  stable candidates.
- V0.0.1 is a lightweight tag at GitHub-verified commit
  `afc3cb3a421746fc66dd55b09a270c750cf536ce`, parents
  `104dcad90073cd8d1e6828b2af19185b60cf3e29` and
  `b9a611abb799b370ad8aed695239fd086ae5a1fa`, tree
  `0a681728dc14ccfa08d532e5a20e85280c2fdb04`, time
  2019-01-18T17:37:16Z. Proxy ZIP SHA-256 is
  `71564aa3cb6e2820ee31e4d9e264e4ed889c7916f958b2f54c6f3004d4fcd8d2`;
  its regular files match Git and sumdb verifies both checksums. Current
  master `b69938bb7c99c1b9deeece78cc99714a94ec933b` is four unreleased,
  compliance-only commits later and retains the blockers.
- The release contains 273 regular files, 204 Go files, 46 test files, and 16
  packages spanning legacy context, dictionary, HTML/charset/IDNA/public-
  suffix, IPv4/IPv6, listener limiting, proxy/SOCKS5, SPDY v3, and WebSocket
  APIs. Build-ignored generators, generated tables/syscall files, examples,
  benchmarks, testdata, and all platform branches were inventoried; there is
  no release CLI or modern fuzz target.
- V0.0.1 has only a module directive and no requirements. Its complete
  production/test closure does not resolve under exact Go 1.26.7 or contained
  Go 1.18.10: production html/charset code and context/html/ipv4/ipv6/
  websocket tests retain obsolete `code.google.com/p/go.net` and
  `code.google.com/p/go.text` imports. Whole-module listing, build, count-one,
  two count-ten repeats, race, vet, and test-symbol scanning therefore fail
  under both SDKs. The mdns-used ipv4/ipv6 production packages build and mdns
  ordinary tests pass on Darwin, but target tests do not resolve and the
  packages fail Linux arm64 and `js/wasm` cross-builds. No complete Go-1.18
  source/test closure qualifies.
- Independent fixture SHA-256
  `46328fe4873c5266f6f8de9435dd96d105760401b20b54732dcbba858308d3d6`
  passes count-ten and race under both SDKs while reproducing negative-limit
  and non-Hijacker panics, caller-IP aliasing, process-global proxy
  registration, and SOCKS5 error-identity loss. Source inspection also
  records zero-limit blocking, unsynchronized global registration, absent
  WebSocket context/timeouts, entropy panic, peer-sized SPDY allocation,
  connection/deadline ownership, tokenizer aliasing, static-table age, and
  platform-stub boundaries. These are decision inputs, not qualification.
- MVS selects v0.0.1 only through exact mdns v1.0.0. There is no direct root,
  why/import results are negative, and both loads contain zero target
  packages. A disposable direct root changes no selected version or load and
  adds only one edge plus the content checksum; tidy removes it. That zero-
  load project projection passes verification, build, repeats, race, vet,
  API/CLI compatibility, host acceptance, and supported cross-builds. No
  projection was applied.
- Fresh primary vulnerability data remains 1,402 records at SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
  with no target entry. Exact OSV and GitHub advisory feeds are empty. Exact-
  Go production scans of ipv4/ipv6 are empty; old-Go findings belong only to
  its standard library. Base and direct-root project scans remain identical at
  30 module, 22 package, and 20 symbol/test-symbol IDs with no target trace.
- All 20 earlier guarded selections, requests, negative why results, zero Go-
  source imports, zero production/test loads, runtime-unreachability, and
  advisory states remain exact. The base stays 234 modules, 3,599 edges, 355
  production entries, 429 complete-test entries, 197 module-backed entries
  across 41 modules, 1,067 sum lines, and a 432-line tidy projection.
  `go.mod`/`go.sum` SHA-256 remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- Dependency metadata stayed unchanged. No changed-selection scorecard
  applies and accepted quality remains 27/27 Q0-Q2 PASS at L2. P7 stopped on
  the single prepared bounded go.net product decision among guarded retention,
  a separately authorized Go-1.18-compatible patch/fork/replacement design,
  or a separately authorized mdns parent/graph-removal decision. That decision
  was not executed in the evaluation session.

Hashicorp go.net product decision (2026-09-20):

- Option 1 is explicitly authorized: retain exact selected, inherited,
  unloaded `github.com/hashicorp/go.net v0.0.1` without changing product
  source, `go.mod`, or `go.sum`. The new target-specific, non-transferable
  exception accepts only the completed unresolved production/test closure;
  obsolete imports; Linux-arm64/wasm gaps; negative-limit and non-Hijacker
  panics; zero-limit blocking; SOCKS5 error-identity loss; `PerHost` aliasing;
  unsynchronized process-global registration; absent WebSocket context/
  timeouts; entropy panic; peer-sized SPDY allocation; and the recorded
  connection, deadline, resource, buffer-aliasing, static-table, allocation,
  non-concurrent-use, API, MVS, vulnerability, and related findings. It
  accepts no new or independently discovered defect.
- The exception remains valid only while exact v0.0.1 and its sole exact mdns
  v1.0.0 request remain unchanged, no direct main-module root or repository Go
  import is added, production and complete-test target loads remain zero, the
  module remains runtime-unreachable, every earlier guard remains intact, and
  no new advisory or independent defect appears. A target version or request
  change, direct root/import/load, runtime reachability, earlier owning-guard
  change, or new advisory or independent defect expires the exception and
  requires the owning fresh dependency and product decision before merge. The
  patch/fork/replacement and mdns parent/graph-removal alternatives are not
  authorized, and no exception transfers to another target.
- Decision revalidation began from clean ordinary and ignored state at HEAD
  `fa9af5e9b168520f024da6f2049065ad186735b8`, parent
  `70867a4a7701f69cbc4551de3d4db53fd7d86fe2`, tree
  `53de66fc8ea39ffdfa68611677564af70f487f8e`. Its exact five-file handoff
  delta, reciprocal 220-archive chain, latest Google UUID implementation
  ancestry, and launcher check pass. Exact Go 1.26.7 binary SHA-256 remains
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
- Go.net remains exact v0.0.1 through only
  `github.com/hashicorp/mdns@v1.0.0 github.com/hashicorp/go.net@v0.0.1`.
  There is no direct root, its why result remains negative, repository Go
  imports are zero, and production and complete-test target loads remain zero.
  All 20 earlier guarded selections and recorded requests remain exact; all 20
  earlier why/import/load results remain negative or zero. Go.net and every
  earlier target remain runtime-unreachable.
- The project remains 234 modules, 3,599 graph edges, 355 production entries,
  429 complete-test entries, 197 module-backed entries across 41 loaded
  modules, 1,067 sum lines, and the 432-line tidy projection. `go.mod` and
  `go.sum` retain SHA-256 values
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact module verification passes; no product source, parent, dependency
  metadata, direct root, Go floor, or unrelated selection changed.
- Fresh primary vulnerability data remains byte-identical at 1,402 records,
  SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
  with no target record. Exact OSV remains empty except for the recorded
  Gorilla and go-retryablehttp pairs; both exact go.net GitHub advisory feeds
  remain empty. No new advisory or independently observed defect appeared.
- No dependency implementation or metadata commit was created. No changed-
  selection scorecard applies, and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. P7 continues only with the prepared bounded evaluation of selected
  exact-path `github.com/hashicorp/golang-lru v0.5.4`; it was not executed in
  this decision session.

Hashicorp golang-lru evaluation and blocked product boundary (2026-09-20):

- No exact-path stable release qualifies. Fresh proxy, sumdb, `go-import`, Git,
  and GitHub evidence resolves Hashicorp's public active unarchived non-fork
  MPL-2.0 repository. Retrievable exact-path releases are v0.5.0-v0.5.4,
  v0.6.0, v1.0.1, and v1.0.2. A stale v1.0.0 proxy-list entry has no
  retrievable endpoint or Git tag. V1.0.2 is proxy latest and the final v0/v1
  release. Current main is the distinct Go-1.19
  `github.com/hashicorp/golang-lru/v2` module, with ARC also split to
  `github.com/hashicorp/golang-lru/arc/v2`; neither alternate major path was
  promoted. There is no retraction or module deprecation.
- Selected v0.5.4 is a lightweight unsigned tag at commit
  `14eae340515388ca95aa8e7b86f0de668e981f54`, parent
  `7f827b33c0f158ec5dfbba01bb0b14a4541fd81d`, tree
  `635407bf580448ac91cc1265a72cac484c273f21`, time
  2020-01-16T13:29:26-05:00. Proxy ZIP SHA-256 is
  `7b2a8b1739c858727fca497a6415323edb801dc97b8aca04f7bac4ab9fb5c66b`;
  proxy and Git files match, and sumdb verifies it. V0.5.4 is an ancestor of
  v0.6.0 and v1.0.2. Latest v1.0.2 is GitHub-verified commit
  `a032ef5a154020ffc7a74ba73702f9c5d3ea11f2`, parent v0.6.0, tree
  `ddabbfd3c054c45dccdc7c3a04378ea85eff781d`, time
  2023-08-08T11:23:39Z; its lightweight tag is unsigned.
- V0.5.4 contains two packages, six production Go files, four test files, 27
  tests, and six benchmarks. It has no command, example, fuzz target,
  testdata, generated file, build-tag/platform branch, cgo, embed, resource,
  environment, network, process, or goroutine boundary. V0.6.0/v1.0.2 add one
  standard-library-only ordinary testing helper file. V1.0.1 has no Go
  package or test and is unusable.
- Every sourceful release has a standard-library-only minimal closure and
  passes verification, listing, build, count-one, two count-ten repeats, race,
  and vet under exact Go 1.26.7 and contained Go 1.18.10. Selected v0.5.4,
  v0.6.0, and v1.0.2 also pass production/test cross-compilation for Darwin,
  Linux amd64/arm64, Windows, FreeBSD, and `js/wasm` under both SDKs. The
  complete dependency closure preserves the Go 1.18 floor.
- The root package exports synchronized Cache, 2Q, and ARC constructors and
  operations; `simplelru` exports its explicitly non-thread-safe LRU,
  interface, callback, and operations. Pinned API comparison proves
  v0.5.0-v0.5.3 remove selected methods or signatures. V0.6.0/v1.0.2 add
  only `DefaultEvictedBufferSize`, but changing `Cache` internals from an
  interface to a pointer plus slices makes the exported type non-comparable.
  Pinned apidiff and an independent `map[lru.Cache]bool` compile fixture prove
  that v0.5.4 compiles while v0.6.0/v1.0.2 fail under both SDKs.
- Independent fixture SHA-256
  `891d2144f6d6adb7acba5788859eb992b7c0238bbff2806f0e9444dfdd93f73e`
  characterizes LRU/2Q/ARC capacity, order, adaptation, resize, purge,
  callbacks, aliasing, nil/panic inputs, allocation, and concurrency. Selected
  v0.5.4 invokes user eviction callbacks while holding `Cache.lock`: callback
  reentrancy through `Len` deadlocks, and callback panic escapes before unlock
  and permanently poisons the mutex under both SDKs. V0.6.0/v1.0.2 invoke
  callbacks after unlocking and pass those fixtures but fail API comparability.
  Earlier releases already fail selected API compatibility.
- Common boundaries are recorded: Get/update refresh recency while Peek/
  Contains do not; caller key/value identity is retained; nil comparable keys
  and nil values work; non-comparable keys and nil/zero receivers panic;
  constructors reject nonpositive size; default `New2Q(1)` and zero ghost
  capacity fail; and `simplelru.Resize(-1)` evicts everything, returns an
  over-capacity count, retains negative capacity, and immediately evicts later
  adds. Cache/2Q/ARC are synchronized, simplelru is not, Get/Peek allocate
  zero, and there is no global mutable state or external resource ownership.
- MVS selects v0.5.4 through eighteen exact graph edges: three v0.5.4 requests
  from Viper v1.15.0/v1.10.1 and crypt v0.4.0; three v0.5.0 requests from
  immutable-radix v1.0.0/v1.3.1 and OpenCensus v0.21.0; and twelve v0.5.1
  requests from OpenCensus v0.22.0 plus the recorded historical Google API
  versions. There is no direct root; why/import results are negative; and
  production/complete-test loads contain zero target packages.
- A disposable direct v0.5.4 root retains the selection and all 234 modules,
  adds one graph edge and one sum line. Direct v0.6.0/v1.0.1/v1.0.2 roots
  move only the target, retain 234 modules and 3,600 edges, and produce 1,069
  sum lines. All retain 355 production and 429 complete-test entries with zero
  target loads and no unrelated version move. Tidy removes every manufactured
  root and returns exactly to common base projection hashes
  `5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
  `b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
  No projection was applied.
- Fresh primary vulnerability data remains 1,402 records at SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
  with no target record. Exact OSV and target GitHub advisory queries are
  empty for every retrievable stable release. Exact-Go isolated scans are
  empty; old-Go findings belong only to its standard library. Base and every
  project projection are byte-identical at 30 module, 22 package, and 20
  symbol/test-symbol IDs with no target assignment or trace.
- All 21 prior guarded selections, requests, negative why results, zero
  imports/loads, runtime-unreachability, and advisory states remain exact. The
  base stays 234 modules, 3,599 edges, 355 production entries, 429 complete-
  test entries, 197 module-backed entries across 41 modules, 1,067 sum lines,
  and a 432-line tidy projection. Base `go.mod`/`go.sum` hashes remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- All four direct-root projections pass exact-Go verify, build, count-one, two
  repeats, race, vet, API/CLI compatibility, host acceptance, and supported
  project cross-builds. The unchanged base additionally passes pinned lint,
  empty-HOME count-two, canonical full preflight, every script/meta pair, all
  eight mutation meta-stages with 80/80 kills, host/snapshot/Docker acceptance,
  and all audit meta-controls. The project retains its unrelated pre-existing
  `js/wasm` chzyer/readline gap. No changed-selection scorecard applies and
  accepted quality remains 27/27 Q0-Q2 PASS at L2.
- No dependency implementation or metadata commit was created. P7 is blocked
  on the single prepared golang-lru product decision among exact v0.5.4
  guarded retention, a separately authorized Go-1.18-compatible exact-path
  and API-preserving remediation design, or a separately authorized parent/
  graph-removal study. That decision was not executed in this evaluation.

Hashicorp golang-lru product decision (2026-09-20):

- Option 1 is explicitly authorized: retain exact selected, inherited,
  unloaded `github.com/hashicorp/golang-lru v0.5.4` without changing product
  source, `go.mod`, or `go.sum`. This bounded direction uses the verified
  zero-load exposure while avoiding a manufactured direct root and the
  exported `Cache` comparability regression in v0.6.0/v1.0.2; it is an
  explicit product choice, not an inference from physical selection. The new
  target-specific, non-transferable
  exception accepts only the completed callback-under-lock reentrancy
  deadlock; callback-panic mutex poisoning; nonpositive-capacity constructor
  errors; default `New2Q(1)` and zero-ghost-ratio failures; negative
  `simplelru.Resize` over-capacity counting, retained negative capacity, and
  immediate later eviction; nil/zero-receiver and non-comparable-key panics;
  caller key/value identity and aliasing; explicitly non-concurrent simplelru
  use; and the recorded API, behavior, MVS, vulnerability, allocation,
  resource, and related completed findings. It accepts no new or independently
  discovered defect.
- The exception remains valid only while exact v0.5.4 and all eighteen
  recorded incoming requests remain unchanged: v0.5.4 from Viper v1.15.0,
  historical Viper v1.10.1, and `sagikazarmark/crypt v0.4.0`; v0.5.0 from
  go-immutable-radix v1.0.0/v1.3.1 and OpenCensus v0.21.0; and v0.5.1 from
  OpenCensus v0.22.0 plus Google API v0.7.0-v0.9.0, v0.13.0-v0.15.0,
  v0.17.0-v0.20.0, and v0.22.0. It also requires no direct root or repository
  import, zero production and complete-test target loads, runtime
  unreachability, every earlier guard remaining intact, and no new advisory or
  independent defect. A target/request/root/import/load/runtime, earlier-
  guard, or new-finding change expires the exception and requires the owning
  fresh decision. The exact-path remediation and parent/graph-removal
  alternatives are not authorized. V0.6.0/v1.0.2's exported `Cache`
  comparability regression, the package-empty v1.0.1 release, v2 paths, main,
  forks, and alternate paths remain unauthorized, and no exception transfers.
- Decision revalidation began from clean ordinary and ignored state at HEAD
  `961c48582ca568f0867d91e7781c6fa0b970dcea`, parent
  `f867a05fa1334c7e7fc0b774edc060e78a9be88d`, tree
  `716bf2f170c296836aa1903414c55ab4f5450109`. Its exact five-file handoff
  delta, reciprocal 222-archive chain, latest Google UUID implementation
  ancestry, and launcher check pass. Exact Go 1.26.7 binary SHA-256 remains
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
- Golang-lru remains exact v0.5.4 through all eighteen recorded requests. It
  has no direct root, its why result is negative, repository Go imports are
  zero, and the 355-entry production and 429-entry complete-test loads contain
  zero target packages. All 21 earlier guarded selections and requests remain
  exact; all 22 why results are negative; imports and guarded loads remain
  zero; and every target remains runtime-unreachable.
- The project remains 234 modules, 3,599 graph edges, 197 module-backed
  complete-test entries across 41 loaded modules, 1,067 sum lines, and the
  432-line tidy projection. `go.mod`/`go.sum` SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact module verification passes; no product source, parent, dependency
  metadata, direct root, Go floor, or unrelated selection changed.
- Fresh primary vulnerability data remains byte-identical at 1,402 records,
  SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
  with no target record. Exact-version OSV remains empty except for the
  recorded Gorilla and go-retryablehttp pairs; both exact golang-lru GitHub
  advisory feeds remain empty. No new advisory or independently observed
  defect appeared.
- No dependency implementation or metadata commit was created. No changed-
  selection scorecard applies, and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. P7 continues only with the prepared bounded evaluation of selected
  exact-path `github.com/hashicorp/mdns v1.0.4`; it was not executed in this
  decision session.

Hashicorp mdns evaluation and blocked product boundary (2026-09-20):

- No exact-path stable release qualifies under the combined Go-1.18 floor,
  complete source/test closure, API, behavior, and project contracts. The
  exact path resolves Hashicorp's public active unarchived non-fork MIT
  repository and eight stable releases v1.0.0-v1.0.7. There is no retraction,
  deprecation, redirect, alternate major line, or qualifying fork/branch.
  Every proxy archive matches Git and sumdb verifies it. Tags are lightweight
  and unsigned; the tagged commits are GitHub-verified.
- Selected v1.0.4 is commit
  `f504e0b3c01cf50e077ab6805546fa95ed9f2933`, parent
  `34dc81282184e8c94b4ed4312182d67b99827a5e`, tree
  `50afb5c00de8ee437c229be31e0d0c5c0671f713`, time
  2021-04-14T13:17:39-05:00. Its proxy ZIP SHA-256 is
  `f7d04f484c5b398018385ad954a922b0d94d0cb26315780bf995bedede166704`;
  Git/proxy bytes match and sumdb verifies it. It is an ancestor of
  v1.0.5-v1.0.7 and main.
- Selected contains one library package, three production files, two tests,
  and no command, example, benchmark, fuzz target, testdata, generated file,
  or build constraint. V1.0.0-v1.0.4 share its exported API; v1.0.5 adds
  compatible address-family controls; v1.0.6/v1.0.7 add context, logger, and
  IPv6-zone API.
- V1.0.0 has a vulnerable DNS closure, real shutdown races, and a vet defect;
  v1.0.1 retains the vet/vulnerability failures. V1.0.2-v1.0.4 preserve the
  Go floor but exact Go 1.26.7 cannot link their Darwin tests through old
  x/net's removed `syscall.recvmsg` reference; v1.0.3/v1.0.4 also fail native
  Go 1.18.10 IPv6 multicast routing. V1.0.5 fixes that floor-SDK condition but
  retains the exact-Go Darwin test-link failure. V1.0.6's loaded miekg/dns
  closure requires Go 1.19. V1.0.7 declares Go 1.25, loads Go-1.24 modules,
  and cannot compile under Go 1.18.
- Independent fixture SHA-256
  `73739d0dcb546dcb5f7edcdb645d9c964444543296f5f1d8e74f9da22011967f`
  records unchecked wrapped ports, trailing-dot-only name validation, IP/TXT
  aliasing, inconsistent records after exported-field mutation, nil panics,
  last-only delivery for multi-entry packets, closed-channel panic, silent
  slow/nil-channel drops, requester-directed multicast responses, caller
  parameter mutation, unsolicited-answer acceptance, timeout without context,
  no goroutine join, partial socket/resource and route dependencies,
  unsynchronized caller-owned state, global logging/environment interaction,
  eight record allocations, and the other completed API/error/concurrency/
  ownership boundaries. Exact Go cannot link this fixture through the same
  old-x/net defect; its Go-1.18 count-one, two repeats, and race runs pass.
- MVS selects v1.0.4 only through Serf v0.9.6; historical Serf v0.8.2 requests
  v1.0.0, which is also the sole request preserving guarded go.net v0.0.1.
  Mdns has no direct root/import, its why result is negative, and zero target
  packages load in 355 production or 429 complete-test entries. Direct
  v1.0.0-v1.0.3 projections broadly collapse guarded selections and fail
  project loading. V1.0.4/v1.0.5 manufacture target/DNS roots but tidy exactly
  to the common base. V1.0.6/v1.0.7 move unrelated x modules; v1.0.7 also
  raises the main Go directive to 1.25 and fails project test/race/vet through
  stricter format checks. No projection was applied.
- Fresh 1,402-record primary vulnerability data remains at SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
  with no mdns entry. Every exact stable OSV result and both selected GitHub
  advisory feeds are empty. V1.0.0/v1.0.1 reach a DNS advisory; later target
  closures have no called advisory. Base/v1.0.4/v1.0.5 project scans remain
  30 module, 22 vulnerable-package, and 20 called/test-symbol IDs with no mdns
  assignment. Later projections reduce counts only by moving unrelated x
  modules.
- All 22 earlier guarded selections, 162 incoming edges, negative why/import/
  load results, runtime-unreachability conditions, and advisory states remain
  exact. The base remains 234 modules, 3,599 edges, 355/429 load entries, 197
  module-backed entries across 41 modules, 1,067 sums, and a 432-line tidy
  projection. `go.mod`/`go.sum` hashes remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- No source or metadata changed, no dependency commit exists, and no changed-
  selection scorecard applies. Accepted quality remains 27/27 Q0-Q2 PASS at
  L2. P7 is blocked on one prepared mdns decision: exact v1.0.4 guarded
  retention, a separately authorized Go-1.18-compatible remediation, or a
  separately authorized Serf/graph-removal study. It was not executed.

Hashicorp mdns product decision (2026-09-20):

- Option 1 is explicitly authorized: retain exact selected, inherited,
  unloaded `github.com/hashicorp/mdns v1.0.4` without changing product source,
  `go.mod`, or `go.sum`. This is an explicit bounded product choice for the
  verified zero-load graph, not an inference from physical selection. It
  avoids manufacturing a direct mdns/DNS root, moving unrelated graph
  selections, or changing the Go floor.
- The new target-specific, non-transferable exception accepts only the
  completed exact-Go Darwin test-link failure through old x/net's removed
  `syscall.recvmsg`; the Go-1.18 IPv6 multicast-routing test failure; unchecked
  wrapped ports; trailing-dot-only name validation; caller IP/TXT and returned
  TXT aliasing; inconsistent records after exported-field mutation; nil
  service receiver, nil config, nil Zone, and internal UDP-address assertion
  panics; last-only multi-entry-packet delivery; closed-result-channel panic;
  silent nil, unbuffered, and slow-channel drops; requester-directed multicast
  responses; caller `QueryParam` mutation; unsolicited-answer acceptance;
  timeout without context or cancellation; absent goroutine lifecycle joins;
  partial socket setup and resource cleanup; interface, route, hostname, DNS,
  process-global logging, and unsynchronized caller-owned-state boundaries;
  eight record allocations; and the recorded API, error, protocol,
  packet/channel, concurrency, platform, ownership, environment, MVS, and
  vulnerability findings. The vulnerability bounds include empty exact
  v1.0.4 and project mdns results and the recorded reachable DNS advisory in
  the historical v1.0.0/v1.0.1 closure. It accepts no new advisory or
  independently discovered defect.
- The exception remains valid only while exact selected v1.0.4, the exact
  Serf v0.9.6 -> mdns v1.0.4 request, the historical Serf v0.8.2 -> mdns
  v1.0.0 request, and the exact mdns v1.0.0 -> go.net v0.0.1 request remain
  unchanged. It additionally requires no direct main-module root or repository
  Go import, zero production and complete-test target loads, runtime
  unreachability, every earlier owning guard remaining intact, and no new
  advisory or independent defect. A target selection or recorded request
  change, direct root/import/load, runtime reachability, earlier-guard change,
  or new advisory or independent defect expires the exception and requires
  the owning fresh dependency and product decision before merge.
- The exact-path remediation and Serf/graph-removal alternatives are not
  authorized. V1.0.5-v1.0.7, main, forks, alternate paths, direct roots,
  patches, parent changes, unrelated selection moves, and a Go-floor change
  remain unauthorized. No exception transfers to another target.
- Guard-only revalidation began from clean ordinary and ignored state at HEAD
  `9edb38c51e0067672e90224d8665bae580f92a05`, parent
  `b641b55ec01f18a79ecf57b8e5e866310dd44562`, tree
  `63a43a0d26649532e760560a9e114fea47feb8d1`. Its exact five-file handoff
  delta, reciprocal 224-archive chain, latest Google UUID implementation
  ancestry, exact Go 1.26.7 binary identity, and launcher check pass.
- Exact mdns v1.0.4, both incoming Serf requests, and the historical mdns-
  v1.0.0-to-go.net request remain exact. Mdns has no direct root, its why
  result is negative, repository Go imports are zero, and 355 production plus
  429 complete-test entries contain zero target packages. All 22 earlier
  guarded selections and their 162 incoming edges remain exact at recorded
  snapshot SHA-256
  `73b2f342ab8b1f5334afedead3c49cee97ee0177bf2413ea13d10f5d45995c63`;
  their why/import/load results remain negative or zero. With mdns included,
  23 guarded selections and 164 incoming edges retain snapshot SHA-256
  `f42304a43a7ade746c5c2bc238dfa9a6daec62594834fd5b606fdf245d44ea2e`.
- The project remains 234 modules, 3,599 graph edges, 197 module-backed
  complete-test entries across 41 loaded modules, 1,067 sum lines, and the
  432-line tidy projection. `go.mod`/`go.sum` SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  Exact module verification passes; no product source, parent, dependency
  metadata, direct root, Go floor, or unrelated selection changed.
- Fresh primary vulnerability data remains byte-identical at 1,402 records,
  SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
  Last-Modified 2026-09-17T17:29:18Z, with no mdns record. Exact-version OSV
  remains empty except for the recorded Gorilla and go-retryablehttp pairs;
  mdns v1.0.4 and both exact GitHub advisory feeds remain empty. No new
  advisory or independently observed defect appeared.
- No dependency implementation or metadata commit was created. No changed-
  selection scorecard applies and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. P7 continues only with the prepared bounded evaluation of selected
  exact-path `github.com/hashicorp/memberlist v0.3.0`; it was not executed in
  this decision session.
- The first memberlist execution attempt began later from clean HEAD
  `a1dee3a81a71e0e6656de09fef603b94cc135bab` and was automatically stopped by
  a cybersecurity policy check without a tracked worktree change. It retained
  provisional public evidence that HashiCorp advisory `HCSEC-2026-18` /
  `CVE-2026-14362` affects releases through v0.5.4, v0.6.0 is the first fixed
  release, and v0.6.0 declares Go 1.25. The active retry is restricted to
  narrow primary-source confirmation, static remediation identity, ordinary
  non-adversarial checks, and a product-decision handoff. It prohibits exploit
  reproduction, custom attack traffic, malformed-packet fixtures, scans,
  credential access, and security-boundary testing. If confirmed, no fixed
  exact-path memberlist release preserves Go 1.18; dependency metadata remains
  unchanged and the successor must decide among a floor/fixed-line move,
  explicit target-specific retention, or a separately measured owning-parent
  or graph-removal study.

Hashicorp memberlist defensive evaluation and blocked product boundary
(2026-09-20):

- No exact-path stable release qualifies. Primary HashiCorp advisory
  `HCSEC-2026-18` and its published `CVE-2026-14362` CNA record confirm that
  selected v0.3.0 is affected, releases through v0.5.4 are affected, and
  v0.6.0 is first fixed. Current latest v0.7.0 is also fixed. Exact OSV and
  GitHub feeds are still empty for the target and do not override the primary
  vendor record. The defensive retry did not reproduce or assess
  exploitability, run the advisory regression, craft packets/attack traffic,
  scan hosts, inspect credentials, or test a security boundary.
- The exact path resolves Hashicorp's public active unarchived non-fork
  MPL-2.0 repository with 24 stable releases v0.1.0-v0.7.0. The suffix tag
  v0.3.1-metrics-labels is a prerelease and was not promoted. There is no
  retraction, deprecation, redirect, alternate major line, or qualifying
  branch/fork. Serious tags are lightweight/unsigned and their commits are
  GitHub-verified. Selected, last-affected, and fixed proxy archives match Git
  byte-for-byte and sumdb verifies them.
- Selected v0.3.0 is commit
  `923f1b205dd4653f2ea35e1c9531088e52053aa0`, parent
  `123f3fbfeacd70bbe467775ab563a2b87e8c5cd8`, tree
  `5c0926d642da064035733e5a591214269c223805`, with proxy ZIP SHA-256
  `77e9fd374266825d7875a21a099cd80675bf9c8d469290f79c0a41bccba9df30`.
  First fixed v0.6.0 is commit
  `371698b4dc493ecf2fc94c4383c75b2c9358f6c9`, parent
  `0000b77c906dc53889a10de7a12fe4d146e84c50`, tree
  `9009580c109364895d764539ede0498fc086254b`, proxy ZIP SHA-256
  `b2c7ebd03f36b0dd22604e65b73b1e87ee1ce75eb36a333b1b6666b1f54497ff`.
  Latest fixed v0.7.0 is commit
  `4e3e17fc20e38ac48982ab6f65fec243fd139ef9`, tree
  `1e2d57c8a19d11776ed294cd4692e91c6e43aac6`, proxy ZIP SHA-256
  `a77ea0b0c1bc0120cdcaf5ba34b1e95971b566d05f9eba7ac36191aedd83ff4d`.
- Static remediation identity is fix commit/v0.6.0
  `371698b4dc493ecf2fc94c4383c75b2c9358f6c9`, which changes only `net.go`
  and `net_test.go` by 67 insertions and one deletion. The email/binary patch
  SHA-256 is
  `f40e2a2034c07db2df2b2b76e76f09e3f67cca1289fe66a986d124a057cdb020`.
  It is absent from v0.5.4 and retained by v0.6.0/v0.7.0.
- No fixed release preserves the floor. Selected v0.3.0 declares Go 1.12 and
  its closure preserves Go 1.18. V0.5.1-v0.5.3 declare Go 1.20, last affected
  v0.5.4 declares Go 1.24.0, and fixed v0.6.0/v0.7.0 declare Go 1.25.0.
  Both fixed closures load Go-1.24/Go-1.25 modules and Go 1.18.10 rejects
  their directives. Exact Go 1.26.7 package listing, source/test and race
  compilation without executing test bodies, vet, and relevant cross-builds
  pass for both fixed releases. Selected v0.3.0 passes those compile/vet
  controls under Go 1.18.10 but exact Go cannot link its root tests through
  old x/net's removed `syscall.recvmsg` reference.
- V0.3.0 is a library-only, no-source-build-tag module. Its exported root API
  covers configuration/defaults and creation; member join/update/query/send/
  leave/shutdown; node, transport, delegate/event, keyring, broadcast queue,
  label, and logging abstractions. Static lifecycle review records background
  listeners/tickers after `Create`, config and returned-node mutation rules,
  delegate concurrency/blocking/borrowed-data rules, repeat-safe `Leave` and
  `Shutdown`, the leave-after-shutdown and oversized-metadata panics, and the
  difference between leave broadcast and shutdown teardown. `internal/retry`
  is internal test support, there is no CLI, and the remediation commit changes
  no exported declaration. Zero project import/load keeps this surface and
  lifecycle out of the product runtime and API/CLI.
- MVS selects v0.3.0 only through Serf v0.9.6; historical Serf v0.8.2 requests
  v0.1.3. There is no direct root, target why/import/load results are negative
  or zero, and the project remains 234 modules, 3,599 edges, 355 production
  entries, 429 complete-test entries, 197 module-backed entries across 41
  loaded modules, and 1,067 sums. A redundant v0.3.0 root changes no
  selection and tidies to the common base projection.
- A disposable v0.6.0 root raises the main directive to Go 1.25.0 and yields
  259 modules, 3,743 edges, and 1,103 sums; it changes four guarded
  selections. V0.7.0 yields 257 modules, 3,709 edges, and 1,095 sums and
  changes five guarded selections. Both retain 355/429 project loads with zero
  target packages, but project test/vet fails existing non-constant format
  calls under the raised floor. Their tidy projections retain Go 1.25 and
  non-base graph/hash results. Nothing was applied.
- All 23 earlier guarded selections and 164 incoming edges remain exact at
  snapshot SHA-256
  `f42304a43a7ade746c5c2bc238dfa9a6daec62594834fd5b606fdf245d44ea2e`;
  all guarded why/import/load results remain negative or zero. Fresh primary
  data remains 1,402 records at SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
  Earlier exact OSV state retains only the recorded Gorilla and
  go-retryablehttp pairs. Base `go.mod`/`go.sum` hashes remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- No source or dependency metadata changed, no dependency commit exists, no
  changed-selection scorecard applies, and accepted quality remains 27/27
  Q0-Q2 PASS at L2. P7 is blocked on the one prepared memberlist product
  decision among an integrated Go-1.25/fixed-line migration, explicit guarded
  retention of affected but unloaded v0.3.0, or a separately measured Serf/
  graph-removal study. That decision was not executed.

Hashicorp memberlist product decision (2026-09-20):

- Option 3 is explicitly authorized. Perform one separate, measurement-only
  owning-parent and graph-removal study for exact historical Serf v0.8.2/
  v0.9.6 and the request population that preserves memberlist. Independently
  measure every serious Serf upgrade, replacement, and removal route and all
  transitive API/CLI, runtime-behavior, MVS, guarded-edge, Go-floor,
  vulnerability, build, test, platform, and project effects. The study may
  prepare one later product decision but may not implement a route or change
  Serf, memberlist, a parent, another guarded module, product source, `go.mod`,
  or `go.sum`.
- This is neither an affected-version risk acceptance nor a fixed-line
  migration. Selected memberlist v0.3.0 remains affected and unqualified; it
  must not be described as secure. V0.6.0/v0.7.0, a Go-floor change, stricter-
  format repairs, direct roots, patches, forks, alternate paths, parent
  changes, and guarded selection moves remain unauthorized. Zero loading
  bounds current exposure but creates no exception.
- The study authorization remains valid only while selected memberlist is
  exact v0.3.0; the Serf-v0.9.6-to-memberlist-v0.3.0 and historical
  Serf-v0.8.2-to-memberlist-v0.1.3 requests remain exact; selected Serf is
  exact v0.10.1; and the five incoming Serf requests remain exact: Consul API
  v1.1.0 -> Serf v0.8.2; Consul API v1.12.0, crypt v0.4.0, and Viper v1.10.1
  -> Serf v0.9.6; and Viper v1.15.0 -> Serf v0.10.1. It also requires no
  direct memberlist root/import/load or runtime reachability, every earlier
  owning guard and module hash remaining exact, and no new memberlist
  advisory, independent defect, affected/fixed boundary, or Go-1.18-compatible
  fixed stable release. Any difference expires the authorization. It also
  expires after the one coherent study handoff and grants no implementation
  authority.
- Guard-only revalidation began from clean ordinary and ignored state at HEAD
  `636e6b4f886ee15947c954deee70365c9eab1873`, parent
  `8a292c9a10f25aeae5be9e30693a7ed27f91912d`, tree
  `43dbc70674569af72990e06d1b16c38cda6f6dd7`. Its exact five-file handoff
  delta, reciprocal 226-archive chain, latest Google UUID implementation
  ancestry, exact Go 1.26.7 identity, and launcher check pass.
- MVS selects memberlist v0.3.0 only through the Serf v0.9.6 request;
  historical Serf v0.8.2 still requests v0.1.3. There is no direct root, its
  why result and repository imports are negative or zero, production/test
  loads are zero, and runtime unreachability holds. All 23 earlier guarded
  selections and 164 incoming edges remain exact at snapshot SHA-256
  `f42304a43a7ade746c5c2bc238dfa9a6daec62594834fd5b606fdf245d44ea2e`;
  all 24 target-plus-earlier why results are negative and guarded imports and
  loads are zero. The project remains 234 modules, 3,599 edges, 355/429 load
  entries, 197 module-backed complete-test entries across 41 modules, 1,067
  sum lines, and the 432-line tidy projection. `go.mod`/`go.sum` hashes remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- Fresh primary revalidation preserves the decision premise. HCSEC-2026-18
  still identifies memberlist through v0.5.4 as affected and v0.6.0 as fixed.
  The PUBLISHED HashiCorp CNA record remains affected below v0.6.0 at response
  SHA-256
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
  The Go vulnerability module index remains byte-identical at 1,402 records,
  SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
  Exact-version OSV remains empty except for the recorded Gorilla and
  go-retryablehttp pairs; empty memberlist OSV does not override the primary
  records. No source or metadata changed, no dependency implementation commit
  exists, and accepted quality remains 27/27 Q0-Q2 PASS at L2. The sole
  prepared successor is the bounded Serf/owning-graph removal study; it was not
  executed. Exact-Go module verification, build, canonical-`umask 022`
  count-one tests, race tests, vet, launcher check, and all 62 launcher
  lifecycle controls pass. Initial test/race runs inherited `umask 077` and
  reproduced only the known fixture-mode assertions before the canonical
  reruns passed.

Hashicorp Serf/owning-graph removal study (2026-09-21):

- Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves
  HashiCorp's public active unarchived non-fork MPL-2.0 Serf repository. It has
  37 stable exact-path releases v0.1.0-v0.11.0 plus one suffix prerelease, no
  retraction, deprecation, redirect, alternate major, or supported
  replacement. Serious proxy archives match Git except for correctly omitted
  v0.8.x vendor trees, and sumdb verifies them. V0.8.2 is an annotated signed
  tag with unknown-key forge verification; other serious tags are lightweight
  with verified commits and form one ancestry.
- The exact request population is owned by direct Viper v1.15.0 -> selected
  Serf v0.10.1 and mvn-pom-mutator v0.2.3 paths through Viper v1.10.1, Consul
  API v1.1.0/v1.12.0, and both crypt modules to historical Serf v0.8.2/v0.9.6.
  MVS selection does not erase historical vertices. Viper's pruned v1.15.0
  graph does not expose selected Serf's descendants, so historical v0.9.6
  continues to select memberlist v0.3.0 and v0.8.2 records v0.1.3.
- No fixed stable Serf release preserves the Go 1.18 closure floor. V0.10.1
  requests affected memberlist v0.5.0 and closes at Go 1.17; v0.10.2 requests
  affected v0.5.2 and closes at Go 1.20; v0.10.3/v0.10.4 request first-fixed
  v0.6.0 and declare/close at Go 1.25; latest v0.11.0 requests v0.7.0,
  declares/closes at Go 1.26, and incompatibly changes exported metrics-label
  fields. Patch-line API/CLI behavior is otherwise measured and ordinary
  source/test, race, vet, and cross-build failures are classified by obsolete
  dependency, SDK, or contained platform boundary.
- Disposable direct Serf roots select memberlist v0.5.0, v0.5.2, v0.6.0, or
  v0.7.0 and produce up to 246 modules/3,652 edges/1,079 sums. V0.10.2+
  changes five guarded selections. Fixed roots raise the main directive and
  make project tests/race/vet fail existing Go-1.25 format checks. Tidy removes
  every manufactured Serf root and restores selected v0.10.1/memberlist
  v0.3.0, so no direct root is durable remediation. Replacements either leave
  the historical selection intact or become unsupported version masquerading
  with a raised floor and much larger graph.
- Dropping mvn-pom-mutator is the only graph-removal projection. It removes
  memberlist and historical Serf, reducing 234 modules/3,599 edges to
  156/2,218, but also removes 15 guarded selections and 91 guarded incoming
  edges. The project imports the parent POM model/parser/mutator API across 21
  source/test files; listing and compilation fail and tidy re-adds v0.2.3.
  V0.2.3 is already latest stable/default-branch head, and inspected
  alternatives are non-drop-in subsystem migrations. Removal is therefore a
  separate product redesign, not a bounded dependency remediation.
- Exact vulnerability guards, all 24 negative why results, zero guarded
  imports/loads, runtime unreachability, 23 guarded selections/164 edges,
  module hashes, and 27/27 Q0-Q2 PASS at L2 remain unchanged. No source,
  `go.mod`, or `go.sum` change was retained. Neither bounded removal nor
  compatible modernization is viable. The sole prepared successor is a final
  product decision between integrated fixed-memberlist migration and explicit
  affected-version risk acceptance; the study recommends the fixed migration
  but does not authorize or implement it. Final exact-Go verify/build/test/
  race/vet, all 15 audit meta-controls, and all 62 launcher lifecycle controls
  pass. A focused raw Q0-Q2 audit retains the expected 21 automated PASS, six
  manual-evidence-bound UNMEASURABLE rows, and zero ratchet regressions; no
  changed-selection scorecard applies, so accepted quality remains 27/27
  Q0-Q2 PASS at L2.

Hashicorp memberlist final product decision (2026-09-21):

- Option 1 is explicitly authorized: perform exactly one separate integrated
  Go-1.25/fixed-memberlist migration. The successor must choose and fully
  qualify exact first-fixed v0.6.0 or latest fixed v0.7.0, establish the fixed
  selection through a supported tidy-stable owning dependency relationship,
  explicitly own raising the module floor from Go 1.18 to Go 1.25.0, and
  repair or separately decide the existing stricter-format failures.
- This is not an affected-version exception. Selected inherited v0.3.0 remains
  affected, unqualified, and not secure. No direct root, replacement, version
  masquerade, fork, patch, silent parent/POM redesign, unrelated dependency
  modernization, or transfer of another exception is authorized. If no
  supported durable selection or fully passing integration exists, the
  successor must retain no projection and stop P7 for a fresh product decision.
- Every changed guarded selection and incoming edge requires a measured fresh
  owning decision. Before retaining implementation, the chosen candidate must
  preserve project API/CLI/runtime behavior and pass complete source/test,
  repeated-test, race, vet, lint, platform/floor, MVS, tidy, vulnerability,
  acceptance, launcher/audit-control, and Q0-Q2 quality gates. Dependency and
  behavior changes remain separated under the roadmap's normal commit rules.
- Guard-only revalidation began from clean ordinary and ignored state at HEAD
  `e1ff4d49490ed94513d54807b5e19cebd6faaacc`, parent
  `36bf3a3f5a05594213e0ccd6f5c45a2e2f3bbf3a`, tree
  `5063148c562d6d298aa8555f52f7e56f14260dde`. Its exact five-file predecessor
  handoff, reciprocal 228-archive chain, latest Google UUID implementation
  ancestry, exact Go 1.26.7 binary, and launcher check pass.
- The unchanged base remains 234 modules, 3,599 edges, 355/429 load entries,
  197 module-backed complete-test entries across 41 modules, 1,067 sum lines,
  and a 432-line tidy projection. All 24 target-plus-earlier why results are
  negative; repository imports and production/complete-test guarded loads are
  zero. All 23 earlier guarded selections and 164 incoming edges retain
  snapshot SHA-256
  `f42304a43a7ade746c5c2bc238dfa9a6daec62594834fd5b606fdf245d44ea2e`.
  The exact five Serf requests and two historical Serf-to-memberlist requests
  remain unchanged. `go.mod`/`go.sum` retain SHA-256 values
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- Primary advisory identity is unchanged. The PUBLISHED HashiCorp CNA response
  still marks memberlist below v0.6.0 affected and v0.6.0 first fixed at
  SHA-256
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
  The 1,402-record Go vulnerability index remains byte-identical at SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
  Last-Modified 2026-09-17T17:29:18Z. Exact OSV remains empty for memberlist
  and all Serf candidates; earlier guards retain only the recorded Gorilla and
  go-retryablehttp pairs. No security investigation was repeated.
- The authorization is non-transferable and requires this exact base, no new
  advisory or independently observed defect, an unchanged affected/fixed
  range, and no Go-1.18-floor-compatible fixed stable route. Expected migration
  deltas are valid only when explicitly measured and owned; any unexplained
  target, request, root, import, load, runtime, hash, guard, advisory, range,
  or route change expires the authorization. It also expires after one
  coherent migration handoff and authorizes neither another group nor P8.
- The decision session changed only continuity documentation and the launcher
  handoff. It did not edit source, `go.mod`, or `go.sum`; choose or select a
  fixed release; raise the floor; change a parent or guarded module;
  risk-accept v0.3.0; or execute the migration successor. Accepted quality
  remains 27/27 Q0-Q2 PASS at L2. Final exact-Go module verification,
  count-one tests, race tests, vet, all 62 launcher lifecycle controls, and all
  15 audit meta-controls pass. The canonical reruns resolved only the already
  documented launcher signal-log timing race and noncanonical Homebrew Go
  1.26.2 path boundary; exact Go 1.26.7 reproduces the stored audit baseline.

Integrated Go-1.25/fixed-memberlist migration outcome (2026-09-21):

- Exact first-fixed memberlist v0.6.0 was preferred over latest fixed v0.7.0.
  It is the primary-vendor first fixed release, avoids v0.7.0's additional
  guarded move, and is reached through Serf v0.10.4 on Go 1.25 rather than
  Serf v0.11.0's Go 1.26 floor and exported metrics-label change. No fixed
  implementation qualified because the repository has no genuine supported
  tidy-stable owner for either fixed release.
- Direct memberlist and Serf roots remain removable by tidy; replacements are
  unsupported version masquerading; excluding memberlist v0.3.0 selects
  historical v0.1.3, and excluding both historical requests removes the
  selected module instead of selecting v0.6.0. Memberlist has no command.
  Blank, build-tag, test-only, or unused tool anchors would manufacture
  ownership and are not remediation.
- A disposable Go 1.25 `tool` declaration for
  `github.com/hashicorp/serf/cmd/serf@v0.10.4` was the only syntactically
  durable route. Tidy retained `go 1.25.0`, `toolchain go1.26.7`, exact Serf
  v0.10.4, and memberlist v0.6.0. Its `go.mod`/`go.sum` SHA-256 values were
  `88b9659119152f3e3d5e84a81be0685db27ee730c266a215743bc947cb4ef3a9` /
  `37dcf9ce84bcd16012189132f37ac839c329d68eb1eca532343a81832f55d768`.
  The repository does not invoke or otherwise need that CLI, so declaring it
  solely to select memberlist is manufactured tool ownership.
- The failed projection expanded the base to 261 modules and 3,826 graph
  edges, with 1,031 sum lines. It removed 51 edges and added 278. Production
  remained 355 entries, but complete-test entries fell from 429 to 402 and
  module-backed entries from 197 to 170 across the same 41 loaded modules.
  Serf, memberlist, and all earlier guarded packages remained absent from the
  product closures; only the 270-package Serf tool closure made their why
  results positive.
- Six guarded selections moved: errwrap v1.1.0, go-multierror v1.1.1,
  go-retryablehttp v0.7.7, go-sockaddr v1.0.7, golang-lru v1.0.2, and mdns
  v1.0.7. Guarded incoming edges grew from 164 to 195 through 31 additions and
  no removal, at candidate snapshot SHA-256
  `bef5bb840a3f67db60fd3fdc2460cbd1dd4e10c70bf7da1065aabe1ec82f9616`.
  The full selection delta also moved unrelated CLI, metrics, compression,
  serialization, Prometheus, crypt, protobuf, terminal, and x/* modules.
  This combines another dependency population with the remediation.
- The candidate failed the genuine-owner and focused-change gates. No owning
  decision was therefore granted for any candidate delta, and no Go-floor,
  stricter-format, API/CLI, behavior, or implementation-only integration was
  retained. The real `go.mod`/`go.sum`, source, all 23 guarded selections and
  164 guarded edges, graph counts, closure counts, and toolchain remain exact.
  No implementation commit exists.
- Fresh CNA, Go-index, OSV, release-route, import/load, and all earlier guard
  checks preserve the starting premise. Memberlist v0.3.0 remains affected,
  unqualified, and not risk-accepted. Exact Go 1.26.7 module verification,
  build, canonical count-one tests, race tests, and vet pass on the unchanged
  tree. No changed-selection scorecard applies and accepted quality remains
  27/27 Q0-Q2 PASS at L2.
- The integrated migration authorization is exhausted. P7 is blocked on one
  fresh documentation-only product decision: explicitly accept exact affected
  unloaded v0.3.0 under new target-specific guards; authorize a separately
  scoped POM subsystem redesign; or stop P7 unresolved. The decision may not
  manufacture a dependency owner, implement its chosen option, combine
  another group, or begin P8.

Hashicorp memberlist affected-version decision (2026-09-21):

- Option 1 is explicitly authorized: retain exact selected, inherited,
  unloaded `github.com/hashicorp/memberlist v0.3.0` without changing product
  source, `go.mod`, or `go.sum`. V0.3.0 is affected by HCSEC-2026-18 /
  CVE-2026-14362 and is not secure. Zero current import/load/runtime
  reachability bounds this explicit risk decision; it does not qualify the
  release or override the primary advisory.
- The decision is target-specific and non-transferable. It remains valid only
  while memberlist v0.3.0, selected Serf v0.10.1, both historical Serf-to-
  memberlist requests, all five incoming Serf requests and their parent
  identities, the exact graph and module hashes, no direct root/import/load/
  runtime reachability, every earlier guard, the exact PUBLISHED HashiCorp CNA
  identity and affected-below-v0.6.0 range, and the absence of a new advisory,
  independent defect, real supported tidy-stable fixed owner, or compatible
  fixed route all remain exact. Any difference immediately expires the
  decision and requires a fresh owning memberlist decision before merge.
- Guard revalidation began from clean HEAD `652baa8`, parent `7bb2629`, tree
  `464b3944`. Its exact five-file handoff, reciprocal 230-archive chain, latest
  Google UUID implementation ancestry, and launcher check passed. Exact Go
  1.26.7 retains 234 modules, 3,599 edges, 355/429 load entries, 197 module-
  backed complete-test entries across 41 modules, 1,067 sum lines, and the
  432-line tidy projection. `go.mod`/`go.sum` retain SHA-256 values
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- All 23 earlier guarded selections and 164 incoming edges retain snapshot
  SHA-256
  `f42304a43a7ade746c5c2bc238dfa9a6daec62594834fd5b606fdf245d44ea2e`;
  all 24 target-plus-earlier why results remain negative and guarded imports/
  loads remain zero. Including memberlist, 24 guarded selections and 166
  incoming edges yield snapshot SHA-256
  `675ec82de2ae67ea293a2229d0a81c49e8523839951d03703c3be788e013dc54`.
- Fresh CNA bytes remain PUBLISHED and byte-identical at SHA-256
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`,
  with memberlist below v0.6.0 affected and v0.6.0 first fixed. No security
  investigation was repeated. Direct roots, the unused Serf CLI tool owner,
  replacements, exclusions, forks, patches, version masquerading, POM
  redesign, Go-floor change, guarded-module change, and P8 remain unauthorized.
- P7 continues only with the reciprocal bounded evaluation of selected
  `github.com/iancoleman/strcase v0.2.0`. It exists through the sole recorded
  `envoyproxy/protoc-gen-validate v0.6.2` request, has negative why and zero
  repository import/load state, and must be evaluated independently without
  transferring the memberlist exception or combining another dependency group.
- Final exact-Go module verification, build, canonical count-one tests, race,
  and vet pass. The reciprocal 231-archive launcher check, all 62 launcher
  lifecycle controls, and all 15 audit meta-controls pass. No changed-
  selection scorecard applies; accepted quality remains 27/27 Q0-Q2 PASS at
  L2.

Iancoleman strcase evaluation (2026-09-21):

- No exact-path stable release qualifies. Fresh public evidence resolves six
  stable tags v0.1.0-v0.3.0 in the public unarchived non-fork MIT repository;
  v0.3.0 is latest and master. There is no later commit, prerelease, alternate
  major line, redirect, retraction, deprecation, or GitHub Release object.
- Selected v0.2.0 and latest v0.3.0 both declare Go 1.16, have no module
  dependencies, expose one library package and the same ten functions, and
  preserve Go 1.18 through their complete minimal source/test closures. Both
  pass upstream tests/repeats/race/vet and five cross-builds under exact Go
  1.26.7 and Go 1.18.10; apidiff reports no exported API change.
- Selected v0.2.0 fails the independent concurrent acronym contract with
  read/write and write/write races followed by `fatal error: concurrent map
  writes`. It also leaves all-uppercase input unchanged. Latest v0.3.0 fixes
  both with `sync.Map` and new uppercase handling but regresses
  `DBClusterParameterGroup` to `DbclusterParameterGroup`. Both releases'
  camel functions silently discard Unicode and malformed bytes. Upstream
  issue 40, issue 43, issue 39, and open PR 49 corroborate the findings.
- MVS selects v0.2.0 solely through protoc-gen-validate v0.6.2; parent and
  target why/import/load results are negative or zero. A disposable v0.3.0
  root moves only strcase and passes project verify/build/test/race/vet,
  API/CLI, and host acceptance, but it manufactures ownership, changes the
  memberlist graph/module-hash guard, and tidy removes it and restores v0.2.0.
  No source or dependency metadata changed.
- Fresh Go-index, OSV, GitHub, and govulncheck evidence has no strcase finding.
  All 24 guarded selections and 166 edges retain snapshot SHA-256
  `675ec82de2ae67ea293a2229d0a81c49e8523839951d03703c3be788e013dc54`;
  all 25 why results remain negative and guarded imports/loads stay zero. The
  memberlist CNA evidence and every earlier guard remain exact.
- Selected v0.2.0 is not qualified or accepted. P7 stops for one explicit
  product decision: target-specific guarded v0.2.0 retention, a separate
  measurement-only protoc-gen-validate owning-parent study, or stop unresolved.
  The decision may not implement a route, add a root, change a guard, combine
  another group, reopen memberlist/Serf/POM work, or begin P8.

Iancoleman strcase product decision (2026-09-21):

- Option 1 is explicitly authorized: retain exact selected, inherited and
  unloaded `github.com/iancoleman/strcase v0.2.0` without changing product
  source, `go.mod`, or `go.sum`. This is a target-specific, non-transferable
  exception, not release qualification. Zero current import/load/runtime
  reachability bounds the accepted exposure but did not itself authorize the
  decision.
- The exception accepts only the completed unsynchronized process-global
  acronym map; read/write and write/write races and fatal concurrent-map
  crash; unchanged all-uppercase input; Unicode/malformed-byte loss; permanent
  acronym state; and the characterized API, initialism, digit, delimiter,
  ignore, allocation, MVS, vulnerability, and related findings. It accepts no
  new advisory, independent defect, or uncharacterized behavior. Latest
  v0.3.0 remains unqualified because it retains Unicode loss and introduces
  the recorded mixed-initialism regression.
- The decision remains valid only while exact v0.2.0, the sole exact
  protoc-gen-validate v0.6.2 request, both historical Viper v1.10.1 and crypt
  v0.4.0 incoming parent requests, no direct root/import/load/runtime
  reachability, exact graph/module/tidy state, every earlier guard, and no new
  advisory, independently observed defect, qualified exact-path stable
  release, genuine supported tidy-stable owner, or compatible qualified route
  remain exact. Any difference expires the exception and requires a fresh
  owning strcase decision before merge.
- Guard revalidation began from clean HEAD `5071025`, parent `c0d76f1`, tree
  `f1625d3`. Its exact five-file handoff, reciprocal 232-archive chain, latest
  Google UUID implementation ancestry, exact Go 1.26.7 identity, and launcher
  check passed. The unchanged project retains 234 modules, 3,599 edges,
  355/429 load entries, 197 module-backed complete-test entries across 41
  modules, 1,067 sum lines, and the 432-line tidy projection. `go.mod`/`go.sum`
  remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- All 25 strcase-plus-earlier guarded why results are negative;
  protoc-gen-validate is also negative; repository imports are zero; and
  production/complete-test closures load zero target, parent, or guarded
  packages. All 24 earlier selections and 166 incoming edges retain snapshot
  SHA-256
  `675ec82de2ae67ea293a2229d0a81c49e8523839951d03703c3be788e013dc54`;
  including strcase yields 25 selections and 167 edges at
  `b6e0bcee17f83b0cc88a370c2f86a51c1872506b113075960c497bf420c48fc9`.
- Fresh narrow checks retain the same six versions and v0.3.0 master. Exact
  v0.2.0/v0.3.0 OSV, GitHub global, and repository advisory results remain
  empty. The Go index remains 1,402 records at SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
  The PUBLISHED memberlist CNA response remains 2,807 bytes at SHA-256
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`,
  with memberlist below v0.6.0 affected. No completed evaluation or parent
  study was repeated.
- Direct roots, replacements, forks, patches, exclusions, version
  masquerading, unused anchors, parent changes, POM redesign, Go-floor or
  guarded-module changes, and P8 remain unauthorized. No exception transfers.
  Its bounded demangle successor was later executed under the defensive scope
  below. Final exact-Go verify/build/test/race/vet, launcher lifecycle, and
  audit meta-controls pass; no changed-selection scorecard applies and
  accepted quality remains 27/27 Q0-Q2 PASS at L2.

Ianlancetaylor demangle evaluation (2026-09-21):

- No exact-path candidate qualifies. The active public unarchived non-fork
  BSD-3-Clause repository has one `master` branch, no Git tags, no GitHub
  Releases, no redirect or alternate major line, and no retraction or module
  deprecation. Selected
  `v0.0.0-20200824232613-28f6c0f3b639` is exact commit
  `28f6c0f3b63983aaa99575ca3b693afff7996387`; current master/proxy `@latest`
  is `v0.0.0-20260724033716-83e58baca724`, exact commit
  `83e58baca7248962d58657affe41b3b6f27ee423`. Selected is 97 commits behind
  current; proxy/Git files and sumdb identities agree.
- Selected source has no `go.mod`; the proxy synthesizes only its module line.
  Current declares Go 1.13. Both have standard-library-only closures and
  preserve Go 1.18 under the completed upstream native/count-ten/race/vet and
  ordinary cross-build checks on exact Go 1.26.7 and contained Go 1.18.10.
  Static apidiff is incompatible because current changes `NoClones` from 2 to
  3 and `Verbose` from 3 to 5; current also compatibly adds Rust demangling and
  new C++/Rust API forms and options.
- Ordinary upstream valid C++ behavior passes on both versions, and current's
  Rust, stdin, and LLVM-style examples pass. Both commands retain unchecked
  `doDemangle` indexing consistent with the already observed empty-argument
  panic. That panic was not reproduced. No malformed, empty, deeply nested,
  oversized, randomized, or adversarial input; fuzzing; stress/resource probe;
  crash reproduction; or security/exploitability analysis was performed.
- MVS selects the target only through thirteen historical Google pprof
  requests: five request the 2018 version and eight, including selected pprof
  `v0.0.0-20210720184732-4bb14d4b1be1`, request selected. Target why is
  negative and repository imports and production/complete-test loads are zero.
  All 25 earlier guarded selections and 167 incoming edges retain snapshot
  SHA-256
  `b6e0bcee17f83b0cc88a370c2f86a51c1872506b113075960c497bf420c48fc9`.
- A disposable current direct root changed only demangle, retained 234 modules
  and 355/429 project loads, added one graph edge and two checksum lines, and
  passed verify/build/count-one/count-ten/race/vet, API/CLI, CLI surface, and
  host acceptance. It manufactured ownership and changed earlier graph/module
  guards. Tidy removed it, restored selected, and returned exact baseline
  hashes
  `5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
  `b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
  No source or dependency metadata changed and no projection was applied.
- Fresh OSV, GitHub, Go-index, and isolated govulncheck evidence contains no
  demangle finding. Base and disposable-current project scans are identical at
  30/22/20/20 module/package/symbol/test-symbol findings and contain zero
  target trace. The Go index remains 1,402 records at SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`;
  guard OSV and memberlist CNA evidence remain exact.
- Selected remains physically present but is not qualified or accepted. P7
  stops for one documentation-only product decision: exact selected-version
  retention under a demangle-specific non-transferable exception, one
  separate measurement-only pprof owning-parent/request study, or an explicit
  unresolved stop. The decision may not reproduce parser failures, add a
  direct root, silently select current, implement a parent route, alter pprof
  or an earlier guard, combine another group, reopen excluded work, or begin
  P8. Every disposable used by the retry remained beneath the launcher's
  managed `${CODEX_SESSION_SCRATCH_ROOT:?}` and must not be retained.

Ianlancetaylor demangle product decision (2026-09-21):

- Option 1 is explicitly authorized: retain exact selected, inherited and
  unloaded `github.com/ianlancetaylor/demangle
  v0.0.0-20200824232613-28f6c0f3b639` without changing product source,
  `go.mod`, or `go.sum`. This is a demangle-specific, non-transferable
  exception, not release qualification. Zero current import/load/runtime
  reachability bounds the accepted exposure but did not itself authorize the
  decision.
- The exception accepts only the already observed empty-argument `c++filt`
  panic, absence of tagged/GitHub release identity, and completed C++ API,
  options, behavior, command, Go-floor, MVS, loading, vulnerability, and
  related findings. It accepts no new advisory, independent defect, or
  uncharacterized behavior. Current
  `v0.0.0-20260724033716-83e58baca724` remains unqualified and unselected
  because it is unreleased, retains the panic, and incompatibly changes the
  exported `NoClones` and `Verbose` values.
- The exception remains valid only while the exact selected pseudo-version,
  all thirteen historical pprof requests, selected pprof
  `v0.0.0-20210720184732-4bb14d4b1be1` and its selected request, no direct
  root/import/load/runtime reachability, exact graph/module/tidy state, every
  earlier guard, and no new advisory, independently observed defect, qualified
  exact-path release, genuine supported tidy-stable owner, or compatible
  qualified route remain exact. Any target, request, parent, pprof, root,
  import, load, runtime, graph, module-hash, tidy, earlier-guard, advisory,
  finding, release, owner, or route change expires the exception and requires
  its fresh owning decision before merge. No exception transfers.
- Guard-only revalidation began from clean HEAD `050c47a`, parent `de47934`,
  tree `fb1f0bd`. Its exact five-file handoff, reciprocal 234-archive chain,
  latest Google UUID implementation ancestry, exact Go 1.26.7 identity, and
  launcher check passed. The unchanged project retains 234 modules, 3,599
  edges, 355/429 load entries, 197 module-backed complete-test entries across
  41 modules, 1,067 sum lines, and the 432-line tidy projection. `go.mod` and
  `go.sum` retain SHA-256 values
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- All 26 demangle-plus-earlier guarded why results are negative; pprof is also
  negative; repository imports are zero; and production/complete-test closures
  load zero target, pprof, or guarded packages. All 25 earlier selections and
  167 edges retain snapshot SHA-256
  `b6e0bcee17f83b0cc88a370c2f86a51c1872506b113075960c497bf420c48fc9`;
  including demangle yields 26 guarded selections and 180 incoming edges at
  `8f39e82800e0abf94ae419186614a2cf709759192aba2614603dbf578ef64121`.
- Fresh exact selected/current demangle OSV, GitHub global, and repository
  advisory results remain empty. The Go index remains 1,402 records at
  SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`;
  guard OSV retains only the recorded Gorilla/go-retryablehttp pairs, and the
  PUBLISHED memberlist CNA response remains 2,807 bytes at SHA-256
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
  No parser, crash, parent-route, security, or exploitability work was
  repeated.
- P7 continues only with the reciprocal bounded evaluation of selected
  `github.com/jonboulle/clockwork v0.1.0`. Current queue observations are one
  historical `mvn-pom-mutator v0.2.3` request, negative why, zero repository
  imports, and zero production/complete-test loads. They are not qualification.
  The successor must preserve demangle, pprof, every earlier guard, the parent,
  and the Go floor; it may not combine another group or begin P8.
- Final exact-Go module verification, build, canonical count-one tests, race,
  vet, reciprocal 235-archive launcher check, and scratch-containment cleanup
  pass. No dependency implementation or changed-selection scorecard applies;
  accepted quality remains 27/27 Q0-Q2 PASS at L2.

Jonboulle clockwork evaluation (2026-09-21):

- No exact-path stable release qualifies. The active public unarchived
  non-fork Apache-2.0 repository exposes eight stable tags v0.1.0-v0.5.0 on
  one linear ancestry, no prerelease, retraction, deprecation, replacement,
  redirect, alternate path, or `/v2`/`/v3` module. Selected signed annotated
  v0.1.0 is commit `2eee05ed794112d45db504eb05aa693efd2b8b09`;
  highest Go-1.18-compatible signed v0.4.0 is commit
  `606c48b92358fcca153952b56fb0d14d6845f84a`; latest/current signed v0.5.0
  is commit `6d8d032a18422c2e3ef651170a8a55012d1f704c`. Proxy/Git regular
  files and sumdb identities agree.
- V0.1.0 has a synthesized module-only `go.mod`; v0.2.0-v0.3.0 declare Go
  1.13 and v0.4.0 declares Go 1.15. Their standard-library-only source/test
  closures preserve Go 1.18 and pass upstream native/repeat/race/vet checks
  under exact Go 1.26.7 and contained Go 1.18.10. Serious versions also pass
  linux/amd64, linux/arm64, windows/amd64, freebsd/amd64, and js/wasm cross-
  compilation. V0.5.0 declares Go 1.21 and imports `slices`, so Go 1.18.10
  cannot compile it.
- The target is one library package with no command, cgo, generated file,
  build tag, platform branch, external resource, testdata, benchmark, fuzz
  target, or actual example function. V0.1.0 exports only the basic real/fake
  clock surface; v0.2-v0.4 add ticker, context injection, timer/AfterFunc,
  cancellable blocking, and reset. V0.5.0 incompatibly changes `FakeClock`
  from interface to struct pointer and both constructor return types while
  adding Until and fake-clock deadline contexts.
- Selected v0.1.0 fails ordinary parity because fake `After` does not fire
  immediately for negative durations; exact-count BlockUntil can also miss a
  skipped waiter count. V0.2.0-v0.3.0 fake tickers do not reject non-positive
  periods. Highest compatible v0.4.0 fails four deterministic contracts: reset
  retains its old period, reset appends a duplicate schedule, Reset-then-Stop
  still fires, and NewTicker/Reset accept non-positive values. V0.5.0 retains
  the duplicate/Stop defect and Reset validation defect; its Context.Err
  blocks before Done, Cancel retains its deadline waiter, and parent deadline
  handling contradicts its documentation. Its count-ten upstream delivery-
  order test also has a scheduling-dependent expectation incompatible with
  documented one-slot ticker drops. Positive timer/ticker/AfterFunc and 64-
  waiter race fixtures otherwise pass.
- MVS selects exact v0.1.0 only through
  `mvn-pom-mutator@v0.2.3 -> clockwork@v0.1.0`. The POM parent is genuinely
  imported/loaded, but clockwork has no direct root, negative why, zero source
  imports, and zero production/complete-test loads. A disposable v0.4.0 root
  yields 234 modules/3,600 edges/1,069 sums; v0.5.0 yields
  234/3,602/1,069 and raises the main directive to Go 1.21. Both remain
  unloaded and pass the exact-Go project gate. Tidy removes either manufactured
  root and restores selected v0.1.0; v0.4.0 returns common baseline tidy hashes
  `5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
  `b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
- Fresh exact candidate OSV/GitHub feeds, Go-index lookup, and isolated
  govulncheck module/package/symbol/test-symbol scans are empty. Base and both
  projections have identical 30/22/20/20 project scan populations and no
  clockwork trace. The 1,402-record Go index, only recorded Gorilla and
  retryablehttp guard pairs, and PUBLISHED memberlist CNA bytes remain exact.
  All 26 earlier guarded selections/180 recorded edges, 27 negative target-
  plus-guarded why results, zero guarded imports/loads, project hashes, and
  earlier decisions remain unchanged.
- Product source and dependency metadata remain unchanged. No changed-
  selection gate applies and accepted quality remains 27/27 Q0-Q2 PASS at L2.
  The reciprocal documentation-only clockwork decision selected option 1 as
  recorded below. Direct roots, v0.4.0/v0.5.0 selection, a floor raise,
  parent/POM or earlier-guard changes, and P8 remain unauthorized.

Jonboulle clockwork product decision (2026-09-21):

- Option 1 is explicitly authorized. Retain exact selected, inherited,
  unloaded `github.com/jonboulle/clockwork v0.1.0` without changing product
  source, `go.mod`, or `go.sum`. This is a clockwork-specific,
  non-transferable exception; v0.1.0 remains unqualified, and zero current
  import/load/runtime reachability bounds exposure without itself qualifying
  or silently authorizing the selection.
- The exception accepts only the completed negative-duration fake-`After`
  behavior, exact-count `BlockUntil` boundary, characterized API, scheduling,
  blocking, concurrency, determinism, caller ownership, resource lifecycle,
  Go-floor, MVS, loading, repeatability, vulnerability, and related findings,
  plus the completed later-release ticker/context/API/Go-floor findings needed
  to establish that no stable exact-path release qualifies. It accepts no
  uncharacterized behavior, new advisory, or independently observed defect.
  V0.4.0 and v0.5.0 remain unqualified and unselected.
- Guard-only revalidation began from clean handoff HEAD `e5fd2d4`, parent
  `bb38c46`, tree `9713362`. Its exact five-file changed set, reciprocal
  236-archive chain, latest Google UUID implementation ancestry, exact
  official Go 1.26.7 archive/binary identities, and launcher check passed.
  The project remains 234 modules, 3,599 graph edges, 355/429 load entries,
  197 module-backed complete-test entries across 41 modules, 1,067 sum lines,
  and the exact 432-line tidy projection. `go.mod`/`go.sum` remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- MVS retains exactly the sole
  `mvn-pom-mutator@v0.2.3 -> clockwork@v0.1.0` request. The parent remains
  genuinely imported and its `pkg/pom` package loads; clockwork has no direct
  root, negative why, zero repository imports, and zero production or
  complete-test loads. All 27 clockwork-plus-earlier guarded why results are
  negative. All 26 earlier selections and 180 incoming edges retain snapshot
  SHA-256
  `8f39e82800e0abf94ae419186614a2cf709759192aba2614603dbf578ef64121`;
  including clockwork yields 27 selections and 181 incoming edges at
  `09258fe1e7425eed2d927f8f0be20a60fa610098677e2f31a10b5b4153b35989`.
- Fresh exact clockwork candidate OSV, GitHub global, and repository advisory
  results remain empty. The 1,402-record Go index remains byte-identical at
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`;
  guard OSV retains only the recorded Gorilla/go-retryablehttp pairs, and the
  2,807-byte PUBLISHED memberlist CNA response remains exact at
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
  No behavior failure, parent route, security, exploitability, fuzz, stress,
  or resource-exhaustion work was repeated.
- The exception remains valid only while exact v0.1.0; the sole exact parent
  request and genuine parent ownership; no direct root/import/load/runtime
  reachability; negative target why; the exact graph, module hashes, tidy
  projection, and every earlier guard; and no new advisory, independent
  defect, qualified stable release, genuine supported tidy-stable owner, or
  compatible qualified route remain exact. Any target, request, parent, root,
  import, load, runtime, graph, module-hash, tidy, earlier-guard, advisory,
  finding, release, owner, or route change expires the exception and requires
  a fresh clockwork dependency and product decision before merge. It
  authorizes no parent/request study or workaround and transfers to no target.
- P7 continues only with the reciprocal bounded evaluation of selected
  `github.com/json-iterator/go v1.1.12`. That successor must preserve
  clockwork and every earlier guard, evaluate only json-iterator, and stop for
  a fresh owning decision if no exact-path stable release qualifies or any
  guard changes. It authorizes neither a combined dependency group nor P8.

Json-iterator Go evaluation (2026-09-21):

- No exact-path stable release qualifies. Canonical MIT repository
  `json-iterator/go` is a non-fork and now archived. The proxy exposes exactly
  eight valid stable releases v1.1.5-v1.1.12 on one linear ancestry, no v2/v3
  module, retraction, deprecation, replacement, or alternate exact path. Tags
  are lightweight/unsigned. Selected/latest v1.1.12 is commit
  `024077e996b048517130b21ea6bf12aa23055d3d`, tree
  `fab0aaa4437b102ed0e2e1f6db945807c4e68318`; its 139-file proxy archive is
  byte-identical to Git and sumdb identities agree. Master is an unreleased
  seven-commit-later pseudo-version and is ineligible for stable promotion.
- All stable minimal source/test closures preserve the Go 1.18 directive floor
  and resolve under exact Go 1.26.7 and contained Go 1.18.10. Compile-only
  tests pass both SDKs for v1.1.6-v1.1.12; v1.1.5 has an exact-Go misnamed-
  example compile failure. Selected production packages cross-build under both
  SDKs for Linux, Windows, and Darwin amd64. There is no command, cgo,
  generated source, embed, testdata, symlink, or external resource.
- The randomized `google/gofuzz`-based `type_tests` identified by the initial
  broad run was excluded from qualification and not executed again. The nine
  admissible packages pass count-one/count-ten/race under Go 1.18.10. Exact
  Go fails those populations only at deterministic control-escape equality;
  vet under both SDKs reports the same three duplicate tags in test fixtures.
- V1.1.5-v1.1.6 accept non-finite floats without error. V1.1.7-v1.1.9 fail
  exact-Go byte compatibility because standard `encoding/json` emits
  `"\\b\\f"` while json-iterator emits `"\\u0008\\u000c"`. V1.1.10-v1.1.12
  additionally encode an ordinary TextMarshaler string-alias map key as
  `{"TEXT_key":"value"}` rather than standard `{"key":"value"}` under
  both SDKs. Master fixes the map key but remains unreleased and retains the
  escape failure. Both escape forms are valid JSON, but the advertised and
  upstream-tested byte-compatibility contract fails.
- Bounded ordinary fixtures otherwise pass the relevant struct, map,
  interface, nil, raw/custom codec, number, validation, decoder/encoder,
  iterator, and stream paths. Static review records documented token/indent
  compatibility gaps, config-specific determinism/precision, borrowed buffer
  lifetimes, unsafe reflection, retained pool buffers, and mutable unprotected
  global codec/extension registration that must complete before concurrent
  cached use. There is no closeable external resource.
- Seven historical Viper, crypt, etcd client/v2, and Prometheus requests make
  v1.1.12 the MVS maximum. Viper v1.15.0 is genuine/direct/imported/loaded;
  json-iterator has no direct root, negative why, zero repository imports,
  zero production/complete-test loads, and no runtime reachability. A direct
  v1.1.12 projection yields 234 modules/3,602 edges/1,070 sums; master yields
  234/3,607/1,071. Both remain unloaded. Tidy removes either manufactured root,
  restores v1.1.12, and returns common baseline hashes
  `5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
  `b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
- Fresh candidate OSV/GitHub/Go-index and isolated pinned govulncheck evidence
  is empty. Base/direct-selected/master project streams remain identical at
  30/22/20/20 with no target trace. The 1,402-record Go index, only recorded
  Gorilla/retryablehttp guard pairs, PUBLISHED memberlist CNA bytes, all 27
  guarded selections/181 edges, negative guarded why/import/load results,
  project hashes/state, and prior decisions remain exact.
- Product source and dependency metadata remain unchanged. Final exact-Go
  module verification, build, count-one, race, vet, reciprocal launcher check,
  and scratch cleanup pass. No changed-selection scorecard applies; accepted
  quality remains 27/27 Q0-Q2 PASS at L2. P7 stops for one reciprocal bounded
  product decision: retain exact selected/unloaded v1.1.12 under a new target-
  specific exception, authorize one measurement-only owning-parent/request
  study, or stop P7 unresolved. It may not alter a parent/guard, select master
  or a lower release, add a direct root, transfer an exception, or begin P8.

Json-iterator Go product decision (2026-09-21):

- Selected option 1. Retain exact selected, inherited, unloaded
  `github.com/json-iterator/go v1.1.12` without changing product source,
  `go.mod`, or `go.sum`. The selection is not qualified. Zero loading bounds
  current exposure but did not qualify or silently authorize it.
- The json-iterator-specific, non-transferable exception accepts only the
  completed map-key, control-escape, older non-finite-number, example/vet,
  repository/archive, API/global-state/buffer, Go-floor, MVS, loading,
  repeatability, vulnerability, and related evaluation findings. It accepts no
  uncharacterized behavior, new advisory, or independently discovered defect;
  it does not promote master, select a lower release, or transfer another
  exception.
- The exception is owned by exact v1.1.12 and all seven historical requests
  with their parent identities: bketelsen crypt requests v1.1.6, Prometheus
  v1.0.0 requests v1.1.6, Prometheus v1.4.0 requests v1.1.9, sagikazarmark
  crypt and Viper v1.10.1/v1.15.0 request v1.1.12, and etcd client/v2
  v2.305.1 requests v1.1.11. Viper v1.15.0 remains the genuine direct,
  imported, and loaded project parent without loading a json-iterator package.
- Its guards require no direct root, negative target why, zero repository
  imports, zero production/complete-test loads, no runtime reachability, the
  exact 234-module/3,599-edge/1,067-sum and 355/429/197/41 load state, the
  exact real module hashes and 432-line tidy projection, every earlier guard,
  and no new advisory, independent defect, qualified exact-path stable release,
  genuine supported tidy-stable owner, or compatible qualified route. Any
  target, request, parent, ownership, root, import, load, runtime, graph,
  module-hash, tidy, earlier-guard, advisory, finding, release, owner, or route
  change expires the exception and requires its fresh owning decision before
  merge.
- Guard-only revalidation from clean handoff HEAD `19c5002`, parent `80ea158`,
  tree `9cec063`, preserved the exact five-file handoff, reciprocal archive
  chain, Google UUID implementation ancestry, exact Go 1.26.7 identity, module
  hashes, seven requests, and launcher check. All 28 target-plus-earlier why
  results are negative and guarded imports/loads are zero. The earlier 27-
  selection/181-edge snapshot remains exact at
  `09258fe1e7425eed2d927f8f0be20a60fa610098677e2f31a10b5b4153b35989`;
  including json-iterator yields 28 selections/188 edges at
  `ca6a8b9d4a0d4aa5c6c436d9edeb27cf6d603e36e29edccef9bc2d9d1b81622e`.
- Fresh selected/master OSV and GitHub global/repository results remain empty.
  The 1,402-record Go index remains byte-identical at
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`;
  guard OSV retains only the recorded Gorilla/retryablehttp pairs, and the
  2,807-byte PUBLISHED memberlist CNA response remains exact at
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
  No behavior result was repeated and no option-2 study was run.
- P7 continues only with the reciprocal bounded evaluation of selected exact-
  path `github.com/jstemmer/go-junit-report v0.9.1`. It must preserve the
  json-iterator exception and every earlier guard, evaluate no other group, and
  stop for a fresh owning decision if no exact-path stable release qualifies or
  any guard changes. P8 remains queued.

Jstemmer Go JUnit Report evaluation (2026-09-21):

- No exact-path stable release qualifies. Fresh `go-import` resolves without
  redirect to the public, unarchived, non-fork MIT repository. The exact-path
  proxy exposes exactly stable v0.9.0, v0.9.1, and v1.0.0 on one ancestry with
  annotated tags, no retraction/deprecation/replacement/alternate exact path,
  byte-identical proxy/Git files, and matching sumdb identities. Latest stable
  v1.0.0 is commit `16c7efad77dbe9cb51f1a8c896c652c8ecdaa20e`, tree
  `2cd81297935f0f139b09b32a7f47f60429208986`; its tag signature is verified.
  The later v1 pseudo changes only tests/workflow. Current master declares the
  different `/v2` module path, whose releases were not promoted.
- Every stable candidate is a one-module, standard-library-only closure that
  resolves and passes upstream count-one/count-ten, race, vet, build, and five-
  platform test cross-compilation under exact Go 1.26.7 and contained Go
  1.18.10. Their exported parser/formatter APIs are identical. Bounded command
  and library fixtures pass ordinary report/XML, flags, stdin/stdout, exit,
  reader error, determinism, ownership, and concurrent-render behavior.
- Every stable release fails the same exported output contract:
  `formatter.JUnitReportXML` ignores errors from buffered writes and `Flush`
  and returns nil. A small deterministic destination returning an ordinary
  write error therefore produces nil under both SDKs for all three stable
  releases. The unreleased v1 branch has no production change and retains the
  defect. This ordinary error/resource-lifecycle failure disqualifies every
  stable exact-path candidate.
- MVS selects v0.9.1 through 25 historical Cloud Go/storage requests. It has
  no direct root, negative why, zero source imports, zero production/complete-
  test loads, and no runtime reachability. Direct v0.9.1 and v1.0.0 projections
  yield 234 modules/3,600 edges and 1,068/1,069 sum lines respectively; v1.0.0
  changes only the target selection, remains unloaded, preserves every earlier
  guard, and passes project/API/CLI/acceptance gates. Tidy removes either
  manufactured root and restores exact baseline v0.9.1 and projection hashes.
  No genuine tidy-stable owner was identified; no projection was retained.
- Fresh candidate OSV/GitHub/Go-index and isolated pinned govulncheck evidence
  is empty. Base/direct-v1.0.0 project streams remain byte-identical at
  30/22/20/20 with no target trace. The 1,402-record Go index, only recorded
  Gorilla/retryablehttp guard pairs, PUBLISHED memberlist CNA bytes, all 28
  guarded selections/188 edges, negative guarded why/import/load results,
  project hashes/state, and every earlier decision remain exact.
- Product source and dependency metadata remain unchanged. Final exact-Go
  module verification, build, count-one, race, vet, reciprocal launcher check,
  and scratch cleanup pass. No changed-selection scorecard applies; accepted
  quality remains 27/27 Q0-Q2 PASS at L2. P7 stops for one reciprocal bounded
  product decision: retain exact selected/unloaded v0.9.1 under a new target-
  specific exception, authorize one measurement-only Cloud Go/storage owning-
  parent/request study, or stop P7 unresolved. It may not add a direct root,
  promote a branch/pseudo-version, select `/v2`, transfer an exception, combine
  another dependency group, or begin P8.

Jstemmer Go JUnit Report product decision (2026-09-21):

- Selected option 1. Retain exact selected, inherited, unloaded
  `github.com/jstemmer/go-junit-report v0.9.1` without changing product source,
  `go.mod`, or `go.sum`. The selection is not qualified. Zero loading bounds
  current exposure but did not qualify or silently authorize it.
- The go-junit-report-specific, non-transferable exception accepts only the
  completed ordinary destination-writer error, repository/archive and exact-
  path release identity, exported parser/formatter API, command surface,
  bounded ordinary behavior, deterministic rendering, caller ownership,
  output resource-lifecycle, Go-floor, MVS, loading, repeatability,
  vulnerability, and related findings. It accepts no uncharacterized behavior,
  new advisory, or independently discovered defect; it does not promote the
  unreleased v1 branch/pseudo-version, select `/v2`, or transfer another
  exception.
- The exception is owned by exact v0.9.1 and all 25 historical requests with
  their genuine parent identities: Cloud Go v0.38.0, v0.44.1-v0.44.3,
  v0.45.1, v0.46.3, and v0.50.0 request the older pseudo-version; Cloud Go
  v0.52.0, v0.53.0, v0.54.0, v0.56.0, v0.57.0, v0.62.0, v0.65.0, v0.72.0,
  v0.74.0, v0.75.0, v0.78.0, v0.79.0, v0.81.0, v0.83.0, v0.84.0, v0.87.0,
  and v0.90.0 plus Cloud Storage v1.5.0 request v0.9.1.
- Its guards require no direct root, negative target why, zero repository
  imports, zero production/complete-test loads, no runtime reachability, the
  exact 234-module/3,599-edge/1,067-sum and 355/429/197/41 load state, the
  exact real module hashes and 432-line tidy projection, every earlier guard,
  and no new advisory, independent defect, qualified exact-path stable release,
  genuine supported tidy-stable owner, or compatible qualified route. Any
  target, request, parent, ownership, root, import, load, runtime, graph,
  module-hash, tidy, earlier-guard, advisory, finding, release, owner, or route
  change expires the exception and requires its fresh owning decision before
  merge.
- Guard-only revalidation from clean handoff HEAD `efbe471`, parent `5946d09`,
  tree `ed26c49`, preserved the exact five-file handoff, reciprocal archive
  chain, Google UUID implementation ancestry, exact Go 1.26.7 identity, module
  hashes, all 25 requests, and launcher check. All 29 target-plus-earlier why
  results are negative and guarded imports/loads are zero. The earlier 28-
  selection/188-edge snapshot remains exact at
  `ca6a8b9d4a0d4aa5c6c436d9edeb27cf6d603e36e29edccef9bc2d9d1b81622e`;
  including go-junit-report yields 29 selections/213 edges at
  `5d2a35a2961c04eb07dd927afc072e77f555d87c7317579ddda486bf128d2c93`.
- Fresh stable/v1-branch OSV and GitHub global/repository results remain empty.
  The 1,402-record Go index remains byte-identical at
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`;
  guard OSV retains only the recorded Gorilla/retryablehttp pairs, and the
  2,807-byte PUBLISHED memberlist CNA response remains exact at
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
  No behavior result was repeated and no option-2 study was run.
- P7 continues only with the reciprocal bounded evaluation of selected exact-
  path `github.com/jtolds/gls v4.20.0+incompatible`. Its current queue facts
  are one historical goconvey v1.6.4 request, negative why, zero repository
  imports, zero production/complete-test loads, and no runtime reachability;
  none is qualification. The evaluation must preserve the go-junit-report
  exception and every earlier guard, evaluate no other group, and stop for a
  fresh owning decision if no exact-path stable release qualifies or any guard
  changes. P8 remains queued.

Jtolds GLS evaluation (2026-09-21):

- No exact-path stable release qualifies. Legacy unsuffixed exact path
  `github.com/jtolds/gls` exposes only v4.2.0+incompatible,
  v4.2.1+incompatible, and selected/latest v4.20.0+incompatible. The old
  GitHub URL and `go-import` metadata redirect to canonical `jtolio/gls`, but
  that alternate path was not promoted. There is no `/v4` module, later
  stable release, prerelease, retraction, deprecation, replacement, or later
  master commit. Tags are lightweight/unsigned; selected v4.20.0 is commit
  `b4936e06046bbecbb94cae9c18127ebe510a2cb9`, tree
  `21801a3091b3d1c1822c0621039366472633fafd`.
- `+incompatible` records v4 tags on a legacy unsuffixed source tree with no
  `go.mod`; the proxy synthesizes only `module github.com/jtolds/gls`.
  Native closure is target plus standard library. Complete build-tag tidy
  resolves the omitted JS import to GopherJS v1.21.0 across 18 modules; that
  dependency declares Go 1.21, so release metadata does not preserve the Go
  1.18 floor even though exact Go 1.18.10 compiles the js/wasm test binary.
- V4.2.0 and v4.2.1 compile, vet, and cross-compile but fail their upstream
  tests/examples under both exact Go 1.26.7 and contained Go 1.18.10 because
  stack tags are optimized away. V4.20.0 adds `//go:noinline` marker guards
  and passes upstream count-one/count-ten/race/vet under both SDKs. All three
  cross-compile their test binaries for five native targets under both SDKs.
- A bounded ordinary selected-release fixture passes value nesting/restoration,
  panic cleanup, ordinary four-goroutine propagation/isolation, and unregister
  behavior. Nested `SetValues` plus `gls.Go` deterministically returns the
  restored outer value instead of the documented point-in-time inner value
  because `getValues` aliases a live map until the child runs. Both SDKs
  reproduce this at count-one/count-ten, and race count-ten reports the
  parent restoration/child copy map race. This ordinary API/aliasing/
  concurrency failure disqualifies v4.20.0.
- MVS selects v4.20.0+incompatible only through goconvey v1.6.4, reached from
  direct mvn-pom-mutator v0.2.3. GLS has no direct root/import/load/runtime
  reachability and negative why. A disposable direct root adds only one graph
  edge and one sum line, remains unloaded, passes project verify/build/tests/
  race/vet, API/CLI/surface and host acceptance, and is removed by tidy back
  to the exact baseline. No root or projection was retained.
- Fresh exact/package OSV, GitHub global/repository, Go-index, and pinned
  govulncheck module/package/symbol/test-symbol evidence is empty for every
  candidate. Base/direct project scans remain identical at 30/22/20/20 with
  no target trace. The 1,402-record index hash, only Gorilla/retryablehttp
  guard pairs, PUBLISHED memberlist CNA bytes, all 29 earlier selections/213
  edges, all 30 negative why results, zero guarded imports/loads, project
  hashes/state, and every earlier decision remain exact.
- Product source and dependency metadata remain unchanged. Final exact-Go
  verify/build/count-one/race/vet, reciprocal launcher check, and scratch
  cleanup pass. No changed-selection scorecard applies; accepted quality
  remains 27/27 Q0-Q2 PASS at L2.
- P7 stops for one reciprocal bounded product decision: explicitly retain
  exact selected/unloaded v4.20.0+incompatible under a new target-specific
  exception, authorize one measurement-only goconvey owning-parent/request
  study, or stop P7 unresolved. It may not repeat the completed fixture, alter
  goconvey/mvn-pom-mutator, select a lower release or alternate path, add a
  direct root, transfer an exception, combine another group, or begin P8.

Jtolds GLS product decision (2026-09-21):

- Option 1 was selected. Exact selected, inherited, unloaded
  `github.com/jtolds/gls v4.20.0+incompatible` is explicitly retained without
  changing product source, `go.mod`, or `go.sum`. It remains unqualified; zero
  loading bounds exposure but is neither qualification nor implicit
  authorization.
- The jtolds/gls-specific, non-transferable exception accepts only the
  completed repository/release/legacy-module, Go-floor/JS-closure,
  older-release upstream incompatibility, selected-release nested
  propagation/live-map alias and race, API/global-state/lifecycle,
  MVS/loading, repeatability, project, vulnerability, and related findings.
  It accepts no uncharacterized behavior, new advisory, or independent defect
  and promotes no lower release, alternate path, fork, branch, pseudo-version,
  or replacement.
- The exception is bounded by exact path/version; the sole
  `goconvey@v1.6.4 -> jtolds/gls@v4.20.0+incompatible` request and genuine
  direct/imported/loaded mvn-pom-mutator v0.2.3 ownership route; negative GLS
  and GoConvey why results; no direct target root/import/load/runtime
  reachability; exact graph/module/tidy/Go-floor and earlier guards; and no new
  advisory, independent finding, release, owner, supported tidy-stable owner,
  or compatible qualified route. Any target version/path, request, ownership,
  root, import, load, runtime, graph, module-hash, tidy, Go-floor,
  earlier-guard, advisory, finding, release, owner, or route change expires
  the exception and requires a fresh jtolds/gls dependency and product
  decision before merge.
- The decision authorizes no parent study, workaround, direct root, alternate
  path, product-source or dependency-metadata change, unrelated selection, or
  implementation. It transfers neither this nor another dependency exception.
- Guard-only revalidation preserved the exact 234-module/3,599-edge/
  355-production/429-complete-test/197-module-backed/41-loaded-module state,
  1,067 sum lines, real module hashes, 432-line tidy projection, exact Go
  1.26.7 identity, sole GLS request/route, and earlier 29-selection/213-edge
  snapshot SHA-256
  `5d2a35a2961c04eb07dd927afc072e77f555d87c7317579ddda486bf128d2c93`.
  All 30 why results and guarded imports/loads remain zero. Including GLS gives
  30 selections/214 incoming edges at sorted snapshot SHA-256
  `22407933faccf68294140b09ff535f552685bfed20e7e806b6e74d552c275014`.
- Fresh exact/package OSV, GitHub, and pinned govulncheck module/package/
  symbol/test-symbol results remain empty. Base project findings remain
  30/22/20/20 without a GLS trace. The exact Go-index, only recorded
  Gorilla/retryablehttp guard pairs, and PUBLISHED memberlist CNA response
  identities remain unchanged. The completed behavior/race fixture was not
  repeated and no option-2 parent study ran.
- Final unchanged-project exact-Go verification/build/count-one/race/vet,
  reciprocal launcher/archive validation, diff checks, and scratch cleanup
  pass. No changed-selection scorecard applies; accepted quality remains
  27/27 Q0-Q2 PASS at L2.
- P7 continues only with the reciprocal bounded evaluation of selected exact-
  path `github.com/julienschmidt/httprouter v1.2.0`. It must independently
  resolve the target's exact ownership and loading, preserve the jtolds/gls
  exception and every earlier guard, evaluate no other group, and stop for a
  fresh product decision if no exact-path stable release qualifies or any
  guard changes. P8 remains queued.

Julienschmidt Httprouter evaluation and blocked ownership boundary
(2026-09-21):

- Stable exact-path v1.3.0 is the highest and only qualified release. Fresh
  exact-path identity resolves the public active unarchived non-fork
  BSD-3-Clause repository and exactly four stable releases v1.0.0-v1.3.0.
  Their lightweight unsigned tags form one ancestry; proxy/Git/sumdb identity
  agrees. There is no retraction, deprecation, replacement, `/v2` module,
  later stable, or prerelease. Alternate `gopkg.in`, branch, and unreleased
  master identities were not promoted.
- Every release has one package and a standard-library-only complete closure.
  V1.0.0-v1.2.0 use proxy-synthesized module files; v1.3.0 declares the exact
  path and `go 1.7`. Exact Go 1.26.7 and Go 1.18.10 resolve all build-tag
  branches and cross-compile test binaries for six targets. There are no
  commands, external module requirements, examples, fuzz targets, testdata,
  generated files, cgo, embeds, symlinks, subprocesses, or production
  external boundaries.
- V1.0.0/v1.1.0 fail obsolete upstream not-found expectations for ordinary
  default redirects. Selected v1.2.0 fails its upstream multibyte Unicode
  case-insensitive lookup tests under both SDKs at count one and count ten and
  is unqualified. V1.3.0 fixes that behavior, empty registration paths,
  wildcard conflicts, parameter counts, and deterministic allowed methods;
  its upstream count-one/repeats/race/vet pass under both SDKs.
- V1.3.0 is exported-API compatible with v1.2.0 and adds only
  `Router.GlobalOPTIONS`. A bounded ordinary fixture passes registration,
  matching, params, lookup/context ownership, redirects and Unicode fixed
  paths, deterministic OPTIONS/405, hooks, in-memory file serving, cleaning,
  documented panics, and concurrent read-only serving under both SDKs.
  Callers must complete unsynchronized router/configuration mutation before
  concurrent serving.
- Baseline MVS selects v1.2.0 only from Prometheus common v0.9.1 and historical
  v0.4.1 requests. The genuine shortest route begins at direct/imported/loaded
  mvn-pom-mutator v0.2.3 through Viper v1.10.1, go-metrics v0.3.10, and common
  v0.9.1. Target why/import/production/test/runtime loads are negative or
  zero.
- A disposable exact v1.3.0 get necessarily manufactures a main-module
  indirect root. It changes only target selection, yields 234 modules/3,600
  edges/1,069 sums while preserving 355/429/197/41 loads and every guard, and
  passes exact-Go verify/build/repeats/race/vet, API/CLI/surface, five cross-
  builds, and all four host acceptance stages. Tidy deletes the root and new
  sums, reselects v1.2.0, and converges byte-for-byte with the common baseline
  projection. No genuine supported tidy-stable v1.3.0 owner was found, so the
  mission's direct-root prohibition prevents implementation. No projection
  was retained and no dependency commit exists.
- Fresh exact/package OSV, GitHub global/repository, Go-index, and pinned
  govulncheck evidence is empty for every release. Base and v1.3.0 project
  populations remain identical at 30/22/20/20 without a target trace. The
  1,402-record index hash, only recorded Gorilla/retryablehttp guard pairs,
  PUBLISHED memberlist CNA bytes, all 30 earlier selections/214 edges,
  negative guarded why/import/load results, project hashes/state, and every
  earlier decision remain exact.
- Product source and dependency metadata remain unchanged. P7 stops for one
  reciprocal bounded product decision: explicitly retain exact inherited,
  unloaded, unqualified v1.2.0 under a target-specific exception; authorize
  one measurement-only owning-parent/request study; or stop P7 unresolved.
  It may not repeat the completed fixtures, add a direct target root,
  implement a parent change, transfer an exception, combine another group, or
  begin P8.
- Final unchanged-project exact-Go verify/build/count-one/race/vet,
  reciprocal 244-archive launcher validation, diff checks, and contained
  scratch cleanup pass. No changed-selection scorecard applies; accepted
  quality remains 27/27 Q0-Q2 PASS at L2.

Julienschmidt Httprouter product decision (2026-09-21):

- Option 1 was selected. Exact selected, inherited, unloaded
  `github.com/julienschmidt/httprouter v1.2.0` is explicitly retained without
  changing product source, `go.mod`, or `go.sum`. It remains unqualified; zero
  loading bounds exposure but is neither qualification nor implicit
  authorization. Stable v1.3.0 remains qualified but has no genuine supported
  tidy-stable project owner.
- The httprouter-specific, non-transferable exception accepts only the
  completed v1.2.0 Unicode lookup failure, obsolete v1.0.0/v1.1.0 redirect-
  test expectations, repository/archive/release/module/Go-floor, API/behavior,
  ownership/concurrency, MVS/loading, repeatability, project, vulnerability,
  qualified-v1.3.0, and related findings. It accepts no uncharacterized
  behavior, new advisory, or independently discovered defect and promotes no
  alternate path, branch, master, pseudo-version, replacement, or direct root.
- The exception is owned by exact v1.2.0; both
  `prometheus/common@v0.9.1 -> httprouter@v1.2.0` and
  `prometheus/common@v0.4.1 -> httprouter@v1.2.0` requests and parent
  identities; the genuine direct/imported/loaded mvn-pom-mutator v0.2.3 ->
  Viper v1.10.1 -> go-metrics v0.3.10 route; and the historical
  client_golang v1.0.0 request.
- Its guards require no direct target root, negative target why, zero
  repository imports, zero production/complete-test loads, no runtime
  reachability, the exact 234-module/3,599-edge/1,067-sum and
  355/429/197/41 load state, the exact real module hashes and 432-line tidy
  projection, the declared Go 1.18 floor and exact Go 1.26.7 identity, every
  earlier guard, and no new advisory, independent finding, release/owner
  change, newly disqualifying v1.3.0 finding, genuine supported tidy-stable
  owner, or compatible genuine route to qualified v1.3.0.
- Any target, request, parent, ownership, root, import, load, runtime, graph,
  module-hash, tidy, Go-floor, earlier-guard, advisory, finding, release,
  owner, or route change expires the exception and requires its fresh owning
  dependency and product decision before merge. The decision authorizes no
  parent study, workaround, direct root, parent change, alternate path,
  product-source/metadata change, unrelated selection, or implementation and
  transfers no exception.
- Guard-only revalidation from clean handoff HEAD `2e85c1e`, parent `cffd78b`,
  tree `0a2a83f`, preserved the exact five-file handoff, reciprocal archive
  chain, Google UUID implementation ancestry, exact Go 1.26.7 identity,
  module hashes, both requests and owner route, and launcher check. All 31
  target-plus-earlier why results and guarded imports/loads are zero. The
  earlier 30-selection/214-edge snapshot remains exact at
  `22407933faccf68294140b09ff535f552685bfed20e7e806b6e74d552c275014`;
  including httprouter yields 31 selections/216 edges at
  `a586cb138c690f082f0e2298fbc6872a06a11dd1e0301c38e86ea1bfd12944c2`.
- Fresh proxy metadata still exposes exactly v1.0.0-v1.3.0, and the GitHub
  repository remains public, active, unarchived, and non-fork. Exact OSV for
  all four releases and GitHub global/repository target results remain empty.
  OSV across all guards retains only the recorded Gorilla/retryablehttp pairs;
  the 1,402-record Go-index and PUBLISHED memberlist CNA response remain byte-
  exact. The completed behavior fixtures were not repeated and no option-2
  owning-parent study ran.
- Final unchanged-project exact-Go verify/build/count-one/race/vet, reciprocal
  245-archive launcher validation, diff checks, and contained scratch cleanup
  pass. No changed-selection scorecard applies; accepted quality remains
  27/27 Q0-Q2 PASS at L2.
- P7 continues only with the reciprocal bounded evaluation of selected exact-
  path `github.com/kisielk/errcheck v1.5.0`. Its current queue facts are one
  historical gogo/protobuf v1.3.2 request, negative why, zero repository
  imports, zero production/complete-test loads, and no runtime reachability;
  none is qualification. The evaluation must preserve the httprouter exception
  and every earlier guard, evaluate no other group, and stop for a fresh owning
  decision if no exact-path stable release qualifies or any guard changes. P8
  remains queued.

Kisielk Errcheck evaluation and product-decision boundary (2026-09-21):

- No exact-path stable release qualifies. Fresh identity evidence resolves
  the public active unarchived non-fork MIT repository, 16 stable releases
  v1.0.0-v1.20.0 plus excluded v1.5.0-alpha, one tag ancestry, matching
  proxy/Git/sumdb content, and no redirect, retraction, deprecation,
  replacement, `/v2` module, promoted fork, or later stable.
- V1.8.0 is the highest stable release preserving Go 1.18. Selected v1.5.0
  resolves a ten-module closure with directives at or below Go 1.14; v1.8.0
  resolves six modules at or below Go 1.18. V1.9.0/v1.10.0 require Go 1.22.0
  and v1.20.0 requires Go 1.25.0.
- Selected v1.5.0 fails its Go 1.18.10 library test in the old go/packages
  loader, fails ordinary exact-Go loading/analysis, and statically leaks every
  diagnostic source file opened by `readfile`; v1.6.1 adds the missing close.
  V1.6.1-v1.8.0 pass Go1.18 count one but fail exact-Go loading or x/tools
  compilation. V1.9.0+ pass exact-Go count one but are floor-ineligible.
- Every v1.6.1-and-later serious candidate fails count-ten/race-count-ten
  because upstream tests retain a `t.TempDir()` loader closure in global state
  and later iterations use the deleted directory. Ordinary v1.8/Go1.18 and
  v1.9/exact-Go command fixtures otherwise agree on deterministic checking,
  filtering, diagnostics, and exit behavior. Serious commands cross-build for
  six ordinary targets on compatible SDKs.
- MVS selects exact v1.5.0 only through gogo/protobuf v1.3.2 on the genuine
  route from direct/imported/loaded Viper v1.15.0. Target why/import/
  production-load/complete-test-load/runtime results are negative or zero;
  gogo/protobuf also loads no package.
- A disposable selected root manufactures three roots and tidy restores the
  common baseline. V1.8.0 moves six unrelated x/* selections and does not
  tidy-converge; v1.9.0 raises the project Go line. No projection or dependency
  commit was retained. The project remains
  234/3,599/355/429/197/41/1,067 with exact hashes and tidy projection. All 32
  why/import/load guards are negative/zero; 217 incoming edges hash to
  `500c57a1b4ffbdf5bb1296ef4dcb1bd16dbae12e54b86664fd1f8a2481a6c0dd`.
- Exact/package OSV and GitHub target results are empty. Pinned govulncheck has
  no errcheck package/symbol trace, but serious candidates' x/mod closures
  carry GO-2026-6180 and GO-2026-6179. Base and direct-selected project
  findings are identical at 30/22/20/20 without a target trace; earlier
  advisory guards remain exact. No exploitability analysis occurred.
- Product source and dependency metadata remain unchanged. P7 stops for one
  reciprocal bounded product decision: explicitly retain inherited unloaded
  unqualified v1.5.0 under a target-specific exception, authorize one
  measurement-only gogo/protobuf owner/request study, or stop P7 unresolved.
  It may not repeat fixtures, add a direct target root, implement a parent
  change, transfer an exception, combine another group, or begin P8. Accepted
  quality remains 27/27 Q0-Q2 PASS at L2.

Kisielk Errcheck product decision (2026-09-21):

- Option 1 was selected. Exact selected, inherited, unloaded
  `github.com/kisielk/errcheck v1.5.0` is explicitly retained without changing
  product source, `go.mod`, or `go.sum`. It remains unqualified and is not
  described as secure. Zero loading bounds current exposure but is neither
  qualification nor implicit authorization. No stable exact-path release
  satisfies all Go-floor, exact-Go, repeatability, resource-ownership,
  closure-advisory, and target-only-selection contracts.
- The errcheck-specific, non-transferable exception accepts only the completed
  v1.5.0 loader/test and source-close failures; later-release floor, exact-Go,
  repeatability, and x/mod closure-advisory blockers; and completed repository,
  archive, release, module, API/command, ownership/global-state, MVS/loading,
  project, vulnerability, and related findings. It accepts no uncharacterized
  behavior, new advisory, or independently discovered defect, promotes no
  alternate identity or direct root, and transfers no other exception.
- Ownership is exact v1.5.0; the sole
  `gogo/protobuf@v1.3.2 -> errcheck@v1.5.0` request; and the genuine route from
  direct/imported/loaded Viper v1.15.0 through unloaded Gogo Protobuf v1.3.2.
  Guards require no direct root, negative errcheck and Gogo Protobuf why,
  zero repository imports and production/complete-test loads, no target runtime
  reachability, exact graph/module/tidy/Go-floor state, every earlier guard,
  and no new advisory, independent finding, release, owner, qualified release,
  supported tidy-stable owner, or compatible genuine qualified route.
- Any target path/version, request, Gogo Protobuf/Viper identity or ownership,
  root, import, load, runtime, graph, module-hash, tidy, Go-floor, earlier-
  guard, advisory, finding, release, owner, qualification, or route change
  expires the exception and requires a fresh errcheck dependency and product
  decision before merge. It authorizes no parent study, workaround, direct
  root, parent change, alternate path, product-source/metadata change,
  unrelated selection, or implementation.
- Guard-only revalidation preserved 234 modules, 3,599 edges, 355 production
  and 429 complete-test entries, 197 module-backed complete-test entries across
  41 modules, 1,067 sums, exact real module hashes, and the 432-line tidy
  projection at
  `3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`.
  All 32 guarded selections and negative why results remain exact; guarded
  imports and loads are zero; and 217 incoming edges reproduce SHA-256
  `500c57a1b4ffbdf5bb1296ef4dcb1bd16dbae12e54b86664fd1f8a2481a6c0dd`.
- Fresh proxy metadata retains 16 stable releases plus v1.5.0-alpha; GitHub
  retains the active unarchived non-fork repository; exact/package OSV and
  GitHub target results remain empty; guard OSV retains only the recorded
  Gorilla/retryablehttp pairs; and x/mod v0.14.0 retains GO-2026-6180 and
  GO-2026-6179. The 1,402-record Go index and PUBLISHED memberlist CNA bytes
  remain exact. No completed fixture or option-2 owner study was run.
- No dependency implementation or metadata commit was created. No changed-
  selection scorecard applies and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. Final unchanged-project exact Go 1.26.7 module verification, build,
  count-one tests, race count-one tests, and vet pass. The reciprocal
  247-archive chain, sole NEXT state, launcher/archive prompt mirror, exact
  documentation-only changed set, diff checks, and launcher check pass. Every
  task-owned scratch artifact was contained beneath the managed session root
  and removed; only its pre-existing launcher-owned Node compile cache remains.
  P7 continues only with the prepared bounded evaluation of selected exact-
  path `github.com/kisielk/gotool v1.0.0`. Current queue facts are four
  historical requests, negative why, zero repository imports, zero production/
  complete-test loads, and no runtime reachability; none is qualification.
  The successor must independently verify them, preserve this errcheck
  exception and every earlier guard, evaluate no other group, and stop for a
  fresh product decision if no stable exact-path release qualifies or a guard
  changes. It was prepared but not executed. P8 remains queued.

Kisielk Gotool evaluation and product-decision boundary (2026-09-21):

- Fresh proxy, sumdb, go-import, Git, and GitHub evidence identifies v1.0.0 as
  the only exact-path stable release. The sole tag, `master`, proxy `@latest`,
  and public active unarchived non-fork MIT repository agree on commit
  `80517062f582ea3340cd4baf70e86d539ae7d84d`. There is no later commit,
  GitHub Release, retraction, replacement, redirect, `/v2` line, or qualified
  alternate path/fork; the tag and commit are unsigned and proxy/Git regular
  files are byte-identical.
- Its two-package, standard-library-only closure has no declared Go floor and
  passes exact Go 1.26.7 and contained Go 1.18.10 upstream tests, repeats,
  race, vet, builds, supported cross-builds, and a bounded ordinary fixture.
  The API is deterministic for read-only calls, but exposes mutable default
  and BuildContext state requiring caller synchronization and reports
  omissions through global stderr.
- V1.0.0 fails the applicable module-aware package-pattern contract. In a
  small ordinary two-module local-replace fixture, `go list` expands the
  dependency module pattern while gotool returns no packages under both SDKs.
  Upstream maintainers state that module support is not viable, recommend
  `golang.org/x/tools/go/packages`, and say gotool can be considered
  deprecated. With no other stable exact-path release, none qualifies.
- Exact MVS ownership remains four v1.0.0 requests from Gogo Protobuf v1.3.2
  and three historical Honnef tools versions. The shortest genuine route is
  direct/imported/loaded Viper v1.15.0 through unloaded Gogo Protobuf v1.3.2.
  Gotool has negative why, zero repository imports, zero production/complete-
  test loads, and no runtime reachability; physical selection is not
  qualification.
- A disposable direct root changes no selection and only manufactures a main
  requirement and source sum; exact-Go project verification passes and tidy
  removes both. No projection or dependency commit was retained. The project
  remains 234 modules, 3,599 edges, 355 production and 429 complete-test
  entries, 197 module-backed complete-test entries across 41 modules, 1,067
  sums, exact real module hashes, and the recorded 432-line tidy projection.
  All 32 earlier why/import/load guards and 217 incoming edges remain exact.
- Fresh exact/package OSV, GitHub, repository-advisory, index, and pinned
  govulncheck evidence has no gotool finding or trace. Base and redundant-root
  project populations are identical at 30/22/20/20; the earlier guard
  advisories, index, and CNA evidence remain exact. No exploitability analysis
  occurred.
- Product source and dependency metadata remain unchanged; no changed-
  selection scorecard applies and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. The final unchanged-project exact Go 1.26.7 module verification,
  build, count-one tests, race count-one tests, and vet pass. P7 stops for one
  reciprocal choice: explicitly retain inherited unloaded unqualified v1.0.0
  under a target-specific exception; authorize exactly one later measurement-
  only study of one existing request/owner; or stop P7 unresolved. Do not
  repeat the evaluation, add a direct root, implement a parent/source change,
  transfer an exception, combine another group, or begin P8.

Kisielk Gotool product decision (2026-09-21):

- Option 1 was selected. Exact selected, inherited, unloaded
  `github.com/kisielk/gotool v1.0.0` is explicitly retained without changing
  product source, `go.mod`, or `go.sum`. It remains unqualified and is not
  described as secure. Physical selection, negative why, and zero loading
  bound current exposure but are neither qualification nor implicit
  authorization.
- The gotool-specific, non-transferable exception accepts only the completed
  module-pattern failure, upstream deprecation/module-support finding, and the
  completed repository/archive/release, module/Go-floor, source/API,
  GOPATH/GOROOT behavior, ownership/global-state, filesystem/stderr, MVS/
  loading, project, vulnerability, and related findings. It accepts no
  uncharacterized behavior, new advisory, or independent defect, promotes no
  alternate identity, direct root, or workaround, and transfers no exception.
- Ownership is exact v1.0.0; all four requests from Gogo Protobuf v1.3.2 and
  Honnef tools 2019.2.3/2020.1.3/2020.1.4; and the genuine Viper/Gogo and
  mvn-pom/Cloud/Storage/BigQuery/Honnef routes. Guards require no direct root,
  negative target/Gogo why, zero target/Gogo/Honnef imports and production/
  complete-test loads, no target runtime reachability, exact graph/module/
  tidy/Go-floor state, every earlier guard, and no new advisory, independent
  finding, repository/release/owner, qualified release, supported tidy-stable
  owner, or compatible genuine route to a qualified release.
- Any target path/version, request, owner identity/route, root, import, load,
  runtime, graph, module-hash, tidy, Go-floor, earlier-guard, advisory,
  finding, repository/release/owner, qualification, or route change expires
  the exception and requires a fresh gotool dependency and product decision
  before merge. It authorizes no owner study, workaround, direct root, parent
  or source/metadata change, alternate path, unrelated selection, or
  implementation.
- Guard-only revalidation preserved 234 modules, 3,599 edges, 355 production
  and 429 complete-test entries, 197 module-backed complete-test entries
  across 41 modules, 1,067 sums, exact real module hashes, and the 432-line
  tidy projection at
  `3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`.
  All four requests/routes and exact Go identity reproduce. All 32 earlier
  guarded selections and negative why results remain exact; guarded imports
  and loads are zero; and 217 incoming edges reproduce SHA-256
  `500c57a1b4ffbdf5bb1296ef4dcb1bd16dbae12e54b86664fd1f8a2481a6c0dd`.
  Including gotool gives 33 selections/221 edges at
  `dc3506a8e687d59a90e5711c347821b01672f8fa99ea060e494dadb505e320d5`.
- Fresh proxy/GitHub metadata retains sole v1.0.0, the active unarchived
  non-fork repository, one tag, no releases, and unchanged master. Target OSV
  and GitHub results remain empty; guard OSV retains only the recorded
  Gorilla/retryablehttp pairs; x/mod v0.14.0 retains GO-2026-6179 and
  GO-2026-6180; and the 1,402-record Go index and PUBLISHED memberlist CNA
  bytes remain exact. No behavior fixture, direct-root projection, or owner
  study ran.
- No dependency implementation or metadata commit was created. No changed-
  selection scorecard applies and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. P7 continues only with the prepared bounded evaluation of selected
  exact-path `github.com/konsorten/go-windows-terminal-sequences v1.0.1`.
  Its current queue facts are two historical Logrus v1.2.0/v1.4.2 requests,
  negative why, zero repository imports, zero production/complete-test loads,
  and no runtime reachability; none is qualification. The successor must
  independently verify every target/owner fact, preserve this exception and
  every earlier guard, evaluate no other group, and stop for a fresh product
  decision if no stable exact-path release qualifies or any guard changes.
  It was prepared but not executed. P8 remains queued.

Konsorten Go Windows Terminal Sequences evaluation (2026-09-22):

- No exact-path stable release qualifies. The exact proxy exposes v1.0.1,
  v1.0.2, and latest v1.0.3. Exact `go-import` points to the public enabled
  unarchived MIT Konsorten repository, which is a fork of the original
  Nine-Lives-Later repository. The three lightweight tags have continuous
  ancestry and master equals v1.0.3. Original-repository v1.0.4 changes the
  module/import identity and is an unauthorized alternate path. Proxy/sumdb/
  Git bytes, commits, trees, license, and repository/release state reproduce.
- Every exact stable is one standard-library-only package with no Go directive
  or requirements and preserves the Go 1.18 floor. Its sole export controls
  the Windows virtual-terminal-processing bit. V1.0.1/v1.0.2 retain an unsafe
  handle conversion that fails Windows vet and has a public ordinary crash
  report; v1.0.3 fixes that conversion. V1.0.2/v1.0.3 add deterministic
  Darwin/Linux error stubs, while v1.0.1 has no non-Windows source.
- All three exact releases retain the disqualifying ordinary behavior: they
  always read stdout's whole console mode, toggle one bit, and write that mode
  to the caller-supplied handle. Because console mode is handle-specific, a
  stderr or alternate-screen-buffer call can overwrite caller-owned mode bits
  with stdout's mode. The upstream stderr test omits state preservation.
  Highest stable v1.0.3 therefore remains unqualified. This is exact static
  source behavior, not an unobserved Windows-runtime test failure.
- MVS selects v1.0.1 through exactly two historical Logrus v1.4.2/v1.2.0
  requests. Both genuine routes begin at direct mvn-pom-mutator v0.2.3, then
  historical Viper v1.10.1, go-metrics v0.3.10, and Prometheus Common v0.9.1;
  the v1.2.0 route adds client_golang v1.0.0 and Common v0.4.1. Direct,
  imported, loaded Logrus v1.9.3 does not request/import the target and uses
  the same handle for mode read/write. Target why is negative, repository
  imports and production/complete-test loads are zero, and it is runtime-
  unreachable. None of these exposure facts is qualification.
- A disposable v1.0.3 root changed only target selection, manufactured one
  main edge and two sums, and preserved loads. Tidy removed it and restored
  inherited v1.0.1; no genuine owner requests v1.0.3. No projection or
  dependency commit was retained. The project remains
  234/3,599/355/429/197/41/1,067 with exact module hashes and 432-line tidy
  projection.
- Fresh exact OSV/GitHub/repository-advisory and isolated govulncheck evidence
  is empty. Base and disposable project populations are identical. All 33
  earlier guarded selections and 221 edges remain exact at SHA-256
  `dc3506a8e687d59a90e5711c347821b01672f8fa99ea060e494dadb505e320d5`;
  including the target gives 34 selections/223 edges at SHA-256
  `8cb329e95dcd81165af873c9bb4f11c85f61f7d68f39794394ab16a49f5cddc6`.
  Every guarded why remains negative and imports/loads zero; earlier advisory
  identities remain exact. No security reproduction or exploitability work
  occurred.
- Product source and dependency metadata remain unchanged. No changed-
  selection scorecard applies; accepted quality remains 27/27 Q0-Q2 PASS at
  L2. Final exact-Go verify/build/count-one/race/vet pass under canonical
  `umask 022`; an inherited-`umask 077` count-one attempt reproduced only the
  known permission-mode harness failures. P7 stops for one
  reciprocal product choice: explicitly retain inherited unloaded
  unqualified v1.0.1 under a target-specific exception; authorize one later
  measurement-only existing historical Logrus owner/request study; or stop P7
  unresolved. Do not repeat behavior, add a direct root, promote alternate-
  path v1.0.4, implement a parent/source change, transfer an exception,
  combine another group, or begin P8.

Konsorten Go Windows Terminal Sequences product decision (2026-09-22):

- Option 1 was selected. Exact selected, inherited, unloaded
  `github.com/konsorten/go-windows-terminal-sequences v1.0.1` is explicitly
  retained without changing product source, `go.mod`, or `go.sum`. It remains
  unqualified and is not described as secure. Physical selection, negative
  why, zero loading, and advisory absence bound exposure but are not
  qualification or implicit authorization.
- The target-specific, non-transferable exception accepts only the completed
  unsafe handle-conversion, non-Windows source/stub, handle-mode-preservation,
  and related repository/release/module/Go-floor/API/behavior, ownership/state,
  MVS/loading, project, and vulnerability findings. It accepts no
  uncharacterized behavior, new advisory, or independent defect, promotes no
  alternate identity, direct root, patch, wrapper, or workaround, and
  transfers no other exception.
- Ownership is exact v1.0.1; both requests from historical Logrus v1.4.2 and
  v1.2.0; and the genuine mvn-pom-mutator/Viper/go-metrics/Prometheus owner
  routes. Direct/imported/loaded Logrus v1.9.3 remains outside those requests
  and does not import the target. Guards require no direct target root,
  negative target why, zero repository imports and production/complete-test
  loads, no runtime reachability, exact graph/module/tidy/Go-floor state, and
  every earlier guard.
- The exception also requires no new target or closure advisory, independent
  finding, repository owner/status change, exact-path stable release,
  requester/owner, qualified release, supported tidy-stable owner, or
  compatible genuine route to a qualified exact-path release. Any target,
  request, owner identity/route, root, import, load, runtime, graph, module-
  hash, tidy, Go-floor, earlier-guard, advisory, finding, repository/release/
  owner, qualification, or route change expires the exception and requires a
  fresh target dependency and product decision before merge. It authorizes no
  owner study, parent/source/metadata change, direct root, alternate path,
  unrelated selection, or implementation.
- Guard-only revalidation preserved 234 modules, 3,599 edges, 355 production
  and 429 complete-test entries, 197 module-backed entries across 41 modules,
  1,067 sums, exact module hashes, and the 432-line tidy projection. Both
  requests/routes reproduce. All 33 earlier guarded selections retain
  negative why and zero imports/loads; their 221 incoming edges reproduce
  `dc3506a8e687d59a90e5711c347821b01672f8fa99ea060e494dadb505e320d5`.
  Including the target gives 34 selections/223 edges at
  `8cb329e95dcd81165af873c9bb4f11c85f61f7d68f39794394ab16a49f5cddc6`.
- Fresh proxy/go-import/GitHub metadata retains the exact release/repository
  identities. Exact target OSV and GitHub results remain empty; guard OSV
  retains only the Gorilla/retryablehttp pairs; x/mod v0.14.0 retains
  GO-2026-6179 and GO-2026-6180; the 1,402-record Go index and PUBLISHED
  memberlist CNA bytes remain exact. No behavior evaluation, direct-root
  projection, or owner study was repeated.
- No dependency implementation or metadata commit was created. No changed-
  selection scorecard applies and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. Final exact-Go module verification, build, count-one tests, race
  count-one tests, vet, the reciprocal 251-archive chain, sole NEXT state,
  launcher/archive prompt mirror, exact five-file documentation-only changed
  set, diff checks, launcher check, and all 62 launcher controls pass. Every
  task-owned scratch artifact was contained beneath the managed session root
  and removed; only its pre-existing launcher-owned Node compile cache
  remains. P7 continues only with the prepared bounded evaluation of selected
  exact-path `github.com/kr/fs v0.1.0`. Its current queue facts are exact
  requests from SFTP v1.10.1/v1.13.1, negative why, zero repository imports,
  zero production/complete-test loads, and no runtime reachability; none is
  qualification. The successor must independently verify them, preserve this
  exception and every earlier guard, evaluate no other group, and stop for a
  fresh product decision if no exact-path stable release qualifies or any
  guard changes. It was prepared but not executed. P8 remains queued.

Kr Fs evaluation (2026-09-22):

- No exact-path stable release qualifies. The proxy exposes only v0.1.0,
  whose exact module file has no Go directive, requirements, retractions,
  deprecation, or replacement. Exact go-import metadata resolves without
  redirect to the public enabled unarchived non-fork BSD-3-Clause `kr/fs`
  repository. Its sole annotated tag is unsigned, peels to commit
  `1455def202f6e05b95cc7bfc7e8ae67ae5141eba`, and matches `master`; there
  are no GitHub Releases.
- Default `main` is one unsigned README-only commit later and resolves only as
  `v0.1.1-0.20210218185759-c64b65e7619f`; all Go source is identical. Its
  `kr.dev/walk` recommendation is a distinct unauthorized module path. Exact
  `/v2` and `/v3` lines do not exist. Proxy and tag files agree byte for byte;
  no pseudo-version, fork, redirect, alternate path, or replacement was
  promoted.
- V0.1.0 is one standard-library-only package, preserves Go 1.18, and exports
  `FileSystem`, `Walker`, `Walk`, `WalkFS`, and Walker methods. Exact Go
  1.26.7 and contained Go 1.18.10 pass upstream build/tests/repeated race/vet,
  supported cross-builds, and all non-disqualifying bounded behavior coverage.
- The stable release fails its documented ordinary behavior. `Walker`
  promises lexical traversal unconditionally, but `WalkFS` never sorts a
  caller filesystem's `ReadDir` result. A bounded filesystem returning
  ordinary root entries `[b, a]` yields `[root, root/b, root/a]`, not
  `[root, root/a, root/b]`, under both SDKs. The main pseudo-version has the
  same source and defect. No fuzzing, stress, adversarial path generation,
  security reproduction, or exploitability work occurred.
- MVS selects v0.1.0 through exact requests from SFTP v1.13.1 and v1.10.1.
  The current route is main -> direct Afero v1.9.4 -> SFTP v1.13.1 -> target;
  the historical route is main -> direct mvn-pom-mutator v0.2.3 -> Viper
  v1.10.1 -> Afero v1.6.0 -> SFTP v1.10.1 -> target. Both SFTP clients
  implement the target interface, call `WalkFS`, and preserve server reply
  order. Neither loads in this project; current Afero loads no SFTP adapter.
- Target why is negative, repository imports and production/complete-test
  loads are zero, and there is no runtime reachability. Disposable stable and
  pseudo direct roots manufactured only one main edge and sums, preserved all
  loads, and were removed by tidy; pseudo tidy restored inherited v0.1.0. No
  direct root, projection, owner change, or dependency commit was retained.
- The project remains 234/3,599/355/429/197/41/1,067 with exact module hashes
  and the recorded 432-line tidy projection. All 34 guarded selections and
  their 223 incoming edges retain SHA-256
  `8cb329e95dcd81165af873c9bb4f11c85f61f7d68f39794394ab16a49f5cddc6`;
  every guarded why remains negative and imports/loads zero.
- Fresh exact OSV/GitHub/repository-advisory and isolated govulncheck evidence
  has no target finding. Base and disposable project finding populations are
  identical. Earlier guard advisories, the exact 1,402-record Go index, and
  PUBLISHED memberlist CNA response remain unchanged. Advisory absence is not
  qualification.
- Product source and dependency metadata remain unchanged. No changed-
  selection scorecard applies and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. Final exact-Go verify/build/count-one/race/vet pass. P7 stops for one
  reciprocal product choice after the reciprocal 252-archive chain, sole NEXT
  state, launcher/archive prompt mirror, exact five-file documentation-only
  changed set, diff checks, launcher check, and all 62 launcher controls pass.
  Every task-owned scratch artifact was contained beneath the managed session
  root and removed; only its pre-existing launcher-owned Node compile cache
  remains. The choice is: retain exact inherited/unloaded unqualified v0.1.0
  under a target-specific exception; authorize one later measurement-only
  study of the current Afero v1.9.4 -> SFTP v1.13.1 owner/request route; or
  stop P7 unresolved. Do not repeat completed behavior, add a direct root,
  promote the pseudo-version or alternate path, change a parent/source/
  metadata, transfer an exception, combine another group, or begin P8.

Kr Fs product decision (2026-09-22):

- Option 1 was selected. Exact selected, inherited, unloaded
  `github.com/kr/fs v0.1.0` is explicitly retained without changing product
  source, `go.mod`, or `go.sum`. It remains unqualified and is not described as
  secure. Physical selection, negative why, zero loading, and advisory absence
  bound exposure but are not qualification or implicit authorization.
- The kr/fs-specific, non-transferable exception accepts only the completed
  lexical-order failure and the related repository/release/module/Go-floor/
  API/ordinary-behavior, ownership/state, MVS/loading, project, vulnerability,
  and related completed evaluation findings. It accepts no uncharacterized
  behavior, new advisory, or independent defect, promotes no alternate path,
  pseudo-version, direct root, patch, wrapper, or workaround, and transfers no
  other exception.
- Ownership is exact v0.1.0; both SFTP v1.13.1/v1.10.1 requests; the current
  direct Afero v1.9.4 -> SFTP v1.13.1 route; and the historical direct
  mvn-pom-mutator v0.2.3 -> Viper v1.10.1 -> Afero v1.6.0 -> SFTP v1.10.1
  route. Guards require no direct target root, negative target why, zero target
  repository imports and production/complete-test loads, no runtime
  reachability, exact graph/module/tidy/Go-floor state, and every earlier
  guard. Both SFTP clients must retain the completed interface/use/order facts,
  neither SFTP may load, and selected Afero may load only its recorded root,
  `internal/common`, and `mem` packages.
- The exception also requires no new target or requester-closure advisory,
  independent finding, repository owner/status change, exact-path stable
  release, requester/owner, qualified release, supported tidy-stable owner, or
  compatible genuine route to a qualified exact-path release. Any target,
  request, owner identity/route/behavior, root, import, load, runtime, graph,
  module-hash, tidy, Go-floor, earlier-guard, advisory, finding, repository/
  release/owner, qualification, or route change expires the exception and
  requires a fresh target dependency and product decision before merge. It
  authorizes no owner study, parent/source/metadata change, direct root,
  alternate path, unrelated selection, or implementation.
- Guard-only revalidation preserved 234 modules, 3,599 edges, 355 production
  and 429 complete-test entries, 197 module-backed entries across 41 modules,
  1,067 sums, exact module hashes, and the 432-line tidy projection. Both
  requests/routes and target/SFTP/Afero load boundaries reproduce. All 34
  earlier guarded selections retain negative why and zero imports/loads; their
  223 incoming edges reproduce
  `8cb329e95dcd81165af873c9bb4f11c85f61f7d68f39794394ab16a49f5cddc6`.
  Including kr/fs gives 35 selections/225 edges at
  `5dae96771b3994a9ce1999f8d0487f152d94adb3b2d0bb467b18054be39dc284`.
- Fresh proxy/GitHub metadata retains the exact release/repository identities.
  Exact target OSV and GitHub results remain empty; guard OSV retains only the
  Gorilla/retryablehttp pairs; x/mod v0.14.0 retains GO-2026-6179 and
  GO-2026-6180; the 1,402-record Go index and PUBLISHED memberlist CNA bytes
  remain exact. No behavior evaluation, direct-root projection, or owner study
  was repeated.
- No dependency implementation or metadata commit was created. No changed-
  selection scorecard applies and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. Final exact-Go module verification, build, count-one tests, race
  count-one tests, vet, the reciprocal 253-archive chain, single NEXT state,
  launcher/archive prompt mirror, exact five-file documentation-only changed
  set, diff checks, and launcher check pass. Every task-owned scratch artifact
  was contained beneath the managed session root and removed; only its pre-
  existing launcher-owned Node compile cache remains. P7 continues only with
  the prepared bounded evaluation of selected exact-path
  `github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515`. Current queue
  observations are three requests from go-logfmt v0.4.0, Prometheus Common
  v0.4.1, and Prometheus TSDB v0.7.1, negative why, zero repository imports,
  zero production/complete-test loads, and no runtime reachability; none is
  qualification. The successor must independently verify them and every
  earlier guard, evaluate no other group, and stop for a fresh product
  decision if no exact-path release qualifies or any guard changes. It was
  prepared but not executed. P8 remains queued.

Kr Logfmt evaluation (2026-09-22):

- No exact-path release qualifies. The proxy version list is empty and the
  repository has no tags or GitHub Releases. Exact selected
  `v0.0.0-20140226030751-b84e30acd515` and proxy `@latest`
  `v0.0.0-20210122060352-19f9bcb100e6` are the only serious candidates.
  Selected has no Go directive or requirements; latest adds only `go 1.14`.
  Neither has retractions, deprecation, or replacement metadata.
- Exact go-import resolves without redirect to the public enabled unarchived
  non-fork MIT `kr/logfmt` repository. Selected is unsigned commit
  `b84e30acd515aadc4b783ad4ff83aff3299bdfe0`; `main`/latest is its direct
  GitHub-verified signed child
  `19f9bcb100e6bcb308b5db29c682de01e9b3f2e6`, adding only `LICENSE` and
  `go.mod`. All Go source/tests are identical. Divergent `maybe` is older and
  untagged; `/v2` and `/v3` do not exist; and go-logfmt is a distinct
  unauthorized path. Proxy and Git files agree byte for byte.
- Both candidates are the same single standard-library-only parser package,
  preserve Go 1.18, and export `Unmarshal`, handler/struct-handler APIs, and
  related errors. There is no encoder/formatter API, cgo, generated source,
  global mutable state, network/subprocess boundary, platform branch, or
  library-owned resource. Exact Go 1.26.7 and contained Go 1.18.10 pass
  upstream build/repeated tests/race/vet, supported cross-builds, and all non-
  disqualifying bounded ordinary fixtures.
- Three exported ordinary contracts fail in both candidates and SDKs: a non-
  struct `&int` target panics despite the promise to return an error; the
  documented all-numeric-types support rejects `uint8`; and ordinary value
  `128` silently wraps to `-128` in `int8`. Latest cannot fix them because its
  Go source is byte-identical. No fuzzing, stress, adversarial malformed
  records, security reproduction, or exploitability work occurred.
- MVS selects exact selected through requests from go-logfmt v0.4.0,
  Prometheus Common v0.4.1, and Prometheus TSDB v0.7.1. Genuine routes all
  begin at direct mvn-pom-mutator v0.2.3: directly through TSDB, or through
  historical Viper v1.10.1 -> go-metrics v0.3.10 -> Common v0.9.1 and
  go-logfmt, or client_golang v1.0.0 -> Common v0.4.1. Go-logfmt imports the
  target only from a build-tagged fuzz file and benchmark; the other two
  requesters have stale indirect requirements and no Go imports.
- Target why is negative, repository imports and production/complete-test
  loads are zero, and there is no runtime reachability. Disposable selected
  and latest direct roots manufactured only a main edge/sums, preserved all
  loads, and were removed by tidy. Latest's only selection delta was selected
  -> latest; tidy restored selected because no genuine owner requests latest.
  No direct root, projection, parent change, or dependency commit was retained.
- The project remains 234/3,599/355/429/197/41/1,067 with exact module hashes
  and the recorded 432-line tidy projection. All 35 earlier guarded selections
  and their 225 incoming edges retain SHA-256
  `5dae96771b3994a9ce1999f8d0487f152d94adb3b2d0bb467b18054be39dc284`;
  all guarded why results remain negative and imports/loads zero.
- Fresh exact OSV/GitHub/repository-advisory and isolated govulncheck evidence
  has no candidate finding. Base and both disposable project populations are
  identical at 30 module, 22 package, 20 symbol, and 20 test-symbol OSVs with
  no target trace. Earlier guard advisories, the exact 1,402-record Go index,
  and PUBLISHED memberlist CNA response remain unchanged. Advisory absence is
  not qualification.
- Product source and dependency metadata remain unchanged. No changed-
  selection scorecard applies and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. Final exact-Go verify/build/count-one/race/vet pass under `umask 022`.
  P7 stops for one reciprocal product decision after the reciprocal 254-
  archive chain, sole NEXT state, launcher/archive prompt mirror, exact five-
  file documentation-only changed set, diff checks, and launcher check pass.
  Every task-owned scratch artifact was contained beneath the managed session
  root and removed; only its pre-existing launcher-owned Node compile cache
  remains. The choice is: retain exact selected/inherited/unloaded
  unqualified kr/logfmt under a target-specific exception; authorize one later
  measurement-only study of the direct mvn-pom-mutator v0.2.3 -> TSDB v0.7.1
  request route; or stop P7 unresolved. Do not repeat completed behavior, add
  a direct root, promote latest/go-logfmt, change a parent/source/metadata,
  transfer an exception, combine another group, or begin P8.

Kr Logfmt product decision (2026-09-22):

- Option 1 was selected. Exact selected, inherited, unloaded
  `github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515` is explicitly
  retained without changing product source, `go.mod`, or `go.sum`. It remains
  unqualified and is not described as secure. Physical selection, negative
  why, zero loading, and advisory absence bound exposure but are not
  qualification or implicit authorization.
- The kr/logfmt-specific, non-transferable exception accepts only the
  completed non-struct panic, `uint8` rejection, and `int8` overflow-wrap
  failures plus the related repository/release/module/Go-floor/API/ordinary-
  behavior, ownership/state, MVS/loading, project, vulnerability, and related
  completed findings. It accepts no uncharacterized behavior, new advisory,
  or independent defect and transfers no other exception.
- Ownership is exact selected; all three requests from go-logfmt v0.4.0,
  Prometheus Common v0.4.1, and Prometheus TSDB v0.7.1; all genuine routes
  beginning at direct mvn-pom-mutator v0.2.3 through TSDB or historical
  Viper/go-metrics/Common/go-logfmt/client_golang; and the recorded requester
  import facts. Guards require no direct target root, negative target why,
  zero target repository imports and production/complete-test loads, no
  runtime reachability, exact graph/module/tidy/Go-floor state, and every
  earlier guard.
- The exception also requires no new target or requester-closure advisory,
  independent finding, repository owner/status change, exact-path release,
  requester/owner, qualified release, supported tidy-stable owner, or
  compatible genuine route to a qualified exact-path release. Any target,
  request/requester import, owner identity/route, root, import, load, runtime,
  graph, module-hash, tidy, Go-floor, earlier-guard, advisory, finding,
  repository/release/owner, qualification, or route change expires the
  exception and requires a fresh target dependency and product decision
  before merge. It authorizes no owner study, parent/source/metadata change,
  direct root, latest, alternate path, fork, replacement, patch, wrapper,
  workaround, unrelated selection, or implementation.
- Guard-only revalidation preserved 234 modules, 3,599 edges, 355 production
  and 429 complete-test entries, 197 module-backed entries across 41 modules,
  1,067 sums, exact module hashes, and the 432-line tidy projection. All three
  requests, routes, requester facts, and target boundaries reproduce. All 35
  earlier guarded selections retain negative why and zero imports/loads; their
  225 incoming edges reproduce
  `5dae96771b3994a9ce1999f8d0487f152d94adb3b2d0bb467b18054be39dc284`.
  Including kr/logfmt gives 36 selections/228 incoming edges at
  `361d0c355a69842518a518e682f81c9728d37acfdeff64f430a4fb253929691e`.
- Fresh proxy/GitHub metadata retains exact release/repository identities.
  Exact target OSV/GitHub results remain empty; guard OSV retains only the
  Gorilla/retryablehttp pairs; x/mod v0.14.0 retains GO-2026-6179 and
  GO-2026-6180; the 1,402-record Go index and PUBLISHED memberlist CNA bytes
  remain exact. No completed behavior, direct-root projection, or owner study
  was repeated.
- No dependency implementation or metadata commit was created. No changed-
  selection scorecard applies and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. Final exact-Go module verification, build, count-one tests, race
  count-one tests, vet, the reciprocal 255-archive chain, sole NEXT state,
  launcher/archive prompt mirror, exact five-file documentation-only changed
  set, diff checks, and launcher check pass. Every task-owned scratch artifact
  was contained beneath the managed session root and removed; only its pre-
  existing launcher-owned Node compile cache remains. P7 continues only with
  the prepared bounded evaluation of selected exact-path
  `github.com/kr/pretty v0.3.1`. Current queue observations are a selected
  v0.3.1 request from Cast v1.5.1, four lower historical requests, a positive
  why chain only through yaml.v2 tests/check.v1, zero repository imports, and
  zero production/complete-test loads; none is qualification. The successor
  must independently verify these facts and every earlier guard, evaluate no
  other group, and stop for a fresh product decision if no exact-path release
  qualifies or any guard changes. It was prepared but not executed. P8 remains
  queued.

Kr Pretty evaluation (2026-09-22):

- No exact-path stable qualifies. The canonical proxy and continuous
  annotated-tag line is v0.1.0, v0.2.0, v0.2.1, v0.3.0, and selected/latest
  v0.3.1. Exact `/v2` and `/v3` lines do not exist. Go-import resolves without
  redirect to the public enabled unarchived non-fork MIT `kr/pretty`
  repository, owned by `kr` with protected default `main`; it has no GitHub
  Releases. Default main, proxy latest, and v0.3.1 are exact commit
  `3cd153a126da607b78d1762779b1e1054f9889fc`.
- Every release names exact `github.com/kr/pretty`. V0.1.0 has no Go directive;
  later stables use `go 1.12`. There are no retractions, deprecation, or
  replacements, and all releases preserve Go 1.18. Selected requires kr/text
  v0.2.0 and go-internal v1.9.0; the selected closure's maximum Go directive
  is 1.17. Proxy ZIPs and peeled Git regular files agree byte for byte.
- The module is one formatting/diff package with no cgo, generated source,
  embed, build tag, platform source, external I/O boundary, library-owned
  resource, or production mutable global. Exact Go 1.26.7 and contained Go
  1.18.10 pass build, vet-disabled repeated upstream tests, supported test
  cross-builds, and bounded positive ordinary/non-mutation/concurrency checks
  for all five stables.
- Every stable fails qualification in both SDKs. Full upstream vet rejects its
  invalid `unsafe.Pointer(uintptr(1))` test construction, and upstream race
  count-ten aborts there with a checkptr bad-pointer fatal error; the first two
  releases also retain an int-to-string diagnostic. Ordinary `Formatter(nil)`
  produces a Format-method reflect panic diagnostic rather than fmt's `<nil>`
  pass-through representation. Diff output order also changes across ten runs
  for an ordinary four-entry map because unsorted map keys are traversed.
  Advisory absence and positive non-disqualifying checks cannot override these
  failures.
- MVS requests are Cast v1.5.1 -> v0.3.1, Consul SDK v0.8.0 -> v0.2.0, and
  client_golang v1.4.0, sergi/go-diff v1.2.0, and errgo.v2 v2.1.0 -> v0.1.0.
  The shortest genuine current owner is main -> selected Cast v1.5.1 ->
  target; Cast's quicktest v1.14.4 test closure imports target, while Cast and
  every lower requester do not. Lower historical routes through mvn-pom-
  mutator/Viper/Consul or client_golang, direct sergi/Assert, and Honnef/go-
  internal/errgo reproduce.
- Target why is positive only through project config -> yaml.v2 tests ->
  check.v1; repository imports and production/complete-test loads are zero,
  with no runtime reachability. A disposable direct v0.3.1 get manufactures
  target/text/go-internal roots and reaches 236 modules/3,606 edges/1,072 sums
  without changing loads; tidy removes them and preserves Cast as owner. Lower
  projections downgrade accepted parents/closures or make the project
  unloadable and are not supported guard-preserving owners. No projection,
  direct root, or dependency commit was retained.
- The real project remains 234 modules, 3,599 edges, 355 production entries,
  429 complete-test entries, 197 module-backed entries across 41 loaded
  modules, 1,067 sums, exact module hashes, and the recorded 432-line tidy
  projection. All 36 earlier guarded selections remain exact with negative
  why and zero imports/loads; their 228 incoming edges reproduce
  `361d0c355a69842518a518e682f81c9728d37acfdeff64f430a4fb253929691e`.
- Exact OSV/GitHub/repository-advisory and isolated govulncheck evidence has no
  target finding. Base and direct-v0.3.1 project populations are identical at
  30/22/20/20 with no target trace. Guard Gorilla/retryablehttp pairs, x/mod
  findings, the 1,402-record Go module index, and the PUBLISHED memberlist CNA
  identity remain exact. The evaluation used only bounded ordinary values and
  upstream tests; it did not fuzz, stress, generate adversarial values,
  reproduce a security issue, or assess exploitability.
- Product source and dependency metadata remain unchanged. No changed-
  selection scorecard applies and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. Final exact-Go verify/build/count-one/race/vet pass. P7 stops for one
  reciprocal bounded product decision: explicitly retain exact selected,
  inherited, unloaded unqualified v0.3.1 under a kr/pretty-specific exception;
  authorize exactly one later measurement-only study of the existing main ->
  Cast v1.5.1 -> target owner/request route; or stop P7 unresolved. Do not
  repeat the evaluation, reopen Cast, add a direct target root, change source/
  metadata or a parent, transfer an exception, combine another group, or begin
  P8. The successor was prepared but not executed.

Kr Pretty product decision (2026-09-22):

- Option 1 was selected. Exact selected, inherited, unloaded
  `github.com/kr/pretty v0.3.1` is explicitly retained without changing
  product source, `go.mod`, or `go.sum`. It remains unqualified and is not
  described as secure. Physical selection, the positive test-only why chain,
  zero target loading, and advisory absence bound exposure but are not
  qualification or implicit authorization.
- The kr/pretty-specific, non-transferable exception accepts only the
  completed nil-Formatter pass-through failure, map-Diff output-order
  nondeterminism, upstream vet failures, upstream race/checkptr failure, and
  completed related evaluation. It accepts no uncharacterized behavior, new
  advisory, or independent defect and transfers no earlier exception.
- The exception requires the exact path/version; all five requests and their
  requested versions; the exact main -> Cast v1.5.1 -> target current owner
  route; Cast's loaded Viper route and quicktest v1.14.4/pretty v0.3.1/
  go-internal v1.9.0 test closure; every lower historical route; and every
  requester/import fact to remain exact. It also requires no direct target
  root, the recorded yaml.v2-tests/check.v1 why chain only, zero target
  repository imports and production/complete-test loads, no runtime
  reachability, exact graph/module/tidy/Go-floor state, and every earlier
  guard.
- Any target/request/requester-import/Cast-closure/owner-route, root/why/
  import/load/runtime, graph/module/tidy/Go-floor, earlier-guard, advisory/
  independent-finding, repository/release/owner, qualification, supported
  tidy-stable-owner, or compatible qualified-route change expires the
  exception and requires a fresh kr/pretty dependency and product decision
  before merge. Any compatible genuine route to a qualified exact-path release
  expressly expires the retention. No owner study, Cast or parent change,
  direct root, source/metadata edit, alternate path, workaround, unrelated
  selection, or implementation is authorized.
- Guard-only revalidation preserved 234 modules, 3,599 edges, 355 production
  and 429 complete-test entries, 197 module-backed entries across 41 modules,
  1,067 sums, exact module hashes, the 432-line tidy projection, and applied
  52/948-line tidy hashes. All five requests, current and historical routes,
  requester facts, and target boundaries reproduce. All 36 earlier guarded
  selections retain negative why and zero imports/loads; their 228 incoming
  edges reproduce
  `361d0c355a69842518a518e682f81c9728d37acfdeff64f430a4fb253929691e`.
  Including kr/pretty gives 37 selections/233 incoming edges at
  `1779751e93109c20a3cc7e3f2215665c9115156eb3947a1e45afb51e5cfe364f`.
- Fresh proxy/GitHub metadata retains the exact release/repository identities.
  Exact target OSV/GitHub results remain empty; guard OSV retains only the
  Gorilla/retryablehttp pairs; x/mod v0.14.0 retains GO-2026-6179 and
  GO-2026-6180; the 1,402-record Go index and PUBLISHED memberlist CNA bytes
  remain exact. No completed behavior, release, closure, direct-root
  projection, or owner study was repeated.
- No dependency implementation or metadata commit was created. No changed-
  selection scorecard applies and accepted quality remains 27/27 Q0-Q2 PASS
  at L2. Final exact-Go module verification, build, count-one tests, race
  count-one tests, vet, reciprocal 257-archive chain, single NEXT state,
  launcher/archive prompt mirror, exact documentation-only changed set, diff
  checks, and launcher check pass. Every task-owned scratch artifact was
  contained beneath the managed session root and removed; only its pre-
  existing launcher-owned Node compile cache remains. P7 continues only with
  the prepared bounded evaluation of selected exact-path
  `github.com/kr/pty v1.1.1`. Current queue observations are one request from
  historical kr/text v0.1.0, negative why, zero repository imports, and zero
  production/complete-test loads; none is qualification. The successor must
  independently verify them and every earlier guard, evaluate no other group,
  and stop for a fresh product decision if no exact-path release qualifies or
  any guard changes. It was prepared but not executed. P8 remains queued.

Kr Pty evaluation (2026-09-22):

- The canonical exact-path proxy line contains exactly v1.0.0, v1.0.1, and
  v1.1.0-v1.1.8; exact `/v2` and `/v3` lines do not exist. Go-import resolves
  to public enabled archived MIT fork `kr/pty`, owned by `kr`, with active
  non-fork `creack/pty` as parent/source and exact v1.1.8 on default `master`.
  There are no GitHub Releases, retractions, or replacements. Proxy/sumdb,
  tag/commit/tree, and Git regular-file identities agree for all eleven
  releases; v1.1.7 is the transferred-parent exact-path release.
- V1.0.0-v1.1.4 and transferred v1.1.7 have no requirement or Go directive;
  v1.1.5/v1.1.6 declare `go 1.12`. V1.1.8 branches from v1.1.6, declares
  `go 1.12`, and is a deprecated shim requiring and delegating to alternate
  path creack/pty v1.1.7. No alternate path, fork, branch, pseudo-version, or
  replacement was promoted.
- The pre-shim release is one PTY package with no upstream tests, cgo, network,
  or mutable shared state. Selected v1.1.1 exports Open/Start and terminal-size
  helpers; v1.1.3 adds StartWithSize. Callers own returned descriptors and
  terminal restoration; Start closes the slave and master-on-error. Small
  bounded fixtures cover Open/close, size operations and input nonmutation,
  four independent concurrent opens, and short-lived `/bin/echo`/`stty`
  processes under exact Go 1.26.7 and contained Go 1.18.10.
- V1.0.0-v1.1.4 pass every applicable ordinary fixture, build, repeated,
  race, vet, and supported cross-build contract under both SDKs. V1.1.5-
  v1.1.8 fail documented Start/StartWithSize with `Setctty set but Ctty not
  valid in child`; repeated and race invocations reproduce the functional
  failure. V1.1.6/v1.1.7 also fail their claimed Solaris closure and cannot
  resolve its unpinned x/sys under Go 1.18; v1.1.8 also fails its Windows shim
  build. V1.1.4 is therefore the highest behavior-qualified exact-path stable
  preserving Go 1.18. Selected v1.1.1 passes but is lower.
- Exact MVS has one target request: historical kr/text v0.1.0 -> kr/pty
  v1.1.1. Requester inspection was ownership-only: kr/text imports target only
  from its unloaded `mc` command. Genuine historical routes extend kr/pretty
  v0.1.0/v0.2.0 routes through direct sergi/Assert, mvn-pom-mutator/Viper/
  Consul or Prometheus, and three Honnef/go-internal/errgo paths. Target why is
  negative; repository imports and production/complete-test target loads are
  zero. Kr/text why is positive only through yaml.v2 tests/check.v1/kr/pretty,
  but neither its command nor target loads. There is no runtime reachability.
- Disposable exact v1.1.1 and v1.1.4 gets each manufacture a redundant target
  root. They produce 234 modules/3,600 edges and 1,101/1,102 sums respectively,
  preserve all 355/429/197/41 load facts, and tidy removes the root and restores
  selected v1.1.1/common tidy state. V1.1.8 also manufactures creack/pty
  v1.1.7, producing 235/3,602/1,104; tidy removes both. No genuine supported
  tidy-stable owner requests v1.1.4, so no projection or dependency commit was
  retained.
- The unchanged project remains 234/3,599/355/429/197/41/1,067 with exact
  module hashes and the recorded 432-line tidy projection. All 37 earlier
  selections and their 233 incoming edges retain SHA-256
  `1779751e93109c20a3cc7e3f2215665c9115156eb3947a1e45afb51e5cfe364f`;
  36 why results are negative and only the closed kr/pretty result is positive,
  with zero guarded imports/loads. Including kr/pty gives 38 selections/234
  edges at
  `d31b2afdb5a11dadac236eec600408cc9ba0ee3deca2793ae2adf7ffc449c320`.
- Exact OSV/GitHub results for all stables and repository advisories are empty.
  Exact-Go pinned govulncheck has no target finding; base/direct-v1.1.4 project
  populations are identical at 30/22/20/20. Guard Gorilla/retryablehttp pairs,
  x/mod v0.14.0 findings, the exact 1,402-record Go index, and PUBLISHED
  memberlist CNA response remain unchanged. Advisory absence does not override
  behavior or ownership.
- No source or dependency metadata changed. Final exact Go 1.26.7 module
  verification, build, count-one tests, race count-one tests, and vet pass
  under `umask 022`; accepted quality remains 27/27 Q0-Q2 PASS at L2. P7 stops
  for one reciprocal choice: retain selected inherited/unloaded v1.1.1 under a
  kr/pty-specific unqualified exception; authorize one later measurement-only
  study of the shortest main -> direct sergi/go-diff v1.2.0 -> historical
  kr/pretty v0.1.0 -> historical kr/text v0.1.0 -> target route; or stop P7
  unresolved. Do not repeat the evaluation, evaluate kr/text, add a direct
  root, select v1.1.4 without an owner, promote v1.1.8's alternate path,
  transfer an exception, combine another group, or begin P8. The decision
  successor was prepared but not executed.

Kr Pty product decision (2026-09-22):

- Option 1 was selected. Exact selected, inherited, unloaded
  `github.com/kr/pty v1.1.1` is explicitly retained without changing product
  source, `go.mod`, or `go.sum`. It remains unqualified and is not described
  as secure. Exact v1.1.4 remains the highest behavior-qualified
  Go-1.18-compatible exact-path stable, but no genuine supported tidy-stable
  owner requests it. Physical selection, negative target why, zero loading,
  and advisory absence are exposure bounds, not qualification.
- The target-specific, non-transferable exception accepts only the completed
  eleven-release identity; selected v1.1.1 passing contracts; v1.1.4
  qualification and unsupported-owner result; v1.1.5-v1.1.8 ordinary Start
  failure; later Solaris/Windows closure failures; sole historical kr/text
  v0.1.0 request and unloaded `mc` import; all sergi/Assert,
  mvn-pom-mutator/Viper/Consul/Prometheus, and Honnef/go-internal/errgo route
  families; why/import/load/runtime boundary; graph/tidy/Go-floor state;
  disposable projections; earlier guards; advisory identities; and completed
  related evaluation. No earlier exception transfers.
- The exception requires exact selected path/version and sole request; the
  requester's target import and every genuine owner route; no supported
  v1.1.4 owner; no direct target root; negative target why; zero target
  repository imports and production/complete-test loads; only the recorded
  positive kr/text why chain through yaml.v2 tests/check.v1/kr/pretty; no
  runtime reachability; exact release/repository/parent/shim identities and
  qualification results; exact project/module/tidy/Go-floor state; every
  earlier guard; and no new advisory, independent defect, release, owner,
  qualified stable, supported tidy-stable owner, or compatible genuine route
  to qualified v1.1.4.
- Any target/request/requester-import/owner-route, root/why/import/load/runtime,
  graph/module/tidy/Go-floor, earlier-guard, advisory/finding,
  repository/release/owner, qualification, supported-owner, or compatible-
  route change expires the exception and requires a fresh kr/pty dependency
  and product decision before merge. Any compatible genuine route to
  qualified v1.1.4 expressly expires retention. No owner study, direct root,
  alternate path, parent/source/metadata change, workaround, unrelated
  selection, implementation, or transferred exception is authorized.
- Guard-only revalidation preserved 234 modules, 3,599 edges, 355 production
  and 429 complete-test entries, 197 module-backed entries across 41 modules,
  1,067 sums, exact module hashes, the 432-line tidy projection, applied
  52/948-line tidy hashes, the sole request, every historical route, and the
  target/text why/import/load boundary. All 37 earlier guarded selections and
  233 incoming edges remain exact at
  `1779751e93109c20a3cc7e3f2215665c9115156eb3947a1e45afb51e5cfe364f`;
  including kr/pty gives 38 selections/234 incoming edges at
  `d31b2afdb5a11dadac236eec600408cc9ba0ee3deca2793ae2adf7ffc449c320`.
- Fresh official Go 1.26.7 archive/binary identities reproduce as
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
  Fresh release/repository/advisory identities remain exact. Guard OSV retains
  only the Gorilla/retryablehttp pairs; x/mod v0.14.0 retains GO-2026-6179 and
  GO-2026-6180; the 1,402-record Go index and PUBLISHED memberlist CNA bytes
  remain exact. No completed PTY behavior, source/release/closure, or
  disposable-project evaluation was repeated, and no owner study ran.
- No dependency implementation or metadata commit was created. No changed-
  selection scorecard applies; accepted quality remains 27/27 Q0-Q2 PASS at
  L2. Final exact-Go module verification, build, count-one tests, race
  count-one tests, vet, reciprocal 259-archive chain, single NEXT state,
  launcher/archive prompt mirror, exact documentation-only changed set, diff
  checks, and launcher check pass. Every task-owned scratch artifact was
  contained beneath the managed session root and removed; only its pre-
  existing launcher-owned Node compile cache remains. P7 continues only with
  the prepared bounded evaluation of selected exact-path
  `github.com/kr/text v0.2.0`. It must preserve the kr/pty exception and stop
  before implementation if any kr/pty expiry guard changes. The successor was
  prepared but not executed. P8 remains queued.

Kr Text evaluation (2026-09-22):

- The canonical exact path has exactly v0.1.0 and v0.2.0; `/v2` and `/v3`
  lines do not exist. Public enabled, unarchived, non-fork MIT repository
  `kr/text`, owned by `kr`, defaults to `main` and has no GitHub Releases.
  Neither stable has a Go directive, retraction, deprecation, or replacement.
  Proxy/sumdb, annotated tag, commit/tree, and Git regular-file identities
  agree. No branch, pseudo-version, alternate path, fork, or replacement was
  promoted.
- V0.1.0 requires kr/pty v1.1.1; v0.2.0 changes only its `mc` import and
  requirement to creack/pty v1.1.9. Both stables expose the root text API,
  colwriter, and `agg`/`mc` commands. Their complete closures preserve Go 1.18.
  Both SDKs pass upstream build, repeated tests, race, vet, positive bounded
  ordinary fixtures, and supported Darwin/Linux/FreeBSD cross-builds. Windows
  and js/wasm stop at the PTY-backed `mc` platform boundary.
- Neither stable qualifies. The documentation says a line can exceed its limit
  only when a single word exceeds it, but both stables join two adjacent
  over-limit words: `Wrap("overlong overlong foo", 4)` returns
  `"overlong overlong\nfoo"`. Exact Go 1.26.7 and Go 1.18.10 repeated and race
  runs reproduce the ordinary failure. Verified post-release commit
  `838204404ccb967534580e2a6efb24062880dd0f` merged the matching PR #9 fix,
  but no later stable exists, so that commit and `main` are ineligible.
- Exact MVS selects v0.2.0 through supported Cast v1.5.1. Historical kr/pretty
  v0.1.0/v0.2.0 request v0.1.0, whose unloaded `mc` command alone requests and
  imports kr/pty v1.1.1. The complete direct sergi/Assert,
  mvn-pom-mutator/Viper/Consul/Prometheus, and Honnef/go-internal/errgo route
  families reproduce. Target why is positive only through yaml.v2 tests/
  check.v1/kr/pretty; target and both PTY modules have zero repository imports
  and production/complete-test loads and no runtime reachability.
- A disposable exact v0.2.0 root adds creack/pty, giving 235 modules, 3,601
  edges, and 1,069 sums; tidy removes both manufactured roots and restores the
  exact base. Its raw graph/module state would expire a kr/pty guard. A
  disposable v0.1.0 root downgrades Cast, Viper, and kr/pretty, giving
  229/3,585/389/425/1,073; tidy retains broad selection and graph drift. It
  violates Cast, kr/pretty, and kr/pty guards. Neither projection was retained,
  and the real project received no exact get.
- The unchanged project remains 234/3,599/355/429/197/41/1,067 with exact
  module hashes and the recorded 432-line/common applied tidy state. All 38
  earlier guarded selections remain exact; their 234 incoming edges retain
  `d31b2afdb5a11dadac236eec600408cc9ba0ee3deca2793ae2adf7ffc449c320`.
  Thirty-seven why results are negative and only closed kr/pretty is positive;
  guarded imports/loads remain zero. The complete kr/pty exception remains
  exact in the real project.
- Exact target OSV/GitHub and repository advisory results are empty. Pinned
  isolated govulncheck v1.8.0 is empty for both stables; base/direct-v0.2.0
  project populations are identical at 30/22/20/20 with no target trace.
  Guard Gorilla/retryablehttp pairs, x/mod findings, the 1,402-record Go index,
  and PUBLISHED memberlist CNA identity remain exact. Advisory absence does not
  override behavior. Final exact-Go module verification, build, count-one,
  race count-one, and vet pass; accepted quality remains 27/27 Q0-Q2 PASS at
  L2. No source or dependency metadata changed. P7 stops for one reciprocal
  choice: exact target-specific unqualified v0.2.0 retention, one later
  measurement-only main -> Cast v1.5.1 -> kr/text v0.2.0 owner/request study,
  or stopping P7 unresolved. The decision successor was prepared but not
  executed; P8 remains queued.

Kr Text product decision (2026-09-22):

- Option 1 was selected. Exact selected, inherited, unloaded
  `github.com/kr/text v0.2.0` is explicitly retained without changing product
  source, `go.mod`, or `go.sum`. It remains unqualified and is not described
  as secure: both exact-path stables fail the completed adjacent-overlong-word
  wrapping contract, and the matching upstream fix is unreleased. Physical
  selection, the test-only why chain, zero loading, and advisory absence bound
  exposure but are not qualification.
- The kr/text-specific, non-transferable exception accepts only the completed
  two-release/repository/source, API/closure/platform/behavior, Cast test-
  closure ownership, current and historical owner/request routes, why/import/
  load/runtime, graph/tidy/projection/Go-floor, earlier-guard, kr/pty-guard,
  and advisory findings. It accepts no new or uncharacterized defect and
  transfers or broadens no exception.
- Retention requires exact v0.2.0; Cast v1.5.1 -> v0.2.0; historical kr/pretty
  v0.1.0/v0.2.0 -> v0.1.0 -> kr/pty v1.1.1; the complete Cast, sergi/Assert,
  mvn-pom-mutator/Viper/Consul/Prometheus, and Honnef/go-internal/errgo route
  families and requester imports; no direct target root; only the yaml.v2-
  tests/check.v1/kr/pretty why chain; negative kr/pty and creack/pty why; zero
  target/PTY imports and loads; and runtime unreachability.
- The exact two-release and repository/tag/commit/tree/archive/sumdb identities
  must remain unchanged; both releases must remain unqualified for the
  completed failure; the post-release fix must remain unreleased; and no
  qualified supported route may appear. Exact project/module/tidy/projection/
  Go-floor state and every earlier and kr/pty guard must also remain exact.
- Any target/request/requester-import/owner-route, root/why/import/load/runtime,
  graph/module/tidy/Go-floor, earlier or kr/pty guard, advisory/finding,
  repository/release/owner, qualification, supported-owner, or compatible-
  route change expires retention and requires a fresh kr/text dependency and
  product decision before merge. Any kr/pty guard change separately requires
  its fresh dependency and product decision. No owner study, direct root,
  unreleased commit, alternate path, dependency edit, workaround, unrelated
  selection, or implementation is authorized.
- Guard-only revalidation preserves 234/3,599/355/429/197/41/1,067, exact
  module and tidy identities, all target requests/routes, and the why/import/
  load boundary. All 38 earlier selections/234 incoming edges remain exact at
  `d31b2afdb5a11dadac236eec600408cc9ba0ee3deca2793ae2adf7ffc449c320`;
  including kr/text gives 39 selections/237 edges at
  `151e72c0b5444acc59ffab45d6ce8fe43e80821b42df6a49a11902b3a6632803`.
  Thirty-seven earlier why results are negative and only closed kr/pretty and
  kr/text are positive; guarded imports/loads remain zero.
- Fresh official Go 1.26.7, target advisory, guard OSV, x/mod, Go-index,
  memberlist CNA, proxy-release, and repository identities reproduce. No
  completed text or PTY behavior/source/closure/projection evaluation was
  repeated, and no owner study ran. Final exact-Go verify/build/count-one/
  race/vet pass; accepted quality remains 27/27 Q0-Q2 PASS at L2.
- P7 continues only with the prepared bounded product decision for selected
  exact-path `github.com/kyokomi/emoji/v2 v2.2.12`. It must preserve all 39
  earlier selections/237 edges and every kr/text and kr/pty expiry guard. The
  successor was prepared but not executed; P8 remains queued.

Kyokomi Emoji V2 dependency evaluation (2026-09-22):

- No canonical exact-path stable qualifies. The proxy's canonical resolvable
  `/v2` line is v2.2.5-v2.2.14; earlier v2.0.0-v2.2.4 tags lack the `/v2`
  module declaration and belong to the unsuffixed `+incompatible` line.
  `/v3` is absent. Public enabled, unarchived, non-fork MIT repository
  `kyokomi/emoji`, ID 21064634, is owned by `kyokomi`, defaults to `main`, and
  has no parent/source. Release/tag/commit/tree/signature, module, ancestry,
  proxy/sumdb, license, and regular-file archive/Git identities agree. There
  are no retractions, deprecations, replacements, or eligible alternate paths.
- V2.2.5-v2.2.13 declare Go 1.14 and have no requirements; latest v2.2.14
  declares Go 1.21 and has no requirements. Exact v2.2.13 is therefore the
  highest floor-compatible stable. The standard-library-only module contains
  one library and one example package, with no cgo, embed, build-tag, platform,
  or runtime network boundary. Generated maps and `ReplacePadding` are exposed
  mutable globals; callers own mutation, synchronization, writers, and errors.
- Upstream build, count-one/count-ten, race-count-ten, vet, and supported
  cross-build gates pass for owner-requested v2.2.8, selected v2.2.12, and
  highest floor-compatible v2.2.13 under exact Go 1.26.7 and Go 1.18.10;
  Go 1.26.7 gates pass v2.2.14. Positive bounded ordinary fixtures pass.
  Every canonical stable nevertheless merges adjacent `Fprintln` operands
  rather than matching documented `fmt.Fprintln`, and implements documented
  `Errorf` with a plain error so `%w` does not wrap as `fmt.Errorf` does. Both
  ordinary failures reproduce under both SDKs and race. V2.2.14 independently
  violates the retained Go floor.
- Exact requests are main -> v2.2.12 and direct supported
  go-term-markdown v0.1.4 -> v2.2.8. Main historically introduced v2.2.8
  with that requester and later raised only its target root. The requester
  imports the target in `renderer.go` and uses `emoji.Sprint`. Main directly
  imports the requester but not the target. Both requester and target have
  positive why results, are production/complete-test loaded, and are runtime
  relevant through main `cmd` -> requester -> target.
- Disposable v2.2.8, v2.2.12, v2.2.13, and v2.2.14 exact gets were not
  retained. V2.2.8 changes only the target root while preserving the base;
  v2.2.12 is a no-op; v2.2.13 changes only the target root and two sum lines
  for 234/3,599/355/429/197/41/1,069; v2.2.14 raises main `go 1.18` to
  `go 1.21` and produces 234/3,601/355/429/197/41/1,069. Comparable tidy
  projections preserve all kr/text and kr/pty guards. The real project remains
  234/3,599/355/429/197/41/1,067 with exact module/common tidy hashes.
- All 39 earlier selections and 237 incoming edges remain exact at
  `151e72c0b5444acc59ffab45d6ce8fe43e80821b42df6a49a11902b3a6632803`;
  37 why results are negative and only closed kr/pretty and kr/text are
  positive, with zero guarded imports/loads. Neither kr/text nor kr/pty expiry
  guard changed.
- Exact target OSV/GitHub/repository results are empty. Pinned isolated
  govulncheck is empty for selected/highest/latest; base and v2.2.13 project
  populations match at 30/22/20/20 with no target finding. Guard OSV, x/mod,
  the exact 1,402-record Go index, and PUBLISHED memberlist CNA identities
  remain unchanged. Final exact-Go verify/build/count-one/race/vet pass;
  accepted quality remains 27/27 Q0-Q2 PASS at L2.
- Product source and dependency metadata remain unchanged. P7 stops for one
  reciprocal bounded product choice: exact target-specific unqualified
  v2.2.12 retention; one later measurement-only study of the existing main ->
  direct go-term-markdown v0.1.4 -> target owner/request route; or stopping P7
  unresolved. Do not repeat the evaluation, run the study during the decision,
  add a root, change either closed exception, combine another group, or begin
  P8. The decision successor was prepared but not executed.

Kyokomi Emoji V2 product decision (2026-09-22):

- Option 1 was selected. Exact selected, direct-indirect, production-loaded,
  runtime-relevant `github.com/kyokomi/emoji/v2 v2.2.12` is explicitly
  retained without changing product source, `go.mod`, or `go.sum`. It remains
  unqualified and is not described as secure: all ten canonical exact-path
  stables fail the completed documented Fprintln and Errorf contracts, and
  v2.2.14 independently violates the Go 1.18 floor. Physical selection,
  positive why, runtime loading, and advisory absence are not qualification.
- The emoji/v2-specific, non-transferable exception accepts only the completed
  ten-release/repository/source, ordinary behavior, API/closure/platform,
  go-term-markdown importer/use and route, why/import/load/runtime, graph/tidy/
  projection/Go-floor, earlier-guard, kr/text/kr/pty-guard, and advisory
  findings. It accepts no new or uncharacterized defect and transfers or
  broadens no exception.
- Retention requires exact v2.2.12 and main's direct-indirect request;
  go-term-markdown v0.1.4 -> v2.2.8; the historical main v2.2.8 -> v2.2.12
  target-root-only change; `renderer.go` target import and `emoji.Sprint` use;
  the main `cmd` -> direct requester -> target route; no direct repository
  target import; positive target/requester why; production/complete-test
  loading of both modules; and target runtime relevance.
- The exact canonical v2.2.5-v2.2.14 line, separation from older unsuffixed
  tags, absent `/v3`, repository/owner/status/license/default-branch and
  release/tag/commit/tree/archive/sumdb identities, nine canonical Release
  objects, and no replacement/retraction/deprecation or eligible alternate
  must remain exact. Every stable must remain unqualified for the completed
  failures, v2.2.14 must remain floor-ineligible, and no future qualified
  Go-1.18-compatible stable or supported route may appear.
- Exact base 234/3,599/355/429/197/41/1,067, real module hashes, the 432-line
  and common 52/948-line tidy identities, every completed disposable
  projection, exact Go 1.18/1.26.7, source/API/CLI/help/launcher/Make/quality
  contracts, accepted 27/27 Q0-Q2 PASS at L2, all 39 earlier selections/237
  edges, and every kr/text and kr/pty expiry guard must remain exact. Including
  emoji/v2 gives 40 selections/239 incoming edges at
  `deb453cac4dcb6df8a6b502f6f20aba72eaeb415ac2ff40f873a677113cf9d4a`.
- Any target/request/requester-import/owner-route, root/why/import/load/runtime,
  graph/module/tidy/Go-floor, earlier or kr/text/kr/pty guard, advisory/finding,
  repository/release/owner, qualification, supported-owner, or compatible-
  route change expires retention and requires a fresh emoji/v2 dependency and
  product decision before merge. Any kr/text or kr/pty change separately
  requires its own fresh decision. No owner study, root change, branch/pseudo-
  version, alternate path, dependency edit, workaround, unrelated selection,
  implementation, or transferred exception is authorized.
- Guard-only revalidation reproduced exact continuity, module/tidy, request/
  route/why/import/load, 39-selection/237-edge and 40-selection/239-edge,
  closed-exception, official Go 1.26.7, proxy/repository/release, target/
  requester/guard advisory, Go-index, memberlist CNA, and final exact-Go gate
  identities. Completed behavior/source/closure/upstream/projection/archive/
  govulncheck work was not repeated and no owner study ran. No source or
  dependency metadata changed; accepted quality remains 27/27 Q0-Q2 PASS at
  L2.
- P7 continued with the bounded evaluation of selected direct exact
  `github.com/manifoldco/promptui v0.9.0`; that evaluation is now complete and
  its product decision is the sole prepared successor. All 40 guarded
  selections/239 edges and the emoji/v2, kr/text, and kr/pty exceptions remain
  exact. P8 remains queued.

Manifoldco Promptui evaluation (2026-09-22):

- No canonical exact-path stable qualifies. Exact v0.9.0 is the latest stable
  and declares Go 1.12, but it discards an ordinary caller-provided output
  error despite Prompt.Run's documented execution-error result and races
  internally during a single supported Prompt run. Both failures reproduce
  under exact Go 1.26.7 and Go 1.18.10. V0.3.2-v0.8.0 share them;
  v0.1.0-v0.3.1 have no release go.mod and their synthesized one-line proxy
  module files cannot resolve a standalone release build/test closure.
- Exact go-import/module identity resolves the public enabled, unarchived,
  non-fork BSD-3-Clause `manifoldco/promptui` repository, ID 107154646,
  defaulting to master. The twelve-release exact stable line is v0.1.0-
  v0.9.0; `/v2` and `/v3` are absent. There are no retractions, deprecations,
  prereleases, or eligible replacements, alternate paths, forks, branches, or
  pseudo-versions. All lightweight tags form continuous ancestry; nine have
  GitHub Release objects; every proxy ZIP byte-matches its exact Git regular-
  file manifest. Selected v0.9.0 is unsigned commit
  `c2e487d3597f59bcf76b24c9e80679740a72212b`, tree
  `becf36d02c10091c82f8eeccba0e67bead88d418`.
- V0.9.0 has root/list/screenbuf packages plus examples, mutable exported
  globals owned by callers, and readline/caller streams as its terminal
  boundary. It has no cgo, generated, embed, subprocess, or network boundary.
  Both exact SDKs pass module verification, build, upstream count-one/count-
  ten/race-count-ten/vet, and supported Darwin/Linux/Windows/FreeBSD cross
  gates. Positive ordinary input/validation/determinism behavior passes, but
  the output-error and single-Prompt-race failures disqualify it.
- The sole target request is main -> v0.9.0; main's direct readline v1.5.1
  wins over Promptui's old readline pseudo-version request. Main history
  records direct v0.8.0 addition/removal and v0.9.0 reintroduction before the
  current source import and prompt functions. `cmd/root.go` is the sole target
  import, configures PromptTemplates, and calls Prompt.Run for value and
  confirmation input. Target why is positive; root/list/screenbuf are
  production and complete-test loaded, and the dependency is runtime relevant.
- A disposable exact v0.9.0 get is a byte no-op. The real project remains
  234/3,599/355/429/197/41/1,067 with main Go 1.18, exact module hashes, the
  432-line tidy identity, and common 52/948-line applied tidy hashes. All 40
  earlier guarded selections/239 incoming edges remain exact at
  `deb453cac4dcb6df8a6b502f6f20aba72eaeb415ac2ff40f873a677113cf9d4a`;
  every emoji/v2, kr/text, and kr/pty expiry guard is unchanged. No projection
  or implementation was retained.
- Exact target and both requested/selected readline OSV plus GitHub global/
  repository advisory results are empty. Pinned govulncheck is empty for
  isolated target scans; base project populations remain 30/22/20/20 with no
  target trace. Guard OSV, x/mod, the
  1,402-record Go index, and PUBLISHED memberlist CNA identities remain exact.
  Final exact-Go verify/build/count-one/race/vet pass; accepted quality remains
  27/27 Q0-Q2 PASS at L2. Advisory absence is not qualification.
- Product source and dependency metadata remain unchanged. P7 stops for one
  reciprocal bounded product choice: exact target-specific unqualified v0.9.0
  retention; one later measurement-only `Ply Direct Promptui Ownership Study`
  of the existing main -> direct target request and `cmd/root.go` runtime
  route; or stopping P7 unresolved. Do not repeat the evaluation, run the
  study during the decision, transfer an exception, combine another group, or
  begin P8. The decision successor was prepared but not executed.

Manifoldco Promptui product decision (2026-09-22):

- Option 1 was selected. Exact selected, direct, production-loaded, runtime-
  relevant `github.com/manifoldco/promptui v0.9.0` is explicitly retained
  without changing product source, `go.mod`, or `go.sum`. It remains
  unqualified: no canonical exact-path stable qualifies. Directness, positive
  why, runtime loading, and advisory absence are not qualification.
- The Promptui-specific, non-transferable exception accepts only the completed
  twelve-release/repository/source, ordinary behavior, API/closure/platform,
  owner/request, `cmd/root.go` route/use, why/import/load/runtime, graph/tidy/
  Go-floor, earlier-guard, closed-exception, and advisory findings. It accepts
  no new defect or advisory and transfers or broadens no exception.
- Retention requires exact v0.9.0 and main's sole direct request; Promptui's
  four exact requirements; main's direct readline v1.5.1 MVS selection; the
  historical v0.8.0 request/removal and v0.9.0 reintroduction; sole import,
  PromptTemplates and two Prompt.Run uses; positive main `cmd` -> target route;
  three production/complete-test target packages; and runtime relevance.
- It also requires the exact v0.1.0-v0.9.0 line, absent `/v2` and `/v3`, exact
  repository/owner/status/license/default-branch and completed release/tag/
  commit/tree/signature/archive/sumdb/module identities, nine Release objects,
  no replacement/retraction/deprecation or eligible alternate, every stable
  remaining unqualified for the completed closure/output-error/race results,
  and no future qualified Go-1.18-compatible stable or supported route.
- Exact base 234/3,599/355/429/197/41/1,067, real module hashes, the 432-line
  and common 52/948-line tidy identities, no-op v0.9.0 projection, exact Go
  1.18/1.26.7, source/API/CLI/help/launcher/Make/quality contracts, accepted
  27/27 Q0-Q2 PASS at L2, all 40 earlier selections/239 edges, and every
  emoji/v2, kr/text, and kr/pty expiry guard must remain exact. Including
  Promptui gives 41 selections/240 incoming edges at
  `a4e6042216821077b7a74472a6fcdc60d70644d07a805f6a6e99db1140a71563`.
- Any target/request/import/owner-route, root/why/import/load/runtime, graph/
  module/tidy/Go-floor, earlier or closed-exception guard, advisory/finding,
  independent defect, repository/release/owner, qualification, supported-
  owner, or compatible-route change expires retention and requires a fresh
  Promptui decision before merge. Any emoji/v2, kr/text, or kr/pty change
  separately requires its own fresh decision. No owner study, root change,
  branch/pseudo-version, alternate path, dependency edit, workaround,
  unrelated selection, implementation, or transferred exception is
  authorized.
- Guard-only revalidation reproduced exact continuity, module/tidy, request/
  route/why/import/load, 40-selection/239-edge and 41-selection/240-edge,
  closed-exception, official Go 1.26.7, proxy/repository/release, target/
  readline/guard advisory, Go-index, memberlist CNA, and final exact-Go gate
  identities. Completed behavior/source/closure/upstream/projection/archive/
  govulncheck work was not repeated and no owner study ran. P7 continues only
  with the prepared bounded go-homedir evaluation; P8 remains queued.

Mitchellh Go Homedir evaluation (2026-09-22):

- No canonical exact-path stable qualifies. The proxy line contains exactly
  v1.0.0 and v1.1.0; v1.1.0 is latest and Go-1.18-compatible. Exact `/v2`
  and `/v3` lines are absent. Public enabled archived non-fork MIT repository
  `mitchellh/go-homedir`, ID 23121609, is owned by `mitchellh`, defaults to
  `main`, and has no parent/source. There are no GitHub Releases, retractions,
  deprecations, replacements, prereleases, or eligible alternate paths.
- Both lightweight tags form continuous ancestry and both proxy ZIPs
  byte-match exact Git regular-file manifests. V1.1.0 is validly signed commit
  `af06845cf3004701891bf4fdb884bfe4920b3727`, tree
  `c70b446f839c30d92fb3bbb491f0f0b540f47529`; v1.0.0's old signature key is
  not verifiable by GitHub. Both modules declare only the exact path, with no
  Go directive, requirement, replacement, or retraction.
- Both standard-library-only releases pass exact Go 1.26.7 and Go 1.18.10
  verify/build/vet and supported Darwin/Linux/Windows/FreeBSD/Plan 9/js-wasm
  cross gates. Small ordinary HOME-backed Dir/Expand/cache/Reset/concurrency
  fixtures pass repeated and race runs under both SDKs. Both releases fail the
  required upstream count-one/count-ten/race-count-ten gate: TestDir unsets
  HOME, Darwin `dscl` returns `eServerError`, the shell fallback cannot `cd`,
  and Dir returns `exit status 1`. The documented API permits a discovery
  error, but the complete upstream gate is mandatory for qualification.
- Seven graph requests select the target. Main has directly requested v1.1.0
  since its initial module commit; mvn-pom-mutator v0.2.3, Viper v1.15.0 and
  historical v1.10.1, go-rootcerts v1.0.2, and crypt v0.4.0 also request
  v1.1.0, while historical go-rootcerts v1.0.0 requests v1.0.0.
  `pkg/config/profiles.go` is the sole production target import and calls Dir;
  three tests import Reset. Target why is positive, one package is production/
  complete-test loaded, and the dependency is runtime relevant.
- A disposable exact v1.1.0 get is a byte no-op. The real project remains
  234/3,599/355/429/197/41/1,067 with main Go 1.18, exact module hashes, the
  432-line tidy identity, and common 52/948-line applied tidy identities. An
  exact v1.0.0 request removes direct mvn-pom-mutator, downgrades Viper,
  shrinks the graph to 181 modules/2,551 edges, breaks project loading, and is
  restored to v1.1.0 by tidy. No projection or implementation was retained.
- All 41 earlier guarded selections and 240 incoming edges remain exact at
  `a4e6042216821077b7a74472a6fcdc60d70644d07a805f6a6e99db1140a71563`.
  Thirty-seven why results are negative; only kr/pretty, kr/text, emoji/v2,
  and Promptui are positive. Promptui is the only guarded repository import;
  only emoji/v2 and Promptui are loaded. Every Promptui, emoji/v2, kr/text,
  and kr/pty expiry guard remains exact.
- Exact target OSV/GitHub/repository results are empty. Exact-Go pinned
  govulncheck is empty for both isolated releases; base and no-op-v1.1.0
  project populations match at 30/22/20/20 with no target trace. Guard OSV,
  x/mod, the 1,402-record Go index, and PUBLISHED memberlist CNA identities
  remain unchanged. Advisory absence does not override the upstream gate.
- Product source and dependency metadata remain unchanged. Final exact-Go
  verify/build/count-one/race/vet pass; accepted quality remains 27/27 Q0-Q2
  PASS at L2. Every task-owned scratch artifact was contained beneath the
  managed root and removed; only the pre-existing launcher-owned Node compile
  cache remains. P7 stops for one reciprocal bounded product choice: exact
  target-specific unqualified v1.1.0 retention; one later measurement-only
  `Ply Direct Go Homedir Ownership Study` of the existing main -> direct
  target request and `pkg/config/profiles.go` runtime route; or stopping P7
  unresolved. Do not repeat the evaluation, transfer an exception, combine
  another group, or begin P8. The decision successor was prepared but not
  executed.

Mitchellh Go Homedir product decision (2026-09-22):

- Option 1 was selected. Exact selected, direct, production-loaded, runtime-
  relevant `github.com/mitchellh/go-homedir v1.1.0` is explicitly retained
  without changing product source, `go.mod`, or `go.sum`. It remains
  unqualified: no canonical exact-path stable qualifies. Directness, positive
  why, runtime loading, and advisory absence are not qualification.
- The go-homedir-specific, non-transferable exception accepts only the
  completed two-release/repository/source, ordinary behavior, API/closure/
  platform, exact owner/request, `pkg/config/profiles.go` route/use, three
  Reset test imports, why/import/load/runtime, graph/tidy/projection/Go-floor,
  earlier-guard, closed-exception, and advisory findings. It accepts no new
  defect or advisory and transfers or broadens no exception.
- Retention requires exact v1.1.0, main's direct request and the other six
  current/historical requests, every current requester selection, the direct
  root since the initial module commit, sole production import and Dir use,
  three Reset imports, positive main `pkg/config` -> target route, one
  production/complete-test target package, runtime relevance, and all
  requester import/why/load/owner-route boundaries.
- It also requires the exact v1.0.0-v1.1.0 line, absent `/v2` and `/v3`, exact
  repository/owner/status/license/default-branch and completed release/tag/
  commit/tree/signature/archive/sumdb/module identities, zero GitHub Releases,
  no replacement/retraction/deprecation or eligible alternate, both stables
  remaining unqualified for the completed upstream-suite result, and no
  future qualified Go-1.18-compatible stable or supported route.
- Exact base 234/3,599/355/429/197/41/1,067, real module hashes, the 432-line
  and common 52/948-line tidy identities, no-op v1.1.0 and rejected v1.0.0
  projections, exact Go 1.18/1.26.7, source/API/CLI/help/launcher/Make/
  quality contracts, accepted 27/27 Q0-Q2 PASS at L2, all 41 earlier
  selections/240 edges, and every Promptui, emoji/v2, kr/text, and kr/pty
  expiry guard must remain exact. Including go-homedir gives 42 guarded
  selections and 247 incoming edges at
  `4caf9e803c1b324d2afcd17b9ba664e6d2c8d220cd74b264484107ea3dd93f42`.
- Any target/request/requester-import/owner-route, root/why/import/load/runtime,
  graph/module/tidy/Go-floor, source/behavior/closure, earlier or closed-
  exception guard, advisory/finding, independent defect, repository/release/
  owner, qualification, supported-owner, or compatible-route change expires
  retention and requires a fresh go-homedir decision before merge. Any
  Promptui, emoji/v2, kr/text, or kr/pty change separately requires its own
  fresh decision. No owner study, root change, lower selection, fork,
  replacement, alternate path, dependency edit, workaround, unrelated
  selection, implementation, or transferred exception is authorized.
- Guard-only revalidation reproduced exact continuity, module/tidy, request/
  route/why/import/load, 41-selection/240-edge and 42-selection/247-edge,
  closed-exception, official Go 1.26.7, proxy/repository/release, target/guard
  advisory, Go-index, memberlist CNA, and final exact-Go gate identities.
  Completed behavior/source/closure/upstream/cross-build/projection/archive/
  govulncheck work was not repeated and no owner study ran. P7 continues only
  with the prepared bounded mapstructure v1.5.0 evaluation; P8 remains queued.

Mitchellh Mapstructure evaluation (2026-09-22):

- No canonical exact-path stable qualifies, and no stable has a genuine
  supported exact-path owner. The proxy line contains exactly 17 stables from
  v1.0.0 through v1.5.0; v1.5.0 is latest and Go-1.18-compatible. V1.0.0-
  v1.1.2 omit a Go directive and v1.2.0-v1.5.0 declare Go 1.14. Exact `/v2`
  and `/v3` lines, GitHub Releases, retractions, deprecations, replacements,
  and prereleases are absent.
- Public enabled archived non-fork MIT repository `mitchellh/mapstructure`,
  ID 10166531, remains the exact go-import owner. All 17 lightweight tags
  form continuous ancestry; the later 16 commits are validly signed while
  v1.0.0 has an unverifiable old key. Selected v1.5.0 is signed commit
  `ab69d8d93410fce4361f4912bb1ff88110a81311`, tree
  `5a1166013faa55170e1acce958c1a3863faf3f80`, and its proxy ZIP byte-matches
  its exact Git manifest. V1.2.0 is the sole historical tag-mutation record:
  its immutable proxy archive matches historical commit `047abd31...`, while
  the current signed tag points to descendant `9e401191...` and differs only
  by two CHANGELOG lines. Selected v1.5.0 is unaffected.
- The owner archive notice says there will be no more tags and names active
  fork `go-viper/mapstructure` as blessed. That fork declares alternate
  module path `github.com/go-viper/mapstructure/v2` with Go 1.18, so it is
  neither an exact-path stable nor a supported owner for the exact line.
- Every stable passes exact Go 1.26.7/Go 1.18.10 verification, build,
  count-one tests, race count-one tests, vet, and ten Darwin/Linux/Windows/
  FreeBSD/Plan 9/js-wasm cross-builds. V1.2.2 additionally fails repeated
  upstream metadata-order assertions under both SDKs. Selected ordinary
  deterministic/error/concurrency fixtures pass repeated and race runs.
- The decisive bounded ordinary slice-replacement fixture fails identically
  at all 17 stables under both SDKs: decoding one source element into a
  pre-populated two-element destination leaves its stale second element.
  Upstream merged the two-line truncation fix in commit `33d262e...` three
  days after v1.5.0 and records it for unreleased v1.5.1. Later branch fixes
  are likewise unreleased and ineligible. No lower stable is a fallback.
- Nine exact requests select the target: main/Viper v1.15.0 -> v1.5.0;
  historical Viper v1.10.1/crypt v0.4.0 -> v1.4.3; Kong plus Consul API
  v1.1.0/v1.12.0 -> v1.1.2; and Serf v0.8.2/v0.9.6 -> the recorded 2016
  pseudo-version. Main marks the request indirect and has no target import.
  Target why is positive through `cmd -> Viper -> target`; one target package
  is production/complete-test loaded. Viper, Consul, and Serf have genuine
  historical imports, while Kong and crypt are metadata-only requesters.
- A disposable exact v1.5.0 get is a byte no-op. The real project remains
  234/3,599/355/429/197/41/1,067 with exact module hashes and the 432-line
  tidy identity. All 42 earlier guarded selections and 247 incoming edges
  remain exact at
  `4caf9e803c1b324d2afcd17b9ba664e6d2c8d220cd74b264484107ea3dd93f42`.
  Every go-homedir, Promptui, emoji/v2, kr/text, and kr/pty expiry guard
  remains exact. No lower projection or dependency implementation was
  retained.
- Exact OSV/GitHub/repository results are empty for all 17 stables. Pinned
  govulncheck is empty for isolated targets; base/no-op project populations
  match at 30/22/20/20 with no target trace. Guard OSV, x/mod, the exact
  1,402-record Go index, and PUBLISHED memberlist CNA identities remain
  unchanged. Advisory absence does not override the ordinary failure or
  unsupported exact-path owner.
- Product source and dependency metadata remain unchanged. Final exact-Go
  verify/build/count-one/race/vet pass; accepted quality remains 27/27 Q0-Q2
  PASS at L2. P7 stops for one reciprocal bounded product choice: exact
  mapstructure-specific unqualified v1.5.0 retention; one later measurement-
  only `Ply Viper Mapstructure Ownership Study` of the existing indirect
  request, Viper import/load route, historical requester provenance, and
  blessed alternate-path owner; or stopping P7 unresolved. Do not repeat the
  evaluation, transfer an exception, combine another group, or begin P8. The
  decision successor was prepared but not executed. Every task-owned scratch
  artifact was contained beneath the managed session root and removed; only
  the pre-existing launcher-owned Node compile cache remains.

Mitchellh Mapstructure product decision (2026-09-22):

- Option 1 was selected. Exact selected indirect, production-loaded
  `github.com/mitchellh/mapstructure v1.5.0` is explicitly retained without
  changing product source, `go.mod`, or `go.sum`. It remains unqualified: no
  canonical exact-path stable qualifies and no stable has a genuine supported
  exact-path owner. Indirectness, positive why, production loading, and
  advisory absence are not qualification.
- The mapstructure-specific, non-transferable exception accepts only the
  completed 17-release/repository/source, ordinary behavior, API/closure/
  platform, exact owner/request, Viper route, why/import/load/runtime, graph/
  tidy/projection/Go-floor, earlier-guard, closed-exception, and advisory
  findings. It accepts no new defect or advisory and transfers or broadens no
  exception.
- Retention requires exact v1.5.0, main's indirect request and the other eight
  current/historical requests, every current requester selection, every
  requester import or metadata-only boundary, zero repository imports, the
  positive `cmd -> Viper -> target` route, Viper production/test imports,
  project Viper config calls, one production/complete-test target package,
  and the recorded load/runtime-use boundary.
- It also requires the exact 17-release line, absent exact-path `/v2` and
  `/v3`, archived exact repository/owner/status/license/default-branch,
  blessed alternate-path fork, and completed release/tag/commit/tree/
  signature/archive/sumdb/module identities; the v1.2.0 tag-mutation record;
  no replacement/retraction/deprecation or eligible exact-path alternate; all
  17 stables remaining unqualified for the stale-slice result; the v1.2.2
  repeated-suite result; and no future qualified exact-path stable, supported
  exact-path owner, or compatible genuine route.
- Exact base 234/3,599/355/429/197/41/1,067, real module hashes, the 432-line
  and common 52/948-line tidy identities, no-op v1.5.0 projection, exact Go
  1.18/1.26.7, source/API/CLI/help/launcher/Make/quality contracts, accepted
  27/27 Q0-Q2 PASS at L2, all 42 earlier selections/247 edges, and every
  go-homedir, Promptui, emoji/v2, kr/text, and kr/pty expiry guard must remain
  exact. Including mapstructure gives 43 selections and 256 incoming edges at
  `8f05ff3b4755e6582db200d1d457b6e8b52e84983d9f634976f76b0d731d298f`.
- Any target/request/requester-import/owner-route, root/why/import/load/runtime,
  graph/module/tidy/Go-floor, source/behavior/closure, repository/release/
  owner, advisory/finding, qualification, supported-owner, earlier guard,
  closed-exception, or compatible-route change expires retention and requires
  a fresh mapstructure dependency and product decision before merge. Any
  other closed-exception change separately requires its corresponding fresh
  decision. This exception is not a security claim and cannot transfer.
- Guard-only revalidation reproduced exact continuity, Go/module/tidy,
  request/route/import/load, 42-selection/247-edge and 43-selection/256-edge,
  closed-exception, owner/release, narrow advisory, Go-index, memberlist CNA,
  and final exact-Go gate identities. Completed source/behavior/closure/
  upstream/cross-build/projection/archive/govulncheck work was not repeated
  and no owner study ran. P7 continues only with the prepared bounded
  modern-go/concurrent evaluation; P8 remains queued.

Modern-go/concurrent evaluation (2026-09-22):

- Upstream `1.0.3` is the highest qualified genuine exact-path stable.
  Selected pseudo-version
  `v0.0.0-20180306012644-bacd9c7ef1dd` is its exact commit and byte-identical
  source, so product source, `go.mod`, `go.sum`, all selections, and the
  transitive closure remain unchanged. The pseudo-version itself is not
  promoted as stable.
- Public enabled unarchived non-fork Apache-2.0 repository
  `modern-go/concurrent`, ID 123229061, defaults to `master` at the exact
  selected commit. Its four genuine GitHub Releases/lightweight tags are
  named `1.0.0`-`1.0.3` without Go-semver `v` prefixes. The proxy therefore
  has no stable list; Go resolves each release name to the corresponding
  pseudo-version. The four unsigned release commits form continuous ancestry,
  and each proxy archive byte-matches its Git regular-file content. Synthetic
  module files contain only the exact module directive. No `/v2` or `/v3`,
  replacement, retraction, deprecation, prerelease, or eligible alternate
  exists.
- Every release passes exact Go 1.26.7 and Go 1.18.10 module verification,
  build, count-one/repeated upstream tests, race, vet, four supported cross-
  builds, and small ordinary Map/executor fixtures. Release 1.0.3 also passes
  repeated/race per-executor panic-callback and documented `runtime.Goexit`
  cleanup fixtures. The standard-library-only package has no cgo, generated,
  embed, network, or filesystem boundary. Release 1.0.3 qualifies with its
  genuine supported unarchived exact-path owner.
- Eight graph edges request the target. Bketelsen/crypt, historical
  sagikazarmark/crypt, Prometheus client v1.0.0, and both Viper vertices are
  metadata-only. Json-iterator v1.1.9/v1.1.11/v1.1.12 genuinely use
  `concurrent.Map`, but json-iterator is not loaded. Selected Viper is the only
  loaded requester and does not import the target. Repository imports, target
  why, production/complete-test/module-backed loads, and runtime relevance
  remain negative. The answered archive records every shortest graph route;
  selected crypt v0.9.0 also keeps a redundant metadata-only target
  requirement omitted from the pruned main graph.
- A disposable exact `@1.0.3` get preserves the selected module set but would
  add an unauthorized main indirect root, one source sum, and one graph edge;
  its tidy projection removes the promotion. Because no selection changed,
  the selection-only implementation contract does not authorize those edits.
  The projection was not retained and changed closure is empty.
- Exact base 234/3,599/355/429/197/41/1,067, module hashes, the 432-line and
  common 52/948-line tidy identities, Go 1.18 floor, source/API/CLI/help/
  launcher/Make/quality contracts, and accepted 27/27 Q0-Q2 PASS at L2 remain
  exact. All 43 earlier guarded selections and 256 incoming edges remain at
  `8f05ff3b4755e6582db200d1d457b6e8b52e84983d9f634976f76b0d731d298f`.
  Including concurrent gives 44 guards and 264 incoming edges at
  `4a83d4e3a015f93da4b4d35bd590137071b01c40315ec9bb9c987688004f9b09`.
  Every closed-exception boundary remains exact.
- All four exact target OSV/GitHub results and target repository advisories
  are empty. Pinned govulncheck is empty for each isolated release; unchanged
  project populations remain 30/22/20/20 with no target trace. Guard OSV,
  x/mod, the 1,402-record Go index, and PUBLISHED memberlist CNA identities
  remain exact. Advisory absence was not used as qualification.
- Final required-umask exact-Go verify/build/count-one/race/vet passes. No
  dependency implementation commit exists. Every task-owned scratch artifact
  is removed before handoff; only the pre-existing launcher-owned Node compile
  cache remains. P7 continues only with the prepared bounded evaluation of
  selected `github.com/modern-go/reflect2 v1.0.2`; it was not executed. P8
  remains queued.

Modern-go/reflect2 evaluation (2026-09-22):

- Canonical `v1.0.2` is the highest qualified Go-1.18-compatible exact-path
  stable and is already the exact graph selection. Public enabled unarchived
  non-fork Apache-2.0 repository `modern-go/reflect2`, ID 123239987, is the
  genuine supported owner. Product source and dependency metadata remain
  unchanged; the minimal changed closure is empty.
- The owner published genuine lightweight-tag/GitHub releases `1.0.0`,
  `1.0.1`, `v1.0.1`, and `v1.0.2`. Only the `v`-prefixed releases appear in
  the proxy stable list; the legacy releases resolve to exact pseudo-versions.
  Their commits form continuous ancestry, the two canonical release commits
  are validly signed, and every proxy ZIP exactly byte-matches its Git regular-
  file set. V1.0.2 declares Go 1.12 with no requirements. No `/v2` or `/v3`,
  replacement, retraction, deprecation, prerelease, or eligible alternate
  owner exists.
- The owner moved tests to `modern-go/reflect2-tests`. Legacy 1.0.1 and
  canonical v1.0.1 reproducibly fault in ordinary unsafe map iteration under
  both exact SDKs. V1.0.2 passes count-one/count-ten upstream tests and vet;
  its independent legacy test harness is checkptr-incompatible under race,
  while a standard-library-only target fixture passes repeated/race tests and
  vet under exact Go 1.26.7 and Go 1.18.10. Direct build/race passes under
  both SDKs; Go 1.18 vet's sole intentional `NoEscape` unsafeptr diagnostic is
  isolated, and all other analyzers pass. Nine supported cross-build targets
  pass under both SDKs.
- The selected standard-library-only package has reflection, unsafe,
  architecture assembly, Go-version constraints, runtime/reflect linknames,
  concurrency-safe caches, and documented caller-owned unsafe pointer
  correctness. It has no cgo, generated, embed, network, filesystem, retained-
  goroutine, resource-lifecycle, or external-state boundary.
- Nine graph edges request the target. Json-iterator v1.1.9/v1.1.11/v1.1.12
  and etcd client/v2 v2.305.1 genuinely import it; Viper v1.10.1/v1.15.0,
  both crypt requesters, and Prometheus client v1.0.0 are metadata-only.
  Selected Viper is the only loaded requester and does not import reflect2.
  Repository imports, target why, production/complete-test/module-backed
  loads, and current runtime relevance remain negative.
- A disposable exact v1.0.2 get preserves the selection and all guards but
  would add an unauthorized indirect main root, source sum, and graph edge.
  Tidy removes that promotion and restores the common 52/948-line projection,
  so no projection was retained. Exact project counts, module hashes, 432-line
  tidy identity, Go floor, and every source/API/CLI/help/launcher/Make/quality
  contract remain unchanged.
- All 44 earlier guarded selections and 264 incoming edges remain exact at
  `4a83d4e3a015f93da4b4d35bd590137071b01c40315ec9bb9c987688004f9b09`.
  Including qualified reflect2 gives 45 selections and 273 incoming edges at
  `7e5c820da743ac628fb18da129b4d91428a36fd374ed248760d1318f2694c525`.
  Every closed-exception and concurrent guard remains exact.
- All exact target OSV/GitHub/repository results and isolated govulncheck
  results are empty. Base/projection govulncheck populations remain identical
  at 30/22/20/20 with no target trace; guard OSV, x/mod, the exact 1,402-record
  Go index, and PUBLISHED memberlist CNA identities remain unchanged.
  Advisory absence was not used as qualification.
- Final required-umask exact-Go verify/build/count-one/race/vet passes. No
  dependency implementation commit exists. P7 continues only with the
  prepared bounded evaluation of selected
  `github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223`;
  it was not executed. P8 remains queued.

- Keep the selected Go 1.26.7 declarations and exact baseline identity aligned;
  reconsidering the Go line requires a separate measured toolchain move.
- Upgrade dependencies in small groups, with `go mod tidy`, build, tests, race,
  vet/lint, API/CLI diff, acceptance, and vulnerability scanning after each.
- Keep dependency-only commits separate from behavior changes.
- The exact-toolchain baseline migration is complete: old/new instruments were
  reproduced with every comparable debt value preserved and both identities
  recorded.

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
root="$(mktemp -d "${CODEX_SESSION_SCRATCH_ROOT:?}/hermetic.XXXXXX")"
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
