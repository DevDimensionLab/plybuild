# Quality Upgrade Handover

Generated: 2026-09-05T04:54:41+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest implementation remains dependency-only go-colorful commit
  `dc4f27c0319c6f023b5750df2148219a4cb32e60`, exact parent
  `f40b0322974eee388d5c064596826898e63f13ca`, clean tree
  `a5f3282ed71f9e5557718d63e5ea8d23020c44e2`. It changes only `go.mod` and
  `go.sum`.
- Properties evaluation began at documentation commit
  `469049fdb547850d4874579f937e92e0418c7e6d`, exact parent dc4f27c. Its
  documentation handoff commit must have exact parent 469049f. Relative to
  accepted go-cmp commit c314bcb, dependency metadata still changes only
  go-colorful v1.2.0 -> v1.4.1 and adds its exact checksum pair.
- The answered properties archive and the mattn/go-runewidth NEXT archive must
  link reciprocally. Only the launcher's mutable header and prompt regions may
  change during handoff. Ordinary and ignored status must end empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and
completed dependency groups through accepted
`github.com/lucasb-eyer/go-colorful v1.4.1`. The HCL, mousetrap, and properties
no-change/rejection decisions are closed. Go-cmp v0.7.0, Viper v1.16.0, Emoji
v2.2.14, ini v1.67.3, fatih/color v1.19.0, fsnotify v1.10.1, latest
gomarkdown `v0.0.0-20260824154242-13c5cf49db8d`, and HCL, mousetrap, or
properties selection changes remain rejected for their recorded floor,
loaded-behavior, release-qualification, self-test, exact-latest, or closure
decisions. P8 remains queued.

## Rejected Properties Group

Fresh Go proxy, checksum-database, GitHub, and upstream Git evidence establish
`github.com/magiconair/properties v1.18.11` as canonical latest stable. The
proxy lists exactly 34 stable versions: v1.0.0 through v1.4.2, v1.5.0 through
v1.5.6, v1.6.0, v1.7.0 through v1.7.6, v1.8.0 through v1.8.10, and v1.18.11.
The unusual v1.18.11 number is nevertheless an ordinary same-major stable
release. `@latest`, `go list -m -u`, the GitHub Release and tag, and default
branch `main` all resolve v1.18.11, with zero later commits. Upstream tags and
proxy versions match exactly; all 34 GitHub Release objects are stable, with
no draft or prerelease objects. There are no prerelease tags, retractions, v2
tags, or major-path changes.

Releases v1.0.0 through v1.8.1 have no Go directive. Versions v1.8.2 through
v1.8.6 declare Go 1.13. Selected v1.8.7, all later v1.8 releases, and
canonical latest v1.18.11 declare Go 1.19. Thus the selected version itself is
unresolved compatibility debt, canonical latest and every higher release
exceed the retained Go 1.18 floor, and v1.8.6 is the highest floor-compatible
stable release.

| Version | Annotated tag / peeled commit / commit UTC / tag UTC | Go | Primary signature status |
| --- | --- | --- | --- |
| `v1.8.6` | `e97104ee21cdf9cc3ac33992fc829cdc3664d238` / `869a5592420f4ff6ebf74500b5c658b29d973172` / 2022-02-23T08:50:35Z / 2022-02-23T08:51:39Z | 1.13 | GitHub tag `bad_email`; commit valid |
| `v1.8.7` | `44cd486a2d048ebef0b05852be982ecb8b5adabe` / `c9a06e8f8f0164e4e16c0d5c4793cbed4ac90264` / 2022-12-08T14:25:49Z / 2022-12-08T14:35:14Z | 1.19 | GitHub tag and commit `bad_email` |
| `v1.8.8` | `46e4ca92331d52456c0b19d476b75e99d7aa0798` / `e64054d847ae0adba8535c3d174ba8deae26fc3c` / 2024-12-07T17:53:51Z / 2024-12-07T17:57:00Z | 1.19 | GitHub tag and commit valid |
| `v1.8.9` | `c95b8ed87e68403841499adea6cc489cd905b3f6` / `d8bdba35b511a72d4c00a47e801dc703328198e8` / 2024-12-07T18:36:23Z / 2024-12-07T18:37:57Z | 1.19 | GitHub tag and commit valid |
| `v1.8.10` | `c2c94dc0cff78d930813a4be45fdffe320d2f45d` / `281f515b93cf755020c587f2ffcbdee2322f3615` / 2025-04-09T19:30:57Z / 2025-04-09T19:31:22Z | 1.19 | GitHub tag and commit valid |
| `v1.18.11` | `0be6f60c94bd49cff1f6c1aa9423ac0d8b9f85d8` / `33f415127a06719be057e464f58f9885c30657aa` / 2026-07-23T07:36:47Z / 2026-07-23T07:41:17Z | 1.19 | GitHub tag and commit valid |

Local verification of older OpenPGP-signed objects was unavailable because
`gpg` is absent. Newer SSH-signed objects lacked a configured allowed-signers
file locally, so the recorded GitHub verification is the available primary
signature status. All six modules have no requirements, and each proxy archive
is byte-identical to its peeled tag-commit archive. Relevant checksum pairs
are:

- Highest compatible v1.8.6:
  `h1:5ibWZ6iY0NctNGWo87LalDlEZ6R41TqbbDamhfG/Qzo=` /
  `h1:y3VJvCyxH9uVvJTWEGAELF3aiYNyPKd5NZ3oSwXrF60=`.
- Selected v1.8.7:
  `h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=` /
  `h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=`.
- Canonical latest v1.18.11:
  `h1:j5ozYZl0zCjG7ahMDH0GWIobOvvUzT0BdAguG0ViKy0=` /
  `h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=`.

The v1.8.6 -> v1.8.7 history has 12 commits across 19 files (+222/-140): it
includes the faster Merge implementation, the Go 1.13 -> 1.19 floor change,
and migration from Travis to GitHub Actions. The seven commits to v1.8.8 add
LoadReader and remove `ioutil`; the five to v1.8.9 skip ignored private fields
and remove old material; the five to v1.8.10 add non-panicking 32-bit numeric
getters and leading-whitespace escaping; the four to v1.18.11 strip UTF-8 BOM
and exercise Go 1.25/1.26 in CI.

Exact `go get github.com/magiconair/properties@v1.8.6` proves that the highest
compatible release does not have an authorized minimal closure. Viper v1.15.0
requires properties v1.8.7, so the command also downgrades direct Viper to
v1.14.0. In total it changes 20 selected versions and removes one selected
module, including cloud, Consul, etcd, gRPC, and related transitive modules.
Modules become 234 -> 233. Graph edges become 3,557 -> 3,558 through 68 removed
and 69 added edges. The native complete-test population remains byte-identical
at 429 packages. The projected go.mod diff has exactly the properties and
Viper lines; go.sum remains byte-identical at 1,018 lines because the v1.8.6
and Viper v1.14.0 pairs already exist. Tidy remains unapplied and changes from
a 310-line projection to 312 lines, with a 163-line old/candidate projection
delta.

One properties package loads through `plybuild/cmd -> viper ->
viper/internal/encoding/javaproperties -> properties`. The codec calls
`NewProperties`, `Properties.WriteComment`, and `Load`. Writable Viper v1.15.0
and v1.14.0 codec tests pass at count 10, as do the three Ply consumer packages
in both states. Candidate v1.8.6 complete tests and vet pass, but complete race
fails identically on two runs: `assert/TestPanicPanicsAndDoesNotPanic` expects
`assert.go:65: did not panic`, while race instrumentation reports line 66.
Selected v1.8.7 reproduces the same brittle self-test defect. This is not a
reported data race, but it still fails the required complete candidate race
contract.

Pinned govulncheck v1.7.0 used a database updated
2026-09-02T19:12:04Z. Its primary modules index contains no properties entry,
and exact old/candidate vulnerability ID sets remain equal at 20 Darwin symbol,
30 Darwin module, and 20 Windows symbol findings. This is parity, not an
improvement.

The decision is
`reject-all-selection-changes-retain-v1.8.7-as-unresolved-floor-debt-no-dependency-commit`.
Canonical latest and every higher release violate the retained floor; highest
compatible v1.8.6 has an unauthorized 21-selection closure and fails complete
race twice. Dependency metadata therefore remained byte-identical. Per the
stop rule, repository build/full tests/race/vet/lint, help/API/CLI, and
post-implementation quality/acceptance/audit gates were not run after this
decision.

## Evidence And Tool Identity

- Properties decision root
  `/private/tmp/ply-p7-properties-selection.469049f.iT6qcB` has a fully
  verified 17,240-entry manifest SHA-256
  `941e805f6edf4e7c4d07a8912dcf8e5203dd2cea48530aa3fe6995cb19da097f`;
  decision-summary SHA-256 is
  `23b8556153b7ca9d606d5156d50d7d7791b4978ec2368ac3171ef6370b4f2c2d`.
- Recovery root `/private/tmp/ply-p7-go1.26.7-recovery.qvk4zs` retains its
  fully verified two-entry manifest SHA-256
  `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.
- Accepted go-colorful review, exact-quality, and regression roots still fully
  verify at 13/236,146/19,011 entries and SHA-256 values
  `6444e55370fb37a6c2c06c2c1c280b8fcbd34fc1f94f110a70bc7d067b05673d`,
  `924500cb4512fa9451cf218b119ce8969e32e7f02a67558807d067e3f41017a6`,
  and `2cf393025cf7768ef90fbfa8c0faf378284c3521cb118220dd74ba6fc4735d0e`.
  Its selection manifest retains its recorded digest and 47,633 entries, but
  one mutable Go telemetry counter beneath the retained HOME no longer
  verifies; the other 47,632 entries do. The stable selection summary remains
  `07467cdfd752a8ee9c64b1b6632370c8e0d75b6e3715090931ec8831b2080658`.
- Preserve the btree correction: its regression manifest retains its recorded
  digest and 24,197 entries, but 163 mutable cache/HOME entries no longer
  verify. Its other three manifests verify, and fresh later evidence
  supersedes that regression root.
- No-change mousetrap and rejected HCL roots remain fully verified at
  29,283/`71f337b114d1383cd103ef53856a5e767d519eda0b5635e9598b26654a78a152`
  and 21,427/`ad2954c3ebcf8bfe40e5e1684fcfcf662610555e6ed113515750db1a91c1acff`.
- Tool identities remain Go 1.26.7
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
  golangci-lint 2.12.2
  `3ba856c13833c4cda2182eb71bdbc7f96ddad339a1cdda8c36f88bfb1e34bd6f`,
  GoReleaser 2.17.1
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  apidiff `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`,
  and govulncheck v1.7.0
  `0db06024e71e82bd3e4444d57b757abf86fd4b7d0b08c1c21ca811597c953cd5`.

## Next Objective

Independently evaluate selected indirect
`github.com/mattn/go-runewidth v0.0.14` as exactly one bounded P7 module group.
Resolve canonical latest and every potentially compatible release from fresh
Go proxy, checksum-database, upstream, and primary Go vulnerability evidence.
Treat latest, Go floor, source identity, closure, loaded status, and behavior
as unknown.

Measure exact old/candidate selections, graph edges, complete package
population, checksums, exact-get diff, tidy projection, loaded packages/path,
focused behavior when loaded, candidate module tests, repository quality,
help/API/CLI identity, and vulnerability populations. Implement only an exact
floor-compatible selection with an explained minimal closure and every gate
passing; otherwise record rejection without editing dependency metadata. Do
not manufacture a dependency commit if v0.0.14 is already canonical latest.

Warm caches from separate Git archives. Never use tidy as implementation or
run `go mod download all` inside a measured tree. Stop before another module,
behavior/source changes, language/toolchain or quality-policy changes,
packaging, publication, or P8. Keep all evidence outside the worktree.
