package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestShortenAndRedirect(t *testing.T) {
	os.Remove("./test.db") // Clean previous test DB
	var err error
	db, err = sql.Open("sqlite3", "./test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	schema, _ := os.ReadFile("db/schema.sql")
	_, err = db.Exec(string(schema))
	if err != nil {
		t.Fatal(err)
	}

	// Step 1: Shorten a URL
	payload := []byte(`{"original_url":"https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	shortenHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp URLMapping
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	// Step 2: Redirect
	req2 := httptest.NewRequest(http.MethodGet, "/r/"+resp.ShortCode, nil)
	w2 := httptest.NewRecorder()
	redirectHandler(w2, req2)

	if w2.Code != http.StatusFound {
		t.Fatalf("Expected redirect 302, got %d", w2.Code)
	}
}
