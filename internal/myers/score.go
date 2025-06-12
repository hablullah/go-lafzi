package myers

import "math"

// Score calculate the scores for query inside the document
// using LCS length and its compactness.
func Score[T comparable](s1, s2 []T) float64 {
	// Calculate shortest edit using Myers algorithm
	edits := Diff(s1, s2, 0, 0)

	// Fetch LCS from the edits
	lcs, lcsIndexes := LCS(s1, edits)

	// Calculate LCS score
	lcsScore := float64(len(lcs)) / float64(len(s2))

	// Calculate LCS compactness
	compactScore := calcCompactScore(lcsIndexes)

	// Return the final score
	return lcsScore * compactScore
}

func calcCompactScore(indexes []int) float64 {
	var sum float64
	for i := range len(indexes) - 1 {
		nextIdx := float64(indexes[i+1])
		currentIdx := float64(indexes[i])
		idxDiff := math.Abs(nextIdx - currentIdx)
		sum += 1.0 / idxDiff
	}

	return sum / float64(len(indexes)-1)
}
