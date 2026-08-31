# Quality Upgrade Handover

Generated: 2026-08-31T19:37:43+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P6 quality implementation commit:
  `097a9f15782c773e47a61d903f5d09f5c96408f0`.
- Its exact parent is the launch continuity commit
  `b8b803c807166b00ed2cd257789b2337b41b9eab`, whose exact parent is the
  Docker acceptance implementation `7e3729321fdd28d7d561e5160fa7065387b68b03`.
- The implementation tree is
  `8f9616f7c4efd90a4d346adb0af4156566e1914d`.
- After this handoff, obtain the new continuity HEAD with
  `git rev-parse HEAD`; its exact parent must be `097a9f1`.
- The focused implementation changed only `Makefile` and executable
  `test/makefile_quality_test.sh`. No production Go or Go tests, dependency,
  Dockerfile, GoReleaser configuration, existing acceptance or mutation
  implementation, inventory, audit/parser/scanner/baseline instrument,
  fixture, publisher, registry, or release configuration changed.
- The implementation and every authoritative P6 measurement had empty
  ordinary and ignored status. No push, merge, publication, release, stash,
  revert, successor launch, evidence deletion, image deletion, or worktree
  removal occurred.

## Continuity Checkpoint

P2A-P6 are complete. P7 is active for one bounded maintained-Go toolchain
baseline move; P8 remains queued. The exact Q0-Q2 apparatus exits 0 at clean
L2 with every required population present and no held material debt.

The P6 archive is answered history and links reciprocally to exactly one NEXT
P7 archive. Only the launcher's mutable header and prompt regions changed; its
stable executable skeleton remains byte-identical.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. Audit evidence,
reports, caches, build contexts, and generated artifacts remain external.

## P6 Implementation And Contract

`make quality` now composes the existing entry points instead of duplicating
them. It requires explicit absolute external regular executables for Go, Bash,
GoReleaser, apidiff, golangci-lint, Docker, and Python; an explicit non-empty
external schema-2 document; a fresh external output root; and external Go build
and module caches.

The target validates and runs, in contractual order:

- complete preflight;
- all eight mutation meta-tests and all eight actual mutation harnesses, each
  with exact `declared=10 killed=10 survived=0 unusable=0` receipts;
- all four host acceptance verifiers;
- fresh snapshot and Docker acceptance; and
- the authoritative audit with exact scope `--only Q0.*,Q1.*,Q2.*`.

It rejects missing or duplicate stages, wrong ordering, empty populations,
local or stale evidence/output, skipped artifact populations, audit exit 1 or
2, missing criteria, non-L2 output, dirty measurements, or held, regressed, or
not-comparable scoped debt. Its executable contract passes all fail-closed
controls and is included in `make test`.

## Current Schema-2 Evidence

The independent review root is
`/private/tmp/ply-p6-quality-review.097a9f1.IfKUcH`. The canonical document is
`manual-evidence-schema-2.json`, SHA-256
`1601eaa449b2435908024aa8675445d5d95d84cc59e03db8a4529ad573a55bdb`.
It is bound to commit `097a9f1`, tree `8f9616f`, module, inventory, audit
instruments, every declared mutation subject, every mutation/meta script, and
the exact non-empty four-script verifier population.

Criterion-object SHA-256 values are:

- Q1.6: `393e8e9d4afcc510b81cf59f443512c48efd5ce292762eb02b6f45ae274717b0`;
- Q1.7: `298e43239ca5d81bbaba08f7ad44be5a658885ccf11da83ec571b745b2e0ee01`;
- Q1.9: `730b6217ba8dc8e76253026bca3a84c89361fa2cd579490ba0878e92099600e1`;
- Q2.4: `a437cf14c0c3733abdda988e13621a634368c524fb1c49ff9f38167db7539c7a`;
- Q2.8: `579d83d4b8c04a768b2458ce70d62848680b7f55583a7e5feedc5e55bd7b2668`;
  and
- Q2.9: `0b781b5184a03ee427d2532870102570df8235c614edd1a87eb0b338724127db`.

Fresh observations compare real input/treatment magnitudes in the declared
direction for install, status, upgrade, and build. All four bad-input probes
exit non-zero while creating or changing zero protected artifacts. The focused
six-criterion audit exits 0 with all receipts valid.

## Exact Quality Gate

The authoritative apparatus root is
`/private/tmp/ply-p6-quality-gate.097a9f1-final.lOV5NN`. Its verified
15,849-entry evidence-manifest SHA-256 is
`b89ac6645f814d40d2444ddca3a808919afacf40c873b576a170bc0e5a4091a6`.
The exact 21-line stage ledger ends with `audit:Q0.*,Q1.*,Q2.*`, and
`make quality` exits 0.

The scorecard SHA-256 is
`bfe32a5e7eb90e9c7e9cc756a47f38bfefd68866e754cbb030653d409fce5322`.
All 27 scoped criteria pass; attained level is L2; manual evidence is valid
with six receipts; mutation and acceptance denominators are 8 and 4; held,
regressed, and current-not-comparable counts are zero; and the measured tree is
clean with zero dirty paths.

The gate produced a fresh `darwin/arm64` snapshot artifact, size 19,389,170,
SHA-256
`b740c5e0df160aedeee034718fce416ec38a3012f936121efec341e0dd147a63`.
It also produced and retained fresh local `linux/arm64` image
`sha256:7aa5b59b2116179f848e749d59841faf14a563c3b8bfaf35bfd54e72fd0c42aa`,
entrypoint `/bin/ply`, size 275,721,016. Snapshot and Docker status, upgrade,
and build verifiers all pass; host install remains separate with its artifact
override unset.

## Independent Regression Gate

The post-gate rerun is retained at
`/private/tmp/ply-p6-regression-gate.097a9f1.YGeuN9`. Its verified 15,458-entry
evidence-manifest SHA-256 is
`85fc11883c3ecd9ce444ae23777ef34c36944273d03fbc1f869b17fa477c7de1`.
It retains passing current receipts for snapshot meta 20/20 and fresh
acceptance, Docker meta 26/26 and a second fresh image, API/CLI and
entry/subprocess compatibility, pinned golangci-lint 2.12.2, complete tests,
race, vet, launcher 62/62, Make contracts, complete preflight, all four host
verifiers, audit meta 15/15, and empty-HOME count-2.

The focused six-criterion and exact Q0-Q2 scorecard SHA-256 values are
`441b773a56141b5cbf8f1a20e129a39c2543aded13280895d109ba4929bb2ced`
and `bfe32a5e7eb90e9c7e9cc756a47f38bfefd68866e754cbb030653d409fce5322`;
both exit 0. The full audit exits the expected 1, never 2, and its scorecard
SHA-256 is
`5e12b8b9348a98aa2d18c153bc1e9ae591fd36333ffd24cb0443ec082af08a9c`.
It attains L2 with 32/36 PASS; only queued L3 rows Q3.1 and Q3.4 fail and Q3.3
and Q3.7 remain unmeasurable. The full report is not an exit gate.

Two retained attempts hit the characterized launcher signal-retention timing
fixture. Isolated unchanged reruns passed launcher 62/62, the complete Make
contract, and complete preflight. No product or quality-apparatus failure
remains.

## Next Objective

Begin P7 with one bounded toolchain-baseline move. Determine from current
primary support evidence whether Go 1.26 or 1.27 is the maintained baseline
supported by the pinned quality tools. Inventory every declared toolchain
identity before editing, then migrate the module, build, CI, release, and
documented declarations together only where they are active or contractual.

Keep this move toolchain-only: do not upgrade dependency versions or change
production/API/CLI behavior. Reproduce old and new quality measurements with
both instrument identities, explain any numeric or compatibility drift, refresh
external evidence that becomes commit-bound, and retain all generated state
outside the worktree.

## Start And Stop

Confirm branch, HEAD, exact ancestry, empty ordinary and ignored status,
reciprocal links, launcher `--check`, and the authorized P7 checkpoint before
editing. Read the P7 roadmap, current toolchain declarations, baseline identity
contract, active and deactivated CI, Docker/release inputs, and pinned-tool
support policy. Verify time-sensitive toolchain support from primary sources.

Stop before dependency upgrades, P8 domain modernization, unrelated behavior,
API/CLI changes, publication, publishers, registries, or release. Do not push,
merge, publish, release, delete retained evidence or either retained local
image, stash, revert, launch a successor, or remove the worktree.
