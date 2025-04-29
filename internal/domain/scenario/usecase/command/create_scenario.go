package command

import (
	"context"
	"errors"
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	scenarioPort "telemafia/internal/domain/scenario/port"
	sharedEntity "telemafia/internal/shared/entity"
)

// CreateScenarioCommand represents the command to create a new scenario
type CreateScenarioCommand struct {
	Requester sharedEntity.User
	ID        string
	Name      string
}

// CreateScenarioHandler handles scenario creation
type CreateScenarioHandler struct {
	scenarioRepo scenarioPort.ScenarioWriter // Use imported ScenarioWriter interface
}

// NewCreateScenarioHandler creates a new CreateScenarioHandler
func NewCreateScenarioHandler(repo scenarioPort.ScenarioWriter) *CreateScenarioHandler {
	return &CreateScenarioHandler{
		scenarioRepo: repo,
	}
}

// Handle processes the create scenario command
func (h *CreateScenarioHandler) Handle(ctx context.Context, cmd CreateScenarioCommand) error {
	// --- Permission Check ---
	if !cmd.Requester.Admin {
		return errors.New("create scenario: admin privilege required")
	}

	scenario := &scenarioEntity.Scenario{
		ID:    cmd.ID,
		Name:  cmd.Name,
		Sides: []scenarioEntity.Side{}, // Use imported Role type
	}
	return h.scenarioRepo.CreateScenario(scenario) // Propagates errors from repo
}
