package pcommon

import "github.com/samber/lo"

type SetType int8

const BINANCE_PAIR SetType = 1

var SET_ARCHIVES = map[SetType][]ArchiveType{
	BINANCE_PAIR: {BINANCE_SPOT_TRADES, BINANCE_FUTURES_TRADES, BINANCE_BOOK_DEPTH, BINANCE_METRICS},
}

func (st SetType) GetSupportedAssets() []AssetType {
	listArchives := SET_ARCHIVES[st]
	if listArchives == nil {
		return nil
	}
	ret := []AssetType{}
	for _, archive := range listArchives {
		ret = append(ret, archive.GetTargetedAssets()...)
	}

	return lo.Uniq(ret)
}
