# Quality Upgrade Handover

Generated: 2026-08-27T02:15:52+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P5.4 implementation commit:
  `da1cc888a9d6059e1853ae9051cfa00e3c1875bb`.
- Its exact parent is the launch continuity commit
  `0b1f81579fa4269ed7c9aa005f2835d86d4a8b51`.
- Measured implementation tree:
  `ebe0ac95221cd47538dae207e0fae1a07739cbd2`.
- Clean status SHA-256:
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
- After handoff, obtain the new continuity HEAD with `git rev-parse HEAD`; its
  exact parent must be `da1cc88`.
- The implementation changes only `scripts/mutate-template` and
  `scripts/test-mutate-template`. No production Go, Go test, inventory, audit,
  parser, scanner, baseline, acceptance, API, CLI, packaging, or dependency
  file changed.
- No push, merge, release, publication, distribution, stash, revert,
  successor launch, or worktree removal was performed.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P5 is active and P6-P8 are queued.
Its active archive is
`docs/plan/agent-sessions/2026-08-27T021552+0200-build-p5-file-shell-harness.md`.
The `template` archive is answered history and links reciprocally to the new
tail. The graph has exactly one NEXT archive. Only the launcher's mutable
header and prompt regions changed; the stable execution region is unchanged.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. The audit counts
ignored and untracked bytes, so ordinary and ignored status must both be empty
at a measured checkpoint.

## P5.4 Result

The fourth of eight mutation subjects, `template`, is complete:

- The executable harness declares ten deterministic, unique mutations in
  `pkg/template/template.go`: source-root selection, target path component
  ordering, copy argument ordering, copy error gating, walk directory
  classification, root-POM detection, nested-POM exemption, render-file
  exemption, returned file ordering, and markdown write bytes.
- The inventory seam label remains byte-exact:
  `6. template copy keeps source before destination`.
- Each declaration binds production syntax occurring exactly once to a
  non-empty exact named `pkg/template` test population. Init code, test files,
  fixtures, generated/vendor code, other subjects, equivalent replacements,
  and known non-compiling replacements are excluded.
- One clean unmodified external control ran every exact selection. Each mutant
  used a fresh external Git archive, cache, HOME, config, and temp root; its
  changed package compiled separately; and its selected population ran under
  `go test -json` with exact run and terminal-action validation. Tooling,
  setup, discovery, selection, compilation, and unrelated failures are
  unusable; only a selected test failure counts as killed.
- Exact totals are `declared=10`, `killed=10`, `survived=0`, and
  `unusable=0`. No reachability, observability, or controllability gap appeared,
  so no production or test repair was made.
- T1-T10 prove failure for empty/duplicate manifests, unauthorized paths,
  zero/multiple replacements, empty or inexact tests, broken control,
  uncompiled/unexercised mutants, false accounting, unclassified survivors,
  repository-local artifacts, and non-deterministic declarations/totals.

Every mutation and killing population is in the retained report
`/private/tmp/ply-template-authoritative.uHrYxm/report.txt`; its SHA-256 is
`e769cb49b2ae6fd3da15b207acab8ef11c402f5bff171dc6e150623b93279e95`.
All retained control/mutant checkouts and logs are under
`/private/tmp/ply-template-authoritative.uHrYxm/work`. The independent T1-T10
meta-log is `/private/tmp/ply-template-authoritative.uHrYxm/meta.log`, SHA-256
`cde7a149274b7f4b968a40b1ec5e4f7769c5d6d4243edcb40783a4fa19ba91bd`.

## Authoritative Measurement

The refreshed external evidence and audits are under
`/private/tmp/ply-template-evidence.dFB57Q`. The exact 245 Q1.6, 34 Q1.7, and
69 Q1.9 named populations were rerun at `da1cc88`; every named test emitted one
run and one pass. Manifest SHA-256 values are
`845cc6bdfaa0e4e84956a96a162e84ed0372238f067f8d2b22d05f7a22307f27`,
`cb028c078c2f61283a974ec308a3a0a2de1a6694cd2a60a130b6f4dd7f357d0d`,
and `c6488bd00c342304d8d1fa46a4af6dce33c9b095ec6ad1572bc9e0de40252034`.
Event-log SHA-256 values are
`333924496a178d0c76963cb5f6b41694bc02cc38222be4a80ea7d82d0b2e7773`,
`678268f152fdd21a54bbb578e9645d8c5099bf73e2ad7ed7716e156f27368b55`,
and `d0737df73a8ff85cc6f0fe0e1fdcd3107247837d4b8481229e638f3800f2bacc`.

The external schema-2 document contains refreshed Q1.6, Q1.7, Q1.9 receipts
and Q2.4 receipts for all four completed P5 subjects. Its SHA-256 is
`f235776be8c9fc454906415a1f983aeebf6ef29adbd4a8665a83e3bed76fe050`.
Evidence-object SHA-256 values are:

- Q1.6: `b4d64092ab5536dd6937fec80eb86a7a71a14af6d29ea11fcc26e3509b7bd41a`;
- Q1.7: `a7e767f194c7a9e125efa1cf4416dab3790aa2c565432c8109d0c32f6cbc57cc`;
- Q1.9: `27c7531b81c5aa70307ec634b4518c5feee44011a6a045b5ede0bfbbf563a390`;
- Q2.4: `bd2899300363740c159c07b41dfa76201371ec5e74b29ac545e20033ebb0a433`.

The focused Q1.6/Q1.7/Q1.9/Q2.4 audit exits 0 with scorecard SHA-256
`28732aa8d367b86cfb74a3e25d42a9909d01e8fc4564d4488d720b04675c4e8f`.
The full authoritative audit exits 1, never 2, with scorecard SHA-256
`fcfddf9b5794dabb1d2bdbbc19b04d594522bf7efbd7db7432a35882a6ee797c`.
It records L0 8/8, all nine L1 rows PASS, Q2.1 improved to 4/8, Q2.4 PASS,
seven improved ratchets, one held, zero regressed, zero non-comparable, and
zero dirty paths. Eight P5-P8 rows remain non-passing.

The separate no-evidence Q2.1-Q2.4 view is under
`/private/tmp/ply-template-evidence.dFB57Q/no-evidence-q2`, scorecard SHA-256
`1f29a6fdd16484ec3c4741772dde20ce3b076025532db90f1231506f7ee137d3`.
It exposes the exact unratcheted state: Q2.1 is 4/8; Q2.2 and Q2.3 identify
only the remaining legacy subjects; Q2.4 is unmeasurable without the external
run receipt.

## Gate Result

Gate logs and external caches are under
`/private/tmp/ply-template-gate.ATRg0s`. API/CLI and explicit entry/subprocess
compatibility, compatibility meta-tests, pinned golangci-lint 2.12.2 with zero
issues, complete uncached tests across all 27 packages, race, vet, `make test`,
all 62 launcher controls, Make distribution/install/lint/preflight contracts,
all four host acceptance flows, exact empty-HOME count-2, the standalone
15-control audit meta-suite, and complete preflight pass. API and CLI report
SHA-256 values remain
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
The complete preflight, acceptance, audit-meta, empty-HOME count-2, and focused
entry/subprocess log SHA-256 values are respectively
`0f1dcb95f1564f12e745f476175f95ea4b971be39e7004da70173974a5cf5175`,
`36b484c628424a75eac04fd1550952ef4b8ecea87123c15e2c6d2b6b02400f00`,
`ae853e4750310bf06d5a580067531fd2a08dd4c869efc04675c9b5fa10272baa`,
`8243e0a98eb80b1f71c36298aceb2d3ededf50a0daaba96e22a8d82dcbfe2cd1`,
and `60b9178c738cbd3c7bb100cad44f2613376564962ddceb477056a7503854371e`.
Ordinary and ignored status were empty before authoritative measurement and
after the complete implementation gate.

## Next Objective

Convert exactly the fifth inventory subject, `file-shell`, whose production
roots are `pkg/file` and `pkg/shell` and whose declared harness is
`scripts/mutate-file-shell`. Replace the existing non-executable P3 seam driver
and its `scripts/test-mutate-file-shell` meta-test with one executable mutation
harness and T1-T10 falsifiability test. Define at least eight meaningful
mutations after inspecting the complete production and test population.
Preserve both seam labels exactly:
`1. git clone keeps URL before target directory` and
`7. git commit keeps target directory before message`.

Use all four completed P5 harnesses as methodology references. Require one
clean external control, fresh external copies and caches, exact test discovery,
successful mutant compilation, selected JSON run/terminal actions, and exact
totals `declared == killed`, `survived == 0`, `unusable == 0`. If a genuine
survivor appears, classify it as reachability, observability, or controllability
before the smallest in-subject test or private seam repair. Do not start a
sixth subject. P5 remains active afterward.

## Start And Stop

Confirm branch, HEAD, exact ancestry, clean and ignored status, reciprocal
links, launcher `--check`, and the authorized checkpoint block before editing.
Read the rolling handover, linked NEXT archive, complete P5 entry/gate, both
designs, `.quality/README.md`, `.quality/inventory`, complete mutation audit and
parser logic, all completed and legacy harness/meta-test pairs relevant to the
move, Make/preflight contracts, the process adapter, and every production/test
path under `pkg/file` and `pkg/shell` before selecting mutations.

Stop before another subject, inventory/audit/parser/scanner/baseline changes,
acceptance expansion, P6-P8, Go/dependency upgrades, exported API/CLI changes,
packaging, publication, or distribution. Keep generated artifacts external.
Do not push, merge, stash, revert, launch a successor, or remove the worktree.
