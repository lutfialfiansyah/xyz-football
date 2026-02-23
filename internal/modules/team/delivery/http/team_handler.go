package http

import (
	"net/http"

	"xyz-football-api/internal/modules/team"
	"xyz-football-api/internal/modules/team/usecase"
	"xyz-football-api/internal/pkg/pagination"
	"xyz-football-api/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamUsecase usecase.TeamUsecase
}

func NewTeamHandler(u usecase.TeamUsecase) *TeamHandler {
	return &TeamHandler{
		teamUsecase: u,
	}
}

func (h *TeamHandler) Create(c *gin.Context) {
	var req team.CreateTeamRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	res, err := h.teamUsecase.CreateTeam(c.Request.Context(), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Team created successfully", res)
}

func (h *TeamHandler) GetAll(c *gin.Context) {
	cursor, limit, q := pagination.ParseParams(c)

	res, nextCursor, hasMore, err := h.teamUsecase.GetAllTeams(c.Request.Context(), cursor, limit, q)
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

func (h *TeamHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.teamUsecase.GetTeamByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Team not found", "Tim tidak ditemukan")
		return
	}
	response.Success(c, http.StatusOK, "Success", res)
}

func (h *TeamHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req team.UpdateTeamRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	res, err := h.teamUsecase.UpdateTeam(c.Request.Context(), id, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Team updated successfully", res)
}

func (h *TeamHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.teamUsecase.DeleteTeam(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Team deleted successfully", nil)
}
