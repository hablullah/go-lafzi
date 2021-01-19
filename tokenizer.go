package lafzi

func tokenizeQuery(query string) []string {
	if len(query) < 3 {
		return nil
	}

	tokenList := []string{}
	for i := 0; i <= len(query)-3; i++ {
		token := query[i : i+3]
		tokenList = append(tokenList, token)
	}

	return tokenList
}
