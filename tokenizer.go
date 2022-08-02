package lafzi

const nGram = 3

func tokenize(s string) []string {
	// Convert string to runes
	chars := []rune(s)

	// Generate token
	var tokens []string
	for i := 0; i <= len(chars)-nGram; i++ {
		token := string(chars[i : i+nGram])
		tokens = append(tokens, token)
	}

	return tokens
}
