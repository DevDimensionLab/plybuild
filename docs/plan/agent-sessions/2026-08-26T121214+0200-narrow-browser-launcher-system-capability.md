# Agent Session: Narrow Browser Launcher System Capability

Status: ANSWERED - HISTORY
Session ID: `2026-08-26T121214+0200-narrow-browser-launcher-system-capability`
Created: `2026-08-26T12:12:14+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `798d847cc6c9bbd806660730cd7f4729b9145bf07ac8c843ffafd1896cfe875d`
Previous: [2026-08-26T114525+0200-narrow-plugin-diagrams-system-capabilities.md](2026-08-26T114525+0200-narrow-plugin-diagrams-system-capabilities.md)
Next: [2026-08-26T124026+0200-narrow-unzip-output-file-provenance.md](2026-08-26T124026+0200-narrow-unzip-output-file-provenance.md)
Outcome: Completed in product commit `0b96f10073c771d8dd86360b916dd296220cebbe`.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 provenance move: give the existing private browser-
launcher production selector a runner-only system process dependency so public
`OpenBrowser` no longer inherits an unused process-dependency standard-output
capability. Preserve exact platform, asynchronous start, URL, command, error,
and public behavior; P3.63; every completed move; and zero comparable ratchet
regressions.

# Authorized Roadmap

P3 remains active, P4-P8 remain queued in the machine-readable block in
`docs/plan/quality-upgrade.md`, and the launcher must remain NEXT until all
authorized checkpoints finish. P3.63 product commit `1bce06f` passed its
truthful complete checkpoint.

This mission authorizes only focused browser-launcher recording contracts,
mechanical preservation of the complete process dependency in the affected
double, and replacement of the one `process.System()` selection inside
`systemBrowserLauncherDependencies()` in `pkg/webservice/api.go` with the
existing exact `process.SystemRunner()`.

It does not authorize a browser-launcher execution or command-flow change; a
real browser launch; server, clock, web handler, port, or shutdown work;
plugin diagrams; profile editor; shell Run or Git; Maven; Unzip; another
process caller; another adapter; inventory; scanner; audit apparatus; mutation
harnesses; or P4-P8 implementation. Do not change `process.System()`,
`process.SystemRunner()`, the filesystem selector, Maven's completed standard-
output behavior, or the completed plugin-diagrams selectors.

# Measurements At Start

Clean product commit `1bce06f` has 388 tests across 19 of 25 packages. Q0.6 has
26 guarded safe-writer sites, 21 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 6/25, Q1.2 is zero, exact Q1.3 is 8/30 with
clock and server absent, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is zero
phrases across 81 Markdown files.

The authoritative clean full audit exits 1 for 15 documented findings, never
2, with L0 8/8, five improved, two held, zero regressed, one non-comparable
ratchet, and zero dirty paths. Its scorecard SHA-256 is
`5cf5c2c171749e42dee5dfea6f5986b09dce96d00510560f4e625c7a6e4e3643`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for this single browser-launcher system-capability move
and its recording contracts, then make the normal separate continuity-only
commit. Do not push, merge, publish, distribute, remove the worktree, stash
inherited changes, revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the existing complete
caller-owned browser-launcher dependency value. The production process
dependency must carry the exact existing system runner and nil `Process.Stdout`;
the separate exact `runtime.GOOS` selection must remain unchanged.

Preserve the exact asynchronous `xdg-open` Linux request, `rundll32` Windows
request with ordered `url.dll,FileProtocolHandler` and URL arguments, and
`open` macOS request. Preserve each arbitrary URL byte, empty directory, nil
stdin/stdout/stderr command streams, `Start: true`, direct process start error,
one exact request attempt, and absence of another request. Preserve the exact
unsupported-platform error and no process attempt, safe zero process
dependency, non-empty population, exported `OpenBrowser` signature and return,
runtime platform selection, caller placement, and every public behavior. Do
not add fallback, retry, wrapping, logging, synchronous execution, cleanup, or
global-state changes.

Focused tests must launch no external program, touch no network, and write no
repository fixture. Preserve complete caller-owned dependency fields in the
affected recording double. Leave all completed process, plugin-diagrams,
profile, shell, Maven, filesystem, server, browser-call-site, and Unzip behavior
unchanged.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact commit `1bce06f`. Read the rolling handover, this
archive, the P3 tail and checkpoint gate, both design documents,
`.quality/inventory`, complete process adapter code/tests and every complete
process double/caller, complete browser-launcher code and all its contracts and
callers, the P3.63 plugin-diagrams contracts, P3.62 profile-editor contracts,
P3.61 shell Git contracts, P3.60 shell Run contracts, P3.58 Maven stdout
contracts, the import-aware scanner, API/CLI contracts, and the T15 repair and
baseline reproduction README.

# Three Moves

1. Start red by strengthening the focused browser-launcher selector and
   recording contracts. Prove the private production selector carries the
   exact system runner and nil process-dependency stdout while retaining exact
   `runtime.GOOS`; prove every complete platform request, URL argument, stream,
   asynchronous attempt, direct error, unsupported-platform result, safe zero
   behavior, complete dependency preservation, non-empty population, and
   absence of another request.
2. Replace only the authorized `process.System()` selection with
   `process.SystemRunner()`. Preserve the complete selector, `openBrowser`,
   exported `OpenBrowser`, every command field, runtime selection, caller, and
   return path. Regenerate exact Q1.3 without changing the scanner or
   broadening the move to force a number.
3. Run focused browser-launcher/process and relevant caller tests; API/CLI and
   subprocess compatibility; launcher and Make contracts; complete tests,
   race, and vet; all four host acceptance flows; the 15-control audit meta-
   suite; focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurement; full clean
   audit; and empty-HOME count-2. Expect Q0.6 to hold at 26 guarded sites. The
   full audit may exit 1 for documented findings but never 2, and comparable
   ratchets must not regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent browser-launcher
system capability move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create exactly one
reciprocally linked NEXT archive for the next coherent authorized roadmap move,
replace only the launcher's mutable regions, run launcher and handoff contracts,
and make the separate `docs: prepare next agent session` commit. Do not launch a
real successor, push, merge, publish, distribute, stash, revert, or remove the
worktree. COMPLETE remains invalid while P3 or P4-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
