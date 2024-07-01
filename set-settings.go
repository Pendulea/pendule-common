package pcommon

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

type SetLittleSetting struct {
	ID    string `json:"id"`
	Value int64  `json:"value"`
}

type SetSettings struct {
	Assets   []AssetSettings    `json:"assets"`
	ID       []string           `json:"id"`
	Settings []SetLittleSetting `json:"settings"`
}

func (s SetSettings) IDString() string {
	return strings.ToLower(strings.Join(s.ID, ""))
}

func (s SetSettings) DBPath() string {
	return filepath.Join(Env.DATABASES_DIR, strings.ToLower(s.IDString()))
}

func (s SetSettings) ContainsAssetAddress(address AssetAddress) bool {
	for _, asset := range s.Assets {
		if asset.Address.BuildAddress() == address {
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

func (s *SetSettings) IsValid() error {

	for _, id := range s.ID {
		id = strings.TrimSpace(id)
		if !isAlphanumeric(id) || id == "" {
			return fmt.Errorf("id contains non-alphanumeric characters: %s", strings.Join(s.ID, "_"))
		}
	}

	addresseFound := map[AssetAddress]bool{}
	for _, asset := range s.Assets {

		if err := asset.Address.IsValid(); err != nil {
			return err
		}

		assetAddress := asset.Address.BuildAddress()

		if _, ok := addresseFound[assetAddress]; ok {
			return fmt.Errorf("duplicate asset address: %s", assetAddress)
		}

		_, err := Format.StrDateToDate(asset.MinDataDate)
		if err != nil {
			return err
		}
		if asset.Decimals < 0 || asset.Decimals > 12 {
			return fmt.Errorf("decimals out of range: %d", asset.Decimals)
		}

		settingsFound := map[string]bool{}
		for _, setting := range s.Settings {
			if settingsFound[setting.ID] {
				return fmt.Errorf("duplicate setting: %s", setting.ID)
			}
			settingsFound[setting.ID] = true
		}
	}

	return nil
}

func (s *SetSettings) GetSetType() (SetType, error) {

	if s.IsValid() != nil {
		return -1, fmt.Errorf("invalid set settings")
	}

	if err := s.IsBinancePair(); err == nil {
		return BINANCE_PAIR, nil
	}

	return -1, fmt.Errorf("type does not exist")

}

func (s *SetSettings) IsBinancePair() error {
	if len(s.ID) == 2 {
		denominatorPair := strings.ToUpper(s.ID[1])
		if !strings.HasPrefix(denominatorPair, "USDC") && !strings.HasPrefix(denominatorPair, "USDT") {
			return fmt.Errorf("unsupported denominator pair")
		}
		if s.HasSettingValue("binance") != 1 {
			return fmt.Errorf("no binance settings")
		}
		supportedAsset := BINANCE_PAIR.GetSupportedAssets()
		for _, asset := range s.Assets {

			if len(asset.Address.Dependencies) == 0 && len(asset.Address.Arguments) == 0 {
				if lo.IndexOf(supportedAsset, asset.Address.AssetType) == -1 {
					return fmt.Errorf("unsupported asset")
				}
				if asset.Decimals > 6 {
					return fmt.Errorf("decimals out of range")
				}
			}
		}
		return nil
	}

	return fmt.Errorf("not a binance pair")
}

func (s *SetSettings) Copy() *SetSettings {
	var r SetSettings

	r.ID = s.ID
	r.Settings = make([]SetLittleSetting, len(s.Settings))
	copy(r.Settings, s.Settings)
	r.Assets = make([]AssetSettings, len(s.Assets))
	copy(r.Assets, s.Assets)
	return &r
}

func (s *SetSettings) BuildArchiveFolderPath(asset AssetType) string {
	return fmt.Sprintf("%s/%s/%s", Env.ARCHIVES_DIR, s.IDString(), asset)
}

func (s *SetSettings) BuildArchiveFilePath(asset AssetType, date string, ext string) string {
	return fmt.Sprintf("%s/%s.%s", s.BuildArchiveFolderPath(asset), date, ext)
}
