mutation_meta_source() {
	local source=$1 harness=$2 fixture=$3
	mkdir "$fixture"
	git -C "$source" archive HEAD | tar -xf - -C "$fixture"
	cp "$harness" "$fixture/scripts/"
	git -C "$fixture" init -q --template=
	git -C "$fixture" add -f -- .
	git -C "$fixture" -c user.name='Mutation meta-test' -c user.email=mutation@example.invalid \
		-c core.hooksPath=/dev/null -c commit.gpgSign=false commit -qm 'Mutation meta-test fixture'

	mkdir "$fixture/.build" "$fixture/.git/info"
	printf '%s\n' ignored >"$fixture/.build/mutation-meta-sentinel"
	printf '/mutation-meta-link\n' >"$fixture/.git/info/exclude"
	ln -s .build/mutation-meta-sentinel "$fixture/mutation-meta-link"
	git -C "$fixture" check-ignore -q .build/mutation-meta-sentinel
	git -C "$fixture" check-ignore -q mutation-meta-link
}

mutation_meta_check_copies() {
	local work=$1 copy
	for copy in "$work/control/repo" "$work"/mutant-*/repo; do
		[[ -d $copy ]] || fail "T9 expected an archived source copy at $copy"
		[[ ! -e $copy/.build/mutation-meta-sentinel && ! -L $copy/.build/mutation-meta-sentinel ]] ||
			fail "T9 archived source copy contains the ignored file: $copy"
		[[ ! -e $copy/mutation-meta-link && ! -L $copy/mutation-meta-link ]] ||
			fail "T9 archived source copy contains the ignored symlink: $copy"
	done
}
