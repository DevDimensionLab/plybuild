# Quality Upgrade Handover

Generated: 2026-08-31T11:38:50+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Docker acceptance implementation commit:
  `7e3729321fdd28d7d561e5160fa7065387b68b03`.
- Its exact parent is the launch continuity commit
  `8fdfb9b07412a5c794e8a003c27d5a6dfd7ddb30`, whose exact parent is
  `36156e4dd74db213cecb02d74197da687ec1bb03`.
- After this handoff, obtain the new continuity HEAD with
  `git rev-parse HEAD`; its exact parent must be `7e37293`.
- The clean implementation tree is
  `b5f46f56ba639755e91f482d6df62cf320f810be`; its empty ordinary and
  ignored status SHA-256 is
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
- The implementation changed only `Dockerfile`, focused `Makefile` and
  distribution-contract wiring, `scripts/accept-docker`, its meta-test, and
  shared private acceptance helpers. No production Go, Go test, dependency,
  audit, baseline, inventory, publisher, registry, release, or completed P5
  harness changed.
- No login, push, publication, remote release, ordinary release,
  `make quality`, schema-2 evidence creation, stash, revert, successor launch,
  or worktree removal occurred.

## Continuity Checkpoint

P2A-P5 are complete. P6 remains active, but both required artifact populations
are now demonstrated: a current clean GoReleaser snapshot and a current clean
local Docker image each ran status, upgrade, and build against local fixtures,
with host install kept separate. The only remaining P6 move is current external
schema-2 Q2.8/Q2.9 evidence plus the scoped `make quality` L2 exit gate. P7-P8
remain queued.

The prior snapshot and Docker-blocker archives are answered history. The
Docker-success archive links reciprocally to exactly one NEXT archive for the
final P6 move. Only the launcher's mutable header and prompt regions changed;
its stable execution region remains byte-identical.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. The audit counts
ignored and untracked bytes, so ordinary and ignored status must both be empty
at every authoritative measurement.

## Live Daemon And Classified Build

The mandatory pre-design probe passed before the Dockerfile was read or any
repository file was edited. Fresh evidence is retained at
`/private/tmp/ply-docker-probe.NL83Jw`; its verified evidence-manifest SHA-256
is `219851ff9c5ee84d8a59ac541445b17f8f28aaf0f095a2e363afa48bb593050e`.
It records:

- Docker client 29.4.0 and Docker Desktop 4.70.0 / Engine 29.4.0 server
  identities;
- selected context `desktop-linux` at
  `unix:///Users/perottochristensen/.docker/run/docker.sock`;
- `docker version` and `docker info` exit 0 with daemon host `linux/arm64`;
  and
- buildx 0.33.0, BuildKit v0.29.0, the `docker` builder driver, and supported
  `linux/arm64` platform.

Executed clean daemon builds under
`/private/tmp/ply-docker-classify.3jFRxJ` distinguished environment failures
from two local Dockerfile defects. Default buildx first required unwritable
user activity state and then a credential helper unavailable to the anonymous
build; the accepted local path therefore uses an external empty Docker config
and the non-publishing daemon builder with `DOCKER_BUILDKIT=0`. That real build
then proved the builder lacked `python3`, so the launcher contract failed, and
proved the Dockerfile's forced `GOARCH=amd64` produced the wrong executable for
an arm64 image. Adding builder `python3` and allowing Go to select the native
target were the smallest repairs. An external clean candidate build passed
before those two lines and the acceptance path were committed.

## Authoritative Docker Evidence

The successful clean run is retained at
`/private/tmp/ply-docker-evidence.7e37293-20260831t1110`. Its 506-entry
manifest verifies; the manifest and report SHA-256 values are
`6128c6a1c9471a0ffb64b57f519053582986fcb59907135d041c5ad1a6ec2bee`
and `f9ae547a1bd98fdcbf49b1f4e222d49d6da5bc83ebfe7bec515372ca039c8b09`.
The report binds all claims to source commit `7e37293`, tree `b5f46f5`, and
archive SHA-256
`3b1f83fa4915d09a6c7018c9e977e67da0d3bcdf4ef2533d7f869a79fc5cee30`.

The unique absent tag
`ply-acceptance:7e3729321fdd-20260831t090408-47800` was created by exactly one
actual `docker build --no-cache --tag ... <external-clean-context>` call. No
login, push, publication, or non-anonymous credential operation occurred. The
tag resolved unambiguously before and after every verifier to immutable image
ID and descriptor digest
`sha256:a3b1f58b9861f23f555cdcfc51b6fc82d3452d9425c3314c48997fec5bd475ff`.
Inspection records:

- platform `linux/arm64`, entrypoint `/bin/ply`, empty default command, no
  leaked cross-target configuration, and image size 275,720,311 bytes;
- regular, non-symlink, executable `/bin/ply`, size 23,797,696 bytes and
  SHA-256
  `972bf6fca791fae76eecf09cd9cfcb04ffa7a994405985fa090b304b48ad58f8`;
- a separate initially missing host `GOBIN`; `make install` / `go install
  ./cmd/ply` produced a 26,504,386-byte host executable with SHA-256
  `ca1d0ead122a8da336298a26d6dbc0574ef6bf6d684bb53dff595520241572fa`;
  `verify-install` passed with its artifact override unset; and
- exactly three Docker verifiers, status, upgrade, and build. Each emitted one
  terminal PASS, two non-help behavioral receipts, and three runtime receipts
  proving the exact immutable image ID was executed. The call ledger records
  54 allowed Docker calls, one build, ten runs, and zero publication calls.

The executable Docker meta-test passes all 26 controls, including no or wrong
build/run calls, publication, stale/ambiguous/retargeted identity, wrong
platform/configuration/entrypoint, invalid `/bin/ply`, skipped or duplicate
verifiers, substitution, verifier failure, duplicate PASS, and local output.

## Current Snapshot And Complete Gate

The current clean snapshot rerun is retained under
`/private/tmp/ply-p6-docker-gate.KpjA7l/snapshot-evidence`. Its report SHA-256
is `830da95218ab85884f040ad62088504a72ddd7d212b368b472c1ce2ea0da581a`.
GoReleaser v2.17.1 ran exact credential-cleared argv
`release --snapshot --clean --skip=publish` at `7e37293` and produced one
`darwin/arm64` artifact, size 19,389,170 and SHA-256
`6283963e4fcd72a2bffb2f574eeff3c6e93b6f993982b99a0f5af7b713eb704f`.
Separate host install passed and status, upgrade, and build each produced one
PASS and two behavioral traces against that exact snapshot.

The complete gate root is `/private/tmp/ply-p6-docker-gate.KpjA7l`. Its
verified 534-entry evidence-manifest SHA-256 is
`b7d1e6d82eada86df900b22556a74239ab793142836fb0d288971b7a3be3abbf`.
It retains passing receipts for:

- Docker meta 26/26 and snapshot meta 20/20;
- actual current snapshot acceptance and all four existing host acceptances;
- API and CLI compatibility, CLI surface, and entry/subprocess packages;
- pinned golangci-lint 2.12.2, complete uncached tests across all 27 packages,
  race, vet, `make test`, launcher check, and Make contracts;
- complete preflight, including all eight mutation meta-tests and all other
  script meta-tests;
- an independent rerun of the 15-control audit meta-suite; and
- exact empty-HOME count-2.

No manual evidence was supplied to the authoritative audit views. The full
audit exits 1, never 2, at clean commit `7e37293`; its scorecard SHA-256 is
`c7c4c9a9cff2831bedd6c5417660551d69771f44e625339d8b01f55b29e1daec`.
L0 is 8/8, L1 has six PASS and three UNMEASURABLE, L2 has seven PASS and three
UNMEASURABLE, seven ratchets improve, one regresses, and dirty paths are empty.
The unmeasurable L2 rows are Q2.4, Q2.8, and Q2.9 as expected without current
schema-2 receipts.

The focused Q2.5-Q2.10 view also exits 1, never 2, and has scorecard SHA-256
`4a76baac730083a683580492e7e1eb4b6a7918641775fc1fdaf9a9e87c5bfe3c`.
Q2.5, Q2.6, Q2.7, and Q2.10 pass; only Q2.8 and Q2.9 remain honestly
unmeasurable. Docker and snapshot orchestration remain outside the four-script
`verify-*` audit denominator.

## Next Objective

Finish P6 with one bounded exit-gate move. Recheck all schema-2 and P6 quality
contracts, then implement the smallest `make quality` target and executable
contract test that run preflight, mutation meta-tests and actual mutation
harnesses, all required host/snapshot/Docker acceptance, and the authoritative
audit scoped exactly to `--only Q0.*,Q1.*,Q2.*`. It must require a non-empty
external schema-2 evidence input, external output roots, a real daemon-backed
Docker population, all required criterion populations, and audit exit 0.

After committing that focused wiring, create one external schema-2 document
bound to the exact clean commit, tree, inventory, instruments, subjects, and
retained executable observations. Refresh the already valid Q1.6, Q1.7, Q1.9,
and Q2.4 receipts and add truthful Q2.8 magnitude and Q2.9 non-zero-plus-zero-
artifact receipts for the exact four core verifier population. Run `make
quality` through fresh snapshot and Docker evidence and require the clean
Q0/Q1/Q2 audit to attain L2 with no held material debt. A separate full audit
may continue to exit 1 for L3 debt.

## Start And Stop

Confirm branch, HEAD, exact ancestry, empty ordinary and ignored status,
reciprocal links, launcher `--check`, and the authorized P6 checkpoint before
editing. Probe `docker version` for real Client and Server identities and
require `docker info` exit 0 before running the final gate. Recheck the Docker
and snapshot manifests rather than assuming retained paths are valid.

Stop before P7-P8 implementation, Go or dependency upgrades, production/API/CLI
behavior, audit/inventory/baseline changes, checked-in evidence, remote
publication, publishers, registries, or unrelated acceptance/harness changes.
Keep evidence, reports, logs, contexts, caches, and generated artifacts outside
the worktree. Do not push, merge, publish, release, delete retained evidence,
stash, revert, launch a successor, or remove the worktree.
