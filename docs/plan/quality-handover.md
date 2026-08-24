# Quality Upgrade Handover

Generated: 2026-08-24T22:13:50+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `1b85711b139c`.
- Restart preparation base: `1b85711b139c`.
- Session head: use `git rev-parse --short=12 HEAD` after launch; the restart
  commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits remain:

```text
5c6f2fa quality: move process exits to main
03d6242 quality: move git processes behind adapter
f5ee37d quality: move Maven processes behind adapter
ee5e9ab quality: cover cloud clone seam
ffc4e77 quality: move Maven metadata HTTP behind adapter
204e222 quality: move template copy behind filesystem adapter
89d0f76 quality: move Spring download behind adapters
dee214c quality: move Spring discovery behind HTTP adapter
789ae23 quality: move Bitbucket JSON behind HTTP adapter
07ac6ce quality: move Kibana POST behind HTTP adapter
a7eb3ef quality: move Wpost behind HTTP and filesystem adapters
```

The separate operational continuity implementation is:

```text
1b85711 quality: supervise non-interactive agent sessions
```

It is not P3.12 and changes no Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains `NEXT` while P3 is active and P4-P8 remain queued
in the machine-readable plan block. Its active archive is
`docs/plan/agent-sessions/2026-08-24T221350+0200-migrate-bitbucket-clone-selection.md`.
Its predecessor is answered history, and the complete reciprocal graph has
exactly one `NEXT` tail.

Normal launch is now a Bash 3.2-compatible, non-interactive supervisor. Each
generation resolves an external Codex executable and invokes exact
`codex exec` arguments for normal service tier, workspace-write, the repository
working directory, JSONL, an external final-message path, and the byte-exact
validated prompt. It does not use interactive mode or resume a thread.

Every turn receives a unique physical log directory outside the worktree. Raw
JSONL is preserved byte-for-byte while concise progress is streamed. The
parser decodes JSON, treats messages only as data, requires the characterized
thread/turn lifecycle and one final `turn.completed`, and rejects empty,
malformed, truncated, contradictory, failed, error, post-terminal, and
non-zero-exit streams.

After a successful child exit, the parent re-reads the on-disk launcher and
validates branch/root, queue, full archive graph, prompt, and normalized stable
skeleton. A next turn requires a clean changed HEAD, an answered old archive,
a changed session ID, exactly one new committed `NEXT` archive, and reciprocal
history. A valid `COMPLETE` requires the whole authorized queue complete.
Signals are forwarded to the active child/parser and never start a successor.

`test/codex_dev_start_test.sh` now has 62 controls. They characterize exact
argv/prompt bytes, logging/progress, malicious event data, two-generation
continuation, `COMPLETE`, all terminal/JSON/process/no-progress/dirty/handoff/
contract failures, and signal interruption without using the real Codex
executable. Mutable header and prompt data remain inert after the stable
execution boundary; the pinned normalized skeleton digest is
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.

## P3 Move 11 Preserved

Exported `http.Wpost`, every caller, and all P2A API/CLI/subprocess contracts
remain unchanged. Wpost routes its complete form POST through the HTTP adapter
with explicit default-client selection and routes create/copy through the
filesystem adapter while retaining request-before-create, file-before-response
close order, ignored close errors, exact logging, and no status rejection.
Its seven contracts remain part of the 130-test suite. No inventory entry,
seam driver, or mutation label changed; Q1.4 remains 7 of 8 and exact Q2.1
remains 0 of 8.

## Measured Quality State

The clean full audit at `1b85711b139c` reports:

- Absolute L0: 8 of 8.
- 130 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 18 guarded safe-writer sites, 13 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 58 direct external sites outside declared adapters of 68 production
  effect sites. Clock and server remain absent, so it is not comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Acceptance scripts: 4 of 4; Q2.5, Q2.6, Q2.7, and Q2.10 pass.
- Full audit: expected exit 1 for 16 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single current non-comparable ratchet.
- Measurement identity: clean at tree
  `9328c0d008527fc18d9dd01dd3877b17316e5d4f`, with zero dirty paths.

Reports were written under `/private/tmp/ply-supervisor-gate.WSL8yu`; no
ignored quality or compatibility artifact entered the measured tree.

## Decisions And Learned Facts

1. Supervisor success is the conjunction of process exit, structured terminal
   stream, and committed repository evidence. Final prose is observable only.
2. The local `codex-cli 0.149.0` event contract used for the implementation is
   characterized by the recording executable; official documentation alone is
   not treated as executable truth.
3. JSONL must be parsed rather than substring-matched. Raw bytes are written
   before decode, and event strings are never evaluated as shell content.
4. Logs and final messages stay outside the repository so observation cannot
   dirty or change task authority.
5. The parent supervisor, not the agent, owns succession. It starts only after
   the child ends and the fresh on-disk handoff passes every progression check.
6. A function-valued effect dependency makes Q1.3 fail closed. Production
   moves use resolvable interfaces and complete dependency values.
7. `filesystem.Exists` preserves the relevant legacy rule: only an
   `os.IsNotExist` error means absent. The Bitbucket boundary must additionally
   prevent its zero value from reaching a real Git operation.
8. The remaining Bitbucket `os.Stat` site decides clone versus pull; Git clone
   and pull already execute through the process adapter and are outside the
   next move.
9. `.quality/inventory` is baseline-checksum-bound. Do not relabel seams or
   claim P5 mutation coverage.
10. Provide `APIDIFF` and `GOLANGCI_LINT` as environment variables for
    `make preflight`; Make command-line values propagate through `MAKEFLAGS`
    and defeat the missing-binary mutant.

## Next Objective

Move only Bitbucket repository-existence selection behind the existing
filesystem adapter. Start with recording contracts in
`pkg/bitbucket/bitbucket_contract_test.go`: require the complete repository
path and complete dependencies, prove missing selects clone, existing and
other stat errors select pull, propagate clone/pull failures, preserve safe
zero-value behavior, and reject empty recorded populations. Characterize all
observable path, log, and ordering behavior needed to prevent drift.

Give the private `Bitbucket` flow a complete, resolvable filesystem dependency
whose production value is selected by `With` and `QueryRepos` as needed. Keep
exported `With`, `QueryRepos`, `SynchronizeAllRepos`, project/repository query
behavior, lowercase selection, warnings, and clone/pull behavior unchanged.
Remove only the direct `os.Stat` from `pkg/bitbucket/bitbucket.go`; do not move
Git execution, other filesystem/HTTP effects, config, clock/server, P4, P5, or
later roadmap work.

Expected direction is one fewer Q1.3 violation and one fewer production effect
site (nominally 57 of 67), but regenerate the exact structured measurement and
accept it only with zero comparable ratchet regressions. Q1.2, Q1.4, and exact
Q2.1 should remain unchanged.

## Verification Notes

Completed from clean implementation commit `1b85711b139c`:

- Red launcher evidence: the recording contract expected the 11-argument
  non-interactive invocation and failed against the old five-argument
  interactive launcher.
- `/bin/bash test/codex_dev_start_test.sh`: PASS, 62 controls.
- `make test-agent-start`, `make test-preflight`, and complete `make preflight`:
  PASS, including 15 audit meta-controls.
- API/CLI and subprocess compatibility: PASS.
- `make test`, `make test-install`, uncached tests, race tests, and `go vet`:
  PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS.
- Focused audit: expected exit 1 with exact start values and no regression.
- Clean full audit: expected exit 1, 16 documented findings, L0 8 of 8, five
  improved, two held, zero regressed, and zero dirty paths.

Tool paths used were `/private/tmp/ply-p2b-api.4umBuM/bin/apidiff` and
`/private/tmp/ply-p2b-lint.SGWVGp/bin/golangci-lint`; probe before reuse.

Environment: host Go 1.26.2 on Darwin arm64, module Go 1.18, `/bin/bash`
3.2.57, PATH Bash 5.3.9, and `codex-cli 0.149.0`.

## Start

From any directory:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh
```

This now starts the non-interactive supervisor and may run successive fresh
missions after valid handoffs. Do not invoke it while merely validating the
handoff; use `--check` or `--print-prompt`. `.agent-task/current.md` is not task
authority.

## Stop Conditions

Stop and report rather than forcing progress when:

- Bitbucket selection behavior or exported API cannot be preserved;
- the filesystem dependency is incomplete, unresolvable, function-valued, or
  unsafe at its zero value;
- the effect cannot be isolated without moving Git execution or another flow;
- a comparable ratchet regresses, the audit exits 2, or the tree cannot be
  measured cleanly;
- the work requires P4-P8 implementation, publication, distribution, or a real
  successor launch; or
- the focused Bitbucket move and its separate automatic handoff are complete.
