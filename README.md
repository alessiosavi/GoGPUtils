# GoGPUtils

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/86513a2282374f87a813110db86f018b)](https://www.codacy.com/manual/alessiosavi/GoGPUtils?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=alessiosavi/GoGPUtils&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/alessiosavi/GoGPUtils)](https://goreportcard.com/report/github.com/alessiosavi/GoGPUtils) [![GoDoc](https://godoc.org/github.com/alessiosavi/GoGPUtils?status.svg)](https://godoc.org/github.com/alessiosavi/GoGPUtils) [![License](https://img.shields.io/github/license/alessiosavi/GoGPUtils)](https://img.shields.io/github/license/alessiosavi/GoGPUtils) [![Version](https://img.shields.io/github/v/tag/alessiosavi/GoGPUtils)](https://img.shields.io/github/v/tag/alessiosavi/GoGPUtils) [![Code size](https://img.shields.io/github/languages/code-size/alessiosavi/GoGPUtils)](https://img.shields.io/github/languages/code-size/alessiosavi/GoGPUtils) [![Repo size](https://img.shields.io/github/repo-size/alessiosavi/GoGPUtils)](https://img.shields.io/github/repo-size/alessiosavi/GoGPUtils) [![Issue open](https://img.shields.io/github/issues/alessiosavi/GoGPUtils)](https://img.shields.io/github/issues/alessiosavi/GoGPUtils)
[![Issue closed](https://img.shields.io/github/issues-closed/alessiosavi/GoGPUtils)](https://img.shields.io/github/issues-closed/alessiosavi/GoGPUtils)

A simple Go general purpouse utility repository for avoid to reinvent the wheel every time that i need to start a new project.


# Benchmark

```text
$ go test . ./...  -bench=. -benchmem -benchtime=3s
?       github.com/alessiosavi/GoGPUtils        [no test files]
?       github.com/alessiosavi/GoGPUtils/array  [no test files]
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/files
BenchmarkCountLinesFile-8                 310632             11014 ns/op           37008 B/op          6 allocs/op
BenchmarkListFile-8                         1057           3149317 ns/op         1700843 B/op       3332 allocs/op
BenchmarkFindFilesSensitive-8               1144           3052561 ns/op         1692683 B/op       3324 allocs/op
BenchmarkFindFilesInsensitive-8             1123           3251126 ns/op         1692685 B/op       3324 allocs/op
BenchmarkGetFileSize-8                   2152818              1757 ns/op             272 B/op          2 allocs/op
BenchmarkGetFileSize2-8                   586221              5907 ns/op             384 B/op          5 allocs/op
PASS
ok      github.com/alessiosavi/GoGPUtils/files  24.021s
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/helper
BenchmarkRandomIntn-8                     401257              8917 ns/op               0 B/op          0 allocs/op
BenchmarkRandomInt32-8                    392035              8954 ns/op               0 B/op          0 allocs/op
BenchmarkRandomInt64-8                    383592              8652 ns/op               0 B/op          0 allocs/op
BenchmarkRandomFloat32-8                  391240              8670 ns/op               0 B/op          0 allocs/op
BenchmarkRandomFloat64-8                  382899              8673 ns/op               0 B/op          0 allocs/op
BenchmarkRandomIntnR-8                  330459334               10.8 ns/op             0 B/op          0 allocs/op
BenchmarkRandomInt32R-8                 421643742                8.53 ns/op            0 B/op          0 allocs/op
BenchmarkRandomInt64R-8                 178519230               20.1 ns/op             0 B/op          0 allocs/op
BenchmarkRandomFloat32R-8               558883695                6.29 ns/op            0 B/op          0 allocs/op
BenchmarkRandomFloat64R-8               530065819                6.81 ns/op            0 B/op          0 allocs/op
BenchmarkRandomIntnRArray-8               239304             13421 ns/op            8192 B/op          1 allocs/op
BenchmarkRandomInt32RArray-8              322059             10385 ns/op            4096 B/op          1 allocs/op
BenchmarkRandomInt64RArray-8              156453             21958 ns/op            8192 B/op          1 allocs/op
BenchmarkRandomFloat32Array-8             409972              8077 ns/op            4096 B/op          1 allocs/op
BenchmarkRandomFloat64RArray-8            391137              8722 ns/op            8192 B/op          1 allocs/op
PASS
ok      github.com/alessiosavi/GoGPUtils/helper 58.285s
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/math
BenchmarkSumIntArray-8           9767974               350 ns/op               0 B/op          0 allocs/op
BenchmarkSumInt32Array-8         9989743               355 ns/op               0 B/op          0 allocs/op
BenchmarkSumInt64Array-8         6300518               563 ns/op               0 B/op          0 allocs/op
BenchmarkSumFloat32Array-8       3322027              1074 ns/op               0 B/op          0 allocs/op
BenchmarkSumFloat64Array-8       3347278              1070 ns/op               0 B/op          0 allocs/op
BenchmarkMaxIntIndex-8           4171531               855 ns/op               0 B/op          0 allocs/op
BenchmarkMaxInt32Index-8         4213944               875 ns/op               0 B/op          0 allocs/op
BenchmarkMaxInt64Index-8         4032759               856 ns/op               0 B/op          0 allocs/op
BenchmarkMaxFloat32Index-8       2856822              1138 ns/op               0 B/op          0 allocs/op
BenchmarkMaxFloat64Index-8       3859327               942 ns/op               0 B/op          0 allocs/op
BenchmarkAverageInt-8            9782680               359 ns/op               0 B/op          0 allocs/op
BenchmarkAverageInt32-8          9919570               358 ns/op               0 B/op          0 allocs/op
BenchmarkAverageInt64-8          9758103               368 ns/op               0 B/op          0 allocs/op
BenchmarkAverageFloat32-8        3361938              1071 ns/op               0 B/op          0 allocs/op
BenchmarkAverageFloat64-8        3357955              1061 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/alessiosavi/GoGPUtils/math   64.756s
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/string
BenchmarkTestIsUpperKO-8        18585994               184 ns/op               0 B/op          0 allocs/op
BenchmarkTestIsUpperOK-8        18989983               184 ns/op               0 B/op          0 allocs/op
BenchmarkTestIsLowerOK-8        18955988               185 ns/op               0 B/op          0 allocs/op
BenchmarkTestIsLowerKO-8        18982801               184 ns/op               0 B/op          0 allocs/op
BenchmarkRemoveFromString-8     297683752               11.5 ns/op             0 B/op          0 allocs/op
BenchmarkRandomString-8           162973             21412 ns/op            5376 B/op          1 allocs/op
PASS
ok      github.com/alessiosavi/GoGPUtils/string 23.062s
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/zip
BenchmarkReadZipFile-8            263382             12236 ns/op            7808 B/op         29 allocs/op
BenchmarkReadZip01-8              270822             12522 ns/op            8144 B/op         31 allocs/op
PASS
ok      github.com/alessiosavi/GoGPUtils/zip    6.894s
```