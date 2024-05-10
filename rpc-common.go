package pcommon

type SetJSON struct {
	Pair            Pair     `json:"pair"`
	Inconsistencies []string `json:"inconsistencies"`
	SetSize         int64    `json:"set_size"`
	Timeframes      []int64  `json:"timeframes"`
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
