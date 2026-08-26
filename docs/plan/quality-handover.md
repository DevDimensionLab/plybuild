# Quality Upgrade Handover

Generated: 2026-08-26T23:33:25+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P5.1 implementation commit:
  `1dbc163b6716b96c3036ca95cadfd5a4a47c669d`.
- Its exact parent is the launch continuity commit
  `8bf2e9aa4953dbb6bda946d8dbbd645a946f7526`.
- Measured implementation tree:
  `23241502a50d38c0bdce711855d0507fa3d798ac`.
- Clean status SHA-256:
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
- After handoff, obtain the new continuity HEAD with `git rev-parse HEAD`; its
  exact parent must be `1dbc163`.
- The implementation adds only `scripts/mutate-cli-context` and
  `scripts/test-mutate-cli-context`. No production Go, Go test, inventory,
  audit, parser, scanner, baseline, acceptance, API, CLI, packaging, or
  dependency file changed.
- No push, merge, release, publication, distribution, stash, revert,
  successor launch, or worktree removal was performed.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P5 is active and P6-P8 are queued.
Its active archive is
`docs/plan/agent-sessions/2026-08-26T233325+0200-build-p5-config-cloud-harness.md`.
The `cli-context` archive is answered history and links reciprocally to the new
tail. The graph has exactly one NEXT archive. Only the launcher's mutable
header and prompt regions changed; the stable execution region is unchanged.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. The audit counts
ignored and untracked bytes, so ordinary and ignored status must both be empty
at a measured checkpoint.

## P5.1 Result

The first of eight mutation subjects, `cli-context`, is complete:

- The executable harness declares ten deterministic, unique mutations across
  the authorized `cmd` and `pkg/context` roots: CLI exit status, interactive
  build endpoint, interactive upgrade endpoint, recursive discovery branch,
  flattened-POM exclusion, target-directory exclusion, stealth-mode selection,
  per-project dry-run gate, root-project selection, and configured Maven
  repository selection.
- Each declaration binds one exact production replacement to a non-empty exact
  named test population. Every original syntax occurs exactly once. Test files,
  generated files, vendor, other subjects, equivalent replacements, and known
  non-compiling replacements are excluded.
- The harness requires a clean HEAD, creates one unmodified control and a fresh
  external Git archive/cache/HOME/config/temp root per mutant, verifies exact
  test discovery, compiles each mutant separately, and requires selected JSON
  run and terminal actions. A non-zero tool, setup, selection, or compile result
  is unusable, not a kill; only a selected test failure kills.
- Exact final totals are `declared=10`, `killed=10`, `survived=0`, and
  `unusable=0`. No reachability, observability, or controllability gap appeared,
  so no production or test repair was made.
- T1-T10 prove failure for empty/duplicate manifests, unauthorized paths,
  zero/multiple replacements, empty or inexact tests, broken control,
  uncompiled/unexercised mutants, false accounting, unclassified survivors,
  repository-local artifacts, and non-deterministic declarations/totals.

Every mutation and killing population is in the retained report
`/private/tmp/ply-cli-context-authoritative.gvdHKm.report.txt`; its SHA-256 is
`53d1489a6db15f18cd2acf38ff49ab3fb0db545a6c8e63639c90311f9bf53d5c`.
All retained mutant checkouts and logs are under
`/private/tmp/ply-cli-context-authoritative.gvdHKm`. The independent T1-T10
meta-log is `/private/tmp/ply-cli-context-meta.XXXXXX.log`, SHA-256
`f1d801fccc97af714d483f5a3c84202fe9868fa19b17038ada0d56dbd6f3e67e`.

## Authoritative Measurement

The refreshed external evidence and audits are under
`/private/tmp/ply-cli-context-evidence.Ij0EsY`. Current Go production and test
bytes are identical to P4 evidence commit `88a95ad`; nevertheless, the exact
245 Q1.6, 34 Q1.7, and 69 Q1.9 named populations were rerun at `1dbc163` and
all named tests emitted one run and one pass. Event-log SHA-256 values are
`369ef7a7acac7d46e5735e7fe8b1ba48383623f791f328bddc847e7d4d7c4a60`,
`426b26697215c46ac494acda287aba1abf06916b1e909fb093fd95d3100d54e9`,
and `8822c2e67a8c588616c81ca0c68379d9810a2d1df9cbb2565f4826aa9bb882c0`.

The external schema-2 document contains refreshed Q1.6, Q1.7, Q1.9 receipts
and one `cli-context` Q2.4 receipt. Its SHA-256 is
`0a6659c3b5a5234ca19134b4e2305de43a094fcb02032f0c6e106f5bbecaabbd`.
Evidence-object SHA-256 values are:

- Q1.6: `7d5a47c65720d218e5cf8bebb41564f06a9fd8280d8f5c419d0a30e896ef7c79`;
- Q1.7: `12878ba1b2a8cb0c5d6faa12f134c3fc5e40b1d99285444520ed3bbd2ec5b64a`;
- Q1.9: `45ab5d3ff070b6b35ac5354110f05f953350bda3cb509e1c1cf78fee08c37aea`;
- Q2.4: `70e2101ca2719d07022dd42e8db56b9ffc76c444473460907a216b9f6f4b124f`.

The focused Q1.6/Q1.7/Q1.9/Q2.4 audit exits 0 with scorecard SHA-256
`c693cc912f8c936c4438661aae48102ec09d05fe1ec4a78ad52f7c3d40abf7b8`.
The full authoritative audit exits 1, never 2, with scorecard SHA-256
`20528637cc01f651f6411484575bf3cec7f1a1ac6cccc9f1b30cdf9b989838ac`.
It records L0 8/8, all nine L1 rows PASS, Q2.1 improved from 0/8 to 1/8,
Q2.4 PASS, seven improved ratchets, one held, zero regressed, zero
non-comparable, and zero dirty paths. Eight P5-P8 rows remain non-passing.

The separate no-evidence Q2.1-Q2.4 view is under
`/private/tmp/ply-cli-context-q2-audit.h6auIu`, scorecard SHA-256
`7500813f38b080a06494b62d70a3def51ddd35bf75d6f804e382ddb2133c4545`.
It exposes the exact unratcheted state: Q2.1 is 1/8; Q2.2 and Q2.3 identify
only the remaining non-executable P3/P4 drivers; Q2.4 is unmeasurable without
the external run receipt.

## Gate Result

Gate logs and external caches are under
`/private/tmp/ply-cli-context-gate.RxKw8d`. API/CLI and CLI-surface
compatibility, compatibility meta-tests, pinned golangci-lint 2.12.2 with zero
issues, complete uncached tests across all 27 packages, race, vet, `make test`,
all 62 launcher controls, Make distribution/install/lint/preflight contracts,
all four host acceptance flows, exact empty-HOME count-2, the standalone
15-control audit meta-suite, and complete preflight pass. API and CLI report
SHA-256 values remain
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
Ordinary and ignored status were empty before authoritative measurement and
after the complete gate.

## Next Objective

Convert exactly the second inventory subject, `config-cloud`, whose production
root is `pkg/config` and whose declared harness is
`scripts/mutate-config-cloud`. Replace the existing non-executable P3 seam
driver and its `scripts/test-mutate-config-cloud` meta-test with one executable
mutation harness and T1-T10 falsifiability test. Define at least eight
meaningful config/cloud mutations after inspecting the complete production and
test populations. Preserve seam label
`3. cloud clone keeps URL before target directory`.

Use the completed `cli-context` harness as the methodology reference. Require
one clean external control, fresh external copies and caches, exact test
discovery, successful mutant compilation, selected JSON run/terminal actions,
and exact totals `declared == killed`, `survived == 0`, `unusable == 0`. If a
genuine survivor appears, classify it as reachability, observability, or
controllability before the smallest in-subject test or private seam repair.
Do not start a third subject. P5 remains active afterward.

## Start And Stop

Confirm branch, HEAD, exact ancestry, clean and ignored status, reciprocal
links, launcher `--check`, and the authorized checkpoint block before editing.
Read the rolling handover, linked NEXT archive, complete P5 entry/gate, both
designs, `.quality/README.md`, `.quality/inventory`, complete mutation audit and
parser logic, both completed and legacy harness/meta-test pairs relevant to the
move, Make/preflight contracts, and every production/test path under
`pkg/config` before selecting mutations.

Stop before another subject, inventory/audit/parser/scanner/baseline changes,
acceptance expansion, P6-P8, Go/dependency upgrades, exported API/CLI changes,
packaging, publication, or distribution. Keep generated artifacts external.
Do not push, merge, stash, revert, launch a successor, or remove the worktree.
