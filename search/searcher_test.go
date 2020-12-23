package search

import (
	"fmt"
	"testing"
)

func Test_bound_contains(t *testing.T) {
	tests := []struct {
		name string
		b    bound
		x    int
		want bool
	}{
		{"Outside on the left", bound{1000, 1500}, 999, false},
		{"Outside on the right", bound{1000, 1500}, 1999, false},
		{"Exactly the start", bound{1000, 1500}, 1000, true},
		{"Just Inside", bound{1000, 1500}, 1250, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.contains(tt.x); got != tt.want {
				t.Errorf("bound.contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearcher_Load(t *testing.T) {

	s := Searcher{}

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"File that does exist", "../completeworks.txt", false},
		{"File that does NOT exist", "awawawawa.txt", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := s.Load(tt.filename); (err != nil) != tt.wantErr {
				t.Errorf("Searcher.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type searchTestCase struct {
	query               string
	expectedOccurrences int
	searcher            *Searcher
}

func (tt *searchTestCase) str() string {
	return fmt.Sprintf("Find '%v'", tt.query)
}

// Put it here because I'm getting a lot of
// indentation

func (tt *searchTestCase) run(t *testing.T) {
	// Get the matches
	data := tt.searcher.Search(tt.query)
	occurrences := len(data)

	// Check we receive the same matches
	if occurrences != tt.expectedOccurrences {
		t.Errorf("s.Search('%v') gave %v matches, expected %v",
			tt.query, occurrences, tt.expectedOccurrences)
	}
}

func TestSearcher_Search(t *testing.T) {
	// Make the base searcher
	// this one will search in ALL the works

	all := Searcher{}
	if err := all.Load("../completeworks.txt"); err != nil {
		// Something bad happened
		t.Errorf("Could not load file to test: %v", err)
	}

	// Execute sub tests

	t.Log("Searching in ALL the works. No sub searcher so far. This is to validate the search works in case insensitive mode")

	t.Run("Search in all", func(t *testing.T) {

		tests := []searchTestCase{
			// Romeo case insensitive
			{"Romeo", 317, &all},
			{"rOmEo", 317, &all},

			// Juliet case insensitive
			{"Juliet", 210, &all},
			{"JULIET", 210, &all},
			{"jUlIeT", 210, &all},

			// Another case insensitive test
			{"MESSENGER", 395, &all},
			{"mEsSeNgEr", 395, &all},
		}

		for _, tt := range tests {
			// Run each test case
			t.Run(tt.str(), tt.run)
		}
	})

	// Make sub searchers here
	hamlet := all
	hamlet.SetLines(30499, 37185)

	// Romeo aNd Juliet
	rnj := all
	rnj.SetLines(121875, 127131)

	t.Log("Now, we will proceed to look inside the sub searchers")

	t.Run("Search in sub searchers", func(t *testing.T) {
		tests := []searchTestCase{
			// Romeo is not present in hamlet, only in R&J
			{"Romeo", 316, &rnj}, // The index does not count
			{"Hamlet", 0, &rnj},

			{"Romeo", 0, &hamlet},
			{"Hamlet", 474, &hamlet},
		}

		for _, tt := range tests {
			// Run each test case
			t.Run(tt.str(), tt.run)
		}

	})

}
