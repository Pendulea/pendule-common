package pcommon

type TickData struct {
	Data any      `json:"data"`
	Time TimeUnit `json:"time"`
}

type TickDataArray []TickData

type Index struct {
	Indicator string        `json:"indicator"`
	Indexes   TickDataArray `json:"indexes"`
}
