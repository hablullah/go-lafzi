package lafzi

import (
	"math"
	"sort"
)

// Database is the core of Lafzi. Used to store the position of tokens within the submitted documents.
type Database struct {
	storage dataStorage
}

// Document is the Arabic document that will be indexed.
type Document struct {
	ID         int64
	ArabicText string
}

type documentTokens struct {
	ID           int64
	TokenCount   int
	TokenIndexes []int
}

type documentScore struct {
	ID                  int64
	TokenCount          int
	NLongestSubSequence int
	SubSequenceDensity  float64
}

// OpenDatabase open and creates database at the specified path.
func OpenDatabase(path string, storageType StorageType) (*Database, error) {
	var err error
	var storage dataStorage

	switch storageType {
	case SQLite:
		storage, err = newSQLiteStorage(path)
	default:
		storage, err = newBoltStorage(path)
	}

	if err != nil {
		return nil, err
	}

	return &Database{storage}, nil
}

// Close closes the database and prevent any read and write.
func (db *Database) Close() {
	db.storage.close()
}

// AddDocuments adds the documents to database.
func (db *Database) AddDocuments(documents ...Document) error {
	return db.storage.saveDocuments(documents...)
}

// Search looks for documents whose transliterations contain the specified keyword.
func (db *Database) Search(keyword string) ([]int64, error) {
	// Convert keyword into tokens
	query := queryFromLatin(keyword)
	tokens := tokenizeQuery(query)
	if len(tokens) == 0 {
		return nil, nil
	}

	// Find documents that contains the tokens
	documents, err := db.storage.findTokens(tokens...)
	if err != nil {
		return nil, err
	}

	// Calculate score and filter the dictionary documents.
	// Here we want at least 3/4 of tokens found in each document.
	countThreshold := int(math.Ceil(float64(len(tokens)) * 3 / 4))
	if countThreshold <= 1 {
		countThreshold = len(tokens)
	}

	documentScores := []documentScore{}
	for _, doc := range documents {
		// Make sure count of token inside this document pass the threshold
		if doc.TokenCount < countThreshold {
			continue
		}

		// Make sure length of longest sub sequence pass the threshold as well
		longestSubSequence := db.getLongestSubSequence(doc.TokenIndexes)
		nLongestSubSequence := len(longestSubSequence)
		if nLongestSubSequence < countThreshold {
			continue
		}

		// Calculate sequence density
		density := db.getSequenceDensity(longestSubSequence)
		if density < 0.5 {
			continue
		}

		documentScores = append(documentScores, documentScore{
			ID:                  doc.ID,
			TokenCount:          doc.TokenCount,
			NLongestSubSequence: nLongestSubSequence,
			SubSequenceDensity:  density,
		})
	}

	// Sort document scores with following order:
	// - token count, descending
	// - sub sequence density, descending
	// - document id, ascending
	sort.Slice(documentScores, func(a, b int) bool {
		scoreA := documentScores[a]
		scoreB := documentScores[b]

		if scoreA.TokenCount != scoreB.TokenCount {
			return scoreA.TokenCount > scoreB.TokenCount
		}

		if scoreA.SubSequenceDensity != scoreB.SubSequenceDensity {
			return scoreA.SubSequenceDensity > scoreB.SubSequenceDensity
		}

		return scoreA.ID < scoreB.ID
	})

	result := make([]int64, len(documentScores))
	for i, score := range documentScores {
		result[i] = score.ID
	}

	return result, nil
}

func (db *Database) getLongestSubSequence(sequence []int) []int {
	var maxStart, maxLength int
	var currentStart, currentLength int

	for i := 1; i < len(sequence); i++ {
		// If current number difference with the previous is less than five,
		// it's still within one sequence.
		if sequence[i]-sequence[i-1] <= 5 {
			currentLength++
			continue
		}

		// If not, then it's a brand new sequence.
		// Check if it's larger than current biggest sub sequence
		if currentLength > maxLength {
			maxStart = currentStart
			maxLength = currentLength
		}

		currentStart = i
		currentLength = 0
	}

	// There are cases where a sequence only have exactly one sub sequence
	// (sequence = sub sequence). In this case, maxLength will be 0, so we need
	// to check it here.
	if currentLength > maxLength {
		maxStart = currentStart
		maxLength = currentLength
	}

	return sequence[maxStart : maxStart+maxLength+1]
}

func (db *Database) getSequenceDensity(sequence []int) float64 {
	var sigma float64
	for i := 0; i < len(sequence)-1; i++ {
		tmp := sequence[i+1] - sequence[i]
		sigma += 1 / float64(tmp)
	}

	switch nSequence := len(sequence); nSequence {
	case 1:
		return 1
	case 0:
		return 0
	default:
		return (1 / float64(nSequence-1)) * sigma
	}
}
