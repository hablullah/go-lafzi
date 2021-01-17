package lafzi

func peek(runes []rune, idx int) (rune, bool) {
	if idx < 0 || idx >= len(runes) {
		return 0, false
	}

	return runes[idx], true
}
