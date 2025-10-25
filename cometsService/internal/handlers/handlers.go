package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/domain"
)

type ICometsHandler interface {
	// Observation handlers
	CreateObservation(c *gin.Context)
	GetObservation(c *gin.Context)
	GetUserObservations(c *gin.Context)
	GetUserObservationsByCometID(c *gin.Context)
	UpdateObservation(c *gin.Context)
	DeleteObservation(c *gin.Context)

	// Comet handlers
	CreateComet(c *gin.Context)
	GetComet(c *gin.Context)
	GetUserComets(c *gin.Context)
	DeleteComet(c *gin.Context)

	// Calculation handlers
	CalculateOrbit(c *gin.Context)
	CalculateCloseApproach(c *gin.Context)
	GetCalculationStatus(c *gin.Context)

	// File upload handler
	UploadObservationPhoto(c *gin.Context)
}

type CometsHandler struct {
	cometsService domain.ICometsService
}

func NewCometsHandler(cometsService domain.ICometsService) *CometsHandler {
	return &CometsHandler{
		cometsService: cometsService,
	}
}

// Observation handlers
func (h *CometsHandler) CreateObservation(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	var req domain.CreateObservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	observation, err := h.cometsService.CreateObservation(c.Request.Context(), userID, &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, observation)
}

func (h *CometsHandler) GetObservation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	observation, err := h.cometsService.GetObservation(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	if observation == nil {
		HandleError(c, domain.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, observation)
}

func (h *CometsHandler) GetUserObservations(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	observations, err := h.cometsService.GetUserObservations(c.Request.Context(), userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, observations)
}

func (h *CometsHandler) GetUserObservationsByCometID(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	cometID, err := strconv.Atoi(c.Param("comet_id"))
	if err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	observations, err := h.cometsService.GetUserObservationsByCometID(c.Request.Context(), cometID, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, observations)
}

func (h *CometsHandler) UpdateObservation(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	var req domain.UpdateObservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	err = h.cometsService.UpdateObservation(c.Request.Context(), userID, id, &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Observation updated successfully"})
}

func (h *CometsHandler) DeleteObservation(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	err = h.cometsService.DeleteObservation(c.Request.Context(), id, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Observation deleted successfully"})
}

// Comet handlers
func (h *CometsHandler) CreateComet(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	var req domain.CreateCometRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	comet, err := h.cometsService.CreateComet(c.Request.Context(), userID, &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	response := domain.CometCreatedResponse{
		ID:     comet.ID,
		UserID: comet.UserID,
		Name:   comet.Name,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *CometsHandler) GetComet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	comet, err := h.cometsService.GetComet(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	if comet == nil {
		HandleError(c, domain.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, comet)
}

func (h *CometsHandler) GetUserComets(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	comets, err := h.cometsService.GetUserComets(c.Request.Context(), userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, comets)
}

func (h *CometsHandler) DeleteComet(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	err = h.cometsService.DeleteComet(c.Request.Context(), id, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comet deleted successfully"})
}

// Calculation handlers
func (h *CometsHandler) CalculateOrbit(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	cometID, err := strconv.Atoi(c.Param("comet_id"))
	if err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	result, err := h.cometsService.CalculateOrbit(c.Request.Context(), userID, cometID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *CometsHandler) CalculateCloseApproach(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	cometID, err := strconv.Atoi(c.Param("comet_id"))
	if err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	result, err := h.cometsService.CalculateCloseApproach(c.Request.Context(), userID, cometID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *CometsHandler) GetCalculationStatus(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	requestID, err := strconv.Atoi(c.Param("request_id"))
	if err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	status, err := h.cometsService.GetCalculationStatus(c.Request.Context(), userID, requestID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, status)
}

// File upload handler
func (h *CometsHandler) UploadObservationPhoto(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	// Чтение файла
	openedFile, err := file.Open()
	if err != nil {
		HandleError(c, err)
		return
	}
	defer openedFile.Close()

	fileData := make([]byte, file.Size)
	_, err = openedFile.Read(fileData)
	if err != nil {
		HandleError(c, err)
		return
	}

	photoURL, err := h.cometsService.UploadObservationPhoto(c.Request.Context(), userID, fileData, file.Filename)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"photo_url": photoURL})
}