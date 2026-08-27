# Agent Session: Resume P6 Docker Daemon Blocker

Status: NEXT
Session ID: `2026-08-27T070412+0200-resume-p6-docker-daemon-blocker`
Created: `2026-08-27T07:04:12+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `e7d0aaa29725be15e70fad5c330ecdb33c8c266ee442c7cae80f2cc957f30735`
Previous: [2026-08-27T064625+0200-build-p6-docker-acceptance.md](2026-08-27T064625+0200-build-p6-docker-acceptance.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Resume P6 with the same single acceptance-evidence target: a fresh local Docker
image. Proceed only after a real local daemon is already reachable. Build from
clean committed source through a non-publishing daemon path, revalidate host
install as the separate `make install` / `go install` contract, and exercise
status, upgrade, and build through the exact immutable image against local
fixtures. Finish only with executable falsifiability controls and retained
external evidence proving every declared verifier ran through that image.

# Authorized Roadmap

P2A-P5 are complete. Snapshot acceptance is complete in `06e3ac4`; P6 remains
active and P7-P8 remain queued. Change only the smallest focused Docker
acceptance orchestration and meta-test under `scripts`, a shared private helper
under `test/acceptance`, and focused `Makefile` or
`test/makefile_distribution_test.sh` wiring needed to expose the target. Change
`Dockerfile` only if an executed clean daemon build exposes a classified local
acceptance defect, and then make only the smallest non-publishing repair.

Do not add `make quality`, reimplement snapshot acceptance, change
`.goreleaser.yml`, production Go or Go tests, `.quality/inventory`, audit or
baseline instruments, completed P5 harnesses, dependencies, API/CLI behavior,
unrelated fixtures, `docker-publish.sh`, publishers, registry state, remote
releases, or distribution outside the local image. Keep contexts, caches,
reports, logs, and generated artifacts outside the worktree. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Measurements At Start

The prior Docker attempt made no implementation change. Its launch continuity
was `9fa49c304b28ceccb1d55dee9c6cf980f75dc160`, exact parent `06e3ac4`, with
clean ordinary and ignored status. After this handoff, the new continuity HEAD
must have exact parent `9fa49c3`.

Docker client 29.4.0 and buildx 0.33.0 selected `desktop-linux` at
`unix:///Users/perottochristensen/.docker/run/docker.sock` on macOS 15.3.1
arm64, but `docker version`, `docker info`, and `docker desktop status` had no
Server and exited 1. The managed workspace sandbox could not write Docker
Desktop's required `~/Library/Containers/com.docker.docker` state, and UI
control was unavailable. Exact evidence remains at
`/private/tmp/ply-docker-blocker.KRsWoi`; its manifest SHA-256 is
`991d7bf4a7a739b9d971da7776a8be7e447d2aad6c931a7cf72c1bbdaeb682ae`.

Before reading the Dockerfile as an implementation input or editing anything,
require `docker version` to contain real Client and Server identities and
require `docker info` to exit 0 for the selected context. Record fresh client,
server, builder, host, and supported-platform identities externally. If no
daemon is available, do not implement from static assumptions; retain the new
probe and hand off the exact blocker again.

# Role And Boundaries

Use a clean external archive or clone at exact HEAD and a unique tag that
cannot name a pre-existing image. Require one actual daemon build, zero
login/push/publish operations, a build-created immutable image ID, exact
unambiguous tag resolution, and inspection of OS, architecture, entrypoint,
creation/configuration metadata, and the regular executable `/bin/ply`. Record
the build argv, commit, context identity, tag, image ID, platform, executable
path, size, and digest. Reject stale or ambiguous identity, tag retargeting,
wrong platform/configuration, artifact substitution, and missing, symlinked, or
non-executable image content.

Keep host install semantically separate in a missing external GOBIN and run
`verify-install` with its artifact override unset. Reuse the existing status,
upgrade, and build verifiers through the smallest runtime bridge. Mount only
the exact external fixture/cache/output state they require. Require the exact
non-empty verifier population, one terminal PASS each, non-help behavior, and
runtime receipts proving each command used the immutable image ID rather than
the tag or a host executable.

Add controls that fail closed for no Docker build/run; wrong build or run argv;
login, push, or publication; stale, missing, or ambiguous image identity; wrong
platform, entrypoint, or configuration; missing, symlinked, or non-executable
`/bin/ply`; skipped or duplicate verifiers; tag retargeting or image
substitution; verifier failure; duplicate terminal PASS; and repository-local
output. Do not weaken assertions or accept help-only, static, grep-only, or
exit-status-only evidence.

# Required Reading

Confirm branch, HEAD, exact ancestry, clean and ignored status, reciprocal
archive links, launcher `--check`, and the authorized checkpoint block. Read
the rolling handover and this archive, then recheck the complete P6/gate,
design, Docker/distribution, snapshot, acceptance, compatibility, fixture, and
Q2.5-Q2.10 contracts before defining the image population.

# Three Moves

Run a clean external daemon probe before choosing the design. Implement one
executable Docker acceptance path, run the real fresh build, separate host
install, and exact status/upgrade/build verifiers, and retain every receipt.
Make one focused implementation commit. Then run focused meta-tests, actual
Docker acceptance, retained snapshot acceptance/meta, API/CLI and
entry/subprocess compatibility, pinned lint, complete tests/race/vet, launcher
and Make contracts, complete preflight, existing host acceptance, audit
meta-suite, full and focused Q2.5-Q2.10 audit views, and empty-HOME count-2.

# Automatic Handoff

After successful Docker evidence, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the final
P6 schema-2 evidence and scoped `make quality` exit-gate move, replace only the
launcher's mutable regions, run launcher/handoff contracts, and make the normal
`docs: prepare next agent session` commit. Do not launch a successor, push,
merge, publish, release, delete retained evidence, stash, revert, or remove the
worktree. P6 remains active until snapshot and Docker evidence plus the scoped
quality audit exit gate all pass.
<!-- CODEX_SESSION_PROMPT_END -->
