package lcs

import "math"

// Score calculate the scores for query inside the document
// using LCS length and its compactness.
func Score(docTokens, queryTokens []string) float64 {
	// Calculate shortest edit using Myers algorithm
	edits := myersDiff(docTokens, queryTokens, 0, 0)

	// Fetch LCS from the edits
	lcs, lcsIndexes := getLCS(docTokens, queryTokens, edits)

	// Calculate LCS score
	lcsScore := float64(len(lcs)) / float64(len(queryTokens))

	// Calculate LCS compactness
	compactScore := calcCompactScore(lcsIndexes)

	// Return the final score
	return lcsScore * compactScore
}

func calcCompactScore(indexes []int) float64 {
	var sum float64
	for i := 0; i < len(indexes)-1; i++ {
		nextIdx := float64(indexes[i+1])
		currentIdx := float64(indexes[i])
		idxDiff := math.Abs(nextIdx - currentIdx)
		sum += 1.0 / idxDiff
	}

	return sum / float64(len(indexes)-1)
}
