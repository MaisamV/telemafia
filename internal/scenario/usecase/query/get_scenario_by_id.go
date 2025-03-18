package query

import (
	"context"
	"telemafia/internal/scenario/entity"
	"telemafia/internal/scenario/repo"
)

// GetScenarioByIDQuery represents the query to get a specific scenario
type GetScenarioByIDQuery struct {
	ID string
}

// GetScenarioByIDHandler handles single scenario queries
type GetScenarioByIDHandler struct {
	scenarioRepo repo.ScenarioReader
}

// NewGetScenarioByIDHandler creates a new GetScenarioByIDHandler
func NewGetScenarioByIDHandler(repo repo.ScenarioReader) *GetScenarioByIDHandler {
	return &GetScenarioByIDHandler{
		scenarioRepo: repo,
	}
}

// Handle processes the get scenario query
func (h *GetScenarioByIDHandler) Handle(ctx context.Context, query GetScenarioByIDQuery) (*entity.Scenario, error) {
	return h.scenarioRepo.GetScenarioByID(query.ID)
}
