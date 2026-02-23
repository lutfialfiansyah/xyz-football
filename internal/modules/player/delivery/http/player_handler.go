package http

import (
	"net/http"

	"xyz-football-api/internal/modules/player"
	"xyz-football-api/internal/modules/player/usecase"
	"xyz-football-api/internal/pkg/pagination"
	"xyz-football-api/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	playerUsecase usecase.PlayerUsecase
}

func NewPlayerHandler(u usecase.PlayerUsecase) *PlayerHandler {
	return &PlayerHandler{
		playerUsecase: u,
	}
}

func (h *PlayerHandler) Create(c *gin.Context) {
	var req player.CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	res, err := h.playerUsecase.CreatePlayer(c.Request.Context(), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Player created successfully", res)
}

func (h *PlayerHandler) GetAll(c *gin.Context) {
	cursor, limit, q := pagination.ParseParams(c)

	res, nextCursor, hasMore, err := h.playerUsecase.GetAllPlayers(c.Request.Context(), cursor, limit, q)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	meta := response.CursorMeta{
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}

	response.SuccessWithMeta(c, http.StatusOK, "Success", res, meta)
}

func (h *PlayerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.playerUsecase.GetPlayerByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Player not found", "Pemain tidak ditemukan")
		return
	}
	response.Success(c, http.StatusOK, "Success", res)
}

func (h *PlayerHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req player.UpdatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	res, err := h.playerUsecase.UpdatePlayer(c.Request.Context(), id, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Player updated successfully", res)
}

func (h *PlayerHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.playerUsecase.DeletePlayer(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Player deleted successfully", nil)
}
