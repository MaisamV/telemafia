package command

import (
	"context"
	"telemafia/internal/scenario/repo"
)

// DeleteScenarioCommand represents the command to delete a scenario
type DeleteScenarioCommand struct {
	ID string
}

// DeleteScenarioHandler handles scenario deletion
type DeleteScenarioHandler struct {
	scenarioRepo repo.Repository
}

// NewDeleteScenarioHandler creates a new DeleteScenarioHandler
func NewDeleteScenarioHandler(repo repo.Repository) *DeleteScenarioHandler {
	return &DeleteScenarioHandler{
		scenarioRepo: repo,
	}
}

// Handle processes the delete scenario command
func (h *DeleteScenarioHandler) Handle(ctx context.Context, cmd DeleteScenarioCommand) error {
	return h.scenarioRepo.DeleteScenario(cmd.ID)
}
