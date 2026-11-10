#!/usr/bin/env bash

# Exercise the existing script's GO boundary without building a second CLI.
set -euo pipefail
repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
if [[ $# -gt 0 ]]; then
	evidence_root=$1
	mkdir -p "$evidence_root"
else
	. "$repo_root/scripts/compat_root.sh"
	compat_run_root 'cli surface evidence' ply-cli-evidence
	evidence_root=$run_root
fi
evidence_root=$(cd "$evidence_root" && pwd -P)
cat >"$evidence_root/go-fixture" <<'FIXTURE'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == build && $2 == -o && $4 == ./cmd/ply && $# == 4 ]]
printf '%s\n' "$3" >"$FIXTURE_RECORD"
printf 'fixture build stdout\n'
printf 'fixture build stderr\n' >&2
if [[ $FIXTURE_MODE == hang ]]; then
	printf '%s\n' "$$" >"$FIXTURE_RECORD.hang"
	exec sleep 60
fi
if [[ $FIXTURE_MODE == build ]]; then
	exit 37
fi
if [[ $FIXTURE_MODE == help ]]; then
	cat >"$3" <<'CLI'
#!/usr/bin/env bash
printf '%s\n' '== version: dev ==' 'Status functionality for a project' \
	'Perform upgrade on existing projects' 'Builds a ply project' "fixture help for $*"
CLI
	chmod +x "$3"
	exit 0
fi
cat >"$3" <<'CLI'
#!/usr/bin/env bash
printf 'fixture CLI stdout\n'
printf 'fixture CLI stderr\n' >&2
exit 9
CLI
chmod +x "$3"
FIXTURE
mkdir -p "$evidence_root/find-bin"
printf '#!/usr/bin/env bash\nexit 42\n' >"$evidence_root/find-bin/find"
chmod +x "$evidence_root/go-fixture" "$evidence_root/find-bin/find"

fixture_parent=$(mktemp -d "${PLY_ROUNDTRIP_ROOT:?set PLY_ROUNDTRIP_ROOT for the surface test fixtures}/ply-cli-evidence root.XXXXXX")
trap 'rm -rf -- "$fixture_parent"; [[ -z ${run_root:-} ]] || compat_cleanup' EXIT
compat_root=$(cd "$repo_root" && pwd -P)/.build/compatibility
listing() { ls -A "$fixture_parent"; ls -A "$compat_root" 2>/dev/null || true; }
listing >"$evidence_root/before"
surface() {
	exec env GO="$evidence_root/go-fixture" FIXTURE_MODE="$scenario" FIXTURE_RECORD="$case_root/binary-path" \
		TMPDIR="$case_root/tmp" PLY_ROUNDTRIP_ROOT="$fixture_root" PATH="$path" \
		bash "$repo_root/test/cli_surface_contract_test.sh" >"$case_root/stdout" 2>"$case_root/stderr"
}
for scenario in build cli find help root hang; do
	case_root="$evidence_root/$scenario"
	mkdir -p "$case_root/tmp"
	printf 'unrelated\n' >"$case_root/tmp/sentinel"
	printf '%s\n' "GO=$evidence_root/go-fixture FIXTURE_MODE=$scenario bash $repo_root/test/cli_surface_contract_test.sh" >"$case_root/argv"
	if [[ $scenario == root ]]; then fixture_root=; else fixture_root=$fixture_parent; fi
	path=$PATH
	[[ $scenario != find ]] || path="$evidence_root/find-bin:$PATH"
	if [[ $scenario == hang ]]; then
		set -m; surface & pid=$!; set +m
		for _ in $(seq 100); do [[ ! -s $case_root/binary-path.hang ]] || break; sleep 0.1; done
		[[ -s $case_root/binary-path.hang ]] || { printf 'hang child never started\n' >&2; exit 1; }
		kill -TERM -- "-$pid"
	fi
	set +e
	if [[ $scenario == hang ]]; then wait "$pid"; else (surface); fi
	status=$?
	set -e
	printf '%s\n' "$status" >"$case_root/exit"
	case $scenario in build) expected_exit=37 ;; cli|find|help) expected_exit=1 ;; root) expected_exit=2 ;; hang) expected_exit=143 ;; esac
	[[ $status == "$expected_exit" ]] || { printf '%s exit: expected %s, got %s\n' "$scenario" "$expected_exit" "$status" >&2; cat "$case_root/stderr" >&2; exit 1; }
	[[ $(ls -A "$case_root/tmp") == sentinel && $(cat "$case_root/tmp/sentinel") == unrelated ]]
	listing | cmp -s - "$evidence_root/before" || { printf '%s left run data behind\n' "$scenario" >&2; exit 1; }
	! grep -F 'retained' "$case_root/stderr"
	case $scenario in
	root)
		[[ ! -e $case_root/binary-path ]]
		grep -F 'cli surface contract: unsuitable PLY_ROUNDTRIP_ROOT' "$case_root/stderr" ;;
	build)
		[[ $(dirname "$(cat "$case_root/binary-path")") == "$(cd "$fixture_parent" && pwd -P)"/ply-cli-surface.* ]]
		grep -Fx 'fixture build stderr' "$case_root/stderr" ;;
	hang)
		! kill -0 "$(cat "$case_root/binary-path.hang")" 2>/dev/null ;;
	help)
		grep -Fx 'cli surface contract: task help lost journal' "$case_root/stderr"
		grep -Fx 'fixture help for workspace task --help' "$case_root/stderr" ;;
	cli|find)
		grep -Fx 'fixture CLI stdout' "$case_root/stderr"
		grep -Fx 'fixture CLI stderr' "$case_root/stderr" ;;
	esac
done
printf 'cli surface evidence: PASS (build exit 37, CLI assertion exit 1, failing find exit 1, older stream shown, missing root exit 2, TERM exit 143, failure streams on stderr, owned roots removed)\n'
