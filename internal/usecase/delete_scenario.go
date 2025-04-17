package usecase

import (
	"context"
	// "telemafia/internal/entity" // Removed unused import
)

// DeleteScenarioCommand represents the command to delete a scenario
type DeleteScenarioCommand struct {
	ID string
}

// DeleteScenarioHandler handles scenario deletion
type DeleteScenarioHandler struct {
	scenarioRepo ScenarioWriter
}

// NewDeleteScenarioHandler creates a new DeleteScenarioHandler
func NewDeleteScenarioHandler(repo ScenarioWriter) *DeleteScenarioHandler {
	return &DeleteScenarioHandler{
		scenarioRepo: repo,
	}
}

// Handle processes the delete scenario command
func (h *DeleteScenarioHandler) Handle(ctx context.Context, cmd DeleteScenarioCommand) error {
	return h.scenarioRepo.DeleteScenario(cmd.ID)
}
