package fileutils

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	stringutils "github.com/alessiosavi/GoGPUtils/string"
)

// ReadFileInArray is delegated to read the file content as tokenize the data by the new line
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
	// Read all the file recursively
	err := filepath.Walk(path, func(file string, f os.FileInfo, err error) error {
		if IsFile(file) {
			fileList = append(fileList, file)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return nil
	}
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

// CountLinesFile return the number of lines in the given file
func CountLinesFile(fileName string, bufferLenght int) (int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return -1, err
	}
	defer file.Close()
	r := bufio.NewReader(file)
	// 32K as buffer
	if bufferLenght == -1 {
		bufferLenght = 32
	}
	buf := make([]byte, bufferLenght*1024)
	count := 0
	lineSep := []byte{'\n'}
	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// FileExists verify that the file exist
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// GetFileContentType is delegated to retrieve the filetype for a given file path
func GetFileContentType(fileName string) (string, error) {
	if !FileExists(fileName) {
		return "", errors.New("File " + fileName + " does not exist!")
	}
	if !IsFile(fileName) {
		if IsDir(fileName) {
			return "directory", nil
		}
		// Maybe a link
		return "not_regular_file", errors.New("Not a regular file [" + fileName + "]")
	}
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()
	// Only the first 128 bytes are used to sniff the content type.
	buffer := make([]byte, 128)
	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}
	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	// In case of application/zip, we have to be sure the file type
	if contentType == "application/zip" {
		//fmt.Println(string(buffer))
		if bytes.Contains(buffer, []byte("mimetypeapplication/vnd.oasis.opendocument.tex")) {
			return "application/odt", nil
		}
		if bytes.Contains(buffer, []byte("rels/.rels")) {
			return "application/docx", nil
		}
	} else if contentType == "application/octet-stream" {
		// A pickle item have an X as 2nd byte
		if buffer[2] == 88 {
			return "application/pickle", nil
		} else if bytes.Contains(buffer, []byte("isomiso2mp41")) ||
			bytes.Contains(buffer, []byte("isomiso2avc1mp41")) {
			return "iso/mp4", nil
		}
		if bytes.Equal(buffer[1:4], []byte("ELF")) {
			return "elf/binary", nil
		}
		for {
			n, err := f.Read(buffer)
			if err == io.EOF || n == 0 {
				break
			} else if bytes.Contains(buffer, []byte("Microsoft Word-D")) {
				return "application/doc", nil
			}
		}
	}
	return contentType, nil
}
