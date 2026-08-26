# Quality Upgrade Handover

Generated: 2026-08-27T00:11:40+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P5.2 implementation commit:
  `e5b4a26db271ac43d514bc0e9b2e19c23cdbc144`.
- Its exact parent is the launch continuity commit
  `be33bb9fa9e6d9831f824c831d9de113123a710b`.
- Measured implementation tree:
  `baa63ed3c4d44f063854a4fbee63728cf336c937`.
- Clean status SHA-256:
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
- After handoff, obtain the new continuity HEAD with `git rev-parse HEAD`; its
  exact parent must be `e5b4a26`.
- The implementation changes only `scripts/mutate-config-cloud` and
  `scripts/test-mutate-config-cloud`. No production Go, Go test, inventory,
  audit, parser, scanner, baseline, acceptance, API, CLI, packaging, or
  dependency file changed.
- No push, merge, release, publication, distribution, stash, revert,
  successor launch, or worktree removal was performed.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P5 is active and P6-P8 are queued.
Its active archive is
`docs/plan/agent-sessions/2026-08-27T001140+0200-build-p5-maven-sorting-harness.md`.
The `config-cloud` archive is answered history and links reciprocally to the
new tail. The graph has exactly one NEXT archive. Only the launcher's mutable
header and prompt regions changed; the stable execution region is unchanged.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. The audit counts
ignored and untracked bytes, so ordinary and ignored status must both be empty
at a measured checkpoint.

## P5.2 Result

The second of eight mutation subjects, `config-cloud`, is complete:

- The executable harness declares ten deterministic, unique mutations in
  `pkg/config/cloud.go`: clone argument ordering, the existing-repository
  branch, pull and clone error gates, Git-hook file classification, template
  config matching and relative naming, example-directory classification,
  default service-environment selection, and valid-template deduplication.
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
`/private/tmp/ply-config-cloud-authoritative.AL9ate/report.txt`; its SHA-256 is
`db76fdf8c624c4326483ae71fa9ec3e7e0f94de4d6d3ae8a4dff185c9a7f23f7`.
All retained mutant checkouts and logs are under
`/private/tmp/ply-config-cloud-authoritative.AL9ate/work`. The independent
T1-T10 meta-log is
`/private/tmp/ply-config-cloud-authoritative.AL9ate/meta.log`, SHA-256
`12b1197d521681b70dea0b481f6b7d8bb9daba2ca21dc68a635868177d43ffec`.

## Authoritative Measurement

The refreshed external evidence and audits are under
`/private/tmp/ply-config-cloud-evidence.uZBhkl`. Current Go production and test
bytes are unchanged by P5.2; nevertheless, the exact 245 Q1.6, 34 Q1.7, and
69 Q1.9 named populations were rerun at `e5b4a26` and
all named tests emitted one run and one pass. Event-log SHA-256 values are
`8f50bbefaa4e9b2d6f981cd21ec5215b08ab438fbc9867d3019b470e027b1e68`,
`c636f81d7d471347ac7be7ef255badfbc35b78c6bedc906e281de2f7b76b42a9`,
and `b468616dfd9751ff91fc2cfe5c71cf412ea08fe3c1af2e1c78fcd07d1428359f`.

The external schema-2 document contains refreshed Q1.6, Q1.7, Q1.9 receipts
and Q2.4 receipts for both completed P5 subjects. Its SHA-256 is
`7665b6d6d15da8dbd6b2da3e20a45902cb4e1076683fb0e1acd52c17b85c4648`.
Evidence-object SHA-256 values are:

- Q1.6: `fba28325058f5fbc728a45e1e9da7d8db6d004053e1709ad84b5178cc0cfd0a0`;
- Q1.7: `4400c6f5d0b8eaff61573fef6de94c4ebc2375fd97fc507702265949f7eaaf90`;
- Q1.9: `e582319bbaa6b2e629ae704915f967924b1b4764266585d0e7b8abc6da01d609`;
- Q2.4: `496ae9cd500b896080a14874dd5ff34bc0bdf521a401cbb7044e5f87604184b5`.

The focused Q1.6/Q1.7/Q1.9/Q2.4 audit exits 0 with scorecard SHA-256
`7ef785e2a17b5371358cc18741a78fe7cbca65c587725970970710d42a7a6b0c`.
The full authoritative audit exits 1, never 2, with scorecard SHA-256
`4a7089ea117a97bdf5b265f3712c534954417e00ff352b24a5874a2d710bf648`.
It records L0 8/8, all nine L1 rows PASS, Q2.1 improved from 1/8 to 2/8,
Q2.4 PASS, seven improved ratchets, one held, zero regressed, zero
non-comparable, and zero dirty paths. Eight P5-P8 rows remain non-passing.

The separate no-evidence Q2.1-Q2.4 view is under
`/private/tmp/ply-config-cloud-evidence.uZBhkl/no-evidence-q2`, scorecard
SHA-256
`5a17bd88fc7a48065883d72b1c6383b037fc2b47c18e47390918ba7cd782ed24`.
It exposes the exact unratcheted state: Q2.1 is 2/8; Q2.2 and Q2.3 identify
only the remaining non-executable P3/P4 drivers; Q2.4 is unmeasurable without
the external run receipt.

## Gate Result

Gate logs and external caches are under
`/private/tmp/ply-config-cloud-gate.1sHPne`. API/CLI and CLI-surface
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

Convert exactly the third inventory subject, `maven-sorting`, whose production
roots are `pkg/maven` and `pkg/sorting` and whose declared harness is
`scripts/mutate-maven-sorting`. Replace the existing non-executable P3 seam
driver and its `scripts/test-mutate-maven-sorting` meta-test with one executable
mutation harness and T1-T10 falsifiability test. Define at least eight
meaningful mutations after inspecting both complete production and test
populations. Preserve both seam labels exactly:
`2. Maven command keeps executable before arguments` and
`4. Maven metadata keeps username before password`.

Use both completed P5 harnesses as methodology references. Require
one clean external control, fresh external copies and caches, exact test
discovery, successful mutant compilation, selected JSON run/terminal actions,
and exact totals `declared == killed`, `survived == 0`, `unusable == 0`. If a
genuine survivor appears, classify it as reachability, observability, or
controllability before the smallest in-subject test or private seam repair.
Do not start a fourth subject. P5 remains active afterward.

## Start And Stop

Confirm branch, HEAD, exact ancestry, clean and ignored status, reciprocal
links, launcher `--check`, and the authorized checkpoint block before editing.
Read the rolling handover, linked NEXT archive, complete P5 entry/gate, both
designs, `.quality/README.md`, `.quality/inventory`, complete mutation audit and
parser logic, both completed and legacy harness/meta-test pairs relevant to the
move, Make/preflight contracts, and every production/test path under
`pkg/maven` and `pkg/sorting` before selecting mutations.

Stop before another subject, inventory/audit/parser/scanner/baseline changes,
acceptance expansion, P6-P8, Go/dependency upgrades, exported API/CLI changes,
packaging, publication, or distribution. Keep generated artifacts external.
Do not push, merge, stash, revert, launch a successor, or remove the worktree.
