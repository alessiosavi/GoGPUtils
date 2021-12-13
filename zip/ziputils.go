package ziputils

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
)

// ReadZip is delegated to extract the files and read the content
func ReadZip(filename string) (map[string]string, error) {
	//log.Println("ReadZip | Opening zipped content from [" + filename + "]")
	zf, err := zip.OpenReader(filename)
	if err != nil {
		//log.Println("ReadZip | Error during read "+filename+" | Err: ", err)
		return nil, err
	}
	defer zf.Close()
	return read(zf)
}

func read(zf *zip.ReadCloser) (map[string]string, error) {
	var filesContent = make(map[string]string)
	for _, file := range zf.File {
		if file.Mode().IsRegular() {
			if data, err := ReadZipFile(file); err != nil {
				continue
			} else {
				filesContent[file.Name] = data
			}
		}
	}
	return filesContent, nil
}

// ReadZipBytes is delegated to extract the files and read the content
func ReadZipBytes(raw io.ReadCloser) (map[string]string, error) {
	bytesData, err := io.ReadAll(raw)
	if err != nil {
		return nil, err
	}
	if err = raw.Close(); err != nil {
		return nil, err
	}
	var filesContent map[string]string
	zf, err := zip.NewReader(bytes.NewReader(bytesData), int64(len(bytesData)))
	if err != nil {
		return nil, err
	}
	filesContent = make(map[string]string)
	for _, file := range zf.File {
		if file.Mode().IsRegular() {
			if data, err := ReadZipFile(file); err != nil {
				log.Printf("Error on file:%s\n", file.Name)
				continue
			} else {
				filesContent[file.Name] = data
			}
		}
	}
	return filesContent, nil
}

// ReadZipFile is a wrapper function for ioutil.ReadAll. It accepts a zip.File as
// its parameter, opens it, reads its content and returns it as a byte slice.
func ReadZipFile(file *zip.File) (string, error) {
	if !file.Mode().IsRegular() {
		return "", errors.New("ReadZipFile | File " + file.Name + " is not a regular!")
	}
	fc, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fc.Close()

	content, err := ioutil.ReadAll(fc)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
