package search

type bound struct {
	start, end int
}

// Convert from Line bounds to Ofset Bounds
func (b *bound) convert(indexes []int) {
	b.start = indexes[b.start]
	b.end = indexes[b.end]
}

func getLineIndexes(txt string) (idxs []int) {
	// Special case: Empty string has no lines
	if len(txt) == 0 {
		return
	}

	// Not empty. The first line is located at zero

	// Reserve some capacity to avoid exesive allocations
	idxs = make([]int, 1, max(len(txt)/32, 1000))
	idxs[0] = 0

	// Use a limit to avoid empty to avoid
	// conflicts at EOF

	end := len(txt) - 1

	// Append all indices where we can find end of line
	for index, character := range txt {
		if character == '\n' && index != end {
			idxs = append(idxs, index+1)
		}
	}

	return
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
