package goutils

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	fileutils "github.com/alessiosavi/GoGPUtils/files"
)

func CreateBenchmarkSignature(codeFile string) (string, error) {
	if !fileutils.IsFile(codeFile) {
		log.Println(fileutils.ListFile("."))
		return "", errors.New("unable to find file " + codeFile)
	}
	text, err := ioutil.ReadFile(codeFile)
	if err != nil {
		return "", err
	}
	if len(text) == 0 {
		return "", errors.New("file " + codeFile + " is an empty file")
	}

	var functions []string
	scanner := bufio.NewScanner(strings.NewReader(string(text)))
	for scanner.Scan() {
		data := scanner.Text()
		if strings.HasPrefix(data, "func ") {
			stopIndex := strings.Index(data, "(")
			if stopIndex != -1 {
				startIndex := strings.Index(data, "func ") + len("func ")
				if startIndex != -1 {
					functions = append(functions, data[startIndex:stopIndex])
				}
			}
		}
	}
	var testfileContent strings.Builder

	if len(functions) > 0 {
		for i := range functions {
			testfileContent.WriteString("func Test" + functions[i] + "(t *testing.T){}\n")
			testfileContent.WriteString("func Benchmark" + functions[i] + "(b *testing.B){for i := 0; i < b.N; i++ {}}\n")
		}
	}
	return testfileContent.String(), nil
}

// GenerateTestSignature is delegated to generate the test signature for every golang code file in the subdirectory
func GenerateTestSignature(fileFolder, outFolder string) error {
	files := fileutils.FindFiles(fileFolder, ".go", true)
	i := 0
	for _, file := range files {
		if !strings.Contains(file, "_test.go") {
			files[i] = file
			i++
		}
	}
	files = files[:i]
	err := fileutils.CreateDir(outFolder)
	if err != nil {
		return err
	}
	for _, item := range files {
		//t.Log("=======" + item + "=======")
		data, err := CreateBenchmarkSignature(item)
		if err != nil {
			log.Println(err)
		} else {
			item = strings.Replace(item, "../", "", -1)
			index := strings.Index(item, "/")
			var codeFolder string
			if index != -1 {
				codeFolder = outFolder + "/" + item[:index]
				item = item[:index] + "_test.go"
			} else {
				codeFolder = outFolder + "/" + item
			}
			if !strings.Contains(item, "_test.go") {
				item = strings.Replace(item, ".go", "_test.go", -1)
			}
			log.Println(outFolder, item, codeFolder)
			err = fileutils.CreateDir(codeFolder)
			if err != nil {
				return err
			}
			folder := codeFolder + "/" + item
			f, err := os.Create(folder)
			if err != nil {
				log.Println("Error during folder creation", err)
			} else {
				_, err = f.WriteString(data)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
