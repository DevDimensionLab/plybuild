# Agent Session: Migrate Project Config Write

Status: NEXT
Session ID: `2026-08-25T085931+0200-migrate-config-project-write`
Created: `2026-08-25T08:59:31+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `72f9a1400e9675e5a922abcbaca6d12d44b37ba937bdde51ac7557734c94292e`
Previous: [2026-08-25T082407+0200-migrate-template-markdown-write.md](2026-08-25T082407+0200-migrate-template-markdown-write.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route
`(*config.ProjectConfiguration).WriteTo`'s direct file write through the
existing filesystem adapter. Preserve the exported method and `ProjectConfig`
contract, exact log and JSON-marshal sequencing, target path, serialized bytes,
mode and returned error, every completed filesystem, file, template, config,
Maven, Bitbucket, Wpost, and supervisor move, and every P2A contract with zero
comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`os.WriteFile(targetFile, data, 0644)` operation in
`(*ProjectConfiguration).WriteTo` in `pkg/config/project.go` and its use of the
existing zero-value-safe filesystem `WriteFile` operation. It authorizes
focused config recording contracts. It does not authorize changing logging,
JSON composition, marshal error handling, `ProjectConfig`, another project or
config method, `SortAndWritePom`, project initialization, cloud config,
template, file, Maven, command composition or sequencing, an adapter or
complete-double extension, another production effect, a function-valued effect
dependency, a new adapter family or public API, mutation harnesses, or later
roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, the complete `pkg/config` implementation and tests,
all callers and interface uses of `(*ProjectConfiguration).WriteTo`,
`ProjectConfig`, `ProjectConfigPath`, `projectConfigFile`, and
`SortAndWritePom`, the relevant command/template/file/Maven behavior, the
filesystem adapter and relevant recording doubles, `.quality/inventory`, and
the continuity and quality-lift designs. Regenerate ignored reports outside
the measured tree or remove them before a clean audit.

Implementation commit `8d49345` has 254 tests across 16 of 25 packages. Q0.6
has 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 9 of 25, Q1.2 is 0, Q1.3 is 34 violations of 51 production
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
implementation commit for the project-config write boundary, config contracts,
and measured planning notes. Then perform the normal separate handoff-only
commit. Do not push, merge, publish, distribute, remove the worktree, stash
inherited changes, revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, public API/CLI behavior, all four host acceptance flows,
`ProjectConfig`, `ProjectConfiguration`, `ProjectConfigPath`,
`projectConfigFile`, `SortAndWritePom`, `CloudConfig`, `MergeTemplates`,
template composition, and every caller's observable behavior.

Use the existing complete filesystem dependency and existing
`WriteFile(name string, data []byte, perm fs.FileMode) error` operation with a
safe zero value; do not add or change an adapter operation, extend another
complete recording double, or store a function-valued effect dependency.
Production must select `filesystem.System()` only for a private complete
project-config-write composition.

Preserve the exact first log call
`log.Infof("writes project config file to %s", targetFile)`, then the exact
`json.MarshalIndent(config, "", "    ")` evaluation and early return of its
error, then one adapter write attempt with the exact caller-supplied
`targetFile`, returned `data`, and `0644`. Preserve nil and non-nil receivers,
empty and non-ASCII fields, nil and empty collections, map-key serialization,
arbitrary target strings, the absence of directory creation or path
normalization, and the exact dependency error without wrapping, logging,
suppression, or substitution. Do not clean or join paths, use `file.Path` or
`filepath.Join`, preflight, retry, inspect or alter permissions, reorder
logging or marshaling, mutate the receiver, or broaden into another config,
template, file, Maven, Bitbucket, HTTP/process, Wpost, clock, server, P4
adapter, P5 harness, dependency, Docker, distribution, or publication move.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete `pkg/config`, all
callers and interface uses named above, the relevant command/template/file and
Maven code and tests, and `internal/adapter/filesystem`. Read the relevant
recording tests and the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit
implementation before editing. Before the full gate, read the complete
launcher contract, relevant Make meta-tests, P2A API/CLI/subprocess contracts,
and all four host acceptance flows.

# Three Moves

1. Start red with focused project-config-write recording contracts. Prove the
   complete dependency, production selection of `filesystem.System()`, exact
   caller-supplied target path, non-empty recorded write populations, exact
   JSON bytes for nil and representative non-nil receivers including empty and
   non-ASCII content, exact `0644` mode, one write attempt, exact write error
   identity, unchanged receiver state, and safe zero behavior without
   developer-path access. Production composition contracts must not mutate the
   real filesystem.

2. Reuse only the existing filesystem adapter `WriteFile` operation. Keep the
   exported `(*ProjectConfiguration).WriteTo(targetFile string) error` as the
   production entry, select `filesystem.System()` only for its private complete
   composition, and replace only its direct `os.WriteFile` with the adapter
   operation. Do not change the log, marshal expression or early return,
   target argument, bytes, mode, return semantics, adapter, another complete
   double, config method or caller, inventory, or public API.

3. Run focused config/filesystem contracts and relevant command, template,
   file, and Maven caller packages' tests, the launcher contract from
   `/bin/bash`, Make preflight meta-contracts, API/CLI compatibility, full Go
   tests and race/vet, all four host acceptance flows, the audit meta-suite,
   focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean
   checkpoint audit, and empty-HOME count-2. Expect Q1.3 to move nominally from
   34 of 51 to 33 of 50 while Q0.6, Q1.2, Q1.4, and exact Q2.1 hold. Regenerate
   exact values; the full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent project-config
write move or record an exact resumable state. Rewrite the rolling handover,
record the measured P3 result, answer this archive, create one linked NEXT
archive for the next coherent P3 effect move, replace only the launcher's
mutable regions, run the launcher contract, and make the separate handoff-only
commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
