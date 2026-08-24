# Quality Upgrade Handover

Generated: 2026-08-24T21:24:08+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `a7eb3ef400f7`.
- Restart preparation base: `a7eb3ef400f7`.
- Session head: use `git rev-parse --short=12 HEAD` after launch; the restart
  commit contains this handover and no implementation changes.
- No push, merge, release, publication, stash, revert, or worktree removal was
  performed.

P3 implementation commits:

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

## Continuity Checkpoint

`codex-dev-start.sh` remains in `NEXT` state while P3 is active and P4-P8
remain authorized in the machine-readable plan block. Eleven P3 moves are
complete. At the user's request, the next task is an operational continuity
move: convert normal launcher execution from one interactive session into a
non-interactive, observable supervisor loop. After that focused launcher move,
resume the remaining P3 production-effect migrations.

The active archive is
`docs/plan/agent-sessions/2026-08-24T212408+0200-supervise-noninteractive-sessions.md`.
Its predecessor is answered history, and the connected graph has exactly one
`NEXT` tail. The launcher is still interactive and one-shot at this handoff;
the active task must characterize and change that behavior rather than assume
it already exists.

The current launcher validates its attached branch/root, authorized queue,
strict session metadata, non-symlink planning inputs, reciprocal archive graph,
historical prompt digests, exact launcher/archive prompt bytes, and external
Codex executable. Mutable prompt data remains comment-encoded after the stable
execution boundary. A dirty worktree produces a static warning; no status
content enters the prompt.

`test/codex_dev_start_test.sh` derives active metadata from the launcher. Its
49 controls include premature `COMPLETE` rejection while authorized work
remains, so ordinary handoff requires no test edit. The next task necessarily
changes the stable launcher skeleton and must update its pinned digest and
recording contracts. The ignored `.agent-task` path is absent and is not task
authority.

## P3 Move 11 Result

The eleventh P3 move completed as one measured implementation change:

1. Exported
   `http.Wpost(downloadUrl, filePath string, formData url.Values) error` and
   every caller remain unchanged. Production now supplies complete
   `httpclient.Dependencies` and `filesystem.Dependencies` values to a
   private, zero-value-safe Wpost boundary.
2. The existing narrow `httpclient.POST` value gained only an explicit
   `UseDefaultClient` selection. Wpost passes the complete URL,
   `[]byte(formData.Encode())`, POST method, and exact
   `application/x-www-form-urlencoded` header through the adapter while
   retaining standard `http.PostForm` use of the mutable
   `http.DefaultClient`. Existing anonymous, basic-auth, bearer JSON, and
   Kibana POST paths retain their previous client selection.
3. Wpost still logs the exact debug message before requesting, receives the
   response before creating/truncating the file, rejects no status, copies the
   body through the filesystem adapter, defers file close before response close
   so the file closes first, and ignores both close errors.
4. Seven new top-level contracts raise the suite from 123 to 130 tests. They
   prove complete request and dependency delivery, form encoding and content
   type, default-client redirects and custom redirect behavior, debug logging,
   request/create/copy errors after recording, non-success copying, exact
   bytes, create/truncate, close order, ignored close errors, safe defaults,
   system behavior, and rejection of empty recorded populations.
5. No inventory entry, seam driver, or mutation label changed. Q1.4 therefore
   remains 7 of 8 and exact Q2.1 remains 0 of 8. Exactly Wpost's request,
   create, close, and copy sites left Q1.3.

## Measured Quality State

The clean full audit at `a7eb3ef400f7` reports:

- Absolute L0: 8 of 8.
- 130 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 18 guarded safe-writer call sites, 13 write and 5 copy, with zero
  unsafe direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 58 direct external sites outside the declared adapters of 68
  production effect sites. Clock and server remain absent, so Q1.3 is not
  comparable. The four direct Wpost effects alone left the population.
- Q1.4: 7 of 8 declared seams covered: `git-process`, `maven-process`,
  `cloud-clone`, `maven-http`, `spring-download`, `template-copy`, and
  `git-commit`.
- Exact Q2.1: 0 of 8 subjects have a real executable harness. The upstream
  filename-only denominator sees five non-executable `mutate-*` paths; the
  exact-path validator records all five as non-executable.
- Acceptance scripts: 4 of 4; Q2.5, Q2.6, Q2.7, and Q2.10 pass.
- Full audit exit: expected 1 for 16 documented findings; no accepted audit
  returned exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single current non-comparable ratchet while two declared adapter paths are
  absent.
- Measurement identity: clean at tree `ad2996070ce4`, with zero dirty paths.

Reports were written under `/private/tmp`; no ignored quality or compatibility
artifact entered the measured tree.

## Decisions And Learned Facts

1. A function-valued effect dependency makes the type-aware Q1.3 scanner fail
   closed. Use resolvable interfaces and complete dependency values rather than
   stored functions.
2. Test recorders must write only through recognized in-memory concrete types
   or guarded `internal/testutil` helpers. The new system file assertion adds
   one recognized safe-writer call site, producing the Q0.6 count of 18.
3. `.quality/inventory` is baseline-checksum-bound. Seam labels stay at their
   declared driver paths; changing them invalidates comparison.
4. Non-executable seam drivers keep exact Q2.1 at 0 of 8, while the pinned
   upstream tool counts their filenames. Q2.2/Q2.3 remain findings until P5.
5. `http.PostForm` delegates to `http.DefaultClient.PostForm`; preserving it
   is observably different from allocating a fresh `http.Client` when callers
   replace the global client's redirect policy. `UseDefaultClient` is
   therefore part of the complete Wpost request value.
6. `httpclient.POST` still distinguishes POST presence even for an empty body.
   Only Wpost selects the default client; Kibana retains a fresh client.
7. Response ownership remains in Wpost. The HTTP adapter executes the request,
   and the filesystem adapter creates and copies, but neither duplicates or
   changes the legacy close/error lifecycle.
8. Removing Wpost's direct `http.PostForm`, `os.Create`, deferred close, and
   `io.Copy` calls removes exactly four violations and four production sites,
   yielding Q1.3 at 58 of 68.
9. The local `codex-cli 0.149.0` exposes `codex exec --json`,
   `--sandbox workspace-write`, `-C`, and `--output-last-message`. The
   [official non-interactive-mode documentation](https://learn.chatgpt.com/docs/non-interactive-mode)
   describes JSONL progress events including completed and failed turns. The
   next launcher must characterize the actual local event stream in recording
   tests before trusting it.
10. Supervisor continuation must be based on committed repository evidence,
    not merely an agent final message: changed HEAD, changed session ID, an
    answered predecessor, one valid NEXT tail, a clean worktree, and a passing
    launcher contract. `COMPLETE` remains valid only after P8.
11. Stream concise event summaries to the terminal and preserve raw JSONL logs
    outside the repository so observability cannot dirty the measured tree.
    Stop on malformed JSON, `turn.failed`, non-zero Codex exit, contract drift,
    dirty/no-progress state, or signal.
12. Provide `APIDIFF` and `GOLANGCI_LINT` as environment variables for
    `make preflight`; Make command-line values propagate through `MAKEFLAGS`
    and defeat the lint meta-test's missing-binary mutant.
13. Isolated `GOTMPDIR` must be created before invoking Go. The first
    `make test` verification attempt failed only because the fresh path did
    not yet exist; the corrected rerun and all product contracts passed.

## Next Objective

Implement one focused, contract-first launcher change. Normal
`codex-dev-start.sh` execution should become a non-interactive supervisor
that invokes exactly one archived mission at a time with
`codex exec --json --sandbox workspace-write` in normal service tier, streams
useful progress, and stores raw logs outside the worktree. Preserve
`--check`, `--print-prompt`, help, branch/root/archive/prompt validation,
Bash 3.2 support, exact prompt bytes, inert mutable regions, and the no-publish
boundary.

After a successful turn, re-read and validate the on-disk launcher and archive
graph. Continue only when the agent process has ended and a clean committed
handoff advances both HEAD and session ID to one new valid `NEXT` tail. Stop
successfully when the validated queue and launcher reach `COMPLETE`. Stop
fail-closed on a failed/error/malformed event, non-zero exit, signal, dirty
tree, unchanged state, invalid handoff, or launcher-contract failure. Agent
final text may be logged but must not override repository evidence.

Begin red by extending `test/codex_dev_start_test.sh` with a recording Codex
executable and JSONL scenarios for exact argv/prompt, progress and raw logging,
successful multi-generation continuation, eventual COMPLETE, turn/non-zero
failure, malformed input, no progress, dirty or invalid handoff, and signal
cleanup. Update `docs/design/agent-session-continuity.md` with the measured
protocol. Make one implementation commit, run the full applicable gate, then
prepare the next P3 session in the usual separate handoff commit.

Do not launch a real next session while developing or finalizing this change.
Do not broaden into the remaining Q1.3 effects, adapter work, mutation
harnesses, P4-P8 implementation, distribution, or publication.

## Start

The active continuity task can be started from any directory:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh
```

This invocation is still interactive and one-shot until the active task lands.
Do not start through `.agent-task/current.md`; it is not task authority.

## Verification Notes

Completed from clean implementation commit `a7eb3ef400f7`:

- Red evidence: focused Wpost/HTTP-adapter tests failed only on the absent
  default-client POST value and private complete-dependency boundary.
- Focused HTTP/httpclient/filesystem plus adjacent Spring, Kibana, Maven, and
  Bitbucket contracts: PASS.
- Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1 focused audit: expected exit 1, 18 guarded
  writes, zero unsafe test writes, 9 untested packages, 0 exits, 58 of 68
  effects, 7 of 8 seams, and exact 0 of 8 executable harnesses.
- API and CLI compatibility plus root/status/upgrade/build and unknown-command
  subprocess surfaces: PASS.
- Host install, status, upgrade, and build acceptance: PASS.
- `make preflight`, `make test`, `make test-install`,
  `make test-agent-start`, uncached tests, race tests, and `go vet ./...`:
  PASS.
- Empty-HOME `go test ./... -count=2` with isolated writable state: PASS.
- Audit meta-suite: PASS, 15 controls and all 228 baseline numeric leaves.
- Clean full audit: expected exit 1, L0 8 of 8, 130 tests, 16 tested packages,
  five improved, two held, zero regressed, one not comparable, and zero dirty
  paths.

Tool paths used were `/private/tmp/ply-p2b-api.4umBuM/bin/apidiff` and
`/private/tmp/ply-p2b-lint.SGWVGp/bin/golangci-lint`; probe before reuse.

Environment on 2026-08-24: host Go 1.26.2 on Darwin arm64, module Go 1.18,
`/bin/bash` 3.2.57, PATH Bash 5.3.9, and `codex-cli 0.149.0`. Docker, public
network, real cloud, real Spring, and distribution execution were outside the
Wpost move.

## Stop Conditions

Stop and report rather than forcing progress when:

- a non-interactive Codex invocation cannot preserve exact prompt bytes,
  workspace-write scope, normal service tier, or actionable progress logs;
- success/continuation cannot be distinguished from failure or no progress
  using both the JSONL stream and committed repository state;
- a loop can replay the same task, continue from dirty/invalid state, or mark
  COMPLETE while authorized checkpoints remain;
- Bash 3.2, archive-history, dirty-worktree, source-archive, or inert-tail
  contracts cannot be preserved;
- the change requires launching a real successor, publishing, distributing, or
  broadening into production-effect or later-roadmap work; or
- the focused launcher-supervisor implementation and its separate automatic
  handoff are complete.
