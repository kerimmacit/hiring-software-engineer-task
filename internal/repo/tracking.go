package repo

import (
	"sync"

	"go.uber.org/zap"

	"sweng-task/internal/model"
)

type TrackingEventRepository interface {
	CreateTrackingEvent(event *model.TrackingEvent) error
}

var _ TrackingEventRepository = (*TrackingEventRepositoryImp)(nil)

type TrackingEventRepositoryImp struct {
	events []*model.TrackingEvent
	mu     sync.RWMutex
	log    *zap.SugaredLogger
}

func NewTrackingEventRepository(log *zap.SugaredLogger) TrackingEventRepository {
	return &TrackingEventRepositoryImp{
		events: make([]*model.TrackingEvent, 0),
		log:    log,
	}
}

func (s *TrackingEventRepositoryImp) CreateTrackingEvent(event *model.TrackingEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, event)
	return nil
}
