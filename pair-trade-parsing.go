package pcommon

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"time"

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

func (trades TradeList) AggregateTradesToCandles(timeframe time.Duration) TickMap {
	buckets := make(TickMap)

	candleTrades := []Trade{}

	for _, trade := range trades {
		timeBucket := (trade.Timestamp / timeframe.Milliseconds()) * int64(timeframe.Seconds())
		tick, exists := buckets[timeBucket]
		if !exists {
			// candleTrades = [trade];
			candleTrades = make([]Trade, 0)
			candleTrades = append(candleTrades, trade)
			// Create new candle
			buckets[timeBucket] = Tick{
				Open:                trade.Price,
				High:                trade.Price,
				Low:                 trade.Price,
				Close:               trade.Price,
				VolumeBought:        conditionalQuantity(trade.IsBuyerMaker, trade.Quantity),
				VolumeSold:          conditionalQuantity(!trade.IsBuyerMaker, trade.Quantity),
				TradeCount:          1,
				MedianVolumeBought:  conditionalQuantity(trade.IsBuyerMaker, trade.Quantity),
				AverageVolumeBought: conditionalQuantity(trade.IsBuyerMaker, trade.Quantity),
				MedianVolumeSold:    conditionalQuantity(!trade.IsBuyerMaker, trade.Quantity),
				AverageVolumeSold:   conditionalQuantity(!trade.IsBuyerMaker, trade.Quantity),
				VWAP:                trade.Price,
				StandardDeviation:   0.0,
			}
			continue
		} else {
			candleTrades = append(candleTrades, trade)
			// Update existing candle
			updateTick(&tick, trade, candleTrades)
		}
		buckets[timeBucket] = tick
	}

	return buckets
}

func conditionalQuantity(condition bool, quantity float64) float64 {
	if condition {
		return quantity
	}
	return 0
}

func updateTick(candle *Tick, currentTrades Trade, cumulatedTrades []Trade) {
	candle.High = math.Max(candle.High, currentTrades.Price)
	candle.Low = math.Min(candle.Low, currentTrades.Price)
	candle.Close = currentTrades.Price
	candle.VolumeBought += conditionalQuantity(currentTrades.IsBuyerMaker, currentTrades.Quantity)
	candle.VolumeSold += conditionalQuantity(!currentTrades.IsBuyerMaker, currentTrades.Quantity)
	candle.TradeCount += 1

	tradeVolumesBought := []float64{}
	tradeVolumesSold := []float64{}
	for _, t := range cumulatedTrades {
		if t.IsBuyerMaker {
			tradeVolumesBought = append(tradeVolumesBought, t.Quantity)
		} else {
			tradeVolumesSold = append(tradeVolumesSold, t.Quantity)
		}
	}

	candle.MedianVolumeBought = Math.SafeMedian(tradeVolumesBought)
	candle.MedianVolumeSold = Math.SafeMedian(tradeVolumesSold)
	candle.AverageVolumeBought = Math.SafeAverage(tradeVolumesBought)
	candle.AverageVolumeSold = Math.SafeAverage(tradeVolumesSold)

	candle.VWAP = calculateVWAP(cumulatedTrades)
	candle.StandardDeviation = Math.CalculateStandardDeviation(append(tradeVolumesBought, tradeVolumesSold...))
}

func calculateVWAP(trades []Trade) float64 {
	if len(trades) == 0 {
		return 0.0 // VWAP is not defined if there are no trades.
	}

	var totalVolume float64
	var vwapNumerator float64

	for _, trade := range trades {
		vwapNumerator += trade.Price * trade.Quantity
		totalVolume += trade.Quantity
	}

	if totalVolume == 0 {
		return 0.0 // Prevent division by zero if total volume is zero.
	}

	vwap := vwapNumerator / totalVolume
	return vwap
}

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

func (p *Pair) ParseTradesFromCSV(date string) ([]Trade, error) {
	file, err := os.Open(p.BuildArchivesFilePath(date, "csv"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // Set the delimiter to comma
	reader.TrimLeadingSpace = true

	var trades []Trade

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

		trade, err := p.TradeType().parseTradeFromCSVLine(fields)
		if err != nil {
			fmt.Printf("Error parsing line: %v\n", err)
			continue // or return nil, err to stop processing on error
		}
		trades = append(trades, trade)
	}

	return trades, nil
}
