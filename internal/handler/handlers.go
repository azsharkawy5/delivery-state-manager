package handler

import (
	"delivery-state-manager/internal/models"
	"delivery-state-manager/internal/usecase"
	"delivery-state-manager/pkg/errs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// errorResponse represents an error response
type errorResponse struct {
	Error string `json:"error"`
}

// Handler holds all use cases
type Handler struct {
	driverUC usecase.DriverUseCase
	orderUC  usecase.OrderUseCase
	debugUC  usecase.DebugUseCase
}

// NewHandler creates a new Handler instance
func NewHandler(driverUC usecase.DriverUseCase, orderUC usecase.OrderUseCase, debugUC usecase.DebugUseCase) *Handler {
	return &Handler{
		driverUC: driverUC,
		orderUC:  orderUC,
		debugUC:  debugUC,
	}
}

// SetupRouter sets up the HTTP router with all handlers
func (h *Handler) SetupRouter() *gin.Engine {
	r := gin.Default()

	// Driver endpoints
	r.POST("/drivers", h.createOrUpdateDriverHandler())
	r.GET("/drivers", h.getAllDriversHandler())
	r.GET("/drivers/:id", h.getDriverHandler())
	r.PATCH("/drivers/:id/status", h.updateDriverStatusHandler())

	return r
}

// createOrUpdateDriverHandler handles POST /drivers
func (h *Handler) createOrUpdateDriverHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var driver models.Driver
		if err := c.ShouldBindJSON(&driver); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse{Error: "Invalid request body"})
			return
		}

		if err := h.driverUC.CreateOrUpdateDriver(&driver); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}

		log.Printf("Driver created/updated: %s (%s)", driver.ID, driver.Name)
		c.JSON(http.StatusOK, driver)
	}
}

// getAllDriversHandler handles GET /drivers
func (h *Handler) getAllDriversHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		drivers := h.driverUC.GetAllDrivers()
		c.JSON(http.StatusOK, drivers)
	}
}

// getDriverHandler handles GET /drivers/:id
func (h *Handler) getDriverHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		driver, err := h.driverUC.GetDriver(id)
		if err != nil {
			c.JSON(http.StatusNotFound, errorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, driver)
	}
}

// updateDriverStatusHandler handles PATCH /drivers/:id/status
func (h *Handler) updateDriverStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req struct {
			Status models.DriverStatus `json:"status"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse{Error: "Invalid request body"})
			return
		}

		if err := h.driverUC.UpdateDriverStatus(id, req.Status); err != nil {
			if err == errs.ErrDriverNotFound {
				c.JSON(http.StatusNotFound, errorResponse{Error: err.Error()})
			} else {
				c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
			}
			return
		}

		driver, _ := h.driverUC.GetDriver(id)
		log.Printf("Driver status updated: %s -> %s", id, req.Status)

		c.JSON(http.StatusOK, driver)
	}
}
