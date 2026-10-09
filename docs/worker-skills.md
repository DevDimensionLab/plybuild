# Pinned worker role guidance

`internal/workerskills` embeds Addy Osmani's selected skills and their referenced
files in the Ply executable. Worker start and correction need no Addy checkout,
plugin installation, or network access to obtain instructions.

The source is [`addyosmani/agent-skills`](https://github.com/addyosmani/agent-skills)
at commit `1401c8b8030e023baeebb31781a6653fe8e93026`. All ten files under
`internal/workerskills/upstream/` are unmodified bytes from that commit: five
skills, four references, and the complete MIT license with Addy's copyright.
No upstream `AGENTS.md` or other repository-level agent instructions are imported.

## Roles and controller integration

| Role | Upstream guidance |
| --- | --- |
| `planning` | Planning and task breakdown; definition of done |
| `developer` | Planning, incremental implementation, TDD; definition of done and testing patterns |
| `correction` | Developer material plus debugging and error recovery |
| `reviewer` | Code review and quality; security and performance checklists |

Each bundle starts with the common Ply adapter and its role adapter. These
explicitly preserve the native Goal and Spec, controller ownership of commits
and native state, bounded automatic correction, and independent read-only review.
An approved task does not acquire another human approval checkpoint from an
upstream example. The planning material within developer/correction attempts
organizes implementation without creating a separate task tracker or stopping
for plan approval.

Reviewer mutation experiments use a separate authorized scratch copy without
Git metadata. Candidate source and evidence remain read-only. All roles preserve
the chosen model, process/write boundaries, disabled plugin/MCP/network surfaces,
and the prohibition on worker subdelegation. The adapter does not grant access;
the controller's runtime policy must still enforce and probe those boundaries.

The package API is:

```go
func Bundle(role string) (string, error)
func Manifest() (ManifestData, error)
func Verify() error
```

Use `RoleDeveloper` for the first implementation attempt, `RoleCorrection` for a
correction after a failed test/review, and `RoleReviewer` for the separate review
session. `RolePlanning` is available for a controller-owned planning attempt.
Role names are exact; unknown values fail instead of selecting a fallback.

The controller prepends the complete `Bundle` result to its task-specific prompt
**before** computing the prompt and attempt hashes. Preserve the complete prompt
in the native attempt evidence. Keeping only a role name, file path, or digest
does not preserve what the worker actually received. Bundle errors must stop
launch; they must not silently omit the skills. The package owns no Task/Run,
queue, process, provider, or delivery state.

## Provenance and integrity

`internal/workerskills/manifest.json` records the source repository and commit,
MIT license path, every embedded file's path, byte count, and SHA-256, plus ordered
role selections. Each upstream entry also carries its original source path.
The five Ply-owned adapters have separate entries and hashes.

`ManifestSHA256` pins the exact manifest bytes in Go source. `Bundle`, `Manifest`,
and `Verify` check this pin, all file hashes and sizes, and the embedded file set.
Missing, changed, or extra files fail closed. `Manifest` returns fresh data, so a
caller cannot mutate later bundles through its result. The emitted prompt
contains the adapter, manifest, license, and exact selected source bytes in a
deterministic order. Relative reference links identify included sections; they
do not require paths in a live checkout. Unbundled named skills are optional
background and do not trigger downloads or broaden runtime permissions.

These checks detect an inconsistent package. They are not a signature or an
independent trust root: qualifying the executable and reviewing updates remain
the controller/release workflow's responsibility.

To update the pin, review a specific upstream commit and retrieve each selected
file with `git show <commit>:<source_path>`. Preserve upstream bytes and license,
include every referenced local file needed by a selected skill, and review the
adapter against any changed instructions. Update manifest provenance, selections,
byte counts, and SHA-256 values, then update `SourceCommit` and `ManifestSHA256`.
Adapter-only changes also require new file and manifest hashes. Never regenerate
these pins during worker startup or from mutable runtime files.

Run the focused contract and integrity tests with:

```sh
go test ./internal/workerskills -v
```

They check exact role content and relative references, MIT provenance, corrupted
or missing material, accidental imported repository instructions, immutable
caller results, and execution of the compiled tests from an empty working
directory with no user environment. Delivery qualification must additionally
check the actual developer/reviewer prompts and runtime policy in a real worker
attempt; package tests alone do not prove that integration.
