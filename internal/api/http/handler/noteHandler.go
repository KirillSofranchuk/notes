package handler

import (
	"Notes/internal/model"
	"Notes/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type NoteHandler struct {
	noteService service.AbstractNoteService
}

type NoteRq struct {
	Title   string    `json:"Title" example:"My Note" binding:"required"`
	Content string    `json:"Content" example:"Note content"`
	Tags    *[]string `json:"Tags" example:"tag1,tag2"`
}

type MoveNoteRq struct {
	FolderId *int `json:"FolderId" example:"1" binding:"required"`
}

func NewNoteHandler(s service.AbstractNoteService) *NoteHandler {
	return &NoteHandler{noteService: s}
}

// CreateNote godoc
// @Summary Create a new note
// @Description Create a new note for the authenticated user
// @Tags notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body NoteRq true "Note creation data"
// @Success 200 {object} int "Returns ID of created note"
// @Failure 400 {object} response "Invalid request data"
// @Failure 401 {object} response "Unauthorized"
// @Failure 500 {object} response "Internal server error"
// @Router /api/notes [post]
func (n *NoteHandler) CreateNote(c *gin.Context) {
	var req NoteRq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s", err.Error()))
		return
	}

	userId := c.MustGet("UserId").(int)

	id, err := n.noteService.CreateNote(userId, req.Title, req.Content, req.Tags)
	if err != nil {
		apiError := model.GetAppropriateApiError(err)
		errorResponseFromApiError(c, apiError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// UpdateNote godoc
// @Summary Update a note
// @Description Update an existing note for the authenticated user
// @Tags notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Note ID"
// @Param input body NoteRq true "Note update data"
// @Success 200 "Note updated successfully"
// @Failure 400 {object} response "Invalid request data or ID"
// @Failure 401 {object} response "Unauthorized"
// @Failure 404 {object} response "Note not found"
// @Failure 500 {object} response "Internal server error"
// @Router /api/notes/{id} [put]
func (n *NoteHandler) UpdateNote(c *gin.Context) {
	var req NoteRq
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

	errUpdate := n.noteService.UpdateNote(userId, idInt, req.Title, req.Content, req.Tags)

	if errUpdate != nil {
		apiError := model.GetAppropriateApiError(errUpdate)
		errorResponseFromApiError(c, apiError)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteNote godoc
// @Summary Delete a note
// @Description Delete an existing note for the authenticated user
// @Tags notes
// @Produce json
// @Security BearerAuth
// @Param id path int true "Note ID"
// @Success 200 "Note deleted successfully"
// @Failure 400 {object} response "Invalid ID"
// @Failure 401 {object} response "Unauthorized"
// @Failure 404 {object} response "Note not found"
// @Failure 500 {object} response "Internal server error"
// @Router /api/notes/{id} [delete]
func (n *NoteHandler) DeleteNote(c *gin.Context) {
	userId := c.MustGet("UserId").(int)

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	errDelete := n.noteService.DeleteNote(userId, idInt)

	if errDelete != nil {
		apiError := model.GetAppropriateApiError(errDelete)
		errorResponseFromApiError(c, apiError)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// GetFavoriteNotes godoc
// @Summary Get favorite notes
// @Description Get all favorite notes for the authenticated user
// @Tags notes
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []model.Note "Returns list of favorite notes"
// @Failure 401 {object} response "Unauthorized"
// @Failure 500 {object} response "Internal server error"
// @Router /api/notes/favorites [get]
func (n *NoteHandler) GetFavoriteNotes(c *gin.Context) {
	userId := c.MustGet("UserId").(int)

	favorites := n.noteService.GetFavoriteNotes(userId)

	c.JSON(http.StatusOK, gin.H{
		"notes": favorites,
	})
}

// FindNotes godoc
// @Summary Search notes
// @Description Search notes by query phrase for the authenticated user
// @Tags notes
// @Produce json
// @Security BearerAuth
// @Param query query string true "Search phrase"
// @Success 200 {object} []model.Note "Returns list of matching notes"
// @Failure 400 {object} response "Empty query parameter"
// @Failure 401 {object} response "Unauthorized"
// @Failure 500 {object} response "Internal server error"
// @Router /api/notes/search [get]
func (n *NoteHandler) FindNotes(c *gin.Context) {
	userId := c.MustGet("UserId").(int)
	queryPhrase := c.Query("query")

	if queryPhrase == "" {
		errorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	notes := n.noteService.FindNotesByQueryPhrase(userId, queryPhrase)
	c.JSON(http.StatusOK, gin.H{
		"notes": notes,
	})

}

// MoveNote godoc
// @Summary Moves note
// @Description Move note from/out folder for the authenticated user
// @Tags notes
// @Produce json
// @Security BearerAuth
// @Param id path int true "Note ID"
// @Param input body MoveNoteRq true "Note update data"
// @Success 200 {object} string "Note updated successfully"
// @Failure 400 {object} response "Empty query parameter"
// @Failure 401 {object} response "Unauthorized"
// @Failure 500 {object} response "Internal server error"
// @Router /api/notes/{id}/move [put]
func (n *NoteHandler) MoveNote(c *gin.Context) {
	var req MoveNoteRq
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

	errMove := n.noteService.MoveToFolder(userId, idInt, req.FolderId)

	if errMove != nil {
		apiError := model.GetAppropriateApiError(errMove)
		errorResponseFromApiError(c, apiError)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// AddToFavorites godoc
// @Summary Add note to favorites
// @Description Add note to favorites for the authenticated user
// @Tags notes
// @Produce json
// @Security BearerAuth
// @Param id path int true "Note ID"
// @Success 200 {object} string "Note updated successfully"
// @Failure 400 {object} response "Empty query parameter"
// @Failure 401 {object} response "Unauthorized"
// @Failure 500 {object} response "Internal server error"
// @Router /api/notes/{id}/favorites [put]
func (n *NoteHandler) AddToFavorites(c *gin.Context) {
	userId := c.MustGet("UserId").(int)

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	errMove := n.noteService.AddToFavorites(userId, idInt)

	if errMove != nil {
		apiError := model.GetAppropriateApiError(errMove)
		errorResponseFromApiError(c, apiError)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteFromFavorites godoc
// @Summary Delete note to favorites
// @Description Delete note to favorites for the authenticated user
// @Tags notes
// @Produce json
// @Security BearerAuth
// @Param id path int true "Note ID"
// @Success 200 {object} string "Note updated successfully"
// @Failure 400 {object} response "Empty query parameter"
// @Failure 401 {object} response "Unauthorized"
// @Failure 500 {object} response "Internal server error"
// @Router /api/notes/{id}/favorites [delete]
func (n *NoteHandler) DeleteFromFavorites(c *gin.Context) {
	userId := c.MustGet("UserId").(int)

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	errDelete := n.noteService.DeleteFromFavorites(userId, idInt)

	if errDelete != nil {
		apiError := model.GetAppropriateApiError(errDelete)
		errorResponseFromApiError(c, apiError)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
