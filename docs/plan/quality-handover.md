# Quality Upgrade Handover

Generated: 2026-08-23T21:47:35+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Measured implementation head: `25689c5`
- Documentation head: use `git rev-parse --short HEAD` after checkout.
- Base: `master` at `5635d50`
- Both worktrees were clean at checkpoint close.
- No push or merge was performed.

## Completed Work

```text
3822f4a build: repair make install
674e0a4 quality: establish audited baseline
80b43ba test: make core fixtures hermetic
25689c5 test: replace skipped Kibana integration
```

Key outcomes:

- `make install` delegates to `go install ./cmd/ply` and has a fail-closed shell
  contract test.
- A real binary installed into a missing `GOBIN` containing spaces and returned
  the expected help output.
- The quality meta-suite passes 15 of 15 controls.
- The full clean-tree audit reports L0 7 of 8, one improved ratchet, six held,
  and zero regressions.
- Q0.6 reports 0 skipped tests out of 37 and 0 unsafe direct test writes.
- `make test`, uncached tests, empty-HOME count-2 tests, race, vet, and the
  install contract passed.

The ignored report is `target/quality-audit/scorecard.json`. Regenerate it; do
not treat an old ignored file as evidence.

## Important Interpretation

Several PASS verdicts are held ratchets, not completed work:

- 15 of 22 packages have no test files.
- 127 process-exiting calls remain outside `main`.
- 0 of 8 seam swaps have executable coverage.
- 0 of 8 mutation harnesses exist.
- 1 repository shell script lacks an independent negative meta-test.

Q1.3 is an explicit FAIL: 80 of 80 measured production effect sites sit outside
the five declared adapters. Q2.5 is an explicit FAIL: all four acceptance
scripts are absent.

The apparatus also emits permanent UNMEASURABLE rows for Q1.6, Q1.7, Q1.9,
Q2.4, Q2.8, and Q2.9. P1B must add commit/tree-bound structured evidence before
an honest L1 or L2 claim is technically possible.

Do not invoke `make release` or `make release-brew`: `.goreleaser.yml`,
`.goreleaser.brews.yml`, and the Makefile still expose Homebrew publication
paths that contradict the inactive Homebrew decision.

## Next Objective

Execute P1 from `docs/plan/quality-upgrade.md`: close absolute L0 and establish
the daily preflight gate in no more than three measured moves.

Recommended first slice:

1. Inspect the supported linter/toolchain matrix without changing dependencies.
2. Add a pinned linter configuration and focused negative fixture.
3. Make `make lint` non-mutating and add a separate `make format` command.
4. Drive scripts without independent meta-tests to zero.
5. Add `make preflight`; do not add `make verify` or `make quality` in P1.
6. Run focused tests, the full gate, and a clean-commit audit.

Do not combine the Go/dependency upgrade with this slice.

## Start

```sh
cd /Users/perottochristensen/github/ply/upgrade-quality
git status --short --branch
sed -n '1,260p' docs/design/quality-lift.md
sed -n '1,320p' docs/plan/quality-upgrade.md
```

For a Codex session, run: `Use $agent-task-handoff`.

After reading the local task file, move or remove `.agent-task/current.md`
before a local audit. Regenerate evidence from a clean commit instead of reading
the ignored target report:

```sh
clone="$(mktemp -d /private/tmp/ply-handoff-audit.XXXXXX)"
git clone --shared --no-checkout . "$clone"
git -C "$clone" checkout --detach HEAD
bash "$clone/.quality/tools/quality-audit.sh" "$clone" \
  --baseline "$clone/.quality/baseline/scorecard.json"
test "$?" -eq 1
python3 -m json.tool "$clone/target/quality-audit/scorecard.json"
```

## Verification Notes

- Use `LC_ALL=C LANG=C` for deterministic shell tooling in this environment.
- Use a fresh `GOCACHE` under `/private/tmp`; the shared cache produced a
  sandbox access error during one install attempt.
- Docker daemon access was unavailable, so no daemon-backed image build was
  recorded.
- `shellcheck` was unavailable. Shell syntax and repository meta-tests passed.
- Real cloud, Spring, and external-service behavior was outside checkpoint 1.
- The audit is exact-toolchain-bound. Parser, scanner, inventory, or Go context
  changes require an explicit baseline migration that preserves measured debt.

## Stop Conditions

Stop and report rather than forcing progress when:

- the audit exits 2;
- a comparable ratchet regresses;
- a compatibility expectation cannot be established;
- the proposed change requires cloud, Spring, or packaging scope from a later
  checkpoint;
- three quality moves have been completed.
