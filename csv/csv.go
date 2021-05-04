package csv

import (
	"bytes"
	"encoding/csv"
	"github.com/alessiosavi/GoGPUtils/files/processing"
	"log"
)

// ReadCSV is delegated to read into a CSV the content of the bytes in input
// []string -> Headers of the CSV
// [][]string -> Content of the CSV
func ReadCSV(buf []byte, separator rune) ([]string, [][]string, error) {
	terminator, err := processing.DetectLineTerminator(bytes.NewReader(buf))
	// Clean file if possible ...
	if err == nil {
		log.Println("Cleaning file ...")
		buf = bytes.ReplaceAll(buf, []byte(terminator), []byte("\n"))
		buf = bytes.ReplaceAll(buf, []byte("\u001D"), []byte{}) // Remove group separator
		buf = bytes.ReplaceAll(buf, []byte("\u000B"), []byte{}) // Remove vertical tab
		buf = bytes.Trim(buf, "\n")
		buf = bytes.TrimSpace(buf)
	}

	csvReader := csv.NewReader(bytes.NewReader(buf))
	csvReader.Comma = separator
	csvReader.LazyQuotes = true
	csvReader.TrimLeadingSpace = true
	csvData, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	headers := csvData[0]
	// Remove the headers from the row data
	csvData = csvData[1:]
	// Remove the latest element due to headers shift
	return headers, csvData, nil
}
func WriteCSV(headers []string, records [][]string, separator rune) ([]byte, error) {
	var buff bytes.Buffer
	writer := csv.NewWriter(&buff)
	defer writer.Flush()
	writer.Comma = separator
	if err := writer.Write(headers); err != nil {
		return nil, err
	}
	if err := writer.WriteAll(records); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(buff.Bytes(), []byte("\n")), nil
}
