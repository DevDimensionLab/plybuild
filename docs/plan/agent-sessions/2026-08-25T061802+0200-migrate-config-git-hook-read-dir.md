# Agent Session: Migrate Git Hook Directory Read

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T061802+0200-migrate-config-git-hook-read-dir`
Created: `2026-08-25T06:18:02+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `58ea70f7d74041aa7c91d3522f66688ea9aee52a0f433bf6bb02adc3a64449b1`
Previous: [2026-08-25T055100+0200-migrate-file-ide-read-dir.md](2026-08-25T055100+0200-migrate-file-ide-read-dir.md)
Next: [2026-08-25T064347+0200-migrate-config-examples-read-dir.md](2026-08-25T064347+0200-migrate-config-examples-read-dir.md)
Outcome: Completed in implementation commit `82e631d`; the clean checkpoint audit exited 1 for 16 documented findings, never 2, with zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route
`GitCloudConfig.GitHookFiles`' direct directory read through the existing
filesystem adapter. Preserve the exact receiver and root construction, sorted
entry order, entry classification, filename-only ordered results, nil result
behavior, returned errors, every completed filesystem, file, Bitbucket, Wpost,
and supervisor move, and every P2A contract with zero comparable ratchet
regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`ioutil.ReadDir(root)` operation in `GitCloudConfig.GitHookFiles` in
`pkg/config/cloud.go` and its use of the existing zero-value-safe filesystem
`ReadDir` operation. It authorizes focused config recording contracts. It does
not authorize `Examples`, `Templates`, another cloud or config method, another
filesystem operation, an adapter or complete-double extension, another
production effect, a new adapter family or public API, mutation harnesses, or
later roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, the complete `pkg/config` implementation and tests,
all callers of `GitHookFiles`, the complete `CloudConfig` interface and its
implementations or doubles, the filesystem adapter and relevant recording
doubles, `.quality/inventory`, and the continuity and quality-lift designs.
Regenerate ignored reports outside the measured tree or remove them before a
clean audit.

Implementation commit `d6cb593` has 224 tests across 16 of 25 packages. Q0.6
has 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 9 of 25, Q1.2 is 0, Q1.3 is 39 violations of 56 production
effect sites with clock and server absent, Q1.4 is 7 of 8, and exact Q2.1 is 0
of 8 executable harnesses. The clean gate passed, the full audit exited 1 for
16 documented findings and never 2, and comparable ratchets were five
improved, two held, and zero regressed.

The launcher has 62 Bash 3.2 contracts and supervises fresh non-interactive
JSONL turns with external raw logs. It continues only after a successful
structured stream and valid clean committed handoff. Preserve its stable
skeleton and do not launch a real successor while developing, testing, or
finalizing this move.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the Git-hook directory-read boundary, config
contracts, and measured planning notes. Then perform the normal separate
handoff-only commit. Do not push, merge, publish, distribute, remove the
worktree, stash inherited changes, revert user work, or run destructive Git
commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, public API/CLI behavior, all four host acceptance flows,
`CloudConfig`, the exported
`GitCloudConfig.GitHookFiles(path string) ([]string, error)` method signature,
its command caller, `GitCloudConfig.Implementation`, `DirConfig.Dir`,
`file.Path`, `Examples`, `Templates`, and every caller's observable behavior.

Use the existing complete filesystem dependency and existing
`ReadDir(path string) ([]fs.FileInfo, error)` operation with a safe zero value;
do not add or change an adapter operation and do not store a function-valued
effect dependency. Production must select `filesystem.System()` only for the
Git-hook-files composition. Preserve exact evaluation of
`root := file.Path("%s/%s", gitCfg.Implementation().Dir(), path)` and pass that
exact root through the adapter. Iterate entries in their delivered order;
production sorting remains the adapter's established `ioutil.ReadDir`
behavior. Preserve the exact `!f.IsDir()` decision and append only each exact
`f.Name()` to the result. An initial read error must return an unchanged nil
result and the exact error. An empty directory or all-directory population
must return a nil result and nil error. Do not normalize, wrap, retry,
preflight, resort, deduplicate, clean paths, change classification, return full
paths, inspect other metadata, or broaden into `Examples`, `Templates`, another
config method, command behavior, file behavior, Bitbucket, HTTP/process
effects, Wpost, clock, server, P4 adapters, P5 harnesses, dependencies, Docker,
cloud distribution, or publication.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete `pkg/config`, all
callers of `GitHookFiles`, the `CloudConfig` interface and every implementation
or double, and `internal/adapter/filesystem`. Read the relevant tests and the
Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation before editing. Before
the full gate, read the complete launcher contract, relevant Make meta-tests,
P2A API/CLI/subprocess contracts, and all four host acceptance flows.

# Three Moves

1. Start red with focused Git-hook-files recording contracts. Prove the
   complete dependency, production selection of `filesystem.System()`, exact
   receiver-derived root, delivered entry order and metadata, exact read
   error, exact non-directory classification including ignored directories,
   exact filename-only ordered results, nil results for empty and all-directory
   populations, safe zero behavior without developer-path access, and
   non-empty recorded directory-read and entry populations. Production
   composition contracts must not mutate the real filesystem.

2. Reuse only the existing filesystem adapter `ReadDir` operation. Keep
   exported `GitHookFiles` as the production entry, select
   `filesystem.System()` only for its private complete composition, and replace
   only its direct `ioutil.ReadDir(root)` with the adapter operation. Do not
   change the root expression or loop body, extend the adapter or complete
   doubles, change `CloudConfig`, touch another config method or caller, or
   change inventory or public API.

3. Run focused config/filesystem contracts and the command caller package's
   relevant tests, the launcher contract from `/bin/bash`, Make preflight
   meta-contracts, API/CLI compatibility, full Go tests and race/vet, all four
   host acceptance flows, the audit meta-suite, focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean checkpoint
   audit, and empty-HOME count-2. Expect Q1.3 to move nominally from 39 of 56
   to 38 of 55 while Q0.6, Q1.2, Q1.4, and exact Q2.1 hold. Regenerate exact
   values; the full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Git-hook
directory-read move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create one linked
NEXT archive for the next coherent P3 effect move, replace only the launcher's
mutable regions, run the launcher contract, and make the separate handoff-only
commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
