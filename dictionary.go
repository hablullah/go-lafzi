package lafzi

import (
	"fmt"
	"math"
	"os"
	"sort"

	"go.etcd.io/bbolt"
)

type Dictionary struct {
	*bbolt.DB
}

type DictionaryEntry struct {
	ID         int64
	ArabicText string
}

type dictionaryEntryScore struct {
	ID                  int64
	TokenCount          int
	NLongestSubSequence int
	SubSequenceDensity  float64
}

func OpenDictionary(path string) (*Dictionary, error) {
	db, err := bbolt.Open(path, os.ModePerm, bbolt.DefaultOptions)
	if err != nil {
		return nil, err
	}

	return &Dictionary{db}, nil
}

func (dict *Dictionary) AddEntries(entries ...DictionaryEntry) error {
	return dict.Update(func(tx *bbolt.Tx) (err error) {
		for _, entry := range entries {
			err = dict.saveEntry(tx, entry)
			if err != nil {
				return
			}
		}
		return
	})
}

func (dict *Dictionary) Lookup(latinText string) error {
	// Convert latin text into tokens
	query := queryFromLatin(latinText)
	tokens := tokenizeQuery(query)
	if len(tokens) == 0 {
		return nil
	}

	// For each token, find the dictionary entry that contains such token,
	// also the position of that token within the dictionary entry.
	entryTokenCount := map[int64]int{}
	entryTokenIndexes := map[int64]map[int]struct{}{}
	dict.View(func(tx *bbolt.Tx) error {
		for _, token := range tokens {
			tokenBucket := tx.Bucket([]byte(token))
			if tokenBucket == nil {
				continue
			}

			tokenBucket.ForEach(func(btEntryID, btIndexes []byte) error {
				entryID := bytesToInt64(btEntryID)
				tokenIndexes := bytesToArrayInt(btIndexes)

				existingIndexes := entryTokenIndexes[entryID]
				if existingIndexes == nil {
					existingIndexes = map[int]struct{}{}
				}

				for _, idx := range tokenIndexes {
					existingIndexes[idx] = struct{}{}
				}

				entryTokenCount[entryID]++
				entryTokenIndexes[entryID] = existingIndexes
				return nil
			})
		}

		return nil
	})

	// Calculate score and filter the dictionary entries.
	// Here we want at least 3/4 of tokens found in each entry.
	countThreshold := int(math.Ceil(float64(len(tokens)) * 3 / 4))
	if countThreshold <= 1 {
		countThreshold = len(tokens)
	}

	entryScores := []dictionaryEntryScore{}
	for entryID, indexes := range entryTokenIndexes {
		// Make sure count of token inside this entry pass the threshold
		tokenCount := entryTokenCount[entryID]
		if tokenCount < countThreshold {
			continue
		}

		// Convert indexes (which is map) to array
		arrIndexes := []int{}
		for idx := range indexes {
			arrIndexes = append(arrIndexes, idx)
		}
		sort.Ints(arrIndexes)

		// Make sure length of longest sub sequence pass the threshold as well
		longestSubSequence := dict.getLongestSubSequence(arrIndexes)
		nLongestSubSequence := len(longestSubSequence)
		if nLongestSubSequence < countThreshold {
			continue
		}

		// Calculate sequence density
		density := dict.getSequenceDensity(longestSubSequence)
		if density < 0.5 {
			continue
		}

		entryScores = append(entryScores, dictionaryEntryScore{
			ID:                  entryID,
			TokenCount:          tokenCount,
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

func (dict *Dictionary) saveEntry(tx *bbolt.Tx, entry DictionaryEntry) error {
	// Create token from entry
	query := queryFromArabic(entry.ArabicText)
	tokens := tokenizeQuery(query)
	if len(tokens) == 0 {
		return nil
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

	// Save to dictionary
	entryID := int64ToBytes(entry.ID)
	for token, indexes := range compactTokens {
		tokenBucket, err := tx.CreateBucketIfNotExists([]byte(token))
		if err != nil {
			return err
		}

		btIndexes := arrayIntToBytes(indexes)
		err = tokenBucket.Put(entryID, btIndexes)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dict *Dictionary) getLongestSubSequence(sequence []int) []int {
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

func (dict *Dictionary) getSequenceDensity(sequence []int) float64 {
	var sigma float64
	for i := 0; i < len(sequence)-1; i++ {
		tmp := sequence[i+1] - sequence[i]
		sigma += 1 / float64(tmp)
	}

	nSequence := len(sequence)
	return (1 / float64(nSequence-1)) * sigma
}
