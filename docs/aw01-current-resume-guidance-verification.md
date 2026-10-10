# AW-01 current resume guidance: developer verification

This record covers the developer attempt on 2026-10-10 for Goal
`ply-aw01-current-resume-guidance-v1`, immutable manifest
`2d802d8a5503fb873f210b0f4e56c4dc1b688406e65ab987e38147da392764eb`.
The native handoff is
`run_9bb0145f869680ca3bb044ba5b10fabd/handoff.16aa14e07f72dee5bf71b25f137b90f922f850f5dbf485f6365582b80f4319dd.json`.
The starting source was `b27f46b7d20da7bbeb4a6488e1e2c9d64a61f3ee`.

## Change and behavioral coverage

`WorkflowShow` projects the known retained human/resume and worker-wait actions
after resolving existing evidence priorities. It reads pending controls using the
same identity check as `ShowWorker`, which retains its existing lock observation.
There are no transition, retry, provider-policy or schema changes. Raw state and
historical actions remain intact. Unrelated actions and diagnostic reasons are
preserved, including unavailable closeout evidence.

The eight required tests use private physical directories, persisted request and
state bindings, actual attempt/control files, and injected native observations.
They do not initialize a workspace, access a production registry, mutate Git,
start a provider or operate on a historical Run. Synthetic handoff/closeout facts
test the read contract, not native delivery qualification.

| Required test | Observed coverage |
| --- | --- |
| `TestAW01GuidanceLegacyRunningReadback` | Both public readers; developer/reviewer phases; directly seeded legacy resume and wait messages; Run, request and attempt identity retained; no liveness claim. |
| `TestAW01GuidanceResumeTransition` | Actual pause consumption through `ServeWorker`, then matching resume consumption under a held owner lock; stale generation rejected; same reviewer attempt and binding. |
| `TestAW01GuidanceWaitingResume` | Recognized synthetic usage-limit result; real pause/resume consumption; injected time; exact deadline preserved; expiry alone does not release the attempt; missing/unknown completion stays blocked; confirmed completion permits ordinary observation. |
| `TestAW01GuidancePauseStopAndBlocked` | Paused/stopped advice and concrete blocked/unknown process/ownership diagnostics remain unchanged. |
| `TestAW01GuidancePendingControls` | Pause/stop/resume identity and safe-boundary meaning in both readers; consumed fields cleared only in readback; foreign/malformed control evidence rejected without writes. |
| `TestAW01GuidanceCloseoutAndDriftPriority` | Complete/pending closeout, unavailable closeout, changed executable binding and missing Goal evidence retain priority over both stale message kinds. |
| `TestAW01GuidanceReadDoesNotWrite` | Repeated public reads compare every actual file/directory entry, bytes and permissions; absent/held/unowned/unknown owner observations; absent owner lock remains absent; injected file/process/effect adapters reject effects. |
| `TestAW01GuidanceInteractiveUnchanged` | Current and legacy interactive messages, unbound startup, startup recovery and drift; worker show still rejects interactive Runs; no file changes. |

## Actual commands and results

The following paths identify this attempt's private output, not native
acceptance or integration receipts:

```sh
aw01_attempt=/Users/perottochristensen/github/.ply/workflow-runs/runs/wfr_946895b4d004bf2abecb373c9f0ac1d81030782da5917ef8862756682895ae6c/worker/attempts/developer-0001
aw01_scratch="$aw01_attempt/tmp/aw01-local"
aw01_go=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go
```

Each Go command used this environment prefix, with its directories already
created under the writable attempt scratch:

```sh
env HOME="$aw01_scratch/home" TMPDIR="$aw01_scratch/tmp" \
  GOCACHE="$aw01_scratch/go-cache" GOTMPDIR="$aw01_scratch/go-tmp" \
  GOMODCACHE=/Users/perottochristensen/go/pkg/mod \
  GOTOOLCHAIN=local GOENV=off GOWORK=off GOPROXY=off GOSUMDB=off
```

1. **Product RED, before production edits:**
   `"$aw01_go" test -json -mod=readonly -count=1 -p=1 -parallel=1 -timeout=120s -run '^TestAW01GuidanceLegacyRunningReadback$' ./internal/taskrun`
   exited 1. Both readers reported
   `stale active worker guidance: {Actor:human Message:Work is preserved; use worker resume for this same Run.}`;
   the stale worker-wait cases also failed. Raw output is
   `output/guidance-red-2.jsonl` and `.stderr`. The earlier `guidance-red.jsonl`
   was fixture setup failure on null arrays and is **not** product RED evidence.
2. The same command passed after the first projection
   (`output/guidance-green-1.jsonl`). The subsequent focused transition command,
   with `-run '^TestAW01Guidance(ResumeTransition|WaitingResume|PauseStopAndBlocked|PendingControls)$'`,
   passed transitions/waits/blocked cases and failed pending-control assertions
   before the shared control readback was added
   (`output/guidance-transitions-red.jsonl`).
3. The focused priority command used
   `-run '^TestAW01Guidance(CloseoutAndDriftPriority|ReadDoesNotWrite|InteractiveUnchanged)$'`.
   It exposed active observation advice replacing the prior action when closeout
   evidence was unavailable (`output/guidance-priority-red.jsonl`). Moving the
   projection after evidence resolution fixed that failure. No-write and
   interactive cases already passed.
4. **Final GREEN:** the same Go test options with the exact selection below
   exited 0. All eight top-level tests ran and passed exactly once, with one
   package pass and zero fail/skip events. Raw output is
   `output/guidance-final.jsonl` and `.stderr`; the observed count summary is
   `output/verification-summary.json`.

```sh
"$aw01_go" test -json -mod=readonly -count=1 -p=1 -parallel=1 -timeout=120s \
  -run '^(TestAW01GuidanceLegacyRunningReadback|TestAW01GuidanceResumeTransition|TestAW01GuidanceWaitingResume|TestAW01GuidancePauseStopAndBlocked|TestAW01GuidancePendingControls|TestAW01GuidanceCloseoutAndDriftPriority|TestAW01GuidanceReadDoesNotWrite|TestAW01GuidanceInteractiveUnchanged)$' \
  ./internal/taskrun
```

`python3 -I -B test/resume_next_action_dispatch_test.py`, under the same private
HOME/TMPDIR, first failed because the exact Task ref entered the ordinary default
(`output/dispatch-red.txt`). After adding dispatch, all three test groups passed
(`output/dispatch-green.txt`). The routing fixtures cover the exact ref, a nearby
ref, another Task, the receiving Epic, main and detached HEAD. Binding fixtures
reject missing/mismatched OID/tree, dirty/untracked source, another/detached ref
and a missing candidate. They deliberately stop before Go execution and cannot
manufacture a passing product acceptance stream.

Additional successful checks:

- `"$aw01_go" vet -mod=readonly ./internal/taskrun` (`output/vet.stderr` is empty).
- `"$aw01_go" build -mod=readonly -o "$aw01_attempt/output/candidate" ./cmd/ply`.
- The exact candidate's `version --format json` exited 0
  (`output/candidate-version.json`). This is executable identity/smoke evidence,
  not CLI guidance behavior coverage.
- `/bin/sh -n test/delivery_acceptance.sh test/resume-next-action-acceptance.sh`.
- `gofmt` on changed Go files and `git --no-optional-locks diff --check`.

The executable exists at the required private output path. Its observed SHA-256
is `7ecdc1eda36e2dadb4245170083f036c380d11f54464e5854d4035169632b0b8`.

## Acceptance boundary

The committed dispatcher selects the focused driver only for
`refs/heads/ply-task/aw-01-current-resume-guidance-038968865f8f`, before any
ordinary tempfile or fullsuite setup. The reviewed driver is copied byte-for-byte
with SHA-256 `a5b3bb5e6f6d1062851f4c926fb5939ab096f89a315506ef06a4f5482321feb4`.
The ordinary suffix remains byte-for-byte unchanged, with SHA-256
`61abee5a1bab34699d07fcc4883c248cbc019597bc4540561917638efe6d68ac`.
Thus its eight actual-run/pass checks, before/after clean source and candidate
binding, deadlines, and interruption cleanup are preserved unchanged.

The actual supplied native entrypoint was invoked with this candidate,
`PLY_CANDIDATE_OID=b27f46b7d20da7bbeb4a6488e1e2c9d64a61f3ee`, and
`PLY_CANDIDATE_TREE=fd19e12e9ab68e050d44710a7a296822b7136942`:

```sh
/bin/sh /Users/perottochristensen/github/.ply/task-executions/7318d0234c8a7b6aa3e3a2fbfc1f38d083edd14eadbc1ead7327b042a4826a42/acceptance.sh
```

It exited 1 with `candidate source must remain clean`
(`output/native-entrypoint-precommit.json`), correctly refusing this uncommitted
developer tree. No clean-source check was bypassed. Successful native pretest,
separate review, verification, commit and local Epic integration belong to the
controller after this worker exits; none is claimed here. No fullsuite, live
provider/service, installation, external message or Git mutation was performed.

## Review correction: pending resume after a parked wait

Correction attempt `developer-0002` starts at
`38ab7da700a38e902ef7074fbfc27146aeeb4c22` (tree
`6428d96029c9335cf0b9df0d0612830beb2141df`). The prior record above and all prior
attempt evidence remain unchanged. The required P2 finding is preserved in
`reviewer-0001/answer.json` alongside these worker attempts: the reason guard in
`workerCurrentAction` accepted a usage-limit reason only in the `waiting` phase.
Actual pause/stop transitions retain that reason, so both public readers still
advised repeating resume after it had already been published.

`TestAW01GuidancePendingControls/waiting_paused` and `waiting_stopped` now enter
the wait using a confirmed synthetic usage-limit result, consume pause/stop
through `ServeWorker`, and publish resume without consuming it. Both public
readers retain ordinary resume advice before publication, then report the
pending request without claiming consumption. The tests check the exact saved
deadline, reason, retry count, Run/request/attempt identity, parked phase, and
control action/timestamp. The control file and complete file inventory, bytes,
and permissions remain unchanged by the reads.

The fix recognizes only the exact usage-limit reason matching a nonempty saved
deadline, including when pause/stop preserved it. Other reasons, run-level
diagnostics, and the blocked phase retain priority. Additional public-reader
cases cover authentication and unknown process completion, a blocked wait,
mismatched or missing deadlines, and a run-level obstacle with pending resume.
No control consumption, retry policy, transition, or acceptance-driver changes
are needed; the existing Task dispatcher runs the extended mandatory test.

The commands below use these correction-specific paths:

```sh
aw01_attempt=/Users/perottochristensen/github/.ply/workflow-runs/runs/wfr_946895b4d004bf2abecb373c9f0ac1d81030782da5917ef8862756682895ae6c/worker/attempts/developer-0002
aw01_scratch="$aw01_attempt/tmp/aw01-correction"
aw01_go=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go
```

Each command used private HOME/TMPDIR and writable build caches, with the
existing local module/toolchain cache and downloads disabled:

```sh
env HOME="$aw01_scratch/home" TMPDIR="$aw01_scratch/tmp" \
  GOCACHE="$aw01_scratch/go-cache" GOTMPDIR="$aw01_scratch/go-tmp" \
  GOMODCACHE=/Users/perottochristensen/go/pkg/mod \
  GOROOT=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64 \
  GOTOOLCHAIN=local GOENV=off GOWORK=off GOPROXY=off GOSUMDB=off
```

The broader verification/build commands also set `GIT_OPTIONAL_LOCKS=0`.
Exact expanded commands are in `output/verification-commands.json`.

1. **RED before the production fix:**
   `"$aw01_go" test -json -mod=readonly -count=1 -p=1 -parallel=1 -timeout=120s -run '^TestAW01GuidancePendingControls/waiting_(paused|stopped)$' ./internal/taskrun`
   exited 1. All four post-publication public-reader assertions failed with
   `pending resume for a parked wait treated as completed or repeated: {Actor:human Message:Work is preserved; use worker resume for this same Run.}`.
   Pre-publication assertions passed. Raw output:
   `output/pending-resume-red.jsonl`, `.stderr`, and `.exit`.
2. **GREEN after the guard fix:** the same Go options with
   `-run '^TestAW01GuidancePendingControls$'` exited 0, including both new
   transitions and all obstacle/control cases, with no failures or skips.
   Raw output: `output/pending-resume-green.jsonl`, `.stderr`, and `.exit`.
3. **All eight mandatory regressions:** the exact eight-name selection recorded
   above, using the same Go options, exited 0. Each top-level test ran and passed
   once; one package pass and zero failed/skipped events were observed.
   Raw output: `output/guidance-all.stdout`, `.stderr`, and `.exit`;
   parsed red/green counts: `output/test-summary.json`.
4. `python3 -I -B test/resume_next_action_dispatch_test.py` passed all three
   groups (`output/dispatch.stderr`). The reviewed driver and ordinary default
   remain byte-identical. This checks dispatch/binding, not native acceptance.
5. `"$aw01_go" vet -mod=readonly ./internal/taskrun` and
   `"$aw01_go" build -mod=readonly -o "$aw01_attempt/output/candidate" ./cmd/ply`
   both exited 0 with empty stderr. The exact candidate was invoked through
   `PLY_CANDIDATE_BINARY` with `version --format json` and exited 0. This is an
   executable smoke check, not public guidance coverage.
6. `gofmt`, `/bin/sh -n test/delivery_acceptance.sh test/resume-next-action-acceptance.sh`,
   and `git --no-optional-locks diff --check` passed.

The rebuilt executable exists at `developer-0002/output/candidate`, with SHA-256
`3e641de87aac18f992f64ea22d620dd40a5225cd56dd7ee3e3af710eaf48b7ed`.
Its identity and smoke output are in `output/candidate-identity.json` and
`output/candidate-version.stdout`.

The supplied native entrypoint was rerun with this candidate and the actual
starting HEAD/tree above in `PLY_CANDIDATE_OID` / `PLY_CANDIDATE_TREE`:

```sh
/bin/sh /Users/perottochristensen/github/.ply/task-executions/7318d0234c8a7b6aa3e3a2fbfc1f38d083edd14eadbc1ead7327b042a4826a42/acceptance.sh
```

It exited 1 with `candidate source must remain clean`, before running tests
(`output/native-entrypoint-precommit.stdout`, `.stderr`, and `.exit`). The
clean-source guard is intact; this is not a successful native acceptance receipt.
The correction is ready for the controller's commit, fresh native acceptance
and independent review. The controller also retains delivery and notification
ownership. No native acceptance, review approval, or local integration is claimed
for this correction attempt.
