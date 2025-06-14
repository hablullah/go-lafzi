package myers

// LCS extracts all elements that are *not* involved in any Delete or Insert
// operations. These are the elements that remain unchanged between s1 and s2.
func LCS[T comparable](s1 []T, edits []Edit) ([]T, []int) {
	// Mark all deleted positions
	deletedIndexes := make(map[int]struct{})
	for _, edit := range edits {
		if edit.Operation == Delete {
			deletedIndexes[edit.OldPosition] = struct{}{}
		}
	}

	// LCS is s1 minus deleted elements
	nLeftover := len(s1) - len(deletedIndexes)
	lcs := make([]T, 0, nLeftover)
	lcsIndexes := make([]int, 0, nLeftover)

	for i, elem := range s1 {
		if _, isDeleted := deletedIndexes[i]; !isDeleted {
			lcs = append(lcs, elem)
			lcsIndexes = append(lcsIndexes, i)
		}
	}

	return lcs, lcsIndexes
}
