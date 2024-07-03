package pcommon

import "fmt"

type AssetAddressParsedWithoutSetID struct {
	AssetType    AssetType      `json:"asset_type"`
	Dependencies []AssetAddress `json:"dependencies"`
	Arguments    []string       `json:"arguments"`
}

type AssetSettings struct {
	Address     AssetAddressParsedWithoutSetID `json:"address"`
	MinDataDate string                         `json:"min_data_date"`
}

func (adp AssetAddressParsedWithoutSetID) AddSetID(setID []string) AssetAddressParsed {
	return AssetAddressParsed{
		SetID:        setID,
		AssetType:    adp.AssetType,
		Dependencies: adp.Dependencies,
		Arguments:    adp.Arguments,
	}
}

func (as AssetSettings) IsValid(setSettings SetSettings) error {

	assetAddress := as.Address.AddSetID(setSettings.ID)
	if err := assetAddress.IsValid(); err != nil {
		return err
	}

	config := DEFAULT_ASSETS[as.Address.AssetType]
	if config.DefaultDecimals < 0 && config.SetUpDecimals == nil {
		return fmt.Errorf("config decimals not set for asset type %s", as.Address.AssetType)
	}

	if !assetAddress.HasDependencies() {
		_, err := Format.StrDateToDate(as.MinDataDate)
		if err != nil {
			return err
		}
	} else if as.MinDataDate != "" {
		return fmt.Errorf("min_data_date should be empty when dependencies are present")
	}

	return nil
}
