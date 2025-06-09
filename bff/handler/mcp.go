package handler

import (
	"mcp/applications/mcp/pb/mcp"
	"mcp/bff/client"
	"mcp/bff/ctxmodel"
	"mcp/core/util/ctxutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

func CreateSession(c *gin.Context) {
	type Request struct {
		Hint string `json:"hint"`
	}

	var req Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	userId, ok := ctxutil.Get[ctxmodel.UserID](c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	result, err := client.MCP.CreateSession(c.Request.Context(), &mcp.CreateSessionRequest{
		UserId: uint64(userId),
		Hint:   req.Hint,
	})

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

func GetSessionList(c *gin.Context) {
	type Request struct {
		Page     uint64 `form:"page"`
		PageSize uint64 `form:"pageSize"`
	}

	var req Request

	logx.Info("getting session list")

	if err := c.ShouldBind(&req); err != nil {
		c.Status(http.StatusBadRequest)
	}

	logx.Info(req)

	userId, ok := ctxutil.Get[ctxmodel.UserID](c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	result, err := client.MCP.GetSessionList(c.Request.Context(), &mcp.GetSessionListRequest{
		UserId:   uint64(userId),
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

func GetSession(c *gin.Context) {
	id := c.Param("id")

	userId, ok := ctxutil.Get[ctxmodel.UserID](c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	result, err := client.MCP.GetSession(c.Request.Context(), &mcp.GetSessionRequest{
		UserId:      uint64(userId),
		SessionUuid: id,
	})

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

func AppendMessage(c *gin.Context) {
	type Request struct {
		SessionUuid string         `json:"session_uuid"`
		Message     []*mcp.Message `json:"messages"`
	}

	userId, ok := ctxutil.Get[ctxmodel.UserID](c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	var req Request

	if err := c.BindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	_, err := client.MCP.AppendMessage(c.Request.Context(), &mcp.AppendMessageRequest{
		UserId:      uint64(userId),
		SessionUuid: req.SessionUuid,
		Messages:    req.Message,
	})

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
