package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/argform/baitfolio-backend/internal/domain"
	"github.com/argform/baitfolio-backend/internal/repository"
	"github.com/argform/baitfolio-backend/internal/service"
	httpresponse "github.com/argform/baitfolio-backend/internal/transport/http/response"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

type CreateReviewRequest struct {
	PointID uint64  `json:"point_id"`
	Score   int16   `json:"score"`
	Content *string `json:"content"`
}

type ReviewResponse struct {
	ReviewID uint64  `json:"review_id"`
	AuthorID *uint64 `json:"author_id"`
	PointID  uint64  `json:"point_id"`
	Score    int16   `json:"score"`
	Content  *string `json:"content"`
}

func newReviewResponse(review *domain.Review) ReviewResponse {
	return ReviewResponse{
		ReviewID: review.ReviewID,
		AuthorID: review.AuthorID,
		PointID:  review.PointID,
		Score:    review.Score,
		Content:  review.Content,
	}
}

func (h *ReviewHandler) GetAllByPointID(c *gin.Context) {
	pointID, err := strconv.ParseUint(c.Param("pointID"), 10, 64)
	if err != nil {
		httpresponse.WriteError(c, http.StatusBadRequest, "invalid point id")
		return
	}

	reviews, err := h.reviewService.GetAllByPointID(c.Request.Context(), pointID)
	if err != nil {
		httpresponse.WriteError(c, http.StatusBadRequest, err.Error())
		return
	}

	response := make([]ReviewResponse, 0, len(reviews))
	for _, review := range reviews {
		response = append(response, newReviewResponse(review))
	}

	c.JSON(http.StatusOK, response)
}

func (h *ReviewHandler) Create(c *gin.Context) {
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

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.WriteError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	review, err := h.reviewService.Create(c.Request.Context(), service.CreateReviewInput{
		AuthorID: &userID,
		PointID:  req.PointID,
		Score:    req.Score,
		Content:  req.Content,
	})
	if err != nil {
		httpresponse.WriteError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, newReviewResponse(review))
}

func (h *ReviewHandler) Delete(c *gin.Context) {
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

	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpresponse.WriteError(c, http.StatusBadRequest, "invalid review id")
		return
	}

	err = h.reviewService.Delete(c.Request.Context(), reviewID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrReviewNotFound) {
			httpresponse.WriteError(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			httpresponse.WriteError(c, http.StatusForbidden, "forbidden")
			return
		}
		httpresponse.WriteError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
