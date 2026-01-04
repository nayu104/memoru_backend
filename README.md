# プロジェクト構造と責務ドキュメント

## 目次
1. [プロジェクト概要](#プロジェクト概要)
2. [フォルダ構造と責務](#フォルダ構造と責務)
3. [アーキテクチャパターン](#アーキテクチャパターン)
4. [依存関係の注入フロー](#依存関係の注入フロー)
5. [リクエスト処理フロー](#リクエスト処理フロー)
6. [ルーティングがmain.goに集約される理由](#ルーティングがmaingoに集約される理由)

---

## プロジェクト概要

感情メモを管理するREST APIです。Go言語で実装されており、クリーンアーキテクチャの原則に従った3層構造（Handler → Service → Repository）を採用しています。

---

## フォルダ構造と責務

```
emotion-memo-api/
├── main.go                          # 【エントリーポイント】アプリケーションの起動・依存性注入・ルーティング定義
├── go.mod                           # Goモジュールの依存関係管理
├── go.sum                           # 依存関係のチェックサム
│
├── internal/                        # 外部から直接アクセスできない内部パッケージ
│   ├── domain/                      # 【ドメイン層】ビジネスエンティティの定義
│   │   └── memo.go                  # Memoエンティティ（ID, UserID, Body, Mood, CreatedAt, UpdatedAt）
│   │
│   ├── handler/                     # 【プレゼンテーション層】HTTPリクエスト/レスポンスの処理
│   │   └── memo_handler.go          # HTTPハンドラー（JSONのパース、HTTPステータスコードの返却）
│   │
│   ├── service/                     # 【ビジネスロジック層】アプリケーション固有のルールとロジック
│   │   └── memo_service.go          # ビジネスロジック（バリデーション、データの加工）
│   │
│   ├── repository/                  # 【データアクセス層】データベース操作の抽象化
│   │   ├── memo_repository.go       # リポジトリインターフェース（抽象化）
│   │   └── postgres_memo_repository.go  # PostgreSQL実装（SQLクエリの実行）
│   │
│   └── db/                          # データベース接続管理（現在は未使用、main.goで直接接続）
│       └── db.go                    # DB接続のファクトリー関数
│
└── migrations/                      # データベーススキーマのバージョン管理
    └── 001_create_tables.sql        # テーブル定義（users, memos）
```

### 各パッケージの詳細責務

#### 1. `main.go` - エントリーポイント

**責務:**
- アプリケーションの起動
- データベース接続の確立
- 各層のインスタンス生成と依存性注入（DI）
- HTTPルーティングの定義
- HTTPサーバーの起動

**依存関係:**
- 全ての内部パッケージに依存（依存性注入の実行者）

**なぜここにルーティングがあるのか:**
- アプリケーション全体の構成を一箇所で管理できる
- Handler層はHTTPの詳細（ルーティング）を知る必要がない
- 単一責任の原則：Handlerは「リクエスト処理」、main.goは「アプリケーション構成」に専念

#### 2. `internal/domain/` - ドメイン層

**責務:**
- ビジネスエンティティの定義
- エンティティの構造とフィールドの定義のみ
- ビジネスロジックは含まない（純粋なデータ構造）

**依存関係:**
- 標準ライブラリのみ（他の内部パッケージに依存しない）

**ファイル:**
- `memo.go`: `Memo`構造体の定義

```go
type Memo struct {
    ID        int64     `json:"id"`
    UserID    string    `json:"user_id"`
    Body      string    `json:"body"`
    Mood      string    `json:"mood"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### 3. `internal/handler/` - プレゼンテーション層

**責務:**
- HTTPリクエストの受信
- JSONのパース/バリデーション（形式チェック）
- HTTPステータスコードの設定
- JSONレスポンスの生成
- Service層への処理委譲

**依存関係:**
- `internal/service/`（下位層）
- `internal/domain/`（ドメインエンティティ）

**ファイル:**
- `memo_handler.go`: 
  - `CreateMemo`: POST /memos の処理
  - `ListMemos`: GET /memos の処理

**重要な原則:**
- HTTPの詳細（ルーティング、URLパス）は知らない
- ビジネスロジックは含まない（Service層に委譲）

#### 4. `internal/service/` - ビジネスロジック層

**責務:**
- アプリケーション固有のビジネスルールの実装
- データのバリデーション（例：bodyが空でないか）
- データの加工（例：CreatedAt, UpdatedAtの設定）
- Repository層への処理委譲

**依存関係:**
- `internal/repository/`（インターフェース）
- `internal/domain/`

**ファイル:**
- `memo_service.go`:
  - `CreateMemo`: メモ作成のビジネスロジック
  - `ListMemos`: メモ一覧取得のビジネスロジック

**重要な原則:**
- HTTPの詳細を知らない
- データベースの実装詳細を知らない（Repositoryインターフェースに依存）

#### 5. `internal/repository/` - データアクセス層

**責務:**
- データベース操作の抽象化
- SQLクエリの実行
- データベースレコード ↔ Go構造体の変換

**依存関係:**
- `internal/domain/`
- `database/sql`（標準ライブラリ）

**ファイル:**
- `memo_repository.go`: `MemoRepository`インターフェースの定義
- `postgres_memo_repository.go`: PostgreSQL実装

**重要な原則:**
- ビジネスロジックを含まない（単純なCRUD操作のみ）
- インターフェースを定義することで、データベース実装を交換可能にする

---

## アーキテクチャパターン

### レイヤードアーキテクチャ（クリーンアーキテクチャ）

```
┌─────────────────────────────────────────────────────────┐
│                    Handler Layer                        │
│              (HTTPリクエスト/レスポンス)                │
└───────────────────────┬─────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────┐
│                    Service Layer                        │
│              (ビジネスロジック)                         │
└───────────────────────┬─────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────┐
│                  Repository Layer                       │
│              (データアクセス)                           │
└───────────────────────┬─────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────┐
│              PostgreSQL Database                        │
└─────────────────────────────────────────────────────────┘
```

### 依存関係の方向

- **上位層は下位層に依存する**
- **下位層は上位層に依存しない**
- **Service層はRepositoryインターフェースに依存（実装に依存しない）**

---

## 依存関係の注入フロー

### main.goでの組み立て順序

```go
// 1. データベース接続の確立
db := sql.Open("postgres", connStr)

// 2. Repository層の作成（最下位層から開始）
memoRepo := repository.NewPostgresMemoRepository(db)
// → PostgresMemoRepository{db: db} が生成される

// 3. Service層の作成（Repositoryを注入）
memoService := memo_service.NewMemoService(memoRepo)
// → MemoService{repo: memoRepo} が生成される
// → memoRepoはMemoRepositoryインターフェースとして扱われる

// 4. Handler層の作成（Serviceを注入）
memoHandler := handler.NewMemoHandler(memoService)
// → MemoHandler{service: memoService} が生成される

// 5. ルーティングの定義
mux.HandleFunc("POST /memos", memoHandler.CreateMemo)
mux.HandleFunc("GET /memos", memoHandler.ListMemos)
```

### 依存関係の図

```
main.go
  │
  ├─→ db (*sql.DB)
  │     │
  │     └─→ PostgresMemoRepository{db: db}
  │            │
  │            └─→ MemoService{repo: PostgresMemoRepository}
  │                   │
  │                   └─→ MemoHandler{service: MemoService}
  │                          │
  │                          └─→ mux.HandleFunc()でルーティング登録
```

---

## リクエスト処理フロー

### POST /memos (メモ作成) の処理フロー

```
1. HTTPリクエスト到着
   POST /memos
   Body: {"user_id": "xxx", "body": "今日は...", "mood": "happy"}
   │
   ▼
2. main.goのルーティング
   mux.HandleFunc("POST /memos", memoHandler.CreateMemo)
   → memoHandler.CreateMemo() が呼ばれる
   │
   ▼
3. Handler層 (memo_handler.go:21)
   - JSONをパース → createMemoRequest構造体に変換
   - バリデーション（JSON形式のチェックのみ）
   - handler.service.CreateMemo() を呼ぶ
   │
   ▼
4. Service層 (memo_service.go:30)
   - ビジネスロジックのバリデーション（bodyが空でないか）
   - 現在時刻の設定
   - domain.Memo構造体を作成（IDはまだ空）
   - service.repo.Create(ctx, memo) を呼ぶ
   │
   ▼
5. Repository層 (postgres_memo_repository.go:19)
   - SQL INSERT文を実行
   - RETURNING句でmemo_idを取得
   - memo.ID に値を書き込む（ポインタなのでService側にも反映）
   │
   ▼
6. 戻り値の伝播
   Repository → Service → Handler
   - 各層でエラーがあればそのまま上位層に返す
   │
   ▼
7. HTTPレスポンス
   - Handler層でJSONにエンコード
   - 201 Created + JSON本文を返す
```

### GET /memos?user_id=xxx (メモ一覧取得) の処理フロー

```
1. HTTPリクエスト到着
   GET /memos?user_id=xxx
   │
   ▼
2. main.goのルーティング
   mux.HandleFunc("GET /memos", memoHandler.ListMemos)
   → memoHandler.ListMemos() が呼ばれる
   │
   ▼
3. Handler層 (memo_handler.go:59)
   - クエリパラメータからuser_idを取得
   - handler.service.ListMemos(ctx, userID) を呼ぶ
   │
   ▼
4. Service層 (memo_service.go:59)
   - 特に処理なし（そのままRepositoryに委譲）
   - service.repo.ListByUserID(ctx, userID) を呼ぶ
   │
   ▼
5. Repository層 (postgres_memo_repository.go:39)
   - SQL SELECT文を実行
   - 各行をdomain.Memo構造体に変換
   - []domain.Memo を返す
   │
   ▼
6. HTTPレスポンス
   - Handler層でJSON配列にエンコード
   - 200 OK + JSON本文を返す
```

---

## ルーティングがmain.goに集約される理由

### 1. 単一責任の原則（Single Responsibility Principle）

**Handler層の責務:**
- HTTPリクエストの処理（JSONのパース、HTTPステータスコードの設定）
- ルーティングの定義は責務外

**main.goの責務:**
- アプリケーション全体の構成（依存性注入、ルーティング定義）

### 2. 関心の分離（Separation of Concerns）

- Handler層: 「どのように処理するか」に集中
- main.go: 「どのURLにどのHandlerを割り当てるか」に集中

### 3. テスト容易性

- Handler層をテストする際、ルーティングの詳細を知る必要がない
- Handlerのメソッドを直接呼び出してテスト可能

### 4. 柔軟性

- ルーティングを変更する際、Handler層のコードを変更する必要がない
- ミドルウェアの追加などもmain.goで一括管理できる

### 5. アプリケーションの全体像が一目でわかる

```go
// main.goを見れば、このアプリケーションが提供するエンドポイントが一覧できる
mux.HandleFunc("POST /memos", memoHandler.CreateMemo)  // メモ作成
mux.HandleFunc("GET /memos", memoHandler.ListMemos)    // メモ一覧
// 将来的に他のエンドポイントを追加する場合も、ここに追加すればよい
```

### 代替案との比較

#### ❌ もしHandler層にルーティングがあったら...

```go
// handler/memo_handler.go（悪い例）
func RegisterRoutes(mux *http.ServeMux, service *MemoService) {
    handler := NewMemoHandler(service)
    mux.HandleFunc("POST /memos", handler.CreateMemo)
    mux.HandleFunc("GET /memos", handler.ListMemos)
}
```

**問題点:**
- Handler層がルーティングの詳細（URLパス）を知る必要がある
- ルーティングの変更時にHandler層のコードを変更する必要がある
- アプリケーション全体のルーティングが分散して見えにくくなる

#### ✅ 現在の設計（main.goに集約）

**メリット:**
- Handler層は処理ロジックに専念できる
- ルーティングの変更がmain.goのみで完結
- アプリケーション全体のエンドポイントが一目でわかる

---

## データの流れ

### メモ作成時のデータ変換

```
HTTPリクエスト (JSON)
  ↓
createMemoRequest構造体 {UserID, Body, Mood}
  ↓
Service層で変換
  ↓
domain.Memo構造体 {ID=0, UserID, Body, Mood, CreatedAt, UpdatedAt}
  ↓
Repository層でDBに保存
  ↓
domain.Memo構造体 {ID=123, UserID, Body, Mood, CreatedAt, UpdatedAt}  ← IDが設定される
  ↓
HTTPレスポンス (JSON)
```

---

## エラーハンドリングの流れ

```
Repository層でエラー発生
  ↓
return fmt.Errorf("❌メモの作成に失敗しました: %w", err)
  ↓
Service層で受け取る
  ↓
return nil, err  // そのまま返す（追加のエラーハンドリングなし）
  ↓
Handler層で受け取る
  ↓
http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
```

**原則:**
- 各層は下位層のエラーをそのまま上位層に返す
- エラーを「処理」するのはHandler層（HTTPステータスコードに変換）

---

## まとめ

### 各層の責務

| 層 | 責務 | 知るべきこと | 知らなくてもよいこと |
|---|---|---|---|
| **Handler** | HTTPリクエスト/レスポンスの処理 | HTTP、JSON、Service層のインターフェース | ルーティングの詳細、ビジネスロジック、DBの詳細 |
| **Service** | ビジネスロジック | ビジネスルール、Repository層のインターフェース | HTTPの詳細、DBの実装詳細 |
| **Repository** | データアクセス | SQL、データベース操作 | HTTPの詳細、ビジネスロジック |
| **main.go** | アプリケーション構成 | 全て（依存性注入とルーティング） | 各層の内部実装詳細 |

### ルーティングがmain.goに集約される理由（再確認）

1. **単一責任の原則**: Handlerは「処理」、main.goは「構成」
2. **関心の分離**: 処理ロジックとルーティング定義を分離
3. **テスト容易性**: Handlerをルーティングから独立してテスト可能
4. **柔軟性**: ルーティング変更がmain.goのみで完結
5. **可読性**: アプリケーション全体のエンドポイントが一覧できる

---

## 参考資料

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

