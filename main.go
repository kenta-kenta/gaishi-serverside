package main

import (
	"github.com/kenta-kenta/gaishi-sserverside/handlers"
	"github.com/kenta-kenta/gaishi-sserverside/migrations"
	_ "github.com/lib/pq" // PostgreSQLの場合
	"log"
	"net/http"
)

func main() {
	// データベースのマイグレーション
	migrations.Migrate()

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Hello)
	mux.HandleFunc("POST /api/memos", handlers.CreateMemo)
	mux.HandleFunc("GET /api/memos/", handlers.GetMemo)
	mux.HandleFunc("PUT /api/memos/", handlers.UpdateMemo)
	mux.HandleFunc("DELETE /api/memos/", handlers.DeleteMemo)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
