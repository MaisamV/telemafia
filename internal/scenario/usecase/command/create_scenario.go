package command

import (
	"context"
	"telemafia/internal/scenario/entity"
	"telemafia/internal/scenario/repo"
)

// CreateScenarioCommand represents the command to create a scenario
type CreateScenarioCommand struct {
	ID   string
	Name string
}

// CreateScenarioHandler handles scenario creation
type CreateScenarioHandler struct {
	scenarioRepo repo.Repository
}

// NewCreateScenarioHandler creates a new CreateScenarioHandler
func NewCreateScenarioHandler(repo repo.Repository) *CreateScenarioHandler {
	return &CreateScenarioHandler{
		scenarioRepo: repo,
	}
}

// Handle processes the create scenario command
func (h *CreateScenarioHandler) Handle(ctx context.Context, cmd CreateScenarioCommand) error {
	scenario := &entity.Scenario{
		ID:   cmd.ID,
		Name: cmd.Name,
	}
	return h.scenarioRepo.CreateScenario(scenario)
}
