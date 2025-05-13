package handler

import (
	"Notes/internal/model"
	"Notes/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FolderHandler struct {
	service service.AbstractFolderService
}

func NewFolderHandler(s service.AbstractFolderService) FolderHandler {
	return FolderHandler{service: s}
}

type FolderReq struct {
	Title string `json:"Title" example:"My Folder" binding:"required"`
}

// CreateFolder godoc
// @Summary Create a new folder
// @Description Create a new folder for the authenticated user
// @Tags folders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body FolderReq true "Folder creation data"
// @Success 200 {object} map[string]interface{} "Returns ID of created folder"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/folder [post]
func (f *FolderHandler) CreateFolder(c *gin.Context) {
	var req FolderReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	userId := c.MustGet("UserId").(int)

	id, err := f.service.CreateFolder(userId, req.Title)

	if err != nil {
		apiError := model.GetAppropriateApiError(err)
		c.JSON(apiError.Code, gin.H{
			"error": apiError.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

// UpdateFolder godoc
// @Summary Update a folder
// @Description Update an existing folder for the authenticated user
// @Tags folders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Folder ID"
// @Param input body FolderReq true "Folder update data"
// @Success 200 "Folder updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data or ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Folder not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/folder/{id} [put]
func (f *FolderHandler) UpdateFolder(c *gin.Context) {
	var req FolderReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	userId := c.MustGet("UserId").(int)

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	errUpdate := f.service.UpdateFolder(userId, idInt, req.Title)

	if errUpdate != nil {
		apiError := model.GetAppropriateApiError(errUpdate)
		c.JSON(apiError.Code, gin.H{
			"error": apiError.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// DeleteFolder godoc
// @Summary Delete a folder
// @Description Delete an existing folder for the authenticated user
// @Tags folders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Folder ID"
// @Success 200 "Folder deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Folder not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/folder/{id} [delete]
func (f *FolderHandler) DeleteFolder(c *gin.Context) {
	userId := c.MustGet("UserId").(int)

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	errDelete := f.service.DeleteFolder(userId, idInt)

	if errDelete != nil {
		apiError := model.GetAppropriateApiError(errDelete)
		c.JSON(apiError.Code, gin.H{
			"error": apiError.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
