---
name: ply-machine-load
description: 'Check how hard this machine is working with `ply machine status` before starting a heavy job (full test suite, `make check`, Playwright, a plycq scan or review of a whole repo, a release build), and wait or run in series when it is under pressure. Advisory: Ply blocks nothing.'
---

# Machine load before heavy jobs

Several agents share this machine. Heavy jobs started at the same time slow each other down and
make load-sensitive tests fail. Read the load before you start one.

## Before a heavy job

A heavy job is a full test suite (`go test ./...`, `make check`), Playwright, a plycq scan or
review of a whole repo, or a release build. A single package test or a small command is not.

1. Run `ply machine status --format json`.
2. The machine is **under pressure** when `load_percent["1m"]` is above 100, or
   `memory_pressure` is `warn` or `critical`.
3. Not under pressure: start the job.
4. Under pressure: wait and read again about every 30 seconds, or run your own heavy jobs one at
   a time instead of together. The 1-minute value lags about a minute after a job ends, so read
   again rather than trusting one value.
5. Still under pressure after about 10 minutes: start the job anyway, and say that the machine
   was under pressure when you report its result.

To see what takes the load before you wait, run `ply machine status --processes`: it adds the
ten processes that use most CPU, by program name only.

`memory_used_percent` and `swap_used_bytes` show how full memory is; growing swap means memory is
the bottleneck even when `memory_pressure` is still `normal`.

`memory_pressure` can be `null` with a reason; then judge by `load_percent` alone. If the command
fails (`unsupported_os`, `read_failed`), start the job and mention that the load was unknown.

## Limits

- The threshold is advice. Do not skip a job because of it, and do not stop or slow other
  agents' processes; reading the load gives no authority over them.
- A test that fails under load and passes alone is still a flaky test. Report it where the
  project registers failures and fix it to wait for state, not for time. Waiting for a quiet
  machine does not fix it.
