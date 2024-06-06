package pcommon

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

type format struct{}

var Format = format{}

// ExtractDateFromTradeZipFile extracts date from filename
func (f format) ExtractDateFromTradeZipFile(filename string) (string, error) {
	regex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})\.zip$`)
	matches := regex.FindStringSubmatch(filename)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", errors.New("no match found")
}

func (f format) TimeFrameToLabel(timeFrame time.Duration) (string, error) {
	if timeFrame > MAX_TIME_FRAME {
		return "", errors.New("time frame is too large")
	}
	if timeFrame < Env.MIN_TIME_FRAME {
		return "", errors.New("time frame is too small")
	}
	if timeFrame%Env.MIN_TIME_FRAME != 0 {
		return "", fmt.Errorf("time frame must be a multiple of %d seconds", Env.MIN_TIME_FRAME/1000)
	}

	if timeFrame%WEEK == 0 {
		return fmt.Sprintf("%dw", int64(timeFrame/WEEK)), nil
	}
	if timeFrame%DAY == 0 {
		return fmt.Sprintf("%dd", int64(timeFrame/DAY)), nil
	}
	if timeFrame%time.Hour == 0 {
		return fmt.Sprintf("%dh", int64(timeFrame.Hours())), nil
	}
	if timeFrame%time.Minute == 0 {
		return fmt.Sprintf("%dm", int64(timeFrame.Minutes())), nil
	}
	if timeFrame%time.Second == 0 {
		return fmt.Sprintf("%ds", int64(timeFrame.Seconds())), nil
	}
	return fmt.Sprintf("%dms", timeFrame.Milliseconds()), nil
}

// StrDateToDate converts a string date to a time.Time object
func (f format) StrDateToDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02T15:04:05Z"
	return time.Parse(layout, dateStr+"T00:00:00Z")
}

// FormatDateStr formats a time.Time object to a string
func (f format) FormatDateStr(date time.Time) string {
	return date.UTC().Format("2006-01-02")
}

// BuildDateStr computes a date string from days ago
func (f format) BuildDateStr(daysAgo int) string {
	now := time.Now()
	pastDate := now.AddDate(0, 0, -daysAgo)
	return f.FormatDateStr(pastDate)
}

// LargeBytesToShortString converts byte size to a human-readable string
func (f format) LargeBytesToShortString(b int64) string {
	switch {
	case b >= 1_000_000_000:
		return fmt.Sprintf("%.2fgb", float64(b)/1_000_000_000)
	case b >= 1_000_000:
		return fmt.Sprintf("%.1fmb", float64(b)/1_000_000)
	case b >= 1_000:
		return fmt.Sprintf("%dkb", b/1_000)
	default:
		return fmt.Sprintf("%db", b)
	}
}

func (f format) LargeNumberToShortString(n int64) string {
	if n >= 1_000_000_000 {
		return fmt.Sprintf("%.2fb", float64(n)/1_000_000_000)
	}
	if n >= 1_000_000 {
		return fmt.Sprintf("%.2fm", float64(n)/1_000_000)
	}
	if n >= 1_000 {
		return fmt.Sprintf("%.1fk", float64(n)/1_000)
	}

	return fmt.Sprintf("%d", n)
}

// AccurateHumanize provides a human-readable representation of time duration in milliseconds
func (f format) AccurateHumanize(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Second*30 {
		return fmt.Sprintf("%.1fs", float64(d.Milliseconds())/1000)
	}
	if d < time.Minute {
		return fmt.Sprintf("%dsecs", int(d.Seconds()))
	}
	if d < time.Hour {
		min := int(d.Minutes())
		sec := int(d.Seconds()) % 60
		if sec < 10 {
			return fmt.Sprintf("%dm0%ds", min, sec)
		}
		return fmt.Sprintf("%dm%ds", min, sec)
	}
	if d < DAY {
		hour := int(d.Hours())
		min := int(d.Minutes()) % 60
		if min < 10 {
			return fmt.Sprintf("%dh0%dm", hour, min)
		}
		return fmt.Sprintf("%dh%dm", hour, min)
	}
	if d < DAY*4 {
		hour := int(d.Hours())
		return fmt.Sprintf("%dh", hour)
	}
	if d < WEEK*4 {
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	}
	weeks := int(d.Hours() / 24 / 7)
	return fmt.Sprintf("%dw", weeks)
}

func (f format) BuildSetID(symbol string, tradeType TradeType) string {
	return symbol + tradeType.Key()
}

func (f format) CuteHash(s string) string {
	var h int
	for i, c := range s {
		h += i + int(c)
	}
	return fmt.Sprintf("%010d", h)
}

func (f format) ArrayDurationToArrInt64(durations []time.Duration) []int64 {
	arr := make([]int64, len(durations))
	for i, d := range durations {
		arr[i] = d.Milliseconds()
	}
	return arr
}

func (f format) DecodeMapIntoStruct(data map[string]interface{}, result interface{}) error {
	decoderConfig := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   result,
		TagName:  "json", // Set the tag name used in your structs if any, default is none
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return err
	}

	return decoder.Decode(data)
}

func (f format) EncodeStructIntoMap(data interface{}) (map[string]interface{}, error) {
	// First, marshal the struct into JSON bytes
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Then, unmarshal JSON bytes into a map
	var resultMap map[string]interface{}
	err = json.Unmarshal(jsonData, &resultMap)
	if err != nil {
		return nil, err
	}

	return resultMap, nil
}

func (f format) Float(val float64, precision int8) string {
	// Format the float with the specified precision
	str := strconv.FormatFloat(val, 'f', int(precision), 64)

	// Remove trailing zeros and the decimal point if it's not needed
	str = strings.TrimRight(str, "0")
	str = strings.TrimRight(str, ".")

	if str == "" {
		return "0"
	}

	return str
}
