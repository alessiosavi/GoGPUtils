// Package fileutils provided a set of method for work with files
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
	"strings"
	"time"

	arrayutils "github.com/alessiosavi/GoGPUtils/array"
)

// Tail is delegated to read the latest lines of the file.
// NOTE: buffer have to be lesser than the minimum string lenght
func Tail(FILE string, BUFF_BYTE int64, START_POS, N_STRING int) string {
	// list of strings readed
	var stringsArray []string = make([]string, N_STRING)
	// Contains the data
	var buff []byte = make([]byte, -BUFF_BYTE)

	if !(START_POS >= 0 && START_POS <= 2) {
		log.Fatal("Wrong argument for Seek ...")
	}

	file, err := os.Open(FILE)
	if err != nil {
		log.Println("Unable to open file: " + FILE + " ERR: " + err.Error())
		log.Fatal(err)
	}
	defer file.Close()

	// Go to end of file
	_, err = file.Seek(BUFF_BYTE, START_POS)
	if err != nil {
		log.Println("Unable to seek to the end of the file: " + FILE + " ERR: " + err.Error())
		log.Fatal(err)
	}
	var linesReaded int = 0
	var nByte int                // Number of byte readed
	var stringBuffer string = "" // Contains the string until we don't found the new line
	var iteration int64 = 1
	var n int64 = -BUFF_BYTE // Just for pass the first check
	var lastPosition int64
	// Until we haven't read all the string
	for linesReaded < N_STRING {
		if n >= -BUFF_BYTE {
			n, err = file.Seek(iteration*BUFF_BYTE, START_POS)
			if err != nil {
				log.Println("2) Error during read of file | Lines readed: ", linesReaded, " Byte readed: ", nByte, " Iteration: ", iteration)
				log.Fatal(err)
			}
			lastPosition = n
		} else {
			// We have not enought data for fill the buffer, seeking to the start of the file
			file.Seek(0, 0)
			buff = make([]byte, lastPosition)
			_, err = file.Read(buff)
			if err != nil {
				log.Println("3) Error during read of file | Lines readed: ", linesReaded, " Byte readed: ", nByte, " Iteration: ", iteration)
				log.Fatal(err)
			}
			stringBuffer = string(buff) + stringBuffer
			stringsArray[N_STRING-linesReaded-1] = stringBuffer
			break
		}

		// Read the string related to the buffer
		nByte, err = file.Read(buff)
		if err != nil {
			log.Println("1) Error during read of file | Lines readed: ", linesReaded, " Byte readed: ", nByte, " Iteration: ", iteration)
			log.Fatal(err)
		}
		// Append the string in initial position
		stringBuffer = string(buff) + stringBuffer
		if strings.Contains(stringBuffer, "\n") {
			stringsArray[N_STRING-linesReaded-1] = stringBuffer
			stringBuffer = ""
			linesReaded++
			// Continue to read, we have not found a new line and we have enough file to read
		}
		iteration++
	}
	stringsArray = stringsArray[linesReaded-1:]
	return arrayutils.JoinStrings(stringsArray, "")
}

// ReadFileInArray is delegated to read the file content as tokenize the data by the new line
func ReadFileInArray(filePath string) []string {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil
	}
	return strings.Split(string(data), "\n")
}

//IsFile verify if the given filepath is a file
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

// CreateDir is delegated to create a new directory if not present
func CreateDir(path string) error {
	var err error
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println("Error during creation of folder ", err)
		}
	}
	return err
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

// ListFile is delegated to find the files from the given directory, recursively for each dir
func ListFile(path string) []string {
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

// FindFiles is delegated to find the files from the given directory, recursively for each dir, and extract only the one that match the input
func FindFiles(path, target string, caseSensitive bool) []string {
	fileList := []string{}
	if caseSensitive {
		// Read all the file recursively, taking care about the case of the string
		err := filepath.Walk(path, func(file string, f os.FileInfo, err error) error {
			if IsFile(file) && strings.Contains(file, target) {
				fileList = append(fileList, file)
			}
			return nil
		})

		if err != nil {
			log.Println(err)
			return nil
		}
	} else {
		// Case insensitive
		// Read all the file recursively, without taking care about the case of the string
		target = strings.ToLower(target)
		err := filepath.Walk(path, func(file string, f os.FileInfo, err error) error {
			if IsFile(file) && strings.Contains(strings.ToLower(file), target) {
				fileList = append(fileList, file)
			}
			return nil
		})

		if err != nil {
			log.Println(err)
			return nil
		}
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
// If called with an empty separator, new line will be used as default
func CountLinesFile(fileName, separator string, bufferLenght int) (int, error) {
	var lineSep []byte
	var buf []byte
	var count int

	file, err := os.Open(fileName)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	// 32K as buffer in case of not provided
	if bufferLenght == -1 {
		bufferLenght = 32
	}
	count = 0
	if len(separator) == 0 {
		lineSep = []byte{'\n'}
	}

	r := bufio.NewReader(file)
	buf = make([]byte, bufferLenght*1024)

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
		if bytes.Contains(buffer, []byte("mimetypeapplication/vnd.oasis.opendocument.text")) {
			return "application/odt", nil
		}
		if bytes.Contains(buffer, []byte("rels/.rels")) {
			return "application/docx", nil
		}
	} else if contentType == "application/octet-stream" {
		// A pickle item have an X as 2nd byte
		if buffer[2] == 88 {
			return "application/pickle", nil
			// Not necessary, mp4 will be catched by the http.DetectContentType
		} else if bytes.Contains(buffer, []byte("isomiso2mp41")) ||
			bytes.Contains(buffer, []byte("isomiso2avc1mp41")) {
			return "iso/mp4", nil
		}
		if bytes.Equal(buffer[1:4], []byte("ELF")) {
			return "elf/binary", nil
		}
		// Read the data until we found Microsoft Word-D
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

// GetFileSize is delegated to return the bytes size of the given file
func GetFileSize(filepath string) (int64, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	// get the size
	return fi.Size(), nil
}

// GetFileSize2 is a less efficient method for calculate the file size
func GetFileSize2(filepath string) (int64, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

// FilterFromFile is delegated to retrieve the lines that contain the target
func FilterFromFile(filename, target string, ignorecase bool) []string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return nil
	}
	if len(data) == 0 {
		return nil
	}

	var result []string
	if ignorecase {
		target = strings.ToLower(target)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		data := strings.ToLower(scanner.Text())
		if strings.Contains(data, target) {
			result = append(result, data)
		}
	}

	if len(result) == 0 {
		return nil
	}
	return result
}
