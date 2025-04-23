package api

import (
	gamePort "telemafia/internal/domain/game/port"
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	scenarioPort "telemafia/internal/domain/scenario/port"
)

// Ensure LocalScenarioClient implements the gamePort.ScenarioClient interface.
var _ gamePort.ScenarioClient = (*LocalScenarioClient)(nil)

// LocalScenarioClient implements the ScenarioClient interface by directly calling
// the Scenario domain's repository reader within the monolith.
// In a microservice architecture, this would be replaced by an HTTP/gRPC client.
type LocalScenarioClient struct {
	scenarioRepo scenarioPort.ScenarioReader
}

// NewLocalScenarioClient creates a new LocalScenarioClient.
func NewLocalScenarioClient(scenarioRepo scenarioPort.ScenarioReader) *LocalScenarioClient {
	return &LocalScenarioClient{scenarioRepo: scenarioRepo}
}

// FetchScenario retrieves a scenario using the injected ScenarioReader.
func (c *LocalScenarioClient) FetchScenario(id string) (*scenarioEntity.Scenario, error) {
	// In a real microservice, this would make an API call to the Scenario service.
	return c.scenarioRepo.GetScenarioByID(id)
}
