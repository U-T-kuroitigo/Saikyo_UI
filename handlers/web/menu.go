package web

import (
	"net/http"

	// 自身のプロジェクトパスに合わせて変更してください
	"github.com/U-T-kuroitigo/Saikyo_UI/handlers/api"

	"github.com/labstack/echo"
)

// メニューページのルートを登録
func RegisterMenuPageRoutes(e *echo.Echo) {
	e.GET("/menu", HandleMenu)
	e.POST("/menu", HandleMenu)
}

// HandleMenuはリクエストをロジック層に渡し、結果をレンダリングする
func HandleMenu(c echo.Context) error {
	// フォームの全パラメータをマップに変換
	formParams, _ := c.FormParams()
	params := make(map[string]string)
	for k, v := range formParams {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	// GETリクエストの場合はparamsが空のマップになる
	// POSTリクエストの場合は送信されたフォームの値が入る

	// ビジネスロジックを呼び出して、表示すべきデータを取得
	data := api.ProcessMenuStep(params)

	// テンプレートとデータを渡してHTMLを生成
	return c.Render(http.StatusOK, "menu.html", data)
}
