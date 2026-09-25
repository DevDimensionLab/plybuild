# Agent Session: Build P6 Docker Acceptance

Status: ANSWERED - HISTORY
Session ID: `2026-08-27T064625+0200-build-p6-docker-acceptance`
Created: `2026-08-27T06:46:25+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `07fbc927bebce327a22c92dcd6404e1a53f1b7d0e7a40081295b7584d54c10ff`
Previous: [2026-08-27T054500+0200-build-p6-snapshot-acceptance.md](2026-08-27T054500+0200-build-p6-snapshot-acceptance.md)
Next: [2026-08-27T070412+0200-resume-p6-docker-daemon-blocker.md](2026-08-27T070412+0200-resume-p6-docker-daemon-blocker.md)
Outcome: Blocked before implementation because the installed Docker client had no reachable daemon and the managed sandbox could not start Docker Desktop; exact external probe evidence is retained.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P6 with exactly one acceptance-evidence target: a fresh local Docker
image. Build it through a real daemon-backed, non-publishing path from clean
committed source, revalidate host install as its separate `make install` /
`go install` contract, and exercise status, upgrade, and build through the
exact immutable image against local fixtures. Finish only with executable
falsifiability controls and retained external evidence proving every declared
verifier ran through that image.

# Authorized Roadmap

P2A-P5 are complete. The host-platform GoReleaser snapshot target is complete
in `06e3ac4`; P6 remains active and P7-P8 remain queued in
`docs/plan/quality-upgrade.md`. This session may change only the smallest
focused Docker acceptance orchestration and meta-test under `scripts`, a shared
private helper under `test/acceptance`, and focused `Makefile` or
`test/makefile_distribution_test.sh` wiring needed to expose the image
acceptance target. Change `Dockerfile` only if an actually executed clean
daemon build exposes a classified local-acceptance defect and only with the
smallest non-publishing repair.

Do not add `make quality`, reimplement snapshot acceptance, change
`.goreleaser.yml`, production Go or Go tests, `.quality/inventory`,
audit/parser/scanner/baseline code, completed P5 harnesses, Go/dependencies,
exported API/CLI behavior, fixtures unrelated to the four core flows,
`docker-publish.sh`, package-manager publishers, registry publication, remote
releases, or distribution outside the local image. Keep build contexts, logs,
caches, reports, and generated artifacts outside the worktree. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Measurements At Start

The clean snapshot implementation is
`06e3ac43282c1f4ccd3b99cec78d02e009a90645`, exact parent
`61763f4b5728fa77f9c947c66fcc04ba5136c208`, tree
`0cb65473791c5eebf12f5aa6df9534be4a955a01`, and clean status SHA-256
`6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
After launch, the new continuity HEAD must have exact parent `06e3ac4`.

The accepted GoReleaser v2.17.1 run used exact argv
`release --snapshot --clean --skip=publish`, cleared all five release
credentials, and produced exactly one `darwin/arm64` executable:
`dist/plybuild_darwin_arm64_v8.0/ply`, size 19,389,170, SHA-256
`28d7012a81e9bdab2b776a8be64f84d12fc9e2f239e52dcbfc2932f7469d2e4f`.
Its report and evidence-manifest SHA-256 values are
`6e759f75ae40c27c6312cfc4f4c1e8a8c4eab56119a64a88fa4c4c62dec40f61`
and `bb33ae40ccab0992284318a123c9462fa86ae6a1e5dd2c7b881d0b302c4ab9a1`.
Host install passed separately with its artifact override unset; status,
upgrade, and build each emitted one terminal PASS and two non-help behavioral
traces against that exact snapshot. All 20 snapshot controls pass.

P5 remains eight executable subjects, 80 declared, 80 killed, 0 survived, and
0 unusable. The refreshed external schema-2 document SHA-256 is
`7a09ec592025564f6600b3edf21fd0c2d56dc28fd73a92de2b3bb80024f32e70`.
The exact Q1.6/Q1.7/Q1.9/Q2.4 audit exits 0 with scorecard SHA-256
`2b8d42fa6457ab02bd6dbed85d3fd13ccd0552cae794d8cc7c6d70644e09cbae`.
The full audit exits 1, never 2, with scorecard SHA-256
`9b75a367fc07ddc698787103f80360227e6f03827009c62c4a80a59dcd547b4d`:
L0 is 8/8, L1 is 9/9, L2 is 8/10, seven ratchets improve, the Q3.4
documentation-phrase ratchet records one regression from immutable prompt
history, dirty paths are empty, and seven rows remain non-passing. The focused
Q2.5-Q2.10 scorecard SHA-256 is
`367819708c9959d6b6bd6085e7a362124c95a2f9bb744db48a7294a1ffa0d482`;
Q2.5, Q2.6, Q2.7, and Q2.10 pass while Q2.8 and Q2.9 remain manual.

API/CLI and entry/subprocess compatibility, pinned lint, complete tests/race/
vet, `make test`, 62 launcher controls, Make contracts, four host acceptance
flows, snapshot meta/acceptance, complete preflight, the 15-control audit
meta-suite, and empty-HOME count-2 pass. API and CLI report SHA-256 values are
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

# Role And Boundaries

Work autonomously on `codex/upgrade-quality`. Before editing, confirm branch,
HEAD, exact ancestry, clean and ignored status, reciprocal archive links,
launcher `--check`, and the authorized checkpoint block. Probe the real Docker
client and daemon before choosing an implementation shape. Record exact client,
server, builder, host, and selected image-platform identities; a static
Dockerfile reading or a recording-only fake is not acceptance evidence.

Use a clean external source archive or clone at exact HEAD and a unique local
tag that cannot name a pre-existing image. Require an actual daemon build,
zero login/push/publish operations, exact unambiguous post-build resolution to
one immutable image ID, a build-created identity rather than a stale tag, and
inspection of OS, architecture, entrypoint, creation/config metadata, and the
regular executable inside the final image. Record command, commit, context
identity, image tag and ID, platform, executable path, size, and digest.
Reject missing or multiple identities, mutable-tag substitution, wrong
platform/configuration, symlinked or non-executable image artifacts, and
images not produced by the commanded run.

Keep host install semantically separate: it remains `make install` /
`go install ./cmd/ply` into a missing temporary GOBIN and `verify-install`
runs with its artifact override unset. Reuse the existing status, upgrade, and
build acceptance verifiers through the smallest artifact/runtime bridge rather
than duplicating their behavioral assertions. Mount only the exact external
fixture/cache/output state needed for those behaviors. Require exact non-empty
verifier selection, one terminal PASS per verifier, non-help behavior, and
runtime receipts proving every selected command used the immutable image ID,
not merely the tag or a host binary.

Add focused meta-controls that fail closed for at least: no Docker build/run;
wrong build or run argv; login, push, or another publishing operation; stale
tag or image; missing or ambiguous identity; wrong platform, entrypoint, or
configuration; missing, symlinked, or non-executable `/bin/ply`; skipped or
duplicate verifier selection; tag retargeting or image substitution; a
verifier failure; duplicate terminal PASS; and repository-local output. Do not
weaken an assertion, accept help-only smoke as behavioral evidence, or infer
runtime use from filenames, static config, grep, or exit status alone.

# Required Reading

Read the rolling handover, this archive, the complete P6 and checkpoint-gate
entries, both design documents, `.quality/README.md`, `.quality/inventory`, and
the complete Q2.5-Q2.10 discovery/parser logic. Read the full Dockerfile,
`docker-publish.sh`, Docker Make targets and their meta-tests, `.dockerignore`
if present, the complete snapshot acceptance orchestration/helper/meta-test,
all four acceptance verifiers and meta-tests, shared acceptance helpers and
fixtures, CLI/API compatibility contracts, and relevant Docker/distribution
archive history before defining the image population.

# Three Moves

1. Define the exact build context, final-image population, immutable identity,
   platform/config/executable provenance, exclusions, and non-publishing build
   and run argv. Run a clean external daemon probe before choosing the smallest
   implementation shape, and classify any Dockerfile defect from executed
   evidence.
2. Implement one executable Docker acceptance path and its falsifiability
   controls. Run a real fresh daemon build, the separate host-install verifier,
   and exact status/upgrade/build verifiers through the immutable image; retain
   every report, inspection record, runtime receipt, and image identity
   externally.
3. Make one focused implementation commit, then run focused meta-tests, actual
   Docker acceptance, retained snapshot acceptance/meta, API/CLI and
   entry/subprocess compatibility, pinned lint, complete tests/race/vet,
   launcher and Make contracts, complete preflight, existing host acceptance,
   audit meta-suite, full and focused Q2.5-Q2.10 audit views, and empty-HOME
   count-2. Record the bounded Docker result and hand off the separate final P6
   evidence/`make quality` move without implementing it.

# Automatic Handoff

Before ending, finish the coherent Docker-image move or record an exact
resumable blocker. Rewrite the rolling handover, update the roadmap, answer
this archive, create exactly one reciprocal NEXT archive for the final P6
evidence and `make quality` exit-gate move, replace only the launcher's mutable
regions, run launcher and handoff contracts, and make the normal
`docs: prepare next agent session` continuity commit after the focused
implementation commit. Do not launch a successor, push, merge, publish, invoke
an ordinary production release, enable a publisher, delete retained evidence,
stash, revert, or remove the worktree. P6 remains active until snapshot and
Docker evidence plus the scoped `make quality` audit exit gate all pass.
<!-- CODEX_SESSION_PROMPT_END -->
