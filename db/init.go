package db

func (db *DB) Init() error {
	// TODO(eefi): Check the schema_history table.
	return db.create()
}
