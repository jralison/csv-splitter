package splitter

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func CsvSplit(file *os.File, separator rune, columns []string, preserveColumn bool, outputFilenameTemplate string) error {
	reader := csv.NewReader(file)
	reader.Comma = separator

	headers, splitIndexes, err := parseColumnsIndexes(reader, columns)
	if err != nil {
		return err
	}

	outputHeaders := preservedValues(headers, splitIndexes, preserveColumn)

	toBeClosedDescriptors := make([]*os.File, 0)
	defer func(files []*os.File) {
		for _, f := range files {
			_ = f.Close()
		}
	}(toBeClosedDescriptors)

	writers := make(map[string]*csv.Writer, 8)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		outputFilename := filenameFromRecord(record, splitIndexes, headers, outputFilenameTemplate)
		if _, hasWriter := writers[outputFilename]; !hasWriter {
			outputFile, err := os.OpenFile(outputFilename, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			toBeClosedDescriptors = append(toBeClosedDescriptors, outputFile)
			writers[outputFilename] = csv.NewWriter(outputFile)
			if err := writers[outputFilename].Write(outputHeaders); err != nil {
				return err
			}
		}

		record = preservedValues(record, splitIndexes, preserveColumn)
		if err := writers[outputFilename].Write(record); err != nil {
			return err
		}
	}

	return nil
}

func filenameFromRecord(record []string, splitIndexes []int, headers []string, filenameTemplate string) string {
	suffix := ""
	for _, i := range splitIndexes {
		suffix += "_" + headers[i] + record[i]
	}
	return strings.Replace(filenameTemplate, "{suffix}", suffix, 1)
}

func preservedValues(record []string, splitIndexes []int, preserve bool) []string {
	if preserve {
		return record
	}

	preserved := make([]string, 0, len(record)-len(splitIndexes))
	left := 0
	for _, right := range splitIndexes {
		preserved = append(preserved, record[left:right]...)
		left = right + 1
	}
	preserved = append(preserved, record[left:]...)

	return preserved
}

func parseColumnsIndexes(reader *csv.Reader, columns []string) ([]string, []int, error) {
	headers, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}

	indexes := make([]int, len(columns))
	for c, column := range columns {
		indexes[c] = -1

		for h, header := range headers {
			if header == column {
				indexes[c] = h
				break
			}
		}

		if indexes[c] == -1 {
			return headers, nil, fmt.Errorf("%v is not a column in the file", column)
		}
	}

	return headers, indexes, nil
}
