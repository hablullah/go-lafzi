package phonetic

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

// FromArabic convert Arabic string into its phonetic.
func FromArabic(s string) string {
	// If string empty, stop early
	if s == "" {
		return s
	}

	// Normalize unicode
	s = norm.NFKD.String(s)
	s = norm.NFKC.String(s)

	// Convert Arabic chars into its phonetic
	var sb strings.Builder
	for _, r := range s {
		replacementRunes := transformArabicRune(r)
		for _, rr := range replacementRunes {
			sb.WriteRune(rr)
		}
	}

	// Normalize the converted phonetic
	return Normalize(sb.String())
}

func transformArabicRune(r rune) []rune {
	switch r {
	case jeem, thal, zain, zah:
		return []rune{'z'}
	case hah, khah, heh:
		return []rune{'h'}
	case hamza, alef, ain,
		alefWithHamzaAbove, alefWithHamzaBelow,
		yehWithHamzaAbove, wawWithHamzaAbove:
		return []rune{'x'}
	case theh, seen, sheen, sad:
		return []rune{'s'}
	case dal, dad:
		return []rune{'d'}
	case tehMarbuta, teh, tah:
		return []rune{'t'}
	case qaf, kaf:
		return []rune{'k'}
	case ghain:
		return []rune{'g'}
	case feh:
		return []rune{'f'}
	case meem:
		return []rune{'m'}
	case noon:
		return []rune{'n'}
	case lam:
		return []rune{'l'}
	case beh:
		return []rune{'b'}
	case yeh:
		return []rune{'y'}
	case waw:
		return []rune{'w'}
	case reh:
		return []rune{'r'}
	case fathatan:
		return []rune{'a', 'n'}
	case dammatan:
		return []rune{'u', 'n'}
	case kasratan:
		return []rune{'i', 'n'}
	case fatha:
		return []rune{'a'}
	case damma:
		return []rune{'u'}
	case kasra:
		return []rune{'i'}
	case sukun:
		return []rune{'0'}
	default:
		return nil
	}
}

const (
	hamza              = '\u0621'
	alefWithHamzaAbove = '\u0623'
	wawWithHamzaAbove  = '\u0624'
	alefWithHamzaBelow = '\u0625'
	yehWithHamzaAbove  = '\u0626'
	alef               = '\u0627'
	beh                = '\u0628'
	tehMarbuta         = '\u0629'
	teh                = '\u062A'
	theh               = '\u062B'
	jeem               = '\u062C'
	hah                = '\u062D'
	khah               = '\u062E'
	dal                = '\u062F'
	thal               = '\u0630'
	reh                = '\u0631'
	zain               = '\u0632'
	seen               = '\u0633'
	sheen              = '\u0634'
	sad                = '\u0635'
	dad                = '\u0636'
	tah                = '\u0637'
	zah                = '\u0638'
	ain                = '\u0639'
	ghain              = '\u063A'
	feh                = '\u0641'
	qaf                = '\u0642'
	kaf                = '\u0643'
	lam                = '\u0644'
	meem               = '\u0645'
	noon               = '\u0646'
	heh                = '\u0647'
	waw                = '\u0648'
	yeh                = '\u064A'
	fathatan           = '\u064B'
	dammatan           = '\u064C'
	kasratan           = '\u064D'
	fatha              = '\u064E'
	damma              = '\u064F'
	kasra              = '\u0650'
	sukun              = '\u0652'
)
