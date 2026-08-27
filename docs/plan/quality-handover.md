# Quality Upgrade Handover

Generated: 2026-08-27T06:46:25+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P6 snapshot implementation commit:
  `06e3ac43282c1f4ccd3b99cec78d02e009a90645`.
- Its exact parent is the launch continuity commit
  `61763f4b5728fa77f9c947c66fcc04ba5136c208`; that commit's exact parent is
  the completed P5.8 implementation `58c5224`.
- Measured implementation tree:
  `0cb65473791c5eebf12f5aa6df9534be4a955a01`.
- Clean status SHA-256:
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
- After handoff, obtain the new continuity HEAD with `git rev-parse HEAD`; its
  exact parent must be `06e3ac4`.
- The implementation changes only `Makefile`, `scripts/accept-snapshot`,
  `scripts/test-accept-snapshot`,
  `test/acceptance/snapshot_acceptance.py`, and
  `test/makefile_distribution_test.sh`.
- No `.goreleaser.yml`, Docker, production Go, Go test, inventory, audit,
  parser, scanner, baseline, dependency, API, CLI, fixture, completed P5, or
  publisher file changed.
- No push, merge, release, publication, registry operation, Docker acceptance,
  `make quality`, stash, revert, successor launch, or worktree removal was
  performed.

## Continuity Checkpoint

P5 is complete and P6 remains active. The host-platform GoReleaser snapshot is
one finished P6 acceptance target; the fresh Docker image is the next separate
target. `codex-dev-start.sh` remains NEXT for that Docker-only move. The
snapshot archive is answered history and links reciprocally to the single new
Docker archive. The graph has exactly one NEXT archive. Only the launcher's
mutable header and prompt regions change; the stable execution region remains
byte-identical.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. The audit counts
ignored and untracked bytes, so ordinary and ignored status must both be empty
at a measured checkpoint.

## Snapshot Acceptance Result

The executable `scripts/accept-snapshot`, its executable
`scripts/test-accept-snapshot` meta-test, and the private
`test/acceptance/snapshot_acceptance.py` helper implement one fail-closed
acceptance path exposed as `make acceptance-snapshot`.

The orchestration requires:

- a clean committed source and a new external evidence root with no stale
  bytes;
- a regular external GoReleaser whose schema pin, `--version`, and Go build
  metadata all identify v2.17.1 for the host;
- the exact Make invocation `release --snapshot --clean --skip=publish` with
  all five release credentials cleared and publication disabled;
- fresh `artifacts.json` and `metadata.json` from the commanded run;
- exactly one Binary record for host `darwin/arm64`, with the path derived
  from metadata rather than guessed from a filename;
- no symlink path component, a regular executable, fresh timestamps, exact
  platform metadata, and embedded VCS revision equal to the source commit with
  `vcs.modified=false`;
- a separate host install into an initially missing external `GOBIN`, followed
  by `verify-install` with `PLY_VERIFY_ARTIFACT` unset; and
- exact status, upgrade, and build selection, one terminal PASS each, non-help
  runtime traces against the snapshot path, and stable before/after artifact
  identity.

The meta-test has 20 controls. It rejects no GoReleaser call; wrong or
publishing argv; leaked credentials; stale output; missing, ambiguous,
symlinked, non-executable, wrong-platform, misplaced, wrong-build-platform, or
foreign artifacts; skipped or duplicate verifier selection; artifact
substitution; verifier failure; duplicate terminal PASS; and repository-local
output.

An initial real target attempt rejected the system `go` command because that
command path was a symlink. No snapshot evidence root was created and the
source remained clean. The implementation now resolves command symlink chains
to their final regular executable before recording identity. The classified
failed log is retained at
`/private/tmp/ply-snapshot-evidence.jfKty7/acceptance-target.log`.

## Retained Snapshot Evidence

The clean tool probe is under `/private/tmp/ply-snapshot-probe.CiDxj0`.
The official `goreleaser_Darwin_arm64.tar.gz` v2.17.1 checksum is
`b65624885c25da9a677b7ad11cf86a02123cc5a56af66f6b4ebb574658eada2e`.
The retained GoReleaser executable is 84,264,062 bytes with SHA-256
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`.
The clean external probe found eight Binary records and exactly one host
candidate; no local-artifact defect required a configuration change.

The authoritative run is under
`/private/tmp/ply-snapshot-evidence.XnZgfB/run`. Its report records:

- source commit `06e3ac43282c1f4ccd3b99cec78d02e009a90645` and tree
  `0cb65473791c5eebf12f5aa6df9534be4a955a01`;
- exact command `release --snapshot --clean --skip=publish`, five cleared
  credentials, and disabled/skipped remote publication;
- exactly one host artifact at
  `source/dist/plybuild_darwin_arm64_v8.0/ply`, target
  `darwin_arm64_v8.0`, size 19,389,170, SHA-256
  `28d7012a81e9bdab2b776a8be64f84d12fc9e2f239e52dcbfc2932f7469d2e4f`;
- separate host-install artifact size 26,504,386, SHA-256
  `ca1d0ead122a8da336298a26d6dbc0574ef6bf6d684bb53dff595520241572fa`;
- `verify-install` PASS with the override unset; and
- status, upgrade, and build PASS, each with two non-help behavioral traces
  against the exact snapshot executable.

The report SHA-256 is
`6e759f75ae40c27c6312cfc4f4c1e8a8c4eab56119a64a88fa4c4c62dec40f61`.
The evidence-manifest SHA-256 is
`bb33ae40ccab0992284318a123c9462fa86ae6a1e5dd2c7b881d0b302c4ab9a1`.
The outer target log SHA-256 is
`54721148f142db090a680bd25d2e87d449271a46e393ef8fcd2667b3a51d32e1`.
The metadata SHA-256 is
`91b1f025906cb79a1637838bcef25bd072469ab454cbfc7952ad547033640e68`.

Verifier trace SHA-256 values are:

- status:
  `a27905a76a9c46cfb2259cc859a7f9d6307e9af727104753e1cff8e755eb2b63`;
- upgrade:
  `532b3e047fafb077b33e461b74b447e7a4dd57e6a03be23ced19928bb9eacc3a`;
- build:
  `1abb49da0f902d6f42cf4a251ecf3fef8ee116606f55f8d4a16e6560b6a42550`.

## Audit Measurement

The external gate and audit root is
`/private/tmp/ply-snapshot-gate.06e3ac4`; its evidence-manifest SHA-256 is
`61d453d3edf60fc582523472b02ca8d89449ced1dd04c54273e19b44de309255`.
The refreshed external schema-2 document is bound to the clean implementation
commit and has SHA-256
`7a09ec592025564f6600b3edf21fd0c2d56dc28fd73a92de2b3bb80024f32e70`.
Its Q1.6, Q1.7, Q1.9, and Q2.4 evidence-object SHA-256 values are:

- `e504d3a4446654b1d403a8c72400693f2f0be9f34ca02497ba4f04e9a03df8e7`;
- `07aaa86eeebd1b6ed8fcc5de69040e12b72c4d808589d8c93602a09ca20e5e18`;
- `e2a5c56643a2b018e6e6e1b75bde168eb9077af5ed67090aa7a4ce33a9643226`;
- `29a4f2387fa72461fd966d0022ba5a29018d3b7c1ab79191efd11ae43c1576fd`.

The focused receipt audit exits 0 with scorecard SHA-256
`2b8d42fa6457ab02bd6dbed85d3fd13ccd0552cae794d8cc7c6d70644e09cbae`.
The full authoritative audit exits 1, never 2, with scorecard SHA-256
`9b75a367fc07ddc698787103f80360227e6f03827009c62c4a80a59dcd547b4d`.
It records L0 8/8, L1 9/9, L2 8/10, seven improved ratchets, one Q3.4
documentation-phrase regression from immutable prompt history, zero dirty
paths, and the same seven non-passing rows: Q2.8, Q2.9, Q3.1, Q3.3, Q3.4,
Q3.7, and Q3.8.

The focused Q2.5-Q2.10 view exits 1, never 2, with scorecard SHA-256
`367819708c9959d6b6bd6085e7a362124c95a2f9bb744db48a7294a1ffa0d482`.
Q2.5, Q2.6, Q2.7, and Q2.10 pass; Q2.8 and Q2.9 remain honestly manual.
Snapshot orchestration is deliberately not named `verify-*`, so it does not
silently change the four-core-flow audit denominator.

## Gate Result

API/CLI and explicit entry/subprocess compatibility, compatibility meta-tests,
pinned golangci-lint 2.12.2 with zero issues, complete uncached tests across all
27 packages, race, vet, `make test`, all 62 launcher controls, Make
distribution/install/lint/preflight contracts, all four host acceptance flows,
the 20-control snapshot meta-test, exact empty-HOME count-2, the standalone
15-control audit meta-suite, and complete preflight pass.

API and CLI report SHA-256 values remain
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
Passing log SHA-256 values include:

- complete preflight:
  `21b41a0a87e07e2eb12b4c469e997a1451ede29972de9d7db373b50e3a774e93`;
- host acceptance:
  `36b484c628424a75eac04fd1550952ef4b8ecea87123c15e2c6d2b6b02400f00`;
- standalone audit meta-suite:
  `ae853e4750310bf06d5a580067531fd2a08dd4c869efc04675c9b5fa10272baa`;
- empty-HOME count-2:
  `d42090481f02f3e17f5c0a4a7eee7f565339e6196110e9831f4056d875a7e1fb`;
- entry/subprocess:
  `fac2f31951e4c4c5e1ad41b80925a4b7c0ddb2bc7bcd951bea9266a7f0e4a096`;
- `make test`:
  `9beeeb53b603276b2e6e6944d726adccc8c1f4d7e0fd0c038cd85e90a61d46d4`;
- snapshot meta-test:
  `9804c5b0d54f5ada0dcfbd1ee401559dc026c8aa93778b2c604afb64adffa9b9`.

The first isolated compatibility attempt used an empty HOME without the
existing external module cache, so the offline API exporter failed before
comparison; the retained unchanged rerun with explicit `GOMODCACHE` passed.
The first complete preflight hit the known launcher signal-retention timing
control after 25 nested assertions. No source changed; the standalone launcher
suite had already passed all 62 controls, and the unchanged isolated complete
preflight rerun passed. Ordinary and ignored status were empty before
authoritative measurement and after the implementation gate.

## Next Objective

Continue P6 with one bounded daemon-backed Docker acceptance move. Build a
fresh local image from clean committed source with a unique non-publishing
identity, inspect and record the exact image ID/platform/configuration and the
regular executable inside it, revalidate host install as the separate
`make install` / `go install` contract, and run the existing status, upgrade,
and build behaviors through the exact image against local fixtures. Require
runtime proof, not static Dockerfile inspection or a help-only smoke.

Add fail-closed controls for missing Docker calls, stale or substituted images,
wrong build/run or publishing operations, ambiguous identities, wrong
platform/configuration, missing/non-executable image artifacts, skipped or
duplicate verifiers, verifier failure, and repository-local evidence. Use a
real daemon probe before choosing the smallest orchestration. Change
`Dockerfile` only if that executed clean probe exposes a classified local
acceptance defect.

Do not reimplement snapshot acceptance and do not add `make quality` in the
Docker move. P6 remains active after Docker evidence until the separate final
quality-gate move resolves Q2.8/Q2.9 evidence and proves the scoped audit exit.

## Start And Stop

Confirm branch, HEAD, exact ancestry, clean and ignored status, reciprocal
links, launcher `--check`, and the authorized checkpoint block before editing.
Read this handover, the linked NEXT archive, the complete P6 and checkpoint
gate, both design documents, the full snapshot implementation and evidence
schema, Dockerfile/Make/publish boundaries and their tests, every acceptance
verifier/meta-test/helper/fixture, and Q2.5-Q2.10 audit/parser logic before
defining the image population.

Stop before `make quality`, P7-P8, Go/dependency upgrades, production/API/CLI
behavior, remote publication, registry pushes, package-manager publishers, or
unrelated audit/inventory/P5/snapshot changes. Keep generated contexts, logs,
reports, and caches external. Do not push, merge, stash, revert, launch a
successor, or remove the worktree.
