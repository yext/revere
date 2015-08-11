/*
Revere is an alerting system for Graphite.

TODO(eefi): Detailed usage documentation.
*/
package main

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"strings"
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
			fmt.Println(err)
			continue
		}
		for subprobe, reading := range readings {
			if reading.State != revere.Normal {
				headers := make(map[string]string)
				headers["To"] = strings.Join(emails, ", ")
				headers["Subject"] = "Revere reported unhealthy state for " + subprobe

				b := new(bytes.Buffer)
				for k, v := range headers {
					fmt.Fprintf(b, "%s: %s\r\n", k, v)
				}
				fmt.Fprintf(b, "\r\nProbe %s reported unhealthy state with message: \n\n%s",
					subprobe, reading.Details.Text())

				err = smtp.SendMail(mailServer, nil, sender, emails, b.Bytes())
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
