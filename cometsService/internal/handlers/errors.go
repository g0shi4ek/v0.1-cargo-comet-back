package handlers

import (
	"errors"
	"net/http"

	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/domain"
	"github.com/gin-gonic/gin"
)

var (
	ErrUserNotAuthenticated = errors.New("user not authenticated")
)

func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, domain.ErrorResponse{
			Error:   "Not Found",
			Message: err.Error(),
		})
	case errors.Is(err, domain.ErrUnauthorized):
		c.JSON(http.StatusForbidden, domain.ErrorResponse{
			Error:   "Forbidden",
			Message: err.Error(),
		})
	case errors.Is(err, domain.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
	case errors.Is(err, domain.ErrNotEnoughObservations):
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
	case errors.Is(err, domain.ErrOrbitNotCalculated):
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
	default:
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "An unexpected error occurred",
		})
	}
}
