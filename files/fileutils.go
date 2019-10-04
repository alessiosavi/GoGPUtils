package fileutils

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	stringutils "github.com/alessiosavi/GoGPUtils/string"
)

// ReadAllFileInArray is delegated to read the file content as tokenize the data by the new line
func ReadFileInArray(filePath string) []string {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil
	}
	return stringutils.Split(string(data))
}

//IsFile verify if a give filepath is a directory
func IsFile(path string) bool {
	fi, err := os.Stat(path)
	if os.IsNotExist(err) || err != nil {
		return false
	}
	return !fi.Mode().IsDir()
}

// IsDir is delegated to verify that the given path is a directory
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || err != nil {
		return false
	}
	return info.IsDir()
}

// GetFileModification return the last modification time of the file in input in a UNIX time format
func GetFileModification(filepath string) int64 {
	f, err := os.Open(filepath)
	if err != nil {
		return -1
	}
	defer f.Close()
	statinfo, err := f.Stat()
	if err != nil {
		return -2
	}
	return statinfo.ModTime().Unix()
}

// GetFileDate is delegated to return the date in a string format in which the file was (latest) modified
func GetFileDate(filepath string) string {
	unixTimestamp := GetFileModification(filepath)
	if unixTimestamp != -1 {
		loc, err := time.LoadLocation("Europe/Rome")
		if err != nil {
			return ""
		}
		currentTime := time.Unix(unixTimestamp, 0).In(loc)
		date := currentTime.Format("2006-01-02 15:04:05")
		return date
	}
	return ""
}

// FindFiles is delegated to find the files from the given directory, recursively for each dir
func FindFiles(path string) []string {
	fileList := []string{}
	// Read all the file recursivly
	filepath.Walk(path, func(file string, f os.FileInfo, err error) error {
		if IsFile(file) {
			fileList = append(fileList, file)
		}
		return nil
	})
	return fileList
}

// VerifyFilesExists is delegated to verify that the given list of file exist in the directory
func VerifyFilesExists(filePath string, files []string) bool {
	if IsDir(filePath) {
		for i := range files {
			filename := path.Join(filePath, files[i])
			if !IsFile(filename) {
				return false
			}
		}
		return true
	}
	return false
}
