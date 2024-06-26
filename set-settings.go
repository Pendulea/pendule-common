package pcommon

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/samber/lo"
)

type SetSettings struct {
	Assets []struct {
		ID          AssetType `json:"id"`
		MinDataDate string    `json:"min_data_date"`
		Decimals    int8      `json:"decimals"`
	} `json:"assets"`
	ID       []string `json:"id"`
	Settings []struct {
		ID    string `json:"id"`
		Value int64  `json:"value"`
	} `json:"settings"`
}

func (s SetSettings) IDString() string {
	return strings.ToLower(strings.Join(s.ID, ""))
}

func (s SetSettings) DBPath() string {
	return filepath.Join(Env.DATABASES_DIR, strings.ToLower(s.IDString()))
}

func (s SetSettings) ContainsAsset(assetID AssetType) bool {
	for _, asset := range s.Assets {
		if asset.ID == assetID {
			return true
		}
	}
	return false
}

func (s *SetSettings) HasSettingValue(id string) int64 {
	if s.Settings == nil {
		return 0
	}
	for _, setting := range s.Settings {
		if setting.ID == id {
			return setting.Value
		}
	}
	return 0
}

func (s *SetSettings) isValid() error {
	if len(s.ID) == 0 {
		return errors.New("ID is empty")
	}

	v := reflect.ValueOf(Asset)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct or a pointer to a struct")
	}

	for _, asset := range s.Assets {
		fieldVal := v.FieldByName(string(asset.ID))
		if !fieldVal.IsValid() {
			return fmt.Errorf("no such field: %s", asset.ID)
		}
		_, err := Format.StrDateToDate(asset.MinDataDate)
		if err != nil {
			return err
		}
		if asset.Decimals < 0 || asset.Decimals > 12 {
			return fmt.Errorf("decimals out of range: %d", asset.Decimals)
		}
	}

	return nil
}

func (s *SetSettings) IsSupportedBinancePair() (bool, []AssetType) {
	d := Asset
	listSupportedAssets := []AssetType{
		d.FUTURES_PRICE, d.SPOT_PRICE,
		d.FUTURES_VOLUME, d.SPOT_VOLUME,
		d.BOOK_DEPTH_P1, d.BOOK_DEPTH_P2, d.BOOK_DEPTH_P3, d.BOOK_DEPTH_P4, d.BOOK_DEPTH_P5, d.BOOK_DEPTH_M1, d.BOOK_DEPTH_M2, d.BOOK_DEPTH_M3, d.BOOK_DEPTH_M4, d.BOOK_DEPTH_M5,
		d.METRIC_SUM_OPEN_INTEREST, d.METRIC_COUNT_TOP_TRADER_LONG_SHORT_RATIO, d.METRIC_SUM_TOP_TRADER_LONG_SHORT_RATIO, d.METRIC_COUNT_LONG_SHORT_RATIO, d.METRIC_SUM_TAKER_LONG_SHORT_VOL_RATIO,
		d.CIRCULATING_SUPPLY,
	}

	if s.isValid() != nil {
		return false, listSupportedAssets
	}

	if len(s.ID) == 3 {
		denominatorPair := strings.ToUpper(s.ID[1])
		if !strings.HasPrefix(denominatorPair, "USDC") && !strings.HasPrefix(denominatorPair, "USDT") {
			return false, listSupportedAssets
		}
		if s.HasSettingValue("binance") != 1 {
			return false, listSupportedAssets
		}
		lastIDKey := s.ID[2]
		if lastIDKey != "_spot" && lastIDKey != "_futures" {
			return false, listSupportedAssets
		}

		for _, asset := range s.Assets {
			if lo.IndexOf(listSupportedAssets, asset.ID) == -1 || asset.Decimals > 6 {
				return false, listSupportedAssets
			}
		}
		return true, listSupportedAssets
	}
	return false, listSupportedAssets
}

func (s *SetSettings) BuildArchiveFolderPath(asset AssetType) string {
	return fmt.Sprintf("%s/%s/%s", Env.ARCHIVES_DIR, s.IDString(), asset)
}

func (s *SetSettings) BuildArchiveFilePath(asset AssetType, date string, ext string) string {
	return fmt.Sprintf("%s/%s.%s", s.BuildArchiveFolderPath(asset), date, ext)
}
