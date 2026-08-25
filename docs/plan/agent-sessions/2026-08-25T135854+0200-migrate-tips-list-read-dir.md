# Agent Session: Migrate Tips List Read Dir

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T135854+0200-migrate-tips-list-read-dir`
Created: `2026-08-25T13:58:54+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `4b20e15b9523ced0ebfe241aa0ba7a8595f2b1ad5e9869d7d7dfccbca89536ae`
Previous: [2026-08-25T131820+0200-migrate-local-config-directory-mkdir.md](2026-08-25T131820+0200-migrate-local-config-directory-mkdir.md)
Next: [2026-08-25T143403+0200-migrate-profile-editor-process.md](2026-08-25T143403+0200-migrate-profile-editor-process.md)
Outcome: completed at `85f4c2b`; Q1.1 improved from 8/25 to 7/25 and Q1.3 from 22/40 to 21/40 with zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the direct
`os.ReadDir(LocalDir(gitCfg))` in `tips.List` through one new exact filesystem
adapter `ReadDirEntries` operation. Preserve exact local tips-directory
selection, the exported `[]os.DirEntry` result, `os.ReadDir` entry and sorting
semantics, error behavior, delivered-entry order and identity, directory and
`.md` filtering, every caller, the completed tips-show read, every completed
filesystem, file, template, config, Maven, structurizr, Bitbucket, Wpost,
local-config, and supervisor move, and every P2A contract with zero comparable
ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`os.ReadDir(LocalDir(gitCfg))` operation in `List` in `pkg/tips/tips.go`; one
distinct zero-value-safe internal filesystem adapter
`ReadDirEntries(path string) ([]fs.DirEntry, error)` operation with exact
`os.ReadDir` system semantics; its focused adapter contracts; the mandatory
method addition to every current complete `filesystem.FileSystem` test double;
and focused tips-list recording and filtering contracts. It does not authorize
changing the established `ReadDir` operation returning `[]fs.FileInfo`, another
tips or command method or effect, another adapter operation or family, a
function-valued effect dependency, a public API, inventory, mutation harnesses,
or later roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, complete `pkg/tips/tips.go` and all tips tests, every
caller and use of `tips.List`, `tips.LocalDir`, `TipsDir`, and `tipsShowCmd`,
complete `cmd/tips.go` and its registration and tests, the relevant config
`CloudConfig`, `Directory`, and local terminal behavior, complete
`internal/adapter/filesystem` source and tests, every complete
`filesystem.FileSystem` implementation and recording double, relevant
filesystem/file/template/config/Maven/structurizr/Bitbucket/Wpost/local-config
code and tests, `internal/adapter/process`, `.quality/inventory`, the continuity
and quality-lift designs, and the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit
implementation. Regenerate ignored reports outside the measured tree or remove
them before a clean audit.

Implementation commit `b2d37cc` has 296 tests across 17 of 25 packages. Q0.6
has 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 8 of 25, Q1.2 is 0, Q1.3 is 22 violations of 40 production
effect sites with clock and server absent, Q1.4 is 7 of 8, and exact Q2.1 is 0
of 8 executable harnesses. The clean gate passed, the full audit exited 1 for
15 documented findings and never 2, and comparable ratchets were five improved,
two held, and zero regressed.

The launcher has 62 Bash 3.2 contracts and supervises fresh non-interactive
JSONL turns with external raw logs. It continues only after a successful
structured stream and valid clean committed handoff. Preserve its stable
skeleton and do not launch a real successor while developing, testing, or
finalizing this move.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the exact filesystem directory-entry read operation,
mandatory complete-double compatibility, tips-list boundary, contracts, and
measured planning notes. Then perform the normal separate handoff-only commit.
Do not push, merge, publish, distribute, remove the worktree, stash inherited
changes, revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, all command objects and registrations, public API/CLI behavior,
all four host acceptance flows, `tips.List`, `tips.LocalDir`, `TipsDir`, the
public `[]os.DirEntry` result, `tipsShowCmd`, local terminal width behavior,
config types, Maven, structurizr, local-config behavior, and every caller's
observable behavior.

Add `ReadDirEntries(string) ([]fs.DirEntry, error)` to the existing complete
`filesystem.FileSystem` interface, add the matching zero-value-safe package
operation returning exact `filesystem.ErrNoFilesystem`, and map the system
implementation directly to `os.ReadDir`. Keep established `ReadDir`, its
`ioutil.ReadDir` mapping, result type, sorting, and contracts unchanged. Update
all 35 current complete filesystem test doubles only as required to remain
complete: the new method must record it for its dedicated contracts or reject
it consistently as unrelated everywhere else. Do not store a function-valued
dependency, split the complete dependency, convert entries to `fs.FileInfo`,
call `Info`, or expose application API.

Preserve exact `LocalDir(gitCfg)` evaluation and its `file.Path("%s/%s", ...)`
behavior, then select `filesystem.System()` only for a private complete tips-
directory-entry-read composition and make one adapter read attempt with that
exact arbitrary string. Preserve exact returned entry objects and error at the
adapter boundary, including partial entries paired with an error. In
`tips.List`, any read error must still return the nil named result with that
exact error without inspecting entries. On success, inspect delivered entries
once in delivered order with the same short-circuit predicate: directories are
excluded without reading their names, non-directories are included only when
their exact name has the case-sensitive suffix `.md`, included entry identities
are unchanged, and empty or all-filtered populations return nil.

Do not clean or join the path, restat entries, call `DirEntry.Info`, reorder,
deduplicate, normalize suffixes, accept directories or `.MD`, wrap an error,
retry, preflight, change `LocalDir`, change command rendering or file reads, or
broaden into config, local-config, Maven, structurizr, template, file, HTTP,
Wpost, process, clock, server, P4, P5, dependency, Docker, distribution, or
publication work.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete `pkg/tips` source
and tests, every tips-list caller, complete `cmd/tips.go` and relevant command
tests, relevant config and terminal-width behavior, complete
`internal/adapter/filesystem` source and tests, every complete filesystem
implementation and double, relevant filesystem/file/template/config/Maven/
structurizr/Bitbucket/Wpost/local-config code and tests,
`internal/adapter/process`, and the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit
implementation before editing. Before the full gate, read the complete launcher
contract, relevant Make meta-tests, P2A API/CLI/subprocess contracts, and all
four host acceptance flows.

# Three Moves

1. Start red with focused filesystem adapter directory-entry-read contracts.
   Prove exact arbitrary paths, exact delivered `fs.DirEntry` order and
   identity, non-empty recorded populations, one attempt, exact nil, empty,
   representative, partial-entry, and arbitrary-error results, safe zero
   behavior without developer-path access, and system `os.ReadDir` semantics in
   temporary directories: filename sorting, entry names and directory
   classification, and exact missing-path failure. Add focused tips-list
   recording contracts that prove the complete dependency, production
   `filesystem.System()` selection, exact `LocalDir` path, one read attempt,
   exact error identity and nil-on-error result, delivered-order filtering and
   entry identity, short-circuit method use, nil empty/all-filtered results,
   non-empty populations, and no unrelated operation. Production-composition
   contracts must not read the real filesystem or run another command.

2. Add only the exact adapter `ReadDirEntries` operation and mandatory method
   implementations on all current complete test doubles. Keep exported
   `tips.List(gitCfg config.CloudConfig) ([]os.DirEntry, error)` as the
   production entry, select `filesystem.System()` only for a private complete
   tips-directory-entry-read composition, and replace only its direct
   `os.ReadDir`. Do not change prior evaluation, filtering, nilness, error
   selection, another tips or command effect, adapter operation, method,
   caller, inventory, or public API.

3. Run focused tips/filesystem contracts and relevant process, command, config,
   file, template, Maven, structurizr, Bitbucket, HTTP, local-config, and caller
   packages' tests, the launcher contract from `/bin/bash`, Make preflight
   meta-contracts, API/CLI compatibility, full Go tests and race/vet, all four
   host acceptance flows, the audit meta-suite, focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean checkpoint
   audit, and empty-HOME count-2. Expect Q1.1 to improve from 8 to 7 of 25 when
   `pkg/tips` gains direct tests and Q1.3 to improve from 22 of 40 to 21 of 40
   because the direct read moves into the exact adapter while its required
   system call remains in the production population. Q0.6, Q1.2, Q1.4, and
   exact Q2.1 must hold. Regenerate exact values; the full audit may exit 1 for
   documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent filesystem
directory-entry-read and tips-list move or record an exact resumable state.
Rewrite the rolling handover, record the measured P3 result, answer this
archive, create one linked NEXT archive for the next coherent P3 effect move,
replace only the launcher's mutable regions, run the launcher contract, and
make the separate handoff-only commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
