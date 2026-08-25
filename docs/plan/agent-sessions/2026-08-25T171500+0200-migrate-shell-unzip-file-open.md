# Agent Session: Migrate Shell Unzip File Open

Status: NEXT
Session ID: `2026-08-25T171500+0200-migrate-shell-unzip-file-open`
Created: `2026-08-25T17:15:00+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `1baacd1f3c5f19a9ccc97a4d2cefed02cbeb194a9acc2846cfa14621a29ce968`
Previous: [2026-08-25T163321+0200-migrate-shell-unzip-directory-creation.md](2026-08-25T163321+0200-migrate-shell-unzip-directory-creation.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the direct
`os.OpenFile` call in the file-entry branch of `pkg/shell.Unzip` through the
established `filesystem.OpenFile` adapter operation. Preserve the public
signature, private production composition, archive-open and traversal order,
filename append order, both completed directory-creation boundaries,
directory short-circuit, parent-directory selection, exact open flags and
entry mode, partial results, error identity and precedence, every remaining
unzip operation, every completed filesystem, Spring, process, browser,
profile, HTTP, tips, config, Maven, structurizr, Bitbucket, Wpost,
local-config, and supervisor move, and every P2A contract with zero comparable
ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only reuse of the
private complete shell unzip composition containing
`filesystem.Dependencies`, its existing production selection of
`filesystem.System()`, focused recording of the existing `OpenFile`
operation, and replacement of the one direct `os.OpenFile` call with the
existing forwarding helper. It does not authorize a filesystem interface or
adapter implementation change, another filesystem operation or caller,
direct `zip.OpenReader`, `f.Open`, `io.Copy`, or close migration, process
work, another shell function, Spring, clock/server adapters, public API,
inventory, mutation harnesses, or later roadmap implementation.

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

Implementation commit `b9fe209` has 330 tests across 19 of 25 packages. Q0.6
has 25 guarded safe-writer sites, 20 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 6 of 25, Q1.2 is 0, Q1.3 is 13 violations of 34 production
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
implementation commit for the one unzip file-open effect and focused private
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
error text, filename append before mutation, both established
`filesystem.MkdirAll` calls with exact paths and `os.ModePerm`, directory
`continue`, `filepath.Dir(fpath)`, file open flags
`os.O_WRONLY|os.O_CREATE|os.O_TRUNC`, entry mode, direct entry open, copy,
file close, entry close, final results, and every error precedence.

Add no filesystem interface operation and change no adapter implementation.
Use the established `filesystem.OpenFile` helper with the exact joined path,
flags, and entry mode. Preserve exact safe-zero behavior in the filesystem
adapter and at the earlier reached directory boundary. Reuse the one private
complete production composition and its single `filesystem.System()`
selection. Only the focused shell double may add recording for this existing
operation; every pre-existing complete double and focused behavior stays
unchanged.

For a directory entry, preserve the exact joined path append, one MkdirAll
attempt, and direct continue without a file-open attempt. For a file entry,
preserve the exact joined path append and successful parent MkdirAll before
one OpenFile attempt. Return the current partial filenames plus the exact open
error before `f.Open`, copy, or close on failure. Do not clean or rewrite paths
beyond the existing join, clean, prefix, and parent operations. Do not add
fallback, retry, logging, wrapping, rollback, cleanup, or error normalization.

Temporary zip inputs and destinations may be created only under guarded test
temporary directories when needed to characterize the legacy entry flow. Do
not write into the repository, invoke a process or network request, change the
process working directory, or record any filesystem operation other than the
already-recorded MkdirAll and newly authorized OpenFile.

Do not move `zip.OpenReader`, `io.Copy`, `outFile.Close`, `rc.Close`, `f.Open`,
path traversal validation, process construction, `shell.Run`, structurizr,
Spring, clock, server, another filesystem caller or operation, P4, P5, Docker,
distribution, or publication work.

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

1. Start red by extending only the focused private shell double and contracts
   to prove the exact joined file path, exact
   `os.O_WRONLY|os.O_CREATE|os.O_TRUNC` flags, exact entry mode, one parent
   MkdirAll followed by one file-open attempt, preserved filename append and
   entry order, exact open error identity and partial result, no later
   file-entry operation after failure, directory-entry no-open short-circuit,
   rejection of an empty recorded OpenFile population, and absence of
   unrelated adapter operations. Use only the existing guarded temporary
   archive fixtures necessary to reach the legacy branch.

2. Reuse the private complete unzip dependency composition and its single
   public production selection of `filesystem.System()`. Preserve the public
   signature and established body order. Replace only direct `os.OpenFile`
   with `filesystem.OpenFile(dependencies.Files, ...)`. Change no filesystem
   interface, adapter implementation, other shell function, remaining unzip
   effect, inventory, or completed caller.

3. Run focused shell/filesystem/Spring and relevant command, context, config,
   HTTP, process, Maven, structurizr, profile, browser, tips, file, template,
   Bitbucket, Wpost, local-config, and caller package tests; the launcher
   contract from `/bin/bash`; Make preflight meta-contracts; API/CLI and
   subprocess compatibility; full Go tests and race/vet; all four host
   acceptance flows; the audit meta-suite; focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/
   Q2.1/Q3.4 measurements; full clean checkpoint audit; and empty-HOME count-2.
   Expect nominal Q1.3 to improve from 13 of 34 to 12 of 33 when the one caller
   site leaves and the existing adapter implementation stays singular. Q0.6,
   Q1.1, Q1.2, Q1.4, and exact Q2.1 must hold. Regenerate exact values; the
   full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent shell unzip
file-open move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create one linked
NEXT archive for the next coherent P3 effect move, replace only the launcher's
mutable regions, run the launcher contract, and make the separate handoff-only
commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
