package pcommon

import (
	"encoding/json"
	"math"
)

type MAState struct {
	Sum      float64
	Buffer   []float64
	Pos      int
	Count    int
	EMA      float64
	WMABuff  []float64
	AltState [][]byte
}

func newEmptyMAState(period int) MAState {
	return MAState{
		Buffer:  make([]float64, period),
		WMABuff: make([]float64, period),
	}
}

func (state *MAState) buildMA(value float64, period int, Type string) float64 {
	switch Type {
	case "SMA":
		return state.buildSMA(value, period)
	case "EMA":
		return state.buildEMA(value, period)
	case "WMA":
		return state.buildWMA(value, period)
	case "HMA":
		return state.buildHMA(value, period)
	default:
		return state.buildSMA(value, period)
	}
}

func (state *MAState) buildSMA(value float64, SMA_PERIOD int) float64 {
	state.Count++

	if state.Count <= SMA_PERIOD {
		state.Sum += value
		state.Buffer[state.Pos] = value
		state.Pos = (state.Pos + 1) % SMA_PERIOD

		if state.Count == SMA_PERIOD {
			return state.Sum / float64(SMA_PERIOD)
		}
		return -1
	}

	// Full buffer, replace the oldest value
	state.Sum = state.Sum - state.Buffer[state.Pos] + value
	state.Buffer[state.Pos] = value
	state.Pos = (state.Pos + 1) % SMA_PERIOD
	return state.Sum / float64(SMA_PERIOD)
}

func (state *MAState) buildEMA(value float64, period int) float64 {
	if state.Count == 0 {
		state.EMA = value
		state.Count++
		return state.EMA
	}

	k := 2.0 / float64(period+1)
	today := value
	state.EMA = today*k + state.EMA*(1-k)
	state.Count++
	return state.EMA
}

func (state *MAState) buildWMA(value float64, WMA_PERIOD int) float64 {
	state.WMABuff[state.Pos] = value
	state.Pos = (state.Pos + 1) % WMA_PERIOD
	state.Count++

	if state.Count < WMA_PERIOD {
		return -1
	}

	weightedSum := 0.0
	weight := 0.0
	for i := 0; i < WMA_PERIOD; i++ {
		w := float64(i + 1)
		weightedSum += state.WMABuff[(state.Pos+i)%WMA_PERIOD] * w
		weight += w
	}
	return weightedSum / weight
}

func (state *MAState) buildHMA(value float64, HMA_PERIOD int) float64 {
	sqrtPeriod := int(math.Sqrt(float64(HMA_PERIOD)))

	if len(state.AltState) == 0 {
		halfState := newEmptyMAState(HMA_PERIOD / 2)
		fullState := newEmptyMAState(HMA_PERIOD)
		wmaSqrtState := newEmptyMAState(sqrtPeriod)

		hb, _ := json.Marshal(halfState)
		fb, _ := json.Marshal(fullState)
		wb, _ := json.Marshal(wmaSqrtState)
		state.AltState = [][]byte{hb, fb, wb}
	}

	halfMaState, _ := parseIndicatorState[MAState](&IndicatorDataBuilder{prevState: state.AltState[0]})
	fullMaState, _ := parseIndicatorState[MAState](&IndicatorDataBuilder{prevState: state.AltState[1]})

	wmaHalfPeriod := halfMaState.buildWMA(value, HMA_PERIOD/2)
	wmaFullPeriod := fullMaState.buildWMA(value, HMA_PERIOD)

	//serializing the state
	state.AltState[0], _ = json.Marshal(halfMaState)
	state.AltState[1], _ = json.Marshal(fullMaState)

	state.Count++
	if state.Count < HMA_PERIOD {
		return -1
	}

	wmaSqrtState, _ := parseIndicatorState[MAState](&IndicatorDataBuilder{prevState: state.AltState[2]})
	close := (2 * float64(wmaHalfPeriod)) - float64(wmaFullPeriod)
	hmaValue := wmaSqrtState.buildWMA(close, sqrtPeriod)
	state.AltState[2], _ = json.Marshal(wmaSqrtState)
	return hmaValue
}
