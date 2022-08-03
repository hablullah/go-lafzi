package lcs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getLCS(t *testing.T) {
	// Helper function
	check := func(s1, s2, expected string) {
		arr1 := strings.Split(s1, "")
		arr2 := strings.Split(s2, "")
		diff := myersDiff(arr1, arr2, 0, 0)
		lcs, _ := getLCS(arr1, arr2, diff)
		strLcs := strings.Join(lcs, "")
		assert.Equal(t, expected, strLcs, "%q and %q", s1, s2)
	}

	// Example from myer's 1986 paper.
	check("ABCABBA", "CBABAC", "BABA")

	// Wikipedia dynamic programming example.
	check("AGCAT", "GAC", "GA")
	check("XMJYAUZ", "MZJAWXU", "MJAU")

	// Longer examples.
	check("ABCADEFGH", "ABCIJKFGH", "ABCFGH")
	check("ABCDEF1234", "PQRST2UV4", "24")
	check("SABCDE", "SC", "SC")
	check("SABCDE", "SSC", "SC")

	// More exhaustive cases.
	check("", "", "")
	check("", "B", "")
	check("B", "", "")
	check("A", "A", "A")
	check("AB", "AB", "AB")
	check("AB", "ABC", "AB")
	check("ABC", "AB", "AB")
	check("AC", "AXC", "AC")
	check("ABC", "ABX", "AB")
	check("ABC", "ABXY", "AB")
	check("ABXY", "AB", "AB")

	// Example where rune and byte results are identical.
	check("日本語", "日本de語", "日本語")
}
