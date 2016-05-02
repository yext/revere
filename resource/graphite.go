package resource

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/juju/errors"
)

// Graphite represents a remote Graphite server. For more information, see
// http://graphite.readthedocs.org/ .
type Graphite struct {
	// Base is the URL where the server can be found. Base should have a
	// trailing slash, with the expectation that the render API endpoint is
	// at Base + "render".
	Base string
}

// GraphiteSeries encapsulates the data returned by Graphite for a particular
// series. The values in Values correspond to the times between Start and End,
// inclusive, that are Step times apart. Null values are represented by NaN.
//
// XXX(eefi): It is possible to store NaN in Graphite, and encoding nulls as NaN
// makes those cases indistinguishable. However, actual NaNs should be rare, and
// encoding nulls as NaN removes the need for an extra layer of pointer
// indirection.
type GraphiteSeries struct {
	Name       string
	Start, End time.Time
	Step       time.Duration
	Values     []float64
}

// QueryRecent retrieves data stored in Graphite for the most recent d time.
func (g Graphite) QueryRecent(target string, d time.Duration) ([]GraphiteSeries, error) {
	// TODO(eefi): Warn if d is not whole seconds?
	from := fmt.Sprintf("-%ds", int64(d/time.Second))
	data, err := g.query(target, from, "now")
	if err != nil {
		return nil, errors.Trace(err)
	}
	return data, nil
}

// Query retrieves data stored in Graphite for the given time period.
func (g Graphite) Query(target string, from, until time.Time) ([]GraphiteSeries, error) {
	fromString := GraphiteTimestamp(from)
	untilString := GraphiteTimestamp(until)
	data, err := g.query(target, fromString, untilString)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return data, nil
}

func (g Graphite) query(target, from, until string) ([]GraphiteSeries, error) {
	url := g.RenderURL([]string{target}, map[string]string{
		"from":   from,
		"until":  until,
		"format": "raw",
	})
	data, err := g.get(url)
	if err != nil {
		return nil, errors.Maskf(err, "query Graphite")
	}

	series, err := parseGraphiteRawRender(data)
	if err != nil {
		return nil, errors.Maskf(err, "parse Graphite raw render response")
	}

	return series, nil
}

func (g Graphite) get(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Get %s: not-OK HTTP status code: %d", url, r.StatusCode)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Annotatef(err, "Get %s", url)
	}

	return b, nil
}

var graphiteRawRenderSeriesFormat = regexp.MustCompile(
	`^(.*),([0-9]+),([0-9]+),([0-9]+)\|([^|]*)$`)

func parseGraphiteRawRender(data []byte) ([]GraphiteSeries, error) {
	var series []GraphiteSeries
	for _, seriesString := range strings.Split(string(data), "\n") {
		if seriesString == "" {
			continue
		}

		parts := graphiteRawRenderSeriesFormat.FindStringSubmatch(seriesString)
		if parts == nil {
			return nil, errors.Errorf("malformed series: %s", seriesString)
		}

		s := GraphiteSeries{}

		s.Name = parts[1]

		startUnix, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return nil, errors.Errorf("could not parse series start time: %s", parts[2])
		}
		s.Start = time.Unix(startUnix, 0)

		endUnix, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return nil, errors.Errorf("could not parse series end time: %s", parts[3])
		}
		s.End = time.Unix(endUnix, 0)

		stepSecond, err := strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			return nil, errors.Errorf("could not parse series step: %s", parts[4])
		}
		s.Step = time.Duration(stepSecond) * time.Second

		valueStrings := strings.Split(parts[5], ",")
		s.Values = make([]float64, len(valueStrings))
		for i, valueString := range valueStrings {
			if valueString == "None" {
				s.Values[i] = math.NaN()
			} else {
				s.Values[i], err = strconv.ParseFloat(valueString, 64)
				if err != nil {
					return nil, errors.Errorf("could not parse datapoint: %s", valueString)
				}
			}
		}

		series = append(series, s)
	}

	return series, nil
}

// RenderURL builds a Graphite render API URL for the given targets. The
// key-value pairs in args are added to the URL.
func (g Graphite) RenderURL(targets []string, args map[string]string) string {
	values := url.Values{"target": targets}
	for k, v := range args {
		values.Add(k, v)
	}

	return g.Base + "render?" + values.Encode()
}

// GraphiteTimestamp converts t to a value suitable for use as the from or until
// argument to the Graphite render API.
func GraphiteTimestamp(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}
