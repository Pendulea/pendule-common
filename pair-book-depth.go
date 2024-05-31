package pcommon

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"
)

type BookDepth struct {
	Timestamp  time.Time
	Percentage int
	Depth      float64
	Notional   float64
}

type BookDepthList []BookDepth

func parseBookDepthFromCSVLine(fields []string) (BookDepth, error) {
	var err error
	bd := BookDepth{}

	bd.Timestamp, err = time.ParseInLocation("2006-01-02 15:04:05", fields[0], time.UTC)
	if err != nil {
		return BookDepth{}, err
	}
	bd.Percentage, err = strconv.Atoi(fields[1])
	if err != nil {
		return BookDepth{}, err
	}
	bd.Depth, err = strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return BookDepth{}, err
	}
	bd.Notional, err = strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return BookDepth{}, err
	}

	return bd, nil
}

func (p *Pair) ParseBookDepthFromCSV(date string) (BookDepthList, error) {
	file, err := os.Open(p.BuildBookDepthArchivesFilePath(date, "csv"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // Set the delimiter to comma
	reader.TrimLeadingSpace = true

	var bookDepths BookDepthList

	// Check if the CSV is empty
	firstRow, err := reader.Read()
	if err == io.EOF {
		// CSV is empty, return an empty slice
		return bookDepths, nil
	}
	if err != nil {
		return nil, err
	}

	// Determine if the first row is a header or a data row
	if isHeader(firstRow) {
		// Read the next row if the first row is a header
		firstRow, err = reader.Read()
		if err == io.EOF {
			// CSV only contains a header, return an empty slice
			return bookDepths, nil
		}
		if err != nil {
			return nil, err
		}
	}

	bd, err := parseBookDepthFromCSVLine(firstRow)
	if err != nil {
		return nil, err
	}
	bookDepths = append(bookDepths, bd)

	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		bd, err := parseBookDepthFromCSVLine(fields)
		if err != nil {
			return nil, err
		}
		bookDepths = append(bookDepths, bd)
	}

	return bookDepths, nil
}
