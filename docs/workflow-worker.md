# Run one Goal without an open terminal

On macOS, `workflow execute --headless` starts one durable Codex worker for an
already registered, queued Goal. The Goal must explicitly authorize automatic
local Epic integration without a human QA or integration gate. Existing
interactive runs retain their contracts and cannot be switched to this mode.

From the registered receiving Epic worktree:

```sh
ply workflow execute --headless --spec SPEC_ID --check --format json
ply workflow execute --headless --spec SPEC_ID --format json
```

The preview creates no Task, worker or provider configuration. Starting reserves
the existing Goal, creates its Task checkout and one native Run, and detaches a
small supervisor. Repeating the same start inspects the same work. The worker
finishes this Goal only; it does not select the next queue item.

Use `workflow week` for an explicitly selected set of dependent Goals and a
named login service; see [durable multi-Spec operation](workflow-week.md).

The default is Codex with `gpt-6-astra` and `max` effort; explicit Goal choices
win. It uses the existing Codex account login and does not introduce an API-key
fallback. The installed CLI must support the frozen command options, including
`exec --ignore-user-config --ignore-rules`, JSONL, output schemas and managed
permissions. Runtime probes must pass before development. Requested model and
effort remain distinct from model facts actually exposed by the provider.

## Observe and control the same Run

Use the Run ID returned by start and the same compatible private Ply binary:

```sh
ply workflow worker show RUN_ID --workspace /absolute/workspace --format json
ply workflow worker pause RUN_ID --workspace /absolute/workspace
ply workflow worker stop RUN_ID --workspace /absolute/workspace
ply workflow worker resume RUN_ID --workspace /absolute/workspace
```

`show` is read-only. It reports stage, developer attempt, controller generation,
pending control request, blocked reason and a current lock observation. A saved
controller PID alone does not prove that a process is alive. Inventory and trace
retain native Task/Run history without inventing a Herdr tab or provider session.

Pause and stop finish the active attempt before parking, retain all work and
can both be resumed. They do not kill a model halfway through a command. A pause
or stop request remains visible while the active attempt holds ownership. The
worker never claims a safe pause from a missing or uncertain process result.
An unknown attempt remains blocked on resume; it is not permission to start a
second writer. Inspect the preserved attempt, streams and native receipts.

The supervisor survives closing the launching terminal. If the controller dies,
its trusted attempt runner or active Git process keeps the inherited owner lock;
the supervisor waits and starts the next controller generation only afterward.
The next generation consumes the same immutable attempt result. This version is
not an installed login service: after host reboot or loss of the supervisor,
run `worker resume`. It still requires positive completion or native recovery
evidence. Keep the frozen control and Codex executables: replacing either blocks
dependent work instead of silently upgrading a Run.

## Development, tests and delivery

Each developer uses a fresh actual Codex session. The native controller commits
the finished Task files, runs acceptance, and starts a distinct read-only Codex
reviewer against the exact candidate and test evidence. Concrete test or review
failures return to a developer with preserved feedback and a regression-before-
fix instruction. A passing review must have no findings and an actual completed
provider turn distinct from every developer session.

The developer implements the Task's committed `test/delivery_acceptance.sh` and
builds the executable at the private output path in its prompt. Tests invoke
`PLY_CANDIDATE_BINARY`. They receive the exact candidate OID/tree and writable
`HOME` and `TMPDIR`; source, Git and evidence remain read-only during acceptance.
Keep test caches and temporary files under those scratch directories.

After pretest and review, existing native delivery verification runs acceptance
again and binds the exact candidate, executable, review and receipt. Only a real
automatic pass permits the agreed local Epic return. Native integration advances
the Epic base, completes the Task and retains its worktree with released
ownership. HumanQA stays absent. The worker performs no remote push, installation
or notification. Existing native completion receipts recover interruptions after
tests or local integration without repeating a completed effect.

Defaults are three developer attempts and 1800 seconds per provider attempt.
Set `--worker-attempts` (1–20) and `--worker-timeout` (10–86400 seconds) at first
start to freeze other bounds. Acceptance uses its Goal's explicit timeout.
Exhausted bounds preserve the Task and report a concrete blocker.

## Execution boundary

The native controller has the invoking user's host authority. Delegate commands
use a probed macOS filesystem/network sandbox: developer Task files and private
output are writable; Git, native state and evidence are read-only. The reviewer
has scratch writes only. Test subprocesses have a separate read-only source
sandbox. These are observed subprocess boundaries, not a sandbox around Ply.

Provider user config and execpolicy rules are ignored for execution. App,
browser, computer-use, plugin, image-generation and agent-delegation features are
explicitly disabled. Discovered MCP servers receive complete inert, disabled
transport overrides. A fresh configuration check verifies these settings before
each launch and preserves only server names/types/booleans, not credentials.
This does not modify global Codex configuration or its existing authentication.

Every actual provider prompt contains its [pinned Addy role package](worker-skills.md)
before the prompt is hashed and preserved in the attempt. Development includes
planning and TDD; corrections add debugging guidance, and review has a separate
read-only adapter. The CLI embeds the exact source files, required references,
manifest and MIT license. No live Addy checkout or network fetch is needed.

## Usage limits and native queue ownership

A positively completed provider attempt with a recognized usage-limit error
parks the same Run as `waiting`, retaining its Task, model, effort, partial files
and attempt streams. `worker show --format json` exposes `retry_not_before`,
`provider_retries` and `last_provider_failure`. A provider reset is honored;
otherwise backoff increases from five minutes to one hour. `--worker-retries`
freezes the retry bound (default 20, range 1–1000), separately from developer
correction attempts. Restart and resume cannot shorten the saved wait. Unknown
processes, truncated evidence and generic failures never trigger blind retry.
Authentication failures remain visible and require recovery of the existing
login; there is no automatic model, account or payment fallback.

The multi-Spec controller can permanently park a concrete failed, quiescent Run
and release its native queue reservation so independent selected work proceeds.
That release is not retry: the released Run cannot resume, even if the controller
was interrupted before recording its local release marker. Native advancement
after integration permits only recovery of the reserved integration and closeout.

The process runner owns child waiting/signalling and records bounded complete
streams and process-group quiescence. It never sends a signal after reaping its
owned child. A deliberately daemonized descendant escaping that process group
is outside this first version's proven containment; background jobs and
delegation are prohibited. This worker is for the selected trusted development
Goal, not a general hostile-code executor.

## Reproduce the isolated journey

With a newly built private binary and existing Codex login:

```sh
go build -o /absolute/private/ply ./cmd/ply
python3 test/fixtures/headless_worker/journey.py --ply /absolute/private/ply --start
```

The script creates a disposable native workspace, local Git Epic, Task, Goal and
queue, then exercises the public start command. Its small executable Spec has
ten baseline behavioral cases. It records command results and `journey.json` and
does not touch a real project, remote or installed Ply. Model calls run only when
`--start` is explicit. Follow its Run with `worker show` and inspect the actual
developer/reviewer JSONL and native delivery receipts.

Unit fixtures label synthetic provider input explicitly. They additionally
exercise inherited-lock controller interruption, failed test correction,
receipt recovery after tests/integration, and rejection of stale candidate,
executable and reviewer bindings. Run `test/delivery_acceptance.sh` for the full
repository checks. The optional configuration-only test can use
`PLY_WORKER_TEST_CODEX=/absolute/codex go test ./internal/taskrun -run
TestWorkerActualProviderToolPolicy -v`; it makes no model or MCP connection.
