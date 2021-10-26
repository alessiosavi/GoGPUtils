package csv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"github.com/alessiosavi/GoGPUtils/files/processing"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"strconv"
)

// ReadCSV is delegated to read into a CSV the content of the bytes in input
// []string -> Headers of the CSV
// [][]string -> Content of the CSV
func ReadCSV(buf []byte, separator rune) ([]string, [][]string, error) {
	terminator, err := processing.DetectLineTerminator(bytes.NewReader(buf))
	// Clean file if possible ...
	if err == nil {
		buf = bytes.Replace(buf, stringutils.BOM, []byte{}, 1)
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
	csvReader.ReuseRecord = true
	csvData, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	csvData = DecodeNonUTF8CSV(csvData)

	if len(csvData) < 2 {
		return nil, nil, errors.New("input data does not contains at least 2 rows (headers + data)")
	}
	// Remove the headers from the row data
	headers := csvData[0]
	// Remove the first element due to headers shift
	csvData = csvData[1:]
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

// GetCSVDataType is delegated to retrieve the data type for every field of the CSV
// Return: headers, csv data, data type, error
func GetCSVDataType(raw []byte, separator rune) ([]string, [][]string, map[string]string, error) {
	headers, data, err := ReadCSV(raw, separator)
	if err != nil {
		return nil, nil, nil, err
	}
	// key = headers ; value = type
	var dataType = make(map[string]string)
	// Iterating headers
	for i, header := range headers {
		// Track if a given type is already tested and returned an error
		// e[<type>] = True -> Error, skip check for the give <type>
		// e[<type>] = False -> Not checked, continue trying to parse the field
		var e = make(map[string]bool)
		for _, row := range data {
			if row[i][0] == '0' {
				// A number that start with 0 is a valid number for golang, but from a data warehouse POV, it has to be saved as is, so it's better to use a string.
				// Example: 00100 will be saved as 100, that is not correct
				dataType[header] = "string"
				break
			}
			// INT was not checked for this header
			if !e["int"] {
				if _, err = strconv.ParseInt(row[i], 10, 64); err == nil {
					dataType[header] = "int"
					continue
				} else {
					// Error, INT can be used as type for this headers
					e["int"] = true
				}
			}
			if !e["float"] {
				if _, err = strconv.ParseFloat(row[i], 64); err == nil {
					dataType[header] = "float"
					continue
				} else {
					e["float"] = true
				}
			}
			if !e["bool"] {
				if _, err = strconv.ParseBool(row[i]); err == nil {
					dataType[header] = "bool"
					continue
				} else {
					e["bool"] = true
				}
			}
			// fallback data type
			dataType[header] = "string"
		}
	}
	return headers, data, dataType, nil
}

func DecodeNonUTF8CSV(data [][]string) [][]string {
	for i := range data {
		for j := range data[i] {
			data[i][j] = processing.DecodeNonUTF8String(data[i][j])
		}
	}
	return data
}
