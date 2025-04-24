package migrations

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Migrate() {
	dsn := "user:password@tcp(localhost:3307)/memo_db?parseTime=true"
	db, err := sql.Open("mysql", dsn)
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
