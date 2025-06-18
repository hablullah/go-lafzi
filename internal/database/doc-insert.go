package database

import (
	"database/sql"
	"fmt"

	"github.com/hablullah/go-lafzi/internal/phonetic"
	"github.com/jmoiron/sqlx"
)

type InsertDocumentArg struct {
	Identifier string
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
	_, err = db.Exec(`DROP INDEX IF EXISTS document_token_idx_token`)
	if err != nil {
		return
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

		// Recreate index
		if err == nil {
			_, err = db.Exec(ddlCreateDocumentTokenIndexToken)
		}
	}()

	// Prepare statement
	stmtGetDoc, err := tx.Preparex(`SELECT id FROM document WHERE identifier = ?`)
	if err != nil {
		return
	}

	stmtInsertDoc, err := tx.Preparex(`
		INSERT INTO document (identifier, arabic)
		VALUES (?, ?)
		ON CONFLICT (identifier) DO UPDATE
		SET arabic = excluded.arabic`)
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
		INSERT INTO document_token (document_id, token, start, end)
		VALUES (?, ?, ?, ?)
		ON CONFLICT DO NOTHING`)
	if err != nil {
		return
	}

	// Insert the document
	for _, arg := range args {
		// Get document ID if it's exist
		var documentID int64
		documentExist := true
		err = stmtGetDoc.Get(&documentID, arg.Identifier)
		if err != nil {
			if err == sql.ErrNoRows {
				documentExist = false
				err = nil
			} else {
				return
			}
		}

		// Save document
		var res sql.Result
		res, err = stmtInsertDoc.Exec(
			arg.Identifier,
			arg.Arabic)
		if err != nil {
			return
		}

		// If document not exist, use ID from last inserted
		if !documentExist {
			documentID, err = res.LastInsertId()
			if err != nil {
				return
			}
		}

		// Remove any token that associated with this document
		_, err = stmtDeleteDocToken.Exec(documentID)
		if err != nil {
			return
		}

		// Save tokens
		for _, token := range arg.Phonetic.Split(3) {
			_, err = stmtInsertDocToken.Exec(
				documentID,
				token.Text,
				token.Start,
				token.End)
			if err != nil {
				return
			}
		}
	}

	// Commit to database
	err = tx.Commit()
	return
}
