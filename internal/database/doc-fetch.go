package database

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// FetchDocuments returns documents with specified IDs.
func FetchDocuments(db *sqlx.DB, ids ...int) (documents []Document, err error) {
	// If there are no IDs submitted, stop early
	if len(ids) == 0 {
		return
	}

	// Start read only transaction
	tx, err := db.Beginx()
	if err != nil {
		err = fmt.Errorf("failed to start transaction: %v", err)
		return
	}

	defer func() {
		if err == sql.ErrNoRows {
			err = nil
		}
		tx.Rollback()
	}()

	// Prepare statement
	stmtGetDoc, err := tx.Preparex(`
		SELECT id, arabic FROM document
		WHERE id = ?`)
	if err != nil {
		return
	}

	// Fetch data
	documents = make([]Document, 0, len(ids))
	for _, id := range ids {
		var doc Document
		err = stmtGetDoc.Get(&doc, id)
		if err != nil && err != sql.ErrNoRows {
			return
		}
		documents = append(documents, doc)
	}

	return
}
