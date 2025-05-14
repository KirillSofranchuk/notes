package handler

import (
	"Notes/internal/model"
	"Notes/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	service service.AbstractUserService
}

// UserReq represents user request structure
// @Description User creation/update request
type UserReq struct {
	Login    string `json:"Login" example:"user123" binding:"required"`
	Password string `json:"Password" example:"securePassword123" binding:"required"`
	Name     string `json:"Name" example:"John"`
	Surname  string `json:"Surname" example:"Doe"`
}

// UserRsp represents user response structure
// @Description User response data
type UserRsp struct {
	Id      int    `json:"Id" example:"1"`
	Login   string `json:"Login" example:"user123"`
	Name    string `json:"Name" example:"John"`
	Surname string `json:"Surname" example:"Doe"`
}

func NewUserHandler(s service.AbstractUserService) *UserHandler {
	return &UserHandler{service: s}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Register a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param input body UserReq true "User registration data"
// @Success 201 {object} int "User created successfully"
// @Failure 400 {object} response "Invalid request data"
// @Failure 409 {object} response "User already exists"
// @Failure 500 {object} response "Internal server error"
// @Router /api/user [post]
func (u UserHandler) CreateUser(c *gin.Context) {
	var req UserReq

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s", err.Error()))
		return
	}

	id, err := u.service.CreateUser(req.Login, req.Password, req.Name, req.Surname)

	if err != nil {
		apiError := model.GetAppropriateApiError(err)
		errorResponseFromApiError(c, apiError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})

	return
}

// GetUser godoc
// @Summary Get user profile
// @Description Get profile information for the authenticated user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserRsp "Returns user profile data"
// @Failure 401 {object} response "Unauthorized"
// @Failure 404 {object} response "User not found"
// @Failure 500 {object} response "Internal server error"
// @Router /api/user [get]
func (u UserHandler) GetUser(c *gin.Context) {
	userId := c.MustGet("UserId").(int)

	user, err := u.service.GetUser(userId)

	if err != nil {
		apiError := model.GetAppropriateApiError(err)
		errorResponseFromApiError(c, apiError)
		return
	}

	userRsp := UserRsp{
		Id:      userId,
		Login:   user.Login,
		Name:    user.Name,
		Surname: user.Surname,
	}
	c.JSON(http.StatusOK, gin.H{
		"user": userRsp,
	})
}

// UpdateUser godoc
// @Summary Update user profile
// @Description Update profile information for the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body UserReq true "User update data"
// @Success 200 "Profile updated successfully"
// @Failure 400 {object} response "Invalid request data"
// @Failure 401 {object} response "Unauthorized"
// @Failure 404 {object} response "User not found"
// @Failure 500 {object} response "Internal server error"
// @Router /api/user [put]
func (u UserHandler) UpdateUser(c *gin.Context) {
	var req UserReq

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s", err.Error()))
		return
	}
	userId := c.MustGet("UserId").(int)

	errUpdate := u.service.UpdateUser(userId, req.Login, req.Password, req.Name, req.Surname)

	if errUpdate != nil {
		apiError := model.GetAppropriateApiError(errUpdate)
		errorResponseFromApiError(c, apiError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// DeleteUser godoc
// @Summary Delete user account
// @Description Delete account for the authenticated user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 "Account deleted successfully"
// @Failure 401 {object} response "Unauthorized"
// @Failure 404 {object} response "User not found"
// @Failure 500 {object} response "Internal server error"
// @Router /api/user [delete]
func (u UserHandler) DeleteUser(c *gin.Context) {
	userId := c.MustGet("UserId").(int)

	err := u.service.DeleteUser(userId)

	if err != nil {
		apiError := model.GetAppropriateApiError(err)
		errorResponseFromApiError(c, apiError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
