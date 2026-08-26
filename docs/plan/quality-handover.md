# Quality Upgrade Handover

Generated: 2026-08-26T18:01:11+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P4.8 product implementation:
  `67344a800db0030e1aceb908562eef8c8153394b`.
- Preceding P4.7 product implementation:
  `4da64d2468307d39ed3e8303ad9183fa18d71a75`.
- Truthful P3.54 audit-apparatus repair:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  the final commit changes only planning, archives, and launcher mutable state.
- No push, merge, release, publication, distribution, stash, revert, successor
  launch, or worktree removal was performed.

All earlier focused P3 and P4 implementation commits and their clean
checkpoints remain recorded in `docs/plan/quality-upgrade.md`. The separate
operational continuity implementation is `1b85711`, and the focused T15 repair
is `ce736a2`; neither changes a Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P4 is active and P5-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T180111+0200-cover-context-package.md`.
The P4.8 logger archive is answered history and links reciprocally to that
archive. The graph has exactly one NEXT tail. The archived prompt SHA-256 is
`d15a436b770b306c5003da989c238153a0822b4b2a4b632e56265bca8ecc676a`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P4 Move 8

Product commit `67344a8` adds only `pkg/logger/logger_test.go`. Both production
logger files, exported signatures, exact package globals and identities,
private logger and collector-hook arrangement, standard-logrus split, levels,
fields, formatter, output selection, callers, logging bytes, and observable
behavior remain unchanged. Six focused contracts bring the suite to 440 tests
across 26 of 27 packages and cover every logger statement, including package
initialization.

The contracts pin the exact ordered `Info` and `Warn` collector levels,
append-only capture of exact entry pointers without copying, and the collector
hook installed only on the private package logger. `DebugLogger` returns an
empty allocated field map and enables debug only on the private logger;
`Context` returns empty data; and `ExternalError` preserves representative
exact multiline and empty text without wrapping. JSON logging selects the
exact zero `logrus.JSONFormatter`; field-logger transitions are exact;
`StdOut` depends only on standard-logrus debug enablement and returns the exact
`os.Stdout` pointer or nil; and `LogEntries` exposes the collector's exact
slice backing and entry pointers without copying.

Every mutating subtest snapshots and restores the exact `fieldLogger`,
`collector`, private `log`, and touched private and standard-logrus level,
formatter, output, hooks, and caller-reporting state. Outer contracts verify
the restored identities and values. The tests use only in-memory logrus
objects and buffers, use no parallelism around globals, and leave no hook or
entry behind. They open no socket, launch no program, sleep no clock, and write
no fixture. Every iterated table and recorded-entry population is required to
be non-empty. Preferred logger ownership, isolation, copying, formatting, or
output design remains a later product decision requiring explicit authority.

## T15 Apparatus Checkpoint

The P3.54 T15 repair remains unchanged. Separate verified execution and
structured replicas preserve historical test effects while measuring a
pristine exact `5635d50` checkout plus the declared inventory overlay. No
output path is excluded or deleted, and every effect stays in the retained
manifest. The complete 15-control proof passes and reproduces old scorecard
`d420887d73aabf496ff276fcc55d13ad379ac49c9322fad58808b5e28fdba7df`,
new scorecard
`5fb3226009cfbf0d29f63fa03592157cce4efcec6e38583b64a86f6288e89490`,
normalized stored raw body
`cf23c9dca987acd3a966693f933c4c7eca9f7d11c51f0fbabfe1d697f3d7497f`,
Q3.9, and all 228 numeric debt leaves. P4.8 did not modify the apparatus,
parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `67344a8` reports:

- Measured commit tree: `e85dc852f5875c812509b8f7e1b08520d7f0251d`.
- Structured scorecard SHA-256:
  `38e938f060312840cb244a575536aa613266b172de98542cb2bb2e0998da576b`.
- Absolute L0: 8 of 8.
- 440 test functions, zero skipped; 26 of 27 packages have tests.
- `pkg/logger`: 100% statement coverage.
- Q0.6: 27 guarded safe-writer sites, 22 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 1 of 27 packages has no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 0 direct external sites outside five valid declared adapters of 27
  production effect sites.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 92 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 13 documented non-passing criteria, never 2.
- Comparable ratchets: six improved, two held, zero regressed, zero
  non-comparable.
- Clean identity: zero dirty paths.

The focused seven-criterion scorecard under
`/private/tmp/ply-p4-logger-focused-audit.8CF6Sv` has SHA-256
`5c640e0717d17f3200b1acfed83171f1cc824529ae2af18ab9b31a7b4fa77ee6`.
The full clean audit is under
`/private/tmp/ply-p4-logger-full-audit.bqTIxi`. Pinned tools remain under
`/private/tmp/ply-p4-tools` and `/private/tmp/ply-p358-tools`.

## Next Objective

Continue P4 with one test-only coverage move for `pkg/context`, the sole
remaining untested package, reducing Q1.1 from 1 of 27 to zero without changing
production code or observable semantics. Characterize the existing Context
state, `FindAndPopulateMavenProjects`, `OnEachMavenProject`, `OnRootProject`,
`LoadProfile`, `GetMavenRepository`, package logger initialization, and
`SetLogger` through their current file, config, Maven, project, and logger
interfaces.

Keep the move inside one new `pkg/context` test file. Use `t.TempDir()`, the
central guarded writer, in-memory logrus values, and isolated test-owned HOME,
config, Maven, and process state. Do not use parallel tests around globals or
process state. Snapshot and restore the private logger's exact identity plus
every environment, working-directory, or other process value touched. Open no
socket, access no network, launch no program, and perform no sleep or timed
wait. Require all iterated table, callback, project, job, and recorded-entry
populations to be non-empty. Leave production context, logger, config, Maven,
file code, callers, scanner, inventory, apparatus, mutation work, manual
evidence, and P5-P8 unchanged.

## Verification Notes

- Focused logger and relevant command and package caller tests pass; the logger
  package reports 100% statement coverage and its focused race, vet, shuffled
  count-10, and pinned-lint runs pass.
- Actual API/CLI, CLI surface, fresh executable subprocess, and compatibility
  meta-contracts pass with reports and caches outside the worktree.
- Complete build and uncached tests, race, vet, exact pinned lint,
  `make preflight`, `make test`, Make contracts, and all 62 launcher controls
  pass.
- The first focused coverage invocation used the sandbox-blocked default Go
  build cache; its unchanged isolated-cache rerun passed.
- Initial standalone launcher and `make test` invocations hit the documented
  partial-raw-log signal-fixture race at control 26 or during its nested
  control-50 run; unchanged complete reruns passed all 62 controls.
- All production-script meta-contracts, all four host acceptance flows, and the
  repaired 15-control audit meta-suite pass.
- Empty-HOME `go test ./... -count=2` passes with isolated writable state under
  `/private/tmp/ply-p4-logger-hermetic.rbdLgq` and the existing module cache.
- Full clean audit: expected exit 1, 13 non-passing criteria, six improved, two
  held, zero regressed, zero non-comparable, and zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the complete P4 entry and
checkpoint gate, both design documents, `.quality/inventory`, both complete
context production files, every context caller, and all reached config,
project, Maven, file, logger, profile, and guarded-writer code. Read
representative restoration, temporary-home, working-directory, exact-log,
partial-result, job-order, dry-run, callback, pointer-identity, and non-empty
contracts, the import-aware scanner, API/CLI contracts, and the T15 repair and
baseline reproduction README before editing. Confirm branch, HEAD, clean
status, reciprocal archive links, launcher `--check`, and exact product commit
`67344a8`.

Make one focused context test implementation commit, then the normal separate
continuity commit. Stop before production context/logger/config/Maven/file
changes, another package, manual evidence, P5-P8, publication, or distribution.
Do not push, merge, stash, revert, launch a successor, or remove the worktree.
