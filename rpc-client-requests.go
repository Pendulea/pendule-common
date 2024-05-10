package pcommon

import "fmt"

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

func (rpc *parserRPCRequests) IsDateParsed(parserRPCClient *RPCClient, settings IsDateParsedRequest) (bool, error) {
	if err := parserRPCClient.CheckConnectedError(); err != nil {
		return false, err
	}

	reqMap, err := Format.EncodeStructIntoMap(&settings)
	if err != nil {
		return false, err
	}

	res, err := parserRPCClient.Request("IsDateParsed", reqMap)
	if err != nil {
		return false, err
	}

	exist := IsDateParsedResponse{}
	Format.DecodeMapIntoStruct(res.Data, &exist)
	return exist.Exist, nil
}

func (rpc *parserRPCRequests) AddTimeframe(parserRPCClient *RPCClient, settings AddTimeFrameRequest) error {
	if err := parserRPCClient.CheckConnectedError(); err != nil {
		return err
	}

	reqMap, err := Format.EncodeStructIntoMap(&settings)
	if err != nil {
		return err
	}

	res, err := parserRPCClient.Request("AddTimeFrame", reqMap)
	if err != nil {
		return err
	}

	resp := AddTimeFrameResponse{}
	Format.DecodeMapIntoStruct(res.Data, &resp)
	if resp.Scheduled {
		return nil
	}
	return fmt.Errorf("failed to schedule")
}

func (rpc *parserRPCRequests) RemoveTimeframe(parserRPCClient *RPCClient, settings RemoveTimeFrameRequest) (int, error) {
	if err := parserRPCClient.CheckConnectedError(); err != nil {
		return 0, err
	}

	reqMap, err := Format.EncodeStructIntoMap(&settings)
	if err != nil {
		return 0, err
	}

	res, err := parserRPCClient.Request("RemoveTimeframe", reqMap)
	if err != nil {
		return 0, err
	}

	resp := RemoveTimeFrameResponse{}
	Format.DecodeMapIntoStruct(res.Data, &resp)
	return resp.Count, nil
}
