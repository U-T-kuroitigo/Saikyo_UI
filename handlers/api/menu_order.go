package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
)

// RegisterMenuOrderRoutes は、かき氷の注文APIのエンドポイントを登録します。
// main.goなどで `api.RegisterMenuOrderRoutes(e)` のように呼び出す必要があります。
func RegisterMenuOrderRoutes(e *echo.Echo) {
	e.POST("/api/orders", createMenuOrder)
}

// createMenuOrder は、クライアントからの注文リクエストを受け取り、外部のAPIへ中継します。
func createMenuOrder(c echo.Context) error {
	// ▼▼▼ .envファイルなどから環境変数を読み込むように修正 ▼▼▼
	storeID := os.Getenv("STORE_ID")
	if storeID == "" {
		log.Println("エラー: 環境変数 STORE_ID が設定されていません。")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "サーバーの設定エラーです。"})
	}
	// ▲▲▲ ここまで修正 ▲▲▲

	// クライアントから送られてきたリクエストのボディ（JSON）を読み込みます
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストボディが不正です"})
	}

	// 中継先となる外部APIのURLを構築します
	url := fmt.Sprintf("https://kakigori-api.fly.dev/v1/stores/%s/orders", storeID)

	// 外部APIへ送信するための新しいリクエストを作成します
	req, err := http.NewRequestWithContext(c.Request().Context(), http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		c.Logger().Errorf("注文APIプロキシ: リクエスト構築エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "リクエストの構築に失敗しました"})
	}

	// 必要なHTTPヘッダーを設定します
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SaikyoUI/1.0 (+echo)")

	// HTTPクライアントを使って、作成したリクエストを外部APIに送信します
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.Logger().Errorf("注文APIプロキシ: 外部APIへのリクエストエラー: %v", err)
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "外部APIへのリクエストに失敗しました"})
	}
	defer resp.Body.Close()

	// 外部APIからのレスポンスを読み込みます
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger().Errorf("注文APIプロキシ: レスポンスボディの読み込みエラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "外部APIのレスポンス読み込みに失敗しました"})
	}

	// 外部APIからのレスポンスを、そのままクライアントに返します
	return c.Blob(resp.StatusCode, "application/json", respBody)
}

