# Agent Session: Decide Json-Iterator Go Product Direction

Status: NEXT
Session ID: `2026-09-21T140451+0200-decide-json-iterator-go-product-direction`
Created: `2026-09-21T14:04:51+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `a42e3d185b9a8af3c56270e38dd8add278a6f1c220413ac2900dfda205148e8a`
Previous: [2026-09-21T132433+0200-evaluate-json-iterator-go-dependency.md](2026-09-21T132433+0200-evaluate-json-iterator-go-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by making one bounded product decision for exact-path
`github.com/json-iterator/go`. The completed independent evaluation found no
stable release that preserves Go 1.18 and passes every applicable ordinary
JSON-compatibility, API, repeatability, and project-ownership contract. Choose
exactly one direction below, record its ownership and expiry guards, and stop.
Do not repeat behavior failures, silently accept selected v1.1.12, implement a
dependency or parent change, evaluate another dependency group, reopen
clockwork/mvn-pom-mutator/demangle/pprof/strcase/memberlist work, or begin P8.

# Defensive Scope

This is a documentation-only product decision for an ordinary dependency-
quality review. Reuse the completed public metadata, static source/API,
admissible upstream-test, bounded deterministic fixture, project-graph, and
advisory evidence. Do not fuzz, run the randomized upstream type package,
stress, probe resource exhaustion, generate oversized or deeply nested input,
or perform security or exploitability analysis. The recorded ordinary
map-key, control-escape, non-finite-number, example, and vet results are final;
do not reproduce them.

Every disposable archive, cache, tool, report, and fixture must remain beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to `/private/tmp`, `/tmp`, a
sibling of the managed root, or another external root. Verify containment and
remove task-owned scratch evidence before handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 is blocked only on this json-iterator decision after
exact Go 1.26.7, every accepted dependency move through Google UUID v1.4.0,
qualified go-cleanhttp v0.5.2, and all target-specific retained-module
decisions through exact inherited unloaded clockwork v0.1.0. Every earlier
outcome is final under its own guards. P8 remains queued.

The evaluation left product source, `go.mod`, and `go.sum` unchanged. Selected
json-iterator remains exact stable v1.1.12, inherited through seven historical
requests. It has negative `go mod why -m`, zero repository imports, zero
production or complete-test package loads, and no runtime reachability. Viper
v1.15.0 is a genuine direct, imported, and loaded project parent, but it does
not load a json-iterator package. Physical MVS selection and zero target
loading are not qualification or risk acceptance. No clockwork, demangle,
strcase, memberlist, or other exception transfers.

# Measurements At Start

The json-iterator evaluation began from clean ordinary and ignored state on
branch `codex/upgrade-quality` at handoff HEAD
`80ea158a425503ea5c881c70a391b1bc23d5efbf`, parent
`e5fd2d47f5e50ff46f3a42e7be6d0bc127720277`, tree
`f9a3ad53b6898f5ad25d00677a2623aa81f67ddd`. That handoff changed exactly the
launcher, answered clockwork decision archive, then-NEXT json-iterator
evaluation archive, rolling handover, and roadmap. Verify the new handoff,
reciprocal archive chain, latest Google UUID implementation ancestry, exact Go
identity, and launcher check rather than assuming these facts.

The unchanged project has 234 modules, 3,599 graph edges, 355 production and
429 complete-test entries, 197 module-backed complete-test entries across 41
loaded modules, 1,067 sum lines, and the recorded 432-line tidy projection.
`go.mod` and `go.sum` SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
All 27 earlier guarded selections and their 181 incoming graph edges retain
sorted snapshot SHA-256
`09258fe1e7425eed2d927f8f0be20a60fa610098677e2f31a10b5b4153b35989`;
all guarded why results are negative and guarded imports/loads are zero.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

# Completed Evaluation

Fresh `go-import` metadata resolves exact path `github.com/json-iterator/go`
without redirect to the canonical MIT `json-iterator/go` repository. It is a
non-fork and is now archived. It has exactly eight proxy-valid stable releases
v1.1.5-v1.1.12 on one linear ancestry, no valid v2/v3 module, retraction,
deprecation, replacement, or alternate exact path. Every stable tag is
lightweight and unsigned. Proxy/Git regular files and sumdb identities agree.

Selected/latest stable v1.1.12 is commit
`024077e996b048517130b21ea6bf12aa23055d3d`, tree
`fab0aaa4437b102ed0e2e1f6db945807c4e68318`. Current master is unreleased
pseudo-version `v1.1.13-0.20220915233716-71ac16282d12`, seven commits later.
It fixes one map-key defect but remains ineligible for stable selection and
retains the exact-Go compatibility failure.

All stable minimal production/test module closures preserve the Go 1.18
directive floor and resolve under exact Go 1.26.7 and contained Go 1.18.10.
Compile-only tests pass under both SDKs for v1.1.6-v1.1.12; v1.1.5 alone has
an exact-Go misnamed-example compilation failure. Selected root/extra packages
cross-build under both SDKs for linux/amd64, windows/amd64, and darwin/amd64.
The tree has no command, cgo, generated source, embed, testdata, symlink, or
external resource. Its strict and opt-in sloppy production branches compile.

The initial broad upstream run was discarded after identifying the randomized
`google/gofuzz`-based `type_tests`; it was not executed again or used as
qualification evidence. The other nine admissible packages pass count-one,
count-ten, and race under Go 1.18.10. Under exact Go 1.26.7 all three runs fail
only standard-library byte-compatibility assertions for backspace/form-feed.
Vet under both SDKs reports the same three duplicate JSON tags in upstream
test structs.

No stable release qualifies:

- V1.1.5-v1.1.6 accept a non-finite float without error under both SDKs;
  v1.1.5 also fails exact-Go example compilation.
- V1.1.7-v1.1.9 fail exact-Go byte compatibility because `encoding/json`
  emits `"\\b\\f"` while json-iterator emits `"\\u0008\\u000c"`.
- V1.1.10-v1.1.12 also encode an ordinary string-alias map key implementing
  `encoding.TextMarshaler` as `{"TEXT_key":"value"}` rather than standard
  `encoding/json` output `{"key":"value"}` under both SDKs.
- Unreleased master fixes that map-key difference and passes the bounded Go
  1.18 fixture, but still fails exact-Go control-escape byte compatibility.

The control encodings are both valid JSON, but json-iterator advertises a
100%-compatible `encoding/json` replacement and its own deterministic tests
require byte equality. The selected release therefore fails an applicable
ordinary compatibility contract. Its API additionally exposes documented
prefix/indent and decoder-token compatibility gaps, aliased iterator/stream
buffers, unsafe reflection, and mutable unprotected global codec/extension
registration that must complete before concurrent cached use. Ordinary
struct/map/interface/nil/raw/custom-codec/number/validation/decoder/encoder/
iterator/stream behavior otherwise passed the bounded fixtures.

MVS selects v1.1.12 through seven historical requests across Viper, two crypt
paths, etcd client/v2, and Prometheus. A disposable direct v1.1.12 root retains
234 modules,
produces 3,602 edges and 1,070 sum lines, leaves target loads zero, and does
not change selection. A disposable master root yields 234/3,607/1,071 and
also leaves target loads zero. Tidy removes either manufactured root, restores
v1.1.12, and returns the common 52-line/948-line hashes
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
Lower direct requests cannot create a tidy-stable lower selection without
changing genuine parents and unrelated selections. No projection was applied.

Fresh stable/master OSV, GitHub global/repository, Go-index lookup, and pinned
govulncheck v1.8.0 isolated module/package/symbol/test-symbol results are
empty. Base and both project projections have byte-identical normalized
finding streams, identical 30/22/20/20 project populations, and no target
trace. The 1,402-record Go index remains at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
Guard OSV retains only the recorded Gorilla/go-retryablehttp pairs, and the
2,807-byte PUBLISHED memberlist CNA response remains exact at
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.

# Required Product Decision

Choose exactly one. Option 1 is recommended because the target is completely
unloaded, no stable release qualifies, and retaining the exact inherited
selection preserves every genuine parent and earlier guard:

1. Explicitly retain exact selected, inherited, unloaded v1.1.12 without
   metadata changes under a json-iterator-specific, non-transferable
   exception. Accept only the completed map-key, control-escape, older non-
   finite-number, example/vet, repository/archive, API/global-state/buffer,
   Go-floor, MVS, loading, repeatability, vulnerability, and related findings.
   Guard exact v1.1.12, all seven requests and genuine parent identities,
   negative why/import/load results, runtime unreachability, exact graph/
   module/tidy state, every earlier guard, and no new advisory, independent
   defect, qualified stable release, genuine supported tidy-stable owner, or
   compatible qualified route.
2. Authorize exactly one separate measurement-only owning-parent/request
   study. It may identify whether genuine compatible Viper, crypt, etcd, or
   Prometheus request changes remove or update json-iterator and measure every
   API, behavior, Go-floor, graph, guarded-selection, vulnerability, and
   project consequence. It may recommend a later route but may not implement
   one, add a direct target root, combine parent upgrades, or alter a finalized
   guard without its fresh owning decision.
3. Stop P7 explicitly with selected v1.1.12 unresolved and unaccepted. Prepare
   no implementation or dependency-evaluation successor and do not begin P8.

If none is acceptable, choose option 3. Do not infer acceptance from zero
loading, describe selected as qualified, promote master, select a lower
release by changing parents, or transfer another dependency's exception.

# Required Reading And Moves

Read this archive and its answered evaluation, the answered clockwork and
demangle decisions/evaluations, pprof and strcase records, memberlist ownership
chain, relevant retained-module records, rolling handover, roadmap, `go.mod`,
and `go.sum`. Reuse the evaluation; revalidate only exact decision guards and
fresh advisory identities needed for the choice.

First verify branch, clean ordinary/ignored state, handoff identity/changed
set, reciprocal chain, exact Go identity, module hashes, all seven target
requests and genuine parents, why/import/load state, every guarded selection/
edge, memberlist CNA identity, scratch containment, and
`./codex-dev-start.sh --check`. Stop if a premise changed.

Second, choose and record exactly one option with explicit ownership and
expiry bounds. This is documentation-only. Do not edit product source,
`go.mod`, or `go.sum`; repeat behavior tests; run the option-2 study; add a
root; implement a workaround; alter a parent/guard; or execute P8.

Third, update the roadmap and rolling handover, answer this archive, and
prepare exactly one reciprocal NEXT mission matching the decision, or a
COMPLETE state if option 3 ends the authorized roadmap. Run applicable final
checks, verify scratch containment, and make only the required local
documentation handoff commit. Do not execute a successor.

# Automatic Handoff

Make only the required local `docs: prepare next agent session` commit. Do not
create an implementation commit, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, remove the worktree, combine another
dependency group, reopen clockwork/mvn-pom-mutator/demangle/pprof/strcase/
memberlist work, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
