# Quality Upgrade Handover

Generated: 2026-08-25T11:31:36+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 89f43aa44096.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d,
d6cb593, 82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438,
527a8b9, 886dff0, and 89f43aa in roadmap order. The separate operational
continuity implementation is 1b85711 and changes no Go quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T112922+0200-migrate-local-config-update-write.md.
The local-config touch predecessor is answered history, and the reciprocal
archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 37 Preserved

Exported `LocalConfigDir.TouchFile`, every caller, and all later config,
context, and profile behavior remain unchanged. It still evaluates exact
`CheckOrCreateConfigDir()` and its early error, `FilePath()`, zero-value
`LocalConfiguration{}`, the exact default cloud URL assignment,
`yaml.Marshal(&config)` and its early error, and the exact creation log in
the same order. It preserves the following `os.Create(configFilePath)`
attempt and its early error. Production then selects `filesystem.System()`
only for a private complete touch-write composition and attempts established
adapter `WriteFile` once with the exact path, YAML bytes, and mode `0644`.
The exact write error returns without a Close attempt; successful writes still
return the exact `f.Close()` result.

Four new top-level config contracts bring the suite to 274 tests. They prove
complete production selection and delivery, exact arbitrary caller-composed
paths, non-empty recorded write populations, empty, representative, non-ASCII,
and arbitrary bytes, exact mode, one attempt, exact private-boundary error
identity, safe zero behavior without developer-path access, and rejection of
every unrelated filesystem operation. Production-composition tests do not
mutate the real filesystem or run another command. No adapter, established
complete double, directory creation, path construction, config initialization,
default, YAML semantics, log, `os.Create`, Close, `UpdateLocalConfig`,
caller, public API, CLI, or other completed effect changed.

## Measured Quality State

The clean full audit from implementation commit 89f43aa reports:

- Absolute L0: 8 of 8.
- 274 test functions, zero skipped; 17 of 25 packages have tests.
- Q0.6: 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 8 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 29 direct external sites outside five declared adapters of 46
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 temporary-state claims across 51 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The clean report is
/private/tmp/ply-local-config-clean-audit.FdNW96/scorecard.json and the focused
report is /private/tmp/ply-local-config-focused.aQrQi3/scorecard.json. All
authoritative audit output and tool caches stayed outside the measured tree.

## Decisions And Learned Facts

1. The local-config touch output write reuses the adapter's established
   WriteFile operation; this move does not extend the adapter or any established
   complete filesystem recording double.
2. Q1.3 moved from 30/47 to 29/46 because the direct TouchFile write disappeared
   while the already-counted adapter operation was reused. Q1.1 held at 8/25.
3. The exact directory/path, config/default, marshal, logging, file creation,
   write error, and final Close ordering remain visible in the unchanged
   production sequence around the one replaced call.
4. The private recording double rejects every unrelated adapter operation and
   an empty write population. Its safe zero case returns exact
   filesystem.ErrNoFilesystem without touching a developer path.
5. The next isolated exact-match flow is
   `LocalConfigDir.UpdateLocalConfig`'s direct
   `os.WriteFile(configFilePath, d, 0644)`. It can reuse established WriteFile
   without extending the adapter or an established complete double, nominally
   moving Q1.3 to 28/45 while Q1.1 holds at 8/25.
6. That move must preserve CheckOrCreateConfigDir, FilePath, the exact
   caller-supplied config, yaml.Marshal evaluation and early error, update log
   ordering, preceding os.Create selection/error, exact bytes/mode/error, and
   the existing final file Close behavior.
7. `tips.List` is not an exact reuse: its public result is `[]os.DirEntry`,
   while the established adapter ReadDir returns `[]fs.FileInfo`. Converting
   or extending that boundary requires a separate compatibility decision.
8. `.quality/inventory` is baseline-checksum-bound. Do not relabel seams or
   claim P5 mutation coverage.
9. Supply APIDIFF and GOLANGCI_LINT as environment variables for preflight,
   and keep GOLANGCI_LINT_CACHE and GOCACHE in writable external directories.
   The pinned tools are under /private/tmp/ply-quality-tools while that
   temporary directory remains available.
10. The authoritative full audit used .quality/baseline/scorecard.json plus
    .quality/baseline/manual-evidence.json and exited 1 with valid clean
    structured output and zero comparable regressions.

## Next Objective

Move only the direct `os.WriteFile(configFilePath, d, 0644)` in
`LocalConfigDir.UpdateLocalConfig` through the existing filesystem WriteFile
operation.

Start red with focused local-config-update recording contracts for the complete
dependency, production system selection, exact caller-composed target path,
non-empty write populations, empty, representative, non-ASCII, and arbitrary
YAML bytes, exact mode `0644`, one attempt, exact write-error identity, safe
zero behavior without developer-path access, and no unrelated adapter
operation. Production-composition tests must not mutate the real filesystem.

Preserve exported config types and methods, every caller, the exact
caller-supplied config, and public API/CLI behavior. Production must select
`filesystem.System()` only for a private complete local-config-update-write
composition. Reuse adapter WriteFile without extending the adapter or
established complete doubles. Do not change directory selection or creation,
file path composition, YAML marshaling, logging, the preceding os.Create, the
final Close, exact bytes, mode or error, TouchFile, Config, Print, Exists,
context/profile flows, inventory, or another production effect.

Expected direction is nominally Q1.3 28/45 with Q1.1 8/25 and Q0.6, Q1.2,
Q1.4, and exact Q2.1 held. Regenerate the structured result and require zero
comparable ratchet regressions.

## Verification Notes

Completed from implementation commit 89f43aa44096:

- Red evidence: with an external Go build cache, focused compilation failed
  only on the absent private local-config-touch dependency, system composition,
  and write composition.
- Focused config/filesystem/process and relevant command, context, file,
  template, Maven, tips, structurizr, Bitbucket, HTTP, Kibana, Spring, and shell
  package tests: PASS.
- Standalone /bin/bash launcher contract: PASS all 62 controls. The first
  handoff-state run hit the documented nested signal-interruption partial-log
  timing flake; its immediate complete rerun passed all 62 controls.
- Make preflight meta-contract and full preflight: PASS, including all 15 audit
  meta-controls and all 62 launcher controls.
- API, CLI, and subprocess compatibility: PASS.
- make test, including the install and launcher contracts: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME go test ./... -count=2: PASS.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable; Q1.3 is exactly 29/46.
- Clean full audit: exit 1, 15 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Implementation commit: 89f43aa440965caf25991e33fa748baf1e47adb8
  (quality: route local config touch write through filesystem).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete local config implementation/tests and callers,
context/profile command tests, relevant filesystem/file/template/config/Maven/
tips/structurizr behavior, existing adapters and relevant doubles, and the
named audits before editing. Confirm branch, HEAD, status, reciprocal links,
and ./codex-dev-start.sh --check. Begin red for only the local-config update
output-write boundary, then finish with one implementation commit and one
separate handoff-only commit.

Stop before CheckOrCreateConfigDir, yaml.Marshal, logging, os.Create, Close,
TouchFile, another local-config operation, another adapter operation or family,
Q1.4 expansion, P4, mutation, Docker, cloud distribution, or publication. Stop
on API/CLI change, comparable ratchet regression, authoritative full-audit exit
2, or failure to isolate implementation from generated and handoff-only state.
