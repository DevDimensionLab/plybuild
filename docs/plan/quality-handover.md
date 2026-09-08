# Quality Upgrade Handover

Generated: 2026-09-09T00:23:51+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation remains exact Fatih Color v1.15.0
  commit `6ca672ef38688b7f6f505cf0cb273d07c4c2ba9a`, exact parent
  `d181fcd6c11fa147e0b44dd007598872875fc9e6`, and tree
  `9ba2fdc442553622028a4a8464536915d345510c`. It changes only `go.mod` and
  `go.sum`, with three insertions and one deletion.
- Ghodss YAML v1.0.0 and Fsnotify v1.6.0 were retained without dependency
  implementation commits. Exact XXHash v2.3.0 implementation
  `e5d6252825d7a1822c01819b9144050f345a6ad4` and Speakeasy v0.2.0
  implementation `41f9561f6ea2f5b6395c5c4d9bcc56a54533133a` remain ancestors.
  Ansimage, Imaging, Fnmatch, Readline, Logex, OpenCensus Proto, and Crypt
  remain retained without dependency edits.
- The answered Ghodss YAML archive and sole NEXT root-GLFW archive link
  reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` was created. No push, merge, publication,
  release, stash, revert, successor launch, or worktree removal occurred.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, accepted
Speakeasy v0.2.0, XXHash v2.3.0, and Fatih Color v1.15.0 moves, and retained
Crypt, OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
and Ghodss YAML selections. P8 remains queued. Earlier outcomes and lifecycle
ancestry are final; do not reopen them or combine another module group.

Every disposable cache, projection, archive, source, report, generated
artifact, evidence tree, and build context must remain beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never create direct `/private/tmp/ply-*`
roots, never run `go mod download all` in a measured tree, and preserve the
launcher's scratch cleanup and reciprocal archive contract.

## Ghodss YAML Decision

Retain exact-path `github.com/ghodss/yaml v1.0.0` without editing dependency
metadata. The proxy exposes only v1.0.0, which is both selected and stable
latest at `2017-03-27T23:54:44Z`. Its source/mod checksum pair is
`h1:wQHKEahhL6wmXdzwWG11gIVCkOv05bNOh+Rxn0yngAk=` /
`h1:4dBDuWmgqj2HViK6kFavaiC9ZROes6MMH2rRYeMEF04=`; the proxy ZIP SHA-256 is
`c3f295d23c02c0b35e4d3b29053586e737cf9642df9615da99c0bda9bbacc624`
and matches the exact Git tag source.

V1.0.0 is a lightweight unsigned tag and unsigned commit
`0ca9ea5df5451ffdf184b4428c902747c2c11cd7`, tree
`252e285a136d503d2913ad39cb3b1669b9a76999`, with parents
`04f313413ffd65ce25f2541bfd2b2ceec5c0908c` and
`a4f8cbd2fd05654d25f651b7e26d612ce3c98cc7`. Go-import resolves the public,
enabled, unarchived, non-fork `ghodss/yaml` repository on default `master`.
There is one GitHub release and no prerelease, retraction, module deprecation,
redirect, or exact-path alternate release.

Current master `d8423dcdf3440d0a5baffc6f90a11e4128545620`, tree
`8df4d8facd33473e866d0f2e41e121c8732d5a6a`, parent
`1e4101787d1907800b0200eb38faa0f5041d8ee6`, is 16 commits after v1.0.0 but
resolves only as unreleased pseudo-version
`v1.0.1-0.20220118164431-d8423dcdf344`. The maintained fork/module
`sigs.k8s.io/yaml` is a distinct identity and was not substituted.

The published module file contains no Go directive or requirements. Isolated
source resolution adds `gopkg.in/yaml.v2 v2.4.0` at Go 1.15; its historical
test closure adds
`gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405`. The complete closure's
native tests pass exact Go 1.26.7 and contained Go 1.18.10, so it preserves the
floor.

Release count-1, two independent count-10 passes, race, and production
cross-builds for Darwin, Linux including 386, Windows, and FreeBSD pass both
SDKs. The code is pure Go with no Cgo, assembly, generation, examples,
testdata, fuzz, or property suite. Historical source/test debt is explicit:
vet reports one unreachable return, and the release test file fails Linux/386
compile because an untyped MaxInt64 overflows int. Production Linux/386
builds, and the first post-release commit fixes the test. Unreleased master is
not a better candidate: Go 1.26 native tests fail through a malformed example
vet error, while the unreachable code remains.

V1.0.0 exports Marshal, Unmarshal, YAMLToJSON, and JSONToYAML, but no strict
API. Master adds options and strict conversion; duplicate keys then fail,
while unknown fields require DisallowUnknownFields. Independent fixtures under
both SDKs cover keys/scalars/numbers, tags, anchors/aliases, duplicates,
unknown fields, embedded and custom marshalers, JSON tags, interfaces,
byte/string handling, invalid inputs, multi-document input, ordering,
escaping, error propagation, determinism, and concurrency.

The YAML-to-JSON pipeline is lossy: duplicate keys are last-wins, later YAML
documents are ignored, numeric/interface types change, tags and aliases lose
identity, and invalid binary bytes become replacement runes. Integer and
boolean keys stringify, while composite/null keys fail. Distinct YAML keys
`1` and `"1"` collapse to one JSON key and the surviving value depends on Go
map iteration. Ordinary non-colliding conversion is deterministic. These are
documented release boundaries, not a serialization-migration scope.

## Projection, Quality, And Vulnerability Measurements

Ply selects Ghodss YAML only through:

`plybuild -> mvn-pom-mutator@v0.2.3 -> viper@v1.10.1 ->
etcd/api/v3@v3.5.1 -> grpc-gateway@v1.16.0 -> ghodss/yaml@v1.0.0`.

`go mod why -m` says the main module does not need it, and no Ghodss package
is loaded in production or complete tests. Ply's direct yaml.v2/yaml.v3
imports are separate. Retention preserves 234 modules, 3,583 graph edges, 429
complete-test entries, 41 loaded modules, 197 loaded module-backed packages,
1,051 checksum lines, and a 383-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, metadata still adds exactly 35 checksum
lines and removes zero.

An exact selected `go get` is not truly inert: it would manufacture a
redundant explicit main edge and add the unused source checksum, producing
3,584 edges, 1,052 sum lines, and a 388-line tidy projection without loading
a package. It was not applied. Master similarly loads no package and would
add an unreleased explicit edge plus two checksums.

Project mod verify, build, count-1/count-10/race/vet, Windows-amd64 build,
pinned golangci-lint 2.12.2, byte-identical core help, API/CLI compatibility,
empty-HOME count-2, and complete preflight pass. Preflight includes all 62
launcher controls, 80/80 killed mutants, and 15/15 audit self-test controls.
Changed-selection-only snapshot, Docker, and exact quality work is
inapplicable. The accepted 27/27 Q0-Q2 L2 scorecard remains
`dae9e51e26353f72d9026e1d6eecbef697bcfa3b06cdc1905164c6d20057dafb`.

A supplemental contained-Go-1.18 project build passes. Full project tests
retain two pre-existing Darwin `pkg/shell` closed-file error-text differences;
Ghodss remains unloaded and retained/master project projections are identical,
so those inherited diagnostics are not Ghodss effects.

Fresh govulncheck v1.7.0 used 1,392 primary module records. Retained/master
canonical populations are identical: 20 IDs/22 traces for Darwin and Windows
symbol scans, 22 Darwin package findings, and 30 Darwin module findings.
Ghodss YAML has no record, finding, or trace. Existing x/image/ansimage
findings remain unrelated.

The 92-entry selected evidence verifies against manifest SHA-256
`cfb5bfcc497db42efebda5323b9e187aa63520192a7354631fd93391970bd771`.
Decision-summary SHA-256 is
`724f7f2ac4394296cd36be360df90539deed1ecaf0b92386c6027b0f2c45f7b9`.

One stopped preflight attempt wrote its two compatibility JSON reports to the
default ignored `target/compatibility` path. The cleanliness gate detected the
artifact; both files and directories were removed immediately. The valid
complete rerun routed all reports beneath session scratch and passed. The
stopped run is excluded from passing evidence.

## Tools And Retained Decisions

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  Contained Go 1.18.10 binary/archive SHA-256 values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
- Pinned golangci-lint 2.12.2 archive SHA-256 is
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`.
  GoReleaser 2.17.1 and portable apidiff receipts remain
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`
  and `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
- Retain accepted Fatih Color v1.15.0, XXHash v2.3.0, and Speakeasy v0.2.0
  moves and every earlier exact decision. Answered archives are authoritative
  detail.

## Next Objective

Independently evaluate the historical root module exact path
`github.com/go-gl/glfw v0.0.0-20190409004039-e6da0acd62b1` as the next single
P7 group. Keep it strictly separate from the nested
`github.com/go-gl/glfw/v3.3/glfw` module.

The root proxy list is empty. Selected commit
`e6da0acd62b1b57ee2799d4d0a76a7d4514dc5bc` has tree
`17ab3d23b59cab5cffcafe1232da9c2af635a2a1`, parent
`39f94f8075907c0c6524d2791e05d625627fe268`, and checksum pair
`h1:QbL/5oDUmRBzO9/Z7Seo6zf912W/a6Sr4Eu0G/3Jho0=` /
`h1:vR7hzQXu2zJy9AVAgeJqvqgH9Q5CA+iKCZ2gyEVpxRU=`. Its synthesized module
file has no Go directive or requirements.

Exact-path latest is floor-ineligible Go-1.19 pseudo-version
`v0.0.0-20260823155953-d41da22a9587`, full commit
`d41da22a9587f777098f96d37014f6cdd35d1afb`. Project MVS reaches selected
only through x/exp, but `go mod why -m` says it is unneeded. Resolve exact root
release/pseudo ancestry, complete closure, Cgo and platform prerequisites,
package loading, projections, tests, and vulnerability identity without
combining the nested v3.3 module or any graphics-stack migration.
