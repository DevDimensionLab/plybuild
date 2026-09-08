# Quality Upgrade Handover

Generated: 2026-09-08T23:00:17+02:00

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
- Fsnotify v1.6.0 was retained without a dependency implementation commit.
  Exact XXHash v2.3.0 implementation `e5d6252825d7a1822c01819b9144050f345a6ad4`
  and Speakeasy v0.2.0 implementation
  `41f9561f6ea2f5b6395c5c4d9bcc56a54533133a` remain ancestors. Ansimage,
  Imaging, Fnmatch, Readline, Logex, OpenCensus Proto, and Crypt remain
  retained without dependency edits.
- The answered Fsnotify archive and sole NEXT Ghodss YAML archive link
  reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` was created. No push, merge, publication,
  release, stash, revert, successor launch, or worktree removal occurred.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, accepted
Speakeasy v0.2.0, XXHash v2.3.0, and Fatih Color v1.15.0 moves, and retained
Crypt, OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, and
Fsnotify selections. P8 remains queued. All earlier outcomes and lifecycle
ancestry are final; do not reopen them or combine another module group.

Every disposable cache, projection, archive, source, report, generated
artifact, evidence tree, and build context must remain beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never create direct `/private/tmp/ply-*`
roots, never run `go mod download all` in a measured tree, and preserve the
launcher's scratch cleanup and reciprocal archive contract.

## Fsnotify Decision

Retain exact-path `github.com/fsnotify/fsnotify v1.6.0`. Stable v1.9.0 is the
highest release whose complete minimal closure preserves Go 1.18, but two
independent complete Go 1.18 repeat suites fail intermittently on Darwin
kqueue. V1.7.0 and v1.8.0 fail other native-suite contracts, and v1.10.x
requires Go 1.23. No exact newer stable therefore passes every contract.

The proxy exposes 40 stable versions. Retracted v1.5.0 and v1.5.3 leave 38
eligible exact versions; old leading-zero repository tags are proxy-absent.
There is no prerelease or module deprecation. Go-import, proxy, sumdb, and Git
all resolve the enabled, unarchived, non-fork upstream on default `main`.
Every relevant tag is a main ancestor and proxy/Git sources are byte-identical.

Relevant release commits/trees are:

- v1.6.0 `5f8c606accbcc6913853fe7e083ee461d181d88d` /
  `566d47ec45f239dd7674a3f0ad54d40fe76c482e`;
- v1.7.0 `cfc9c4f277ea6ec18de92444b31983b183deb4fb` /
  `3af1f5b0f8a3c0e869fd2b2a6e662e5f68415dcd`;
- v1.8.0 `a9bc2e01792f868516acf80817f7d7d7b3315409` /
  `316db70673ed511404435ea2ca0843183b470561`;
- v1.9.0 `ae0e7923765f64fb8061396db7edebb558cf6093` /
  `e983f596e89398e1d512b70eca573328662dac46`;
- v1.10.0 `8d01d7b9cbe0199e4a1e60fbd965fb05dbb42123` /
  `82a0947703590aaaf9738eb6f6c57abca74e6f5c`; and
- v1.10.1 `76b01a6e8f502187fecedea8b025e79e5a86085c` /
  `372b850d90678bf373c51640f701212ba17332ae`.

V1.6.0, v1.7.0, and v1.9.0 have annotated unsigned tags, and v1.8.0 is
lightweight; their commits have GitHub-valid SSH signatures. V1.10.0/v1.10.1
tag objects and commits have valid PGP signatures. Unreleased main resolves
as `v1.10.2-0.20260511064106-20b1e15ef3c7` and still declares Go 1.23. No
redirect, fork, alternate path, or pseudo-version outranks stable latest.

Selected v1.6.0 declares Go 1.16 and x/sys
`v0.0.0-20220908164124-27713097b956`; v1.7.0 declares Go 1.17 and x/sys
v0.4.0; v1.8.0/v1.9.0 declare Go 1.17 and x/sys v0.13.0. All corresponding
x/sys modules declare no higher than Go 1.17. V1.10.x declares Go 1.23, and
contained Go 1.18 cannot load its test dependency on standard `slices`.

The v1.6.0-to-v1.9.0 API adds only compatible `NewBufferedWatcher`,
`AddWith`, `WithBufferSize`, and `ErrClosed` symbols. `Op.Has` changes combined
masks from all-bit to any-bit matching at v1.7; Viper uses only single masks.
Independent fixtures cover event strings/masks, paths and symlinks,
non-recursion, event types, blocked consumers, close/channel ordering,
WatchList, errors, and concurrent Add/Remove. Corrected old/v1.9 fixtures pass
count-1, two independent count-10 runs, race, and vet under both SDKs.

Linux/inotify has kernel watch, instance, and overflow limits. Darwin/BSD
kqueue consumes descriptors per watched object and scans directories. Windows
uses ReadDirectoryChangesW with a 64-KiB default buffer, lacks chmod, and
retains renamed watches. Illumos/Solaris use FEN; Solaris is not executed in
upstream CI. Unsupported targets compile an error backend. There is no Cgo,
assembly, or generated Go; v1.9's C file is an external kqueue diagnostic.

V1.7's Go 1.18 count-10 suite raises descriptor limits, creates roughly
58,000 watches, hits `EMFILE`, and cascades. V1.8 deterministically fails a
2024-only `TestDiffMatch` in 2026 and separately observed a multiple-write
timing miss. V1.9 passes count-1, vet, race, and Go 1.26 repeats, but two
complete Go 1.18 count-10 runs fail `TestRace/add_and_remove_watches` with an
unexpected kqueue `bad file descriptor`. Later isolated and full repeats pass,
proving the fault intermittent. V1.9 already contains mitigation commit
`0023e08`; its message calls the failure longstanding and explicitly says the
change is incomplete. V1.10 has further kqueue descriptor-lifecycle fixes but
is independently Go-1.23-only. Retries were not used to hide the failure.
Selected v1.6's resource-heavy stress-suite debt remains characterized rather
than misreported as clean.

Ply's production path is `plybuild/cmd -> spf13/viper -> fsnotify`. Viper
v1.15.0 full/focused-repeat/race/vet checks pass old and candidate under both
SDKs. It watches the config parent directory and handles writes, creates,
symlink replacements, removes, errors, rereads, and callbacks. Ply itself does
not call WatchConfig or OnConfigChange.

## Projection, Quality, And Vulnerability Measurements

Retained metrics remain 234 modules, 3,583 graph edges, 429 native
complete-test entries, 41 loaded modules, 197 loaded module-backed packages,
1,051 checksum lines, and a 383-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, metadata adds exactly 35 checksum lines and
removes zero.

Exact old get is inert. A v1.9 projection changes only the main Fsnotify edge,
adds its checksum pair, and keeps project x/sys v0.30.0. It has 430 test
entries, 198 loaded packages, 1,053 sum lines, and 387 tidy-diff lines; the
extra package is `fsnotify/internal`. V1.10.1 also moves main Go 1.18 to 1.23,
has 3,585 graph edges, and a 391-line tidy projection. No projection or tidy
result was applied.

The retained projection passes mod verify, build, count-1/count-10/race/vet,
Windows-amd64 build, pinned golangci-lint 2.12.2, byte-identical root/status/
upgrade/build help, API/CLI compatibility, and empty-HOME count-2. Since the
selection is unchanged, changed-selection-only snapshot, Docker, acceptance,
audit, and exact quality work is inapplicable. The current tree retains the
accepted 27/27 Q0-Q2 L2 scorecard
`dae9e51e26353f72d9026e1d6eecbef697bcfa3b06cdc1905164c6d20057dafb`.

Fresh govulncheck v1.7.0 used 1,392 primary module records. Old/candidate
canonical populations are identical: 20 IDs/22 traces for Darwin and Windows
symbol scans, 22 Darwin package findings, and 30 Darwin module findings.
Fsnotify has no record, finding, or trace. Existing x/image/ansimage findings
remain unrelated.

The 899-entry selected session evidence verifies against manifest SHA-256
`e03b23b7a574cb7d4ec4213fc7a28d785fc5a70517a864a338fa64f227db11db`.
Decision-summary SHA-256 is
`6ce37518efc54801eea666954aa6f850142909d99c347cad21d75aa408133a34`.

One misdirected disposable `go get` briefly touched primary metadata before
being detected. Exact original blobs were restored immediately, and ordinary
and ignored status were proven empty before further measurement. Superseded
malformed diagnostics and first-version fixture assumptions do not count as
passing evidence.

## Tools And Retained Decisions

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  Contained Go 1.18.10 binary SHA-256 is
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
- Pinned golangci-lint 2.12.2 archive SHA-256 is
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`.
  GoReleaser 2.17.1 and apidiff receipts remain
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`
  and `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
- Retain accepted Fatih Color v1.15.0, XXHash v2.3.0, and Speakeasy v0.2.0
  moves and every earlier exact decision. Answered archives are authoritative
  detail.

## Next Objective

Independently evaluate exact-path `github.com/ghodss/yaml v1.0.0` as the next
single P7 group. The proxy has only selected/latest stable v1.0.0, released
2017-03-27 at lightweight tag commit
`0ca9ea5df5451ffdf184b4428c902747c2c11cd7`, with checksum pair
`h1:wQHKEahhL6wmXdzwWG11gIVCkOv05bNOh+Rxn0yngAk=` /
`h1:4dBDuWmgqj2HViK6kFavaiC9ZROes6MMH2rRYeMEF04=`. Its module file contains
no `go` directive or requirements, so prove imported source/test closure
rather than assuming a floor.

Project MVS selects it through the declared grpc-gateway v1.16.0 edge, but
`go mod why -m` says Ply does not need it and no package is loaded. Keep that
MVS-only selection separate from Ply's direct YAML libraries and from any
maintained successor identity. Resolve current master
`d8423dcdf3440d0a5baffc6f90a11e4128545620`, release ancestry, exact
conversion semantics, tests, closure, projection, and vulnerability facts.
