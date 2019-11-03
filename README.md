# GoGPUtils

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/86513a2282374f87a813110db86f018b)](https://www.codacy.com/manual/alessiosavi/GoGPUtils?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=alessiosavi/GoGPUtils&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/alessiosavi/GoGPUtils)](https://goreportcard.com/report/github.com/alessiosavi/GoGPUtils) [![GoDoc](https://godoc.org/github.com/alessiosavi/GoGPUtils?status.svg)](https://godoc.org/github.com/alessiosavi/GoGPUtils) [![License](https://img.shields.io/github/license/alessiosavi/GoGPUtils)](https://img.shields.io/github/license/alessiosavi/GoGPUtils) [![Version](https://img.shields.io/github/v/tag/alessiosavi/GoGPUtils)](https://img.shields.io/github/v/tag/alessiosavi/GoGPUtils) [![Code size](https://img.shields.io/github/languages/code-size/alessiosavi/GoGPUtils)](https://img.shields.io/github/languages/code-size/alessiosavi/GoGPUtils) [![Repo size](https://img.shields.io/github/repo-size/alessiosavi/GoGPUtils)](https://img.shields.io/github/repo-size/alessiosavi/GoGPUtils) [![Issue open](https://img.shields.io/github/issues/alessiosavi/GoGPUtils)](https://img.shields.io/github/issues/alessiosavi/GoGPUtils)
[![Issue closed](https://img.shields.io/github/issues-closed/alessiosavi/GoGPUtils)](https://img.shields.io/github/issues-closed/alessiosavi/GoGPUtils)

A simple Go general purpouse utility repository for avoid to reinvent the wheel every time that i need to start a new project.

## Benchmark

```text
$ go test -bench=. -benchmem -benchtime=5s ./...
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/files
BenchmarkCountLinesFile-8         	  465236	     12395 ns/op	   37008 B/op	       6 allocs/op
BenchmarkListFile-8               	    3960	   1443884 ns/op	  747355 B/op	    1672 allocs/op
BenchmarkFindFilesSensitive-8     	    4386	   1436463 ns/op	  743291 B/op	    1665 allocs/op
BenchmarkFindFilesInsensitive-8   	    4170	   1430643 ns/op	  743291 B/op	    1665 allocs/op
PASS
ok  	github.com/alessiosavi/GoGPUtils/files	30.439s
2019/11/03 20:02:32 Function before filter ->  [ExtractTextFromQuery CheckPresence IsUpper IsLower ContainsLetter ContainsOnlyLetter CreateJSON RemoveDoubleWhiteSpace RemoveWhiteSpace IsASCII IsASCIIRune RemoveFromString Split CountLines ExtractString ReplaceAtIndex RemoveNonASCII IsBlank Trim RandomString CheckPalindrome ReverseString] Len: 22
2019/11/03 20:02:32 Function after filter ->  [ContainsLetter IsASCIIRune ReplaceAtIndex IsBlank CheckPalindrome] Len: 5
PASS
ok  	github.com/alessiosavi/GoGPUtils/goutils	0.005s
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/helper
BenchmarkRandomIntn-8                       	  662750	      9232 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomInt32-8                      	  662794	      9219 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomInt64-8                      	  641769	      9249 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomFloat32-8                    	  639522	      9203 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomFloat64-8                    	  659672	      9216 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomIntnR-8                      	514047970	        12.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomInt32R-8                     	614509504	         9.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomInt64R-8                     	277488235	        21.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomFloat32R-8                   	861068836	         7.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomFloat64R-8                   	788966076	         7.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomIntnRArray-8                 	  382156	     14776 ns/op	    8192 B/op	       1 allocs/op
BenchmarkRandomInt32RArray-8                	  540450	     11214 ns/op	    4096 B/op	       1 allocs/op
BenchmarkRandomInt64RArray-8                	  259094	     23550 ns/op	    8192 B/op	       1 allocs/op
BenchmarkRandomFloat32Array-8               	  707756	      9149 ns/op	    4096 B/op	       1 allocs/op
BenchmarkRandomFloat64RArray-8              	  655885	      9764 ns/op	    8192 B/op	       1 allocs/op
PASS
ok  	github.com/alessiosavi/GoGPUtils/helper	102.007s
2019/11/03 20:04:15 listen tcp 127.0.0.1:8080: bind: address already in use
exit status 1
FAIL	github.com/alessiosavi/GoGPUtils/http	0.452s
2019/11/03 20:04:15 Arrays:  [1 1 2 3 4] [0 9 3 3 3]
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/math
BenchmarkSumIntArray-8             	13760162	       392 ns/op	       0 B/op	       0 allocs/op
BenchmarkSumInt32Array-8           	15577788	       385 ns/op	       0 B/op	       0 allocs/op
BenchmarkSumInt64Array-8           	 9920103	       604 ns/op	       0 B/op	       0 allocs/op
BenchmarkSumFloat32Array-8         	 5285958	      1147 ns/op	       0 B/op	       0 allocs/op
BenchmarkSumFloat64Array-8         	 5295811	      1134 ns/op	       0 B/op	       0 allocs/op
BenchmarkMaxIntIndex-8             	 6041853	       934 ns/op	       0 B/op	       0 allocs/op
BenchmarkMaxInt32Index-8           	 6419010	       920 ns/op	       0 B/op	       0 allocs/op
BenchmarkMaxInt64Index-8           	 5319060	      1005 ns/op	       0 B/op	       0 allocs/op
BenchmarkMaxFloat32Index-8         	 5110430	      1206 ns/op	       0 B/op	       0 allocs/op
BenchmarkMaxFloat64Index-8         	 5772138	       996 ns/op	       0 B/op	       0 allocs/op
BenchmarkAverageInt-8              	15230282	       381 ns/op	       0 B/op	       0 allocs/op
BenchmarkAverageInt32-8            	16072452	       382 ns/op	       0 B/op	       0 allocs/op
BenchmarkAverageInt64-8            	15377040	       389 ns/op	       0 B/op	       0 allocs/op
BenchmarkAverageFloat32-8          	 5311312	      1122 ns/op	       0 B/op	       0 allocs/op
BenchmarkAverageFloat64-8          	 5419072	      1112 ns/op	       0 B/op	       0 allocs/op
BenchmarkInitRandomMatrix-8        	  532336	     11199 ns/op	    6352 B/op	      13 allocs/op
BenchmarkMultiplySumArray1000-8    	 3155178	      1944 ns/op	    8192 B/op	       1 allocs/op
BenchmarkMultiplyMatrix100x100-8   	    1134	   4874564 ns/op	18101931 B/op	   20201 allocs/op
BenchmarkIsPrime-8                 	 1441786	      4136 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/alessiosavi/GoGPUtils/math	135.749s
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/search
BenchmarkLinearSearchInt-8                 	  288487	     20348 ns/op	       0 B/op	       0 allocs/op
BenchmarkLinearSearchParallelInt-8         	367461364	        16.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsStringByte-8              	  569611	      9844 ns/op	   22153 B/op	      15 allocs/op
BenchmarkContainsStringsByte-8             	  259568	     23238 ns/op	   58194 B/op	      33 allocs/op
BenchmarkContainsWhichStrings-8            	    5176	   1145881 ns/op	   58365 B/op	      36 allocs/op
PASS
ok  	github.com/alessiosavi/GoGPUtils/search	32.975s
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/string
BenchmarkContainsOnlyLetter-8       	249640588	        24.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkRemoveFromString-8         	   58827	    103743 ns/op	 1114116 B/op	       2 allocs/op
BenchmarkRandomString-8             	  258217	     22776 ns/op	    5376 B/op	       1 allocs/op
BenchmarkExtractTextFromQuery-8     	      80	  66453374 ns/op	21759789 B/op	  102049 allocs/op
BenchmarkRemoveWhiteSpace-8         	    2109	   2828598 ns/op	  557057 B/op	       1 allocs/op
BenchmarkIsASCII-8                  	162925543	        35.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkSplit-8                    	    3993	   1547392 ns/op	 2447640 B/op	   14363 allocs/op
BenchmarkSplitBuiltin-8             	   17044	    354168 ns/op	  319488 B/op	       1 allocs/op
BenchmarkExtractString-8            	   10000	    522806 ns/op	       0 B/op	       0 allocs/op
BenchmarkRemoveNonASCII-8           	     996	   5905851 ns/op	 1589252 B/op	       3 allocs/op
BenchmarkCreateJSON-8               	      15	 348614019 ns/op	2998418187 B/op	    9916 allocs/op
BenchmarkJoin-8                     	   10000	    621044 ns/op	 2930657 B/op	      29 allocs/op
BenchmarkTrim-8                     	    3049	   1984919 ns/op	  557063 B/op	       1 allocs/op
BenchmarkRemoveDoubleWhiteSpace-8   	    1935	   3060925 ns/op	  557057 B/op	       1 allocs/op
BenchmarkCountLines-8               	   18690	    332069 ns/op	    4128 B/op	       2 allocs/op
BenchmarkReverseString-8            	    4476	   1299330 ns/op	 2914327 B/op	      32 allocs/op
PASS
ok  	github.com/alessiosavi/GoGPUtils/string	132.674s
goos: linux
goarch: amd64
pkg: github.com/alessiosavi/GoGPUtils/zip
BenchmarkReadZipFile-8   	  419516	     13929 ns/op	    7808 B/op	      29 allocs/op
BenchmarkReadZip01-8     	  412348	     14240 ns/op	    8144 B/op	      31 allocs/op
PASS
ok  	github.com/alessiosavi/GoGPUtils/zip	12.016s
```
