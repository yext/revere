package test

import (
	"testing"
	"time"

	"github.com/yext/revere"
)

var startTime = time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC)

var testDurations = map[time.Duration]string{
	0:                     "0 min.",
	10 * time.Minute:      "10 min.",
	75 * time.Minute:      "75 min.",
	20 * time.Hour:        "20 hours",
	30 * time.Hour:        "30 hours",
	(24 + 8) * time.Hour:  "1.5 days",
	(24 + 12) * time.Hour: "1.5 days",
	(24 + 18) * time.Hour: "2 days",
	(48 + 1) * time.Hour:  "2 days",
}

func TestDurationStrings(t *testing.T) {
	for d, e := range testDurations {
		a := revere.GetFmtEnteredState(startTime, startTime.Add(d))
		if e != a {
			t.Errorf("Duration: %v Expected: %s, Actual: %s", d, e, a)
		}
	}
}
