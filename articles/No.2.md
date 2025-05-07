## 現在の問題点

1. データベース関連の設定とハンドラーが混在している
2. 共通の構造体や定数が散在している
3. 各ハンドラーが独立しているが、同じファイルにまとめられている

## リファクタリングの目的

1. 関心の分離
2. 拡張性
3. テストのしやすさ
4. クリーンアーキテクチャへの移行

## 実際に修正した点

今回利用しているアーキテクチャは「3層アーキテクチャ」または「レイヤードアーキテクチャ」と呼ばれているものです。
具体的には以下の3層に分かれています。

1. プレゼンテーション層
* HTTPリクエストの処理とレスポンスの整形を担当
* ユーザーインターフェースとのやり取りを管理

2. ビジネスロジック層
* アプリケーションのビジネスロジックを実装
* データの検証や処理の制御を担当

3. データアクセス層
* データベースとのやり取りを担当
* データの永続化を管理

具体的に見ていきましょう
### モデル
```golang: models/memo.go
package models

type Memo struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CreateMemoRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateMemoRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
} 
```
モデルは、データ構造の定義を行います。
今回は、メモのデータ構造を定義しました。

### データアクセス層
```golang: repositories/memo/memo.go
package memo

import (
	"database/sql"
	"fmt"

	"gaishi-serverside/db"
	"gaishi-serverside/models"
)

// Repository はメモのデータアクセスを担当する構造体です
type Repository struct {
	db *sql.DB
}

// NewRepository は新しいRepositoryインスタンスを作成します
func NewRepository() (*Repository, error) {
	db, err := db.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get db connection: %w", err)
	}
	return &Repository{db: db}, nil
}

func (r *Repository) CreateMemo(memo models.CreateMemoRequest) (*models.Memo, error) {
	const sqlStr = "INSERT INTO memos(title, content) VALUES ($1, $2) RETURNING id, title, content, created_at, updated_at"
	
	var result models.Memo
	err := r.db.QueryRow(sqlStr, memo.Title, memo.Content).Scan(
		&result.ID,
		&result.Title,
		&result.Content,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create memo: %w", err)
	}
	
	return &result, nil
}

func (r *Repository) GetMemoByID(id int64) (*models.Memo, error) {
	const sqlStr = "SELECT id, title, content, created_at, updated_at FROM memos WHERE id = $1"
	
	var memo models.Memo
	err := r.db.QueryRow(sqlStr, id).Scan(
		&memo.ID,
		&memo.Title,
		&memo.Content,
		&memo.CreatedAt,
		&memo.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get memo: %w", err)
	}
	
	return &memo, nil
}

func (r *Repository) UpdateMemo(id int64, memo models.UpdateMemoRequest) (*models.Memo, error) {
	const sqlStr = "UPDATE memos SET title = $1, content = $2 WHERE id = $3 RETURNING id, title, content, created_at, updated_at"
	
	var result models.Memo
	err := r.db.QueryRow(sqlStr, memo.Title, memo.Content, id).Scan(
		&result.ID,
		&result.Title,
		&result.Content,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update memo: %w", err)
	}
	
	return &result, nil
}

func (r *Repository) DeleteMemo(id int64) error {
	const sqlStr = "DELETE FROM memos WHERE id = $1"
	
	result, err := r.db.Exec(sqlStr, id)
	if err != nil {
		return fmt.Errorf("failed to delete memo: %w", err)
	}
	
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affected == 0 {
		return nil
	}
	
	return nil
} 
```
データアクセス層は、データベースとのやり取りを担当します。
前回作成したものから増えたものを説明します。
それが以下のコードです。
```golang: 
// MemoHandler はメモ関連のHTTPハンドラーを管理する構造体です
type MemoHandler struct {
	service *memo.Service
}

// NewMemoHandler は新しいMemoHandlerインスタンスを作成します
func NewMemoHandler() (*MemoHandler, error) {
	service, err := memo.NewService()
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}
	return &MemoHandler{service: service}, nil
}
```
これは、構造体を作成するためのコードです。
このコードを追加した理由は、依存関係がわかりやすいことや、
テストのしやすさを考慮したためです。

### ビジネスロジック層
```golang: services/memo.go
package services

import (
	"fmt"

	"gaishi-serverside/models"
	"gaishi-serverside/repositories/memo"
)

type Service struct {
	repo *memo.Repository
}

func NewService() (*Service, error) {
	repo, err := memo.NewRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}
	return &Service{repo: repo}, nil
}

func (s *Service) CreateMemo(req models.CreateMemoRequest) (*models.Memo, error) {
	return s.repo.Create(req)
}

func (s *Service) GetMemo(id int64) (*models.Memo, error) {
	return s.repo.GetByID(id)
}

func (s *Service) UpdateMemo(id int64, req models.UpdateMemoRequest) (*models.Memo, error) {
	return s.repo.Update(id, req)
}

func (s *Service) DeleteMemo(id int64) error {
	return s.repo.Delete(id)
} 
```
これは、ビジネスロジックを実装するためのコードです。
前回作成したものから増えたものは先ほど説明したので省略します。

ビジネスロジック層が担うのは、データの検証や処理の制御です。
つまり、今回は特に行っていることはありません。

### プレゼンテーション層
```golang: handlers/memo.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq" 
	"net/http"
	"strconv"
	"strings"

	"gaishi-serverside/models"
	"gaishi-serverside/services/memo"
)

const (
	dbUser     = "user"
	dbPassword = "postgres"
	dbName     = "memo_db"
	dsn        = "user=user password=postgres dbname=memo_db sslmode=disable"
)

// MemoHandler はメモ関連のHTTPハンドラーを管理する構造体です
type MemoHandler struct {
	service *memo.Service
}

// NewMemoHandler は新しいMemoHandlerインスタンスを作成します
func NewMemoHandler() (*MemoHandler, error) {
	service, err := memo.NewService()
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}
	return &MemoHandler{service: service}, nil
}

// getDB はデータベース接続を生成するヘルパー関数です
func getDB() (*sql.DB, error) {
	return sql.Open("postgres", dsn)
}

// Hello はテスト用のハンドラーです
func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello, World")
}

func (h *MemoHandler) CreateMemo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("decode err: %v", err), http.StatusBadRequest)
		return
	}

	memo, err := h.service.CreateMemo(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("service err: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(memo); err != nil {
		http.Error(w, fmt.Sprintf("encode err: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *MemoHandler) GetMemo(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得（例：/memo/{id}）
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid memo ID", http.StatusBadRequest)
		return
	}

	memo, err := h.service.GetMemo(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("service err: %v", err), http.StatusInternalServerError)
		return
	}
	if memo == nil {
		http.Error(w, "memo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(memo); err != nil {
		http.Error(w, fmt.Sprintf("encode err: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *MemoHandler) UpdateMemo(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得（例：/memo/{id}）
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid memo ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateMemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("decode err: %v", err), http.StatusBadRequest)
		return
	}

	memo, err := h.service.UpdateMemo(id, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("service err: %v", err), http.StatusInternalServerError)
		return
	}
	if memo == nil {
		http.Error(w, "memo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(memo); err != nil {
		http.Error(w, fmt.Sprintf("encode err: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *MemoHandler) DeleteMemo(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得（例：/memo/{id}）
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid memo ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteMemo(id); err != nil {
		http.Error(w, fmt.Sprintf("service err: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
```
プレゼンテーション層では、HTTPリクエストの処理とレスポンスの整形を担当します。



## まとめ

今回は、全ての処理が同じファイルにあり、可読性が低いという問題がありました。
そのため、関心の分離を行い、可読性を高めるためにリファクタリングを行いました。

## 次回

次回は、JWTを利用して認証を行います。

