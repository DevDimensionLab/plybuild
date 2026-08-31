#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

makefile_under_test=${MAKEFILE_UNDER_TEST:-"$repo_root/Makefile"}
config_under_test=${GORELEASER_CONFIG_UNDER_TEST:-"$repo_root/.goreleaser.yml"}
standalone_brew_config=${GORELEASER_BREW_CONFIG_UNDER_TEST:-"$repo_root/.goreleaser.brews.yml"}
fake_goreleaser="$tmp_dir/recording goreleaser"
calls_file="$tmp_dir/goreleaser-calls"
make_workspace="$tmp_dir/make workspace"
snapshot_scripts="$tmp_dir/snapshot scripts"
snapshot_calls="$tmp_dir/snapshot-acceptance-calls"
docker_scripts="$tmp_dir/docker scripts"
docker_calls="$tmp_dir/docker-acceptance-calls"

mkdir -p "$make_workspace" "$snapshot_scripts" "$docker_scripts"

fail() {
	printf 'make distribution contract: %s\n' "$*" >&2
	exit 1
}

assert_local_config() {
	local config=$1
	local active_publishers

	[[ -f "$config" && ! -L "$config" ]] ||
		fail "default GoReleaser config is missing or not a regular file: $config"

	active_publishers=$(awk '
		/^(brews|homebrew_casks|snapcrafts|snaps):[[:space:]]*($|#)/ {
			print NR ":" $0
		}
	' "$config")
	[[ -z "$active_publishers" ]] ||
		fail "default GoReleaser config contains an active package-manager publisher: $active_publishers"

	awk '
		/^[[:alnum:]_-]+:[[:space:]]*($|#)/ {
			section = $1
			sub(/:.*/, "", section)
		}
		section == "release" &&
			/^[[:space:]]+disable:[[:space:]]*(true|"true"|'\''true'\'')([[:space:]#]|$)/ {
			disabled = 1
		}
		END { exit(disabled ? 0 : 1) }
	' "$config" || fail 'default GoReleaser config does not disable remote releases'
}

assert_rejected_config() {
	local config=$1
	local label=$2
	local output=$3

	if (
		GORELEASER_CONFIG_UNDER_TEST="$config" \
		GORELEASER_BREW_CONFIG_UNDER_TEST="$tmp_dir/absent-brew-config" \
		"${BASH_SOURCE[0]}"
	) >"$output" 2>&1; then
		fail "contract accepted $label"
	fi
}

assert_local_config "$config_under_test"

if [[ -e "$standalone_brew_config" || -L "$standalone_brew_config" ]]; then
	fail "standalone Homebrew publisher config remains callable: $standalone_brew_config"
fi

cp "$config_under_test" "$tmp_dir/homebrew.yml"
printf '\nbrews:\n  - name: forbidden\n' >>"$tmp_dir/homebrew.yml"
assert_rejected_config "$tmp_dir/homebrew.yml" 'an active Homebrew publisher' "$tmp_dir/homebrew-output"

cp "$config_under_test" "$tmp_dir/homebrew-cask.yml"
printf '\nhomebrew_casks:\n  - name: forbidden\n' >>"$tmp_dir/homebrew-cask.yml"
assert_rejected_config "$tmp_dir/homebrew-cask.yml" 'an active Homebrew cask publisher' "$tmp_dir/homebrew-cask-output"

cp "$config_under_test" "$tmp_dir/snapcraft.yml"
printf '\nsnapcrafts:\n  - id: forbidden\n' >>"$tmp_dir/snapcraft.yml"
assert_rejected_config "$tmp_dir/snapcraft.yml" 'an active Snap publisher' "$tmp_dir/snapcraft-output"

cat >"$fake_goreleaser" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
{
	printf 'goreleaser'
	printf ' <%s>' "$@"
	printf '\n'
} >>"$FAKE_GORELEASER_CALLS"

for credential in GITHUB_TOKEN GITLAB_TOKEN GITEA_TOKEN HOMEBREW_TAP_GITHUB_TOKEN SNAPCRAFT_STORE_CREDENTIALS; do
	[[ -z ${!credential-} ]] || {
		printf 'recording goreleaser received credential %s\n' "$credential" >&2
		exit 97
	}
done
EOF
chmod +x "$fake_goreleaser"

cat >"$snapshot_scripts/accept-snapshot" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
printf 'GORELEASER=%s\n' "${PLY_SNAPSHOT_GORELEASER-}" >>"$FAKE_SNAPSHOT_ACCEPTANCE_CALLS"
EOF
chmod +x "$snapshot_scripts/accept-snapshot"

cat >"$docker_scripts/accept-docker" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
printf 'accept-docker\n' >>"$FAKE_DOCKER_ACCEPTANCE_CALLS"
EOF
chmod +x "$docker_scripts/accept-docker"

run_distribution_target() {
	local target=$1
	local output=$2
	: >"$calls_file"
	(
		cd "$make_workspace"
		FAKE_GORELEASER_CALLS="$calls_file" \
			GITHUB_TOKEN=github-secret \
			GITLAB_TOKEN=gitlab-secret \
			GITEA_TOKEN=gitea-secret \
			HOMEBREW_TAP_GITHUB_TOKEN=homebrew-secret \
			SNAPCRAFT_STORE_CREDENTIALS=snap-secret \
			make --no-print-directory -f "$makefile_under_test" \
				GORELEASER="$fake_goreleaser" "$target"
	) >"$output" 2>&1
}

expected_call='goreleaser <release> <--snapshot> <--clean> <--skip=publish>'
for target in snapshot release; do
	if ! run_distribution_target "$target" "$tmp_dir/$target-output"; then
		sed -n '1,160p' "$tmp_dir/$target-output" >&2
		fail "$target did not run the local-only distribution path"
	fi
	[[ $(wc -l <"$calls_file" | tr -d '[:space:]') -eq 1 ]] ||
		fail "$target did not invoke exactly one recording executable"
	grep -Fqx "$expected_call" "$calls_file" ||
		fail "$target did not record the exact local-only GoReleaser invocation"
done

: >"$snapshot_calls"
: >"$calls_file"
if ! (
	cd "$make_workspace"
	FAKE_SNAPSHOT_ACCEPTANCE_CALLS="$snapshot_calls" \
		make --no-print-directory -f "$makefile_under_test" \
			BASH=/bin/bash GORELEASER="$fake_goreleaser" \
			SCRIPTS_DIR="$snapshot_scripts" acceptance-snapshot
) >"$tmp_dir/acceptance-snapshot-output" 2>&1; then
	sed -n '1,160p' "$tmp_dir/acceptance-snapshot-output" >&2
	fail 'acceptance-snapshot did not invoke the focused orchestration'
fi
[[ $(wc -l <"$snapshot_calls" | tr -d '[:space:]') -eq 1 ]] ||
	fail 'acceptance-snapshot did not invoke exactly one orchestration'
grep -Fqx "GORELEASER=$fake_goreleaser" "$snapshot_calls" ||
	fail 'acceptance-snapshot did not pass the configured GoReleaser tool'
[[ ! -s "$calls_file" ]] ||
	fail 'acceptance-snapshot bypassed its orchestration and invoked GoReleaser directly'

: >"$docker_calls"
if ! (
	cd "$make_workspace"
	FAKE_DOCKER_ACCEPTANCE_CALLS="$docker_calls" \
		make --no-print-directory -f "$makefile_under_test" \
			BASH=/bin/bash SCRIPTS_DIR="$docker_scripts" acceptance-docker
) >"$tmp_dir/acceptance-docker-output" 2>&1; then
	sed -n '1,160p' "$tmp_dir/acceptance-docker-output" >&2
	fail 'acceptance-docker did not invoke the focused orchestration'
fi
[[ $(wc -l <"$docker_calls" | tr -d '[:space:]') -eq 1 ]] ||
	fail 'acceptance-docker did not invoke exactly one orchestration'
grep -Fqx 'accept-docker' "$docker_calls" ||
	fail 'acceptance-docker invoked the wrong focused orchestration'

if run_distribution_target release-brew "$tmp_dir/release-brew-output"; then
	fail 'release-brew remained callable'
fi
[[ ! -s "$calls_file" ]] || fail 'release-brew reached the recording executable'
grep -F 'Homebrew distribution is inactive' "$tmp_dir/release-brew-output" >/dev/null ||
	fail 'release-brew did not explain its fail-closed boundary'

printf 'make distribution contract: PASS (default config, snapshot argv, credential isolation, snapshot/Docker acceptance wiring, and inactive publishers checked)\n'
