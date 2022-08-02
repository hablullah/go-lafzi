package lafzi

import (
	"fmt"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Storage is the container for storing reverse indexes for
// Arabic documents that will be searched later. Use sqlite3
// database as back-end.
type Storage struct {
	db *sqlx.DB
}

// Document is the Arabic document that will be indexed.
type Document struct {
	ID     int
	Arabic string
}

// OpenStorage lafzi storage on specified path.
func OpenStorage(path string) (storage *Storage, err error) {
	// Prepare DSN
	q := url.Values{}
	q.Add("_foreign_keys", "1")
	dsn := path

	// Connect database
	db, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	db.SetConnMaxLifetime(time.Minute)

	// Create transaction
	var tx *sqlx.Tx
	tx, err = db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}

	// If error ever happened, rollback and close database
	defer func() {
		if err != nil {
			if tx != nil {
				tx.Rollback()
			}
			if db != nil {
				db.Close()
			}
			storage = nil
		}
	}()

	// Generate tables
	ddlQueries := []string{
		ddlCreateDocument,
		ddlCreateToken,
		ddlCreateTokenIndex}

	for _, query := range ddlQueries {
		_, err = tx.Exec(query)
		if err != nil {
			return
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return
	}

	// Return the storage
	storage = &Storage{db: db}
	return
}

// AddDocuments save the Arabic document into storage.
func (s *Storage) AddDocuments(docs ...Document) (err error) {
	// Start transaction
	tx, err := s.db.Beginx()
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
	stmtInsertDocument, err := tx.Preparex(`
		INSERT INTO document (id, content)
		VALUES (?, ?)
		ON CONFLICT (id) DO UPDATE
		SET content = excluded.content`)
	if err != nil {
		return
	}

	stmtInsertToken, err := tx.Preparex(`
		INSERT INTO token (key, document_id)
		VALUES (?, ?)
		ON CONFLICT DO NOTHING`)
	if err != nil {
		return
	}

	// Process each document
	for _, doc := range docs {
		// Create content and tokens of document
		content := ArabicToPhonetic(doc.Arabic)
		content = NormalizePhonetic(content)
		tokens := tokenize(content)

		// Save content
		_, err = stmtInsertDocument.Exec(doc.ID, content)
		if err != nil {
			return
		}

		// Save tokens
		for _, token := range tokens {
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

// SearchTokens look for documents which contains the specified tokens, then count
// how many tokens occured in each document.
func (s *Storage) SearchTokens(keys ...string) (map[int]int, error) {
	// Prepare query
	sql := `SELECT document_id, COUNT(*) n
		FROM token WHERE key IN (?)
		GROUP BY document_id`
	sql, args, err := sqlx.In(sql, keys)
	if err != nil {
		return nil, fmt.Errorf("failed to expand query: %v", err)
	}

	// Fetch the documents
	result := make(map[int]int)
	if err = s.db.Select(&result, sql, args...); err != nil {
		return nil, fmt.Errorf("failed to search token: %v", err)
	}

	return result, nil
}

const ddlCreateDocument = `
CREATE TABLE IF NOT EXISTS document (
	id          INT  NOT NULL,
	content     TEXT NOT NULL,
	PRIMARY KEY (id))`

const ddlCreateToken = `
CREATE TABLE IF NOT EXISTS token (
	key        TEXT NOT NULL,
	document_id INT  NOT NULL,
	PRIMARY KEY (key, document_id),
	CONSTRAINT token_document_fk FOREIGN KEY (document_id) REFERENCES document (id))`

const ddlCreateTokenIndex = `
CREATE INDEX IF NOT EXISTS token_key_idx ON token (key)`
