package web

import (
	"net/http"
	"github.com/labstack/echo"
)

// RegisterTermsRoutes は利用規約関連のルートを登録します
func RegisterTermsRoutes(e *echo.Echo) {
	e.GET("/", TermsPage)
	e.GET("/terms", TermsPage)
	e.GET("/agreed", AgreedPage)
	e.GET("/rejected", RejectedPage)
}

// TermsPage は利用規約ページを描画します
func TermsPage(c echo.Context) error {
	return c.Render(http.StatusOK, "terms.html", nil)
}

// AgreedPage は同意後の仮ページを描画します
func AgreedPage(c echo.Context) error {
	return c.Render(http.StatusOK, "agreed_page.html", nil)
}

// RejectedPage は「ココ」をクリックした後の仮ページを描画します
func RejectedPage(c echo.Context) error {
	return c.Render(http.StatusOK, "rejected_page.html", nil)
}

