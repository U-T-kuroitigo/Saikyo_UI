package web

import (
	"net/http"

	"github.com/labstack/echo"
)

// RegisterTermsRoutes は利用規約関連のルートを登録します
func RegisterTermsRoutes(e *echo.Echo) {
	e.GET("/", TermsPage)
	e.GET("/terms", TermsPage)
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
	// menu.html を直接かき氷選択画面で表示するためのデータを準備
	data := map[string]interface{}{
		"CurrentStep": "ice_flavor",
		"Flavors":     []string{"いちご", "メロン", "ブルーハワイ", "オレンジ"},
	}
	return c.Render(http.StatusOK, "menu.html", data)
}
