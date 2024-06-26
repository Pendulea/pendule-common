package pcommon

type parserRPCRequests struct{}

type StatusHTML struct {
	AssetID string `json:"asset_id"`
	HTML    string `json:"html"`
}

type CSVStatus struct {
	BuildID        string     `json:"build_id"`
	RequestTime    int64      `json:"request_time"`
	Status         string     `json:"status"`
	Size           int64      `json:"size"`
	Percent        float64    `json:"percent"`
	From           TimeUnit   `json:"from"`
	To             TimeUnit   `json:"to"`
	TimeframeLabel string     `json:"timeframe_label"`
	Assets         [][]string `json:"assets"`
}

type GetStatusResponse struct {
	CountPendingTasks  int          `json:"count_pending_tasks"`
	CountRunningTasks  int          `json:"count_running_tasks"`
	CSVStatuses        []CSVStatus  `json:"csv_statuses"`
	HTMLStatuses       []StatusHTML `json:"html_statuses"`
	CPUCount           int          `json:"cpu_count"`
	AvailableMemory    uint64       `json:"available_memory"`
	AvailableDiskSpace uint64       `json:"available_disk_space"`
	MinTimeframe       int64        `json:"min_timeframe"`
}

type GetSetListsResponse struct {
	SetList []SetJSON `json:"set_list"`
}

func (rpc *parserRPCRequests) FetchAvailableSetList(parserRPCClient *RPCClient) ([]SetJSON, error) {
	if err := parserRPCClient.CheckConnectedError(); err != nil {
		return nil, err
	}

	res, err := parserRPCClient.Request("GetSetList", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	setList := GetSetListsResponse{}
	Format.DecodeMapIntoStruct(res.Data, &setList)
	return setList.SetList, nil
}

func (rpc *parserRPCRequests) FetchStatus(parserRPCClient *RPCClient) (*GetStatusResponse, error) {
	if err := parserRPCClient.CheckConnectedError(); err != nil {
		return nil, err
	}

	res, err := parserRPCClient.Request("GetStatus", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	ret := GetStatusResponse{}
	Format.DecodeMapIntoStruct(res.Data, &ret)
	return &ret, nil
}
