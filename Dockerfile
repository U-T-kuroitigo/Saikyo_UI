# =========================
# Builder stage（Goでビルド）
# =========================
FROM golang:1.25.1-alpine AS builder

# 必要ツール
RUN apk add --no-cache git tzdata ca-certificates && update-ca-certificates
# go.mod の go バージョンが新しめでも自動追従できるように
ENV GOTOOLCHAIN=auto

WORKDIR /app

# go.mod / go.sum だけ先にコピー → 依存取得をキャッシュ
COPY go.mod go.sum ./
RUN go mod download

# 残りのソースをコピー
COPY . .

# 最適化ビルド（静的リンク）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./main.go

# =========================
# Runtime stage（実行用）
# =========================
FROM alpine:3.20

# タイムゾーンと証明書
RUN apk add --no-cache tzdata ca-certificates && update-ca-certificates
ENV TZ=Asia/Tokyo

# 実行ユーザ
RUN adduser -D -H -u 10001 appuser
WORKDIR /home/appuser

# アプリ本体配置
COPY --from=builder /app/app /usr/local/bin/app
# テンプレートと静的ファイルも配置
COPY --from=builder /app/views  /home/appuser/views
COPY --from=builder /app/public /home/appuser/public

# コンテナ外へ公開するアプリのポート（例：8080）
# 実際のListenは main.go 側で 0.0.0.0:8080 にしておく
EXPOSE 8080

USER appuser

# ※ Echo は 0.0.0.0 で Listen していればOK
CMD ["/usr/local/bin/app"]
