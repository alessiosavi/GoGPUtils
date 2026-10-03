# Benchmark targets; see docs/benchmarks.md.
SHELL := bash
.DEFAULT_GOAL := help
BENCH_VARS = PKG='$(value PKG)' BENCH='$(value BENCH)' COUNT='$(value COUNT)' BENCHTIME='$(value BENCHTIME)' BASE='$(value BASE)'

.PHONY: help bench bench-compare bench-profile bench-cov

help:
	@printf '%s\n' \
	  'make bench [PKG=pattern] [BENCH=.] [COUNT=1] [BENCHTIME=1s]' \
	  'make bench-compare BASE=ref [PKG=pattern] [BENCH=.] [COUNT=6] [BENCHTIME=500ms]' \
	  'make bench-profile PKG=package [BENCH=.] [BENCHTIME=1s]' \
	  'make bench-cov' \
	  'See docs/benchmarks.md for details; filter with PKG and BENCH first.'

bench:
	@$(BENCH_VARS) scripts/bench.sh run

bench-compare:
	@$(BENCH_VARS) scripts/bench.sh compare

bench-profile:
	@$(BENCH_VARS) scripts/bench.sh profile

bench-cov:
	go test -count=1 -v ./internal/benchcov
