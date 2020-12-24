package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"pulley.com/shakesearch/search"
)

func Test_searcherGroup_searchIn(t *testing.T) {
	/*
		For this test we will use the text in complete works and
		2 specific works: One for Romeo & Juliet
		and another for hamlet
	*/

	s := search.Searcher{}
	if err := s.Load("completeworks.txt"); err != nil {
		t.Fatalf("Could not load text from file: %v", err)
	}

	// Make the db from that searcher
	db := newSG(s)

	// Use all the works
	t.Run("Global search in ALL the works", func(t *testing.T) {

		tests := []searchTestCase{
			{"Romeo", "all", 317, true},
			{"rOmEo", "all", 317, true},

			{"Juliet", "all", 210, true},
			{"JULIET", "", 210, true},
			{"jUlIeT", "all", 210, true},

			{"MESSENGER", "", 395, true},
			{"mEsSeNgEr", "", 395, true},

			// 'loooooooool' is obviously not part
			// of the works, but should give ok becuse
			// 'all' is allowed
			{"loooooooool", "all", 0, true},
			{"Shakira", "all", 0, true},

			// Find in works that DO NOT exist
			{"Juliet", "The Avengers", 0, false},
			{"Hamlet", "The Avengers", 0, false},
		}

		for _, tt := range tests {
			t.Run(tt.name(), func(t *testing.T) {
				// Run this test with the example db
				tt.runWithDB(&db, t)
			})
		}
	})

	// Use only Romeo and Juliet

	// Register the searchers
	rnj := s
	rnj.SetLines(121875, 127131)

	db.searchers["rnj"] = rnj

	h := s
	h.SetLines(30500, 37186)

	db.searchers["hamlet"] = h

	t.Run("Search only in Romeo and Juliet", func(t *testing.T) {

		tests := []searchTestCase{
			{"Romeo", "rnj", 316, true},
			{"rOmEo", "rnj", 316, true},

			// There are about 20 juliets before in
			// the corpus
			{"Juliet", "rnj", 210 - 20, true},
			{"JULIET", "rnj", 210 - 20, true},
			{"jUlIeT", "rnj", 210 - 20, true},

			// Hamlet should NOT be in this work
			{"Hamlet", "rnj", 0, true},
			{"hAmLeT", "rnj", 0, true},

			// Nor should

			// 'loooooooool' is obviously not part
			// of the works, but should give ok becuse
			// 'all' is allowed
			{"loooooooool", "rnj", 0, true},
			{"Shakira", "rnj", 0, true},

			// Find in works that DO NOT exist
			{"Juliet", "The Avengers", 0, false},
			{"Hamlet", "The Avengers", 0, false},
		}

		// I should somehow specify in the db object
		// where do certain works beguin and end

		for _, tt := range tests {
			t.Run(tt.name(), func(t *testing.T) {
				// Run this test with the example db
				tt.runWithDB(&db, t)
			})
		}
	})

	t.Run("Search only in Hamlet", func(t *testing.T) {

		tests := []searchTestCase{
			{"Romeo", "hamlet", 0, true},
			{"rOmEo", "hamlet", 0, true},

			// Hamlet should NOT be in this work
			{"Hamlet", "hamlet", 474, true},
			{"hAmLeT", "hamlet", 474, true},
		}

		// I should somehow specify in the db object
		// where do certain works beguin and end

		for _, tt := range tests {
			t.Run(tt.name(), func(t *testing.T) {
				// Run this test with the example db
				tt.runWithDB(&db, t)
			})
		}
	})
}

type searchTestCase struct {
	query, work         string
	expectedOccurrences int
	expectedOk          bool
}

// Helper for the name
func (tt *searchTestCase) name() string {
	return fmt.Sprintf("Find '%v' in the work '%v'", tt.query, tt.work)
}

// A helper to avoid indentation
func (tt *searchTestCase) runWithDB(db *searcherGroup, t *testing.T) {
	data, ok := db.searchIn(tt.work, tt.query)

	// Check the OK variable
	if ok != tt.expectedOk {
		t.Fatalf("Expected ok = %t, got %t", tt.expectedOk, ok)
	}

	occurrences := len(data)

	if occurrences != tt.expectedOccurrences {
		t.Fatalf("Expected occurences = %d, got %d", tt.expectedOccurrences, occurrences)
	}

}

// Another helper to avoid indentation
func (tt *searchTestCase) runHTTP(db *searcherGroup, t *testing.T) {
	// Make a test request and its response writer
	uri := fmt.Sprintf("/search?q=%s&work=%s", tt.query, tt.work)
	req := httptest.NewRequest(http.MethodGet, uri, nil)

	w := httptest.NewRecorder()

	// Make the 'request'

	db.handleSearch(w, req)

}

func Test_searcherGroup_handleSearch(t *testing.T) {
	// Make a pseudo db with 2 works properly identified
	s := search.Searcher{}
	if err := s.Load("completeworks.txt"); err != nil {
		t.Fatalf("Could not load text from file: %v", err)
	}

	// Make the db from that searcher
	db := newSG(s)

	// Identify 2 works
	// Register the searchers
	rnj := s
	rnj.SetLines(121875, 127131)

	db.searchers["rnj"] = rnj

	h := s
	h.SetLines(30500, 37186)

	db.searchers["hamlet"] = h

	// Test bad or malformed urls
	t.Run("Bad URL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/search", nil)
		w := httptest.NewRecorder()

		db.handleSearch(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	// Test looking in works that do not exist -> 404
	t.Run("Test 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/search?q=Juliet&work=TheAvengers", nil)
		w := httptest.NewRecorder()

		db.handleSearch(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	// Create some samle queries
	tests := []searchTestCase{
		{"Romeo", "all", 317, true},
		{"rOmEo", "all", 317, true},

		{"Juliet", "all", 210, true},
		{"JULIET", "", 210, true},
		{"jUlIeT", "all", 210, true},

		{"MESSENGER", "", 395, true},
		{"mEsSeNgEr", "", 395, true},

		// 'loooooooool' is obviously not part
		// of the works, but should give ok becuse
		// 'all' is allowed
		{"loooooooool", "all", 0, true},
		{"Shakira", "all", 0, true},
	}

	for _, tt := range tests {
		tt.runHTTP(&db, t)
	}

}

func Test_searcherGroup_identifyWorks(t *testing.T) {
	// Make a pseudo db with 2 works properly identified
	s := search.Searcher{}
	if err := s.Load("completeworks.txt"); err != nil {
		t.Fatalf("Could not load text from file: %v", err)
	}

	// Make the db from that searcher
	db := newSG(s)

	// Identify Rome and Juliet, and hamlet
	db.identifyWorks(map[string][2]int{
		"hamlet": {30500, 37186},
		"rj":     {121875, 127131},
	})

	t.Run("Test Romeo And Juliet", func(t *testing.T) {
		_, ok := db.searchers["rj"]

		if !ok {
			t.Error("Romeo and Juliet not registered")
		}
	})

	t.Run("Test Hamlet", func(t *testing.T) {
		_, ok := db.searchers["hamlet"]

		if !ok {
			t.Error("Romeo and Juliet not registered")
		}
	})

}
