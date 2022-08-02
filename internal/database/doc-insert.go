package database

import (
	"fmt"

	"github.com/hablullah/go-lafzi/internal/tokenizer"
	"github.com/jmoiron/sqlx"
)

// InsertDocuments save the documents into the database.
func InsertDocuments(db *sqlx.DB, docs ...Document) (err error) {
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
	stmtInsertDoc, err := tx.Preparex(`
		INSERT INTO document (id, content) VALUES (?, ?)
		ON CONFLICT (id) DO UPDATE
		SET content = excluded.content`)
	if err != nil {
		return
	}

	stmtDeleteToken, err := tx.Preparex(`
		DELETE FROM token
		WHERE document_id = ?`)
	if err != nil {
		return
	}

	stmtInsertToken, err := tx.Preparex(`
		INSERT INTO token (key, document_id) VALUES (?, ?)
		ON CONFLICT DO NOTHING`)
	if err != nil {
		return
	}

	// Process each document
	for _, doc := range docs {
		// Remove any token that associated with this document
		_, err = stmtDeleteToken.Exec(doc.ID)
		if err != nil {
			return
		}

		// Save document
		_, err = stmtInsertDoc.Exec(doc.ID, doc.Content)
		if err != nil {
			return
		}

		// Save tokens
		for _, token := range tokenizer.Split(doc.Content) {
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
