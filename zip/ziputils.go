package ziputils

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"testing"
)

// ReadZip is delegated to extract the files and read the content
func ReadZip(filename string) (map[string]string, error) {
	var filesContent map[string]string

	//log.Println("ReadZip | Opening zipped content from [" + filename + "]")
	zf, err := zip.OpenReader(filename)
	if err != nil {
		//log.Println("ReadZip | Error during read "+filename+" | Err: ", err)
		return nil, err
	}
	defer zf.Close()
	filesContent = make(map[string]string)
	for _, file := range zf.File {
		if file.Mode().IsRegular() {
			//log.Println("ReadZip | Unzipping regular file " + file.Name)
			data, err := ReadZipFile(file)
			if err != nil {
				//log.Println("ReadZip | Unable to unzip file " + file.Name)
			} else {
				//log.Println("ReadZip | File unzipped successfully!")
				filesContent[file.Name] = data
			}
		}
	}
	//log.Println("ReadZip | Unzipped ", len(filesContent), " files")
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

func TestReadZip(t *testing.T) {}
func BenchmarkReadZip(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
