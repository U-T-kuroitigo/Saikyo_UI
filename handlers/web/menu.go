package web

import (
	"net/http"

	"github.com/U-T-kuroitigo/Saikyo_UI/handlers/api"
	"github.com/labstack/echo"
)

// RegisterMenuPageRoutes はメニューページのルートを登録します。
func RegisterMenuPageRoutes(e *echo.Echo) {
	e.GET("/menu", HandleMenu)
	e.POST("/menu", HandleMenu)
}

// HandleMenuはリクエストをロジック層に渡し、結果をレンダリングします。
func HandleMenu(c echo.Context) error {
	// フォームの全パラメータをマップに変換します
	// c.FormParams() はPOSTリクエストの application/x-www-form-urlencoded 形式のボディを解析します
	formParams, err := c.FormParams()
	if err != nil {
		c.Logger().Errorf("フォームパラメータの解析に失敗しました: %v", err)
		return c.String(http.StatusBadRequest, "リクエストの形式が正しくありません。")
	}

	params := make(map[string]string)
	for k, v := range formParams {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	// ビジネスロジックを呼び出して、表示すべきデータを取得します
	data := api.ProcessMenuStep(params)

	// main.goで設定されたレンダラーを使い、テンプレートとデータを渡してHTMLを生成します
	// これにより、リクエストごとにファイルを読み込む必要がなくなります
	if err := c.Render(http.StatusOK, "menu.html", data); err != nil {
		c.Logger().Errorf("テンプレートのレンダリングに失敗しました: %v", err)
		return err
	}

	return nil
}

