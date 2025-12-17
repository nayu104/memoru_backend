package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// %wとは。。。
func New() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")

	 // 環境変数が空なら、そもそも接続先が分からないのでエラーにする
    if dsn == "" {
        return nil, fmt.Errorf("DATABASE_URL is not set")
    }

	conn, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, fmt.Errorf("sql.Open failed: %w", err)
	}

	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("db.Ping failed: %w", err)
	}
	return conn, nil
}