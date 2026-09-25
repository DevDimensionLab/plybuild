#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
script_under_test=${SCRIPT_UNDER_TEST:-"$repo_root/scripts/search-replace.sh"}
temp_root=$(mktemp -d)
cleanup_temp_root() {
	find "$temp_root" -type d -exec chmod u+rwx {} + 2>/dev/null || true
	rm -rf -- "$temp_root"
}
trap cleanup_temp_root EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM

fail() {
	printf 'search-replace meta-test: %s\n' "$*" >&2
	exit 1
}

[[ -f "$script_under_test" && ! -L "$script_under_test" ]] ||
	fail 'script under test must be a regular, non-symlink file'
/bin/bash -n "$script_under_test" || fail 'script does not parse'
grep -Fx 'set -euo pipefail' "$script_under_test" >/dev/null ||
	fail 'script does not enable errexit, nounset, and pipefail'

fake_bin="$temp_root/bin"
mkdir -p "$fake_bin"

cat >"$fake_bin/jq" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
[[ ${1-} == '-r' && ${2-} == .* && $# -eq 2 ]] || exit 91
key=${2#.}
payload=$(command cat)
printf '%s\n' "$payload" | command sed -n "s/.*\"$key\":\"\([^\"]*\)\".*/\1/p"
EOF

cat >"$fake_bin/gsed" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
[[ ${1-} == '-i' && $# -eq 3 ]] || exit 92
expression=$2
target_file=$3
line_number=${expression%%s/*}
substitution=${expression#*s/}
old_value=${substitution%%/*}
substitution=${substitution#*/}
new_value=${substitution%%/*}
temporary_file="$target_file.recording-gsed"
command awk -v line="$line_number" -v old="$old_value" -v new="$new_value" '
	NR == line { gsub(old, new) }
	{ print }
' "$target_file" >"$temporary_file"
command mv "$temporary_file" "$target_file"
EOF

chmod +x "$fake_bin/jq" "$fake_bin/gsed"

write_pom() {
	local target=$1
	cat >"$target" <<'EOF'
<dependency>
  <groupId>example.group</groupId>
  <artifactId>example-artifact</artifactId>
  <version>100</version>
</dependency>
EOF
}

outdated_record='{"type":"outdated dependency","versionIsProperty":"false","artifactId":"example-artifact","groupId":"example.group","oldVersion":"100","newVersion":"200"}'
current_record='{"type":"current dependency","versionIsProperty":"false","artifactId":"example-artifact","groupId":"example.group","oldVersion":"100","newVersion":"200"}'

exercise_replacement() {
	local script=$1
	local name=$2
	local workspace="$temp_root/$name"
	mkdir -p "$workspace"
	write_pom "$workspace/pom.xml"
	(
		cd "$workspace"
		PATH="$fake_bin:$PATH" /bin/bash "$script" <<<"$outdated_record"
	) >"$workspace/output"
	grep -F '=> [outdated dependency] example-artifact:example.group (100 -> 200) in pom.xml' \
		"$workspace/output" >/dev/null || return 1
	grep -F '<version>200</version>' "$workspace/pom.xml" >/dev/null || return 1
	if grep -F '<version>100</version>' "$workspace/pom.xml" >/dev/null; then
		return 1
	fi
}

exercise_replacement "$script_under_test" working ||
	fail 'script did not replace the selected dependency version'

irrelevant_workspace="$temp_root/irrelevant"
mkdir -p "$irrelevant_workspace"
write_pom "$irrelevant_workspace/pom.xml"
cp "$irrelevant_workspace/pom.xml" "$irrelevant_workspace/expected.xml"
(
	cd "$irrelevant_workspace"
	PATH="$fake_bin:$PATH" /bin/bash "$script_under_test" <<<"$current_record"
) >"$irrelevant_workspace/output"
cmp -s "$irrelevant_workspace/expected.xml" "$irrelevant_workspace/pom.xml" ||
	fail 'script rewrote a non-outdated dependency'
[[ ! -s "$irrelevant_workspace/output" ]] ||
	fail 'script reported a non-outdated dependency'

pipefail_workspace="$temp_root/pipefail"
mkdir -p "$pipefail_workspace"
write_pom "$pipefail_workspace/pom.xml"
if (
	echo() { return 42; }
	export -f echo
	cd "$pipefail_workspace"
	PATH="$fake_bin:$PATH" /bin/bash "$script_under_test" <<<"$current_record"
) >"$pipefail_workspace/output" 2>"$pipefail_workspace/error"; then
	fail 'script hid an upstream pipeline failure'
fi

mutant="$temp_root/search-replace-noop.sh"
awk '
	index($0, "gsed -i ") {
		print "          : # mutant removes the replacement"
		mutations++
		next
	}
	{ print }
	END { if (mutations != 1) exit 1 }
' "$script_under_test" >"$mutant" || fail 'could not construct the no-op mutant'
chmod +x "$mutant"
if exercise_replacement "$mutant" mutant; then
	fail 'meta-test accepted a no-op replacement mutant'
fi

printf 'search-replace meta-test: PASS (parse/options, selection, replacement, pipeline failure, and no-op mutant checked)\n'
