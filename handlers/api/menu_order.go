package api

import (
	"bytes"
	"encoding/json" // JSONを扱うためにインポートを追加
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
)

// RegisterMenuOrderRoutes は、かき氷の注文APIのエンドポイントを登録します。
// 注文作成(POST)と状況確認(GET)の両方を登録します。
func RegisterMenuOrderRoutes(e *echo.Echo) {
	e.GET("/api/orders/id", getMenuOrderStatus)
	e.POST("/api/orders", createMenuOrder)
}

// createMenuOrder は、クライアントからの注文リクエストを受け取り、外部のAPIへ中継します。
func createMenuOrder(c echo.Context) error {
	// .envファイルなどから環境変数を読み込むように修正
	storeID := os.Getenv("STORE_ID")
	if storeID == "" {
		log.Println("エラー: 環境変数 STORE_ID が設定されていません。")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "サーバーの設定エラーです。"})
	}

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

	// デバッグのため、外部APIからのレスポンスと注文IDをログに出力します
	log.Printf("DEBUG: 外部APIからのレスポンスボディ: %s", string(respBody))

	// レスポンスから注文IDをパースしてログに出力
	var orderResponse struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBody, &orderResponse); err == nil {
		log.Printf("DEBUG: 外部APIから正常に注文ID(%s)を取得しました。", orderResponse.ID)
	} else {
		log.Printf("DEBUG: 注文IDのJSONパースに失敗しました: %v", err)
	}

	// 外部APIからのレスポンスを、そのままクライアントに返します
	return c.Blob(resp.StatusCode, "application/json", respBody)
}

// getMenuOrderStatus は、特定の注文IDの状況を外部APIに問い合わせて返します。
func getMenuOrderStatus(c echo.Context) error {
	// storeID := os.Getenv("STORE_ID")
	// if storeID == "" {
	// 	log.Println("エラー: 環境変数 STORE_ID が設定されていません。")
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "サーバーの設定エラーです。"})
	// }

	// id := c.Param("id")

	// if id == "" {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": "注文IDが指定されていません"})
	// }

	// url := fmt.Sprintf("https://kakigori-api.fly.dev/v1/stores/%s/orders/%s", storeID, id)
	// fmt.Print(url)

	// req, err := http.NewRequestWithContext(c.Request().Context(), http.MethodGet, url, nil)
	// if err != nil {
	// 	c.Logger().Errorf("注文状況確認プロキシ: リクエスト構築エラー: %v", err)
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "リクエストの構築に失敗しました"})
	// }
	// req.Header.Set("User-Agent", "SaikyoUI/1.0 (+echo)")

	// client := &http.Client{Timeout: 10 * time.Second}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	c.Logger().Errorf("注文状況確認プロキシ: 外部APIへのリクエストエラー: %v", err)
	// 	return c.JSON(http.StatusBadGateway, map[string]string{"error": "外部APIへのリクエストに失敗しました"})
	// }
	// defer resp.Body.Close()

	// respBody, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	c.Logger().Errorf("注文状況確認プロキシ: レスポンスボディの読み込みエラー: %v", err)
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "外部APIのレスポンス読み込みに失敗しました"})
	// }

	// return c.Blob(resp.StatusCode, "application/json", respBody)
	return c.JSON(http.StatusOK, map[string]string{"status": "waiting-pickup"})
}

