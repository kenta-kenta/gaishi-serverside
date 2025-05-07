package main

import (
	"github.com/kenta-kenta/gaishi-serverside/handlers"
	"github.com/kenta-kenta/gaishi-serverside/migrations"
	"github.com/kenta-kenta/gaishi-serverside/repositories"
	_ "github.com/lib/pq" // PostgreSQLの場合
	"log"
	"net/http"
)

func main() {
	// データベースのマイグレーション
	migrations.Migrate()

	// データベース接続の初期化
	repo, err := repositories.NewRepository()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Hello)
	mux.HandleFunc("POST /api/memos", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateMemo(w, r, repo)
	})
	mux.HandleFunc("GET /api/memos/", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetMemoByID(w, r, repo)
	})
	mux.HandleFunc("PUT /api/memos/", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateMemo(w, r, repo)
	})
	mux.HandleFunc("DELETE /api/memos/", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteMemo(w, r, repo)
	})

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
