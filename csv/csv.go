package csvutils

import (
	"bytes"
	"encoding/csv"
	"errors"
	processingutils "github.com/alessiosavi/GoGPUtils/files/processing"
	"github.com/schollz/progressbar/v3"
	"strconv"
	"time"
)

// FIXME: Migrate the module to map[headers][]row
type Headers []string
type CSVData [][]string

type ExplodeCSV func(int, []string) [][]string
type ApplyCSV func(int, []string) []string
type ApplyHeader func(int, *string) string

func (c *Headers) Apply(fn ApplyHeader, inplace bool) Headers {
	var res Headers
	if inplace {
		res = *c
	} else {
		res = make(Headers, len(*c))
		copy(res, *c)
	}
	for i := range res {
		res[i] = fn(i, &res[i])
	}
	return res
}

func (c *CSVData) Explode(fn ExplodeCSV) CSVData {
	var res = make(CSVData, 0, len(*c))
	bar := progressbar.Default(int64(len(res)))
	for i := range *c {
		bar.Add(1)
		res = append(res, fn(i, (*c)[i])...)
	}
	*c = res
	bar.Close()
	return *c
}
func (c *CSVData) Apply(fn ApplyCSV, inplace bool) CSVData {
	var res CSVData
	if inplace {
		res = *c
	} else {
		res = make(CSVData, len(*c))
		copy(res, *c)
	}
	bar := progressbar.Default(int64(len(res)))
	for i := range res {
		bar.Add(1)
		res[i] = fn(i, res[i])
	}
	bar.Close()
	return res
}

// ReadCSV is delegated to read into a CSV the content of the bytes in input
// []string -> Headers of the CSV
// [][]string -> Content of the CSV
func ReadCSV(data []byte, separator rune, hasHeaders bool) (Headers, CSVData, error) {
	buf, err := processingutils.ToUTF8(data)
	if err != nil {
		return nil, nil, err
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

	var headers Headers
	if hasHeaders {
		if len(csvData) < 2 {
			return nil, nil, errors.New("input data does not contain at least 2 rows (headers + data)")
		}
		// Remove the headers from the row data
		headers = csvData[0]
		// Remove the first element due to headers shift
		csvData = csvData[1:]
	}
	return headers, csvData, nil
}
func WriteCSV(headers Headers, records CSVData, separator rune) ([]byte, error) {
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

// GetCol is delegated to filter a given column from the csvData
func GetCol(csvData CSVData, index int) []string {
	if len(csvData) == 0 || len(csvData[0]) < index {
		return nil
	}
	var ret = make([]string, len(csvData))
	for i := range csvData {
		ret[i] = csvData[i][index]
	}
	return ret
}

// GetCSVDataType is delegated to retrieve the data type for every field of the CSV
// Return: headers, csv data, data type, error
func GetCSVDataType(raw []byte, separator rune) (Headers, CSVData, map[string]string, error) {
	headers, data, err := ReadCSV(raw, separator, true)
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
			// INT was not checked for this header
			if !e["int"] {
				if _, err = strconv.ParseInt(row[i], 10, 64); err == nil {
					dataType[header] = "int"
					continue
				} else {
					// Error, INT can be used as type for this headers
					//log.Println(fmt.Sprintf("Line %d break the int rule [%s] for header %s", lineN+2, row[i], header))
					e["int"] = true
				}
			}
			if !e["float"] {
				if _, err = strconv.ParseFloat(row[i], 64); err == nil {
					dataType[header] = "float"
					continue
				} else {
					//log.Println(fmt.Sprintf("Line %d break the float rule [%s] for header %s", lineN+2, row[i], header))
					e["float"] = true
				}
			}
			if !e["bool"] {
				if _, err = strconv.ParseBool(row[i]); err == nil {
					dataType[header] = "bool"
					continue
				} else {
					////log.Println(fmt.Sprintf("Line %d break the bool rule [%s] for header %s", lineN+2, row[i], header))
					e["bool"] = true
				}
			}
			if row[i] == "" {
				dataType[header] = "string"
				//log.Println("Found empty value in line", lineN+2, "of column", header, "| Setting text type")
				continue
			}
			if row[i][0] == '0' {
				// A number that start with 0 is a valid number for golang, but from a data warehouse POV, it has to be saved as is, so it's better to use a string.
				// Example: 00100 will be saved as 100, that is not correct
				dataType[header] = "string"
				break
			}
			if !e["time"] {
				for _, layout := range []string{
					"2006-01-02T15:04:05Z",
					"2006-01-02",
					"2006-01-02 15:04:05",
					"2006-01-02T15:04:05.999999Z07:00",
					"2006-01-02 15:04:05.999999Z07",
					time.ANSIC,
					time.UnixDate,
					time.RubyDate,
					time.RFC822,
					time.RFC822Z,
					time.RFC850,
					time.RFC1123,
					time.RFC1123Z,
					time.RFC3339,
					time.RFC3339Nano,
					time.Kitchen,
					time.Stamp,
					time.StampMilli,
					time.StampMicro,
					time.StampNano} {

					if _, err := time.Parse(layout, row[i]); err == nil {
						dataType[header] = "time"
						e["time"] = false
						break
					}
					//log.Println(fmt.Sprintf("Line %d break the time rule [%s] for header %s", lineN+2, row[i], header))
					e["time"] = true
				}
				if !e["time"] {
					continue
				}
			}
			// fallback data type
			dataType[header] = "string"
		}

	}
	return headers, data, dataType, nil
}

func DecodeNonUTF8CSV(data CSVData) CSVData {
	for i := range data {
		for j := range data[i] {
			data[i][j] = processingutils.DecodeNonUTF8String(data[i][j])
		}
	}
	return data
}
