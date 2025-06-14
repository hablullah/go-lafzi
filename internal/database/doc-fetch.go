package database

import (
	"fmt"

	"github.com/hablullah/go-lafzi/internal/phonetic"
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
		err = tx.Rollback()
	}()

	// Prepare statement
	stmtGetDoc, err := tx.Preparex(`
		SELECT id, arabic FROM document
		WHERE id = ?`)
	if err != nil {
		return
	}

	stmtGetTokens, err := tx.Preparex(`
		SELECT token, start, end FROM document_token
		WHERE document_id = ?
		ORDER BY start, end`)
	if err != nil {
		return
	}

	// Fetch data
	documents = make([]Document, 0, len(ids))
	for _, id := range ids {
		var doc Document
		err = stmtGetDoc.Get(&doc, id)
		if err != nil {
			return
		}

		var tokens []DocumentToken
		err = stmtGetTokens.Select(&tokens, id)
		if err != nil {
			return
		}

		doc.Tokens = make([]phonetic.NGram, len(tokens))
		for i, token := range tokens {
			doc.Tokens[i] = phonetic.NGram{
				Text:  token.Token,
				Start: token.Start,
				End:   token.End,
			}
		}

		documents = append(documents, doc)
	}

	return
}
