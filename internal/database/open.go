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
	q.Add("_pragma", "synchronous(0)")
	q.Add("_pragma", "journal_mode(MEMORY)")
	q.Add("_pragma", "foreign_keys(1)")
	dsn := "file:" + path + "?" + q.Encode()

	// Connect database
	db, err = sqlx.Connect("sqlite", dsn)
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
		ddlCreateDocumentToken,
		ddlCreateDocumentTokenIndexToken}

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
	id       INT  NOT NULL,
	arabic   TEXT NOT NULL,
	PRIMARY KEY (id))`

const ddlCreateDocumentToken = `
CREATE TABLE IF NOT EXISTS document_token (
	document_id INT  NOT NULL,
	token       TEXT NOT NULL,
	start       INT  NOT NULL,
	end         INT  NOT NULL,
	CONSTRAINT token_document_fk FOREIGN KEY (document_id) REFERENCES document (id))`

const ddlCreateDocumentTokenIndexToken = `
CREATE INDEX IF NOT EXISTS document_token_idx_token ON document_token (token)`
