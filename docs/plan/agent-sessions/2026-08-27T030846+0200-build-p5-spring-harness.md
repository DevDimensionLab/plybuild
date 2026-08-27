# Agent Session: Build P5 Spring Harness

Status: NEXT
Session ID: `2026-08-27T030846+0200-build-p5-spring-harness`
Created: `2026-08-27T03:08:46+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `a02ef002468a1ac189d0175c2be589ee095e293cefc0c156746b8457c8e84368`
Previous: [2026-08-27T021552+0200-build-p5-file-shell-harness.md](2026-08-27T021552+0200-build-p5-file-shell-harness.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P5 with exactly one mutation-evidence subject: `spring`. Convert its
existing P3 seam driver and meta-test into a real executable mutation harness
and T1-T10 falsifiability test. Finish only if at least eight meaningful
mutations are actually compiled and exercised in disposable external copies and
the final totals are `declared == killed`, `survived == 0`, and `unusable == 0`.

# Authorized Roadmap

P2A-P4 are complete. P5 is active with `cli-context`, `config-cloud`,
`maven-sorting`, `template`, and `file-shell` complete as the first five of
eight subjects; P6-P8 remain queued in `docs/plan/quality-upgrade.md`. This
session may change only `scripts/mutate-spring`,
`scripts/test-mutate-spring`, and the smallest focused test or private seam
inside `pkg/spring` if an actually executed survivor exposes a classified gap.
Preserve the inventory seam label exactly:
`5. Spring download keeps URL before archive path`.

Do not start a seventh mutation subject or change `.quality/inventory`, any
completed P5 harness, audit/parser/scanner/baseline code, acceptance, P6-P8,
Go/dependencies, exported APIs, CLI behavior, packaging, publication, or
distribution. Keep all checkouts, caches, reports, and generated artifacts
outside the worktree. Never create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Measurements At Start

The clean P5.5 implementation is
`7bb94acf47e5971bf331f09e3e2e17ff94e137f6`, exact parent
`f575f315bdee668c62ed217f67f4d0cf0a455c1f`, tree
`00a6db27697d337df9cb6e1bba59f3f791fea923`, and clean status SHA-256
`6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
After launch, the new continuity HEAD must have exact parent `7bb94ac`.

The `cli-context` report and T1-T10 meta-log SHA-256 values are
`53d1489a6db15f18cd2acf38ff49ab3fb0db545a6c8e63639c90311f9bf53d5c`
and `f1d801fccc97af714d483f5a3c84202fe9868fa19b17038ada0d56dbd6f3e67e`.
The `config-cloud` values are
`db76fdf8c624c4326483ae71fa9ec3e7e0f94de4d6d3ae8a4dff185c9a7f23f7`
and `12b1197d521681b70dea0b481f6b7d8bb9daba2ca21dc68a635868177d43ffec`.
The `maven-sorting` values are
`6b360f604802c047b4946b474c35b2860d47bda4581b3cb33f5af45653dc111e`
and `fda3185a71ebd842f3e924ae507ad0daa9d59d4b30da3699d114c0ff0d2d53b3`.
The `template` values are
`e769cb49b2ae6fd3da15b207acab8ef11c402f5bff171dc6e150623b93279e95`
and `cde7a149274b7f4b968a40b1ec5e4f7769c5d6d4243edcb40783a4fa19ba91bd`.
The `file-shell` values are
`ea12b122e6210f2b2c8c64a7792ce74a9886e2273e9ddac6176df6552afd1c29`
and `759ad7a28b3d3ab06d42d42e8e1679e0df516d6abdfd90f3fbac1e92e52bfaad`.
Each subject has exact totals 10 declared, 10 killed, 0 survived, and 0
unusable.

The refreshed external schema-2 document SHA-256 is
`c5ec79633274ad618b7c78568e38891b9955352e020053c8b3d987ad6348fd5e`;
its Q2.4 receipt SHA-256 is
`732083b2add74a30bf5bed70ef9446682bbba20ebb4d1dfc89d0b904db81b922`.
The focused audit exits 0 with scorecard SHA-256
`5261b60d5e20338c5f1a88870b60c55e54411200a7fa84a0d80791dd96597001`.
The full audit exits 1, never 2, with scorecard SHA-256
`7a42d1f40ec96641fce8dab1188a32a73a2f6f790cedea60e98e5b05b1833502`:
L0 is 8/8, L1 is 9/9, Q2.1 is 5/8, Q2.4 passes, seven ratchets improve, one
holds, none regress, dirty paths are empty, and eight P5-P8 rows remain
non-passing. All complete checkpoint gates pass.

# Role And Boundaries

Work autonomously on `codex/upgrade-quality`. Before editing, confirm branch,
HEAD, exact ancestry, clean and ignored status, reciprocal archive links,
launcher `--check`, and the authorized checkpoint block. Inspect every
production and test file under `pkg/spring`, plus the HTTP client and filesystem
adapters used by its download and archive seams, before defining the mutation
population. Read all five completed P5 harness/meta-test pairs as the
methodology reference and the current non-executable `spring` driver/meta-test
as history to replace.

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
mutation discovery and Q2.1-Q2.4 logic in the vendored audit and parser, all
five completed P5 script pairs, the existing `spring` scripts, relevant
Make/preflight population contracts, the HTTP client and filesystem adapters,
and the complete `pkg/spring` production and test population. Do not infer
reachability or a kill from names, grep, compilation failure, or exit status
alone.

# Three Moves

1. Define the explicit `spring` scope, exclusions, and manifest of at least
   eight exact mutations with named killing tests.
2. Convert the two existing scripts to the executable harness and T1-T10
   meta-test; run the clean control and every mutant externally; repair only a
   classified in-subject gap; require zero survived and zero unusable.
3. Make one focused implementation commit, then run the harness/meta-test,
   automated Q2.1-Q2.4 view, refreshed external receipt audit, API/CLI and
   subprocess compatibility, pinned lint, complete tests/race/vet, launcher
   and Make contracts, complete preflight, host acceptance, audit meta-suite,
   and empty-HOME count-2. Record the result and hand off only the seventh P5
   subject with one continuity-only commit.

# Automatic Handoff

Before ending, finish the coherent `spring` move or record an exact resumable
blocker. Rewrite the rolling handover, update the roadmap, answer this archive,
create exactly one reciprocal NEXT archive, replace only the launcher's mutable
regions, run launcher and handoff contracts, and make the normal
`docs: prepare next agent session` continuity commit after the focused
implementation commit. Do not launch a successor, push, merge, publish,
distribute, stash, revert, or remove the worktree. P5 remains active.
<!-- CODEX_SESSION_PROMPT_END -->
