package delivery

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"psPro-task/internal/models"
	"strconv"
)

func (h *Handler) createCommand(c *gin.Context) {
	var input models.Command
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "неверное содержание json")
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

func (h *Handler) listCommands(c *gin.Context) {
	commands, err := h.services.GetAllCommands()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Ошибка во время работы сервера")
		return
	}
	c.JSON(http.StatusOK, commands)
}

func (h *Handler) oneCommand(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Неверный парамметр id")
		return
	}
	command, err := h.services.GetOneCommand(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Ошибка во время работы сервера")
		return
	}
	c.JSON(http.StatusOK, command)
}

func (h *Handler) stopCommand(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Неверный параметр id")
		return
	}
	err = h.services.StopCommand(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "Сигнал отправлен"})
}

func (h *Handler) startCommand(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Неверный параметр id")
		return
	}
	err = h.services.StartCommand(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "Сигнал отправлен"})
}

func (h *Handler) killCommand(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Неверный параметр id")
		return
	}
	err = h.services.KillCommand(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "Сигнал отправлен"})
}
