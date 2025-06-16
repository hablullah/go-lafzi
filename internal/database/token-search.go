package database

import (
	"cmp"
	"database/sql"
	"fmt"
	"math"
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
	DocumentID  int
	LastTokenID int
	Start       int
	End         int
	Count       int
	Positions   []int
	TokenScore  float64
	Compactness float64
	Confidence  float64
}

type SearchResult struct {
	DocumentID int
	Start      int
	End        int
	Confidence float64
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
	current := TokenLocationGroup{
		DocumentID:  firstTL.DocumentID,
		LastTokenID: firstTL.TokenID,
		Start:       firstTL.Start,
		End:         firstTL.End,
		Count:       1,
		Positions:   []int{firstTL.Start},
	}

	for i := 1; i < nTokenLocations; i++ {
		tl := flatTokenLocations[i]
		isSameGroup := tl.DocumentID == current.DocumentID &&
			tl.TokenID > current.LastTokenID

		if isSameGroup {
			current.Count++
			current.End = tl.End
			current.LastTokenID = tl.TokenID
			current.Positions = append(current.Positions, tl.Start)
		} else {
			// We landed on a new group, so save the current one
			current.TokenScore = float64(current.Count) / float64(nToken)
			current.Compactness = calcCompactness(current.Positions)
			current.Confidence = current.TokenScore * current.Compactness
			if current.Confidence >= minConfidence {
				groups = append(groups, current)
			}

			// Once saved, reset the current group with the current token
			current = TokenLocationGroup{
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
	current.TokenScore = float64(current.Count) / float64(nToken)
	current.Compactness = calcCompactness(current.Positions)
	current.Confidence = current.TokenScore * current.Compactness
	if current.Confidence >= minConfidence {
		groups = append(groups, current)
	}

	// Sort by the best score
	slices.SortFunc(groups, func(a, b TokenLocationGroup) int {
		if a.Confidence != b.Confidence {
			return -cmp.Compare(a.Confidence, b.Confidence)
		}

		return cmp.Compare(a.DocumentID, b.DocumentID)
	})

	// Create the final result
	results = make([]SearchResult, len(groups))
	for i, g := range groups {
		results[i] = SearchResult{
			DocumentID: g.DocumentID,
			Start:      g.Start,
			End:        g.End,
			Confidence: g.Confidence,
		}
	}

	return
}

// calcCompactness calculates the compactness of an integer slice, in this case the
// position of tokens.
//
// The method is based on the coefficient of variation (CV = σ/μ), which measures
// relative dispersion by normalizing standard deviation against the mean. Since
// higher CV indicates more spread-out data (higher CV = less compact), we invert
// this relationship to create a compactness score (lower CV = more compact).
//
// Mathematical foundation:
// - Coefficient of Variation: CV = σ/μ (standard deviation / mean)
// - Compactness Score: 1 / (1 + CV) = 1 / (1 + σ/μ)
//
// Why 1/(1 + σ/μ) instead of 1/(σ/μ)?
//  1. Avoids division by zero when σ = 0 (perfectly compact data)
//  2. Bounds output to [0, 1] range for intuitive interpretation
//  3. Provides smooth, continuous behavior without extreme values
//  4. When σ/μ = 0 (perfect compactness): result = 1.0
//  5. When σ/μ → ∞ (very spread out): result → 0.0
//
// Returns compactness score between 0.0 and 1.0.
//
// Edge cases:
// - Empty or single-element slices return 1.0 (considered perfectly compact)
// - Zero mean returns 1.0 (avoids division by zero in CV calculation)
func calcCompactness(slice []int) float64 {
	// Handle edge cases: empty slice or single element
	// Single elements have no variance, so they're perfectly compact
	nSlice := len(slice)
	if nSlice <= 1 {
		return 1.0
	}

	// Calculate arithmetic mean (μ)
	sum := 0
	for _, v := range slice {
		sum += v
	}
	mean := float64(sum) / float64(nSlice)

	// Calculate standard deviation
	var sumSquaredDiffs float64
	for _, v := range slice {
		diff := float64(v) - mean
		sumSquaredDiffs += diff * diff
	}

	variance := sumSquaredDiffs / float64(nSlice)
	stdDev := math.Sqrt(variance)

	// Handle edge case: zero mean would make coefficient of variation undefined
	// In practice, this is rare with typical integer sequences
	if mean == 0 {
		return 1.0
	}

	// Calculate compactness using inverse coefficient of variation
	// CV = σ/μ (coefficient of variation)
	// Compactness = 1 / (1 + CV)
	//
	// This transformation:
	// - Maps [0, ∞) CV values to (0, 1] compactness values
	// - Ensures higher dispersion → lower compactness score
	// - Provides intuitive 0-1 scale for comparison
	coefficientOfVariation := stdDev / mean
	return 1.0 / (1.0 + coefficientOfVariation)
}
