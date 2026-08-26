# Quality Upgrade Plan

Last measured checkpoint: 2026-08-26, commit `1bce06f`.

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

Status: active. Moves 1 through 63 are complete.

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
Continue with only the private browser-launcher production selector: make its
complete process dependency use `process.SystemRunner()` while preserving
runtime platform selection, the exact asynchronous commands and URLs for
Linux, Windows, and macOS, direct start errors, unsupported-platform behavior,
safe zero behavior, and public `OpenBrowser` behavior. Leave server behavior,
Maven's completed full dependency and stdout selection, `process.System()`,
profile, shell Run and Git, plugin diagrams, Unzip, clock/server work,
inventory, scanner, and P4-P8 unchanged.

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
