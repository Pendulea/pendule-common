package pcommon

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/samber/lo"
)

type Point struct {
	Value float64 `json:"v"`
}

type PointTime struct {
	Point
	Time TimeUnit `json:"time"`
}

type PointTimeArray []PointTime

func (lst PointTimeArray) Reverse() DataList {
	ret := make(PointTimeArray, len(lst))
	for i, v := range lst {
		ret[len(lst)-1-i] = v
	}
	return ret
}

func (lst PointTimeArray) Aggregate(timeframe time.Duration, newTime TimeUnit) Data {
	log.Fatal("no aggregation for points data")
	return PointTime{}
}

func (lst PointTimeArray) Map() []Data {
	ret := make([]Data, len(lst))
	for i, v := range lst {
		ret[i] = v
	}
	return ret
}

func (list PointTimeArray) ToJSON(columns []ColumnName) ([]map[ColumnName]interface{}, error) {
	for _, col := range columns {
		if lo.IndexOf(POINT.Columns(), col) == -1 {
			return nil, fmt.Errorf("column %s not found", col)
		}
	}

	return filterToMap(list, columns)
}

func (lst PointTimeArray) Append(pt Data) DataList {
	return append(lst, pt.(PointTime))
}

func (lst PointTimeArray) Prepend(pt Data) DataList {
	return append(PointTimeArray{pt.(PointTime)}, lst...)
}

func (lst PointTimeArray) RemoveFirstN(n int) DataList {
	if n >= len(lst) {
		return PointTimeArray{}
	}
	return lst[n:]
}

func (lst PointTimeArray) First() Data {
	if len(lst) == 0 {
		return nil
	}
	return &lst[0]
}

func (lst PointTimeArray) Last() Data {
	if len(lst) == 0 {
		return nil
	}
	return &lst[len(lst)-1]
}

func (lst PointTimeArray) Len() int {
	if lst == nil {
		return 0
	}
	return len(lst)
}

func (lst PointTimeArray) ToRaw(decimal int8) map[TimeUnit][]byte {
	ret := make(map[TimeUnit][]byte)
	for _, v := range lst {
		ret[v.Time] = v.ToRaw(decimal)
	}
	return ret
}

func newPoint(v float64) Point {
	if v == 0.00 {
		return Point{}
	}

	return Point{Value: v}
}

func (m Point) Type() DataType {
	return POINT
}

func (m Point) IsEmpty() bool {
	return m.Value == 0.00
}

func ParseRawPoint(d []byte) (Point, error) {
	if len(d) == 0 {
		return Point{}, nil
	}
	v, err := strconv.ParseFloat(string(d), 64)
	if err != nil {
		return Point{}, err
	}
	return newPoint(v), nil
}

func (p Point) ToTime(time TimeUnit) PointTime {
	return PointTime{Point: p, Time: time}
}

func (p Point) ToRaw(decimals int8) []byte {
	return []byte(Format.Float(p.Value, decimals))
}

func (p PointTime) GetTime() TimeUnit {
	return p.Time
}

func (m PointTime) CSVLine(volumeDecimals int8, requirement CSVCheckListRequirement) []string {
	ret := []string{}

	if requirement[ColumnType.TIME] {
		if m.Time > 0 {
			if Env.MIN_TIME_FRAME >= time.Second {
				ret = append(ret, strconv.FormatInt(m.Time.ToTime().Unix(), 10))
			} else {
				ret = append(ret, m.Time.String())
			}
		} else {
			ret = append(ret, "")
		}
	}

	if requirement[ColumnType.VALUE] {
		if m.Value != 0.00 {
			ret = append(ret, Format.Float(m.Value, volumeDecimals))
		} else {
			ret = append(ret, "")
		}
	}

	return ret
}
