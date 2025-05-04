package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/kenta-kenta/gaishi-sserverside/migrations"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq" // PostgreSQLの場合
)

// DB接続情報
const (
	dbUser     = "user"
	dbPassword = "postgres"
	dbName     = "memo_db"

	dsn = "user=user password=postgres dbname=memo_db sslmode=disable"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello, World")
}

func createMemo(w http.ResponseWriter, r *http.Request) {
	var memo struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var responseMemo struct {
		ID        int64  `json:"id"`
		Title     string `json:"title"`
		Content   string `json:"content"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&memo); err != nil {
		fmt.Fprintf(w, "decode err: %v", err)
	}
	// データベースに保存
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "db open err: %v", err)
		return
	}
	const sqlStr = "INSERT INTO memos(title, content) VALUES ($1, $2) RETURNING id, title, content, created_at, updated_at"
	err = db.QueryRow(sqlStr, memo.Title, memo.Content).Scan(&responseMemo.ID, &responseMemo.Title, &responseMemo.Content, &responseMemo.CreatedAt, &responseMemo.UpdatedAt)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "db exec err: %v", err)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(responseMemo); err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "encode err: %v", err)
		return
	}
}

func getMemo(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得
	path := r.URL.Path
	id := strings.Split(path, "/")[3]
	fmt.Println(path)
	fmt.Println(id)

	var memo struct {
		ID        int64  `json:"id"`
		Title     string `json:"title"`
		Content   string `json:"content"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "db open err: %v", err)
		return
	}
	const sqlStr = "SELECT * FROM memos WHERE id = $1"
	err = db.QueryRow(sqlStr, id).Scan(&memo.ID, &memo.Title, &memo.Content, &memo.CreatedAt, &memo.UpdatedAt)
	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "memo not found")
		return
	} else if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "db exec err: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(memo); err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
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
	// データベースのマイグレーション
	migrations.Migrate()

	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)
	mux.HandleFunc("POST /api/memos", createMemo)
	mux.HandleFunc("GET /api/memos/", getMemo)
	mux.HandleFunc("PUT /api/memos/", updateMemo)
	mux.HandleFunc("DELETE /api/memos/", deleteMemo)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
