package pcommon

type IndicatorJSON struct {
	Indicator     string `json:"indicator"`
	LastIndexTime int64  `json:"last_index_time"`
}

type TimeframeJSON struct {
	Timeframe     int64 `json:"timeframe"`
	LastIndexTime int64 `json:"last_index_time"`
	Indicators    []IndicatorJSON
}

type SetJSON struct {
	Pair            Pair            `json:"pair"`
	Inconsistencies []string        `json:"inconsistencies"`
	SetSize         int64           `json:"set_size"`
	Timeframes      []TimeframeJSON `json:"timeframes"`
}

func (s *SetJSON) IsConsistent() bool {
	return len(s.Inconsistencies) == 0
}

// Request represents the parameters of the RPC request.
type RemoveTimeFrameRequest struct {
	Symbol    string `json:"symbol"`
	TimeFrame int64  `json:"timeframe"`
}

type RemoveTimeFrameResponse struct {
	Scheduled bool `json:"scheduled"`
}

type IsDateParsedRequest struct {
	SetID string `json:"set_id"`
	Date  string `json:"date"`
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

type TickData struct {
	Data any   `json:"data"`
	Time int64 `json:"time"`
}

type TickDataArray []TickData

type Index struct {
	Indicator string        `json:"indicator"`
	Indexes   TickDataArray `json:"indexes"`
}

type GetIndexesResponse struct {
	Timeframe int64         `json:"timeframe"`
	Candles   TickTimeArray `json:"candles"`
	Indexes   []Index       `json:"indexes"`
}

type GetIndexesRequest struct {
	GetCandlesRequest
	Indicators  string `json:"indicators"`
	WithCandles bool   `json:"with_candles"`
}
