package pcommon

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Tick struct {
	Open                float64                `json:"open"`
	High                float64                `json:"high"`
	Low                 float64                `json:"low"`
	Close               float64                `json:"close"`
	VolumeBought        float64                `json:"volume_bought"`
	VolumeSold          float64                `json:"volume_sold"`
	TradeCount          int64                  `json:"trade_count"`
	MedianVolumeBought  float64                `json:"median_volume_bought"`
	AverageVolumeBought float64                `json:"average_volume_bought"`
	MedianVolumeSold    float64                `json:"median_volume_sold"`
	AverageVolumeSold   float64                `json:"average_volume_sold"`
	VWAP                float64                `json:"vwap"`
	StandardDeviation   float64                `json:"standard_deviation"`
	AbsolutePriceSum    float64                `json:"absolute_price_sum"`
	PrevBookDepth       *FullBookDepthTickTime `json:"prev_book_depth"`
}

func (t *Tick) OpenString() string {
	return Format.Float(t.Open, -1)
}

func (t *Tick) HighString() string {
	return Format.Float(t.High, -1)
}

func (t *Tick) LowString() string {
	return Format.Float(t.Low, -1)
}

func (t *Tick) CloseString() string {
	return Format.Float(t.Close, -1)
}

func (t *Tick) AbsolutePriceSumString() string {
	return Format.Float(t.AbsolutePriceSum, -1)
}

func (t *Tick) VolumeBoughtString(decimals int8) string {
	return Format.Float(t.VolumeBought, decimals)
}

func (t *Tick) VolumeSoldString(decimals int8) string {
	return Format.Float(t.VolumeSold, decimals)
}

func (t *Tick) TradeCountString() string {
	return strconv.FormatInt(t.TradeCount, 10)
}

func (t *Tick) MedianVolumeBoughtString(decimals int8) string {
	return Format.Float(t.MedianVolumeBought, decimals)
}

func (t *Tick) AverageVolumeBoughtString(decimals int8) string {
	return Format.Float(t.AverageVolumeBought, decimals)
}

func (t *Tick) MedianVolumeSoldString(decimals int8) string {
	return Format.Float(t.MedianVolumeSold, decimals)
}

func (t *Tick) AverageVolumeSoldString(decimals int8) string {
	return Format.Float(t.AverageVolumeSold, decimals)
}

func (t *Tick) VWAPString() string {
	return Format.Float(t.VWAP, 5)
}

func (t *Tick) StandardDeviationString() string {
	return Format.Float(t.StandardDeviation, 3)
}

type TickMap map[TimeUnit]Tick

type TickTime struct {
	Tick
	Time TimeUnit `json:"time"`
}

type TickTimeArray []TickTime

func (t *Tick) ToTickTime(time TimeUnit) TickTime {
	return TickTime{
		Tick: *t,
		Time: time,
	}
}

func (t *TickTime) ToTick() Tick {
	return t.Tick
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
		if time.ToTime().UnixNano() >= t0.UnixNano() && time.ToTime().UnixNano() < t1.UnixNano() {
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
		if time.ToTime().UnixNano() >= t0.UnixNano() && time.ToTime().UnixNano() < t1.UnixNano() {
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
	ret += Format.Float(tick.Open, -1) + "|"
	ret += Format.Float(tick.High, -1) + "|"
	ret += Format.Float(tick.Low, -1) + "|"
	ret += Format.Float(tick.Close, -1) + "|"
	ret += Format.Float(tick.VolumeBought, decimals) + "|"
	ret += Format.Float(tick.VolumeSold, decimals) + "|"
	ret += strconv.FormatInt(tick.TradeCount, 10) + "|"
	ret += Format.Float(tick.MedianVolumeBought, decimals) + "|"
	ret += Format.Float(tick.AverageVolumeBought, decimals) + "|"
	ret += Format.Float(tick.MedianVolumeSold, decimals) + "|"
	ret += Format.Float(tick.AverageVolumeSold, decimals) + "|"
	ret += Format.Float(tick.VWAP, 5) + "|"
	ret += Format.Float(tick.StandardDeviation, 3) + "|"
	ret += Format.Float(tick.AbsolutePriceSum, -1)
	if tick.PrevBookDepth != nil {
		ret += "|" + tick.PrevBookDepth.Time().String()
	}

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

	var pbd *FullBookDepthTickTime = nil
	if len(split) == 15 {
		prevBookDepthTime, _ := strconv.ParseInt(split[14], 10, 64)
		pbd = newEmptyFullBookDepthTickTime(NewTimeUnit(prevBookDepthTime))
	}

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
		PrevBookDepth:       pbd,
	}
}

func (candles TickTimeArray) CalculateVWAP() float64 {
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
