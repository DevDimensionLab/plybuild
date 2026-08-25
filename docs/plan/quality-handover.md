# Quality Upgrade Handover

Generated: 2026-08-25T13:18:20+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 6928a72e90e8.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d,
d6cb593, 82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438,
527a8b9, 886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, and 6928a72 in roadmap
order. The separate operational continuity implementation is 1b85711 and
changes no Go quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T131820+0200-migrate-local-config-directory-mkdir.md.
The local-config update-create predecessor is answered history, and the
reciprocal archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 41 Preserved

Exported `LocalConfigDir.UpdateLocalConfig`, every caller, and all later config,
context, and profile behavior remain unchanged. It still evaluates exact
`CheckOrCreateConfigDir`, `FilePath`, the caller-supplied configuration, YAML
marshaling, its early error, and the update log in the same order. Production
selects `filesystem.System()` only for a private complete update-create
composition and reuses established adapter `Create`, preserving the exact
`filesystem.File` and error for the exact arbitrary path.

Four new top-level config contracts bring the suite to 290 tests. They prove
complete production selection and delivery, exact arbitrary file paths,
non-empty create populations, one attempt, exact file and error identity for
nil, successful, arbitrary-error, and unusual combined results, safe zero
behavior without developer-path access, the dedicated file double's exact
Close result, and rejection of every unrelated filesystem operation.
Production-composition tests do not mutate the real filesystem or run another
command. The completed update-write composition still receives the exact path,
YAML bytes, and `0644`; a create error returns before it, a write error returns
without Close, and write success returns the exact final Close result. No
adapter, established complete double, caller, public API, CLI, or other
completed effect changed.

## Measured Quality State

The clean full audit from implementation commit 6928a72 reports:

- Absolute L0: 8 of 8.
- 290 test functions, zero skipped; 17 of 25 packages have tests.
- Q0.6: 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 8 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 23 direct external sites outside five declared adapters of 40
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 temporary-state claims across 55 Markdown files at the clean
  implementation checkpoint.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The clean report with manual evidence is
/private/tmp/ply-local-config-update-create-clean-audit-manual.CsghaX/report/scorecard.json
and the focused report is
/private/tmp/ply-local-config-update-create-focused.BfOxnH/report/scorecard.json.
All authoritative audit output and tool caches stayed outside the measured
tree.

## Decisions And Learned Facts

1. The local-config update create reuses the adapter's established `Create`
   operation; this move does not extend the adapter or any established complete
   filesystem recording double.
2. Q1.3 moved from 25/42 to 23/40 because the direct create and the scanner's
   concrete `*os.File.Close` classification disappeared while the lifecycle
   call remained. Q1.1 held at 8/25.
3. Exact directory and path selection, caller-supplied config, YAML bytes,
   logging, completed write, write error, Close selection, and Close result
   remain visible around the one replaced call.
4. The private recording double rejects every unrelated adapter operation and
   an empty create population. Its safe zero case returns exact
   `filesystem.ErrNoFilesystem` without touching a developer path; its
   dedicated file double preserves exact Close results.
5. The remaining local-config effect is the direct `os.Mkdir(dir, 0755)` in
   `CheckOrCreateConfigDir`. Established `MkdirAll` is not equivalent because
   it creates missing parents and succeeds for an existing directory.
6. The next move therefore needs a distinct internal filesystem adapter
   `Mkdir(path string, mode fs.FileMode) error` operation that maps exactly to
   `os.Mkdir`, has a safe zero value, and does not alter `MkdirAll`.
7. Extending `filesystem.FileSystem` requires the system implementation and all
   current complete filesystem test doubles to add the method. The existing 34
   test doubles should reject or record it consistently without gaining any
   unrelated behavior.
8. The config production entry must still make the established adapter stat
   first and apply exact `os.IsNotExist(err)`. Only that branch may select the
   private complete directory-create composition and attempt one exact mkdir.
9. Replacing the one direct mkdir should move Q1.3 nominally from 23/40 to
   22/39. Regenerate the scanner result instead of assuming it.
10. `tips.List` is still not an exact reuse: its public result is
    `[]os.DirEntry`, while established `ReadDir` returns `[]fs.FileInfo`.
11. `.quality/inventory` is baseline-checksum-bound. Do not relabel seams or
    claim P5 mutation coverage.
12. Supply APIDIFF and GOLANGCI_LINT as environment variables for preflight,
    and keep GOLANGCI_LINT_CACHE and GOCACHE in writable external directories.
    The pinned tools remain under /private/tmp/ply-quality-tools while present.
13. The authoritative full audit used `.quality/baseline/scorecard.json` plus
    `.quality/baseline/manual-evidence.json` and exited 1 with valid clean
    structured output and zero comparable regressions.

## Next Objective

Move only the direct `os.Mkdir(dir, 0755)` in
`LocalConfigDir.CheckOrCreateConfigDir` through one new exact filesystem
adapter `Mkdir` operation.

Start red at the adapter boundary for safe zero behavior, exact arbitrary path,
exact mode, one call, exact error identity, and production `os.Mkdir` semantics:
missing parents fail, an existing directory fails, and no recursive creation
occurs. Then add focused local-config-directory-create recording contracts for
the complete dependency, production system selection, exact path and `0755`,
non-empty populations, one attempt, exact error, and no unrelated adapter
operation. Production-composition contracts must not mutate the real
filesystem or run another command.

Preserve exported config types and methods, every caller, exact directory
selection, completed stat dependency and result, exact `os.IsNotExist`
branching, early errors, and all later TouchFile and UpdateLocalConfig
lifecycle. Add only the internal adapter operation, its system mapping and
contracts, and the mandatory complete-double method implementations. Do not
change `MkdirAll`, stat, file-path selection, YAML, logging, create/write/Close,
another config operation, inventory, public API/CLI, or another production
effect.

Expected direction is nominally Q1.3 22/39 with Q1.1 8/25 and Q0.6, Q1.2,
Q1.4, and exact Q2.1 held. Regenerate the structured result and require zero
comparable ratchet regressions.

## Verification Notes

Completed from implementation commit 6928a72e90e8:

- Red evidence: with an external Go build cache, focused compilation failed
  only on the absent private local-config-update-create dependency, system
  composition, and create composition.
- Focused config/filesystem/process and relevant command, context, file,
  template, Maven, tips, structurizr, Bitbucket, HTTP, Kibana, Spring, and shell
  package tests: PASS.
- Standalone `/bin/bash` launcher contract: PASS all 62 controls on the complete
  rerun. Its first run hit the documented nested signal-interruption partial-
  raw-log timing flake after 25 controls.
- Make preflight meta-contract and full preflight: PASS, including all 15 audit
  meta-controls and all 62 launcher controls.
- API, CLI, and subprocess compatibility: PASS.
- Full make test and install targets: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS from the clean commit.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable; Q1.3 is exactly 23/40.
- Clean full audit: exit 1, 15 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Handoff-state Q3.4 audit: PASS with zero state-claim phrases across 56
  Markdown files and zero ratchet regressions.
- The first handoff-state launcher run hit the same documented nested signal
  partial-raw-log timing flake; its immediate complete rerun passed all 62
  controls.
- Implementation commit: 6928a72e90e8e5be7117e10c22988ae59d566012
  (quality: route local config update create through filesystem).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete local config implementation/tests and callers,
context/profile command tests, relevant filesystem/file/template/config/Maven/
tips/structurizr behavior, the complete filesystem and process adapters, every
complete filesystem recording double, and the named audits before editing.
Confirm branch, HEAD, status, reciprocal links, and
`./codex-dev-start.sh --check`. Begin red for only the filesystem adapter
`Mkdir` contract and local-config directory-create boundary, then finish with
one implementation commit and one separate handoff-only commit.

Stop before changing the completed directory stat or `os.IsNotExist` branch,
`MkdirAll`, TouchFile, UpdateLocalConfig, another local-config operation,
another adapter operation or family, Q1.4 expansion, P4, mutation, Docker,
cloud distribution, or publication. Stop on API/CLI change, comparable ratchet
regression, authoritative full-audit exit 2, or failure to isolate
implementation from generated and handoff-only state.
