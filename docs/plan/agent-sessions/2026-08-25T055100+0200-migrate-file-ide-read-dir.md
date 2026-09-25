# Agent Session: Migrate Non-Recursive IDE Directory Read

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T055100+0200-migrate-file-ide-read-dir`
Created: `2026-08-25T05:51:00+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `acef2a39adf67431fe1887e8e553a01e20ce749c4b51ec53c986f9e836f5a4a6`
Previous: [2026-08-25T052243+0200-migrate-file-grep-recursive.md](2026-08-25T052243+0200-migrate-file-grep-recursive.md)
Next: [2026-08-25T061802+0200-migrate-config-git-hook-read-dir.md](2026-08-25T061802+0200-migrate-config-git-hook-read-dir.md)
Outcome: Completed in implementation commit `d6cb5939cccefc3fd7ae8e80ef81c6bd5e922342`; the clean checkpoint audit exited 1 for 16 documented findings, never 2, with zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route non-recursive IDE
cleanup's direct directory read through the existing filesystem adapter.
Preserve the exact target path, sorted entry order, entry classification,
ordered dry-run or deletion behavior, counts, report, returned errors,
recursive selection, every completed filesystem, Bitbucket, Wpost, and
supervisor move, and every P2A contract with zero comparable ratchet
regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`ioutil.ReadDir(targetDir)` operation in non-recursive `removeIntellijFile` in
`pkg/file/ide.go` and one matching `ReadDir` operation on the existing
zero-value-safe filesystem adapter. It authorizes the required updates to
every complete filesystem recording double. It does not authorize recursive
IDE cleanup, `FindAll`, `GrepRecursive`, `logAndDelete`, `DeleteAll`, another
file operation or Q1.3 flow, another adapter operation, a new adapter family
or public API, mutation harnesses, or later roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, the complete `pkg/file` implementation and tests,
all callers of `file.RemoveIntellijFiles`, the filesystem adapter and every
complete filesystem recording double, `.quality/inventory`, and the continuity
and quality-lift designs. Regenerate ignored reports outside the measured tree
or remove them before a clean audit.

Implementation commit `d03e96d` has 216 tests across 16 of 25 packages. Q0.6
has 20 guarded safe-writer sites and zero unsafe direct test writes, Q1.1 is 9
of 25, Q1.2 is 0, Q1.3 is 40 violations of 56 production effect sites with
clock and server absent, Q1.4 is 7 of 8, and exact Q2.1 is 0 of 8 executable
harnesses. The clean gate passed, the full audit exited 1 for 16 documented
findings and never 2, and comparable ratchets were five improved, two held,
and zero regressed.

The launcher has 62 Bash 3.2 contracts and supervises fresh non-interactive
JSONL turns with external raw logs. It continues only after a successful
structured stream and valid clean committed handoff. Preserve its stable
skeleton and do not launch a real successor while developing, testing, or
finalizing this move.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the non-recursive IDE directory-read boundary,
adapter contract, complete recording doubles, file contracts, and measured
planning notes. Then perform the normal separate handoff-only commit. Do not
push, merge, publish, distribute, remove the worktree, stash inherited
changes, revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, public API/CLI behavior, all four host acceptance flows, the
exported
`file.RemoveIntellijFiles(targetDir string, recursive bool, dryRun bool) (string, error)`
signature, its command caller, recursive selection and cleanup, `FindAll`,
`GrepRecursive`, `logAndDelete`, `DeleteAll`, `Path`, and every caller's
observable behavior.

Extend the existing complete filesystem dependency with one resolvable
`ReadDir(path string) ([]fs.FileInfo, error)` operation and a safe zero value;
do not store a function-valued effect dependency. Production must use the
same `ioutil.ReadDir` behavior behind `filesystem.System()` so names remain
filename-sorted and the returned metadata and errors remain exact. Pass the
exact `targetDir` through the adapter and iterate the dependency's entries in
their delivered order. Preserve both existing ordered, independent
classification checks: a non-directory whose exact name contains `.iml`, and
a directory whose exact name contains `.idea`. Preserve each exact
`Path("%s/%s", targetDir, f.Name())`, the call to `logAndDelete`, first-error
short circuit, counts, and exact `Iml files: %d, .idea dirs: %d` report. Return
an initial read error unchanged with an empty report. Do not normalize, wrap,
retry, preflight, resort, deduplicate, clean paths, change substring matching,
coalesce conditions, inspect other metadata, change logging or deletion, add
selection behavior, or broaden into recursive cleanup, another file operation,
Bitbucket, HTTP/process effects, config behavior, Wpost, clock, server, P4
adapters, P5 harnesses, dependencies, Docker, cloud, distribution, or
publication.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, `pkg/file`, all callers of
`file.RemoveIntellijFiles`, and `internal/adapter/filesystem`. Read the relevant
tests, every filesystem test double that implements the complete adapter, and
the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation before editing.
Before the full gate, read the complete launcher contract, relevant Make
meta-tests, P2A API/CLI/subprocess contracts, and all four host acceptance
flows.

# Three Moves

1. Start red with adapter and non-recursive IDE cleanup recording contracts.
   Prove the complete dependency, production selection of
   `filesystem.System()`, exact target path, ordered entries and metadata,
   filename sorting in the system adapter, exact read error, exact ordered
   `.iml` file and `.idea` directory substring decisions including ignored
   opposite types, exact dry-run paths, logs, counts and report, safe zero
   behavior without developer-path access, and non-empty recorded directory-read
   and entry populations. Keep fixture mutation confined to guarded temporary
   directories; production-composition contracts must not mutate the real
   filesystem.

2. Extend only the existing filesystem adapter and every complete filesystem
   recording double with `ReadDir`. Keep exported `RemoveIntellijFiles` as the
   production entry, select `filesystem.System()` only for the non-recursive
   helper, and replace only its direct `ioutil.ReadDir(targetDir)` with the
   adapter operation. Do not change the loop body, recursive branch,
   `logAndDelete`, `DeleteAll`, another operation, inventory, or public API.

3. Run focused file/filesystem contracts and the command caller package's
   relevant tests, all adjacent complete-double packages, the launcher contract
   from `/bin/bash`, Make preflight meta-contracts, API/CLI compatibility, full
   Go tests and race/vet, all four host acceptance flows, the audit meta-suite,
   focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean
   checkpoint audit, and empty-HOME count-2. Expect Q1.3 to move nominally from
   40 of 56 to 39 of 56 while Q1.2, Q1.4, and exact Q2.1 hold. Regenerate exact
   values; the full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent non-recursive
IDE directory-read move or record an exact resumable state. Rewrite the
rolling handover, record the measured P3 result, answer this archive, create
one linked NEXT archive for the next coherent P3 effect move, replace only the
launcher's mutable regions, run the launcher contract, and make the separate
handoff-only commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
