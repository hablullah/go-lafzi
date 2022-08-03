package phonetic

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	rxVowelPrefix   = regexp.MustCompile(`(?i)^[aiu]`)
	rxHamzahI       = regexp.MustCompile(`(?i)i([au])`)
	rxHamzahU       = regexp.MustCompile(`(?i)u([ai])`)
	rxHamzahPrefix  = regexp.MustCompile(`(?i)^([^aiu0])?([^aiu0])0?([^aiu0])([aiu])`)
	rxMaddaA        = regexp.MustCompile(`(?i)ax([^aiu]|$)`)
	rxMaddaI        = regexp.MustCompile(`(?i)iy([^aiu]|$)`)
	rxMaddaU        = regexp.MustCompile(`(?i)uw([^aiu]|$)`)
	rxTajweedIkhfa  = regexp.MustCompile(`(?i)n0?g`)
	rxTajweedIqlab  = regexp.MustCompile(`(?i)n0?b`)
	rxTajweedIdgham = regexp.MustCompile(`(?i)n0?([ynmwlr])`)
	rxAlifLamSyams  = regexp.MustCompile(`(?i)x([aiu]?)l([zsdtnlr])`)
	rxUnusedX       = regexp.MustCompile(`(?i)x([^aiu0])`)
	rxHangingVowel  = regexp.MustCompile(`(?i)[aiu]$`)

	unicodeNormalizer = transform.Chain(
		norm.NFKD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFKC,
	)

	similarSoundingRunesCleaner = runes.Map(func(r rune) rune {
		switch r {
		case 'o':
			return 'a'
		case 'e':
			return 'i'
		case 'v', 'p':
			return 'f'
		case 'q':
			return 'k'
		case 'j':
			return 'z'
		case '\'', '`', '‘':
			return 'x'
		default:
			return r
		}
	})

	invalidPhoneticRunesCleaner = runes.Remove(
		runes.Predicate(func(r rune) bool {
			switch r {
			case 'z', 'h', 'x', 's', 'd', 't', 'k', 'g',
				'f', 'm', 'n', 'l', 'b', 'y', 'w', 'r',
				'a', 'u', 'i', '0':
				return false
			default:
				return true
			}
		}),
	)
)

// Normalize the phonetics by using several heuristics.
func Normalize(s string) string {
	// Normalize unicode
	normal, _, err := transform.String(unicodeNormalizer, s)
	if err == nil {
		s = normal
	}

	// Convert string to lowercase
	s = strings.ToLower(s)

	// Normalize similar sounding runes, e.g. 'p' => 'f', 'e' => 'i'
	s = similarSoundingRunesCleaner.String(s)

	// Mark possible hamzah location
	s = rxHamzahI.ReplaceAllString(s, "ix$1")
	s = rxHamzahU.ReplaceAllString(s, "ux$1")

	// Normalize space and its vowels, e.g. 'amman antum' => 'xammanxantum'
	s = normalizeSpaces(s)

	// Remove invalid (or disallowed) phonetic runes
	s = invalidPhoneticRunesCleaner.String(s)

	// Normalize alif or hamzah 'x' in prefix
	s = rxHamzahPrefix.ReplaceAllString(s, "${1}${4}${2}${3}${4}")

	// Remove invisible 'al' or in tahweed called 'alif lam syamsiah'
	s = rxAlifLamSyams.ReplaceAllString(s, "${1}${2}")

	// Remove unused alif or hamzah 'x', e.g. 'rabixlxalamin' => 'rabilalamin'
	s = rxUnusedX.ReplaceAllString(s, "${1}")
	s = strings.TrimPrefix(s, "x")

	// Remove madda i.e. vowel that spelled a bit long
	s = rxMaddaA.ReplaceAllString(s, "a${1}")
	s = rxMaddaI.ReplaceAllString(s, "i${1}")
	s = rxMaddaU.ReplaceAllString(s, "u${1}")

	// Apply tajweed rules, i.e. ikhfa, iqlab, and idgham
	s = rxTajweedIkhfa.ReplaceAllString(s, "n0")
	s = rxTajweedIqlab.ReplaceAllString(s, "m0b")
	s = rxTajweedIdgham.ReplaceAllString(s, "$1")

	// Remove sukun (stop mark)
	s = strings.ReplaceAll(s, "0", "")

	// Remove hanging vowel, e.g. 'amanu' => 'aman'
	s = rxHangingVowel.ReplaceAllString(s, "")

	// Normalize diphtong
	s = strings.ReplaceAll(s, "ai", "ay")
	s = strings.ReplaceAll(s, "au", "aw")
	s = strings.ReplaceAll(s, "sh", "s")
	s = strings.ReplaceAll(s, "ts", "s")
	s = strings.ReplaceAll(s, "sy", "s")
	s = strings.ReplaceAll(s, "kh", "h")
	s = strings.ReplaceAll(s, "ch", "h")
	s = strings.ReplaceAll(s, "zh", "z")
	s = strings.ReplaceAll(s, "dz", "z")
	s = strings.ReplaceAll(s, "dh", "d")
	s = strings.ReplaceAll(s, "th", "t")
	s = strings.ReplaceAll(s, "gh", "g")

	// Merge identic adjacent runes, e.g. 'amman' => 'aman'
	s = mergeIdenticAdjacentRunes(s)

	return s
}

func normalizeSpaces(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		// If word started with 'a', 'i', 'u', add 'x' in front of it.
		if rxVowelPrefix.MatchString(word) {
			words[i] = "x" + word
			continue
		}
	}
	return strings.Join(words, "")
}

func mergeIdenticAdjacentRunes(s string) string {
	var sb strings.Builder

	src := []rune(s)
	for i := 0; i < len(src)-1; i++ {
		if src[i] == src[i+1] {
			continue
		}
		sb.WriteRune(src[i])
	}

	sb.WriteRune(src[len(src)-1])
	return sb.String()
}
