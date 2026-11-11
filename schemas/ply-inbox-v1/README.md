# Ply Inbox producer contract v1 (`inbox_v1`)

The public contract of `ply inbox`: Ply stores the canonical full agent reply,
its title and the per-item read marks of one owner in one registered project.
plyapp-server is the only intended writer in v1 and keeps no mark store of its
own. Mode `inbox_v1`, `schema_version: 1`.

`ply capabilities --format json` advertises `inbox.put`, `inbox.show`,
`inbox.list` and `inbox.mark` in the always-present top-level array
`inbox_extensions`. Treat a missing entry or another decoder version as
`source_unavailable`, never as an empty Inbox.

## Fixtures

Every file except `put-request.json` and `derived-title-cases.json` is produced
by the real store with a fixed clock (`go test ./internal/inbox ./cmd`), and the
tests decode each one back with the production types:

| File | Content |
| --- | --- |
| `put-request.json` | `InboxPutRequest@1` (`ply inbox put --file`) |
| `message.json` | stored `InboxMessage@1`; its `payload_sha256` is the golden RFC 8785 vector of the request |
| `item.json` | `InboxItem@1` from `show`, read |
| `list.json` | `InboxList@1`, limit 2 of 3 items: agent and derived titles, `next_cursor`, counts `filename_unverified` |
| `list-unavailable.json` | `InboxList@1` when `read/` cannot be read: every `read` and `unread_count` null |
| `mark.json` | `InboxMarkResult@1` of the first `mark --value read` |
| `error-conflict.json` | `InboxError@1` of a replay with a different body |
| `capabilities-excerpt.json` | the exact `inbox_extensions` array |
| `derived-title-cases.json` | `derived-title-v1` input/output table |

## Rules a consumer relies on

- Identity is (`conversation_id`, `message_id`). The immutable payload is
  `InboxPayload@1` = project, owner, identity, canonical `replied_at_utc`
  (nine fraction digits), `title`, `body`, `truncated`, `full_sha256` and
  `answer_source`; `payload_sha256` is `sha256:` + hex SHA-256 of its RFC 8785
  bytes. `agent_label`, `provider` and `write_kind` are first-write metadata.
- `put` gives `created` or `unchanged`; both confirm durable storage of exactly
  this payload. A different payload under the same identity is
  `inbox_conflict` with `conflicting_fields` and both hashes, never values.
- `title` is `null` or one line: a control character or line separator anywhere
  is rejected; only surrounding whitespace is trimmed, so `" Hello "` and
  `"Hello"` are the same payload.
- A stored message file that cannot be read is `inbox_item_invalid` (not
  retryable) in `put`, `show` and as a `list` `invalid` entry.
- `inbox_busy` and `inbox_io_error` are retryable with the same payload; every
  other class is not. Workspace project and path classes pass through.
- `list` orders by `replied_at_utc` descending, then `conversation_id`, then
  `message_id`. The cursor is opaque and bound to project and owner.
  `unread_count_exact` is true only when every counted message was validated in
  this call; otherwise the count is a lower bound from file names
  (`filename_unverified`). A corrupt mark never reads as read.
- Marks are never pruned. Only `mark --value unread` removes one.

Nothing here starts, steers or answers agents, sends notifications, or uses the
network.
