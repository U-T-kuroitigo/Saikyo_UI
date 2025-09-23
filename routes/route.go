package routes

import (
	"github.com/U-T-kuroitigo/Saikyo_UI/functions"
	"github.com/labstack/echo"

	apiHandlers "github.com/U-T-kuroitigo/Saikyo_UI/handlers/api"
	webHandlers "github.com/U-T-kuroitigo/Saikyo_UI/handlers/web"
)

func userRoutes(e *echo.Echo) {
	e.GET("api/v2/users", functions.GetAllUsers)  // GetAll Users
	e.GET("api/v2/user", functions.GetUser)       // GET one user
	e.POST("api/v2/user", functions.CreateUser)   // CREATE
	e.PUT("api/v2/user", functions.UpdateUser)    // UPDATE
	e.DELETE("api/v2/user", functions.DeleteUser) // DELETE

	// //web追加
	// e.GET("/", webHandlers.TermsPage)
    // e.GET("/terms", webHandlers.TermsPage)
    // e.GET("/agreed", webHandlers.AgreedPage)
    // e.GET("/rejected", webHandlers.RejectedPage)
}

func StartRoutes(e *echo.Echo) {
	userRoutes(e)

	// Webページのルート
	e.GET("/", webHandlers.TermsPage)
	e.GET("/terms", webHandlers.TermsPage)
	e.GET("/agreed", webHandlers.AgreedPage)
	e.GET("/rejected", webHandlers.RejectedPage)
	webHandlers.RegisterTestRoutes(e)
	webHandlers.RegisterOrderPageRoutes(e)

	// APIのルート
	apiHandlers.RegisterMenuRoutes(e)
	apiHandlers.RegisterOrderAPIRoutes(e)

	// 静的ファイルの配信設定を追加
	e.Static("/css", "public/css")
	e.Static("/js", "public/js")
}
