package db

import (
	"fmt"
	"strings"

	"github.com/juju/errors"

	"github.com/yext/revere/state"
)

func (db *DB) create() error {
	// Necessary to make foreign key names unique database-wide.
	noDBPrefix := db.prefix[strings.Index(db.prefix, ".")+1:]

	for _, table := range createTables {
		query := fmt.Sprintf("CREATE TABLE pfx_%s (", table.name)
		query += strings.Join(table.rowsAndKeys, ", ")
		query += ") ENGINE=InnoDB CHARACTER SET=utf8mb4"
		query = strings.Replace(query, "nodbpfx_", noDBPrefix, -1)
		query = cq(db, query)

		_, err := db.Exec(query)
		if err != nil {
			return errors.Maskf(err, "create table %s", table.name)
		}
	}
	for i, query := range createExtra {
		_, err := db.Exec(cq(db, query))
		if err != nil {
			return errors.Maskf(err, "run query %d", i)
		}
	}
	return nil
}

var createTables = []struct {
	name        string
	rowsAndKeys []string
}{
	{
		name: "monitors",
		rowsAndKeys: []string{
			"monitorid INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY",
			"name VARCHAR(30) NOT NULL",
			"owner VARCHAR(60) NOT NULL DEFAULT ''",
			"description TEXT NOT NULL",
			"response TEXT NOT NULL",
			"probetype SMALLINT NOT NULL",
			"probe TEXT NOT NULL",
			"changed DATETIME NOT NULL",
			"version INTEGER NOT NULL DEFAULT 0",
			"archived DATETIME DEFAULT NULL",
			"KEY idx_changed (changed)",
		},
	},
	{
		name: "subprobes",
		rowsAndKeys: []string{
			"subprobeid INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY",
			"monitorid INTEGER UNSIGNED NOT NULL",
			"name VARCHAR(150) NOT NULL",
			"archived DATETIME DEFAULT NULL",
			"UNIQUE KEY idx_monitorid_name (monitorid, name)",
			"CONSTRAINT nodbpfx_subprobes_fk_monitorid FOREIGN KEY (monitorid) REFERENCES pfx_monitors (monitorid) ON DELETE CASCADE",
		},
	},
	{
		name: "subprobe_statuses",
		rowsAndKeys: []string{
			"subprobeid INTEGER UNSIGNED PRIMARY KEY",
			"recorded DATETIME NOT NULL",
			"state TINYINT NOT NULL",
			"silenced BOOLEAN NOT NULL DEFAULT FALSE",
			"enteredstate DATETIME NOT NULL",
			"lastnormal DATETIME NOT NULL",
			"KEY idx_state_silenced_enteredstate_recorded (state, silenced, enteredstate, recorded)",
			"CONSTRAINT nodbpfx_subprobe_statuses_fk_subprobeid FOREIGN KEY (subprobeid) REFERENCES pfx_subprobes (subprobeid) ON DELETE CASCADE",
		},
	},
	{
		name: "readings",
		rowsAndKeys: []string{
			"readingid BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY",
			"subprobeid INTEGER UNSIGNED NOT NULL",
			"recorded DATETIME NOT NULL",
			"state TINYINT NOT NULL",
			"KEY idx_subprobeid_recorded_readingid (subprobeid, recorded, readingid)",
			"CONSTRAINT nodbpfx_readings_fk_subprobeid FOREIGN KEY (subprobeid) REFERENCES pfx_subprobes (subprobeid) ON DELETE CASCADE",
		},
	},
	{
		name: "triggers",
		rowsAndKeys: []string{
			"triggerid INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY",
			fmt.Sprintf("level TINYINT NOT NULL DEFAULT %d", state.Error),
			"triggeronexit BOOLEAN NOT NULL DEFAULT TRUE",
			"periodmilli INTEGER NOT NULL DEFAULT 0",
			"targettype SMALLINT NOT NULL",
			"target TEXT NOT NULL",
		},
	},
	{
		name: "monitor_triggers",
		rowsAndKeys: []string{
			"monitorid INTEGER UNSIGNED NOT NULL",
			"subprobes TEXT NOT NULL",
			"triggerid INTEGER UNSIGNED NOT NULL",
			"PRIMARY KEY (monitorid, triggerid)",
			"UNIQUE KEY idx_triggerid (triggerid)",
			"CONSTRAINT nodbpfx_monitor_triggers_fk_monitorid FOREIGN KEY (monitorid) REFERENCES pfx_monitors (monitorid) ON DELETE CASCADE",
			"CONSTRAINT nodbpfx_monitor_triggers_fk_triggerid FOREIGN KEY (triggerid) REFERENCES pfx_triggers (triggerid) ON DELETE CASCADE",
		},
	},
	{
		name: "labels",
		rowsAndKeys: []string{
			"labelid INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY",
			"name VARCHAR(30) NOT NULL",
			"description TEXT NOT NULL",
		},
	},
	{
		name: "label_triggers",
		rowsAndKeys: []string{
			"labelid INTEGER UNSIGNED NOT NULL",
			"triggerid INTEGER UNSIGNED NOT NULL",
			"PRIMARY KEY (labelid, triggerid)",
			"UNIQUE KEY idx_triggerid (triggerid)",
			"CONSTRAINT nodbpfx_label_triggers_fk_labelid FOREIGN KEY (labelid) REFERENCES pfx_labels (labelid) ON DELETE CASCADE",
			"CONSTRAINT nodbpfx_label_triggers_fk_triggerid FOREIGN KEY (triggerid) REFERENCES pfx_triggers (triggerid) ON DELETE CASCADE",
		},
	},
	{
		name: "labels_monitors",
		rowsAndKeys: []string{
			"labelid INTEGER UNSIGNED NOT NULL",
			"monitorid INTEGER UNSIGNED NOT NULL",
			"subprobes TEXT NOT NULL",
			"PRIMARY KEY (labelid, monitorid)",
			"KEY idx_monitorid (monitorid)",
			"CONSTRAINT nodbpfx_labels_monitors_fk_labelid FOREIGN KEY (labelid) REFERENCES pfx_labels (labelid) ON DELETE CASCADE",
			"CONSTRAINT nodbpfx_labels_monitors_fk_monitorid FOREIGN KEY (monitorid) REFERENCES pfx_monitors (monitorid) ON DELETE CASCADE",
		},
	},
	{
		name: "silences",
		rowsAndKeys: []string{
			"silenceid INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY",
			"monitorid INTEGER UNSIGNED NOT NULL",
			"subprobes TEXT NOT NULL",
			"start DATETIME NOT NULL",
			"end DATETIME NOT NULL",
			"KEY idx_monitorid_end_start (monitorid, end, start)",
			"CONSTRAINT nodbpfx_silences_fk_monitorid FOREIGN KEY (monitorid) REFERENCES pfx_monitors (monitorid) ON DELETE CASCADE",
		},
	},
	{
		name: "resources",
		rowsAndKeys: []string{
			"resourceid INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY",
			"resourcetype SMALLINT NOT NULL",
			"resource TEXT NOT NULL",
		},
	},
	{
		name: "settings",
		rowsAndKeys: []string{
			"settingid INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY",
			"settingtype SMALLINT NOT NULL",
			"setting TEXT NOT NULL",
		},
	},
	{
		name: "schema_history",
		rowsAndKeys: []string{
			"version TINYINT UNSIGNED PRIMARY KEY",
			"migrationstarted DATETIME NOT NULL",
			"migrationcompleted DATETIME",
		},
	},
}

var createExtra = []string{
	"INSERT INTO pfx_schema_history (version, migrationstarted, migrationcompleted) " +
		"VALUES (1, UTC_TIMESTAMP(), UTC_TIMESTAMP())",
}
