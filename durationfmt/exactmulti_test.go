package durationfmt

import (
	"fmt"
	"testing"
	"time"
)

var exactMultiTests = []test{
	{0, "0 s"},

	{999 * time.Nanosecond, "999 ns"},

	{1 * time.Microsecond, "1 us"},
	{1*time.Microsecond + 1*time.Nanosecond, "1.001 us"},
	{1*time.Microsecond + 100*time.Nanosecond, "1.1 us"},
	{999*time.Microsecond + 999*time.Nanosecond, "999.999 us"},

	{1 * time.Millisecond, "1 ms"},
	{1*time.Millisecond + 1*time.Nanosecond, "1.000001 ms"},
	{1*time.Millisecond + 100*time.Microsecond, "1.1 ms"},
	{999*time.Millisecond + 999*time.Microsecond + 999*time.Nanosecond, "999.999999 ms"},

	{1 * time.Second, "1 s"},
	{1*time.Second + 1*time.Nanosecond, "1.000000001 s"},
	{1*time.Second + 100*time.Millisecond, "1.1 s"},

	{1 * time.Minute, "1 min"},
	{1*time.Minute + 1*time.Nanosecond, "1 min, 0.000000001 s"},
	{1*time.Minute + 1*time.Second, "1 min, 1 s"},

	{1 * time.Hour, "1 hr"},
	{1*time.Hour + 1*time.Second, "1 hr, 1 s"},
	{1*time.Hour + 1*time.Minute, "1 hr, 1 min"},
	{1*time.Hour + 1*time.Minute + 1*time.Second, "1 hr, 1 min, 1 s"},

	{24 * time.Hour, "1 d"},
	{25 * time.Hour, "1 d, 1 hr"},

	{7 * 24 * time.Hour, "1 wk"},
	{(7 + 1) * 24 * time.Hour, "1 wk, 1 d"},

	{365 * 24 * time.Hour, "1 yr"},
	{(365 + 7 + 1) * 24 * time.Hour, "1 yr, 1 wk, 1 d"},

	{
		-(12*time.Hour + 34*time.Minute + 56*time.Second + 789*time.Millisecond),
		"-12 hr, 34 min, 56.789 s",
	},
}

func TestExactMulti(t *testing.T) {
	f := ExactMulti()
	testFormatter(t, f, "ExactMulti", exactMultiTests)

	const (
		minYears       = 292
		minWeeks       = 24
		minDays        = 3
		minHours       = 23
		minMinutes     = 47
		minSeconds     = 16
		minNanoseconds = 854775808
	)
	testFormatter(t, f, "ExactMulti", []test{{
		minDuration,
		fmt.Sprintf("-%d yr, %d wk, %d d, %d hr, %d min, %d.%09d s",
			minYears, minWeeks, minDays, minHours, minMinutes,
			minSeconds, minNanoseconds),
	}})
	const calcMinDuration = -minYears*time.Hour*24*365 +
		-minWeeks*time.Hour*24*7 +
		-minDays*time.Hour*24 +
		-minHours*time.Hour +
		-minMinutes*time.Minute +
		-minSeconds*time.Second +
		-minNanoseconds*time.Nanosecond
	if calcMinDuration != minDuration {
		t.Errorf("ExactMulti().Format(minDuration) returns %s, want %s",
			calcMinDuration, minDuration)
	}
}
