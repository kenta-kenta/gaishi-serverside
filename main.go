package main

import (
	"log"
	"net/http"

	"github.com/kenta-kenta/gaishi-serverside/handlers"
	"github.com/kenta-kenta/gaishi-serverside/migrations"
	_ "github.com/lib/pq" // PostgreSQLの場合
)

func main() {
	// データベースのマイグレーション
	migrations.Migrate()

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Hello)
	mux.HandleFunc("POST /api/memos", handlers.CreateMemo)
	mux.HandleFunc("GET /api/memos/", handlers.GetMemoByID)
	mux.HandleFunc("PUT /api/memos/", handlers.UpdateMemo)
	mux.HandleFunc("DELETE /api/memos/", handlers.DeleteMemo)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
