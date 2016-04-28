package db

import (
	"github.com/juju/errors"
)

type SubprobeID int32

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
