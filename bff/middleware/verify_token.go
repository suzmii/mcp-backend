package middleware

import (
	"mcp/applications/auth/pb/auth"
	"mcp/bff/client"
	"mcp/bff/ctxmodel"
	"mcp/core/util/ctxutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func VerifyToken(c *gin.Context) {
	authStr := c.GetHeader("Authorization")
	if len(authStr) == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	parts := strings.SplitN(authStr, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	result, err := client.Auth.VerifyAccessToken(c.Request.Context(), &auth.VerifyAccessTokenRequest{
		AccessToken: parts[1],
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Request = c.Request.WithContext(ctxutil.Set(c.Request.Context(), ctxmodel.UserID(result.Claims.UserId)))
	c.Next()
}
