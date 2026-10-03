// Command benchfingerprint prints the content fingerprint that
// scripts/bench.sh records for a checkout: the SHA-256 of every tracked or
// untracked, non-ignored file, or with -suite only the benchmark suite's
// files (test files and internal/benchkit).
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alessiosavi/GoGPUtils/internal/benchcov"
)

func main() {
	dir := flag.String("dir", ".", "checkout to fingerprint")
	suite := flag.Bool("suite", false, "only the benchmark suite's files")
	flag.Parse()
	var include func(string) bool
	if *suite {
		include = benchcov.SuiteFile
	}
	fp, err := benchcov.Fingerprint(*dir, include)
	if err != nil {
		fmt.Fprintln(os.Stderr, "benchfingerprint:", err)
		os.Exit(1)
	}
	fmt.Println(fp)
}
