# Agent Session: Build P6 Snapshot Acceptance

Status: NEXT
Session ID: `2026-08-27T054500+0200-build-p6-snapshot-acceptance`
Created: `2026-08-27T05:45:00+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `a5280f7ef856db7a1af4feb2a9619d1bee6cfeb4deaeda44be616e983be6d068`
Previous: [2026-08-27T044626+0200-build-p5-interactive-build-harness.md](2026-08-27T044626+0200-build-p5-interactive-build-harness.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Begin P6 with exactly one acceptance-evidence target: the host-platform
GoReleaser snapshot binary. Build a real fresh snapshot through the established
credential-cleared non-publishing path, revalidate host install as its separate
`make install` / `go install` contract, and exercise status, upgrade, and build
against the exact snapshot executable. Finish only with executable
falsifiability controls and retained external evidence that proves every
declared verifier ran against that artifact.

# Authorized Roadmap

P2A-P5 are complete. P6 is active; P7-P8 remain queued in
`docs/plan/quality-upgrade.md`. This session may change only the smallest
focused acceptance orchestration and meta-test under `scripts`, a shared
private helper under `test/acceptance`, and focused `Makefile` or
`test/makefile_distribution_test.sh` wiring needed to expose the snapshot
acceptance target. Change `.goreleaser.yml` only if an actually executed clean
snapshot exposes a classified local-artifact defect and only with the smallest
non-publishing repair.

Do not start Docker acceptance, add `make quality`, change production Go or Go
tests, alter `.quality/inventory`, audit/parser/scanner/baseline code, completed
P5 harnesses, Go/dependencies, exported API/CLI behavior, fixtures unrelated to
the four core flows, package-manager publishers, remote releases, publication,
or distribution outside the local snapshot. Keep all snapshot output, tools,
caches, reports, and generated artifacts outside the worktree. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Measurements At Start

The clean P5.8 implementation is
`58c5224183308e8e48aa4d91b25b65c1e4021a96`, exact parent
`5bcd6c48a8ce6a09f0553368e1f5c34e2bf09cef`, tree
`8cc32082df0fbc4af2f10810c8ecf5274fecbc63`, and clean status SHA-256
`6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
After launch, the new continuity HEAD must have exact parent `58c5224`.

P5 has eight executable subjects, 80 declared mutations, 80 killed, 0
survived, and 0 unusable. The `interactive-build` report and T1-T10 meta-log
SHA-256 values are
`109619d5d336a52f01d5f751d3b89b5bcfa85f2ba512cdc5073fbc58e1954d5f`
and `22a0ac95820514fec641d72289c6ef1a3a56023cc6471f4c2e154087199c9f81`.
The external eight-subject schema-2 document SHA-256 is
`b12d4c262ecd16c9529887968cc33324552b08a7cd1d300ff0f1dff6e34acbba`;
its Q2.4 evidence-object SHA-256 is
`2565dd9824e224520004774521164f0d24ff771b134f7662125a337ae27dd09d`.

The focused Q1.6/Q1.7/Q1.9/Q2.4 audit exits 0 with scorecard SHA-256
`a5478b88f4dbd99dfa5561eafa09117c649a3f3998b79b34818b54b3970d6077`.
The full audit exits 1, never 2, with scorecard SHA-256
`9032b0e3763c64f820564f1608fd6ef52b0a751311c4dcbcd9a2fca8b1eef972`:
L0 is 8/8, L1 is 9/9, Q2.1-Q2.4 pass, seven ratchets improve, the Q3.4
documentation-phrase ratchet records one regression from immutable prompt
history, dirty paths are empty, and seven P6-P8 rows remain non-passing.

API/CLI and entry/subprocess compatibility, pinned lint, complete tests/race/
vet, `make test`, 62 launcher controls, Make contracts, four host acceptance
flows, complete preflight, the 15-control audit meta-suite, and empty-HOME
count-2 pass. API and CLI report SHA-256 values are
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

# Role And Boundaries

Work autonomously on `codex/upgrade-quality`. Before editing, confirm branch,
HEAD, exact ancestry, clean and ignored status, reciprocal archive links,
launcher `--check`, and the authorized checkpoint block. Establish the exact
GoReleaser version and host artifact naming from the accepted config and tool
metadata. Use a pinned external GoReleaser binary; do not install it into the
repository or change module dependencies.

Run the exact local-only snapshot argv with every release credential cleared.
Require a clean HEAD source, a fresh external output root, zero remote
publication, and exact unambiguous discovery of one regular executable for the
host GOOS/GOARCH. Reject stale pre-existing output, zero or multiple candidates,
symlinks, wrong platform paths, and artifacts not produced by the commanded
run. Record tool identity, command, commit, artifact path, size, and digest.

Keep host install semantically separate: it remains `make install` /
`go install ./cmd/ply` into a missing temporary GOBIN. Reuse the existing
status, upgrade, and build acceptance verifiers through their artifact override
rather than duplicating their behavioral assertions. Require exact non-empty
verifier selection and terminal PASS records, and prove each selected verifier
received the snapshot path. A help-only smoke is not evidence.

Add focused meta-controls that fail closed for at least: no GoReleaser call;
wrong or publishing argv; leaked credentials; stale, missing, ambiguous,
symlinked, or non-executable artifact; skipped or duplicate verifier selection;
artifact substitution; a verifier failure; and repository-local output. Do not
weaken an assertion or accept a recording-only distribution run as real
snapshot evidence.

# Required Reading

Read the rolling handover, this archive, the complete P6 and checkpoint-gate
entries, both design documents, `.quality/README.md`, `.quality/inventory`, and
the complete Q2.5-Q2.10 discovery/parser logic. Read the full `.goreleaser.yml`,
Make snapshot/release/acceptance targets, distribution and Make meta-tests, all
four acceptance verifiers and meta-tests, shared acceptance helpers and
fixtures, CLI/API compatibility contracts, and relevant archive history before
choosing the orchestration. Do not infer artifact production or verifier
execution from filenames, static config, grep, or exit status alone.

# Three Moves

1. Define the exact snapshot artifact population, host-platform selection,
   provenance record, exclusions, and non-publishing invocation. Run a clean
   external probe before choosing the smallest implementation shape.
2. Implement one executable snapshot acceptance path and its falsifiability
   controls. Run a real fresh external snapshot, the separate host-install
   verifier, and exact status/upgrade/build verifiers against the snapshot;
   retain every report and artifact externally.
3. Make one focused implementation commit, then run focused meta-tests, actual
   snapshot acceptance, API/CLI and entry/subprocess compatibility, pinned
   lint, complete tests/race/vet, launcher and Make contracts, complete
   preflight, existing host acceptance, audit meta-suite, full and focused
   Q2.5-Q2.10 audit views, and empty-HOME count-2. Record the bounded P6 result
   and hand off the separate Docker acceptance move without implementing it.

# Automatic Handoff

Before ending, finish the coherent snapshot-binary move or record an exact
resumable blocker. Rewrite the rolling handover, update the roadmap, answer
this archive, create exactly one reciprocal NEXT archive for the separate P6
Docker acceptance move, replace only the launcher's mutable regions, run
launcher and handoff contracts, and make the normal
`docs: prepare next agent session` continuity commit after the focused
implementation commit. Do not launch a successor, push, merge, publish, invoke
an ordinary production release, enable a publisher, stash, revert, or remove
the worktree. P6 remains active until both snapshot and Docker evidence plus
the `make quality` exit gate pass.
<!-- CODEX_SESSION_PROMPT_END -->
