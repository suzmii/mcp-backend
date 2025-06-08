package router

import (
	"mcp/bff/handler"
	"mcp/bff/middleware"

	"github.com/gin-gonic/gin"
)

var Router = gin.Default()
var apiRouter = Router.Group("/v1")

func init() {
	{
		apiRouter.POST("/users/register", handler.Register)
		apiRouter.POST("/users/login", handler.Login)
		apiRouter.POST("/auth/refresh", handler.RefreshAccessToken)
	}
	{
		g := apiRouter.Group("/mcp", middleware.VerifyToken)
		g.POST("/sessions", handler.CreateSession)
		g.GET("/sessions", handler.GetSessionList)
		g.GET("/sessions/:id", handler.GetSession)
		g.POST("/messages", handler.AppendMessage)
	}
}
