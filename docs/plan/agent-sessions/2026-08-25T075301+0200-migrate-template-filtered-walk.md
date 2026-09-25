# Agent Session: Migrate Template Filtered Walk

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T075301+0200-migrate-template-filtered-walk`
Created: `2026-08-25T07:53:01+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `d22416d4473eac8221bc517d999fe963894a2b6f4274a35f4e8eeec24562684a`
Previous: [2026-08-25T071258+0200-migrate-config-templates-walk.md](2026-08-25T071258+0200-migrate-config-templates-walk.md)
Next: [2026-08-25T082407+0200-migrate-template-markdown-write.md](2026-08-25T082407+0200-migrate-template-markdown-write.md)
Outcome: Routed only the filtered template walk through the existing filesystem adapter with exact callback behavior and zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route
`template.filteredFilesFromTemplate`'s direct tree walk through the existing
filesystem adapter. Preserve the exact source root, system walk order,
callback paths, metadata and incoming-error behavior, directory handling,
root-file detection, ordered filtering and render exception, logging, ordered
and partial results, returned errors, every completed filesystem, file,
config, Bitbucket, Wpost, and supervisor move, and every P2A contract with zero
comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`filepath.Walk(sourceDir, callback)` operation in
`template.filteredFilesFromTemplate` in `pkg/template/template.go` and its use
of the existing zero-value-safe filesystem `Walk` operation. It authorizes
focused template recording contracts. It does not authorize
`SaveTemplateListMarkdown`, `MergeTemplate` or `merge` composition or
sequencing, `getIgnores`, copy, search, render, Maven, another template,
config, or file operation, an adapter or complete-double extension, another
production effect, a function-valued effect dependency, a new adapter family
or public API, mutation harnesses, or later roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, the complete `pkg/template` implementation and
tests, all callers of `filteredFilesFromTemplate`, `merge`, and
`MergeTemplate`, the relevant config/file/Maven behavior, the filesystem
adapter and relevant recording doubles, `.quality/inventory`, and the
continuity and quality-lift designs. Regenerate ignored reports outside the
measured tree or remove them before a clean audit.

Implementation commit `dfcfa75` has 243 tests across 16 of 25 packages. Q0.6
has 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 9 of 25, Q1.2 is 0, Q1.3 is 36 violations of 53 production
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
implementation commit for the filtered template walk boundary, template
contracts, and measured planning notes. Then perform the normal separate
handoff-only commit. Do not push, merge, publish, distribute, remove the
worktree, stash inherited changes, revert user work, or run destructive Git
commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, public API/CLI behavior, all four host acceptance flows,
`CloudConfig`, `MergeTemplates`, `MergeTemplate`, `merge`,
`templateCopyDependencies`, `templateFiles`, `packageTemplateFiles`, the
private `filteredFilesFromTemplate(sourceDir string, filter []string)
(files []string, err error)` signature and callers, `getIgnores`, and every
caller's observable behavior.

Use the existing complete filesystem dependency and existing
`Walk(root string, walkFn filepath.WalkFunc) error` operation with a safe zero
value; do not add or change an adapter operation, extend
`templateCopyDependencies`, or store a function-valued effect dependency.
Production must select `filesystem.System()` only for a private complete
filtered-walk composition. Preserve exact evaluation of the supplied
`sourceDir` and pass that exact string and the existing callback through the
adapter.

Preserve every callback observation and decision. The callback ignores its
incoming `err` argument and first evaluates `info.IsDir()`. A directory returns
nil without calling `info.Name()`, evaluating filters, logging, or appending a
path. A nil `info`, including the ordinary system-walk error shape, therefore
retains the exact current panic and must not become a returned, wrapped, or
suppressed error. A non-nil `info` with a non-nil incoming error is processed
exactly like any other non-directory entry; do not branch on or propagate the
incoming error.

For a non-directory, preserve exact evaluation of
`rootDir := sourceDir == strings.ReplaceAll(path, "/"+info.Name(), "")`.
Iterate `filter` in its delivered order. Preserve the exact
`ignore == "pom.xml" && !rootDir` continue for only that filter item. For every
other item, preserve the exact `strings.Contains(path, ignore)` match and
global `strings.Contains(path, ".render")` exception, the exact
`log.Debugf("ignoring %s in %s", path, filter)` call, and nil callback return.
Otherwise preserve exact `files = append(files, path)` and
`log.Debugf("fileToCopy: %s", path)`. Preserve delivered order, system lexical
walk order, nil results, partial results before a final walk error, and the
exact final walk error. Do not normalize separators, clean paths, use Rel,
change ReplaceAll, precompile or reorder filters, deduplicate, sort, retry,
preflight, add bounds or nil checks, alter panic behavior, inspect other
metadata, return full or relative paths differently, or broaden into merge,
copy/search/render, config, Maven, Bitbucket, HTTP/process effects, Wpost,
clock, server, P4 adapters, P5 harnesses, dependencies, Docker, cloud
distribution, or publication.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete `pkg/template`, all
callers of `filteredFilesFromTemplate`, `merge`, and `MergeTemplate`, the
relevant config/file/Maven code and tests, and `internal/adapter/filesystem`.
Read the relevant recording tests and the
Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation before editing. Before
the full gate, read the complete launcher contract, relevant Make meta-tests,
P2A API/CLI/subprocess contracts, and all four host acceptance flows.

# Three Moves

1. Start red with focused filtered-walk recording contracts. Prove the complete
   dependency, production selection of `filesystem.System()`, exact sourceDir
   walk root, non-empty recorded walk-root and callback populations, delivered
   callback order, paths, metadata and incoming errors, exact directory
   handling and metadata observations, nil-info incoming-error panic behavior,
   non-nil-info incoming-error processing, slash-based root detection, ordered
   filter behavior, non-root pom.xml handling, substring matches, the .render
   exception, exact logging, ordered and nil results, exact final walk errors,
   partial results, and safe zero behavior without developer-path access.
   Production composition contracts must not mutate the real filesystem.

2. Reuse only the existing filesystem adapter `Walk` operation. Keep private
   `filteredFilesFromTemplate` as the production entry, select
   `filesystem.System()` only for its private complete composition, and replace
   only its direct `filepath.Walk(sourceDir, callback)` with the adapter
   operation. Do not change the sourceDir argument or callback body, extend the
   adapter, templateCopyDependencies, or complete doubles, touch merge or
   another template/config/file method or caller, or change inventory or public
   API.

3. Run focused template/filesystem contracts and relevant config, file, Maven,
   and command caller packages' tests, the launcher contract from `/bin/bash`,
   Make preflight meta-contracts, API/CLI compatibility, full Go tests and
   race/vet, all four host acceptance flows, the audit meta-suite, focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean checkpoint
   audit, and empty-HOME count-2. Expect Q1.3 to move nominally from 36 of 53
   to 35 of 52 while Q0.6, Q1.2, Q1.4, and exact Q2.1 hold. Regenerate exact
   values; the full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent filtered
template walk move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create one linked
NEXT archive for the next coherent P3 effect move, replace only the launcher's
mutable regions, run the launcher contract, and make the separate handoff-only
commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
