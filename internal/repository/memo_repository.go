package repository

import (
    "context"

    "emotion-memo-api/internal/domain"
)

// MemoRepository: DB操作を抽象化（serviceはDBの詳細を知らない）
// ctx: キャンセル/タイムアウトをDB処理まで伝える
// memo *domain.Memo: 保存したいメモ本体（ポインタ＝必要なら中身を書き換えられる）
type MemoRepository interface {
    Create(ctx context.Context, memo *domain.Memo) error
    ListByUserID(ctx context.Context, userID string) ([]domain.Memo, error)
}