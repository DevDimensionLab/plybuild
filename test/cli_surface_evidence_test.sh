#!/usr/bin/env bash

# Exercise the existing script's GO boundary without building a second CLI.
set -euo pipefail
repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
if [[ $# -gt 0 ]]; then
	evidence_root=$1
	mkdir -p "$evidence_root"
else
	evidence_root=$(mktemp -d "${TMPDIR:-/tmp}/ply-cli-evidence.XXXXXX")
	trap 'status=$?; if [[ $status -eq 0 ]]; then rm -rf -- "$evidence_root"; else printf "evidence regression retained at %s\n" "$evidence_root" >&2; fi' EXIT
fi
evidence_root=$(cd "$evidence_root" && pwd -P)
cat >"$evidence_root/go-fixture" <<'FIXTURE'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == build && $2 == -o && $4 == ./cmd/ply && $# == 4 ]]
printf '%s\n' "$3" >"$FIXTURE_RECORD"
printf 'fixture build stdout\n'
printf 'fixture build stderr\n' >&2
if [[ $FIXTURE_MODE == build ]]; then
	exit 37
fi
cat >"$3" <<'CLI'
#!/usr/bin/env bash
printf 'fixture CLI stdout\n'
printf 'fixture CLI stderr\n' >&2
exit 9
CLI
chmod +x "$3"
FIXTURE
chmod +x "$evidence_root/go-fixture"

for scenario in build cli; do
	case_root="$evidence_root/$scenario"
	mkdir -p "$case_root/tmp"
	printf 'unrelated\n' >"$case_root/tmp/sentinel"
	printf '%s\n' "GO=$evidence_root/go-fixture FIXTURE_MODE=$scenario bash $repo_root/test/cli_surface_contract_test.sh" >"$case_root/argv"
	set +e
	GO="$evidence_root/go-fixture" FIXTURE_MODE="$scenario" FIXTURE_RECORD="$case_root/binary-path" \
		TMPDIR="$case_root/tmp" bash "$repo_root/test/cli_surface_contract_test.sh" \
		>"$case_root/stdout" 2>"$case_root/stderr"
	status=$?
	set -e
	printf '%s\n' "$status" >"$case_root/exit"
	if [[ $scenario == build ]]; then expected_exit=37; else expected_exit=1; fi
	[[ $status == "$expected_exit" ]] || { printf '%s exit: expected %s, got %s\n' "$scenario" "$expected_exit" "$status" >&2; exit 1; }
	retained_root=$(dirname "$(cat "$case_root/binary-path")")
	[[ -d $retained_root ]]
	grep -Fx "cli surface contract: failure evidence retained at $retained_root" "$case_root/stderr"
	[[ $(cat "$case_root/tmp/sentinel") == unrelated ]]
	[[ $(find "$case_root/tmp" -type d -name 'ply-cli-surface.*' | wc -l | tr -d ' ') == 1 ]]
	if [[ $(uname -s) == Darwin ]]; then mode=$(stat -f '%Lp' "$retained_root"); else mode=$(stat -c '%a' "$retained_root"); fi
	[[ $mode == 700 ]]
	grep -Fx 'fixture build stdout' "$retained_root/build.stdout"
	grep -Fx 'fixture build stderr' "$retained_root/build.stderr"
	if [[ $scenario == build ]]; then
		grep -Fx 'fixture build stderr' "$case_root/stderr"
	else
		grep -Fx 'fixture CLI stdout' "$retained_root/root.stdout"
		grep -Fx 'fixture CLI stderr' "$retained_root/root.stderr"
	fi
done
printf 'cli surface evidence: PASS (build exit 37, CLI assertion exit 1, retained streams and private roots)\n'
