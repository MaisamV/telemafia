package command

import (
	"context"
	"errors"
	scenarioPort "telemafia/internal/domain/scenario/port"
	sharedEntity "telemafia/internal/shared/entity"
)

// DeleteScenarioCommand represents the command to delete a scenario
type DeleteScenarioCommand struct {
	Requester sharedEntity.User
	ID        string
}

// DeleteScenarioHandler handles scenario deletion
type DeleteScenarioHandler struct {
	scenarioRepo scenarioPort.ScenarioWriter // Use imported ScenarioWriter interface
}

// NewDeleteScenarioHandler creates a new DeleteScenarioHandler
func NewDeleteScenarioHandler(repo scenarioPort.ScenarioWriter) *DeleteScenarioHandler {
	return &DeleteScenarioHandler{
		scenarioRepo: repo,
	}
}

// Handle processes the delete scenario command
func (h *DeleteScenarioHandler) Handle(ctx context.Context, cmd DeleteScenarioCommand) error {
	// --- Permission Check ---
	if !cmd.Requester.Admin {
		return errors.New("delete scenario: admin privilege required")
	}
	return h.scenarioRepo.DeleteScenario(cmd.ID) // Propagates errors from repo
}
