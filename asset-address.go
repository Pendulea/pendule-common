package pcommon

import (
	"crypto"
	"fmt"
	"strconv"
	"strings"
)

type AssetAddress string

func (adp AssetAddressParsed) IsValid() error {

	for _, id := range adp.SetID {
		id = strings.TrimSpace(id)
		if !isAlphanumeric(id) || id == "" {
			return fmt.Errorf("id contains non-alphanumeric characters: %s", strings.Join(adp.SetID, "_"))
		}
	}

	if _, ok := AssetTypeMap[string(adp.AssetType)]; !ok {
		return fmt.Errorf("asset type is invalid")
	}

	for at := range ArchivesIndex {
		for _, assetType := range at.GetTargetedAssets() {
			if adp.AssetType == assetType {
				if len(adp.Dependencies) > 0 || len(adp.Arguments) > 0 {
					return fmt.Errorf("asset %s has dependencies or arguments but it should not", adp.AssetType)
				}
			}
		}
	}

	config := DEFAULT_ASSETS[adp.AssetType]
	if len(adp.Dependencies) != len(config.RequiredDependencyDataTypes) {
		return fmt.Errorf("dependencies length mismatch")
	}
	for i, depAddr := range adp.Dependencies {
		dep, err := depAddr.Parse()
		if err != nil {
			return err
		}
		if err := dep.IsValid(); err != nil {
			return err
		}
		if DEFAULT_ASSETS[dep.AssetType].DataType != config.RequiredDependencyDataTypes[i] {
			return fmt.Errorf("invalid dependency type")
		}
	}

	if len(adp.Arguments) != len(config.RequiredArgumentTypes) {
		return fmt.Errorf("arguments length mismatch")
	}
	for i, arg := range adp.Arguments {
		requiredArgType := config.RequiredArgumentTypes[i].String()
		if strings.HasPrefix(requiredArgType, "int") {
			if _, err := strconv.ParseInt(arg, 10, 64); err != nil {
				return fmt.Errorf("invalid argument type")
			}
		} else if strings.HasPrefix(requiredArgType, "bool") {
			if _, err := strconv.ParseBool(arg); err != nil {
				return fmt.Errorf("invalid argument type")
			}
		} else if strings.HasPrefix(requiredArgType, "float") {
			if _, err := strconv.ParseFloat(arg, 64); err != nil {
				return fmt.Errorf("invalid argument type")
			}
		} else {
			if requiredArgType != "string" {
				return fmt.Errorf("not supported argument type")
			}
		}
	}

	return nil

}

// setID is a list of alpa-numeric strings (string array is for easier comparison)
func (adp AssetAddressParsed) BuildAddress() AssetAddress {

	setIDStr := strings.ToLower(strings.Join(adp.SetID, "_"))
	argumentsStr := strings.Join(adp.Arguments, "_")

	dependenciesSliceStr := make([]string, len(adp.Dependencies))
	for i, dep := range adp.Dependencies {
		dependenciesSliceStr[i] = string(dep)
	}
	dependenciesStr := strings.Join(dependenciesSliceStr, "=")

	assetAddress := fmt.Sprintf("%s;%s;[%s];%s", setIDStr, adp.AssetType, dependenciesStr, argumentsStr)
	return AssetAddress(assetAddress)
}

type AssetAddressParsed struct {
	SetID        []string       `json:"set_id"`
	AssetType    AssetType      `json:"asset_type"`
	Dependencies []AssetAddress `json:"dependencies"`
	Arguments    []string       `json:"arguments"`
}

func (address AssetAddress) Parse() (*AssetAddressParsed, error) {
	// Step 1: Find the main parts
	parts, err := splitMainParts(string(address))
	if err != nil {
		return nil, err
	}

	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid asset address format")
	}

	setIDStr := parts[0]
	assetTypeStr := parts[1]
	dependenciesStr := parts[2]
	argumentsStr := strings.TrimSpace(parts[3])

	dependencies, err := parseDependencies(dependenciesStr)
	if err != nil {
		return nil, err
	}

	var arguments []string = nil
	if argumentsStr != "" {
		arguments = strings.Split(argumentsStr, "_")
	}

	return &AssetAddressParsed{
		SetID:        strings.Split(setIDStr, "_"),
		AssetType:    AssetType(assetTypeStr),
		Arguments:    arguments,
		Dependencies: dependencies,
	}, nil
}

func (address AssetAddress) Sha256() []byte {
	//createe a sha of the string address
	h := crypto.SHA256.New()
	h.Write([]byte(address))
	return h.Sum(nil)
}

func splitMainParts(address string) ([]string, error) {
	var parts []string
	var currentPart strings.Builder
	var bracketLevel int

	for i, r := range address {
		switch r {
		case ';':
			if bracketLevel == 0 {
				parts = append(parts, currentPart.String())
				currentPart.Reset()
				continue
			}
		case '[':
			bracketLevel++
		case ']':
			bracketLevel--
			if bracketLevel < 0 {
				return nil, fmt.Errorf("unbalanced brackets at position %d", i)
			}
		}
		currentPart.WriteRune(r)
	}

	if bracketLevel != 0 {
		return nil, fmt.Errorf("unbalanced brackets")
	}

	parts = append(parts, currentPart.String())
	return parts, nil
}

func parseDependencies(dependenciesStr string) ([]AssetAddress, error) {
	if dependenciesStr == "[]" {
		return nil, nil
	}

	dependenciesStr = strings.Trim(dependenciesStr, "[]")
	dependenciesSliceStr, err := splitDependencies(dependenciesStr)
	if err != nil {
		return nil, err
	}

	dependencies := make([]AssetAddress, len(dependenciesSliceStr))
	for i, dep := range dependenciesSliceStr {
		dependencies[i] = AssetAddress(dep)
	}
	return dependencies, nil
}

func splitDependencies(dependenciesStr string) ([]string, error) {
	var parts []string
	var currentPart strings.Builder
	var bracketLevel int

	for i, r := range dependenciesStr {
		switch r {
		case '=':
			if bracketLevel == 0 {
				parts = append(parts, currentPart.String())
				currentPart.Reset()
				continue
			}
		case '[':
			bracketLevel++
		case ']':
			bracketLevel--
			if bracketLevel < 0 {
				return nil, fmt.Errorf("unbalanced brackets at position %d", i)
			}
		}
		currentPart.WriteRune(r)
	}

	if bracketLevel != 0 {
		return nil, fmt.Errorf("unbalanced brackets")
	}

	parts = append(parts, currentPart.String())
	return parts, nil
}
