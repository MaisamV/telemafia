package main

import (
	"log"
	"telemafia/pkg/event"
)

// EventPublisher implements both room and user event publishers
type EventPublisher struct{}

func (p *EventPublisher) Publish(event interface{}) error {
	// For now, just log the events
	log.Printf("Event published: %+v\n", event)
	return nil
}

// RoomEventPublisher adapts EventPublisher to room.EventPublisher
type RoomEventPublisher struct {
	publisher *EventPublisher
}

func (p *RoomEventPublisher) Publish(event event.Event) error {
	return p.publisher.Publish(event)
}

func main() {
	// Load Configuration
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Dependencies
	botHandler, err := InitializeDependencies(config)
	if err != nil {
		log.Fatal(err)
	}

	// Register bot handlers
	botHandler.RegisterHandlers()

	log.Println("Bot is running...")
	botHandler.Start()
}

// Helper function to check if a string is in a slice
// func contains(slice []string, str string) bool {
// 	for _, s := range slice {
// 		if s == str {
// 			return true
// 		}
// 	}
// 	return false
// }
