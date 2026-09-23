# Agent Session: Decide Posener Complete Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-23T092359+0200-decide-posener-complete-product-direction`
Created: `2026-09-23T09:23:59+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `3bb44d86808558eccec2f9df95e3bf6163fccf7189985e7ba0906f2e9ff6f003`
Previous: [2026-09-23T082105+0200-evaluate-posener-complete-dependency.md](2026-09-23T082105+0200-evaluate-posener-complete-dependency.md)
Next: [2026-09-23T100201+0200-evaluate-prometheus-client-golang-dependency.md](2026-09-23T100201+0200-evaluate-prometheus-client-golang-dependency.md)
Outcome: Option 1 is final. Exact selected, inherited, indirect, unloaded `github.com/posener/complete v1.2.3` remains unchanged under a Complete-specific, unqualified, non-transferable exception. It is latest but not qualified or fixed; no study or dependency change was authorized, and one bounded Prometheus client_golang evaluation was prepared.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Choose exactly one bounded product direction for graph-selected exact
`github.com/posener/complete v1.2.3` from the completed evaluation recorded in
this archive. Exactly one of the three choices below may be selected. Do not
repeat the evaluation, change a dependency, run a study, combine another
dependency group, or begin P8.

# Defensive Scope

This is an ordinary dependency-quality product decision. Use only the answered
evaluation and narrow guard-only public metadata, graph, build, and advisory
checks. Do not fuzz, stress, probe resource exhaustion, create oversized,
deeply nested, cyclic, malformed, adversarial, or escape-sequence payloads,
reproduce a security issue, or perform security or exploitability analysis.

Every disposable cache, tool, report, or response must remain beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to `/private/tmp`, `/tmp`, a
sibling of the managed root, or another external root. Verify containment and
remove task-owned scratch evidence before handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active only at this Complete product decision;
P8 remains queued. No dependency implementation or owner/request study is
authorized by this archive. Exactly one of the three bounded choices below may
be selected.

# Measurements At Start

The completed evaluation found six canonical exact-path stable versions,
v1.1.1, v1.1.2, and v1.2.0 through latest/selected v1.2.3. The public,
enabled, unarchived, non-fork MIT repository `posener/complete`, GitHub ID
90418143, remains owned by `posener`. Default branch `v1` explicitly directs
current development to `master`; all canonical v1 tags are lightweight,
valid-signed commits and ancestors of `v1`. The separate canonical `/v2` line
through v2.1.0 is an ineligible alternate module path; `/v3` is absent. There
are no exact-path prereleases, replacements, retractions, or deprecations.

Selected v1.2.3 is commit
`05b68ffc813dd10c420993cb1cf927b346c057b8`, tree
`58ea2366a478ec75ecf97d4d30097a3f16709fab`. Exact go-import, release, tag,
commit/tree/signature/ancestry, proxy/sumdb, archive-to-Git, module, and license
identities are recorded in the answered evaluation and remain final unless a
guard-only check proves an input changed. V1.2.3 declares Go 1.13 and requires
go-multierror v1.0.0 and Testify v1.4.0.

V1.1.1 and v1.1.2 fail build, tests, race, and vet under exact Go 1.18.10 and
Go 1.26.7 because their synthesized module metadata omits imported
`github.com/hashicorp/go-multierror`. V1.2.0 through v1.2.3 preserve Go 1.18
and pass module verification, build, count-one/count-ten tests, race, vet, and
the recorded supported cgo-disabled cross-build/test-compilation matrix under
both exact SDKs. They nevertheless do not qualify: every buildable v1.2
release implements exported `cmd/install.Uninstall` by creating a hard-coded
`/tmp/complete-*` temporary file, copying it back, and never removing it after
successful completion. Static source identity is sufficient; the path was not
executed because it would write outside the managed scratch root. Passing
upstream gates and advisory absence do not override this ordinary cleanup
failure. No canonical exact-path stable qualifies.

Exactly three graph requests target the module path. Historical Mitchellh CLI
v1.0.0 and v1.1.0 each request v1.1.1 and genuinely import the target in their
production command/autocomplete implementation; Hashicorp Serf v0.9.6 requests
selected v1.2.3 only in module metadata. Preserve the three exact requester
versions, requested target versions, genuine-import/metadata-only boundaries,
and every current or historical route recorded by the answered evaluation.

The shortest selected route is main -> direct mvn-pom-mutator v0.2.3 ->
historical Viper v1.10.1 -> Serf v0.9.6 -> selected target. The two v1.1.1
routes continue through the recorded Viper/crypt, Consul API, Serf, and CLI
vertices. Target, selected Serf, and selected CLI why results are negative.
The repository has zero target, Serf, or CLI imports; their packages are absent
from production and complete-test loads; target module-backed load and runtime
relevance are zero; and no current or historical main target root exists.
Physical graph selection or a route does not establish qualification,
ownership, loading, or runtime relevance.

A disposable exact selected get changes no selected version. It manufactures
three unused indirect roots/main edges and three source sums, producing 234
modules, 3,602 graph edges, unchanged 355/429/197/41 loads, and 1,070 sum
lines. Normal tidy removes those roots and sums and returns the exact common
234-module/3,557-edge, 52/948-line projection. No projection was retained.

The real project remains 234 modules, 3,599 edges, 355 production entries,
429 complete-test entries, 197 module-backed entries across 41 loaded modules,
and 1,067 sum lines. Its `go.mod`, `go.sum`, and graph identities, common tidy
projection, Go 1.18 floor, 47-selection/276-edge pre-Goe guard, separate Goe,
pkg/errors, SFTP, and go-difflib guards, and every earlier decision remain
exact. Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Exact target OSV, GitHub global, repository advisory, and isolated pinned
govulncheck findings are empty. The unchanged project retains exact
30/22/20/20 module/package/symbol/test-symbol advisory populations without a
target trace. Guard OSV remains limited to the recorded Gorilla WebSocket and
go-retryablehttp pairs; x/mod v0.14.0 retains GO-2026-6179 and GO-2026-6180.
The 518,501-byte/1,402-record Go index and PUBLISHED 2,807-byte memberlist CNA
response remain byte-exact. Advisory absence is not qualification.

# Earlier Decisions Remain Closed

Exact go-difflib v1.0.0, SFTP v1.13.1, and pkg/errors v0.9.1 remain retained
only under their own explicit target-specific, unqualified, non-transferable
exceptions. Preserve their exact requests, requester boundaries, routes,
why/import/load/runtime/root facts, source/behavior/closure/projection/advisory
identities, mutual guards, and expiry conditions. Do not run their rejected
owner studies, reopen them, or transfer an exception.

Exact Goe, ULID, go-conntrack, mapstructure, go-homedir, Promptui, emoji/v2,
kr/text, kr/pty, kr/pretty, kr/logfmt, kr/fs,
go-windows-terminal-sequences, gotool, errcheck, httprouter, GLS,
go-junit-report, json-iterator, and clockwork exceptions remain final,
target-specific, unqualified, and non-transferable. Concurrent, reflect2,
Cast, Viper, memberlist, every earlier selection/owner/request/route/fact, and
all qualification/expiry guards remain final. Do not reopen, broaden, or
transfer any decision.

# Role And Boundaries

Act only as the bounded Posener Complete product decision-maker. Treat the
answered evaluation as final, use guard-only checks to detect changed inputs,
and stop for a fresh owning evaluation if any input changed. Do not perform new
behavior, closure, projection, owner-study, or dependency implementation work.

# Authorized Choice

Choose exactly one:

1. Explicitly retain exact selected, inherited, indirect, unloaded
   `github.com/posener/complete v1.2.3` without source or dependency-metadata
   changes under a Complete-specific, unqualified, non-transferable exception.
   State that it is latest but not qualified or fixed. Bind retention to all
   three exact requests, both genuine CLI import boundaries, the Serf metadata-
   only boundary, every route, negative target/requester why, zero repository
   import/load/runtime/root state, all release/source/behavior/closure/
   projection/advisory identities, every earlier guard, and explicit expiry on
   any changed input. Then prepare only the next bounded selected dependency
   queue evaluation named by the roadmap.
2. Authorize exactly one later, measurement-only **Mvn-Pom-Mutator Complete
   Owner/Request Study**. It may measure only whether a supported stable update
   of the existing direct mvn-pom-mutator owner can eliminate all three target
   requests through its recorded historical Viper/crypt/Consul API/Serf/CLI
   graph while preserving product behavior, every direct root, the Go 1.18
   floor, and every earlier guard. It may not implement, add a target root,
   patch/vendor/fork upstream, change an unrelated requester, select a branch,
   pseudo-version, prerelease, replacement, or alternate module path, execute
   the leaking uninstall path, or call v1.2.3 qualified. Any viable measured
   route still requires fresh owning dependency evaluations and a product
   decision before implementation.
3. Stop P7 unresolved without source, dependency, owner/requester, roadmap-
   queue, or exception changes and prepare no dependency evaluation.

Do not silently retain v1.2.3, transfer an earlier exception, or use physical
selection, zero loading, or advisory absence as acceptance. A decision is
explicit only when its exact scope, unqualified status, guards, and expiry
conditions are recorded.

# Required Reading

Read this archive, the answered Complete evaluation, answered go-difflib,
SFTP, pkg/errors, Goe, ULID, and go-conntrack decisions/evaluations, rolling
handover, P7/P8 roadmap, `go.mod`, and `go.sum`. Perform guard-only verification
of branch/HEAD/parent/tree, exact changed set, ancestry, ordinary and ignored
cleanliness, reciprocal archive chain, launcher check, exact Go identities,
project hashes/counts/tidy state, all target request/route/why/import/load/
runtime facts, every earlier guard, and narrow fresh advisory identities. Do
not repeat completed upstream behavior, closure, race, cross-build, archive,
projection, or govulncheck work unless a guard-only check proves an input
changed; stop for a fresh owning evaluation if one did.

# Three Moves

First, verify the completed record and choose exactly one authorized option.
Second, record only that decision without source or dependency-metadata
changes and without executing a study. Third, update roadmap and rolling
handover, answer this archive, prepare at most the one reciprocal successor
required by the chosen option, verify scratch containment and cleanup, and
make the local handoff commit.

# Automatic Handoff

Do not launch a successor or study; push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, edit a dependency, transfer an
exception, reopen Complete, go-difflib, SFTP, pkg/errors, or an earlier group,
write outside the managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Option 1 is final. Exact selected, inherited, indirect, unloaded
`github.com/posener/complete v1.2.3` is explicitly retained without product-
source or dependency-metadata changes under a Complete-specific, unqualified,
non-transferable exception. It is the latest exact-path v1 stable but is not
qualified or fixed: every buildable v1.2 stable leaves the hard-coded
`/tmp/complete-*` temporary file created by exported `cmd/install.Uninstall`
after successful completion. The operation was not executed. Physical graph
selection, zero loading, and advisory absence are not acceptance. Option 2's
Mvn-Pom-Mutator Complete Owner/Request Study was not authorized or run, and
option 3 was not selected.

Retention is bound to exactly three requests: Mitchellh CLI v1.0.0 and v1.1.0
each request Complete v1.1.1 and genuinely import it in production command/
autocomplete code, while Serf v0.9.6 requests selected v1.2.3 only in module
metadata. It requires the shortest selected route main -> direct mvn-pom-
mutator v0.2.3 -> historical Viper v1.10.1 -> Serf v0.9.6 -> v1.2.3 and both
recorded v1.1.1 routes through the historical Viper/crypt, Consul API, Serf,
and CLI vertices. Target, selected Serf, and selected CLI why remain negative;
repository target/Serf/CLI imports, production and complete-test loads, target
module-backed load and runtime relevance, and current or historical main
target roots remain zero. Any changed request, requester import boundary,
route, owner, why/import/load/runtime fact, root, or supported owner route
expires the exception.

The exact six-release v1 line, separate ineligible `/v2` line, absent `/v3`,
repository/owner/status/license/default-branch, tag/commit/tree/signature/
ancestry, proxy/sumdb/archive/module identities, and completed source/API/
behavior/closure/race/vet/cross-build results remain guards. Selected v1.2.3
remains commit `05b68ffc813dd10c420993cb1cf927b346c057b8`, tree
`58ea2366a478ec75ecf97d4d30097a3f16709fab`. The selected-get
234/3,602/355/429/197/41/1,070 state and its ordinary tidy return to the common
234/3,557, 52/948-line projection remain guards; neither projection is
retained.

The real project remains 234 modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed entries across 41 loaded
modules, and 1,067 sum lines. Its exact `go.mod`, `go.sum`, graph, common-tidy,
Go-1.18-floor, 47-selection/276-edge pre-Goe, separate Goe/pkg-errors/SFTP/go-
difflib, and every earlier decision guard remain exact. Exact target advisory
results remain empty; the project 30/22/20/20 populations have no target trace;
the recorded Gorilla WebSocket, go-retryablehttp, x/mod, Go-index, and
memberlist-CNA advisory identities remain guards. Advisory absence did not
qualify the target.

Guard-only revalidation reproduced the clean evaluation handoff, exact changed
set and ancestry, reciprocal archive chain, launcher state, both exact SDK
identities, project hashes/counts/tidy state, every target request/route/why/
import/load/runtime fact, earlier guards, repository/release status, and narrow
fresh advisory identities. Final unchanged-project exact Go 1.26.7 module
verification, build, count-one tests, race count-one tests, and vet pass. No
completed upstream behavior, closure, projection, archive, cross-build, or
govulncheck work was repeated. One accidental bootstrap-toolchain invocation
completed only build and count-one tests before it was stopped; all authorized
final gates were then rerun with the verified exact Go 1.26.7 binary.

Any release/support, source/behavior/closure, request/route/relevance,
projection/project/Go-floor, earlier-guard, advisory/finding, independent
defect, qualification, or compatible supported-route change expires retention
and requires a fresh owning Complete evaluation and explicit product decision
before merge. Go-difflib, SFTP, pkg/errors, Goe, ULID, go-conntrack, and every
earlier exception remain separate, unqualified, and untransferred.

No source, `go.mod`, or `go.sum` change remains. P7 continues only with the
prepared bounded evaluation of exact selected
`github.com/prometheus/client_golang v1.4.0`; it was not launched or executed.
P8 remains queued.
