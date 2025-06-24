package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"time"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String() // используем github.com/google/uuid
		start := time.Now()

		var userID interface{}
		if uid, exists := c.Get("UserId"); exists {
			userID = uid.(int)
		} else {
			userID = 0
		}

		c.Set("traceID", traceID)

		c.Next()

		log.Println(
			fmt.Sprintf("HTTP REQUEST. Method: %s, path: %s, traceId: %s, userId: %v, status: %d, duration: %d",
				c.Request.Method, c.Request.URL.Path, traceID, userID, c.Writer.Status(), time.Since(start)))
	}
}
