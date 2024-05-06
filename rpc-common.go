package pcommon

type EmptyStruct struct{}

type Action struct {
	Id      string                 `json:"id"`
	Method  string                 `json:"method"`
	Payload map[string]interface{} `json:"payload"`
}

type Response struct {
	Id    string `json:"id"`
	Data  any    `json:"data"`
	Error string `json:"error"`
}

type Request map[string]interface{}

type SetJSON struct {
	Pair       Pair    `json:"pair"`
	Consistent bool    `json:"consistent"`
	Timeframes []int64 `json:"timeframes"`
}

// Request represents the parameters of the RPC request.
type RemoveTimeFrameRequest struct {
	Symbol    string `json:"symbol"`
	TimeFrame int64  `json:"timeframe"`
}

type RemoveTimeFrameResponse struct {
	Count int `json:"count"`
}

type IsDateParsedRequest struct {
	SetID     string `json:"set_id"`
	Date      string `json:"date"`
	TimeFrame int64  `json:"timeframe"`
}

type IsDateParsedResponse struct {
	Exist bool `json:"exist"`
}

type GetSetListsResponse struct {
	SetList []SetJSON `json:"set_list"`
}

type AddTimeFrameRequest struct {
	Symbol    string `json:"symbol"`
	TimeFrame int64  `json:"timeframe"`
}

type AddTimeFrameResponse struct {
	Scheduled bool `json:"scheduled"`
}

type GetCandlesRequest struct {
	Limit          int    `json:"limit"`
	OffsetUnixTime int64  `json:"offset_unix_time"`
	Symbol         string `json:"symbol"`
	TimeFrame      int64  `json:"timeframe"`
}

type GetCandlesResponse struct {
	Candles TickMap `json:"candles"`
}

type Tick struct {
	Open                float64 `json:"open"`
	High                float64 `json:"high"`
	Low                 float64 `json:"low"`
	Close               float64 `json:"close"`
	VolumeBought        float64 `json:"volume_bought"`
	VolumeSold          float64 `json:"volume_sold"`
	TradeCount          int64   `json:"trade_count"`
	MedianVolumeBought  float64 `json:"median_volume_bought"`
	AverageVolumeBought float64 `json:"average_volume_bought"`
	MedianVolumeSold    float64 `json:"median_volume_sold"`
	AverageVolumeSold   float64 `json:"average_volume_sold"`
	VWAP                float64 `json:"vwap"`
	StandardDeviation   float64 `json:"standard_deviation"`
}

type TickMap map[int64]Tick
type TickArray []Tick
