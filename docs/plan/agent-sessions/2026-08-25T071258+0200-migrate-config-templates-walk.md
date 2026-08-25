# Agent Session: Migrate Config Templates Walk

Status: NEXT
Session ID: `2026-08-25T071258+0200-migrate-config-templates-walk`
Created: `2026-08-25T07:12:58+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `d306fe4c05ba6de66c618f42e04d1032a7cb1da0c4c73a6eee2807e61b0d2e1b`
Previous: [2026-08-25T064347+0200-migrate-config-examples-read-dir.md](2026-08-25T064347+0200-migrate-config-examples-read-dir.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route
`GitCloudConfig.Templates`' direct tree walk through the existing filesystem
adapter. Preserve the exact receiver and templates-root construction, system
walk order, callback paths and metadata, incoming walk-error handling, project
configuration matching and loading, relative template names, partial ordered
results, returned errors, every completed filesystem, file, config, Bitbucket,
Wpost, and supervisor move, and every P2A contract with zero comparable ratchet
regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`filepath.Walk(root, callback)` operation in `GitCloudConfig.Templates` in
`pkg/config/cloud.go` and its use of the existing zero-value-safe filesystem
`Walk` operation. It authorizes focused config recording contracts. It does
not authorize `Examples`, `GitHookFiles`, another cloud or config method,
`InitProjectFromDirectory`, another filesystem operation, an adapter or
complete-double extension, another production effect, a function-valued effect
dependency, a new adapter family or public API, mutation harnesses, or later
roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, the complete `pkg/config` implementation and tests,
all callers of `Templates`, the complete `CloudConfig` interface and its
implementations or doubles, the filesystem adapter and relevant recording
doubles, `.quality/inventory`, and the continuity and quality-lift designs.
Regenerate ignored reports outside the measured tree or remove them before a
clean audit.

Implementation commit `cb94f81` has 236 tests across 16 of 25 packages. Q0.6
has 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 9 of 25, Q1.2 is 0, Q1.3 is 37 violations of 54 production
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
implementation commit for the Templates walk boundary, config contracts, and
measured planning notes. Then perform the normal separate handoff-only commit.
Do not push, merge, publish, distribute, remove the worktree, stash inherited
changes, revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, public API/CLI behavior, all four host acceptance flows,
`CloudConfig`, the exported
`GitCloudConfig.Templates() (templates []CloudTemplate, err error)` method
signature, its config and command callers, `GitCloudConfig.Implementation`,
`DirConfig.Dir`, `file.Path`, `Examples`, `GitHookFiles`,
`InitProjectFromDirectory`, and every caller's observable behavior.

Use the existing complete filesystem dependency and existing
`Walk(root string, walkFn filepath.WalkFunc) error` operation with a safe zero
value; do not add or change an adapter operation or store a function-valued
effect dependency. Production must select `filesystem.System()` only for the
Templates composition. Preserve exact evaluation of
`root := file.Path("%s/templates", gitCfg.Implementation().Dir())` and pass
that exact root and the existing callback through the adapter.

Preserve every callback observation and decision. For an incoming non-nil walk
error, return nil without dereferencing `info`, adding a template, wrapping or
propagating that incoming error, or changing subsequent system-walk behavior.
For a nil incoming error, preserve the exact `info.Name()` comparisons with
`projectConfigFileName` and `legacyProjectConfigFileName`; ignore every other
entry without inspecting additional metadata. For a matched configuration,
preserve the exact path splits, exact relative template name, exact
`InitProjectFromDirectory(relPath[0])` call, exact logged and returned project
load error, and exact `CloudTemplate{Name: name[1], Project: project}` append.
Preserve delivered order, system lexical walk order, partial results before an
error, and the exact final walk error. Do not normalize, wrap, retry, preflight,
resort, deduplicate, clean paths, change matching, add bounds checks, return full
paths, alter panic behavior, inspect other metadata, or broaden into another
config method, project loading, file behavior, Bitbucket, HTTP/process effects,
Wpost, clock, server, P4 adapters, P5 harnesses, dependencies, Docker, cloud
distribution, or publication.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete `pkg/config`, all
callers of `Templates`, the `CloudConfig` interface and every implementation
or double, and `internal/adapter/filesystem`. Read the relevant tests and the
Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation before editing. Before
the full gate, read the complete launcher contract, relevant Make meta-tests,
P2A API/CLI/subprocess contracts, and all four host acceptance flows.

# Three Moves

1. Start red with focused Templates recording contracts. Prove the complete
   dependency, production selection of `filesystem.System()`, exact
   receiver-derived templates root, non-empty recorded walk-root and callback
   populations, delivered callback order, paths, metadata, and incoming errors,
   exact suppression of incoming walk errors, exact current and legacy project
   filename matching, ignored entries, exact relative names and loaded projects,
   exact project-load and walk errors, partial ordered results, and safe zero
   behavior without developer-path access. Production composition contracts
   must not mutate the real filesystem.

2. Reuse only the existing filesystem adapter `Walk` operation. Keep exported
   `Templates` as the production entry, select `filesystem.System()` only for
   its private complete composition, and replace only its direct
   `filepath.Walk(root, callback)` with the adapter operation. Do not change
   the root expression or callback body, extend the adapter or complete doubles,
   change `CloudConfig`, touch another config method or caller, move project
   loading, or change inventory or public API.

3. Run focused config/filesystem contracts and command caller packages'
   relevant tests, the launcher contract from `/bin/bash`, Make preflight
   meta-contracts, API/CLI compatibility, full Go tests and race/vet, all four
   host acceptance flows, the audit meta-suite, focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean checkpoint
   audit, and empty-HOME count-2. Expect Q1.3 to move nominally from 37 of 54
   to 36 of 53 while Q0.6, Q1.2, Q1.4, and exact Q2.1 hold. Regenerate exact
   values; the full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Templates walk
move or record an exact resumable state. Rewrite the rolling handover, record
the measured P3 result, answer this archive, create one linked NEXT archive for
the next coherent P3 effect move, replace only the launcher's mutable regions,
run the launcher contract, and make the separate handoff-only commit
`docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
