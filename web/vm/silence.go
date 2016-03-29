package vm

import (
	"database/sql"
	"fmt"

	"github.com/yext/revere"
)

type Silence struct {
	*revere.Silence
	AllMonitors []*Monitor
}

func NewSilence(db *sql.DB, id int) (*Silence, error) {
	viewmodel, err := baseSilence(db)
	if err != nil {
		return nil, err
	}

	silence, err := revere.LoadSilence(db, uint(id))
	if err != nil {
		return nil, err
	}
	if silence == nil {
		return nil, fmt.Errorf("Error loading silence with id: %d", id)
	}

	viewmodel.Silence = silence

	return viewmodel, nil
}

func BlankSilence(db *sql.DB) (*Silence, error) {
	viewmodel, err := baseSilence(db)
	if err != nil {
		return nil, err
	}
	viewmodel.Silence = new(revere.Silence)

	return viewmodel, nil
}

func baseSilence(db *sql.DB) (*Silence, error) {
	viewmodel := new(Silence)

	var err error
	viewmodel.AllMonitors, err = AllMonitors(db)
	if err != nil {
		return nil, err
	}

	return viewmodel, nil
}
