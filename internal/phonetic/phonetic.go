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
	rxHamzahA       = regexp.MustCompile(`(?i)a([iu])`)
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
		case '\'', '`',
			'\u2019', // right single quotation mark
			'\u02bc', // modifier letter apostrophe
			'\u02bb', // modifier letter turned comma
			'\u055a', // armenian apostrophe
			'\ua78c', // latin small letter saltillo
			'\u2032', // prime
			'\u2035', // reversed prime
			'\u02b9', // modifier letter prime
			'\uff07', // fullwidth apostrophe
			'\u2018': // left single quotation mark
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
	s = rxHamzahA.ReplaceAllString(s, "ax$1")
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

	// Normalize diphtong
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
	// Convert common separator into spaces
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")

	// Split sentence into words
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
	// Make sure string not empty
	runes := []rune(s)
	if len(runes) == 0 {
		return s
	}

	// Create strings builder and init with the first rune
	var sb strings.Builder
	sb.WriteRune(runes[0])

	// Add the next runes
	for i := 1; i < len(runes); i++ {
		if runes[i] != runes[i-1] {
			sb.WriteRune(runes[i])
		}
	}

	return sb.String()
}
