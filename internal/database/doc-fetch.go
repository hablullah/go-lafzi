package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// FetchDocuments returns documents with specified IDs.
func FetchDocuments(db *sqlx.DB, ids ...int) ([]Document, error) {
	// If there are no IDs submitted, stop early
	if len(ids) == 0 {
		return nil, nil
	}

	// Prepare query
	query, args, err := sqlx.In(`
		SELECT id, content FROM document
		WHERE id IN (?)`, ids)
	if err != nil {
		return nil, err
	}

	// Fetch data
	result := []Document{}
	err = db.Select(&result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return result, nil
}
