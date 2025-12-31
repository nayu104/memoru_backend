package repository

import (
    "context"
    "database/sql"
    "fmt"

    "emotion-memo-api/internal/domain"
)

type PostgresMemoRepository struct {
    db *sql.DB
}

func NewPostgresMemoRepository(db *sql.DB) *PostgresMemoRepository {
    return &PostgresMemoRepository{db: db}
}

func (repo *PostgresMemoRepository) Create(ctx context.Context, memo *domain.Memo) error {
    const query = `
    INSERT INTO "memo" (user_id, body, mood, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING memo_id
    `

    if err := repo.db.QueryRowContext(ctx, query,
        memo.UserID,
        memo.Body,
        memo.Mood,
        memo.CreatedAt,
        memo.UpdatedAt,
    ).Scan(&memo.ID); err != nil {
        return fmt.Errorf("❌メモの作成に失敗しました: %w", err)
    }
    return nil
}

// ListByUserID: 指定されたユーザーIDのメモを取得する
func (r *PostgresMemoRepository) ListByUserID(ctx context.Context, userID string) ([]domain.Memo, error) {
    const query = `
    SELECT memo_id, user_id, body, mood, created_at, updated_at
    FROM "memo"
    WHERE user_id = $1
    ORDER BY created_at DESC
    `

    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("❌エラーが発生しました: %w", err)
    }
    defer rows.Close()

	// メモリ上にメモのリストを作る.0は要素が入っていないという。
    memos := make([]domain.Memo, 0)

    for rows.Next() {
        var memo domain.Memo
        // Scan は、「データベースから取ってきた値を、Goの変数にコピーする機能」
		// この err は、エラーチェック以外では使わないから、if文の中に閉じ込めてしまおう」 という意図がある。
        if err := rows.Scan(
            &memo.ID,
            &memo.UserID,
            &memo.Body,
            &memo.Mood,
            &memo.CreatedAt,
            &memo.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("❌データの読み込みに失敗しました: %w", err)
        }
        memos = append(memos, memo)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("❌データ取得中にエラーが発生しました: %w", err)
    }

    return memos, nil
}