package query

import (
	"context"
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	scenarioPort "telemafia/internal/domain/scenario/port"
)

// GetAllScenariosQuery represents the query to get all scenarios
type GetAllScenariosQuery struct{}

// GetAllScenariosHandler handles queries for all scenarios
type GetAllScenariosHandler struct {
	scenarioRepo scenarioPort.ScenarioReader // Use imported ScenarioReader interface
}

// NewGetAllScenariosHandler creates a new GetAllScenariosHandler
func NewGetAllScenariosHandler(repo scenarioPort.ScenarioReader) *GetAllScenariosHandler {
	return &GetAllScenariosHandler{
		scenarioRepo: repo,
	}
}

// Handle processes the get all scenarios query
// Note: Query struct is empty, context is unused.
func (h *GetAllScenariosHandler) Handle(ctx context.Context, query GetAllScenariosQuery) ([]*scenarioEntity.Scenario, error) {
	return h.scenarioRepo.GetAllScenarios() // Propagates results/errors from repo
}
