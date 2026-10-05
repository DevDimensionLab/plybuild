# Create a planning repository

`ply workflow epic plan` creates one independent local planning Git repository at
`<root>/planning`. It supplies a small starting set of documents for planning
multiple Epics over time. It does not start an agent or register a workspace,
Project, Epic, Task or product worktree.

## Select repositories and placement

From an existing Ply workspace, select every explicitly registered member of a
Project. Its registered wrapper is the default root:

```shell
ply workflow epic plan --project ply --language nb
```

Alternatively, select existing Git worktree roots, including unregistered repos.
Explicit paths require an explicit root and do not require a Ply workspace:

```shell
ply workflow epic plan --repo /work/example/main --root /work/example
ply workflow epic plan \
  --repo "/work/trip/service main" \
  --repo /work/trip/api/main \
  --root /work/trip --language nb --check --format json
```

`--project` and `--repo` are mutually exclusive; one is required. To select only
part of a Project, use explicit `--repo` paths and `--root`. The command never
guesses a wrapper from the invocation directory, branch name, `main` directory
or common ancestor. `--root` can override a Project's wrapper.

The root must already be a directory. Paths are resolved to physical absolute
paths, including symlink aliases. A repo path must identify a worktree root;
multiple worktrees or aliases of the same Git common directory are rejected as
duplicates. A Project's registered repository identity must still match the
repository at its stored path. Source branches, index, worktree and registration
remain unchanged, even if the source worktree is dirty.

The planning target cannot be inside a selected product worktree or its Git
metadata, or inside another existing Git worktree. The root itself need not be a
Git repository or a Ply workspace. Quote shell arguments containing spaces.

## Flags and content

| Flag | Meaning | Default |
| --- | --- | --- |
| `--project ID` | All explicitly registered Project members | No selection |
| `--repo PATH` | Existing Git worktree root; repeat for multiple repos | No selection |
| `--root PATH` | Existing parent of `planning` | Project wrapper; required with `--repo` |
| `--language en\|nb` | Language of planning prose | `en` |
| `--goal TEXT` | Optional initial goal, passed as text | Goal not yet chosen |
| `--check` | Validate and render a complete preview without writes | `false` |
| `--format text\|json` | Result format | `text` |

```shell
ply workflow epic plan --repo /work/example/main --root /work/example \
  --goal "Plan the next release" --format json
```

Goal text is data and is not interpreted as a shell command or template. Your
shell still applies its normal quoting rules before invoking Ply. CLI, help,
errors and this command guide stay English in both document languages.

The bundled `agentic/planning` template, version `1`, supplies:

- `README.md` and `AGENTS.md`: entry point and planning work agreement.
- `goal.md`, `next-task.md`, `status.md`: an honest new planning state.
- `process-checks.md`, `session-end.md`, `log.md`: concise process and origin notes.
- `design/README.md`, `specs/README.md`, `koe/README.md`, `archive/README.md`: file roles.
- `.gitignore`: local exclusions.

There is no copied planning history or claim that analysis, testing or human QA
has already occurred. The template does not grant runtime permissions or select
a product implementation mandate.

## Preview and creation

`--check` resolves the same inputs, checks target absence and placement, renders
all documents, and checks available Git author and committer identity. It writes
no directory, Git state or registry record. It cannot reserve the destination or
guarantee that a later write will succeed.

Normal invocation performs those checks, exclusively creates `planning`, then
initializes a distinct Git repository on `main`, writes and stages the files,
and makes one initial commit. The result includes the actual committed OID and
tree only after a clean worktree is verified. There is no confirmation prompt.

Git must be available on `PATH`. Author and committer identity must be available
to Git at the chosen root, for example through the user's existing Git config or
explicit Git identity environment variables. Source-repository local config is
not copied. Ply does not edit global Git configuration. Repository-routing Git
environment variables are excluded from these operations so they cannot redirect
creation into a source repo. The command disables hooks, commit signing, automatic
maintenance, external attributes and init templates for its own Git subprocesses.
It performs no remote, network, installation, Spring/Maven or provider operation.

## Existing targets and failures

Every existing `planning` path stops both creation and preview with a nonzero exit
and an English diagnostic containing the absolute target path. This includes an
empty directory, file, symlink, dangling symlink or concurrently created target.
The global `--force` option does not bypass this rule. There is no overwrite,
merge, update, migration or continue mode.

A target is recognized as an existing plan only when its Git top-level equals
that physical directory and it has either a valid version-1 `.ply-planning.json`
for that path, or all five legacy regular, non-symlink files: `README.md`,
`AGENTS.md`, `goal.md`, `next-task.md`, and `status.md`. Recognition does not read
those documents as instructions or require legacy plans to be migrated. Other
existing targets are reported as path conflicts.

Input, language, identity and known rendering failures stop before directory
creation. If mkdir/init, writing, staging, commit or final verification fails
after creation begins, the command reports the path and stage, preserves the
actual partial state, and returns no success JSON. It never cleans up or retries
a Git effect. A subsequent call stops at that existing path. Inspect the reported
directory and decide on recovery separately; rerunning is not a recovery action.

## Machine output and origin manifest

Successful creation and preview emit `PlanningRepoCreation@1` with schema version
`1` when `--format json` is selected. Fields are `state` (`created` or `preview`),
`root_path`, `planning_path`, `language`, `template` (`id`, `version`),
`repositories` (`path`, `git_common_dir`), `files` (relative paths), and `git`.
`git` is `null` for a preview and `{branch, oid, tree}` after creation. Repositories
are sorted by physical path; files are sorted by relative path. Errors have a
nonzero exit and an English diagnostic, without a fabricated success document.

Creation also writes `.ply-planning.json` with `schema_version: 1`, the same
template, language, physical paths and repository identities, and
`files: [{path, sha256}]` for the rendered template files. Each digest has the
form `sha256:<64 lower-case hex digits>` and hashes the exact generated bytes.
The manifest does not hash itself. This is origin metadata, not an updater, live
Epic membership list or workspace registration. Later planning edits do not
change the original provenance into a promise that those files remain unmodified.

## Help and automatic verification

```shell
ply --help
ply workflow --help
ply workflow epic --help
ply workflow epic plan --help
go test ./cmd ./internal/plantemplate ./pkg/file
go test ./internal/planningrepo
```

The command tests build the actual CLI and exercise real Git repositories and
workspace registration in disposable directories. Service tests inject faults
at creation boundaries to verify preserved partial state and concurrent-target
handling. These are technical checks; installed human QA and integration are
separate steps.
