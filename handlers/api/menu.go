package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
)

// RegisterMenuRoutes registers API proxy routes.
func RegisterMenuRoutes(e *echo.Echo) {
	e.GET("/api/menu", proxyMenu)
}

func proxyMenu(c echo.Context) error {
	storeID := os.Getenv("STORE_ID")
	if storeID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "STORE_ID is empty"})
	}

	url := fmt.Sprintf("https://kakigori-api.fly.dev/v1/stores/%s/menu", storeID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": err.Error()})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.Blob(resp.StatusCode, "application/json", body)
}
