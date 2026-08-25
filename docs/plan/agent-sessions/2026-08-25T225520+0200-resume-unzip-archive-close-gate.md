# Agent Session: Resume Unzip Archive Close Gate

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T225520+0200-resume-unzip-archive-close-gate`
Created: `2026-08-25T22:55:20+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `229b7729de108c595723ef276892d0c75c93adf0b629e7530e0b87c280530a97`
Previous: [2026-08-25T220107+0200-resume-shell-unzip-archive-close.md](2026-08-25T220107+0200-resume-shell-unzip-archive-close.md)
Next: [2026-08-25T232218+0200-resume-p354-audit-harness-blocker.md](2026-08-25T232218+0200-resume-p354-audit-harness-blocker.md)
Outcome: paired complete meta-suites at `70bee0e` and `3bd07e9` fail identically in T15, proving an inherited audit-harness blocker; no product or apparatus code changed

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Finish the checkpoint gate for the already-committed focused P3 Unzip archive-
close move at `3bd07e9`. Preserve that implementation exactly. Determine why
the unchanged audit meta-suite's T15 baseline-reproduction control now sees nine
ignored Maven/template fixture paths in its pinned pre-hermetic baseline clone,
although T15 passed at predecessor checkpoint `70bee0e`. Restore a truthful
green gate only through environment or invocation correction already supported
by the repository. If the same failure reproduces on the predecessor under the
same environment, record the inherited blocker precisely and stop without
changing the audit, baseline, inventory, fixtures, product code, or tests.

# Authorized Roadmap

P3 remains active, P4-P8 remain queued in the machine-readable block in
`docs/plan/quality-upgrade.md`, and the launcher remains NEXT. This mission
authorizes read-only diagnosis, temporary clean shared clones under
`/private/tmp`, exact reruns of the established audit meta-suite and checkpoint
gate, truthful P3.54 planning/handover corrections, and the normal separate
handoff-only commit.

It does not authorize another production-effect move, an audit or scanner
change, baseline or inventory change, test-fixture change, cleanup added inside
T15, identity weakening, ignored-path exclusion, a Go/product/test change,
commit amendment, revert, adapter work, plugin-diagram work, clock/server work,
mutation work, or later-roadmap implementation. Do not manufacture a green
result with a wrapper, hook, patched temporary instrument, filesystem watcher,
manual mid-test cleanup, cached result, or altered baseline checkout.

# Measurements At Start

Implementation commit `3bd07e99c3170fd0f1c0dda076252c6eb5833db7`
routes only the deferred archive close through
`filesystem.CloseReader(dependencies.Files, r)` and widens that existing
operation from `io.ReadCloser` to `io.Closer` mechanically. It has 360 tests
across 19 of 25 packages. Q0.6 is 26 guarded sites, 21 write and 5 copy, with
zero unsafe direct test writes. Q1.1 is 6/25, Q1.2 is zero, Q1.3 is 10/32,
Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is zero phrases across 69 Markdown
files. The focused audit exited 1 with four improved, two held, zero regressed,
and one non-comparable selected criterion.

The authoritative clean full audit at `3bd07e9` is valid: exit 1 for 15
documented findings, L0 8/8, five improved, two held, zero regressed, one
non-comparable, and zero dirty paths. Its structured report is
`/private/tmp/ply-quality-move54-clean-audit/scorecard.json` while retained.
The handoff-only Q3.4 audit adds this archive to the denominator and passes with
zero phrases across 70 Markdown files and zero ratchet regressions.

Focused and caller tests, API/CLI/subprocess compatibility, the standalone
62-control launcher, Make meta-contracts, complete `make test`, install,
uncached tests, race, vet, all four host flows and their meta-contracts, and
empty-HOME count-2 pass. Full preflight passes every component until audit meta
T15. Three complete audit-meta attempts pass T1-T14 and fail T15 with exit 2:
`AUDIT BROKEN: inventory overlay mode permits only untracked
.quality/inventory`. Supplying the existing external `GOMODCACHE` removes
read-only module-cache cleanup noise but leaves the same T15 result.

The retained diagnostic clone at
`/private/tmp/test-ply-quality-audit.qEdLnB/baseline-repository`, if still
present, showed exactly these extra measured paths after the old upstream test
run:

- six files under `pkg/maven/test/analyze/target/`;
- `pkg/template/test/target-simple-template/src/main/java/no/ply/template/target/DummyConfiguration.kt`;
- `pkg/template/test/target-test-template/test.properties`;
- `pkg/template/test/target-test-template/textfile.txt`;
- plus the authorized `.quality/inventory` overlay.

The pinned source commit is `5635d50`, before hermetic fixture commit
`80b43ba`. Directly running its template tests creates the three template
paths. The checked-in T15 calls the old vendor audit before its structured
identity check and contains no intervening cleanup. Do not treat these learned
facts as permission to reshape the instrument.

# Role And Boundaries

Work autonomously in `/Users/perottochristensen/github/ply/upgrade-quality` on
`codex/upgrade-quality`. Begin at the clean handoff commit following `3bd07e9`.
Do not change or amend `3bd07e9`. Do not push, merge, publish, distribute,
stash, revert, remove the worktree, or launch a successor.

Preserve public `Unzip(string, string) ([]string, error)`, the exact deferred
CloseReader request, the `io.Closer` widening, every complete double, archive
and entry behavior, all other product/test code, and `.quality/inventory`.
Keep Go 1.18 and Bash 3.2 compatibility. Keep reports, clones, Go state, and
diagnostic output outside the measured tree.

# Required Reading

Read the rolling handover, this archive, P3.54 and the checkpoint gate in
`docs/plan/quality-upgrade.md`, both design documents, `.quality/inventory`,
complete `.quality/tools/test-quality-audit.sh`, the wrapper, parser identity
logic, vendored audit test execution, baseline migration metadata, baseline
README, commits `5635d50`, `80b43ba`, `70bee0e`, and `3bd07e9`, and the exact
prior verification notes. Confirm branch, HEAD, clean status, archive links,
and `./codex-dev-start.sh --check` before diagnosis.

# Three Moves

1. Reproduce T15 without product mutation. Use fresh clean shared clones or
   detached temporary worktrees below `/private/tmp` to compare the complete
   audit meta-suite at `70bee0e` and `3bd07e9` under byte-identical environment,
   toolchain, PATH, HOME/module-cache policy, TMPDIR, GOTMPDIR, and GOCACHE
   setup. Capture exact heads, Go version/environment differences, exit codes,
   dirty-path populations, and whether the old checkout actually succeeded.
   Do not modify either measured clone between its old upstream run and parser.

2. If an established supported environment/invocation difference explains the
   regression, apply only that external invocation correction and rerun full
   `make preflight`, the standalone audit meta-suite, and the clean checkpoint
   audit. Do not encode an environment workaround in tracked files. If both
   commits reproduce the same T15 failure, stop product work and document that
   inherited harness blocker; do not weaken or edit the audit to pass.

3. If and only if the complete gate becomes green, update P3.54 from
   implementation measurement to a clean checkpoint, retain the exact clean
   audit result, answer this archive, and create one linked NEXT archive for
   the next isolated P3 effect. Otherwise preserve P3.54 as implemented but
   awaiting its gate, rewrite the handover with exact reproduction evidence,
   and leave one valid resumable NEXT archive. In both cases run the launcher
   contract and `--check`, verify Q3.4 and the archive graph, and make only the
   separate `docs: prepare next agent session` commit.

# Automatic Handoff

Before the session ends, leave the worktree clean and committed with a truthful
rolling handover and exactly one NEXT archive. Do not launch it. Keep P3 active
and P4-P8 queued. COMPLETE is invalid until every authorized checkpoint through
P8 is complete.
<!-- CODEX_SESSION_PROMPT_END -->
