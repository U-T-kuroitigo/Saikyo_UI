package web

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
)

// RegisterTestRoutes registers web page routes.
func RegisterTestRoutes(e *echo.Echo) {
	e.GET("/test", getTest)
}

func getTest(c echo.Context) error {
	storeID := os.Getenv("STORE_ID")
	return c.Render(http.StatusOK, "test.html", map[string]any{
		"StoreID": storeID,
	})
}
