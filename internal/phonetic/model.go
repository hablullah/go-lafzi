package phonetic

import (
	"slices"
	"strings"
)

type Data struct {
	Rune rune
	Pos  int
}

type Group []Data

type NGram []Group

func (g Group) String() string {
	var sb strings.Builder
	for _, d := range g {
		sb.WriteRune(d.Rune)
	}
	return sb.String()
}

func (g Group) Boundary() (int, int) {
	if len(g) == 0 {
		return -1, -1
	}

	pos := g[0].Pos
	start, end := pos, pos
	for i := 1; i < len(g); i++ {
		pos = g[i].Pos
		if pos < start {
			start = pos
		}
		if pos > end {
			end = pos
		}
	}

	return start, end + 1
}

// Split splits the group into several n-grams of specified size
func (g Group) Split(n int) NGram {
	// Make sure n is not zero
	if n <= 0 {
		return nil
	}

	// Make sure group is longer than n
	if len(g) < n {
		return nil
	}

	// Pre-allocate slice with exact capacity needed
	numNGrams := len(g) - n + 1
	ngrams := make([]Group, 0, numNGrams)

	for i := 0; i <= len(g)-n; i++ {
		ngram := slices.Clone(g[i : i+n])
		ngrams = append(ngrams, ngram)
	}

	return ngrams
}

func (ng NGram) Texts() []string {
	texts := make([]string, len(ng))
	for i, token := range ng {
		texts[i] = token.String()
	}
	return texts
}
