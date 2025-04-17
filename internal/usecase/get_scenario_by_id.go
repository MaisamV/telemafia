package usecase

import (
	"context"
	"telemafia/internal/entity"
)

// GetScenarioByIDQuery represents the query to get a scenario by ID
type GetScenarioByIDQuery struct {
	ID string
}

// GetScenarioByIDHandler handles queries for a specific scenario
type GetScenarioByIDHandler struct {
	scenarioRepo ScenarioReader
}

// NewGetScenarioByIDHandler creates a new GetScenarioByIDHandler
func NewGetScenarioByIDHandler(repo ScenarioReader) *GetScenarioByIDHandler {
	return &GetScenarioByIDHandler{
		scenarioRepo: repo,
	}
}

// Handle processes the get scenario by ID query
func (h *GetScenarioByIDHandler) Handle(ctx context.Context, query GetScenarioByIDQuery) (*entity.Scenario, error) {
	return h.scenarioRepo.GetScenarioByID(query.ID)
}
