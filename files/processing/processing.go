package processing

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
)

type LineTerminatorType string

const (
	LF   LineTerminatorType = "\n"   // 0a -- 10
	CR   LineTerminatorType = "\r"   // 0d -- 13
	CRLF LineTerminatorType = "\r\n" // 0d,0a -- 13,10
	LFCR LineTerminatorType = "\n\r" // 0a,0d -- 10,13
	RS   LineTerminatorType = "\036" // 1e -- 30
	ND   LineTerminatorType = `unable to detect line terminator`
)

// FIXME: Select the line terminator by reading all the file and finding the one that occurs more time
func DetectLineTerminator(reader io.Reader) (LineTerminatorType, error) {
	buff := make([]byte, 1024*10)

	for {
		if _, err := reader.Read(buff); err != nil {
			if err != io.EOF {
				return "", err
			} else {
				return "", errors.New(string(ND))
			}
		}
		if bytes.Contains(buff, []byte("\r\n")) {
			return CRLF, nil
		}
		if bytes.Contains(buff, []byte("\n\r")) {
			return LFCR, nil
		}
		if bytes.Contains(buff, []byte("\r")) {
			return CR, nil
		}
		if bytes.Contains(buff, []byte("\n")) {
			return LF, nil
		}
		if bytes.Contains(buff, []byte("\036")) {
			return RS, nil
		}
	}
}

// ReplaceLineTerminator is delegated to find the line terminator of the given byte array and replace them without the one provided in input
func ReplaceLineTerminator(data, newLineTerminator []byte) ([]byte, error) {
	terminator, err := DetectLineTerminator(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	newData := bytes.ReplaceAll(data, []byte(terminator), newLineTerminator)
	newData = bytes.TrimSpace(newData)
	return bytes.Trim(newData, string(newLineTerminator)), nil
}

// ReplaceLineTerminator is delegated to find the line terminator of the given byte array and replace them without the one provided in input
func ReplaceLineTerminatorBytesReader(data *bytes.Reader, newLineTerminator []byte) ([]byte, error) {

	terminator, err := DetectLineTerminator(bufio.NewReader(data))
	if err != nil {
		return nil, err
	}

	// Read all the file
	newData, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}
	// Reset to the start of the file
	data.Seek(0, io.SeekStart)
	newData = bytes.ReplaceAll(newData, []byte(terminator), newLineTerminator)
	newData = bytes.TrimSpace(newData)
	return bytes.Trim(newData, string(newLineTerminator)), nil
}
