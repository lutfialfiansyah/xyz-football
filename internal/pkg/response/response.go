package response

import (
	"net/http"

	"xyz-football-api/internal/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

type ErrorMessage struct {
	En string `json:"en"`
	Id string `json:"id"`
}

type ErrorResponse struct {
	Success bool         `json:"success"`
	ErrorID int          `json:"error_id"`
	Message ErrorMessage `json:"message"`
}

type CursorMeta struct {
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}

func Success(c *gin.Context, code int, message string, data any) {
	c.JSON(code, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(c *gin.Context, code int, message string, data any, meta any) {
	c.JSON(code, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Error sends a bilingual error response with explicit EN and ID messages.
func Error(c *gin.Context, code int, en string, id string) {
	c.JSON(code, ErrorResponse{
		Success: false,
		ErrorID: code,
		Message: ErrorMessage{
			En: en,
			Id: id,
		},
	})
}

// HandleError auto-detects *apperror.AppError vs generic error.
// AppError → uses its Code, En, Id.
// Generic error → defaults to 500 with generic bilingual message.
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperror.AppError); ok {
		Error(c, appErr.Code, appErr.En, appErr.Id)
		return
	}

	Error(c, http.StatusInternalServerError,
		"Internal server error",
		"Kesalahan server internal",
	)
}
