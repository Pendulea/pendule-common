package pcommon

import (
	"strconv"
	"strings"
	"time"
)

// the key is the percentage of the book depth
type FullBookDepthTick map[int]SingleBookDepth
type FullBookDepthTickTime map[int]SingleBookDepthTime

type SingleBookDepth struct {
	Percent int `json:"percent"`

	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Avg    float64 `json:"avg"`
	Median float64 `json:"median"`
	Count  int     `json:"count"`
}

type SingleBookDepthTime struct {
	SingleBookDepth
	Time int64 `json:"time"`
}

func newEmptyFullBookDepthTickTime(time int64) *FullBookDepthTickTime {
	ret := make(FullBookDepthTickTime, 10)
	for i := -5; i <= 5; i++ {
		if i != 0 {
			ret[i] = SingleBookDepthTime{
				Time: time,
				SingleBookDepth: SingleBookDepth{
					Percent: i,
					Count:   0,
				},
			}
		}
	}
	return &ret
}

// the key is the timestamp
type FullBookDepthTickMap map[int64]FullBookDepthTick

func (fdbtp *FullBookDepthTickMap) Keys(sortAsc bool) []int64 {
	keys := make([]int64, len(*fdbtp))
	i := 0
	for k := range *fdbtp {
		keys[i] = k
		i++
	}
	return Sort[int64](keys, !sortAsc)
}

// returns the percentaages of the book depth sorted in ascending order
func (fdbtp *FullBookDepthTick) Keys() []int {
	keys := make([]int, len(*fdbtp))
	i := 0
	for k := range *fdbtp {
		keys[i] = k
		i++
	}
	return Sort[int](keys, false)
}

func (t *SingleBookDepth) ToTime(time int64) SingleBookDepthTime {
	return SingleBookDepthTime{
		SingleBookDepth: *t,
		Time:            time,
	}
}

func (t *FullBookDepthTick) ToTime(time int64) FullBookDepthTickTime {
	ret := make(FullBookDepthTickTime, len(*t))
	for t, sbd := range *t {
		ret[t] = sbd.ToTime(time)
	}
	return ret
}

func (fbdtm *FullBookDepthTickMap) ToTime() []FullBookDepthTickTime {
	ret := make([]FullBookDepthTickTime, len(*fbdtm))

	i := 0
	for time, fbdt := range *fbdtm {
		ret[i] = fbdt.ToTime(time)
		i++
	}
	return ret
}

func (fbdtt *FullBookDepthTickTime) IsFilled() bool {
	if len(*fbdtt) != 10 {
		return false
	}

	for _, sbdt := range *fbdtt {
		if sbdt.Count == 0 {
			return false
		}
	}

	return true
}

func (fbdtt *FullBookDepthTickTime) Time() time.Time {
	for _, sbdt := range *fbdtt {
		if sbdt.Time > 0 {
			return time.Unix(sbdt.Time, 0)
		}
	}
	return time.Unix(0, 0)
}

func (fbd FullBookDepthTick) stringify(decimals int8) string {
	ret := make([]string, 10)
	keys := fbd.Keys()
	for i, k := range keys {
		ret[i] = fbd[k].stringify(decimals)
	}

	return strings.Join(ret, "@")
}

func parseFullBookDepthTick(str string) FullBookDepthTick {
	ret := make(FullBookDepthTick, 10)
	split := strings.Split(str, "@")
	for i, s := range split {
		ret[i] = parseSingleBookDepthTick(s)
	}
	return ret
}

func (tick SingleBookDepth) stringify(decimals int8) string {
	ret := ""
	if tick.Count == 1 {
		ret += strconv.Itoa(tick.Percent) + "|"
		ret += Format.Float(tick.Open, decimals) + "|"
	} else if tick.Count > 1 {
		ret += strconv.Itoa(tick.Percent) + "|"
		ret += Format.Float(tick.Open, decimals) + "|"
		ret += Format.Float(tick.High, decimals) + "|"
		ret += Format.Float(tick.Low, decimals) + "|"
		ret += Format.Float(tick.Close, decimals) + "|"
		ret += Format.Float(tick.Avg, decimals) + "|"
		ret += Format.Float(tick.Median, decimals) + "|"
		ret += strconv.Itoa(tick.Count) + "|"
	}
	return ret
}

func parseSingleBookDepthTick(str string) SingleBookDepth {
	split := ChunkString(str, 2)
	percent, _ := strconv.Atoi(split[0])
	open, _ := strconv.ParseFloat(split[1], 64)

	if len(split) == 2 {
		return SingleBookDepth{
			Percent: percent,
			Open:    open,
			High:    open,
			Low:     open,
			Close:   open,
			Avg:     open,
			Median:  open,
			Count:   1,
		}
	}

	high, _ := strconv.ParseFloat(split[2], 64)
	low, _ := strconv.ParseFloat(split[3], 64)
	close, _ := strconv.ParseFloat(split[4], 64)
	avg, _ := strconv.ParseFloat(split[5], 64)
	median, _ := strconv.ParseFloat(split[6], 64)
	count, _ := strconv.Atoi(split[7])

	return SingleBookDepth{
		Percent: percent,
		Open:    open,
		High:    high,
		Low:     low,
		Close:   close,
		Avg:     avg,
		Median:  median,
		Count:   count,
	}
}
