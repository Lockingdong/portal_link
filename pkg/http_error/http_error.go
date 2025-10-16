package http_error

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ErrInternal = "ErrInternal"

	ErrInvalidParams = "ErrInvalidParams"

	ErrForbidden = "ErrForbidden"

	ErrNotFound = "ErrNotFound"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ResponseInternalServerError 回應 Internal Server Error
func ResponseInternalServerError(c *gin.Context, errorResponse *ErrorResponse) {
	code := ErrInternal
	message := "Internal server error"
	if errorResponse != nil {
		if errorResponse.Code != "" {
			code = errorResponse.Code
		}
		if errorResponse.Message != "" {
			message = errorResponse.Message
		}
	}

	c.JSON(http.StatusInternalServerError, &ErrorResponse{
		Code:    code,
		Message: message,
	})
}

// ResponseBadRequest 回應 Bad Request
func ResponseBadRequest(c *gin.Context, errorResponse *ErrorResponse) {
	code := ErrInvalidParams
	message := "Invalid request parameters"
	if errorResponse != nil {
		if errorResponse.Code != "" {
			code = errorResponse.Code
		}
		if errorResponse.Message != "" {
			message = errorResponse.Message
		}
	}

	c.JSON(http.StatusBadRequest, &ErrorResponse{
		Code:    code,
		Message: message,
	})
}

// ResponseForbidden 回應 Forbidden
func ResponseForbidden(c *gin.Context, errorResponse *ErrorResponse) {
	code := ErrForbidden
	message := "You do not have permission"
	if errorResponse != nil {
		if errorResponse.Code != "" {
			code = errorResponse.Code
		}
		if errorResponse.Message != "" {
			message = errorResponse.Message
		}
	}

	c.JSON(http.StatusForbidden, &ErrorResponse{
		Code:    code,
		Message: message,
	})
}

// ResponseNotFound 回應 Not Found
func ResponseNotFound(c *gin.Context, errorResponse *ErrorResponse) {
	code := ErrNotFound
	message := "Resource not found"
	if errorResponse != nil {
		if errorResponse.Code != "" {
			code = errorResponse.Code
		}
	}

	c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    code,
		Message: message,
	})
}
