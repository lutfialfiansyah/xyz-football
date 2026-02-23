package http

import (
	"log"
	"net/http"

	"xyz-football-api/internal/modules/admin"
	"xyz-football-api/internal/modules/admin/usecase"
	"xyz-football-api/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminUsecase usecase.AdminUsecase
}

func NewAdminHandler(u usecase.AdminUsecase) *AdminHandler {
	return &AdminHandler{
		adminUsecase: u,
	}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req admin.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "Format request tidak valid")
		return
	}

	res, err := h.adminUsecase.Login(c.Request.Context(), req)
	if err != nil {
		log.Printf("Failed login attempt for username: %s, error: %v", req.Username, err)
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Login successful", res)
}
