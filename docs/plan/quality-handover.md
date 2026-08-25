# Quality Upgrade Handover

Generated: 2026-08-25T15:54:09+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 9fcdfb54a084.
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
85f4c2b, 0e10282, and 9fcdfb5 in roadmap order. The separate operational
continuity implementation is 1b85711 and changes no Go quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T155409+0200-migrate-spring-archive-working-directory.md.
The browser process-start predecessor is answered history, and the reciprocal
archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 45 Preserved

`process.Command` now carries boolean `Start` in addition to its established
name, arguments, directory, stdin, stdout, and stderr. The system runner maps
every established field exactly as before, returns direct `exec.Cmd.Start()`
only for true, and keeps the single direct `exec.Cmd.Run()` attempt for false.
`Runner`, `Dependencies`, `Execute`, `System`, exact errors, safe-zero behavior,
and all earlier process callers remain unchanged.

`OpenBrowser(string) error` still selects exact `runtime.GOOS` through a private
complete composition. Linux requests `xdg-open` plus the URL, Windows requests
`rundll32` plus exact `url.dll,FileProtocolHandler` and the URL, and Darwin
requests `open` plus the URL. Each supported request has start mode true, empty
directory, nil streams, one adapter attempt, exact URL bytes, and exact error
identity. Unsupported platforms still return exact `unsupported platform`
without a process attempt. Server globals, handlers, startup, shutdown, every
caller, public API, CLI, and every completed effect remain unchanged.

Eight focused contracts bring the suite to 316 tests and add direct tests to
`pkg/webservice`. They prove start-mode identity, false synchronous behavior,
true asynchronous behavior, direct start errors, complete production process
and platform selection, all three command mappings, arbitrary URL bytes, one
attempt, exact errors, nil streams, unsupported no-attempt behavior, non-empty
populations, and safe zero values without a real browser or server.

## Measured Quality State

The clean full audit from implementation commit 9fcdfb5 reports:

- Absolute L0: 8 of 8.
- 316 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 24 guarded safe-writer sites, 19 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 16 direct external sites outside five declared adapters of 36
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 temporary-state claims across 59 Markdown files at the clean
  implementation checkpoint.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative clean report is
/private/tmp/ply-open-browser-tools/reports/full-clean-9fcdfb5/scorecard.json
and the focused implementation report is
/private/tmp/ply-open-browser-tools/reports/focused-dirty/scorecard.json.
Compatibility reports, tool binaries, Go and linter caches, audit output, and
empty-HOME state stayed outside the measured tree.

## Decisions And Learned Facts

1. Start-versus-run is request data, not a new process operation. One boolean
   preserves all zero-valued synchronous callers and the established runner
   interface.
2. A private browser composition can carry exact runtime selection and the
   complete process dependency without introducing an OS adapter or
   function-valued dependency.
3. The three former browser construction/start sites leave and one direct
   system-adapter start site enters, producing the predicted Q1.3 result of
   16/36. Direct webservice tests produce the predicted Q1.1 result of 6/25.
4. The remaining process violations are three command constructions in
   `cmd/plugin_diagrams.go` and the generic construction in
   `pkg/shell/command.go`. The exported structurizr command signatures and
   `shell.Run`'s exact `exec.Cmd.String()` logging make those broader than the
   next isolated move.
5. Clock and server adapter creation belongs to queued P4. Do not move either
   `time.Now` call, Kibana sleep behavior, or webservice server operation during
   P3.
6. The narrow next existing-adapter effect is `archivePath`'s single direct
   `os.Getwd()` in `pkg/spring/io.go`. It occurs before the preserved direct
   `time.Now().Unix()` selection and exact `file.Path` formatting.
7. One distinct working-directory read on the complete filesystem adapter can
   preserve the private `archivePath() (string, error)` signature. Its wrapper
   must return `ErrNoFilesystem` safely for zero dependencies, and its system
   implementation must return direct `os.Getwd()` values and errors.
8. Every complete filesystem double must implement the added method only to
   remain complete. Existing focused doubles should reject it as unrelated;
   dedicated adapter and Spring archive-path doubles should record it.
9. Expected next direction is nominal Q1.3 from 16/36 to 15/36 when the caller
   site leaves and one exact adapter site enters. Q0.6, Q1.1, Q1.2, Q1.4, and
   exact Q2.1 should hold; regenerate every value.
10. `.quality/inventory` remains baseline-checksum-bound. Do not relabel seams,
    make mutation harnesses executable, or claim P4/P5 coverage.
11. Full preflight uses APIDIFF and golangci-lint v2.12.2 from
    `/private/tmp/ply-open-browser-tools/bin`, with reports and caches under the
    same external root. Every repository control passed.

## Next Objective

Move only `pkg/spring.archivePath`'s direct working-directory read through one
new distinct operation on the existing filesystem adapter.

Start red at the adapter boundary for exact path/error delivery, one attempt,
direct system `os.Getwd()` behavior, safe zero dependencies, and non-empty
recorded populations. Add private Spring archive-path contracts for complete
system dependency selection, an arbitrary delivered working directory, exact
error identity and empty path on error, one attempt, preserved
`spring-<Unix>.zip` construction under the delivered directory, safe zero
behavior, and no real file creation, download, unzip, delete, process, HTTP, or
clock seam.

Add only the working-directory read to the complete filesystem interface,
forwarder, safe-zero selection, and system implementation. Keep
`archivePath() (string, error)` as the production entry, select
`filesystem.System()` only in its private complete composition, and replace
only `os.Getwd()`. Preserve exact Getwd-before-time ordering, error short
circuit, direct `time.Now().Unix()`, `file.Path("%s/spring-%d.zip", ...)`, every
caller, Spring download/unzip/delete behavior, public API/CLI, inventory, and
every completed effect. Do not migrate a process, another filesystem caller,
clock, server, P4, or P5 work.

Expected direction is nominal Q1.3 15/36 with Q0.6 at 24 guarded sites,
Q1.1 6/25, Q1.2 zero, Q1.4 7/8, and exact Q2.1 0/8. Require zero comparable
ratchet regressions and regenerate exact values.

## Verification Notes

Completed from implementation commit 9fcdfb54a08479656d99a185750c75afb6919c6b:

- Valid red evidence: isolated focused compilation failed only on absent
  `process.Command.Start` and the authorized private browser-launcher symbols.
  An earlier host-cache permission failure was rejected as red evidence.
- Focused process/webservice and relevant command, context, config, HTTP,
  Maven, shell, structurizr, profile, tips, filesystem, file, template,
  Bitbucket, Kibana, Spring, local-config, and caller package tests: PASS.
- Standalone `/bin/bash` launcher contract: PASS all 62 controls.
- Make preflight meta-contracts, pinned linter with zero issues, one complete
  preflight, and all 15 audit meta-controls: PASS with external tools, reports,
  and caches.
- API, CLI, and subprocess compatibility: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance plus their meta-contracts:
  PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS.
- Focused seven-criterion audit: exit 1, four improved, two held, zero
  regressed, one non-comparable; Q1.3 is exactly 16/36.
- Clean full audit: exit 1, 15 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Handoff-state Q3.4 audit: exit 0 with zero state-claim phrases across 60
  Markdown files and zero ratchet regressions.
- Implementation commit: 9fcdfb54a08479656d99a185750c75afb6919c6b
  (refactor: route browser launch through process adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete Spring implementation/tests/callers, complete
filesystem adapter/tests and every complete double, relevant file, HTTP,
shell, process, config, Maven, structurizr, context, Bitbucket, Wpost,
local-config, browser, and audit contracts before editing. Confirm branch,
HEAD, status, reciprocal links, and `./codex-dev-start.sh --check`. Begin red
only for the working-directory adapter operation and private Spring archive
path, then finish with one implementation commit and one separate handoff-only
commit.

Stop before changing time selection, path formatting, another Spring function,
another filesystem caller or operation, a process/clock/server adapter, public
API, inventory, Q1.4, P4, mutation, Docker, distribution, or publication. Stop
on API/CLI change, comparable ratchet regression, authoritative full-audit exit
2, or failure to isolate implementation from generated and handoff-only state.
