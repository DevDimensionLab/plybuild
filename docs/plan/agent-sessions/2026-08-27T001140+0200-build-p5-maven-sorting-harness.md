# Agent Session: Build P5 Maven Sorting Harness

Status: NEXT
Session ID: `2026-08-27T001140+0200-build-p5-maven-sorting-harness`
Created: `2026-08-27T00:11:40+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `6610270cc624f0eca268c46c3f1268d4064c3340ab5c6058ac32e1e0b6311f64`
Previous: [2026-08-26T233325+0200-build-p5-config-cloud-harness.md](2026-08-26T233325+0200-build-p5-config-cloud-harness.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P5 with exactly one mutation-evidence subject: `maven-sorting`. Convert
its existing P3 seam driver and meta-test into a real executable mutation
harness and T1-T10 falsifiability test. Finish only if at least eight meaningful
mutations are actually compiled and exercised in disposable external copies and
the final totals are `declared == killed`, `survived == 0`, and `unusable == 0`.

# Authorized Roadmap

P2A-P4 are complete. P5 is active with `cli-context` and `config-cloud`
complete as the first two of eight subjects; P6-P8 remain queued in
`docs/plan/quality-upgrade.md`. This session may change only
`scripts/mutate-maven-sorting`, `scripts/test-mutate-maven-sorting`, and the
smallest focused test or private seam inside `pkg/maven` or `pkg/sorting` if an
actually executed survivor exposes a classified gap. Preserve both inventory
seam labels exactly:
`2. Maven command keeps executable before arguments`
`4. Maven metadata keeps username before password`.

Do not start a fourth mutation subject or change `.quality/inventory`, either
completed P5 harness, audit/parser/scanner/baseline code, acceptance, P6-P8,
Go/dependencies, exported APIs, CLI behavior, packaging, publication, or
distribution. Keep all checkouts, caches, reports, and generated artifacts
outside the worktree. Never create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Measurements At Start

The clean P5.2 implementation is
`e5b4a26db271ac43d514bc0e9b2e19c23cdbc144`, exact parent
`be33bb9fa9e6d9831f824c831d9de113123a710b`, tree
`baa63ed3c4d44f063854a4fbee63728cf336c937`, and clean status SHA-256
`6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
After launch, the new continuity HEAD must have exact parent `e5b4a26`.

The `cli-context` report and T1-T10 meta-log SHA-256 values are
`53d1489a6db15f18cd2acf38ff49ab3fb0db545a6c8e63639c90311f9bf53d5c`
and `f1d801fccc97af714d483f5a3c84202fe9868fa19b17038ada0d56dbd6f3e67e`.
The `config-cloud` report and T1-T10 meta-log SHA-256 values are
`db76fdf8c624c4326483ae71fa9ec3e7e0f94de4d6d3ae8a4dff185c9a7f23f7`
and `12b1197d521681b70dea0b481f6b7d8bb9daba2ca21dc68a635868177d43ffec`.
Each subject has exact totals 10 declared, 10 killed, 0 survived, and 0
unusable.

The refreshed external schema-2 document SHA-256 is
`7665b6d6d15da8dbd6b2da3e20a45902cb4e1076683fb0e1acd52c17b85c4648`;
its Q2.4 receipt SHA-256 is
`496ae9cd500b896080a14874dd5ff34bc0bdf521a401cbb7044e5f87604184b5`.
The focused audit exits 0 with scorecard SHA-256
`7ef785e2a17b5371358cc18741a78fe7cbca65c587725970970710d42a7a6b0c`.
The full audit exits 1, never 2, with scorecard SHA-256
`4a7089ea117a97bdf5b265f3712c534954417e00ff352b24a5874a2d710bf648`:
L0 is 8/8, L1 is 9/9, Q2.1 is 2/8, Q2.4 passes, seven ratchets
improve, one holds, none regress, dirty paths are empty, and eight P5-P8 rows
remain non-passing. All complete checkpoint gates pass.

# Role And Boundaries

Work autonomously on `codex/upgrade-quality`. Before editing, confirm branch,
HEAD, exact ancestry, clean and ignored status, reciprocal archive links,
launcher `--check`, and the authorized checkpoint block. Inspect every
production and test file under `pkg/maven` and `pkg/sorting` before defining the
mutation population. Read both completed P5 harness/meta-test pairs as the
methodology reference and the current non-executable `maven-sorting`
driver/meta-test as history to replace.

Declare at least eight deterministic, unique, behaviorally meaningful
mutations. Bind each to exact production syntax that occurs once and a
non-empty exact named test population expected to kill it. Record exclusions.
Do not manufacture easy mutations around assertions or weaken a mutation until
it passes.

Require one clean unmodified control and one fresh external copy per mutant.
Verify exact non-empty test discovery, compile every changed package
separately, and require selected JSON run and terminal actions. Compilation,
selection, setup, tooling, or unrelated package failure is unusable, never a
kill. Count a kill only from a selected test failure. Report every mutation
with its killing population and exact declared/killed/survived/unusable totals.

Implement T1-T10 and fail closed on empty/duplicate manifests, unauthorized
paths, zero/multiple replacements, empty or inexact test selection, broken
control, uncompiled/unexercised mutants, false accounting, an unclassified
survivor, repository-local artifacts, and non-deterministic declarations or
totals. If a real survivor appears, first classify it as reachability,
observability, or controllability; then make only the smallest in-subject test
or private seam repair and rerun the complete harness.

# Required Reading

Read the rolling handover, this archive, complete P5 entry and checkpoint gate,
both design documents, `.quality/README.md`, `.quality/inventory`, complete
mutation discovery and Q2.1-Q2.4 logic in the vendored audit and parser, both
completed P5 script pairs, the existing `maven-sorting` scripts, relevant
Make/preflight population contracts, and complete `pkg/maven` and `pkg/sorting`
production and test populations. Do not infer reachability or a kill from
names, grep, compilation failure, or exit status alone.

# Three Moves

1. Define the explicit `maven-sorting` scope, exclusions, and manifest of at
   least eight exact mutations with named killing tests.
2. Convert the two existing scripts to the executable harness and T1-T10
   meta-test; run the clean control and every mutant externally; repair only a
   classified in-subject gap; require zero survived and zero unusable.
3. Make one focused implementation commit, then run the harness/meta-test,
   automated Q2.1-Q2.4 view, refreshed external receipt audit, API/CLI and
   subprocess compatibility, pinned lint, complete tests/race/vet, launcher
   and Make contracts, complete preflight, host acceptance, audit meta-suite,
   and empty-HOME count-2. Record the result and hand off only the fourth P5
   subject with one continuity-only commit.

# Automatic Handoff

Before ending, finish the coherent `maven-sorting` move or record an exact
resumable blocker. Rewrite the rolling handover, update the roadmap, answer
this archive, create exactly one reciprocal NEXT archive, replace only the
launcher's mutable regions, run launcher and handoff contracts, and make the
normal `docs: prepare next agent session` continuity commit after the focused
implementation commit. Do not launch a successor, push, merge, publish,
distribute, stash, revert, or remove the worktree. P5 remains active.
<!-- CODEX_SESSION_PROMPT_END -->
