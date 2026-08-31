# Quality Upgrade Handover

Generated: 2026-08-31T09:05:07+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Launch continuity commit:
  `20dbf71a17f2108a6a66c99af9226ce2fbf6a92c`.
- Its exact parent is the second Docker-blocker continuity commit
  `ddea23b1dbdf5a3277abcf0152d1dd187048c83e`.
- After handoff, obtain the new continuity HEAD with `git rev-parse HEAD`; its
  exact parent must be `20dbf71`.
- The snapshot implementation tree is
  `0cb65473791c5eebf12f5aa6df9534be4a955a01`.
- Clean status SHA-256 is
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
- All three Docker attempts made no implementation, test, Makefile, Dockerfile, or
  distribution change because no daemon-backed build could be started.
- No Docker build or run, login, push, publication, GoReleaser run, ordinary
  release, `make quality`, stash, revert, successor launch, or worktree removal
  was performed.

## Continuity Checkpoint

P5 is complete and P6 remains active. The host-platform GoReleaser snapshot is
a finished P6 acceptance target. Docker acceptance remains wholly unfinished
because the installed Docker client could not reach or start a daemon from the
managed workspace sandbox. P7-P8 remain queued and `make quality` remains
deferred.

The snapshot archive, first Docker archive, and two resume archives are
answered history. The second resume archive links reciprocally to one new
resume-only Docker archive, which is the sole NEXT tail. Only the launcher's
mutable header and prompt regions change; its stable execution region remains
byte-identical.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. The audit counts
ignored and untracked bytes, so ordinary and ignored status must both be empty
at a measured checkpoint.

## Docker Daemon Blocker

The clean launch state was verified before probing:

- branch `codex/upgrade-quality` at `20dbf71`, exact parent `ddea23b`;
- no ordinary, untracked, or ignored status entries;
- reciprocal archive graph and authorized P6 checkpoint valid; and
- `./codex-dev-start.sh --check` PASS for the Docker session.

The rolling handover, all reciprocal Docker archives, complete P6 and
checkpoint gate, and launcher contracts were rechecked before the fresh probe.
The daemon gate failed before any Dockerfile or implementation input was read,
so no design or image population was defined.

The retained blocker evidence is
`/private/tmp/ply-docker-blocker.KRsWoi`. Its evidence-manifest SHA-256 is
`991d7bf4a7a739b9d971da7776a8be7e447d2aad6c931a7cf72c1bbdaeb682ae`.
It records:

- Docker client 29.4.0, API 1.54, Go 1.26.1, commit `9d7ad9f`,
  `darwin/arm64`;
- buildx `v0.33.0-desktop.1`, commit
  `7f91f038ac14cbf5c4b2a6b76470860814424da1`;
- selected context `desktop-linux`, endpoint
  `unix:///Users/perottochristensen/.docker/run/docker.sock`;
- host macOS 15.3.1 build 24D70 on Darwin 24.3.0 arm64;
- Docker Desktop app metadata version 4.70.0 build 224270; and
- `docker version` exit 1, `docker info` exit 1, and `docker desktop status`
  exit 1 because no server answered.

The installed app could not be started from this environment. `open` exited 1
with LaunchServices `kLSNoExecutableErr`; the bundled backend exists and is
signed, but both it and `docker desktop status/start` must create logs under
`~/Library/Containers/com.docker.docker`, which is outside the session's
writable roots and failed with `operation not permitted`. Apple Events were
unavailable and Computer Use approval for `com.docker.docker` was denied. No
alternate Colima, OrbStack, or Podman executable is installed. The daemon must
therefore be started by the user or another process outside this sandbox before
the next session.

The latest repeated probe is retained at
`/private/tmp/ply-docker-probe.ffXBUb`. Its 49-file evidence manifest verifies,
and the manifest SHA-256 is
`16ed6fed1d517f8585815032363e2ba8b2a57fa6f1e811ad4bde4da892ec2a52`.
It records:

- Docker client 29.4.0 and buildx `v0.33.0-desktop.1`;
- selected context `desktop-linux` and the same Unix socket endpoint;
- macOS 15.3.1 build 24D70 on Darwin 24.3.0 arm64;
- `docker version` exit 1 with JSON `Server: null`, `docker info` exit 1,
  and `docker desktop status` exit 1;
- a named `desktop-linux` builder with no driver or supported-platform
  identity because buildx could not reach the daemon; and
- the same sandbox denial while Docker Desktop tried to open its host log.

No Docker build or run, login, push, publication, Dockerfile read, design
choice, or repository implementation change occurred in the latest attempt.

This remains an external execution blocker, not a Dockerfile defect. A static
read cannot classify the current Dockerfile, and no orchestration shape was
selected or implemented. The successor must first require `docker version` to
contain a real Server identity and `docker info` to exit 0; it must stop before
editing if those conditions do not hold.

## Retained Snapshot Evidence

The authoritative snapshot run remains under
`/private/tmp/ply-snapshot-evidence.XnZgfB/run`. It records clean source commit
`06e3ac4`, exact credential-cleared argv
`release --snapshot --clean --skip=publish`, and exactly one `darwin/arm64`
artifact at `source/dist/plybuild_darwin_arm64_v8.0/ply`, size 19,389,170,
SHA-256
`28d7012a81e9bdab2b776a8be64f84d12fc9e2f239e52dcbfc2932f7469d2e4f`.

The report SHA-256 is
`6e759f75ae40c27c6312cfc4f4c1e8a8c4eab56119a64a88fa4c4c62dec40f61`
and the evidence-manifest SHA-256 is
`bb33ae40ccab0992284318a123c9462fa86ae6a1e5dd2c7b881d0b302c4ab9a1`.
Host install passed separately with the artifact override unset. Status,
upgrade, and build each emitted exactly one terminal PASS and two non-help
behavioral traces against that snapshot. All 20 snapshot meta-controls pass.

## Prior Gate And Audit Measurement

The last complete implementation gate is the clean `06e3ac4` snapshot gate.
API/CLI and entry/subprocess compatibility, pinned golangci-lint 2.12.2,
complete uncached tests across 27 packages, race, vet, `make test`, all 62
launcher controls, Make contracts, all four host acceptance flows, snapshot
meta/acceptance, complete preflight, the 15-control audit meta-suite, and
empty-HOME count-2 passed.

The external schema-2 document SHA-256 is
`7a09ec592025564f6600b3edf21fd0c2d56dc28fd73a92de2b3bb80024f32e70`.
The exact Q1.6/Q1.7/Q1.9/Q2.4 audit exits 0 with scorecard SHA-256
`2b8d42fa6457ab02bd6dbed85d3fd13ccd0552cae794d8cc7c6d70644e09cbae`.
The full audit exits 1, never 2, with scorecard SHA-256
`9b75a367fc07ddc698787103f80360227e6f03827009c62c4a80a59dcd547b4d`:
L0 is 8/8, L1 is 9/9, L2 is 8/10, seven ratchets improve, Q3.4 records one
prompt-history regression, dirty paths are empty, and seven rows remain
non-passing.

The focused Q2.5-Q2.10 view exits 1, never 2, with scorecard SHA-256
`367819708c9959d6b6bd6085e7a362124c95a2f9bb744db48a7294a1ffa0d482`.
Q2.5, Q2.6, Q2.7, and Q2.10 pass; Q2.8 and Q2.9 remain honestly manual.
Docker orchestration must remain outside the `verify-*` audit denominator.

## Next Objective

Before starting the next Codex session, start Docker Desktop outside the managed
workspace sandbox and verify that `docker version` includes both Client and
Server records and that `docker info` exits 0 for `desktop-linux`.

Then resume the same bounded P6 Docker target. Run a clean external daemon
probe before choosing the smallest design. Build a fresh local image from clean
committed source with a unique non-publishing tag; bind it to one immutable
image ID; inspect platform, entrypoint, configuration, creation metadata, and
regular executable provenance; revalidate host install separately; and run the
existing status, upgrade, and build behaviors through the immutable image with
runtime receipts. Add the required fail-closed meta-controls and retain all
evidence externally.

Change `Dockerfile` only if the executed clean build exposes a classified local
acceptance defect. Do not add `make quality` during Docker acceptance. After
successful Docker evidence and its full gate, hand off the separate final P6
schema-2 evidence and scoped `make quality` exit-gate move.

## Start And Stop

Confirm branch, HEAD, exact ancestry, clean and ignored status, reciprocal
links, launcher `--check`, and the authorized checkpoint block before editing.
Recheck the retained blocker facts, then require a live daemon before choosing
or implementing Docker acceptance.

Stop before `make quality`, P7-P8, Go/dependency upgrades, production/API/CLI
behavior, remote publication, registry pushes, package-manager publishers, or
unrelated audit/inventory/P5/snapshot changes. Keep generated contexts, logs,
reports, and caches external. Do not push, merge, stash, revert, launch a
successor, or remove the worktree.
