# Quality Upgrade Handover

Generated: 2026-08-25T12:22:24+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: d2e330b1bb1b.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d,
d6cb593, 82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438,
527a8b9, 886dff0, 89f43aa, fd45ebf, and d2e330b in roadmap order. The separate
operational continuity implementation is 1b85711 and changes no Go quality
denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T122224+0200-migrate-local-config-touch-create.md.
The local-config directory-stat predecessor is answered history, and the
reciprocal archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 39 Preserved

Exported `LocalConfigDir.CheckOrCreateConfigDir`, every caller, and all later
config, context, and profile behavior remain unchanged. It still evaluates the
exact `localCfgDir.Implementation().Path` result and attempts one stat with
that arbitrary string. Production selects `filesystem.System()` only for a
private complete directory-stat composition and reuses established adapter
`Stat`, preserving its exact `fs.FileInfo` and error. The enclosing code still
applies exact `os.IsNotExist(err)`: only a not-exist error attempts exact
`os.Mkdir(dir, 0755)` and returns its exact non-nil error; every other stat
result still returns nil.

Four new top-level config contracts bring the suite to 282 tests. They prove
complete production selection and delivery, exact arbitrary directory paths,
non-empty recorded stat populations, one attempt, exact metadata and error
identity for nil, existing, not-exist, and arbitrary error results, safe zero
behavior without developer-path access, and rejection of every unrelated
filesystem operation. Production-composition tests do not mutate the real
filesystem or run another command. No adapter, established complete double,
directory evaluation or creation, predicate, mode, error, local-config write,
caller, public API, CLI, or other completed effect changed.

## Measured Quality State

The clean full audit from implementation commit d2e330b reports:

- Absolute L0: 8 of 8.
- 282 test functions, zero skipped; 17 of 25 packages have tests.
- Q0.6: 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 8 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 27 direct external sites outside five declared adapters of 44
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 temporary-state claims across 53 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The clean report is
/private/tmp/ply-local-config-directory-stat-clean-audit/scorecard.json and the
focused report is
/private/tmp/ply-local-config-directory-stat-focused/scorecard.json. All
authoritative audit output and tool caches stayed outside the measured tree.

## Decisions And Learned Facts

1. The local-config directory stat reuses the adapter's established `Stat`
   operation; this move does not extend the adapter or any established complete
   filesystem recording double.
2. Q1.3 moved from 28/45 to 27/44 because the direct stat disappeared while the
   already-counted adapter operation was reused. Q1.1 held at 8/25.
3. Exact directory selection, stat result, not-exist branching, direct mkdir,
   mode, mkdir error, and every later local-config lifecycle step remain visible
   around the one replaced call.
4. The private recording double rejects every unrelated adapter operation and
   an empty stat population. Its safe zero case returns exact
   `filesystem.ErrNoFilesystem` without touching a developer path.
5. The next isolated exact-match flow is `LocalConfigDir.TouchFile`'s direct
   `os.Create(configFilePath)`. It can reuse established `Create` without
   extending the adapter or an established complete double.
6. Because established `filesystem.Create` returns the adapter-owned
   `filesystem.File` interface, replacing that direct create also removes the
   scanner's concrete `*os.File.Close` classification while retaining the same
   final `f.Close()` lifecycle call. The nominal Q1.3 result is therefore
   25/42; regenerate it rather than assuming it.
7. That move must preserve `CheckOrCreateConfigDir`, exact `FilePath`, config
   initialization/default, YAML marshal and error, logging, one create attempt
   and exact create result, the completed adapter write and its error, no Close
   after a write error, and exact final Close result after write success.
8. `os.Mkdir(dir, 0755)` is not an exact reuse of established `MkdirAll`; moving
   it requires a separate adapter-contract decision and is not part of the next
   create move.
9. `tips.List` is not an exact reuse: its public result is `[]os.DirEntry`,
   while established `ReadDir` returns `[]fs.FileInfo`.
10. `.quality/inventory` is baseline-checksum-bound. Do not relabel seams or
    claim P5 mutation coverage.
11. Supply APIDIFF and GOLANGCI_LINT as environment variables for preflight,
    and keep GOLANGCI_LINT_CACHE and GOCACHE in writable external directories.
    The pinned tools remain under /private/tmp/ply-quality-tools while present.
12. The authoritative full audit used `.quality/baseline/scorecard.json` plus
    `.quality/baseline/manual-evidence.json` and exited 1 with valid clean
    structured output and zero comparable regressions.

## Next Objective

Move only the direct `os.Create(configFilePath)` in `LocalConfigDir.TouchFile`
through the existing filesystem `Create` operation.

Start red with focused local-config-touch-create recording contracts for the
complete dependency, production system selection, exact arbitrary paths,
non-empty create populations, one attempt, exact `filesystem.File` and error
identity for nil, successful, arbitrary-error, and unusual combined results,
safe zero behavior without developer-path access, and no unrelated adapter
operation. Use a dedicated file double to preserve exact Close results without
mutating the real filesystem.

Preserve exported config types and methods, every caller, exact directory and
file path selection, default config and YAML bytes, logging, completed touch
write composition and exact error, final Close selection and exact result, and
public API/CLI behavior. Production must select `filesystem.System()` only for
a private complete local-config-touch-create composition. Reuse adapter Create
without extending the adapter or established complete doubles. Do not change
the directory stat/mkdir, `UpdateLocalConfig`, another local-config operation,
inventory, or another production effect.

Expected direction is nominally Q1.3 25/42 with Q1.1 8/25 and Q0.6, Q1.2,
Q1.4, and exact Q2.1 held. Regenerate the structured result and require zero
comparable ratchet regressions.

## Verification Notes

Completed from implementation commit d2e330b1bb1b:

- Red evidence: with an external Go build cache, focused compilation failed
  only on the absent private local-config-directory-stat dependency, system
  composition, and stat composition.
- Focused config/filesystem/process and relevant command, context, file,
  template, Maven, tips, structurizr, Bitbucket, HTTP, Kibana, Spring, and shell
  package tests: PASS.
- Standalone `/bin/bash` launcher contract: PASS all 62 controls on complete
  reruns; two separate runs hit the documented nested signal-interruption
  partial-raw-log timing flake before immediate complete passes.
- Make preflight meta-contract and full preflight: PASS, including all 15 audit
  meta-controls and all 62 launcher controls.
- API, CLI, and subprocess compatibility: PASS.
- make test, install, and launcher targets: PASS on completed runs.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS from the clean commit.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable; Q1.3 is exactly 27/44.
- Clean full audit: exit 1, 15 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Handoff-state Q3.4 audit: PASS with zero state-claim phrases across 54
  Markdown files and zero ratchet regressions.
- Implementation commit: d2e330b1bb1bb2df3e4cb77b8b9b001a4c77d7a8
  (quality: route local config directory stat through filesystem).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete local config implementation/tests and callers,
context/profile command tests, relevant filesystem/file/template/config/Maven/
tips/structurizr behavior, existing adapters and relevant doubles, and the
named audits before editing. Confirm branch, HEAD, status, reciprocal links,
and `./codex-dev-start.sh --check`. Begin red for only the local-config touch
create boundary, then finish with one implementation commit and one separate
handoff-only commit.

Stop before the directory mkdir, the completed touch write, Close semantics,
UpdateLocalConfig create or write, another local-config operation, another
adapter operation or family, Q1.4 expansion, P4, mutation, Docker, cloud
distribution, or publication. Stop on API/CLI change, comparable ratchet
regression, authoritative full-audit exit 2, or failure to isolate
implementation from generated and handoff-only state.
