package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/U-T-kuroitigo/Saikyo_UI/configuration"
	"github.com/U-T-kuroitigo/Saikyo_UI/routes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

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
