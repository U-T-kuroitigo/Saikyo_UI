package web

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
)

// RegisterOrderPageRoutes registers the order page route.
func RegisterOrderPageRoutes(e *echo.Echo) {
	e.GET("/order", getOrderPage)
}

func getOrderPage(c echo.Context) error {
	// UI上で表示したい場合のみ（必須ではない）
	storeID := os.Getenv("STORE_ID")
	return c.Render(http.StatusOK, "order.html", map[string]any{
		"StoreID": storeID,
	})
}
