package pcommon

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"
)

func (tmap *TickMap) FilterInRange(t0 time.Time, t1 time.Time) TickMap {
	ret := make(TickMap)
	for time, tick := range *tmap {
		if time >= t0.Unix() && time < t1.Unix() {
			ret[time] = tick
		}
	}
	return ret
}

func (m *TickMap) DeleteInRange(t0 time.Time, t1 time.Time) {
	for time := range *m {
		if time >= t0.Unix() && time < t1.Unix() {
			delete(*m, time)
		}
	}
}

func (m *TickMap) ToJSON(tick Tick) (string, error) {
	type ITickWithTime struct {
		Tick
		Time int64 `json:"time"`
	}

	var tickArray []ITickWithTime = make([]ITickWithTime, len(*m))

	i := 0
	for time, tick := range *m {
		tickArray[i] = ITickWithTime{
			Tick: tick,
			Time: time,
		}
		i++
	}

	tickArrayJSON, err := json.Marshal(tickArray)
	if err != nil {
		return "", err
	}

	return string(tickArrayJSON), nil
}

func (tick Tick) Stringify() string {
	ret := ""
	ret += strconv.FormatFloat(tick.Open, 'f', -1, 64) + "|"
	ret += strconv.FormatFloat(tick.High, 'f', -1, 64) + "|"
	ret += strconv.FormatFloat(tick.Low, 'f', -1, 64) + "|"
	ret += strconv.FormatFloat(tick.Close, 'f', -1, 64) + "|"
	ret += strconv.FormatFloat(tick.VolumeBought, 'f', -1, 64) + "|"
	ret += strconv.FormatFloat(tick.VolumeSold, 'f', -1, 64) + "|"
	ret += strconv.FormatInt(tick.TradeCount, 10) + "|"
	ret += strconv.FormatFloat(tick.MedianVolumeBought, 'f', -1, 64) + "|"
	ret += strconv.FormatFloat(tick.AverageVolumeBought, 'f', -1, 64) + "|"
	ret += strconv.FormatFloat(tick.MedianVolumeSold, 'f', -1, 64) + "|"
	ret += strconv.FormatFloat(tick.AverageVolumeSold, 'f', -1, 64) + "|"
	ret += strconv.FormatFloat(tick.VWAP, 'f', -1, 64) + "|"
	ret += strconv.FormatFloat(tick.StandardDeviation, 'f', -1, 64)
	return ret
}

func ParseTick(str string) Tick {
	split := strings.Split(str, "|")
	open, _ := strconv.ParseFloat(split[0], 64)
	high, _ := strconv.ParseFloat(split[1], 64)
	low, _ := strconv.ParseFloat(split[2], 64)
	close, _ := strconv.ParseFloat(split[3], 64)
	volumeBought, _ := strconv.ParseFloat(split[4], 64)
	volumeSold, _ := strconv.ParseFloat(split[5], 64)
	tradeCount, _ := strconv.ParseInt(split[6], 10, 64)
	medianVolumeBought, _ := strconv.ParseFloat(split[7], 64)
	averageVolumeBought, _ := strconv.ParseFloat(split[8], 64)
	medianVolumeSold, _ := strconv.ParseFloat(split[9], 64)
	averageVolumeSold, _ := strconv.ParseFloat(split[10], 64)
	vwap, _ := strconv.ParseFloat(split[11], 64)
	standardDeviation, _ := strconv.ParseFloat(split[12], 64)

	return Tick{
		Open:                open,
		High:                high,
		Low:                 low,
		Close:               close,
		VolumeBought:        volumeBought,
		VolumeSold:          volumeSold,
		TradeCount:          tradeCount,
		MedianVolumeBought:  medianVolumeBought,
		AverageVolumeBought: averageVolumeBought,
		MedianVolumeSold:    medianVolumeSold,
		AverageVolumeSold:   averageVolumeSold,
		VWAP:                vwap,
		StandardDeviation:   standardDeviation,
	}
}

func (candles TickArray) AggregateCandlesToCandle() Tick {

	aggregateCandle := Tick{
		Open:                candles[0].Open,
		High:                candles[0].High,
		Low:                 candles[0].Low,
		Close:               candles[len(candles)-1].Close,
		VolumeBought:        0,
		VolumeSold:          0,
		TradeCount:          0,
		MedianVolumeBought:  0,
		AverageVolumeBought: 0,
		MedianVolumeSold:    0,
		AverageVolumeSold:   0,
		VWAP:                0,
		StandardDeviation:   0,
	}

	tradeVolumesBought := []float64{}
	tradeVolumesSold := []float64{}
	for _, c := range candles {
		aggregateCandle.High = math.Max(aggregateCandle.High, c.High)
		aggregateCandle.Low = math.Min(aggregateCandle.Low, c.Low)
		aggregateCandle.VolumeBought += c.VolumeBought
		aggregateCandle.VolumeSold += c.VolumeSold
		aggregateCandle.TradeCount += c.TradeCount

		tradeVolumesBought = append(tradeVolumesBought, c.VolumeBought)
		tradeVolumesSold = append(tradeVolumesSold, c.VolumeSold)
	}

	aggregateCandle.MedianVolumeBought = SafeMedian(tradeVolumesBought)
	aggregateCandle.MedianVolumeSold = SafeMedian(tradeVolumesSold)
	aggregateCandle.AverageVolumeBought = SafeAverage(tradeVolumesBought)
	aggregateCandle.AverageVolumeSold = SafeAverage(tradeVolumesSold)

	aggregateCandle.VWAP = candles.calculateVWAP()
	aggregateCandle.StandardDeviation = CalculateStandardDeviation(append(tradeVolumesBought, tradeVolumesSold...))

	return aggregateCandle
}

func (candles TickArray) calculateVWAP() float64 {
	if len(candles) == 0 {
		return 0.0 // VWAP is not defined if there are no trades.
	}

	var totalVolume float64
	var vwapNumerator float64

	for _, candle := range candles {
		vwapNumerator += candle.VWAP * (candle.VolumeBought + candle.VolumeSold)
		totalVolume += candle.VolumeBought + candle.VolumeSold
	}

	if totalVolume == 0 {
		return 0.0 // Prevent division by zero if total volume is zero.
	}

	vwap := vwapNumerator / totalVolume
	return vwap
}
