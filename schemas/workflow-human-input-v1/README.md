# Human input contract v1 (`cooperative_agent_v1`)

The public contract of `ply workflow input`: an ordinary Claude or Codex agent
publishes one necessary human question; a configured consumer (the app)
discovers it, records one immutable answer and reads the agent's reported
receipt and use. Mode `cooperative_agent_v1`, `schema_version: 1`. It is
additive and separate from the schema 2 `week_codex_same_owner_turn_v1`
answer contract (`ply workflow execute question|answer`), whose behavior and
journals are unchanged.

- [Schema](workflow-human-input-v1.schema.json) (JSON Schema 2020-12) has eight
  closed roots under `$defs`: `question_draft`, `question`, `list`, `preview`,
  `request`, `status`, `wait` and `error`.
- [Examples](examples.json) are generated from the actual commands with a fixed
  clock: every valid result with its command line and exit code, labeled
  schema-invalid values, and RFC8785 canonical request vectors.
- `ply capabilities --format json` advertises each operation below under
  `workflow_extensions` with `mode: cooperative_agent_v1`, exact `command`,
  `selectors`, `effect` and `result_schemas`. Require both the capability and
  this bundle's digests before executing a write.

Nothing here starts, wakes or resumes an agent, dispatches a provider, creates a
Task, or grants Git, process or tool authority. CLI data alone attests no
physical passkey; local CLI use keeps the OS user's existing authority.

## Operations

Every consumer command takes the explicit physical `--workspace`, never
inferred from cwd. For `answer apply` it must also equal the request's immutable
`workspace`; a mismatch is refused (`not_recorded`, `human_input_request_mismatch`)
before anything is read or changed. Use `--format json` (`text` prints the same
JSON).

| Command | Effect | Result |
| --- | --- | --- |
| `ask --workspace WS --file ABS_DRAFT` | `human_input_question_publication` | `question` |
| `list --workspace WS --project ID [--state S] [--limit 1-20] [--cursor C]` | `read` | `list` |
| `show QUESTION_ID --workspace WS` | `read` | `question` |
| `answer preview QUESTION_ID --workspace WS --request-id ID (--option-id ID \| --text-file ABS)` | `read` | `preview` (contains `request`) |
| `answer apply --workspace WS --file ABS_REQUEST` | `human_input_answer_record` | `status` |
| `answer status REQUEST_ID --question-id QUESTION_ID --workspace WS` | `read` | `status` |
| `wait QUESTION_ID --workspace WS --timeout-seconds 1-600` | `read` | `wait` |
| `acknowledge QUESTION_ID --workspace WS --request-id ID --request-sha256 D` | `agent_reported_receipt` | `question` |
| `continue QUESTION_ID --workspace WS --request-id ID --request-sha256 D --summary TEXT` | `agent_reported_continuation` | `question` |
| `withdraw QUESTION_ID --workspace WS --reason TEXT` | `human_input_question_withdrawal` | `question` |

Exit and output:

- Success prints one result and exits 0. `wait` exits 0 for every outcome.
- A semantic refusal exits nonzero and prints one `error` result (or, for
  `answer apply`, a `status`). It carries the current `question` when known, so
  a refusal never invents a next actor.
- An argument error (`human_input_invalid_arguments`) exits nonzero, prints
  nothing on stdout and touches no workspace state.
- An exit code alone never proves that no earlier effect happened.
- One result is at most 4 MiB, also with every text field at its maximum.

Reads (`list`, `show`, `answer preview`, `answer status`, `wait`) create no
file, directory or lock. Writes hold one store lock only while deciding, never
while waiting.

## Identity and digests

- `question_id` is `hiq_` plus 32 hex characters. It is derived from
  (workspace, `project_id`, `publication_key`) and is stable across restarts
  and resumed agents.
- `question_sha256` binds the workspace, ID, key, `scope`
  (project/repo/Task/activity ID and title), `text`, `reason`, `answer_kind` and
  `options`. Source facts (`provider`, `session_id`, `origin_cwd`) and status
  observations are outside it and never change it.
- `ask` with the same key and the same content returns the original question,
  including its original source and `created_at_utc`. Other content under the
  same key is `human_input_publication_conflict`; a changed question needs a new
  key, and the original question is kept.
- A `request` binds `purpose: human_input_answer_v1`, the workspace,
  `question_id`, `question_sha256`, `scope`, the consumer's `request_id` and the
  exact `answer`.
- `binding_sha256` is the SHA-256 of the RFC8785 canonical request without
  `binding_sha256`. `request_sha256` is the SHA-256 of the complete RFC8785
  canonical request. Bind the physical confirmation to those exact bytes.
- Request identity is the pair (`question_id`, `request_id`). There is no global
  request register.

## Text and limits

Limits are in Unicode codepoints:

| Field | Limit |
| --- | --- |
| `text` (the question) | 2000 |
| `reason`, option `consequence`, `continue --summary`, `withdraw --reason` | 1000 |
| option `label` | 120 |
| `activity_title` | 160 |
| text answer | 4000 codepoints and 16000 UTF-8 bytes |

Choice questions have 2-6 options with unique IDs and labels. Text questions
have an empty `options` array.

- Blank means empty or only Unicode whitespace, and is rejected. Invalid UTF-8,
  lone surrogate escapes, duplicate or unknown JSON members, and NUL are also
  rejected.
- Agent-written prose allows tab and newline as its only control characters.
  Labels and titles are single lines.
- A text answer is stored and returned exactly as written: no trimming, Unicode
  normalization or line-ending change. `--text-file` bytes are the answer.
  Treat answer text as data: never as shell or a filename.

## States and next actor

`waiting_for_answer -> answered -> agent_received -> agent_reported_continued`.
`withdrawn` is terminal for an unanswered question. `next_action` concerns this
question only:

| State | Actor |
| --- | --- |
| `waiting_for_answer` | `human` |
| `answered`, `agent_received` | `agent` |
| `agent_reported_continued`, `withdrawn` | `none` |

The progress of the whole assignment and any later human steps are not tracked
here.

Receipt, continuation and withdrawal have `provenance: reported_by_agent`. They
are the agent's own reports, not native process observation, not the schema 2
`continuation_completed`, and not Task, Epic or main delivery. There is no
liveness attestation. A stopped agent is not woken; the stored answer persists
until a resumed agent reads it.

## Answer outcome and recovery

`status.state` is one of the following:

| State | Meaning | Fields |
| --- | --- | --- |
| `recorded` | This request is the question's one answer. | `request_sha256` and `winner` are set. |
| `other_winner` | Another request is the answer. | `winner` names it; `human_input_request_conflict` when the same `request_id` names other bytes. |
| `not_recorded` | The question is readable and this request is not its answer. | `reasons` says why when `apply` refused it before any change. |
| `unavailable` | The question or workspace cannot be found. | |
| `unknown` | Stored state is unreadable or inconsistent. | |

`unavailable` and `unknown` are not `not_recorded` and open no new write.
Nothing is repaired silently.

- Applying the same request again is idempotent. Another request never
  overwrites the winner.
- After a lost reply, read `answer status` for the same pair.
- The same immutable request may be applied again, with fresh confirmation, only
  when all of these hold: authoritative `not_recorded`, a fresh `show` in
  `waiting_for_answer`, and no winner. Each grant permits at most one apply.
  Never create a new request to retry.
- Requests have no native expiry. Consumer passkey or page validity is separate
  and never withdraws the question.
- When a project, repository or Task leaves the registers, new previews and
  applies are refused (`human_input_scope_unavailable`). History stays readable
  and listed.

## Discovery

`list --project ID` covers every state in immutable `question_id` order. It does
not age questions out.

- `--limit` is 1-20 (default 20). Pass `next_cursor` (opaque) as `--cursor`
  until `complete` is true.
- `complete` means the page scan finished successfully. It is not a
  transactional snapshot. A new poll starts from the beginning and also finds
  questions published meanwhile.
- A failed scan returns an `error` and no rows. A missing row never means
  withdrawn or answered.

## Publication scope

`ask` validates the following against the registers, and never creates or
rewrites a registration:

- The project is registered.
- The optional `repo_id` belongs to the project.
- The optional existing `task_id` belongs to the same project and repository,
  and then `repo_id` is required.
- `origin_cwd` is inside the project's wrapper, repositories, companions or
  registered Epic or Task worktrees.

A named activity without a Task is valid.

## Error codes

All codes start with `human_input_`:

`invalid_arguments`, `invalid_draft`, `invalid_request`, `invalid_answer`,
`scope_invalid`, `scope_unavailable`, `publication_conflict`, `not_found`,
`request_mismatch`, `request_conflict`, `already_answered`, `withdrawn`,
`not_answered`, `answer_mismatch`, `not_acknowledged`, `report_conflict`,
`state_invalid` and `outcome_unknown`.
