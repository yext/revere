package durationfmt

import (
	"math"
	"testing"
	"time"
)

const (
	minDuration time.Duration = math.MinInt64
)

type test struct {
	d        time.Duration
	expected string
}

func testFormatter(t *testing.T, f Formatter, name string, tests []test) {
	for _, test := range tests {
		actual := f.Format(test.d)
		if actual != test.expected {
			t.Errorf("%s().Format(%s) == %q, want %q",
				name, test.d, actual, test.expected)
		}
	}
}
