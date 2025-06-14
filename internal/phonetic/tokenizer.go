package phonetic

// NGrams splits a string into n-grams of specified size
func NGrams(s string, n int) []string {
	// Make sure n is not zero
	if n <= 0 {
		return nil
	}

	// Make sure string is longer than n
	runes := []rune(s)
	if len(runes) < n {
		return nil
	}

	// Pre-allocate slice with exact capacity needed
	numNGrams := len(runes) - n + 1
	ngrams := make([]string, 0, numNGrams)

	for i := 0; i <= len(runes)-n; i++ {
		ngram := string(runes[i : i+n])
		ngrams = append(ngrams, ngram)
	}

	return ngrams
}
