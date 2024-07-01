package pcommon

import (
	"sort"

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

var UNIT_COLUNMS = []ColumnName{ColumnType.TIME, ColumnType.OPEN, ColumnType.HIGH, ColumnType.LOW, ColumnType.CLOSE, ColumnType.AVERAGE, ColumnType.MEDIAN, ColumnType.ABSOLUTE_SUM, ColumnType.COUNT}
var QUANTITY_COLUMNS = []ColumnName{ColumnType.TIME, ColumnType.PLUS, ColumnType.MINUS, ColumnType.PLUS_AVERAGE, ColumnType.MINUS_AVERAGE, ColumnType.PLUS_MEDIAN, ColumnType.MINUS_MEDIAN, ColumnType.PLUS_COUNT, ColumnType.MINUS_COUNT}
var POINT_COLUMNS = []ColumnName{ColumnType.TIME, ColumnType.VALUE}
