package database

import (
	"cmp"
	"database/sql"
	"fmt"
	"slices"

	"github.com/jmoiron/sqlx"
)

type TokenLocation struct {
	DocumentID int    `db:"document_id"`
	TokenID    int    `db:"token_id"`
	Token      string `db:"token"`
	Start      int    `db:"start"`
	End        int    `db:"end"`
}

type TokenLocationGroup struct {
	DocumentID   int
	LastTokenID  int
	Start        int
	End          int
	Count        int
	Positions    []int
	Completeness float64
	Compactness  float64
	Confidence   float64
}

type SearchResult struct {
	DocumentID int
	Confidence float64
	Positions  [][2]int
}

// SearchTokens look for document ids which contains the specified tokens,
// then count how many tokens occured in each document.
func SearchTokens(db *sqlx.DB, minConfidence float64, tokens ...string) (results []SearchResult, err error) {
	// If there are no tokens submitted, stop early
	nToken := len(tokens)
	if nToken == 0 {
		return
	}

	// Start read only transaction
	tx, err := db.Beginx()
	if err != nil {
		err = fmt.Errorf("failed to start transaction: %v", err)
		return
	}

	defer func() {
		if err == sql.ErrNoRows {
			err = nil
		}

		tx.Rollback()
	}()

	// Prepare query
	stmtSearchToken, err := tx.Preparex(`
		SELECT document_id, token, start, end
		FROM document_token
		WHERE token = ?
		ORDER BY document_id, start`)
	if err != nil {
		return
	}

	// Search per token
	tokenLocations := make([][]TokenLocation, len(tokens))
	for i, token := range tokens {
		err = stmtSearchToken.Select(&tokenLocations[i], token)
		if err != nil && err != sql.ErrNoRows {
			return
		}

		for j := range tokenLocations[i] {
			tokenLocations[i][j].TokenID = i
		}
	}

	// Merge all token locations into one
	var nTokenLocations int
	for i := range tokenLocations {
		nTokenLocations += len(tokenLocations[i])
	}

	if nTokenLocations == 0 {
		return
	}

	flatTokenLocations := make([]TokenLocation, 0, nTokenLocations)
	for i := range tokenLocations {
		flatTokenLocations = append(flatTokenLocations, tokenLocations[i]...)
	}

	// Sort the flattened token locations
	slices.SortFunc(flatTokenLocations, func(a, b TokenLocation) int {
		if a.DocumentID != b.DocumentID {
			return cmp.Compare(a.DocumentID, b.DocumentID)
		}

		if a.Start != b.Start {
			return cmp.Compare(a.Start, b.Start)
		}

		return cmp.Compare(a.TokenID, b.TokenID)
	})

	// Compact the sorted token locations
	flatTokenLocations = slices.CompactFunc(flatTokenLocations, func(e1, e2 TokenLocation) bool {
		return e1.DocumentID == e2.DocumentID &&
			e1.Start == e2.Start &&
			e1.End == e2.End
	})

	// Check again after compact. If there are no tokens, stop
	nTokenLocations = len(flatTokenLocations)
	if nTokenLocations == 0 {
		return
	}

	// Create group from token locations
	groups := make([]TokenLocationGroup, 0, nTokenLocations)
	firstTL := flatTokenLocations[0]
	currentGroup := TokenLocationGroup{
		DocumentID:  firstTL.DocumentID,
		LastTokenID: firstTL.TokenID,
		Start:       firstTL.Start,
		End:         firstTL.End,
		Count:       1,
		Positions:   []int{firstTL.Start},
	}

	for i := 1; i < nTokenLocations; i++ {
		tl := flatTokenLocations[i]
		isSameGroup := tl.DocumentID == currentGroup.DocumentID &&
			tl.TokenID > currentGroup.LastTokenID

		if isSameGroup {
			currentGroup.Count++
			currentGroup.End = tl.End
			currentGroup.LastTokenID = tl.TokenID
			currentGroup.Positions = append(currentGroup.Positions, tl.Start)
		} else {
			// We landed on a new group, so save the current one
			currentGroup.Completeness = calcCompleteness(currentGroup.Count, nToken)
			currentGroup.Compactness = calcCompactness(currentGroup.Positions)
			currentGroup.Confidence = currentGroup.Completeness * currentGroup.Compactness
			if currentGroup.Confidence >= minConfidence {
				groups = append(groups, currentGroup)
			}

			// Once saved, reset the current group with the current token
			currentGroup = TokenLocationGroup{
				DocumentID:  tl.DocumentID,
				LastTokenID: tl.TokenID,
				Start:       tl.Start,
				End:         tl.End,
				Count:       1,
				Positions:   []int{tl.Start},
			}
		}
	}

	// Save the last group
	currentGroup.Completeness = calcCompleteness(currentGroup.Count, nToken)
	currentGroup.Compactness = calcCompactness(currentGroup.Positions)
	currentGroup.Confidence = currentGroup.Completeness * currentGroup.Compactness
	if currentGroup.Confidence >= minConfidence {
		groups = append(groups, currentGroup)
	}

	// If there are no groups, stop early
	nGroups := len(groups)
	if nGroups == 0 {
		return
	}

	// Sort the groups based on document ID and its start
	slices.SortFunc(groups, func(a, b TokenLocationGroup) int {
		if a.DocumentID != b.DocumentID {
			return cmp.Compare(a.DocumentID, b.DocumentID)
		}

		if a.Start != b.Start {
			return cmp.Compare(a.Start, b.Start)
		}

		return -cmp.Compare(a.Confidence, b.Confidence)
	})

	// Create the final result
	results = make([]SearchResult, 0, nGroups)
	firstGroup := groups[0]
	currentResult := SearchResult{
		DocumentID: firstGroup.DocumentID,
		Confidence: firstGroup.Confidence,
		Positions:  [][2]int{{firstGroup.Start, firstGroup.End}},
	}

	for i := 1; i < nGroups; i++ {
		gi := groups[i]
		giPos := [2]int{gi.Start, gi.End}

		// Same document as before, so merge it
		if currentResult.DocumentID == gi.DocumentID {
			currentResult.Confidence = max(currentResult.Confidence, gi.Confidence)
			currentResult.Positions = append(currentResult.Positions, giPos)
		} else {
			// We reach different document, so save the current result
			results = append(results, currentResult)

			// Reset the value of current
			currentResult = SearchResult{
				DocumentID: gi.DocumentID,
				Confidence: gi.Confidence,
				Positions:  [][2]int{giPos},
			}
		}
	}

	// Save the leftover result
	results = append(results, currentResult)

	// Sort by best confidence
	slices.SortFunc(results, func(a, b SearchResult) int {
		if a.Confidence != b.Confidence {
			return -cmp.Compare(a.Confidence, b.Confidence)
		}

		return cmp.Compare(a.DocumentID, b.DocumentID)
	})

	return
}

func calcCompleteness(currentCount, expectedCount int) float64 {
	// Penalize when completeness is too small
	score := float64(currentCount) / float64(expectedCount)
	if score <= 0.5 {
		score *= 0.5
	}
	return score
}

func calcCompactness(positions []int) float64 {
	// Handle edge cases: empty positions or single element
	// Single elements have no gaps, so they're perfectly compact
	nPosition := len(positions)
	if nPosition <= 1 {
		return 1.0
	}

	// Calculate gaps average
	var gapSum int
	nGap := nPosition - 1
	for i := range nGap {
		gapSum += positions[i+1] - positions[i]
	}

	// Calculate mean of gaps
	gapMean := float64(gapSum) / float64(nGap)
	if gapMean == 0 {
		return 1
	}

	// Calculate compactness by comparing the mean with ideal gap value. Ideally,
	// gap between token position is at most 3.
	return min(1, 3.0/gapMean)
}
