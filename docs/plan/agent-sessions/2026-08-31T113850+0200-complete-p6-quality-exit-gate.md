# Agent Session: Complete P6 Quality Exit Gate

Status: NEXT
Session ID: `2026-08-31T113850+0200-complete-p6-quality-exit-gate`
Created: `2026-08-31T11:38:50+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `e61cefa211102b8b71d8a0b163153025d872614ffb12af9c8faef81857c9edec`
Previous: [2026-08-31T092124+0200-resume-p6-docker-daemon-blocker.md](2026-08-31T092124+0200-resume-p6-docker-daemon-blocker.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Finish P6 with the one remaining exit-gate move: create current external
schema-2 evidence for the completed acceptance population and add the smallest
scoped `make quality` target that proves Q0-Q2 as one executable apparatus.
Require a fresh clean snapshot and fresh clean local Docker image during the
real gate. Finish only when the exact `--only Q0.*,Q1.*,Q2.*` audit exits 0 at
L2 with no missing population or held material debt.

# Authorized Roadmap

P2A-P5 are complete. P6 is active; current snapshot and Docker acceptance are
complete at `7e37293`. P7-P8 remain queued until this move succeeds. Change
only the smallest focused `Makefile` quality wiring and executable Make
contract test, plus a private test helper only if the target cannot otherwise
be falsified cleanly. Preserve the existing preflight, mutation, host,
snapshot, Docker, and audit entry points rather than reimplementing them.

Do not change production Go or Go tests, `Dockerfile`, `.goreleaser.yml`,
existing acceptance or mutation implementations, dependencies,
`.quality/inventory`, audit/parser/scanner/baseline instruments, API/CLI
behavior, fixtures, publishers, `docker-publish.sh`, registries, or release
configuration. Keep the schema-2 document, audit output, build contexts,
caches, reports, logs, and generated artifacts outside the worktree. Never
create `.agent-task/current.md` or `.quality/manual-evidence.json`.

# Measurements At Start

The Docker implementation commit is
`7e3729321fdd28d7d561e5160fa7065387b68b03`, exact parent
`8fdfb9b07412a5c794e8a003c27d5a6dfd7ddb30`, clean tree
`b5f46f56ba639755e91f482d6df62cf320f810be`. After this handoff the new
continuity HEAD must have exact parent `7e37293` and empty ordinary and ignored
status.

The current Docker evidence is
`/private/tmp/ply-docker-evidence.7e37293-20260831t1110`. Its verified
506-entry manifest SHA-256 is
`6128c6a1c9471a0ffb64b57f519053582986fcb59907135d041c5ad1a6ec2bee`.
The fresh image is `linux/arm64`, entrypoint `/bin/ply`, immutable ID
`sha256:a3b1f58b9861f23f555cdcfc51b6fc82d3452d9425c3314c48997fec5bd475ff`;
status, upgrade, and build each have one PASS, two behavioral receipts, and
three immutable-ID runtime receipts. Host install passed separately with its
artifact override unset. The meta-test passes 26/26.

The complete current gate is `/private/tmp/ply-p6-docker-gate.KpjA7l`; its
verified 534-entry manifest SHA-256 is
`b7d1e6d82eada86df900b22556a74239ab793142836fb0d288971b7a3be3abbf`.
It contains the current snapshot rerun, all requested compatibility,
lint/test/race/vet/Make/launcher/preflight/host/meta/empty-HOME receipts, and
the full and focused audits. With no manual evidence, Q2.5, Q2.6, Q2.7, and
Q2.10 pass while Q2.8 and Q2.9 remain unmeasurable. The full audit also leaves
Q2.4 unmeasurable, as expected.

The last valid external schema-2 document is
`/private/tmp/ply-snapshot-gate.06e3ac4/manual-evidence-q16-q17-q19-q24.json`,
SHA-256
`7a09ec592025564f6600b3edf21fd0c2d56dc28fd73a92de2b3bb80024f32e70`.
It is prior-commit evidence and must not be reused as current evidence; use it
only as a schema and receipt-population input to an independent fresh review.

# Role And Boundaries

The new target must invoke, through existing entry points, complete preflight;
every mutation meta-test and actual mutation harness; the four host acceptance
verifiers; fresh snapshot acceptance; fresh Docker acceptance; and the
authoritative audit scoped exactly to `--only Q0.*,Q1.*,Q2.*`. It must accept
explicit external tool, evidence, cache, and output paths, require a non-empty
external schema-2 document, and reject repository-local generated state.

The target must fail closed for a missing or duplicate required stage, an
empty mutation or acceptance population, wrong stage argv or ordering where
ordering is contractual, missing or repository-local evidence/output, a
skipped fresh artifact population, a missing criterion population, an audit
scope other than exact Q0-Q2, audit exit 1 or 2, or a report that does not
attain L2. It may not turn the separate full L3 report into an exit gate.

# Schema-2 Evidence Contract

Make the focused wiring commit before producing current receipts. Then create
one external schema-2 document bound to that exact clean commit, commit tree,
module, inventory, audit instruments, and every declared subject or acceptance
script digest. Refresh the already proved Q1.6, Q1.7, Q1.9, and Q2.4 receipts
from fresh current observations. Add Q2.8 observations that compare a real
input magnitude with the treatment in the declared direction, and Q2.9
observations that jointly prove bad input exits non-zero and creates or changes
zero artifacts. Cover the exact non-empty four-script core verifier population;
do not infer a receipt from grep, counts, help output, or exit status alone.

Canonicalize the JSON exactly as required by the parser, record the whole-file
and criterion-object SHA-256 values, and pass it explicitly with
`--manual-evidence`. First require the focused Q1.6/Q1.7/Q1.9/Q2.4/Q2.8/Q2.9
audit to exit 0. Then require the exact Q0/Q1/Q2 audit used by `make quality`
to exit 0 and attain L2 with no absent, duplicate, stale, dirty, malformed, or
held evidence.

# Required Reading

Confirm branch, HEAD, exact ancestry, clean and ignored status, reciprocal
archive links, launcher `--check`, and the authorized checkpoint block. Read
the rolling handover and this archive; recheck the complete P6 quality-target,
schema-2, scorecard, baseline, preflight, mutation, host acceptance,
snapshot/Docker, Make, compatibility, fixture, and Q0-Q2 contracts before
choosing the target shape.

# Three Moves

Probe Docker first. Require `docker version` to contain real Client and Server
identities and `docker info` to exit 0 for the selected context before relying
on Docker in the final apparatus. Revalidate both retained manifests. Make one
focused implementation commit, produce fresh external schema-2 evidence at
that exact commit, and run focused contract controls plus the actual complete
`make quality` path with fresh snapshot and Docker evidence.

Then rerun current snapshot/Docker meta and acceptance, API/CLI and
entry/subprocess compatibility, pinned lint, complete tests/race/vet,
launcher and Make contracts, complete preflight, existing host acceptance,
the standalone audit meta-suite, focused and full audits, and empty-HOME
count-2. Retain executable evidence for every declared stage and verifier.

# Automatic Handoff

After the exact Q0-Q2 gate exits 0, mark P6 complete and P7 active while P8
remains queued. Rewrite the rolling handover and roadmap, answer this archive,
create exactly one reciprocal NEXT archive for the first bounded P7 toolchain
baseline move, replace only the launcher's mutable regions, run launcher and
handoff contracts, and make the normal `docs: prepare next agent session`
commit. Do not implement P7 in this session.

Do not launch a successor, push, merge, publish, release, delete retained
evidence or the local image, stash, revert, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
