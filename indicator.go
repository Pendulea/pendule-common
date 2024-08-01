package pcommon

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type IndicatorDataBuilder struct {
	assetType AssetType
	prevState []byte
	arguments []string
	Precision int8

	cachedParsedState interface{}
	cachedArguments   []interface{}
}

func NewIndicatorDataBuilder(assetType AssetType, prevState []byte, arguments []string, precision int8) *IndicatorDataBuilder {
	b := &IndicatorDataBuilder{
		assetType: assetType,
		prevState: prevState,
		arguments: arguments,
		Precision: precision,

		cachedParsedState: nil,
		cachedArguments:   nil,
	}

	indicatorConfig := DEFAULT_ASSETS[assetType]
	ret := make([]interface{}, len(arguments))
	for i, arg := range arguments {
		requiredArgType := indicatorConfig.RequiredArgumentTypes[i].String()
		if strings.HasPrefix(requiredArgType, "int") {
			ret[i], _ = strconv.ParseInt(arg, 10, 64)
		} else if strings.HasPrefix(requiredArgType, "bool") {
			ret[i], _ = strconv.ParseBool(arg)
		} else if strings.HasPrefix(requiredArgType, "float") {
			ret[i], _ = strconv.ParseFloat(arg, 64)
		} else {
			ret[i] = arg
		}
	}
	b.cachedArguments = ret
	return b
}

func parseIndicatorState[T any](b *IndicatorDataBuilder) (*T, error) {
	if b.cachedParsedState != nil {
		return b.cachedParsedState.(*T), nil
	}

	data := b.prevState

	var zeroValue T

	if len(data) == 0 {
		// Check if T is a struct type with reflect, then instantiate it and return
		t := reflect.TypeOf(zeroValue)
		if t.Kind() == reflect.Struct {
			b.cachedParsedState = zeroValue
			return &zeroValue, nil
		}
		return nil, errors.New("provided type is not a struct")
	}

	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	b.cachedParsedState = &result
	return b.cachedParsedState.(*T), nil
}

func (b *IndicatorDataBuilder) saveState(state interface{}) {
	data, err := json.Marshal(state)
	if err != nil {
		log.Fatal(err)
	}
	b.prevState = data
	b.cachedParsedState = state
}

func (b *IndicatorDataBuilder) PrevState() []byte {
	if b.prevState == nil {
		return nil
	}
	stateCopy := make([]byte, len(b.prevState))
	copy(stateCopy, b.prevState)
	return stateCopy
}

func (b *IndicatorDataBuilder) ComputeUnsafe(dataList ...Data) (*Point, error) {

	var p *Point = nil

	if b.assetType == Asset.RSI {
		state, err := parseIndicatorState[rsiState](b)
		if err != nil {
			return nil, err
		}
		defer b.saveState(state)
		p = state.buildRSI(dataList[0].(UnitTime).Close, b.cachedArguments[0].(int64))
	}

	if b.assetType == Asset.RSI2 {
		state, err := parseIndicatorState[rsiState](b)
		if err != nil {
			return nil, err
		}
		defer b.saveState(state)
		column := b.cachedArguments[0].(string)
		v, err := dataList[0].ValueAt(ColumnName(column))
		if err != nil {
			return nil, err
		}
		p = state.buildRSI(v, b.cachedArguments[1].(int64))
	}

	if b.assetType == Asset.SMA || b.assetType == Asset.EMA || b.assetType == Asset.WMA || b.assetType == Asset.HMA {
		var state *maState = nil
		if b.prevState == nil && b.cachedParsedState == nil {
			s := newEmptyMAState(int(b.cachedArguments[1].(int64)))
			state = &s
		} else {
			var err error
			state, err = parseIndicatorState[maState](b)
			if err != nil {
				return nil, err
			}
		}

		defer b.saveState(state)
		column := b.cachedArguments[0].(string)
		v, err := dataList[0].ValueAt(ColumnName(column))
		if err != nil {
			return nil, err
		}
		switch b.assetType {
		case Asset.SMA:
			p = state.buildSMA(v, int(b.cachedArguments[1].(int64)))
		case Asset.EMA:
			p = state.buildEMA(v, int(b.cachedArguments[1].(int64)))
		case Asset.WMA:
			p = state.buildWMA(v, int(b.cachedArguments[1].(int64)))
		case Asset.HMA:
			p = state.buildHMA(v, int(b.cachedArguments[1].(int64)))
		}
	}

	if p == nil {
		return nil, errors.New("not implemented")
	}

	if b.Precision >= 0 {
		p.Value = Math.RoundFloat(p.Value, uint(b.Precision))
	}

	return p, nil
}
