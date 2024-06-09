package pcommon

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"
)

type FuturesMetrics struct {
	CreateTime                   time.Time `json:"create_time"`
	Symbol                       string    `json:"symbol"`
	SumOpenInterest              float64   `json:"sum_open_interest"`
	SumOpenInterestValue         float64   `json:"sum_open_interest_value"`
	CountTopTraderLongShortRatio float64   `json:"count_toptrader_long_short_ratio"`
	SumTopTraderLongShortRatio   float64   `json:"sum_toptrader_long_short_ratio"`
	CountLongShortRatio          float64   `json:"count_long_short_ratio"`
	SumTakerLongShortVolRatio    float64   `json:"sum_taker_long_short_vol_ratio"`
}

type FuturesMetricsList []FuturesMetrics

func parseMetricsFromCSVLine(fields []string) (FuturesMetrics, error) {
	var err error
	m := FuturesMetrics{}

	m.CreateTime, err = time.ParseInLocation("2006-01-02 15:04:05", fields[0], time.UTC)
	if err != nil {
		return FuturesMetrics{}, err
	}
	m.CreateTime.Truncate(time.Minute)
	m.Symbol = fields[1]
	if fields[2] != "" {
		m.SumOpenInterest, err = strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return FuturesMetrics{}, err
		}
	} else {
		m.SumOpenInterest = 0
	}

	if fields[3] != "" {
		m.SumOpenInterestValue, err = strconv.ParseFloat(fields[3], 64)
		if err != nil {
			return FuturesMetrics{}, err
		}
	} else {
		m.SumOpenInterestValue = 0
	}

	if fields[4] != "" {
		m.CountTopTraderLongShortRatio, err = strconv.ParseFloat(fields[4], 64)
		if err != nil {
			return FuturesMetrics{}, err
		}
	} else {
		m.CountTopTraderLongShortRatio = 0
	}

	if fields[5] != "" {
		m.SumTopTraderLongShortRatio, err = strconv.ParseFloat(fields[5], 64)
		if err != nil {
			return FuturesMetrics{}, err
		}
	} else {
		m.SumTopTraderLongShortRatio = 0
	}

	if fields[6] != "" {
		m.CountLongShortRatio, err = strconv.ParseFloat(fields[6], 64)
		if err != nil {
			return FuturesMetrics{}, err
		}
	} else {
		m.CountLongShortRatio = 0
	}

	if fields[7] != "" {
		m.SumTakerLongShortVolRatio, err = strconv.ParseFloat(fields[7], 64)
		if err != nil {
			return FuturesMetrics{}, err
		}
	} else {
		m.SumTakerLongShortVolRatio = 0
	}

	return m, nil
}

func (p *Pair) ParseMetricsFromCSV(date string) (FuturesMetricsList, error) {
	file, err := os.Open(p.BuildFuturesMetricsArchivesFilePath(date, "csv"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // Set the delimiter to comma
	reader.TrimLeadingSpace = true

	var metrics FuturesMetricsList

	// Check if the CSV is empty
	firstRow, err := reader.Read()
	if err == io.EOF {
		// CSV is empty, return an empty slice
		return metrics, nil
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
			return metrics, nil
		}
		if err != nil {
			return nil, err
		}
	}

	m, err := parseMetricsFromCSVLine(firstRow)
	if err != nil {
		return nil, err
	}
	metrics = append(metrics, m)

	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		m, err := parseMetricsFromCSVLine(fields)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}

	return metrics, nil
}
