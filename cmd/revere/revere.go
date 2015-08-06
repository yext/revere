/*
Revere is an alerting system for Graphite.

TODO(eefi): Detailed usage documentation.
*/
package main

import (
	"fmt"
	"net/smtp"
	"os"
	"time"

	"github.com/yext/revere"
	"github.com/yext/revere/probes"
)

const (
	mailServer string = "localhost:25"
	sender     string = "revere@yext.com"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Not enough arguments.")
		return
	}
	probeSettings := os.Args[1]
	emails := os.Args[2:]

	probe, err := probes.NewGraphiteThreshold(probeSettings)
	if err != nil {
		fmt.Println(err)
		return
	}

	ticker := time.Tick(5 * time.Minute)
	for _ = range ticker {
		readings, err := probe.Check()
		if err != nil {
			fmt.Print(err)
		}
		for subprobe, reading := range readings {
			if reading.State != revere.Normal {
				message := "Subject: Revere reported unhealthy state for " +
					subprobe +
					"\n" +
					"Probe " +
					subprobe +
					" reported unhealthy state with message: \n\n" +
					reading.Details.Text()
				smtp.SendMail(
					mailServer,
					nil,
					sender,
					emails,
					[]byte(message))
			}
		}
	}
}
