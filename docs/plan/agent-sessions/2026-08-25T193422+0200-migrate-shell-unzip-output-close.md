# Agent Session: Migrate Shell Unzip Output Close

Status: NEXT
Session ID: `2026-08-25T193422+0200-migrate-shell-unzip-output-close`
Created: `2026-08-25T19:34:22+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `043c870dae7f5888815f748604e4505bed8da5f6bf64589b4a5c5255c39110f2`
Previous: [2026-08-25T185749+0200-migrate-shell-unzip-archive-open.md](2026-08-25T185749+0200-migrate-shell-unzip-archive-open.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the direct
`err = outFile.Close()` request in private `pkg/shell.Unzip` through one narrow
Close operation on the existing filesystem adapter. Preserve the public
signature, exact opened output identity, close attempt count and placement,
archive and entry traversal, entry order and bytes, filenames and partial
results, zip-slip behavior, every completed Unzip operation and error
precedence, and every earlier shell, filesystem, Spring, process, browser,
profile, HTTP, tips, config, Maven, structurizr, Bitbucket, Wpost,
local-config, Kibana, and supervisor move, with zero comparable ratchet
regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only one
`Close(filesystem.File) error` method on the existing
`filesystem.FileSystem`, its zero-safe forwarding helper and exact system
implementation, the mechanical method addition required by every complete
filesystem double, focused adapter and private Unzip output-close recording
contracts, and replacement of the one direct `outFile.Close()` call with the
helper.

It does not authorize archive open or deferred Close, entry-reader Close,
Copy, OpenFile, MkdirAll, entry-open injection, another unzip branch, shell
Run or Git changes, process or HTTP work, Spring, Maven, structurizr, plugin
diagrams, clock/server adapters, public API, inventory, mutation harnesses,
audit reshaping, or later-roadmap implementation.

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

Implementation commit `2f0a072` has 346 tests across 19 of 25 packages. Q0.6
has 26 guarded safe-writer sites, 21 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 6 of 25, Q1.2 is 0, Q1.3 is 10 violations of 32
production effect sites with clock and server absent, Q1.4 is 7 of 8, and
exact Q2.1 is 0 of 8 executable harnesses. The clean full audit exited 1 for
15 documented findings and never 2; comparable ratchets were five improved,
two held, and zero regressed.

The direct archive-open caller is gone and exact Q1.3 improved to 10/32. The
scanner previously retained filesystem provenance through Copy's opened-file
and reader arguments. It may likewise retain provenance through a Close helper
receiving the opened file while the new in-boundary system implementation adds
one production site. Regenerate the exact value without changing the scanner
or broadening this move to force an expected number.

The launcher has 62 Bash 3.2 contracts and supervises fresh non-interactive
JSONL turns with external raw logs. It continues only after a successful
structured stream and valid clean committed handoff. Preserve its stable
skeleton and do not launch a real successor while developing or finalizing.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the Unzip output-close move, adapter contracts, and
focused private recording contracts, with measured planning notes. Then make
the normal separate handoff-only commit. Do not push, merge, publish,
distribute, remove the worktree, stash inherited changes, revert user work, or
run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve public
`Unzip(string, string) ([]string, error)`, the private complete
`unzipDependencies`, production `filesystem.System()` selection, exact opened
output identity, one output-close request per reached file after Copy, the
assignment and following error branch, exact close error and partial
filenames, suppression of entry-reader Close and later traversal after an
output-close error, and continued entry Close and later traversal after
success.

Reuse the existing zero-safe adapter shape and pass the complete Unzip
dependency composition whole. A zero filesystem dependency must return
`filesystem.ErrNoFilesystem` without closing a developer-supplied file. Do not
add fallback, retry, logging, cleanup, new defers, environment or
working-directory behavior, or error normalization. Any fixture write must
stay guarded below `t.TempDir()`; launch no process, touch no network, and
change no working directory.

Do not change `OpenZipReader`, the deferred archive Close, `filesystem.Copy`,
`rc.Close()`, OpenFile, MkdirAll, entry open, shell Run, Git composition,
structurizr command construction, Maven process composition, plugin diagrams,
another adapter operation or caller, public API, or `.quality/inventory`.

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

1. Start red with only focused filesystem-adapter and private Unzip
   output-close recording contracts. Prove exact file identity, one request,
   exact injected error and partial filenames, Copy-before-close order,
   success traversal, error short-circuit before entry Close and later entries,
   complete dependency selection, zero-safe behavior, rejection of an empty
   recording population, and absence of unrelated filesystem requests. Invoke
   no real process or network request.

2. Add only `Close(filesystem.File) error` to the existing filesystem
   interface, zero-safe helper, exact system implementation, and complete
   doubles. Replace only `err = outFile.Close()` with
   `err = filesystem.Close(dependencies.Files, outFile)`. Keep the surrounding
   Copy expression, assignment, error branch, entry Close, and all earlier and
   later operations unchanged.

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

Before this agent session ends, finish and commit the coherent Unzip
output-close move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create one linked
NEXT archive for the next coherent P3 effect move, replace only the launcher's
mutable regions, run the launcher contract, and make the separate handoff-only
commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
