package pcommon

import (
	"errors"
	"sort"
	"time"

	"github.com/samber/lo"
)

type ColumnName string

var ColumnType = struct {
	TIME ColumnName

	//Quantity
	PLUS          ColumnName
	MINUS         ColumnName
	PLUS_AVERAGE  ColumnName
	MINUS_AVERAGE ColumnName
	PLUS_MEDIAN   ColumnName
	MINUS_MEDIAN  ColumnName
	PLUS_COUNT    ColumnName
	MINUS_COUNT   ColumnName

	//Unit
	OPEN         ColumnName
	HIGH         ColumnName
	LOW          ColumnName
	CLOSE        ColumnName
	AVERAGE      ColumnName
	MEDIAN       ColumnName
	ABSOLUTE_SUM ColumnName
	COUNT        ColumnName

	//Point
	VALUE ColumnName
}{
	TIME:          "time",
	PLUS:          "plus",
	MINUS:         "minus",
	PLUS_AVERAGE:  "plus_average",
	MINUS_AVERAGE: "minus_average",
	PLUS_MEDIAN:   "plus_median",
	MINUS_MEDIAN:  "minus_median",
	PLUS_COUNT:    "plus_count",
	MINUS_COUNT:   "minus_count",

	OPEN:         "open",
	HIGH:         "high",
	LOW:          "low",
	CLOSE:        "close",
	AVERAGE:      "average",
	MEDIAN:       "median",
	ABSOLUTE_SUM: "absolute_sum",
	COUNT:        "count",

	VALUE: "value",
}

type CSVCheckListRequirement map[ColumnName]bool

func (c CSVCheckListRequirement) Columns() []ColumnName {
	result := lo.MapToSlice(c, func(k ColumnName, v bool) ColumnName {
		if v {
			return k
		}
		return ""
	})
	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	return result
}

type Data interface {
	CSVLine(volumeDecimals int8, requirement CSVCheckListRequirement) []string
	ToRaw(decimals int8) []byte
	IsEmpty() bool
	Type() DataType
	GetTime() TimeUnit
}

type DataList interface {
	Aggregate(timeframe time.Duration, newTime TimeUnit) Data
	First() Data
	ToRaw(decimals int8) map[TimeUnit][]byte
	Append(pt Data) DataList
	Prepend(pt Data) DataList
	Len() int
	RemoveFirstN(n int) DataList
}

type DataType int8

var UNIT_COLUNMS = []ColumnName{ColumnType.TIME, ColumnType.OPEN, ColumnType.HIGH, ColumnType.LOW, ColumnType.CLOSE, ColumnType.AVERAGE, ColumnType.MEDIAN, ColumnType.ABSOLUTE_SUM, ColumnType.COUNT}
var QUANTITY_COLUMNS = []ColumnName{ColumnType.TIME, ColumnType.PLUS, ColumnType.MINUS, ColumnType.PLUS_AVERAGE, ColumnType.MINUS_AVERAGE, ColumnType.PLUS_MEDIAN, ColumnType.MINUS_MEDIAN, ColumnType.PLUS_COUNT, ColumnType.MINUS_COUNT}
var POINT_COLUMNS = []ColumnName{ColumnType.TIME, ColumnType.VALUE}

func NewTypeTime(t DataType, value float64, valueTime TimeUnit) Data {
	if t == UNIT {
		return NewUnit(value).ToTime(valueTime)
	}
	if t == QUANTITY {
		return NewQuantity(value).ToTime(valueTime)
	}
	if t == POINT {
		return newPoint(value).ToTime(valueTime)
	}
	return nil
}

func NewTypeTimeArray(t DataType) DataList {
	if t == UNIT {
		return UnitTimeArray{}
	}
	if t == QUANTITY {
		return QuantityTimeArray{}
	}
	if t == POINT {
		return PointTimeArray{}
	}
	return nil
}

func ParseTypeData(t DataType, d []byte, dataTime TimeUnit) (Data, error) {
	if t == UNIT {
		return ParseRawUnit(d).ToTime(dataTime), nil
	}
	if t == QUANTITY {
		return ParseRawQuantity(d).ToTime(dataTime), nil
	}
	if t == POINT {
		p, err := ParseRawPoint(d)
		if err != nil {
			return nil, err
		}
		return p.ToTime(dataTime), nil
	}
	return nil, errors.New("unknown data type")
}

// units are data that can be aggregated around a candle (open, close, high, low, etc)
const UNIT DataType = 1

// quantities are data that can be summed up (volume, open interest, etc)
const QUANTITY DataType = 2

// points are simple data (a float64) that cannot be aggregated or summed
const POINT DataType = 3

func (d DataType) Columns() []ColumnName {
	if d == UNIT {
		return UNIT_COLUNMS
	}
	if d == QUANTITY {
		return QUANTITY_COLUMNS
	}
	if d == POINT {
		return POINT_COLUMNS
	}
	return []ColumnName{}
}

func (q DataType) Header(prefix string, requirement CSVCheckListRequirement) []string {
	list := []string{}
	for _, column := range q.Columns() {
		if requirement[column] {
			if column == ColumnType.VALUE {
				list = append(list, prefix)
			} else {
				list = append(list, prefix+"_"+string(column))
			}
		}
	}
	return list
}
