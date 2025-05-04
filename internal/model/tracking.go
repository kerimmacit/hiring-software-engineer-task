package model

import "time"

// TrackingEventType represents the type of tracking event
type TrackingEventType string

const (
	TrackingEventTypeImpression TrackingEventType = "impression"
	TrackingEventTypeClick      TrackingEventType = "click"
	TrackingEventTypeConversion TrackingEventType = "conversion"
)

// TrackingEvent represents a user interaction with an ad
type TrackingEvent struct {
	EventType  TrackingEventType `json:"event_type" validate:"required,oneof=impression click conversion"`
	LineItemID string            `json:"line_item_id" validate:"required"`
	Timestamp  time.Time         `json:"timestamp,omitempty"`
	Placement  string            `json:"placement,omitempty"`
	UserID     string            `json:"user_id,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}
