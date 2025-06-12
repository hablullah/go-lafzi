package lafzi

import (
	"sort"

	"github.com/hablullah/go-lafzi/internal/arabic"
	"github.com/hablullah/go-lafzi/internal/database"
	"github.com/hablullah/go-lafzi/internal/myers"
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

// Result contains id of the suitable document and its confidence level.
type Result struct {
	DocumentID int
	Confidence float64
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
	docIDs := make([]int, 0)
	scores := make(map[int]float64)
	for _, loc := range tokenLocations {
		score := float64(loc.Count) / nUniqueToken
		if score >= st.minConfidence {
			docIDs = append(docIDs, loc.DocID)
			scores[loc.DocID] = score
		}
	}

	// Fetch the document from database
	docs, err := database.FetchDocuments(st.db, docIDs...)
	if err != nil {
		return nil, err
	}

	// Create final result by scoring each document using LCS
	results := make([]Result, 0)
	for _, doc := range docs {
		docTokens := tokenizer.Split(doc.Content)
		score := myers.Score(docTokens, tokens)

		if score >= st.minConfidence {
			results = append(results, Result{
				DocumentID: doc.ID,
				Confidence: score,
			})
		}
	}

	// Sort results by its confidence
	sort.SliceStable(results, func(a, b int) bool {
		return results[a].Confidence > results[b].Confidence
	})

	return results, nil
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
