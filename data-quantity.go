package pcommon

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

type Quantity struct {
	Plus  float64 `json:"plus"`
	Minus float64 `json:"minus"`

	PlusAvg  float64 `json:"plus_average"`
	MinusAvg float64 `json:"minus_average"`

	PlusMed  float64 `json:"plus_median"`  // median
	MinusMed float64 `json:"minus_median"` // median

	PlusCount  int64 `json:"plus_count"`  // count
	MinusCount int64 `json:"minus_count"` // count
}

func (m Quantity) Type() DataType {
	return QUANTITY
}

type QuantityTime struct {
	Quantity
	Time TimeUnit `json:"time"`
}

type QuantityTimeArray []QuantityTime

func (lst QuantityTimeArray) Reverse() DataList {
	ret := make(QuantityTimeArray, len(lst))
	for i, v := range lst {
		ret[len(lst)-1-i] = v
	}
	return ret
}

func (a QuantityTimeArray) ToRaw(decimal int8) map[TimeUnit][]byte {
	ret := make(map[TimeUnit][]byte)
	for _, v := range a {
		ret[v.Time] = v.ToRaw(decimal)
	}
	return ret
}

func NewQuantity(v float64) Quantity {
	if v == 0.00 {
		return Quantity{}
	}

	ret := Quantity{}
	vAbs := math.Abs(v)

	if v > 0 {
		ret.Plus = vAbs
		ret.PlusAvg = vAbs
		ret.PlusMed = vAbs
		ret.PlusCount = 1
	} else {
		ret.Minus = vAbs
		ret.MinusAvg = vAbs
		ret.MinusMed = vAbs
		ret.MinusCount = 1
	}
	return ret
}

func (lst QuantityTimeArray) Append(pt Data) DataList {
	return append(lst, pt.(QuantityTime))
}

func (lst QuantityTimeArray) Prepend(pt Data) DataList {
	return append(QuantityTimeArray{pt.(QuantityTime)}, lst...)
}

func (lst QuantityTimeArray) First() Data {
	if len(lst) == 0 {
		return nil
	}
	return &lst[0]
}

func (lst QuantityTimeArray) Last() Data {
	if len(lst) == 0 {
		return nil
	}
	return &lst[len(lst)-1]
}

func (lst QuantityTimeArray) Len() int {
	if lst == nil {
		return 0
	}
	return len(lst)
}

func (lst QuantityTimeArray) RemoveFirstN(n int) DataList {
	if n >= len(lst) {
		return PointTimeArray{}
	}
	return lst[n:]
}

func (list QuantityTimeArray) ToJSON(columns []ColumnName) ([]map[ColumnName]interface{}, error) {
	for _, col := range columns {
		if lo.IndexOf(QUANTITY.Columns(), col) == -1 {
			return nil, fmt.Errorf("column %s not found", col)
		}
	}

	return filterToMap(list, columns)
}

func (lst QuantityTimeArray) Map() []Data {
	ret := make([]Data, len(lst))
	for i, v := range lst {
		ret[i] = v
	}
	return ret
}

func (list QuantityTimeArray) Aggregate(timeframe time.Duration, newTime TimeUnit) Data {
	ret := QuantityTime{Time: newTime}

	amountsPlus := []float64{}
	amountMinus := []float64{}

	for _, q := range list {
		if q.Plus > 0 {
			ret.Plus += q.Plus
			ret.PlusCount++
			amountsPlus = append(amountsPlus, q.Plus)
		}
		if q.Minus > 0 {
			ret.Minus += q.Minus
			ret.MinusCount++
			amountMinus = append(amountMinus, q.Minus)
		}
	}
	ret.MinusAvg = Math.SafeAverage(amountMinus)
	ret.PlusAvg = Math.SafeAverage(amountsPlus)

	ret.PlusMed = Math.SafeMedian(amountsPlus)
	ret.MinusMed = Math.SafeMedian(amountMinus)
	return ret
}

func (p QuantityTime) ValueAt(column ColumnName) (float64, error) {
	switch column {
	case ColumnType.PLUS:
		return p.Plus, nil
	case ColumnType.MINUS:
		return p.Minus, nil
	case ColumnType.PLUS_AVERAGE:
		return p.PlusAvg, nil
	case ColumnType.MINUS_AVERAGE:
		return p.MinusAvg, nil
	case ColumnType.PLUS_MEDIAN:
		return p.PlusMed, nil
	case ColumnType.MINUS_MEDIAN:
		return p.MinusMed, nil
	case ColumnType.PLUS_COUNT:
		return float64(p.PlusCount), nil
	case ColumnType.MINUS_COUNT:
		return float64(p.MinusCount), nil
	}
	return 0.00, fmt.Errorf("column %s not found", column)
}

func (m Quantity) IsEmpty() bool {
	return m.MinusCount == 0 && m.PlusCount == 0
}

func ParseRawQuantity(raw []byte) Quantity {
	s := string(raw)

	splited := strings.Split(s, "@")
	if len(splited) == 1 {
		v, err := strconv.ParseFloat(splited[0], 64)
		if err != nil {
			log.Fatal("Invalid float format")
		}
		return NewQuantity(v)
	}

	if len(splited) != 8 {
		log.Fatal("Invalid quantity format")
	}

	plus, _ := strconv.ParseFloat(splited[0], 64)
	minus, _ := strconv.ParseFloat(splited[1], 64)

	plusAvg, _ := strconv.ParseFloat(splited[2], 64)
	minusAvg, _ := strconv.ParseFloat(splited[3], 64)

	plusMed, _ := strconv.ParseFloat(splited[4], 64)
	minusMed, _ := strconv.ParseFloat(splited[5], 64)

	plusCount, _ := strconv.ParseInt(splited[6], 10, 64)
	minusCount, _ := strconv.ParseInt(splited[7], 10, 64)

	return Quantity{
		Plus:       plus,
		Minus:      minus,
		PlusAvg:    plusAvg,
		MinusAvg:   minusAvg,
		PlusMed:    plusMed,
		MinusMed:   minusMed,
		PlusCount:  plusCount,
		MinusCount: minusCount,
	}
}

func (q Quantity) ToRaw(decimals int8) []byte {
	if q.MinusCount+q.PlusCount == 1 {
		if q.Plus > 0 {
			return []byte(Format.Float(q.Plus, decimals))
		}
		return []byte(Format.Float(q.Minus*-1, decimals))
	}
	ret := fmt.Sprintf("%s@%s@%s@%s@%s@%s@%d@%d",
		Format.Float(q.Plus, decimals), Format.Float(q.Minus, decimals),
		Format.Float(q.PlusAvg, decimals), Format.Float(q.MinusAvg, decimals),
		Format.Float(q.PlusMed, decimals), Format.Float(q.MinusMed, decimals),
		q.PlusCount, q.MinusCount)
	return []byte(ret)
}

func (q Quantity) ToTime(time TimeUnit) QuantityTime {
	return QuantityTime{
		Quantity: q,
		Time:     time,
	}
}

func (p QuantityTime) Min() float64 {
	return -p.Minus
}

func (p QuantityTime) Max() float64 {
	return p.Plus
}

func (p QuantityTime) GetTime() TimeUnit {
	return p.Time
}

func (q QuantityTime) CSVLine(volumeDecimals int8, requiremment CSVCheckListRequirement) []string {
	ret := []string{}

	if requiremment[ColumnType.TIME] {
		if q.Time > 0 {
			if Env.MIN_TIME_FRAME >= time.Second {
				ret = append(ret, strconv.FormatInt(q.Time.ToTime().Unix(), 10))
			} else {
				ret = append(ret, q.Time.String())
			}
		} else {
			ret = append(ret, "")
		}
	}

	if requiremment[ColumnType.PLUS] {
		if q.Plus != 0 {
			ret = append(ret, Format.Float(q.Plus, volumeDecimals))
		} else {
			ret = append(ret, "")
		}
	}

	if requiremment[ColumnType.MINUS] {
		if q.Minus != 0 {
			ret = append(ret, Format.Float(q.Minus, volumeDecimals))
		} else {
			ret = append(ret, "")
		}
	}

	if requiremment[ColumnType.PLUS_AVERAGE] {
		if q.PlusAvg != 0 {
			ret = append(ret, Format.Float(q.PlusAvg, volumeDecimals))
		} else {
			ret = append(ret, "")
		}
	}

	if requiremment[ColumnType.MINUS_AVERAGE] {
		if q.MinusAvg != 0 {
			ret = append(ret, Format.Float(q.MinusAvg, volumeDecimals))
		} else {
			ret = append(ret, "")
		}
	}

	if requiremment[ColumnType.PLUS_MEDIAN] {
		if q.PlusMed != 0 {
			ret = append(ret, Format.Float(q.PlusMed, volumeDecimals))
		} else {
			ret = append(ret, "")
		}
	}

	if requiremment[ColumnType.MINUS_MEDIAN] {
		if q.MinusMed != 0 {
			ret = append(ret, Format.Float(q.MinusMed, volumeDecimals))
		} else {
			ret = append(ret, "")
		}
	}

	if requiremment[ColumnType.PLUS_COUNT] {
		if q.PlusCount != 0 {
			ret = append(ret, strconv.FormatInt(q.PlusCount, 10))
		} else {
			ret = append(ret, "")
		}
	}

	if requiremment[ColumnType.MINUS_COUNT] {
		if q.MinusCount != 0 {
			ret = append(ret, strconv.FormatInt(q.MinusCount, 10))
		} else {
			ret = append(ret, "")
		}
	}

	return ret
}

func (qty QuantityTime) String() string {
	return fmt.Sprintf("[%d] Plus: %f Minus: %f PlusAvg: %f MinusAvg: %f PlusMed: %f MinusMed: %f PlusCount: %d MinusCount: %d", qty.Time.ToTime().Unix(), qty.Plus, qty.Minus, qty.PlusAvg, qty.MinusAvg, qty.PlusMed, qty.MinusMed, qty.PlusCount, qty.MinusCount)
}
