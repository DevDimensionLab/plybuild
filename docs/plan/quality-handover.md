# Quality Upgrade Handover

Generated: 2026-08-24T10:05:25+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `eb987fd8db58`.
- Restart preparation base: `eb987fd8db58`.
- Session head: use `git rev-parse --short=12 HEAD` after launch; the restart
  commit contains this handover and no implementation changes.
- No push or merge was performed.

Earlier focused commits:

```text
3822f4a build: repair make install
674e0a4 quality: establish audited baseline
80b43ba test: make core fixtures hermetic
25689c5 test: replace skipped Kibana integration
1cc3361 docs: establish quality upgrade plan
166bff5 chore: ignore local agent handoffs
3e80994 dev: add restartable Codex launcher
2972d11 quality: pin static analysis gate
e956da7 test: add search-replace meta-test
ab9a5c6 quality: add non-publishing preflight
99cebaa quality: bind manual criterion receipts
4887222 quality: enforce manual evidence precedence
eb987fd quality: migrate manual-evidence baseline
```

## Continuity Checkpoint

`codex-dev-start.sh` is in terminal `COMPLETE` state. The P1B prompt remains
byte-exact inert history below the execution boundary so archive integrity can
still be checked, but normal start and `--print-prompt` fail closed because no
approved next task exists. Its stable executable section precedes
`CODEX_STABLE_EXECUTION_END`; both mutable regions are comment-encoded `#|`
data after an unconditional `exit`.

Startup validates:

- attached branch and repository root;
- required design and plan files;
- strict session ID/archive path grammar and range-checked timestamps whose
  `Created` value matches the ID;
- non-symlink archive directories and files;
- one connected reciprocal archive graph;
- exactly one `NEXT` tail, or no `NEXT` when `COMPLETE`;
- SHA-256 integrity for every historical prompt block;
- byte-equality among decoded launcher prompt, active archive prompt block,
  `--print-prompt`, and the Codex argument, including terminal LF;
- the configured Codex executable.

A dirty worktree is allowed. The launcher emits one static stderr warning and
does not place status output or filenames in the prompt.

`test/codex_dev_start_test.sh` derives session metadata from the launcher, so a
normal restart does not require a test edit. Its 48 controls cover Bash 3.2,
stable-skeleton drift, raw-command injection, source-snapshot behavior, failed
Git status, exact bytes, normal-service override, external-executable
resolution, recording-Codex invocation, traversal, source and runtime symlinks,
hidden archives, timestamp ranges and binding, prompt digests, malformed
metadata, duplicate `NEXT`, disconnected history, a cycle, a nested
second-generation source-archive run without `.git`, missing and non-reciprocal
links, terminal `COMPLETE`, and a hermetic Git fixture. The checked-in graph is
executed only after its source types are validated and it is copied into that
synthetic repository, so Docker and source archives do not depend on the host
worktree's absolute `.git` pointer. Nested
mode is an internal argument, ambient `CDPATH` and Git configuration are
sanitized, Git's repository-local environment is cleared, and both normal and
child runs assert their exact control count.

`make test` runs this target through `/bin/bash`.

The terminal archive is
`docs/plan/agent-sessions/2026-08-24T085458+0200-make-manual-evidence-reachable.md`.
It and its predecessor are answered history; both prompt blocks and digests
remain unchanged. The connected graph has no `NEXT` archive.
The ignored `.agent-task/current.md` source was retired and is absent.

## Measured Quality State

The clean full audit at `eb987fd8db58` reports:

- Absolute L0: 8 of 8.
- Q0.6: 0 skipped tests out of 37 and 0 unsafe direct test writes.
- Q0.8: 0 of 1 production scripts lack a meta-test.
- Ratchets: two improved, five held, zero regressed.
- Packages with tests: 7 of 22; fifteen remain without test files.
- Process-exiting calls outside `main`: 127.
- Direct production effects outside the five declared adapters: 80 of 80.
- Declared seam swap tests: 0 of 8.
- Mutation harnesses: 0 of 8.
- Acceptance scripts: 0 of 4.
- Manual L1/L2 reachability: 6 of 6 rows with synthetic non-empty fixtures;
  the current mutation and acceptance populations remain zero and non-passing.
- Baseline migration: old and new denominator/criterion objects are equal and
  all 228 numeric debt leaves match; Q3.9 is PASS under both instruments.

Correction to the prior handover: its structured 7-of-8 count included Q0.8 as
a held ratchet PASS. The starting raw absolute count was 6 of 8 because Q0.3
and Q0.8 both failed their underlying criteria. P1 closed both rather than
treating the held Q0.8 debt as complete.

Regenerate `target/quality-audit/scorecard.json`; it is ignored output, not
persistent evidence. The full audit may exit 1 for measured findings. Exit 2
invalidates the checkpoint.

## Decisions Learned

Contract controls established through P0A and P1:

1. Mutable data is after the execution boundary, not executable assignments.
2. Tests derive session IDs and archive paths and exercise a second generation.
3. Archive paths have strict basename grammar and the whole chain is checked.
4. SHA-256 selection supports Alpine's `sha256sum` and macOS `shasum`.
5. Git-status failures stop startup, and one launcher-source snapshot supplies
   both validation and the Codex prompt.
6. Calendar/timezone ranges, `Created`/ID equality, terminal links, and prompt
   digests make archive lifecycle drift fail closed.
7. Contract execution uses synthetic Git state and separately starts from a
   no-`.git` source archive; portable `make test` never depends on host worktree
   metadata.
8. Source types are validated before copying, the launcher itself cannot be a
   symlink, internal test mode cannot be selected through ambient state, and
   `CODEX_BIN` must resolve to an external executable file.
9. Startup overrides only `service_tier` to `default`; global Fast mode and all
   other user-profile settings remain untouched.
10. `make lint` resolves an explicit `GOLANGCI_LINT` before `PATH`, verifies
    exactly `2.12.2`, forces offline module resolution, validates the v2 config,
    and never requests fixes.
11. `make preflight` runs build, uncached tests, vet/lint, launcher/lint/install
    contracts, every `scripts/test-*`, and the quality meta-suite. Its contract
    rejects missing, empty, incomplete, orphaned, and symlinked populations.
12. Manual criterion receipts bind module, commit, commit tree, measured status,
    inventory, complete instrument identity, criterion kind, and a canonical
    evidence digest. Criterion-specific numeric/boolean contracts reject a
    receipt whose declared PASS is false.
13. Receipts can alter only upstream `UNMEASURABLE`. Automated PASS and FAIL
    retain precedence, and independently measured test, mutation-harness, or
    acceptance-script populations must be non-zero.

The archive stores the actual prompt argument, not a template with runtime
substitutions. Dirty state therefore cannot make the archive and Codex input
disagree.

The archive graph, launcher regions, and historical prompt files may change
only after the exact restart trigger. During ordinary work, correct false
claims in their owning design/plan/handover document and queue prompt changes
in this handover.

Do not invoke `make release` or `make release-brew`. Active Homebrew
publication remains outside the accepted distribution matrix and is owned by
P2B.

## Next Objective

P1B is complete in three moves. P2 compatibility, distribution, cloud, Spring,
packaging, dependency, and publishing work remains outside this session's
approved scope. There is no approved follow-up objective, so the launcher and
archive chain are terminal. A future task requires explicit user approval and
a newly prepared session rather than replaying this answered prompt.

## Start

There is no next task to start. From any directory, the lifecycle check remains
available:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh --check
```

These commands intentionally fail closed in terminal state:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh --print-prompt
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh
```

Do not start through `.agent-task/current.md`; it is not task authority.

The exact restart trigger was consumed to prepare this terminal transition and
one local allowlisted commit named `docs: prepare next agent session`. It did
not authorize push, merge, release, stash, revert, worktree removal, or staging
unrelated changes.

## Verification Notes

Completed during this checkpoint:

- Terminal restart remeasurement at clean `eb987fd8db58`: `make preflight`
  PASS with the pinned golangci-lint `2.12.2`, zero lint issues, all 48 launcher
  controls, script/install contracts, and all 15 quality controls.
- Focused manual-evidence suite and empty-HOME `go test ./... -count=2`: PASS
  with isolated writable state under `/private/tmp`.
- Full structured audit at `eb987fd8db58`: expected exit 1 for 20 documented
  findings, absolute L0 8 of 8, two ratchets improved, five held, and zero
  regressed. No audit returned exit 2.

- Schema-2 manual-evidence focused tests: PASS. Synthetic non-empty fixtures
  reach all six Q1/L2 rows; stale, dirty, duplicate, empty, wrong-kind,
  wrong-digest, false-PASS, precedence, and zero-population controls all pass.
- Old/new baseline reproduction: both instruments exit 1 only for source
  findings; old and new raw bodies match; scorecards reproduce byte-for-byte;
  denominator and criterion objects and all 228 numeric debt leaves match;
  Q3.9 remains PASS. The old/new scorecard SHA-256 values are
  `d420887d73aabf496ff276fcc55d13ad379ac49c9322fad58808b5e28fdba7df`
  and `5fb3226009cfbf0d29f63fa03592157cce4efcec6e38583b64a86f6288e89490`.
- Expanded `bash .quality/tools/test-quality-audit.sh`: PASS, 15 controls,
  including criterion-bound evidence and dual-instrument T15 reproduction.

- Restart remeasurement at clean `bf189a1c04ea`: `make preflight` PASS after
  isolating both Go and golangci-lint caches under `/private/tmp`.
- Restart full quality audit at `bf189a1c04ea`: expected exit 1, L0 8 of 8,
  20 findings above L0, and zero comparable ratchet regressions.
- Prepared two-generation archive graph: launcher `--check`, exact prompt
  digest/mirror, stable skeleton, and all 48 lifecycle controls PASS.
- Post-transition `make preflight`: PASS, including zero lint issues and all
  15 quality-audit meta-controls; no publication path ran.
- `make preflight`: PASS from clean commit `ab9a5c6`; no publishing or packaging
  path ran and the tree remained clean.
- Pinned golangci-lint `2.12.2`: config verification and `make lint` PASS with
  zero issues. The binary used for measurement was installed outside the repo
  at `/private/tmp/ply-golangci-lint-v2.12.2/golangci-lint`.
- `test/makefile_lint_test.sh`, `scripts/test-search-replace.sh`, and
  `test/makefile_preflight_test.sh`: PASS, including their negative controls.
- `/bin/bash test/codex_dev_start_test.sh`: PASS, 48 controls.
- `make test-agent-start`: PASS, using macOS Bash 3.2.57.
- `make test`: PASS with isolated `GOCACHE` and `GOTMPDIR`.
- `make test-install`, `go test ./... -count=1`,
  `go test -race ./... -count=1`, and `go vet ./...`: PASS.
- Empty-HOME `go test ./... -count=2` with isolated writable state: PASS.
- `bash .quality/tools/test-quality-audit.sh`: PASS, 15 controls.
- Full quality audit: expected exit 1 for 20 findings above L0; absolute L0 is
  8 of 8, Q0.8 is 0 of 1, two ratchets improved, five held, and zero regressed.
  Q3.4 measured zero state-claim phrases.
- Exported `go doc -all` output matched `3e80994` byte-for-byte after the static
  cleanup; no exported Go declaration changed.
- `codex --strict-config -c 'service_tier="default"' --help`: PASS without a
  model request.
- Contract runs poisoned with `CDPATH`, global Git fsmonitor, and `GIT_DIR`:
  PASS, 48 controls each.
- `CODEX_BIN=/usr/bin/true ./codex-dev-start.sh --check`: PASS with the
  expected dirty-tree warning.
- `git diff --check`: PASS.

Environment probes on 2026-08-24:

- `/bin/bash` is 3.2.57; the `bash` on `PATH` is 5.3.9.
- Docker CLI exists, but `docker info` exits 1 because the daemon socket is
  unavailable. No image build was claimed.
- `shellcheck` is unavailable. No shellcheck result was claimed.
- Host Go is 1.26.2 on Darwin arm64; the module declaration remains Go 1.18.
- golangci-lint is still unavailable on `PATH`; use the explicit temporary path
  above if it remains present, or install the official `v2.12.2` binary outside
  the repository before rerunning `make lint` or `make preflight`. Set both
  `GOCACHE` and `GOLANGCI_LINT_CACHE` to fresh `/private/tmp` directories in
  this sandbox; the default shared caches returned operation-not-permitted
  errors during restart remeasurement.

## Stop Conditions

Stop and report rather than forcing progress when:

- the audit exits 2;
- a comparable ratchet regresses;
- public CLI or Go API compatibility cannot be established;
- a change needs cloud, Spring, publication, or packaging scope assigned to a
  later checkpoint;
- a required external tool is missing and no fail-closed evidence can replace
  it;
- three moves in the active checkpoint are complete.
