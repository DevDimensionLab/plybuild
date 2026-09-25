# Agent Session: Evaluate Armon Consul API Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-07T214638+0200-evaluate-armon-consul-api-dependency`
Created: `2026-09-07T21:46:38+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `45a46dd3188cfabbad82b9ff2bf7c6e295924f1ae11876d46ab479059b3f7760`
Previous: [2026-09-07T194916+0200-evaluate-armon-circbuf-dependency.md](2026-09-07T194916+0200-evaluate-armon-circbuf-dependency.md)
Next: [2026-09-07T223740+0200-evaluate-armon-go-metrics-dependency.md](2026-09-07T223740+0200-evaluate-armon-go-metrics-dependency.md)
Outcome: Retained the already-selected canonical-latest unreleased pseudo-version without dependency edits; exact get changed no selection and the module's repeated/race tests plus vet failed required gates.

# Answer

## Decision

Retain exact-path `github.com/armon/consul-api` at selected pseudo-version
`v0.0.0-20180202201655-eb2c6b5be1b6` without changing `go.mod` or `go.sum`.
It is the canonical latest, `master` head, and highest exact-path version whose
declarations preserve Go 1.18, but it is not a stable release and it does not
qualify for a dependency change: exact get changes no selected version, and
the module's complete repeated/race test and vet contracts fail. No dependency
implementation commit was created.

## Canonical Resolution, History, And Identity

- The exact proxy version list is empty: there are zero published stable or
  prerelease versions. Proxy and exact Go `@latest`, exact Go `@master`, and
  the selected query all resolve
  `v0.0.0-20180202201655-eb2c6b5be1b6` at
  2018-02-02T20:16:55Z. Exact `@v0` has no matching versions. The synthetic
  proxy `go.mod` has no retract directive. There are no tags, tag objects, tag
  signatures, or GitHub Releases, so this remains an unreleased pseudo-version.
- The authoritative `https://github.com/armon/consul-api` repository is public,
  enabled, operationally unarchived, and non-fork. It has default branch
  `master`, one branch, 46 commits, zero tags, and zero Releases. The sole
  branch head is selected merge commit
  `eb2c6b5be1b66bab83016e0b05f01b8d5496ffbd`, with parents
  `dcfedd50ed5334f96adee43fc88518a4f095e15c` and
  `f746cfef698c8f40b4b4ef0396be5984dcd7acbb`, tree
  `aeb2299aaf107d0823ce91f057798119b821e81b`, and author/commit time
  2018-02-02T12:16:55-08:00.
- GitHub reports the commit signature verified/valid, verified at
  2024-11-05T00:08:31Z. Independent verification finds a cryptographically
  valid embedded GitHub web-flow signature for fingerprint
  `5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`; the key is now expired, so
  current `git verify-commit` exits nonzero after emitting `VALIDSIG`,
  `KEYEXPIRED`, and `EXPKEYSIG`. The official key file SHA-256 is
  `6e8af687f60cf3f403151c8fb1b26e95e6f9e424ca60cc8f3787bd4466a3ef84`.
- All 21 regular proxy ZIP files byte-match the exact GitHub commit archive and
  Git tree. Both normalized source manifests have SHA-256
  `358c1cfd7c63c682fc02e069ef0ab20f494bc44d9789e62ca39e62fc4026cbee`.
  Proxy ZIP, proxy `go.mod`, and GitHub archive SHA-256 values are
  `091b79667f16ae245785956c490fe05ee26970a89f8ecdbe858ae3510d725088`,
  `5efec71e3cc865139136973305b438abc97908a545997521b7ed2aa48b1b4dd6`,
  and `ba5e02e82cb70d6d7cb7970a8f81535df3266cfeae2b2d3b9e38de7ca2d76b4a`.
- The checksum pair is
  `h1:G1bPvciwNyF7IUmKXNt9Ak3m6u9DE1rF+RmtIkBpVdA=` /
  `h1:grANhF5doyWs3UAsr3K4I6qtAmlQcZDesFNEHPZAzj8=`. Independent
  sum.golang.org lookup agrees at tree 62473365.
- The repository is unarchived but source-deprecated: its README directs users
  to `github.com/hashicorp/consul/api`. That is a distinct module/import path
  in the Consul monorepo, not an alternate exact-path version. Its current
  latest v1.34.4 declares Go 1.26 and is outside this bounded decision; no
  `/v2`, `/v3`, renamed, or gopkg.in exact-path continuation exists.

## Floor, Closure, And Module Gates

The selected source predates modules; its proxy-generated `go.mod` declares
only `module github.com/armon/consul-api`, with no Go version and no
requirements. Go therefore applies the historical implicit Go 1.16 module
version, and the complete declared minimal closure is Consul API alone. It
does not raise the retained Go 1.18 floor. Its standalone tests add no module
requirements and remain separate from project MVS.

Exact Go 1.26.7 proxy and exact-commit forms independently pass `go mod verify`
and `go list`, remain unchanged, and fail identically:

- count-1: exit 1, 33 failed tests;
- count-10: exit 1, 330 failed tests;
- race: exit 1, 33 failed tests; and
- vet: exit 1 with `kv_test.go:187:4` and `kv_test.go:238:4` calling
  `(*testing.T).Fatalf` from non-test goroutines.

Every test failure is a refused connection to the required external Consul
agent at `http://127.0.0.1:8500`. Normalized proxy/commit count-1, count-10,
and race SHA-256 values are respectively
`60095eb84d9f0d2d2e38ce8bf2de49a971879febe6b08c6e7de0b32f3b0ba052`,
`15dc8ac4c780ebcfaab521cb89352c0c460b675dd5a199857df10e065da8862d`,
and the same count-1 hash. Vet output SHA-256 is
`59298258e3b1d6c8f9116b21ea14b5d2f8a1aa03e735e9b6e23036e6be6dfb67`.
Exact Go 1.18.10 compiles the package, then reproduces the 33/330/33 failures
and two vet diagnostics. These are candidate stop-rule failures, not hidden
runner errors.

## Project Projection And Consumers

The original graph has one request for the module:
`github.com/devdimensionlab/mvn-pom-mutator@v0.2.3` requests the selected
version. That requester carries the requirement indirectly but its source does
not import Consul API. Consul API loads in zero Ply packages, Ply has no source
import, and `go mod why -m` says the main module does not need it.

Exact selected-version get exits 0 but changes no selection. It projects:

- selected modules 234 -> 234;
- graph edges 3,581 -> 3,582, solely a redundant main-to-Consul edge;
- complete-test packages 429 -> 429;
- loaded modules 41 -> 41 and loaded Consul packages 0 -> 0;
- `go.sum` lines 1,045 -> 1,046, solely the already-selected full checksum;
  and
- unapplied tidy diff 361 -> 363 lines.

It also adds a redundant indirect `go.mod` requirement. Tidy removes that
requirement, main edge, and full checksum while retaining the inherited
go.mod checksum and all pre-existing debt. Accepted measurements and the exact
29 checksum lines added relative to accepted go-cmp commit c314bcb therefore
remain unchanged.

The selected closure's historical real consumer is
`github.com/xordataexchange/crypt@v0.0.3-0.20170626215501-b2862e3d0a77/backend/consul`.
It uses `DefaultConfig`, `Config.Address`, `NewClient`, `Client.KV`, `KV.Get`,
`KV.List`, `KV.Put`, `KVPair.Key`, `KVPair.Value`,
`QueryOptions.WaitIndex`, and `QueryMeta.LastIndex`. An external Go-1.18
fixture with a local HTTP test server exercises those exact configuration,
client, KV, query, metadata, key, and value flows and passes count-10, race,
and vet. This is consumer evidence, not a claim that Ply executes the module.

## Vulnerabilities And Applicable Quality Scope

Govulncheck v1.7.0, built with exact Go 1.26.7, queried primary data updated
2026-09-02T19:12:04Z. The 1,392-entry module index has no Consul API record;
none of the reachable traces contains Consul API. Original and exact-get
projection results are byte-identical:

- Darwin and Windows reachable-symbol IDs are
  `GO-2023-1989`, `GO-2023-1990`, `GO-2024-2937`, `GO-2024-3333`,
  `GO-2025-3595`, `GO-2026-4440`, `GO-2026-4441`, `GO-2026-4815`,
  `GO-2026-4961`, `GO-2026-5025`, `GO-2026-5027`, `GO-2026-5028`,
  `GO-2026-5029`, `GO-2026-5030`, `GO-2026-5031`, `GO-2026-5032`,
  `GO-2026-5061`, `GO-2026-5062`, `GO-2026-5066`, and `GO-2026-6222`.
  Each platform has 20 IDs and 22 traces; normalized SHA-256 is
  `5bebff017082945899b65922b3262a8b91affcda7abfe2fade54dee439749dfe`.
- Darwin module scanning adds `GO-2023-1988`, `GO-2023-2102`,
  `GO-2024-2687`, `GO-2025-3503`, `GO-2026-4918`, `GO-2026-4962`,
  `GO-2026-5024`, `GO-2026-5026`, `GO-2026-5942`, and `GO-2026-5970`,
  for 30 exact IDs/findings. Normalized SHA-256 is
  `f1cc7393d1a1d82da88379fceb830d926f8c1f0fe553e4af339e76ac6f7a1bc3`.
- The 22 precise reachable traces comprise the 20 IDs above: TIFF decode
  chains for 1989/1990/2937/4815/5032/5062/5066, HTML parse chains for
  3333/3595/4440/4441/5025/5027/5028/5029/5030, WebP decode chains for
  4961/5061, BMP decode for 5031, and VP8L-through-WebP decode for 6222. All
  decode/parse chains continue through `image.Decode` or `html.Parse`,
  ansimage/go-term-markdown/gomarkdown, Ply command initialization, Cobra, and
  `cmd.ExecuteE`. The two additional traces are 4815 through
  `tiff.ReadAt -> io.Read -> io.ReadAll -> getJsonWithAccessToken ->
  GetJsonWithAccessToken`, and 5061 through
  `webp.init -> ansimage.init -> go-term-markdown.init -> cmd.init`. The full
  normalized trace objects are covered by the reachable summary hash above.

Because no selected version changed, exact get manufactured only tidy-removable
metadata, and the dependency module stop rule failed, downstream repository,
snapshot/Docker, acceptance, and Q0-Q2 gates were inapplicable and were not
claimed. The previously accepted Circbuf quality baseline remains unchanged.

## Evidence And Corrections

Exact Go 1.26.7 was recreated below session scratch; its archive SHA-256 is
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`
and binary SHA-256 is
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Go 1.18.10 archive/binary SHA-256 values are
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`
and `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
GOENV/GOWORK/GOTOOLCHAIN remained off/off/local with no ambient GOFLAGS.
Govulncheck's rebuilt v1.7.0 binary is a nonportable receipt at SHA-256
`3d287f71f53e4dfaf03a7b9eb998ca7b3c8363f7db7dbc6c7dcf9f4162a394c1`.

The first Go-1.18 CGO setup attempted an Xcode resolver cache outside scratch;
it was superseded by final runs using the SDK and clang paths directly plus
scratch-local TMPDIR/GOTMPDIR. Initial patterned module-mode govulncheck calls
were invalid and superseded by correct no-pattern module scans. Concatenated
JSON was decoded before filtering trace-bearing symbol findings. No correction
mutated repository files or concealed a candidate failure.

The sealed evidence tree has 315 verified entries. Evidence-manifest SHA-256
is `1dbe7063e72dade7bc31e9c8967da78a60ef97859b68562ffa1a07b75a0b3b0b`;
decision-summary SHA-256 is
`0f5413813a179949c2dbce2a029becfceb9f650f7a603ac8b19755c3f4d48733`.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6` as one
bounded dependency group. Resolve canonical latest, release qualification,
default-branch history, and the highest Go-1.18-floor-compatible candidate
from primary evidence. Make an exact dependency selection only if it changes a
selected version, preserves the retained floor through the complete minimal
closure, and passes every quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through accepted Circbuf pseudo-version
`v0.0.0-20190214190532-5111143e8da2`. All earlier rejections, no-change
decisions, accepted closures, and evidence corrections remain final. Do not
revisit Circbuf or combine another module group. P8 remains queued.

Project MVS already selects Consul API's 2018 pseudo-version although it is not
an explicit `go.mod` requirement. A minimal post-Circbuf survey finds an empty
exact stable proxy list. Proxy and exact Go `@latest` and `@master` resolve the
same selected pseudo-version at 2018-02-02T20:16:55Z; exact `@v0` has no
matching version. The authoritative repository currently reports public,
enabled, unarchived, non-fork status, default branch `master`, one branch, zero
tags, and zero GitHub Releases. Treat canonical qualification, source identity,
signatures, release and alternate-path history, complete closure, self-tests,
loaded consumers, actual symbols, and vulnerability effect as unknown until
independently resolved. Do not call an untagged pseudo-version a stable release.

# Measurements At Start

Latest dependency implementation is Circbuf commit
`3be2183ee310ccdc358ce4ed372c0785de25b88b`, exact parent
`7a0ca4e2caba6d2fff20a9c169b181f9187b40dc`, tree
`a51846840b5f252a30359530dfc950f811398431`, changing only `go.mod` and
`go.sum`. Operator-authorized lifecycle repair
`f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between
the earlier Repr implementation and direct-child Repr handoff. Preserve it.

Circbuf advanced from 2015 pseudo-version
`v0.0.0-20150827004946-bbbad097214e` to canonical latest and master-head
`v0.0.0-20190214190532-5111143e8da2`, commit
`5111143e8da2e98b4ea6a8f32b9065ea1821c191`, tree
`2ab2d9cf2632ab7b549f7da7f081dbe868a697db`, at
2019-02-14T19:05:32Z. It is an unreleased pseudo-version, not a stable release.
The exact path has no stable or prerelease versions, tags, GitHub Releases,
retractions, or authoritative alternate module path. The candidate adds only
one `go.mod` line relative to the selected source; every Go source and test
file is identical. It declares no Go version and no requirements, so its
complete minimal closure is Circbuf alone and preserves Go 1.18.

Candidate proxy and exact-commit forms independently pass verify/list,
count-1/count-10/race tests, and vet without mutation. Exact Go 1.18.10 also
passes count-10/race/vet. Circbuf loads in zero Ply packages; the historical
and selected Serf consumers plus an external Go-1.18 fixture exercise
`NewBuffer`, `Write`, `TotalWritten`, `Size`, `String`, and `Bytes`.

The exact get added one indirect requirement, one main graph edge, and the
candidate checksum pair. Current accepted measurements are 234 selected
modules, 3,581 graph edges, 429 native complete-test packages, 41 loaded
modules, 1,045 `go.sum` lines, and a 361-line unapplied tidy projection. Tidy
removes the explicit Circbuf pin and its two candidate checksum lines while
leaving inherited debt unchanged. Relative to accepted go-cmp commit c314bcb,
accepted metadata adds exactly 29 checksum lines. The main module retains Go
1.18 and toolchain Go 1.26.7. Ordinary and ignored status must be empty.

Circbuf's 430-entry evidence-manifest SHA-256 is
`2684f8c269376320bcfca3404bfbf10c638c81753e6af36b5be59c4df2370d62`;
decision-summary SHA-256 is
`17441163b7924e2cd61cb167b739562da7672fa5abfa10089da36ef906f07ff9`.
Exact quality passed all 27 Q0-Q2 rows at L2 with 80/80 mutations killed and
zero held, regressed, not-comparable, or dirty counts. Full audit exited 1 only
for queued L3 rows Q3.1/Q3.3/Q3.4/Q3.7. Fresh primary vulnerability evidence
has no Circbuf record or trace and preserves exact 20/30/20 Darwin-symbol/
Darwin-module/Windows-symbol populations.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256 `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in PATH, keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and inject
no ambient GOFLAGS. Recreate pinned tools beneath scratch as needed. Portable
receipts remain golangci-lint 2.12.2 archive
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and apidiff
`0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
Treat rebuilt golangci-lint and govulncheck hashes as nonportable receipts;
prove versions and functionality.

# Role And Boundaries

From fresh external archives and caches, resolve Consul API versions through
the Go proxy, checksum database, authoritative upstream repository, and
primary Go vulnerability data. Record exact selected and candidate commits/
times, module Go declarations and requirements, checksum pairs, source
identity, tag and commit signature status, release history, and archived/
deprecated state. Explicitly distinguish stable versions, prereleases, pseudo-
versions, retractions, forks, alternate module paths, branch heads, and
unreleased commits. Do not treat the selected pseudo-version as a stable
release without authoritative release evidence.

Prove canonical latest and the highest qualified exact-path version compatible
with Go 1.18 from declarations and the complete closure. Measure old versus
candidate modules, graph edges, complete packages, checksums, loaded packages
and paths, explicit exact-get diff, and `go mod tidy -diff`. Explain every
selection, edge, and checksum change. Identify real consumers and exercise the
actually used symbols.

Require candidate module complete, repeated, and race-enabled self-tests plus
vet; repository verify/build, complete tests, race, vet, Windows build, pinned
lint, byte-identical public help, API/CLI reports, and primary vulnerability
identity. Preserve the retained Go 1.18 floor across the complete changed
closure. Keep standalone test-only closure separate from project MVS.

Reject or retain without dependency edits if canonical qualification, source
identity, floor compatibility, closure, module tests, loaded behavior, or any
applicable quality contract fails. Evaluate an unreleased canonical pseudo-
version under the recorded dependency policy without relabeling it as stable.
Do not manufacture metadata when exact get changes no selected version. Record
precise primary old/candidate vulnerability IDs and traces.

# Required Reading

Verify branch, clean ordinary and ignored status, exact ancestry, reciprocal
archive history, `./codex-dev-start.sh --check`, and P7/P8 state before work.
Read this archive, the rolling handover, P7 roadmap, `go.mod`, `go.sum`, the
answered Circbuf, Optional, Template, Colour, Chroma, Kong, Repr, Assert, and
Units archives, retained Kingpin/Resty/Errgo/Check/YAML outcomes, and the
toolchain, quality, baseline, compatibility, snapshot/Docker, acceptance,
audit, and lifecycle contracts. Do not reopen earlier decisions.

# Three Moves

If and only if a qualified candidate changes the selected version and its
complete minimal closure preserves Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/armon/consul-api@<candidate>` for one dependency-only
commit. Do not hand-edit metadata or use tidy as implementation. Preserve every
retained version, the toolchain declarations, production source, quality
apparatus, and release inputs.

If canonical qualification confirms the already-selected version and exact
get changes no selection, retain it without adding a redundant direct or
indirect requirement, main edge, or full checksum solely for metadata. A no-
change decision gets no dependency implementation commit.

After a changed selection, run the complete P7 dependency gate: module and
consumer tests, graph/path/checksum/tidy explanation, repository verify/build/
tests/race/vet/Windows/pinned lint, help/API/CLI, launcher and Make contracts,
preflight, host plus fresh snapshot/Docker meta and acceptance, audit meta,
focused and exact Q0-Q2 audits, separate full audit, vulnerability comparison,
empty-HOME count-2, and final ordinary/ignored cleanliness. Exact
`make quality` must exit 0 with all 27 rows at L2 and zero held, regressed,
not-comparable, or dirty counts. Full audit may exit 1 only for established
queued L3 rows, never 2.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`. Never
create direct `/private/tmp/ply-*` roots, never bypass scratch cleanup, and
never run `go mod download all` inside a measured tree. Warm only exact needed
closures from a separate scratch archive. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the Consul API decision, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
single P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
