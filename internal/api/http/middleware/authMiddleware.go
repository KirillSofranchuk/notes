package middleware

import (
	"Notes/internal/model"
	"Notes/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(service service.AbstractAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No token provided"})
			c.Abort()

			return
		}

		tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

		claims, err := service.ValidateToken(tokenString)

		if err != nil {
			apiError := model.GetAppropriateApiError(err)

			if apiError.Code == http.StatusBadRequest {
				c.JSON(apiError.Code, gin.H{"message": apiError.Message})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "пользователь не найден"})
			}

			c.Abort()

			return
		}

		c.Set("UserId", claims.UserId)
		c.Next()
	}
}
