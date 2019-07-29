#!/bin/sh
#
# Build and test Git
#

. ${0%/*}/lib.sh

case "$CI_OS_NAME" in
windows*) cmd //c mklink //j t\\.prove "$(cygpath -aw "$cache_dir/.prove")";;
*) ln -s "$cache_dir/.prove" t/.prove;;
esac

make
case "$jobname" in
linux-gcc)
	make test

	make -C t/trace_schema_validator
	export GIT_TRACE2_EVENT=$(mktemp)
	export GIT_TEST_SPLIT_INDEX=yes
	export GIT_TEST_FULL_IN_PACK_ARRAY=true
	export GIT_TEST_OE_SIZE=10
	export GIT_TEST_OE_DELTA_SIZE=5
	export GIT_TEST_COMMIT_GRAPH=1
	export GIT_TEST_MULTI_PACK_INDEX=1
	make test

	t/trace_schema_validator/trace_schema_validator \
		--trace2_event_file=${GIT_TRACE2_EVENT} \
		--schema_file=t/trace_schema_validator/strict_schema.json \
		--progress=10000
	;;
linux-gcc-4.8)
	# Don't run the tests; we only care about whether Git can be
	# built with GCC 4.8, as it errors out on some undesired (C99)
	# constructs that newer compilers seem to quietly accept.
	;;
*)
	make test
	;;
esac

check_unignored_build_artifacts

save_good_tree
