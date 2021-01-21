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

	tx.MustExec(`CREATE TABLE IF NOT EXISTS entry_token (
		entry_id  INTEGER NOT NULL,
		token_id  INTEGER NOT NULL,
		position  INTEGER NOT NULL,
		CONSTRAINT entry_token_pk PRIMARY KEY (entry_id, token_id, position),
		CONSTRAINT entry_token_id_fk FOREIGN KEY (token_id) REFERENCES token (id))`)

	// Generate indexes
	tx.MustExec(`CREATE INDEX IF NOT EXISTS entry_token_token_id_idx ON entry_token(token_id)`)

	// Commit transaction
	err = tx.Commit()
	panicError(err)

	storage = &sqliteStorage{db}
	return
}

func (ss *sqliteStorage) close() {
	ss.Close()
}

func (ss *sqliteStorage) saveEntries(entries ...DatabaseEntry) (err error) {
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

	stmtInsertEntryToken, err := tx.Preparex(`INSERT OR IGNORE INTO entry_token
		(entry_id, token_id, position) VALUES (?, ?, ?)`)
	panicError(err)

	// Proses each entry
	for _, entry := range entries {
		// Create token from entry
		query := queryFromArabic(entry.ArabicText)
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

			// Save entry token
			for _, idx := range indexes {
				stmtInsertEntryToken.MustExec(entry.ID, tokenID, idx)
			}
		}
	}

	// Commit transaction
	err = tx.Commit()
	panicError(err)

	return
}

func (ss *sqliteStorage) findTokens(tokens ...string) (result []dbEntryTokens, err error) {
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
	type tokenEntryData struct {
		EntryID  int64 `db:"entry_id"`
		Position int   `db:"position"`
	}

	stmtGetTokenID, err := tx.Preparex(
		`SELECT id FROM token WHERE value = ?`)
	panicError(err)

	stmtSelectEntry, err := tx.Preparex(`
		SELECT et.entry_id, et.position
		FROM entry_token et
		LEFT JOIN token t ON et.token_id = t.id
		WHERE t.value = ?`)
	panicError(err)

	// For each token, find the dictionary entry that contains such token,
	// also the position of that token within the dictionary entry.
	entryTokenCount := map[int64]int{}
	entryTokenIndexes := map[int64]map[int]struct{}{}

	for _, token := range tokens {
		// Get token ID
		var tokenID int64
		err = stmtGetTokenID.Get(&tokenID, token)
		panicError(err)

		if tokenID == 0 {
			continue
		}

		// Get token entries
		tokenEntries := []tokenEntryData{}
		err = stmtSelectEntry.Select(&tokenEntries, token)
		panicError(err)

		for _, te := range tokenEntries {
			existingIndexes := entryTokenIndexes[te.EntryID]
			if existingIndexes == nil {
				existingIndexes = map[int]struct{}{}
			}
			existingIndexes[te.Position] = struct{}{}

			entryTokenCount[te.EntryID]++
			entryTokenIndexes[te.EntryID] = existingIndexes
		}
	}

	// Convert map of token count and indexes to array
	result = []dbEntryTokens{}
	for entryID, indexes := range entryTokenIndexes {
		arrIndexes := []int{}
		for idx := range indexes {
			arrIndexes = append(arrIndexes, idx)
		}
		sort.Ints(arrIndexes)

		result = append(result, dbEntryTokens{
			ID:           entryID,
			TokenCount:   entryTokenCount[entryID],
			TokenIndexes: arrIndexes,
		})
	}

	return
}
