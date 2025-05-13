package handler

import (
	"Notes/internal/model"
	"Notes/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	a service.AbstractAuthService
}

// AuthReq represents authentication request structure
// @Description User authentication credentials
type AuthReq struct {
	Login    string `json:"Login" example:"user123" binding:"required"`
	Password string `json:"Password" example:"securePassword123" binding:"required"`
}

func NewAuthHandler(service service.AbstractAuthService) *AuthHandler {
	return &AuthHandler{a: service}
}

// Login godoc
// @Summary Authenticate user
// @Description Login user and get authentication token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body AuthReq true "User credentials"
// @Success 200 {object} map[string]interface{} "Returns JWT token"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/auth/login [post]
func (a *AuthHandler) Login(c *gin.Context) {
	var req AuthReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	token, err := a.a.AuthUser(req.Login, req.Password)

	if err != nil {
		apiError := model.GetAppropriateApiError(err)
		c.JSON(apiError.Code, gin.H{
			"error": apiError.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

	return
}
