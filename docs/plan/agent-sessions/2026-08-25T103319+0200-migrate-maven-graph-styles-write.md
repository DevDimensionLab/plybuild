# Agent Session: Migrate Maven Graph Styles Write

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T103319+0200-migrate-maven-graph-styles-write`
Created: `2026-08-25T10:33:19+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `b319c1d8f3c99869bb499096eda2e0d175ab3a900a05e2c5f26617b73240fd43`
Previous: [2026-08-25T100637+0200-migrate-structurizr-output-write.md](2026-08-25T100637+0200-migrate-structurizr-output-write.md)
Next: [2026-08-25T110346+0200-migrate-local-config-touch-write.md](2026-08-25T110346+0200-migrate-local-config-touch-write.md)
Outcome: Maven graph-style output write routed through the filesystem adapter at `886dff0`; clean P3.36 gate passed.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route the direct
`os.WriteFile(stylesFile, jsonStyles, 0644)` in `maven.WriteGraphStyles`
through the existing filesystem adapter. Preserve exact JSON marshaling, early
error, target path, existence selection, no-write-when-present behavior, bytes,
mode, write error, Graph and command behavior, every completed filesystem,
file, template, config, Maven, tips, structurizr, Bitbucket, Wpost, and
supervisor move, and every P2A contract with zero comparable ratchet
regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`os.WriteFile(stylesFile, jsonStyles, 0644)` operation in
`WriteGraphStyles` in `pkg/maven/graph.go` and its use of the existing
zero-value-safe filesystem `WriteFile` operation. It authorizes focused Maven
graph-style recording contracts for that write boundary. It does not authorize
changing `file.Exists`, JSON marshaling, path construction, Graph filtering or
style mutation, Graph arguments, `RunOn`, Maven/dot execution, callers, an
adapter or established complete-double extension, another production effect, a
function-valued effect dependency, a new adapter family or public API,
mutation harnesses, or later roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, complete `pkg/maven/graph.go` and all Maven tests,
every caller and use of `WriteGraphStyles`, `Graph`, `GraphDefaultStyles`, and
`GraphArgs`, complete `cmd/plugin_maven.go` and its registration and relevant
command tests, relevant Maven/process/filesystem/file/template/config/tips/
structurizr behavior, the filesystem and process adapters and relevant
recording doubles, `.quality/inventory`, and the continuity and quality-lift
designs. Regenerate ignored reports outside the measured tree or remove them
before a clean audit.

Implementation commit `527a8b9` has 266 tests across 17 of 25 packages. Q0.6
has 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 8 of 25, Q1.2 is 0, Q1.3 is 31 violations of 48 production
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
implementation commit for the Maven graph-style output-write boundary, package
contracts, and measured planning notes. Then perform the normal separate
handoff-only commit. Do not push, merge, publish, distribute, remove the
worktree, stash inherited changes, revert user work, or run destructive Git
commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, all command objects and registrations, public API/CLI behavior,
all four host acceptance flows, `WriteGraphStyles`, `GraphDefaultStyles`,
`GraphArgs`, `Graph`, `RunOn`, plugin-Maven behavior, ProjectConfig,
CloudConfig, LocalConfig, structurizr and tips behavior, template/file/config/
Maven behavior, and every caller's observable behavior.

Use the existing complete filesystem dependency and existing
`WriteFile(name string, data []byte, perm fs.FileMode) error` operation with a
safe zero value; do not add or change an adapter operation, extend an
established complete recording double, or store a function-valued effect
dependency. Production must select `filesystem.System()` only for a private
complete Maven graph-style-write composition.

Preserve exact `json.MarshalIndent(styles, "", "    ")` first and its early
error return, then exact
`file.Path("%s/target/dependency-graph-styles.json", projectPath)`, then the
existing `file.Exists(stylesFile)` selection and nil return without a write
when it is true. When absent, make one adapter write attempt with exact
`stylesFile`, exact `jsonStyles`, and exact mode `0644`; return the exact write
error. Preserve nil, empty, representative, non-ASCII, and arbitrary map
content and ordering semantics supplied by `encoding/json`, arbitrary project
paths, and Graph's style/filter/argument and later `RunOn` behavior. Do not
clean or join paths, preflight, retry, normalize, validate, create directories,
inspect permissions, log, wrap or discard the write error, write twice, or
broaden into the existence boundary, Maven execution, structurizr, tips,
config, template, file, HTTP, Wpost, clock, server, P4, P5, dependency, Docker,
distribution, or publication work.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete
`pkg/maven/graph.go` and all Maven tests, every graph caller, complete
`cmd/plugin_maven.go` and relevant registration/command tests, and relevant
Maven/process/filesystem/file/template/config/tips/structurizr code and tests.
Read `internal/adapter/filesystem`, `internal/adapter/process`, their relevant
complete recording doubles, and the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit
implementation before editing. Before the full gate, read the complete
launcher contract, relevant Make meta-tests, P2A API/CLI/subprocess contracts,
and all four host acceptance flows.

# Three Moves

1. Start red with focused Maven graph-style-output recording contracts. Prove
   the complete dependency, production selection of `filesystem.System()`,
   exact caller-composed output path, non-empty recorded write populations,
   empty, representative, non-ASCII, and arbitrary output bytes, exact `0644`
   mode, one write attempt, exact write-error identity at the private boundary,
   safe zero behavior without developer-path access, and no unrelated adapter
   operation. Production composition contracts must not mutate the real
   filesystem or run Maven, dot, or another command.

2. Reuse only the existing filesystem adapter `WriteFile` operation. Keep
   exported `WriteGraphStyles(styles GraphStyles, projectPath string) error` as
   the production entry, select `filesystem.System()` only for a private
   complete graph-style-output-write composition, and replace only its direct
   `os.WriteFile` with the adapter operation. Do not change marshal evaluation,
   early error, path construction, `file.Exists` evaluation or result, exact
   bytes, mode, error, adapter, another complete double, function or caller,
   inventory, or public API.

3. Run focused Maven/filesystem contracts and relevant process, command, file,
   template, config, tips, and structurizr caller packages' tests, the launcher
   contract from `/bin/bash`, Make preflight meta-contracts, API/CLI
   compatibility, full Go tests and race/vet, all four host acceptance flows,
   the audit meta-suite, focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4
   measurements, full clean checkpoint audit, and empty-HOME count-2. Expect
   Q1.1 to hold at 8 of 25 and Q1.3 to move nominally from 31 of 48 to 30 of 47
   while Q0.6, Q1.2, Q1.4, and exact Q2.1 hold. Regenerate exact values; the
   full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Maven
graph-style output-write move or record an exact resumable state. Rewrite the
rolling handover, record the measured P3 result, answer this archive, create
one linked NEXT archive for the next coherent P3 effect move, replace only the
launcher's mutable regions, run the launcher contract, and make the separate
handoff-only commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
