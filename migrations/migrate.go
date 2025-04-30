package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Migrate() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	defer db.Close()
	// テーブルが存在しない場合は作成する
	sqlStr := `
	CREATE TABLE IF NOT EXISTS memos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT content,
    NULL TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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
