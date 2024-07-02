package pcommon

import (
	"log"

	"github.com/samber/lo"
)

type SetType int8

const BINANCE_PAIR SetType = 1

func (st SetType) String() string {
	switch st {
	case BINANCE_PAIR:
		return "BINANCE PAIR"
	}
	log.Fatal("unknown set type")
	return "UNKNOWN"
}

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

type SetTypeJSON struct {
	Type          SetType           `json:"type"`
	TypeName      string            `json:"type_name"`
	ArchiveChilds []ArchiveTypeJSON `json:"archive_childs"`
}

func (st SetType) JSON() SetTypeJSON {
	archiveChilds := []ArchiveTypeJSON{}
	for _, archive := range SET_ARCHIVES[st] {
		archiveChilds = append(archiveChilds, archive.ToJSON())
	}
	return SetTypeJSON{st, st.String(), archiveChilds}
}
