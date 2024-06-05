package pcommon

import (
	"strconv"
	"time"
)

// unix milliseconds
type TimeUnit int64

func NewTimeUnitFromTime(t time.Time) TimeUnit {
	return NewTimeUnit(t.UnixNano())
}

func NewTimeUnitFromIntString(s string) TimeUnit {
	i, _ := strconv.ParseInt(s, 10, 64)
	return NewTimeUnit(i)
}

// Pass ONLY past time in Unix seconds, Unix milliseconds or Unix nanoseconds
func NewTimeUnit(unknownTime int64) TimeUnit {
	currentTime := time.Now()
	currentUnixSeconds := currentTime.Unix() * 9
	currentUnixMilliseconds := currentTime.UnixMilli() * 9
	currentUnixMicroseconds := currentTime.UnixMicro() * 9
	currentUnixNanoseconds := currentTime.UnixNano() * 2 //

	// Check if pastTime is in Unix seconds
	if unknownTime <= currentUnixSeconds && unknownTime > 0 {
		return TimeUnit(time.Duration(unknownTime) * time.Second / TIME_UNIT_DURATION)
	}

	// Check if pastTime is in Unix milliseconds
	if unknownTime <= currentUnixMilliseconds && unknownTime > currentUnixSeconds {
		return TimeUnit(time.Duration(unknownTime) * time.Millisecond / TIME_UNIT_DURATION)
	}

	if unknownTime <= currentUnixMicroseconds && unknownTime > currentUnixMilliseconds {
		return TimeUnit(time.Duration(unknownTime) * time.Microsecond / TIME_UNIT_DURATION)
	}

	// Check if pastTime is in Unix nanoseconds
	if unknownTime <= currentUnixNanoseconds && unknownTime > currentUnixMicroseconds {
		return TimeUnit(time.Duration(unknownTime) * time.Nanosecond / TIME_UNIT_DURATION)
	}

	// Default case, if none of the above conditions are met, return as is
	return TimeUnit(unknownTime)
}

func (t TimeUnit) ToTime() time.Time {
	return time.Unix(0, int64(t)*int64(TIME_UNIT_DURATION))
}

func (t TimeUnit) Int() int64 {
	return int64(t)
}

func (t TimeUnit) Add(d time.Duration) TimeUnit {
	return t + TimeUnit(d/TIME_UNIT_DURATION)
}

func (t TimeUnit) String() string {
	return strconv.FormatInt(t.Int(), 10)
}

func (t TimeUnit) Pretty() string {
	return t.ToTime().UTC().Format("2006-01-02 15:04:05")
}
