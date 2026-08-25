# Agent Session: Migrate Open Browser Process Start

Status: NEXT
Session ID: `2026-08-25T151736+0200-migrate-open-browser-process-start`
Created: `2026-08-25T15:17:36+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `3b58e8754104eb4030bccfa56b4390b2c1ca1aaea9e0e6fa103edd50584089b0`
Previous: [2026-08-25T143403+0200-migrate-profile-editor-process.md](2026-08-25T143403+0200-migrate-profile-editor-process.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the three direct
`exec.Command(...).Start()` paths in `pkg/webservice.OpenBrowser` through the
existing process adapter. Preserve exact runtime platform selection,
executable and argument values, asynchronous start semantics, empty working
directory, nil stdin/stdout/stderr, one start attempt, exact start errors, the
unsupported-platform error, every server operation and caller, every completed
process, profile, filesystem, HTTP, tips, config, Maven, structurizr,
Bitbucket, Wpost, local-config, and supervisor move, and every P2A contract
with zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the three
chained process constructions and `Start` calls in `OpenBrowser`; one boolean
start-mode addition to the existing internal `process.Command` request and its
exact direct system selection of `exec.Cmd.Start()` instead of the established
`exec.Cmd.Run()`; focused process-start contracts; and focused private browser-
launcher recording contracts. It does not authorize another process caller,
`StartWebServer`, `StopWebServer`, an OS/runtime adapter, a function-valued
dependency, public API, inventory, mutation harnesses, clock/server seams, or
later roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, complete `pkg/webservice/api.go` and its package,
every `OpenBrowser` caller and relevant build command/test, complete
`internal/adapter/process` source and tests, every `process.Command`,
`process.Dependencies`, `process.Runner`, and `process.Execute` caller and
recording double, the completed profile-editor stdin move, Maven and Git
process flows, relevant structurizr, shell, context, config, HTTP, tips,
filesystem, Bitbucket, Wpost, and local-config code/tests, `.quality/inventory`,
the continuity and quality-lift designs, and the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/
Q2.1/Q3.4 audit implementation. Regenerate ignored reports outside the
measured tree or remove them before a clean audit.

Implementation commit `0e10282` has 308 tests across 18 of 25 packages. Q0.6
has 24 guarded safe-writer sites, 19 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 7 of 25, Q1.2 is 0, Q1.3 is 19 violations of 38 production
effect sites with clock and server absent, Q1.4 is 7 of 8, and exact Q2.1 is 0
of 8 executable harnesses. The clean gate passed, the full audit exited 1 for
15 documented findings and never 2, and comparable ratchets were five
improved, two held, and zero regressed.

The launcher has 62 Bash 3.2 contracts and supervises fresh non-interactive
JSONL turns with external raw logs. It continues only after a successful
structured stream and valid clean committed handoff. Preserve its stable
skeleton and do not launch a real successor while developing, testing, or
finalizing this move.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for exact process start representation and the browser
launcher boundary and contracts, with measured planning notes. Then perform
the normal separate handoff-only commit. Do not push, merge, publish,
distribute, remove the worktree, stash inherited changes, revert user work, or
run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `OpenBrowser(string) error`, every caller,
all three supported platform branches and their order, `StartWebServer`,
`StopWebServer`, server globals and handlers, command registrations, public
API/CLI behavior, all four host acceptance flows, process stdin behavior,
profile behavior, Maven and Git process behavior, and every caller's
observable behavior.

Add only a boolean start-mode value to the existing complete
`process.Command` request. In the system runner, retain exact command
construction and field mapping, call `cmd.Start()` and return its exact error
only when that value is true, and otherwise retain the existing single
`cmd.Run()` attempt. Keep name, args, directory, stdin, stdout, stderr, safe
zero behavior, `Runner`, `Dependencies`, `Execute`, `System`, and every existing
caller unchanged. Do not add a new interface method, wait for or kill a started
process, retry, invoke a shell, split the process operation, or change error
selection.

Preserve exact browser selection. Linux must request `xdg-open` with the URL as
its only argument; Windows must request `rundll32` with exact first argument
`url.dll,FileProtocolHandler` and the URL second; Darwin must request `open`
with the URL as its only argument. Each supported request must use start mode,
empty `Dir`, and nil `Stdin`, `Stdout`, and `Stderr`, make one adapter attempt,
return the exact start error, and never wait or run synchronously. An unsupported
platform must return the same exact error text without a process attempt.

Do not clean, parse, validate, escape, or rewrite the URL; use a shell; add
stream forwarding; change supported-platform order or spelling; preflight an
executable; retry; wrap an error; broaden into HTTP serving, shutdown context,
runtime/clock/server adapters, shell.Run, structurizr processes, Maven, Git,
profile, filesystem, P4, P5, Docker, distribution, or publication work.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete
`pkg/webservice/api.go` and related source/tests/callers, relevant build command
registration and acceptance contracts, complete `internal/adapter/process`
source/tests and every complete process caller/double, the completed profile,
Maven, and Git process contracts, relevant shell/structurizr/HTTP/context/config
code and tests, and the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation
before editing. Before the full gate, read the complete launcher contract,
relevant Make meta-tests, P2A API/CLI/subprocess contracts, and all four host
acceptance flows.

# Three Moves

1. Start red by extending focused process adapter contracts to prove exact
   start-mode identity reaches a recording runner; false still performs one
   direct synchronous run; true performs one direct asynchronous start with
   exact name, args, directory, streams, and errors; and zero dependencies stay
   safe. Add focused browser-launcher contracts for complete production
   selection, all three exact platform requests, arbitrary URL bytes, one
   attempt, exact errors, unsupported-platform no-attempt behavior, non-empty
   recorded populations, safe zero behavior, and no real browser or server.

2. Add only the boolean start mode and its direct system branch. Keep
   `OpenBrowser` as the production entry, select `process.System()` only for a
   private complete browser composition, preserve runtime selection, and
   replace only the three direct construction/start chains. Do not change
   another webservice function, process call, adapter operation, caller,
   inventory, or public API.

3. Run focused webservice/process and relevant command, context, config, HTTP,
   Maven, shell, structurizr, profile, tips, filesystem, Bitbucket, Wpost,
   local-config, and caller package tests; the launcher contract from
   `/bin/bash`; Make preflight meta-contracts; API/CLI and subprocess
   compatibility; full Go tests and race/vet; all four host acceptance flows;
   the audit meta-suite; focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4
   measurements; full clean checkpoint audit; and empty-HOME count-2. Expect
   Q1.1 to improve from 7 to 6 of 25 when `pkg/webservice` gains direct tests
   and nominal Q1.3 to improve from 19 of 38 to 16 of 36 when three caller
   sites leave and one exact adapter start site enters. Q0.6, Q1.2, Q1.4, and
   exact Q2.1 must hold. Regenerate exact values; the full audit may exit 1 for
   documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent process-start
and open-browser move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create one linked
NEXT archive for the next coherent P3 effect move, replace only the launcher's
mutable regions, run the launcher contract, and make the separate handoff-only
commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
