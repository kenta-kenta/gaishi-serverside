package repositories

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/kenta-kenta/gaishi-serverside/models"
)

const (
	dbUser     = "user"
	dbPassword = "postgres"
	dbName     = "memo_db"
	dsn        = "user=user password=postgres dbname=memo_db sslmode=disable"
)

func getDB() (*sql.DB, error) {
	return sql.Open("postgres", dsn)
}

func CreateMemo(memo models.CreateMemoRequest) (*models.Memo, error) {
	db, err := getDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get db connection: %w", err)
	}
	defer db.Close()

	const sqlStr = "INSERT INTO memos(title, content) VALUES ($1, $2) RETURNING id, title, content, created_at, updated_at"

	var result models.Memo
	err = db.QueryRow(sqlStr, memo.Title, memo.Content).Scan(
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

func GetMemoByID(id int) (*models.Memo, error) {
	db, err := getDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get db connection: %w", err)
	}
	defer db.Close()

	const sqlStr = "SELECT id, title, content, created_at, updated_at FROM memos WHERE id = $1"

	var memo models.Memo
	err = db.QueryRow(sqlStr, id).Scan(
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

func UpdateMemo(id int, memo models.UpdateMemoRequest) (*models.Memo, error) {
	db, err := getDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get db connection: %w", err)
	}
	defer db.Close()

	const sqlStr = "UPDATE memos SET title = $1, content = $2 WHERE id = $3 RETURNING id, title, content, created_at, updated_at"

	var result models.Memo
	err = db.QueryRow(sqlStr, memo.Title, memo.Content, id).Scan(
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

func DeleteMemo(id int) error {
	db, err := getDB()
	if err != nil {
		return fmt.Errorf("failed to get db connection: %w", err)
	}
	defer db.Close()

	const sqlStr = "DELETE FROM memos WHERE id = $1"

	result, err := db.Exec(sqlStr, id)
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
