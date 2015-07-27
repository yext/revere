/*
Revere is an alerting system for Graphite.

TODO(eefi): Detailed usage documentation.
*/
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/yext/revere"
	"github.com/yext/revere/probes"
)

func main() {
	probeSettings := os.Args[1]
	emails := os.Args[2:]

	probe := probes.NewGraphiteThreshold(probeSettings)

	ticker := time.Tick(5 * time.Minute)
	for _ = range ticker {
		readings := probe.Check()
		for subprobe, reading := range readings {
			if reading.State != revere.Normal {
				// TODO(eefi): Email.
				fmt.Printf("%v %v\n", subprobe, emails)
			}
		}
	}
}
