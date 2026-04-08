package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/argform/baitfolio-backend/internal/domain"
	"github.com/argform/baitfolio-backend/internal/geo"
	"github.com/argform/baitfolio-backend/internal/repository"
	"github.com/argform/baitfolio-backend/internal/service"
	httpresponse "github.com/argform/baitfolio-backend/internal/transport/http/response"
)

type PointHandler struct {
	pointService *service.PointService
}

func NewPointHandler(pointService *service.PointService) *PointHandler {
	return &PointHandler{
		pointService: pointService,
	}
}

type CreatePointRequest struct {
	Name string `json:"name"`
	Description *string `json:"description"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type PointResponse struct {
	PointID uint64 `json:"point_id"`
	CreatedBy *uint64 `json:"created_by"`
	Name string `json:"name"`
	Description *string `json:"description"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func newPointResponse(point *domain.Point) PointResponse {
	return PointResponse{
		PointID:     point.PointID,
		CreatedBy:   point.CreatedBy,
		Name:        point.Name,
		Description: point.Description,
		Lat:         point.Lat,
		Lon:         point.Lon,
	}
}

func (h *PointHandler) Create(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		httpresponse.WriteError(c, http.StatusBadRequest, "missing user context")
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		httpresponse.WriteError(c, http.StatusBadRequest, "invalid user context")
		return
	}

	var req CreatePointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.WriteError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	point, err := h.pointService.Create(c.Request.Context(), service.CreatePointInput{
		CreatedBy:   &userID,
		Name:        req.Name,
		Description: req.Description,
		Lat:         req.Lat,
		Lon:         req.Lon,
	})
	if err != nil {
		httpresponse.WriteError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, newPointResponse(point))
}

func (h *PointHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpresponse.WriteError(c, http.StatusBadRequest, "invalid point id")
		return
	}

	point, err := h.pointService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrPointNotFound) {
			httpresponse.WriteError(c, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, newPointResponse(point))
}

func (h *PointHandler) GetAllInsideTile(c *gin.Context) {
	x, err := strconv.Atoi(c.Query("x"))
	if err != nil {
		httpresponse.WriteError(c, http.StatusBadRequest, "invalid tile x")
		return
	}

	y, err := strconv.Atoi(c.Query("y"))
	if err != nil {
		httpresponse.WriteError(c, http.StatusBadRequest, "invalid tile y")
		return
	}

	z, err := strconv.Atoi(c.Query("z"))
	if err != nil {
		httpresponse.WriteError(c, http.StatusBadRequest, "invalid tile z")
		return
	}

	points, err := h.pointService.GetAllInsideTile(c.Request.Context(), geo.Tile{
		X: x,
		Y: y,
		Z: z,
	})
	if err != nil {
		httpresponse.WriteError(c, http.StatusBadRequest, err.Error())
		return
	}

	response := make([]PointResponse, 0, len(points))
	for _, point := range points {
		response = append(response, newPointResponse(point))
	}

	c.JSON(http.StatusOK, response)
}
