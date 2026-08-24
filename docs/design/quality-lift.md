# Quality Lift Design

Status: accepted for the `codex/upgrade-quality` worktree.

## Problem

The repository had no maintained quality baseline, `make install` no longer
matched the executable layout, tests wrote into committed fixtures, and two
Kibana tests were permanently skipped. The upgrade must improve those facts
without silently changing the CLI or public Go API.

## Design

The lift uses four complementary controls:

1. `.quality/inventory` declares the subjects, seams, user-facing features,
   and intended adapter boundaries.
2. `.quality/tools/quality-audit.sh` produces the authoritative structured
   scorecard and compares measurements with the stored baseline.
3. `internal/testutil` owns repository-write prevention for test fixtures.
4. Small commits are measured from clean checkouts. A checkpoint stops after
   at most three quality moves for review and replanning.

The scorecard is evidence, not a substitute for judgment. A criterion can be
PASS because a ratchet held its baseline while substantial debt remains. The
plan therefore records both the verdict and the underlying number.

## Decisions Not To Reopen

- The target is a staged, honest L2 lift before broad modernization.
- Command-line and public Go API compatibility are preserved.
- Go 1.26 or 1.27 adoption follows L2 characterization and acceptance work.
- Cloud migration to `ply-config` is later work and retains cache-first behavior.
- Spring repair follows the four core CLI surfaces.
- Binary and Docker are active distribution targets.
- Homebrew and Snap are inactive. The default GoReleaser configuration disables
  remote releases and contains no package-manager publishers; supported
  `snapshot` and ordinary `release` Make targets are local-only.

## Core Surfaces

| Surface | Contract to preserve | Artifact evidence |
| --- | --- | --- |
| Host install | `make install` / `go install ./cmd/ply` produces an executable named `ply`. This is distinct from the `ply install` command group. | Fresh host install into a missing temporary `GOBIN`. |
| Status | `ply status` remains read-only and reports dependency state with compatible flags, output, and exit behavior. | Fresh snapshot binary and Docker image against local fixtures. |
| Upgrade | `ply upgrade <scope>` rewrites only the expected local Maven fixture and reports partial failures by content. | Fresh snapshot binary and Docker image against local fixtures and cache. |
| Build | `ply build` (alias `generate`) creates a project with compatible flags and output without requiring public network access in the acceptance path. | Fresh snapshot binary and Docker image using local fixture/cache/loopback inputs. |

The binary distribution artifact means a GoReleaser snapshot after inactive
publishers have been removed. It does not mean an unreviewed production
`make release` invocation.

## Premises

| Premise | Evidence recorded 2026-08-23 | Recheck | Consequence if false |
| --- | --- | --- | --- |
| The executable entry point is `./cmd/ply`. | Real `make install` produced `ply`; `ply --help` showed the command tree. | `make test-install` and install into a temporary `GOBIN`. | Stop release work and repair the entry point contract. |
| Tests can run without user state or repository writes. | Empty-HOME run passed twice; Q0.4, Q0.5, and Q0.6 passed; 0 of 37 tests were skipped. | Run the hermetic and tree-identity gates in the plan. | Treat the checkpoint as a regression. |
| The stored baseline is reproducible. | Quality meta-suite passed 15 of 15 controls, including byte-identical T15 reproduction. | `bash .quality/tools/test-quality-audit.sh`. | Do not use scorecard deltas until the apparatus is repaired. |
| CLI and public Go API compatibility are required. | User decision for this lift; no intentional compatibility break was included in checkpoint 1. | Characterization and acceptance scripts before each refactor. | Split any required break into an explicit migration decision. |
| Binary and Docker are the active distribution targets. | User decision; binary install was exercised. Docker execution was unavailable in the session environment. | Build and smoke-test the Docker image when a daemon is available. | Do not claim Docker support from static inspection alone. |
| Cloud work will use cache-first behavior and later move to `ply-config`. | User decision; real cloud behavior was outside checkpoint 1. | Record fixture and real-boundary evidence before migration. | Keep the existing boundary until equivalent behavior is proven. |
| Spring repair follows the core CLI flows. | User decision; Spring behavior was outside checkpoint 1. | Revisit after install, status, upgrade, and build have acceptance coverage. | Reorder only with a documented production need. |

## Evidence Status

| State | Evidence |
| --- | --- |
| Verified | Host install, status, upgrade, and local build through fresh host artifacts; API/CLI compatibility; command help; local-only distribution configuration and recording contract; uncached, empty-HOME, and race tests; vet, tree identity, and quality meta-tests. |
| Not verified | Docker build/run, GoReleaser snapshot, public network, real cloud, Spring end-to-end behavior. |
| Known and accepted for the next plan | Fifteen untested packages, 127 process exits outside `main`, 80 direct effect sites, zero seam swaps, zero mutation harnesses, and no snapshot/Docker acceptance yet. |
| Resolved at P2B | Remote releases are disabled in `.goreleaser.yml`, Homebrew/Snap publisher sections are absent, the standalone brew config is removed, and `make release-brew` fails closed. |

## Invariants

- Preserve command names, flags, output contracts, and public Go API symbols
  and signatures. Preserve exit behavior unless the Q1.2 exit-only-main rule
  proves incompatible; that case requires an explicit migration decision.
- Add external effects only through one of the five declared adapters.
- Test doubles record arguments, have safe defaults, and fail on empty
  populations where an empty result could create a false green.
- Tests do not write into the worktree and do not depend on a developer HOME.
- A quality move does not regress any comparable ratchet.
- Reports from dirty trees or broken audit tooling are not release evidence.
- Manual receipts are commit-, measured-tree-, inventory-, instrument-, and
  digest-bound. They can resolve only an upstream `UNMEASURABLE` criterion with
  an independently non-empty project population; automated PASS and FAIL
  verdicts retain precedence.
- Homebrew and Snap remain inactive until they receive explicit scope and
  acceptance coverage.

## Boundaries

Checkpoint 1 did not upgrade the declared Go version or dependencies, migrate
cloud configuration, repair Spring behavior, validate real network services,
build a Docker image, or invoke GoReleaser. Those are ordered workstreams in
the upgrade plan.
