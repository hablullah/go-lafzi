package database

import (
	"fmt"

	"github.com/hablullah/go-lafzi/internal/phonetic"
	"github.com/jmoiron/sqlx"
)

type InsertDocumentArg struct {
	DocumentID int
	Arabic     string
	Phonetic   phonetic.Group
}

// InsertDocuments save the documents into the database.
func InsertDocuments(db *sqlx.DB, args ...InsertDocumentArg) (err error) {
	// If there are no args submitted, stop early
	if len(args) == 0 {
		return nil
	}

	// Remove index, and create it once it over
	_, err = db.Exec(`DROP INDEX IF EXISTS document_token_idx`)
	if err != nil {
		return err
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
			err = tx.Rollback()
		}

		// Recreate index
		if err == nil {
			_, err = db.Exec(ddlCreateDocumentTokenIndex)
		}
	}()

	// Prepare statement
	stmtInsertDoc, err := tx.Preparex(`
		INSERT INTO document (id, arabic, phonetic)
		VALUES (?, ?, ?)
		ON CONFLICT (id) DO UPDATE
		SET arabic   = excluded.arabic,
			phonetic = excluded.phonetic`)
	if err != nil {
		return
	}

	stmtDeleteDocToken, err := tx.Preparex(`
		DELETE FROM document_token
		WHERE document_id = ?`)
	if err != nil {
		return
	}

	stmtInsertDocToken, err := tx.Preparex(`
		INSERT INTO document_token (document_id, token)
		VALUES (?, ?)
		ON CONFLICT DO NOTHING`)
	if err != nil {
		return
	}

	// Insert the document
	for _, arg := range args {
		// Save document
		_, err = stmtInsertDoc.Exec(arg.DocumentID, arg.Arabic, arg.Phonetic.String())
		if err != nil {
			return
		}

		// Remove any token that associated with this document
		_, err = stmtDeleteDocToken.Exec(arg.DocumentID)
		if err != nil {
			return
		}

		// Save tokens
		for _, token := range arg.Phonetic.Split(3) {
			_, err = stmtInsertDocToken.Exec(arg.DocumentID, token.String())
			if err != nil {
				return
			}
		}
	}

	// Commit to database
	err = tx.Commit()
	return
}
