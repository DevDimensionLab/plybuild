# Agent Session: Migrate Shell Unzip Copy

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T181827+0200-migrate-shell-unzip-copy`
Created: `2026-08-25T18:18:27+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `59b9f50aa0dc33690dc9d4026c684344bea9b197583570bddfec0fa6d76de322`
Previous: [2026-08-25T174800+0200-migrate-shell-run-process.md](2026-08-25T174800+0200-migrate-shell-run-process.md)
Next: [2026-08-25T185749+0200-migrate-shell-unzip-archive-open.md](2026-08-25T185749+0200-migrate-shell-unzip-archive-open.md)
Outcome: Move 50 committed as `89918bd`; the direct `io.Copy` is routed through the existing filesystem Copy adapter with 339 tests and zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the direct
`io.Copy(outFile, rc)` payload transfer inside the file-entry branch of public
`pkg/shell.Unzip` through the established `filesystem.Copy` adapter operation.
Preserve the public signature, ordered filenames and partial results, exact
archive-entry bytes and destination identity, archive and entry traversal,
copy-result handling, file and entry close order and errors, every earlier
error precedence, and every completed shell Run, unzip, filesystem, Spring,
process, browser, profile, HTTP, tips, config, Maven, structurizr, Bitbucket,
Wpost, local-config, and supervisor move, with zero comparable ratchet
regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes focused reshaping
of the existing private unzip recording filesystem only to record the already
established `Copy(filesystem.File, io.Reader) (int64, error)` request and
replacement of the one direct file-entry `io.Copy` call with
`filesystem.Copy(dependencies.Files, outFile, rc)`.

It does not authorize a filesystem interface or adapter implementation change,
another filesystem operation or caller, archive-open injection, output-file or
entry-reader close injection, another unzip branch, shell Run or Git changes,
process or HTTP work, Spring, Maven, structurizr, plugin diagrams, clock/server
adapters, public API, inventory, mutation harnesses, or later-roadmap
implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, complete `pkg/shell` source/tests/callers, complete
`internal/adapter/filesystem` source/tests and every complete filesystem double,
all unzip callers and completed unzip contracts, relevant Spring download,
command, context, config, HTTP, process, Maven, structurizr, Bitbucket, Wpost,
local-config, browser, profile, tips, file, template, and Kibana code/tests,
`.quality/inventory`, both design documents, and the focused audit
implementation. Regenerate ignored reports outside the measured tree or remove
them before a clean audit.

Implementation commit `da7eebf` has 336 tests across 19 of 25 packages. Q0.6
has 25 guarded safe-writer sites, 20 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 6 of 25, Q1.2 is 0, Q1.3 is 11 violations of 32 production
effect sites with clock and server absent, Q1.4 is 7 of 8, and exact Q2.1 is 0
of 8 executable harnesses. The clean full audit exited 1 for 15 documented
findings and never 2; comparable ratchets were five improved, two held, and
zero regressed.

The launcher has 62 Bash 3.2 contracts and supervises fresh non-interactive
JSONL turns with external raw logs. It continues only after a successful
structured stream and valid clean committed handoff. Preserve its stable
skeleton and do not launch a real successor while developing or finalizing.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the unzip file-entry copy move and focused private
recording contracts, with measured planning notes. Then make the normal
separate handoff-only commit. Do not push, merge, publish, distribute, remove
the worktree, stash inherited changes, revert user work, or run destructive Git
commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve public
`Unzip(string, string) ([]string, error)`, the private complete
`unzipDependencies`, production `filesystem.System()` selection, exact joined
paths and zip-slip behavior, filename append order, directory short-circuit,
parent creation, file open flags and mode, entry open, and one payload-transfer
attempt per reached file entry. Characterize the legacy ignored `io.Copy`
count and error before choosing assignments: do not return, wrap, log, or give
precedence to the adapter copy error if the current flow ignores it. Preserve
the following output-file Close and entry-reader Close attempts, their order,
their exact errors and precedence, and later entry traversal.

Reuse the existing zero-safe filesystem adapter and its complete `Copy`
operation. Pass the existing complete unzip dependency composition whole. Do
not add fallback, retry, buffering, extra logging, cleanup, deferred closes,
new environment or working-directory behavior, or error normalization. The
focused double may record only existing unzip filesystem requests, compare the
exact destination identity, read delivered source bytes in memory, and return
an arbitrary count/error without writing the destination. Any archive fixture
write must stay guarded below `t.TempDir()`; launch no process, touch no
network, and change no working directory.

Do not change `zip.OpenReader`, either direct Close, OpenFile, MkdirAll, shell
Run, Git composition, structurizr command construction, Maven process
composition, plugin diagrams, the filesystem adapter, another caller, public
API, or `.quality/inventory`.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete `pkg/shell` source
and tests, every Unzip caller, complete `internal/adapter/filesystem` source and
tests and every complete double, completed unzip/Spring/filesystem/process
contracts, and relevant caller packages before editing. Before the full gate,
read the complete launcher contract, Make meta-tests, P2A API/CLI and subprocess
contracts, and all four host acceptance flows.

# Three Moves

1. Start red with only focused private unzip Copy recording contracts. Prove
   the exact opened destination identity, one Copy request for each reached
   file entry, exact ordered arbitrary entry bytes delivered through the
   source reader, established operation order after parent creation, OpenFile,
   and entry open, legacy ignored copy count/error behavior, continued close
   and later-entry behavior, rejection of an empty Copy recording population,
   no Copy for directories or earlier failures, and absence of unrelated
   filesystem requests. Invoke no real process or network request.

2. Change only the direct `io.Copy(outFile, rc)` expression to the existing
   `filesystem.Copy(dependencies.Files, outFile, rc)` helper. Preserve the
   ignored two-value result exactly, keep the surrounding operation and return
   order unchanged, and change no adapter, dependency selection, caller, or
   other effect.

3. Run focused shell/filesystem/Spring and relevant process, Maven, command,
   context, config, HTTP, structurizr, profile, browser, tips, file, template,
   Bitbucket, Wpost, local-config, Kibana, and caller package tests; the launcher
   contract from `/bin/bash`; Make preflight meta-contracts; API/CLI and
   subprocess compatibility; full Go tests and race/vet; all four host
   acceptance flows; the audit meta-suite; focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/
   Q2.1/Q3.4 measurements; full clean checkpoint audit; and empty-HOME count-2.
   Expect nominal Q1.3 to improve from 11 of 32 to 10 of 31 when the one direct
   copy caller site leaves and the existing adapter implementation stays
   singular. Regenerate exact values; the full audit may exit 1 for documented
   findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent unzip Copy move
or record an exact resumable state. Rewrite the rolling handover, record the
measured P3 result, answer this archive, create one linked NEXT archive for the
next coherent P3 effect move, replace only the launcher's mutable regions, run
the launcher contract, and make the separate handoff-only commit
`docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
