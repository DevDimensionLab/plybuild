# Quality Upgrade Handover

Generated: 2026-08-26T18:40:54+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P4.9 product implementation:
  `746a5abbdc5b879201bc83690b3c633831dca81b`.
- Preceding P4.8 product implementation:
  `67344a800db0030e1aceb908562eef8c8153394b`.
- Truthful P3.54 audit-apparatus repair:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  the final commit changes only planning, archives, and launcher mutable state.
- No push, merge, release, publication, distribution, stash, revert, successor
  launch, or worktree removal was performed.

All earlier focused P3 and P4 implementation commits and their clean
checkpoints remain recorded in `docs/plan/quality-upgrade.md`. The separate
operational continuity implementation is `1b85711`, and the focused T15 repair
is `ce736a2`; neither changes a Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P4 is active and P5-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T184054+0200-close-loopback-seam.md`.
The P4.9 context archive is answered history and links reciprocally to that
archive. The graph has exactly one NEXT tail. The archived prompt SHA-256 is
`6811bbf0d83e0e2edc38fd734098d10243397caf9943ac37babefae32db941c9`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P4 Move 9

Product commit `746a5ab` adds only `pkg/context/context_test.go`. Both production
Context files, exported fields and method signatures, private package logger,
config and project ownership, file and Maven calls, callers, errors, logging,
writes, and observable behavior remain unchanged. Nine focused contracts bring
the suite to 449 tests across all 27 packages and cover `pkg/context` at 97.6%
statement coverage.

The contracts characterize exact package logger initialization and `SetLogger`
identity restoration; zero and directly owned public state; non-recursive and
recursive project discovery; append, filesystem traversal, and partial-result
order; missing-root and invalid-project behavior; empty project handling;
cloud-default merge and errors; type skips; stealth, dry-run, dirty-repository,
job, nil-job, and write behavior; root selection; profile assignment, creation,
retention, and failure; configured authenticated and anonymous repositories;
legacy default repository selection; and representative exact logs and errors.
Value copies, shared pointer fields, callback order, early returns,
continuations, and non-empty iterated and recorded populations are explicit.

All test-created directories and files stay below `t.TempDir()` and use the two
central guarded writers. Each logger, environment, homedir-cache, and
working-directory mutation restores the exact original identity and value. No
parallel process-state test, socket, network access, external program, sleep,
timed wait, repository fixture, leaked hook, entry, file, environment, or
process-state change is introduced.

Two production statements remain unreachable through deterministic current
interfaces: `file.FindAll` suppresses callback walk errors before Context can
observe them, and `maven.DefaultRepository` uses
`os/user.Current().HomeDir` rather than `HOME` and exposes no error seam. The
tests preserve both existing designs. Preferred walk-error propagation and a
caller-owned Maven home or repository selector require later explicit product
scope.

## T15 Apparatus Checkpoint

The P3.54 T15 repair remains unchanged. Separate verified execution and
structured replicas preserve historical test effects while measuring a
pristine exact `5635d50` checkout plus the declared inventory overlay. No
output path is excluded or deleted, and every effect stays in the retained
manifest. The complete 15-control proof passes and reproduces old scorecard
`d420887d73aabf496ff276fcc55d13ad379ac49c9322fad58808b5e28fdba7df`,
new scorecard
`5fb3226009cfbf0d29f63fa03592157cce4efcec6e38583b64a86f6288e89490`,
normalized stored raw body
`cf23c9dca987acd3a966693f933c4c7eca9f7d11c51f0fbabfe1d697f3d7497f`,
Q3.9, and all 228 numeric debt leaves. P4.9 did not modify the apparatus,
parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `746a5ab` reports:

- Measured commit tree: `7bfabae149bffc5dfd2ad59f90abb08e29c982f2`.
- Structured scorecard SHA-256:
  `48a58e58abf53254ec318de93266a683f60e6ac4a194e02b187c76c170a823b2`.
- Absolute L0: 8 of 8.
- 449 test functions, zero skipped; all 27 packages have tests.
- `pkg/context`: 97.6% statement coverage; all deterministically reachable
  public-method statements covered.
- Q0.6: 29 guarded safe-writer sites, 23 write and 6 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 0 of 27 packages has no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 0 direct external sites outside five valid declared adapters of 27
  production effect sites.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects has an executable harness.
- Q3.4: 0 state-claim phrases across 93 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 13 documented non-passing criteria, never 2.
- Comparable ratchets: six improved, two held, zero regressed, zero
  non-comparable.
- Clean identity: zero dirty paths.

The focused seven-criterion scorecard under
`/private/tmp/ply-p4-context-focused-audit.uaxupD` has SHA-256
`f937c9be84913156be11dfa9a937565959cfcde6ba1b9389bea7b2f337e0d7c6`.
The full clean audit is under
`/private/tmp/ply-p4-context-full-audit.iJ6p0e`. Pinned tools remain under
`/private/tmp/ply-p4-tools` and `/private/tmp/ply-p358-tools`. The empty-HOME
count-2 proof is under `/private/tmp/ply-p4-context-hermetic.Dzf0oW`.

## Next Objective

Continue P4 with the eighth declared seam: change only the private interactive
HTTP server address in `pkg/webservice/api.go` from `:7999` to exact IPv4
loopback `127.0.0.1:7999`, preserving port 7999, the exact package server
identity, handler order, adapter selection, listen/log/shutdown behavior,
browser endpoint and every caller. Prove the address without opening a socket.

Add the missing non-executable `scripts/mutate-interactive-build` P4 seam
driver and reciprocal `scripts/test-mutate-interactive-build` meta-contract,
following the exact seven established seam-driver patterns. The driver contains
the inventory label `8. interactive server binds only to loopback` exactly
once, declares no mutation operation, and runs a non-empty exact set of command,
webservice, server-adapter, API, and callback contracts. This move must improve
Q1.4 from 7/8 to 8/8. It does not authorize an executable P5 mutation harness,
T1-T10 work, inventory changes, manual evidence, another seam, or another
production behavior.

## Verification Notes

- Focused context and relevant root/command caller tests pass; context coverage
  is 97.6%, with focused shuffled count-10, race, vet, and pinned lint green.
- Actual API/CLI, CLI surface, fresh executable subprocess, and compatibility
  meta-contracts pass with reports and caches outside the worktree.
- Complete build and uncached tests, race, vet, exact pinned lint,
  `make preflight`, `make test`, install, Make contracts, and all 62 launcher
  controls pass.
- All production-script meta-contracts, all four host acceptance flows, and the
  repaired 15-control audit meta-suite pass.
- The first focused audit exposed three direct test directory writes. Their
  test-only replacement with the central guarded copy writer was included by
  amending the single product commit before final measurement; exact Q0.6 then
  passed with zero unsafe writes.
- The first `make test` invocation hit the documented partial-raw-log
  signal-fixture race at control 26; its immediate unchanged standalone
  launcher run and full Make rerun passed all 62 controls.
- Empty-HOME `go test ./... -count=2` passes with isolated writable state and
  the existing read-only module cache.
- Full clean audit: expected exit 1, 13 non-passing criteria, six improved, two
  held, zero regressed, zero non-comparable, and zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the complete P4 entry and
checkpoint gate, both design documents, `.quality/inventory`, complete
interactive webservice and command code, the server adapter and all related
contracts, every existing seam driver/meta-test, the import-aware scanner,
API/CLI contracts, and the T15 repair and baseline reproduction README before
editing. Confirm branch, HEAD, clean status, reciprocal archive links, launcher
`--check`, and exact product commit `746a5ab`.

Make one focused loopback-seam implementation commit, then the normal separate
continuity commit. Stop before an executable P5 harness, T1-T10, manual
evidence, another package or seam, P6-P8, publication, or distribution. Do not
push, merge, stash, revert, launch a successor, or remove the worktree.
