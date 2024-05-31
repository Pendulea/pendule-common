package pcommon

import "strconv"

type BookDepthTick struct {
	Percent int `json:"percent"`

	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Avg    float64 `json:"avg"`
	Median float64 `json:"median"`
	Count  int     `json:"count"`
}

type BookDepthTickMap map[int64]BookDepthTick

type BookDepthTickTime struct {
	BookDepthTick
	Time int64 `json:"time"`
}

type BookDepthTickTimeArray []BookDepthTickTime

func (t *BookDepthTick) ToTickTime(time int64) BookDepthTickTime {
	return BookDepthTickTime{
		BookDepthTick: *t,
		Time:          time,
	}
}

func (t *BookDepthTickTime) ToTick() BookDepthTick {
	return t.BookDepthTick
}

func (tmap *BookDepthTickMap) ToTickTimeArray() *BookDepthTickTimeArray {
	tickTimeArray := make(BookDepthTickTimeArray, len(*tmap))
	i := 0
	for time, tick := range *tmap {
		tickTimeArray[i] = tick.ToTickTime(time)
		i++
	}
	return &tickTimeArray
}

func (tta *BookDepthTickTimeArray) Sort(asc bool) BookDepthTickTimeArray {
	if asc {
		ret := make(BookDepthTickTimeArray, len(*tta))
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
		ret := make(BookDepthTickTimeArray, len(*tta))
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

func (tick BookDepthTick) Stringify(decimals int8) string {
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

func ParseBookDepthTick(str string) BookDepthTick {
	split := ChunkString(str, 2)
	percent, _ := strconv.Atoi(split[0])
	open, _ := strconv.ParseFloat(split[1], 64)

	if len(split) == 2 {
		return BookDepthTick{
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

	return BookDepthTick{
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
