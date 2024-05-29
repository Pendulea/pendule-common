package pcommon

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"
)

type Tick struct {
	Open                float64 `json:"open"`
	High                float64 `json:"high"`
	Low                 float64 `json:"low"`
	Close               float64 `json:"close"`
	VolumeBought        float64 `json:"volume_bought"`
	VolumeSold          float64 `json:"volume_sold"`
	TradeCount          int64   `json:"trade_count"`
	MedianVolumeBought  float64 `json:"median_volume_bought"`
	AverageVolumeBought float64 `json:"average_volume_bought"`
	MedianVolumeSold    float64 `json:"median_volume_sold"`
	AverageVolumeSold   float64 `json:"average_volume_sold"`
	VWAP                float64 `json:"vwap"`
	StandardDeviation   float64 `json:"standard_deviation"`
	AbsolutePriceSum    float64 `json:"absolute_price_sum"`
}

func formatFloat(val float64, precision int8) string {
	// Format the float with the specified precision
	str := strconv.FormatFloat(val, 'f', int(precision), 64)

	// Remove trailing zeros and the decimal point if it's not needed
	str = strings.TrimRight(str, "0")
	str = strings.TrimRight(str, ".")

	return str
}

func (t *Tick) OpenString() string {
	return formatFloat(t.Open, -1)
}

func (t *Tick) HighString() string {
	return formatFloat(t.High, -1)
}

func (t *Tick) LowString() string {
	return formatFloat(t.Low, -1)
}

func (t *Tick) CloseString() string {
	return formatFloat(t.Close, -1)
}

func (t *Tick) AbsolutePriceSumString() string {
	return formatFloat(t.AbsolutePriceSum, -1)
}

func (t *Tick) VolumeBoughtString(decimals int8) string {
	if t.VolumeBought == 0 {
		return "0"
	}
	return formatFloat(t.VolumeBought, decimals)
}

func (t *Tick) VolumeSoldString(decimals int8) string {
	if t.VolumeSold == 0 {
		return "0"
	}
	return formatFloat(t.VolumeSold, decimals)
}

func (t *Tick) TradeCountString() string {
	return strconv.FormatInt(t.TradeCount, 10)
}

func (t *Tick) MedianVolumeBoughtString(decimals int8) string {
	if t.MedianVolumeBought == 0 {
		return "0"
	}
	return formatFloat(t.MedianVolumeBought, decimals)
}

func (t *Tick) AverageVolumeBoughtString(decimals int8) string {
	if t.AverageVolumeBought == 0 {
		return "0"
	}
	return formatFloat(t.AverageVolumeBought, decimals)
}

func (t *Tick) MedianVolumeSoldString(decimals int8) string {
	if t.MedianVolumeSold == 0 {
		return "0"
	}
	return formatFloat(t.MedianVolumeSold, decimals)
}

func (t *Tick) AverageVolumeSoldString(decimals int8) string {
	if t.AverageVolumeSold == 0 {
		return "0"
	}
	return formatFloat(t.AverageVolumeSold, decimals)
}

func (t *Tick) VWAPString() string {
	if t.VWAP == 0 {
		return "0"
	}

	return formatFloat(t.VWAP, 5)
}

func (t *Tick) StandardDeviationString() string {
	if t.StandardDeviation == 0 {
		return "0"
	}
	return strconv.FormatFloat(t.StandardDeviation, 'f', 3, 64)
}

type TickMap map[int64]Tick

type TickTime struct {
	Open                float64 `json:"open"`
	High                float64 `json:"high"`
	Low                 float64 `json:"low"`
	Close               float64 `json:"close"`
	VolumeBought        float64 `json:"volume_bought"`
	VolumeSold          float64 `json:"volume_sold"`
	TradeCount          int64   `json:"trade_count"`
	MedianVolumeBought  float64 `json:"median_volume_bought"`
	AverageVolumeBought float64 `json:"average_volume_bought"`
	MedianVolumeSold    float64 `json:"median_volume_sold"`
	AverageVolumeSold   float64 `json:"average_volume_sold"`
	VWAP                float64 `json:"vwap"`
	StandardDeviation   float64 `json:"standard_deviation"`
	AbsolutePriceSum    float64 `json:"absolute_price_sum"`
	AbsolumeVolumeSum   float64 `json:"absolute_volume_sum"`
	Time                int64   `json:"time"`
}

type TickTimeArray []TickTime

func (t *Tick) ToTickTime(time int64) TickTime {
	return TickTime{
		Open:                t.Open,
		High:                t.High,
		Low:                 t.Low,
		Close:               t.Close,
		VolumeBought:        t.VolumeBought,
		VolumeSold:          t.VolumeSold,
		TradeCount:          t.TradeCount,
		MedianVolumeBought:  t.MedianVolumeBought,
		AverageVolumeBought: t.AverageVolumeBought,
		MedianVolumeSold:    t.MedianVolumeSold,
		AverageVolumeSold:   t.AverageVolumeSold,
		VWAP:                t.VWAP,
		StandardDeviation:   t.StandardDeviation,
		AbsolutePriceSum:    t.AbsolutePriceSum,
		Time:                time,
	}
}

func (t *TickTime) ToTick() Tick {
	return Tick{
		Open:                t.Open,
		High:                t.High,
		Low:                 t.Low,
		Close:               t.Close,
		VolumeBought:        t.VolumeBought,
		VolumeSold:          t.VolumeSold,
		TradeCount:          t.TradeCount,
		MedianVolumeBought:  t.MedianVolumeBought,
		AverageVolumeBought: t.AverageVolumeBought,
		MedianVolumeSold:    t.MedianVolumeSold,
		AverageVolumeSold:   t.AverageVolumeSold,
		VWAP:                t.VWAP,
		StandardDeviation:   t.StandardDeviation,
		AbsolutePriceSum:    t.AbsolutePriceSum,
	}
}

func (tmap *TickMap) ToTickTimeArray() *TickTimeArray {
	tickTimeArray := make(TickTimeArray, len(*tmap))
	i := 0
	for time, tick := range *tmap {
		tickTimeArray[i] = tick.ToTickTime(time)
		i++
	}
	return &tickTimeArray
}

func (tta *TickTimeArray) Sort(asc bool) TickTimeArray {
	if asc {
		ret := make(TickTimeArray, len(*tta))
		copy(ret, *tta)
		for i := 0; i < len(ret); i++ {
			for j := i + 1; j < len(ret); j++ {
				if ret[i].Time > ret[j].Time {
					ret[i], ret[j] = ret[j], ret[i]
				}
			}
		}
		*tta = ret
	} else {
		ret := make(TickTimeArray, len(*tta))
		copy(ret, *tta)
		for i := 0; i < len(ret); i++ {
			for j := i + 1; j < len(ret); j++ {
				if ret[i].Time < ret[j].Time {
					ret[i], ret[j] = ret[j], ret[i]
				}
			}
		}
		*tta = ret
	}
	return *tta
}

func (tmap *TickMap) FilterInRange(t0 time.Time, t1 time.Time) TickMap {
	ret := make(TickMap)
	for time, tick := range *tmap {
		if time >= t0.Unix() && time < t1.Unix() {
			ret[time] = tick
		}
	}
	return ret
}

func (tmap *TickMap) Merge(t TickMap) TickMap {
	for time, tick := range t {
		(*tmap)[time] = tick
	}
	return *tmap
}

func (m *TickMap) DeleteInRange(t0 time.Time, t1 time.Time) {
	for time := range *m {
		if time >= t0.Unix() && time < t1.Unix() {
			delete(*m, time)
		}
	}
}

func (list *TickTime) ToJSON(tick Tick) (string, error) {
	tickArrayJSON, err := json.Marshal(*list)
	if err != nil {
		return "", err
	}

	return string(tickArrayJSON), nil
}

func (tick Tick) Stringify(decimals int8) string {
	ret := ""
	ret += formatFloat(tick.Open, -1) + "|"
	ret += formatFloat(tick.High, -1) + "|"
	ret += formatFloat(tick.Low, -1) + "|"
	ret += formatFloat(tick.Close, -1) + "|"
	ret += formatFloat(tick.VolumeBought, decimals) + "|"
	ret += formatFloat(tick.VolumeSold, decimals) + "|"
	ret += strconv.FormatInt(tick.TradeCount, 10) + "|"
	ret += formatFloat(tick.MedianVolumeBought, decimals) + "|"
	ret += formatFloat(tick.AverageVolumeBought, decimals) + "|"
	ret += formatFloat(tick.MedianVolumeSold, decimals) + "|"
	ret += formatFloat(tick.AverageVolumeSold, decimals) + "|"
	ret += formatFloat(tick.VWAP, 5) + "|"
	ret += formatFloat(tick.StandardDeviation, 3)
	ret += formatFloat(tick.AbsolutePriceSum, -1) + "|"
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
	absolutePriceSum, _ := strconv.ParseFloat(split[13], 64)

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
		AbsolutePriceSum:    absolutePriceSum,
	}
}

func (candles TickTimeArray) AggregateCandlesToCandle() Tick {

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
		AbsolutePriceSum:    0,
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
		aggregateCandle.AbsolutePriceSum += c.AbsolutePriceSum
	}

	aggregateCandle.MedianVolumeBought = Math.SafeMedian(tradeVolumesBought)
	aggregateCandle.MedianVolumeSold = Math.SafeMedian(tradeVolumesSold)
	aggregateCandle.AverageVolumeBought = Math.SafeAverage(tradeVolumesBought)
	aggregateCandle.AverageVolumeSold = Math.SafeAverage(tradeVolumesSold)

	aggregateCandle.VWAP = candles.calculateVWAP()
	aggregateCandle.StandardDeviation = Math.CalculateStandardDeviation(append(tradeVolumesBought, tradeVolumesSold...))

	return aggregateCandle
}

func (candles TickTimeArray) calculateVWAP() float64 {
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
