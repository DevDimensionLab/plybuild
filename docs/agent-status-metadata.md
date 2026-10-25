# SM01 public status metadata

Run `ply workflow agents --metadata-only --format json` with the ordinary CLI.
The result is `AgentStatusMetadataReadback@1`, schema_version 1, specified in
`schemas/agent-status-metadata-v1/schema.json`; representative public responses
and a raw sessionless source are in its `examples` directory. Text output is
also supported. `--workspace` remains accepted for CLI compatibility but is not
read or discovered on this path: the source scope is the current Herdr user,
not the Ply workspace. There are no Task/Run, session, terminal or binding fields.

The flag selects this path before constructing the full AgentWork Service.
Each invocation executes one public `herdr agent list`, followed by bounded
`herdr tab list --workspace <id>` calls for optional actual tab names. Both
share the request deadline. It reads no native inventory/history, human questions,
follow-ups, bindings, dispatch, transcript or process probes, and
uses no cross-request cache. Existing `workflow agents` without the flag keeps
its `AgentWorkReadback@1` execution completeness and fresh dispatch vetoes.

The supported enumeration is the nonpaginated public envelope
`{"result":{"type":"agent_list","agents":[...]}}` with optional transport `id` and absent or null `error`.
Only `agent`, `agent_status`, `pane_id`, `tab_id`, and `workspace_id` are decoded
from rows; execution fields are opaque, including wrongly typed session fields.
Codex and Claude rows are included regardless of Ply binding or native session.
Other valid provider names are outside scope. A missing/invalid provider cannot
prove it is outside scope and makes coverage partial. Unknown status text maps
to `unknown`; it does not invent a known state or activity timestamp.

The single source has independently qualified `coverage`:

- `complete`: a recognized `agent_list` enumeration, at most 4096 source rows,
  with valid unique metadata for all recognized agents.
- `partial`: valid enumeration with malformed rows, duplicate pane identities,
  limits exceeded, or additional envelope/result members that could signal narrowed,
  paginated, unknown or truncated enumeration. Every ambiguous pane is omitted;
  valid rows are preserved, ordered by pane. At most the first 4096 source rows
  are considered; exceeding that bound never reports complete coverage.
- `unknown`: no qualified enumeration, including missing/unknown `result.type`,
  malformed/duplicate envelope members, null agents, unsuccessful transport,
  truncated JSON, oversized stdout/stderr or timeout. No agents are emitted.

`status` is `available` for a complete nonempty result, `empty` for a complete
successful empty projection, `partial` for partial coverage, and `unavailable`
for unknown coverage. Thus successful empty cannot stand in for failed reading.
`freshness=fresh` means a response was received by this invocation, even when
metadata coverage is partial; it does not claim that an agent just acted.
`source_observed_at_utc` is the completed response observation, null when the
source could not be qualified. `observed_at_utc` is captured after transport
returns, including failed responses. These clocks never claim last activity.

The five base fields remain required; optional `tab_label` is display-only.
Labels join on the exact workspace/tab pair, never pane titles or execution
identity. Missing, rejected, duplicate or failed label reads omit only the label;
membership/status/coverage remain those of agent-list. Every refresh reads names
anew, including renames and disappearances. Multiple agents in a tab share its
name, while equal names never merge identities.

Labels retain usable Unicode verbatim (including `cmd`, `1`, accents and inert
HTML-like text). In addition to the public text safety floor below, reject
Unicode Cf, U+2028/U+2029, U+FFFD, malformed Unicode and case-insensitive
`session_id`, `terminal_id`, `cwd=`, `token=`, `secret=` and `password=`.
No raw rejected value or provider diagnostic is emitted. Full AgentWork keeps
its existing display policy.

All returned fields must be strings;
pane must be nonempty. Tab/workspace and status can be empty; empty status maps
to unknown. Public text uses the existing safety policy: no paths, control text,
surrounding whitespace, backslashes, URLs or @ text. Each string is limited to
256 UTF-16 code units. JSON Schema `maxLength` counts Unicode scalar values,
so its 256 bound is necessary but insufficient for astral characters; runtime
validation and conformance tests enforce 128 emoji = 256 units, rejecting an
additional character. No truncation turns invalid metadata into valid metadata.

Membership transport caps stdout and stderr independently at 8 MiB. Labels
reuse the bounds of 32 sorted workspaces, four parallel calls per request,
256 KiB per stream per workspace, and 4096 tab rows per workspace. More
workspaces receive no label; an oversized tab response contributes no labels.
These limits are independent: 4096 long distinct names may exceed the byte
limit even below the row limit. Membership remains complete when qualified.
There is no global admission scheduler or cache; four concurrent callers each
retain this per-request bound (at most sixteen concurrent label children).

One 2-second budget starts before membership transport. The transport deadline
reserves 50 ms for pipe cleanup, shared by all label workers. A smaller supplied
adapter timeout and parent cancellation are honored; no new label deadline
starts after membership. On Unix, when Ply is a process-group leader, Herdr
inherits that group so a caller's final group cleanup can reach descendants
even after Ply is hard-stopped. **Group leadership is not evidence that a
cleanup owner exists.** Ply bounds its direct child and pipes in this branch;
without caller cleanup, descendants may survive even a normal timeout/return.
The identical isolated reproducer confirmed this preexisting limitation on
installed clean F02 `7eb42e741e4cd0646c6447fdb2f5f615fbb14620` and the producer
candidate. No standalone cleanup pass is claimed.

When Ply is not group leader, each read owns a child group and terminates it on
cancellation/completion. An unmanaged, uncatchable hard kill cannot execute
Ply's deferred cleanup. Windows
bounds the direct child/pipes but does not claim descendant cleanup. Decode and
serialization are bounded by the byte/row limits; actual full-process timings
must still qualify against the consumer's unchanged 2600 ms read timeout and
3000 ms whole-HTTP boundary. No retries occur. Source failure returns exit 0 and
honest unavailable JSON; optional label failure preserves base membership.

## Consumer migration handoff

The consumer invokes only this ordinary CLI and validates the versioned result,
limits and source qualification. It does not import Ply packages, read `.ply`,
call Herdr directly, or interpret an empty array alone as success. To preserve
status-v1's binary source state, map only complete + fresh available/empty to
Herdr `available`; map partial/unavailable/unknown to `unavailable` and an empty
agent list. Map the five agent fields without adding execution identity. This
is a conservative migration rule; richer partial coverage requires a separately
agreed future consumer API version.

Keep GET `/api/v1/status`, limited OPTIONS, `version: "1.0"`, and response
`observed_at` generated when the consumer creates its response. Do not replace
it with an agent activity clock. Preserve maximum 4096 agents and 256 UTF-16
units per string, nonempty pane, `no-store`/`Vary`, only Firebase origin
`https://plyapp-pwa.web.app`, and exact trusted-proxy Tailscale login
`perottobc@gmail.com`. Keep existing HTTPS/auth/CORS; the current consumer read timeout is 2600 ms
and whole-response acceptance is 3000 ms.
No browser bearer/cookie or private control route is added. Consumer code,
activation, installation and installed timing belong to the coordinator and
PWA lead, outside this producer worktree's implementation boundary.

Run the exact-candidate conformance check with
`python3 test/status_metadata_acceptance.py /absolute/private/ply`.
It uses isolated public-source fixtures, checks JSON privacy and limits, and
reports actual elapsed times and margin to 2600 ms for success and negative
sources, including a hanging child that retains output pipes.

Run `test/status_metadata_labels_acceptance.py /absolute/private/ply` using the
pinned test-only dependencies in `test/read_contract_requirements.txt` for the
schema/negative fixtures and complete-name timing matrix (cold-first/repeated,
9/4096 agents, multiple workspaces, four concurrent reads, stress and recovery).
Every sample includes process startup, full JSON read and semantic validation;
all samples and the maximum are printed. This is producer evidence, not an HTTP
pass. JSON Schema scalar maxLength is supplemented by the explicit UTF-16 check.

The external-group samples explicitly perform caller `killpg` after Ply exits;
they qualify that isolated composition only. They do not prove standalone
cleanup or identify a cleanup owner from process topology. The existing bound
server adapter's Setpgid/group.finish behavior still needs exact-artifact
qualification for timeout, TERM, hard kill/inherited pipes, permits/recovery,
four supported concurrent reads, complete healthy membership/status/names and
full-body HTTP <=3000 ms. Group escape and regressions are not exempted.

Open baseline limitation **T03079C2A / r001-f01** is tracked by Ply-plan in
`planning/koe/metadata-standalone-gruppeleder-cleanup.md`, not in a second product
backlog. The consumer clarification is
`planning/design/evidence/metadata-tab-label-20261025/standalone-cleanup-scope.md`.
These are source-bound review evidence, not edits to the frozen Goal/Spec or a
replacement for independent review. The retained red reproducer compares the
exact installed F02 binary with each reviewed candidate and records timeout,
descendant survival and external cleanup separately. The standalone limitation
remains open; no generic runner, watchdog, process inspection or new ownership
protocol is introduced by this task.

Compatibility hold: the old strict server rejects the new optional field.
Deliver the clean source, exact schema/examples and retained staged binary first.
Do not install or activate changed default emission until the coordinator has
observed separately authorized compatible server activation. Parent/root owns
whole HTTP qualification, installation and activation; none is implied by a
producer test pass.
