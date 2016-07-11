package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/juju/errors"
	"github.com/yext/revere/state"
)

type SubprobeID int32

type Subprobe struct {
	SubprobeID SubprobeID
	MonitorID  MonitorID
	Name       string
	Archived   *time.Time
}

type SubprobeWithStatusInfo struct {
	MonitorID   MonitorID
	MonitorName string
	Name        string
	Archived    *time.Time
	*SubprobeStatus
}

func (db *DB) LoadSubprobe(subprobeID SubprobeID) (*Subprobe, error) {
	var s Subprobe
	err := db.Get(&s, cq(db, "SELECT * FROM pfx_subprobes WHERE subprobeid = ?"), subprobeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Trace(err)
	}
	return &s, nil
}

func (db *DB) LoadSubprobeWithStatusInfo(subprobeID SubprobeID) (*SubprobeWithStatusInfo, error) {
	return loadSubprobeWithStatusInfo(db, subprobeID)
}

func (tx *Tx) LoadSubprobeWithStatusInfo(subprobeID SubprobeID) (*SubprobeWithStatusInfo, error) {
	return loadSubprobeWithStatusInfo(tx, subprobeID)
}

func loadSubprobeWithStatusInfo(dt dbOrTx, subprobeID SubprobeID) (*SubprobeWithStatusInfo, error) {
	results, err := loadSubprobesWithStatus(dt, fmt.Sprintf(`WHERE s.subprobeid = %d`, subprobeID))
	if err != nil || len(results) != 1 {
		return nil, errors.Trace(err)
	}
	return results[0], nil
}

func (tx *Tx) LoadSubprobesByName(monitorID MonitorID) ([]*SubprobeWithStatusInfo, error) {
	return loadSubprobesWithStatus(tx, fmt.Sprintf(
		`WHERE s.monitorid = %d
		ORDER BY s.name`, monitorID))
}

func (tx *Tx) LoadSubprobesBySeverity() ([]*SubprobeWithStatusInfo, error) {
	return loadSubprobesWithStatus(tx, fmt.Sprintf(
		`WHERE ss.state != %d
		ORDER BY ss.state DESC, ss.enteredstate, s.name`, state.Normal))
}

func (tx *Tx) LoadSubprobesBySeverityForLabel(labelID LabelID) ([]*SubprobeWithStatusInfo, error) {
	return loadSubprobesWithStatus(tx, fmt.Sprintf(
		`JOIN pfx_labels_monitors lm USING (monitorid)
		WHERE ss.state != %d AND lm.labelid = %d
		ORDER BY ss.state DESC, ss.enteredstate, s.name`, state.Normal, labelID))
}

func loadSubprobesWithStatus(dt dbOrTx, condition string) ([]*SubprobeWithStatusInfo, error) {
	var subprobes []*SubprobeWithStatusInfo
	q := fmt.Sprintf(
		`SELECT s.monitorid, s.name, s.archived, ss.*, m.name AS monitorname FROM pfx_subprobes s
	     LEFT JOIN pfx_subprobe_statuses ss ON s.subprobeid = ss.subprobeid
		 JOIN pfx_monitors m USING (monitorid)
		 %s`, condition)
	err := dt.Select(&subprobes, cq(dt, q))
	if err != nil {
		return nil, errors.Trace(err)
	}
	return subprobes, nil
}

func (tx *Tx) InsertSubprobe(monitorID MonitorID, name string) (SubprobeID, error) {
	q := `INSERT INTO pfx_subprobes (monitorid, name) VALUES (?, ?)`
	result, err := tx.Exec(cq(tx, q), monitorID, name)
	if err != nil {
		return 0, errors.Trace(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Trace(err)
	}

	return SubprobeID(id), nil
}
