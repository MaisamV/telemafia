package command

import (
	"context"
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	scenarioPort "telemafia/internal/domain/scenario/port"
)

// CreateScenarioCommand represents the command to create a new scenario
type CreateScenarioCommand struct {
	ID   string
	Name string
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
	scenario := &scenarioEntity.Scenario{
		ID:    cmd.ID,
		Name:  cmd.Name,
		Roles: []scenarioEntity.Role{}, // Use imported Role type
	}
	return h.scenarioRepo.CreateScenario(scenario) // Propagates errors from repo
}
