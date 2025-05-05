package service

import (
	"go.uber.org/zap"

	"sweng-task/internal/model"
	"sweng-task/internal/repo"
)

type TrackingService struct {
	repo repo.TrackingEventRepository
	log  *zap.SugaredLogger
}

func NewTrackingService(repo repo.TrackingEventRepository, log *zap.SugaredLogger) *TrackingService {
	return &TrackingService{
		repo: repo,
		log:  log,
	}
}

// Track is uses repository level to persist event, and logs it.
// Suggestion: We can add async to this stage, and let TrackingService.Track to be producer,
// while consumers will persist events.
func (s *TrackingService) Track(event *model.TrackingEvent) error {
	// Future: Budget consumption logic should be added here
	err := s.repo.CreateTrackingEvent(event)
	if err != nil {
		return err
	}
	s.log.Infow("tracking event stored",
		"type", event.EventType,
		"line_item", event.LineItemID,
		"placement", event.Placement,
	)
	return nil
}
