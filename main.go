package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	// PostgreSQLドライバ（これがないとDBに繋がりません）
	_ "github.com/lib/pq"

	"emotion-memo-api/internal/handler"
	"emotion-memo-api/internal/repository"
	memo_service "emotion-memo-api/internal/service"
)

func main() {
	// 1. データベースに接続する
	// ⚠️ "password" の部分は、あなたのPostgreSQLのパスワードに書き換えてください
	connStr := "postgres://postgres:password@localhost:5432/emotion_memo?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ DB接続設定のエラー:", err)
	}
	defer db.Close()

	// 実際に繋がるかテスト
	if err := db.Ping(); err != nil {
		log.Fatal("❌ DBに繋がりません:", err)
	}
	fmt.Println("✅ DB接続成功")

	// 2. 部品を組み立てる（依存性の注入）
	// DB担当 (Repository) を作る
	memoRepo := repository.NewPostgresMemoRepository(db)

	// ルール担当 (Service) を作る（DB担当を渡す）
	memoService := memo_service.NewMemoService(memoRepo)

	// 窓口担当 (Handler) を作る（ルール担当を渡す）
	memoHandler := handler.NewMemoHandler(memoService)

	// 3. URLと処理を紐付ける（ルーティング）
	mux := http.NewServeMux()

	// POST /memos -> メモ作成
	mux.HandleFunc("POST /memos", memoHandler.CreateMemo)

	// GET /memos -> メモ一覧
	mux.HandleFunc("GET /memos", memoHandler.ListMemos)

	// 4. サーバー起動
	fmt.Println("🚀 サーバー起動: http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("❌ サーバー終了:", err)
	}
}
