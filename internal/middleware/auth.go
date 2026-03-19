package middleware

import (
	"net/http"
	"peer-link-server/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, 40100, "unauthorized")
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			response.Fail(c, http.StatusUnauthorized, 40100, "unauthorized")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Fail(c, http.StatusUnauthorized, 40100, "unauthorized")
			return
		}

		if uid, ok := claims["user_id"]; ok {
			c.Set("user_id", uid)
		}

		c.Next()
	}
}
