# Quality Upgrade Handover

Generated: 2026-09-01T15:46:46+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P7 toolchain implementation commit:
  `16ecb67eebb6450226d3264f30a64a14391b5245`.
- Its exact parent is the launch continuity commit
  `9852ed15b6e22ec0263bc4f81a973f4e088d7cf5`, whose exact parent is the P6
  implementation `097a9f15782c773e47a61d903f5d09f5c96408f0`.
- The implementation tree is
  `94262fbfd426692fdc44f9eb21f85065ed40a6e2`.
- After this handoff, obtain the new continuity HEAD with
  `git rev-parse HEAD`; its exact parent must be `16ecb67`.
- The implementation and every accepted measurement had empty ordinary and
  ignored status. No push, merge, publication, release, stash, revert,
  successor launch, retained-evidence deletion, image deletion, or worktree
  removal occurred.

## Continuity Checkpoint

P2A-P6 are complete. P7 remains active after its bounded maintained-toolchain
baseline move; dependency groups remain. P8 is queued.

The answered toolchain archive links reciprocally to exactly one NEXT archive
for the first small dependency group. Only the launcher's mutable header and
prompt regions changed; its stable executable skeleton remains byte-identical.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. Downloaded tools,
module/build caches, reports, generated files, build contexts, and audit
evidence remain external.

## Maintained Toolchain Decision

The selected baseline is exact Go 1.26.7. Go 1.27.0 was rejected for this move
because pinned golangci-lint 2.12.2 was built with Go 1.26.2 and its official
support policy does not claim targets newer than the build Go line. Pinned
GoReleaser 2.17.1 was built with Go 1.26.5. Updating either quality tool to
admit Go 1.27 would have mixed dependency upgrades into the toolchain move.

Primary Go releases/toolchain/module documentation, pinned golangci-lint and
GoReleaser module/release metadata, setup-go v2's official manifest, Docker
Official Image source, and registry manifests are recorded in
`.quality/baseline/toolchain-migration.json` and the P7 roadmap record.

The retained official candidate is
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`:

- version: `go version go1.26.7 darwin/arm64`;
- executable SHA-256:
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
- official archive SHA-256:
  `020a1e82d85d1df183e0d3c2c6ddf181c295327fca19b3d8d614855b6cb49f6d`.

Contractual declarations now agree:

- `go.mod` retains the characterized language floor `go 1.18` and adds exact
  preferred `toolchain go1.26.7`;
- the Docker builder is
  `golang:1.26.7-alpine3.24@sha256:28d89ee9cc0ff9fec75c82ca201e6bf7fdf9a679d4b7b24dfa04f2bb766bb468`;
- the deactivated release workflow and active README name 1.26.7; and
- the exact baseline instrument identity names 1.26.7.

The deactivated lint workflow has no separate Go selector. The inactive Snap
publisher remains a P8 scope decision. Raising the module language floor to
1.26.7 was rejected because it introduced ten pinned-lint findings from newer
language-analysis semantics; retaining the documented `go`/`toolchain` split
keeps this move behavior- and toolchain-only.

## Baseline Reproduction And Compatibility

Old Go 1.26.2 reproduced the stored pre-P7 scorecard at
`/private/tmp/ply-p7-baseline-old-go1.26.2.iHSEG4`, SHA-256
`5fb3226009cfbf0d29f63fa03592157cce4efcec6e38583b64a86f6288e89490`.
Go 1.26.7 reproduced the migrated scorecard at
`/private/tmp/ply-p7-baseline-new-go1.26.7.xY13SW`, SHA-256
`9f044510ae95d1df85b6c0868323c11453f97015c16fc46721d70c9da0b3a463`.
Only `tool.go_build.version` differs: the raw body, criteria, denominators, and
all 228 numeric debt leaves are identical.

Both toolchains select the same 233 modules and 143 non-standard packages,
emit byte-identical core help, and preserve the same existing 207-line
`go mod tidy -diff` projection. No direct or indirect version or `go.sum`
entry changed. API/CLI and entry/subprocess compatibility, generated archive
and Docker contracts, and all host behavior remain green.

Audit-meta T15 composes old/new instrument reproduction with the pre-P7/current
Go scorecards. It requires byte-exact authoritative baselines and compares all
228 numeric leaves across both axes.

## Accepted External Evidence

The commit-bound schema-2 review root is
`/private/tmp/ply-p7-quality-review.16ecb67.Ugm5Bh`. Its verified 22,002-entry
manifest SHA-256 is
`9fc89ac3d4296c2ad7f9a30b6f9decc308113a6b5008afecd12b741983a7e6eb`.
The canonical evidence SHA-256 is
`f3ea505bc6bd414f74bcdbc5b5f6b747bfadf7dc48ddc16ad9977937b0c4958d`;
all six manual receipts validate and the focused audit exits 0.

The exact complete apparatus root is
`/private/tmp/ply-p7-quality-gate.16ecb67-final.hdvr9q`. Its verified
246,560-entry manifest SHA-256 is
`15136fa28368eb4ada0b81dfddc4f36daa19f17b7ab00d3a4ea85d496066cb8e`.
The exact 21-line stage ledger ends with `audit:Q0.*,Q1.*,Q2.*`, and
`make quality` exits 0. The Q0-Q2 scorecard SHA-256 is
`0cf6f16797ec73ed8d51944a77dc0cde4de85e174526b915afb58542e98d4cec`:
all 27 rows pass at L2, manual evidence is valid, mutation and acceptance
populations are 8/8 and 4/4, and held, regressed, current-not-comparable, and
dirty counts are zero.

The independent accepted regression root is
`/private/tmp/ply-p7-regression-gate.16ecb67-final.4nMFvg`. Its verified
67,726-entry manifest SHA-256 is
`48be0c3b828c60702ff4f06c513b8720dba4259c9df39d99cf23b0f655dffded`.
All 36 ledger stages pass, including declaration/graph contracts,
compatibility, pinned lint, tests/race/vet, launcher and Make contracts,
complete preflight, host/snapshot/Docker meta and acceptance, standalone audit
meta, focused and Q0-Q2 audits, and empty-HOME count-2. The focused and Q0-Q2
scorecards are
`00ff8e287212841d18de0bb7107aa7af49a58f105f7280597059b09e7a0b4fab`
and
`0cf6f16797ec73ed8d51944a77dc0cde4de85e174526b915afb58542e98d4cec`.
The full scorecard is
`980caa86d74519abec73de7a3896b77a5c9c51c7754c3e9f901ee9ea4f9d3960`;
it exits the expected 1, never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7.

The complete cross-root index is
`/private/tmp/ply-p7-complete-move.16ecb67.tGLtNi/evidence-index.json`, SHA-256
`68e9a7cf649a02b882703bd079d37fafab69caba166a97efa164616308b9578b`.
Its verified one-entry manifest SHA-256 is
`8da5e8a79e26cbb8e4279150f820c8e0bf88da1fb032860015b079d8da718607`.

For exact replay, put the candidate SDK directory first in `PATH` as well as
passing its absolute `GO`; T15's historical instruments invoke literal `go`.
Do not inject ambient `GOFLAGS`, because the audit meta-contract intentionally
changes that variable in its probes. Warm a fresh external `GOMODCACHE` from a
Git archive outside the worktree: `go mod download all` adds 208 historical
checksums if run in the real tree. Diagnostic roots from those characterized
environment mistakes are not accepted evidence.

## Next Objective

Implement only the first small P7 dependency group: direct
`golang.org/x/term v0.5.0` and the coupled `golang.org/x/sys v0.5.0` selection
required by it. Verify current releases, module requirements, Go support, and
security information from primary sources before choosing versions.

Capture old/new module graphs and `go mod tidy -diff` projections in external
state first. Admit no unrelated direct or indirect version change; if minimum
version selection requires a wider group, stop and record the decision instead
of broadening scope. Keep the Go 1.26.7 declarations and exact baseline
identity unchanged unless a measured incompatibility requires stopping.

Make one dependency-only implementation commit. Prove pinned lint,
tests/race/vet, API/CLI and entry/subprocess compatibility, generated artifact
and host/snapshot/Docker behavior, audit reproduction, exact Q0-Q2 L2, the
separate queued-L3 result, and empty-HOME count-2 from clean external state.
Do not change production behavior to accommodate an upgrade in the same move.

## Start And Stop

Confirm branch, exact ancestry, clean ordinary and ignored status, reciprocal
links, launcher `--check`, the P7/P8 queue, implementation commit/tree, and all
accepted manifests before editing. Read the active archive, this handover, the
P7 roadmap, module graph, dependency callers/tests, toolchain contract,
compatibility, snapshot/Docker, and audit contracts.

Stop before any dependency outside the terminal pair, behavior/API/CLI change,
quality-tool upgrade, P8 domain work, inactive packaging work, publication,
publisher/registry/credential change, or release. Do not push, merge, publish,
release, delete retained evidence or images, stash, revert, launch a successor,
or remove the worktree.
