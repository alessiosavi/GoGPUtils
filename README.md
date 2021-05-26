# GoGPUtils

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/86513a2282374f87a813110db86f018b)](https://www.codacy.com/manual/alessiosavi/GoGPUtils?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=alessiosavi/GoGPUtils&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/alessiosavi/GoGPUtils)](https://goreportcard.com/report/github.com/alessiosavi/GoGPUtils) [![GoDoc](https://godoc.org/github.com/alessiosavi/GoGPUtils?status.svg)](https://godoc.org/github.com/alessiosavi/GoGPUtils) [![License](https://img.shields.io/github/license/alessiosavi/GoGPUtils)](https://img.shields.io/github/license/alessiosavi/GoGPUtils) [![Version](https://img.shields.io/github/v/tag/alessiosavi/GoGPUtils)](https://img.shields.io/github/v/tag/alessiosavi/GoGPUtils) [![Code size](https://img.shields.io/github/languages/code-size/alessiosavi/GoGPUtils)](https://img.shields.io/github/languages/code-size/alessiosavi/GoGPUtils) [![Repo size](https://img.shields.io/github/repo-size/alessiosavi/GoGPUtils)](https://img.shields.io/github/repo-size/alessiosavi/GoGPUtils) [![Issue open](https://img.shields.io/github/issues/alessiosavi/GoGPUtils)](https://img.shields.io/github/issues/alessiosavi/GoGPUtils)
[![Issue closed](https://img.shields.io/github/issues-closed/alessiosavi/GoGPUtils)](https://img.shields.io/github/issues-closed/alessiosavi/GoGPUtils)

A simple Go general-purpose utility repository for avoid reinventing the wheel every time that I need to start a new
project.

## Benchmark

```text
$ go test -bench=. -benchmem -benchtime=5s `go list ./... | grep -v aws`
?       github.com/alessiosavi/GoGPUtils        [no test files]
PASS
ok      github.com/alessiosavi/GoGPUtils/array  0.058s
goos: windows
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/byte
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkTestIsUpperByteOK-8    1000000000               3.058 ns/op           0 B/op          0 allocs/op
BenchmarkTestIsLowerByteKO-8    1000000000               2.492 ns/op           0 B/op          0 allocs/op
PASS
ok      github.com/alessiosavi/GoGPUtils/byte   6.199s
PASS
ok      github.com/alessiosavi/GoGPUtils/datastructure/binarytree       0.028s
?       github.com/alessiosavi/GoGPUtils/datastructure/stack    [no test files]
goos: windows
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/files
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkCountLinesFile-8                 131373             91284 ns/op           37576 B/op          5 allocs/op
BenchmarkListFile-8                          880          13606368 ns/op          706878 B/op       5906 allocs/op
BenchmarkFindFilesSensitive-8                776          13658224 ns/op          698715 B/op       5898 allocs/op
BenchmarkFindFilesInsensitive-8              864          15583262 ns/op          698716 B/op       5898 allocs/op
PASS
ok      github.com/alessiosavi/GoGPUtils/files  53.376s
goos: windows
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/helper
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkRandomIntn-8                    1710613              6795 ns/op               0 B/op          0 allocs/op
BenchmarkRandomInt32-8                   1752652              6734 ns/op               0 B/op          0 allocs/op
BenchmarkRandomInt64-8                   1781472              6791 ns/op               0 B/op          0 allocs/op
BenchmarkRandomFloat32-8                 1778323              6732 ns/op               0 B/op          0 allocs/op
BenchmarkRandomFloat64-8                 1767133              6795 ns/op               0 B/op          0 allocs/op
BenchmarkRandomIntnR-8                  1000000000               8.555 ns/op           0 B/op          0 allocs/op
BenchmarkRandomInt32R-8                 1000000000               6.796 ns/op           0 B/op          0 allocs/op
BenchmarkRandomInt64R-8                 772450699               15.48 ns/op            0 B/op          0 allocs/op
BenchmarkRandomFloat32R-8               1000000000               5.457 ns/op           0 B/op          0 allocs/op
BenchmarkRandomFloat64R-8               1000000000               5.402 ns/op           0 B/op          0 allocs/op
BenchmarkRandomIntnRArray-8              1000000             10747 ns/op            8192 B/op          1 allocs/op
BenchmarkRandomInt32RArray-8             1441555              8346 ns/op            4096 B/op          1 allocs/op
BenchmarkRandomInt64RArray-8              688870             17605 ns/op            8192 B/op          1 allocs/op
BenchmarkRandomFloat32Array-8            1755104              7007 ns/op            4096 B/op          1 allocs/op
BenchmarkRandomFloat64RArray-8           1561810              7646 ns/op            8192 B/op          1 allocs/op
PASS
ok      github.com/alessiosavi/GoGPUtils/helper 219.105s
2021/02/26 16:34:58 Request data ->
 GET / HTTP/1.1
Host: localhost:38281
Accept-Encoding: gzip
User-Agent: Go-http-client/1.1


PASS
ok      github.com/alessiosavi/GoGPUtils/http   3.716s
goos: windows
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/math
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkSumIntArray-8                          34900125               338.8 ns/op             0 B/op          0 allocs/op
BenchmarkSumInt32Array-8                        34916972               338.9 ns/op             0 B/op          0 allocs/op
BenchmarkSumInt64Array-8                        34876167               338.9 ns/op             0 B/op          0 allocs/op
BenchmarkSumFloat32Array-8                      13985671               870.8 ns/op             0 B/op          0 allocs/op
BenchmarkSumFloat64Array-8                      13976918               859.1 ns/op             0 B/op          0 allocs/op
BenchmarkMaxIntIndex-8                          14721426               889.0 ns/op             0 B/op          0 allocs/op
BenchmarkMaxInt32Index-8                        14443495               878.8 ns/op             0 B/op          0 allocs/op
BenchmarkMaxInt64Index-8                        13456785               901.6 ns/op             0 B/op          0 allocs/op
BenchmarkMaxFloat32Index-8                      13920338               852.4 ns/op             0 B/op          0 allocs/op
BenchmarkMaxFloat64Index-8                      14205249               878.6 ns/op             0 B/op          0 allocs/op
BenchmarkAverageInt-8                           42384150               283.8 ns/op             0 B/op          0 allocs/op
BenchmarkAverageInt32-8                         35384409               340.5 ns/op             0 B/op          0 allocs/op
BenchmarkAverageInt64-8                         34569390               351.3 ns/op             0 B/op          0 allocs/op
BenchmarkAverageFloat32-8                       14001175               858.9 ns/op             0 B/op          0 allocs/op
BenchmarkAverageFloat64-8                       13993460               858.3 ns/op             0 B/op          0 allocs/op
BenchmarkInitRandomMatrix-8                      1333974              8958 ns/op            6352 B/op         13 allocs/op
BenchmarkMultiplySumArray1000-8                  5363206              2146 ns/op            8192 B/op          1 allocs/op
BenchmarkMultiplyMatrixLegacy100x100-8              2569           4702579 ns/op        18012327 B/op      20101 allocs/op
BenchmarkMultiplyMatrix100x100-8                    8821           1388659 ns/op           92288 B/op        101 allocs/op
BenchmarkIsPrime-8                               4391913              2754 ns/op               0 B/op          0 allocs/op
BenchmarkCosineSimilarity-8                     100000000              111.1 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/alessiosavi/GoGPUtils/math   277.852s
goos: windows
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/search
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkLinearSearchInt-8                796713             14862 ns/op               0 B/op          0 allocs/op
BenchmarkLinearSearchParallelInt-8      473144830               25.39 ns/op           16 B/op          1 allocs/op
BenchmarkContainsStringByte-8            1278368              9348 ns/op           22144 B/op         15 allocs/op
BenchmarkContainsStringsByte-8            576740             20916 ns/op           58193 B/op         33 allocs/op
BenchmarkContainsWhichStrings-8            13597            882602 ns/op           58298 B/op         36 allocs/op
PASS
ok      github.com/alessiosavi/GoGPUtils/search 81.237s
goos: windows
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/string
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkLevenshteinDistanceLegacy-8         147          76146592 ns/op        51286032 B/op       2502 allocs/op
BenchmarkLevenshteinDistance-8               970          12378208 ns/op           40960 B/op          2 allocs/op
BenchmarkDiceCoefficient-8                239919             49929 ns/op               0 B/op          0 allocs/op
BenchmarkJaroDistance-8                    10000           1055681 ns/op            5376 B/op          2 allocs/op
BenchmarkContainsOnlyLetter-8           712269397               16.94 ns/op            0 B/op          0 allocs/op
BenchmarkRemoveFromString-8               108150            107150 ns/op         1163271 B/op          2 allocs/op
BenchmarkRandomString-8                   761596             15854 ns/op            5376 B/op          1 allocs/op
BenchmarkExtractTextFromQuery-8              242          49268281 ns/op        21767194 B/op     102069 allocs/op
BenchmarkCheckPresence-8                1000000000               2.675 ns/op           0 B/op          0 allocs/op
BenchmarkIsUpper-8                      1000000000               3.530 ns/op           0 B/op          0 allocs/op
BenchmarkIsLower-8                      1000000000               3.086 ns/op           0 B/op          0 allocs/op
BenchmarkRemoveWhiteSpace-8                 5170           2288161 ns/op          581633 B/op          1 allocs/op
BenchmarkIsASCII-8                      460007290               26.04 ns/op            0 B/op          0 allocs/op
BenchmarkSplit-8                            9020           1394357 ns/op         2446877 B/op      14363 allocs/op
BenchmarkSplitBuiltin-8                    40594            293564 ns/op          319488 B/op          1 allocs/op
BenchmarkExtractString-8                   29148            413375 ns/op               0 B/op          0 allocs/op
BenchmarkRemoveNonASCII-8                   2469           4889266 ns/op         1662981 B/op          3 allocs/op
BenchmarkTestIsUpperOK-8                1000000000               3.532 ns/op           0 B/op          0 allocs/op
BenchmarkTestIsLowerOK-8                1000000000               3.096 ns/op           0 B/op          0 allocs/op
BenchmarkCreateJSON-8                         30         367769360 ns/op        3096448136 B/op    10052 allocs/op
BenchmarkJoin-8                            22886            534203 ns/op         2931697 B/op         30 allocs/op
BenchmarkTrim-8                             4222           2883190 ns/op         1745819 B/op         13 allocs/op
BenchmarkRemoveDoubleWhiteSpace-8           4723           2533971 ns/op          581633 B/op          1 allocs/op
BenchmarkCountLines-8                      49856            243237 ns/op            4128 B/op          2 allocs/op
BenchmarkReverseString-8                   10000           1087407 ns/op         2914313 B/op         32 allocs/op
PASS
ok      github.com/alessiosavi/GoGPUtils/string 291.349s
goos: windows
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/zip
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkReadZipFile-8            273814             44309 ns/op            6784 B/op         27 allocs/op
BenchmarkReadZip01-8              263980             44548 ns/op            7120 B/op         29 allocs/op
PASS
ok      github.com/alessiosavi/GoGPUtils/zip    24.834s
```
