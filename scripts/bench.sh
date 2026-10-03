#!/usr/bin/env bash
# bench.sh runs, compares, and profiles the repository's benchmarks.
# Use it through the Makefile; docs/benchmarks.md documents every variable.
#
#   run      PKG BENCH COUNT(1) BENCHTIME(1s)
#   compare  BASE PKG BENCH COUNT(6) BENCHTIME(500ms)
#   profile  PKG(one package) BENCH BENCHTIME(1s)
#   report   RUN_DIR (re-run the comparison report of a bench-compare run)
#
# Results files contain only `go test` output, so benchstat pairs runs by
# package, benchmark name and environment. Provenance (commit, dirty flag,
# content and suite fingerprints, toolchain, command) lives in meta.txt.
set -euo pipefail
export LC_ALL=C

readonly BENCHSTAT=golang.org/x/perf/cmd/benchstat@v0.0.0-20260929162123-406019bb8b68

repo=$(git rev-parse --show-toplevel)
cd "$repo"
BENCH_WT=""

die() {
	echo "bench.sh: $*" >&2
	exit 2
}

usage() {
	sed -n '2,14p' "$0" >&2
	exit 2
}

# selected_pkgs DIR prints the in-scope import paths of the checkout at DIR,
# narrowed to the PKG patterns when PKG is set. A requested package that the
# active build cannot load is an error, never dropped silently.
selected_pkgs() {
	local dir=$1 all listing bad want line
	all=$(go run ./internal/benchcov/cmd/benchscope -root "$dir")
	if [[ -z ${PKG:-} ]]; then
		printf '%s\n' "$all"
		return
	fi
	# shellcheck disable=SC2086 # PKG may hold several patterns.
	listing=$(cd "$dir" && go list -e -f '{{.ImportPath}}{{"\t"}}{{if .Error}}{{printf "%q" .Error.Err}}{{end}}' $PKG)
	bad=$(awk -F '\t' '$2 != ""' <<<"$listing")
	if [[ -n $bad ]]; then
		while IFS= read -r line; do
			printf 'bench.sh: package unavailable in %s: %s\n' "$dir" "$line" >&2
		done <<<"$bad"
		return 1
	fi
	want=$(cut -f1 <<<"$listing" | sed '/^$/d' | sort -u)
	comm -13 <(sort <<<"$all") <(printf '%s\n' "$want") | sed '/^$/d; s/^/bench.sh: not in benchmark scope: /' >&2
	comm -12 <(sort <<<"$all") <(printf '%s\n' "$want")
}

# fingerprint DIR [suite] prints the SHA-256 content fingerprint of the
# checkout at DIR, restricted to the suite's files when $2 is "suite".
fingerprint() {
	local -a args=(-dir "$1")
	if [[ ${2:-} == suite ]]; then
		args+=(-suite)
	fi
	go run ./internal/benchcov/cmd/benchfingerprint "${args[@]}"
}

cpu_model() {
	sysctl -n machdep.cpu.brand_string 2>/dev/null ||
		sed -n 's/^model name[[:space:]]*: //p' /proc/cpuinfo 2>/dev/null | head -n 1 ||
		echo unknown
}

# write_meta DIR FILE CMD records provenance and environment for one run.
write_meta() {
	local dir=$1 file=$2 cmd=$3
	{
		echo "commit: $(git -C "$dir" rev-parse HEAD)"
		if [[ -n $(git -C "$dir" status --porcelain) ]]; then echo "dirty: true"; else echo "dirty: false"; fi
		echo "content-fingerprint: $(fingerprint "$dir")"
		echo "suite-fingerprint: $(fingerprint "$dir" suite)"
		paste -d ' ' <(printf '%s\n' goversion: goflags: cgo_enabled: goos: goarch:) \
			<(cd "$dir" && go env GOVERSION GOFLAGS CGO_ENABLED GOOS GOARCH)
		echo "gomaxprocs: ${GOMAXPROCS:-$(getconf _NPROCESSORS_ONLN)}"
		echo "cpu: $(cpu_model)"
		echo "command: $cmd"
	} >"$file"
}

new_run_dir() {
	local d
	d="$repo/bench-out/$(date +%Y%m%dT%H%M%S)-$$-$RANDOM"
	mkdir -p "$d"
	printf '%s\n' "$d"
}

# has_measurements FILE succeeds when FILE holds at least one measurement: a
# benchmark name, a positive iteration count, and complete value/unit pairs.
# Skipped, truncated, and zero-iteration lines are not measurements.
has_measurements() {
	awk '
		$1 !~ /^Benchmark[^[:space:]]*$/ || $2 !~ /^[0-9]+$/ || $2 + 0 <= 0 || NF < 4 || NF % 2 { next }
		{
			valid = 1
			for (i = 3; i <= NF; i += 2) {
				if ($i !~ /^[+-]?([0-9]+([.][0-9]*)?|[.][0-9]+)([eE][+-]?[0-9]+)?$/) valid = 0
			}
			if (valid) found = 1
		}
		END { exit !found }' "$1"
}

benchstat_run() { go run "$BENCHSTAT" "$@"; }

# run_bench DIR OUT runs the selected benchmarks of the checkout at DIR and
# writes OUT/results.txt (pure go test output), OUT/meta.txt and
# OUT/packages.txt.
run_bench() {
	local dir=$1 out=$2 pkgs
	mkdir -p "$out"
	pkgs=$(selected_pkgs "$dir")
	printf '%s\n' "$pkgs" | sed '/^$/d' >"$out/packages.txt"
	if [[ -z $pkgs ]]; then
		echo "bench.sh: no in-scope package selected in $dir" >&2
		: >"$out/results.txt"
		write_meta "$dir" "$out/meta.txt" "(nothing selected)"
		return 0
	fi
	local -a cmd=(go test -run '^$' -bench "$BENCH" -benchmem -count "$COUNT" -benchtime "$BENCHTIME")
	local p
	while IFS= read -r p; do cmd+=("$p"); done <"$out/packages.txt"
	write_meta "$dir" "$out/meta.txt" "${cmd[*]}"
	(cd "$dir" && "${cmd[@]}") | tee "$out/results.txt"
}

# check_unchanged DIR OUT appends "valid: true" to OUT/meta.txt, or marks the
# run invalid and fails when the content of DIR changed during the run.
check_unchanged() {
	local dir=$1 out=$2 before after
	before=$(sed -n 's/^content-fingerprint: //p' "$out/meta.txt")
	after=$(fingerprint "$dir")
	if [[ $before != "$after" ]]; then
		echo "valid: false" >>"$out/meta.txt"
		echo "bench.sh: sources changed during the run; $out is invalid for comparison" >&2
		exit 1
	fi
	echo "valid: true" >>"$out/meta.txt"
}

cmd_run() {
	local out
	out=$(new_run_dir)
	run_bench "$repo" "$out"
	has_measurements "$out/results.txt" || die "no benchmark matched (PKG=${PKG:-} BENCH=$BENCH)"
	check_unchanged "$repo" "$out"
	echo "results: $out/results.txt"
	echo "meta:    $out/meta.txt"
}

cleanup_worktree() {
	local status=$?
	trap - EXIT INT TERM
	if [[ -n $BENCH_WT && -e $BENCH_WT ]]; then
		if ! git -C "$repo" worktree remove --force "$BENCH_WT" >/dev/null 2>&1; then
			echo "bench.sh: could not remove $BENCH_WT; run: git worktree remove --force $BENCH_WT" >&2
		fi
	fi
	exit "$status"
}

# validate_side DIR LABEL fails unless DIR holds a complete, valid run: all
# three files, exactly one validity value equal to true, and the environment
# keys.
validate_side() {
	local dir=$1 label=$2 key validity
	[[ -f $dir/meta.txt && -f $dir/results.txt && -f $dir/packages.txt ]] || die "$label: incomplete run in $dir"
	validity=$(sed -n 's/^valid: //p' "$dir/meta.txt")
	[[ $validity == true ]] || die "$label: run is not valid (see $dir/meta.txt)"
	for key in goversion goflags cgo_enabled; do
		grep -q "^$key:" "$dir/meta.txt" || die "$label: $dir/meta.txt lacks $key"
	done
}

# pair_counts reads the pinned benchstat's CSV output, produced with the file
# labels BASE and HEAD, and prints "paired base_only head_only" metric rows.
# Each table's revision columns come from its own file-label header (a
# one-revision table holds only BASE or only HEAD); the mapping resets between
# tables. Configuration, header, and geomean rows are not counted. Benchmark
# names and dimension values never contain commas (see conventions).
pair_counts() {
	awk -F, '
		NF == 0 { bi = hi = 0; next }
		$1 == "" && ($2 == "BASE" || $2 == "HEAD") {
			bi = hi = 0
			for (i = 2; i <= NF; i++) {
				if ($i == "BASE") bi = i
				if ($i == "HEAD") hi = i
			}
			next
		}
		$1 == "" || $1 == "geomean" || (!bi && !hi) { next }
		{ b = (bi && $bi != ""); h = (hi && $hi != "") }
		b && h { p++ }
		b && !h { bo++ }
		!b && h { ho++ }
		END { printf "%d %d %d\n", p, bo, ho }'
}

compare_results() {
	local base=$1/base head=$1/head key a b mismatch=0
	validate_side "$base" base
	validate_side "$head" head
	if ! has_measurements "$base/results.txt" && ! has_measurements "$head/results.txt"; then
		die "no benchmark matched on either side (PKG=${PKG:-} BENCH=$BENCH)"
	fi
	echo "== base meta"
	cat "$base/meta.txt"
	echo "== head meta"
	cat "$head/meta.txt"
	comm -23 "$base/packages.txt" "$head/packages.txt" | sed 's/^/base only package: /'
	comm -13 "$base/packages.txt" "$head/packages.txt" | sed 's/^/head only package: /'

	for key in goversion goflags cgo_enabled; do
		a=$(sed -n "s/^$key: //p" "$base/meta.txt")
		b=$(sed -n "s/^$key: //p" "$head/meta.txt")
		if [[ $a != "$b" ]]; then
			echo "environment differs: $key base='$a' head='$b'"
			mismatch=1
		fi
	done
	if ((mismatch)); then
		if has_measurements "$base/results.txt"; then benchstat_run "$base/results.txt"; fi
		if has_measurements "$head/results.txt"; then benchstat_run "$head/results.txt"; fi
		echo "paired rows: 0 - no performance comparison possible (environment differs)"
		return 0
	fi

	local csv counts paired base_only head_only
	benchstat_run "$base/results.txt" "$head/results.txt"
	csv=$(benchstat_run -format csv "BASE=$base/results.txt" "HEAD=$head/results.txt" 2>/dev/null)
	counts=$(pair_counts <<<"$csv")
	read -r paired base_only head_only <<<"$counts"
	echo "paired rows: $paired; base only: $base_only; head only: $head_only"
	if ((paired == 0)); then
		echo "no performance comparison possible (no paired rows)"
	elif [[ $(sed -n 's/^suite-fingerprint: //p' "$base/meta.txt") != "$(sed -n 's/^suite-fingerprint: //p' "$head/meta.txt")" ]]; then
		echo "warning: the benchmark suite differs between revisions; paired rows may measure different workloads"
	fi
}

cmd_compare() {
	[[ -n ${BASE:-} ]] || die "BASE is required, e.g. make bench-compare BASE=master"
	local base_sha run
	base_sha=$(git rev-parse --verify --quiet "${BASE}^{commit}") || die "unknown revision: $BASE"
	run=$(new_run_dir)
	BENCH_WT="$repo/.worktrees/bench-base-${base_sha:0:7}-$$-$RANDOM"
	trap cleanup_worktree EXIT
	trap 'exit 130' INT
	trap 'exit 143' TERM
	git worktree add --detach --quiet "$BENCH_WT" "$base_sha"

	echo "== base: $BASE ($base_sha)" >&2
	run_bench "$BENCH_WT" "$run/base"
	check_unchanged "$BENCH_WT" "$run/base"
	echo "== head: working tree" >&2
	run_bench "$repo" "$run/head"
	check_unchanged "$repo" "$run/head"
	compare_results "$run"
	echo "results: $run"
}

cmd_profile() {
	[[ -n ${PKG:-} ]] || die "PKG is required, e.g. make bench-profile PKG=./sliceutil BENCH=BenchmarkFilter"
	local pkgs out
	pkgs=$(selected_pkgs "$repo")
	if [[ -z $pkgs || $(wc -l <<<"$pkgs") -ne 1 ]]; then
		die "PKG must select exactly one in-scope package; selected: $(tr '\n' ' ' <<<"$pkgs")"
	fi
	out=$(new_run_dir)
	printf '%s\n' "$pkgs" >"$out/packages.txt"
	local -a cmd=(go test -run '^$' -bench "$BENCH" -benchmem -count 1 -benchtime "$BENCHTIME"
		-cpuprofile "$out/cpu.out" -memprofile "$out/mem.out" -o "$out/pkg.test" "$pkgs")
	write_meta "$repo" "$out/meta.txt" "${cmd[*]}"
	"${cmd[@]}" | tee "$out/results.txt"
	has_measurements "$out/results.txt" || die "no benchmark matched BENCH=$BENCH"
	check_unchanged "$repo" "$out"
	echo "inspect: go tool pprof $out/pkg.test $out/cpu.out"
}

sub=${1:-}
BENCH=${BENCH:-.}
case $sub in
run)
	COUNT=${COUNT:-1}
	BENCHTIME=${BENCHTIME:-1s}
	cmd_run
	;;
compare)
	COUNT=${COUNT:-6}
	BENCHTIME=${BENCHTIME:-500ms}
	cmd_compare
	;;
profile)
	BENCHTIME=${BENCHTIME:-1s}
	cmd_profile
	;;
report)
	[[ -n ${2:-} && -d ${2:-} ]] || die "usage: scripts/bench.sh report RUN_DIR (a bench-compare run directory)"
	compare_results "$2"
	;;
*) usage ;;
esac
