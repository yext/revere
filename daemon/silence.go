package daemon

import (
	"regexp"

	"github.com/juju/errors"

	"github.com/yext/revere/db"
)

type silence struct {
	subprobes *regexp.Regexp
}

func newSilence(dbSilence db.Silence) (silence, error) {
	subprobes, err := regexp.Compile(dbSilence.Subprobes)
	if err != nil {
		return silence{}, errors.Maskf(err, "compile regexp")
	}
	return silence{subprobes}, nil
}

func (s silence) silences(subprobe *subprobe) bool {
	return s.subprobes.MatchString(subprobe.name)
}
