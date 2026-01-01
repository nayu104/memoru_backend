package service

import (
	"context"
	"errors"
	"time"

	"emotion-memo-api/internal/domain"
	"emotion-memo-api/internal/repository"
)

type MemoService struct {
	// このサービスが仕事をするには、DB担当者（Repository）が必要です。
	// なので、repo という名前で MemoRepository（インターフェース）を持たせます。
	repo repository.MemoRepository
}

// NewMemoService は、MemoService（担当者）を新しく雇うための関数（コンストラクタ）です。
// 引数で「DB担当者(repo)」を受け取り、それを装備した「Service担当者」を返します。
func NewMemoService(dbRepo repository.MemoRepository) *MemoService {
	// 構文: 「構造体のフィールド名: 入れる変数の名前」
	// 左側の repo: MemoServiceが持っているポケットの名前
	// 右側の repo: 引数で渡されてきた変数の名前
	return &MemoService{repo: dbRepo}
}

// CreateMemo: メモを作成するメソッドです。
//
//	①誰の技？              ②入力データは？           ③出力データは？
func (service *MemoService) CreateMemo(ctx context.Context, userID, body, mood string) (*domain.Memo, error) {
	if body == "" {
		return nil, errors.New("❌メモの内容は必須です😡")
	}

	// DBに入れる前に、プログラム側で現在時刻を決めます。
	now := time.Now()

	// domain.Memo の構造体（データの塊）を作ります。
	// まだDBに入れていないので、IDは空っぽのままです。
	memo := &domain.Memo{
		UserID:    userID,
		Body:      body,
		Mood:      mood,
		CreatedAt: now, // 作成日時
		UpdatedAt: now, // 更新日時（作成時は同じ）
	}

	// 3. DB担当者（Repository）に仕事を依頼する
	// service.repo : 私(service)が持っている、DB担当者(repo)
	// .Create: その担当者の「Create」技を使う
	// ここで初めてSQLが実行され、memo.ID に数値が入ります
	if err := service.repo.Create(ctx, memo); err != nil {
		return nil, err
	}
	// 成功したら、IDが入った状態のメモデータを返します。
	return memo, nil
}

func (service *MemoService) ListMemos(ctx context.Context, userID string) ([]domain.Memo, error) {
	// 今回は特に「加工」や「チェック」のルールがないので、
	// そのまま部下である DB担当者(service.repo) に丸投げしています。
	// 「DBから取ってきて！」と頼んでいるだけです。
	return service.repo.ListByUserID(ctx, userID)
}
