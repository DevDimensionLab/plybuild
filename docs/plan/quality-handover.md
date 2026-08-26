# Quality Upgrade Handover

Generated: 2026-08-27T01:00:33+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P5.3 implementation commit:
  `1a40e1146f27619d1b46e53f2e24b155760bcc7c`.
- Its exact parent is the launch continuity commit
  `e0f9a4efe38349caf97d61f335cace017d52b4f3`.
- Measured implementation tree:
  `44ab7771e8e9c4ea3f9a2c44b248db4500567d34`.
- Clean status SHA-256:
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
- After handoff, obtain the new continuity HEAD with `git rev-parse HEAD`; its
  exact parent must be `1a40e11`.
- The implementation changes only `scripts/mutate-maven-sorting` and
  `scripts/test-mutate-maven-sorting`. No production Go, Go test, inventory,
  audit, parser, scanner, baseline, acceptance, API, CLI, packaging, or
  dependency file changed.
- No push, merge, release, publication, distribution, stash, revert,
  successor launch, or worktree removal was performed.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P5 is active and P6-P8 are queued.
Its active archive is
`docs/plan/agent-sessions/2026-08-27T010033+0200-build-p5-template-harness.md`.
The `maven-sorting` archive is answered history and links reciprocally to the
new tail. The graph has exactly one NEXT archive. Only the launcher's mutable
header and prompt regions changed; the stable execution region is unchanged.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. The audit counts
ignored and untracked bytes, so ordinary and ignored status must both be empty
at a measured checkpoint.

## P5.3 Result

The third of eight mutation subjects, `maven-sorting`, is complete:

- The executable harness declares ten deterministic, unique production
  mutations across `pkg/maven/command.go`, `pkg/maven/query.go`, and
  `pkg/sorting/sort.go`: Maven executable/argument ordering; metadata
  username/password ordering; coordinate path ordering; authentication and
  metadata-error gates; dependency swap and comparison directions; sort-key
  field ordering; group-match gating; and compile-scope priority.
- Both inventory seam labels remain byte-exact:
  `2. Maven command keeps executable before arguments` and
  `4. Maven metadata keeps username before password`.
- Each declaration binds syntax occurring exactly once to a non-empty exact
  named test population. Test files, vendor/generated code, other subjects,
  equivalent replacements, and known non-compiling replacements are excluded.
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
`/private/tmp/ply-maven-sorting-authoritative.kZBNCd/report.txt`; its SHA-256 is
`6b360f604802c047b4946b474c35b2860d47bda4581b3cb33f5af45653dc111e`.
All retained mutant checkouts and logs are under
`/private/tmp/ply-maven-sorting-authoritative.kZBNCd/work`. The independent
T1-T10 meta-log is
`/private/tmp/ply-maven-sorting-authoritative.kZBNCd/meta.log`, SHA-256
`fda3185a71ebd842f3e924ae507ad0daa9d59d4b30da3699d114c0ff0d2d53b3`.

## Authoritative Measurement

The refreshed external evidence and audits are under
`/private/tmp/ply-maven-sorting-evidence.bMq3EJ`. Current Go production and test
bytes are identical to the P4 evidence commit. The exact 245 Q1.6, 34 Q1.7,
and 69 Q1.9 named populations were nevertheless rerun at `1a40e11`; every
named test emitted one run and one pass. Event-log SHA-256 values are
`177ca7de050d567bf25fd51dea9920288ede89c66925c17b1ca44b67335dd929`,
`0bf842d18f2ab4ada5ceb518ddc618800d1a40a86a156086618c75d44b8fd46a`,
and `9d401949937f63087d10a6e8183ad5176b6b2bdefccba768895d1bf9c3139188`.

The external schema-2 document contains refreshed Q1.6, Q1.7, Q1.9 receipts
and Q2.4 receipts for all three completed P5 subjects. Its SHA-256 is
`a279e8ac3402cacad23eb757bc3adbf7e9e1e105115135f7183ebf5353ff1bb4`.
Evidence-object SHA-256 values are:

- Q1.6: `1be89cac435fc4183ef0e748a18ffc4510d86a0d976e0fb064b376f8554045c9`;
- Q1.7: `fa8129fc0c3edbc32b808af8408c277b5f93e03c285a5820a9868f6b7596dba5`;
- Q1.9: `82c1bd6452a45b6f302404df9b61ca53180304a439636421cce9f03c1af79ce9`;
- Q2.4: `53d27607c73ae97ea00dcfce384a373c5c421fcd1e900433521fbfbb1e4bf689`.

The focused Q1.6/Q1.7/Q1.9/Q2.4 audit exits 0 with scorecard SHA-256
`07d9f680bba06cb88c589fe87962feaf74ca6add94fda8b933679d6a98bc313c`.
The full authoritative audit exits 1, never 2, with scorecard SHA-256
`122deef2c884c1e3f6edb3052b20d47c0fad82ba5a4dc44f0b3696afe947d9f4`.
It records L0 8/8, all nine L1 rows PASS, Q2.1 improved to 3/8, Q2.4 PASS,
seven improved ratchets, one held, zero regressed, zero non-comparable, and
zero dirty paths. Eight P5-P8 rows remain non-passing.

The separate no-evidence Q2.1-Q2.4 view is under
`/private/tmp/ply-maven-sorting-evidence.bMq3EJ/no-evidence-q2`, scorecard
SHA-256
`f6e6c4ecaf67549fde6c913f848a1630c79237659a000c540b308b5f8fb8e724`.
It exposes the exact unratcheted state: Q2.1 is 3/8; Q2.2 and Q2.3 identify
only the remaining non-executable P3/P4 drivers; Q2.4 is unmeasurable without
the external run receipt.

## Gate Result

Gate logs and external caches are under
`/private/tmp/ply-maven-sorting-gate.JPO0hA`. API/CLI and CLI-surface
compatibility, compatibility meta-tests, pinned golangci-lint 2.12.2 with zero
issues, complete uncached tests across all 27 packages, race, vet, `make test`,
all 62 launcher controls, Make distribution/install/lint/preflight contracts,
all four host acceptance flows, exact empty-HOME count-2, the standalone
15-control audit meta-suite, and complete preflight pass. API and CLI report
SHA-256 values remain
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
The complete preflight log SHA-256 is
`34a09c1a7c0d57ceca6cc8832492c77d638d5b59504d148da7615673c51dbd16`.

A few PTY or output-piped launcher runs missed the known signal-fixture
partial-log timing assertion. The unchanged launcher suite passed all 62
controls in the direct non-TTY handoff rerun, full non-TTY `make test`, and
complete preflight. Ordinary and ignored status were empty before authoritative
measurement and after the complete gate.

## Next Objective

Convert exactly the fourth inventory subject, `template`, whose production
root is `pkg/template` and whose declared harness is `scripts/mutate-template`.
Replace the existing non-executable P3 seam driver and its
`scripts/test-mutate-template` meta-test with one executable mutation harness
and T1-T10 falsifiability test. Define at least eight meaningful mutations
after inspecting the complete production and test population. Preserve the
seam label exactly: `6. template copy keeps source before destination`.

Use all three completed P5 harnesses as methodology references. Require one
clean external control, fresh external copies and caches, exact test discovery,
successful mutant compilation, selected JSON run/terminal actions, and exact
totals `declared == killed`, `survived == 0`, `unusable == 0`. If a genuine
survivor appears, classify it as reachability, observability, or controllability
before the smallest in-subject test or private seam repair. Do not start a fifth
subject. P5 remains active afterward.

## Start And Stop

Confirm branch, HEAD, exact ancestry, clean and ignored status, reciprocal
links, launcher `--check`, and the authorized checkpoint block before editing.
Read the rolling handover, linked NEXT archive, complete P5 entry/gate, both
designs, `.quality/README.md`, `.quality/inventory`, complete mutation audit and
parser logic, all completed and legacy harness/meta-test pairs relevant to the
move, Make/preflight contracts, and every production/test path under
`pkg/template` before selecting mutations.

Stop before another subject, inventory/audit/parser/scanner/baseline changes,
acceptance expansion, P6-P8, Go/dependency upgrades, exported API/CLI changes,
packaging, publication, or distribution. Keep generated artifacts external.
Do not push, merge, stash, revert, launch a successor, or remove the worktree.
