package pcommon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type Pair struct {
	Binance          bool   `json:"binance"`
	Symbol0          string `json:"symbol0"`
	Symbol1          string `json:"symbol1"`
	MinHistoricalDay string `json:"min_historical_day"`
	Indicators       string `json:"indicators"`
	Futures          bool   `json:"futures"`
}

func (p *Pair) TradeType() TradeType {
	if p.Futures {
		return FUTURES_TRADE
	}
	return SPOT_TRADE
}

func (p *Pair) BuildArchiveFolderPath() string {
	path := fmt.Sprintf("%s/%s/%s", Env.ARCHIVES_DIR, p.BuildBinanceSymbol(), p.TradeType().Key())
	return path
}

func (p *Pair) BuildArchivesFilePath(date string, ext string) string {
	fp := p.BuildArchiveFolderPath()
	symbol := p.BuildBinanceSymbol()
	if ext != "csv" && ext != "zip" {
		log.Fatal("invalid extension for archive file")
	}
	return fmt.Sprintf("%s/%s-trades-%s.%s", fp, symbol, date, ext)
}

func (p *Pair) BuildDBPath() string {
	dbPath := filepath.Join(Env.DATABASES_DIR, strings.ToLower(p.BuildSetID()))
	return dbPath
}

func (p *Pair) JSON() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Pair) Copy() Pair {
	return Pair{
		Binance:          p.Binance,
		Symbol0:          p.Symbol0,
		Symbol1:          p.Symbol1,
		Indicators:       p.Indicators,
		MinHistoricalDay: p.MinHistoricalDay,
		Futures:          p.Futures,
	}
}

func (p Pair) BuildSetID() string {
	if p.Binance {
		return Format.BuildSetID(p.BuildBinanceSymbol(), p.TradeType())
	}
	return ""
}

func (p Pair) IsBinanceValid() bool {
	return p.Binance && p.Symbol0 != "" && p.Symbol1 != ""
}

func (p Pair) ErrorFilter(allowedStablePairs []string) error {

	if p.MinHistoricalDay == "" {
		return fmt.Errorf("min_historical_day is required")
	}

	if p.Binance {
		if p.Symbol0 == "" || p.Symbol1 == "" {
			return fmt.Errorf("symbol0 and symbol1 are required for binance pairs")
		}

		if len(allowedStablePairs) > 0 {
			symb1 := strings.ToUpper(p.Symbol1)
			for _, pair := range allowedStablePairs {
				if symb1 == strings.ToUpper(pair) {
					return nil
				}
			}
		}
		return fmt.Errorf("pair %s not allowed for symbol1: only %s", p.Symbol1, strings.Join(allowedStablePairs, ", "))
	}

	return nil
}

func (p Pair) BuildBinanceSymbol() string {
	if p.Binance {
		return strings.ToUpper(p.Symbol0 + p.Symbol1)
	}
	return ""
}

func (pair Pair) BuildBinanceArchiveURL() string {
	symbol := pair.BuildBinanceSymbol()
	if symbol == "" {
		return ""
	}

	date := pair.MinHistoricalDay
	futures := pair.Futures

	fileName := fmt.Sprintf("%s-trades-%s.zip", symbol, date)
	if futures {
		return fmt.Sprintf("https://data.binance.vision/data/futures/um/daily/trades/%s/%s", symbol, fileName)
	} else {
		return fmt.Sprintf("https://data.binance.vision/data/spot/daily/trades/%s/%s", symbol, fileName)
	}
}

func (pair Pair) CheckBinanceSymbolWorks() (bool, error) {
	symbol := pair.BuildBinanceSymbol()
	if symbol == "" {
		return false, nil
	}

	url := pair.BuildBinanceArchiveURL()
	resp, err := http.Head(url) // Perform a HEAD request
	if err != nil {
		return false, err
	}
	defer resp.Body.Close() // Ensure we close the response body
	return resp.StatusCode == 200, nil
}

func (pair Pair) ParseIndicators(supportedIndicatorList []string) []string {
	if pair.Indicators == "" {
		return []string{}
	}
	if pair.Indicators == "*" {
		return supportedIndicatorList
	}

	allowedIndicators := strings.Split(pair.Indicators, ",")
	finalList := make([]string, 0) // Use a separate slice for the results to avoid duplication

	for _, indicator := range allowedIndicators {
		if strings.Contains(indicator, "*") {
			// Handle wildcard entries
			sp := strings.Split(indicator, "*")
			if len(sp) != 2 {
				continue // Skip invalid formats
			}
			prefix, suffix := sp[0], sp[1]
			for _, ind := range supportedIndicatorList {
				if strings.HasPrefix(ind, prefix) && strings.HasSuffix(ind, suffix) {
					finalList = appendIfNotExists(finalList, ind)
				}
			}
		} else {
			// Handle exact matches
			if contains(supportedIndicatorList, indicator) {
				finalList = appendIfNotExists(finalList, indicator)
			}
		}
	}

	return finalList
}

// Helper function to check if a slice contains a specific string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// Append to the slice if the element does not already exist
func appendIfNotExists(slice []string, str string) []string {
	if !contains(slice, str) {
		slice = append(slice, str)
	}
	return slice
}
