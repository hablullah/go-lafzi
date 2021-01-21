package lafzi

import (
	"fmt"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteStorage struct {
	*sqlx.DB
}

func newSQLiteStorage(path string) (storage *sqliteStorage, err error) {
	// Connect to database
	dsn := fmt.Sprintf("file:%s?_fk=1", path)
	db, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	db.SetConnMaxLifetime(time.Minute)

	// Create transaction
	tx, err := db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Make sure to rollback if panic ever happened
	defer func() {
		if r := recover(); r != nil {
			panicErr, _ := r.(error)
			tx.Rollback()

			storage = nil
			err = panicErr
		}
	}()

	// Generate tables
	tx.MustExec(`CREATE TABLE IF NOT EXISTS token (
		id         INTEGER NOT NULL,
		value      TEXT    NOT NULL,
		CONSTRAINT token_pk PRIMARY KEY (id),
		CONSTRAINT token_value_UNIQUE UNIQUE (value))`)

	tx.MustExec(`CREATE TABLE IF NOT EXISTS document_token (
		document_id INTEGER NOT NULL,
		token_id    INTEGER NOT NULL,
		position    INTEGER NOT NULL,
		CONSTRAINT document_token_pk PRIMARY KEY (document_id, token_id, position),
		CONSTRAINT document_token_id_fk FOREIGN KEY (token_id) REFERENCES token (id))`)

	// Generate indexes
	tx.MustExec(`CREATE INDEX IF NOT EXISTS document_token_token_id_idx ON document_token(token_id)`)

	// Commit transaction
	err = tx.Commit()
	panicError(err)

	storage = &sqliteStorage{db}
	return
}

func (ss *sqliteStorage) close() {
	ss.Close()
}

func (ss *sqliteStorage) saveDocuments(documents ...Document) (err error) {
	// Create transaction
	tx, err := ss.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	// Make sure to rollback if panic ever happened
	defer func() {
		if r := recover(); r != nil {
			panicErr, _ := r.(error)
			tx.Rollback()
			err = panicErr
		}
	}()

	// Prepare statements
	stmtGetToken, err := tx.Preparex(`SELECT id FROM token WHERE value = ?`)
	panicError(err)

	stmtInsertToken, err := tx.Preparex(`INSERT INTO token (value) VALUES (?)`)
	panicError(err)

	stmtInsertDocumentToken, err := tx.Preparex(`INSERT OR IGNORE INTO document_token
		(document_id, token_id, position) VALUES (?, ?, ?)`)
	panicError(err)

	// Proses each document
	for _, doc := range documents {
		// Create token from document
		query := queryFromArabic(doc.ArabicText)
		tokens := tokenizeQuery(query)
		if len(tokens) == 0 {
			continue
		}

		// Compact the tokens
		compactTokens := map[string][]int{}
		for idx, token := range tokens {
			if existingIndexes, exist := compactTokens[token]; !exist {
				compactTokens[token] = []int{idx}
			} else {
				existingIndexes = append(existingIndexes, idx)
				compactTokens[token] = existingIndexes
			}
		}

		// Save all tokens
		for token, indexes := range compactTokens {
			// Check if token already saved before
			var tokenID int64
			err = stmtGetToken.Get(&tokenID, token)
			panicError(err)

			// If not, save it to database
			if tokenID == 0 {
				res := stmtInsertToken.MustExec(token)
				tokenID, _ = res.LastInsertId()
			}

			// Save document token
			for _, idx := range indexes {
				stmtInsertDocumentToken.MustExec(doc.ID, tokenID, idx)
			}
		}
	}

	// Commit transaction
	err = tx.Commit()
	panicError(err)

	return
}

func (ss *sqliteStorage) findTokens(tokens ...string) (result []documentTokens, err error) {
	// Create read only transaction
	tx, err := ss.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			panicErr, _ := r.(error)
			tx.Rollback()

			result = nil
			err = panicErr
		}
	}()

	// Preare type and statement
	type tokenDocumentData struct {
		DocumentID int64 `db:"document_id"`
		Position   int   `db:"position"`
	}

	stmtGetTokenID, err := tx.Preparex(
		`SELECT id FROM token WHERE value = ?`)
	panicError(err)

	stmtSelectDocument, err := tx.Preparex(`
		SELECT dt.document_id, dt.position
		FROM document_token dt
		LEFT JOIN token t ON dt.token_id = t.id
		WHERE t.value = ?`)
	panicError(err)

	// For each token, find the documents that contains such token,
	// also the position of that token within the document.
	docTokenCount := map[int64]int{}
	docTokenIndexes := map[int64]map[int]struct{}{}

	for _, token := range tokens {
		// Get token ID
		var tokenID int64
		err = stmtGetTokenID.Get(&tokenID, token)
		panicError(err)

		if tokenID == 0 {
			continue
		}

		// Get token documents
		tokenDocuments := []tokenDocumentData{}
		err = stmtSelectDocument.Select(&tokenDocuments, token)
		panicError(err)

		for _, td := range tokenDocuments {
			existingIndexes := docTokenIndexes[td.DocumentID]
			if existingIndexes == nil {
				existingIndexes = map[int]struct{}{}
			}
			existingIndexes[td.Position] = struct{}{}

			docTokenCount[td.DocumentID]++
			docTokenIndexes[td.DocumentID] = existingIndexes
		}
	}

	// Convert map of token count and indexes to array
	result = []documentTokens{}
	for docID, indexes := range docTokenIndexes {
		arrIndexes := []int{}
		for idx := range indexes {
			arrIndexes = append(arrIndexes, idx)
		}
		sort.Ints(arrIndexes)

		result = append(result, documentTokens{
			ID:           docID,
			TokenCount:   docTokenCount[docID],
			TokenIndexes: arrIndexes,
		})
	}

	err = nil
	return
}
