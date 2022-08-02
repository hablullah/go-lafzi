package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// DeleteDocuments remove documents in database.
func DeleteDocuments(db *sqlx.DB, ids ...int) (err error) {
	// Start transaction
	tx, err := db.Beginx()
	if err != nil {
		err = fmt.Errorf("failed to start transaction: %v", err)
		return
	}

	// Make sure to rollback if error ever happened
	defer func() {
		if err != nil && tx != nil {
			tx.Rollback()
		}
	}()

	// Prepare query
	sqlDoc, docArgs, err := sqlx.In(`
		DELETE FROM document
		WHERE id IN (?)`, ids)
	if err != nil {
		return
	}

	sqlToken, tokenArgs, err := sqlx.In(`
		DELETE FROM token
		WHERE document_id IN (?)`, ids)
	if err != nil {
		return
	}

	// Execute query
	_, err = tx.Exec(sqlToken, tokenArgs...)
	if err != nil {
		return
	}

	_, err = tx.Exec(sqlDoc, docArgs...)
	if err != nil {
		return
	}

	// Commit to database
	err = tx.Commit()
	return
}
