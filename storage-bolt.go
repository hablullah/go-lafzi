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

func (bs *boltStorage) saveDocuments(documents ...Document) error {
	return bs.Update(func(tx *bbolt.Tx) (err error) {
		for _, doc := range documents {
			err = bs.saveDocument(tx, doc)
			if err != nil {
				return
			}
		}
		return
	})
}

func (bs *boltStorage) findTokens(tokens ...string) ([]documentTokens, error) {
	// For each token, find the documents that contains such token,
	// also the position of that token within the document.
	docTokenCount := map[int64]int{}
	docTokenIndexes := map[int64]map[int]struct{}{}

	bs.View(func(tx *bbolt.Tx) error {
		for _, token := range tokens {
			tokenBucket := tx.Bucket([]byte(token))
			if tokenBucket == nil {
				continue
			}

			tokenBucket.ForEach(func(btDocID, _ []byte) error {
				docBucket := tokenBucket.Bucket(btDocID)
				if docBucket == nil {
					return nil
				}

				tokenIndexes := []int{}
				docBucket.ForEach(func(btIdx, _ []byte) error {
					idx := int(bytesToInt64(btIdx))
					tokenIndexes = append(tokenIndexes, idx)
					return nil
				})

				docID := bytesToInt64(btDocID)
				existingIndexes := docTokenIndexes[docID]
				if existingIndexes == nil {
					existingIndexes = map[int]struct{}{}
				}

				for _, idx := range tokenIndexes {
					existingIndexes[idx] = struct{}{}
				}

				docTokenCount[docID]++
				docTokenIndexes[docID] = existingIndexes
				return nil
			})
		}

		return nil
	})

	// Convert map of token count and indexes to array
	result := []documentTokens{}
	for docID, indexes := range docTokenIndexes {
		arrIndexes := []int{}
		for idx := range indexes {
			arrIndexes = append(arrIndexes, idx)
		}
		sort.Ints(arrIndexes)

		result = append(result, documentTokens{
			ID:           docID,
			TokenCount:   docTokenCount[docID],
			TokenIndexes: arrIndexes,
		})
	}

	return result, nil
}

func (bs *boltStorage) saveDocument(tx *bbolt.Tx, doc Document) error {
	// Create token from document
	query := queryFromArabic(doc.ArabicText)
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
	docID := int64ToBytes(doc.ID)
	for token, indexes := range compactTokens {
		tokenBucket, err := tx.CreateBucketIfNotExists([]byte(token))
		if err != nil {
			return err
		}

		docBucket, err := tokenBucket.CreateBucketIfNotExists(docID)
		if err != nil {
			return err
		}

		for _, idx := range indexes {
			err = docBucket.Put(int64ToBytes(int64(idx)), nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
