package processing

import (
	"bytes"
	"errors"
	"io"
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
