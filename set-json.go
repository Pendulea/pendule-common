package pcommon

import (
	"encoding/json"
	"fmt"
)

type SetJSON struct {
	Settings SetSettings `json:"settings"`
	Size     int64       `json:"size"`
	Assets   []AssetJSON `json:"assets"`
	Type     SetType     `json:"type"`
}

func (s SetJSON) PrintJSON() {
	jsonData, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}
