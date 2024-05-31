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
	// Read the header row (and ignore it if necessary)
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

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
