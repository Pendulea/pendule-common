package pcommon

type parserRPCRequests struct{}

func (rpc *parserRPCRequests) FetchCandleList(parserRPCClient *RPCClient, settings GetCandlesRequest) (TickTimeArray, error) {
	if err := parserRPCClient.CheckConnectedError(); err != nil {
		return nil, err
	}

	reqMap, err := Format.EncodeStructIntoMap(&settings)
	if err != nil {
		return nil, err
	}
	res, err := parserRPCClient.Request("GetCandles", reqMap)
	if err != nil {
		return nil, err
	}
	responseJSON := GetCandlesResponse{}

	if err := Format.DecodeMapIntoStruct(res.Data, &responseJSON); err != nil {
		return nil, err
	}

	return responseJSON.Candles, nil
}

func (rpc *parserRPCRequests) FetchAvailablePairSetList(parserRPCClient *RPCClient) ([]SetJSON, error) {
	if err := parserRPCClient.CheckConnectedError(); err != nil {
		return nil, err
	}

	res, err := parserRPCClient.Request("GetSetLists", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	setList := GetSetListsResponse{}
	Format.DecodeMapIntoStruct(res.Data, &setList)
	return setList.SetList, nil
}
