# 1. ビルド用のステージ（料理を作る場所）
FROM golang:1.23-alpine AS builder

# 作業ディレクトリを作成
WORKDIR /app

# 依存関係のファイルをコピーしてダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# Goのアプリをビルド（mainという名前の実行ファイルを作る）
RUN go build -o main ./main.go

# 2. 実行用のステージ（料理を提供するお皿）
# alpineという超軽量Linuxを使う
FROM alpine:latest

WORKDIR /root/

# ビルド用ステージから、完成した実行ファイルだけを持ってくる
COPY --from=builder /app/main .

# ポート8080を開ける
EXPOSE 8080

# アプリを起動
CMD ["./main"]