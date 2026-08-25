# Quality Upgrade Handover

Generated: 2026-08-25T09:35:10+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 775457562cf1.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d,
d6cb593, 82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, and 7754575 in roadmap
order. The separate operational continuity implementation is 1b85711 and
changes no Go quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T093228+0200-migrate-tips-show-read.md. The
project-config-write predecessor is answered history, and the reciprocal
archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 33 Preserved

Exported (*ProjectConfiguration).WriteTo(targetFile string) error retains its
signature and remains the production entry satisfying ProjectConfig. Only its
private complete project-config-write composition selects filesystem.System().
The exact caller target, json.MarshalIndent result, and mode 0644 pass through
the existing filesystem.WriteFile operation. No adapter operation, established
complete recording double, caller, or public contract changed.

The private composition retains the exact first
log.Infof("writes project config file to %s", targetFile) call, then the exact
json.MarshalIndent(config, "", "    ") evaluation and early error return,
then one write attempt. It does not clean, join, normalize, preflight, create
directories, retry, inspect permissions, mutate the receiver, wrap the write
error, or add logging. Nil receivers serialize to exact bytes `null` and
non-nil receivers retain the established indented JSON representation.

Four new top-level config contracts bring the suite to 258 tests. They prove
complete production selection and delivery, exact arbitrary targets, non-empty
recorded write populations, exact nil, zero-value, empty, and non-ASCII JSON
bytes including map-key ordering, mode 0644, one write, exact error identity,
log-before-write sequencing, unchanged receiver state, and safe defaults
without developer-path access. The production composition tests do not mutate
the real filesystem.

## Measured Quality State

The clean full audit from implementation commit 7754575 reports:

- Absolute L0: 8 of 8.
- 258 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 33 direct external sites outside five declared adapters of 50
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 temporary-state claims across 47 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The clean report is
/private/tmp/ply-project-config-clean-audit.nVbvrs/scorecard.json and the
focused report is
/private/tmp/ply-project-config-focused.agRgF7/scorecard.json. All
authoritative audit output and tool caches stayed outside the measured tree.

## Decisions And Learned Facts

1. ProjectConfiguration.WriteTo reuses the adapter's established WriteFile
   operation; this move does not extend the adapter or any established complete
   filesystem recording double.
2. Q1.3 moved from 34/51 to 33/50 because the direct project-config write
   disappeared while the already-counted adapter operation was reused.
3. The exported production entry alone selects filesystem.System() for its
   private complete composition; ProjectConfig and every caller remain
   unchanged.
4. The log, marshal, early return, target, serialized bytes, mode, single
   attempt, exact error, and receiver state are preserved without additional
   effects.
5. The focused recording double rejects every unrelated adapter operation and
   rejects an empty write population. Its safe zero case returns the exact
   filesystem.ErrNoFilesystem without touching a developer path.
6. The next isolated flow is cmd/tips.go's os.ReadFile(tipsPath). It can reuse
   the established ReadFile operation without extending the adapter or an
   established complete double, nominally moving Q1.3 to 32/49.
7. The next move must preserve the missing-argument warning/delegation, exact
   name and file.Path evaluation, read-error wrapping with %w, returned bytes,
   terminal-config lookup, markdown rendering, stdout, local-source log,
   global-cloud-config lookup, and cloud-source log sequencing.
8. The next production composition should isolate only the read dependency;
   focused contracts must not run the command or read a real filesystem path.
9. .quality/inventory is baseline-checksum-bound. Do not relabel seams or
   claim P5 mutation coverage.
10. Supply APIDIFF and GOLANGCI_LINT as environment variables for preflight,
    and keep GOLANGCI_LINT_CACHE and GOCACHE in writable external directories.
    The pinned tools are under /private/tmp/ply-quality-tools while that
    temporary directory remains available.
11. The complete preflight, make-test target, and independent launcher runs
    passed all 62 launcher controls without a retry.
12. A discarded clean-audit wrapper invocation used the nonexistent path
    .quality/baseline.json and failed closed with exit 2 after the upstream
    scorecard. The authoritative rerun used .quality/baseline/scorecard.json
    plus .quality/baseline/manual-evidence.json and exited 1 with valid
    structured output and zero regressions.

## Next Objective

Move only the direct os.ReadFile(tipsPath) in tipsShowCmd.RunE through the
existing filesystem ReadFile operation.

Start red with focused command recording contracts for the complete dependency,
production system selection, exact delivered path, non-empty read populations,
empty and non-ASCII bytes, one attempt, exact private-boundary error identity,
and safe zero behavior without developer-path access. Production-composition
tests must not read the real filesystem or run the command.

Preserve all command objects, registration, the zero-argument branch, exact
name and path composition, the existing fmt.Errorf text and %w identity,
returned source bytes, later terminal/render/stdout/log/cloud sequencing, and
every caller. Production must select filesystem.System() only for a private
complete tips-show-read composition. Reuse adapter ReadFile without extending
the adapter or established complete doubles. Do not change list, sync, profile,
config, template/file/Maven behavior, inventory, public API, another production
effect, or enter clock/server, P4, P5, or later roadmap work.

Expected direction is nominally Q1.3 32/49 with Q0.6, Q1.2, Q1.4, and exact
Q2.1 held. Regenerate the structured result and require zero comparable
ratchet regressions.

## Verification Notes

Completed from implementation commit 775457562cf1:

- Red evidence: focused compilation failed only on the absent private
  project-config dependency, system composition, and writeTo composition.
- Focused filesystem, config, file, template, Maven, context, and command
  caller package tests: PASS.
- Standalone /bin/bash launcher contract: PASS all 62 controls.
- Make preflight meta-contracts and full preflight: PASS, including all 15
  audit meta-controls and all 62 launcher controls.
- API, CLI, and subprocess compatibility: PASS.
- make test and make test-install: PASS without a retry.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME go test ./... -count=2: PASS.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable; Q1.3 is exactly 33/50.
- Clean full audit: exit 1, 15 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Implementation commit: 775457562cf184ff6c8c049581dadcd97b38ceb0
  (quality: route project config write through filesystem adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete tips command and package code/tests, relevant
context/config/file/template/Maven callers, the filesystem adapter and doubles,
and the named audits before editing. Confirm branch, HEAD, status, reciprocal
links, and ./codex-dev-start.sh --check. Begin red for only the tips-show read
boundary, then finish with one implementation commit and one separate
handoff-only commit.

Stop before another tips branch or command, argument/registration/path/error/
render/output/log/cloud changes, another config or filesystem operation,
another adapter operation or family, Q1.4 expansion, P4, mutation, Docker,
cloud distribution, or publication. Stop on API/CLI change, comparable ratchet
regression, authoritative full-audit exit 2, or failure to isolate
implementation from generated and handoff-only state.
