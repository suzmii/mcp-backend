package handler

import (
	"mcp/applications/auth/pb/auth"
	"mcp/applications/user/pb/user"
	"mcp/bff/client"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

func Login(c *gin.Context) {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	result, err := client.User.Login(c.Request.Context(), &user.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

func Register(c *gin.Context) {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	_, err := client.User.Register(c.Request.Context(), &user.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}

func RefreshAccessToken(c *gin.Context) {
	type Request struct {
		RefreshToken string `json:"refresh_token"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	claims, err := client.Auth.VerifyRefreshToken(c.Request.Context(), &auth.VerifyRefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	result, err := client.Auth.IssueAccessToken(c.Request.Context(), &auth.IssueAccessTokenRequest{
		UserId:      claims.Claims.UserId,
		Permissions: claims.Claims.Permissions,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": result.AccessToken,
	})
}

func VerifyAccessToken(c *gin.Context) {
	type Request struct {
		AccessToken string `json:"access_token"`
	}

	var req Request

	if err := c.ShouldBind(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	_, err := client.Auth.VerifyAccessToken(c.Request.Context(), &auth.VerifyAccessTokenRequest{
		AccessToken: req.AccessToken,
	})

	if err != nil {
		logx.Info(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)

}
