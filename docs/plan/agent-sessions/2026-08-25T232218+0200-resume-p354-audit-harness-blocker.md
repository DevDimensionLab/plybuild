# Agent Session: Resume P3.54 Audit Harness Blocker

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T232218+0200-resume-p354-audit-harness-blocker`
Created: `2026-08-25T23:22:18+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `89d8951db6d47546fac025421da17bc22f398e5a1b2c685a74750c615802570e`
Previous: [2026-08-25T225520+0200-resume-unzip-archive-close-gate.md](2026-08-25T225520+0200-resume-unzip-archive-close-gate.md)
Next: [2026-08-26T061848+0200-resume-p354-audit-harness-repair.md](2026-08-26T061848+0200-resume-p354-audit-harness-repair.md)
Outcome: user authorized the recommended minimal T15 audit-apparatus repair boundary; no repair was attempted in this session

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Keep the focused P3.54 Unzip archive-close implementation at `3bd07e9` exactly
as committed and keep its checkpoint gate blocked. The complete unchanged audit
meta-suite now fails identically at predecessor `70bee0e` and implementation
`3bd07e9`; this is an inherited T15 apparatus failure, not a product regression.
Resume only after the user explicitly authorizes an audit-apparatus repair
boundary. Without new authorization, report the required scope decision and
stop without mutation.

# Authorized Roadmap

P3 remains active with moves 1-53 complete and move 54 implemented but awaiting
a clean complete checkpoint gate. P4-P8 remain queued in the machine-readable
block in `docs/plan/quality-upgrade.md`; the launcher remains NEXT. Public
`Unzip(string, string) ([]string, error)`, the exact deferred
`filesystem.CloseReader(dependencies.Files, r)` request, its `io.Closer`
widening, every complete double, all archive/entry behavior, product code,
tests, fixtures, audit, baseline, inventory, and scanner are unchanged.

# Measurements At Start

Implementation commit `3bd07e99c3170fd0f1c0dda076252c6eb5833db7`
has 360 tests across 19 of 25 packages. Its authoritative clean full audit is
valid: exit 1 for 15 documented findings, L0 8/8, five improved, two held, zero
regressed, one non-comparable, and zero dirty paths. Q0.6 is 26 guarded sites,
21 write and 5 copy, with zero unsafe direct test writes; Q1.1 is 6/25, Q1.2 is
zero, Q1.3 is 10/32, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is zero phrases.

# Inherited Blocker Evidence

Fresh detached shared clones at exact heads `70bee0ed16af4b1e2dec179d6a0ba1ae4237c9e0`
and `3bd07e99c3170fd0f1c0dda076252c6eb5833db7` were run sequentially
from the same clone pathname and the same freshly recreated external environment
pathname. PATH, HOME, XDG config, TMPDIR, GOTMPDIR, GOCACHE, existing
GOMODCACHE, locale, Go selectors, and toolchain were identical. Exported
environment evidence has SHA-256
`963ad106f66f7770704a13b3ef619079ae19748b1aafe678700f9233bdc833ef`;
toolchain evidence has SHA-256
`8fc0505f60a4950583fd56222b64a0a175b44fbfc55aff409467602d19592c5f`.
Both used Go 1.26.2 darwin/arm64, `CGO_ENABLED=0`, `GOENV=off`, `GOWORK=off`,
empty `GOFLAGS`, and the same external module cache.

Both complete meta-suites exit 1, pass T1-T14, and have byte-identical normalized
output bodies with SHA-256
`bc83831ae61962b92c3233de04df52f70e3f2a2fdd6d871fb4d0db3f451b733b`.
T15's old upstream baseline audit actually succeeds with expected exit 1; its
following old structured parser exits 2 with `AUDIT BROKEN: inventory overlay
mode permits only untracked .quality/inventory`. Each exact `5635d50` checkout
contains the authorized inventory overlay plus three ignored template outputs.

The retained developer-environment clone additionally contains six Maven
compiler/status/class files below `pkg/maven/test/analyze/target/`. Ambient PATH
exposes SDKMAN Maven; the paired controlled PATH does not. That explains why
earlier attempts saw nine ignored paths, but removing Maven cannot restore a
green T15 because the old template tests still create three ignored outputs.
The paired evidence is retained under `/private/tmp/ply-p354-gate.B7akzM` while
present.

# Role And Boundaries

The present scope does not authorize the apparatus change required to make the
baseline reproduction hermetic while retaining exact source, debt, instrument,
and measured-tree identity. Obtain explicit user direction that names the
allowed repair boundary before editing any audit, parser, baseline recipe,
baseline metadata, fixture, test, product, or inventory file. Do not infer that
authority from this NEXT archive.

Do not manufacture green evidence through a wrapper, hook, patched temporary
instrument, filesystem watcher, manual mid-test cleanup, cached result, ignored-
path exclusion, weakened identity, altered `5635d50` checkout, baseline change,
or inventory change. Do not begin another P3 effect, plugin diagrams,
clock/server work, mutation work, or P4-P8 while P3.54 remains blocked.

# Required Reading

Read `docs/plan/quality-handover.md`, this archive, P3.54 and the checkpoint gate,
both design documents, `.quality/inventory`, the complete audit meta-suite,
wrapper, parser identity logic, vendored test execution, baseline migration
metadata and README, and commits `5635d50`, `80b43ba`, `70bee0e`, and `3bd07e9`.
Confirm branch, HEAD, clean status, reciprocal archive links, and
`./codex-dev-start.sh --check`.

# Three Moves

1. Determine whether a new user instruction explicitly authorizes an apparatus
   repair and names its boundary. Do not treat this archive as that authority.
2. If no new scope exists, make no tracked or untracked repository change;
   report the inherited T15 blocker and request the missing decision.
3. If new scope exists, re-read it against the retained paired evidence before
   acting. Preserve exact source, debt, instrument, and measured-tree identity,
   and keep any apparatus repair separate from product implementation.

# Automatic Handoff

Without new explicit scope, do not create another successor or repeat the
blocked gate. Leave this committed NEXT state intact and report the blocker.
Do not launch a successor, push, merge, publish, distribute, stash, revert, or
remove the worktree. COMPLETE remains invalid while P3 and P4-P8 are unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
