# Quality Upgrade Handover

Generated: 2026-08-24T11:03:54+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `14764fda60f0`.
- Restart preparation base: `14764fda60f0`.
- Session head: use `git rev-parse --short=12 HEAD` after launch; the restart
  commit contains this handover and no implementation changes.
- No push, merge, release, publication, stash, revert, or worktree removal was
  performed.

P2A commits:

```text
5e8232d dev: keep authorized roadmap sessions restartable
224cbcc quality: gate public Go API compatibility
06b2a24 quality: bind Cobra compatibility contract
5f4ba29 quality: add host acceptance skeletons
14764fd quality: keep CLI exporter within exit ratchet
```

## Continuity Checkpoint

`codex-dev-start.sh` is in `NEXT` state for P2B. P2A is complete; P2B is active
and P3-P8 remain queued in the machine-readable plan block. The active archive
is
`docs/plan/agent-sessions/2026-08-24T110354+0200-make-distribution-local.md`.
Its predecessor is answered history, and the connected graph has exactly one
`NEXT` tail.

The launcher validates its attached branch/root, the authorized queue, strict
session metadata, non-symlink planning inputs, a reciprocal archive graph,
every historical prompt digest, exact launcher/archive prompt bytes, and the
external Codex executable. Mutable prompt data remains comment-encoded after
the stable execution boundary. A dirty worktree is allowed with a static
warning; no status content enters the prompt.

`test/codex_dev_start_test.sh` derives active metadata from the launcher. Its 49
controls include premature `COMPLETE` rejection while authorized work remains,
so normal handoff does not require a test edit. The ignored `.agent-task` path
is absent and is not task authority.

## P2A Result

P2A completed its three measured moves:

1. The official `golang.org/x/exp/cmd/apidiff` is pinned externally at
   `v0.0.0-20260709172345-9ea1abe57597`. The v1.0.1 comparison permits only the
   compatible `cmd/ply` package addition and rejects every incompatible delta.
2. `test/compat/cli-v1.0.1.json` stores the complete normalized Cobra tree.
   The explicit CLI allowlist is empty because the current tree has zero
   deltas. Subprocess contracts preserve root, status, upgrade, build, and
   unknown-command help/output/exit behavior.
3. Four host verifiers build or install fresh artifacts and exercise install,
   read-only status, targeted dependency upgrade, and local build using copied
   fixtures, isolated HOME/cache state, and loopback Maven metadata. Their four
   matching meta-tests run 26 controls that kill no-op artifacts, missing or
   wrong output, broad writes, read-only violations, help failures, and
   bad-input exit/artefact regressions.

The first clean full audit exposed one gate regression from `os.Exit(1)` in the
build-ignored CLI exporter. Commit `14764fd` replaced that explicit exit with a
panic on the impossible encoder failure. CLI compatibility stayed green and
Q1.2 returned from 128 to its held baseline of 127.

## Measured Quality State

The clean full audit at `14764fda60f0` reports:

- Absolute L0: 8 of 8.
- 37 test functions, zero skipped; 7 of 22 packages have tests.
- Q0.8: 0 of 7 production scripts lack a meta-test.
- Ratchets: two improved, five held, zero regressed.
- Process-exiting calls outside `main`: 127.
- Direct production effects outside the five declared adapters: 80 of 80.
- Declared seam swap tests: 0 of 8.
- Mutation harnesses: 0 of 8.
- Acceptance scripts: 4 of 4; Q2.5, Q2.6, Q2.7, and Q2.10 pass.
- Q2.8 and Q2.9 remain honestly `UNMEASURABLE` pending criterion-bound manual
  evidence; the executable magnitude, bad-input, and read-only controls exist.
- Full audit exit: expected 1 for 16 documented non-passing criteria; no run
  returned exit 2.

Regenerate `target/quality-audit/scorecard.json`; it is ignored output, not
persistent evidence. The full audit may exit 1 for measured findings. Exit 2
invalidates the checkpoint.

## Decisions And Learned Facts

1. Compatibility tools must stay outside `go.mod`; their current toolchain
   requirements are newer than the module's retained Go 1.18 declaration.
2. Machine-readable API and CLI allowlists are the review boundary. Raw output
   alone is not accepted evidence.
3. The quality scanner includes build-ignored Go helpers, so those helpers must
   not introduce explicit process-exit calls outside the product boundary.
4. An empty status target is accepted today and exits 0 with a no-project-type
   warning. The negative status fixture instead uses a `src/` directory without
   language configuration, which reliably exits non-zero without writes.
5. Acceptance checks assert exact file-set and content boundaries: status is
   byte-identical, upgrade changes only `pom.xml`, and build adds only
   `ply.json` while preserving every original non-POM file.
6. `make preflight` now discovers seven production scripts and runs all seven
   matching meta-tests. Real acceptance remains the explicit `make acceptance`
   target until P6 adds snapshot and Docker coverage to the single quality
   entry point.
7. Docker CLI exists but the daemon is unavailable. No image evidence is
   claimed. GoReleaser was not invoked during P2A.
8. `.goreleaser.yml` still has active GitHub release and Homebrew publisher
   configuration, `.goreleaser.brews.yml` remains active, `make release` can
   publish, and `make release-brew` directly invokes it. This contradiction is
   the entire P2B scope.
9. Never invoke `make release` or `make release-brew` while those paths remain
   active.

## Next Objective

P2B is one measured move: make default distribution local-only. Begin with a
negative contract test, remove inactive Homebrew/Snap publishers from the
default path, disable or explicitly guard the standalone brew path, and add a
snapshot command whose recorded GoReleaser invocation cannot publish GitHub
releases or package-manager metadata. Keep production release execution out of
the test; use a recording executable and static configuration checks.

Do not broaden P2B into Go/dependency upgrades, Docker acceptance, cloud,
Spring, process-exit refactoring, adapters, or mutation work. Those remain in
P3-P8.

## Start

The active P2B task can be started from any directory:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh
```

Do not start through `.agent-task/current.md`; it is not task authority. Every
session prepares its successor automatically before stopping.

## Verification Notes

Completed at the P2A boundary:

- `make acceptance`: PASS for install, status, upgrade, and build.
- Four `scripts/test-verify-*`: PASS, 26 controls total.
- `make preflight` with explicit apidiff and golangci-lint `2.12.2`: PASS.
- `make test`, `make test-install`, `make test-agent-start`, uncached tests,
  race tests, and `go vet ./...`: PASS.
- Empty-HOME `go test ./... -count=2` with isolated writable state: PASS.
- `bash .quality/tools/test-quality-audit.sh`: PASS, 15 controls.
- Full audit at `14764fda60f0`: expected exit 1, L0 8 of 8, two improved,
  five held, zero regressed ratchets.
- API compatibility, CLI compatibility, their meta-tests, and CLI subprocess
  surface contracts: PASS after the exporter correction.
- The compatibility binaries used were
  `/private/tmp/ply-p2a-api.9mgcdg/bin/apidiff` and
  `/private/tmp/ply-golangci-lint-v2.12.2/golangci-lint`; probe them before
  reuse because temporary paths are not persistent dependencies.

Environment on 2026-08-24: host Go 1.26.2 on Darwin arm64, module Go 1.18,
`/bin/bash` 3.2.57, PATH Bash 5.3.9, Docker daemon unavailable, and shellcheck
unavailable.

## Stop Conditions

Stop and report rather than forcing progress when:

- the audit exits 2;
- a comparable ratchet regresses;
- public CLI or Go API compatibility cannot be established;
- a safe local-only snapshot cannot be proven without running a publisher;
- P2B would require Docker, cloud, Spring, dependency, or packaging scope; or
- the one measured P2B move is complete.
