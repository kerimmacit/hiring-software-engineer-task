package service

import (
	"sync"

	"go.uber.org/zap"

	"sweng-task/internal/model"
)

type TrackingService struct {
	events []*model.TrackingEvent
	mu     sync.RWMutex
	log    *zap.SugaredLogger
}

func NewTrackingService(log *zap.SugaredLogger) *TrackingService {
	return &TrackingService{log: log}
}

func (s *TrackingService) Track(event *model.TrackingEvent) {
	s.mu.Lock()
	s.events = append(s.events, event)
	s.mu.Unlock()
	s.log.Infow("tracking event stored",
		"type", event.EventType,
		"line_item", event.LineItemID,
		"placement", event.Placement,
	)
}
