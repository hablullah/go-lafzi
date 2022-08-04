package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// SearchTokens look for document ids which contains the specified tokens,
// then count how many tokens occured in each document.
func SearchTokens(db *sqlx.DB, keys ...string) ([]TokenLocation, error) {
	// If there are no keys submitted, stop early
	if len(keys) == 0 {
		return nil, nil
	}

	// Prepare query
	query, args, err := sqlx.In(`
		SELECT document_id, COUNT(*) n
		FROM token WHERE key IN (?)
		GROUP BY document_id
		ORDER BY n DESC, document_id ASC`, keys)
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
