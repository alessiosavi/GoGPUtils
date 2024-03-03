<div align="center">

# GoGPUtils

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/86513a2282374f87a813110db86f018b)](https://www.codacy.com/manual/alessiosavi/GoGPUtils?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=alessiosavi/GoGPUtils&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/alessiosavi/GoGPUtils)](https://goreportcard.com/report/github.com/alessiosavi/GoGPUtils) [![GoDoc](https://godoc.org/github.com/alessiosavi/GoGPUtils?status.svg)](https://godoc.org/github.com/alessiosavi/GoGPUtils) [![License](https://img.shields.io/github/license/alessiosavi/GoGPUtils)](https://img.shields.io/github/license/alessiosavi/GoGPUtils) [![Version](https://img.shields.io/github/v/tag/alessiosavi/GoGPUtils)](https://img.shields.io/github/v/tag/alessiosavi/GoGPUtils) [![Code size](https://img.shields.io/github/languages/code-size/alessiosavi/GoGPUtils)](https://img.shields.io/github/languages/code-size/alessiosavi/GoGPUtils) [![Repo size](https://img.shields.io/github/repo-size/alessiosavi/GoGPUtils)](https://img.shields.io/github/repo-size/alessiosavi/GoGPUtils) [![Issue open](https://img.shields.io/github/issues/alessiosavi/GoGPUtils)](https://img.shields.io/github/issues/alessiosavi/GoGPUtils)
[![Issue closed](https://img.shields.io/github/issues-closed/alessiosavi/GoGPUtils)](https://img.shields.io/github/issues-closed/alessiosavi/GoGPUtils)

A simple Go general-purpose utility repository for avoid reinventing the wheel every time that I need to start a new
project.


## Benchmark

```text
$ go test -bench=. -benchmem -benchtime=5s `go list ./... | grep -v "aws\|sftp\|http"`
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/byte
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkTestIsUpperByteOK-8   	1000000000	         2.174 ns/op	       0 B/op	       0 allocs/op
BenchmarkTestIsLowerByteKO-8   	1000000000	         2.185 ns/op	       0 B/op	       0 allocs/op

goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/files
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkCountLinesFile-8         	  137064	     48990 ns/op	   37000 B/op	       5 allocs/op
BenchmarkListFile-8               	    2167	   2795762 ns/op	  634139 B/op	    7390 allocs/op
BenchmarkFindFilesSensitive-8     	    2142	   2821969 ns/op	  601264 B/op	    7378 allocs/op
BenchmarkFindFilesInsensitive-8   	    2152	   2789374 ns/op	  601175 B/op	    7378 allocs/op

goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/helper
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkRandomIntn-8            	406849989	        14.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomInt32-8           	417671403	        14.45 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomInt64-8           	317977765	        18.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomFloat32-8         	483547185	        12.38 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomFloat64-8         	450756500	        13.33 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomIntnR-8           	854001973	         6.984 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomInt32R-8          	1000000000	         5.935 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomInt64R-8          	382641378	        15.59 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomFloat32R-8        	1000000000	         4.375 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomFloat64R-8        	1000000000	         4.039 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomIntnRArray-8      	  706681	     11181 ns/op	    8192 B/op	       1 allocs/op
BenchmarkRandomInt32RArray-8     	  720908	      7630 ns/op	    4096 B/op	       1 allocs/op
BenchmarkRandomInt64RArray-8     	  282798	     18387 ns/op	    8192 B/op	       1 allocs/op
BenchmarkRandomFloat32Array-8    	 1000000	      6615 ns/op	    4096 B/op	       1 allocs/op
BenchmarkRandomFloat64RArray-8   	  973923	     14226 ns/op	    8192 B/op	       1 allocs/op
BenchmarkRandomString-8          	  372956	     16016 ns/op	    5376 B/op	       1 allocs/op

goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/math
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkSumIntArray-8                   	21259684	       276.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkSumInt32Array-8                 	21679155	       276.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkSumInt64Array-8                 	21689191	       276.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkSumFloat32Array-8               	 7149135	       837.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkSumFloat64Array-8               	 7130800	       838.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkMaxIntIndex-8                   	 9218755	       660.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkMaxInt32Index-8                 	 9081741	       663.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkMaxInt64Index-8                 	 6498460	       905.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkMaxFloat32Index-8               	 7288107	       807.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkMaxFloat64Index-8               	 6400204	       919.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkAverageInt-8                    	21770044	       275.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkAverageInt32-8                  	21702019	       279.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkAverageInt64-8                  	21416380	       280.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkAverageFloat32-8                	 7214160	       832.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkAverageFloat64-8                	 7128327	       831.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkInitRandomMatrix-8              	  752602	     10999 ns/op	    6352 B/op	      13 allocs/op
BenchmarkMultiplySumArray1000-8          	 1000000	      5371 ns/op	    8192 B/op	       1 allocs/op
BenchmarkMultiplyMatrixLegacy100x100-8   	     408	  15497204 ns/op	18012339 B/op	   20101 allocs/op
BenchmarkMultiplyMatrix100x100-8         	    4159	   1447695 ns/op	   92288 B/op	     101 allocs/op
BenchmarkIsPrime-8                       	 2652973	      2249 ns/op	       0 B/op	       0 allocs/op
BenchmarkCosineSimilarity-8              	62873203	        94.24 ns/op	       0 B/op	       0 allocs/op

goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/search
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkLinearSearchInt-8           	  761466	      7683 ns/op	       0 B/op	       0 allocs/op
BenchmarkLinearSearchParallelInt-8   	212407228	        30.05 ns/op	      16 B/op	       1 allocs/op
BenchmarkContainsStringByte-8        	  266810	     31484 ns/op	   22146 B/op	      15 allocs/op
BenchmarkContainsStringsByte-8       	  118256	     50534 ns/op	   58196 B/op	      33 allocs/op
BenchmarkContainsWhichStrings-8      	    7248	    830028 ns/op	   58332 B/op	      36 allocs/op

goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/string
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkLevenshteinDistanceLegacy-8   	     100	  52979171 ns/op	51286027 B/op	    2502 allocs/op
BenchmarkLevenshteinDistance-8         	     452	  13153222 ns/op	   40960 B/op	       2 allocs/op
BenchmarkDiceCoefficient-8             	  137314	     43159 ns/op	       0 B/op	       0 allocs/op
BenchmarkJaroDistance-8                	    5758	   1032185 ns/op	    5376 B/op	       2 allocs/op
BenchmarkContainsOnlyLetter-8          	301701169	        19.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkRemoveFromString-8            	   36577	    193441 ns/op	 1114117 B/op	       2 allocs/op
BenchmarkExtractTextFromQuery-8        	     126	  45763565 ns/op	21765488 B/op	  102064 allocs/op
BenchmarkCheckPresence-8               	1000000000	         2.426 ns/op	       0 B/op	       0 allocs/op
BenchmarkIsUpper-8                     	1000000000	         3.163 ns/op	       0 B/op	       0 allocs/op
BenchmarkIsLower-8                     	1000000000	         3.486 ns/op	       0 B/op	       0 allocs/op
BenchmarkRemoveWhiteSpace-8            	    2773	   2275941 ns/op	  557058 B/op	       1 allocs/op
BenchmarkIsASCII-8                     	206237712	        28.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkSplit-8                       	    3511	   3395051 ns/op	 2446870 B/op	   14363 allocs/op
BenchmarkSplitBuiltin-8                	    6294	    896215 ns/op	  319488 B/op	       1 allocs/op
BenchmarkExtractString-8               	   15354	    392595 ns/op	       0 B/op	       0 allocs/op
BenchmarkRemoveNonASCII-8              	    1456	   4392784 ns/op	 1589249 B/op	       3 allocs/op
BenchmarkTestIsUpperOK-8               	1000000000	         3.142 ns/op	       0 B/op	       0 allocs/op
BenchmarkTestIsLowerOK-8               	1000000000	         3.492 ns/op	       0 B/op	       0 allocs/op
BenchmarkCreateJSON-8                  	      13	 665361057 ns/op	2998426824 B/op	   10020 allocs/op
BenchmarkJoin-8                        	   10275	    739491 ns/op	 2931697 B/op	      30 allocs/op
BenchmarkTrim-8                        	    1995	   2988434 ns/op	 1672093 B/op	      13 allocs/op
BenchmarkRemoveDoubleWhiteSpace-8      	    2520	   2329044 ns/op	  557056 B/op	       1 allocs/op
BenchmarkCountLines-8                  	   28654	    209826 ns/op	    4128 B/op	       2 allocs/op
BenchmarkReverseString-8               	    6412	   3824754 ns/op	 2914315 B/op	      32 allocs/op

goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/zip
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkReadZipFile-8   	  499692	     10876 ns/op	    6320 B/op	      27 allocs/op
BenchmarkReadZip01-8     	  880110	     12334 ns/op	    6656 B/op	      29 allocs/op
```

</div>