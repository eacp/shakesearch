package search

import (
	"bytes"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"strings"
)

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
	lineIndexes   []int // Cache for the location of the lines
	limits        bound
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)

	// Create a lower case version of the data
	// and use that for the sufix array

	s.SuffixArray = suffixarray.New(bytes.ToLower(dat))

	// Calculate the line indices to be used
	// when requesting individual works
	s.lineIndexes = getLineIndexes(s.CompleteWorks)

	// By default the end is the end of the file
	s.limits.end = len(s.CompleteWorks) - 1

	return nil
}

func (s *Searcher) Search(query string) []string {
	idxs := s.SuffixArray.Lookup([]byte(strings.ToLower(query)), -1)
	results := []string{}
	for _, idx := range idxs {

		// Ignore if not in range
		if !s.limits.contains(idx) {
			continue
		}

		results = append(results, s.CompleteWorks[idx-250:idx+250])
	}
	return results
}

func (b *bound) contains(x int) bool {
	return b.start <= x && x < b.end
}

// SetLines sets the limits of this searcher
// so they start and end at specific lines,
// in order to restrict the search
func (s *Searcher) SetLines(start, end int) {
	s.limits.start = start
	s.limits.end = end

	s.limits.convert(s.lineIndexes)
}
