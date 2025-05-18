package handler

import (
	"Notes/internal/model"
	"github.com/gin-gonic/gin"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func errorResponseFromApiError(c *gin.Context, apiError *model.ApiError) {
	c.AbortWithStatusJSON(apiError.Code, response{apiError.Message})
}

func errorResponse(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, response{message})
}
