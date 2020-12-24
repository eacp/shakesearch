package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"pulley.com/shakesearch/search"
)

func main() {
	searcher := search.Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Create the pseudo db
	db := newSG(searcher)

	// Properly identify the data
	db.identifyWorks(map[string][2]int{
		"sonnets":              {132, 2909},
		"alls-well":            {2910, 7872},
		"antony-and-cleopatra": {7873, 14513},
		"as-you-like-it":       {14514, 17306},
		"the-comedy-of-errors": {17306, 20506},
		"coriolanus":           {20507, 24623},
		"cymbeline":            {24624, 30498},
		"hamlet":               {30500, 37186},
		"henry-iv-1":           {37186, 41902},
		"henry-iv-2":           {41902, 45311},
		"henry-v":              {45311, 50245},
		"henry-vi-1":           {50246, 53518},
		// I could add more data here,
		// but it would get too long

		// These last 2 are classics

		"macbeth":          {80513, 84660},
		"romeo-and-juliet": {121875, 127131},
	})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", db.handleSearch)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

}

// A search group is a pseudo db containing multiple works

type searcherGroup struct {
	defaultSearcher search.Searcher
	searchers       map[string]search.Searcher
}

func (sg *searcherGroup) searchIn(work, query string) ([]string, bool) {
	// If no parameter is passed or we look in all works, use
	// the default searcher
	if len(work) == 0 || work == "all" {
		return sg.defaultSearcher.Search(query), true
	}

	// Look the work in the map
	s, ok := sg.searchers[work]

	// If not ok then the work is not found
	if !ok {
		log.Printf("Requested work '%v' does not exist", work)
		return nil, false
	}

	// The work is found: execute the search in that work
	return s.Search(query), true
}

func (sg *searcherGroup) handleSearch(w http.ResponseWriter, r *http.Request) {
	// Get the work and the query
	work := r.FormValue("work")
	query := r.FormValue("q")

	log.Printf("Received query=%v work=%v", query, work)

	// the query is mandatory
	if len(query) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing search query in URL params"))
		return
	}

	results, found := sg.searchIn(work, query)

	// If not found then 404
	if !found {
		http.NotFound(w, r)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(results); err != nil {
		// Log what happened
		log.Println(err)
		// Return internal server error
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "\"encoding failure\"")
	}

}

func newSG(s search.Searcher) (sg searcherGroup) {
	// Make this the default searcher
	sg.defaultSearcher = s

	// Allocate the default map
	sg.searchers = make(map[string]search.Searcher, 100)
	return
}

func (sg *searcherGroup) identifyWorks(workBounds map[string][2]int) {
	for title, bounds := range workBounds {
		// Copy the default searcher
		workSearcher := sg.defaultSearcher
		// Set the lines
		workSearcher.SetLines(bounds[0], bounds[1])
		// Register
		sg.searchers[title] = workSearcher
	}
}
