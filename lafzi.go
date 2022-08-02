package lafzi

import (
	"github.com/hablullah/go-lafzi/internal/arabic"
	"github.com/hablullah/go-lafzi/internal/database"
	"github.com/hablullah/go-lafzi/internal/phonetic"
	"github.com/hablullah/go-lafzi/internal/tokenizer"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Document is the Arabic document that will be indexed.
type Document struct {
	ID     int
	Arabic string
}

// Storage is the container for storing reverse indexes for
// Arabic documents that will be searched later. Use sqlite3
// as database engine.
type Storage struct {
	db *sqlx.DB
}

// OpenStorage open the reverse indexes database in the specified path.
func OpenStorage(path string) (*Storage, error) {
	db, err := database.Open(path)
	if err != nil {
		return nil, err
	}

	st := &Storage{db}
	return st, nil
}

// AddDocuments save and index the documents into the storage.
func (st *Storage) AddDocuments(docs ...Document) error {
	// Convert Arabic text to phonetics
	dbDocs := make([]database.Document, len(docs))
	for i, doc := range docs {
		dbDocs[i] = database.Document{
			ID:      doc.ID,
			Content: arabic.ToPhonetic(doc.Arabic),
		}
	}

	// Save documents to database
	return database.InsertDocuments(st.db, dbDocs...)
}

// DeleteDocuments remove the documents in the storage.
func (st *Storage) DeleteDocuments(ids ...int) error {
	return database.DeleteDocuments(st.db, ids...)
}

// Search for suitable documents using the specified query.
// On success will return the ids of suitable documents.
func (st *Storage) Search(query string) ([]int, error) {
	// Normalize query
	query = phonetic.Normalize(query)

	// Convert query to trigram tokens
	tokens := tokenizer.Split(query)
	nUniqueToken := float64(countUniqueTokens(tokens...))

	// Search tokens in database
	tokenLocations, err := database.SearchTokens(st.db, tokens...)
	if err != nil {
		return nil, err
	}

	// Score the documents using count of the tokens
	// Here we use 40% as the minimum threshold
	scores := make(map[int]float64)
	for _, loc := range tokenLocations {
		score := float64(loc.Count) / nUniqueToken
		if score >= 0.4 {
			scores[loc.DocID] = score
		}
	}

	// Score the document using LCS and its compactness
	// TODO:

	return nil, nil
}

func countUniqueTokens(tokens ...string) int {
	var nUniqueToken int
	tracker := make(map[string]struct{})

	for _, t := range tokens {
		if _, exist := tracker[t]; !exist {
			tracker[t] = struct{}{}
			nUniqueToken++
		}
	}

	return nUniqueToken
}
