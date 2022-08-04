package lcs

// Score calculate the scores for query inside the document
// using LCS length and its compactness.
func Score(docTokens, queryTokens []string) float64 {
	// Calculate shortest edit using Myers algorithm
	edits := myersDiff(docTokens, queryTokens, 0, 0)

	// Fetch LCS from the edits
	lcs, _ := getLCS(docTokens, queryTokens, edits)

	// Calculate LCS score
	lcsScore := float64(len(lcs)) / float64(len(queryTokens))
	return lcsScore
}
