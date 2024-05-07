package pcommon

type EmptyStruct struct{}

type RPCAction struct {
	Id      string                 `json:"id"`
	Method  string                 `json:"method"`
	Payload map[string]interface{} `json:"payload"`
}

type RPCResponse struct {
	Id    string                 `json:"id"`
	Data  map[string]interface{} `json:"data"`
	Error string                 `json:"error"`
}

type RPCRequest map[string]interface{}

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
	Descending     bool   `json:"descending"`
}

type GetCandlesResponse struct {
	Candles TickTimeArray `json:"candles"`
}
