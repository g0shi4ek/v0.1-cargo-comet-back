package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/domain"
	"github.com/gin-gonic/gin"
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
	UploadCometPhoto(c *gin.Context)

	// Calculation handlers
	CalculateOrbit(c *gin.Context)
	CalculateCloseApproach(c *gin.Context)
	GetCalculationStatus(c *gin.Context)
	GetTrajectory(c *gin.Context)
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
		// Логируем детали ошибки валидации
		log.Printf("Validation error: %v", err)
		
		// Читаем body для отладки
		body, _ := c.GetRawData()
		log.Printf("Request body: %s", string(body))
		
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
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
// Comet handlers
func (h *CometsHandler) CreateComet(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Получаем название из form-data
	name := c.PostForm("name")
	if name == "" {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	// Получаем файл из form-data
	var fileData []byte
	var fileName string

	file, err := c.FormFile("photo")
	if err == nil {
		// Файл присутствует, читаем его
		openedFile, err := file.Open()
		if err != nil {
			HandleError(c, err)
			return
		}
		defer openedFile.Close()

		fileData = make([]byte, file.Size)
		_, err = openedFile.Read(fileData)
		if err != nil {
			HandleError(c, err)
			return
		}
		fileName = file.Filename
	}
	// Если файла нет - это нормально, photoURL будет пустым

	comet, err := h.cometsService.CreateComet(c.Request.Context(), userID, name, fileData, fileName)
	if err != nil {
		c.Error(err) // Логируем в Gin
		HandleError(c, err)
		return
	}

	response := domain.CometCreatedResponse{
		ID:       comet.ID,
		UserID:   comet.UserID,
		Name:     comet.Name,
		PhotoURL: comet.PhotoURL,
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

// File upload handler для кометы
func (h *CometsHandler) UploadCometPhoto(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Получаем comet_id из параметров URL
	cometID, err := strconv.Atoi(c.Param("comet_id"))
	if err != nil {
		HandleError(c, domain.ErrInvalidInput)
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

	// Используем новый метод для загрузки фото кометы, передаем cometID
	comet, err := h.cometsService.UploadCometPhoto(c.Request.Context(), userID, cometID, fileData, file.Filename)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Photo uploaded successfully",
		"photo_url": comet.PhotoURL,
		"comet":     comet,
	})
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

// GetTrajectory получает траекторию кометы для визуализации
func (h *CometsHandler) GetTrajectory(c *gin.Context) {
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

	var req domain.GetTrajectoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	// Парсим временные параметры
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	// Валидация временного диапазона
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		HandleError(c, domain.ErrInvalidInput)
		return
	}

	// Максимальный диапазон - 1 год
	maxDuration := 365 * 24 * time.Hour
	if endTime.Sub(startTime) > maxDuration {
		HandleError(c, errors.New("time range too large, maximum is 1 year"))
		return
	}

	trajectory, err := h.cometsService.GetTrajectory(c.Request.Context(), userID, cometID, startTime, endTime, req.NumPoints)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, trajectory)
}

// func (h *CometsHandler) GetCalculationStatus(c *gin.Context) {
// 	userID, err := GetUserIDFromContext(c)
// 	if err != nil {
// 		HandleError(c, err)
// 		return
// 	}

// 	requestID, err := strconv.Atoi(c.Param("request_id"))
// 	if err != nil {
// 		HandleError(c, domain.ErrInvalidInput)
// 		return
// 	}

// 	status, err := h.cometsService.GetCalculationStatus(c.Request.Context(), userID, requestID)
// 	if err != nil {
// 		HandleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, status)
// }
