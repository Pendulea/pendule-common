package pcommon

import (
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type AssetType string
type AllAssetTypes struct {
	SPOT_PRICE  AssetType
	SPOT_VOLUME AssetType

	FUTURES_PRICE  AssetType
	FUTURES_VOLUME AssetType

	BOOK_DEPTH_P1 AssetType
	BOOK_DEPTH_P2 AssetType
	BOOK_DEPTH_P3 AssetType
	BOOK_DEPTH_P4 AssetType
	BOOK_DEPTH_P5 AssetType

	BOOK_DEPTH_M1 AssetType
	BOOK_DEPTH_M2 AssetType
	BOOK_DEPTH_M3 AssetType
	BOOK_DEPTH_M4 AssetType
	BOOK_DEPTH_M5 AssetType

	METRIC_SUM_OPEN_INTEREST                 AssetType
	METRIC_COUNT_TOP_TRADER_LONG_SHORT_RATIO AssetType
	METRIC_SUM_TOP_TRADER_LONG_SHORT_RATIO   AssetType
	METRIC_COUNT_LONG_SHORT_RATIO            AssetType
	METRIC_SUM_TAKER_LONG_SHORT_VOL_RATIO    AssetType

	CIRCULATING_SUPPLY AssetType
	RSI                AssetType
}

var Asset = AllAssetTypes{
	SPOT_PRICE:  "spot_price",
	SPOT_VOLUME: "spot_volume",

	FUTURES_PRICE:  "futures_price",
	FUTURES_VOLUME: "futures_volume",

	BOOK_DEPTH_P1: "bd-p1",
	BOOK_DEPTH_P2: "bd-p2",
	BOOK_DEPTH_P3: "bd-p3",
	BOOK_DEPTH_P4: "bd-p4",
	BOOK_DEPTH_P5: "bd-p5",

	BOOK_DEPTH_M1: "bd-m1",
	BOOK_DEPTH_M2: "bd-m2",
	BOOK_DEPTH_M3: "bd-m3",
	BOOK_DEPTH_M4: "bd-m4",
	BOOK_DEPTH_M5: "bd-m5",

	METRIC_SUM_OPEN_INTEREST:                 "metrics_sum_open_interest",
	METRIC_COUNT_TOP_TRADER_LONG_SHORT_RATIO: "metrics_count_toptrader_long_short_ratio",
	METRIC_SUM_TOP_TRADER_LONG_SHORT_RATIO:   "metrics_sum_toptrader_long_short_ratio",
	METRIC_COUNT_LONG_SHORT_RATIO:            "metrics_count_long_short_ratio",
	METRIC_SUM_TAKER_LONG_SHORT_VOL_RATIO:    "metrics_sum_taker_long_short_vol_ratio",

	CIRCULATING_SUPPLY: "circulating_supply",
	RSI:                "rsi",
}

var AssetTypeMap = Asset.ToMap()

func (a AllAssetTypes) ToMap() map[string]bool {
	v := reflect.ValueOf(Asset)
	if v.Kind() != reflect.Struct {
		log.Fatal("expected a struct")
	}
	m := make(map[string]bool)
	for i := 0; i < v.NumField(); i++ {
		m[v.Field(i).String()] = true
	}
	return m
}

func (asset AssetType) GetBookDepthAssetPercentage() (int, error) {
	if asset == Asset.BOOK_DEPTH_M1 || asset == Asset.BOOK_DEPTH_M2 || asset == Asset.BOOK_DEPTH_M3 || asset == Asset.BOOK_DEPTH_M4 || asset == Asset.BOOK_DEPTH_M5 ||
		asset == Asset.BOOK_DEPTH_P1 || asset == Asset.BOOK_DEPTH_P2 || asset == Asset.BOOK_DEPTH_P3 || asset == Asset.BOOK_DEPTH_P4 || asset == Asset.BOOK_DEPTH_P5 {
		lastChar := asset[len(asset)-1:]
		percent, err := strconv.Atoi(string(lastChar))
		if err != nil {
			return 0, err
		}

		isPlus := strings.HasPrefix(string(asset), "bd-p")
		isMinus := strings.HasPrefix(string(asset), "bd-m")
		if isPlus {
			return percent, nil
		} else if isMinus {
			return -percent, nil
		}
	}
	return 0, errors.New("invalid asset")
}

type AssetStateConfig struct {
	ID       AssetType
	DataType DataType

	RequiredDependencyDataTypes []DataType
	RequiredArgumentTypes       []reflect.Type
	Label                       string
	Description                 string
}

type AvailableAssets map[AssetType]AssetStateConfig

var DEFAULT_ASSETS = AvailableAssets{
	//binance spot trades
	Asset.SPOT_PRICE: {
		Asset.SPOT_PRICE, UNIT, nil, nil,
		"", "",
	},
	Asset.SPOT_VOLUME: {Asset.SPOT_VOLUME, QUANTITY, nil, nil,
		"", "",
	},

	//binance book depth
	Asset.BOOK_DEPTH_P1: {Asset.BOOK_DEPTH_P1, UNIT, nil, nil,
		"", "",
	},
	Asset.BOOK_DEPTH_P2: {Asset.BOOK_DEPTH_P2, UNIT, nil, nil,
		"", "",
	},
	Asset.BOOK_DEPTH_P3: {Asset.BOOK_DEPTH_P3, UNIT, nil, nil,
		"", "",
	},
	Asset.BOOK_DEPTH_P4: {Asset.BOOK_DEPTH_P4, UNIT, nil, nil,
		"", "",
	},
	Asset.BOOK_DEPTH_P5: {Asset.BOOK_DEPTH_P5, UNIT, nil, nil,
		"", "",
	},
	Asset.BOOK_DEPTH_M1: {Asset.BOOK_DEPTH_M1, UNIT, nil, nil,
		"", "",
	},
	Asset.BOOK_DEPTH_M2: {Asset.BOOK_DEPTH_M2, UNIT, nil, nil,
		"", "",
	},
	Asset.BOOK_DEPTH_M3: {Asset.BOOK_DEPTH_M3, UNIT, nil, nil,
		"", "",
	},
	Asset.BOOK_DEPTH_M4: {Asset.BOOK_DEPTH_M4, UNIT, nil, nil,
		"", "",
	},
	Asset.BOOK_DEPTH_M5: {Asset.BOOK_DEPTH_M5, UNIT, nil, nil,
		"", "",
	},

	Asset.METRIC_SUM_OPEN_INTEREST: {Asset.METRIC_SUM_OPEN_INTEREST, UNIT, nil, nil,
		"", "",
	},

	Asset.METRIC_COUNT_TOP_TRADER_LONG_SHORT_RATIO: {Asset.METRIC_COUNT_TOP_TRADER_LONG_SHORT_RATIO, UNIT, nil, nil,
		"", "",
	},
	Asset.METRIC_SUM_TOP_TRADER_LONG_SHORT_RATIO: {Asset.METRIC_SUM_TOP_TRADER_LONG_SHORT_RATIO, UNIT, nil, nil,
		"", "",
	},
	Asset.METRIC_COUNT_LONG_SHORT_RATIO: {Asset.METRIC_COUNT_LONG_SHORT_RATIO, UNIT, nil, nil,
		"", "",
	},
	Asset.METRIC_SUM_TAKER_LONG_SHORT_VOL_RATIO: {Asset.METRIC_SUM_TAKER_LONG_SHORT_VOL_RATIO, UNIT, nil, nil,
		"", "",
	},

	Asset.CIRCULATING_SUPPLY: {Asset.CIRCULATING_SUPPLY, UNIT, nil, nil,
		"", "",
	},

	Asset.FUTURES_PRICE: {Asset.FUTURES_PRICE, UNIT, nil, nil,
		"", "",
	},
	Asset.FUTURES_VOLUME: {Asset.FUTURES_VOLUME, QUANTITY, nil, nil,
		"", "",
	},

	Asset.RSI: {Asset.RSI, POINT, []DataType{UNIT}, []reflect.Type{reflect.TypeOf(int64(0))},
		"RSI", "",
	},
}

type AvailableAssetJSON struct {
	AssetType AssetType `json:"asset_type"`
	DataType  DataType  `json:"data_type"`

	Dependencies        []DataType   `json:"dependencies"`
	ArgumentTypes       []string     `json:"argument_types"`
	Label               string       `json:"label"`
	Description         string       `json:"description"`
	DataTypeName        string       `json:"data_type_name"`
	DataTypeColor       string       `json:"data_type_color"`
	DataTypeColumns     []ColumnName `json:"data_type_columns"`
	DataTypeDescription string       `json:"data_type_description"`
}

func (aa AvailableAssets) JSON() []AvailableAssetJSON {
	ret := []AvailableAssetJSON{}
	for _, v := range aa {
		argumentTypes := []string{}
		for _, arg := range v.RequiredArgumentTypes {
			argumentTypes = append(argumentTypes, arg.String())
		}
		ret = append(ret, AvailableAssetJSON{
			AssetType:           v.ID,
			DataType:            v.DataType,
			DataTypeName:        v.DataType.String(),
			DataTypeColor:       v.DataType.Color(),
			DataTypeColumns:     v.DataType.Columns(),
			DataTypeDescription: v.DataType.Description(),
			Dependencies:        v.RequiredDependencyDataTypes,
			ArgumentTypes:       argumentTypes,
			Label:               v.Label,
			Description:         v.Description,
		})
	}
	return ret
}
