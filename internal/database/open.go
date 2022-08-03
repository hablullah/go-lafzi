package database

import (
	"fmt"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
)

// Open SQLite database in specified path.
func Open(path string) (db *sqlx.DB, err error) {
	// Prepare DSN
	q := url.Values{}
	q.Add("_sync", "0")
	q.Add("_journal", "MEMORY")
	q.Add("_foreign_keys", "1")
	dsn := "file:" + path + "?" + q.Encode()

	// Connect database
	db, err = sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	db.SetConnMaxLifetime(time.Minute)

	// Create transaction
	var tx *sqlx.Tx
	tx, err = db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}

	// If error ever happened, rollback and close database
	defer func() {
		if err != nil {
			if tx != nil {
				tx.Rollback()
			}
			if db != nil {
				db.Close()
			}
			db = nil
		}
	}()

	// Generate tables
	ddlQueries := []string{
		ddlCreateDocument,
		ddlCreateToken,
		ddlCreateTokenIndex}

	for _, query := range ddlQueries {
		_, err = tx.Exec(query)
		if err != nil {
			return
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return
	}

	return
}

const ddlCreateDocument = `
CREATE TABLE IF NOT EXISTS document (
	id          INT  NOT NULL,
	content     TEXT NOT NULL,
	PRIMARY KEY (id))`

const ddlCreateToken = `
CREATE TABLE IF NOT EXISTS token (
	key         TEXT NOT NULL,
	document_id INT  NOT NULL,
	PRIMARY KEY (key, document_id),
	CONSTRAINT token_document_fk FOREIGN KEY (document_id) REFERENCES document (id))`

const ddlCreateTokenIndex = `
CREATE INDEX IF NOT EXISTS token_key_idx ON token (key)`
