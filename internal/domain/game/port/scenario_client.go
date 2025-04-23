package port

import (
	scenarioEntity "telemafia/internal/domain/scenario/entity"
)

// ScenarioClient defines an interface for the Game domain to fetch Scenario data.
// Implementations could be local (monolith) or remote (microservice).
type ScenarioClient interface {
	FetchScenario(id string) (*scenarioEntity.Scenario, error)
}
