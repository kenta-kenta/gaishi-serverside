
こちらをGoの標準パッケージのみを利用して作成していきます。
理解を重視するのではじめはあまりきれいなコードではありませんが、後程リファクタリングしていきます。
この記事では流れを理解することをゴールとしています。
https://gaishishukatsu.com/engineer/drill/2/1

Go言語は標準パッケージが強力で、ほとんどの基本的なことは標準パッケージで十分です。

## 要求項目

以下の内容を満たすようにAPIを作成していきます。いわゆるCRUD操作というものです。

### メモの作成

新しいメモを作成するAPIを実装します。

* HTTPメソッド: POST
* エンドポイント: /api/memos
* リクエストボディ:
{ "title": "メモのタイトル", "content": "メモの内容" }
* レスポンスステータス: 201 Created
* レスポンスボディ:
{ "id": 1, "title": "メモのタイトル", "content": "メモの内容", "createdAt": "現在の時刻", "updatedAt": "現在の時刻" }

### メモの取得

指定されたIDのメモを取得するAPIを実装します。

* HTTPメソッド: GET
* エンドポイント: /api/memos/:id
* レスポンスステータス: 200 OK
* レスポンスボディ:
{ "id": 1, "title": "メモのタイトル", "content": "メモの内容", "createdAt": "作成した時刻", "updatedAt": "更新した時刻" }
* 404 Not Found: 存在しないIDが指定された場合、404 Not Foundを返してください。

### メモの更新

既存のメモを更新するAPIを実装します。

* HTTPメソッド: PUT
* エンドポイント: /api/memos/:id
* リクエストボディ:
{ "title": "更新されたメモのタイトル", "content": "更新されたメモの内容" }
* レスポンスステータス: 200 OK
* レスポンスボディ:
{ "id": 1, "title": "更新されたメモのタイトル", "content": "更新されたメモの内容", "createdAt": "作成した時刻", "updatedAt": "現在の時刻" }
* 404 Not Found: 存在しないIDが指定された場合、404 Not Foundを返してください。

### メモの削除

指定されたIDのメモを削除するAPIを実装します。

* HTTPメソッド: DELETE
* エンドポイント: /api/memos/:id
* レスポンスステータス: 204 No Content
* 404 Not Found: 存在しないIDが指定された場合、404 Not Foundを返してください。

## サーバーの立て方
早速ですが、まずサーバーの立て方です。Go言語は公式のパッケージを利用するだけでサーバーを立てることができます。
```golang:main.go
package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello, World")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```
Goのファイルを実行するためには以下のコマンドを入力します。どちらでも大丈夫です。
```terminal
$ go run main.go
$ go run .
```

## 各ハンドラ関数の作成
次に、各ハンドラ関数を作成していきます。
まず、コードを紹介します。後程説明します。
```golang:main.go
package main

import (
	"encoding/json"
	"fmt"
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
	path := r.URL.Path
	id := strings.Split(path, "/")[3]
	fmt.Println(path)
	fmt.Println(id)

	type memoType struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	memo := memoType{
		Title:   "Sample Title",
		Content: "Sample Content",
	}

	fmt.Fprintf(w, "id: %s, title: %s, content: %s", id, memo.Title, memo.Content)
}

func updateMemo(w http.ResponseWriter, r *http.Request) {
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
	fmt.Fprintf(w, "delete")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)
	mux.HandleFunc("POST /api/memos", createMemo)
	mux.HandleFunc("GET /api/memos/", getMemo)
	mux.HandleFunc("PUT /api/memos/", updateMemo)
	mux.HandleFunc("DELETE /memos/", deleteMemo)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

```

それぞれを簡単に解説します。

- createMemo関数：Requestされたjsonを構造体に変換する。
- getMemo関数：パスパラメータで指定されたパスを受け取る。構造体をjsonに変換する。
- updateMemo関数：パスパラメータで指定されたパスを受け取る。jsonを構造体に変換する。
- deleteMemo関数：パスパラメータで指定されたパスを受け取る。メッセージを表示。

## データベースとの接続
次に、データベースとの接続を行います。
今回はDockerでPostgreSQLを利用します。
```yaml:docker-compose.yml
version: '3'

services:
  db:
    image: postgres:14
    container_name: postgres_pta
    ports:
      - 5432:5432
    volumes:
      - db-store:/var/lib/postgresql/data
    environment:
      POSTGRES_USER: 'user'
      POSTGRES_PASSWORD: 'postgres'
      POSTGRES_DB: 'memo_db'
volumes:
  db-store:
```
マイグレーションを行います。
マイグレーションとは、テーブルを作ったり、カラムを追加したりする変更を、ファイルとして記録していく仕組みのことです。
```golang:migrations/migrate.go
package migrations

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func Migrate() {
	dbUser := "user"
	dbPassword := "postgres"
	dbName := "memo_db"

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	defer db.Close()
	// テーブルが存在しない場合は作成する
	sqlStr := `
	CREATE TABLE IF NOT EXISTS memos (
		id         INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		title      VARCHAR(255) NOT NULL,
		context    TEXT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`
	_, err = db.Exec(sqlStr)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
}
```
各ハンドラで利用する際に利用するために必要な設定を行っています。
```golang:handlers/handlers.go
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
```
## 各関数をデータベースに接続
関数をデータベースに接続するうえでコードを別ファイルに分割しました。
先にコードを見せます。後程簡単に説明を書きます。

全ての関数に共通する大まかな流れを説明します。
1. リクエストとレスポンスの構造体を定義する
2. リクエストをデコードする
3. データベースに接続する
4. SQLを実行する
5. レスポンスにjsonエンコードする

```golang:handlers/handlers.go
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

```

```golang:handlers/handlers.go
func GetMemo(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得（例：api/memo/{id}）
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
```

```golang:handlers/handlers.go
func UpdateMemo(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得（例：api/memo/{id}）
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
```

```golang:handlers/handlers.go
func DeleteMemo(w http.ResponseWriter, r *http.Request) {
	// パスパラメータを取得（例：api/memo/{id}）
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

	w.WriteHeader(http.StatusNoContent)
}
```

以上のコードで要求されていた事項は満たされました。

## 参考文献