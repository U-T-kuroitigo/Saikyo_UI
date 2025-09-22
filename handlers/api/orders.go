package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
)

// RegisterOrderAPIRoutes registers order API proxy routes.
func RegisterOrderAPIRoutes(e *echo.Echo) {
	e.POST("/api/orders", createOrder)
}

func createOrders(c echo.Context) error {
	storeID := os.Getenv("STORE_ID")
	if storeID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "STORE_ID is empty"})
	}

	// クライアントからの JSON ボディをそのまま上流へ
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	url := fmt.Sprintf("https://kakigori-api.fly.dev/v1/stores/%s/orders", storeID)

	req, err := http.NewRequestWithContext(c.Request().Context(), http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		c.Logger().Errorf("order proxy: build request error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "build request failed"})
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SaikyoUI/1.0 (+echo)")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.Logger().Errorf("order proxy: upstream request error: %v", err)
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "upstream request failed"})
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger().Errorf("order proxy: read body error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to read upstream body"})
	}

	// 上流のステータス/本文を透過
	return c.Blob(resp.StatusCode, "application/json", respBody)
}
