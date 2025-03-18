package query

import (
	"context"
	"telemafia/internal/scenario/entity"
	"telemafia/internal/scenario/repo"
)

// GetAllScenariosQuery represents the query to get all scenarios
type GetAllScenariosQuery struct{}

// GetAllScenariosHandler handles queries for all scenarios
type GetAllScenariosHandler struct {
	scenarioRepo repo.ScenarioReader
}

// NewGetAllScenariosHandler creates a new GetAllScenariosHandler
func NewGetAllScenariosHandler(repo repo.ScenarioReader) *GetAllScenariosHandler {
	return &GetAllScenariosHandler{
		scenarioRepo: repo,
	}
}

// Handle processes the get all scenarios query
func (h *GetAllScenariosHandler) Handle(ctx context.Context, query GetAllScenariosQuery) ([]*entity.Scenario, error) {
	return h.scenarioRepo.GetAllScenarios()
}
