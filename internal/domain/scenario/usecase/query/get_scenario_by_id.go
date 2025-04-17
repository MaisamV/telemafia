package query

import (
	"context"
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	scenarioPort "telemafia/internal/domain/scenario/port"
)

// GetScenarioByIDQuery represents the query to get a scenario by ID
type GetScenarioByIDQuery struct {
	ID string
}

// GetScenarioByIDHandler handles queries for a specific scenario
type GetScenarioByIDHandler struct {
	scenarioRepo scenarioPort.ScenarioReader // Use imported ScenarioReader interface
}

// NewGetScenarioByIDHandler creates a new GetScenarioByIDHandler
func NewGetScenarioByIDHandler(repo scenarioPort.ScenarioReader) *GetScenarioByIDHandler {
	return &GetScenarioByIDHandler{
		scenarioRepo: repo,
	}
}

// Handle processes the get scenario by ID query
func (h *GetScenarioByIDHandler) Handle(ctx context.Context, query GetScenarioByIDQuery) (*scenarioEntity.Scenario, error) {
	return h.scenarioRepo.GetScenarioByID(query.ID) // Propagates results/errors from repo
}
