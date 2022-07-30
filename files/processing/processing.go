package processingutils

import (
	"bytes"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"unicode/utf8"
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
	buff := make([]byte, 1024*1000)
	var counts = make(map[LineTerminatorType]int)
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
	var maxKey = ND
	for k, v := range counts {
		if v > maxV {
			maxV = v
			maxKey = k
		}
	}
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

// ReplaceLineTerminatorBytesReader is delegated to find the line terminator of the given byte array and replace them without the one provided in input
func ReplaceLineTerminatorBytesReader(data *bytes.Reader, newLineTerminator []byte) ([]byte, error) {
	terminator, err := DetectLineTerminator(data)
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

func DecodeNonUTF8String(data string) string {
	if !utf8.ValidString(data) {
		v := make([]rune, 0, len(data))
		for k, r := range data {
			if r == utf8.RuneError {
				if _, size := utf8.DecodeRuneInString(data[k:]); size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		return string(v)
	}
	return data
}

func DecodeNonUTF8Bytes(data []byte) []byte {
	if !utf8.Valid(data) {
		v := make([]rune, 0, len(data))
		for k, r := range string(data) {
			if r == utf8.RuneError {
				if _, size := utf8.DecodeRune(data[k:]); size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		return []byte(string(v))
	}
	return data
}

func ToUTF8(data []byte) ([]byte, error) {
	// Clean file if possible ...
	if terminator, err := DetectLineTerminator(bytes.NewReader(data)); err == nil {
		data = bytes.ReplaceAll(data, []byte(terminator), []byte("\n"))
		// Remove other line terminator if present
		for _, ter := range []LineTerminatorType{CR, CRLF, LFCR, RS} {
			data = bytes.ReplaceAll(data, []byte(ter), []byte(" "))
		}
	}
	var err error
	encoder, _, ok := charset.DetermineEncoding(data, http.DetectContentType(data))
	if ok {
		r := transform.NewReader(bytes.NewReader(data), encoder.NewDecoder())
		data, err = ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
	}

	for _, BOM := range stringutils.BOMS {
		data = bytes.TrimPrefix(data, BOM)
	}
	//buf = bytes.ReplaceAll(buf, []byte("\u001D"), []byte{}) // Remove group separator
	//buf = bytes.ReplaceAll(buf, []byte("\u000B"), []byte{}) // Remove vertical tab
	data = []byte(html.UnescapeString(string(data)))
	data = bytes.TrimSpace(data)
	data = bytes.Trim(data, "\n")

	return DecodeNonUTF8Bytes(data), nil
}
