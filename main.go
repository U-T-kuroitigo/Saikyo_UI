package main

import (
    "fmt"
    "html/template"
    "io"
    // "net/http"
    "os"

    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"

    // 注文システムのハンドラをインポート
    "github.com/U-T-kuroitigo/Saikyo_UI/handlers/web"
)

// Echo用テンプレートレンダラ
type Template struct{ templates *template.Template }

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
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

    // 注文システム用ルート
    // "/"に全てのロジックを統合
    e.GET("/", web.HandleMenu)
    e.POST("/", web.HandleMenu)

    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }

    err := e.Start("0.0.0.0:" + port)
    if err != nil {
        fmt.Printf("Error, could not run server: %v", err)
    }
}