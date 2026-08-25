# Quality Upgrade Handover

Generated: 2026-08-25T10:09:11+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 2566438af9d5.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d,
d6cb593, 82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, and 2566438 in
roadmap order. The separate operational continuity implementation is 1b85711
and changes no Go quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T100637+0200-migrate-structurizr-output-write.md.
The tips-show-read predecessor is answered history, and the reciprocal archive
graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 34 Preserved

All tips command objects and registrations remain unchanged.
`tipsShowCmd.RunE` still warns and delegates identically for zero arguments.
For a non-empty argument list, it retains exact `name := args[0]`, then exact
`file.Path("%s/%s.md", tips.LocalDir(ctx.CloudConfig), name)`, and only its
private complete tips-show-read composition selects `filesystem.System()`.
The exact caller-composed path passes through the existing
`filesystem.ReadFile` operation. No adapter operation, established complete
recording double, caller, or public contract changed.

A read failure retains the exact
`fmt.Errorf("failed to find any tips file for [%s]: %s: %w", name, tipsPath,
err)` text and wrapped identity. Success retains the exact source bytes, then
terminal-config lookup, markdown rendering, stdout, local-source log,
global-cloud-config lookup, and cloud-source log in their established order.
No path cleaning, joining, normalization, validation, preflight, retry, extra
read, global mutation, or early logging was added.

Four new top-level command contracts bring the suite to 262 tests. They prove
complete production selection and delivery, exact arbitrary caller-composed
paths, non-empty recorded read populations, exact empty, representative,
non-ASCII, arbitrary, and partial returned bytes, one read, exact dependency
error identity at the private boundary, safe defaults without developer-path
access, and no unrelated filesystem operation. The production composition
tests neither read the real filesystem nor run the command.

## Measured Quality State

The clean full audit from implementation commit 2566438 reports:

- Absolute L0: 8 of 8.
- 262 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 32 direct external sites outside five declared adapters of 49
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 temporary-state claims across 48 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The clean report is
/private/tmp/ply-tips-clean-audit.zRAWpK/scorecard.json and the focused report
is /private/tmp/ply-tips-focused.v3etp1/scorecard.json. All authoritative audit
output and tool caches stayed outside the measured tree.

## Decisions And Learned Facts

1. The tips-show read reuses the adapter's established ReadFile operation; this
   move does not extend the adapter or any established complete filesystem
   recording double.
2. Q1.3 moved from 33/50 to 32/49 because the direct tips-show read disappeared
   while the already-counted adapter operation was reused.
3. The production RunE entry alone selects filesystem.System() for its private
   complete read composition; all command objects, registrations, and callers
   remain unchanged.
4. The zero-argument branch, exact name and path evaluation, wrapped read error,
   returned bytes, render/output/log/cloud order, and source-selection behavior
   are preserved without additional effects.
5. The focused recording double rejects every unrelated adapter operation and
   rejects an empty read population. Its safe zero case returns exact
   filesystem.ErrNoFilesystem without touching a developer path.
6. The next isolated flow is pkg/structurizr/structurizr.go's ignored
   os.WriteFile(outputFile, out.Bytes(), 0644). It can reuse established
   WriteFile without extending the adapter or an established complete double,
   nominally moving Q1.3 to 31/48.
7. Focused tests in the previously untested pkg/structurizr package should move
   Q1.1 from 9/25 to 8/25 and packages with tests from 16/25 to 17/25.
8. The next move must preserve stdout/stderr buffer assignment, one command
   run, exact early process error with no write, exact output path/stdout bytes/
   mode on success, intentionally discarded write error, and nil return.
   Private production-composition contracts must not run a command or mutate
   the real filesystem.
9. .quality/inventory is baseline-checksum-bound. Do not relabel seams or claim
   P5 mutation coverage.
10. Supply APIDIFF and GOLANGCI_LINT as environment variables for preflight,
    and keep GOLANGCI_LINT_CACHE and GOCACHE in writable external directories.
    The pinned tools are under /private/tmp/ply-quality-tools while that
    temporary directory remains available.
11. The unchanged launcher signal-interruption fixture transiently failed in
    three isolated runs: initial standalone validation, the first full
    preflight, and the first handoff validation. Immediate standalone reruns
    passed all 62 controls, the full preflight rerun passed, and make-test and
    test-agent-start runs passed without source changes.
12. The authoritative full audit used .quality/baseline/scorecard.json plus
    .quality/baseline/manual-evidence.json and exited 1 with valid clean
    structured output and zero comparable regressions.

## Next Objective

Move only the direct ignored
os.WriteFile(outputFile, out.Bytes(), 0644) in
structurizr.RunWithOutputToFile through the existing filesystem WriteFile
operation.

Start red with focused structurizr recording contracts for the complete
dependency, production system selection, exact caller output path, non-empty
write populations, empty and non-ASCII output bytes, mode 0644, one attempt,
exact private-boundary write-error identity, safe zero behavior without
developer-path access, and no unrelated adapter operation. Production
composition tests must not mutate the real filesystem or run a command.

Preserve the exported functions and every caller. Production must select
filesystem.System() only for a private complete structurizr-output-write
composition. Reuse adapter WriteFile without extending the adapter or
established complete doubles. Do not change buffers, command.Run, process
errors, command construction, output path/bytes/mode, discarded write error,
return behavior, plugin diagrams, tips, config/template/file/Maven behavior,
inventory, public API, another production effect, or enter clock/server, P4,
P5, or later roadmap work.

Expected direction is nominally Q1.1 8/25 and Q1.3 31/48 with Q0.6, Q1.2,
Q1.4, and exact Q2.1 held. Regenerate the structured result and require zero
comparable ratchet regressions.

## Verification Notes

Completed from implementation commit 2566438af9d5:

- Red evidence: focused compilation failed only on the absent private
  tips-show-read dependency, system composition, and read composition.
- Focused command/filesystem and relevant tips, context, config, file,
  template, and Maven caller package tests: PASS.
- Standalone /bin/bash launcher contract: PASS all 62 controls on immediate
  reruns after the unchanged signal-interruption fixture's transient failures.
- Make preflight meta-contracts and full preflight: PASS, including all 15
  audit meta-controls and all 62 launcher controls; the unchanged signal
  fixture required the recorded rerun.
- API, CLI, and subprocess compatibility: PASS.
- make test, make test-install, and make test-agent-start: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME go test ./... -count=2: PASS.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable; Q1.3 is exactly 32/49.
- Clean full audit: exit 1, 15 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Implementation commit: 2566438af9d52505b54b92428913fb6b89f4fd38
  (quality: route tips show read through filesystem adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete structurizr code/tests and callers,
plugin-diagram command and registration tests, relevant process/filesystem and
command/file/template/config/Maven/tips behavior, both existing adapters and
relevant doubles, and the named audits before editing. Confirm branch, HEAD,
status, reciprocal links, and ./codex-dev-start.sh --check. Begin red for only
the structurizr output-write boundary, then finish with one implementation
commit and one separate handoff-only commit.

Stop before either command.Run operation, Run, plugin command construction or
registration, buffer/error/return changes, another structurizr or filesystem
operation, another adapter operation or family, Q1.4 expansion, P4, mutation,
Docker, cloud distribution, or publication. Stop on API/CLI change, comparable
ratchet regression, authoritative full-audit exit 2, or failure to isolate
implementation from generated and handoff-only state.
