# Agent Session: Migrate Shell Unzip Archive Close

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T204152+0200-migrate-shell-unzip-archive-close`
Created: `2026-08-25T20:41:52+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `7865b0caac88090284817a37cc96e9ef86c4f6b5a98fbc11ec0fc4bf68bba981`
Previous: [2026-08-25T200905+0200-migrate-shell-unzip-entry-close.md](2026-08-25T200905+0200-migrate-shell-unzip-entry-close.md)
Next: [2026-08-25T220107+0200-resume-shell-unzip-archive-close.md](2026-08-25T220107+0200-resume-shell-unzip-archive-close.md)
Outcome: Blocked because `*zip.ReadCloser` implements `io.Closer`, not `io.ReadCloser`; the user authorized the successor session to widen only the existing CloseReader parameter to `io.Closer`.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the deferred direct
`_ = r.Close()` request in private `pkg/shell.Unzip` through the existing narrow
filesystem `CloseReader` operation. Preserve the public signature, exact opened
archive-reader identity, deferred placement, one close attempt after every
successful archive open, ignored archive-close result, archive and entry
traversal, entry order and bytes, filenames and partial results, zip-slip
behavior, every completed Unzip operation and error precedence, and every
earlier shell, filesystem, Spring, process, browser, profile, HTTP, tips,
config, Maven, structurizr, Bitbucket, Wpost, local-config, Kibana, and
supervisor move, with zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only focused
private Unzip archive-close recording contracts and replacement of the one
deferred direct archive close with
`defer func() { _ = filesystem.CloseReader(dependencies.Files, r) }()`.

Reuse the existing `CloseReader(io.ReadCloser) error` interface operation,
zero-safe helper, exact system implementation, and complete doubles unchanged.
It does not authorize another filesystem method, adapter reshaping, archive
open, entry-reader Close, output-file Close, Copy, OpenFile, MkdirAll,
entry-open injection, another unzip branch, shell Run or Git changes, process
or HTTP work, Spring, Maven, structurizr, plugin diagrams, clock/server
adapters, public API, inventory, mutation harnesses, audit reshaping, or
later-roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, complete `pkg/shell` source/tests/callers, complete
`internal/adapter/filesystem` source/tests and every complete filesystem
double, all Unzip callers and completed unzip contracts, relevant Spring
download, command, context, config, HTTP, process, Maven, structurizr,
Bitbucket, Wpost, local-config, browser, profile, tips, file, template, and
Kibana code/tests, `.quality/inventory`, both design documents, and the
import-aware effect scanner. Regenerate ignored reports outside the measured
tree or remove them before a clean audit.

Implementation commit `70bee0e` has 358 tests across 19 of 25 packages. Q0.6
has 26 guarded safe-writer sites, 21 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 6 of 25, Q1.2 is 0, Q1.3 is 10 violations of 32
production effect sites with clock and server absent, Q1.4 is 7 of 8, and
exact Q2.1 is 0 of 8 executable harnesses. The clean full audit exited 1 for
15 documented findings and never 2; comparable ratchets were five improved,
two held, and zero regressed.

The direct interface-typed entry-reader Close is gone. Exact Q1.3 holds at
10/32 because its adapter request is also absent from the reported violation
set. The deferred concrete archive-reader Close is presently absent from that
set as well. Replacing it with the existing helper may change or hold the
measured population; regenerate the exact value without changing the scanner
or broadening this move to force an expected number.

The launcher has 62 Bash 3.2 contracts and supervises fresh non-interactive
JSONL turns with external raw logs. It continues only after a successful
structured stream and valid clean committed handoff. Preserve its stable
skeleton and do not launch a real successor while developing or finalizing.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the deferred Unzip archive-close move and focused
private recording contracts, with measured planning notes. Then make the
normal separate handoff-only commit. Do not push, merge, publish, distribute,
remove the worktree, stash inherited changes, revert user work, or run
destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve public
`Unzip(string, string) ([]string, error)`, the private complete
`unzipDependencies`, production `filesystem.System()` selection, the exact
`*zip.ReadCloser` returned by `OpenZipReader`, and the current defer immediately
after a successful archive-open request. Make exactly one filesystem
CloseReader request for that archive on every later return, after all reached
entry operations and entry-reader closes. Continue to ignore its exact result
without changing the named return values or any earlier error.

Pass the complete Unzip dependency composition whole. Do not add fallback,
retry, logging, cleanup, another close attempt, a new defer, environment or
working-directory behavior, or error normalization. Do not close an archive
when archive open fails. Any fixture write must stay guarded below
`t.TempDir()`; launch no process, touch no public network, and change no working
directory.

Do not change `OpenZipReader`, `filesystem.CloseReader` or its system
implementation, entry-reader Close, `filesystem.Close` or output-file Close,
Copy, OpenFile, MkdirAll, entry open, shell Run, Git composition, structurizr
command construction, Maven process composition, plugin diagrams, another
adapter operation or caller, public API, or `.quality/inventory`.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete `pkg/shell` source
and tests, every Unzip caller, complete `internal/adapter/filesystem` source
and tests and every complete double, completed Unzip/Spring/filesystem/process
contracts, the audit's import-aware effect scanner, and relevant caller
packages before editing. Before the full gate, read the complete launcher
contract, Make meta-tests, P2A API/CLI and subprocess contracts, and all four
host acceptance flows.

# Three Moves

1. Start red with only focused private Unzip archive-close recording
   contracts. Prove exact archive identity, one deferred request after every
   successful archive open, request placement after successful traversal and
   entry-reader closes, execution on zip-slip and every established later
   error return, no request after archive-open failure, ignored injected
   archive-close error without result or error-precedence change, complete
   dependency selection, rejection of an empty archive-close recording
   population, distinction from entry-reader-close requests, and absence of
   unrelated filesystem requests. Invoke no real process or public network
   request.

2. Replace only `defer func() { _ = r.Close() }()` with
   `defer func() { _ = filesystem.CloseReader(dependencies.Files, r) }()`.
   Keep the defer placement, anonymous function, ignored assignment, archive
   variable, surrounding open branch, traversal, and every earlier and later
   operation unchanged. Do not change the filesystem adapter or doubles except
   the focused private recording behavior needed by the red contracts.

3. Run focused shell/filesystem/Spring and relevant process, Maven, command,
   context, config, HTTP, structurizr, profile, browser, tips, file, template,
   Bitbucket, Wpost, local-config, Kibana, and caller package tests; the
   launcher contract from `/bin/bash`; Make preflight meta-contracts; API/CLI
   and subprocess compatibility; full Go tests and race/vet; all four host
   acceptance flows; the audit meta-suite; focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/
   Q2.1/Q3.4 measurements; full clean checkpoint audit; and empty-HOME count-2.
   Expect Q0.6 to hold at 26 guarded sites and regenerate exact Q1.3 without
   changing the scanner; the full audit may exit 1 for documented findings but
   never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent deferred Unzip
archive-close move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create one linked
NEXT archive for the next coherent P3 effect move, replace only the launcher's
mutable regions, run the launcher contract, and make the separate handoff-only
commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
