package pcommon

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	"strconv"
)

type TradeType int8

const SPOT_TRADE TradeType = 1
const FUTURES_TRADE TradeType = 2

func (t TradeType) Key() string {
	switch t {
	case SPOT_TRADE:
		return "_spot"
	case FUTURES_TRADE:
		return "_futures"
	}
	log.Fatalf("invalid trade type")
	return ""
}

// Trade represents a trade in the system.
type Trade struct {
	TradeID      int64
	Price        float64
	Quantity     float64
	Total        float64
	Timestamp    int64
	IsBuyerMaker bool
	IsBestMatch  bool
}

type TradeList []Trade

func (tradeType TradeType) parseTradeFromCSVLine(fields []string) (Trade, error) {
	var err error
	trade := Trade{}

	trade.TradeID, err = strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return Trade{}, err
	}
	trade.Price, err = strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return Trade{}, err
	}
	trade.Quantity, err = strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return Trade{}, err
	}
	trade.Total, err = strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return Trade{}, err
	}
	trade.Timestamp, err = strconv.ParseInt(fields[4], 10, 64)
	if err != nil {
		return Trade{}, err
	}
	trade.IsBuyerMaker = fields[5] == "True"
	trade.IsBestMatch = When[bool](tradeType == SPOT_TRADE).Then(fields[6] == "True").Else(false)
	return trade, nil
}

func (p *Pair) ParseTradesFromCSV(date string) (TradeList, error) {
	file, err := os.Open(p.BuildTradesArchivesFilePath(date, "csv"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // Set the delimiter to comma
	reader.TrimLeadingSpace = true

	var trades TradeList

	// Check if the CSV is empty
	firstRow, err := reader.Read()
	if err == io.EOF {
		// CSV is empty, return an empty slice
		return trades, nil
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
			return trades, nil
		}
		if err != nil {
			return nil, err
		}
	}

	trade, err := p.TradeType().parseTradeFromCSVLine(firstRow)
	if err != nil {
		return nil, err
	}
	trades = append(trades, trade)

	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		trade, err := p.TradeType().parseTradeFromCSVLine(fields)
		if err != nil {
			return nil, err
		}
		trades = append(trades, trade)

		// Convert fields to a trade and append to trades
	}

	return trades, nil
}

// Example function to determine if a row is a header
func isHeader(row []string) bool {
	for _, field := range row {
		if strings.Contains(field, "time") || strings.Contains(field, "date") || strings.Contains(field, "id") {
			return true
		}
	}
	return false
}
