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
		content    TEXT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`
	_, err = db.Exec(sqlStr)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
}

func Drop() {
	dsn := "user:password@tcp(localhost:3307)/memo_db?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	defer db.Close()
	// テーブルが存在しない場合は作成する
	sqlStr := `DROP TABLE IF EXISTS memos`
	_, err = db.Exec(sqlStr)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
}
