package importer

import (
	"encoding/csv"
	"io"
	"time"
)

func Import(r io.Reader) (*Table, error) {
	csvReader := csv.NewReader(r)
	// read header
	_, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	var rows []Row
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		date, err := time.Parse("2006-01-02 15:04:05", record[1])
		if err != nil {
			return nil, err
		}
		rows = append(rows, Row{
			From: record[0],
			Date: date.Local(),
			Text: record[2],
		})
	}

	table := &Table{
		Rows: rows,
	}

	return table, nil
}
