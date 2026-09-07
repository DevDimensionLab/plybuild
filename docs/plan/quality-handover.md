# Quality Upgrade Handover

Generated: 2026-09-07T23:45:55+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- Latest dependency implementation remains Circbuf commit
  `3be2183ee310ccdc358ce4ed372c0785de25b88b`, exact parent
  `7a0ca4e2caba6d2fff20a9c169b181f9187b40dc`, tree
  `a51846840b5f252a30359530dfc950f811398431`, changing only `go.mod` and
  `go.sum`. The Consul API and Go Metrics evaluations made no dependency edit
  or implementation commit.
- Operator-authorized lifecycle repair
  `f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between
  Repr implementation `6f4d02eb9c86ec2df8a488a85ed973aababe1f38` and
  direct-child Repr handoff `342c7ece82eae3cd726aa658462089cfb876fbe0`.
  All remain ancestors.
- The answered Go Metrics archive and sole NEXT Go Radix archive link
  reciprocally. Relative to accepted go-cmp commit `c314bcb`, accepted
  dependency metadata still adds exactly 29 checksum lines.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, successor launch, or
  worktree removal occurred.

## Lifecycle And Retained Roadmap

Every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context must remain beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Preserve launcher lifecycle and scratch-only
mutation, compatibility, snapshot, Docker, acceptance, and audit rules. Never
run `go mod download all` in a measured tree.

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through retained Go Metrics v0.4.0. All earlier accepted,
rejected, and no-change decisions and evidence corrections remain final. P8
remains queued.

## Retained Armon Go Metrics Group

Retain exact-path `github.com/armon/go-metrics v0.4.0` without changing
dependency metadata. The proxy has 22 exact stable versions and no
prereleases. Exact Go `@latest` and `@v0` resolve v0.6.1 at
2026-07-29T13:06:13Z; `@master` resolves unreleased pseudo-version
`v0.6.2-0.20260907064447-465585286d74` at 2026-09-07T06:44:47Z. Both declare
Go 1.25.0 and module path `github.com/hashicorp/go-metrics`, so neither is an
exact-path or retained-floor candidate.

The old GitHub URL and Git transport redirect to authoritative public,
enabled, unarchived, non-fork `hashicorp/go-metrics`. It reports default branch
`master`, 20 branches, 21 current tags, 15 non-draft/non-prerelease GitHub
Releases, and 263 default-branch commits. Exact old/new Git refs are identical.
The README records the v0.5.0 module rename; the HashiCorp module/import path is
a distinct migration, not an in-place exact-path update.

Proxy v0.4.2 is the first renamed declaration. It aliases v0.5.0 commit
`aee7470331bc2a027cb2711f759f6f527a557ed8` and time
2023-06-08T13:48:53Z, but its alias tag is no longer present and it has no
GitHub Release. Exact get rejects it because it declares
`github.com/hashicorp/go-metrics` while required as the Armon path. Versions v0.5.0
through v0.5.4 retain Go 1.12 but have the wrong module identity; v0.6.0,
v0.6.1, and master additionally exceed the retained floor. The only retraction
is v0.3.11, for an undocumented breaking change to the metrics-sink interface.

The selected stable Release is commit
`129ee86de65934631a7fdbeb8c5aa0ec08bfdb6c`, tree
`dc8cd4f53669534104e66bf06f5ebaaf26d6feb4`, at
2022-05-25T15:01:32Z. Highest exact-path candidate v0.4.1 is stable Release
commit `b6d5c860c07ef6eeec89f4a662c7b452dd4d0c93`, tree
`d582c4e222a01bfe89e45a22b64273970b2008a8`, at
2022-09-08T12:00:52Z, directly descending from v0.4.0. All current tags are
lightweight and therefore have no tag-object signatures. The selected and
candidate merge commits independently verify against GitHub web-flow
fingerprint `5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`; its key is currently expired.
Master commit `465585286d749b5d38c457587f9b79022cc16883` verifies against current
fingerprint `968479A1AFF927E37D1A566BB5690EEEBB952194`.

The selected and candidate checksum pairs are respectively
`h1:yCQqn7dwca4ITXb+CbubHmedzaQYHhNhrEXLYUeEe8Q=` /
`h1:E6amYzXo6aW1tqzoZGT755KkbgrJsSdpwZ+3JqfkOG4=` and
`h1:hR91U9KYmb6bLBYLQjyM+3j+rcd/UhE+G78SFnF8gJA=` /
`h1:E6amYzXo6aW1tqzoZGT755KkbgrJsSdpwZ+3JqfkOG4=`; sum.golang.org agrees.
Each proxy ZIP has 31 regular files and exactly matches its tag source.
Normalized source-manifest SHA-256 values are
`6ddf7f35171d330f6fe6821c6f0eeee13dab3e974e0de91612268158480fba5c`
and `7d5df92d26827a5af7db4ebc0eff46129c3c4161d9c497d2d88d50977e94f84e`.
The candidate changes only `prometheus/prometheus.go` and its test, with 86
insertions and four deletions.

Both versions declare Go 1.12 and the same 11 requirements. Independent proxy
and exact-tag candidate forms select the same 52-module standalone complete
minimal closure and list four packages. Its maximum explicit Go declaration
is 1.12, so the closure preserves Go 1.18. Keep that standalone test closure
separate from project MVS.

## Candidate Tests And Consumers

Candidate qualification fails. Under exact Go 1.26.7, proxy and exact-tag
forms verify and list, but they do not all pass count-1 or race tests and both
fail count-10. The flaky Statsite/DogStatsD UDP harnesses panic on a nil
connection or compare interleaved packets. Vet fails identically in substance
with 46 diagnostics: 29 `testing.T.Fatalf` calls from non-test goroutines, 12
unkeyed Label literals, one copied lock, and four empty appends. Exact Go
1.18.10 builds and passes one race run but reproduces count-10 failure and 44
vet diagnostics. These mandatory module stop-rule failures reject v0.4.1.

Project MVS receives v0.4.0 from Viper v1.15.0. Older Viper, Crypt, Serf, and
Memberlist graph edges request v0.3.10 or the 2018 pseudo-version. Ply has no
Go Metrics source import, the selected module loads in zero packages, and
`go mod why -m` says the main module does not need it. Viper is loaded but its
source does not import Metrics.

Selected historical consumers Serf v0.10.1 and Memberlist v0.3.0 import the
module. Serf uses Label, labeled counter/timer/sample calls, in-memory,
Statsite, StatsD, fanout, default-config, and global setup symbols. Memberlist
uses `SetGauge`, `IncrCounter`, and `MeasureSince`. Under candidate v0.4.1 and
exact Go 1.18, both consumer source trees compile and vet; Serf coordinate and
Memberlist awareness focused tests pass count-10 and race. Their full legacy
suites separately encounter Darwin loopback/test-harness limitations: the
Memberlist suite cannot bind hard-coded 127.0.0.x addresses, and Serf times out
searching that space. These are consumer evidence, not loaded Ply behavior or
candidate-module test waivers.

## Projection And Vulnerabilities

Exact candidate get changes only the selected Go Metrics version, but lazy
module loading also materializes three indirect requirements, 14 graph edges,
and four full checksums. It projects 234 -> 234 selected modules, 3,581 ->
3,595 graph edges, 429 -> 429 complete packages, 41 -> 41 loaded modules, zero
-> zero loaded Go Metrics packages, 1,045 -> 1,049 `go.sum` lines, and 361 ->
386 tidy-diff lines. The edges are the main module's three materialized
requirements plus the candidate's 11 declared requirements; there are no
removed edges. The checksum additions are the candidate pair plus full sums
for already-selected immutable-radix v1.3.1 and golang-lru v0.5.4. Tidy removes
all projected requirements/checksums and restores the inherited v0.4.0
selection. No metadata or dependency implementation commit was created.

Govulncheck v1.7.0 used primary data updated 2026-09-02T19:12:04Z. Its 1,392-
record module index has no old- or new-path Go Metrics record, and no trace
contains either module. Original/candidate projections are identical after
normalization at 20 IDs/22 traces for Darwin reachable symbols, 30 Darwin
module IDs, and 20 IDs/22 traces for Windows reachable symbols. Independently
sorted full reachable and module-finding SHA-256 values are
`7219be1e8592d7177c964e8f3f37cbc357644daa40baf00c520394a073d0a681`
and `390f7bf685e19fc5d10602da7e35394c97f8fd460fe40f92c4e6726c00627d29`.
The exact IDs and 22 traces remain the previously accepted population.

Because the candidate fails module stop rules, downstream repository,
snapshot/Docker, acceptance, and Q0-Q2 gates were inapplicable and were not
claimed. The accepted Circbuf quality baseline and exact 29-line checksum
delta from accepted go-cmp remain unchanged.

Go Metrics decision evidence contains 336 verified entries. Evidence-manifest
SHA-256 is
`5afb2b5fb26824c5e4c1db7497b8a6e6dbd24bc20194df7ca7fc5f4385e0dde1`;
decision-summary SHA-256 is
`c5b47a2a7d6f881e7a5aad5895d197027556e238c683bca685f45fff3b3c0592`.

## Tools And Corrections

- Exact Go 1.26.7 binary SHA-256 is
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
  official archive SHA-256 is
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  Exact Go 1.18.10 archive/binary SHA-256 values are
  `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`
  and `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
  Go 1.26.7 remained first in PATH with GOENV=off, GOWORK=off,
  GOTOOLCHAIN=local, and no ambient GOFLAGS.
- Preserve portable golangci-lint 2.12.2 archive receipt
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
  GoReleaser 2.17.1 receipt
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  and apidiff receipt
  `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
  Rebuilt govulncheck v1.7.0 is a nonportable receipt at SHA-256
  `6ede6b57bd7fe1deb3f707e3c0f3008f5e93e854516687e6769afc4e27d9dacb`;
  version, builder, database identity, and functional scans are the proof.
- Supersede the first Go-1.18 CGO setup, which attempted an Xcode resolver cache
  outside scratch, with final direct SDK/clang and scratch TMPDIR/GOTMPDIR
  runs. Supersede invalid `-mode=module` govulncheck calls with correct no-
  pattern `-scan=module` scans. Concatenated JSON was decoded before trace-
  depth filtering. No correction mutated repository files or hid a failure.
- Preserve all earlier compatibility-cache, source/signature, checksum-delta,
  vulnerability-order, Docker-path, scratch-local BSD-mktemp, preflight,
  locale, and exact-Go-1.18 correction records.

## Retained Decisions

- Retain exact selections `github.com/alecthomas/repr v0.5.4`,
  `github.com/alecthomas/assert v1.0.0`,
  `github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b`,
  `github.com/alecthomas/chroma v0.10.0`,
  `github.com/alecthomas/colour v0.1.0`,
  `github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751`,
  `github.com/antihax/optional v1.0.0`, accepted Circbuf candidate, retained
  Consul API pseudo-version, retained Go Metrics v0.4.0,
  `gopkg.in/alecthomas/kingpin.v2 v2.2.6`, `gopkg.in/resty.v1 v1.12.0`,
  `gopkg.in/errgo.v2 v2.1.0`,
  `gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15`,
  `gopkg.in/yaml.v2 v2.4.0`, `gopkg.in/yaml.v3 v3.0.1`, and
  `go.yaml.in/yaml/v3 v3.0.5`. Retain selected Kong pseudo-version and every
  other prior bounded decision.
- Preserve all recorded corrections from prior dependency groups. The
  deliberately purged former two-entry recovery manifest remains historical
  SHA-256
  `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.

## Next Objective

Independently evaluate selected exact-path `github.com/armon/go-radix v1.0.0`
at 2018-08-24T02:57:28Z as the next single P7 group. Project MVS selects it
without an explicit `go.mod` requirement. The exact stable proxy list contains
only v1.0.0; exact `@latest` and `@v1` resolve that selected stable version,
while exact `@master` resolves unreleased pseudo-version
`v1.0.1-0.20221118154546-54df44f2176c` at 2022-11-18T15:45:46Z. The selected,
latest, and master metadata expose no Go declaration in the minimal survey.

The public, enabled, unarchived, non-fork repository reports default branch
`master`, one branch, one current tag, and zero GitHub Releases. Master is
commit `54df44f2176c4a553657a4f0dbe6fdb108288be3`; tag v1.0.0 points to
`1a2de0c21c94309923825da3df33a4381872c795`. Treat canonical source identity,
tag and commit signatures, release qualification, history, retractions,
declarations, complete closure, tests, consumers, loaded paths, exact-get
projection, and vulnerability effect as unknown until independently proved.

Current accepted measurements remain 234 modules, 3,581 graph edges, 429
complete-test packages, 41 loaded modules, 1,045 `go.sum` lines, a 361-line
unapplied tidy projection, and exact 20/30/20 Darwin-symbol/Darwin-module/
Windows-symbol vulnerability populations. Change only Go Radix's exact
selected version and explained minimal closure if a qualified changed
selection preserves Go 1.18 and passes every applicable contract. If v1.0.0
is the highest qualified selection, do not manufacture redundant metadata.
Do not combine another dependency group or P8.
