# Quality Upgrade Handover

Generated: 2026-08-27T04:01:08+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P5.6 implementation commit:
  `9c20e526288430559c249c67e3884633a7c0de05`.
- Its exact parent is the launch continuity commit
  `0759882df846f4d44d05dd89309bf53b6bce0787`.
- Measured implementation tree:
  `46eb561b44e8e029d1e09f197ed14cf639e55e87`.
- Clean status SHA-256:
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
- After handoff, obtain the new continuity HEAD with `git rev-parse HEAD`; its
  exact parent must be `9c20e52`.
- The implementation changes only `scripts/mutate-spring` and
  `scripts/test-mutate-spring`. No production Go, Go test, inventory, audit,
  parser, scanner, baseline, acceptance, API, CLI, packaging, or dependency
  file changed.
- No push, merge, release, publication, distribution, stash, revert,
  successor launch, or worktree removal was performed.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P5 is active and P6-P8 are queued.
Its active archive is
`docs/plan/agent-sessions/2026-08-27T040108+0200-build-p5-http-harness.md`.
The `spring` archive is answered history and links reciprocally to the new
tail. The graph has exactly one NEXT archive. Only the launcher's mutable
header and prompt regions changed; the stable execution region is unchanged.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. The audit counts
ignored and untracked bytes, so ordinary and ignored status must both be empty
at a measured checkpoint.

## P5.6 Result

The sixth of eight mutation subjects, `spring`, is complete:

- The executable harness declares ten deterministic, unique mutations in
  `pkg/spring/io.go`: download and unzip argument ordering; download and unzip
  error gates; archive deletion selection; injected Unix-second ordering; root
  and dependency discovery URLs; empty-dependency validation; and valid
  dependency equality.
- The inventory seam label remains byte-exact:
  `5. Spring download keeps URL before archive path`.
- Each declaration binds production syntax occurring exactly once to a
  non-empty exact named `pkg/spring` test population. `pkg/spring/init.go`,
  `pkg/spring/types.go`, test files, fixtures, generated/vendor code, other
  subjects, unreached `UrlValuesFrom` and `DeleteDemoFiles` branches,
  equivalent replacements, and known non-compiling replacements are excluded.
- One clean unmodified external control ran every exact selection. Each mutant
  used a fresh external Git archive, cache, HOME, config, and temp root; the
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
`/private/tmp/ply-spring-evidence.x8UaSF/report.txt`; its SHA-256 is
`0aa2bdcb32cf20372f8a7b5c232dcefc93668259eb4ed76c45f8ca339b1dbf25`.
All retained control/mutant checkouts and logs are under
`/private/tmp/ply-spring-evidence.x8UaSF/harness-work`. The independent T1-T10
meta-log is `/private/tmp/ply-spring-evidence.x8UaSF/meta.log`, SHA-256
`f373190fa04abba73b61453b1693f892aced845d931ce267b3674ea84961fa06`.

## Authoritative Measurement

The refreshed external evidence and audits are under
`/private/tmp/ply-spring-evidence.x8UaSF`. The exact 245 Q1.6, 34 Q1.7, and 69
Q1.9 named populations were rerun at `9c20e52`; every named test emitted one
run and one pass. Manifest SHA-256 values are
`845cc6bdfaa0e4e84956a96a162e84ed0372238f067f8d2b22d05f7a22307f27`,
`cb028c078c2f61283a974ec308a3a0a2de1a6694cd2a60a130b6f4dd7f357d0d`,
and `c6488bd00c342304d8d1fa46a4af6dce33c9b095ec6ad1572bc9e0de40252034`.
Event-log SHA-256 values are
`60adda4a3acddcbd073bfd19c191f138920cefc9bc410a6a176eab1809843c98`,
`84035c3347e7fb7261c71f0d640a5086655f21c7163ee31879ed2311fadda1e6`,
and `0810add71fc180784bf4578711b842b1b362f1936d304125ac40a3f4e92190d9`.

The external schema-2 document contains refreshed Q1.6, Q1.7, Q1.9 receipts
and Q2.4 receipts for all six completed P5 subjects. Its SHA-256 is
`1eaabad8e769ad3bdc8eedfcf145aa521cb982da1d64697c06535df1e4b5ff14`.
Evidence-object SHA-256 values are:

- Q1.6: `35224a745edee899f8c1663f8560e0fbd9de52e47f4b2b61767ba0caca92a0e6`;
- Q1.7: `bfedca3f5dd2d6fec5907767233fbcb17046384cad226fd18de03e583f5f569c`;
- Q1.9: `f2228bca545937354fca3d8faa9c5ed632108ff3181f30bb5cdf2f8e1b4f9300`;
- Q2.4: `61731b25435e57c08c0c813d34d41cfce24644defabb13483fbc5cf88ee1194c`.

The focused Q1.6/Q1.7/Q1.9/Q2.4 audit exits 0 with scorecard SHA-256
`ab347e874db50c8df2682c30ad86f95f47c109bed71342dc131e90c3bb645a0c`.
The full authoritative audit exits 1, never 2, with scorecard SHA-256
`50ea67476a0e01cf7e31fa58e8a6c4590934d6691939b5c304e835a8af6aaf67`.
It records L0 8/8, all nine L1 rows PASS, Q2.1 improved to 6/8, Q2.4 PASS,
seven improved ratchets, one held, zero regressed, zero non-comparable, and
zero dirty paths. Eight P5-P8 rows remain non-passing.

The separate no-evidence Q2.1-Q2.4 view is under
`/private/tmp/ply-spring-evidence.x8UaSF/no-evidence-q2`, scorecard SHA-256
`a05b6401db516d8dc09e213b6545aa080a58fdb48ca710de29a2d1dff9e49d1f`.
It exposes the exact unratcheted state: Q2.1 is 6/8; Q2.2 and Q2.3 identify
only the legacy `interactive-build` recorder; Q2.4 is unmeasurable without the
external run receipt.

## Gate Result

Gate logs and external caches are under
`/private/tmp/ply-spring-gate.cmNq3k`. API/CLI and explicit entry/subprocess
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
`7afd548b287e549338276fc381cb747464403c66b9cff4135576dcd90b922123`,
`36b484c628424a75eac04fd1550952ef4b8ecea87123c15e2c6d2b6b02400f00`,
`ae853e4750310bf06d5a580067531fd2a08dd4c869efc04675c9b5fa10272baa`,
`70d9c7af646cfbc80dbe94a7193c02eb34ad0eb2f732578cf1473eb92eb2599f`,
and `a1147ffd4da8a0228684c4c6cd9c4363764aa4c38548b97150361c98d24ea85d`.

The first API compatibility invocation used the isolated HOME without binding
the existing module cache while the checker intentionally had `GOPROXY=off`;
the corrected invocation passed. One launcher contract run then missed the
known final signal-fixture partial-log timing assertion after its first 25
assertions. No source changed; the unchanged isolated rerun passed all 62
controls, and complete preflight passed the same suite. The passing rerun log
SHA-256 is
`bb0054078e20aca6a546ad399b0b55120d045c414523fa50cf36f1038d55fd48`.
Ordinary and ignored status were empty before authoritative measurement and
after the complete implementation gate.

## Next Objective

Create exactly the seventh inventory subject, `http`, whose production root is
`pkg/http` and whose declared but currently absent harness is
`scripts/mutate-http`. Add only that executable mutation harness and its
`scripts/test-mutate-http` T1-T10 falsifiability test. Define at least eight
meaningful mutations after inspecting the complete production and test
population. Keep `.quality/inventory` unchanged; unlike the seam-bearing
subjects, `http` has no separate inventory seam label.

Use all six completed P5 harnesses as methodology references. Require one clean
external control, fresh external copies and caches, exact test discovery,
successful mutant compilation, selected JSON run/terminal actions, and exact
totals `declared == killed`, `survived == 0`, `unusable == 0`. If a genuine
survivor appears, classify it as reachability, observability, or controllability
before the smallest in-subject test or private seam repair. Do not start the
eighth `interactive-build` subject. P5 remains active afterward.

## Start And Stop

Confirm branch, HEAD, exact ancestry, clean and ignored status, reciprocal
links, launcher `--check`, and the authorized checkpoint block before editing.
Read the rolling handover, linked NEXT archive, complete P5 entry/gate, both
designs, `.quality/README.md`, `.quality/inventory`, complete mutation audit and
parser logic, all completed harness/meta-test pairs, Make/preflight contracts,
the complete `pkg/http` production and test population, and the
`internal/adapter/httpclient` and `internal/adapter/filesystem` production and
test populations before selecting mutations.

Stop before another subject, inventory/audit/parser/scanner/baseline changes,
acceptance expansion, P6-P8, Go/dependency upgrades, exported API/CLI changes,
packaging, publication, or distribution. Keep generated artifacts external.
Do not push, merge, stash, revert, launch a successor, or remove the worktree.
