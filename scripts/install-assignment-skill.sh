#!/bin/sh
# Usage: install-assignment-skill.sh <source SKILL.md> <destination directory>
# Copies the repository-owned skill to <destination>/SKILL.md. There is no
# default destination. Symlinks are never followed, nothing is overwritten and
# other files in the destination are not touched. An identical SKILL.md is a
# no-op; anything else already there is existing_install_conflict. A failure
# to write while SKILL.md is still missing is install_io_error.

set -eu

usage() {
	printf '%s\n' 'usage: install-assignment-skill.sh <source SKILL.md> <destination directory>' >&2
	exit 2
}

io_error() {
	printf 'install_io_error: %s\n' "$*" >&2
	exit 1
}

conflict() {
	printf 'existing_install_conflict: %s\n' "$*" >&2
	exit 1
}

[ "$#" -eq 2 ] || usage
source_file=$1
destination=$2
[ -n "$source_file" ] && [ -n "$destination" ] || usage
# A trailing slash would make the symlink checks follow the link.
while [ "$destination" != / ] && [ "${destination%/}" != "$destination" ]; do
	destination=${destination%/}
done
[ -f "$source_file" ] && [ ! -L "$source_file" ] || {
	printf '%s\n' 'install-assignment-skill: source is not a regular file' >&2
	exit 2
}
target="$destination/SKILL.md"

check_destination() {
	if [ -L "$destination" ]; then
		conflict "$destination is a symlink"
	fi
	if [ -e "$destination" ] && [ ! -d "$destination" ]; then
		conflict "$destination is not a directory"
	fi
}

# Exit 0 when target already holds the source bytes; conflict on anything else.
check_existing() {
	if [ -L "$target" ]; then
		conflict "$target is a symlink"
	fi
	if [ -e "$target" ]; then
		if [ -f "$target" ] && cmp -s "$source_file" "$target"; then
			printf 'unchanged %s\n' "$target"
			exit 0
		fi
		conflict "$target exists with other content"
	fi
}

check_destination
mkdir -p "$destination"
check_destination
check_existing

# Write a private temporary copy, then create the target name with POSIX
# link(1): exactly one link(2), which fails when the name already exists as a
# file, directory or symlink and never links into a directory or follows a
# symlink. A concurrent writer is therefore detected, never written into.
temporary="$destination/.SKILL.md.install.$$"
# Create the temporary file exclusively on descriptor 3 (a failed exclusive
# create exits before anything is armed), arm its cleanup at once, then write
# through the descriptor so the name is not looked up again.
umask 022
set -C
# POSIX shells exit on a failed exec redirection; the guard covers bash.
if ! exec 3>"$temporary"; then
	io_error "cannot create a temporary copy in $destination"
fi
set +C
trap 'rm -f "$temporary"' EXIT
if ! cat "$source_file" >&3 || ! exec 3>&-; then
	io_error "cannot write a temporary copy in $destination"
fi
chmod 0644 "$temporary"
if ! link "$temporary" "$target" 2>/dev/null; then
	# Only a name that now exists is a conflict (or an identical no-op).
	if [ -e "$target" ] || [ -L "$target" ]; then
		check_existing
		conflict "$target appeared during installation"
	fi
	io_error "cannot create $target"
fi
printf 'installed %s\n' "$target"
