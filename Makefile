# Benchmark targets; see docs/benchmarks.md.
SHELL := bash
BENCH_VARS = PKG='$(PKG)' BENCH='$(BENCH)' COUNT='$(COUNT)' BENCHTIME='$(BENCHTIME)' BASE='$(BASE)'

.PHONY: bench bench-compare bench-profile bench-cov

bench:
	@$(BENCH_VARS) scripts/bench.sh run

bench-compare:
	@$(BENCH_VARS) scripts/bench.sh compare

bench-profile:
	@$(BENCH_VARS) scripts/bench.sh profile

bench-cov:
	go test -count=1 -v ./internal/benchcov
