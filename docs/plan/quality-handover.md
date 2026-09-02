# Quality Upgrade Handover

Generated: 2026-09-02T14:23:12+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P7 terminal-pair implementation commit:
  `0f93a527d1f95cc043a695101746d359c74c9c54`.
- Its exact parent is the launch continuity commit
  `e25a98546c20f8707ffe4f71b1d544fc8fb4944f`, whose exact parent is the P7
  toolchain implementation `16ecb67eebb6450226d3264f30a64a14391b5245`.
- The implementation tree is
  `ebbc91e08907b37175faf45ddc3283e8bd0ed771`.
- After this handoff, obtain the new continuity HEAD with
  `git rev-parse HEAD`; its exact parent must be `0f93a52`.
- The implementation and every accepted measurement had empty ordinary and
  ignored status. No push, merge, publication, release, stash, revert,
  successor launch, retained-evidence deletion, image deletion, or worktree
  removal occurred.

## Continuity Checkpoint

P2A-P6 are complete. P7 remains active after its maintained-toolchain move and
first terminal dependency group; further dependency groups remain. P8 is
queued.

The answered terminal archive links reciprocally to exactly one NEXT archive
for direct `github.com/sirupsen/logrus`. Only the launcher's mutable header and
prompt regions changed; its stable executable skeleton remains byte-identical.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. Downloaded tools,
module/build caches, reports, generated files, build contexts, and audit
evidence remain external.

## Terminal-Pair Decision

The bounded dependency selection is:

- direct `golang.org/x/term v0.5.0` -> `v0.29.0`;
- coupled indirect `golang.org/x/sys v0.5.0` -> `v0.30.0`.

As verified on 2026-09-02, the latest releases are x/term v0.45.0 and x/sys
v0.47.0, both declaring Go 1.25. The selected x/term v0.29.0 declares Go 1.18
and requires x/sys v0.30.0, which also declares Go 1.18. The immediately newer
x/term v0.30.0 and x/sys v0.31.0 raise their module floors to Go 1.23. The
selected pair is therefore the highest release pair compatible with the
retained `go 1.18` language floor, while remaining fully supported by preferred
Go 1.26.7 and the pinned golangci-lint 2.12.2 and GoReleaser 2.17.1 tools.

The external old/candidate comparison selected 233 modules, 143 non-standard
packages, and 3,551 graph edges in both states. The only selection changes are
the two authorized modules and their version-bearing edges. The existing tidy
projection changes from 207 to 211 lines only because of those pair checksums;
no unrelated version or checksum enters the candidate projection. The
`golang.org/x/term.ReadPassword` implementations used by `prompt.go` remain
byte-identical for the supported platforms.

Primary module repositories and the Go vulnerability database record
GO-2026-5024 for `golang.org/x/sys/windows` before v0.44.0. Exact host and
Windows scans preserve the pre-existing 22 reachable findings and do not
report this advisory as symbol-reachable; the module-only scan reports it in
both old and new selections. Eliminating that module-level advisory would
require a later x/sys line and a higher language floor, so it was not silently
mixed into this bounded pair move.

## Implementation And Compatibility

The exact Go command performed
`go get golang.org/x/term@v0.29.0`; dependency metadata was not hand-edited.
Commit `0f93a52` changes only `go.mod` and `go.sum` with six insertions and two
deletions. It does not change source, public Go API, CLI output or semantics,
toolchain declarations, Docker/release inputs, quality tools or thresholds,
baselines, compatibility allowlists, acceptance/mutation populations,
packaging, publishers, registries, credentials, or P8 code.

The external candidate and committed implementation pass the focused terminal
callers, complete tests, race, vet, host/Linux/Windows builds, pinned
lint, public help, API/CLI and entry/subprocess compatibility, launcher and
Make contracts, complete preflight, host acceptance, snapshot and Docker meta
and acceptance, audit meta, focused/Q0-Q2/full audits, vulnerability equality,
and empty-HOME count-2.

The declared and verified toolchain remains exact Go 1.26.7. The retained
official executable is
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in `PATH` as well as passing exact `GO`; literal `go`
subprocesses must resolve to it. Keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, and do not inject ambient `GOFLAGS`.

Warm fresh external module caches from Git archives outside the worktree. A
real-tree `go mod download all` materializes historical checksum debt and is
not an authorized tidy. The accepted regression includes both current and
v1.0.1 archive bootstraps for compatibility.

## Accepted External Evidence

The terminal selection root is
`/private/tmp/ply-p7-terminal-selection.ctFhJk`. Its verified 76,252-entry
manifest SHA-256 is
`c5cb5e2706ecccfc8eea7c80cadaa747c458e605b2919ab80cd378b263f84aca`;
its summary SHA-256 is
`371e40242b52af838354ea375a82ad13d63ffaf8731a4ba53ad91ad70e0bcfef`.

The commit-bound schema-2 review root is
`/private/tmp/ply-p7-terminal-quality-review.0f93a52.nzAZ5z`. Its verified
252,440-entry manifest SHA-256 is
`5213d578d639a85bae51c2b85cba73ff72fc380e950afa23dabc1728e3370d4a`.
The manual evidence SHA-256 is
`952ffd64db2d5207dc12e0c25c2e3cf58844c6e744b2e7c5f9cfba401fac228f`;
all 242 reviewed subjects and 207 declared files bind to `0f93a52`, all manual
receipts validate, and the focused scorecard SHA-256 is
`52800b59a800bed8597fa0e29959d9a72ca850fa55dcd0ddb57a1d6c1e33d56a`.

The exact complete quality apparatus root is
`/private/tmp/ply-p7-terminal-quality-gate.0f93a52.V0TBoc`. Its verified
268,572-entry manifest SHA-256 is
`e8027bcb9e2a30ff2fbbf3e3bb8c1a1072abe98a7f55ff6a1fb313c4200d91df`.
Exact `make quality` exits 0. The Q0-Q2 scorecard SHA-256 is
`1ec4203b099e49017f67a906f19a2ceeb69f1e86b7ca33bd6f8cd1bf23a27cc5`:
all 27 rows pass at L2, mutation and acceptance populations are 8/8 and 4/4,
and held, regressed, current-not-comparable, and dirty counts are zero. Its
Docker artifact's internal evidence manifest is verified at SHA-256
`05ab1e065d8fa18d06eb6e4127fc0fffcea223eb68be627b744c7f4006976`.

The independent accepted regression root is
`/private/tmp/ply-p7-terminal-regression-gate.0f93a52-final.hIYQFq`. Its
verified 68,961-entry manifest SHA-256 is
`1a4414491a83262c4e7684dcc50d7c8372310a66d41e67cf0b9b9b30aca4964a`.
All 40 ledger stages pass. The focused and Q0-Q2 scorecards are the hashes
above. The full scorecard SHA-256 is
`e8d7dea22e12e08e64cbee667a4a24d1959d17b4ec1b93f9b0ad1bc08e6f808e`;
it exits the expected 1, never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7.
Its fresh Docker artifact's internal evidence manifest is verified at SHA-256
`6ae30d5c3aa2447c3bb655662f049acbe66ea72bf697d1b9bc400fcb77b8d305`.

Two timing diagnostics are retained but are not accepted gates: a first
preflight launcher subprocess and a first Make-contract launcher subprocess
received signals after producing partial logs. Immediate isolated reruns and
the complete quality and final regression ledgers pass the exact contracts.
The earlier incomplete regression root is likewise diagnostic only; the final
root above is authoritative.

## Next Objective

Implement only direct `github.com/sirupsen/logrus v1.9.0` -> `v1.9.3`.
Reverify release and vulnerability evidence before editing. The retained next
selection root is
`/private/tmp/ply-p7-next-group-selection.0f93a52.Nl74oE`; its verified
31,623-entry manifest SHA-256 is
`e91b0260398550ed95d15063dbfbb96e01c2467b9f878ea70a0e459f25ec9b01`
and its summary SHA-256 is
`b822c4102c01ddcad1032271c99ef0f7c92e76cb7dc9bbbdc02f8e1753a26c4b`.

As probed on 2026-09-02, current logrus v1.10.0-v1.10.2 require Go 1.23;
v1.9.3 declares Go 1.13 and is the highest release retaining the module's Go
1.18 floor. The external `go get github.com/sirupsen/logrus@v1.9.3` projection
changes only the logrus selection and adds its two checksums; current x/sys
v0.30.0 already exceeds logrus's lower requirement, so no coupled module moves.

Make one dependency-only implementation commit. Stop if independent MVS replay
requires another module, or if the candidate changes behavior, API/CLI,
acceptance, artifacts, vulnerability reachability, or baseline debt. Do not
combine Cobra, Viper, another dependency, a source fix, or P8 work with this
group.

## Start And Stop

Confirm branch, exact ancestry, clean ordinary and ignored status, reciprocal
links, launcher `--check`, the P7/P8 queue, implementation commit/tree, and all
accepted manifests before editing. Read the active archive, this handover, the
P7 roadmap, module graph, logrus callers/tests, toolchain contract,
compatibility, snapshot/Docker, quality, and audit contracts.

Stop before any dependency outside logrus, behavior/API/CLI change,
quality-tool upgrade, P8 domain work, inactive packaging work, publication,
publisher/registry/credential change, or release. Do not push, merge, publish,
release, delete retained evidence or images, stash, revert, launch a successor,
or remove the worktree.
