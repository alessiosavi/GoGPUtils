package processing

import (
	"bufio"
	"bytes"
	"github.com/alessiosavi/GoGPUtils/helper"
	"io"
	"io/ioutil"
	"log"
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

func DetectLineTerminator(reader io.Reader) (LineTerminatorType, error) {
	// Read 1mb of file
	log.SetFlags(log.Llongfile | log.LstdFlags)

	buff := make([]byte, 1024*1000)
	var counts map[LineTerminatorType]int = make(map[LineTerminatorType]int)
	for {
		if _, err := reader.Read(buff); err != nil {
			if err != io.EOF {
				return ND, err
			} else {
				break
			}
		}
		counts[CRLF] = bytes.Count(buff, []byte("\r\n"))
		counts[LFCR] = bytes.Count(buff, []byte("\n\r"))
		counts[CR] = bytes.Count(buff, []byte("\r"))
		counts[LF] = bytes.Count(buff, []byte("\n"))
		counts[RS] = bytes.Count(buff, []byte("\036"))

	}

	counts[CR] -= counts[CRLF] + counts[LFCR]
	counts[LF] -= counts[CRLF] + counts[LFCR]
	maxV := 0
	var maxKey LineTerminatorType = ND
	for k, v := range counts {
		if v > maxV {
			maxV = v
			maxKey = k
		}
	}
	log.Println(helper.MarshalIndent(counts))
	return maxKey, nil
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
