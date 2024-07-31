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

	RSI AssetType

	RSI2 AssetType
	SMA  AssetType
	EMA  AssetType
	WMA  AssetType
	HMA  AssetType
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

	RSI2: "rsi2",
	SMA:  "sma",
	EMA:  "ema",
	WMA:  "wma",
	HMA:  "hma",
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
	SetUpDecimals func(priceUSDA, priceUSDB float64) int8

	ID       AssetType
	DataType DataType

	RequiredDependencyDataTypes []DataType
	RequiredArgumentTypes       []reflect.Type
	Label                       string
	Description                 string
	Color                       string
}

type AvailableAssets map[AssetType]AssetStateConfig

var DEFAULT_ASSETS = AvailableAssets{
	// Binance spot trades
	Asset.SPOT_PRICE: {
		func(priceUSDA, priceUSDB float64) int8 {
			return priceDecimals(priceUSDA / priceUSDB)
		},
		Asset.SPOT_PRICE, UNIT, nil, nil,
		"Spot Price", "The current price at which an asset is bought or sold in the spot market on Binance.",
		"#5F9EA0",
	},
	Asset.SPOT_VOLUME: {
		func(priceUSDA, priceUSDB float64) int8 {
			return countDivisionsTo(priceUSDA/priceUSDB, 0.01)
		},
		Asset.SPOT_VOLUME, QUANTITY, nil, nil,
		"Spot Volume", "The total amount of an asset traded in the spot market on Binance.",
		"#228B22",
	},

	// Binance order book depth
	Asset.BOOK_DEPTH_P1: {
		func(priceUSDA, priceUSDB float64) int8 {
			return countDivisionsTo(priceUSDA/priceUSDB, 10)
		},
		Asset.BOOK_DEPTH_P1, UNIT, nil, nil,
		"Liquidity +1% Price", "Available liquidity at a price level 1% above the current market price on Binance.",
		"#044f56",
	},
	Asset.BOOK_DEPTH_P2: {
		func(priceUSDA, priceUSDB float64) int8 {
			return countDivisionsTo(priceUSDA/priceUSDB, 10)
		},
		Asset.BOOK_DEPTH_P2, UNIT, nil, nil,
		"Liquidity +2% Price", "Available liquidity at a price level 2% above the current market price on Binance.",
		"#07636c",
	},
	Asset.BOOK_DEPTH_P3: {
		func(priceUSDA, priceUSDB float64) int8 {
			return countDivisionsTo(priceUSDA/priceUSDB, 10)
		},
		Asset.BOOK_DEPTH_P3, UNIT, nil, nil,
		"Liquidity +3% Price", "Available liquidity at a price level 3% above the current market price on Binance.",
		"#0a7882",
	},
	Asset.BOOK_DEPTH_P4: {
		func(priceUSDA, priceUSDB float64) int8 {
			return countDivisionsTo(priceUSDA/priceUSDB, 10)
		},
		Asset.BOOK_DEPTH_P4, UNIT, nil, nil,
		"Liquidity +4% Price", "Available liquidity at a price level 4% above the current market price on Binance.",
		"#0e8d99",
	},
	Asset.BOOK_DEPTH_P5: {
		func(priceUSDA, priceUSDB float64) int8 {
			return countDivisionsTo(priceUSDA/priceUSDB, 10)
		},
		Asset.BOOK_DEPTH_P5, UNIT, nil, nil,
		"Liquidity +5% Price", "Available liquidity at a price level 5% above the current market price on Binance.",
		"#12a3b0",
	},
	Asset.BOOK_DEPTH_M1: {
		func(priceUSDA, priceUSDB float64) int8 {
			return countDivisionsTo(priceUSDA/priceUSDB, 10)
		},
		Asset.BOOK_DEPTH_M1, UNIT, nil, nil,
		"Liquidity -1% Price", "Available liquidity at a price level 1% below the current market price on Binance.",
		"#0b186b",
	},
	Asset.BOOK_DEPTH_M2: {
		func(priceUSDA, priceUSDB float64) int8 {
			return countDivisionsTo(priceUSDA/priceUSDB, 10)
		},
		Asset.BOOK_DEPTH_M2, UNIT, nil, nil,
		"Liquidity -2% Price", "Available liquidity at a price level 2% below the current market price on Binance.",
		"#0b186b",
	},
	Asset.BOOK_DEPTH_M3: {
		func(priceUSDA, priceUSDB float64) int8 {
			return countDivisionsTo(priceUSDA/priceUSDB, 10)
		},
		Asset.BOOK_DEPTH_M3, UNIT, nil, nil,
		"Liquidity -3% Price", "Available liquidity at a price level 3% below the current market price on Binance.",
		"#08135c",
	},
	Asset.BOOK_DEPTH_M4: {
		func(priceUSDA, priceUSDB float64) int8 {
			return countDivisionsTo(priceUSDA/priceUSDB, 10)
		},
		Asset.BOOK_DEPTH_M4, UNIT, nil, nil,
		"Liquidity -4% Price", "Available liquidity at a price level 4% below the current market price on Binance.",
		"#060f4e",
	},
	Asset.BOOK_DEPTH_M5: {
		func(priceUSDA, priceUSDB float64) int8 {
			return countDivisionsTo(priceUSDA/priceUSDB, 10)
		},
		Asset.BOOK_DEPTH_M5, UNIT, nil, nil,
		"Liquidity -5% Price", "Available liquidity at a price level 5% below the current market price on Binance.",
		"#040a3f",
	},

	// Metrics
	Asset.METRIC_SUM_OPEN_INTEREST: {
		func(priceUSDA, priceUSDB float64) int8 {
			return countDivisionsTo(priceUSDA/priceUSDB, 10)
		},
		Asset.METRIC_SUM_OPEN_INTEREST, UNIT, nil, nil,
		"Open Interest", "The total number of outstanding derivative contracts, such as options or futures, that have not been settled on Binance.",
		"#f1ae8a",
	},

	Asset.METRIC_COUNT_TOP_TRADER_LONG_SHORT_RATIO: {
		func(priceUSDA, priceUSDB float64) int8 {
			return 4
		},
		Asset.METRIC_COUNT_TOP_TRADER_LONG_SHORT_RATIO, UNIT, nil, nil,
		"Taker Long/Short Ratio", "The ratio of long to short positions taken by top traders on Binance.",
		"#5e4e29",
	},
	Asset.METRIC_SUM_TOP_TRADER_LONG_SHORT_RATIO: {
		func(priceUSDA, priceUSDB float64) int8 {
			return 4
		},
		Asset.METRIC_SUM_TOP_TRADER_LONG_SHORT_RATIO, UNIT, nil, nil,
		"Top Trader Long/Short Ratio", "The ratio of the sum of long to short positions taken by top traders on Binance.",
		"#6d5e3d",
	},
	Asset.METRIC_COUNT_LONG_SHORT_RATIO: {
		func(priceUSDA, priceUSDB float64) int8 {
			return 4
		},
		Asset.METRIC_COUNT_LONG_SHORT_RATIO, UNIT, nil, nil,
		"Long/Short Ratio", "The overall ratio of long to short positions taken by all traders on Binance.",
		"#b09763",
	},
	Asset.METRIC_SUM_TAKER_LONG_SHORT_VOL_RATIO: {
		func(priceUSDA, priceUSDB float64) int8 {
			return 4
		},
		Asset.METRIC_SUM_TAKER_LONG_SHORT_VOL_RATIO, UNIT, nil, nil,
		"Taker Long/Short Volume Ratio", "The ratio of long to short volumes taken by top traders on Binance.",
		"#b8a173",
	},

	// Supply
	Asset.CIRCULATING_SUPPLY: {
		func(priceUSDA, priceUSDB float64) int8 {
			return 0
		},
		Asset.CIRCULATING_SUPPLY, UNIT, nil, nil,
		"Circulating Supply", "The total number of tokens that are currently available in circulation.",
		"#eaeaea",
	},

	// Binance futures
	Asset.FUTURES_PRICE: {
		func(priceUSDA, priceUSDB float64) int8 {
			return priceDecimals(priceUSDA / priceUSDB)
		},
		Asset.FUTURES_PRICE, UNIT, nil, nil,
		"Futures Price", "The current price at which a futures contract is trading on Binance.",
		"#386061",
	},
	Asset.FUTURES_VOLUME: {
		func(priceUSDA, priceUSDB float64) int8 {
			return countDivisionsTo(priceUSDA/priceUSDB, 0.01)
		},
		Asset.FUTURES_VOLUME, QUANTITY, nil, nil,
		"Futures Volume", "The total amount of futures contracts traded on Binance.",
		"#0b430b",
	},

	// Technical Indicators
	Asset.RSI: {
		func(priceUSDA, priceUSDB float64) int8 {
			return 2
		},
		Asset.RSI, POINT, []DataType{UNIT}, []reflect.Type{reflect.TypeOf(int64(0))},
		"Relative Strength Index (RSI)", "A momentum oscillator that measures the speed and change of price movements, indicating overbought or oversold conditions.",
		"#614C97",
	},

	Asset.RSI2: {
		func(priceUSDA, priceUSDB float64) int8 {
			return 2
		},
		Asset.RSI2, POINT, []DataType{-1}, []reflect.Type{reflect.TypeOf(""), reflect.TypeOf(int64(0))},
		"Relative Strength Index 2 (RSI)", "A momentum oscillator that measures the speed and change of price movements, indicating overbought or oversold conditions.",
		"#614C97",
	},
}

type AvailableAssetJSON struct {
	AssetType       AssetType `json:"asset_type"`
	DataType        DataType  `json:"data_type"`
	DefaultDecimals int8      `json:"default_decimals"`

	Dependencies  []DataType `json:"dependencies"`
	ArgumentTypes []string   `json:"argument_types"`
	Label         string     `json:"label"`
	Description   string     `json:"description"`
	Color         string     `json:"color"`

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

		dataType := v.DataType

		ret = append(ret, AvailableAssetJSON{
			AssetType:           v.ID,
			DataType:            v.DataType,
			DataTypeName:        When[string](dataType == -1).Then("").Else(dataType.String()),
			DataTypeColor:       When[string](dataType == -1).Then("").Else(dataType.Color()),
			DataTypeColumns:     When[[]ColumnName](dataType == -1).Then([]ColumnName{}).Else(dataType.Columns()),
			DataTypeDescription: When[string](dataType == -1).Then("").Else(dataType.Description()),
			Dependencies:        v.RequiredDependencyDataTypes,
			ArgumentTypes:       argumentTypes,
			Label:               v.Label,
			Description:         v.Description,
			Color:               v.Color,
		})
	}
	return ret
}
