// Command benchscope prints the import paths of the packages covered by the
// benchmark suite, sorted, one per line. scripts/bench.sh and CI use it so
// that the coverage gate and the benchmark runs share one scope definition.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alessiosavi/GoGPUtils/internal/benchcov"
)

func main() {
	dir := flag.String("root", ".", "a directory inside the module to inspect")
	flag.Parse()
	root, modulePath, err := benchcov.FindModuleRoot(*dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "benchscope:", err)
		os.Exit(1)
	}
	pkgs, err := benchcov.Scope(root, modulePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "benchscope:", err)
		os.Exit(1)
	}
	for _, p := range pkgs {
		fmt.Println(p.ImportPath)
	}
}
