package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQLの場合
	"net/http"
	"strings"
)

// DB接続情報
const (
	dbUser     = "user"
	dbPassword = "postgres"
	dbName     = "memo_db"
	dsn        = "user=user password=postgres dbname=memo_db sslmode=disable"
)

// getDB はデータベース接続を生成するヘルパー関数です
func getDB() (*sql.DB, error) {
	return sql.Open("postgres", dsn)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello, World")
}

func CreateMemo(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, fmt.Sprintf("decode err: %v", err), http.StatusBadRequest)
		return
	}

	db, err := getDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("db open err: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	const sqlStr = "INSERT INTO memos(title, content) VALUES ($1, $2) RETURNING id, title, content, created_at, updated_at"
	err = db.QueryRow(sqlStr, memo.Title, memo.Content).Scan(&responseMemo.ID, &responseMemo.Title, &responseMemo.Content, &responseMemo.CreatedAt, &responseMemo.UpdatedAt)
	if err != nil {
		http.Error(w, fmt.Sprintf("db exec err: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(responseMemo); err != nil {
		http.Error(w, fmt.Sprintf("encode err: %v", err), http.StatusInternalServerError)
		return
	}
}

func GetMemo(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得（例：/memo/{id}）
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}
	id := parts[3]

	var memo struct {
		ID        int64  `json:"id"`
		Title     string `json:"title"`
		Content   string `json:"content"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}

	db, err := getDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("db open err: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	const sqlStr = "SELECT id, title, content, created_at, updated_at FROM memos WHERE id = $1"
	err = db.QueryRow(sqlStr, id).Scan(&memo.ID, &memo.Title, &memo.Content, &memo.CreatedAt, &memo.UpdatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "memo not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("db exec err: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(memo); err != nil {
		http.Error(w, fmt.Sprintf("encode err: %v", err), http.StatusInternalServerError)
		return
	}
}

func UpdateMemo(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得（例：/memo/{id}）
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}
	id := parts[3]

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
		http.Error(w, fmt.Sprintf("decode err: %v", err), http.StatusBadRequest)
		return
	}

	db, err := getDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("db open err: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	const sqlStr = "UPDATE memos SET title = $1, content = $2 WHERE id = $3 RETURNING id, title, content, created_at, updated_at"
	err = db.QueryRow(sqlStr, memo.Title, memo.Content, id).Scan(&responseMemo.ID, &responseMemo.Title, &responseMemo.Content, &responseMemo.CreatedAt, &responseMemo.UpdatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "memo not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("db exec err: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(responseMemo); err != nil {
		http.Error(w, fmt.Sprintf("encode err: %v", err), http.StatusInternalServerError)
		return
	}
}

func DeleteMemo(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得（例：/memo/{id}）
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}
	id := parts[3]

	db, err := getDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("db open err: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	const sqlStr = "DELETE FROM memos WHERE id = $1"
	result, err := db.Exec(sqlStr, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("db exec err: %v", err), http.StatusInternalServerError)
		return
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "memo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}
