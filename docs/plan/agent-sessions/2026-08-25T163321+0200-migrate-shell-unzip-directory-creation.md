# Agent Session: Migrate Shell Unzip Directory Creation

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T163321+0200-migrate-shell-unzip-directory-creation`
Created: `2026-08-25T16:33:21+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `3663db78d196aeca13987afb38f908915cf1b083e221f7d8d3392f426a6177de`
Previous: [2026-08-25T155409+0200-migrate-spring-archive-working-directory.md](2026-08-25T155409+0200-migrate-spring-archive-working-directory.md)
Next: [2026-08-25T171500+0200-migrate-shell-unzip-file-open.md](2026-08-25T171500+0200-migrate-shell-unzip-file-open.md)
Outcome: completed in `b9fe209`; both direct unzip directory creations now use the existing filesystem adapter with exact legacy order, arguments, errors, and results, and the clean gate has zero ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the two direct
`os.MkdirAll` calls in `pkg/shell.Unzip` through the established
`filesystem.MkdirAll` adapter operation. Preserve the public signature, exact
archive-open and traversal order, filename append order, directory
short-circuit, parent-directory selection, exact `os.ModePerm`, partial
results, error identity and precedence, every remaining unzip operation, every
completed filesystem, Spring, process, browser, profile, HTTP, tips, config,
Maven, structurizr, Bitbucket, Wpost, local-config, and supervisor move, and
every P2A contract with zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only one private
complete shell unzip composition containing `filesystem.Dependencies`,
production selection of `filesystem.System()`, focused recording contracts,
and replacement of the two direct `os.MkdirAll` calls with the existing
forwarding helper. It does not authorize a filesystem interface change,
another filesystem operation or caller, direct `zip.OpenReader`, `os.OpenFile`,
`io.Copy`, or close migration, process work, another shell function, Spring,
clock/server adapters, public API, inventory, mutation harnesses, or later
roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, complete `pkg/shell/command.go` and its package
tests and callers, complete `internal/adapter/filesystem` source/tests and
every complete double, completed unzip/Spring download contracts, relevant
file, HTTP, process, command, config, Maven, structurizr, context, Bitbucket,
Wpost, local-config, browser, profile, tips, and Spring code/tests,
`.quality/inventory`, the continuity and quality-lift designs, and the
Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation. Regenerate ignored
reports outside the measured tree or remove them before a clean audit.

Implementation commit `d982f63` has 324 tests across 19 of 25 packages. Q0.6
has 24 guarded safe-writer sites, 19 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 6 of 25, Q1.2 is 0, Q1.3 is 15 violations of 36 production
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
implementation commit for the two unzip directory-creation effects and private
recording contracts, with measured planning notes. Then perform the normal
separate handoff-only commit. Do not push, merge, publish, distribute, remove
the worktree, stash inherited changes, revert user work, or run destructive Git
commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve public
`Unzip(string, string) ([]string, error)`, every caller, direct
`zip.OpenReader(src)` selection and error, deferred archive close, entry order,
`filepath.Join(dest, f.Name)`, debug logging, zip-slip validation and exact
error text, filename append before mutation, directory `continue`,
`filepath.Dir(fpath)`, direct file open flags
`os.O_WRONLY|os.O_CREATE|os.O_TRUNC`, entry mode, direct entry open, copy, file
close, entry close, final results, and every error precedence.

Add no filesystem interface operation and change no adapter implementation.
Use the established `filesystem.MkdirAll` helper with exact `os.ModePerm`.
Zero dependencies at a reached directory-creation boundary must return exact
`filesystem.ErrNoFilesystem` without attempting another adapter operation.
Otherwise make exactly one MkdirAll attempt for each reached legacy site and
return its exact error. Select `filesystem.System()` only once through a
private complete production composition. Only a focused shell double may
record this existing operation; every pre-existing complete double and focused
behavior stays unchanged.

For a directory entry, preserve the exact joined path append before one
MkdirAll attempt and return the current partial filenames plus the exact error
on failure; success must continue directly to the next archive entry. For a
file entry, preserve the exact joined path append before one MkdirAll attempt
on `filepath.Dir(fpath)` and return the current partial filenames plus the
exact error before file open on failure. Do not clean or rewrite paths beyond
the existing join, clean, prefix, and parent operations. Do not add fallback,
retry, logging, wrapping, rollback, cleanup, or error normalization.

Temporary zip inputs and destinations may be created only under guarded test
temporary directories when needed to characterize the legacy entry flow. Do
not write into the repository, invoke a process or network request, change the
process working directory, or broaden the filesystem recording population.

Do not move `zip.OpenReader`, `os.OpenFile`, `io.Copy`, `outFile.Close`,
`rc.Close`, `f.Open`, path traversal validation, process construction,
`shell.Run`, structurizr, Spring, clock, server, another filesystem caller or
operation, P4, P5, Docker, distribution, or publication work.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete
`pkg/shell/command.go` and related source/tests/callers, complete
`internal/adapter/filesystem` source/tests and every complete double, the
completed Spring archive/discovery/download contracts, relevant
file/HTTP/process/context/config/Maven/structurizr/Bitbucket/Wpost/local-config/
browser/profile/tips code and tests, and the
Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation before editing.
Before the full gate, read the complete launcher contract, relevant Make
meta-tests, P2A API/CLI/subprocess contracts, and all four host acceptance
flows.

# Three Moves

1. Start red with focused private shell contracts proving complete production
   filesystem selection, exact directory-entry and file-parent paths, exact
   `os.ModePerm`, one attempt per reached site, preserved filename append and
   entry ordering, exact error identity and partial result on both failures,
   safe zero behavior at the reached directory boundary, rejection of empty
   recorded populations, and absence of unrelated adapter operations. Use only
   guarded temporary archive fixtures necessary to reach the legacy branches.

2. Add one private complete unzip dependency composition and select
   `filesystem.System()` only in the public production entry. Preserve the
   public signature and move the established body behind the private
   composition without reordering it. Replace only both direct `os.MkdirAll`
   calls with `filesystem.MkdirAll(dependencies.Files, ..., os.ModePerm)`.
   Change no filesystem interface, adapter implementation, other shell
   function, remaining unzip effect, inventory, or completed caller.

3. Run focused shell/filesystem/Spring and relevant command, context, config,
   HTTP, process, Maven, structurizr, profile, browser, tips, file, template,
   Bitbucket, Wpost, local-config, and caller package tests; the launcher
   contract from `/bin/bash`; Make preflight meta-contracts; API/CLI and
   subprocess compatibility; full Go tests and race/vet; all four host
   acceptance flows; the audit meta-suite; focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/
   Q2.1/Q3.4 measurements; full clean checkpoint audit; and empty-HOME count-2.
   Expect nominal Q1.3 to improve from 15 of 36 to 13 of 34 when the two caller
   sites leave and the existing adapter implementation stays singular. Q0.6,
   Q1.1, Q1.2, Q1.4, and exact Q2.1 must hold. Regenerate exact values; the full
   audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent shell unzip
directory-creation move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create one linked
NEXT archive for the next coherent P3 effect move, replace only the launcher's
mutable regions, run the launcher contract, and make the separate handoff-only
commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
