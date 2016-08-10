package durationfmt

import (
	"fmt"
	"time"
)

// MostSigUnit returns a formatter that produces a rounded representation of a
// duration based on its largest significant unit of time
func MostSigUnit() Formatter {
	return mostSigUnit{}
}

type mostSigUnit struct{}

func (f mostSigUnit) Format(d time.Duration) string {
	if d <= 75*time.Minute {
		return fmt.Sprintf("%d min.", int(d.Minutes()))
	}

	if d <= 30*time.Hour {
		return fmt.Sprintf("%d hours", int(d.Hours()))
	}

	days := int(d.Hours()) / 24
	r := int(d.Hours()) % 24
	if r < 6 {
		return fmt.Sprintf("%d days", days)
	}
	if 6 <= r && r < 18 {
		return fmt.Sprintf("%d.5 days", days)
	}

	return fmt.Sprintf("%d days", days+1)
}
