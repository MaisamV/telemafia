package usecase

import (
	"context"
	"telemafia/internal/entity"
)

// CreateScenarioCommand represents the command to create a new scenario
type CreateScenarioCommand struct {
	ID   string
	Name string
}

// CreateScenarioHandler handles scenario creation
type CreateScenarioHandler struct {
	scenarioRepo ScenarioWriter
}

// NewCreateScenarioHandler creates a new CreateScenarioHandler
func NewCreateScenarioHandler(repo ScenarioWriter) *CreateScenarioHandler {
	return &CreateScenarioHandler{
		scenarioRepo: repo,
	}
}

// Handle processes the create scenario command
func (h *CreateScenarioHandler) Handle(ctx context.Context, cmd CreateScenarioCommand) error {
	scenario := &entity.Scenario{
		ID:    cmd.ID,
		Name:  cmd.Name,
		Roles: []entity.Role{},
	}
	return h.scenarioRepo.CreateScenario(scenario)
}
