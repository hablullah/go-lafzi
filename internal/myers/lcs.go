package myers

func LCS[T comparable](s1 []T, edits []Edit) ([]T, []int) {
	var i, j int
	lcs := make([]T, 0, len(s1))
	lcsIndexes := make([]int, 0, len(s1))

	for _, e := range edits {
		// Add unchanged elements before this edit position
		for i < e.OldPosition {
			lcs = append(lcs, s1[i])
			lcsIndexes = append(lcsIndexes, j)
			i, j = i+1, j+1
		}

		// Handle the edit at position i
		if e.OldPosition == i {
			if e.Operation == Delete {
				i++ // Skip deleted element in s1
			}
			j++ // Always advance new position
		}
	}

	// Add remaining unchanged elements
	for i < len(s1) {
		lcs = append(lcs, s1[i])
		lcsIndexes = append(lcsIndexes, j)
		i, j = i+1, j+1
	}

	return lcs, lcsIndexes
}
