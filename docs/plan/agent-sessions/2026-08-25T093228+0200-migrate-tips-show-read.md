# Agent Session: Migrate Tips Show Read

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T093228+0200-migrate-tips-show-read`
Created: `2026-08-25T09:32:28+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `e3aee95deb8606bfa07dec7e14ea0ace9179a458a0847f0ffe6e1c615eaf9642`
Previous: [2026-08-25T085931+0200-migrate-config-project-write.md](2026-08-25T085931+0200-migrate-config-project-write.md)
Next: [2026-08-25T100637+0200-migrate-structurizr-output-write.md](2026-08-25T100637+0200-migrate-structurizr-output-write.md)
Outcome: completed by implementation commit `2566438af9d52505b54b92428913fb6b89f4fd38`; clean full audit recorded Q1.3 at 32/49 with zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route the direct
`os.ReadFile(tipsPath)` in `tipsShowCmd.RunE` through the existing filesystem
adapter. Preserve command registration and delegation, exact tip-name and path
composition, read/error/render/output/log/cloud sequencing, bytes and errors,
every completed filesystem, file, template, config, Maven, Bitbucket, Wpost,
and supervisor move, and every P2A contract with zero comparable ratchet
regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`os.ReadFile(tipsPath)` operation in `tipsShowCmd.RunE` in `cmd/tips.go` and
its use of the existing zero-value-safe filesystem `ReadFile` operation. It
authorizes focused command recording contracts for that read boundary. It does
not authorize changing another tips branch or command, argument handling,
command registration, profile initialization or sync, path composition,
markdown rendering, terminal config, stdout, logging, cloud-source lookup,
error wrapping, another command composition, an adapter or established
complete-double extension, another production effect, a function-valued
effect dependency, a new adapter family or public API, mutation harnesses, or
later roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, complete `cmd/tips.go` and its registration and
command tests, every caller and use of `tipsCmd`, `tipsListCmd`,
`tipsShowCmd`, `tips.LocalDir`, `tips.List`, `file.Path`, local terminal config,
cloud config and `GlobalCloudConfig`, the complete `pkg/tips` implementation
and tests, relevant command/context/config/file/template/Maven behavior, the
filesystem adapter and relevant recording doubles, `.quality/inventory`, and
the continuity and quality-lift designs. Regenerate ignored reports outside
the measured tree or remove them before a clean audit.

Implementation commit `7754575` has 258 tests across 16 of 25 packages. Q0.6
has 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 9 of 25, Q1.2 is 0, Q1.3 is 33 violations of 50 production
effect sites with clock and server absent, Q1.4 is 7 of 8, and exact Q2.1 is 0
of 8 executable harnesses. The clean gate passed, the full audit exited 1 for
15 documented findings and never 2, and comparable ratchets were five
improved, two held, and zero regressed.

The launcher has 62 Bash 3.2 contracts and supervises fresh non-interactive
JSONL turns with external raw logs. It continues only after a successful
structured stream and valid clean committed handoff. Preserve its stable
skeleton and do not launch a real successor while developing, testing, or
finalizing this move.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the tips-show read boundary, command contracts, and
measured planning notes. Then perform the normal separate handoff-only commit.
Do not push, merge, publish, distribute, remove the worktree, stash inherited
changes, revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, all command objects and registrations, public API/CLI behavior,
all four host acceptance flows, `ProjectConfig`, `CloudConfig`, `LocalConfig`,
`tips.LocalDir`, `tips.List`, `file.Path`, markdown rendering, template/file/
config/Maven behavior, and every caller's observable behavior.

Use the existing complete filesystem dependency and existing
`ReadFile(name string) ([]byte, error)` operation with a safe zero value; do
not add or change an adapter operation, extend an established complete
recording double, or store a function-valued effect dependency. Production
must select `filesystem.System()` only for a private complete tips-show-read
composition.

Preserve the zero-argument warning and delegation exactly. For a non-empty
argument list, preserve exact `name := args[0]`, then exact
`file.Path("%s/%s.md", tips.LocalDir(ctx.CloudConfig), name)` evaluation, then
one adapter read attempt with that exact `tipsPath`. Preserve the read error's
exact wrapping and identity through
`fmt.Errorf("failed to find any tips file for [%s]: %s: %w", name, tipsPath,
err)`. On success preserve the exact returned bytes, then terminal-config
lookup, `markdown.Render(string(source), terminalConfig.Width, 2)`, stdout,
local-source log, global-cloud-config lookup, and cloud-source log in their
current order. Preserve arbitrary argument and path content, empty and
non-ASCII bytes, and nil and non-nil dependency errors. Do not clean or join
paths, add extensions differently, preflight, retry, normalize, validate,
mutate globals, read twice, log on read failure, or broaden into list, sync,
profile, config, template, file, Maven, HTTP/process, Wpost, clock, server, P4,
P5, dependency, Docker, distribution, or publication work.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete `cmd/tips.go`, all
tips command registration and caller tests, complete `pkg/tips`, and relevant
context, config, file, template, and Maven code and tests. Read
`internal/adapter/filesystem`, its relevant complete recording doubles, and
the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation before editing.
Before the full gate, read the complete launcher contract, relevant Make
meta-tests, P2A API/CLI/subprocess contracts, and all four host acceptance
flows.

# Three Moves

1. Start red with focused tips-show-read recording contracts. Prove the
   complete dependency, production selection of `filesystem.System()`, exact
   caller-composed path delivery, non-empty recorded read populations, empty,
   representative, and non-ASCII returned bytes, one read attempt, exact read
   error identity at the private boundary, safe zero behavior without
   developer-path access, and no unrelated adapter operation. Production
   composition contracts must not read the real filesystem or run the command.

2. Reuse only the existing filesystem adapter `ReadFile` operation. Keep all
   command objects and `RunE` behavior as production entries, select
   `filesystem.System()` only for a private complete tips-show-read
   composition, and replace only the direct `os.ReadFile(tipsPath)` with the
   adapter operation. Do not change the missing-argument branch, path
   expression, error composition, returned bytes, later command sequencing,
   adapter, another complete double, command, caller, inventory, or public API.

3. Run focused command/filesystem contracts and relevant tips, context,
   config, file, template, and Maven caller packages' tests, the launcher
   contract from `/bin/bash`, Make preflight meta-contracts, API/CLI
   compatibility, full Go tests and race/vet, all four host acceptance flows,
   the audit meta-suite, focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4
   measurements, full clean checkpoint audit, and empty-HOME count-2. Expect
   Q1.3 to move nominally from 33 of 50 to 32 of 49 while Q0.6, Q1.2, Q1.4,
   and exact Q2.1 hold. Regenerate exact values; the full audit may exit 1 for
   documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent tips-show read
move or record an exact resumable state. Rewrite the rolling handover, record
the measured P3 result, answer this archive, create one linked NEXT archive
for the next coherent P3 effect move, replace only the launcher's mutable
regions, run the launcher contract, and make the separate handoff-only commit
`docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
