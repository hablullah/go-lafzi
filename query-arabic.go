package lafzi

import "strings"

func QueryFromArabic(arabicText string) string {
	return queryFromArabic(arabicText)
}

func queryFromArabic(arabicText string) string {
	query := strings.Join(strings.Fields(arabicText), "")
	query = adjustEndSentence(query)
	query = replaceTanween(query)
	query = toPhonetic(query)
	query = removeUnvoweled(query)
	query = removeMadd(query)
	query = strings.ReplaceAll(query, "0", "")
	query = removeShadda(query)
	query = adjustTajweed(query)
	return query
}

func adjustEndSentence(str string) string {
	strRunes := []rune(str)
	if len(strRunes) <= 1 {
		return str
	}

	lastIdx := len(strRunes) - 1
	lastRune := strRunes[lastIdx]

	switch lastRune {
	case alef:
		prev, exist := peek(strRunes, lastIdx-1)
		if exist && prev == fathatan {
			newRunes := append(strRunes[:lastIdx-1], fatha)
			return string(newRunes)
		}

	case tehMarbuta:
		newRunes := append(strRunes[:lastIdx], heh)
		return string(newRunes)

	case fatha, damma, kasra:
		prev, exist := peek(strRunes, lastIdx-1)
		if exist && prev != alef && prev != alefMaksura {
			newRunes := append(strRunes[:lastIdx], sukun)
			return string(newRunes)
		}
	}

	return str
}

func replaceTanween(str string) string {
	noonSukun := string(noon) + string(sukun)
	str = strings.ReplaceAll(str, string(fathatan), string(fatha)+noonSukun)
	str = strings.ReplaceAll(str, string(dammatan), string(damma)+noonSukun)
	str = strings.ReplaceAll(str, string(kasratan), string(kasra)+noonSukun)
	return str
}

func toPhonetic(str string) string {
	runes := []rune(str)
	phoneticRunes := []rune{}

	for _, r := range runes {
		if phonetic, exist := phonetics[r]; exist {
			phoneticRunes = append(phoneticRunes, phonetic)
		}
	}

	return string(phoneticRunes)
}

func removeMadd(str string) string {
	str = strings.ReplaceAll(str, "iy0", "i")
	str = strings.ReplaceAll(str, "uw0", "u")
	return str
}

func removeShadda(str string) string {
	runes := []rune(str)
	newRunes := []rune{}
	for idx, r := range runes {
		nextRune, exist := peek(runes, idx+1)
		if exist && r == nextRune {
			continue
		}

		newRunes = append(newRunes, r)
	}

	return string(newRunes)
}

func removeUnvoweled(str string) string {
	runes := []rune(str)
	newRunes := []rune{}
	for idx, r := range runes {
		// If it's already vowel, just put it back
		switch r {
		case 'a', 'i', 'u', '0':
			newRunes = append(newRunes, r)
			continue
		}

		// If it's not vowel and the next is not vowel as well, skip
		nextRune, exist := peek(runes, idx+1)
		if exist {
			switch nextRune {
			case 'a', 'i', 'u', '0':
			default:
				continue
			}
		}

		newRunes = append(newRunes, r)
	}

	return string(newRunes)
}

func adjustTajweed(str string) string {
	// Replace iqlab
	str = strings.ReplaceAll(str, "nb", "mb")

	// Replace idgham
	str = strings.ReplaceAll(str, "ny", "y")
	str = strings.ReplaceAll(str, "nn", "n")
	str = strings.ReplaceAll(str, "nm", "m")
	str = strings.ReplaceAll(str, "nw", "w")
	str = strings.ReplaceAll(str, "nl", "l")
	str = strings.ReplaceAll(str, "nr", "r")

	return str
}
