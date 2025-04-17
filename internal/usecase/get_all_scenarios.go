package usecase

import (
	"context"
	"telemafia/internal/entity"
)

// GetAllScenariosQuery represents the query to get all scenarios
type GetAllScenariosQuery struct{}

// GetAllScenariosHandler handles queries for all scenarios
type GetAllScenariosHandler struct {
	scenarioRepo ScenarioReader
}

// NewGetAllScenariosHandler creates a new GetAllScenariosHandler
func NewGetAllScenariosHandler(repo ScenarioReader) *GetAllScenariosHandler {
	return &GetAllScenariosHandler{
		scenarioRepo: repo,
	}
}

// Handle processes the get all scenarios query
func (h *GetAllScenariosHandler) Handle(ctx context.Context) ([]*entity.Scenario, error) {
	return h.scenarioRepo.GetAllScenarios()
}
