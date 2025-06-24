package handler

import (
	"Notes/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NotebookHandler struct {
	notebookService service.AbstractNotebookService
}

func NewNotebookHandler(s service.AbstractNotebookService) *NotebookHandler {
	return &NotebookHandler{notebookService: s}
}

// GetNotebook godoc
// @Summary Get user's notebook
// @Description Get the notebook data for the authenticated user
// @Tags notebooks
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.Notebook "Returns user's notebook data"
// @Failure 401 {object} response "Unauthorized"
// @Router /api/notebook [get]
func (n *NotebookHandler) GetNotebook(c *gin.Context) {
	userId := c.MustGet("UserId").(int)

	notebook := n.notebookService.GetUserNotebook(userId)

	c.JSON(http.StatusOK, gin.H{
		"notebook": notebook,
	})
}
