# Quality Upgrade Handover

Generated: 2026-08-25T15:17:36+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 0e10282eeb3e.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d,
d6cb593, 82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438,
527a8b9, 886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc,
85f4c2b, and 0e10282 in roadmap order. The separate operational continuity
implementation is 1b85711 and changes no Go quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T151736+0200-migrate-open-browser-process-start.md.
The profile-editor process predecessor is answered history, and the reciprocal
archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 44 Preserved

`process.Command` now carries `Stdin io.Reader` in addition to its established
name, arguments, directory, stdout, and stderr. The system runner maps it
straight to `exec.Cmd.Stdin` before the unchanged stdout/stderr mappings and
single `cmd.Run`. `Runner`, `Dependencies`, `Execute`, `System`, zero-value
behavior, errors, and all earlier callers remain unchanged.

The profile edit branch still reads `EDITOR` once, replaces only empty with
exact `vim`, and evaluates `ctx.LocalConfig.FilePath()` once before selecting
its private production composition. That composition contains exact
`process.System()`, `os.Stdin`, and `os.Stdout`; the request carries the exact
editor as `Name`, the exact arbitrary config path as its only argument, empty
`Dir`, exact stdin/stdout identities, and nil stderr. One adapter attempt is
made. Its exact error returns before sync, reset, or print; success preserves
all later branches, conditions, ordering, and results.

Four focused profile-editor contracts bring the suite to 308 tests. Together
with the extended process contracts they prove complete production selection,
exact editor and fallback behavior, arbitrary path and stream identity, one
attempt, exact errors, direct system stdin delivery, non-empty populations,
safe zero behavior, and no unrelated or real process. No public API, command
object or registration, caller, inventory, seam, mutation harness, completed
adapter operation, or completed effect changed.

## Measured Quality State

The clean full audit from implementation commit 0e10282 reports:

- Absolute L0: 8 of 8.
- 308 test functions, zero skipped; 18 of 25 packages have tests.
- Q0.6: 24 guarded safe-writer sites, 19 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 7 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 19 direct external sites outside five declared adapters of 38
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 temporary-state claims across 58 Markdown files at the clean
  implementation checkpoint.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative clean report is
/private/tmp/ply-profile-editor-tools/audit-full-clean/scorecard.json and the
focused implementation report is
/private/tmp/ply-profile-editor-tools/audit-focused-final-dirty/scorecard.json.
Compatibility reports, tool binaries, Go and linter caches, audit output, and
empty-HOME state stayed outside the measured tree.

## Decisions And Learned Facts

1. Exact stdin is request data, not a new process operation. Adding one reader
   field and one direct system assignment preserves every established caller.
2. The private profile composition must carry process, stdin, and stdout as one
   complete injected value. Selecting the two OS streams there avoids creating
   a new direct stream-capability site at the adapter call.
3. Both former profile sites leave the exact scanner population, producing the
   predicted Q1.3 result of 19/38.
4. The process recording contract checks reader identity without reading it;
   the system contract proves the delivered input reaches the child unchanged.
5. The next coherent synchronous boundary is not chosen because `shell.Run`
   exposes `exec.Cmd.String()` logging semantics. The narrower next process
   flow is the three asynchronous browser starts in `pkg/webservice/api.go`.
6. `OpenBrowser` has one exact command per supported runtime: Linux `xdg-open`
   plus URL, Windows `rundll32` plus `url.dll,FileProtocolHandler` and URL, and
   Darwin `open` plus URL. Each calls `Start`, never `Run` or `Wait`.
7. A boolean start mode on `process.Command` can preserve every zero-valued
   synchronous caller and the existing runner interface while selecting exact
   `exec.Cmd.Start` only for browser requests.
8. Unsupported platforms return exact `unsupported platform` without a process
   attempt; server startup, shutdown, handlers, runtime selection, and callers
   remain outside the next scope.
9. Expected next direction is Q1.1 from 7/25 to 6/25 when `pkg/webservice`
   gains direct tests and nominal Q1.3 from 19/38 to 16/36 when three caller
   sites leave and one exact adapter start site enters. Regenerate both.
10. `.quality/inventory` remains baseline-checksum-bound. Do not relabel seams,
    make mutation harnesses executable, or claim P5 coverage.
11. Full preflight needs external APIDIFF, GOCACHE, GOTMPDIR,
    GOLANGCI_LINT_CACHE, and golangci-lint v2.12.2. The final run passed every
    repository control. One additional launcher-only run hit the signal probe's
    intermittent partial-log timing failure; its isolated retry passed all 62.

## Next Objective

Move only the three supported-platform `exec.Command(...).Start()` chains in
`pkg/webservice.OpenBrowser` through the existing process adapter.

Start red at the adapter boundary for exact start-mode identity, one true
asynchronous start, unchanged false synchronous run, exact request fields and
errors, and safe zero behavior. Add direct webservice recording contracts for
complete production selection, all three exact executable/argument mappings,
arbitrary URL bytes, empty directory and nil streams, one attempt, exact error,
unsupported-platform no-attempt behavior, non-empty populations, and no real
browser or server.

Add only one boolean start-mode request field and the direct system choice
between `cmd.Start()` and the established `cmd.Run()`. Preserve
`OpenBrowser(string) error`, runtime branch order, every caller, unsupported
error, server operations, prior process requests, profile stdin, Maven and Git
process behavior, public API/CLI, inventory, and every completed effect. Do not
migrate shell, structurizr, server, clock, filesystem, P4, or P5 work.

Expected direction is Q1.1 6/25 and nominal Q1.3 16/36 with Q0.6, Q1.2,
Q1.4, and exact Q2.1 held. Require zero comparable ratchet regressions and
regenerate exact values.

## Verification Notes

Completed from implementation commit 0e10282eeb3edf7db372fedd466e17ae030ca87b:

- Valid red evidence: isolated focused compilation failed only on absent
  `process.Command.Stdin` and the authorized private profile-editor symbols.
  An earlier host-cache setup failure was rejected as red evidence.
- Focused process/command and relevant config, context, Maven, tips, filesystem,
  file, template, structurizr, Bitbucket, HTTP, Kibana, Spring, shell,
  local-config, and caller package tests: PASS.
- Standalone `/bin/bash` launcher contract: PASS all 62 controls in both final
  preflight runs and the isolated final retry. One extra combined-target run
  hit only the signal probe's intermittent partial-log timing failure.
- Make preflight meta-contract, two complete preflights, and all 15 audit
  meta-controls: PASS with external tools, reports, and caches.
- API, CLI, and subprocess compatibility: PASS.
- Full make test, install, and agent-start targets: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance plus their meta-contracts:
  PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS.
- Focused seven-criterion audit: exit 1, four improved, two held, zero
  regressed, one non-comparable; Q1.3 is exactly 19/38.
- Clean full audit: exit 1, 15 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Handoff-state Q3.4 audit: exit 0 with zero state-claim phrases across 59
  Markdown files and zero ratchet regressions.
- Implementation commit: 0e10282eeb3edf7db372fedd466e17ae030ca87b
  (refactor: route profile editor through process adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete webservice implementation/tests and callers,
relevant build command and acceptance contracts, complete process adapter and
tests, every complete process caller/double, completed profile/Maven/Git
process contracts, and named audits before editing. Confirm branch, HEAD,
status, reciprocal links, and `./codex-dev-start.sh --check`. Begin red only for
process start mode and the browser-launch boundary, then finish with one
implementation commit and one separate handoff-only commit.

Stop before changing runtime selection, URL bytes, executable/argument values,
async start semantics, unsupported error, another webservice function or
process caller, another adapter operation or family, inventory, Q1.4, P4,
mutation, Docker, distribution, or publication. Stop on API/CLI change,
comparable ratchet regression, authoritative full-audit exit 2, or failure to
isolate implementation from generated and handoff-only state.
