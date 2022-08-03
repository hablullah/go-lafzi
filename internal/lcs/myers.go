package lcs

import "fmt"

type EditOp uint8

const (
	Delete EditOp = iota
	Insert
)

type Edit struct {
	Op    EditOp
	Index int
	From  int
}

// myersDiff returns minimal edit operation to convert s1 into s2,
// requring O(min(len(e),len(f))) space and O(min(len(e),len(f)) * D)
// worst-case execution time where D is the number of differences.
// Taken from http://blog.robertelder.org/diff-algorithm/
func myersDiff(s1, s2 []string, i, j int) []Edit {
	N := len(s1)
	M := len(s2)
	L := N + M
	Z := 2*min(N, M) + 2

	if N > 0 && M > 0 {
		w := N - M
		g := make([]int, Z)
		p := make([]int, Z)

		for h := 0; h < L/2+L%2+1; h++ {
			for r := 0; r < 2; r++ {
				var c, d []int
				var o, m int
				if r == 0 {
					c, d, o, m = g, p, 1, 1
				} else {
					c, d, o, m = p, g, 0, -1
				}

				for k := -(h - 2*max(0, h-M)); k < h-2*max(0, h-N)+1; k += 2 {
					var a int
					if k == -h || k != h && getInt(c, (k-1)%Z) < getInt(c, (k+1)%Z) {
						a = getInt(c, (k+1)%Z)
					} else {
						a = getInt(c, (k-1)%Z) + 1
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

					if L%2 == o && z >= -(h-o) && z <= h-o && getInt(c, k%Z)+getInt(d, z%Z) >= N {
						var D, x, y, u, v int
						if o == 1 {
							D, x, y, u, v = 2*h-1, s, t, a, b
						} else {
							D, x, y, u, v = 2*h, N-a, M-b, N-s, M-t
						}

						if D > 1 || (x != u && y != v) {
							return append(
								myersDiff(s1[0:x], s2[0:y], i, j),
								myersDiff(s1[u:N], s2[v:M], i+u, j+v)...)
						} else if M > N {
							return myersDiff(nil, s2[N:M], i+N, j+N)
						} else if M < N {
							return myersDiff(s1[M:N], nil, i+M, j+M)
						} else {
							return nil
						}
					}
				}
			}
		}
	} else if N > 0 {
		var edits []Edit
		for n := 0; n < N; n++ {
			edits = append(edits, Edit{
				Op:    Delete,
				Index: i + n,
			})
		}
		return edits
	} else {
		var edits []Edit
		for n := 0; n < M; n++ {
			edits = append(edits, Edit{
				Op:    Insert,
				Index: i,
				From:  j + n,
			})
		}
		return edits
	}

	panic(fmt.Errorf("this should be unreachable"))
}

func getLCS(s1, s2 []string, edits []Edit) ([]string, []int) {
	var i, j int
	var lcs []string
	var lcsIndexes []int

	for _, e := range edits {
		for i < e.Index {
			lcs = append(lcs, s1[i])
			lcsIndexes = append(lcsIndexes, j)
			i, j = i+1, j+1
		}

		if e.Index == i {
			if e.Op == Delete {
				i++
			}
			j++
		}
	}

	for i < len(s1) {
		lcs = append(lcs, s1[i])
		lcsIndexes = append(lcsIndexes, j)
		i, j = i+1, j+1
	}

	return lcs, lcsIndexes
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getInt(arr []int, idx int) int {
	if idx < 0 {
		return arr[len(arr)+idx]
	}
	return arr[idx]
}

func getItem(arr []string, idx int) string {
	if idx < 0 {
		return arr[len(arr)+idx]
	}
	return arr[idx]
}
