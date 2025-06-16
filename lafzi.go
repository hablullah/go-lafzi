package lafzi

import (
	"github.com/hablullah/go-lafzi/internal/database"
	"github.com/hablullah/go-lafzi/internal/phonetic"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// Document is the Arabic document that will be indexed.
type Document struct {
	ID     int
	Arabic string
}

// Result contains id of the suitable document and its confidence level.
type Result struct {
	DocumentID int
	Text       string
	Confidence float64
	Start      int
	End        int
}

// Storage is the container for storing reverse indexes for
// Arabic documents that will be searched later. Use sqlite3
// as database engine.
type Storage struct {
	db            *sqlx.DB
	minConfidence float64
}

// OpenStorage open the reverse indexes database in the specified path.
func OpenStorage(path string) (*Storage, error) {
	db, err := database.Open(path)
	if err != nil {
		return nil, err
	}

	return &Storage{db, 0.4}, nil
}

// AddDocuments save and index the documents into the storage.
func (st *Storage) AddDocuments(docs ...Document) error {
	// Convert Arabic text to phonetics
	dbDocs := make([]database.InsertDocumentArg, len(docs))
	for i, doc := range docs {
		dbDocs[i] = database.InsertDocumentArg{
			DocumentID: doc.ID,
			Arabic:     doc.Arabic,
			Phonetic:   phonetic.FromArabic(doc.Arabic),
		}
	}

	// Save documents to database
	return database.InsertDocuments(st.db, dbDocs...)
}

// DeleteDocuments remove the documents in the storage.
func (st *Storage) DeleteDocuments(ids ...int) error {
	return database.DeleteDocuments(st.db, ids...)
}

// SetMinConfidence set the minimum confidence score for
// the search result. Default is 40%.
func (st *Storage) SetMinConfidence(f float64) {
	switch {
	case f > 1:
		st.minConfidence = 1
	case f <= 0:
		st.minConfidence = 0.4 // default is 40%
	default:
		st.minConfidence = f
	}
}

// Search for suitable documents using the specified query.
func (st *Storage) Search(query string) ([]Result, error) {
	// Normalize query
	query = phonetic.NormalizeString(query)

	// Convert query to trigram tokens
	tokens := phonetic.NGrams(query, 3)

	// Search tokens in database
	searchResults, err := database.SearchTokens(st.db, st.minConfidence, tokens...)
	if err != nil {
		return nil, err
	}

	// Fetch the documents from database
	docIDs := make([]int, len(searchResults))
	for i, loc := range searchResults {
		docIDs[i] = loc.DocumentID
	}

	docs, err := database.FetchDocuments(st.db, docIDs...)
	if err != nil {
		return nil, err
	}

	// Create final result
	results := make([]Result, len(docs))
	for i, doc := range docs {
		// Save the result
		results[i] = Result{
			DocumentID: doc.ID,
			Text:       doc.Arabic,
			Confidence: searchResults[i].Confidence,
			Start:      searchResults[i].Start,
			End:        searchResults[i].End,
		}
	}

	return results, nil
}
