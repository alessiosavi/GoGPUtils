# GoGPUtils

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/86513a2282374f87a813110db86f018b)](https://www.codacy.com/manual/alessiosavi/GoGPUtils?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=alessiosavi/GoGPUtils&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/alessiosavi/GoGPUtils)](https://goreportcard.com/report/github.com/alessiosavi/GoGPUtils) [![GoDoc](https://godoc.org/github.com/alessiosavi/GoGPUtils?status.svg)](https://godoc.org/github.com/alessiosavi/GoGPUtils) [![License](https://img.shields.io/github/license/alessiosavi/GoGPUtils)](https://img.shields.io/github/license/alessiosavi/GoGPUtils) [![Version](https://img.shields.io/github/v/tag/alessiosavi/GoGPUtils)](https://img.shields.io/github/v/tag/alessiosavi/GoGPUtils) [![Code size](https://img.shields.io/github/languages/code-size/alessiosavi/GoGPUtils)](https://img.shields.io/github/languages/code-size/alessiosavi/GoGPUtils) [![Repo size](https://img.shields.io/github/repo-size/alessiosavi/GoGPUtils)](https://img.shields.io/github/repo-size/alessiosavi/GoGPUtils) [![Issue open](https://img.shields.io/github/issues/alessiosavi/GoGPUtils)](https://img.shields.io/github/issues/alessiosavi/GoGPUtils)
[![Issue closed](https://img.shields.io/github/issues-closed/alessiosavi/GoGPUtils)](https://img.shields.io/github/issues-closed/alessiosavi/GoGPUtils)

A simple Go general-purpose utility repository for avoid reinventing the wheel every time that I need to start a new project.

## Benchmark

```text
$ go test -bench=. -benchmem ./... -benchtime=10s
?   	github.com/alessiosavi/GoGPUtils	[no test files]
PASS
ok  	github.com/alessiosavi/GoGPUtils/byte	7.698s
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/files
BenchmarkCountLinesFile-8         	 1000000	     10585 ns/op	   37008 B/op	       6 allocs/op
BenchmarkListFile-8               	    6351	   1893321 ns/op	 1123163 B/op	    2385 allocs/op
BenchmarkFindFilesSensitive-8     	    6394	   1891287 ns/op	 1115003 B/op	    2377 allocs/op
BenchmarkFindFilesInsensitive-8   	    6086	   1878835 ns/op	 1115003 B/op	    2377 allocs/op
PASS
ok  	github.com/alessiosavi/GoGPUtils/files	52.247s
PASS
ok  	github.com/alessiosavi/GoGPUtils/goutils	0.005s
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/helper
BenchmarkRandomIntn-8                       	 1436161	      8319 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomInt32-8                      	 1438686	      8366 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomInt64-8                      	 1430132	      8397 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomFloat32-8                    	 1435780	      8368 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomFloat64-8                    	 1441004	      8365 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomInt64R-8                     	619096278	        19.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomIntnRArray-8                 	  883638	     13384 ns/op	    8192 B/op	       1 allocs/op
BenchmarkRandomInt32RArray-8                	 1000000	     10587 ns/op	    4096 B/op	       1 allocs/op
BenchmarkRandomInt64RArray-8                	  554470	     21839 ns/op	    8192 B/op	       1 allocs/op
BenchmarkRandomFloat32Array-8               	 1464657	      8212 ns/op	    4096 B/op	       1 allocs/op
BenchmarkRandomFloat64RArray-8              	 1363227	      8763 ns/op	    8192 B/op	       1 allocs/op
PASS
ok  	github.com/alessiosavi/GoGPUtils/helper	231.464s
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/math
BenchmarkSumIntArray-8             	34685700	       355 ns/op	       0 B/op	       0 allocs/op
BenchmarkSumInt32Array-8           	34214290	       347 ns/op	       0 B/op	       0 allocs/op
BenchmarkSumInt64Array-8           	34907217	       347 ns/op	       0 B/op	       0 allocs/op
BenchmarkSumFloat32Array-8         	11595754	      1037 ns/op	       0 B/op	       0 allocs/op
BenchmarkSumFloat64Array-8         	11718639	      1040 ns/op	       0 B/op	       0 allocs/op
BenchmarkMaxIntIndex-8             	10471840	      1155 ns/op	       0 B/op	       0 allocs/op
BenchmarkMaxInt32Index-8           	11901468	       959 ns/op	       0 B/op	       0 allocs/op
BenchmarkMaxInt64Index-8           	12267844	       999 ns/op	       0 B/op	       0 allocs/op
BenchmarkMaxFloat32Index-8         	11957236	       920 ns/op	       0 B/op	       0 allocs/op
BenchmarkMaxFloat64Index-8         	10477437	      1209 ns/op	       0 B/op	       0 allocs/op
BenchmarkAverageInt-8              	34238516	       351 ns/op	       0 B/op	       0 allocs/op
BenchmarkAverageInt32-8            	22153039	       545 ns/op	       0 B/op	       0 allocs/op
BenchmarkAverageInt64-8            	33419302	       353 ns/op	       0 B/op	       0 allocs/op
BenchmarkAverageFloat32-8          	11351874	      1042 ns/op	       0 B/op	       0 allocs/op
BenchmarkAverageFloat64-8          	11446464	      1040 ns/op	       0 B/op	       0 allocs/op
BenchmarkInitRandomMatrix-8        	 1000000	     10555 ns/op	    6352 B/op	      13 allocs/op
BenchmarkMultiplySumArray1000-8    	 6499338	      1761 ns/op	    8192 B/op	       1 allocs/op
BenchmarkMultiplyMatrix100x100-8   	    2817	   4258960 ns/op	18101931 B/op	   20201 allocs/op
BenchmarkIsPrime-8                 	 3234999	      3714 ns/op	       0 B/op	       0 allocs/op
BenchmarkCosineSimilarity-8        	78645279	       155 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/alessiosavi/GoGPUtils/math	260.188s
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/search
BenchmarkLinearSearchInt-8                 	  626892	     17645 ns/op	       0 B/op	       0 allocs/op
BenchmarkLinearSearchParallelInt-8         	832372522	        14.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsStringByte-8              	 1000000	     12689 ns/op	   22152 B/op	      15 allocs/op
BenchmarkContainsStringsByte-8             	  424674	     25588 ns/op	   58193 B/op	      33 allocs/op
BenchmarkContainsWhichStrings-8            	   10000	   1055418 ns/op	   58312 B/op	      36 allocs/op
PASS
ok  	github.com/alessiosavi/GoGPUtils/search	60.964s
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/string
BenchmarkLevenshteinDistanceLegacy-8   	 1459297	      8234 ns/op	   14464 B/op	      43 allocs/op
BenchmarkLevenshteinDistance-8         	 3231523	      3701 ns/op	     704 B/op	       2 allocs/op
BenchmarkContainsOnlyLetter-8          	566363220	        21.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkRemoveFromString-8            	  142773	     83081 ns/op	 1114116 B/op	       2 allocs/op
BenchmarkRandomString-8                	  580832	     20787 ns/op	    5376 B/op	       1 allocs/op
BenchmarkExtractTextFromQuery-8        	     194	  61931772 ns/op	21760005 B/op	  102049 allocs/op
BenchmarkRemoveWhiteSpace-8            	    4521	   2668036 ns/op	  557057 B/op	       1 allocs/op
BenchmarkIsASCII-8                     	588510955	        20.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkSplit-8                       	    8779	   1387231 ns/op	 2447637 B/op	   14363 allocs/op
BenchmarkSplitBuiltin-8                	   27868	    471431 ns/op	  319488 B/op	       1 allocs/op
BenchmarkExtractString-8               	   24660	    488150 ns/op	       0 B/op	       0 allocs/op
BenchmarkRemoveNonASCII-8              	    2176	   5486612 ns/op	 1589252 B/op	       3 allocs/op
BenchmarkCreateJSON-8                  	      39	 288444418 ns/op	2998417935 B/op	    9913 allocs/op
BenchmarkJoin-8                        	   22789	    526067 ns/op	 2930657 B/op	      29 allocs/op
BenchmarkTrim-8                        	    6651	   1786511 ns/op	  557064 B/op	       1 allocs/op
BenchmarkRemoveDoubleWhiteSpace-8      	    4118	   2911370 ns/op	  557057 B/op	       1 allocs/op
BenchmarkCountLines-8                  	   40255	    303206 ns/op	    4128 B/op	       2 allocs/op
BenchmarkReverseString-8               	   10000	   1144097 ns/op	 2914327 B/op	      32 allocs/op
PASS
ok  	github.com/alessiosavi/GoGPUtils/string	289.470s
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/zip
BenchmarkReadZipFile-8   	  947400	     12870 ns/op	    7808 B/op	      29 allocs/op
BenchmarkReadZip01-8     	  922921	     13144 ns/op	    8144 B/op	      31 allocs/op
PASS
```
