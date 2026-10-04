---
title: Benchmarks
nav_order: 6
---

# Benchmarks

Every exported function or method of the utility packages (everything except `aws/...` and `internal/...`) has a canonical benchmark or a reviewed exemption. A test in `internal/benchcov` enforces this in CI.

## Conventions

- **File and names.** Each package has one `bench_test.go`. Function `Foo` is benchmarked by `BenchmarkFoo`, method `T.M` by `BenchmarkT_M`.
- **Sub-benchmarks** are `key=value` dimensions: `n=` (input size), `type=`, `case=`, `mode=`, `order=`, `shards=`, `key=`, `bins=`, `files=`, `languages=`, `read=`, `words=`. Benchmark names and dimension values contain neither spaces nor commas, because paired-row counting reads benchstat's CSV output.
- **Inputs** come from `internal/benchkit`. They are deterministic and seeded, with standard sizes `16, 1024, 65536`, or `8, 64, 512` for quadratic work.
- **Loops.** Serial benchmarks use `for b.Loop()` with `b.ReportAllocs()`. Parallel ones use `b.RunParallel`.
- **State.** Every iteration measures the same documented state. A benchmark that consumes its input restores it inside the loop, and its doc comment says so.
- **Exemptions** are limited to presentation methods on error and enum types, and to truly transparent delegation whose reason names the delegate. Any conversion, allocation or different option disqualifies a delegation exemption. The cap is `100 × exemptions ≤ 5 × APIs`. They live in `internal/benchcov/exempt.txt`.

## Commands

**Requirements:** Go 1.26.4 or newer, Git with worktree support, Make, and Bash. The comparison runs the pinned benchstat through `go run`; its first use needs the module proxy and checksum service unless the module is already cached.

| Command | What it does |
|---|---|
| `make bench [PKG=./sliceutil] [BENCH=Filter] [COUNT=1] [BENCHTIME=1s]` | Run benchmarks; results in `bench-out/<run>/` |
| `make bench-compare BASE=master [PKG] [BENCH] [COUNT=6] [BENCHTIME=500ms]` | Run BASE and the working tree, then `benchstat` |
| `make bench-profile PKG=./sliceutil BENCH=BenchmarkFilter` | CPU and memory profiles for one package |
| `make bench-cov` | Coverage report from `internal/benchcov` |
| `scripts/bench.sh report RUN_DIR` | Re-print a comparison report from a saved `bench-compare` run directory |

Each run directory holds `results.txt` (raw `go test` output, ready for `benchstat`), `meta.txt` (commit, dirty flag, content and suite fingerprints, Go version, `GOFLAGS`, `CGO_ENABLED`, OS/arch, CPU, exact command) and `packages.txt`.

## Comparing revisions

- **Filter first.** A full unfiltered comparison is estimated to take about three hours at the defaults (500 ms per benchmark, six samples, two revisions); pass `PKG` and `BENCH` for anything routine.
- **Samples.** Six samples (the default) are fine for a quick look. Use `COUNT=10` or more before claiming an improvement.
- **Workloads.** Compare only paired rows whose workload is unchanged. After this suite lands, benchmark names and workloads are frozen; a changed workload gets a new case name.
- **Environment.** `bench-compare` refuses to pair runs whose Go version, `GOFLAGS` or `CGO_ENABLED` differ.
- **Flag experiments.** For deliberate experiments (PGO, `-gcflags`), run `make bench` once per variant with the same `COUNT`, then run `benchstat` on the two `results.txt` files.
- **Don't edit during a run.** A run is marked invalid if the working tree's content changes while it measures.
- **Historical revisions.** Each revision runs its own benchmark files, so a revision from before this suite can have no paired rows. That result is informational, not a speed comparison.
