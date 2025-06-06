package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Quote struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

var (
	mu     sync.Mutex
	quotes = make([]Quote, 0)
	nextID = 1
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/quotes", quotesHandler)
	http.HandleFunc("/quotes/", quoteByIDHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func quotesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleAddQuote(w, r)
	case http.MethodGet:
		handleGetQuotes(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func quoteByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/quotes/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		handleDeleteQuote(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleAddQuote(w http.ResponseWriter, r *http.Request) {
	var q Quote
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if q.Author == "" || q.Quote == "" {
		http.Error(w, "Missing author or quote field", http.StatusBadRequest)
		return
	}

	mu.Lock()
	q.ID = nextID
	nextID++
	quotes = append(quotes, q)
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(q)
}

func handleGetQuotes(w http.ResponseWriter, r *http.Request) {
	authorFilter := r.URL.Query().Get("author")

	mu.Lock()
	defer mu.Unlock()

	if authorFilter == "random" {

		if len(quotes) == 0 {
			http.Error(w, "No quotes found", http.StatusNotFound)
			return
		}
		quote := quotes[rand.Intn(len(quotes))]
		json.NewEncoder(w).Encode(quote)
		return
	}

	if authorFilter != "" {

		filtered := make([]Quote, 0)
		for _, q := range quotes {
			if strings.EqualFold(q.Author, authorFilter) {
				filtered = append(filtered, q)
			}
		}
		json.NewEncoder(w).Encode(filtered)
		return
	}

	json.NewEncoder(w).Encode(quotes)
}

func handleDeleteQuote(w http.ResponseWriter, id int) {
	mu.Lock()
	defer mu.Unlock()

	for i, q := range quotes {
		if q.ID == id {
			quotes = append(quotes[:i], quotes[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Quote not found", http.StatusNotFound)
}
