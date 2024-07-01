package pcommon

import (
	"errors"
	"time"
)

type DataType int8

// units are data that can be aggregated around a candle (open, close, high, low, etc)
const UNIT DataType = 1

// quantities are data that can be summed up (volume, open interest, etc)
const QUANTITY DataType = 2

// points are simple data (a float64) that cannot be aggregated or summed
const POINT DataType = 3

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
	Map() []Data
	ToJSON(columns []ColumnName) ([]map[ColumnName]interface{}, error)
}

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
