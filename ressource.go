package pcommon

type RessourcesJSON struct {
	AvailableAssets   []AvailableAssetJSON `json:"available_assets"`
	AvailableSetTypes []SetTypeJSON        `json:"available_set_types"`
}

func BuildRessources() RessourcesJSON {
	res := RessourcesJSON{}
	for s := range SET_ARCHIVES {
		res.AvailableSetTypes = append(res.AvailableSetTypes, s.JSON())
	}
	res.AvailableAssets = DEFAULT_ASSETS.JSON()
	return res
}
