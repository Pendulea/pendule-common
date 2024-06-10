package pcommon

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Metric struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Avg    float64 `json:"avg"`
	Median float64 `json:"median"`
	Count  int     `json:"count"`
}

func aggregateMetrics(list []Metric) Metric {
	ret := Metric{}
	closes := []float64{}
	for _, metric := range list {
		if metric.Count == 0 || metric.Open <= 0.00 {
			continue
		}
		ret.Open = metric.Open
		ret.High = math.Max(ret.High, metric.High)
		if ret.Low <= 0.00 {
			ret.Low = metric.Low
		} else {
			ret.Low = math.Min(ret.Low, metric.Low)
		}
		ret.Close = metric.Close
		ret.Count += metric.Count
		closes = append(closes, metric.Close)
	}

	ret.Avg = Math.SafeAverage(closes)
	ret.Median = Math.SafeMedian(closes)
	return ret
}

func (m *Metric) IsEmpty() bool {
	return m.Count == 0
}

func newMetric(v float64) Metric {
	if v <= 0.00 {
		return Metric{}
	}

	return Metric{
		Open:   v,
		High:   v,
		Low:    v,
		Close:  v,
		Avg:    v,
		Median: v,
		Count:  1,
	}
}

type Metrics struct {
	SumOpenInterest              Metric `json:"sum_open_interest"`
	CountTopTraderLongShortRatio Metric `json:"count_toptrader_long_short_ratio"`
	SumTopTraderLongShortRatio   Metric `json:"sum_toptrader_long_short_ratio"`
	CountLongShortRatio          Metric `json:"count_long_short_ratio"`
	SumTakerLongShortVolRatio    Metric `json:"sum_taker_long_short_vol_ratio"`
}

func NewMetrics(pm FuturesMetrics) Metrics {
	return Metrics{
		SumOpenInterest:              newMetric(pm.SumOpenInterest),
		CountTopTraderLongShortRatio: newMetric(pm.CountTopTraderLongShortRatio),
		SumTopTraderLongShortRatio:   newMetric(pm.SumTopTraderLongShortRatio),
		CountLongShortRatio:          newMetric(pm.CountLongShortRatio),
		SumTakerLongShortVolRatio:    newMetric(pm.SumTakerLongShortVolRatio),
	}
}

type MetricsTime struct {
	Metrics
	Time TimeUnit `json:"time"`
}

type MetricsMap map[TimeUnit]Metrics

type MetricsArray []MetricsTime

func (fdbtp *MetricsMap) Keys(sortAsc bool) []TimeUnit {
	keys := make([]TimeUnit, len(*fdbtp))
	i := 0
	for k := range *fdbtp {
		keys[i] = k
		i++
	}
	return Sort(keys, !sortAsc)
}

func (m *Metrics) ToMetricsTime(time TimeUnit) MetricsTime {
	return MetricsTime{
		Metrics: *m,
		Time:    time,
	}
}

func (mm *MetricsMap) ToMetricsArray() *MetricsArray {
	mmArray := make(MetricsArray, len(*mm))
	i := 0
	for time, metric := range *mm {
		mmArray[i] = metric.ToMetricsTime(time)
		i++
	}
	return &mmArray
}

func (t *MetricsTime) ToMetrics() Metrics {
	return t.Metrics
}

func (m *Metric) stringify(decimals int8) string {
	if m.Count == 1 {
		return Format.Float(m.Open, decimals)
	}
	return fmt.Sprintf("%s@%s@%s@%s@%s@%s@%d", Format.Float(m.Open, decimals), Format.Float(m.High, decimals), Format.Float(m.Low, decimals), Format.Float(m.Close, decimals), Format.Float(m.Avg, decimals), Format.Float(m.Median, decimals), m.Count)
}

func parseMetric(s string) Metric {
	splited := strings.Split(s, "@")
	if len(splited) == 1 {
		v, err := strconv.ParseFloat(splited[0], 64)
		if err != nil {
			return Metric{}
		}
		return newMetric(v)
	}
	if len(splited) != 7 {
		return Metric{}
	}
	open, _ := strconv.ParseFloat(splited[0], 64)
	high, _ := strconv.ParseFloat(splited[1], 64)
	low, _ := strconv.ParseFloat(splited[2], 64)
	close, _ := strconv.ParseFloat(splited[3], 64)
	avg, _ := strconv.ParseFloat(splited[4], 64)
	median, _ := strconv.ParseFloat(splited[5], 64)
	count, _ := strconv.Atoi(splited[6])

	return Metric{
		Open:   open,
		High:   high,
		Low:    low,
		Close:  close,
		Avg:    avg,
		Median: median,
		Count:  count,
	}
}

func (m *Metrics) Stringify(decimals int8) string {
	return fmt.Sprintf("%s|%s|%s|%s|%s", m.SumOpenInterest.stringify(decimals), m.CountTopTraderLongShortRatio.stringify(decimals), m.SumTopTraderLongShortRatio.stringify(decimals), m.CountLongShortRatio.stringify(decimals), m.SumTakerLongShortVolRatio.stringify(decimals))
}

func ParseMetrics(str string) *Metrics {
	splited := strings.Split(str, "|")
	if len(splited) != 5 {
		return nil
	}
	return &Metrics{
		SumOpenInterest:              parseMetric(splited[0]),
		CountTopTraderLongShortRatio: parseMetric(splited[1]),
		SumTopTraderLongShortRatio:   parseMetric(splited[2]),
		CountLongShortRatio:          parseMetric(splited[3]),
		SumTakerLongShortVolRatio:    parseMetric(splited[4]),
	}
}

func (t *Metric) OpenString(decimals int8) string {
	return Format.Float(t.Open, decimals)
}

func (t *Metric) HighString(decimals int8) string {
	return Format.Float(t.High, decimals)
}

func (t *Metric) LowString(decimals int8) string {
	return Format.Float(t.Low, decimals)
}

func (t *Metric) CloseString(decimals int8) string {
	return Format.Float(t.Close, decimals)
}

func (t *Metric) AvgString(decimals int8) string {
	return Format.Float(t.Avg, decimals)
}

func (t *Metric) MedianString(decimals int8) string {
	return Format.Float(t.Median, decimals)
}

func (t *Metric) CountString() string {
	return strconv.Itoa(t.Count)
}

func (list MetricsArray) Aggregate() Metrics {

	CountLongShortRatios := make([]Metric, len(list))
	CountTopTraderLongShortRatios := make([]Metric, len(list))
	SumOpenInterests := make([]Metric, len(list))
	SumTopTraderLongShortRatios := make([]Metric, len(list))
	SumTakerLongShortVolRatios := make([]Metric, len(list))

	for i, m := range list {
		CountLongShortRatios[i] = m.CountLongShortRatio
		CountTopTraderLongShortRatios[i] = m.CountTopTraderLongShortRatio
		SumOpenInterests[i] = m.SumOpenInterest
		SumTopTraderLongShortRatios[i] = m.SumTopTraderLongShortRatio
		SumTakerLongShortVolRatios[i] = m.SumTakerLongShortVolRatio
	}

	return Metrics{
		CountLongShortRatio:          aggregateMetrics(CountLongShortRatios),
		CountTopTraderLongShortRatio: aggregateMetrics(CountTopTraderLongShortRatios),
		SumOpenInterest:              aggregateMetrics(SumOpenInterests),
		SumTopTraderLongShortRatio:   aggregateMetrics(SumTopTraderLongShortRatios),
		SumTakerLongShortVolRatio:    aggregateMetrics(SumTakerLongShortVolRatios),
	}
}
