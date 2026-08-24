#!/usr/bin/env bash

acceptance_fail() {
	printf 'acceptance support: %s\n' "$*" >&2
	return 1
}

acceptance_build_artifact() {
	local repo_root=$1
	local work_root=$2
	local output=$3
	local go_bin=${GO:-go}
	mkdir -p "$(dirname "$output")" "$work_root/gocache"
	(
		cd "$repo_root"
		GOCACHE="$work_root/gocache" GOENV=off GOWORK=off \
			"$go_bin" build -o "$output" ./cmd/ply
	)
}

acceptance_start_server() {
	local repo_root=$1
	local work_root=$2
	local release=$3
	local port_file="$work_root/loopback-url"
	python3 "$repo_root/test/acceptance/loopback_repository.py" \
		"$port_file" "$release" >"$work_root/loopback.stdout" \
		2>"$work_root/loopback.stderr" &
	ACCEPTANCE_SERVER_PID=$!
	local attempt=0
	while [[ ! -s "$port_file" && "$attempt" -lt 100 ]]; do
		if ! kill -0 "$ACCEPTANCE_SERVER_PID" 2>/dev/null; then
			sed -n '1,80p' "$work_root/loopback.stderr" >&2
			return 1
		fi
		sleep 0.05
		attempt=$((attempt + 1))
	done
	[[ -s "$port_file" ]] || return 1
	IFS= read -r ACCEPTANCE_SERVER_URL <"$port_file"
	export ACCEPTANCE_SERVER_PID ACCEPTANCE_SERVER_URL
}

acceptance_stop_server() {
	if [[ -n "${ACCEPTANCE_SERVER_PID:-}" ]]; then
		kill "$ACCEPTANCE_SERVER_PID" 2>/dev/null || true
		wait "$ACCEPTANCE_SERVER_PID" 2>/dev/null || true
		ACCEPTANCE_SERVER_PID=''
	fi
}

acceptance_prepare_home() {
	local home=$1
	local repository_url=$2
	local profile="$home/.ply/profiles/default"
	mkdir -p "$profile/cloud-config"
	printf 'default\n' >"$home/.ply/profiles/.active_profile"
	printf 'cloudConfig:\n  git:\n    url: ""\nnexus:\n  url: "%s"\n' \
		"$repository_url" >"$profile/local-config.yaml"
	printf '%s\n' \
		'{"type":"deprecated","data":{"dependencies":[]}}' \
		>"$profile/cloud-config/deprecated.json"
	printf '%s\n' \
		'{"type":"defaults","settings":{"disableDependencySort":false,"disableSpringBootUpgrade":false,"disableKotlinUpgrade":false}}' \
		>"$profile/cloud-config/project-defaults.json"
}

acceptance_tree_digest() {
	local root=$1
	(
		cd "$root"
		find . -type f -print | LC_ALL=C sort | while IFS= read -r path; do
			printf '%s\0' "$path"
			shasum -a 256 "$path" | awk '{ print $1 }'
		done
	) | shasum -a 256 | awk '{ print $1 }'
}

acceptance_copy_order_fixture() {
	local repo_root=$1
	local destination=$2
	mkdir -p "$destination"
	cp -R "$repo_root/test/order/." "$destination/"
}
