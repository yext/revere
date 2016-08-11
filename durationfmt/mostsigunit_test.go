package durationfmt

import (
	"testing"
	"time"
)

var mostSigUnitTests = []test{
	{0, "0 min."},

	{10 * time.Minute, "10 min."},
	{75 * time.Minute, "75 min."},

	{20 * time.Hour, "20 hours"},
	{30 * time.Hour, "30 hours"},

	{(24 + 8) * time.Hour, "1.5 days"},
	{(24 + 12) * time.Hour, "1.5 days"},

	{(24 + 18) * time.Hour, "2 days"},
	{(48 + 1) * time.Hour, "2 days"},
}

func TestMostSigUnit(t *testing.T) {
	f := MostSigUnit()
	testFormatter(t, f, "MostSigUnit", mostSigUnitTests)
}
