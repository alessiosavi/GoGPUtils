package goutils

import (
	"bufio"
	"errors"
	arrayutils "github.com/alessiosavi/GoGPUtils/array"
	fileutils "github.com/alessiosavi/GoGPUtils/files"
	"log"
	"os"
	"strings"
)

// ExtractFunctionFromFile is delegated to filter the function present in the input file that have the given prefix
func ExtractFunctionFromFile(codeFile, prefix string) ([]string, error) {
	// functions will store the method present in the file
	var functions []string
	if fileutils.IsFile(codeFile) {
		text, err := os.ReadFile(codeFile)
		if err != nil {
			return nil, err
		}
		if len(text) == 0 {
			return nil, errors.New("file [" + codeFile + "] is an empty file")
		}
		// Read the file and filter the function
		scanner := bufio.NewScanner(strings.NewReader(string(text)))
		for scanner.Scan() {
			data := scanner.Text()
			if strings.HasPrefix(data, "func "+prefix) {
				stopIndex := strings.Index(data, "(")
				if stopIndex != -1 {
					startIndex := strings.Index(data, "func ") + len("func ")
					if startIndex != -1 {
						functions = append(functions, data[startIndex:stopIndex])
					}
				}
			}
		}
	} else {
		return nil, errors.New("file [" + codeFile + "] not found")
	}
	return functions, nil
}

// CreateBenchmarkSignature Is delegated to create the benchmark test signature for the input codeFile
func CreateBenchmarkSignature(codeFile string) (string, error) {
	// function will save the method present in the file
	functions, err := ExtractFunctionFromFile(codeFile, "")
	if err != nil {
		return "", err
	}
	testFile := strings.Replace(codeFile, ".go", "_test.go", 1)
	// Extract the benchmark that are already present
	benchAlreadyPresent, err := ExtractFunctionFromFile(testFile, "Benchmark")
	if err != nil {
		return "", err
	}
	// Remove the initial benchmark prefix
	for i := range benchAlreadyPresent {
		benchAlreadyPresent[i] = strings.TrimPrefix(benchAlreadyPresent[i], "Benchmark")
	}
	functions = arrayutils.RemoveStrings(functions, benchAlreadyPresent)
	var testfileContent strings.Builder

	if len(functions) > 0 {
		for i := range functions {
			testfileContent.WriteString("func Benchmark" + functions[i] + "(b *testing.B){for i := 0; i < b.N; i++ {}}\n")
		}
	} else {
		log.Println("Test Already generated")
	}
	return testfileContent.String(), nil
}

// CreateTestSignature Is Delegated to create the benchmark test signature for the input codeFile
func CreateTestSignature(codeFile string) (string, error) {
	// function will save the method present in the file
	functions, err := ExtractFunctionFromFile(codeFile, "")
	if err != nil {
		return "", err
	}
	testFile := strings.Replace(codeFile, ".go", "_test.go", 1)
	// Extract the benchmark that are already present
	testAlreadyPresent, err := ExtractFunctionFromFile(testFile, "")
	if err != nil {
		return "", err
	}
	// Remove the initial benchmark prefix
	for i := range testAlreadyPresent {
		testAlreadyPresent[i] = strings.TrimPrefix(testAlreadyPresent[i], "Test")
	}
	functions = arrayutils.RemoveStrings(functions, testAlreadyPresent)
	var testfileContent strings.Builder
	if len(functions) > 0 {
		for i := range functions {
			testfileContent.WriteString("func Test" + functions[i] + "(t *testing.T){}\n")
		}
	} else {
		log.Println("Test Already generated")
	}
	return testfileContent.String(), nil
}
