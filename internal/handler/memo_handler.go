package handler

import (
	"encoding/json"
	"net/http"

	memo_service "emotion-memo-api/internal/service"
)

// MemoHandler は HTTP リクエストを受け取り MemoService に処理を委譲する窓口
type MemoHandler struct {
	service *memo_service.MemoService
}

// NewMemoHandler は MemoService を注入してハンドラを初期化する
func NewMemoHandler(memoService *memo_service.MemoService) *MemoHandler {
	return &MemoHandler{service: memoService}
}

// CreateMemo は POST /memos を処理し、新しいメモを登録する
func (handler *MemoHandler) CreateMemo(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	// 1. クライアントから届く JSON を受け取る一時的な箱を定義する
	type createMemoRequest struct {
		UserID string `json:"user_id"` // 誰のメモか
		Body   string `json:"body"`    // メモ本文
		Mood   string `json:"mood"`    // 気分（任意）
	}

	var requestPayload createMemoRequest

	// 2. HTTP ボディ(JSON) → Go構造体へ変換する
	//    Decode が失敗したら「JSONが壊れている」と判断し、400 Bad Request を返す
	if err := json.NewDecoder(httpRequest.Body).Decode(&requestPayload); err != nil {
		http.Error(responseWriter, "invalid JSON", http.StatusBadRequest)
		return
	}

	// 3. Service 層に「メモを作って」と依頼する
	//    Context を渡すことでキャンセルやタイムアウトを伝播できる
	createdMemo, err := handler.service.CreateMemo(
		httpRequest.Context(),
		requestPayload.UserID,
		requestPayload.Body,
		requestPayload.Mood,
	)
	if err != nil {
		// DBエラーなどは 500 Internal Server Error を返す
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. 成功したら 201 Created と JSON 本文を返す
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(responseWriter).Encode(createdMemo) // エンコード失敗はログだけで十分
}

// ListMemos は GET /memos?user_id=xxx を処理し、メモ一覧を返す
func (handler *MemoHandler) ListMemos(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	// 1. クエリパラメータから user_id を必須取得
	userID := httpRequest.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(responseWriter, "user_id is required", http.StatusBadRequest)
		return
	}

	// 2. Service に一覧取得を依頼
	memoList, err := handler.service.ListMemos(httpRequest.Context(), userID)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. JSON で結果を返す
	responseWriter.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(responseWriter).Encode(memoList)
}
