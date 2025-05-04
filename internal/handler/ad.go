package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"sweng-task/internal/service"
)

type AdHandler struct {
	service *service.AdService
	log     *zap.SugaredLogger
}

func NewAdHandler(service *service.AdService, log *zap.SugaredLogger) *AdHandler {
	return &AdHandler{
		service: service,
		log:     log,
	}
}

func (h *AdHandler) GetWinningAds(c *fiber.Ctx) error {
	placement := c.Query("placement")
	if placement == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "placement is required query string",
			})
	}
	category := c.Query("category")
	keyword := c.Query("keyword")
	limitStr := c.Query("limit", "1")
	limit, err := strconv.Atoi(limitStr)
	if err == nil && (limit < 0 || limit > 10) {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "limit should be an integer between 1 and 10",
			})
	}

	ads, err := h.service.GetWinningAds(service.AdQuery{
		Placement: placement,
		Category:  category,
		Keyword:   keyword,
		Limit:     limit,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "failed to retrieve line item",
				"details": err.Error(),
			})
	}
	return c.Status(fiber.StatusOK).JSON(ads)
}
