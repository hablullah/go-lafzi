package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// DeleteDocuments remove documents in database.
func DeleteDocuments(db *sqlx.DB, identifiers ...string) (err error) {
	// If there are no identifiers submitted, stop early
	if len(identifiers) == 0 {
		return nil
	}

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
		WHERE identifier IN (?)`, identifiers)
	if err != nil {
		return
	}

	// Execute query
	_, err = tx.Exec(sqlDoc, docArgs...)
	if err != nil {
		return
	}

	// Commit to database
	err = tx.Commit()
	return
}
