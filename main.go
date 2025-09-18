package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/U-T-kuroitigo/Saikyo_UI/configuration"
	"github.com/U-T-kuroitigo/Saikyo_UI/routes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Echo用テンプレートレンダラ
// views/*.html を読み込む
// main.go 側で初期化のみ行い、ルート処理は routes/handlers に分離
type Template struct{ templates *template.Template }

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// データベースの初期化
	db := configuration.InitDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	// 静的ファイル配信（/static 配下で public/ を提供）
	e.Static("/static", "public")

	// テンプレートレンダラ登録（views/*.html を対象）
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t

	// ===== ヘルスチェック用エンドポイント =====
	// /health: DB依存なし → コンテナのライブネス確認に使う
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	routes.StartRoutes(e)

	// ===== ポート設定 =====
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // デフォルト値
	}

	err := e.Start("0.0.0.0:" + port)
	if err != nil {
		fmt.Printf("Error, could not run server: %v", err)
	}
}
