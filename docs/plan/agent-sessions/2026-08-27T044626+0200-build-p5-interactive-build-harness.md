# Agent Session: Build P5 Interactive-Build Harness

Status: NEXT
Session ID: `2026-08-27T044626+0200-build-p5-interactive-build-harness`
Created: `2026-08-27T04:46:26+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `fcf7ab60e29c92c85cb4344eebaa35dae7cda31c8b933feb5bdfc31a041201f7`
Previous: [2026-08-27T040108+0200-build-p5-http-harness.md](2026-08-27T040108+0200-build-p5-http-harness.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P5 with exactly one mutation-evidence subject: `interactive-build`.
Convert its P3 recording driver and meta-test into a regular executable mutation
harness and T1-T10 falsifiability test. Finish only if at least eight meaningful
mutations are actually compiled and exercised in disposable external copies and
the final totals are `declared == killed`, `survived == 0`, and `unusable == 0`.

# Authorized Roadmap

P2A-P4 are complete. P5 is active with `cli-context`, `config-cloud`,
`maven-sorting`, `template`, `file-shell`, `spring`, and `http` complete as the
first seven of eight subjects; P6-P8 remain queued in
`docs/plan/quality-upgrade.md`. This session may change only
`scripts/mutate-interactive-build`, `scripts/test-mutate-interactive-build`, and
the smallest focused test or private seam inside `cmd/build.go` or
`pkg/webservice` if an actually executed survivor exposes a classified gap.
Preserve the inventory seam label byte-exact:
`8. interactive server binds only to loopback`.

Do not start P6 or change `.quality/inventory`, any completed P5 harness,
audit/parser/scanner/baseline code, acceptance, P6-P8 implementation,
Go/dependencies, exported APIs, CLI behavior, packaging, publication, or
distribution. Keep all checkouts, caches, reports, and generated artifacts
outside the worktree. Never create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Measurements At Start

The clean P5.7 implementation is
`5829939035cf6b0b5dd2c5ea49c5b4120a40dcb7`, exact parent
`0ae8531f83e5c5edadffdbfc1c71cf6f518fa542`, tree
`ab4f2290a4fa0950abeb7ad259814679908fdeaf`, and clean status SHA-256
`6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
After launch, the new continuity HEAD must have exact parent `5829939`.

The report and T1-T10 meta-log SHA-256 pairs for the seven completed subjects
are:

- `cli-context`: `53d1489a6db15f18cd2acf38ff49ab3fb0db545a6c8e63639c90311f9bf53d5c`,
  `f1d801fccc97af714d483f5a3c84202fe9868fa19b17038ada0d56dbd6f3e67e`;
- `config-cloud`: `db76fdf8c624c4326483ae71fa9ec3e7e0f94de4d6d3ae8a4dff185c9a7f23f7`,
  `12b1197d521681b70dea0b481f6b7d8bb9daba2ca21dc68a635868177d43ffec`;
- `maven-sorting`: `6b360f604802c047b4946b474c35b2860d47bda4581b3cb33f5af45653dc111e`,
  `fda3185a71ebd842f3e924ae507ad0daa9d59d4b30da3699d114c0ff0d2d53b3`;
- `template`: `e769cb49b2ae6fd3da15b207acab8ef11c402f5bff171dc6e150623b93279e95`,
  `cde7a149274b7f4b968a40b1ec5e4f7769c5d6d4243edcb40783a4fa19ba91bd`;
- `file-shell`: `ea12b122e6210f2b2c8c64a7792ce74a9886e2273e9ddac6176df6552afd1c29`,
  `759ad7a28b3d3ab06d42d42e8e1679e0df516d6abdfd90f3fbac1e92e52bfaad`;
- `spring`: `0aa2bdcb32cf20372f8a7b5c232dcefc93668259eb4ed76c45f8ca339b1dbf25`,
  `f373190fa04abba73b61453b1693f892aced845d931ce267b3674ea84961fa06`;
- `http`: `865d29bee9baad2053e603e188d8ccdd5688ecbcbe053979f586d535898645c7`,
  `6b932bf00ed9ae0f14152b821ee5cfcd4636c09c739fe0b4d21ef7974b900270`.

Each completed subject has exact totals 10 declared, 10 killed, 0 survived,
and 0 unusable.

The refreshed external schema-2 document SHA-256 is
`c4f5d9c8c97aa653afd1b4ef8095f1e1bd35873bc40acc5ab9e6d9f308544713`;
its Q2.4 receipt SHA-256 is
`6a462402ad3de49c3f14b3907664e04edd39af7b360421b6497a170a1a620de8`.
The focused audit exits 0 with scorecard SHA-256
`a9034e35989bf4d1508edd5deda4b59fc4ffc34172db12052996a27c39d56383`.
The full audit exits 1, never 2, with scorecard SHA-256
`c3d8cebe64fd4aebc10596852ec9ece963937bc079858e9fab0300a41d8fb968`:
L0 is 8/8, L1 is 9/9, Q2.1 is 7/8, Q2.4 passes, seven ratchets improve,
the Q3.4 documentation-phrase ratchet records one regression from required
literal prompt-history text, dirty paths are empty, and nine P5-P8 rows remain
non-passing. All complete implementation checkpoint gates pass.

# Role And Boundaries

Work autonomously on `codex/upgrade-quality`. Before editing, confirm branch,
HEAD, exact ancestry, clean and ignored status, reciprocal archive links,
launcher `--check`, and the authorized checkpoint block. Inspect every
production and test file under `cmd` that covers `cmd/build.go`, every file
under `pkg/webservice`, and the complete server, HTTP-client, filesystem, and
process adapters used by the interactive build path before defining the
mutation population. Read all seven completed P5 harness/meta-test pairs as the
methodology reference. Read the complete legacy interactive-build recording
driver and meta-test history before replacing that pair.

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
seven completed P5 script pairs, the legacy interactive-build driver/meta-test
history, relevant Make/preflight population contracts, the complete
`cmd/build.go` and `pkg/webservice` production and test population, and the
complete adapter production and test populations used by the interactive path.
Do not infer reachability or a kill from names, grep, compilation failure, or
exit status alone.

# Three Moves

1. Define the explicit `interactive-build` scope, exclusions, and manifest of
   at least eight exact mutations with named killing tests.
2. Convert the executable harness and T1-T10 meta-test; run the clean control
   and every mutant externally; repair only a classified in-subject gap;
   require zero survived and zero unusable.
3. Make one focused implementation commit, then run the harness/meta-test,
   automated Q2.1-Q2.4 view, refreshed external receipt audit, API/CLI and
   subprocess compatibility, pinned lint, complete tests/race/vet, launcher
   and Make contracts, complete preflight, host acceptance, audit meta-suite,
   and empty-HOME count-2. Record the result, mark P5 complete, and hand off P6
   without implementing it, using one continuity-only commit.

# Automatic Handoff

Before ending, finish the coherent `interactive-build` move or record an exact
resumable blocker. Rewrite the rolling handover, update the roadmap, answer
this archive, create exactly one reciprocal NEXT archive for P6, replace only
the launcher's mutable regions, run launcher and handoff contracts, and make
the normal `docs: prepare next agent session` continuity commit after the
focused implementation commit. Do not launch a successor, push, merge,
publish, distribute, stash, revert, or remove the worktree. P5 becomes complete
only if all eight subjects and the checkpoint gate pass.
<!-- CODEX_SESSION_PROMPT_END -->
