package myers

import (
	"fmt"
	"slices"
)

type Operation uint8

const (
	Delete Operation = iota + 1
	Insert
)

type Edit[T comparable] struct {
	Operation   Operation
	OldPosition int
	NewElement  T
}

// Diff returns minimal edit operation to convert s1 into s2,
// requring O(min(len(e),len(f))) space and O(min(len(e),len(f)) * D)
// worst-case execution time where D is the number of differences.
// Taken from http://blog.robertelder.org/diff-algorithm/
func Diff[T comparable](s1, s2 []T, i, j int) []Edit[T] {
	N := len(s1)
	M := len(s2)
	L := N + M
	Z := 2*min(N, M) + 2

	// Early returns for base case

	// Case 1: s1 empty, which mean all inserts
	if N == 0 {
		edits := make([]Edit[T], M)
		for m := range M {
			edits[m] = Edit[T]{
				Operation:   Insert,
				OldPosition: i,
				NewElement:  s2[m],
			}
		}
		return edits
	}

	// Case 2: s2 empty, which mean all deletes
	if M == 0 {
		edits := make([]Edit[T], N)
		for n := range N {
			edits[n] = Edit[T]{
				Operation:   Delete,
				OldPosition: i + n,
			}
		}
		return edits
	}

	// Case 3: s1 and s2 equal, nothing todo
	if N == M && slices.Equal(s1, s2) {
		return nil
	}

	// The real Myers algorithm implementation
	if N > 0 && M > 0 {
		w := N - M
		g := make([]int, Z)
		p := make([]int, Z)

		for h := range L/2 + L%2 + 1 {
			for r := range 2 {
				var c, d []int
				var o, m int
				if r == 0 {
					c, d, o, m = g, p, 1, 1
				} else {
					c, d, o, m = p, g, 0, -1
				}

				for k := -(h - 2*max(0, h-M)); k < h-2*max(0, h-N)+1; k += 2 {
					var a int
					if k == -h || k != h && getItem(c, (k-1)%Z) < getItem(c, (k+1)%Z) {
						a = getItem(c, (k+1)%Z)
					} else {
						a = getItem(c, (k-1)%Z) + 1
					}

					b := a - k
					s, t := a, b
					for a < N && b < M && getItem(s1, (1-o)*N+m*a+(o-1)) == getItem(s2, (1-o)*M+m*b+(o-1)) {
						a, b = a+1, b+1
					}

					cIdx := k % Z
					if cIdx < 0 {
						cIdx = len(c) + cIdx
					}

					c[cIdx] = a
					z := -(k - w)

					if L%2 == o && z >= -(h-o) && z <= h-o && getItem(c, k%Z)+getItem(d, z%Z) >= N {
						var D, x, y, u, v int
						if o == 1 {
							D, x, y, u, v = 2*h-1, s, t, a, b
						} else {
							D, x, y, u, v = 2*h, N-a, M-b, N-s, M-t
						}

						if D > 1 || (x != u && y != v) {
							left := Diff(s1[0:x], s2[0:y], i, j)
							right := Diff(s1[u:N], s2[v:M], i+u, j+v)

							edits := make([]Edit[T], 0, len(left)+len(right))
							edits = append(edits, left...)
							edits = append(edits, right...)
							return edits
						} else if M > N {
							return Diff(nil, s2[N:M], i+N, j+N)
						} else if M < N {
							return Diff(s1[M:N], nil, i+M, j+M)
						} else {
							return nil
						}
					}
				}
			}
		}
	} else if N > 0 {
		edits := make([]Edit[T], 0, N)
		for n := range N {
			edits = append(edits, Edit[T]{
				Operation:   Delete,
				OldPosition: i + n,
			})
		}
		return edits
	} else {
		edits := make([]Edit[T], 0, M)
		for m := range M {
			edits = append(edits, Edit[T]{
				Operation:   Insert,
				OldPosition: i,
				NewElement:  s2[m],
			})
		}
		return edits
	}

	panic(fmt.Errorf("this should be unreachable"))
}

func getItem[T comparable](arr []T, idx int) T {
	if idx < 0 {
		return arr[len(arr)+idx]
	}
	return arr[idx]
}
