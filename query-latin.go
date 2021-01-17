package lafzi

import (
	"regexp"
	"strings"
)

var (
	rxHamza1 = regexp.MustCompile(`(?i)(^|\s)([aiu])`)
	rxHamza2 = regexp.MustCompile(`(?i)i([au])`)
	rxHamza3 = regexp.MustCompile(`(?i)u([ai])`)
	rxIdgham = regexp.MustCompile(`(?i)n(\s*[ynmwlr])`)
)

func queryFromLatin(latinText string) string {
	// Substitute vowels
	query := strings.ToLower(latinText)
	query = strings.ReplaceAll(query, "o", "a")
	query = strings.ReplaceAll(query, "e", "i")

	// Merge adjacent letters
	runes := []rune(query)
	newRunes := []rune{}
	for idx, r := range runes {
		nextRune, exist := peek(runes, idx+1)
		if exist && r == nextRune {
			continue
		}

		newRunes = append(newRunes, r)
	}
	query = string(newRunes)

	// Substitute diphtong
	query = strings.ReplaceAll(query, "ai", "ay")
	query = strings.ReplaceAll(query, "au", "aw")

	// Mark hamza around vowels
	query = rxHamza1.ReplaceAllString(query, "$1x$2")
	query = rxHamza2.ReplaceAllString(query, "ix$1")
	query = rxHamza3.ReplaceAllString(query, "ux$1")

	// Substitute ikhfa
	query = strings.ReplaceAll(query, "ng", "n")

	// Substitute iqlab
	query = strings.ReplaceAll(query, "nb", "mb")

	// Substitute idgham
	query = rxIdgham.ReplaceAllString(query, "$1")

	// Convert to match with common phonetics
	query = strings.ReplaceAll(query, "sh", "s")
	query = strings.ReplaceAll(query, "ts", "s")
	query = strings.ReplaceAll(query, "sy", "s")
	query = strings.ReplaceAll(query, "kh", "h")
	query = strings.ReplaceAll(query, "ch", "h")
	query = strings.ReplaceAll(query, "zh", "z")
	query = strings.ReplaceAll(query, "dz", "z")
	query = strings.ReplaceAll(query, "dh", "d")
	query = strings.ReplaceAll(query, "th", "t")
	query = strings.ReplaceAll(query, "gh", "g")
	query = strings.ReplaceAll(query, "v", "f")
	query = strings.ReplaceAll(query, "p", "f")
	query = strings.ReplaceAll(query, "q", "k")
	query = strings.ReplaceAll(query, "j", "z")
	query = strings.ReplaceAll(query, "‘", "x")
	query = strings.ReplaceAll(query, "`", "x")
	query = strings.ReplaceAll(query, "'", "x")

	// Remove spaces
	query = strings.Join(strings.Fields(query), "")
	return query
}
