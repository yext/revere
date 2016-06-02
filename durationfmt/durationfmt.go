// Package durationfmt formats durations for easier human consumption.
package durationfmt

import (
	"fmt"
	"time"
)

// Formatter defines an interface for formatting a duration for easier human
// consumption. Depending on the goals of the particular formatter, Format may
// return a description that is only an approximation of the provided duration.
type Formatter interface {
	Format(time.Duration) string
}

func formatFraction(v int64, places int) string {
	if v == 0 {
		return ""
	}

	for v%10 == 0 {
		v /= 10
		places -= 1
	}

	return fmt.Sprintf(".%0*d", places, v)
}
