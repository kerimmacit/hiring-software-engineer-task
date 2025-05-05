package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"sweng-task/internal/model"
	"sweng-task/internal/service"
)

type TrackingHandler struct {
	service *service.TrackingService
	log     *zap.SugaredLogger
}

func NewTrackingHandler(service *service.TrackingService, log *zap.SugaredLogger) *TrackingHandler {
	return &TrackingHandler{service: service, log: log}
}

func (h *TrackingHandler) TrackEvent(c *fiber.Ctx) error {
	var event model.TrackingEvent
	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	err := h.service.Track(&event)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "failed to consume tracking event",
				"details": err.Error(),
			})
	}
	return c.JSON(fiber.Map{"success": true})
}
