package pcommon

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

type SetSettings struct {
	Assets   []AssetSettings  `json:"assets"`
	ID       []string         `json:"id"`
	Settings map[string]int64 `json:"settings"`
}

func (s SetSettings) IDString() string {
	return strings.ToLower(strings.Join(s.ID, ""))
}

func (s SetSettings) DBPath() string {
	return filepath.Join(Env.DATABASES_DIR, strings.ToLower(s.IDString()))
}

func (s SetSettings) ContainsAssetAddress(address AssetAddress) bool {
	for _, asset := range s.Assets {
		if asset.Address.AddSetID(s.ID).BuildAddress() == address {
			return true
		}
	}
	return false
}

func (s *SetSettings) HasSettingValue(id string) int64 {
	if s.Settings == nil {
		return 0
	}
	if _, ok := s.Settings[id]; ok {
		return s.Settings[id]
	}
	return 0
}

func (s *SetSettings) IsValid() error {

	// Check for empty ID and invalid characters
	for _, id := range s.ID {
		id = strings.TrimSpace(id)
		if !isAlphanumeric(id) || id == "" {
			return fmt.Errorf("id contains non-alphanumeric characters: %s", strings.Join(s.ID, "_"))
		}
	}

	// Check for duplicate settings
	settingsFound := map[string]bool{}
	for key := range s.Settings {
		if settingsFound[key] {
			return fmt.Errorf("duplicate setting: %s", key)
		}
		settingsFound[key] = true
	}

	// Check for duplicate asset addresses and validate assets
	addresseFound := map[AssetAddress]bool{}
	for _, asset := range s.Assets {

		if err := asset.IsValid(*s); err != nil {
			return err
		}
		assetAddress := asset.Address.AddSetID(s.ID).BuildAddress()
		if _, ok := addresseFound[assetAddress]; ok {
			return fmt.Errorf("duplicate asset address: %s", assetAddress)
		}
		addresseFound[assetAddress] = true
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

			assetAddress := asset.Address.AddSetID(s.ID)
			if !assetAddress.HasArguments() && !assetAddress.HasDependencies() {
				if lo.IndexOf(supportedAsset, assetAddress.AssetType) == -1 {
					return fmt.Errorf("unsupported asset")
				}
				if asset.Decimals > 12 {
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
	r.Settings = make(map[string]int64, len(s.Settings))
	for k, v := range s.Settings {
		r.Settings[k] = v
	}
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
