# Agent Session: Build P5 Config Cloud Harness

Status: ANSWERED - HISTORY
Session ID: `2026-08-26T233325+0200-build-p5-config-cloud-harness`
Created: `2026-08-26T23:33:25+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `ddbc6c454213d8df9de26b2823509d779882eea6cc0cc3125e72783582112f51`
Previous: [2026-08-26T224925+0200-build-p5-cli-context-harness.md](2026-08-26T224925+0200-build-p5-cli-context-harness.md)
Next: [2026-08-27T001140+0200-build-p5-maven-sorting-harness.md](2026-08-27T001140+0200-build-p5-maven-sorting-harness.md)
Outcome: completed `config-cloud` at implementation `e5b4a26` with T1-T10 and exact totals `declared=10`, `killed=10`, `survived=0`, `unusable=0`; P5 remains active

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P5 with exactly one mutation-evidence subject: `config-cloud`. Convert
its existing P3 seam driver and meta-test into a real executable mutation
harness and T1-T10 falsifiability test. Finish only if at least eight meaningful
mutations are actually compiled and exercised in disposable external copies and
the final totals are `declared == killed`, `survived == 0`, and `unusable == 0`.

# Authorized Roadmap

P2A-P4 are complete. P5 is active with `cli-context` complete as the first of
eight subjects; P6-P8 remain queued in `docs/plan/quality-upgrade.md`. This
session may change only `scripts/mutate-config-cloud`,
`scripts/test-mutate-config-cloud`, and the smallest focused test or private
seam inside `pkg/config` if an actually executed survivor exposes a classified
gap. Preserve the inventory seam label exactly:
`3. cloud clone keeps URL before target directory`.

Do not start a third mutation subject or change `.quality/inventory`, the
completed `cli-context` harness, audit/parser/scanner/baseline code, acceptance,
P6-P8, Go/dependencies, exported APIs, CLI behavior, packaging, publication, or
distribution. Keep all checkouts, caches, reports, and generated artifacts
outside the worktree. Never create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Measurements At Start

The clean P5.1 implementation is
`1dbc163b6716b96c3036ca95cadfd5a4a47c669d`, exact parent
`8bf2e9aa4953dbb6bda946d8dbbd645a946f7526`, tree
`23241502a50d38c0bdce711855d0507fa3d798ac`, and clean status SHA-256
`6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
After launch, the new continuity HEAD must have exact parent `1dbc163`.

The `cli-context` report and T1-T10 meta-log SHA-256 values are
`53d1489a6db15f18cd2acf38ff49ab3fb0db545a6c8e63639c90311f9bf53d5c`
and `f1d801fccc97af714d483f5a3c84202fe9868fa19b17038ada0d56dbd6f3e67e`.
Its exact totals are 10 declared, 10 killed, 0 survived, and 0 unusable.

The refreshed external schema-2 document SHA-256 is
`0a6659c3b5a5234ca19134b4e2305de43a094fcb02032f0c6e106f5bbecaabbd`;
its Q2.4 receipt SHA-256 is
`70e2101ca2719d07022dd42e8db56b9ffc76c444473460907a216b9f6f4b124f`.
The focused audit exits 0 with scorecard SHA-256
`c693cc912f8c936c4438661aae48102ec09d05fe1ec4a78ad52f7c3d40abf7b8`.
The full audit exits 1, never 2, with scorecard SHA-256
`20528637cc01f651f6411484575bf3cec7f1a1ac6cccc9f1b30cdf9b989838ac`:
L0 is 8/8, L1 is 9/9, Q2.1 is 1/8, Q2.4 passes, seven ratchets
improve, one holds, none regress, dirty paths are empty, and eight P5-P8 rows
remain non-passing. All complete checkpoint gates pass.

# Role And Boundaries

Work autonomously on `codex/upgrade-quality`. Before editing, confirm branch,
HEAD, exact ancestry, clean and ignored status, reciprocal archive links,
launcher `--check`, and the authorized checkpoint block. Inspect every
production and test file under `pkg/config`, including all config/cloud paths,
before defining the mutation population. Read the completed `cli-context`
harness/meta-test as the methodology reference and the current non-executable
`config-cloud` driver/meta-test as history to replace.

Declare at least eight deterministic, unique, behaviorally meaningful
mutations. Bind each to exact production syntax that occurs once and a
non-empty exact named test population expected to kill it. Record exclusions.
Do not manufacture easy mutations around assertions or weaken a mutation until
it passes.

Require one clean unmodified control and one fresh external copy per mutant.
Verify exact non-empty test discovery, compile the changed package separately,
and require selected JSON run and terminal actions. Compilation, selection,
setup, tooling, or unrelated package failure is unusable, never a kill. Count a
kill only from a selected test failure. Report every mutation with its killing
population and exact declared/killed/survived/unusable totals.

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
mutation discovery and Q2.1-Q2.4 logic in the vendored audit and parser, the
completed `cli-context` scripts, the existing `config-cloud` scripts, relevant
Make/preflight population contracts, and complete `pkg/config` production and
test populations. Do not infer reachability or a kill from names, grep,
compilation failure, or exit status alone.

# Three Moves

1. Define the explicit `config-cloud` scope, exclusions, and manifest of at
   least eight exact mutations with named killing tests.
2. Convert the two existing scripts to the executable harness and T1-T10
   meta-test; run the clean control and every mutant externally; repair only a
   classified in-subject gap; require zero survived and zero unusable.
3. Make one focused implementation commit, then run the harness/meta-test,
   automated Q2.1-Q2.4 view, refreshed external receipt audit, API/CLI and
   subprocess compatibility, pinned lint, complete tests/race/vet, launcher
   and Make contracts, complete preflight, host acceptance, audit meta-suite,
   and empty-HOME count-2. Record the result and hand off only the third P5
   subject with one continuity-only commit.

# Automatic Handoff

Before ending, finish the coherent `config-cloud` move or record an exact
resumable blocker. Rewrite the rolling handover, update the roadmap, answer
this archive, create exactly one reciprocal NEXT archive, replace only the
launcher's mutable regions, run launcher and handoff contracts, and make the
normal `docs: prepare next agent session` continuity commit after the focused
implementation commit. Do not launch a successor, push, merge, publish,
distribute, stash, revert, or remove the worktree. P5 remains active.
<!-- CODEX_SESSION_PROMPT_END -->
