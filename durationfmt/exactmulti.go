package durationfmt

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// ExactMulti returns a formatter that produces an exact accounting of durations
// divided into standard units down to seconds.
//
// The larger units used ignore leap seconds and leap days.
func ExactMulti() Formatter {
	return exactMulti{}
}

type exactMulti struct{}

func (f exactMulti) Format(d time.Duration) string {
	if d == 0 {
		return "0 s"
	}

	if d == math.MinInt64 {
		return "-292 yr, 24 wk, 3 d, 23 hr, 47 min, 16.854775808 s"
	}

	if d < 0 {
		return "-" + f.Format(-d)
	}

	if d < time.Microsecond {
		return fmt.Sprintf("%d ns", d.Nanoseconds())
	}
	if d < time.Millisecond {
		micro := d / time.Microsecond
		frac := formatFraction(int64(d%time.Microsecond), 3)
		return fmt.Sprintf("%d%s us", micro, frac)
	}
	if d < time.Second {
		milli := d / time.Millisecond
		frac := formatFraction(int64(d%time.Millisecond), 6)
		return fmt.Sprintf("%d%s ms", milli, frac)
	}

	var parts []string

	years, extra := d/(365*24*time.Hour), d%(365*24*time.Hour)
	if years > 0 {
		parts = append(parts, fmt.Sprintf("%d yr", years))
	}

	weeks, extra := extra/(7*24*time.Hour), extra%(7*24*time.Hour)
	if weeks > 0 {
		parts = append(parts, fmt.Sprintf("%d wk", weeks))
	}

	days, extra := extra/(24*time.Hour), extra%(24*time.Hour)
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d d", days))
	}

	hours, extra := extra/time.Hour, extra%time.Hour
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d hr", hours))
	}

	minutes, extra := extra/time.Minute, extra%time.Minute
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d min", minutes))
	}

	if extra > 0 {
		seconds, extra := extra/time.Second, extra%time.Second
		frac := formatFraction(int64(extra), 9)
		parts = append(parts, fmt.Sprintf("%d%s s", seconds, frac))
	}

	return strings.Join(parts, ", ")
}
