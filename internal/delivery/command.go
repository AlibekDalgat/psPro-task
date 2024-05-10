package delivery

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"psPro-task/internal/models"
)

func (h *Handler) createCommand(c *gin.Context) {
	var input models.Command
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.CreateCommand(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	go h.services.ExecuteCommand(id, input.Script)
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
