package pcommon

/* RELATIVE STRENGTH INDEX : RSI */
type rsiState struct {
	AvgGain   float64 `json:"avg_gain"`
	AvgLoss   float64 `json:"avg_loss"`
	LastRSI   float64 `json:"last_rsi"`
	LastClose float64 `json:"prev_close"`
	Pos       int64   `json:"pos"`
}

func (state *rsiState) buildRSI(value float64, RSI_PERIOD int64) *Point {
	defer func() {
		state.LastClose = value
	}()

	if state.LastClose <= 0 {
		return &Point{Value: -1}
	}

	change := value - state.LastClose

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
