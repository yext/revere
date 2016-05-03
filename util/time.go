package util

import (
	"time"
)

func GetPeriodAndType(periodMilli int64) (int64, string) {
	ms := time.Duration(periodMilli) * time.Millisecond
	switch {
	case ms == 0:
		return 0, ""
	case ms%(time.Hour*24) == 0:
		return int64(ms / (time.Hour * 24)), "day"
	case ms%time.Hour == 0:
		return int64(ms / time.Hour), "hour"
	case ms%time.Minute == 0:
		return int64(ms / time.Minute), "minute"
	case ms%time.Second == 0:
		return int64(ms / time.Second), "second"
	default:
		return 0, ""
	}
}

func GetMs(Period int64, PeriodType string) int64 {
	switch PeriodType {
	case "day":
		return (Period * int64(time.Hour) * 24) / int64(time.Millisecond)
	case "hour":
		return (Period * int64(time.Hour)) / int64(time.Millisecond)
	case "minute":
		return (Period * int64(time.Minute)) / int64(time.Millisecond)
	case "second":
		return (Period * int64(time.Second)) / int64(time.Millisecond)
	default:
		return 0
	}
}
