package port

import (
	scenarioEntity "telemafia/internal/domain/scenario/entity"
)

// ScenarioReader defines the interface for reading scenario data
type ScenarioReader interface {
	GetScenarioByID(id string) (*scenarioEntity.Scenario, error)
	GetAllScenarios() ([]*scenarioEntity.Scenario, error)
}

// ScenarioWriter defines the interface for writing scenario data
type ScenarioWriter interface {
	CreateScenario(scenario *scenarioEntity.Scenario) error
	DeleteScenario(id string) error
}

// ScenarioRepository defines the interface for scenario persistence
type ScenarioRepository interface {
	ScenarioReader
	ScenarioWriter
}
