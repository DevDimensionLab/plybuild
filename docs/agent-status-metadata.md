# SM01 public status metadata

Run `ply workflow agents --metadata-only --format json` with the ordinary CLI.
The result is `AgentStatusMetadataReadback@1`, schema_version 1, specified in
`schemas/agent-status-metadata-v1/schema.json`; representative public responses
and a raw sessionless source are in its `examples` directory. Text output is
also supported. `--workspace` remains accepted for CLI compatibility but is not
read or discovered on this path: the source scope is the current Herdr user,
not the Ply workspace. There are no Task/Run, session, terminal or binding fields.

The flag selects this path before constructing the full AgentWork Service.
Each invocation executes exactly one public `herdr agent list` through the
existing bounded adapter. It reads no native inventory/history, human questions,
follow-ups, bindings, dispatch, tab metadata, transcript or process probes, and
uses no cross-request cache. Existing `workflow agents` without the flag keeps
its `AgentWorkReadback@1` execution completeness and fresh dispatch vetoes.

The supported enumeration is the nonpaginated public envelope
`{"result":{"type":"agent_list","agents":[...]}}` with absent or null `error`.
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
  limits exceeded, or additional result members that could signal narrowed,
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

Only the five allowlisted fields are returned per agent. All must be strings;
pane must be nonempty. Tab/workspace and status can be empty; empty status maps
to unknown. Public text uses the existing safety policy: no paths, control text,
surrounding whitespace, backslashes, URLs or @ text. Each string is limited to
256 UTF-16 code units. JSON Schema `maxLength` counts Unicode scalar values,
so its 256 bound is necessary but insufficient for astral characters; runtime
validation and conformance tests enforce 128 emoji = 256 units, rejecting an
additional character. No truncation turns invalid metadata into valid metadata.

Transport caps stdout and stderr independently at 8 MiB. The metadata deadline
is 2 seconds (a smaller supplied adapter timeout is honored) with the existing
maximum 1-second pipe cleanup allowance. No retries occur. This leaves a nominal
2-second transport margin within the consumer's unchanged 5-second deadline;
actual end-to-end CLI timing must be measured on the installed bytes. A source
failure is a successful CLI readback with exit 0 and `unavailable` JSON, rather
than stderr containing raw provider data.

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
`perottobc@gmail.com`. Keep existing HTTPS/auth/CORS and the 5-second deadline.
No browser bearer/cookie or private control route is added. Consumer code,
activation, installation and installed timing belong to the coordinator and
PWA lead, outside this producer worktree's implementation boundary.

Run the exact-candidate conformance check with
`python3 test/status_metadata_acceptance.py /absolute/private/ply`.
It uses isolated public-source fixtures, checks JSON privacy and limits, and
reports actual elapsed times and margin to 5 seconds for success and negative
sources, including a hanging child that retains output pipes.
