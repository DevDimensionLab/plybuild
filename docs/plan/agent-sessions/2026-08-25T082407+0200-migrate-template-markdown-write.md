# Agent Session: Migrate Template Markdown Write

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T082407+0200-migrate-template-markdown-write`
Created: `2026-08-25T08:24:07+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `56749b2d3a17e145c6a9b766a661e60e7f0735faa249e93f1fb8f76a9e674cb0`
Previous: [2026-08-25T075301+0200-migrate-template-filtered-walk.md](2026-08-25T075301+0200-migrate-template-filtered-walk.md)
Next: [2026-08-25T085931+0200-migrate-config-project-write.md](2026-08-25T085931+0200-migrate-config-project-write.md)
Outcome: Routed only the template markdown write through the existing filesystem adapter with exact path, bytes, mode, return behavior, and zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route
`template.SaveTemplateListMarkdown`'s direct file write through the existing
filesystem adapter. Preserve the exported signature, exact receiver-derived
README path, document bytes, mode, returned path and write error, every
completed filesystem, file, template, config, Bitbucket, Wpost, and supervisor
move, and every P2A contract with zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`os.WriteFile(readmePath, []byte(markdownDocument), 0644)` operation in
`SaveTemplateListMarkdown` in `pkg/template/template.go` and its use of the
existing zero-value-safe filesystem `WriteFile` operation. It authorizes
focused template recording contracts. It does not authorize
`ListAsMarkdown`, `createTemplateListRenderingModel`,
`filteredFilesFromTemplate`, `MergeTemplate` or `merge` composition or
sequencing, `getIgnores`, copy, search, render, Maven, another template,
config, or file operation, an adapter or complete-double extension, another
production effect, a function-valued effect dependency, a new adapter family
or public API, mutation harnesses, or later roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, the complete `pkg/template` implementation and
tests, all callers of `SaveTemplateListMarkdown`, `ListAsMarkdown`,
`filteredFilesFromTemplate`, `merge`, and `MergeTemplate`, the relevant
config/file/Maven behavior, the filesystem adapter and relevant recording
doubles, `.quality/inventory`, and the continuity and quality-lift designs.
Regenerate ignored reports outside the measured tree or remove them before a
clean audit.

Implementation commit `170b0ae` has 250 tests across 16 of 25 packages. Q0.6
has 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 9 of 25, Q1.2 is 0, Q1.3 is 35 violations of 52 production
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
implementation commit for the template markdown write boundary, template
contracts, and measured planning notes. Then perform the normal separate
handoff-only commit. Do not push, merge, publish, distribute, remove the
worktree, stash inherited changes, revert user work, or run destructive Git
commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, public API/CLI behavior, all four host acceptance flows,
`CloudConfig`, `MergeTemplates`, `MergeTemplate`, `merge`,
`templateCopyDependencies`, `templateFiles`, `packageTemplateFiles`,
`filteredFilesFromTemplate`, `getIgnores`, and every caller's observable
behavior.

Use the existing complete filesystem dependency and existing
`WriteFile(name string, data []byte, perm fs.FileMode) error` operation with a
safe zero value; do not add or change an adapter operation, extend
`templateCopyDependencies`, or store a function-valued effect dependency.
Production must select `filesystem.System()` only for a private complete
markdown-write composition.

Preserve exact evaluation of
`readmePath := gitCfg.Implementation().Dir() + "/" + TemplatesDir +
"/README.md"`. Pass that exact path, `[]byte(markdownDocument)`, and `0644`
through the adapter. Preserve the exact `return readmePath, err`, including the
complete path when the write fails and the exact dependency error without
wrapping, logging, suppression, or substitution. Preserve empty strings,
arbitrary bytes represented by the Go string, nil and non-nil errors, and the
existing single write attempt. Do not clean or join paths, normalize
separators, use `file.Path` or `filepath.Join`, create directories, preflight,
retry, inspect or alter permissions, change `TemplatesDir`, reorder receiver
evaluation, add validation or logging, return an empty path on error, or
broaden into list rendering, filtered walking, merge, copy/search/render,
config, Maven, Bitbucket, HTTP/process effects, Wpost, clock, server, P4
adapters, P5 harnesses, dependencies, Docker, cloud distribution, or
publication.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete `pkg/template`, all
callers of `SaveTemplateListMarkdown`, `ListAsMarkdown`,
`filteredFilesFromTemplate`, `merge`, and `MergeTemplate`, the relevant
config/file/Maven code and tests, and `internal/adapter/filesystem`. Read the
relevant recording tests and the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit
implementation before editing. Before the full gate, read the complete
launcher contract, relevant Make meta-tests, P2A API/CLI/subprocess contracts,
and all four host acceptance flows.

# Three Moves

1. Start red with focused markdown-write recording contracts. Prove the
   complete dependency, production selection of `filesystem.System()`, exact
   receiver-derived README path, non-empty recorded write populations, exact
   document bytes including empty and non-ASCII content, exact `0644` mode,
   one write attempt, exact returned path on success and failure, exact write
   error identity, and safe zero behavior without developer-path access.
   Production composition contracts must not mutate the real filesystem.

2. Reuse only the existing filesystem adapter `WriteFile` operation. Keep the
   exported `SaveTemplateListMarkdown(gitCfg config.CloudConfig,
   markdownDocument string) (string, error)` as the production entry, select
   `filesystem.System()` only for its private complete composition, and
   replace only its direct `os.WriteFile` with the adapter operation. Do not
   change the path expression, byte conversion, mode, return statement,
   adapter, templateCopyDependencies, complete doubles, another template or
   config/file method or caller, inventory, or public API.

3. Run focused template/filesystem contracts and relevant config, file, Maven,
   and command caller packages' tests, the launcher contract from `/bin/bash`,
   Make preflight meta-contracts, API/CLI compatibility, full Go tests and
   race/vet, all four host acceptance flows, the audit meta-suite, focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean checkpoint
   audit, and empty-HOME count-2. Expect Q1.3 to move nominally from 35 of 52
   to 34 of 51 while Q0.6, Q1.2, Q1.4, and exact Q2.1 hold. Regenerate exact
   values; the full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent template
markdown write move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create one linked
NEXT archive for the next coherent P3 effect move, replace only the launcher's
mutable regions, run the launcher contract, and make the separate handoff-only
commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
