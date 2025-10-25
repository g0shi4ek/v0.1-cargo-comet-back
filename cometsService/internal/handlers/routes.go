package handlers

import (
	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/domain"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, cometsService domain.ICometsService, authClient domain.IAuthClient) {
	handler := NewCometsHandler(cometsService)

	// Группа маршрутов, требующих аутентификации
	authGroup := router.Group("/api/v1")
	authGroup.Use(AuthMiddleware(authClient))
	{
		// Observation routes
		observations := authGroup.Group("/observations")
		{
			observations.POST("", handler.CreateObservation)
			observations.GET("", handler.GetUserObservations)
			observations.GET("/:id", handler.GetObservation)
			observations.PUT("/:id", handler.UpdateObservation)
			observations.DELETE("/:id", handler.DeleteObservation)
		}

		// Comet routes
		comets := authGroup.Group("/comets")
		{
			comets.POST("", handler.CreateComet)
			comets.GET("", handler.GetUserComets)
			comets.GET("/:id", handler.GetComet)
			comets.DELETE("/:id", handler.DeleteComet)
		}

		// Calculation routes
		calculations := authGroup.Group("/calculations")
		{
			calculations.POST("/:comet_id/orbit", handler.CalculateOrbit)
			calculations.POST("/:comet_id/close-approach", handler.CalculateCloseApproach)
		}

		// File upload routes
		files := authGroup.Group("/files")
		{
			files.POST("/observations/photo", handler.UploadObservationPhoto)
		}

		// Specific observation routes by comet
		authGroup.GET("/observations/comets/:comet_id", handler.GetUserObservationsByCometID)
	}
}
