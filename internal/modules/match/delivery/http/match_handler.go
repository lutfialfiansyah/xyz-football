package http

import (
	"log"
	"net/http"

	"xyz-football-api/internal/modules/match"
	"xyz-football-api/internal/modules/match/usecase"
	"xyz-football-api/internal/pkg/pagination"
	"xyz-football-api/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	matchUsecase usecase.MatchUsecase
}

func NewMatchHandler(u usecase.MatchUsecase) *MatchHandler {
	return &MatchHandler{
		matchUsecase: u,
	}
}

func (h *MatchHandler) CreateMatch(c *gin.Context) {
	var req match.CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("CreateMatch Binding Error: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	res, err := h.matchUsecase.CreateMatch(c.Request.Context(), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Match schedule created successfully", res)
}

func (h *MatchHandler) GetAllMatches(c *gin.Context) {
	cursor, limit, q := pagination.ParseParams(c)
	status := c.Query("status")

	res, nextCursor, hasMore, err := h.matchUsecase.GetAllMatches(c.Request.Context(), cursor, limit, status, q)
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

func (h *MatchHandler) GetMatchByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.matchUsecase.GetMatchByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Match not found", "Pertandingan tidak ditemukan")
		return
	}
	response.Success(c, http.StatusOK, "Success", res)
}

func (h *MatchHandler) DeleteMatch(c *gin.Context) {
	id := c.Param("id")
	err := h.matchUsecase.DeleteMatch(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Match deleted successfully", nil)
}

func (h *MatchHandler) ChangeMatchStatus(c *gin.Context) {
	id := c.Param("id")
	var req match.ChangeMatchStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	res, err := h.matchUsecase.ChangeMatchStatus(c.Request.Context(), id, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Match status updated successfully", res)
}

// Optional, since MatchEvents also automatically update scores
func (h *MatchHandler) ReportMatchScore(c *gin.Context) {
	id := c.Param("id")
	var req match.ReportMatchScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	res, err := h.matchUsecase.ReportMatchScore(c.Request.Context(), id, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Match score reported successfully", res)
}

func (h *MatchHandler) AddMatchEvent(c *gin.Context) {
	matchID := c.Param("id")
	var req match.AddMatchEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("AddMatchEvent Binding Error: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	res, err := h.matchUsecase.AddMatchEvent(c.Request.Context(), matchID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Match event added successfully", res)
}

func (h *MatchHandler) GetMatchEvents(c *gin.Context) {
	matchID := c.Param("id")
	res, err := h.matchUsecase.GetMatchEvents(c.Request.Context(), matchID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Success", res)
}
