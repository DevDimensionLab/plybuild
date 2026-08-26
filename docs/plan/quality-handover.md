# Quality Upgrade Handover

Generated: 2026-08-26T06:18:48+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head: 3bd07e99c317.
- After launch, obtain the handoff head with `git rev-parse --short=12 HEAD`;
  the restart commit contains only continuity/planning state after the focused
  implementation commit.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this move are 5c6f2fa, 03d6242, f5ee37d,
ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce, a7eb3ef,
e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597, 9e2d669,
a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d, d6cb593,
82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438, 527a8b9,
886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc, 85f4c2b,
0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, da7eebf, 89918bd, 2f0a072,
0128cd0, 70bee0e, and 3bd07e9 in roadmap order. The separate operational
continuity implementation is 1b85711 and changes no Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because the user has now explicitly authorized
a minimal repair of the inherited P3.54 audit-meta apparatus failure and P4-P8
remain queued. Its active archive is
`docs/plan/agent-sessions/2026-08-26T061848+0200-resume-p354-audit-harness-repair.md`.
The blocker archive is answered history, reciprocal links are connected, and
there is exactly one NEXT tail. T15 orchestration is the first and narrowest
authorized repair surface; directly related reproduction recipe, metadata,
apparatus-owned fixtures, and structured parser are conditional on evidence of
necessity.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 54 Implemented

Public `pkg/shell.Unzip(string, string) ([]string, error)` still selects the
complete private `unzipDependencies` composition with `filesystem.System()`.
Only its deferred archive close changed from direct `_ = r.Close()` to
`_ = filesystem.CloseReader(dependencies.Files, r)` inside the same anonymous
defer immediately after successful archive open.

The existing CloseReader parameter widened mechanically from `io.ReadCloser`
to `io.Closer` in the filesystem interface, zero-safe helper, exact system
implementation, adapter recording double, and every complete filesystem double.
No operation was added. This accepts the exact `*zip.ReadCloser` returned by
`OpenZipReader`, invokes one dependency request on every later return, and
continues to ignore the exact close result. Archive-open failure makes no close
request.

Archive and entry traversal, entry order and bytes, zip-slip behavior,
filenames and partial results, MkdirAll, OpenFile, Copy, output-file Close,
entry-reader Close, every error and close precedence, every caller, public API,
and inventory remain unchanged. There is no fallback, retry, logging, extra
defer, cleanup, normalization, environment change, or working-directory change.

Two new top-level tests bring the suite to 360. The complete private Unzip
contracts record exact archive identity and prove one final deferred request on
successful traversal, after reached entry-reader closes, on zip-slip, and on
every established later error; no request after open failure; ignored injected
archive-close errors without result or precedence changes; entry/archive-close
distinction; complete dependency selection; non-empty archive-close recording;
and no unrelated filesystem request. The guarded archive fixture remains below
`t.TempDir()` and the focused red/green contracts launch no process or network
request.

## Measured Quality State

The authoritative clean full audit from implementation commit `3bd07e9` is
valid and reports:

- Absolute L0: 8 of 8.
- 360 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 10 direct external sites outside five declared adapters of 32
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 69 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The clean report is
`/private/tmp/ply-quality-move54-clean-audit/scorecard.json`; the final focused
report is
`/private/tmp/ply-quality-move54-focused-audit-final/scorecard.json` while
retained. Reports, tools, caches, diagnostic clones, and empty-HOME state stay
outside the measured tree.

## Inherited Gate Blocker

All move-specific and product gates pass except audit meta-control T15:

- Focused shell/filesystem/Spring and all relevant caller tests: PASS.
- API, CLI, and subprocess compatibility: PASS.
- Standalone launcher: PASS all 62 controls after one first-run signal/partial-
  log timing miss; later complete runs also pass.
- Make meta-contracts, pinned lint, build, tests, vet, and preflight components
  through audit T14: PASS.
- Complete `make test`, separate install, uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build flows and meta-contracts: PASS 4/4.
- Empty-HOME count-2: PASS.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable; Q1.3 is exactly 10/32.
- Clean full audit: valid exit 1, 15 findings, zero dirty paths and zero
  comparable regressions.
- Audit meta-suite at the implementation checkpoint: prior complete attempts
  pass T1-T14, then T15's old structured parser exits 2.

T15 reports `inventory overlay mode permits only untracked
.quality/inventory`. Its temporary pinned `5635d50` checkout contains the
authorized inventory plus nine ignored paths created by the old upstream test
run: six Maven compiler/status/class outputs and three template merge outputs.
That commit predates hermetic fixture commit `80b43ba`. Supplying the existing
external `GOMODCACHE` eliminates read-only cleanup noise but does not change the
identity failure. The audit, baseline, inventory, fixtures, Maven/template code,
and scanner are unchanged by move 54, so the mission boundary did not permit a
tracked workaround or cleanup inside T15.

Fresh paired detached clones now prove the failure is inherited. `70bee0e` and
`3bd07e9` were each checked out sequentially at the same pathname and run from
the same freshly recreated external environment pathname with identical PATH,
HOME, XDG config, TMPDIR, GOTMPDIR, GOCACHE, existing GOMODCACHE, Go selectors,
locale, and toolchain. The exported-environment files are byte-identical with
SHA-256 `963ad106f66f7770704a13b3ef619079ae19748b1aafe678700f9233bdc833ef`;
the Go/toolchain evidence files are byte-identical with SHA-256
`8fc0505f60a4950583fd56222b64a0a175b44fbfc55aff409467602d19592c5f`.
Both use Go 1.26.2 darwin/arm64, `CGO_ENABLED=0`, `GOENV=off`, `GOWORK=off`,
empty `GOFLAGS`, and the same existing module cache.

Both complete meta-suites exit 1, pass T1-T14, and produce identical normalized
log bodies with SHA-256
`bc83831ae61962b92c3233de04df52f70e3f2a2fdd6d871fb4d0db3f451b733b`.
In each T15 run the old upstream baseline audit really succeeds with expected
exit 1; the following old structured parser exits 2 on the overlay identity.
Each pinned clone is exactly `5635d50` and has the inventory overlay plus these
three ignored paths:

- `pkg/template/test/target-simple-template/src/main/java/no/ply/template/target/DummyConfiguration.kt`;
- `pkg/template/test/target-test-template/test.properties`;
- `pkg/template/test/target-test-template/textfile.txt`.

The controlled PATH intentionally contains only standard supported tool
locations and has no `mvn`. The retained developer-environment clone confirms
that ambient SDKMAN Maven adds exactly six files below
`pkg/maven/test/analyze/target/`, explaining why prior attempts saw nine ignored
paths. Removing Maven explains that population difference but is not a green
invocation correction because the three unconditional template outputs remain.
The paired evidence is retained below
`/private/tmp/ply-p354-gate.B7akzM` while present.

## Decisions And Learned Facts

1. `*zip.ReadCloser` implements `io.Closer` but not `io.ReadCloser`; the
   user-authorized parameter widening is the smallest existing-boundary change.
2. Red evidence was focused: ten private Unzip tests rejected an empty archive-
   close population before the one production expression changed. No process or
   public network request ran.
3. CloseReader receives the exact archive identity once and after every reached
   entry operation. Injected close errors are ignored on both success and
   earlier-error paths.
4. Q1.3 holds at 10/32 because both the old concrete archive Close and the new
   adapter request are absent from the reported violation set. No scanner or
   inventory change forced a number.
5. The clean full audit proves the implementation commit itself is measurable,
   clean, and regression-free even though the separate meta-harness T15 cannot
   reproduce its much older source checkout in this run.
6. T15 must not be made green by a wrapper, hook, manual mid-run cleanup,
   patched temporary instrument, ignored-path exclusion, weakened identity, or
   altered baseline checkout.
7. The complete predecessor and implementation meta-suite runs fail
   identically. The earlier recorded green predecessor gate is not reproducible
   from the preserved repository and environment evidence; the blocker is not
   caused by move 54.
8. PATH exposure of real Maven explains six of the nine ignored outputs, but a
   no-Maven PATH still fails on three template outputs. No supported external
   invocation correction found by this mission can make the existing T15
   identity truthful.
9. The user supplied explicit scope for a minimal P3.54 audit-apparatus repair:
   start with T15 orchestration and expand only when necessary to the directly
   related recipe, metadata, apparatus fixtures, or structured parser. Product,
   historical source, baseline debt, identities, and inventory remain locked.
10. Plugin diagrams remain the next likely isolated production-effect family,
   but no such move may begin until P3.54's checkpoint gate is truthful.

## Next Objective

Repair the inherited T15 apparatus failure within the user-authorized boundary.
Begin with `.quality/tools/test-quality-audit.sh` orchestration; expand only when
focused evidence proves a directly related reproduction recipe, metadata,
apparatus-owned fixture, or structured-parser change is necessary. Preserve the
exact `5635d50` source checkout, stored baseline debt and identities, inventory,
and measured-tree truth rather than excluding or cleaning evidence after fact.

Keep product implementation `3bd07e9` unchanged. Do not claim P3.54 as a clean
checkpoint or begin plugin-diagram, clock/server, mutation, or later-roadmap
work until the repaired complete checkpoint passes truthfully.

## Verification Notes

Completed from implementation commit
`3bd07e99c3170fd0f1c0dda076252c6eb5833db7`:

- Focused red and green Unzip/filesystem/Spring evidence: recorded and PASS.
- All relevant caller packages: PASS.
- API/CLI/subprocess compatibility: PASS.
- Launcher contract: PASS all 62 controls on complete reruns.
- Make contracts, full test, install, pinned lint, uncached/race tests, vet:
  PASS outside the T15 preflight tail.
- Four host flows and their meta-contracts: PASS.
- Empty-HOME count-2: PASS.
- Focused audit: valid exit 1, Q1.3 10/32, zero regressions.
- Clean full audit: valid exit 1, 15 findings, L0 8/8, five improved, two held,
  zero regressed, one non-comparable, zero dirty paths.
- Audit meta T1-T14: PASS; T15: repeated exit 2 for the pinned baseline clone's
  nine ignored fixture outputs plus its authorized inventory overlay.
- Paired predecessor/implementation audit meta: both exit 1; T1-T14 PASS; old
  upstream baseline exit 1; old structured baseline exit 2; normalized output
  bodies and controlled environments are byte-identical.
- Controlled dirty population: inventory overlay plus three template outputs;
  developer PATH adds six Maven outputs for nine ignored paths total.
- No full preflight or clean audit rerun followed the paired failure because the
  mission explicitly requires stopping on the inherited predecessor blocker.
- Authorization handoff-only Q3.4: exit 0, zero phrases across 73 Markdown
  files, zero ratchet regressions.
- Implementation commit: `3bd07e9` (`refactor: route unzip archive close
  through filesystem adapter`).

## Start And Stop

Read this handover, the linked NEXT archive, P3.54 and the gate, both design
documents, inventory, audit meta-suite/wrapper/parser/vendor, baseline migration
metadata and README, and commits `5635d50`, `80b43ba`, `70bee0e`, and `3bd07e9`.
Confirm branch, HEAD, clean status, reciprocal links, launcher `--check`, and
the explicit authorization in the active archive and `.agent-task/current.md`.

Repair only the authorized T15 apparatus surface. Stop before product mutation,
historical source, baseline debt, instrument identity, inventory, or unrelated
scanner changes, another P3 effect, plugin diagrams, clock/server work, Q1.4,
P4, mutation, Docker, distribution, or publication. Always stop on any proposal
that trades away measured-tree identity to make T15 pass.
