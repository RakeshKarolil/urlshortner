package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"urlshortener/metrics"

	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// URLMapping represents a short URL entry
type URLMapping struct {
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./urlshortener.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create schema if not exists
	schema, _ := os.ReadFile("db/schema.sql")
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalf("Error creating schema: %v", err)
	}

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/r/", redirectHandler)

	// Expose Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// shortenHandler handles URL shortening requests
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST method required", http.StatusMethodNotAllowed)
		return
	}

	var req URLMapping
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	shortCode := generateShortCode()
	_, err := db.Exec("INSERT INTO urls (short_code, original_url) VALUES (?, ?)", shortCode, req.OriginalURL)
	if err != nil {
		http.Error(w, "Failed to store URL", http.StatusInternalServerError)
		return
	}

	metrics.URLsShortened.Inc()

	resp := URLMapping{ShortCode: shortCode, OriginalURL: req.OriginalURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// redirectHandler handles redirection from short URLs
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/r/"):]
	var originalURL string
	err := db.QueryRow("SELECT original_url FROM urls WHERE short_code = ?", code).Scan(&originalURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	metrics.Redirects.Inc()
	http.Redirect(w, r, originalURL, http.StatusFound)
}

// generateShortCode creates a random short code
func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
