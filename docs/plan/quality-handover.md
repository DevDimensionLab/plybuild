# Quality Upgrade Handover

Generated: 2026-08-25T10:33:19+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 527a8b9941b2.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d,
d6cb593, 82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438, and
527a8b9 in roadmap order. The separate operational continuity implementation
is 1b85711 and changes no Go quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T103319+0200-migrate-maven-graph-styles-write.md.
The structurizr-output-write predecessor is answered history, and the
reciprocal archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 35 Preserved

`structurizr.Run` is unchanged. `structurizr.RunWithOutputToFile` still creates
`var out bytes.Buffer`, then `var stderr bytes.Buffer`, assigns those exact
buffers to command stdout and stderr, invokes `command.Run()` once, and returns
the exact non-nil process error before any write. On success, it selects
`filesystem.System()` only for a private complete structurizr-output-write
composition and attempts the existing adapter `WriteFile` operation once with
the exact caller `outputFile`, exact captured `out.Bytes()`, and mode `0644`.
It still discards the write error and returns nil.

Four new top-level contracts bring the suite to 266 tests and add
`pkg/structurizr` to the tested population. They prove complete production
selection and delivery, exact arbitrary caller paths, non-empty recorded write
populations, empty, representative, non-ASCII, and arbitrary bytes, exact mode,
one attempt, exact private-boundary dependency-error identity, safe zero
behavior without developer-path access, and rejection of every unrelated
filesystem operation. Production-composition tests neither mutate the real
filesystem nor run a command. No adapter, established complete double, caller,
command construction or execution, buffer behavior, process error, public API,
plugin-diagram behavior, or other completed effect changed.

## Measured Quality State

The clean full audit from implementation commit 527a8b9 reports:

- Absolute L0: 8 of 8.
- 266 test functions, zero skipped; 17 of 25 packages have tests.
- Q0.6: 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 8 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 31 direct external sites outside five declared adapters of 48
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 temporary-state claims across 49 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The clean report is
/private/tmp/ply-structurizr-clean-audit.xuhFch/scorecard.json and the focused
report is /private/tmp/ply-structurizr-focused.Vpw7yt/scorecard.json. All
authoritative audit output and tool caches stayed outside the measured tree.

## Decisions And Learned Facts

1. The structurizr output write reuses the adapter's established WriteFile
   operation; this move does not extend the adapter or any established complete
   filesystem recording double.
2. Q1.3 moved from 32/49 to 31/48 because the direct structurizr write
   disappeared while the already-counted adapter operation was reused.
3. Focused contracts in the previously untested `pkg/structurizr` package moved
   Q1.1 from 9/25 to 8/25 and the tested population from 16/25 to 17/25.
4. The exact command buffer allocation, assignment, execution, early process
   error, caller path, stdout bytes, mode, discarded write error, and return
   behavior remain visible in the unchanged production sequence.
5. The private recording double rejects every unrelated adapter operation and
   an empty write population. Its safe zero case returns exact
   filesystem.ErrNoFilesystem without touching a developer path.
6. `tips.List` is not the next exact reuse: its public result is
   `[]os.DirEntry`, while the existing adapter ReadDir operation returns
   `[]fs.FileInfo`. Converting or extending that boundary would require a
   separately authorized compatibility decision.
7. The next isolated exact-match flow is `maven.WriteGraphStyles`' direct
   `os.WriteFile(stylesFile, jsonStyles, 0644)`. It can reuse established
   WriteFile without extending the adapter or an established complete double,
   nominally moving Q1.3 to 30/47 while Q1.1 holds at 8/25.
8. The next move must preserve `json.MarshalIndent(styles, "", "    ")` and
   its early error, exact `file.Path("%s/target/dependency-graph-styles.json",
   projectPath)`, the existing `file.Exists` selection, no write when the file
   exists, and exact bytes, mode, write error, and caller behavior when absent.
9. `.quality/inventory` is baseline-checksum-bound. Do not relabel seams or
   claim P5 mutation coverage.
10. Supply APIDIFF and GOLANGCI_LINT as environment variables for preflight,
    and keep GOLANGCI_LINT_CACHE and GOCACHE in writable external directories.
    The pinned tools are under /private/tmp/ply-quality-tools while that
    temporary directory remains available.
11. The authoritative full audit used .quality/baseline/scorecard.json plus
    .quality/baseline/manual-evidence.json and exited 1 with valid clean
    structured output and zero comparable regressions.

## Next Objective

Move only the direct `os.WriteFile(stylesFile, jsonStyles, 0644)` in
`maven.WriteGraphStyles` through the existing filesystem WriteFile operation.

Start red with focused Maven graph-style recording contracts for the complete
dependency, production system selection, exact caller-composed target path,
non-empty write populations, empty, representative, non-ASCII, and arbitrary
JSON bytes, exact mode `0644`, one attempt, exact write-error identity, safe
zero behavior without developer-path access, and no unrelated adapter
operation. Production-composition tests must not mutate the real filesystem or
run Maven, dot, or any command.

Preserve exported `WriteGraphStyles`, `GraphDefaultStyles`, `GraphArgs`,
`Graph`, `RunOn`, every caller, and public API/CLI behavior. Production must
select `filesystem.System()` only for a private complete graph-style-write
composition. Reuse adapter WriteFile without extending the adapter or
established complete doubles. Do not change JSON marshal indentation or early
error, target path composition, existing `file.Exists` check or its behavior,
no-write-when-present selection, exact bytes/mode/error, Graph argument/filter/
style selection, Maven and command sequencing, structurizr, tips, config,
template, file, inventory, or another production effect.

Expected direction is nominally Q1.3 30/47 with Q1.1 8/25 and Q0.6, Q1.2,
Q1.4, and exact Q2.1 held. Regenerate the structured result and require zero
comparable ratchet regressions.

## Verification Notes

Completed from implementation commit 527a8b9941b2:

- Red evidence: focused compilation failed only on the absent private
  structurizr-output dependency, system composition, and write composition.
- Focused structurizr/filesystem/process and relevant command, file, template,
  config, Maven, tips, and context package tests: PASS.
- Standalone /bin/bash launcher contract: PASS all 62 controls.
- Make preflight meta-contracts and full preflight: PASS, including all 15
  audit meta-controls and all 62 launcher controls.
- API, CLI, and subprocess compatibility: PASS.
- make test, make test-install, and make test-agent-start: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME go test ./... -count=2: PASS.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable; Q1.3 is exactly 31/48.
- Clean full audit: exit 1, 15 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Implementation commit: 527a8b9941b27136f5c6cecf36cddd73fedfb6fb
  (quality: route structurizr output write through filesystem).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete Maven graph implementation/tests and callers,
plugin-Maven command and registration tests, relevant filesystem/process,
command/file/template/config/tips/structurizr behavior, existing adapters and
relevant doubles, and the named audits before editing. Confirm branch, HEAD,
status, reciprocal links, and ./codex-dev-start.sh --check. Begin red for only
the Maven graph-style output-write boundary, then finish with one implementation
commit and one separate handoff-only commit.

Stop before the existing file-existence boundary, JSON marshal or path changes,
Graph filtering/style/argument logic, RunOn or command execution, another Maven
operation, another adapter operation or family, Q1.4 expansion, P4, mutation,
Docker, cloud distribution, or publication. Stop on API/CLI change, comparable
ratchet regression, authoritative full-audit exit 2, or failure to isolate
implementation from generated and handoff-only state.
