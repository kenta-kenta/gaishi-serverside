package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello, World")
}

func createMemo(w http.ResponseWriter, r *http.Request) {
	var memo struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&memo); err != nil {
		fmt.Fprintf(w, "decode err: %v", err)
	}
	fmt.Fprintf(w, "title: %s, content: %s", memo.Title, memo.Content)
}

func getMemo(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得
	path := r.URL.Path
	id := strings.Split(path, "/")[3]
	fmt.Println(path)
	fmt.Println(id)

	memo := struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}{
		Title:   "Sample Title",
		Content: "Sample Content",
	}

	if err := json.NewEncoder(w).Encode(memo); err != nil {
		fmt.Fprintf(w, "encode err: %v", err)
		return
	}
}

func updateMemo(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得
	path := r.URL.Path
	id := strings.Split(path, "/")[3]
	fmt.Println(path)
	fmt.Println(id)

	var memo struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&memo); err != nil {
		fmt.Fprintf(w, "decode err: %v", err)
	}
	fmt.Fprintf(w, "title: %s, content: %s", memo.Title, memo.Content)
}

func deleteMemo(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得
	path := r.URL.Path
	id := strings.Split(path, "/")[3]
	fmt.Println(path)
	fmt.Println(id)

	fmt.Fprintf(w, "delete")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)
	mux.HandleFunc("POST /api/memos", createMemo)
	mux.HandleFunc("GET /api/memos/", getMemo)
	mux.HandleFunc("PUT /api/memos/", updateMemo)
	mux.HandleFunc("DELETE /api/memos/", deleteMemo)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
