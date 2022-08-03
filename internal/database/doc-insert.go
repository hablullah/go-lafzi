package database

import (
	"fmt"

	"github.com/hablullah/go-lafzi/internal/tokenizer"
	"github.com/jmoiron/sqlx"
)

// InsertDocuments save the documents into the database.
func InsertDocuments(db *sqlx.DB, docs ...Document) (err error) {
	// Remove index, and create it once it over
	db.Exec(`DROP INDEX IF EXISTS token_key_idx`)
	defer db.Exec(ddlCreateTokenIndex)

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

	// Prepare statement
	stmtDeleteToken, err := tx.Preparex(`
		DELETE FROM token
		WHERE document_id = ?`)
	if err != nil {
		return
	}

	stmtInsertDoc, err := tx.Preparex(`
		INSERT INTO document (id, content) VALUES (?, ?)
		ON CONFLICT (id) DO UPDATE
		SET content = excluded.content`)
	if err != nil {
		return
	}

	stmtInsertToken, err := tx.Preparex(`
		INSERT INTO token (key, document_id) VALUES (?, ?)
		ON CONFLICT DO NOTHING`)
	if err != nil {
		return
	}

	// Remove any token that associated with this document
	for _, doc := range docs {
		_, err = stmtDeleteToken.Exec(doc.ID)
		if err != nil {
			return
		}
	}

	// Insert the document
	for _, doc := range docs {
		// Save document
		_, err = stmtInsertDoc.Exec(doc.ID, doc.Content)
		if err != nil {
			return
		}

		// Save tokens
		tokens := tokenizer.Split(doc.Content)
		for _, token := range getUniqueTokens(tokens) {
			_, err = stmtInsertToken.Exec(token, doc.ID)
			if err != nil {
				return
			}
		}
	}

	// Commit to database
	err = tx.Commit()
	return
}

func getUniqueTokens(tokens []string) []string {
	var uniqueTokens []string
	tracker := make(map[string]struct{})

	for _, t := range tokens {
		if _, exist := tracker[t]; !exist {
			tracker[t] = struct{}{}
			uniqueTokens = append(uniqueTokens, t)
		}
	}

	return uniqueTokens
}
