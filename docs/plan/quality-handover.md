# Quality Upgrade Handover

Generated: 2026-08-27T04:46:26+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P5.7 implementation commit:
  `5829939035cf6b0b5dd2c5ea49c5b4120a40dcb7`.
- Its exact parent is the launch continuity commit
  `0ae8531f83e5c5edadffdbfc1c71cf6f518fa542`.
- Measured implementation tree:
  `ab4f2290a4fa0950abeb7ad259814679908fdeaf`.
- Clean status SHA-256:
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
- After handoff, obtain the new continuity HEAD with `git rev-parse HEAD`; its
  exact parent must be `5829939`.
- The implementation changes only `scripts/mutate-http` and
  `scripts/test-mutate-http`. No production Go, Go test, inventory, audit,
  parser, scanner, baseline, acceptance, API, CLI, packaging, or dependency
  file changed.
- No push, merge, release, publication, distribution, stash, revert,
  successor launch, or worktree removal was performed.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P5 is active and P6-P8 are queued.
Its active archive is
`docs/plan/agent-sessions/2026-08-27T044626+0200-build-p5-interactive-build-harness.md`.
The `http` archive is answered history and links reciprocally to the new tail.
The graph has exactly one NEXT archive. Only the launcher's mutable header and
prompt regions changed; the stable execution region is unchanged.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. The audit counts
ignored and untracked bytes, so ordinary and ignored status must both be empty
at a measured checkpoint.

## P5.7 Result

The seventh of eight mutation subjects, `http`, is complete:

- The executable harness declares ten deterministic, unique mutations in
  `pkg/http/client.go`: anonymous JSON request URL; response dependency and
  status gates; response-body source; basic-auth credential order; bearer URL
  component order and token selection; Wget request URL and destination path;
  and Wpost request URL.
- The inventory remains byte-exact. The `http` subject has no separate seam
  label.
- Each declaration binds production syntax occurring exactly once to a
  non-empty exact named `pkg/http` test population. Other production packages,
  tests, fixtures, generated/vendor code, completed subjects, equivalent
  replacements, and known non-compiling replacements are excluded.
- One clean unmodified external control ran every exact selection. Each mutant
  used a fresh external Git archive, cache, HOME, config, and temp root; the
  changed package compiled separately; and its selected population ran under
  `go test -json` with exact run and terminal-action validation. Tooling,
  setup, discovery, selection, compilation, and unrelated failures are
  unusable; only a selected test failure counts as killed.
- Exact totals are `declared=10`, `killed=10`, `survived=0`, and
  `unusable=0`. No reachability, observability, or controllability gap appeared,
  so no production or test repair was made.
- T1-T10 prove failure for empty/duplicate manifests, unauthorized paths,
  zero/multiple replacements, empty or inexact tests, broken control,
  uncompiled/unexercised mutants, false accounting, unclassified survivors,
  repository-local artifacts, and non-deterministic declarations/totals.

Every mutation and killing population is in the retained report
`/private/tmp/ply-http-evidence.N97pVm/report.txt`; its SHA-256 is
`865d29bee9baad2053e603e188d8ccdd5688ecbcbe053979f586d535898645c7`.
All retained control/mutant checkouts and logs are under
`/private/tmp/ply-http-evidence.N97pVm/harness-work`. The independent T1-T10
meta-log is `/private/tmp/ply-http-evidence.N97pVm/meta.log`, SHA-256
`6b932bf00ed9ae0f14152b821ee5cfcd4636c09c739fe0b4d21ef7974b900270`.

## Authoritative Measurement

The refreshed external evidence and audits are under
`/private/tmp/ply-http-evidence.N97pVm`. The exact 245 Q1.6, 34 Q1.7, and 69
Q1.9 named populations were rerun at `5829939`; every named test emitted one
run and one pass. Manifest SHA-256 values are
`845cc6bdfaa0e4e84956a96a162e84ed0372238f067f8d2b22d05f7a22307f27`,
`cb028c078c2f61283a974ec308a3a0a2de1a6694cd2a60a130b6f4dd7f357d0d`,
and `c6488bd00c342304d8d1fa46a4af6dce33c9b095ec6ad1572bc9e0de40252034`.
Event-log SHA-256 values are
`483056e858ef43190c69735655d2c582b932e328f5780b1e521436068cfc05d6`,
`894cacaec604899884a6c45eb5a85932c5159087339eb224bff5d50031394e55`,
and `9a7bd71ed085e9dab2fd9d1fe5982f895577c1ea17f7c468b9ece8c66ab21833`.

The external schema-2 document contains refreshed Q1.6, Q1.7, Q1.9 receipts
and Q2.4 receipts for all seven completed P5 subjects. Its SHA-256 is
`c4f5d9c8c97aa653afd1b4ef8095f1e1bd35873bc40acc5ab9e6d9f308544713`.
Evidence-object SHA-256 values are:

- Q1.6: `e2809553c057a07e0fddf1bf95c6104354590e0c897856e2c78c1c539f06a446`;
- Q1.7: `5c0cd66a2e0715b23a4f297375de298d1a0ec092eaf8578d0b3c87523b41bdbc`;
- Q1.9: `86442a0855506f139300b19eb8a9d8ff61520ad8aa811759255c5fa0d1273f08`;
- Q2.4: `6a462402ad3de49c3f14b3907664e04edd39af7b360421b6497a170a1a620de8`.

The focused Q1.6/Q1.7/Q1.9/Q2.4 audit exits 0 with scorecard SHA-256
`a9034e35989bf4d1508edd5deda4b59fc4ffc34172db12052996a27c39d56383`.
The full authoritative audit exits 1, never 2, with scorecard SHA-256
`c3d8cebe64fd4aebc10596852ec9ece963937bc079858e9fab0300a41d8fb968`.
It records L0 8/8, all nine L1 rows PASS, Q2.1 improved to 7/8, Q2.4 PASS,
seven improved ratchets, one Q3.4 regression, and zero dirty paths. The three
Q3.4 findings are literal prompt-history phrases in the tracked rolling
handover and byte-exact active archive at measurement time. The implementation
introduced no product or apparatus regression, and immutable prompt history
was not rewritten. Nine P5-P8 rows remain non-passing.

The separate no-evidence Q2.1-Q2.4 view is under
`/private/tmp/ply-http-evidence.N97pVm/no-evidence-q2`, scorecard SHA-256
`46a5070d9538576053e5b8889fac01b87f90b1bcfa9e5bc9326c7ab334a98ed2`.
It exposes the exact unratcheted state: Q2.1 is 7/8; Q2.2 and Q2.3 identify
only the legacy `interactive-build` recorder; Q2.4 is unmeasurable without the
external run receipt.

The report and meta-log SHA-256 pairs for the seven completed subjects are:

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

## Gate Result

Gate logs and external caches are under `/private/tmp/ply-http-gate.gzWxRY`.
API/CLI and explicit entry/subprocess compatibility, compatibility meta-tests,
pinned golangci-lint 2.12.2 with zero issues, complete uncached tests across all
27 packages, race, vet, `make test`, all 62 launcher controls, Make
distribution/install/lint/preflight contracts, all four host acceptance flows,
exact empty-HOME count-2, the standalone 15-control audit meta-suite, and
complete preflight pass. API and CLI report SHA-256 values remain
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
The complete passing preflight, acceptance, audit-meta, empty-HOME count-2,
focused entry/subprocess, and `make test` log SHA-256 values are respectively
`b277ee920f8f5bb1ba54b97587035bb8fcb2855ba11d98469e9b12560bbc843b`,
`36b484c628424a75eac04fd1550952ef4b8ecea87123c15e2c6d2b6b02400f00`,
`ae853e4750310bf06d5a580067531fd2a08dd4c869efc04675c9b5fa10272baa`,
`b21e978e6f9e6343dc89959bf8c3854ecf7d5561ca6e4c2be941de773a3a968f`,
`1efb43c081f7d055e1879ce89efadabb0dfb181583483f82304ec4260fbd00c4`,
and `f1518049dd46ecd91a666b9b68ab93259e4befe68e31b6359863f574a52ff5c1`.

The first complete preflight hit the previously recorded final signal-fixture
partial-log timing miss after its first 25 launcher assertions. No source
changed; the standalone launcher suite had already passed all 62 controls, and
the unchanged complete preflight rerun passed the entire gate. Ordinary and
ignored status were empty before authoritative measurement and after the
complete implementation gate.

## Next Objective

Convert exactly the eighth inventory subject, `interactive-build`, from its P3
recording driver into a regular executable mutation harness and T1-T10
falsifiability test. Define at least eight meaningful mutations after
inspecting the complete `cmd/build.go` and `pkg/webservice` production and test
populations plus every adapter used by the interactive server path. Preserve
the inventory seam label byte-exact:
`8. interactive server binds only to loopback`.

Use all seven completed P5 harnesses as methodology references. Require one
clean external control, fresh external copies and caches, exact test discovery,
successful mutant compilation, selected JSON run/terminal actions, and exact
totals `declared == killed`, `survived == 0`, `unusable == 0`. If a genuine
survivor appears, classify it as reachability, observability, or controllability
before the smallest in-subject test or private seam repair. Do not start P6.
P5 becomes complete only after this eighth subject and its full gate pass.

## Start And Stop

Confirm branch, HEAD, exact ancestry, clean and ignored status, reciprocal
links, launcher `--check`, and the authorized checkpoint block before editing.
Read the rolling handover, linked NEXT archive, complete P5 entry/gate, both
designs, `.quality/README.md`, `.quality/inventory`, complete mutation audit and
parser logic, all seven completed harness/meta-test pairs, the full legacy
interactive-build driver/meta-test history, Make/preflight contracts, complete
`cmd/build.go` and `pkg/webservice` production and test populations, and the
complete adapter populations used by that path before selecting mutations.

Stop before inventory/audit/parser/scanner/baseline changes, acceptance
expansion, P6-P8, Go/dependency upgrades, exported API/CLI changes, packaging,
publication, or distribution. Keep generated artifacts external. Do not push,
merge, stash, revert, launch a successor, or remove the worktree.
