package pcommon

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type Unit struct {
	Open        float64 `json:"open"`
	High        float64 `json:"high"`
	Low         float64 `json:"low"`
	Close       float64 `json:"close"`
	Average     float64 `json:"average"`
	Median      float64 `json:"median"`
	AbsoluteSum float64 `json:"absolute_sum"`
	Count       int64   `json:"count"`
}

type UnitTime struct {
	Unit
	Time TimeUnit `json:"time"`
}

type UnitTimeArray []UnitTime

func (lst UnitTimeArray) Reverse() DataList {
	ret := make(UnitTimeArray, len(lst))
	for i, v := range lst {
		ret[len(lst)-1-i] = v
	}
	return ret
}

func (a UnitTimeArray) ToRaw(decimal int8) map[TimeUnit][]byte {
	ret := make(map[TimeUnit][]byte)
	for _, v := range a {
		ret[v.Time] = v.ToRaw(decimal)
	}
	return ret
}

func (m Unit) Type() DataType {
	return UNIT
}

func (m Unit) IsEmpty() bool {
	return m.Count == 0
}

func NewUnit(v float64) Unit {
	if v == 0.00 {
		return Unit{}
	}

	return Unit{
		Open:    v,
		High:    v,
		Low:     v,
		Close:   v,
		Average: v,
		Median:  v,
		Count:   1,
	}
}

func getPrecision(val float64) int {
	// Convert the float64 to a string with high precision
	str := Format.Float(val, -1)
	// Find the position of the decimal point
	decimalPos := strings.Index(str, ".")
	if decimalPos == -1 {
		// No decimal point found, so the precision is 0
		return 0
	}

	//trim all 0's at end
	for i := len(str) - 1; i > decimalPos; i-- {
		if str[i] != '0' {
			break
		}
		str = str[:len(str)-1]
	}

	// The precision is the number of characters after the decimal point
	return len(str) - decimalPos - 1
}

func (lst UnitTimeArray) Append(pt Data) DataList {
	return append(lst, pt.(UnitTime))
}

func (lst UnitTimeArray) Prepend(pt Data) DataList {
	return append(UnitTimeArray{pt.(UnitTime)}, lst...)
}

func (lst UnitTimeArray) First() Data {
	if len(lst) == 0 {
		return nil
	}
	ret := lst[0]
	return &ret
}

func (list UnitTimeArray) ToJSON(columns []ColumnName) ([]map[ColumnName]interface{}, error) {
	for _, col := range columns {
		if lo.IndexOf(UNIT.Columns(), col) == -1 {
			return nil, fmt.Errorf("column %s not found", col)
		}
	}

	return filterToMap(list, columns)
}

func (lst UnitTimeArray) Map() []Data {
	ret := make([]Data, len(lst))
	for i, v := range lst {
		ret[i] = v
	}
	return ret
}

func (lst UnitTimeArray) RemoveFirstN(n int) DataList {
	if n >= len(lst) {
		return PointTimeArray{}
	}
	return lst[n:]
}

func (lst UnitTimeArray) Len() int {
	if lst == nil {
		return 0
	}
	return len(lst)
}

func (list UnitTimeArray) Aggregate(timeframe time.Duration, newTime TimeUnit) Data {
	ret := UnitTime{Time: newTime}
	closes := []float64{}

	absoluteSumDecimals := 0
	absoluteSum := decimal.NewFromFloat(0.00)
	maxClosePrecision := 0

	for i, unit := range list {
		if unit.Count == 0 || unit.Open == 0.00 {
			continue
		}
		currentUnitPricision := getPrecision(unit.Close)
		if currentUnitPricision > maxClosePrecision {
			maxClosePrecision = currentUnitPricision
		}

		if i > 0 && timeframe == Env.MIN_TIME_FRAME {
			prevValue := list[i-1].Close
			currentValue := unit.Close

			if prevValue != currentValue {
				if currentUnitPricision > absoluteSumDecimals {
					absoluteSumDecimals = currentUnitPricision
				}

				max := decimal.NewFromFloat(math.Max(currentValue, prevValue))
				min := decimal.NewFromFloat(math.Min(currentValue, prevValue))

				priceDiff := max.Sub(min)
				absoluteSum = absoluteSum.Add(priceDiff)
			}
		} else if timeframe != Env.MIN_TIME_FRAME {
			absoluteSumDecimals = int(math.Max(float64(absoluteSumDecimals), float64(getPrecision(unit.AbsoluteSum))))
			absoluteSum = absoluteSum.Add(decimal.NewFromFloat(unit.AbsoluteSum))
		}

		if ret.Open == 0.00 {
			ret.Open = unit.Open
		}

		if ret.High == 0.00 {
			ret.High = unit.High
		} else {
			ret.High = math.Max(ret.High, unit.High)
		}
		if ret.Low == 0.00 {
			ret.Low = unit.Low
		} else {
			ret.Low = math.Min(ret.Low, unit.Low)
		}

		ret.Close = unit.Close
		ret.Count += unit.Count
		closes = append(closes, unit.Close)
	}

	ret.AbsoluteSum, _ = absoluteSum.Round(int32(absoluteSumDecimals)).Float64()
	ret.Average = Math.RoundFloat(Math.SafeAverage(closes), uint(maxClosePrecision))
	ret.Median = Math.SafeMedian(closes)
	return ret
}

func ParseRawUnit(raw []byte) Unit {
	s := string(raw)
	splited := strings.Split(s, "@")
	if len(splited) == 1 {
		v, err := strconv.ParseFloat(splited[0], 64)
		if err != nil {
			return Unit{}
		}
		return NewUnit(v)
	}
	if len(splited) != 8 {
		return Unit{}
	}
	open, _ := strconv.ParseFloat(splited[0], 64)
	high, _ := strconv.ParseFloat(splited[1], 64)
	low, _ := strconv.ParseFloat(splited[2], 64)
	close, _ := strconv.ParseFloat(splited[3], 64)
	avg, _ := strconv.ParseFloat(splited[4], 64)
	median, _ := strconv.ParseFloat(splited[5], 64)
	absoluteSum, _ := strconv.ParseFloat(splited[6], 64)
	count, _ := strconv.ParseInt(splited[7], 10, 64)

	return Unit{
		Open:        open,
		High:        high,
		Low:         low,
		Close:       close,
		Average:     avg,
		Median:      median,
		AbsoluteSum: absoluteSum,
		Count:       count,
	}
}

func (p Unit) ToRaw(decimals int8) []byte {
	if p.Count == 1 {
		return []byte(Format.Float(p.Open, decimals))
	}
	return []byte(fmt.Sprintf("%f@%f@%f@%f@%f@%f@%f@%d", p.Open, p.High, p.Low, p.Close, p.Average, p.Median, p.AbsoluteSum, p.Count))
}

func (p Unit) ToTime(time TimeUnit) UnitTime {
	return UnitTime{
		Unit: p,
		Time: time,
	}
}

func (p UnitTime) GetTime() TimeUnit {
	return p.Time
}

func (q UnitTime) CSVLine(decimals int8, requirement CSVCheckListRequirement) []string {
	ret := []string{}

	if requirement[ColumnType.TIME] {
		if q.Time > 0 {
			if Env.MIN_TIME_FRAME >= time.Second && Env.MIN_TIME_FRAME%time.Second == 0 {
				ret = append(ret, strconv.FormatInt(q.Time.ToTime().Unix(), 10))
			} else {
				ret = append(ret, q.Time.String())
			}
		} else {
			ret = append(ret, "")
		}
	}

	if requirement[ColumnType.OPEN] {
		if q.Count >= 1 {
			ret = append(ret, Format.Float(q.Open, decimals))
		} else {
			ret = append(ret, "")
		}
	}

	if requirement[ColumnType.HIGH] {
		if q.Count >= 1 {
			ret = append(ret, Format.Float(q.High, decimals))
		} else {
			ret = append(ret, "")
		}
	}

	if requirement[ColumnType.LOW] {
		if q.Count >= 1 {
			ret = append(ret, Format.Float(q.Low, decimals))
		} else {
			ret = append(ret, "")
		}
	}

	if requirement[ColumnType.CLOSE] {
		if q.Count >= 1 {
			ret = append(ret, Format.Float(q.Close, decimals))
		} else {
			ret = append(ret, "")
		}
	}

	if requirement[ColumnType.AVERAGE] {
		if q.Count >= 1 {
			ret = append(ret, Format.Float(q.Average, decimals))
		} else {
			ret = append(ret, "")
		}
	}

	if requirement[ColumnType.MEDIAN] {
		if q.Count >= 1 {
			ret = append(ret, Format.Float(q.Median, decimals))
		} else {
			ret = append(ret, "")
		}
	}

	if requirement[ColumnType.ABSOLUTE_SUM] {
		if q.AbsoluteSum != 0.00 {
			ret = append(ret, Format.Float(q.AbsoluteSum, decimals))
		} else {
			ret = append(ret, "")
		}
	}

	if requirement[ColumnType.COUNT] {
		if q.Count > 0 {
			ret = append(ret, strconv.FormatInt(q.Count, 10))
		} else {
			ret = append(ret, "")
		}
	}
	return ret
}

func (u UnitTime) String() string {
	return fmt.Sprintf("[%d] Open: %s High: %s Low: %s Close: %s Average: %s Median: %s AbsoluteSum: %s Count: %d", u.Time.ToTime().Unix(),
		Format.Float(u.Open, -1),
		Format.Float(u.High, -1),
		Format.Float(u.Low, -1),
		Format.Float(u.Close, -1),
		Format.Float(u.Average, -1),
		Format.Float(u.Median, -1),
		Format.Float(u.AbsoluteSum, -1),
		u.Count)
}
