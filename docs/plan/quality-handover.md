# Quality Upgrade Handover

Generated: 2026-08-27T03:08:46+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P5.5 implementation commit:
  `7bb94acf47e5971bf331f09e3e2e17ff94e137f6`.
- Its exact parent is the launch continuity commit
  `f575f315bdee668c62ed217f67f4d0cf0a455c1f`.
- Measured implementation tree:
  `00a6db27697d337df9cb6e1bba59f3f791fea923`.
- Clean status SHA-256:
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
- After handoff, obtain the new continuity HEAD with `git rev-parse HEAD`; its
  exact parent must be `7bb94ac`.
- The implementation changes only `scripts/mutate-file-shell` and
  `scripts/test-mutate-file-shell`. No production Go, Go test, inventory,
  audit, parser, scanner, baseline, acceptance, API, CLI, packaging, or
  dependency file changed.
- No push, merge, release, publication, distribution, stash, revert,
  successor launch, or worktree removal was performed.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P5 is active and P6-P8 are queued.
Its active archive is
`docs/plan/agent-sessions/2026-08-27T030846+0200-build-p5-spring-harness.md`.
The `file-shell` archive is answered history and links reciprocally to the new
tail. The graph has exactly one NEXT archive. Only the launcher's mutable
header and prompt regions changed; the stable execution region is unchanged.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. The audit counts
ignored and untracked bytes, so ordinary and ignored status must both be empty
at a measured checkpoint.

## P5.5 Result

The fifth of eight mutation subjects, `file-shell`, is complete:

- The executable harness declares ten deterministic, unique mutations across
  `pkg/file/file.go`, `pkg/shell/command.go`, and `pkg/shell/git.go`: clone,
  pull, commit, command, unzip, and move argument ordering; the Git-add error
  gate; first-match suffix selection; all-match exclusion selection; and
  clear-directory exclusion gating.
- Both inventory seam labels remain byte-exact:
  `1. git clone keeps URL before target directory` and
  `7. git commit keeps target directory before message`.
- Each declaration binds production syntax occurring exactly once to a
  non-empty exact named `pkg/file` or `pkg/shell` test population. Other
  production paths, test files, fixtures, generated/vendor code, other
  subjects, equivalent replacements, and known non-compiling replacements are
  excluded.
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
`/private/tmp/ply-file-shell-evidence.y2RSai/report.txt`; its SHA-256 is
`ea12b122e6210f2b2c8c64a7792ce74a9886e2273e9ddac6176df6552afd1c29`.
All retained control/mutant checkouts and logs are under
`/private/tmp/ply-file-shell-evidence.y2RSai/harness-work`. The independent
T1-T10 meta-log is `/private/tmp/ply-file-shell-evidence.y2RSai/meta.log`,
SHA-256
`759ad7a28b3d3ab06d42d42e8e1679e0df516d6abdfd90f3fbac1e92e52bfaad`.

## Authoritative Measurement

The refreshed external evidence and audits are under
`/private/tmp/ply-file-shell-evidence.y2RSai`. The exact 245 Q1.6, 34 Q1.7,
and 69 Q1.9 named populations were rerun at `7bb94ac`; every named test emitted
one run and one pass. Manifest SHA-256 values are
`845cc6bdfaa0e4e84956a96a162e84ed0372238f067f8d2b22d05f7a22307f27`,
`cb028c078c2f61283a974ec308a3a0a2de1a6694cd2a60a130b6f4dd7f357d0d`,
and `c6488bd00c342304d8d1fa46a4af6dce33c9b095ec6ad1572bc9e0de40252034`.
Event-log SHA-256 values are
`dc4416f59295f5eb7d1afe5d8d19347a3cea9ca85c0d11bf79bcb6ff24b15c2f`,
`930baf6259e1a68086d9b20ade1e4149693fb4fa9d5a5c1737663ab7f562cab4`,
and `640cf82d6bf13149ee4015ebfe72987554a87820cacc83251fb5e54276e1f88a`.

The external schema-2 document contains refreshed Q1.6, Q1.7, Q1.9 receipts
and Q2.4 receipts for all five completed P5 subjects. Its SHA-256 is
`c5ec79633274ad618b7c78568e38891b9955352e020053c8b3d987ad6348fd5e`.
Evidence-object SHA-256 values are:

- Q1.6: `da6bf3b4016a872a8e49f293d6c206b711611162c8bc720467bc9f17a838edfe`;
- Q1.7: `f2a479eb893092f7a13133659496ac083bae844348540ef2076780a182f85334`;
- Q1.9: `610079011a185004fe8f2056e1cfe8ac456e0e3f2e8a407d14b1c879dcdc9194`;
- Q2.4: `732083b2add74a30bf5bed70ef9446682bbba20ebb4d1dfc89d0b904db81b922`.

The focused Q1.6/Q1.7/Q1.9/Q2.4 audit exits 0 with scorecard SHA-256
`5261b60d5e20338c5f1a88870b60c55e54411200a7fa84a0d80791dd96597001`.
The full authoritative audit exits 1, never 2, with scorecard SHA-256
`7a42d1f40ec96641fce8dab1188a32a73a2f6f790cedea60e98e5b05b1833502`.
It records L0 8/8, all nine L1 rows PASS, Q2.1 improved to 5/8, Q2.4 PASS,
seven improved ratchets, one held, zero regressed, zero non-comparable, and
zero dirty paths. Eight P5-P8 rows remain non-passing.

The separate no-evidence Q2.1-Q2.4 view is under
`/private/tmp/ply-file-shell-evidence.y2RSai/no-evidence-q2`, scorecard SHA-256
`fcbb44afd1ec5a92ce8d4cd6dd5b8b6a12de363fbc96b33dab07e345c87a4338`.
It exposes the exact unratcheted state: Q2.1 is 5/8; Q2.2 and Q2.3 identify
only `spring`, `http`, and `interactive-build`; Q2.4 is unmeasurable without
the external run receipt.

## Gate Result

Gate logs and external caches are under
`/private/tmp/ply-file-shell-gate.maocfI`. API/CLI and explicit
entry/subprocess compatibility, compatibility meta-tests, pinned
golangci-lint 2.12.2 with zero issues, complete uncached tests across all 27
packages, race, vet, `make test`, all 62 launcher controls, Make
distribution/install/lint/preflight contracts, all four host acceptance flows,
exact empty-HOME count-2, the standalone 15-control audit meta-suite, and
complete preflight pass. API and CLI report SHA-256 values remain
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
The complete preflight, acceptance, audit-meta, empty-HOME count-2, and focused
entry/subprocess log SHA-256 values are respectively
`c6bb96c2c5ee506a0dff170c7465f69ed39a647073be9414d87e3944fdedc350`,
`36b484c628424a75eac04fd1550952ef4b8ecea87123c15e2c6d2b6b02400f00`,
`ae853e4750310bf06d5a580067531fd2a08dd4c869efc04675c9b5fa10272baa`,
`e67e7862f1d7442c49235decba4c5be419def5f72e8717cd36a22820903e62c9`,
and `835051c0605b2c60d6f8cc7c5eb0244cfd4d310c514fc098486603b49b94e0f6`.
The first final handoff-contract invocation hit the known nested signal-fixture
partial-log timing assertion after its first 50 controls. No source changed;
the unchanged suite passed all 62 controls in a fresh external temp root. That
passing handoff log SHA-256 is
`98e969fd470c7d39c42b9eac0cb7d138c86fed9306acae3286e1168e2e62e841`.
Ordinary and ignored status were empty before authoritative measurement and
after the complete implementation gate.

## Next Objective

Convert exactly the sixth inventory subject, `spring`, whose production root
is `pkg/spring` and whose declared harness is `scripts/mutate-spring`. Replace
the existing non-executable P3 seam driver and its
`scripts/test-mutate-spring` meta-test with one executable mutation harness and
T1-T10 falsifiability test. Define at least eight meaningful mutations after
inspecting the complete production and test population. Preserve the seam
label exactly: `5. Spring download keeps URL before archive path`.

Use all five completed P5 harnesses as methodology references. Require one
clean external control, fresh external copies and caches, exact test discovery,
successful mutant compilation, selected JSON run/terminal actions, and exact
totals `declared == killed`, `survived == 0`, `unusable == 0`. If a genuine
survivor appears, classify it as reachability, observability, or controllability
before the smallest in-subject test or private seam repair. Do not start a
seventh subject. P5 remains active afterward.

## Start And Stop

Confirm branch, HEAD, exact ancestry, clean and ignored status, reciprocal
links, launcher `--check`, and the authorized checkpoint block before editing.
Read the rolling handover, linked NEXT archive, complete P5 entry/gate, both
designs, `.quality/README.md`, `.quality/inventory`, complete mutation audit and
parser logic, all completed and legacy harness/meta-test pairs relevant to the
move, Make/preflight contracts, the HTTP client and filesystem adapters used by
the Spring seam, and every production/test path under `pkg/spring` before
selecting mutations.

Stop before another subject, inventory/audit/parser/scanner/baseline changes,
acceptance expansion, P6-P8, Go/dependency upgrades, exported API/CLI changes,
packaging, publication, or distribution. Keep generated artifacts external.
Do not push, merge, stash, revert, launch a successor, or remove the worktree.
