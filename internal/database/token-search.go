package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type TokenLocation struct {
	DocumentID int `db:"document_id"`
	Count      int `db:"n"`
}

// SearchTokens look for document ids which contains the specified tokens,
// then count how many tokens occured in each document.
func SearchTokens(db *sqlx.DB, tokens ...string) ([]TokenLocation, error) {
	// If there are no tokens submitted, stop early
	if len(tokens) == 0 {
		return nil, nil
	}

	// Prepare query
	query, args, err := sqlx.In(`
		SELECT document_id, COUNT(*) n
		FROM document_token WHERE token IN (?)
		GROUP BY document_id
		ORDER BY n DESC, document_id ASC`, tokens)
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
