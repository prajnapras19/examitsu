package lib

import (
	"encoding/csv"
	"io"
)

func ReadCSV(openedFile io.ReadCloser) (header []string, records []map[string]string, err error) {
	reader := csv.NewReader(openedFile)

	header, err = reader.Read()
	if err == io.EOF {
		return nil, nil, err
	}

	records = []map[string]string{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		if len(record) != len(header) {
			return nil, nil, ErrCSVRecordNotMatchedWithHeader
		}
		recordMap := map[string]string{}
		for i := range record {
			recordMap[header[i]] = record[i]
		}
		records = append(records, recordMap)
	}

	return header, records, nil
}
