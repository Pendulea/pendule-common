package pcommon

import (
	"strconv"
	"time"
)

const TIME_UNIT_DURATION = time.Millisecond

// unix milliseconds
type TimeUnit int64

// Pass ONLY past time in Unix seconds, Unix milliseconds or Unix nanoseconds
func NewTimeUnit(pastTime int64) TimeUnit {
	currentTime := time.Now()
	currentUnixSeconds := currentTime.Unix()
	currentUnixMilliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	currentUnixNanoseconds := currentTime.UnixNano()

	// Check if pastTime is in Unix seconds
	if pastTime <= currentUnixSeconds && pastTime > 0 {
		return TimeUnit(time.Duration(pastTime) * time.Second / TIME_UNIT_DURATION)
	}

	// Check if pastTime is in Unix milliseconds
	if pastTime <= currentUnixMilliseconds && pastTime > currentUnixSeconds*1000 {
		return TimeUnit(time.Duration(pastTime) * time.Millisecond / TIME_UNIT_DURATION)
	}

	// Check if pastTime is in Unix nanoseconds
	if pastTime <= currentUnixNanoseconds && pastTime > currentUnixMilliseconds*1000 {
		return TimeUnit(time.Duration(pastTime) * time.Nanosecond / TIME_UNIT_DURATION)
	}

	// Default case, if none of the above conditions are met, return as is
	return TimeUnit(pastTime)
}

func (t TimeUnit) ToTime() time.Time {
	return time.Unix(0, int64(t)*int64(TIME_UNIT_DURATION))
}

func (t TimeUnit) Int() int64 {
	return int64(t)
}

func TimeToUnit(t time.Time) TimeUnit {
	return TimeUnit(t.UnixNano() / int64(TIME_UNIT_DURATION))
}

func (t TimeUnit) Add(d time.Duration) TimeUnit {
	return t + TimeUnit(d/time.Millisecond)
}

func (t TimeUnit) String() string {
	return strconv.FormatInt(t.Int(), 10)
}

func (t TimeUnit) Pretty() string {
	return t.ToTime().Format("2006-01-02 15:04:05")
}
