package processing

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
)

type LineTerminatorType string

const (
	LF   LineTerminatorType = `\n`
	CR   LineTerminatorType = `\r`
	CRLF LineTerminatorType = `\r\n`
	LFCR LineTerminatorType = `\n\r`
	RS   LineTerminatorType = `\036`
	ND   LineTerminatorType = `unable to detect line terminator`
)

func DetectLineTerminator(filename string) (LineTerminatorType, error) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	buff := make([]byte, 1024*10)

	for {
		if _, err = reader.Read(buff); err != nil {
			if err != io.EOF {
				return "", err
			} else {
				return "", errors.New(string(ND))
			}
		}

		if bytes.Contains(buff, []byte{'\r', '\n'}) {
			return CRLF, nil
		}

		if bytes.Contains(buff, []byte{'\n', '\r'}) {
			return LFCR, nil
		}
		if bytes.Contains(buff, []byte{'\r'}) {
			return CR, nil
		}

		if bytes.Contains(buff, []byte{'\n'}) {
			return LF, nil
		}
		if bytes.Contains(buff, []byte{'\036'}) {
			return RS, nil
		}
	}
}
