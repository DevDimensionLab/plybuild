# Agent Session: Migrate Spring Archive Working Directory

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T155409+0200-migrate-spring-archive-working-directory`
Created: `2026-08-25T15:54:09+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `8755537d5ed21a2f39f7c5c684336f8a2db5e832c1098171740a5b05ebe3efd7`
Previous: [2026-08-25T151736+0200-migrate-open-browser-process-start.md](2026-08-25T151736+0200-migrate-open-browser-process-start.md)
Next: [2026-08-25T163321+0200-migrate-shell-unzip-directory-creation.md](2026-08-25T163321+0200-migrate-shell-unzip-directory-creation.md)
Outcome: Completed by implementation commit `d982f63`; Q1.3 improved to 15/36 with 324 tests and zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the direct
`os.Getwd()` in private `pkg/spring.archivePath` through one new distinct
working-directory operation on the existing filesystem adapter. Preserve the
private signature, exact Getwd-before-time ordering, empty path and exact error
on failure, direct `time.Now().Unix()` selection, exact `file.Path` format,
every caller and Spring operation, every completed process, browser, profile,
filesystem, HTTP, tips, config, Maven, structurizr, Bitbucket, Wpost,
local-config, and supervisor move, and every P2A contract with zero comparable
ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only one complete
`WorkingDirectory() (string, error)` addition to the existing internal
`filesystem.FileSystem` interface, its forwarding helper, safe-zero result,
direct system `os.Getwd()` implementation, focused adapter contracts,
mechanical completeness updates to existing filesystem doubles, one private
complete Spring archive-path composition, focused archive-path recording
contracts, and replacement of only the direct `os.Getwd()` call. It does not
authorize another filesystem caller or operation, another Spring function,
the adjacent clock call, a clock/server adapter, process work, public API,
inventory, mutation harnesses, or later roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, complete `pkg/spring/io.go` and its package tests
and callers, complete `internal/adapter/filesystem` source/tests, every complete
`filesystem.FileSystem` double and dependency caller, the completed Spring
discovery/download contracts, relevant file, HTTP, shell, process, config,
Maven, structurizr, context, Bitbucket, Wpost, local-config, browser-launcher,
and profile code/tests, `.quality/inventory`, the continuity and quality-lift
designs, and the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation.
Regenerate ignored reports outside the measured tree or remove them before a
clean audit.

Implementation commit `9fcdfb5` has 316 tests across 19 of 25 packages. Q0.6
has 24 guarded safe-writer sites, 19 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 6 of 25, Q1.2 is 0, Q1.3 is 16 violations of 36 production
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
implementation commit for exact working-directory representation and the
Spring archive-path boundary and contracts, with measured planning notes. Then
perform the normal separate handoff-only commit. Do not push, merge, publish,
distribute, remove the worktree, stash inherited changes, revert user work, or
run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve private `archivePath() (string, error)`,
every caller, `CreateSpringInitializer`, Spring discovery, download, unzip,
delete and demo-file behavior, base URL and request construction, exact direct
`time.Now().Unix()`, exact path format `file.Path("%s/spring-%d.zip", curDir,
now)`, public API/CLI behavior, all four host acceptance flows, process start
mode, browser and profile behavior, and every completed caller's observable
behavior.

Add only `WorkingDirectory() (string, error)` to the complete filesystem
interface. Its forwarding helper must return `"", filesystem.ErrNoFilesystem`
for zero dependencies, otherwise make one interface attempt and return its
exact string and error. The system implementation must directly return
`os.Getwd()` with no cleaning, normalization, fallback, retry, logging, or error
wrapping. Keep every established filesystem method, dependency selection,
error sentinel, system operation, and caller unchanged. Update every existing
complete filesystem double only to remain complete and reject the unrelated
operation; only focused adapter and Spring doubles may record it.

Keep `archivePath` as the production entry and select `filesystem.System()`
only for one private complete composition. On Getwd error it must return the
same empty path and exact error before evaluating time or formatting a path. On
success it must pass the arbitrary working-directory bytes unchanged to the
established `file.Path` call, evaluate direct `time.Now().Unix()` in the same
position, and return the same `spring-<Unix>.zip` path with nil error. Make one
working-directory attempt. Do not create, stat, read, write, or delete a file;
change the process working directory; invoke a shell or process; add a clock
seam; or change error selection.

Do not clean, parse, validate, resolve, or rewrite the delivered directory;
change the archive filename or timestamp unit; move `time.Now`; use `os.UserHomeDir`;
broaden into Spring download/discovery, shell.Unzip, file.Path, HTTP, Maven,
structurizr, browser, profile, another filesystem operation, P4, P5, Docker,
distribution, or publication work.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete `pkg/spring/io.go`
and related source/tests/callers, complete `internal/adapter/filesystem`
source/tests and every complete double, the completed Spring discovery and
download contracts, relevant file/shell/HTTP/process/context/config/Maven/
structurizr/Bitbucket/Wpost/local-config/browser/profile code and tests, and
the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation before editing.
Before the full gate, read the complete launcher contract, relevant Make
meta-tests, P2A API/CLI/subprocess contracts, and all four host acceptance
flows.

# Three Moves

1. Start red by extending focused filesystem adapter contracts to prove exact
   arbitrary working-directory values and errors reach and return from one
   recording attempt, the direct system operation matches `os.Getwd`, zero
   dependencies return exact `ErrNoFilesystem`, existing operations stay
   unchanged, and empty recorded populations fail. Add focused private Spring
   contracts for complete production selection, one attempt, arbitrary
   directory bytes, exact error identity and empty result on error, preserved
   filename/timestamp construction, safe zero behavior, non-empty populations,
   and no real file, HTTP, process, download, unzip, delete, or clock seam.

2. Add only the working-directory interface operation, forwarder, safe-zero
   branch, and direct system implementation. Keep `archivePath` as the
   production entry, select `filesystem.System()` only for a private complete
   composition, preserve its signature and exact sequencing, and replace only
   `os.Getwd()`. Update complete filesystem doubles mechanically without
   changing their focused behaviors. Do not change another Spring function,
   filesystem caller/operation, inventory, public API, or completed effect.

3. Run focused Spring/filesystem and relevant command, context, config, HTTP,
   process, Maven, shell, structurizr, profile, browser, tips, file, template,
   Bitbucket, Wpost, local-config, and caller package tests; the launcher
   contract from `/bin/bash`; Make preflight meta-contracts; API/CLI and
   subprocess compatibility; full Go tests and race/vet; all four host
   acceptance flows; the audit meta-suite; focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/
   Q2.1/Q3.4 measurements; full clean checkpoint audit; and empty-HOME count-2.
   Expect nominal Q1.3 to improve from 16 of 36 to 15 of 36 when the caller
   site leaves and one exact adapter site enters. Q0.6, Q1.1, Q1.2, Q1.4, and
   exact Q2.1 must hold. Regenerate exact values; the full audit may exit 1 for
   documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent filesystem
working-directory and Spring archive-path move or record an exact resumable
state. Rewrite the rolling handover, record the measured P3 result, answer this
archive, create one linked NEXT archive for the next coherent P3 effect move,
replace only the launcher's mutable regions, run the launcher contract, and
make the separate handoff-only commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
