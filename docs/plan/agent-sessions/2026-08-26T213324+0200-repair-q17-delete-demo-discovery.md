# Agent Session: Repair Q1.7 Delete Demo Discovery

Status: NEXT
Session ID: `2026-08-26T213324+0200-repair-q17-delete-demo-discovery`
Created: `2026-08-26T21:33:24+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `f3b98dbaad1e7c0f7839914054a2ddb71e27b297ec56a9f3d7cc62e9fdb591c7`
Previous: [2026-08-26T204147+0200-repair-q17-partial-failures.md](2026-08-26T204147+0200-repair-q17-partial-failures.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P4 with one narrowly bounded Q1.7 controllability repair. Add the
smallest private seam needed to exercise `spring.DeleteDemoFiles`' discovery
failure, then add an exact executable contract for its warning content and
continuation to the fixed demo-file loop. Reclassify all 17 Q1.7 operations
from executable control flow before deciding whether Q1.7 is supported.

# Authorized Roadmap

P3 is complete. P4 remains active for this one Q1.7 repair and later combined
manual evidence; P5-P8 remain queued in the machine-readable block in
`docs/plan/quality-upgrade.md`. Q1.6 is independently supported by its complete
57-struct, 74-field review, and Q1.9 is independently supported by its complete
169-range classification. No schema-2 receipt is authorized in this seam
session. The launcher must remain NEXT until P4 and P5-P8 finish.

This mission authorizes exactly one private, unexported function-parameter
seam in `pkg/spring/io.go` around the `file.FindFirst` call made by
`DeleteDemoFiles`; one focused exact discovery-failure test and its test-only
helpers; read-only reclassification of all 17 Q1.7 operations; one focused
production-and-test implementation commit; clean automated measurement; and
the normal separate continuity commit. Prefer a private helper to which the
public `DeleteDemoFiles` delegates with `file.FindFirst`. Do not add a
dependency struct or inject deletion behavior.

This mission does not authorize an exported API or signature change, any
other production change or seam, manual evidence or receipts, Q1.6 or Q1.9
changes, audit/parser/scanner changes, `.quality/inventory` or baseline
changes, executable P5 harnesses, T1-T10 changes, acceptance changes, P6-P8
work, or a Q1.7 PASS claim unless the complete classified population supports
it.

# Measurements At Start

The Q1.7 characterization move produced test-only commit
`7be93e7e822f05bdb007e4e5a1a00402e620a6cb`, exact tree
`0f0955f8bbe3e49f5c8d1b4f98454dcb04ea62f1`, on unchanged production commit
`c382caa38be167fe17f847370ad8a12270644de3`. Ten new top-level contracts bring
the suite from 453 to 463 tests. The exact affected packages, complete
33-contract focused population, complete uncached tests, race, vet, API/CLI
and subprocess compatibility, pinned lint, launcher and Make contracts,
complete preflight, all four host acceptance flows, the 15-control audit
meta-suite, and empty-HOME count-2 pass.

The clean no-evidence scorecard SHA-256 is
`b2ccc93dbb045f06404f955ac4b122a614a654f2c0d983d48941fe7026b352ec`.
It records 463 tests across all 27 packages, L0 8/8, L1 six PASS and
Q1.6/Q1.7/Q1.9 UNMEASURABLE, six improved ratchets, two held, zero regressed,
zero non-comparable, and zero dirty paths. The audit exits 1 for 13 documented
non-passing criteria, never 2.

The exhaustive 17-operation Q1.7 record now classifies 16 operations as having
exact failure content plus continuation or exact ordered partial-result
coverage. Its SHA-256 is
`4b247f3343a651374b96245d6c15b563f573a152799f3329713c578ace863e57`.
The exact 33-contract focused population manifest SHA-256 is
`43391020bb74aa696a750c68cbdefad40ceee4de5928c4f547356f581ddcae22`.

The sole remaining included branch is `DeleteDemoFiles`' discovery warning.
The public function hard-codes `file.FindFirst`. That function hard-codes the
system filesystem, its walk callback ignores each incoming walk error and
returns only `nil` or `io.EOF`, and final `io.EOF` is normalized to nil. Safe
test-only approaches therefore cannot make the checked discovery error
non-nil. Deletion failures in the same operation already have exact warning
and later-item continuation coverage.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Confirm the new
continuity HEAD and exact ancestry before editing. Preserve the public
`DeleteDemoFiles(string, config.ProjectConfiguration)` signature and make it
delegate to one unexported helper accepting only the discovery function. Pass
`file.FindFirst` in production. Preserve the current suffix selection, lookup
arguments, exact logs including existing spelling, deletion behavior, item
ordering, side effects, and Go 1.18 compatibility.

The new test must inject a sentinel discovery error, assert the exact warning
`Unable to find testfile, fileSuffix=.kt`, prove the discovery function
received the exact suffix and test-directory path, and prove the later
`HELP.md`, `mvnw`, and `mvnw.cmd` items were deleted. Preserve every existing
test and process-state restoration. Make one focused implementation commit,
then the separate `docs: prepare next agent session` continuity commit.

Do not push, merge, publish, distribute, remove the worktree, stash inherited
changes, revert user work, or run destructive Git commands. Keep every
generated report outside the worktree. Set `API_COMPAT_REPORT_OUT` and
`CLI_COMPAT_REPORT_OUT` to external paths whenever running compatibility. The
measured-tree identity counts ignored, untracked, and tracked bytes, so confirm
both ordinary and ignored status before every authoritative audit.

# Required Reading

Before editing, confirm branch, HEAD, clean and ignored status, reciprocal
archive links, launcher `--check`, exact ancestry, and the authorized
checkpoint block. Read the rolling handover, this archive, the complete P4
entry and gate, both design documents, `.quality/README.md`, the Q1.7 criterion
and schema-2 validation path, all 17 production operations and enclosing call
flows, the complete 17-row classification, every focused partial-failure test
and helper, `pkg/spring/io.go`, `pkg/file/file.go`, the complete audit wrapper,
T15 repair, API/CLI and launcher contracts, and Make and acceptance gates. Do
not classify an operation by name or search hit alone.

# Three Moves

1. Regenerate the complete 17-operation Q1.7 record from the clean current
   tree. Implement only the private discovery-function seam described above.
   Add the exact discovery warning, lookup-argument, and fixed-file continuation
   contract. Inspect every included and excluded branch in executable context.
2. Run the exact Spring tests, the complete 17-operation focused population,
   complete uncached tests, race, and vet. Reclassify all 17 operations and
   require exact failure content plus continuation or exact ordered partial
   result for every included path. If any included path remains uncovered,
   keep Q1.7 blocked. Do not create manual evidence in this session.
3. Run API/CLI and subprocess compatibility with external report paths; pinned
   lint; launcher and Make contracts; complete preflight; all four host
   acceptance flows; the 15-control audit meta-suite; empty-HOME count-2; and a
   clean no-evidence audit. Require audit exit never 2 and zero dirty paths.
   Record the implementation commit, exact population, scorecard, and whether
   the next P4 move is evidence-only or another exact blocker repair.

# Automatic Handoff

Before this session ends, finish and commit the narrow Q1.7 seam and contract
or record an exact resumable blocker. Rewrite the rolling handover, record the
P4 result in the roadmap, answer this archive, create exactly one reciprocally
linked NEXT archive for the next coherent P4 move, replace only the launcher's
mutable regions, run launcher and handoff contracts, and make the separate
continuity-only commit. If all 17 operations support Q1.7, the next coherent
P4 move is combined schema-2 manual evidence for Q1.6, Q1.7, and Q1.9; do not
create that evidence in this session. Do not launch a real successor, push,
merge, publish, distribute, stash, revert, or remove the worktree. COMPLETE
remains invalid while P4 or P5-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
