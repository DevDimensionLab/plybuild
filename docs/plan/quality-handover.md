# Quality Upgrade Handover

Generated: 2026-08-27T05:45:00+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P5.8 implementation commit:
  `58c5224183308e8e48aa4d91b25b65c1e4021a96`.
- Its exact parent is the launch continuity commit
  `5bcd6c48a8ce6a09f0553368e1f5c34e2bf09cef`.
- Measured implementation tree:
  `8cc32082df0fbc4af2f10810c8ecf5274fecbc63`.
- Clean status SHA-256:
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
- After handoff, obtain the new continuity HEAD with `git rev-parse HEAD`; its
  exact parent must be `58c5224`.
- The implementation changes only `scripts/mutate-interactive-build` and
  `scripts/test-mutate-interactive-build`. No production Go, Go test,
  inventory, audit, parser, scanner, baseline, acceptance, API, CLI, packaging,
  or dependency file changed.
- No push, merge, release, publication, distribution, stash, revert,
  successor launch, or worktree removal was performed.

## Continuity Checkpoint

P5 is complete and P6 is active. `codex-dev-start.sh` remains NEXT for the
first bounded P6 snapshot-binary acceptance move. Its active archive is the
single new tail under `docs/plan/agent-sessions`; the `interactive-build`
archive is answered history and links reciprocally to it. The graph has exactly
one NEXT archive. Only the launcher's mutable header and prompt regions change;
the stable execution region remains byte-identical.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. The audit counts
ignored and untracked bytes, so ordinary and ignored status must both be empty
at a measured checkpoint.

## P5.8 Result

The eighth mutation subject, `interactive-build`, and P5 are complete:

- The executable harness declares ten deterministic, unique mutations across
  `cmd/build.go`, `pkg/webservice/init.go`, and `pkg/webservice/api.go`:
  project configuration, cloud configuration, root response, and generation
  endpoint wiring; URI host and asynchronous server start; loopback bind;
  generate route; listen-error gate; and Darwin browser platform selection.
- The inventory remains byte-exact, including
  `8. interactive server binds only to loopback`. Each declaration binds
  production syntax occurring exactly once to a non-empty exact named test
  population. Other production, tests, fixtures, generated/vendor code,
  adapters, completed subjects, equivalent replacements, and known
  non-compiling replacements are excluded.
- One clean unmodified external control ran every exact selection. Each mutant
  used a fresh external Git archive, cache, HOME, config, and temp root; every
  changed package compiled separately; and selected tests ran under
  `go test -json` with exact discovery, run, and terminal-action validation.
  Setup, discovery, selection, compilation, tooling, and unrelated failures
  are unusable; only a selected test failure counts as killed.
- Exact totals are `declared=10`, `killed=10`, `survived=0`, and
  `unusable=0`. No reachability, observability, or controllability gap appeared,
  so no production or test repair was made.
- T1-T10 prove failure for empty/duplicate manifests, unauthorized paths,
  zero/multiple replacements, empty or inexact tests, broken control,
  uncompiled/unexercised mutants, false accounting, unclassified survivors,
  repository-local artifacts, and non-deterministic declarations/totals.

Every mutation and killing population is in
`/private/tmp/ply-interactive-evidence.A49cuF/report.txt`, SHA-256
`109619d5d336a52f01d5f751d3b89b5bcfa85f2ba512cdc5073fbc58e1954d5f`.
Retained control/mutant copies and logs are under
`/private/tmp/ply-interactive-evidence.A49cuF/harness-work`. The independent
T1-T10 meta-log is `/private/tmp/ply-interactive-evidence.A49cuF/meta.log`,
SHA-256
`22a0ac95820514fec641d72289c6ef1a3a56023cc6471f4c2e154087199c9f81`.

## Authoritative Measurement

The refreshed external evidence and audits are under
`/private/tmp/ply-interactive-evidence.A49cuF`. The exact 245 Q1.6, 34 Q1.7,
and 69 Q1.9 named populations were rerun at `58c5224`; every named test emitted
one top-level run and one pass. Manifest SHA-256 values are
`845cc6bdfaa0e4e84956a96a162e84ed0372238f067f8d2b22d05f7a22307f27`,
`cb028c078c2f61283a974ec308a3a0a2de1a6694cd2a60a130b6f4dd7f357d0d`,
and `c6488bd00c342304d8d1fa46a4af6dce33c9b095ec6ad1572bc9e0de40252034`.
Event-log SHA-256 values are
`2a1747ea2bd24e3c136b5a3e7687b32b2a9b570371fa5c133b2864cf2da67d88`,
`36c426ebb1ae611078b78061c0a9b5761060c7c0f68166ee23902aed6911a962`,
and `6610731598fe7180cce3911db8a4ed8e4cca88c3fdbab3f8783f769dbe910a44`.

The external schema-2 document contains refreshed Q1.6, Q1.7, Q1.9 receipts
and Q2.4 receipts for all eight completed P5 subjects. Its SHA-256 is
`b12d4c262ecd16c9529887968cc33324552b08a7cd1d300ff0f1dff6e34acbba`.
Evidence-object SHA-256 values are:

- Q1.6: `94412396d46985af08ca81a649717d7947ea9061d420c89a675d9968d7646e65`;
- Q1.7: `106af984ce687f78a5356e7191747851a9509f40406bba9a507395397cfa4337`;
- Q1.9: `afc77255af92c22c1ed817bcfcf06cac8e43fe9340968afc2a153b64d2bd6df5`;
- Q2.4: `2565dd9824e224520004774521164f0d24ff771b134f7662125a337ae27dd09d`.

The focused Q1.6/Q1.7/Q1.9/Q2.4 audit exits 0 with scorecard SHA-256
`a5478b88f4dbd99dfa5561eafa09117c649a3f3998b79b34818b54b3970d6077`.
The full authoritative audit exits 1, never 2, with scorecard SHA-256
`9032b0e3763c64f820564f1608fd6ef52b0a751311c4dcbcd9a2fca8b1eef972`.
It records L0 8/8, L1 9/9, Q2.1-Q2.4 PASS, seven improved ratchets, one
Q3.4 documentation-phrase regression, zero dirty paths, and seven non-passing
P6-P8 rows. Immutable prompt history was not rewritten.

The no-evidence Q2 view is under
`/private/tmp/ply-interactive-evidence.A49cuF/no-evidence-q2`, scorecard
SHA-256
`4154c02e7520e2021cbcb4ffee777b00d1448c6f451717691d5a5bceb9f8c221`.
It exposes the unratcheted state: Q2.1 is 8/8, Q2.2 and Q2.3 pass, and Q2.4 is
unmeasurable without the external run receipt.

The report and meta-log SHA-256 pairs for all eight subjects are:

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
  `6b932bf00ed9ae0f14152b821ee5cfcd4636c09c739fe0b4d21ef7974b900270`;
- `interactive-build`:
  `109619d5d336a52f01d5f751d3b89b5bcfa85f2ba512cdc5073fbc58e1954d5f`,
  `22a0ac95820514fec641d72289c6ef1a3a56023cc6471f4c2e154087199c9f81`.

## Gate Result

Gate logs and isolated tools/caches are under
`/private/tmp/ply-interactive-gate.58c5224`. API/CLI and explicit
entry/subprocess compatibility, compatibility meta-tests, pinned
golangci-lint 2.12.2 with zero issues, complete uncached tests across all 27
packages, race, vet, `make test`, all 62 launcher controls, Make
distribution/install/lint/preflight contracts, all four host acceptance flows,
exact empty-HOME count-2, the standalone 15-control audit meta-suite, and
complete preflight pass. API and CLI report SHA-256 values remain
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
The complete passing preflight, acceptance, audit-meta, empty-HOME count-2,
entry/subprocess, and `make test` log SHA-256 values are respectively
`a85d99220dc8054c56cad12c34d6bd629f9760e7068ac3451c6485fa211c226f`,
`36b484c628424a75eac04fd1550952ef4b8ecea87123c15e2c6d2b6b02400f00`,
`ae853e4750310bf06d5a580067531fd2a08dd4c869efc04675c9b5fa10272baa`,
`9c1e7a2d0e60897d50363ac1bf58997aadccbbd53d3fc8ba8ebee2eeef0bc7c8`,
`6802980545289c487914d32b2fad16f9e74a9739bf7f49d25f022b1c841d0c78`,
and `57fa7347b37727909495066130b35f88cf2276febc06937c5ad470febc0e98af`.

The first complete preflight hit the known signal-retention timing control
after 25 launcher assertions. No source changed; the standalone launcher suite
had already passed all 62 controls, and the unchanged complete preflight rerun
passed the entire gate. Ordinary and ignored status were empty before
authoritative measurement and after the complete implementation gate.

## Next Objective

Begin P6 with one bounded snapshot-binary acceptance move. Produce a real fresh
GoReleaser snapshot through the credential-cleared, non-publishing path in an
external directory. Revalidate host install as the separate `make install` /
`go install` contract, then run the existing status, upgrade, and build
acceptance behavior against the exact host-platform snapshot executable.
Require exact artifact discovery, provenance, executable identity, non-empty
behavioral verdicts, and falsifiability controls that reject missing, stale,
ambiguous, or substituted artifacts and skipped verifier populations.

Do not start Docker acceptance in the same move. Do not add `make quality`
until both the snapshot and Docker acceptance populations exist and the P6
exit conditions can be measured. Preserve the non-publishing distribution
contract, inactive Homebrew/Snap boundary, CLI/API behavior, Go/dependency
baseline, and all completed P5 evidence.

## Start And Stop

Confirm branch, HEAD, exact ancestry, clean and ignored status, reciprocal
links, launcher `--check`, and the authorized checkpoint block before editing.
Read this handover, the linked NEXT archive, the full P6 and checkpoint-gate
entries, both design documents, the complete GoReleaser/Make distribution
path and its tests, every acceptance verifier/meta-test and shared helper, and
the Q2.5-Q2.10 audit/parser logic before defining the snapshot population.

Stop before Docker work, `make quality`, P7-P8, Go/dependency upgrades,
production/API/CLI behavior, publication, remote releases, package-manager
publishers, or unrelated audit/inventory/P5 changes. Keep all generated
artifacts external. Do not push, merge, stash, revert, launch a successor, or
remove the worktree.
