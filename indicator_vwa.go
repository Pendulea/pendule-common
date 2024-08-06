package pcommon

import "errors"

type emptyState struct{}
type Compacted map[TimeUnit][]float64

func (state *emptyState) buildVWA(c Compacted) (*Point, error) {
	var totalVolume float64
	var vwapNumerator float64

	if len(c) == 0 {
		return &Point{Value: 0.0}, nil
	}

	for _, data := range c {
		if len(data) != 2 {
			return nil, errors.New("Invalid data format")
		}
		first := data[0]
		second := data[1]

		vwapNumerator += first * second
		totalVolume += second
	}

	if totalVolume == 0 {
		return &Point{Value: 0.0}, nil
	}

	return &Point{Value: vwapNumerator / totalVolume}, nil
}
