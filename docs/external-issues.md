# Offline external issue references

Epics and Tasks can optionally link any number of Jira or GitHub Issues.
These are supplied metadata, always `offline_unverified`: Ply never contacts a
provider, checks existence, reads credentials, synchronizes status or writes Git.
An issue may be owned by only one Epic or Task across the entire workspace,
including different Projects. Ordinary work and delivery need no external issue.

```sh
ply workspace epic external-issue add my-epic --provider jira --key PLY-123 \
  --url http://jira.example.invalid/jira/browse/PLY-123 --actor human
ply workspace epic external-issue show my-epic --format json
# Copy the returned ext_... ID and current workspace revision (1 in this example).
ply workspace epic external-issue primary my-epic ext_ID --expected-revision 1 --actor human
ply workspace task external-issue add my-task --provider github \
  --host git.example.invalid --owner owner --repository repo --number 7 \
  --url https://git.example.invalid/owner/repo/issues/7 --actor human \
  --primary --expected-revision 2
ply workspace task external-issue list my-task --format json
ply workspace task external-issue primary my-task --clear --expected-revision 3 --actor human
ply workspace task external-issue unlink my-task ext_TASK_ID --expected-revision 4 --actor human
```

Replace IDs and revisions with actual readback values. `list` aliases `show` for
one selected work item. Both commands show refs, explicit primary, effective
commit key and origin, filtered change history and the shared workspace revision.
Normal Epic/Task `show` also projects linked metadata; absent refs preserve the
existing views. Dedicated metadata reads do not inspect a source checkout or
require it to exist or be clean.

No first ref becomes primary automatically. `--primary` on add explicitly
selects that ref. Task primary wins; otherwise its parent Epic primary is the
fallback; otherwise effective is null. Keys are `PLY-123` and `owner/repo#7`.
Clearing or unlinking the selected Task primary exposes any Epic fallback.
Unlink names exactly one ref and releases its workspace ownership, retaining
history. Primary is only a read-only policy input: it does not enforce or rewrite
commit subjects, install hooks, or authorize PRs or merges.

## Normalization and limits

Only absolute HTTP and HTTPS URLs up to 2048 characters are accepted. Scheme and
ASCII DNS host case are folded, a final DNS dot is removed and the scheme's
default port (`80` or `443`) is omitted. Other numeric ports are normalized to
1..65535. IPv4 hosts are accepted as written; IPv6 and international hostnames
must use another supported DNS authority (punycode is ASCII). Paths accept
bounded ASCII unreserved segments, decode percent aliases once and remove one
trailing slash. Empty/dot segments, encoded separators, backslashes, controls,
credentials, queries (including empty `?`), fragments and opaque URLs are refused.
No routing, redirect, provider or issue existence is guessed.

Jira uses `/<optional-context>/browse/KEY` routing. Other routes are unsupported
and fail explicitly. Keys are folded to uppercase and must match
`[A-Z][A-Z0-9_]*-[1-9][0-9]*` within 128 characters; the last decoded URL segment
must agree. The case-sensitive base context before `/browse` is carried as
`jira_context_path` and participates in identity, together with normalized
host/port and key. Thus `/jira-a` and `/jira-b` on one authority are distinct
instances. Jira URLs cannot carry GitHub identity fields.

GitHub requires explicit host, owner, repository and a positive signed 64-bit
issue number. The URL must be exactly `/owner/repository/issues/number` and
agree with all supplied fields. Owner and repository case are folded to lowercase;
leading decimal zeros in the URL number normalize away. Owner is 1..100 ASCII
letters/digits/hyphens starting with a letter or digit; repository is 1..100
letters/digits/dots/underscores/hyphens excluding `.` and `..`. Pull request routes
are unsupported. Identity is provider, normalized host/port, owner, repo and number.
Host input must agree with the normalized URL authority. Enterprise authorities
are supplied per ref; there is no provider registry.

HTTP and HTTPS URLs with the same normalized authority/context/key share issue
identity; scheme remains in the canonical URL. A retry with equivalent normalized
fields is a no-op. Different canonical supplied fields for an already linked
identity are a conflict, including a different retained scheme. Ply never silently
replaces a selected issue or its URL. Read current metadata and explicitly unlink
before moving or changing a binding.

## Publication and compatibility

The workspace owns `.ply/external-issues.json`, an optional strict version-1
metadata event store. Absence means empty refs, null primary and revision zero;
reads create nothing. First explicit mutation creates version 1. Unknown versions,
fields, duplicate keys and invalid event histories are refused. There is no global
migration or registry format change; existing registry formats and historical
execution inputs are retained. Old readers can continue their ordinary registry
operations but do not understand or enforce external issue ownership. Use this
reader for issue metadata; do not manually edit the sidecar or mix incompatible
metadata writers. One existing work-item writer lock serializes validation and
atomic replacement of the complete event history across all work items.

`--expected-revision` is optional for additive non-primary links, which serialize
safely under that lock. It is required for primary changes and unlinking. It guards
the shared metadata revision, including changes on other items. An exact request
that already has its desired state writes nothing even after the revision advanced.
A stale request that would change state is refused; it cannot restore a primary
cleared or changed by someone else. Conflicts are ordinary actionable request
errors, with no Task/Run stop or persistent reservation. There is no preview/apply
gate. After an interrupted publication, inspect the same item's metadata and
repeat the same request only when its effect is known; no new request identity or
blind replacement is needed.

The store uses the existing 16 MiB metadata limit and durable replacement helper.
Its history stores actor claims, sequence, timestamp and event digest; claims are
not authentication or human QA. Metadata changes leave identity registries, queues,
Goals/Runs and Git untouched, even for live/dirty/unavailable source checkouts.

Run the bounded disposable fixture with either a make-built or installed binary:

```sh
python3 test/external_issues_cli.py /absolute/path/to/ply
```
