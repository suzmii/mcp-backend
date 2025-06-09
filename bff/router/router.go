package router

import (
	"mcp/bff/handler"
	"mcp/bff/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Router = gin.Default()

func init() {
	Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 设置允许的来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	apiRouter := Router.Group("/v1")
	{
		apiRouter.POST("/user/register", handler.Register)
		apiRouter.POST("/user/login", handler.Login)
		apiRouter.POST("/user/token/refresh", handler.RefreshAccessToken)
		apiRouter.POST("/user/token/verify/access", handler.VerifyAccessToken)
	}
	{
		g := apiRouter.Group("/mcp", middleware.VerifyToken)
		g.POST("/session", handler.CreateSession)
		g.GET("/sessions", handler.GetSessionList)
		g.GET("/sessions/:id", handler.GetSession)
		g.POST("/message", handler.AppendMessage)
	}
}
