package lafzi

import (
	"os"
	"sort"

	"go.etcd.io/bbolt"
)

type boltStorage struct {
	*bbolt.DB
}

func newBoltStorage(path string) (*boltStorage, error) {
	db, err := bbolt.Open(path, os.ModePerm, bbolt.DefaultOptions)
	if err != nil {
		return nil, err
	}

	return &boltStorage{db}, nil
}

func (bs *boltStorage) close() {
	bs.Close()
}

func (bs *boltStorage) saveEntries(entries ...DatabaseEntry) error {
	return bs.Update(func(tx *bbolt.Tx) (err error) {
		for _, entry := range entries {
			err = bs.saveEntry(tx, entry)
			if err != nil {
				return
			}
		}
		return
	})
}

func (bs *boltStorage) findTokens(tokens ...string) ([]dbEntryTokens, error) {
	// For each token, find the dictionary entry that contains such token,
	// also the position of that token within the dictionary entry.
	entryTokenCount := map[int64]int{}
	entryTokenIndexes := map[int64]map[int]struct{}{}

	bs.View(func(tx *bbolt.Tx) error {
		for _, token := range tokens {
			tokenBucket := tx.Bucket([]byte(token))
			if tokenBucket == nil {
				continue
			}

			tokenBucket.ForEach(func(btEntryID, _ []byte) error {
				entryBucket := tokenBucket.Bucket(btEntryID)
				if entryBucket == nil {
					return nil
				}

				tokenIndexes := []int{}
				entryBucket.ForEach(func(btIdx, _ []byte) error {
					idx := int(bytesToInt64(btIdx))
					tokenIndexes = append(tokenIndexes, idx)
					return nil
				})

				entryID := bytesToInt64(btEntryID)
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

	// Convert map of token count and indexes to array
	result := []dbEntryTokens{}
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

	return result, nil
}

func (bs *boltStorage) saveEntry(tx *bbolt.Tx, entry DatabaseEntry) error {
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

	// Save to storage
	entryID := int64ToBytes(entry.ID)
	for token, indexes := range compactTokens {
		tokenBucket, err := tx.CreateBucketIfNotExists([]byte(token))
		if err != nil {
			return err
		}

		entryBucket, err := tokenBucket.CreateBucketIfNotExists(entryID)
		if err != nil {
			return err
		}

		for _, idx := range indexes {
			err = entryBucket.Put(int64ToBytes(int64(idx)), nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
