package vm

import (
	"database/sql"
	"fmt"

	"github.com/yext/revere"
)

var (
	allMonitors []*Monitor
)

type Silence struct {
	*revere.Silence
}

func NewSilence(db *sql.DB, id int) (*Silence, error) {
	silence, err := revere.LoadSilence(db, uint(id))
	if err != nil {
		return nil, err
	}
	if silence == nil {
		return nil, fmt.Errorf("Error loading silence with id: %d", id)
	}

	return newSilence(silence), nil
}

func BlankSilence(db *sql.DB) (*Silence, error) {
	silence := new(revere.Silence)

	return newSilence(silence), nil
}

func newSilence(s *revere.Silence) *Silence {
	viewmodel := new(Silence)
	viewmodel.Silence = s
	return viewmodel
}

func AllSilences(db *sql.DB) ([]*Silence, error) {
	revereSilences, err := revere.LoadSilences(db)
	if err != nil {
		return nil, err
	}

	silences := make([]*Silence, len(revereSilences))
	for i, revereSilence := range revereSilences {
		silences[i] = newSilence(revereSilence)
	}

	return silences, nil
}
