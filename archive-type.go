package pcommon

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

type ArchiveType string

const BINANCE_SPOT_TRADES ArchiveType = "binance_spot_trades"
const BINANCE_FUTURES_TRADES ArchiveType = "binance_futures_trades"
const BINANCE_BOOK_DEPTH ArchiveType = "binance_book_depth"
const BINANCE_METRICS ArchiveType = "binance_metrics"

func (at ArchiveType) GetArchiveZipPath(date string, set *SetJSON) string {
	archiveDir := filepath.Join(
		Env.ARCHIVES_DIR,
		strings.ToUpper(set.Settings.IDString()),
		"__archives",
	)

	return filepath.Join(archiveDir, string(at), fmt.Sprintf("%s.zip", date))
}

func (at ArchiveType) GetURL(date string, set *SetJSON) (string, error) {

	if err := set.Settings.IsBinancePair(); err == nil {
		symbol := strings.ToUpper(set.Settings.IDString())
		switch at {
		case BINANCE_SPOT_TRADES:
			fileName := fmt.Sprintf("%s-trades-%s.zip", symbol, date)
			return fmt.Sprintf("https://data.binance.vision/data/spot/daily/trades/%s/%s", symbol, fileName), nil
		case BINANCE_FUTURES_TRADES:
			fileName := fmt.Sprintf("%s-trades-%s.zip", symbol, date)
			return fmt.Sprintf("https://data.binance.vision/data/futures/um/daily/trades/%s/%s", symbol, fileName), nil
		case BINANCE_BOOK_DEPTH:
			fileName := fmt.Sprintf("%s-bookDepth-%s.zip", symbol, date)
			return fmt.Sprintf("https://data.binance.vision/data/futures/um/daily/bookDepth/%s/%s", symbol, fileName), nil
		case BINANCE_METRICS:
			fileName := fmt.Sprintf("%s-metrics-%s.zip", symbol, date)
			return fmt.Sprintf("https://data.binance.vision/data/futures/um/daily/metrics/%s/%s", symbol, fileName), nil
		}
	}

	return "", fmt.Errorf("archive type for set")
}

func (at ArchiveType) GetTargetedAssets() []AssetType {
	tree, ok := ArchivesIndex[at]
	if !ok {
		log.Fatal("archive tree not found")
		return nil
	}
	return lo.Map(tree.Columns, func(ab AssetBranch, index int) AssetType {
		return ab.Asset
	})
}

func GetRequiredArchiveType(assetType AssetType) *ArchiveType {
	for arcT := range ArchivesIndex {
		if lo.IndexOf(arcT.GetTargetedAssets(), assetType) != -1 {
			return &arcT
		}
	}
	return nil
}
