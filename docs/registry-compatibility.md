# Shared work-item registry compatibility

The work-item registry is shared by all projects in a workspace. Ply strictly
decodes and validates it before filtering a project. Unknown authorization
fields are rejected; ordinary writes preserve IDs, Spec/Goal references, Task
results, integration/acceptance evidence, queue events and Epic-base history.

| Reader/writer | Reads | Ordinary writes | Format transition |
| --- | --- | --- | --- |
| Historical installed CLI, SHA `25006d60ef7b7d7b623373295edd72cbc8144e5bb2845f9ea54cd50ebf926ec7` | Historical formats 1–4 within its older field set; rejects newer acceptance fields | Its older contract only | Update executable before using the reported shared registry |
| Private control built from `c401914` | Formats 1–4, including AA-01 acceptance fields | Same legacy format contracts | Diagnostic/bootstrap tool; not an ordinary delivered release |
| CB-01A delivered CLI | Formats 1–5, including existing format-4 acceptance and amendment history | Keeps the current format; strict validation and native history guards remain active | Existing acknowledged 1/2→3 and 1/2/3→4 operations; explicit delivered-reader-backed 4→5 below |
| Future format greater than 5 | Rejected as `registry_unsupported_format` | No writes | Install a supporting release through the independent installer |

Formats 1–3 retain their documented native upgrades: Task creation/content uses
the observed `--upgrade-store`/`registry_upgrade` acknowledgment, and queue/base
operations use their exact `registry_upgrade`/`--upgrade-registry` contract.
Format 4 already contains the reported acceptance additions. CB-01A reads and
preserves that existing history without retroactively changing its meaning.
The incompatible historical addition cannot be fixed inside the old binary;
the bundled installer works without asking that binary to read a workspace.

Format 5 establishes an explicit release boundary. Its records have the same
meaning as format 4, with an immutable `format_boundary` identifying the exact
previous registry bytes and delivered compatible reader. New writers cannot
publish this format merely by setting a version or copying an old preview.
The store itself requires the native upgrade permit, unchanged registry bytes,
unchanged records, and a still-available package whose manifest and executable
identity match. Normal compatible writes need no additional confirmation.

Do not migrate a shared workspace while any supported active consumer still
uses a format-4-only reader. Install the qualified supporting release for those
consumers first. **The CB-01A consumer journey can remain on format 4**; reading
the reported registry does not require migration.

## Explicit transition

Use the permanent, already qualified release package and its trusted manifest
digest. The default and `--check` are read-only:

```sh
ply workspace registry upgrade \
  --package /absolute/releases/ply-v1.1.0-cb01a.1-darwin-arm64 \
  --sha256 sha256:<manifest-sha256> --check --format json
```

Inspect `WorkspaceRegistryUpgrade@1`, including the workspace, old registry
SHA-256, reader identity and byte-for-byte backup location. To authorize that
exact transition, rerun with `--apply --confirm sha256:<confirmation>` instead
of `--check`. No human product pass or Task integration authority is created.
The writer rechecks both the registry and package under the native work-item
lock, and checks again immediately before atomic publication. A stale preview,
competing write, missing reader or changed package is rejected. Re-preview after
a competing ordinary write; never reuse its obsolete confirmation.

After interruption, repeat the same package, manifest digest and confirmation.
Before replacement the old bytes remain intact. After replacement, retry reports
`already_upgraded` without another migration. Backups reside in the native
content-addressed `task-content/v1/backups/sha256` store. A file-sync/rename failure
cannot publish half a YAML file; an unconfirmed directory sync remains an error
until a retry observes the durable state. No automatic downgrade is provided.

## Diagnostics and evolution

Commands using `--format json` return one `WorkspaceRegistryDiagnostic@1` JSON
object on **stderr** on registry decode failures, with a nonzero exit and no
additional text on that stream. Stdout is empty or contains the command's single
ordinary result (for example an unsuccessful execute preview). Parse the streams
separately; a result on stdout does not override the nonzero exit:

```sh
ply workspace epic list --format json >result.json 2>diagnostic.json
```

`code` and `next_action` are stable fields:

| Code | Meaning | Supported action |
| --- | --- | --- |
| `registry_unsupported_format` | Valid YAML declares an unsupported positive integer format | Install a supporting prebuilt release; preserve data |
| `registry_unknown_field` | A supported format contains a field the strict reader does not recognize | Preserve evidence and correct the producer; never strip authorization fields |
| `registry_malformed_yaml` | Syntax is invalid, empty, or contains multiple documents | Run `ply version --format json`; retain the diagnostic, registry and backups for producer repair |
| `registry_invalid_content` | Known-format shape, type or semantic/reference validation failed | Run `ply version --format json`; retain the diagnostic and unchanged publication/backup for producer repair |
| `registry_upgrade_required` | No valid upgrade permit or available compatible reader | Inspect a fresh native upgrade preview with a verified package |
| `registry_upgrade_conflict` | Confirmation, history, registry or release basis changed | Re-read and reconcile; do not bypass the native write guard |

Syntax, duplicate-key/alias checks and integer format validation precede version
classification. A malformed file or unknown authorization field is never simply
assumed to come from a newer CLI. Read rejection changes no registry bytes.

A future change that introduces fields or semantics unreadable by a supported
reader must use a new explicit format boundary, preserve legacy decoding, and
ship its reader **before** enabling writes. Update the reader/writer matrix,
package capability declaration, store publication guard and transition acceptance
together. Do not relax `KnownFields`, reuse format 4/5 for an incompatible shape,
or publish data from a private development binary as a substitute for delivery.
