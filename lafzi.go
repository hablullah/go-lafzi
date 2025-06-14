package lafzi

import (
	"cmp"
	"slices"

	"github.com/hablullah/go-lafzi/internal/database"
	"github.com/hablullah/go-lafzi/internal/myers"
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

	// Get unique tokens
	uniqueTokens := slices.Clone(tokens)
	slices.Sort(uniqueTokens)
	uniqueTokens = slices.Compact(uniqueTokens)
	nUniqueToken := float64(len(uniqueTokens))

	// Search tokens in database
	tokenLocations, err := database.SearchTokens(st.db, uniqueTokens...)
	if err != nil {
		return nil, err
	}

	// Remove the tokens that doesn't meet the minimum threshold
	tokenLocations = slices.DeleteFunc(tokenLocations, func(loc database.TokenLocation) bool {
		score := float64(loc.Count) / nUniqueToken
		return score < st.minConfidence
	})

	// Fetch the documents from database
	docIDs := make([]int, len(tokenLocations))
	for i, loc := range tokenLocations {
		docIDs[i] = loc.DocumentID
	}

	docs, err := database.FetchDocuments(st.db, docIDs...)
	if err != nil {
		return nil, err
	}

	// Create final result by scoring each document using LCS
	nDocs := len(docs)
	results := make([]Result, 0, nDocs)
	for _, doc := range docs {
		// Create string of doc tokens
		docTokens := make([]string, len(doc.Tokens))
		for i, t := range doc.Tokens {
			docTokens[i] = t.Text
		}

		// Calculate final score of the doc
		score, lcsIndexes := myers.Score(docTokens, tokens)
		if score < st.minConfidence {
			continue
		}

		// Get boundary of the match
		start, end := getMatchBoundary(doc.Tokens, lcsIndexes)

		// Save the result
		results = append(results, Result{
			DocumentID: doc.ID,
			Text:       doc.Arabic,
			Confidence: score,
			Start:      start,
			End:        end,
		})
	}

	// Sort results by its confidence
	slices.SortStableFunc(results, func(a, b Result) int {
		return cmp.Compare(b.Confidence, a.Confidence)
	})

	return results, nil
}

func getMatchBoundary(trigrams []phonetic.NGram, lcsIndexes []int) (int, int) {
	if len(trigrams) == 0 || len(lcsIndexes) == 0 {
		return -1, -1
	}

	start := trigrams[lcsIndexes[0]].Start
	end := trigrams[lcsIndexes[0]].End

	for i := 1; i < len(lcsIndexes); i++ {
		idx := lcsIndexes[i]

		if iStart := trigrams[idx].Start; iStart < start {
			start = iStart
		}

		if iEnd := trigrams[idx].End; iEnd > end {
			end = iEnd
		}
	}

	return start, end
}
