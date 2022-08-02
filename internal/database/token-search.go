package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// SearchTokens look for document ids which contains the specified tokens,
// then count how many tokens occured in each document.
func SearchTokens(db *sqlx.DB, keys ...string) ([]TokenLocation, error) {
	// Prepare query
	query, args, err := sqlx.In(`
		SELECT document_id, COUNT(*) n
		FROM token WHERE key IN (?)
		GROUP BY document_id`, keys)
	if err != nil {
		return nil, err
	}

	// Fetch the count
	result := []TokenLocation{}
	err = db.Select(&result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return result, nil
}