package event

import "time"

// Event represents a domain event
type Event interface {
	EventName() string
	OccurredAt() time.Time
}

// Publisher EventPublisher defines the interface for publishing domain events
type Publisher interface {
	Publish(event Event) error
}
