package handler

import (
	"Notes/internal/model"
	"Notes/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FolderHandler struct {
	folderService service.AbstractFolderService
}

func NewFolderHandler(s service.AbstractFolderService) *FolderHandler {
	return &FolderHandler{folderService: s}
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
// @Success 200 {object} int "Returns ID of created folder"
// @Failure 400 {object} response
// @Failure 401 {object} response
// @Failure 500 {object} response
// @Router /api/folder [post]
func (f *FolderHandler) CreateFolder(c *gin.Context) {
	var req FolderReq

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s", err.Error()))
		return
	}

	userId := c.MustGet("UserId").(int)

	id, err := f.folderService.CreateFolder(userId, req.Title)

	if err != nil {
		apiError := model.GetAppropriateApiError(err)
		errorResponseFromApiError(c, apiError)
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
// @Failure 400 {object} response "Invalid request data or ID"
// @Failure 401 {object} response "Unauthorized"
// @Failure 404 {object} response "Folder not found"
// @Failure 500 {object} response "Internal server error"
// @Router /api/folder/{id} [put]
func (f *FolderHandler) UpdateFolder(c *gin.Context) {
	var req FolderReq

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s", err.Error()))
		return
	}

	userId := c.MustGet("UserId").(int)

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	errUpdate := f.folderService.UpdateFolder(userId, idInt, req.Title)

	if errUpdate != nil {
		apiError := model.GetAppropriateApiError(errUpdate)
		errorResponseFromApiError(c, apiError)
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
// @Failure 400 {object} response "Invalid ID"
// @Failure 401 {object} response "Unauthorized"
// @Failure 404 {object} response "Folder not found"
// @Failure 500 {object} response "Internal server error"
// @Router /api/folder/{id} [delete]
func (f *FolderHandler) DeleteFolder(c *gin.Context) {
	userId := c.MustGet("UserId").(int)

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	errDelete := f.folderService.DeleteFolder(userId, idInt)

	if errDelete != nil {
		apiError := model.GetAppropriateApiError(errDelete)
		errorResponseFromApiError(c, apiError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
