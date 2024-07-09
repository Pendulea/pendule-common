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
		p = state.buildRSI(dataList[0].(UnitTime), b.cachedArguments[0].(int64))
	}

	if p == nil {
		return nil, errors.New("not implemented")
	}

	if b.Precision >= 0 {
		p.Value = Math.RoundFloat(p.Value, uint(b.Precision))
	}

	return p, nil
}

/* RELATIVE STRENGTH INDEX : RSI */
type rsiState struct {
	AvgGain   float64 `json:"avg_gain"`
	AvgLoss   float64 `json:"avg_loss"`
	LastRSI   float64 `json:"last_rsi"`
	LastClose float64 `json:"prev_close"`
	Pos       int64   `json:"pos"`
}

func (state *rsiState) buildRSI(newTick UnitTime, RSI_PERIOD int64) *Point {
	defer func() {
		state.LastClose = newTick.Close
	}()

	if state.LastClose <= 0 {
		return &Point{Value: -1}
	}

	change := newTick.Close - state.LastClose

	var gain, loss float64
	if change > 0 {
		gain = change
		loss = 0.0
	} else {
		gain = 0.0
		loss = -change
	}

	state.Pos += 1

	if state.Pos <= RSI_PERIOD {
		// For the first 14 ticks, just accumulate the gains and losses
		state.AvgGain += gain
		state.AvgLoss += loss
		if state.Pos < RSI_PERIOD {
			return &Point{-1} // Not enough data to compute RSI yet
		}
	}

	if state.Pos == RSI_PERIOD {
		// Calculate the average gain and loss for the first 14 ticks
		state.AvgGain /= float64(RSI_PERIOD)
		state.AvgLoss /= float64(RSI_PERIOD)
	} else {
		state.AvgGain = (state.AvgGain*(float64(RSI_PERIOD-1)) + gain) / float64(RSI_PERIOD)
		state.AvgLoss = (state.AvgLoss*(float64(RSI_PERIOD-1)) + loss) / float64(RSI_PERIOD)
	}

	if state.AvgLoss == 0 {
		state.LastRSI = 100
	} else {
		rs := state.AvgGain / state.AvgLoss
		state.LastRSI = float64(100 - (100 / (1 + rs)))
	}

	return &Point{state.LastRSI}
}
