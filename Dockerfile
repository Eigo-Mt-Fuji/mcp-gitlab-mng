# ビルドステージ
FROM golang:1.23 AS builder

WORKDIR /app

# go.mod, go.sum を先にコピーして依存解決
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# サーバーバイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# 実行ステージ
FROM debian:stable-slim

WORKDIR /app

# ビルドしたバイナリのみコピー
COPY --from=builder /app/server .

# 必要に応じて証明書などをインストール
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# 環境変数（例: GITLAB_TOKEN, GITLAB_BASE_URL）は起動時に渡す
ENV GITLAB_TOKEN= GITLAB_BASE_URL=

# デフォルトコマンド
CMD ["./server"]