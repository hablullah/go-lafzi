package lafzi

import (
	"fmt"
	"math"
	"sort"
)

type Database struct {
	storage dataStorage
}

type DatabaseEntry struct {
	ID         int64
	ArabicText string
}

type dbEntryTokens struct {
	ID           int64
	TokenCount   int
	TokenIndexes []int
}

type dbEntryScore struct {
	ID                  int64
	TokenCount          int
	NLongestSubSequence int
	SubSequenceDensity  float64
}

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

func (db *Database) Close() {
	db.storage.close()
}

func (db *Database) AddEntries(entries ...DatabaseEntry) error {
	return db.storage.saveEntries(entries...)
}

func (db *Database) Search(transliteration string) error {
	// Convert latin text into tokens
	query := queryFromLatin(transliteration)
	tokens := tokenizeQuery(query)
	if len(tokens) == 0 {
		return nil
	}

	// Find entries that contains the tokens
	entries, err := db.storage.findTokens(tokens...)
	if err != nil {
		return err
	}

	// Calculate score and filter the dictionary entries.
	// Here we want at least 3/4 of tokens found in each entry.
	countThreshold := int(math.Ceil(float64(len(tokens)) * 3 / 4))
	if countThreshold <= 1 {
		countThreshold = len(tokens)
	}

	entryScores := []dbEntryScore{}
	for _, entry := range entries {
		// Make sure count of token inside this entry pass the threshold
		if entry.TokenCount < countThreshold {
			continue
		}

		// Make sure length of longest sub sequence pass the threshold as well
		longestSubSequence := db.getLongestSubSequence(entry.TokenIndexes)
		nLongestSubSequence := len(longestSubSequence)
		if nLongestSubSequence < countThreshold {
			continue
		}

		// Calculate sequence density
		density := db.getSequenceDensity(longestSubSequence)
		if density < 0.5 {
			continue
		}

		entryScores = append(entryScores, dbEntryScore{
			ID:                  entry.ID,
			TokenCount:          entry.TokenCount,
			NLongestSubSequence: nLongestSubSequence,
			SubSequenceDensity:  density,
		})
	}

	// Sort entry score with following order:
	// - token count, descending
	// - sub sequence density, descending
	// - entry id, ascending
	sort.Slice(entryScores, func(a, b int) bool {
		scoreA := entryScores[a]
		scoreB := entryScores[b]

		if scoreA.TokenCount != scoreB.TokenCount {
			return scoreA.TokenCount > scoreB.TokenCount
		}

		if scoreA.SubSequenceDensity != scoreB.SubSequenceDensity {
			return scoreA.SubSequenceDensity > scoreB.SubSequenceDensity
		}

		return scoreA.ID < scoreB.ID
	})

	for _, score := range entryScores {
		fmt.Println(score.ID, len(tokens),
			score.TokenCount,
			score.NLongestSubSequence,
			score.SubSequenceDensity)
	}

	return nil
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

	nSequence := len(sequence)
	return (1 / float64(nSequence-1)) * sigma
}